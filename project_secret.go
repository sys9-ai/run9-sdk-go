package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ProjectSecrets(ctx context.Context, creds Credentials, projectCID string) ([]ProjectSecretView, error) {
	var views []ProjectSecretView
	err := c.do(ctx, http.MethodGet, projectPath(projectCID, "/secrets"), creds, requestOptions{result: &views})
	return views, err
}

func (c *Client) CreateProjectSecret(ctx context.Context, creds Credentials, projectCID string, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.do(ctx, http.MethodPost, projectPath(projectCID, "/secrets"), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) UpdateProjectSecret(ctx context.Context, creds Credentials, projectCID string, secretID string, req UpdateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.do(ctx, http.MethodPatch, projectPath(projectCID, "/secrets/"+url.PathEscape(strings.TrimSpace(secretID))), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) DeleteProjectSecret(ctx context.Context, creds Credentials, projectCID string, secretID string) error {
	return c.do(ctx, http.MethodDelete, projectPath(projectCID, "/secrets/"+url.PathEscape(strings.TrimSpace(secretID))), creds, requestOptions{})
}

func (c *Client) BoxSecrets(ctx context.Context, creds Credentials, boxID string) ([]ProjectSecretView, error) {
	var views []ProjectSecretView
	err := c.do(ctx, http.MethodGet, c.workspacePath(boxSecretPath(boxID, "")), creds, requestOptions{result: &views})
	return views, err
}

func (c *Client) CreateBoxSecret(ctx context.Context, creds Credentials, boxID string, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.do(ctx, http.MethodPost, c.workspacePath(boxSecretPath(boxID, "")), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) UpdateBoxSecret(ctx context.Context, creds Credentials, boxID string, secretID string, req UpdateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.do(ctx, http.MethodPatch, c.workspacePath(boxSecretPath(boxID, secretID)), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) DeleteBoxSecret(ctx context.Context, creds Credentials, boxID string, secretID string) error {
	return c.do(ctx, http.MethodDelete, c.workspacePath(boxSecretPath(boxID, secretID)), creds, requestOptions{})
}

func boxSecretPath(boxID string, secretID string) string {
	path := "/boxes/" + url.PathEscape(strings.TrimSpace(boxID)) + "/secrets"
	secretID = strings.TrimSpace(secretID)
	if secretID != "" {
		path += "/" + url.PathEscape(secretID)
	}
	return path
}
