# Go Wrangler

## Motivation

Building HTTP APIs in Go often requires extracting data from incoming requests (headers, query parameters, path parameters) and binding it to struct fields, followed by validation. Manually writing this binding and validation code is repetitive, error-prone, and time-consuming—especially as your API grows.

Go Wrangler solves this problem by providing a code generation tool that automates the creation of binding and validation functions based on struct tags. It eliminates boilerplate code, reduces bugs from manual implementation, and keeps your handlers clean and focused on business logic.

Use Go Wrangler when:
- You're building REST APIs or web services in Go
- You need to bind HTTP request data to structs
- You want to enforce validation rules on incoming data
- You prefer generated, consistent code over manual implementation
- You're using frameworks that don't provide automatic binding/validation

Go Wrangler is a stupid simple Go library for parsing struct tags and generating
HTTP request binding and validation code. It automates extracting data from HTTP
requests (headers, query parameters, path parameters, **form data**) and
validating it against struct field constraints.

## Features

- Parse Go struct tags for `bind` and `validate` directives
- Generate Go functions to bind HTTP request data to structs
- Support for header, query, path, and **form** parameter binding
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

## CLI Tool

Go Wrangler also provides a command-line tool to generate binding and validation code for entire packages.

### Installation

Build the CLI:

```bash
go build -o wrangler .
```

### Usage

Run the tool against one or more package directories:

```bash
./wrangler [flags] <directory> [directories...]
```

### Flags

- `--strategy`: Package strategy (`same`, `per`, `single`). Default: `same`
- `--target-pkg`: Target package name for `single` strategy
- `--target-dir`: Target directory for `per` or `single` strategy
- `--target-pkgs`: Target package names for `per` strategy (space-separated)

### Strategies

- `same`: Generate code in the same package, creating `<package_name>_bindings.go` in each input directory
- `per`: Generate separate packages for each input, requires `--target-dir` and `--target-pkgs` (must match number of inputs)
- `single`: Combine all structs into one package, requires `--target-dir` and `--target-pkg`

### Examples

```bash
# Generate in same package
./wrangler examples

# Generate separate packages
./wrangler --strategy per --target-dir ./gen --target-pkgs "bindings1 bindings2" pkg1 pkg2

# Generate all in one package
./wrangler --strategy single --target-dir ./gen --target-pkg combined pkg1 pkg2
```

### Using with `go tool` (Go 1.24+)

In Go 1.24 and later, you can add Go Wrangler as a tool dependency:

```bash
go get -tool github.com/pangobit/go-wrangler
```

Then run the tool using `go tool`:

```bash
go tool github.com/pangobit/go-wrangler [flags] <directories...>
```

This ensures the tool is available in your development environment without needing to build it separately.

### Using with `go generate`

You can integrate Go Wrangler into your build process using `go generate` to automatically generate binding code.

First, ensure you have Go Wrangler built or available. Then add a `go:generate` comment in your Go files:

```go
package main

//go:generate go run github.com/pangobit/go-wrangler --strategy same .

type User struct {
    Name     string `bind:"header,required"`
    Email    string `bind:"query"`
    Password string `bind:"form,required"`
    Age      int    `validate:"min=18"`
}
```

Run `go generate` in the directory to generate the code.

For more advanced usage, you can build the tool and use it directly:

```bash
go build -o wrangler github.com/pangobit/go-wrangler
# Then in go:generate
//go:generate ./wrangler --strategy same .
```

Note: When using `go run`, ensure the module is available in your GOPATH or use Go modules properly.

## Supported Tags

### Bind Tags

- `bind:"header"` - Bind from HTTP header
- `bind:"query"` - Bind from URL query parameter
- `bind:"path"` - Bind from URL path parameter
- `bind:"form"` - Bind from POST form data (application/x-www-form-urlencoded or multipart/form-data)
- `bind:"header=user_name"` - Bind from HTTP header with custom name "user_name"
- `bind:"query=email"` - Bind from URL query parameter with custom name "email"
- `bind:"path=id"` - Bind from URL path parameter with custom name "id"
- `bind:"form=password"` - Bind from POST form data with custom name "password"
- `bind:"header,required"` - Required header binding
- `bind:"form=user_name,required"` - Required form binding with custom name

### Validate Tags

- `validate:"min=18"` - Minimum value for integers
- `validate:"max=120"` - Maximum value for integers
- `validate:"min=10,max=100"` - Both min and max

## Form Data Binding

Form data binding is useful for POST requests that submit data as `application/x-www-form-urlencoded` or `multipart/form-data` (typically from HTML forms or API clients). The generated code uses Go's `http.Request.FormValue()` method, which automatically parses form data when called.

Form binding works with both single-part form data and multipart uploads, making it suitable for:
- Traditional HTML form submissions
- API requests with form-encoded data
- File uploads (when combined with multipart parsing)

**Example:**
```go
type CreateUserRequest struct {
    Username string `bind:"form=username,required"`
    Email    string `bind:"form=email,required"`
    Age      int    `bind:"form=age" validate:"min=13,max=120"`
}
```

This generates code that extracts values from `r.FormValue("username")`, `r.FormValue("email")`, and `r.FormValue("age")`.

## Testing

Run tests:

```bash
go test ./...
```

## License

Licensed under the MIT License. See LICENSE file for details.
