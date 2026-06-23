package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ExecByID(ctx context.Context, creds Credentials, execID string) (ExecView, error) {
	var view ExecView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/execs/"+url.PathEscape(strings.TrimSpace(execID))), creds, requestOptions{result: &view})
	return view, err
}
