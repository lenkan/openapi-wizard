package openapi

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ApiInfo struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

type OperationDefinition struct {
	Description string                        `yaml:"description"`
	OperationId string                        `yaml:"operationId"`
	Summary     string                        `yaml:"summary"`
	Parameters  []ParameterDefinition         `yaml:"parameters"`
	Responses   map[string]ResponseDefinition `yaml:"responses"`
}

type ResponseDefinition struct {
	Description string                       `yaml:"description"`
	Content     map[string]ContentDefinition `yaml:"content"`
}

type ContentDefinition struct {
	Schema JsonSchemaDefinition `yaml:"schema"`
}

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

type ParameterDefinition struct {
	In       string               `yaml:"in"`
	Name     string               `yaml:"name"`
	Required bool                 `yaml:"required"`
	Schema   JsonSchemaDefinition `yaml:"schema"`
}

type PathDefinition struct {
	Get    OperationDefinition `yaml:"get,omitempty"`
	Post   OperationDefinition `yaml:"post,omitempty"`
	Put    OperationDefinition `yaml:"put,omitempty"`
	Delete OperationDefinition `yaml:"delete,omitempty"`
	Patch  OperationDefinition `yaml:"patch,omitempty"`
}

type ApiDefinition struct {
	Openapi    string                    `yaml:"openapi"`
	Info       ApiInfo                   `yaml:"info"`
	Paths      map[string]PathDefinition `yaml:"paths"`
	Components ComponentsDefinition      `yaml:"components"`
}

type ComponentsDefinition struct {
	Responses map[string]ResponseDefinition   `yaml:"responses"`
	Schemas   map[string]JsonSchemaDefinition `yaml:"schemas"`
}

type ApiOperation struct {
	Path       string
	Method     string
	Definition OperationDefinition
}

type ApiSchema struct {
	Name   string
	Schema JsonSchemaDefinition
}

func (schema *ApiDefinition) ListOperations() []ApiOperation {
	result := []ApiOperation{}

	for path, pathDefinition := range schema.Paths {
		if pathDefinition.Get.OperationId != "" {
			result = append(result, ApiOperation{Path: path, Method: "get", Definition: pathDefinition.Get})
		}

		if pathDefinition.Post.OperationId != "" {
			result = append(result, ApiOperation{Path: path, Method: "post", Definition: pathDefinition.Post})
		}
	}

	return result
}

func (schema *ApiDefinition) ListSchemas() []ApiSchema {
	result := []ApiSchema{}

	for name, schema := range schema.Components.Schemas {
		result = append(result, ApiSchema{
			Name:   name,
			Schema: schema,
		})
	}

	return result

}

func (schema *ApiDefinition) Print() string {
	result, error := yaml.Marshal(schema)

	if error != nil {
		panic(error)
	}

	return string(result)
}

func Load(filename string) *ApiDefinition {
	if filename == "" {
		panic("No filename provided")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	spec := &ApiDefinition{}
	{
		yaml.Unmarshal(data, spec)
		if err != nil {
			panic(err)
		}
	}

	return spec
}
