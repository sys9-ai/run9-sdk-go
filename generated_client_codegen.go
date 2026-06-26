package run9

// The SDK vendors its own portal Swagger snapshot so the generated layer can be
// refreshed from the standalone run9-sdk-go repository. A narrow normalization
// pass patches the snapshot before codegen to preserve PATCH request semantics
// that go-swagger would otherwise flatten away.
//go:generate ./scripts/generate-swagger-client.sh
