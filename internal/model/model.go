package model

// このファイルには、MCPツールのリクエストパラメータを表す構造体や、
// レスポンス整形用の補助的な構造体を定義します。
// Goのパッケージ、型、関数などの情報は `golang.org/x/tools/go/packages` の型を直接利用します。

// TODO: 各MCPツールのリクエストパラメータに対応する構造体を定義する
// 例:
// type InspectPackageParams struct {
//     PackageName     string `json:"packageName"`
//     IncludeComments bool   `json:"includeComments"`
// }

// PackageInfo represents information about a Go package
type PackageInfo struct {
	Name       string `json:"name"`        // Package name
	ImportPath string `json:"import_path"` // Import path
	Comment    string `json:"comment"`     // Package comment
}

// StructSummary represents a summary of a struct
type StructSummary struct {
	Name    string `json:"name"`    // Struct name
	Comment string `json:"comment"` // Struct comment
}

// FuncSummary represents a summary of a function
type FuncSummary struct {
	Name    string `json:"name"`    // Function name
	Comment string `json:"comment"` // Function comment
}

// MethodSummary represents a summary of a method
type MethodSummary struct {
	ReceiverType string `json:"receiver_type"` // Receiver type
	Name         string `json:"name"`          // Method name
	Comment      string `json:"comment"`       // Method comment
}

// FieldDoc represents documentation for a struct field
type FieldDoc struct {
	Name       string `json:"name"`        // Field name
	Type       string `json:"type"`        // Field type
	Comment    string `json:"comment"`     // Field comment
	IsExported bool   `json:"is_exported"` // Whether the field is exported
}

// MethodDoc represents documentation for a method
type MethodDoc struct {
	Name      string `json:"name"`      // Method name
	Signature string `json:"signature"` // Method signature
	Comment   string `json:"comment"`   // Method comment
}

// Example represents an example for a function or method
type Example struct {
	Name   string `json:"name"`   // Example name
	Code   string `json:"code"`   // Example code
	Output string `json:"output"` // Example output
}

// ConstDoc represents documentation for a constant
type ConstDoc struct {
	Name    string `json:"name"`    // Constant name
	Type    string `json:"type"`    // Constant type
	Value   string `json:"value"`   // Constant value
	Comment string `json:"comment"` // Constant comment
}

// VarDoc represents documentation for a variable
type VarDoc struct {
	Name    string `json:"name"`    // Variable name
	Type    string `json:"type"`    // Variable type
	Comment string `json:"comment"` // Variable comment
}

// ListPackagesResponse represents the response for list_packages
type ListPackagesResponse struct {
	Packages []PackageInfo `json:"packages"`
}

// InspectPackageResponse represents the response for inspect_package
type InspectPackageResponse struct {
	Package   PackageInfo     `json:"package"`
	Structs   []StructSummary `json:"structs"`
	Functions []FuncSummary   `json:"functions"`
	Methods   []MethodSummary `json:"methods"`
}

// StructDocResponse represents the response for get_doc_struct
type StructDocResponse struct {
	Name    string      `json:"name"`
	Comment string      `json:"comment"`
	Fields  []FieldDoc  `json:"fields"`
	Methods []MethodDoc `json:"methods"`
}

// FuncDocResponse represents the response for get_doc_func
type FuncDocResponse struct {
	Name      string    `json:"name"`
	Signature string    `json:"signature"`
	Comment   string    `json:"comment"`
	Examples  []Example `json:"examples"`
}

// MethodDocResponse represents the response for get_doc_method
type MethodDocResponse struct {
	ReceiverType string    `json:"receiver_type"`
	Name         string    `json:"name"`
	Signature    string    `json:"signature"`
	Comment      string    `json:"comment"`
	Examples     []Example `json:"examples"`
}

// ConstAndVarResponse represents the response for get_doc_const_and_var
type ConstAndVarResponse struct {
	Constants []ConstDoc `json:"constants"`
	Variables []VarDoc   `json:"variables"`
}
