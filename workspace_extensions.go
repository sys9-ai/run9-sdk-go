package run9

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) UpdateBox(ctx context.Context, creds Credentials, boxID string, req UpdateBoxRequest) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodPatch, c.workspacePath("/boxes/"+url.PathEscape(strings.TrimSpace(boxID))), creds, requestOptions{
		body:   req,
		result: &view,
	})
	return view, err
}

func (c *Client) SnapTree(ctx context.Context, creds Credentials, snapID string) (SnapTreeView, error) {
	var view SnapTreeView
	err := c.do(ctx, http.MethodGet, c.workspacePath("/snaps/"+url.PathEscape(strings.TrimSpace(snapID))+"/tree"), creds, requestOptions{result: &view})
	return view, err
}

func (c *Client) DownloadExecLog(ctx context.Context, creds Credentials, execID string) (io.ReadCloser, error) {
	resp, err := c.doRaw(ctx, http.MethodGet, c.workspacePath("/execs/"+url.PathEscape(strings.TrimSpace(execID))+"/log-download"), creds, requestOptions{})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
