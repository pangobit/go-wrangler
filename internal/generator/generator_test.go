package generator

import (
	"testing"

	"github.com/pangobit/go-wrangler/internal/parse"
)

func TestGenerateBindFunction(t *testing.T) {
	structInfo := parse.StructInfo{
		Name: "User",
		Tags: []parse.TagInfo{
			{
				FieldName: "Name",
				FieldType: "string",
				Bind: &parse.BindTag{
					Type:     "header",
					Required: true,
				},
			},
			{
				FieldName: "Email",
				FieldType: "string",
				Bind: &parse.BindTag{
					Type:     "query",
					Required: false,
				},
			},
			{
				FieldName: "Age",
				FieldType: "int",
				Validate: &parse.ValidateTag{
					Min: &[]int{18}[0],
					Max: &[]int{120}[0],
				},
			},
		},
	}

	result := GenerateBindFunction(structInfo)

	expected := `import (
	"fmt"
	"net/http"
)

func BindUser(r *http.Request, pathParams map[string]string, s *User) error {
	s.Name = r.Header.Get("Name")
	if s.Name == "" {
		return fmt.Errorf("Name is required")
	}
	s.Email = r.URL.Query().Get("Email")
	if s.Age < 18 {
		return fmt.Errorf("Age must be at least 18")
	}
	if s.Age > 120 {
		return fmt.Errorf("Age must be at most 120")
	}
	return nil
}
`

	if result != expected {
		t.Errorf("GenerateBindFunction() = %v, want %v", result, expected)
	}
}