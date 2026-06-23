package run9

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SharedSnaps lists shared snaps visible to the caller.
func (c *Client) SharedSnaps(ctx context.Context, creds Credentials) ([]SharedSnapLineView, error) {
	var views []SharedSnapLineView
	err := c.do(ctx, http.MethodGet, "/shared-snaps", creds, requestOptions{result: &views})
	return views, err
}

// SharedSnapDetail loads one shared snap and its versions.
func (c *Client) SharedSnapDetail(ctx context.Context, creds Credentials, name string) (SharedSnapDetailView, error) {
	var view SharedSnapDetailView
	err := c.do(ctx, http.MethodGet, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name)), creds, requestOptions{result: &view})
	return view, err
}

// PublishSharedSnap publishes one snap into the shared snap catalog.
func (c *Client) PublishSharedSnap(ctx context.Context, creds Credentials, req PublishSharedSnapRequest) (SharedSnapVersionView, error) {
	var view SharedSnapVersionView
	err := c.do(ctx, http.MethodPost, "/shared-snaps", creds, requestOptions{body: req, result: &view})
	return view, err
}

// DeleteSharedSnap deletes one shared snap name and all of its versions.
func (c *Client) DeleteSharedSnap(ctx context.Context, creds Credentials, name string) error {
	return c.do(ctx, http.MethodDelete, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name)), creds, requestOptions{})
}

// DeleteSharedSnapVersion deletes one shared snap version.
func (c *Client) DeleteSharedSnapVersion(ctx context.Context, creds Credentials, name string, version int) error {
	return c.do(ctx, http.MethodDelete, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/versions/"+strconv.Itoa(version), creds, requestOptions{})
}

// CreateBoxFromSharedSnap creates a box from one shared snap.
func (c *Client) CreateBoxFromSharedSnap(ctx context.Context, creds Credentials, name string, req CreateBoxFromSharedSnapRequest) (BoxView, error) {
	var view BoxView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/boxes"), creds, requestOptions{body: req, result: &view})
	return view, err
}

// CreateSnapFromSharedSnap creates a snap from one shared snap.
func (c *Client) CreateSnapFromSharedSnap(ctx context.Context, creds Credentials, name string, req CreateSnapFromSharedSnapRequest) (SnapView, error) {
	var view SnapView
	err := c.do(ctx, http.MethodPost, c.workspacePath("/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/snaps"), creds, requestOptions{body: req, result: &view})
	return view, err
}
