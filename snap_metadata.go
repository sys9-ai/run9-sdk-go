package run9

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// ForkSnapWithOptions captures the same filesystem as ForkSnap and atomically
// assigns the supplied labels to the new Snap. Source labels are not inherited.
func (c *Client) ForkSnapWithOptions(ctx context.Context, snapID string, req ForkSnapRequest) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodPost, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID))+"/fork", requestOptions{body: req, result: &view})
	return view, err
}

// UpdateSnap replaces Snap metadata without changing its filesystem or last-used time.
func (c *Client) UpdateSnap(ctx context.Context, snapID string, req UpdateSnapRequest) (SnapView, error) {
	var view SnapView
	err := c.doWorkspace(ctx, http.MethodPatch, "/snaps/"+url.PathEscape(strings.TrimSpace(snapID)), requestOptions{body: req, result: &view})
	return view, err
}
