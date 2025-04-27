package model

import (
	"encoding/json"
	"testing"
)

func TestFormatPackageList(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		packages []PackageInfo
		want     string
	}{
		"multiple packages": {
			packages: []PackageInfo{
				{
					Name:       "pkg1",
					ImportPath: "github.com/example/pkg1",
					Comment:    "Package 1",
				},
				{
					Name:       "pkg2",
					ImportPath: "github.com/example/pkg2",
					Comment:    "Package 2",
				},
			},
			want: `{"packages":[{"name":"pkg1","import_path":"github.com/example/pkg1","comment":"Package 1"},{"name":"pkg2","import_path":"github.com/example/pkg2","comment":"Package 2"}]}`,
		},
		"empty packages": {
			packages: []PackageInfo{},
			want:     `{"packages":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatPackageList(tt.packages)
			// JSONの整形を統一するため、一度パースしてから比較
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatPackageList() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatPackageList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatPackageInspection(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pkg             PackageInfo
		structs         []StructSummary
		funcs           []FuncSummary
		methods         []MethodSummary
		includeComments bool
		want            string
	}{
		"package with all elements": {
			pkg: PackageInfo{
				Name:       "testpkg",
				ImportPath: "github.com/example/testpkg",
				Comment:    "Test package",
			},
			structs: []StructSummary{
				{
					Name:    "TestStruct",
					Comment: "Test struct",
				},
			},
			funcs: []FuncSummary{
				{
					Name:    "TestFunc",
					Comment: "Test function",
				},
			},
			methods: []MethodSummary{
				{
					ReceiverType: "TestStruct",
					Name:         "TestMethod",
					Comment:      "Test method",
				},
			},
			includeComments: true,
			want:            `{"package":{"name":"testpkg","import_path":"github.com/example/testpkg","comment":"Test package"},"structs":[{"name":"TestStruct","comment":"Test struct"}],"functions":[{"name":"TestFunc","comment":"Test function"}],"methods":[{"receiver_type":"TestStruct","name":"TestMethod","comment":"Test method"}]}`,
		},
		"empty package": {
			pkg: PackageInfo{
				Name:       "emptypkg",
				ImportPath: "github.com/example/emptypkg",
				Comment:    "",
			},
			structs:         []StructSummary{},
			funcs:           []FuncSummary{},
			methods:         []MethodSummary{},
			includeComments: false,
			want:            `{"package":{"name":"emptypkg","import_path":"github.com/example/emptypkg","comment":""},"structs":[],"functions":[],"methods":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatPackageInspection(tt.pkg, tt.structs, tt.funcs, tt.methods, tt.includeComments)
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatPackageInspection() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatPackageInspection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatStructDoc(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sName   string
		comment string
		fields  []FieldDoc
		methods []MethodDoc
		want    string
	}{
		"struct with fields and methods": {
			sName:   "TestStruct",
			comment: "Test struct documentation",
			fields: []FieldDoc{
				{
					Name:       "Field1",
					Type:       "string",
					Comment:    "Field1 documentation",
					IsExported: true,
				},
			},
			methods: []MethodDoc{
				{
					Name:      "Method1",
					Signature: "func (t *TestStruct) Method1() error",
					Comment:   "Method1 documentation",
				},
			},
			want: `{"name":"TestStruct","comment":"Test struct documentation","fields":[{"name":"Field1","type":"string","comment":"Field1 documentation","is_exported":true}],"methods":[{"name":"Method1","signature":"func (t *TestStruct) Method1() error","comment":"Method1 documentation"}]}`,
		},
		"empty struct": {
			sName:   "EmptyStruct",
			comment: "",
			fields:  []FieldDoc{},
			methods: []MethodDoc{},
			want:    `{"name":"EmptyStruct","comment":"","fields":[],"methods":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatStructDoc(tt.sName, tt.comment, tt.fields, tt.methods)
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatStructDoc() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatStructDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFuncDoc(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fName     string
		signature string
		comment   string
		examples  []Example
		want      string
	}{
		"function with examples": {
			fName:     "TestFunc",
			signature: "func TestFunc(arg string) error",
			comment:   "Test function documentation",
			examples: []Example{
				{
					Name:   "Example",
					Code:   "result := TestFunc(\"test\")\nfmt.Println(result)",
					Output: "nil",
				},
			},
			want: `{"name":"TestFunc","signature":"func TestFunc(arg string) error","comment":"Test function documentation","examples":[{"name":"Example","code":"result := TestFunc(\"test\")\nfmt.Println(result)","output":"nil"}]}`,
		},
		"function without examples": {
			fName:     "EmptyFunc",
			signature: "func EmptyFunc()",
			comment:   "",
			examples:  []Example{},
			want:      `{"name":"EmptyFunc","signature":"func EmptyFunc()","comment":"","examples":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatFuncDoc(tt.fName, tt.signature, tt.comment, tt.examples)
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatFuncDoc() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatFuncDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatMethodDoc(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		receiverType string
		mName        string
		signature    string
		comment      string
		examples     []Example
		want         string
	}{
		"method with examples": {
			receiverType: "TestStruct",
			mName:        "TestMethod",
			signature:    "func (t *TestStruct) TestMethod() error",
			comment:      "Test method documentation",
			examples: []Example{
				{
					Name:   "Example",
					Code:   "t := &TestStruct{}\nresult := t.TestMethod()",
					Output: "nil",
				},
			},
			want: `{"receiver_type":"TestStruct","name":"TestMethod","signature":"func (t *TestStruct) TestMethod() error","comment":"Test method documentation","examples":[{"name":"Example","code":"t := &TestStruct{}\nresult := t.TestMethod()","output":"nil"}]}`,
		},
		"method without examples": {
			receiverType: "EmptyStruct",
			mName:        "EmptyMethod",
			signature:    "func (e *EmptyStruct) EmptyMethod()",
			comment:      "",
			examples:     []Example{},
			want:         `{"receiver_type":"EmptyStruct","name":"EmptyMethod","signature":"func (e *EmptyStruct) EmptyMethod()","comment":"","examples":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatMethodDoc(tt.receiverType, tt.mName, tt.signature, tt.comment, tt.examples)
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatMethodDoc() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatMethodDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatConstAndVarDoc(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		constants []ConstDoc
		variables []VarDoc
		want      string
	}{
		"constants and variables": {
			constants: []ConstDoc{
				{
					Name:    "TestConst",
					Type:    "string",
					Value:   "\"test\"",
					Comment: "Test constant documentation",
				},
			},
			variables: []VarDoc{
				{
					Name:    "TestVar",
					Type:    "int",
					Comment: "Test variable documentation",
				},
			},
			want: `{"constants":[{"name":"TestConst","type":"string","value":"\"test\"","comment":"Test constant documentation"}],"variables":[{"name":"TestVar","type":"int","comment":"Test variable documentation"}]}`,
		},
		"empty constants and variables": {
			constants: []ConstDoc{},
			variables: []VarDoc{},
			want:      `{"constants":[],"variables":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := FormatConstAndVarDoc(tt.constants, tt.variables)
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Errorf("FormatConstAndVarDoc() invalid JSON = %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Errorf("want invalid JSON = %v", err)
				return
			}
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			if string(gotStr) != string(wantStr) {
				t.Errorf("FormatConstAndVarDoc() = %v, want %v", got, tt.want)
			}
		})
	}
}
