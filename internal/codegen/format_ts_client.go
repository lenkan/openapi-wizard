package codegen

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/lenkan/openapi-wizard/internal/openapi"
)

func formatParameterName(name string) string {
	return strcase.ToLowerCamel(name)
}

func formatParameterInterface(name string, params []openapi.ParameterDefinition) string {
	result := []string{}

	for _, param := range params {
		name := formatParameterName(param.Name)

		if !param.Required {
			name += "?"
		}

		result = append(result, name+": "+FormatSchemaShape(param.Schema))
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

		if param.In == "path" && param.Schema.Type == "string" {
			val := "params." + paramName
			code = append(code, "url.pathname = url.pathname.replace(\"{"+param.Name+"}\", "+val+");")
		}
	}

	responseShape := FormatSchemaShape(operation.Responses["200"].Content["application/json"].Schema)

	paramsSuffix := ""
	if len(operation.Parameters) == 0 {
		paramsSuffix = "?"
	}

	return strings.Join([]string{
		"  async " + name + "(params" + paramsSuffix + ": " + interfaceName + "): Promise<" + responseShape + "> {", //
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
		result = append(result, "export type "+strcase.ToCamel(schema.Name)+" = "+FormatSchemaShape(schema.Schema))
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
