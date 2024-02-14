package dot

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/fs"
)

func GenerateGraph() string {
	configContents, err := os.ReadFile("config.toml")
	if err != nil {
		return fmt.Sprintf("read config: %s\n", err.Error())
	}

	var conf config.Config
	if _, err := toml.Decode(string(configContents), &conf); err != nil {
		return fmt.Sprintf("decode config: %s\n", err.Error())
	}

	project := component.NewProject()

	walk := fs.NewWalk(conf.ProjectDirectory, project)

	if err := walk.FindComponentsAndImports(); err != nil {
		return err.Error()
	}

	if conf.HideThirdPartyImports {
		project.ExcludeThirdPartyImports()
	}

	if len(conf.IncludeOnlyNextPackageNamespaces) > 0 {
		project.IncludeOnlyNextPackageNamespaces(conf.IncludeOnlyNextPackageNamespaces)
	}

	dotExporter := newExporter()

	return dotExporter.export(project.Packages())
}
