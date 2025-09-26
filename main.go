package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pangobit/go-wrangler/internal/generator"
	"github.com/pangobit/go-wrangler/internal/parse"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <directory> [directories...]\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "\nGo Wrangler CLI tool for generating binding and validation code.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nStrategies:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  same    - Generate code in the same package (default)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  per     - Generate separate packages for each input\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  single  - Combine all in one package\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  %s examples\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "  %s --strategy per --target-dir ./gen --target-pkgs \"ofoo obar\" foo bar\n", os.Args[0])
}

func main() {
	flag.Usage = usage
	strategy := flag.String("strategy", "same", "Package strategy: same, per, single")
	targetPkg := flag.String("target-pkg", "", "Target package name for single strategy")
	targetDir := flag.String("target-dir", "", "Target directory for per or single strategy")
	targetPkgs := flag.String("target-pkgs", "", "Target package names for per strategy (space-separated)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("usage: go run main.go [flags] <directory> [directories...]")
	}

	dirs := args

	switch *strategy {
	case "same":
		processSame(dirs)
	case "per":
		if *targetDir == "" || *targetPkgs == "" {
			log.Fatal("per strategy requires --target-dir and --target-pkgs")
		}
		pkgList := strings.Fields(*targetPkgs)
		if len(pkgList) != len(dirs) {
			log.Fatal("number of target packages must match number of input directories")
		}
		processPer(dirs, pkgList, *targetDir)
	case "single":
		if *targetPkg == "" || *targetDir == "" {
			log.Fatal("single strategy requires --target-pkg and --target-dir")
		}
		processSingle(dirs, *targetPkg, *targetDir)
	default:
		log.Fatalf("Unknown strategy: %s", *strategy)
	}
}

func processSame(dirs []string) {
	for _, dir := range dirs {
		structs, pkgName, err := parse.ParsePackage(dir)
		if err != nil {
			log.Fatalf("Failed to parse package %s: %v", dir, err)
		}

		if len(structs) == 0 {
			fmt.Printf("No structs in %s\n", dir)
			continue
		}

		for _, s := range structs {
			fmt.Printf("Parsed struct: %s\n", s.Name)
		}

		outDir := dir
		outPkg := pkgName
		filePath := filepath.Join(outDir, pkgName+"_bindings.go")

		code := generator.GeneratePackage(structs, outPkg)

		err = os.WriteFile(filePath, []byte(code), 0644)
		if err != nil {
			log.Fatalf("Failed to write generated file: %v", err)
		}

		fmt.Printf("Generated code written to %s\n", filePath)
	}
}

func processPer(dirs []string, targetPkgs []string, targetDir string) {
	for i, dir := range dirs {
		structs, _, err := parse.ParsePackage(dir)
		if err != nil {
			log.Fatalf("Failed to parse package %s: %v", dir, err)
		}

		if len(structs) == 0 {
			fmt.Printf("No structs in %s\n", dir)
			continue
		}

		for _, s := range structs {
			fmt.Printf("Parsed struct: %s\n", s.Name)
		}

		outPkg := targetPkgs[i]
		outDir := filepath.Join(targetDir, outPkg)
		filePath := filepath.Join(outDir, "generated.go")

		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		code := generator.GeneratePackage(structs, outPkg)

		err = os.WriteFile(filePath, []byte(code), 0644)
		if err != nil {
			log.Fatalf("Failed to write generated file: %v", err)
		}

		fmt.Printf("Generated code written to %s\n", filePath)
	}
}

func processSingle(dirs []string, targetPkg, targetDir string) {
	allStructs := []parse.StructInfo{}
	for _, dir := range dirs {
		structs, _, err := parse.ParsePackage(dir)
		if err != nil {
			log.Fatalf("Failed to parse package %s: %v", dir, err)
		}

		if len(structs) == 0 {
			fmt.Printf("No structs in %s\n", dir)
			continue
		}

		for _, s := range structs {
			fmt.Printf("Parsed struct: %s\n", s.Name)
		}

		allStructs = append(allStructs, structs...)
	}

	if len(allStructs) == 0 {
		fmt.Println("No structs with bind or validate tags found.")
		return
	}

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	filePath := filepath.Join(targetDir, "generated.go")
	code := generator.GeneratePackage(allStructs, targetPkg)

	err = os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		log.Fatalf("Failed to write generated file: %v", err)
	}

	fmt.Printf("Generated code written to %s\n", filePath)
}