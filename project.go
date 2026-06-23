package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// ListProjects lists projects visible to the caller.
func (c *Client) ListProjects(ctx context.Context) ([]ProjectView, error) {
	var views []ProjectView
	err := c.do(ctx, http.MethodGet, "/projects", requestOptions{result: &views})
	return views, err
}

// CreateProject creates one project.
func (c *Client) CreateProject(ctx context.Context, req CreateProjectRequest) (ProjectView, error) {
	var view ProjectView
	err := c.do(ctx, http.MethodPost, "/projects", requestOptions{body: req, result: &view})
	return view, err
}

// GetProject loads one project by CID.
func (c *Client) GetProject(ctx context.Context, projectCID string) (ProjectView, error) {
	var view ProjectView
	err := c.do(ctx, http.MethodGet, "/projects/"+url.PathEscape(strings.TrimSpace(projectCID)), requestOptions{result: &view})
	return view, err
}
