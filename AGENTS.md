# AGENTS.md

Read this file when a task changes any of these:

- `run9-sdk-go` public API
- `run9-sdk-go` README, godoc, examples, or release tags
- `run9-cli`'s SDK dependency or release workflow

## API Rules

- Design from caller ergonomics, not from internal route names.
- Bind long-lived context once:
  - endpoint + credentials at `NewClient(...)`
  - project scope at `client.WithProject(...)`
- Keep simple identity inputs positional.
- Use request structs for optional filters and mutable payloads.
- Stream APIs must return typed readers, not raw transport bodies.
- Do not leak CLI-only concepts into the SDK.

## Documentation Rules

- Every exported type, function, and method must have a doc comment.
- README must explain the actual supported calling model, not the HTTP routes.
- Keep at least one compiling example for godoc via `go test`.
- If public API changes, update README and examples in the same change.

## Versioning And Release

- The SDK must use semantic version tags: `vX.Y.Z`.
- Tag only commits that are already on `main`.
- Before tagging, run `go test ./...`.
- Before releasing `run9-cli`, first release the SDK, then update CLI to that exact latest SDK tag.
- `run9-cli` main must track the latest released `run9-sdk-go` tag.

## Branch And Integration

- `main` is the only long-lived release branch.
- Feature work may use temporary branches, but SDK tags must point to `main`.
- When a run9 monorepo change touches both SDK and CLI:
  - land SDK API/docs/tests first
  - merge SDK to its own `main`
  - create and push the new SDK tag
  - update CLI to that exact tag
  - rerun CLI tests and release smoke
