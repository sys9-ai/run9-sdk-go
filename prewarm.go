package run9

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MaxPrewarmRecordingRuntime is the largest workload duration accepted for one recording.
const MaxPrewarmRecordingRuntime = 12 * time.Hour

// ListPrewarmProfiles lists recorded prewarm profiles in the selected project.
func (c *Client) ListPrewarmProfiles(ctx context.Context) ([]PrewarmProfileView, error) {
	var profiles []PrewarmProfileView
	err := c.doWorkspace(ctx, http.MethodGet, "/prewarm-profiles", requestOptions{result: &profiles})
	return profiles, err
}

// GetPrewarmProfile loads one recorded prewarm profile by name or id.
func (c *Client) GetPrewarmProfile(ctx context.Context, nameOrID string) (PrewarmProfileView, error) {
	var profile PrewarmProfileView
	err := c.doWorkspace(ctx, http.MethodGet, "/prewarm-profiles/"+url.PathEscape(strings.TrimSpace(nameOrID)), requestOptions{result: &profile})
	return profile, err
}

// SetPrewarmProfileEnabled enables or disables automatic selection of one profile.
func (c *Client) SetPrewarmProfileEnabled(ctx context.Context, nameOrID string, enabled bool) (PrewarmProfileView, error) {
	var profile PrewarmProfileView
	err := c.doWorkspace(ctx, http.MethodPatch, "/prewarm-profiles/"+url.PathEscape(strings.TrimSpace(nameOrID)), requestOptions{
		body: map[string]bool{"enabled": enabled}, result: &profile,
	})
	return profile, err
}

// StartPrewarmRecording creates an enabled profile and starts its standard workload stream.
// Closing or cancelling the stream cancels only the workload; the service still finalizes collected blocks.
func (c *Client) StartPrewarmRecording(ctx context.Context, req RecordPrewarmProfileRequest) (*PrewarmRecordingStream, error) {
	if req.MaxRuntime < 0 {
		return nil, errors.New("prewarm recording max runtime must not be negative")
	}
	if req.MaxRuntime > MaxPrewarmRecordingRuntime {
		return nil, errors.New("prewarm recording max runtime must not exceed 12 hours")
	}
	var maxRuntimeSeconds uint64
	if req.MaxRuntime > 0 {
		if req.MaxRuntime < time.Second {
			return nil, errors.New("prewarm recording max runtime must be at least one second")
		}
		maxRuntimeSeconds = uint64(req.MaxRuntime / time.Second)
		if req.MaxRuntime%time.Second != 0 {
			maxRuntimeSeconds++
		}
	}
	resp, err := c.doWorkspaceRaw(ctx, http.MethodPost, "/prewarm-profiles/record", requestOptions{
		body: struct {
			Name              string   `json:"name"`
			BaseSnapID        string   `json:"base_snap_id"`
			Command           []string `json:"command"`
			MaxRuntimeSeconds uint64   `json:"max_runtime_seconds,omitempty"`
		}{
			Name: req.Name, BaseSnapID: req.BaseSnapID, Command: req.Command,
			MaxRuntimeSeconds: maxRuntimeSeconds,
		},
		headers: map[string]string{"Accept": "application/x-ndjson"},
	})
	if err != nil {
		return nil, err
	}
	profileID := strings.TrimSpace(resp.Header.Get("X-Run9-Prewarm-Profile-ID"))
	execID := strings.TrimSpace(resp.Header.Get("X-Run9-Exec-ID"))
	if profileID == "" || execID == "" {
		_ = resp.Body.Close()
		return nil, errors.New("prewarm recording response omitted profile or exec id")
	}
	return &PrewarmRecordingStream{ProfileID: profileID, ExecID: execID, Stream: newExecStream(execID, resp.Body)}, nil
}
