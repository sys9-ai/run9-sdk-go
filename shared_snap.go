package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/shared_snaps"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// ListSharedSnaps lists shared snaps visible to the caller.
func (c *Client) ListSharedSnaps(ctx context.Context) ([]SharedSnapLineView, error) {
	result, err := c.portal.SharedSnaps.ListSharedSnapsContext(ctx, &shared_snaps.ListSharedSnapsParams{}, c.auth)
	if err != nil {
		return nil, generatedError(err)
	}
	return remarshalJSON[[]SharedSnapLineView](result.GetPayload())
}

// GetSharedSnap loads one shared snap and its versions.
func (c *Client) GetSharedSnap(ctx context.Context, name string) (SharedSnapDetailView, error) {
	result, err := c.portal.SharedSnaps.GetSharedSnapContext(ctx, &shared_snaps.GetSharedSnapParams{
		Name: strings.TrimSpace(name),
	}, c.auth)
	if err != nil {
		return SharedSnapDetailView{}, generatedError(err)
	}
	return remarshalJSON[SharedSnapDetailView](result.GetPayload())
}

// PublishSharedSnap publishes one snap into the shared snap catalog.
func (c *Client) PublishSharedSnap(ctx context.Context, req PublishSharedSnapRequest) (SharedSnapVersionView, error) {
	payload, err := remarshalJSON[*genmodels.PublishSharedSnapPayload](req)
	if err != nil {
		return SharedSnapVersionView{}, err
	}

	result, err := c.portal.SharedSnaps.PublishSharedSnapContext(ctx, &shared_snaps.PublishSharedSnapParams{
		Request: payload,
	}, c.auth)
	if err != nil {
		return SharedSnapVersionView{}, generatedError(err)
	}
	return remarshalJSON[SharedSnapVersionView](result.GetPayload())
}

// DeleteSharedSnap deletes one shared snap name and all of its versions.
func (c *Client) DeleteSharedSnap(ctx context.Context, name string) error {
	_, err := c.portal.SharedSnaps.DeleteSharedSnapContext(ctx, &shared_snaps.DeleteSharedSnapParams{
		Name: strings.TrimSpace(name),
	}, c.auth)
	return generatedError(err)
}

// DeleteSharedSnapVersion deletes one shared snap version.
func (c *Client) DeleteSharedSnapVersion(ctx context.Context, name string, version int) error {
	_, err := c.portal.SharedSnaps.DeleteSharedSnapVersionContext(ctx, &shared_snaps.DeleteSharedSnapVersionParams{
		Name:    strings.TrimSpace(name),
		Version: int64(version),
	}, c.auth)
	return generatedError(err)
}

// CreateBoxFromSharedSnap creates a box from one shared snap.
func (c *Client) CreateBoxFromSharedSnap(ctx context.Context, name string, req CreateBoxFromSharedSnapRequest) (BoxView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return BoxView{}, err
	}

	payload, err := remarshalJSON[*genmodels.ConsumeSharedSnapToBoxPayload](req)
	if err != nil {
		return BoxView{}, err
	}

	result, err := c.portal.SharedSnaps.ConsumeSharedSnapToBoxContext(ctx, &shared_snaps.ConsumeSharedSnapToBoxParams{
		Name:       strings.TrimSpace(name),
		ProjectCid: projectCID,
		Request:    payload,
	}, c.auth)
	if err != nil {
		return BoxView{}, generatedError(err)
	}
	return remarshalJSON[BoxView](result.GetPayload())
}

// CreateSnapFromSharedSnap creates a snap from one shared snap.
func (c *Client) CreateSnapFromSharedSnap(ctx context.Context, name string, req CreateSnapFromSharedSnapRequest) (SnapView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return SnapView{}, err
	}

	payload, err := remarshalJSON[*genmodels.ConsumeSharedSnapToSnapPayload](req)
	if err != nil {
		return SnapView{}, err
	}

	result, err := c.portal.SharedSnaps.ConsumeSharedSnapToSnapContext(ctx, &shared_snaps.ConsumeSharedSnapToSnapParams{
		Name:       strings.TrimSpace(name),
		ProjectCid: projectCID,
		Request:    payload,
	}, c.auth)
	if err != nil {
		return SnapView{}, generatedError(err)
	}
	return remarshalJSON[SnapView](result.GetPayload())
}
