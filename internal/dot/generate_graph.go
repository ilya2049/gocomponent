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

	componentGraph := project.CreateComponentGraph()

	if !conf.IncludeThirdPartyComponents {
		componentGraph = componentGraph.RemoveThirdPartyComponents()
	}

	if len(conf.IncludeParentComponents) > 0 {
		componentGraph = componentGraph.RemoveParentComponents(
			component.NewNamespaces(conf.IncludeParentComponents),
		)
	}

	if len(conf.IncludeChildComponents) > 0 {
		componentGraph = componentGraph.RemoveChildComponents(
			component.NewNamespaces(conf.IncludeChildComponents),
		)
	}

	dotExporter := newExporter(conf.ComponentColors)

	return dotExporter.export(componentGraph)
}
