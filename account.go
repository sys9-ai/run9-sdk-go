package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/org_access"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// UpdateAccount updates mutable fields on the caller's user account.
func (c *Client) UpdateAccount(ctx context.Context, req UpdateMeRequest) (MeView, error) {
	payload, err := remarshalJSON[*genmodels.UpdateMePayload](req)
	if err != nil {
		return MeView{}, err
	}

	return generatedResult[MeView](c.portal.OrgAccess.UpdateAccountContext(ctx, &org_access.UpdateAccountParams{
		Request: payload,
	}, c.auth))
}

// ListSSHKeys lists SSH keys registered on the caller's account.
func (c *Client) ListSSHKeys(ctx context.Context) ([]SSHKeyView, error) {
	return generatedResult[[]SSHKeyView](c.portal.OrgAccess.ListAccountSSHKeysContext(ctx, &org_access.ListAccountSSHKeysParams{}, c.auth))
}

// CreateSSHKey registers one SSH public key on the caller's account.
func (c *Client) CreateSSHKey(ctx context.Context, req CreateSSHKeyRequest) (SSHKeyView, error) {
	payload, err := remarshalJSON[*genmodels.CreateSSHKeyPayload](req)
	if err != nil {
		return SSHKeyView{}, err
	}

	return generatedResult[SSHKeyView](c.portal.OrgAccess.CreateAccountSSHKeyContext(ctx, &org_access.CreateAccountSSHKeyParams{
		Request: payload,
	}, c.auth))
}

// DeleteSSHKey removes one SSH key from the caller's account.
func (c *Client) DeleteSSHKey(ctx context.Context, sshKeyID string) (SSHKeyView, error) {
	return generatedResult[SSHKeyView](c.portal.OrgAccess.DeleteAccountSSHKeyContext(ctx, &org_access.DeleteAccountSSHKeyParams{
		ID: strings.TrimSpace(sshKeyID),
	}, c.auth))
}
