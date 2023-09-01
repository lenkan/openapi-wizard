package codegen

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lenkan/openapi-wizard/internal/openapi"
	"golang.org/x/exp/slices"
)

func formatParameterName(name string) string {
	return strcase.ToLowerCamel(name)
}

func formatSchemaShape(schema openapi.JsonSchemaDefinition) string {
	if schema.Type == "string" {
		return "string"
	}

	if schema.Type == "object" {
		props := []string{}

		for propertyName, propertyDefinition := range schema.Properties {
			suffix := ""
			if !slices.Contains(propertyDefinition.Required, propertyName) {
				suffix += "?"
			}

			props = append(props, propertyName+suffix+": "+formatSchemaShape(propertyDefinition))
		}

		return strings.Join([]string{"{", strings.Join(props, ";"), "}"}, "")
	}

	if schema.Ref != "" {
		return strcase.ToCamel(strings.Replace(schema.Ref, "#/components/schemas/", "", 1))
	}

	return "unknown"
}

func formatParameterInterface(name string, params []openapi.ParameterDefinition) string {
	result := []string{}

	for _, param := range params {
		name := formatParameterName(param.Name)

		if !param.Required {
			name += "?"
		}

		result = append(result, name+": "+formatSchemaShape(param.Schema))
	}

	return strings.Join([]string{
		"export interface " + strcase.ToCamel(name+"_params") + " {",
		"  " + strings.Join(result, "\n  "),
		"}",
	}, "\n")
}

func formatOperation(path string, method string, operation *openapi.OperationDefinition) string {
	name := strcase.ToLowerCamel(operation.OperationId)
	interfaceName := strcase.ToCamel(name + "_params")

	code := []string{}
	for _, param := range operation.Parameters {
		paramName := formatParameterName(param.Name)
		if param.In == "query" && param.Schema.Type == "string" {
			val := "params." + paramName
			code = append(code, []string{
				"if (" + val + " !== undefined) {",
				"  url.searchParams.set(\"" + param.Name + "\", " + val + ")",
				"}",
			}...)
		}

		if param.In == "header" && param.Schema.Type == "string" {
			val := "params." + paramName
			code = append(code, []string{
				"if (" + val + " !== undefined) {",
				"  headers.set(\"" + param.Name + "\", " + val + ")",
				"}",
			}...)
		}
	}

	responseShape := formatSchemaShape(operation.Responses["200"].Content["application/json"].Schema)

	return strings.Join([]string{
		"  async " + name + "(params: " + interfaceName + "): Promise<" + responseShape + "> {", //
		"    const headers = new Headers();",
		"    const url = new URL(\"" + path + "\", this.baseUrl)",
		"    " + strings.Join(code, "\n  "),
		"    const response = await fetch(url, { headers });",
		"    const body = await response.json()",
		"    return body;",
		"  }",
	}, "\n")
}

func formatMethods(operations []openapi.ApiOperation, indent string) string {
	result := []string{}

	for _, operation := range operations {
		result = append(result, formatOperation(operation.Path, "get", &operation.Definition))
	}

	return strings.Join(result, "\n"+indent)
}

func formatInterfaces(operations []openapi.ApiOperation, indent string) string {
	result := []string{}

	for _, operation := range operations {
		result = append(result, formatParameterInterface(operation.Definition.OperationId, operation.Definition.Parameters))
	}

	return strings.Join(result, "\n"+indent)
}

func formatSchemas(schemas []openapi.ApiSchema, indent string) string {
	result := []string{}

	for _, schema := range schemas {
		result = append(result, "export type "+strcase.ToCamel(schema.Name)+" = "+formatSchemaShape(schema.Schema))
	}

	return strings.Join(result, "\n"+indent)
}

func FormatTypescriptClient(spec *openapi.ApiDefinition) string {
	operations := spec.ListOperations()
	schemas := spec.ListSchemas()

	result := []string{
		formatSchemas(schemas, ""),
		formatInterfaces(operations, ""),
		"export class Api {",
		"  constructor(private baseUrl: string = window.location.origin) {",
		"  }",
		"  " + formatMethods(operations, "  "),
		"}",
	}

	return strings.Join(result, "\n\n")
}
