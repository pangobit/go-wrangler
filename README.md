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
