package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatPackageList formats a list of packages into a JSON string
func FormatPackageList(packages []PackageInfo) string {
	response := ListPackagesResponse{
		Packages: packages,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format package list: %v"}`, err)
	}
	return string(jsonBytes)
}

// FormatPackageInspection formats package inspection results into a JSON string
func FormatPackageInspection(pkg PackageInfo, structs []StructSummary, funcs []FuncSummary, methods []MethodSummary, includeComments bool) string {
	response := InspectPackageResponse{
		Package:   pkg,
		Structs:   structs,
		Functions: funcs,
		Methods:   methods,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format package inspection: %v"}`, err)
	}
	return string(jsonBytes)
}

// FormatStructDoc formats struct documentation into a JSON string
func FormatStructDoc(name, comment string, fields []FieldDoc, methods []MethodDoc) string {
	response := StructDocResponse{
		Name:    name,
		Comment: comment,
		Fields:  fields,
		Methods: methods,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format struct documentation: %v"}`, err)
	}
	return string(jsonBytes)
}

// FormatFuncDoc formats function documentation into a JSON string
func FormatFuncDoc(name, signature, comment string, examples []Example) string {
	response := FuncDocResponse{
		Name:      name,
		Signature: signature,
		Comment:   comment,
		Examples:  examples,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format function documentation: %v"}`, err)
	}
	return string(jsonBytes)
}

// FormatMethodDoc formats method documentation into a JSON string
func FormatMethodDoc(receiverType, name, signature, comment string, examples []Example) string {
	response := MethodDocResponse{
		ReceiverType: receiverType,
		Name:         name,
		Signature:    signature,
		Comment:      comment,
		Examples:     examples,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format method documentation: %v"}`, err)
	}
	return string(jsonBytes)
}

// FormatConstAndVarDoc formats constant and variable documentation into a JSON string
func FormatConstAndVarDoc(constants []ConstDoc, variables []VarDoc) string {
	response := ConstAndVarResponse{
		Constants: constants,
		Variables: variables,
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to format constant and variable documentation: %v"}`, err)
	}
	return string(jsonBytes)
}

// formatPackageList formats a list of packages into a markdown string
func FormatPackageListMarkdown(packages []PackageInfo) string {
	var sb strings.Builder
	sb.WriteString("# Packages\n\n")
	for _, pkg := range packages {
		sb.WriteString(fmt.Sprintf("## %s\n", pkg.Name))
		sb.WriteString(fmt.Sprintf("Import Path: `%s`\n\n", pkg.ImportPath))
		if pkg.Comment != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", pkg.Comment))
		}
	}
	return sb.String()
}

// formatPackageInspection formats package inspection results into a markdown string
func FormatPackageInspectionMarkdown(pkg PackageInfo, structs []StructSummary, funcs []FuncSummary, methods []MethodSummary, includeComments bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Package: %s\n\n", pkg.Name))
	sb.WriteString(fmt.Sprintf("Import Path: `%s`\n\n", pkg.ImportPath))
	if pkg.Comment != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", pkg.Comment))
	}

	if len(structs) > 0 {
		sb.WriteString("## Structs\n\n")
		for _, s := range structs {
			sb.WriteString(fmt.Sprintf("### %s\n", s.Name))
			if includeComments && s.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", s.Comment))
			}
		}
	}

	if len(funcs) > 0 {
		sb.WriteString("## Functions\n\n")
		for _, f := range funcs {
			sb.WriteString(fmt.Sprintf("### %s\n", f.Name))
			if includeComments && f.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", f.Comment))
			}
		}
	}

	if len(methods) > 0 {
		sb.WriteString("## Methods\n\n")
		for _, m := range methods {
			sb.WriteString(fmt.Sprintf("### %s.%s\n", m.ReceiverType, m.Name))
			if includeComments && m.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", m.Comment))
			}
		}
	}

	return sb.String()
}

// formatStructDoc formats struct documentation into a markdown string
func FormatStructDocMarkdown(name, comment string, fields []FieldDoc, methods []MethodDoc) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Struct: %s\n\n", name))
	if comment != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", comment))
	}

	if len(fields) > 0 {
		sb.WriteString("## Fields\n\n")
		for _, f := range fields {
			sb.WriteString(fmt.Sprintf("### %s\n", f.Name))
			sb.WriteString(fmt.Sprintf("Type: `%s`\n", f.Type))
			if f.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", f.Comment))
			}
		}
	}

	if len(methods) > 0 {
		sb.WriteString("## Methods\n\n")
		for _, m := range methods {
			sb.WriteString(fmt.Sprintf("### %s\n", m.Name))
			sb.WriteString(fmt.Sprintf("Signature: `%s`\n", m.Signature))
			if m.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", m.Comment))
			}
		}
	}

	return sb.String()
}

// formatFuncDoc formats function documentation into a markdown string
func FormatFuncDocMarkdown(name, signature, comment string, examples []Example) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Function: %s\n\n", name))
	sb.WriteString(fmt.Sprintf("Signature: `%s`\n\n", signature))
	if comment != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", comment))
	}

	if len(examples) > 0 {
		sb.WriteString("## Examples\n\n")
		for _, e := range examples {
			sb.WriteString(fmt.Sprintf("### %s\n", e.Name))
			sb.WriteString("```go\n")
			sb.WriteString(e.Code)
			sb.WriteString("\n```\n")
			if e.Output != "" {
				sb.WriteString("Output:\n```\n")
				sb.WriteString(e.Output)
				sb.WriteString("\n```\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// formatMethodDoc formats method documentation into a markdown string
func FormatMethodDocMarkdown(receiverType, name, signature, comment string, examples []Example) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Method: %s.%s\n\n", receiverType, name))
	sb.WriteString(fmt.Sprintf("Signature: `%s`\n\n", signature))
	if comment != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", comment))
	}

	if len(examples) > 0 {
		sb.WriteString("## Examples\n\n")
		for _, e := range examples {
			sb.WriteString(fmt.Sprintf("### %s\n", e.Name))
			sb.WriteString("```go\n")
			sb.WriteString(e.Code)
			sb.WriteString("\n```\n")
			if e.Output != "" {
				sb.WriteString("Output:\n```\n")
				sb.WriteString(e.Output)
				sb.WriteString("\n```\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// formatConstAndVarDoc formats constant and variable documentation into a markdown string
func FormatConstAndVarDocMarkdown(constants []ConstDoc, variables []VarDoc) string {
	var sb strings.Builder

	if len(constants) > 0 {
		sb.WriteString("# Constants\n\n")
		for _, c := range constants {
			sb.WriteString(fmt.Sprintf("## %s\n", c.Name))
			sb.WriteString(fmt.Sprintf("Type: `%s`\n", c.Type))
			sb.WriteString(fmt.Sprintf("Value: `%s`\n", c.Value))
			if c.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", c.Comment))
			}
		}
	}

	if len(variables) > 0 {
		sb.WriteString("# Variables\n\n")
		for _, v := range variables {
			sb.WriteString(fmt.Sprintf("## %s\n", v.Name))
			sb.WriteString(fmt.Sprintf("Type: `%s`\n", v.Type))
			if v.Comment != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", v.Comment))
			}
		}
	}

	return sb.String()
}
