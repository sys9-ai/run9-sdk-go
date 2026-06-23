package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// UpdateAccount updates mutable fields on the caller's user account.
func (c *Client) UpdateAccount(ctx context.Context, creds Credentials, req UpdateMeRequest) (MeView, error) {
	var view MeView
	err := c.do(ctx, http.MethodPatch, "/account", creds, requestOptions{body: req, result: &view})
	return view, err
}

// SSHKeys lists SSH keys registered on the caller's account.
func (c *Client) SSHKeys(ctx context.Context, creds Credentials) ([]SSHKeyView, error) {
	var views []SSHKeyView
	err := c.do(ctx, http.MethodGet, "/account/ssh-keys", creds, requestOptions{result: &views})
	return views, err
}

// CreateSSHKey registers one SSH public key on the caller's account.
func (c *Client) CreateSSHKey(ctx context.Context, creds Credentials, req CreateSSHKeyRequest) (SSHKeyView, error) {
	var view SSHKeyView
	err := c.do(ctx, http.MethodPost, "/account/ssh-keys", creds, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteSSHKey removes one SSH key from the caller's account.
func (c *Client) DeleteSSHKey(ctx context.Context, creds Credentials, sshKeyID string) (SSHKeyView, error) {
	var view SSHKeyView
	err := c.do(ctx, http.MethodDelete, "/account/ssh-keys/"+url.PathEscape(strings.TrimSpace(sshKeyID)), creds, requestOptions{result: &view})
	return view, err
}
