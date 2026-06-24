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

	"github.com/stretchr/testify/require"
)

func TestClientPullBackgroundExecOutputReadsBinaryBodyAndHeaders(t *testing.T) {
	fixtureTime := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	responseBody, err := encodeBackgroundExecOutputEvents([]BackgroundExecOutputEvent{
		{Seq: 1, Type: BackgroundExecOutputEventStarted},
		{Seq: 2, Type: BackgroundExecOutputEventStdout, Data: []byte("hello")},
		{Seq: 3, Type: BackgroundExecOutputEventExit, ExitCode: intPtr(0)},
	})
	require.NoError(t, err)

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
		_, err := w.Write(responseBody)
		require.NoError(t, err)
	}))
	defer server.Close()

	result, err := newProjectTestClient(t, server.URL, "default").PullBackgroundExecOutput(context.Background(), "exec-1", PullBackgroundExecOutputRequest{
		Cursor: "cursor-1",
		Wait:   2 * time.Second,
	})
	require.NoError(t, err)
	require.Len(t, result.Events, 3)
	require.Equal(t, BackgroundExecOutputEventStarted, result.Events[0].Type)
	require.Equal(t, BackgroundExecOutputEventStdout, result.Events[1].Type)
	require.Equal(t, []byte("hello"), result.Events[1].Data)
	require.Equal(t, BackgroundExecOutputEventExit, result.Events[2].Type)
	require.NotNil(t, result.Events[2].ExitCode)
	require.Equal(t, 0, *result.Events[2].ExitCode)
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

func TestBackgroundExecPullOutputWriteOutputAndTerminalResult(t *testing.T) {
	result := BackgroundExecPullOutput{
		Events: []BackgroundExecOutputEvent{
			{Seq: 1, Type: BackgroundExecOutputEventStarted},
			{Seq: 2, Type: BackgroundExecOutputEventStdout, Data: []byte("hello\n")},
			{Seq: 3, Type: BackgroundExecOutputEventStderr, Data: []byte("warn\n")},
			{Seq: 4, Type: BackgroundExecOutputEventGap, GapBytes: 12},
			{Seq: 5, Type: BackgroundExecOutputEventTruncated},
			{Seq: 6, Type: BackgroundExecOutputEventExit, ExitCode: intPtr(23)},
		},
		State: "failed",
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := result.WriteOutput(ExecOutputWriters{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	require.NoError(t, err)
	require.Equal(t, "hello\n", stdout.String())
	require.Contains(t, stderr.String(), "warn\n")
	require.Contains(t, stderr.String(), "background exec omitted 12 bytes")
	require.Contains(t, stderr.String(), "background exec output was truncated")

	stdout.Reset()
	stderr.Reset()
	err = result.WriteMergedOutput(&stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, "hello\nwarn\n", stdout.String())
	require.Contains(t, stderr.String(), "background exec omitted 12 bytes")
	require.Contains(t, stderr.String(), "background exec output was truncated")

	terminal := result.TerminalResult()
	require.NotNil(t, terminal)
	require.Equal(t, ExecTerminalStatusExited, terminal.Status)
	require.NotNil(t, terminal.ExitCode)
	require.Equal(t, 23, *terminal.ExitCode)
}

func TestBackgroundExecFollowerReadSkipsStartedOnlyPolls(t *testing.T) {
	startedBody, err := encodeBackgroundExecOutputEvents([]BackgroundExecOutputEvent{
		{Seq: 1, Type: BackgroundExecOutputEventStarted},
	})
	require.NoError(t, err)

	finalBody, err := encodeBackgroundExecOutputEvents([]BackgroundExecOutputEvent{
		{Seq: 2, Type: BackgroundExecOutputEventStdout, Data: []byte("done\n")},
		{Seq: 3, Type: BackgroundExecOutputEventExit, ExitCode: intPtr(0)},
	})
	require.NoError(t, err)

	callIndex := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/projects/default/workspace/execs/exec-1/pull-output", r.URL.Path)
		callIndex++

		switch callIndex {
		case 1:
			w.Header().Set("X-Run9-Next-Cursor", "cursor-1")
			w.Header().Set("X-Run9-Exec-State", "running")
			_, err := w.Write(startedBody)
			require.NoError(t, err)
		case 2:
			w.Header().Set("X-Run9-Next-Cursor", "cursor-2")
			w.Header().Set("X-Run9-Exec-State", "failed")
			w.Header().Set("X-Run9-Exit-Code", "0")
			_, err := w.Write(finalBody)
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected pull-output call %d", callIndex)
		}
	}))
	defer server.Close()

	follower := newProjectTestClient(t, server.URL, "default").FollowBackgroundExec("exec-1")
	result, err := follower.Read(context.Background(), 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, "cursor-2", follower.Cursor())
	require.Len(t, result.Events, 2)
	require.Equal(t, BackgroundExecOutputEventStdout, result.Events[0].Type)
	require.Equal(t, []byte("done\n"), result.Events[0].Data)
	require.NotNil(t, result.TerminalResult())
	require.Equal(t, 2, callIndex)
}

func intPtr(value int) *int {
	return &value
}
