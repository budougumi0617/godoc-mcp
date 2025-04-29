package handler

import (
	"context"
	"fmt"
	"go/types"
	"strings"

	godoc "github.com/budougumi0617/godoc-mcp"
	"github.com/budougumi0617/godoc-mcp/internal/model"
	"github.com/budougumi0617/godoc-mcp/internal/parser"
	mcp "github.com/ktr0731/go-mcp"
)

// ToolHandler is a handler structure that processes MCP tool requests.
type ToolHandler struct {
	parser *parser.Parser
}

// NewToolHandler creates a new ToolHandler instance.
func NewToolHandler(p *parser.Parser) *ToolHandler {
	return &ToolHandler{
		parser: p,
	}
}

// HandleToolGolangListPackages returns a list of all loaded packages.
func (h *ToolHandler) HandleToolGolangListPackages(ctx context.Context, req *godoc.ToolGolangListPackagesRequest) (*mcp.CallToolResult, error) {
	pkgs := h.parser.GetAllPackages()
	if len(pkgs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.CallToolContent{
				mcp.TextContent{Text: "No packages loaded."},
			},
		}, nil
	}

	var packages []model.PackageInfo
	for _, p := range pkgs {
		// Get package comment
		packages = append(packages, model.PackageInfo{
			Name:       p.Name,
			ImportPath: p.PkgPath,
			Comment:    parser.GetPackageComment(p),
		})
	}

	// Format in markdown
	mdContent := model.FormatPackageListMarkdown(packages)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGolangInspectPackage lists exported structs, methods, and functions in the specified package.
func (h *ToolHandler) HandleToolGolangInspectPackage(ctx context.Context, req *godoc.ToolGolangInspectPackageRequest) (*mcp.CallToolResult, error) {
	pkg, err := h.parser.GetPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	// Create package info
	pkgInfo := model.PackageInfo{
		Name:       pkg.Name,
		ImportPath: pkg.PkgPath,
		Comment:    parser.GetPackageComment(pkg),
	}

	// Collect struct, function, and method information
	scope := pkg.Types.Scope()
	var structs []model.StructSummary
	var funcs []model.FuncSummary
	var methods []model.MethodSummary

	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}

		switch obj := obj.(type) {
		case *types.TypeName:
			// Get type definition
			if _, ok := obj.Type().Underlying().(*types.Struct); ok {
				structs = append(structs, model.StructSummary{
					Name:    obj.Name(),
					Comment: parser.GetComment(pkg, obj), // Use public function from parser package
				})
			}
		case *types.Func:
			sig, ok := obj.Type().(*types.Signature)
			if !ok {
				continue
			}

			// For methods (with receiver)
			if sig.Recv() != nil {
				recvType := sig.Recv().Type().String()
				// Remove * for pointer types
				recvType = strings.TrimPrefix(recvType, "*")
				methods = append(methods, model.MethodSummary{
					ReceiverType: recvType,
					Name:         obj.Name(),
					Comment:      parser.GetComment(pkg, obj), // Use public function from parser package
				})
			} else {
				// For functions
				funcs = append(funcs, model.FuncSummary{
					Name:    obj.Name(),
					Comment: parser.GetComment(pkg, obj), // Use public function from parser package
				})
			}
		}
	}

	// Format in markdown
	mdContent := model.FormatPackageInspectionMarkdown(pkgInfo, structs, funcs, methods, req.IncludeComments)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGolangGetStructDoc returns information about the specified struct.
func (h *ToolHandler) HandleToolGolangGetStructDoc(ctx context.Context, req *godoc.ToolGolangGetStructDocRequest) (*mcp.CallToolResult, error) {
	structInfo, err := h.parser.GetStructInfo(req.PackageName, req.StructName)
	if err != nil {
		return nil, fmt.Errorf("failed to get struct info: %w", err)
	}

	// Convert field and method information
	var fields []model.FieldDoc
	var methods []model.MethodDoc

	for _, f := range structInfo.Fields {
		fields = append(fields, model.FieldDoc{
			Name:       f.Name,
			Type:       f.Type,
			Comment:    f.Comment,
			IsExported: f.IsExported,
		})
	}

	for _, m := range structInfo.Methods {
		methods = append(methods, model.MethodDoc{
			Name:      m.Name,
			Signature: m.Signature,
			Comment:   m.Comment,
		})
	}

	// Format in markdown
	mdContent := model.FormatStructDocMarkdown(structInfo.Name, structInfo.Comment, fields, methods)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGolangGetFuncDoc returns information about the specified function.
func (h *ToolHandler) HandleToolGolangGetFuncDoc(ctx context.Context, req *godoc.ToolGolangGetFuncDocRequest) (*mcp.CallToolResult, error) {
	funcInfo, err := h.parser.GetFuncInfo(req.PackageName, req.FuncName)
	if err != nil {
		return nil, fmt.Errorf("failed to get function info: %w", err)
	}

	// Convert examples
	var examples []model.Example
	for _, e := range funcInfo.Examples {
		examples = append(examples, model.Example{
			Name:   e.Name,
			Code:   e.Code,
			Output: e.Output,
		})
	}

	// Format in markdown
	mdContent := model.FormatFuncDocMarkdown(funcInfo.Name, funcInfo.Signature, funcInfo.Comment, examples)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGolangGetMethodDoc returns information about the specified method of a struct.
func (h *ToolHandler) HandleToolGolangGetMethodDoc(ctx context.Context, req *godoc.ToolGolangGetMethodDocRequest) (*mcp.CallToolResult, error) {
	methodInfo, err := h.parser.GetMethodInfo(req.PackageName, req.StructName, req.MethodName)
	if err != nil {
		return nil, fmt.Errorf("failed to get method info: %w", err)
	}

	// Convert examples
	var examples []model.Example
	for _, e := range methodInfo.Examples {
		examples = append(examples, model.Example{
			Name:   e.Name,
			Code:   e.Code,
			Output: e.Output,
		})
	}

	// Format in markdown
	mdContent := model.FormatMethodDocMarkdown(req.StructName, methodInfo.Name, methodInfo.Signature, methodInfo.Comment, examples)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}

// HandleToolGolangGetConstAndVarDoc returns information about constants and variables in the specified package.
func (h *ToolHandler) HandleToolGolangGetConstAndVarDoc(ctx context.Context, req *godoc.ToolGolangGetConstAndVarDocRequest) (*mcp.CallToolResult, error) {
	constInfos, varInfos, err := h.parser.GetConstAndVarInfo(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get constant and variable info: %w", err)
	}

	// Convert constant and variable information
	var constants []model.ConstDoc
	var variables []model.VarDoc

	for _, c := range constInfos {
		constants = append(constants, model.ConstDoc{
			Name:    c.Name,
			Type:    c.Type,
			Value:   c.Value,
			Comment: c.Comment,
		})
	}

	for _, v := range varInfos {
		variables = append(variables, model.VarDoc{
			Name:    v.Name,
			Type:    v.Type,
			Comment: v.Comment,
		})
	}

	// Format in markdown
	mdContent := model.FormatConstAndVarDocMarkdown(constants, variables)

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: mdContent},
		},
	}, nil
}
