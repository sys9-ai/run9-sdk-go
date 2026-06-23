package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// GetExec loads one exec by ID from the current project.
func (c *Client) GetExec(ctx context.Context, execID string) (ExecView, error) {
	var view ExecView
	err := c.doWorkspace(ctx, http.MethodGet, "/execs/"+url.PathEscape(strings.TrimSpace(execID)), requestOptions{result: &view})
	return view, err
}
