package codegen

import (
	"testing"

	"github.com/lenkan/openapi-wizard/internal/openapi"
	"github.com/stretchr/testify/assert"
)

func TestObjectWithRequiredProperty(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		Type: "object",
		Properties: map[string]openapi.JsonSchemaDefinition{
			"name": {Type: "string"},
		},
		Required: []string{"name"},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "{name: string}", result)
}

func TestObjectWithOptionalProperty(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		Type: "object",
		Properties: map[string]openapi.JsonSchemaDefinition{
			"name": {Type: "string"},
		},
		Required: []string{},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "{name?: string}", result)
}

func TestObjectWithNoProperties(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		Type:       "object",
		Properties: map[string]openapi.JsonSchemaDefinition{},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "Record<string, never>", result)
}

func TestObjectWithAdditionalProperties(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		Type:                 "object",
		Properties:           map[string]openapi.JsonSchemaDefinition{},
		AdditionalProperties: true,
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "Record<string, unknown>", result)
}

func TestUnion(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		OneOf: []openapi.JsonSchemaDefinition{
			{Type: "string"},
			{Type: "number"},
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(string | number)", result)
}

func TestUnionWithObjects(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		OneOf: []openapi.JsonSchemaDefinition{
			{
				Type: "object", Properties: map[string]openapi.JsonSchemaDefinition{
					"name": {Type: "string"},
				},
				Required: []string{"name"},
			},
			{Type: "number"},
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "({name: string} | number)", result)
}

func TestIntersection(t *testing.T) {
	schema := openapi.JsonSchemaDefinition{
		AllOf: []openapi.JsonSchemaDefinition{
			{Type: "string"},
			{Type: "number"},
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(string & number)", result)
}
