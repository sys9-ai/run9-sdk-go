package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/projects"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// UpdateProject updates mutable fields on the current project.
func (c *Client) UpdateProject(ctx context.Context, req UpdateProjectRequest) (ProjectView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateProjectPayload](req)
	if err != nil {
		return ProjectView{}, err
	}

	return projectGeneratedResult[ProjectView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.UpdateProjectContext(ctx, &projects.UpdateProjectParams{
			ProjectCid: projectCID,
			Request:    payload,
		}, c.auth)
	})
}

// DeleteProject deletes the current project.
func (c *Client) DeleteProject(ctx context.Context) (DeleteProjectResult, error) {
	return projectGeneratedResult[DeleteProjectResult](c, func(projectCID string) (any, error) {
		return c.portal.Projects.DeleteProjectContext(ctx, &projects.DeleteProjectParams{
			ProjectCid: projectCID,
		}, c.auth)
	})
}

// ListProjectMembers lists members in the current project.
func (c *Client) ListProjectMembers(ctx context.Context) ([]ProjectMembershipView, error) {
	return projectGeneratedResult[[]ProjectMembershipView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.ListProjectMembersContext(ctx, &projects.ListProjectMembersParams{
			ProjectCid: projectCID,
		}, c.auth)
	})
}

// UpdateProjectMember updates one member in the current project.
func (c *Client) UpdateProjectMember(ctx context.Context, userID string, req UpdateProjectMembershipRequest) (ProjectMembershipView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateProjectMembershipPayload](req)
	if err != nil {
		return ProjectMembershipView{}, err
	}

	return projectGeneratedResult[ProjectMembershipView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.UpdateProjectMemberContext(ctx, &projects.UpdateProjectMemberParams{
			ProjectCid: projectCID,
			UserID:     strings.TrimSpace(userID),
			Request:    payload,
		}, c.auth)
	})
}

// DeleteProjectMember removes one member from the current project.
func (c *Client) DeleteProjectMember(ctx context.Context, userID string) error {
	return projectGeneratedAction(c, func(projectCID string) (any, error) {
		return c.portal.Projects.RemoveProjectMemberContext(ctx, &projects.RemoveProjectMemberParams{
			ProjectCid: projectCID,
			UserID:     strings.TrimSpace(userID),
		}, c.auth)
	})
}
