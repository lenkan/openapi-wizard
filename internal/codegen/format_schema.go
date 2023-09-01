package codegen

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lenkan/openapi-wizard/internal/openapi"
	"golang.org/x/exp/slices"
)

func MapShapes(schemas []openapi.JsonSchemaDefinition) []string {
	types := []string{}

	for _, s := range schemas {
		result := FormatSchemaShape(s)

		if !slices.Contains(types, result) {
			types = append(types, result)
		}
	}

	return types
}

func FormatSchemaShape(schema openapi.JsonSchemaDefinition) string {
	if len(schema.AllOf) > 0 {
		types := MapShapes(schema.AllOf)
		return "(" + strings.Join(types, " & ") + ")"
	}

	if len(schema.OneOf) > 0 {
		types := MapShapes(schema.OneOf)
		return "(" + strings.Join(types, " | ") + ")"
	}

	if schema.Type == "string" {
		return "string"
	}

	if schema.Type == "number" || schema.Type == "integer" {
		return "number"
	}

	if schema.Type == "object" {
		props := []string{}

		for propertyName, propertyDefinition := range schema.Properties {
			suffix := ""

			if slices.Contains(schema.Required, propertyName) == false {
				suffix += "?"
			}

			props = append(props, propertyName+suffix+": "+FormatSchemaShape(propertyDefinition))
		}

		if len(props) == 0 && schema.AdditionalProperties == true {
			return "Record<string, unknown>"
		}

		if len(props) == 0 {
			return "Record<string, never>"
		}

		return strings.Join([]string{"{", strings.Join(props, ";"), "}"}, "")
	}

	if schema.Ref != "" {
		return strcase.ToCamel(strings.Replace(schema.Ref, "#/components/schemas/", "", 1))
	}

	return "unknown"
}
