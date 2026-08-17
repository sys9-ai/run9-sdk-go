package run9

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientStartPrewarmRecordingReturnsTypedStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/prewarm-profiles/record", r.URL.Path)
		var request RecordPrewarmProfileRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, RecordPrewarmProfileRequest{
			Name: "typescript", BaseSnapID: "snap-base", Command: []string{"npx", "tsc", "--version"},
		}, request)
		w.Header().Set("X-Run9-Exec-ID", "exec-recording")
		w.Header().Set("X-Run9-Prewarm-Profile-ID", "profile-recording")
		_, err := w.Write([]byte("{\"type\":\"exit\",\"exit_code\":0,\"prewarm_recording\":{\"profile_id\":\"profile-recording\",\"generation\":\"generation-1\",\"sha256\":\"digest\",\"blocks\":7,\"bytes\":8192}}\n"))
		require.NoError(t, err)
	}))
	defer server.Close()

	recording, err := newProjectTestClient(t, server.URL, "default").StartPrewarmRecording(context.Background(), RecordPrewarmProfileRequest{
		Name: "typescript", BaseSnapID: "snap-base", Command: []string{"npx", "tsc", "--version"},
	})
	require.NoError(t, err)
	require.Equal(t, "profile-recording", recording.ProfileID)
	require.Equal(t, "exec-recording", recording.ExecID)
	result, err := recording.Stream.Pump(context.Background(), ExecOutputWriters{})
	require.NoError(t, err)
	require.NotNil(t, result.PrewarmRecording)
	require.Equal(t, uint64(7), result.PrewarmRecording.Blocks)
}

func TestClientManagesPrewarmProfileEnabledState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/projects/default/workspace/prewarm-profiles/typescript", r.URL.Path)
		var request map[string]bool
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, map[string]bool{"enabled": false}, request)
		writeJSONResponse(t, w, http.StatusOK, PrewarmProfileView{Name: "typescript", State: PrewarmProfileStateReady})
	}))
	defer server.Close()

	profile, err := newProjectTestClient(t, server.URL, "default").SetPrewarmProfileEnabled(context.Background(), "typescript", false)
	require.NoError(t, err)
	require.Equal(t, "typescript", profile.Name)
	require.False(t, profile.Enabled)
}
