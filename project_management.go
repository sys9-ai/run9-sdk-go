package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// UpdateProject updates mutable fields on one project.
func (c *Client) UpdateProject(ctx context.Context, creds Credentials, projectCID string, req UpdateProjectRequest) (ProjectView, error) {
	var view ProjectView
	err := c.do(ctx, http.MethodPatch, projectPath(projectCID, ""), creds, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteProject deletes one project.
func (c *Client) DeleteProject(ctx context.Context, creds Credentials, projectCID string) (DeleteProjectResult, error) {
	var view DeleteProjectResult
	err := c.do(ctx, http.MethodDelete, projectPath(projectCID, ""), creds, requestOptions{result: &view})
	return view, err
}

// ProjectMembers lists members in one project.
func (c *Client) ProjectMembers(ctx context.Context, creds Credentials, projectCID string) ([]ProjectMembershipView, error) {
	var views []ProjectMembershipView
	err := c.do(ctx, http.MethodGet, projectPath(projectCID, "/members"), creds, requestOptions{result: &views})
	return views, err
}

// UpdateProjectMember updates one project member.
func (c *Client) UpdateProjectMember(ctx context.Context, creds Credentials, projectCID string, userID string, req UpdateProjectMembershipRequest) (ProjectMembershipView, error) {
	var view ProjectMembershipView
	err := c.do(ctx, http.MethodPatch, projectPath(projectCID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{body: req, result: &view})
	return view, err
}

// RemoveProjectMember removes one member from a project.
func (c *Client) RemoveProjectMember(ctx context.Context, creds Credentials, projectCID string, userID string) error {
	return c.do(ctx, http.MethodDelete, projectPath(projectCID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{})
}

func projectPath(projectCID string, suffix string) string {
	base := "/projects/" + url.PathEscape(strings.TrimSpace(projectCID))
	if strings.TrimSpace(suffix) == "" {
		return base
	}
	return base + "/" + strings.TrimLeft(strings.TrimSpace(suffix), "/")
}
