package generator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pangobit/go-wrangler/internal/parse"
)

func TestE2ECompilation(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name: "simple bind",
			source: `package main

type User struct {
	Name string ` + "`bind:\"header\"`" + `
}`,
		},
		{
			name: "bind and validate",
			source: `package main

type User struct {
	Name string ` + "`bind:\"header\"`" + `
	Age  int   ` + "`bind:\"path\" validate:\"min=18\"`" + `
}`,
		},
		{
			name: "required field",
			source: `package main

type User struct {
	Name string ` + "`bind:\"header,required\"`" + `
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the struct
			structInfo, err := parse.ParseStruct(tt.source)
			if err != nil {
				t.Fatalf("ParseStruct failed: %v", err)
			}

			// Generate the bind function
			code := GenerateBindFunction(structInfo)

			// Check that the generated code contains expected elements
			expectedContains := []string{
				fmt.Sprintf("func Bind%s", structInfo.Name),
				"r *http.Request",
				"pathParams map[string]string",
				fmt.Sprintf("s *%s", structInfo.Name),
				"return nil",
			}

			for _, expected := range expectedContains {
				if !strings.Contains(code, expected) {
					t.Errorf("Generated code does not contain expected string: %q", expected)
				}
			}

			// Check specific bindings based on the test case
			for _, tag := range structInfo.Tags {
				if tag.Bind != nil {
					switch tag.Bind.Type {
					case "header":
						if !strings.Contains(code, fmt.Sprintf("r.Header.Get(\"%s\")", tag.FieldName)) {
							t.Errorf("Expected header binding for %s", tag.FieldName)
						}
					case "query":
						if !strings.Contains(code, fmt.Sprintf("r.URL.Query().Get(\"%s\")", tag.FieldName)) {
							t.Errorf("Expected query binding for %s", tag.FieldName)
						}
					case "path":
						if !strings.Contains(code, fmt.Sprintf("pathParams[\"%s\"]", tag.FieldName)) {
							t.Errorf("Expected path binding for %s", tag.FieldName)
						}
					}
					if tag.Bind.Required {
						if !strings.Contains(code, fmt.Sprintf("%s is required", tag.FieldName)) {
							t.Errorf("Expected required check for %s", tag.FieldName)
						}
					}
				}
				if tag.Validate != nil {
					if tag.FieldType == "int" {
						if tag.Validate.Min != nil {
							if !strings.Contains(code, fmt.Sprintf("%s must be at least %d", tag.FieldName, *tag.Validate.Min)) {
								t.Errorf("Expected min validation for %s", tag.FieldName)
							}
						}
						if tag.Validate.Max != nil {
							if !strings.Contains(code, fmt.Sprintf("%s must be at most %d", tag.FieldName, *tag.Validate.Max)) {
								t.Errorf("Expected max validation for %s", tag.FieldName)
							}
						}
					}
				}
			}
		})
	}
}