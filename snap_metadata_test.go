package run9

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnapLabelsRequests(t *testing.T) {
	labels := map[string]string{"owner": "chord", "quoted": "a\"=b", "empty": ""}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/projects/default/workspace/snaps", r.URL.Path)
			require.ElementsMatch(t, []string{"owner=chord", "quoted=a\"=b", "empty="}, r.URL.Query()["label"])
			writeJSONResponse(t, w, 200, []SnapView{{SnapID: "snap", Labels: labels}})
		case http.MethodPost:
			require.Equal(t, "/projects/default/workspace/snaps/base/fork", r.URL.Path)
			var req ForkSnapRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			require.Equal(t, labels, req.Labels)
			writeJSONResponse(t, w, 201, SnapView{SnapID: "snap", Labels: req.Labels})
		case http.MethodPatch:
			require.Equal(t, "/projects/default/workspace/snaps/snap", r.URL.Path)
			var req UpdateSnapRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			require.NotNil(t, req.Labels, "empty map must not be omitted when clearing labels")
			require.Empty(t, *req.Labels)
			writeJSONResponse(t, w, 200, SnapView{SnapID: "snap"})
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()
	client := newProjectTestClient(t, server.URL, "default")
	child, err := client.ForkSnapWithOptions(t.Context(), "base", ForkSnapRequest{Labels: labels})
	require.NoError(t, err)
	require.Equal(t, labels, child.Labels)
	found, err := client.ListSnaps(t.Context(), ListSnapsRequest{Labels: labels})
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Equal(t, labels, found[0].Labels)
	empty := map[string]string{}
	_, err = client.UpdateSnap(t.Context(), child.SnapID, UpdateSnapRequest{Labels: &empty})
	require.NoError(t, err)
}
