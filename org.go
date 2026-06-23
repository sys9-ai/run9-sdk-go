package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) UpdateOrg(ctx context.Context, creds Credentials, orgID string, req UpdateOrgRequest) (OrgView, error) {
	var view OrgView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, ""), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) DeleteOrg(ctx context.Context, creds Credentials, orgID string) (DeleteOrgResult, error) {
	var view DeleteOrgResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, ""), creds, requestOptions{result: &view})
	return view, err
}

func (c *Client) OrgMembers(ctx context.Context, creds Credentials, orgID string) ([]MembershipView, error) {
	var views []MembershipView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/members"), creds, requestOptions{result: &views})
	return views, err
}

func (c *Client) UpdateOrgMember(ctx context.Context, creds Credentials, orgID string, userID string, req UpdateMembershipRequest) (MembershipView, error) {
	var view MembershipView
	err := c.do(ctx, http.MethodPatch, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) RemoveOrgMember(ctx context.Context, creds Credentials, orgID string, userID string) error {
	return c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/members/"+url.PathEscape(strings.TrimSpace(userID))), creds, requestOptions{})
}

func (c *Client) Invitations(ctx context.Context, creds Credentials, orgID string) ([]InvitationView, error) {
	var views []InvitationView
	err := c.do(ctx, http.MethodGet, c.orgPath(orgID, "/invitations"), creds, requestOptions{result: &views})
	return views, err
}

func (c *Client) CreateInvitation(ctx context.Context, creds Credentials, orgID string, req CreateInvitationRequest) (InvitationView, error) {
	var view InvitationView
	err := c.do(ctx, http.MethodPost, c.orgPath(orgID, "/invitations"), creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) RevokeInvitation(ctx context.Context, creds Credentials, orgID string, invitationID string) (DeleteInvitationResult, error) {
	var view DeleteInvitationResult
	err := c.do(ctx, http.MethodDelete, c.orgPath(orgID, "/invitations/"+url.PathEscape(strings.TrimSpace(invitationID))), creds, requestOptions{result: &view})
	return view, err
}

func (c *Client) APIKeys(ctx context.Context, creds Credentials) ([]APIKeyView, error) {
	var views []APIKeyView
	err := c.do(ctx, http.MethodGet, "/api-keys", creds, requestOptions{result: &views})
	return views, err
}

func (c *Client) CreateAPIKey(ctx context.Context, creds Credentials, req CreateAPIKeyRequest) (CreatedAPIKeyView, error) {
	var view CreatedAPIKeyView
	err := c.do(ctx, http.MethodPost, "/api-keys", creds, requestOptions{body: req, result: &view})
	return view, err
}

func (c *Client) RevokeAPIKey(ctx context.Context, creds Credentials, apiKeyID string) (APIKeyView, error) {
	var view APIKeyView
	err := c.do(ctx, http.MethodDelete, "/api-keys/"+url.PathEscape(strings.TrimSpace(apiKeyID)), creds, requestOptions{result: &view})
	return view, err
}

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
