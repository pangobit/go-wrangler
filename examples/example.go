package main

import (
	"fmt"
	"log"

	"github.com/pangobit/go-wrangler/internal/parse"
)

func main() {
	// Example struct with bind and validate tags
	source := `
package main

type User struct {
	Name  string ` + "`bind:\"header,required\"`" + `
	Email string ` + "`bind:\"query\"`" + `
	Age   int    ` + "`validate:\"min=18,max=120\"`" + `
	ID    string ` + "`bind:\"path,required\" validate:\"max=10\"`" + `
}
`

	// Parse the struct
	tags, err := parse.ParseStruct(source)
	if err != nil {
		log.Fatalf("Failed to parse struct: %v", err)
	}

	// Print the parsed tag information
	fmt.Println("Parsed struct tags:")
	for _, tag := range tags {
		fmt.Printf("Field: %s\n", tag.FieldName)
		if tag.Bind != nil {
			fmt.Printf("  Bind: %s (required: %v)\n", tag.Bind.Type, tag.Bind.Required)
		}
		if tag.Validate != nil {
			if tag.Validate.Min != nil {
				fmt.Printf("  Validate Min: %d\n", *tag.Validate.Min)
			}
			if tag.Validate.Max != nil {
				fmt.Printf("  Validate Max: %d\n", *tag.Validate.Max)
			}
		}
		fmt.Println()
	}
}