package main

import (
	"fmt"
	"gocomponent/internal/fswalk"
)

const (
	projectDir    = ""
	rootNamespace = "internal/"
)

func main() {
	walk := fswalk.New(projectDir, rootNamespace)
	if err := walk.FindComponents(); err != nil {
		fmt.Println(err)

		return
	}

	walk.PrintComponents()
}
