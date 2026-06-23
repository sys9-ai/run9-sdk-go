package run9

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ListSharedSnaps lists shared snaps visible to the caller.
func (c *Client) ListSharedSnaps(ctx context.Context) ([]SharedSnapLineView, error) {
	var views []SharedSnapLineView
	err := c.do(ctx, http.MethodGet, "/shared-snaps", requestOptions{result: &views})
	return views, err
}

// GetSharedSnap loads one shared snap and its versions.
func (c *Client) GetSharedSnap(ctx context.Context, name string) (SharedSnapDetailView, error) {
	var view SharedSnapDetailView
	err := c.do(ctx, http.MethodGet, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name)), requestOptions{result: &view})
	return view, err
}

// PublishSharedSnap publishes one snap into the shared snap catalog.
func (c *Client) PublishSharedSnap(ctx context.Context, req PublishSharedSnapRequest) (SharedSnapVersionView, error) {
	var view SharedSnapVersionView
	err := c.do(ctx, http.MethodPost, "/shared-snaps", requestOptions{body: req, result: &view})
	return view, err
}

// DeleteSharedSnap deletes one shared snap name and all of its versions.
func (c *Client) DeleteSharedSnap(ctx context.Context, name string) error {
	return c.do(ctx, http.MethodDelete, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name)), requestOptions{})
}

// DeleteSharedSnapVersion deletes one shared snap version.
func (c *Client) DeleteSharedSnapVersion(ctx context.Context, name string, version int) error {
	return c.do(ctx, http.MethodDelete, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/versions/"+strconv.Itoa(version), requestOptions{})
}

// CreateBoxFromSharedSnap creates a box from one shared snap.
func (c *Client) CreateBoxFromSharedSnap(ctx context.Context, name string, req CreateBoxFromSharedSnapRequest) (BoxView, error) {
	var view BoxView
	err := c.doWorkspace(ctx, http.MethodPost, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/boxes", requestOptions{body: req, result: &view})
	return view, err
}

// CreateSnapFromSharedSnap creates a snap from one shared snap.
func (c *Client) CreateSnapFromSharedSnap(ctx context.Context, name string, req CreateSnapFromSharedSnapRequest) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodPost, "/shared-snaps/"+url.PathEscape(strings.TrimSpace(name))+"/snaps", requestOptions{body: req, result: &view})
	return view, err
}
