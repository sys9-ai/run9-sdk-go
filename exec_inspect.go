package run9

import (
	"context"
	"strings"

	"github.com/sys9-ai/run9-sdk-go/internal/generated/client/execs"
)

// GetExec loads one exec by ID from the current project.
func (c *Client) GetExec(ctx context.Context, execID string) (ExecView, error) {
	return projectGeneratedResult[ExecView](c, func(projectCID string) (any, error) {
		return c.portal.Execs.GetExecContext(ctx, &execs.GetExecParams{
			ID:         strings.TrimSpace(execID),
			ProjectCid: projectCID,
		}, c.auth)
	})
}
