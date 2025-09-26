package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessSame(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test.go file with a struct
	testFile := filepath.Join(tempDir, "test.go")
	content := `package testpkg

type TestStruct struct {
	Name string ` + "`bind:\"query\"`" + `
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	processSame([]string{tempDir})

	// Check file created
	expectedFile := filepath.Join(tempDir, "testpkg_bindings.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s not created", expectedFile)
	}

	// Read content
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	str := string(data)

	if !strings.Contains(str, "package testpkg") {
		t.Errorf("Expected package testpkg in generated file")
	}
	if !strings.Contains(str, "func BindTestStruct") {
		t.Errorf("Expected BindTestStruct function in generated file")
	}
}

func TestProcessPer(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := t.TempDir()

	// Create a test.go file with a struct
	testFile := filepath.Join(tempDir, "test.go")
	content := `package testpkg

type TestStruct struct {
	Name string ` + "`bind:\"query\"`" + `
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	processPer([]string{tempDir}, []string{"otarget"}, targetDir)

	// Check file created
	expectedFile := filepath.Join(targetDir, "otarget", "generated.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s not created", expectedFile)
	}

	// Read content
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	str := string(data)

	if !strings.Contains(str, "package otarget") {
		t.Errorf("Expected package otarget in generated file")
	}
	if !strings.Contains(str, "func BindTestStruct") {
		t.Errorf("Expected BindTestStruct function in generated file")
	}
}

func TestProcessSingle(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := t.TempDir()

	// Create a test.go file with a struct
	testFile := filepath.Join(tempDir, "test.go")
	content := `package testpkg

type TestStruct struct {
	Name string ` + "`bind:\"query\"`" + `
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	processSingle([]string{tempDir}, "bindings", targetDir)

	// Check file created
	expectedFile := filepath.Join(targetDir, "generated.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s not created", expectedFile)
	}

	// Read content
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}
	str := string(data)

	if !strings.Contains(str, "package bindings") {
		t.Errorf("Expected package bindings in generated file")
	}
	if !strings.Contains(str, "func BindTestStruct") {
		t.Errorf("Expected BindTestStruct function in generated file")
	}
}