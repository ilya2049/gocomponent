package main

import (
	"flag"
	"fmt"

	"github.com/ilya2049/gocomponent/internal/fswalk"
)

func main() {
	projectDir := flag.String("project-dir", "", "project directory")
	rootNamespace := flag.String("root-namespace", "internal", "root namespace")

	flag.Parse()

	walk := fswalk.New(*projectDir, *rootNamespace)
	if err := walk.FindComponents(); err != nil {
		fmt.Println(err)

		return
	}

	walk.PrintDotGraph()
}
