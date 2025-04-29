package parser

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Parser is a structure that holds loaded package information
type Parser struct {
	pkgs map[string]*packages.Package
}

// New creates a Parser instance by loading Go packages from the specified directory.
// rootDir is the base directory, and pkgDir specifies the pattern to load.
// If pkgDir is empty, "./..." is used as the pattern.
func New(rootDir string, pkgDir string) (*Parser, error) {
	patterns := []string{"./..."}
	if pkgDir != "" {
		patterns = []string{pkgDir}
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir:   rootDir,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	parser := &Parser{
		pkgs: make(map[string]*packages.Package),
	}

	// Store packages in the map
	for _, pkg := range pkgs {
		parser.pkgs[pkg.PkgPath] = pkg
	}

	return parser, nil
}

// GetAllPackages returns all loaded packages
func (p *Parser) GetAllPackages() []*packages.Package {
	result := make([]*packages.Package, 0, len(p.pkgs))
	for _, pkg := range p.pkgs {
		result = append(result, pkg)
	}
	return result
}

// GetPackage returns a package by its package path.
// Returns an error if the package is not found.
func (p *Parser) GetPackage(pkgPath string) (*packages.Package, error) {
	pkg, ok := p.pkgs[pkgPath]
	if !ok {
		return nil, fmt.Errorf("package not found: %s", pkgPath)
	}
	return pkg, nil
}

// StructInfo represents information about a struct
type StructInfo struct {
	Name    string   // Struct name
	Comment string   // Struct comment
	Fields  []Field  // List of fields
	Methods []Method // List of methods
}

// Field represents information about a struct field
type Field struct {
	Name       string // Field name
	Type       string // Field type
	Comment    string // Field comment
	IsExported bool   // Whether the field is exported
}

// Method represents information about a struct method
type Method struct {
	Name      string    // Method name
	Signature string    // Method signature
	Comment   string    // Method comment
	Examples  []Example // Method examples
}

// GetStructInfo returns information about a struct in the specified package
func (p *Parser) GetStructInfo(pkgPath, structName string) (*StructInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// Get type information from the package
	scope := pkg.Types.Scope()
	obj := scope.Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct not found: %s in package %s", structName, pkgPath)
	}

	// Check if the type is a struct
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	// Build struct information
	info := &StructInfo{
		Name:    structName,
		Comment: GetComment(pkg, obj),
		Fields:  make([]Field, 0, structType.NumFields()),
		Methods: make([]Method, 0),
	}

	// Get field information
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		info.Fields = append(info.Fields, Field{
			Name:       field.Name(),
			Type:       field.Type().String(),
			Comment:    GetComment(pkg, field),
			IsExported: field.Exported(),
		})
	}

	// Get method information
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		info.Methods = append(info.Methods, Method{
			Name:      method.Name(),
			Signature: method.Type().String(),
			Comment:   GetComment(pkg, method),
		})
	}

	return info, nil
}

// FuncInfo represents information about a function
type FuncInfo struct {
	Name      string    // Function name
	Signature string    // Function signature
	Comment   string    // Function comment
	Examples  []Example // Function examples
}

// Example represents an example for a function
type Example struct {
	Name   string // Example name
	Code   string // Example code
	Output string // Example output
}

// GetFuncInfo returns information about a function in the specified package
func (p *Parser) GetFuncInfo(pkgPath, funcName string) (*FuncInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// Get type information from the package
	scope := pkg.Types.Scope()
	obj := scope.Lookup(funcName)
	if obj == nil {
		return nil, fmt.Errorf("function not found: %s in package %s", funcName, pkgPath)
	}

	// Check if the type is a function
	fn, ok := obj.(*types.Func)
	if !ok {
		return nil, fmt.Errorf("not a function: %s in package %s", funcName, pkgPath)
	}

	// Build function information
	info := &FuncInfo{
		Name:      funcName,
		Signature: fn.Type().String(),
		Comment:   GetComment(pkg, obj),
		Examples:  make([]Example, 0),
	}

	// Get examples
	info.Examples = getExamples(pkg, funcName)

	return info, nil
}

// GetMethodInfo returns information about a method in the specified package
func (p *Parser) GetMethodInfo(pkgPath, structName, methodName string) (*Method, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, err
	}

	// Get type information from the package
	scope := pkg.Types.Scope()
	obj := scope.Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct not found: %s in package %s", structName, pkgPath)
	}

	// Check if the type is a struct
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("not a struct: %s in package %s", structName, pkgPath)
	}

	// Find the method
	var method *types.Func
	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		if m.Name() == methodName {
			method = m
			break
		}
	}

	if method == nil {
		return nil, fmt.Errorf("method not found: %s.%s in package %s", structName, methodName, pkgPath)
	}

	// Build method information
	info := &Method{
		Name:      methodName,
		Signature: method.Type().String(),
		Comment:   GetComment(pkg, method),
		Examples:  make([]Example, 0),
	}

	// Get examples
	info.Examples = getExamples(pkg, methodName)

	return info, nil
}

// ConstInfo represents information about a constant
type ConstInfo struct {
	Name    string // Constant name
	Type    string // Constant type
	Value   string // Constant value
	Comment string // Constant comment
}

// VarInfo represents information about a variable
type VarInfo struct {
	Name    string // Variable name
	Type    string // Variable type
	Comment string // Variable comment
}

// GetConstAndVarInfo returns information about constants and variables in the specified package
func (p *Parser) GetConstAndVarInfo(pkgPath string) ([]ConstInfo, []VarInfo, error) {
	pkg, err := p.GetPackage(pkgPath)
	if err != nil {
		return nil, nil, err
	}

	// Get type information from the package
	scope := pkg.Types.Scope()

	// Get constant information
	constants := make([]ConstInfo, 0)
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		// Check if it's a constant
		if constObj, ok := obj.(*types.Const); ok {
			constants = append(constants, ConstInfo{
				Name:    name,
				Type:    constObj.Type().String(),
				Value:   constObj.Val().String(),
				Comment: GetComment(pkg, obj),
			})
		}
	}

	// Get variable information
	variables := make([]VarInfo, 0)
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		// Check if it's a variable
		if varObj, ok := obj.(*types.Var); ok {
			variables = append(variables, VarInfo{
				Name:    name,
				Type:    varObj.Type().String(),
				Comment: GetComment(pkg, obj),
			})
		}
	}

	return constants, variables, nil
}

// getExamples returns examples for a function
func getExamples(pkg *packages.Package, funcName string) []Example {
	var examples []Example

	// Search AST nodes in the package
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			// Look for function declarations
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Check if it's an example function
			if !strings.HasPrefix(funcDecl.Name.Name, "Example") {
				return true
			}

			// Check if it's an example for the target function
			if !strings.HasSuffix(funcDecl.Name.Name, funcName) {
				return true
			}

			// Get example information
			example := Example{
				Name: funcDecl.Name.Name,
			}

			// Get example code
			if funcDecl.Body != nil {
				example.Code = getNodeString(pkg.Fset, funcDecl.Body)
			}

			// Get example output
			if funcDecl.Doc != nil {
				for _, comment := range funcDecl.Doc.List {
					if strings.HasPrefix(comment.Text, "// Output:") {
						example.Output = strings.TrimPrefix(comment.Text, "// Output:")
						example.Output = strings.TrimSpace(example.Output)
					}
				}
			}

			examples = append(examples, example)
			return true
		})
	}

	return examples
}

// getNodeString returns the string representation of a node
func getNodeString(fset *token.FileSet, node ast.Node) string {
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return ""
	}
	return buf.String()
}

// GetComment returns the comment for an object
func GetComment(pkg *packages.Package, obj types.Object) string {
	var comment string

	// Search AST nodes in the package
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeSpec:
				// For type declarations
				if node.Name.Name == obj.Name() {
					if node.Doc != nil {
						comment = node.Doc.Text()
					}
					return false
				}
			case *ast.Field:
				// For fields
				for _, name := range node.Names {
					if name.Name == obj.Name() {
						if node.Doc != nil {
							comment = node.Doc.Text()
						} else if node.Comment != nil {
							comment = node.Comment.Text()
						}
						return false
					}
				}
			case *ast.FuncDecl:
				// For methods
				if node.Name.Name == obj.Name() {
					if node.Doc != nil {
						comment = node.Doc.Text()
					}
					return false
				}
			case *ast.ValueSpec:
				// For constants and variables
				for _, name := range node.Names {
					if name.Name == obj.Name() {
						if node.Doc != nil {
							comment = node.Doc.Text()
						} else if node.Comment != nil {
							comment = node.Comment.Text()
						}
						return false
					}
				}
			}
			return true
		})
	}

	// Trim whitespace from the comment
	return strings.TrimSpace(comment)
}

// GetPackageComment returns the package comment.
// Package comments are typically comment blocks before the package declaration.
func GetPackageComment(pkg *packages.Package) string {
	if pkg == nil || len(pkg.Syntax) == 0 {
		return ""
	}

	var comment string
	// Search for package comments in each file
	for _, file := range pkg.Syntax {
		if file.Doc != nil && file.Doc.Text() != "" {
			// If multiple files have package comments, use the first non-empty comment
			comment = file.Doc.Text()
			break
		}
	}

	return strings.TrimSpace(comment)
}
