package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/shared_snaps"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// ListSharedSnaps lists shared snaps visible to the caller.
func (c *Client) ListSharedSnaps(ctx context.Context) ([]SharedSnapLineView, error) {
	return generatedResult[[]SharedSnapLineView](c.portal.SharedSnaps.ListSharedSnapsContext(ctx, &shared_snaps.ListSharedSnapsParams{}, c.auth))
}

// GetSharedSnap loads one shared snap and its versions.
func (c *Client) GetSharedSnap(ctx context.Context, name string) (SharedSnapDetailView, error) {
	return generatedResult[SharedSnapDetailView](c.portal.SharedSnaps.GetSharedSnapContext(ctx, &shared_snaps.GetSharedSnapParams{
		Name: strings.TrimSpace(name),
	}, c.auth))
}

// PublishSharedSnap publishes one snap into the shared snap catalog.
func (c *Client) PublishSharedSnap(ctx context.Context, req PublishSharedSnapRequest) (SharedSnapVersionView, error) {
	payload, err := remarshalJSON[*genmodels.PublishSharedSnapPayload](req)
	if err != nil {
		return SharedSnapVersionView{}, err
	}

	return generatedResult[SharedSnapVersionView](c.portal.SharedSnaps.PublishSharedSnapContext(ctx, &shared_snaps.PublishSharedSnapParams{
		Request: payload,
	}, c.auth))
}

// DeleteSharedSnap deletes one shared snap name and all of its versions.
func (c *Client) DeleteSharedSnap(ctx context.Context, name string) error {
	return generatedAction(c.portal.SharedSnaps.DeleteSharedSnapContext(ctx, &shared_snaps.DeleteSharedSnapParams{
		Name: strings.TrimSpace(name),
	}, c.auth))
}

// DeleteSharedSnapVersion deletes one shared snap version.
func (c *Client) DeleteSharedSnapVersion(ctx context.Context, name string, version int) error {
	return generatedAction(c.portal.SharedSnaps.DeleteSharedSnapVersionContext(ctx, &shared_snaps.DeleteSharedSnapVersionParams{
		Name:    strings.TrimSpace(name),
		Version: int64(version),
	}, c.auth))
}

// CreateBoxFromSharedSnap creates a box from one shared snap.
func (c *Client) CreateBoxFromSharedSnap(ctx context.Context, name string, req CreateBoxFromSharedSnapRequest) (BoxView, error) {
	payload, err := remarshalJSON[*genmodels.ConsumeSharedSnapToBoxPayload](req)
	if err != nil {
		return BoxView{}, err
	}

	return projectGeneratedResult[BoxView](c, func(projectCID string) (any, error) {
		return c.portal.SharedSnaps.ConsumeSharedSnapToBoxContext(ctx, &shared_snaps.ConsumeSharedSnapToBoxParams{
			Name:       strings.TrimSpace(name),
			ProjectCid: projectCID,
			Request:    payload,
		}, c.auth)
	})
}

// CreateSnapFromSharedSnap creates a snap from one shared snap.
func (c *Client) CreateSnapFromSharedSnap(ctx context.Context, name string, req CreateSnapFromSharedSnapRequest) (SnapView, error) {
	payload, err := remarshalJSON[*genmodels.ConsumeSharedSnapToSnapPayload](req)
	if err != nil {
		return SnapView{}, err
	}

	return projectGeneratedResult[SnapView](c, func(projectCID string) (any, error) {
		return c.portal.SharedSnaps.ConsumeSharedSnapToSnapContext(ctx, &shared_snaps.ConsumeSharedSnapToSnapParams{
			Name:       strings.TrimSpace(name),
			ProjectCid: projectCID,
			Request:    payload,
		}, c.auth)
	})
}
