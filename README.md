# run9-sdk-go

Go SDK for the run9 control-plane API.

This repository is the public home for the shared run9 Go client used by
`run9-cli` and other programmatic integrations.

## Scope

The SDK intentionally exposes the public control-plane contract:

- API key authenticated HTTP client
- project-scoped workspace requests
- box / snap / org / project / shared-snap views and mutations
- foreground exec attach and inline stream paths
- background exec control paths
- archive upload / download helpers

It does not expose CLI-only concerns such as local config persistence, text
output formatting, or shell completion.

## Example

```go
package main

import (
	"context"
	"log"

	run9 "github.com/sys9-ai/run9-sdk-go"
)

func main() {
	client := run9.NewClient("https://api.run.sys9.ai").WithProject("default")
	boxes, err := client.Boxes(context.Background(), run9.Credentials{
		AK: "ak-...",
		SK: "sk-...",
	}, "", "", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d boxes", len(boxes))
}
```
