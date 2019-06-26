package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Parse the V3 API docs
// Usage: parsev3 <v3_api_docs_dir>

func usageAndExit() {
	fmt.Printf("Usage: parsev3 <v3_api_docs_dir>\n")
	os.Exit(1)
}

func processMarkdownFile(filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Couldn't open markdown file '%s', error: '%s'", filename, err)
	}

	exts := parser.CommonExtensions // parser.OrderedListStart | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(exts)
	doc := markdown.Parse(d, p)
	fmt.Printf("AST of file '%s':\n", filename)
	ast.PrintWithPrefix(os.Stdout, doc, " ")
	fmt.Print("\n")

	return nil
}

func main() {
	nFiles := len(os.Args) - 1
	if nFiles != 1 {
		usageAndExit()
	}

	v3ApiDocsRootDir := os.Args[1]
	v3ApiResourcesRootDir := filepath.Join(v3ApiDocsRootDir, "source/includes/resources")

	if _, err := os.Stat(v3ApiResourcesRootDir); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Can't find directory '%s'", v3ApiResourcesRootDir))
		os.Exit(1)
	}

	err := filepath.Walk(v3ApiResourcesRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(info.Name(), ".md") {
			return processMarkdownFile(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the the v3 API dir %q: %v\n", v3ApiResourcesRootDir, err)
		os.Exit(1)
	}
}

