package main

import (
	"flag"
	"fmt"

	"github.com/ilya2049/gocomponent/internal/fs"
)

func main() {
	projectDir := flag.String("dir", "./", "project directory")
	filterInProjectComponents := flag.Bool("in-project", false, "filter the 'in project' components")

	flag.Parse()

	walk := fs.NewWalk(*projectDir)

	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	graph := walk.ConvertComponentsAndImportsToDotGraphDotGraph(*filterInProjectComponents)

	fmt.Println(graph)
}
