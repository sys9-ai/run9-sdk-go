package run9

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientPullBackgroundExecOutputReadsBinaryBodyAndHeaders(t *testing.T) {
	fixtureTime := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/execs/exec-1/pull-output", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "cursor-1", body["cursor"])
		require.EqualValues(t, 2000, body["wait_ms"])

		w.Header().Set("X-Run9-Next-Cursor", "cursor-2")
		w.Header().Set("X-Run9-Exec-State", "running")
		w.Header().Set("X-Run9-Idle-Deadline-At", fixtureTime.Format(time.RFC3339Nano))
		_, err := w.Write([]byte("binary-body"))
		require.NoError(t, err)
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").PullBackgroundExecOutput(context.Background(), "exec-1", PullBackgroundExecOutputRequest{
		Cursor: "cursor-1",
		Wait:   2 * time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, []byte("binary-body"), result.Body)
	require.Equal(t, "cursor-2", result.NextCursor)
	require.Equal(t, "running", result.State)
	require.NotNil(t, result.IdleDeadlineAt)
	require.True(t, result.IdleDeadlineAt.Equal(fixtureTime))
}

func TestClientWriteBackgroundExecStdinSendsOctetStreamAndEOFHeader(t *testing.T) {
	fixtureTime := time.Date(2026, 3, 28, 12, 5, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/execs/exec-1/write-stdin", r.URL.Path)
		require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
		require.Equal(t, "true", r.Header.Get("X-Run9-Close-Stdin"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "hello", string(body))

		w.Header().Set("X-Run9-Idle-Deadline-At", fixtureTime.Format(time.RFC3339Nano))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	idleDeadlineAt, err := newProjectTestClient(t, server.URL, "default").WriteBackgroundExecStdin(context.Background(), "exec-1", WriteBackgroundExecStdinRequest{
		Data:       []byte("hello"),
		CloseStdin: true,
	})
	require.NoError(t, err)
	require.NotNil(t, idleDeadlineAt)
	require.True(t, idleDeadlineAt.Equal(fixtureTime))
}

func TestClientPullBackgroundExecOutputReturnsJSONErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/execs/exec-1/pull-output", r.URL.Path)
		writeJSONResponse(t, w, http.StatusConflict, map[string]string{
			"error": "background exec owner is not available",
		})
	}))
	defer server.Close()

	_, err := newProjectTestClient(t, server.URL, "default").PullBackgroundExecOutput(context.Background(), "exec-1", PullBackgroundExecOutputRequest{
		Wait: 2 * time.Second,
	})
	require.Error(t, err)

	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusConflict, apiErr.StatusCode)
	require.Equal(t, "background exec owner is not available", apiErr.Message)
}

func TestClientKillBackgroundExecAcceptsEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/execs/exec-1/kill", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := newProjectTestClient(t, server.URL, "default").KillBackgroundExec(context.Background(), "exec-1")
	require.NoError(t, err)
}
