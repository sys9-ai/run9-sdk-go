package run9

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	foregroundExecRecoveryPollInterval = 2 * time.Second
	foregroundExecRecoveryTimeout      = 5 * time.Second

	foregroundExecLogUnavailableReason = "foreground exec log archive is not available"
	execOutputUnavailableReason        = "exec output is no longer available"
)

// ExecCapture describes one foreground exec after the SDK has observed a terminal result.
//
// Transcript is the final merged log-download snapshot when it is available. It follows the
// control-plane log-download contract, so stdout and stderr are already merged and the body
// may include omission or truncation markers for bounded retention.
type ExecCapture struct {
	// ExecID is the durable foreground exec identifier.
	ExecID string
	// Terminal is the normalized terminal outcome.
	Terminal ExecTerminalResult
	// Transcript is the final merged transcript when log-download is available.
	Transcript []byte
	// TranscriptUnavailableReason explains why Transcript could not be loaded after the exec reached a terminal state.
	TranscriptUnavailableReason string
}

// RunExecCapture starts one foreground exec, waits for a terminal result, and loads the final merged log snapshot.
//
// The helper uses the normal foreground live stream to start the command. If that transport breaks
// after the exec has been accepted, it polls the durable exec record and still returns a terminal
// result when one becomes visible. It does not change foreground exec disconnect semantics, and the
// returned Transcript intentionally follows the merged log-download view instead of preserving stdout and
// stderr as separate streams.
func (c *Client) RunExecCapture(ctx context.Context, boxID string, req ExecRequest) (ExecCapture, error) {
	stream, err := c.StartExecStream(ctx, boxID, req)
	if err != nil {
		return ExecCapture{}, err
	}
	defer stream.Close()

	execID := strings.TrimSpace(stream.ExecID)
	for {
		event, err := stream.ReadEvent()
		if err != nil {
			if execID == "" {
				return ExecCapture{}, errors.New("foreground exec stream missing exec id")
			}
			return c.recoverExecCapture(ctx, execID, err)
		}
		if execID == "" {
			execID = strings.TrimSpace(event.ExecID)
		}

		terminal, ok := terminalResultFromExecEvent(event)
		if !ok {
			continue
		}
		if execID == "" {
			return ExecCapture{}, errors.New("foreground exec stream missing exec id")
		}
		return c.finishExecCapture(ctx, execID, terminal)
	}
}

func terminalResultFromExecEvent(event ExecStreamEvent) (ExecTerminalResult, bool) {
	switch event.Type {
	case "exit":
		exitCode := int(event.ExitCode)
		return ExecTerminalResult{
			Status:   ExecTerminalStatusExited,
			ExitCode: &exitCode,
		}, true
	case "cancelled":
		return ExecTerminalResult{
			Status: ExecTerminalStatusCancelled,
			Reason: event.CancelReason,
		}, true
	case "error":
		return ExecTerminalResult{
			Status: ExecTerminalStatusError,
			Reason: event.FailureReason,
		}, true
	default:
		return ExecTerminalResult{}, false
	}
}

func (c *Client) recoverExecCapture(ctx context.Context, execID string, streamErr error) (ExecCapture, error) {
	view, err := c.waitExecTerminal(ctx, execID)
	if err != nil {
		if ctx != nil && ctx.Err() != nil {
			return ExecCapture{}, ctx.Err()
		}
		return ExecCapture{}, fmt.Errorf("recover foreground exec %s after stream failure: %w", execID, err)
	}

	terminal := view.TerminalResult()
	if terminal == nil {
		return ExecCapture{}, fmt.Errorf("foreground exec %s never reached a terminal result", execID)
	}
	result, err := c.finishExecCapture(ctx, execID, *terminal)
	if err != nil {
		return ExecCapture{}, err
	}
	return result, nil
}

func (c *Client) waitExecTerminal(ctx context.Context, execID string) (ExecView, error) {
	view, err := c.pollExecUntilTerminal(ctx, execID)
	if err == nil {
		return view, nil
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		return ExecView{}, err
	}

	recoveryCtx, cancel := foregroundExecRecoveryContext(ctx)
	defer cancel()
	view, recoveryErr := c.pollExecUntilTerminal(recoveryCtx, execID)
	if recoveryErr == nil {
		return view, nil
	}
	if ctx != nil && ctx.Err() != nil {
		return ExecView{}, ctx.Err()
	}
	return ExecView{}, recoveryErr
}

func (c *Client) pollExecUntilTerminal(ctx context.Context, execID string) (ExecView, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ticker := time.NewTicker(foregroundExecRecoveryPollInterval)
	defer ticker.Stop()

	for {
		view, err := c.GetExec(ctx, execID)
		if err != nil {
			return ExecView{}, err
		}
		if view.TerminalResult() != nil {
			return view, nil
		}

		select {
		case <-ctx.Done():
			return ExecView{}, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *Client) finishExecCapture(ctx context.Context, execID string, terminal ExecTerminalResult) (ExecCapture, error) {
	outputCtx, cancel := foregroundExecRecoveryContext(ctx)
	defer cancel()

	result := ExecCapture{
		ExecID:                      execID,
		Terminal:                    terminal,
		Transcript:                  nil,
		TranscriptUnavailableReason: "",
	}

	body, err := c.DownloadExecLog(outputCtx, execID)
	if err != nil {
		if reason, ok := execLogDownloadUnavailableReason(err); ok {
			result.TranscriptUnavailableReason = reason
			return result, nil
		}
		return ExecCapture{}, err
	}
	defer body.Close()

	output, err := io.ReadAll(body)
	if err != nil {
		return ExecCapture{}, err
	}
	result.Transcript = output
	return result, nil
}

func foregroundExecRecoveryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), foregroundExecRecoveryTimeout)
	}
	return context.WithTimeout(context.WithoutCancel(ctx), foregroundExecRecoveryTimeout)
}

func execLogDownloadUnavailableReason(err error) (string, bool) {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return "", false
	}

	message := strings.TrimSpace(apiErr.Message)
	switch message {
	case foregroundExecLogUnavailableReason, execOutputUnavailableReason:
		return message, true
	default:
		return "", false
	}
}
