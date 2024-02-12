package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/dot"
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

	project := component.NewProject()

	walk := fs.NewWalk(conf.ProjectDirectory, project)

	if err := walk.FindComponentsAndImports(); err != nil {
		fmt.Println(err)

		return
	}

	if conf.HideThirdPartyImports {
		project.ExcludeThirdPartyImports()
	}

	dotExporter := dot.NewExporter()

	fmt.Println(dotExporter.Export(project.Packages()))
}
