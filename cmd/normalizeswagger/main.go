package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

const codegenPatchImportPath = "github.com/sys9-ai/run9-sdk-go/internal/codegenpatch"

func main() {
	inputPath := flag.String("in", "", "input swagger file")
	outputPath := flag.String("out", "", "output normalized swagger file")
	flag.Parse()

	if err := run(*inputPath, *outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(inputPath string, outputPath string) error {
	if inputPath == "" {
		return errors.New("missing -in")
	}
	if outputPath == "" {
		return errors.New("missing -out")
	}

	doc, err := loads.Spec(inputPath)
	if err != nil {
		return err
	}

	swagger := doc.Spec()
	if err := normalize(swagger); err != nil {
		return err
	}

	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, data, 0o644)
}

func normalize(swagger *spec.Swagger) error {
	if swagger == nil {
		return errors.New("missing swagger document")
	}

	if err := normalizePatchPayloads(swagger); err != nil {
		return err
	}
	return nil
}

func normalizePatchPayloads(swagger *spec.Swagger) error {
	if err := expectStringArrayProperty(swagger, "UpdateProjectSecretPayload", "allowed_hosts"); err != nil {
		return err
	}
	if err := setPropertyCodegenType(swagger, "UpdateProjectSecretPayload", "allowed_hosts", "StringSlice", "array"); err != nil {
		return err
	}
	if err := expectStringMapProperty(swagger, "UpdateBoxPayload", "labels"); err != nil {
		return err
	}
	if err := setPropertyCodegenType(swagger, "UpdateBoxPayload", "labels", "StringMap", "object"); err != nil {
		return err
	}
	return nil
}

func setPropertyCodegenType(swagger *spec.Swagger, definitionName string, propertyName string, typeName string, kind string) error {
	definition, property, err := definitionProperty(swagger, definitionName, propertyName)
	if err != nil {
		return err
	}

	applyCodegenType(property, typeName, kind)
	definition.Properties[propertyName] = *property
	swagger.Definitions[definitionName] = *definition
	return nil
}

func definitionProperty(swagger *spec.Swagger, definitionName string, propertyName string) (*spec.Schema, *spec.Schema, error) {
	definition, ok := swagger.Definitions[definitionName]
	if !ok {
		return nil, nil, fmt.Errorf("missing definition %q", definitionName)
	}

	property, ok := definition.Properties[propertyName]
	if !ok {
		return nil, nil, fmt.Errorf("missing property %q on definition %q", propertyName, definitionName)
	}

	return &definition, &property, nil
}

func expectStringArrayProperty(swagger *spec.Swagger, definitionName string, propertyName string) error {
	_, property, err := definitionProperty(swagger, definitionName, propertyName)
	if err != nil {
		return err
	}
	if propertyType(property) != "array" {
		return fmt.Errorf("unexpected %s.%s type %q", definitionName, propertyName, propertyType(property))
	}
	if property.Items == nil || property.Items.Schema == nil || propertyType(property.Items.Schema) != "string" {
		return fmt.Errorf("unexpected %s.%s items shape", definitionName, propertyName)
	}
	return nil
}

func expectStringMapProperty(swagger *spec.Swagger, definitionName string, propertyName string) error {
	_, property, err := definitionProperty(swagger, definitionName, propertyName)
	if err != nil {
		return err
	}
	if propertyType(property) != "object" {
		return fmt.Errorf("unexpected %s.%s type %q", definitionName, propertyName, propertyType(property))
	}
	if property.AdditionalProperties == nil || property.AdditionalProperties.Schema == nil || propertyType(property.AdditionalProperties.Schema) != "string" {
		return fmt.Errorf("unexpected %s.%s additionalProperties shape", definitionName, propertyName)
	}
	return nil
}

func applyCodegenType(schema *spec.Schema, typeName string, kind string) {
	if schema.Extensions == nil {
		schema.Extensions = spec.Extensions{}
	}
	schema.Extensions["x-omitempty"] = true
	schema.Extensions["x-go-type"] = map[string]any{
		"type": typeName,
		"import": map[string]any{
			"package": codegenPatchImportPath,
		},
		"hints": map[string]any{
			"kind":         kind,
			"noValidation": true,
		},
	}
}

func propertyType(schema *spec.Schema) string {
	if schema == nil || len(schema.Type) == 0 {
		return ""
	}
	return schema.Type[0]
}
