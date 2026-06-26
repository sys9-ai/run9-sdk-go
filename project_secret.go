package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/boxes"
	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/projects"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// ListProjectSecrets lists project-scoped secrets for the current project.
func (c *Client) ListProjectSecrets(ctx context.Context) ([]ProjectSecretView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return nil, err
	}

	result, err := c.portal.Projects.ListProjectSecretsContext(ctx, &projects.ListProjectSecretsParams{
		ProjectCid: projectCID,
	}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]ProjectSecretView](result.GetPayload())
}

// CreateProjectSecret creates one project-scoped secret.
func (c *Client) CreateProjectSecret(ctx context.Context, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectSecretView{}, err
	}

	payload, err := remarshalJSON[*genmodels.CreateProjectSecretPayload](req)
	if err != nil {
		return ProjectSecretView{}, err
	}

	result, err := c.portal.Projects.CreateProjectSecretContext(ctx, &projects.CreateProjectSecretParams{
		ProjectCid: projectCID,
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectSecretView{}, generatedError(err)
	}
	return remarshalJSON[ProjectSecretView](result.GetPayload())
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
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return err
	}

	_, err = c.portal.Projects.DeleteProjectSecretContext(ctx, &projects.DeleteProjectSecretParams{
		ProjectCid: projectCID,
		SecretID:   strings.TrimSpace(secretID),
	}, c.auth)
	return generatedError(err)
}

// ListBoxSecrets lists box-scoped secrets for one box.
func (c *Client) ListBoxSecrets(ctx context.Context, boxID string) ([]ProjectSecretView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return nil, err
	}

	result, err := c.portal.Boxes.ListBoxSecretsContext(ctx, &boxes.ListBoxSecretsParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(boxID),
	}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]ProjectSecretView](result.GetPayload())
}

// CreateBoxSecret creates one box-scoped secret.
func (c *Client) CreateBoxSecret(ctx context.Context, boxID string, req CreateProjectSecretRequest) (ProjectSecretView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectSecretView{}, err
	}

	payload, err := remarshalJSON[*genmodels.CreateProjectSecretPayload](req)
	if err != nil {
		return ProjectSecretView{}, err
	}

	result, err := c.portal.Boxes.CreateBoxSecretContext(ctx, &boxes.CreateBoxSecretParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(boxID),
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectSecretView{}, generatedError(err)
	}
	return remarshalJSON[ProjectSecretView](result.GetPayload())
}

// UpdateBoxSecret updates one box-scoped secret.
func (c *Client) UpdateBoxSecret(ctx context.Context, boxID string, secretID string, req UpdateProjectSecretRequest) (ProjectSecretView, error) {
	var view ProjectSecretView
	err := c.doWorkspace(ctx, http.MethodPatch, boxSecretPath(boxID, secretID), requestOptions{body: req, result: &view})
	return view, err
}

// DeleteBoxSecret deletes one box-scoped secret.
func (c *Client) DeleteBoxSecret(ctx context.Context, boxID string, secretID string) error {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return err
	}

	_, err = c.portal.Boxes.DeleteBoxSecretContext(ctx, &boxes.DeleteBoxSecretParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(boxID),
		SecretID:   strings.TrimSpace(secretID),
	}, c.auth)
	return generatedError(err)
}

func boxSecretPath(boxID string, secretID string) string {
	path := "/boxes/" + url.PathEscape(strings.TrimSpace(boxID)) + "/secrets"
	secretID = strings.TrimSpace(secretID)
	if secretID != "" {
		path += "/" + url.PathEscape(secretID)
	}
	return path
}
