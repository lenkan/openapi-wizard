package openapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOperations(t *testing.T) {
	spec := ApiDefinition{
		Paths: map[string]PathDefinition{
			"/p1": {Get: OperationDefinition{
				Description: "Hello",
				OperationId: "get_p1",
			}},
			"/p2": {Get: OperationDefinition{
				Description: "Hello",
				OperationId: "get_p2",
			}},
		},
	}

	operations := spec.ListOperations()

	assert.Len(t, operations, 2)
}

func TestListSchemas(t *testing.T) {
	spec := ApiDefinition{
		Paths: map[string]PathDefinition{
			"/p1": {Get: OperationDefinition{
				Description: "Hello",
				OperationId: "get_p1",
			}},
		},
		Components: ComponentsDefinition{
			Schemas: map[string]JsonSchemaDefinition{
				"User": {
					Type: "object",
					Properties: map[string]JsonSchemaDefinition{
						"name": {Type: "string"},
					},
				},
			},
		},
	}

	schemas := spec.ListSchemas()

	assert.Len(t, schemas, 1)
}
