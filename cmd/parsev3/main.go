package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	v3 "github.com/cloudfoundry-community/capiclientgen/pkg/v3"
)

// Parse the V3 API docs
// Usage: parsev3 <v3_api_docs_dir>

func usageAndExit() {
	fmt.Printf("Usage: parsev3 <v3_api_docs_dir>\n")
	os.Exit(1)
}

func filenameToResourceName(filename string) string {
	// get the leaf level dir name
	d := strings.TrimSuffix(filename, filepath.Base(filename))
	return filepath.Base(d)
}

func processMarkdownFile(filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Couldn't open markdown file '%s', error: '%s'", filename, err)
	}

	v3ApiRendererOpts := v3.RendererOptions{
		ResourceName: filenameToResourceName(filename),
	}
	v3ApiRenderer := v3.NewRenderer(v3ApiRendererOpts)
	exts := parser.CommonExtensions // parser.OrderedListStart | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(exts)
	doc := markdown.Parse(d, p)
	markdown.Render(doc, v3ApiRenderer)

	for _, e := range v3ApiRenderer.ParsedEndpoints() {
		fmt.Println(fmt.Sprintf("%+v", e))
	}

	// fmt.Printf("AST of file '%s':\n", filename)
	// ast.PrintWithPrefix(os.Stdout, doc, " ")
	// fmt.Print("\n")

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
