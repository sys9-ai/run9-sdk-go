package run9

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientStartPrewarmRecordingReturnsTypedStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/prewarm-profiles/record", r.URL.Path)
		var request struct {
			Name              string   `json:"name"`
			BaseSnapID        string   `json:"base_snap_id"`
			Command           []string `json:"command"`
			MaxRuntimeSeconds uint64   `json:"max_runtime_seconds"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "typescript", request.Name)
		require.Equal(t, "snap-base", request.BaseSnapID)
		require.Equal(t, []string{"npx", "tsc", "--version"}, request.Command)
		require.Equal(t, uint64(91), request.MaxRuntimeSeconds)
		w.Header().Set("X-Run9-Exec-ID", "exec-recording")
		w.Header().Set("X-Run9-Prewarm-Profile-ID", "profile-recording")
		_, err := w.Write([]byte("{\"type\":\"exit\",\"exit_code\":0,\"prewarm_recording\":{\"profile_id\":\"profile-recording\",\"generation\":\"generation-1\",\"sha256\":\"digest\",\"blocks\":7,\"bytes\":8192}}\n"))
		require.NoError(t, err)
	}))
	defer server.Close()

	recording, err := newProjectTestClient(t, server.URL, "default").StartPrewarmRecording(context.Background(), RecordPrewarmProfileRequest{
		Name: "typescript", BaseSnapID: "snap-base", Command: []string{"npx", "tsc", "--version"}, MaxRuntime: 90*time.Second + time.Millisecond,
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

func TestClientRejectsPrewarmRecordingAboveTwelveHours(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	_, err := newProjectTestClient(t, server.URL, "default").StartPrewarmRecording(context.Background(), RecordPrewarmProfileRequest{
		Name: "server", BaseSnapID: "snap-base", Command: []string{"./server"}, MaxRuntime: 12*time.Hour + time.Second,
	})
	require.EqualError(t, err, "prewarm recording max runtime must not exceed 12 hours")
}
