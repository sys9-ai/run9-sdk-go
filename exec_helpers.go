package run9

import (
	"context"
	"io"
)

// ExecOutputWriters selects where exec stdout, stderr, and output notices are written.
//
// Nil writers are treated as io.Discard.
type ExecOutputWriters struct {
	Stdout io.Writer
	Stderr io.Writer
}

// ExecTerminalStatus identifies how one exec stream reached a terminal state.
type ExecTerminalStatus string

const (
	// ExecTerminalStatusExited reports a normal process exit with an exit code.
	ExecTerminalStatusExited ExecTerminalStatus = "exited"
	// ExecTerminalStatusCancelled reports a terminal cancellation.
	ExecTerminalStatusCancelled ExecTerminalStatus = "cancelled"
	// ExecTerminalStatusError reports a terminal runtime or control-plane failure.
	ExecTerminalStatusError ExecTerminalStatus = "error"
)

// ExecTerminalResult describes the terminal outcome of one foreground, attached, or background exec stream.
type ExecTerminalResult struct {
	// Status identifies whether the exec exited, was cancelled, or failed.
	Status ExecTerminalStatus
	// ExitCode is set when the terminal outcome includes a process exit code.
	ExitCode *int
	// Reason carries the terminal cancellation or failure reason when one is available.
	Reason string
	// CacheAccess reports JuiceFS cache and object-store deltas when the runtime exposes them.
	CacheAccess *CacheAccessStats
	// PrewarmProfileID identifies the profile selected for this exec.
	PrewarmProfileID string
	// PrewarmProfileName is the profile name selected for this exec.
	PrewarmProfileName string
	// PrewarmRecording describes the artifact finalized by a recording exec.
	PrewarmRecording *PrewarmRecordingResult
}

type execEventReader interface {
	ReadEvent() (ExecStreamEvent, error)
	Close() error
}

func normalizeExecOutputWriters(writers ExecOutputWriters) ExecOutputWriters {
	if writers.Stdout == nil {
		writers.Stdout = io.Discard
	}
	if writers.Stderr == nil {
		writers.Stderr = io.Discard
	}
	return writers
}

func pumpExecEvents(ctx context.Context, reader execEventReader, writers ExecOutputWriters) (ExecTerminalResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	writers = normalizeExecOutputWriters(writers)
	stop := context.AfterFunc(ctx, func() {
		_ = reader.Close()
	})
	defer stop()
	defer reader.Close()

	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if ctx.Err() != nil {
				return ExecTerminalResult{}, ctx.Err()
			}
			return ExecTerminalResult{}, err
		}

		switch event.Type {
		case "keepalive", "started":
			continue
		case "stdout":
			if _, err := writers.Stdout.Write(event.Data); err != nil {
				return ExecTerminalResult{}, err
			}
		case "stderr":
			if _, err := writers.Stderr.Write(event.Data); err != nil {
				return ExecTerminalResult{}, err
			}
		default:
			if terminal, ok := terminalResultFromExecEvent(event); ok {
				return terminal, nil
			}
		}
	}
}
