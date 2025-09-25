// Package parse parses the incoming struct tags
package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// TagInfo represents the extracted tag information
type TagInfo struct {
	FieldName string
	Bind      *BindTag
	Validate  *ValidateTag
}

// BindTag represents bind tag information
// Type refers to the one of three possible options:
// - Header: http header params
// - Path: Path parameters, e.g., in /user/{id}, {id} would be the path parameter
// - Query: Query params from the URI
// Required is an optional tag, and is used to specify that a parameter must be present
// in order for the parameter validation to pass.
type BindTag struct {
	Type     string
	Required bool
}

// ValidateTag represents validate tag information for min and max validation on incoming int values
// Min specifies the minimum value (inclusive), nil if not specified
// Max specifies the maximum value (inclusive), nil if not specified
type ValidateTag struct {
	Min *int
	Max *int
}

// ParseStruct parses a Go struct source code and extracts tag information
func ParseStruct(source string) ([]TagInfo, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source: %w", err)
	}

	var tags []TagInfo

	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range structType.Fields.List {
			if tagInfo, ok := processField(field); ok {
				tags = append(tags, tagInfo)
			}
		}
		return true
	})

	return tags, nil
}

// processField processes a single struct field and extracts tag information
func processField(field *ast.Field) (TagInfo, bool) {
	if field.Tag == nil {
		return TagInfo{}, false
	}
	tag := strings.Trim(field.Tag.Value, "`")
	tagInfo := TagInfo{}

	if len(field.Names) > 0 {
		tagInfo.FieldName = field.Names[0].Name
	}

	if bindStr := extractTagValue(tag, "bind"); bindStr != "" {
		bindTag, err := parseBindTag(bindStr)
		if err != nil {
			// Skip invalid bind tags
			return TagInfo{}, false
		}
		tagInfo.Bind = bindTag
	}

	if validateStr := extractTagValue(tag, "validate"); validateStr != "" {
		validateTag, err := parseValidateTag(validateStr)
		if err != nil {
			// Skip invalid validate tags
			return TagInfo{}, false
		}
		tagInfo.Validate = validateTag
	}

	if tagInfo.Bind != nil || tagInfo.Validate != nil {
		return tagInfo, true
	}
	return TagInfo{}, false
}

// extractTagValue extracts the value for a specific tag key from the tag string
func extractTagValue(tagStr, key string) string {
	// Simple parsing - in a real implementation you'd want more robust parsing
	// Tags are in format: `key:"value" other:"value"`
	parts := strings.Fields(tagStr)
	for _, part := range parts {
		if strings.HasPrefix(part, key+":") {
			// Extract value between quotes
			if colonIdx := strings.Index(part, ":"); colonIdx != -1 {
				value := part[colonIdx+1:]
				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					return value[1 : len(value)-1]
				}
			}
		}
	}
	return ""
}

// parseBindTag parses the bind tag value
func parseBindTag(value string) (*BindTag, error) {
	parts := strings.Split(value, ",")
	if len(parts) == 0 || parts[0] == "" {
		return nil, fmt.Errorf("empty bind tag")
	}

	bindTag := &BindTag{}

	// First part is the type
	bindTag.Type = strings.TrimSpace(parts[0])
	switch bindTag.Type {
	case "header", "path", "query":
		// Valid
	default:
		return nil, fmt.Errorf("invalid bind type: %s", bindTag.Type)
	}

	// Required is implicit: present means required, absent means optional
	bindTag.Required = false
	if len(parts) > 1 {
		requiredStr := strings.TrimSpace(parts[1])
		if requiredStr == "required" {
			bindTag.Required = true
		} else {
			return nil, fmt.Errorf("invalid option: %s", requiredStr)
		}
	}

	return bindTag, nil
}

// parseValidateTag parses the validate tag value for min and max validation
func parseValidateTag(value string) (*ValidateTag, error) {
	parts := strings.Split(value, ",")
	validateTag := &ValidateTag{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if after, ok := strings.CutPrefix(part, "min="); ok {
			minStr := after
			minVal, err := strconv.Atoi(minStr)
			if err != nil {
				return nil, fmt.Errorf("invalid min value: %s", minStr)
			}
			validateTag.Min = &minVal
		} else if after, ok = strings.CutPrefix(part, "max="); ok {
			maxStr := after
			maxVal, err := strconv.Atoi(maxStr)
			if err != nil {
				return nil, fmt.Errorf("invalid max value: %s", maxStr)
			}
			validateTag.Max = &maxVal
		} else {
			return nil, fmt.Errorf("unsupported validation rule: %s", part)
		}
	}

	return validateTag, nil
}
