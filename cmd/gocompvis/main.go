package main

import (
	"flag"
	"fmt"

	"github.com/ilya2049/gocomponent/internal/fs"
)

func main() {
	projectDir := flag.String("dir", "./", "project directory")
	rootNamespace := flag.String("root", "internal", "root namespace")
	componentsHaveDoubleName := flag.Bool("double", false, "components have double name")

	flag.Parse()

	walk := fs.NewWalk(*projectDir, *rootNamespace)
	if *componentsHaveDoubleName {
		walk.ComponentsHaveDoubleName()
	}

	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	dotGraph := walk.ConvertComponentsAndImportsToDotGraphDotGraph()

	fmt.Println(dotGraph)
}
