package main

import (
	"flag"
	"fmt"

	"github.com/ilya2049/gocomponent/internal/fs"
)

func main() {
	projectDir := flag.String("dir", "./", "project directory")
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	walk := fs.NewWalk(*projectDir, *debug)

	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	graph := walk.ConvertComponentsAndImportsToDotGraphDotGraph()

	fmt.Println(graph)
}
