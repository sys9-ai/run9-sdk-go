package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// UpdateAccount updates mutable fields on the caller's user account.
func (c *Client) UpdateAccount(ctx context.Context, req UpdateMeRequest) (MeView, error) {
	var view MeView
	err := c.do(ctx, http.MethodPatch, "/account", requestOptions{body: req, result: &view})
	return view, err
}

// ListSSHKeys lists SSH keys registered on the caller's account.
func (c *Client) ListSSHKeys(ctx context.Context) ([]SSHKeyView, error) {
	var views []SSHKeyView
	err := c.do(ctx, http.MethodGet, "/account/ssh-keys", requestOptions{result: &views})
	return views, err
}

// CreateSSHKey registers one SSH public key on the caller's account.
func (c *Client) CreateSSHKey(ctx context.Context, req CreateSSHKeyRequest) (SSHKeyView, error) {
	var view SSHKeyView
	err := c.do(ctx, http.MethodPost, "/account/ssh-keys", requestOptions{body: req, result: &view})
	return view, err
}

// DeleteSSHKey removes one SSH key from the caller's account.
func (c *Client) DeleteSSHKey(ctx context.Context, sshKeyID string) (SSHKeyView, error) {
	var view SSHKeyView
	err := c.do(ctx, http.MethodDelete, "/account/ssh-keys/"+url.PathEscape(strings.TrimSpace(sshKeyID)), requestOptions{result: &view})
	return view, err
}
