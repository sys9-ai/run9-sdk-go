package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// ListProjectSecrets lists project-scoped secrets for the current project.
func (c *Client) ListProjectSecrets(ctx context.Context) ([]ProjectSecretView, error) {
	var views []ProjectSecretView
	path, err := projectPath(c.projectCID, "/secrets")
	if err != nil {
		return nil, err
	}
	err = c.do(ctx, http.MethodGet, path, requestOptions{result: &views})
	return views, err
}

// CreateProjectSecret creates one project-scoped secret.
func (c *Client) CreateProjectSecret(ctx context.Context, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	path, err := projectPath(c.projectCID, "/secrets")
	if err != nil {
		return view, err
	}
	err = c.do(ctx, http.MethodPost, path, requestOptions{body: req, result: &view})
	return view, err
}

// UpdateProjectSecret updates one project-scoped secret.
func (c *Client) UpdateProjectSecret(ctx context.Context, secretID string, req UpdateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	path, err := projectPath(c.projectCID, "/secrets/"+url.PathEscape(strings.TrimSpace(secretID)))
	if err != nil {
		return view, err
	}
	err = c.do(ctx, http.MethodPatch, path, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteProjectSecret deletes one project-scoped secret.
func (c *Client) DeleteProjectSecret(ctx context.Context, secretID string) error {
	path, err := projectPath(c.projectCID, "/secrets/"+url.PathEscape(strings.TrimSpace(secretID)))
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, requestOptions{})
}

// ListBoxSecrets lists box-scoped secrets for one box.
func (c *Client) ListBoxSecrets(ctx context.Context, boxID string) ([]ProjectSecretView, error) {
	var views []ProjectSecretView
	err := c.doWorkspace(ctx, http.MethodGet, boxSecretPath(boxID, ""), requestOptions{result: &views})
	return views, err
}

// CreateBoxSecret creates one box-scoped secret.
func (c *Client) CreateBoxSecret(ctx context.Context, boxID string, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.doWorkspace(ctx, http.MethodPost, boxSecretPath(boxID, ""), requestOptions{body: req, result: &view})
	return view, err
}

// UpdateBoxSecret updates one box-scoped secret.
func (c *Client) UpdateBoxSecret(ctx context.Context, boxID string, secretID string, req UpdateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.doWorkspace(ctx, http.MethodPatch, boxSecretPath(boxID, secretID), requestOptions{body: req, result: &view})
	return view, err
}

// DeleteBoxSecret deletes one box-scoped secret.
func (c *Client) DeleteBoxSecret(ctx context.Context, boxID string, secretID string) error {
	return c.doWorkspace(ctx, http.MethodDelete, boxSecretPath(boxID, secretID), requestOptions{})
}

func boxSecretPath(boxID string, secretID string) string {
	path := "/boxes/" + url.PathEscape(strings.TrimSpace(boxID)) + "/secrets"
	secretID = strings.TrimSpace(secretID)
	if secretID != "" {
		path += "/" + url.PathEscape(secretID)
	}
	return path
}
