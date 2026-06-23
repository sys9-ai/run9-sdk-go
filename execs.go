package run9

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ListExecsRequest describes optional filters for ListExecs.
type ListExecsRequest struct {
	BoxID          string
	State          string
	Creator        string
	AcceptedAfter  *time.Time
	AcceptedBefore *time.Time
	Order          string
	Paged          bool
	Limit          *int
	Cursor         string
}

// ListExecsResult describes the response returned by ListExecs.
type ListExecsResult struct {
	Execs      []ExecView `json:"items"`
	NextCursor string     `json:"next_cursor"`
}

// ListExecs lists execs in the current project with optional filters.
func (c *Client) ListExecs(ctx context.Context, req ListExecsRequest) (ListExecsResult, error) {
	query := map[string]string{}
	if strings.TrimSpace(req.BoxID) != "" {
		query["box_id"] = strings.TrimSpace(req.BoxID)
	}
	if strings.TrimSpace(req.State) != "" {
		query["state"] = strings.TrimSpace(req.State)
	}
	if strings.TrimSpace(req.Creator) != "" {
		query["creator"] = strings.TrimSpace(req.Creator)
	}
	if req.AcceptedAfter != nil && !req.AcceptedAfter.IsZero() {
		query["accepted_after"] = req.AcceptedAfter.UTC().Format(time.RFC3339Nano)
	}
	if req.AcceptedBefore != nil && !req.AcceptedBefore.IsZero() {
		query["accepted_before"] = req.AcceptedBefore.UTC().Format(time.RFC3339Nano)
	}
	if strings.TrimSpace(req.Order) != "" {
		query["order"] = strings.TrimSpace(req.Order)
	}
	if req.Paged {
		query["paged"] = "true"
	}
	if req.Limit != nil {
		query["limit"] = strconv.Itoa(*req.Limit)
	}
	if strings.TrimSpace(req.Cursor) != "" {
		query["cursor"] = strings.TrimSpace(req.Cursor)
	}

	resp, err := c.doWorkspaceRaw(ctx, http.MethodGet, "/execs", requestOptions{query: query})
	if err != nil {
		return ListExecsResult{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListExecsResult{}, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return ListExecsResult{}, fmt.Errorf("portal api returned empty response body")
	}

	var views []ExecView
	if err := json.Unmarshal(data, &views); err != nil {
		return ListExecsResult{}, err
	}
	return ListExecsResult{
		Execs:      views,
		NextCursor: strings.TrimSpace(resp.Header.Get("X-Run9-Next-Cursor")),
	}, nil
}
