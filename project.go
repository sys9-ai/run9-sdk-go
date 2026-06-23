package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Projects lists projects visible to the caller.
func (c *Client) Projects(ctx context.Context, creds Credentials) ([]ProjectView, error) {
	var views []ProjectView
	err := c.do(ctx, http.MethodGet, "/projects", creds, requestOptions{result: &views})
	return views, err
}

// CreateProject creates one project.
func (c *Client) CreateProject(ctx context.Context, creds Credentials, req CreateProjectRequest) (ProjectView, error) {
	var view ProjectView
	err := c.do(ctx, http.MethodPost, "/projects", creds, requestOptions{body: req, result: &view})
	return view, err
}

// ProjectByCID loads one project by CID.
func (c *Client) ProjectByCID(ctx context.Context, creds Credentials, projectCID string) (ProjectView, error) {
	var view ProjectView
	err := c.do(ctx, http.MethodGet, "/projects/"+url.PathEscape(strings.TrimSpace(projectCID)), creds, requestOptions{result: &view})
	return view, err
}
