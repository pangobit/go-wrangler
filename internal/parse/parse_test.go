package parse

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"
)

func TestParseStruct(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected StructInfo
	}{
		{
			name: "single field with bind tag",
			source: `package main

type User struct {
	Name string ` + "`" + `bind:"header,required"` + "`" + `
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "Name",
						FieldType: "string",
						Bind: &BindTag{
							Type:     "header",
							Required: true,
						},
					},
				},
			},
		},
		{
			name: "multiple fields with different tags",
			source: `package main

type User struct {
	Name  string ` + "`" + `bind:"header,required"` + "`" + `
	Email string ` + "`" + `bind:"query"` + "`" + `
	}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "Name",
						FieldType: "string",
						Bind: &BindTag{
							Type:     "header",
							Required: true,
						},
					},
					{
						FieldName: "Email",
						FieldType: "string",
						Bind: &BindTag{
							Type:     "query",
							Required: false,
						},
					},
				},
			},
		},
		{
			name: "bind tag with path type",
			source: `package main

type User struct {
	ID string ` + "`" + `bind:"path,required"` + "`" + `
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "ID",
						FieldType: "string",
						Bind: &BindTag{
							Type:     "path",
							Required: true,
						},
					},
				},
			},
		},
		{
			name: "no tags",
			source: `package main

type User struct {
	Name string
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{},
			},
		},
		{
			name: "only validate tag",
			source: `package main

type User struct {
	Name string ` + "`" + `validate:"min=10"` + "`" + `
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "Name",
						FieldType: "string",
						Validate:  &ValidateTag{Min: &[]int{10}[0]},
					},
				},
			},
		},
		{
			name: "field with both bind and validate tags",
			source: `package main

type User struct {
	ID int ` + "`" + `bind:"path,required" validate:"min=1,max=100"` + "`" + `
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "ID",
						FieldType: "int",
						Bind: &BindTag{
							Type:     "path",
							Required: true,
						},
						Validate: &ValidateTag{
							Min: &[]int{1}[0],
							Max: &[]int{100}[0],
						},
					},
				},
			},
		},
		{
			name: "field with multiple tags including unrelated",
			source: `package main

type User struct {
	Name string ` + "`" + `bind:"header,required" json:"name" db:"users.name"` + "`" + `
}`,
			expected: StructInfo{
				Name: "User",
				Tags: []TagInfo{
					{
						FieldName: "Name",
						FieldType: "string",
						Bind: &BindTag{
							Type:     "header",
							Required: true,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseStruct(tt.source)
			if err != nil {
				t.Fatalf("ParseStruct() error = %v", err)
			}

			if result.Name != tt.expected.Name {
				t.Errorf("ParseStruct() name = %v, want %v", result.Name, tt.expected.Name)
			}

			if len(result.Tags) != len(tt.expected.Tags) {
				t.Errorf("ParseStruct() got %d results, want %d", len(result.Tags), len(tt.expected.Tags))
				return
			}

			for i, expected := range tt.expected.Tags {
				actual := result.Tags[i]

				if actual.FieldName != expected.FieldName {
					t.Errorf("FieldName[%d] = %v, want %v", i, actual.FieldName, expected.FieldName)
				}

				if !reflect.DeepEqual(actual.Bind, expected.Bind) {
					t.Errorf("Bind[%d] = %v, want %v", i, actual.Bind, expected.Bind)
				}

				if !reflect.DeepEqual(actual.Validate, expected.Validate) {
					t.Errorf("Validate[%d] = %v, want %v", i, actual.Validate, expected.Validate)
				}
			}
		})
	}
}

func TestParseBindTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *BindTag
		hasError bool
	}{
		{
			name:  "header required",
			input: "header,required",
			expected: &BindTag{
				Type:     "header",
				Required: true,
			},
		},
		{
			name:  "query default optional",
			input: "query",
			expected: &BindTag{
				Type:     "query",
				Required: false,
			},
		},
		{
			name:  "path default optional",
			input: "path",
			expected: &BindTag{
				Type:     "path",
				Required: false,
			},
		},
		{
			name:     "invalid required value",
			input:    "header,maybe",
			hasError: true,
		},
		{
			name:     "explicit optional",
			input:    "query,optional",
			hasError: true,
		},
		{
			name:     "invalid bind type",
			input:    "invalid,required",
			hasError: true,
		},
		{
			name:     "empty",
			input:    "",
			hasError: true,
		},
		{
			name:  "header with spaces",
			input: " header ",
			expected: &BindTag{
				Type:     "header",
				Required: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseBindTag(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("parseBindTag() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseBindTag() error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseBindTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseValidateTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ValidateTag
		hasError bool
	}{
		{
			name:  "min only",
			input: "min=18",
			expected: &ValidateTag{
				Min: &[]int{18}[0],
			},
		},
		{
			name:  "max only",
			input: "max=65",
			expected: &ValidateTag{
				Max: &[]int{65}[0],
			},
		},
		{
			name:  "min and max",
			input: "min=10,max=20",
			expected: &ValidateTag{
				Min: &[]int{10}[0],
				Max: &[]int{20}[0],
			},
		},
		{
			name:  "max and min",
			input: "max=30,min=5",
			expected: &ValidateTag{
				Min: &[]int{5}[0],
				Max: &[]int{30}[0],
			},
		},
		{
			name:     "invalid min value",
			input:    "min=abc",
			hasError: true,
		},
		{
			name:     "unsupported rule",
			input:    "required",
			hasError: true,
		},
		{
			name:     "empty",
			input:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValidateTag(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("parseValidateTag() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseValidateTag() error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseValidateTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessField(t *testing.T) {
	// Helper to create a field with tag
	createField := func(name, tag string) *ast.Field {
		return &ast.Field{
			Names: []*ast.Ident{{Name: name}},
			Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`" + tag + "`"},
		}
	}

	tests := []struct {
		name     string
		field    *ast.Field
		expected TagInfo
		hasTag   bool
	}{
		{
			name:  "field with bind tag",
			field: createField("Name", `bind:"header,required"`),
			expected: TagInfo{
				FieldName: "Name",
				FieldType: "string",
				Bind: &BindTag{
					Type:     "header",
					Required: true,
				},
			},
			hasTag: true,
		},
		{
			name:  "field with validate tag",
			field: createField("Age", `validate:"min=18"`),
			expected: TagInfo{
				FieldName: "Age",
				FieldType: "int",
				Validate:  &ValidateTag{Min: &[]int{18}[0]},
			},
			hasTag: true,
		},
		{
			name:  "field with both tags",
			field: createField("ID", `bind:"path,required" validate:"max=100"`),
			expected: TagInfo{
				FieldName: "ID",
				FieldType: "int",
				Bind: &BindTag{
					Type:     "path",
					Required: true,
				},
				Validate: &ValidateTag{Max: &[]int{100}[0]},
			},
			hasTag: true,
		},
		{
			name:     "field with no tag",
			field:    &ast.Field{Names: []*ast.Ident{{Name: "Name"}}},
			expected: TagInfo{},
			hasTag:   false,
		},
		{
			name:     "field with invalid bind tag",
			field:    createField("Name", `bind:"invalid"`),
			expected: TagInfo{},
			hasTag:   false,
		},
		{
			name:     "field with invalid validate tag",
			field:    createField("Name", `validate:"invalid"`),
			expected: TagInfo{},
			hasTag:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := processField(tt.field)

			if ok != tt.hasTag {
				t.Errorf("processField() ok = %v, want %v", ok, tt.hasTag)
				return
			}

			if !tt.hasTag {
				return
			}

			if result.FieldName != tt.expected.FieldName {
				t.Errorf("FieldName = %v, want %v", result.FieldName, tt.expected.FieldName)
			}

			if !reflect.DeepEqual(result.Bind, tt.expected.Bind) {
				t.Errorf("Bind = %v, want %v", result.Bind, tt.expected.Bind)
			}

			if !reflect.DeepEqual(result.Validate, tt.expected.Validate) {
				t.Errorf("Validate = %v, want %v", result.Validate, tt.expected.Validate)
			}
		})
	}
}
