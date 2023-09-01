package openapi

import (
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
)

type JsonSchemaDefinition struct {
	Type                 string                          `yaml:"type"`
	Properties           map[string]JsonSchemaDefinition `yaml:"properties"`
	Ref                  string                          `yaml:"$ref"`
	Required             []string                        `yaml:"required"`
	AdditionalProperties bool                            `default:"false" yaml:"additionalProperties"`
	OneOf                []JsonSchemaDefinition          `yaml:"oneOf"`
	AllOf                []JsonSchemaDefinition          `yaml:"allOf"`
	Enum                 []string                        `yaml:"enum"`
	Items                *JsonSchemaDefinition           `yaml:"items"`
}

func MapShapes(schemas []JsonSchemaDefinition) []string {
	types := []string{}

	for _, s := range schemas {
		result := FormatSchemaShape(s)

		if !slices.Contains(types, result) {
			types = append(types, result)
		}
	}

	return types
}

func FormatSchemaShape(schema JsonSchemaDefinition) string {
	if len(schema.AllOf) > 0 {
		types := MapShapes(schema.AllOf)
		return "(" + strings.Join(types, " & ") + ")"
	}

	if len(schema.OneOf) > 0 {
		types := MapShapes(schema.OneOf)
		return "(" + strings.Join(types, " | ") + ")"
	}

	if schema.Items != nil {
		return "(" + FormatSchemaShape(*schema.Items) + ")[]"
	}

	if len(schema.Enum) > 0 {
		enum := []string{}
		for _, value := range schema.Enum {
			enum = append(enum, "\""+value+"\"")
		}
		return "(" + strings.Join(enum, " | ") + ")"
	}

	if schema.Type == "boolean" {
		return "boolean"
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
