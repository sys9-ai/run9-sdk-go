package run9

// The SDK vendors its own portal Swagger snapshot so standalone run9-sdk-go
// releases can regenerate without the portal source tree. In the monorepo
// go generate ./portal/api refreshes this snapshot from
// portal/api/swagger/codegen.yaml. A narrow
// normalization pass preserves PATCH slice/map omission and nullable reference
// pointers that go-swagger otherwise flattens away. Enum PATCH semantics stay
// in the portal doc comments.
//go:generate ./scripts/generate-swagger-client.sh
