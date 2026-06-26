package run9

import (
	"context"
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
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectSecretView{}, err
	}

	payload, err := remarshalJSON[*genmodels.UpdateProjectSecretPayload](req)
	if err != nil {
		return ProjectSecretView{}, err
	}

	result, err := c.portal.Projects.UpdateProjectSecretContext(ctx, &projects.UpdateProjectSecretParams{
		ProjectCid: projectCID,
		SecretID:   strings.TrimSpace(secretID),
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectSecretView{}, generatedError(err)
	}
	return remarshalJSON[ProjectSecretView](result.GetPayload())
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
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectSecretView{}, err
	}

	payload, err := remarshalJSON[*genmodels.UpdateProjectSecretPayload](req)
	if err != nil {
		return ProjectSecretView{}, err
	}

	result, err := c.portal.Boxes.UpdateBoxSecretContext(ctx, &boxes.UpdateBoxSecretParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(boxID),
		SecretID:   strings.TrimSpace(secretID),
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectSecretView{}, generatedError(err)
	}
	return remarshalJSON[ProjectSecretView](result.GetPayload())
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
