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
							"network_mode": {
								SchemaProps: spec.SchemaProps{
									AllOf: []spec.Schema{{SchemaProps: spec.SchemaProps{Ref: spec.MustCreateRef("#/definitions/api.BoxNetworkMode")}}},
								},
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{"x-nullable": true},
								},
							},
							"security_mode": {
								SchemaProps: spec.SchemaProps{
									AllOf: []spec.Schema{{SchemaProps: spec.SchemaProps{Ref: spec.MustCreateRef("#/definitions/api.BoxSecurityMode")}}},
								},
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{"x-nullable": true},
								},
							},
						},
					},
				},
				"api.BoxNetworkMode": {
					SchemaProps: spec.SchemaProps{
						Type: []string{"string"},
						Enum: []any{"normal", "managed"},
					},
				},
				"api.BoxSecurityMode": {
					SchemaProps: spec.SchemaProps{
						Type: []string{"string"},
						Enum: []any{"standard", "restricted"},
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

	networkMode := swagger.Definitions["UpdateBoxPayload"].Properties["network_mode"]
	require.Empty(t, networkMode.AllOf)
	require.Equal(t, spec.StringOrArray{"string"}, networkMode.Type)
	require.Equal(t, []any{"normal", "managed"}, networkMode.Enum)
	require.Equal(t, true, networkMode.Extensions["x-nullable"])
	require.Equal(t, "BoxNetworkMode", xGoTypeName(t, networkMode))

	securityMode := swagger.Definitions["UpdateBoxPayload"].Properties["security_mode"]
	require.Empty(t, securityMode.AllOf)
	require.Equal(t, spec.StringOrArray{"string"}, securityMode.Type)
	require.Equal(t, []any{"standard", "restricted"}, securityMode.Enum)
	require.Equal(t, true, securityMode.Extensions["x-nullable"])
	require.Equal(t, "BoxSecurityMode", xGoTypeName(t, securityMode))
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
							"network_mode": {
								SchemaProps: spec.SchemaProps{},
							},
							"security_mode": {
								SchemaProps: spec.SchemaProps{},
							},
						},
					},
				},
				"api.BoxNetworkMode":  {},
				"api.BoxSecurityMode": {},
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
