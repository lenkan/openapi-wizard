package openapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectWithRequiredProperty(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type: "object",
		Properties: map[string]JsonSchemaDefinition{
			"name": {Type: "string"},
		},
		Required: []string{"name"},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "{name: string}", result)
}

func TestObjectWithOptionalProperty(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type: "object",
		Properties: map[string]JsonSchemaDefinition{
			"name": {Type: "string"},
		},
		Required: []string{},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "{name?: string}", result)
}

func TestObjectWithNoProperties(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type:       "object",
		Properties: map[string]JsonSchemaDefinition{},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "Record<string, never>", result)
}

func TestObjectWithAdditionalProperties(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type:                 "object",
		Properties:           map[string]JsonSchemaDefinition{},
		AdditionalProperties: true,
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "Record<string, unknown>", result)
}

func TestUnion(t *testing.T) {
	schema := JsonSchemaDefinition{
		OneOf: []JsonSchemaDefinition{
			{Type: "string"},
			{Type: "number"},
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(string | number)", result)
}

func TestUnionWithObjects(t *testing.T) {
	schema := JsonSchemaDefinition{
		OneOf: []JsonSchemaDefinition{
			{
				Type: "object", Properties: map[string]JsonSchemaDefinition{
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
	schema := JsonSchemaDefinition{
		AllOf: []JsonSchemaDefinition{
			{Type: "string"},
			{Type: "number"},
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(string & number)", result)
}

func TestArraySchema(t *testing.T) {
	schema := JsonSchemaDefinition{
		Items: &JsonSchemaDefinition{
			Type: "string",
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(string)[]", result)
}

func TestEnumSchema(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type: "string",
		Enum: []string{
			"abc",
			"def",
		},
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "(\"abc\" | \"def\")", result)
}

func TestBooleanSchema(t *testing.T) {
	schema := JsonSchemaDefinition{
		Type: "boolean",
	}

	result := FormatSchemaShape(schema)

	assert.Equal(t, "boolean", result)
}
