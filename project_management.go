package run9

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// UpdateProject updates mutable fields on the current project.
func (c *Client) UpdateProject(ctx context.Context, req UpdateProjectRequest) (ProjectView, error) {
	var view ProjectView
	path, err := projectPath(c.projectCID, "")
	if err != nil {
		return view, err
	}
	err = c.do(ctx, http.MethodPatch, path, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteProject deletes the current project.
func (c *Client) DeleteProject(ctx context.Context) (DeleteProjectResult, error) {
	var view DeleteProjectResult
	path, err := projectPath(c.projectCID, "")
	if err != nil {
		return view, err
	}
	err = c.do(ctx, http.MethodDelete, path, requestOptions{result: &view})
	return view, err
}

// ListProjectMembers lists members in the current project.
func (c *Client) ListProjectMembers(ctx context.Context) ([]ProjectMembershipView, error) {
	var views []ProjectMembershipView
	path, err := projectPath(c.projectCID, "/members")
	if err != nil {
		return nil, err
	}
	err = c.do(ctx, http.MethodGet, path, requestOptions{result: &views})
	return views, err
}

// UpdateProjectMember updates one member in the current project.
func (c *Client) UpdateProjectMember(ctx context.Context, userID string, req UpdateProjectMembershipRequest) (ProjectMembershipView, error) {
	var view ProjectMembershipView
	path, err := projectPath(c.projectCID, "/members/"+url.PathEscape(strings.TrimSpace(userID)))
	if err != nil {
		return view, err
	}
	err = c.do(ctx, http.MethodPatch, path, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteProjectMember removes one member from the current project.
func (c *Client) DeleteProjectMember(ctx context.Context, userID string) error {
	path, err := projectPath(c.projectCID, "/members/"+url.PathEscape(strings.TrimSpace(userID)))
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodDelete, path, requestOptions{})
}

func projectPath(projectCID string, suffix string) (string, error) {
	projectCID = strings.TrimSpace(projectCID)
	if projectCID == "" {
		return "", errors.New("missing project cid: use client.WithProject(...) for project-scoped APIs")
	}
	base := "/projects/" + url.PathEscape(strings.TrimSpace(projectCID))
	if strings.TrimSpace(suffix) == "" {
		return base, nil
	}
	return base + "/" + strings.TrimLeft(strings.TrimSpace(suffix), "/"), nil
}
