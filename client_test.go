package run9

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestNewClientNormalizesBaseURLAndUsesCallerContextOnly(t *testing.T) {
	client := NewClient(" http://example.com/ ")

	require.Equal(t, "http://example.com", client.baseURL)
	require.Zero(t, client.http.Timeout)
}

func TestClientBoxesPreservesBaseURLPathPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/boxes", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		writeJSONResponse(t, w, http.StatusOK, []BoxView{})
	}))
	defer server.Close()

	views, err := NewClient(server.URL+"/api").Boxes(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "", "", "")
	require.NoError(t, err)
	require.Empty(t, views)
}

func TestClientWithProjectBoxesUsesWorkspacePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/projects/default/workspace/boxes", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		writeJSONResponse(t, w, http.StatusOK, []BoxView{})
	}))
	defer server.Close()

	views, err := NewClient(server.URL+"/api").WithProject("default").Boxes(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "", "", "")
	require.NoError(t, err)
	require.Empty(t, views)
}

func TestClientCreateProjectPostsCanonicalPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		var req CreateProjectRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.Equal(t, CreateProjectRequest{
			DisplayName: "Sandbox",
			Description: "Isolated experiments",
		}, req)

		writeJSONResponse(t, w, http.StatusCreated, ProjectView{
			ProjectID:   "proj-sandbox",
			OrgID:       "org-1",
			ProjectCID:  "sandbox123abc",
			DisplayName: "Sandbox",
			Description: "Isolated experiments",
			Role:        ProjectRole("admin"),
		})
	}))
	defer server.Close()

	view, err := NewClient(server.URL).CreateProject(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, CreateProjectRequest{
		DisplayName: "Sandbox",
		Description: "Isolated experiments",
	})
	require.NoError(t, err)
	require.Equal(t, "sandbox123abc", view.ProjectCID)
	require.Equal(t, "Sandbox", view.DisplayName)
}

func TestClientUpdateAccountUsesAccountRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/account", r.URL.Path)
		require.Equal(t, http.MethodPatch, r.Method)

		var req UpdateMeRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.NotNil(t, req.DisplayName)
		require.Equal(t, "Alice CLI", *req.DisplayName)

		writeJSONResponse(t, w, http.StatusOK, MeView{
			UserID:       "user-1",
			PrimaryEmail: "alice@example.com",
			DisplayName:  "Alice CLI",
		})
	}))
	defer server.Close()

	name := "Alice CLI"
	view, err := NewClient(server.URL).UpdateAccount(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, UpdateMeRequest{DisplayName: &name})
	require.NoError(t, err)
	require.Equal(t, "Alice CLI", view.DisplayName)
}

func TestClientUpdateBoxCanClearLabelsWithoutSendingNull(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/workspace/boxes/box-1", r.URL.Path)
		require.Equal(t, http.MethodPatch, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"labels":{},"desired_shape":"2c4g","network_mode":"managed","security_mode":"restricted"}`, string(body))

		writeJSONResponse(t, w, http.StatusOK, BoxView{
			BoxID:        "box-1",
			DesiredShape: "2c4g",
			Labels:       map[string]string{},
		})
	}))
	defer server.Close()

	labels := map[string]string{}
	shape := "2c4g"
	networkMode := BoxNetworkModeManaged
	securityMode := BoxSecurityModeRestricted
	view, err := NewClient(server.URL).WithProject("default").UpdateBox(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "box-1", UpdateBoxRequest{
		Labels:       &labels,
		DesiredShape: &shape,
		NetworkMode:  &networkMode,
		SecurityMode: &securityMode,
	})
	require.NoError(t, err)
	require.Equal(t, "box-1", view.BoxID)
	require.Empty(t, view.Labels)
}

func TestClientCreateBoxFromSharedSnapUsesWorkspaceRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/sandbox/workspace/shared-snaps/python-dev/boxes", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		var req CreateBoxFromSharedSnapRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		require.Equal(t, "box-1", req.BoxID)
		require.Equal(t, "2c4g", req.DesiredShape)

		writeJSONResponse(t, w, http.StatusCreated, BoxView{
			BoxID:        "box-1",
			DesiredShape: "2c4g",
		})
	}))
	defer server.Close()

	view, err := NewClient(server.URL).WithProject("sandbox").CreateBoxFromSharedSnap(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "python-dev", CreateBoxFromSharedSnapRequest{
		BoxID:        "box-1",
		DesiredShape: "2c4g",
	})
	require.NoError(t, err)
	require.Equal(t, "box-1", view.BoxID)
}

func TestClientExecsIncludesExtendedFilters(t *testing.T) {
	acceptedAfter := time.Unix(1_700_000_000, 123_456_789).UTC()
	acceptedBefore := time.Unix(1_700_000_000, 987_654_321).UTC()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/workspace/execs", r.URL.Path)
		require.Equal(t, "box-1", r.URL.Query().Get("box_id"))
		require.Equal(t, "running", r.URL.Query().Get("state"))
		require.Equal(t, "user-1", r.URL.Query().Get("creator"))
		require.Equal(t, acceptedAfter.Format(time.RFC3339Nano), r.URL.Query().Get("accepted_after"))
		require.Equal(t, acceptedBefore.Format(time.RFC3339Nano), r.URL.Query().Get("accepted_before"))
		require.Equal(t, "asc", r.URL.Query().Get("order"))
		require.Equal(t, "true", r.URL.Query().Get("paged"))
		require.Equal(t, "10", r.URL.Query().Get("limit"))
		require.Equal(t, "cursor-1", r.URL.Query().Get("cursor"))
		w.Header().Set("X-Run9-Next-Cursor", "cursor-2")
		writeJSONResponse(t, w, http.StatusOK, []ExecView{})
	}))
	defer server.Close()

	limit := 10
	result, err := NewClient(server.URL).WithProject("default").Execs(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, ExecListRequest{
		BoxID:          "box-1",
		State:          "running",
		Creator:        "user-1",
		AcceptedAfter:  &acceptedAfter,
		AcceptedBefore: &acceptedBefore,
		Order:          "asc",
		Paged:          true,
		Limit:          &limit,
		Cursor:         "cursor-1",
	})
	require.NoError(t, err)
	require.Empty(t, result.Execs)
	require.Equal(t, "cursor-2", result.NextCursor)
}

func TestClientDownloadExecLogReturnsRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/workspace/execs/exec-1/log-download", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("line one\nline two\n"))
		require.NoError(t, err)
	}))
	defer server.Close()

	body, err := NewClient(server.URL).WithProject("default").DownloadExecLog(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "exec-1")
	require.NoError(t, err)
	defer body.Close()

	data, err := io.ReadAll(body)
	require.NoError(t, err)
	require.Equal(t, "line one\nline two\n", string(data))
}

func TestClientBoxesReturnsJSONErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/boxes", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid state filter",
		}))
	}))
	defer server.Close()

	_, err := NewClient(server.URL).Boxes(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "", "", "broken")
	require.Error(t, err)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	require.Equal(t, "invalid state filter", apiErr.Message)
}

func TestClientWhoAmIReturnsErrorOnEmptyResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/whoami", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(" \n\t "))
		require.NoError(t, err)
	}))
	defer server.Close()

	_, err := NewClient(server.URL).WhoAmI(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	})
	require.EqualError(t, err, "portal api returned empty response body")
}

func TestClientForkSnapPostsToCanonicalEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/snaps/snap-1/fork", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		ak, sk, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "ak-1", ak)
		require.Equal(t, "sk-1", sk)
		writeJSONResponse(t, w, http.StatusCreated, SnapView{
			SnapID:      "snap-2",
			ParentChain: []string{"snap-1"},
		})
	}))
	defer server.Close()

	view, err := NewClient(server.URL).ForkSnap(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, " snap-1 ")
	require.NoError(t, err)
	require.Equal(t, "snap-2", view.SnapID)
	require.Equal(t, []string{"snap-1"}, view.ParentChain)
}

func TestClientDownloadArchiveFallsBackToRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/boxes/box-1/files/download", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte("gateway overloaded"))
		require.NoError(t, err)
	}))
	defer server.Close()

	_, err := NewClient(server.URL).DownloadArchive(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "box-1", "/work/result.txt")
	require.Error(t, err)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
	require.Equal(t, "gateway overloaded", apiErr.Message)
}

func TestClientUploadArchivePreservesExplicitTarContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/boxes/box-1/files/upload", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/work/result.txt", r.URL.Query().Get("box_abs_path"))
		require.Equal(t, "application/x-tar", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "tar-body", string(body))

		writeJSONResponse(t, w, http.StatusOK, RuntimeRequestView{
			RuntimeRequestID: "runtime-upload-1",
			State:            "prepared",
		})
	}))
	defer server.Close()

	view, err := NewClient(server.URL).UploadArchive(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "box-1", "/work/result.txt", bytes.NewBufferString("tar-body"))
	require.NoError(t, err)
	require.Equal(t, "runtime-upload-1", view.RuntimeRequestID)
	require.Equal(t, "prepared", view.State)
}

func TestClientExecStreamFollowsRedirectAndKeepsFinalExecIDHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/boxes/box-1/execs/stream":
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "application/x-ndjson", r.Header.Get("Accept"))
			require.Equal(t, "inline", r.Header.Get("X-Run9-Exec-Stream-Mode"))
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "application/x-ndjson", r.Header.Get("Accept"))
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Run9-Exec-ID", "exec-redirected")
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "started"}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "exit", ExitCode: 0}))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	execID, body, err := NewClient(server.URL).ExecStream(context.Background(), Credentials{
		AK: "ak-1",
		SK: "sk-1",
	}, "box-1", ExecBoxRequest{Command: []string{"printf", "hello"}})
	require.NoError(t, err)
	require.Equal(t, "exec-redirected", execID)
	defer body.Close()

	payload, err := io.ReadAll(body)
	require.NoError(t, err)
	require.Contains(t, string(payload), `"type":"started"`)
	require.Contains(t, string(payload), `"type":"exit"`)
}

func TestClientExecAttachURLOverridesDefaultWebsocketHandshakeTimeout(t *testing.T) {
	oldTimeout := websocket.DefaultDialer.HandshakeTimeout
	websocket.DefaultDialer.HandshakeTimeout = time.Millisecond
	t.Cleanup(func() {
		websocket.DefaultDialer.HandshakeTimeout = oldTimeout
	})

	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foreground-relay/execs/ticket-1/exec-attach", r.URL.Path)
		time.Sleep(50 * time.Millisecond)
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		require.NoError(t, conn.Close())
	}))
	defer server.Close()

	socket, err := NewClient(server.URL).ExecAttachURL(context.Background(), "/foreground-relay/execs/ticket-1/exec-attach")
	require.NoError(t, err)
	require.NoError(t, socket.Close())
}

func TestClientExecAttachURLResolvesRelativePathAgainstBaseURL(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foreground-relay/execs/ticket-1/exec-attach", r.URL.Path)
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		require.NoError(t, conn.Close())
	}))
	defer server.Close()

	socket, err := NewClient(server.URL).ExecAttachURL(context.Background(), "foreground-relay/execs/ticket-1/exec-attach")
	require.NoError(t, err)
	require.NoError(t, socket.Close())
}

func TestClientExecAttachURLReturnsJSONHandshakeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/foreground-relay/execs/ticket-1/exec-attach", r.URL.Path)
		writeJSONResponse(t, w, http.StatusUnauthorized, map[string]string{
			"error": "invalid API key",
		})
	}))
	defer server.Close()

	_, err := NewClient(server.URL).ExecAttachURL(context.Background(), "/foreground-relay/execs/ticket-1/exec-attach")
	require.Error(t, err)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	require.Equal(t, "invalid API key", apiErr.Message)
}

func writeJSONResponse(t *testing.T, w http.ResponseWriter, status int, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	require.NoError(t, json.NewEncoder(w).Encode(value))
}
