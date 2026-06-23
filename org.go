package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// UpdateOrg updates mutable fields on one organization.
func (c *Client) UpdateOrg(ctx context.Context, creds Credentials, orgID string, req UpdateOrgRequest) (OrgView, error) {
	var view OrgView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, ""), creds, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteOrg deletes one organization.
func (c *Client) DeleteOrg(ctx context.Context, creds Credentials, orgID string) (DeleteOrgResult, error) {
	var view DeleteOrgResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, ""), creds, requestOptions{result: &view})
	return view, err
}

// OrgMembers lists members in one organization.
func (c *Client) OrgMembers(ctx context.Context, creds Credentials, orgID string) ([]MembershipView, error) {
	var views []MembershipView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/members"), creds, requestOptions{result: &views})
	return views, err
}

// UpdateOrgMember updates one organization member.
func (c *Client) UpdateOrgMember(ctx context.Context, creds Credentials, orgID string, userID string, req UpdateMembershipRequest) (MembershipView, error) {
	var view MembershipView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{body: req, result: &view})
	return view, err
}

// RemoveOrgMember removes one member from an organization.
func (c *Client) RemoveOrgMember(ctx context.Context, creds Credentials, orgID string, userID string) error {
	return c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{})
}

// Invitations lists invitations in one organization.
func (c *Client) Invitations(ctx context.Context, creds Credentials, orgID string) ([]InvitationView, error) {
	var views []InvitationView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/invitations"), creds, requestOptions{result: &views})
	return views, err
}

// CreateInvitation creates one invitation in an organization.
func (c *Client) CreateInvitation(ctx context.Context, creds Credentials, orgID string, req CreateInvitationRequest) (InvitationView, error) {
	var view InvitationView
	err := c.do(ctx, http.MethodPost, c.orgPath(orgID, "/invitations"), creds, requestOptions{body: req, result: &view})
	return view, err
}

// RevokeInvitation revokes one organization invitation.
func (c *Client) RevokeInvitation(ctx context.Context, creds Credentials, orgID string, invitationID string) (DeleteInvitationResult, error) {
	var view DeleteInvitationResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/invitations/"+url.PathEscape(strings.TrimSpace(invitationID))), creds, requestOptions{result: &view})
	return view, err
}

// APIKeys lists API keys visible to the caller in the current organization.
func (c *Client) APIKeys(ctx context.Context, creds Credentials) ([]APIKeyView, error) {
	var views []APIKeyView
	err := c.do(ctx, http.MethodGet, "/api-keys", creds, requestOptions{result: &views})
	return views, err
}

// CreateAPIKey creates one API key in the current organization.
func (c *Client) CreateAPIKey(ctx context.Context, creds Credentials, req CreateAPIKeyRequest) (CreatedAPIKeyView, error) {
	var view CreatedAPIKeyView
	err := c.do(ctx, http.MethodPost, "/api-keys", creds, requestOptions{body: req, result: &view})
	return view, err
}

// RevokeAPIKey revokes one API key in the current organization.
func (c *Client) RevokeAPIKey(ctx context.Context, creds Credentials, apiKeyID string) (APIKeyView, error) {
	var view APIKeyView
	err := c.do(ctx, http.MethodDelete, "/api-keys/"+url.PathEscape(strings.TrimSpace(apiKeyID)), creds, requestOptions{result: &view})
	return view, err
}

// OrgHosts lists runtime hosts assigned to the current organization.
func (c *Client) OrgHosts(ctx context.Context, creds Credentials) (OrgHostsView, error) {
	var view OrgHostsView
	err := c.do(ctx, http.MethodGet, "/org-runtime/hosts", creds, requestOptions{result: &view})
	return view, err
}

func (c *Client) orgPath(orgID string, suffix string) string {
	base := "/orgs/" + url.PathEscape(strings.TrimSpace(orgID))
	if strings.TrimSpace(suffix) == "" {
		return base
	}
	return base + "/" + strings.TrimLeft(strings.TrimSpace(suffix), "/")
}
