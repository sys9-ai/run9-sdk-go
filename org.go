package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/org_access"
	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/org_members"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// UpdateOrg updates mutable fields on one organization.
func (c *Client) UpdateOrg(ctx context.Context, orgID string, req UpdateOrgRequest) (OrgView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateOrgPayload](req)
	if err != nil {
		return OrgView{}, err
	}

	return generatedResult[OrgView](c.portal.OrgMembers.UpdateOrgContext(ctx, &org_members.UpdateOrgParams{
		ID:      strings.TrimSpace(orgID),
		Request: payload,
	}, c.auth))
}

// DeleteOrg deletes one organization.
func (c *Client) DeleteOrg(ctx context.Context, orgID string) (DeleteOrgResult, error) {
	return generatedResult[DeleteOrgResult](c.portal.OrgMembers.DeleteOrgContext(ctx, &org_members.DeleteOrgParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth))
}

// ListOrgMembers lists members in one organization.
func (c *Client) ListOrgMembers(ctx context.Context, orgID string) ([]MembershipView, error) {
	return generatedResult[[]MembershipView](c.portal.OrgMembers.ListOrgMembersContext(ctx, &org_members.ListOrgMembersParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth))
}

// UpdateOrgMember updates one organization member.
func (c *Client) UpdateOrgMember(ctx context.Context, orgID string, userID string, req UpdateMembershipRequest) (MembershipView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateMembershipPayload](req)
	if err != nil {
		return MembershipView{}, err
	}

	return generatedResult[MembershipView](c.portal.OrgMembers.UpdateOrgMemberContext(ctx, &org_members.UpdateOrgMemberParams{
		ID:      strings.TrimSpace(orgID),
		UserID:  strings.TrimSpace(userID),
		Request: payload,
	}, c.auth))
}

// DeleteOrgMember removes one member from an organization.
func (c *Client) DeleteOrgMember(ctx context.Context, orgID string, userID string) error {
	return generatedAction(c.portal.OrgMembers.RemoveOrgMemberContext(ctx, &org_members.RemoveOrgMemberParams{
		ID:     strings.TrimSpace(orgID),
		UserID: strings.TrimSpace(userID),
	}, c.auth))
}

// ListInvitations lists invitations in one organization.
func (c *Client) ListInvitations(ctx context.Context, orgID string) ([]InvitationView, error) {
	return generatedResult[[]InvitationView](c.portal.OrgMembers.ListOrgInvitationsContext(ctx, &org_members.ListOrgInvitationsParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth))
}

// CreateInvitation creates one invitation in an organization.
func (c *Client) CreateInvitation(ctx context.Context, orgID string, req CreateInvitationRequest) (InvitationView, error) {
	payload, err := remarshalJSON[*genmodels.CreateInvitationPayload](req)
	if err != nil {
		return InvitationView{}, err
	}

	return generatedResult[InvitationView](c.portal.OrgMembers.CreateOrgInvitationContext(ctx, &org_members.CreateOrgInvitationParams{
		ID:      strings.TrimSpace(orgID),
		Request: payload,
	}, c.auth))
}

// RevokeInvitation revokes one organization invitation.
func (c *Client) RevokeInvitation(ctx context.Context, orgID string, invitationID string) (DeleteInvitationResult, error) {
	return generatedResult[DeleteInvitationResult](c.portal.OrgMembers.RevokeOrgInvitationContext(ctx, &org_members.RevokeOrgInvitationParams{
		ID:           strings.TrimSpace(orgID),
		InvitationID: strings.TrimSpace(invitationID),
	}, c.auth))
}

// ListAPIKeys lists API keys visible to the caller in the current organization.
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKeyView, error) {
	return generatedResult[[]APIKeyView](c.portal.OrgAccess.ListAPIKeysContext(ctx, &org_access.ListAPIKeysParams{}, c.auth))
}

// CreateAPIKey creates one API key in the current organization.
func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (CreatedAPIKeyView, error) {
	payload, err := remarshalJSON[*genmodels.CreateAPIKeyPayload](req)
	if err != nil {
		return CreatedAPIKeyView{}, err
	}

	return generatedResult[CreatedAPIKeyView](c.portal.OrgAccess.CreateAPIKeyContext(ctx, &org_access.CreateAPIKeyParams{
		Request: payload,
	}, c.auth))
}

// RevokeAPIKey revokes one API key in the current organization.
func (c *Client) RevokeAPIKey(ctx context.Context, apiKeyID string) (APIKeyView, error) {
	return generatedResult[APIKeyView](c.portal.OrgAccess.RevokeAPIKeyContext(ctx, &org_access.RevokeAPIKeyParams{
		ID: strings.TrimSpace(apiKeyID),
	}, c.auth))
}

// GetOrgHosts loads runtime hosts assigned to the current organization.
func (c *Client) GetOrgHosts(ctx context.Context) (OrgHostsView, error) {
	return generatedResult[OrgHostsView](c.portal.OrgAccess.ListOrgHostsContext(ctx, &org_access.ListOrgHostsParams{}, c.auth))
}
