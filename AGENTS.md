# Go Wrangler Project

This document outlines the conventions and commands for the `go-wrangler` project.

## Commands

- **Build:** `go build ./...`
- **Test:** `go test ./...`
- **Run a single test:** `go test -run ^TestMyFunction$`
- **Lint:** `golangci-lint run` (assuming `golangci-lint` is installed)
- **Format:** `gofmt -w .` or `goimports -w .`

## Code Style

- **Imports:** Group standard library, third-party, and internal packages
separately. Use `goimports` to automate this.
- **Formatting:** Use `gofmt` or `goimports` to format code.
- **Types:** Use descriptive type names. Avoid single-letter type names except
for receivers.
- **Naming Conventions:**
  - Use `camelCase` for variables and functions.
  - Use `PascalCase` for exported identifiers.
  - Use `UPPER_SNAKE_CASE` for constants.
- **Error Handling:** Use `if err != nil` for error handling. Don't discard
errors with `_`.
- **Comments:** Provide at least one package-level comment, add comments to exported functions, and to all struct fields, in a go doc friendly format. Do not use inline comments to explain code. Code should be self-explanatory.
