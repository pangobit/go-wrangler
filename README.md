# Go Wrangler

Go Wrangler is a stupid simple Go library for parsing struct tags and generating
HTTP request binding and validation code. It automates extracting data from HTTP
requests (headers, query parameters, path parameters) and
validating it against struct field constraints.

## Features

- Parse Go struct tags for `bind` and `validate` directives
- Generate Go functions to bind HTTP request data to structs
- Support for header, query, and path parameter binding
- Validation for min/max values on integer fields
- Required field enforcement

## Installation

```bash
go get github.com/pangobit/go-wrangler
```

## Usage

1. Define a struct with `bind` and `validate` tags
2. Parse the struct using the parser
3. Generate binding code using the generator
4. Use the generated function in your HTTP handlers

See the `examples/` directory for usage examples.

## Supported Tags

### Bind Tags

- `bind:"header"` - Bind from HTTP header
- `bind:"query"` - Bind from URL query parameter
- `bind:"path"` - Bind from URL path parameter
- `bind:"header,required"` - Required header binding

### Validate Tags

- `validate:"min=18"` - Minimum value for integers
- `validate:"max=120"` - Maximum value for integers
- `validate:"min=10,max=100"` - Both min and max

## Testing

Run tests:

```bash
go test ./...
```

## License

Licensed under the MIT License. See LICENSE file for details.
