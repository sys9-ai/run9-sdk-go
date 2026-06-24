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

func TestExecViewTerminalResult(t *testing.T) {
	exitCode := 23

	succeeded := ExecView{State: "succeeded", ExitCode: &exitCode}
	require.NotNil(t, succeeded.TerminalResult())
	require.Equal(t, ExecTerminalStatusExited, succeeded.TerminalResult().Status)
	require.NotNil(t, succeeded.TerminalResult().ExitCode)
	require.Equal(t, 23, *succeeded.TerminalResult().ExitCode)

	failed := ExecView{State: "failed", ExitCode: &exitCode}
	require.NotNil(t, failed.TerminalResult())
	require.Equal(t, ExecTerminalStatusExited, failed.TerminalResult().Status)
	require.NotNil(t, failed.TerminalResult().ExitCode)
	require.Equal(t, 23, *failed.TerminalResult().ExitCode)

	cancelled := ExecView{State: "cancelled", Reason: "explicit_cancel"}
	require.NotNil(t, cancelled.TerminalResult())
	require.Equal(t, ExecTerminalStatusCancelled, cancelled.TerminalResult().Status)
	require.Equal(t, "explicit_cancel", cancelled.TerminalResult().Reason)

	running := ExecView{State: "running"}
	require.Nil(t, running.TerminalResult())
}

func TestClientRunExecCaptureReturnsTerminalResultAndMergedLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-1/execs/stream":
			require.Equal(t, http.MethodPost, r.Method)
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Run9-Exec-ID", "exec-1")
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "started"}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "stdout", Data: []byte("hello\n")}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "stderr", Data: []byte("warn\n")}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "exit", ExitCode: 0}))
		case "/projects/default/workspace/execs/exec-1/log-download":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("hello\nwarn\n"))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").RunExecCapture(context.Background(), "box-1", ExecRequest{
		Command: []string{"printf", "hello"},
	})
	require.NoError(t, err)
	require.Equal(t, "exec-1", result.ExecID)
	require.Equal(t, ExecTerminalStatusExited, result.Terminal.Status)
	require.NotNil(t, result.Terminal.ExitCode)
	require.Equal(t, 0, *result.Terminal.ExitCode)
	require.Equal(t, []byte("hello\nwarn\n"), result.Transcript)
	require.Empty(t, result.TranscriptUnavailableReason)
}

func TestClientRunExecCaptureRecoversFromStreamDisconnect(t *testing.T) {
	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-1/execs/stream":
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Run9-Exec-ID", "exec-1")
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "started"}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "stdout", Data: []byte("hello\n")}))
		case "/projects/default/workspace/execs/exec-1":
			pollCount++
			writeJSONResponse(t, w, http.StatusOK, ExecView{
				ExecID:    "exec-1",
				State:     "succeeded",
				ExitCode:  intPtr(0),
				Reason:    "",
				BoxID:     "box-1",
				ProjectID: "proj-1",
			})
		case "/projects/default/workspace/execs/exec-1/log-download":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("hello\nworld\n"))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").RunExecCapture(context.Background(), "box-1", ExecRequest{
		Command: []string{"printf", "hello"},
	})
	require.NoError(t, err)
	require.Equal(t, "exec-1", result.ExecID)
	require.Equal(t, ExecTerminalStatusExited, result.Terminal.Status)
	require.NotNil(t, result.Terminal.ExitCode)
	require.Equal(t, 0, *result.Terminal.ExitCode)
	require.Equal(t, []byte("hello\nworld\n"), result.Transcript)
	require.Equal(t, 1, pollCount)
}

func TestClientRunExecCaptureUsesRecoveryLookupAfterContextCancel(t *testing.T) {
	firstPoll := make(chan struct{}, 1)
	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-1/execs/stream":
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Run9-Exec-ID", "exec-1")
		case "/projects/default/workspace/execs/exec-1":
			pollCount++
			if pollCount == 1 {
				writeJSONResponse(t, w, http.StatusOK, ExecView{ExecID: "exec-1", State: "running"})
				firstPoll <- struct{}{}
				return
			}
			writeJSONResponse(t, w, http.StatusOK, ExecView{
				ExecID:   "exec-1",
				State:    "succeeded",
				ExitCode: intPtr(0),
			})
		case "/projects/default/workspace/execs/exec-1/log-download":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("done\n"))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-firstPoll
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	result, err := newProjectTestClient(t, server.URL, "default").RunExecCapture(ctx, "box-1", ExecRequest{
		Command: []string{"printf", "hello"},
	})
	require.NoError(t, err)
	require.Equal(t, "exec-1", result.ExecID)
	require.Equal(t, []byte("done\n"), result.Transcript)
	require.Equal(t, 2, pollCount)
}

func TestClientRunExecCaptureReturnsOutputUnavailableReason(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-1/execs/stream":
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Run9-Exec-ID", "exec-1")
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "exit", ExitCode: 0}))
		case "/projects/default/workspace/execs/exec-1/log-download":
			writeJSONResponse(t, w, http.StatusConflict, map[string]string{
				"error": foregroundExecLogUnavailableReason,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").RunExecCapture(context.Background(), "box-1", ExecRequest{
		Command: []string{"printf", "hello"},
	})
	require.NoError(t, err)
	require.Equal(t, "exec-1", result.ExecID)
	require.Nil(t, result.Transcript)
	require.Equal(t, foregroundExecLogUnavailableReason, result.TranscriptUnavailableReason)
}

func TestClientRunExecCaptureUsesEventExecIDWhenHeaderIsMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/default/workspace/boxes/box-1/execs/stream":
			http.Redirect(w, r, "/foreground-relay/execs/ticket-1/exec-stream", http.StatusSeeOther)
		case "/foreground-relay/execs/ticket-1/exec-stream":
			w.Header().Set("Content-Type", "application/x-ndjson")
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "started", ExecID: "exec-1"}))
			require.NoError(t, json.NewEncoder(w).Encode(ExecStreamEvent{Type: "exit", ExitCode: 0}))
		case "/projects/default/workspace/execs/exec-1/log-download":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("done\n"))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").RunExecCapture(context.Background(), "box-1", ExecRequest{
		Command: []string{"printf", "hello"},
	})
	require.NoError(t, err)
	require.Equal(t, "exec-1", result.ExecID)
	require.Equal(t, []byte("done\n"), result.Transcript)
}
