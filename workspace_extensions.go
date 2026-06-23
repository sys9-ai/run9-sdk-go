package run9

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// UpdateBox updates mutable fields on one box.
func (c *Client) UpdateBox(ctx context.Context, boxID string, req UpdateBoxRequest) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodPatch, "/boxes/"+url.PathEscape(strings.TrimSpace(boxID)), requestOptions{
		body:   req,
		result: &view,
	})
	return view, err
}

// GetSnapTree loads the ancestry tree for one snap.
func (c *Client) GetSnapTree(ctx context.Context, snapID string) (SnapTreeView, error) {
	var view SnapTreeView
	err := c.doWorkspace(ctx, http.MethodGet, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID))+"/tree", requestOptions{result: &view})
	return view, err
}

// DownloadExecLog downloads the stored log for one exec.
func (c *Client) DownloadExecLog(ctx context.Context, execID string) (io.ReadCloser, error) {
	resp, err := c.doWorkspaceRaw(ctx, http.MethodGet, "/execs/"+url.PathEscape(strings.TrimSpace(execID))+"/log-download", requestOptions{})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
