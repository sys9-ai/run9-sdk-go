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

	result, err := c.portal.OrgMembers.UpdateOrgContext(ctx, &org_members.UpdateOrgParams{
		ID:      strings.TrimSpace(orgID),
		Request: payload,
	}, c.auth)
	if err != nil {
		return OrgView{}, generatedError(err)
	}
	return remarshalJSON[OrgView](result.GetPayload())
}

// DeleteOrg deletes one organization.
func (c *Client) DeleteOrg(ctx context.Context, orgID string) (DeleteOrgResult, error) {
	result, err := c.portal.OrgMembers.DeleteOrgContext(ctx, &org_members.DeleteOrgParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth)
	if err != nil {
		return DeleteOrgResult{}, generatedError(err)
	}
	return remarshalJSON[DeleteOrgResult](result.GetPayload())
}

// ListOrgMembers lists members in one organization.
func (c *Client) ListOrgMembers(ctx context.Context, orgID string) ([]MembershipView, error) {
	result, err := c.portal.OrgMembers.ListOrgMembersContext(ctx, &org_members.ListOrgMembersParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]MembershipView](result.GetPayload())
}

// UpdateOrgMember updates one organization member.
func (c *Client) UpdateOrgMember(ctx context.Context, orgID string, userID string, req UpdateMembershipRequest) (MembershipView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateMembershipPayload](req)
	if err != nil {
		return MembershipView{}, err
	}

	result, err := c.portal.OrgMembers.UpdateOrgMemberContext(ctx, &org_members.UpdateOrgMemberParams{
		ID:      strings.TrimSpace(orgID),
		UserID:  strings.TrimSpace(userID),
		Request: payload,
	}, c.auth)
	if err != nil {
		return MembershipView{}, generatedError(err)
	}
	return remarshalJSON[MembershipView](result.GetPayload())
}

// DeleteOrgMember removes one member from an organization.
func (c *Client) DeleteOrgMember(ctx context.Context, orgID string, userID string) error {
	_, err := c.portal.OrgMembers.RemoveOrgMemberContext(ctx, &org_members.RemoveOrgMemberParams{
		ID:     strings.TrimSpace(orgID),
		UserID: strings.TrimSpace(userID),
	}, c.auth)
	return generatedError(err)
}

// ListInvitations lists invitations in one organization.
func (c *Client) ListInvitations(ctx context.Context, orgID string) ([]InvitationView, error) {
	result, err := c.portal.OrgMembers.ListOrgInvitationsContext(ctx, &org_members.ListOrgInvitationsParams{
		ID: strings.TrimSpace(orgID),
	}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]InvitationView](result.GetPayload())
}

// CreateInvitation creates one invitation in an organization.
func (c *Client) CreateInvitation(ctx context.Context, orgID string, req CreateInvitationRequest) (InvitationView, error) {
	payload, err := remarshalJSON[*genmodels.CreateInvitationPayload](req)
	if err != nil {
		return InvitationView{}, err
	}

	result, err := c.portal.OrgMembers.CreateOrgInvitationContext(ctx, &org_members.CreateOrgInvitationParams{
		ID:      strings.TrimSpace(orgID),
		Request: payload,
	}, c.auth)
	if err != nil {
		return InvitationView{}, generatedError(err)
	}
	return remarshalJSON[InvitationView](result.GetPayload())
}

// RevokeInvitation revokes one organization invitation.
func (c *Client) RevokeInvitation(ctx context.Context, orgID string, invitationID string) (DeleteInvitationResult, error) {
	result, err := c.portal.OrgMembers.RevokeOrgInvitationContext(ctx, &org_members.RevokeOrgInvitationParams{
		ID:           strings.TrimSpace(orgID),
		InvitationID: strings.TrimSpace(invitationID),
	}, c.auth)
	if err != nil {
		return DeleteInvitationResult{}, generatedError(err)
	}
	return remarshalJSON[DeleteInvitationResult](result.GetPayload())
}

// ListAPIKeys lists API keys visible to the caller in the current organization.
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKeyView, error) {
	result, err := c.portal.OrgAccess.ListAPIKeysContext(ctx, &org_access.ListAPIKeysParams{}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]APIKeyView](result.GetPayload())
}

// CreateAPIKey creates one API key in the current organization.
func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (CreatedAPIKeyView, error) {
	payload, err := remarshalJSON[*genmodels.CreateAPIKeyPayload](req)
	if err != nil {
		return CreatedAPIKeyView{}, err
	}

	result, err := c.portal.OrgAccess.CreateAPIKeyContext(ctx, &org_access.CreateAPIKeyParams{
		Request: payload,
	}, c.auth)
	if err != nil {
		return CreatedAPIKeyView{}, generatedError(err)
	}
	return remarshalJSON[CreatedAPIKeyView](result.GetPayload())
}

// RevokeAPIKey revokes one API key in the current organization.
func (c *Client) RevokeAPIKey(ctx context.Context, apiKeyID string) (APIKeyView, error) {
	result, err := c.portal.OrgAccess.RevokeAPIKeyContext(ctx, &org_access.RevokeAPIKeyParams{
		ID: strings.TrimSpace(apiKeyID),
	}, c.auth)
	if err != nil {
		return APIKeyView{}, generatedError(err)
	}
	return remarshalJSON[APIKeyView](result.GetPayload())
}

// GetOrgHosts loads runtime hosts assigned to the current organization.
func (c *Client) GetOrgHosts(ctx context.Context) (OrgHostsView, error) {
	result, err := c.portal.OrgAccess.ListOrgHostsContext(ctx, &org_access.ListOrgHostsParams{}, c.auth)
	if err != nil {
		return OrgHostsView{}, generatedError(err)
	}
	return remarshalJSON[OrgHostsView](result.GetPayload())
}
