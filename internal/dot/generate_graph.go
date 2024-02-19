package dot

import (
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/fs"
	"github.com/ilya2049/gocomponent/internal/project"
)

func GenerateGraph() string {
	conf, err := config.Read()
	if err != nil {
		return err.Error()
	}

	project := project.New()

	walk := fs.NewWalk(conf.ProjectDirectory, project)

	if err := walk.FindComponentsAndImports(); err != nil {
		return err.Error()
	}

	componentGraph := project.CreateComponentGraph()

	if !conf.IncludeThirdPartyComponents {
		componentGraph = componentGraph.RemoveThirdPartyComponents()
	}

	if len(conf.CustomComponents) > 0 {
		componentGraph = componentGraph.CreateCustomComponents(
			component.NewNamespaces(conf.CustomComponents),
		)
	}

	if len(conf.IncludeParentComponents) > 0 {
		componentGraph = componentGraph.IncludeParentComponents(
			component.NewNamespaces(conf.IncludeParentComponents),
		)
	}

	if len(conf.IncludeChildComponents) > 0 {
		componentGraph = componentGraph.IncludeChildComponents(
			component.NewNamespaces(conf.IncludeChildComponents),
		)
	}

	if len(conf.ExcludeParentComponents) > 0 {
		componentGraph = componentGraph.ExcludeParentComponents(
			component.NewNamespaces(conf.ExcludeParentComponents),
		)
	}

	if len(conf.ExcludeChildComponents) > 0 {
		componentGraph = componentGraph.ExcludeChildComponents(
			component.NewNamespaces(conf.ExcludeChildComponents),
		)
	}

	if len(conf.ComponentColors) > 0 {
		componentGraph.Colorize(component.NewNamespaceColorMap(conf.ComponentColors))
	}

	if conf.ThirdPartyComponentsColor != "" {
		componentGraph.ColorizeThirdParty(component.NewColor(conf.ThirdPartyComponentsColor))
	}

	componentGraph.MakeUniqueComponentIDs()

	if len(conf.ExtendComponentIDs) > 0 {
		componentGraph.ExtendComponentIDs(conf.ExtendComponentIDs)
	}

	dotExporter := newExporter()

	return dotExporter.export(componentGraph)
}
