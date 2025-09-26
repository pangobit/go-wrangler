package main

import (
	"fmt"
	"log"

	"github.com/pangobit/go-wrangler/internal/generator"
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
	structInfo, err := parse.ParseStruct(source)
	if err != nil {
		log.Fatalf("Failed to parse struct: %v", err)
	}

	// Print the parsed tag information
	fmt.Printf("Struct Name: %s\n", structInfo.Name)
	fmt.Println("Parsed struct tags:")
	for _, tag := range structInfo.Tags {
		fmt.Printf("Field: %s (%s)\n", tag.FieldName, tag.FieldType)
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

	// Generate the bind and validate functions
	bindCode, _ := generator.GenerateBindFunction(structInfo)
	validateCode, _ := generator.GenerateValidateFunction(structInfo)

	// Print the generated bind function
	fmt.Println("Generated bind function:")
	fmt.Println(bindCode)

	// Print the generated validate function
	fmt.Println("Generated validate function:")
	fmt.Println(validateCode)
}