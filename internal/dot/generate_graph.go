package dot

import (
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/fs"
)

func GenerateGraph() string {
	conf, err := config.Read()
	if err != nil {
		return err.Error()
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

	dotExporter := newExporter(conf.NamespaceColors)

	return dotExporter.export(project.Packages())
}
