package main

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestNormalizePatchPayloads(t *testing.T) {
	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: map[string]spec.Schema{
				"UpdateProjectSecretPayload": {
					SchemaProps: spec.SchemaProps{
						Properties: map[string]spec.Schema{
							"allowed_hosts": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"array"},
									Items: &spec.SchemaOrArray{
										Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
									},
								},
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{"x-nullable": true},
								},
							},
						},
					},
				},
				"UpdateBoxPayload": {
					SchemaProps: spec.SchemaProps{
						Properties: map[string]spec.Schema{
							"labels": {
								SchemaProps: spec.SchemaProps{
									Type: []string{"object"},
									AdditionalProperties: &spec.SchemaOrBool{
										Allows: true,
										Schema: &spec.Schema{SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
									},
								},
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{"x-nullable": true},
								},
							},
						},
					},
				},
			},
		},
	}

	require.NoError(t, normalize(swagger))

	allowedHosts := swagger.Definitions["UpdateProjectSecretPayload"].Properties["allowed_hosts"]
	require.Equal(t, true, allowedHosts.Extensions["x-omitempty"])
	require.Equal(t, "StringSlice", xGoTypeName(t, allowedHosts))

	labels := swagger.Definitions["UpdateBoxPayload"].Properties["labels"]
	require.Equal(t, true, labels.Extensions["x-omitempty"])
	require.Equal(t, "StringMap", xGoTypeName(t, labels))
}

func TestNormalizeRejectsUnexpectedShapes(t *testing.T) {
	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: map[string]spec.Schema{
				"UpdateProjectSecretPayload": {
					SchemaProps: spec.SchemaProps{
						Properties: map[string]spec.Schema{
							"allowed_hosts": {
								SchemaProps: spec.SchemaProps{Type: []string{"string"}},
							},
						},
					},
				},
				"UpdateBoxPayload": {
					SchemaProps: spec.SchemaProps{
						Properties: map[string]spec.Schema{
							"labels": {
								SchemaProps: spec.SchemaProps{Type: []string{"object"}},
							},
						},
					},
				},
			},
		},
	}

	err := normalize(swagger)
	require.Error(t, err)
}

func xGoTypeName(t *testing.T, schema spec.Schema) string {
	t.Helper()

	raw, ok := schema.Extensions["x-go-type"]
	require.True(t, ok)
	value, ok := raw.(map[string]any)
	require.True(t, ok)
	typeName, ok := value["type"].(string)
	require.True(t, ok)
	return typeName
}
