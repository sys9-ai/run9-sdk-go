package run9

import (
	"context"
	"errors"
	"strings"
	"time"
)

// BackgroundExecFollower tracks one background exec cursor and hides repeated pull-output polling.
type BackgroundExecFollower struct {
	client *Client
	execID string
	cursor string
}

// FollowBackgroundExec returns a cursor-tracking follower for one background exec.
func (c *Client) FollowBackgroundExec(execID string) *BackgroundExecFollower {
	if c == nil {
		return nil
	}
	return &BackgroundExecFollower{
		client: c,
		execID: strings.TrimSpace(execID),
	}
}

// Cursor returns the current replay cursor.
func (f *BackgroundExecFollower) Cursor() string {
	if f == nil {
		return ""
	}
	return f.cursor
}

// SetCursor overwrites the current replay cursor.
func (f *BackgroundExecFollower) SetCursor(cursor string) {
	if f == nil {
		return
	}
	f.cursor = strings.TrimSpace(cursor)
}

// Reset clears the current replay cursor so the next read starts from the beginning once.
func (f *BackgroundExecFollower) Reset() {
	f.SetCursor("")
}

// Read waits for the next meaningful output window.
//
// When wait is positive, Read skips empty and started-only polls while the exec stays active.
func (f *BackgroundExecFollower) Read(ctx context.Context, wait time.Duration) (BackgroundExecPullOutput, error) {
	if f == nil || f.client == nil {
		return BackgroundExecPullOutput{}, errors.New("nil background exec follower")
	}

	deadline := time.Time{}
	if wait > 0 {
		deadline = time.Now().Add(wait)
	}

	requestCursor := f.cursor
	for {
		requestWait := wait
		if !deadline.IsZero() {
			remaining := time.Until(deadline)
			if remaining < 0 {
				remaining = 0
			}
			requestWait = remaining
		}

		result, err := f.client.PullBackgroundExecOutput(ctx, f.execID, PullBackgroundExecOutputRequest{
			Cursor: requestCursor,
			Wait:   requestWait,
		})
		if err != nil {
			return BackgroundExecPullOutput{}, err
		}

		f.cursor = strings.TrimSpace(result.NextCursor)
		if !shouldContinueBackgroundExecRead(result, requestCursor, deadline) {
			return result, nil
		}
		requestCursor = f.cursor
	}
}

// Pump reads one meaningful output window and writes stdout to Stdout and stderr plus notices to Stderr.
func (f *BackgroundExecFollower) Pump(ctx context.Context, wait time.Duration, writers ExecOutputWriters) (BackgroundExecPullOutput, error) {
	result, err := f.Read(ctx, wait)
	if err != nil {
		return BackgroundExecPullOutput{}, err
	}
	if err := result.WriteOutput(writers); err != nil {
		return BackgroundExecPullOutput{}, err
	}
	return result, nil
}

func shouldContinueBackgroundExecRead(result BackgroundExecPullOutput, requestCursor string, deadline time.Time) bool {
	if result.TerminalResult() != nil {
		return false
	}
	if deadline.IsZero() {
		return false
	}
	if time.Until(deadline) < time.Millisecond {
		return false
	}
	if !backgroundExecStateActive(result.State) {
		return false
	}
	nextCursor := strings.TrimSpace(result.NextCursor)
	if nextCursor == "" {
		return false
	}
	if len(result.Events) == 0 {
		return true
	}
	return backgroundExecEventsOnlyStarted(result.Events) && nextCursor != strings.TrimSpace(requestCursor)
}

func backgroundExecStateActive(state string) bool {
	switch strings.TrimSpace(state) {
	case "pending", "running":
		return true
	default:
		return false
	}
}

func backgroundExecEventsOnlyStarted(events []BackgroundExecOutputEvent) bool {
	if len(events) == 0 {
		return false
	}
	for _, event := range events {
		if event.Type != BackgroundExecOutputEventStarted {
			return false
		}
	}
	return true
}
