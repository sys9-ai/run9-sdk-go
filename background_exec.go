package run9

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/execs"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

const idleDeadlineAtHeader = "X-Run9-Idle-Deadline-At"

type backgroundExecPullOutputRequest struct {
	Cursor string `json:"cursor,omitempty"`
	WaitMS int64  `json:"wait_ms,omitempty"`
}

// StartBackgroundExec starts one background exec in a box.
func (c *Client) StartBackgroundExec(ctx context.Context, boxID string, req ExecRequest) (ExecView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ExecView{}, err
	}

	payload, err := remarshalJSON[*genmodels.BackgroundExecBoxPayload](req)
	if err != nil {
		return ExecView{}, err
	}

	result, err := c.portal.Execs.BackgroundExecBoxContext(ctx, &execs.BackgroundExecBoxParams{
		ID:         strings.TrimSpace(boxID),
		ProjectCid: projectCID,
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ExecView{}, generatedError(err)
	}
	return remarshalJSON[ExecView](result.GetPayload())
}

// PullBackgroundExecOutput polls output and state transitions for one background exec.
func (c *Client) PullBackgroundExecOutput(ctx context.Context, execID string, req PullBackgroundExecOutputRequest) (BackgroundExecPullOutput, error) {
	var result BackgroundExecPullOutput
	request := backgroundExecPullOutputRequest{
		Cursor: strings.TrimSpace(req.Cursor),
	}
	if req.Wait > 0 {
		request.WaitMS = req.Wait.Milliseconds()
	}
	resp, err := c.doWorkspaceRaw(ctx, http.MethodPost, "/execs/"+url.PathEscape(strings.TrimSpace(execID))+"/pull-output", requestOptions{
		body: request,
	})
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	if len(body) != 0 {
		result.Events, err = decodeBackgroundExecOutputEvents(body)
		if err != nil {
			return result, err
		}
	}
	result.NextCursor = strings.TrimSpace(resp.Header.Get("X-Run9-Next-Cursor"))
	result.State = strings.TrimSpace(resp.Header.Get("X-Run9-Exec-State"))
	if raw := strings.TrimSpace(resp.Header.Get("X-Run9-Exit-Code")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return result, err
		}
		result.ExitCode = &value
	}
	result.Reason = strings.TrimSpace(resp.Header.Get("X-Run9-Reason"))
	result.IdleDeadlineAt, err = parseOptionalIdleDeadlineAt(resp.Header)
	if err != nil {
		return result, err
	}
	return result, nil
}

// WriteBackgroundExecStdin writes stdin bytes into one background exec.
func (c *Client) WriteBackgroundExecStdin(ctx context.Context, execID string, req WriteBackgroundExecStdinRequest) (*time.Time, error) {
	resp, err := c.doWorkspaceRaw(ctx, http.MethodPost, "/execs/"+url.PathEscape(strings.TrimSpace(execID))+"/write-stdin", requestOptions{
		body: req.Data,
		headers: map[string]string{
			"Content-Type":       "application/octet-stream",
			"X-Run9-Close-Stdin": strconv.FormatBool(req.CloseStdin),
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	idleDeadlineAt, err := parseOptionalIdleDeadlineAt(resp.Header)
	if err != nil {
		return nil, err
	}
	return idleDeadlineAt, nil
}

// KillBackgroundExec requests termination of one background exec.
func (c *Client) KillBackgroundExec(ctx context.Context, execID string) error {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return err
	}

	_, err = c.portal.Execs.KillBackgroundExecContext(ctx, &execs.KillBackgroundExecParams{
		ID:         strings.TrimSpace(execID),
		ProjectCid: projectCID,
	}, c.auth)
	return generatedError(err)
}

func parseOptionalIdleDeadlineAt(headers http.Header) (*time.Time, error) {
	raw := strings.TrimSpace(headers.Get(idleDeadlineAtHeader))
	if raw == "" {
		return nil, nil
	}
	value, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
