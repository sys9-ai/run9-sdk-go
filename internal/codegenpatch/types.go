// Package codegenpatch provides external schema types injected into the
// generated swagger client to preserve PATCH request semantics.
package codegenpatch

// StringSlice preserves omit-vs-empty semantics for optional string slice PATCH
// fields in generated request models.
type StringSlice []string

// StringMap preserves omit-vs-empty semantics for optional string map PATCH
// fields in generated request models.
type StringMap map[string]string
