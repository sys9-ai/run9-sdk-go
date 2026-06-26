package run9

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/boxes"
	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/snaps"
	genmodels "github.com/sys9-ai/run9-sdk-go/internal/generated/models"
)

// UpdateBox updates mutable fields on one box.
func (c *Client) UpdateBox(ctx context.Context, boxID string, req UpdateBoxRequest) (BoxView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return BoxView{}, err
	}

	payload, err := remarshalJSON[*genmodels.UpdateBoxPayload](req)
	if err != nil {
		return BoxView{}, err
	}

	result, err := c.portal.Boxes.UpdateBoxContext(ctx, &boxes.UpdateBoxParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(boxID),
		Request:    payload,
	}, c.auth)
	if err != nil {
		return BoxView{}, generatedError(err)
	}
	return remarshalJSON[BoxView](result.GetPayload())
}

// GetSnapTree loads the ancestry tree for one snap.
func (c *Client) GetSnapTree(ctx context.Context, snapID string) (SnapTreeView, error) {
	projectCID, err := c.requireProjectCID()
	if err != nil {
		return SnapTreeView{}, err
	}

	result, err := c.portal.Snaps.GetSnapTreeContext(ctx, &snaps.GetSnapTreeParams{
		ProjectCid: projectCID,
		ID:         strings.TrimSpace(snapID),
	}, c.auth)
	if err != nil {
		return SnapTreeView{}, generatedError(err)
	}
	return remarshalJSON[SnapTreeView](result.GetPayload())
}

// DownloadExecLog downloads the stored log for one exec.
func (c *Client) DownloadExecLog(ctx context.Context, execID string) (io.ReadCloser, error) {
	resp, err := c.doWorkspaceRaw(ctx, http.MethodGet, "/execs/"+url.PathEscape(strings.TrimSpace(execID))+"/log-download", requestOptions{})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
