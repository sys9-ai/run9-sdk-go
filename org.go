package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// UpdateOrg updates mutable fields on one organization.
func (c *Client) UpdateOrg(ctx context.Context, orgID string, req UpdateOrgRequest) (OrgView, error) {
	var view OrgView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, ""), requestOptions{body: req, result: &view})
	return view, err
}

// DeleteOrg deletes one organization.
func (c *Client) DeleteOrg(ctx context.Context, orgID string) (DeleteOrgResult, error) {
	var view DeleteOrgResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, ""), requestOptions{result: &view})
	return view, err
}

// ListOrgMembers lists members in one organization.
func (c *Client) ListOrgMembers(ctx context.Context, orgID string) ([]MembershipView, error) {
	var views []MembershipView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/members"), requestOptions{result: &views})
	return views, err
}

// UpdateOrgMember updates one organization member.
func (c *Client) UpdateOrgMember(ctx context.Context, orgID string, userID string, req UpdateMembershipRequest) (MembershipView, error) {
	var view MembershipView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), requestOptions{body: req, result: &view})
	return view, err
}

// DeleteOrgMember removes one member from an organization.
func (c *Client) DeleteOrgMember(ctx context.Context, orgID string, userID string) error {
	return c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), requestOptions{})
}

// ListInvitations lists invitations in one organization.
func (c *Client) ListInvitations(ctx context.Context, orgID string) ([]InvitationView, error) {
	var views []InvitationView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/invitations"), requestOptions{result: &views})
	return views, err
}

// CreateInvitation creates one invitation in an organization.
func (c *Client) CreateInvitation(ctx context.Context, orgID string, req CreateInvitationRequest) (InvitationView, error) {
	var view InvitationView
	err := c.do(ctx, http.MethodPost, c.orgPath(orgID, "/invitations"), requestOptions{body: req, result: &view})
	return view, err
}

// RevokeInvitation revokes one organization invitation.
func (c *Client) RevokeInvitation(ctx context.Context, orgID string, invitationID string) (DeleteInvitationResult, error) {
	var view DeleteInvitationResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/invitations/"+url.PathEscape(strings.TrimSpace(invitationID))), requestOptions{result: &view})
	return view, err
}

// ListAPIKeys lists API keys visible to the caller in the current organization.
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKeyView, error) {
	var views []APIKeyView
	err := c.do(ctx, http.MethodGet, "/api-keys", requestOptions{result: &views})
	return views, err
}

// CreateAPIKey creates one API key in the current organization.
func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (CreatedAPIKeyView, error) {
	var view CreatedAPIKeyView
	err := c.do(ctx, http.MethodPost, "/api-keys", requestOptions{body: req, result: &view})
	return view, err
}

// RevokeAPIKey revokes one API key in the current organization.
func (c *Client) RevokeAPIKey(ctx context.Context, apiKeyID string) (APIKeyView, error) {
	var view APIKeyView
	err := c.do(ctx, http.MethodDelete, "/api-keys/"+url.PathEscape(strings.TrimSpace(apiKeyID)), requestOptions{result: &view})
	return view, err
}

// GetOrgHosts loads runtime hosts assigned to the current organization.
func (c *Client) GetOrgHosts(ctx context.Context) (OrgHostsView, error) {
	var view OrgHostsView
	err := c.do(ctx, http.MethodGet, "/org-runtime/hosts", requestOptions{result: &view})
	return view, err
}

func (c *Client) orgPath(orgID string, suffix string) string {
	base := "/orgs/" + url.PathEscape(strings.TrimSpace(orgID))
	if strings.TrimSpace(suffix) == "" {
		return base
	}
	return base + "/" + strings.TrimLeft(strings.TrimSpace(suffix), "/")
}
