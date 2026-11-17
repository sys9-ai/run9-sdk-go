package run9

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientCreatePrivateNetworkSendsTailscaleConfigAndDecodesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/private-networks", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"name":"corp","type":"tailscale","tailscale":{"oauth_client_secret":"tskey-client-k1","tags":["tag:run9"]},"box_label":{"key":"team","value":"web"}}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = io.WriteString(w, `{"private_network_id":"pnet1","project_id":"p1","name":"corp","hostname":"corp","type":"tailscale","tailscale":{"tags":["tag:run9"]},"box_label":{"key":"team","value":"web"},"status":{"state":"connected","hostname":"corp.example.ts.net","dns_domains":["example.ts.net"],"routes":["100.64.0.0/10"],"updated_at":"2026-11-17T10:00:00Z"}}`)
		require.NoError(t, err)
	}))
	defer server.Close()

	view, err := newProjectTestClient(t, server.URL, "default").CreatePrivateNetwork(context.Background(), CreatePrivateNetworkRequest{
		Name: "corp",
		Type: PrivateNetworkTypeTailscale,
		Tailscale: CreatePrivateNetworkTailscaleRequest{
			OAuthClientSecret: "tskey-client-k1",
			Tags:              []string{"tag:run9"},
		},
		BoxLabel: &BoxLabel{Key: "team", Value: "web"},
	})
	require.NoError(t, err)
	reportedAt := time.Date(2026, 11, 17, 10, 0, 0, 0, time.UTC)
	require.Equal(t, PrivateNetworkView{
		PrivateNetworkID: "pnet1",
		ProjectID:        "p1",
		Name:             "corp",
		Hostname:         "corp",
		Type:             PrivateNetworkTypeTailscale,
		Tailscale:        PrivateNetworkTailscaleView{Tags: []string{"tag:run9"}},
		BoxLabel:         &BoxLabel{Key: "team", Value: "web"},
		Status: PrivateNetworkStatusView{
			State:      PrivateNetworkStateConnected,
			Hostname:   "corp.example.ts.net",
			DNSDomains: []string{"example.ts.net"},
			Routes:     []string{"100.64.0.0/10"},
			UpdatedAt:  &reportedAt,
		},
	}, view)
}

func TestClientUpdatePrivateNetworkSendsOnlyTheSetFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/private-networks/pnet1", r.URL.Path)
		require.Equal(t, http.MethodPatch, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"tailscale":{"oauth_client_secret":"tskey-client-k2"}}`, string(body))

		writeJSONResponse(t, w, http.StatusOK, PrivateNetworkView{PrivateNetworkID: "pnet1"})
	}))
	defer server.Close()

	secret := "tskey-client-k2"
	view, err := newProjectTestClient(t, server.URL, "default").UpdatePrivateNetwork(context.Background(), "pnet1", UpdatePrivateNetworkRequest{
		Tailscale: &UpdatePrivateNetworkTailscaleRequest{OAuthClientSecret: &secret},
	})
	require.NoError(t, err)
	require.Equal(t, "pnet1", view.PrivateNetworkID)
}

func TestClientListGetAndDeletePrivateNetworks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "GET /projects/default/private-networks":
			writeJSONResponse(t, w, http.StatusOK, []PrivateNetworkView{{PrivateNetworkID: "pnet1"}, {PrivateNetworkID: "pnet2"}})
		case "GET /projects/default/private-networks/pnet1":
			writeJSONResponse(t, w, http.StatusOK, PrivateNetworkView{PrivateNetworkID: "pnet1", Name: "corp"})
		case "DELETE /projects/default/private-networks/pnet1":
			writeJSONResponse(t, w, http.StatusOK, map[string]string{"status": "removed"})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	client := newProjectTestClient(t, server.URL, "default")

	networks, err := client.ListPrivateNetworks(context.Background())
	require.NoError(t, err)
	require.Len(t, networks, 2)

	network, err := client.GetPrivateNetwork(context.Background(), "pnet1")
	require.NoError(t, err)
	require.Equal(t, "corp", network.Name)

	require.NoError(t, client.DeletePrivateNetwork(context.Background(), "pnet1"))
}

func TestClientReportsAMissingPrivateNetworkGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"this deployment has no private network gateway","code":"private_network_gateway_not_configured"}`))
	}))
	defer server.Close()

	_, err := newProjectTestClient(t, server.URL, "default").CreatePrivateNetwork(context.Background(), CreatePrivateNetworkRequest{Name: "corp"})
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, ErrorCodePrivateNetworkGatewayNotConfigured, apiErr.Code)
}

func TestClientDecodesTheBoxPrivateNetwork(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/projects/default/workspace/boxes/box-a", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, `{"box_id":"box-a","private_network":"pnet1"}`)
		require.NoError(t, err)
	}))
	defer server.Close()

	box, err := newProjectTestClient(t, server.URL, "default").GetBox(context.Background(), "box-a")
	require.NoError(t, err)
	require.Equal(t, "pnet1", box.PrivateNetwork)
}
