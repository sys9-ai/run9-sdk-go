package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/projects"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// ListProjects lists projects visible to the caller.
func (c *Client) ListProjects(ctx context.Context) ([]ProjectView, error) {
	return generatedResult[[]ProjectView](c.portal.Projects.ListProjectsContext(ctx, &projects.ListProjectsParams{}, c.auth))
}

// CreateProject creates one project.
func (c *Client) CreateProject(ctx context.Context, req CreateProjectRequest) (ProjectView, error) {
	payload, err := remarshalJSON[*genmodels.CreateProjectPayload](req)
	if err != nil {
		return ProjectView{}, err
	}

	return generatedResult[ProjectView](c.portal.Projects.CreateProjectContext(ctx, &projects.CreateProjectParams{
		Request: payload,
	}, c.auth))
}

// GetProject loads one project by CID.
func (c *Client) GetProject(ctx context.Context, projectCID string) (ProjectView, error) {
	return generatedResult[ProjectView](c.portal.Projects.GetProjectContext(ctx, &projects.GetProjectParams{
		ProjectCid: strings.TrimSpace(projectCID),
	}, c.auth))
}
