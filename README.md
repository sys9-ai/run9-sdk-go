# run9-sdk-go

`run9-sdk-go` is the public Go SDK for the run9 control-plane API.

The module path is `github.com/sys9-ai/run9-sdk-go`. The package name is `run9`.

## Documentation

- Godoc: [pkg.go.dev/github.com/sys9-ai/run9-sdk-go](https://pkg.go.dev/github.com/sys9-ai/run9-sdk-go)
- Local docs: `go doc github.com/sys9-ai/run9-sdk-go`
- API examples: `go test ./...` compiles the examples shipped with the package

If the public API changes, README, doc comments, and examples are updated in the same change.

## Install

```bash
go get github.com/sys9-ai/run9-sdk-go@latest
```

The SDK uses semantic version tags. For reproducible builds, pin one released tag in your `go.mod` instead of an unpublished commit.

## Quick Start

```go
package main

import (
	"context"
	"log"

	run9 "github.com/sys9-ai/run9-sdk-go"
)

func main() {
	client, err := run9.NewClient("https://api.run.sys9.ai", run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	})
	if err != nil {
		log.Fatal(err)
	}

	project := client.WithProject("default")
	boxes, err := project.ListBoxes(context.Background(), run9.ListBoxesRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("loaded %d boxes", len(boxes))
}
```

## Calling Model

This SDK is intentionally opinionated. It uses one stable calling style across the package:

1. Create one client from endpoint and credentials.
2. Use the base client for global APIs such as account, org, shared snaps, and project discovery.
3. Use `client.WithProject(...)` before calling any project-scoped API.
4. Keep simple identity inputs positional, and use request structs for optional filters or mutable payloads.
5. Use typed streaming helpers instead of raw NDJSON bodies or websocket frames.

The SDK models the public control-plane contract. It does not include local config persistence, terminal rendering, shell completion, or local cursor persistence.

## Common Workflows

Load the current authenticated identity:

```go
identity, err := client.WhoAmI(ctx)
```

List accessible projects:

```go
projects, err := client.ListProjects(ctx)
```

Create one box in a project:

```go
project := client.WithProject("sandbox")

box, err := project.CreateBox(ctx, run9.CreateBoxRequest{
	BoxID:          "devbox",
	DesiredShape:   "2c4g",
	SourceImageRef: "public.ecr.aws/docker/library/alpine:3.20",
})
```

Run one foreground exec and stream its output:

```go
result, err := project.RunExec(ctx, "devbox", run9.ExecRequest{
	Command: []string{"/bin/sh", "-lc", "echo hello"},
}, run9.ExecOutputWriters{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
})
```

`result` tells you whether the exec exited, was cancelled, or failed. When you need direct event control, `StartExecStream(...)` and `ReadEvent()` are still available.

Run one foreground exec and wait for the final merged transcript:

```go
capture, err := project.RunExecCapture(ctx, "devbox", run9.ExecRequest{
	Command: []string{"/bin/sh", "-lc", "echo hello"},
})
if err != nil {
	return err
}
if capture.TranscriptUnavailableReason != "" {
	return fmt.Errorf("exec finished but its final transcript is unavailable: %s", capture.TranscriptUnavailableReason)
}
fmt.Printf("exit=%d transcript=%s", *capture.Terminal.ExitCode, capture.Transcript)
```

`RunExecCapture(...)` is the stable completion helper. It uses the live foreground stream to start the command, then falls back to durable exec truth plus `log-download` when the stream transport breaks near terminal. The returned transcript always follows the merged `log-download` view, so it does not preserve stdout and stderr as separate streams. If the exec finished but the transcript archive is unavailable, the helper still returns the terminal result and reports the transcript problem through `TranscriptUnavailableReason`.

Start one background exec and follow its output with an internal cursor:

```go
execView, err := project.StartBackgroundExec(ctx, "devbox", run9.ExecRequest{
	Command:      []string{"/bin/sh", "-lc", "long task"},
	StdinEnabled: true,
})
if err != nil {
	return err
}

follower := project.FollowBackgroundExec(execView.ExecID)
result, err := follower.Pump(ctx, 2*time.Second, run9.ExecOutputWriters{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
})
```

`result.NextCursor` is still exposed for explicit replay use cases, but callers that only want to keep tailing can let the follower manage it.

If you want one merged transcript in event order, call `Read(...)` and then `WriteMergedOutput(...)`:

```go
result, err := follower.Read(ctx, 2*time.Second)
if err != nil {
	return err
}
if err := result.WriteMergedOutput(os.Stdout, os.Stderr); err != nil {
	return err
}
```

## API Surface

Global APIs live on the base client:

- account identity and SSH keys
- org metadata, members, invitations, API keys, and org hosts
- project list, inspect, and create
- shared snap catalog list, inspect, publish, and delete

Project-scoped APIs require `WithProject(...)`:

- boxes
- snaps
- execs
- file archive upload and download
- project members
- project and box secrets
- create-from-shared-snap flows

If a project-scoped method is called without `WithProject(...)`, the SDK returns an explicit error instead of guessing a default project.

## Errors

Non-2xx responses are returned as `*run9.Error`.

```go
if err != nil {
	var apiErr *run9.Error
	if errors.As(err, &apiErr) {
		log.Printf("run9 request failed: status=%d message=%q", apiErr.StatusCode, apiErr.Message)
	}
}
```

The SDK avoids silent fallbacks. Missing credentials, invalid endpoints, and missing project scope are returned as direct errors.

## Compatibility

Before `v1.0.0`, breaking API cleanup is allowed when it materially improves clarity or ergonomics. Every exported API change must land together with:

- updated doc comments
- updated README examples
- passing `go test ./...`
