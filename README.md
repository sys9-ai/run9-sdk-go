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
go get github.com/sys9-ai/run9-sdk-go@v0.1.1
```

The SDK uses semantic version tags. Depend on one released tag instead of an unpublished commit.

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

Run one foreground exec and read typed events:

```go
stream, err := project.StartExecStream(ctx, "devbox", run9.ExecRequest{
	Command: []string{"/bin/sh", "-lc", "echo hello"},
})
if err != nil {
	return err
}
defer stream.Close()

for {
	event, err := stream.ReadEvent()
	if err != nil {
		return err
	}

	switch event.Type {
	case "stdout":
		// handle stdout bytes
	case "stderr":
		// handle stderr bytes
	case "exit":
		return nil
	}
}
```

Start one background exec and poll merged output:

```go
execView, err := project.StartBackgroundExec(ctx, "devbox", run9.ExecRequest{
	Command:      []string{"/bin/sh", "-lc", "long task"},
	StdinEnabled: true,
})
if err != nil {
	return err
}

result, err := project.PullBackgroundExecOutput(ctx, execView.ExecID, run9.PullBackgroundExecOutputRequest{
	Wait: 2 * time.Second,
})
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
