package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPackageInfoJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pkg  PackageInfo
		want string
	}{
		"basic package info": {
			pkg: PackageInfo{
				Name:       "testpkg",
				ImportPath: "github.com/example/testpkg",
				Comment:    "This is a test package",
			},
			want: `{"name":"testpkg","import_path":"github.com/example/testpkg","comment":"This is a test package"}`,
		},
		"empty comment": {
			pkg: PackageInfo{
				Name:       "emptypkg",
				ImportPath: "github.com/example/emptypkg",
				Comment:    "",
			},
			want: `{"name":"emptypkg","import_path":"github.com/example/emptypkg","comment":""}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Marshal
			got, err := json.Marshal(tt.pkg)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
			}

			// Unmarshal
			var unmarshaled PackageInfo
			err = json.Unmarshal([]byte(tt.want), &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(unmarshaled, tt.pkg) {
				t.Errorf("json.Unmarshal() = %v, want %v", unmarshaled, tt.pkg)
			}
		})
	}
}

func TestListPackagesResponseJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		response ListPackagesResponse
		want     string
	}{
		"multiple packages": {
			response: ListPackagesResponse{
				Packages: []PackageInfo{
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
			},
			want: `{"packages":[{"name":"pkg1","import_path":"github.com/example/pkg1","comment":"Package 1"},{"name":"pkg2","import_path":"github.com/example/pkg2","comment":"Package 2"}]}`,
		},
		"empty packages": {
			response: ListPackagesResponse{
				Packages: []PackageInfo{},
			},
			want: `{"packages":[]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Marshal
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
			}

			// Unmarshal
			var unmarshaled ListPackagesResponse
			err = json.Unmarshal([]byte(tt.want), &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(unmarshaled, tt.response) {
				t.Errorf("json.Unmarshal() = %v, want %v", unmarshaled, tt.response)
			}
		})
	}
}

func TestStructDocResponseJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		response StructDocResponse
		want     string
	}{
		"struct with fields and methods": {
			response: StructDocResponse{
				Name:    "TestStruct",
				Comment: "Test struct documentation",
				Fields: []FieldDoc{
					{
						Name:       "Field1",
						Type:       "string",
						Comment:    "Field1 documentation",
						IsExported: true,
					},
				},
				Methods: []MethodDoc{
					{
						Name:      "Method1",
						Signature: "func (t *TestStruct) Method1() error",
						Comment:   "Method1 documentation",
					},
				},
			},
			want: `{"name":"TestStruct","comment":"Test struct documentation","fields":[{"name":"Field1","type":"string","comment":"Field1 documentation","is_exported":true}],"methods":[{"name":"Method1","signature":"func (t *TestStruct) Method1() error","comment":"Method1 documentation"}]}`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Marshal
			got, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.want)
			}

			// Unmarshal
			var unmarshaled StructDocResponse
			err = json.Unmarshal([]byte(tt.want), &unmarshaled)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(unmarshaled, tt.response) {
				t.Errorf("json.Unmarshal() = %v, want %v", unmarshaled, tt.response)
			}
		})
	}
}
