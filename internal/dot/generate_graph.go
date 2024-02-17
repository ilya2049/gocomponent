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

	if len(conf.IncludeOnlyNextPackageNamespaces) > 0 {
		project.IncludeOnlyNextPackageNamespaces(conf.IncludeOnlyNextPackageNamespaces)
	}

	componentGraph := project.CreateComponentGraph()

	if !conf.IncludeThirdPartyComponents {
		componentGraph = componentGraph.RemoveThirdPartyComponents()
	}

	dotExporter := newExporter(conf.NamespaceColors)

	return dotExporter.export(componentGraph)
}
