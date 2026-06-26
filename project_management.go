package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/projects"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// UpdateProject updates mutable fields on the current project.
func (c *Client) UpdateProject(ctx context.Context, req UpdateProjectRequest) (ProjectView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectView{}, err
	}

	payload, err := remarshalJSON[*genmodels.UpdateProjectPayload](req)
	if err != nil {
		return ProjectView{}, err
	}

	result, err := c.portal.Projects.UpdateProjectContext(ctx, &projects.UpdateProjectParams{
		ProjectCid: projectCID,
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectView{}, generatedError(err)
	}
	return remarshalJSON[ProjectView](result.GetPayload())
}

// DeleteProject deletes the current project.
func (c *Client) DeleteProject(ctx context.Context) (DeleteProjectResult, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return DeleteProjectResult{}, err
	}

	result, err := c.portal.Projects.DeleteProjectContext(ctx, &projects.DeleteProjectParams{
		ProjectCid: projectCID,
	}, c.auth)
	if err != nil {
		return DeleteProjectResult{}, generatedError(err)
	}
	return remarshalJSON[DeleteProjectResult](result.GetPayload())
}

// ListProjectMembers lists members in the current project.
func (c *Client) ListProjectMembers(ctx context.Context) ([]ProjectMembershipView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return nil, err
	}

	result, err := c.portal.Projects.ListProjectMembersContext(ctx, &projects.ListProjectMembersParams{
		ProjectCid: projectCID,
	}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]ProjectMembershipView](result.GetPayload())
}

// UpdateProjectMember updates one member in the current project.
func (c *Client) UpdateProjectMember(ctx context.Context, userID string, req UpdateProjectMembershipRequest) (ProjectMembershipView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return ProjectMembershipView{}, err
	}

	payload, err := remarshalJSON[*genmodels.UpdateProjectMembershipPayload](req)
	if err != nil {
		return ProjectMembershipView{}, err
	}

	result, err := c.portal.Projects.UpdateProjectMemberContext(ctx, &projects.UpdateProjectMemberParams{
		ProjectCid: projectCID,
		UserID:     strings.TrimSpace(userID),
		Request:    payload,
	}, c.auth)
	if err != nil {
		return ProjectMembershipView{}, generatedError(err)
	}
	return remarshalJSON[ProjectMembershipView](result.GetPayload())
}

// DeleteProjectMember removes one member from the current project.
func (c *Client) DeleteProjectMember(ctx context.Context, userID string) error {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return err
	}

	_, err = c.portal.Projects.RemoveProjectMemberContext(ctx, &projects.RemoveProjectMemberParams{
		ProjectCid: projectCID,
		UserID:     strings.TrimSpace(userID),
	}, c.auth)
	return generatedError(err)
}
