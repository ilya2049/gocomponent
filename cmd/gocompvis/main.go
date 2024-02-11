package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/fs"
)

func main() {
	configContents, err := os.ReadFile("config.toml")
	if err != nil {
		fmt.Printf("read config: %s\n", err.Error())

		return
	}

	var conf config.Config
	if _, err := toml.Decode(string(configContents), &conf); err != nil {
		fmt.Printf("decode config: %s\n", err.Error())

		return
	}

	walk := fs.NewWalk(conf.ProjectDirectory)

	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	graph := walk.ConvertComponentsAndImportsToDotGraphDotGraph(conf.ShowThirdPartyImports)

	fmt.Println(graph)
}
