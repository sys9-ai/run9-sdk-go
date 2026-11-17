package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/projects"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// ListPrivateNetworks lists the current project's private networks, oldest
// first.
func (c *Client) ListPrivateNetworks(ctx context.Context) ([]PrivateNetworkView, error) {
	return projectGeneratedResult[[]PrivateNetworkView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.ListPrivateNetworksContext(ctx, &projects.ListPrivateNetworksParams{
			ProjectCid: projectCID,
		}, c.auth)
	})
}

// GetPrivateNetwork loads one private network of the current project.
func (c *Client) GetPrivateNetwork(ctx context.Context, privateNetworkID string) (PrivateNetworkView, error) {
	return projectGeneratedResult[PrivateNetworkView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.GetPrivateNetworkContext(ctx, &projects.GetPrivateNetworkParams{
			ProjectCid:       projectCID,
			PrivateNetworkID: strings.TrimSpace(privateNetworkID),
		}, c.auth)
	})
}

// CreatePrivateNetwork adds a private network to the current project. It
// requires project admin, and returns ErrorCodePrivateNetworkGatewayNotConfigured
// when the deployment cannot run private networks.
func (c *Client) CreatePrivateNetwork(ctx context.Context, req CreatePrivateNetworkRequest) (PrivateNetworkView, error) {
	payload, err := remarshalJSON[*genmodels.CreatePrivateNetworkPayload](req)
	if err != nil {
		return PrivateNetworkView{}, err
	}

	return projectGeneratedResult[PrivateNetworkView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.CreatePrivateNetworkContext(ctx, &projects.CreatePrivateNetworkParams{
			ProjectCid: projectCID,
			Request:    payload,
		}, c.auth)
	})
}

// UpdatePrivateNetwork changes the set fields of one private network. It
// requires project admin.
func (c *Client) UpdatePrivateNetwork(ctx context.Context, privateNetworkID string, req UpdatePrivateNetworkRequest) (PrivateNetworkView, error) {
	payload, err := remarshalJSON[*genmodels.UpdatePrivateNetworkPayload](req)
	if err != nil {
		return PrivateNetworkView{}, err
	}

	return projectGeneratedResult[PrivateNetworkView](c, func(projectCID string) (any, error) {
		return c.portal.Projects.UpdatePrivateNetworkContext(ctx, &projects.UpdatePrivateNetworkParams{
			ProjectCid:       projectCID,
			PrivateNetworkID: strings.TrimSpace(privateNetworkID),
			Request:          payload,
		}, c.auth)
	})
}

// DeletePrivateNetwork removes one private network from the current project.
// It requires project admin.
func (c *Client) DeletePrivateNetwork(ctx context.Context, privateNetworkID string) error {
	return projectGeneratedAction(c, func(projectCID string) (any, error) {
		return c.portal.Projects.DeletePrivateNetworkContext(ctx, &projects.DeletePrivateNetworkParams{
			ProjectCid:       projectCID,
			PrivateNetworkID: strings.TrimSpace(privateNetworkID),
		}, c.auth)
	})
}
