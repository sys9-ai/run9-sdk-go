package run9

// The SDK vendors its own portal Swagger snapshot so standalone run9-sdk-go
// releases can regenerate without the portal source tree. In the monorepo
// go generate ./portal/api refreshes this snapshot from
// portal/api/swagger/codegen.yaml. A narrow
// normalization pass patches only PATCH slice/map fields that go-swagger still
// flattens away, while enum PATCH semantics stay in the portal doc comments.
//go:generate ./scripts/generate-swagger-client.sh
