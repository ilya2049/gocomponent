package main

import (
	"flag"
	"fmt"

	"github.com/ilya2049/gocomponent/internal/fs"
)

func main() {
	projectDir := flag.String("project-dir", "", "project directory")
	rootNamespace := flag.String("root-namespace", "internal", "root namespace")

	flag.Parse()

	walk := fs.New(*projectDir, *rootNamespace)
	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	dotGraph := walk.ConvertComponentsAndImportsToDotGraphDotGraph()

	fmt.Println(dotGraph)
}
