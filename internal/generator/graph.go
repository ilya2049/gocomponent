package generator

import (
	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
)

type fsWalker interface {
	CreateComponentGraph() (*component.Graph, error)
}

func GenerateGraph(conf *config.Config, walker fsWalker) (*component.Graph, error) {
	componentGraph, err := walker.CreateComponentGraph()
	if err != nil {
		return nil, err
	}

	if !conf.IncludeThirdPartyComponents {
		componentGraph = componentGraph.RemoveThirdPartyComponents()
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

	if len(conf.CustomComponents) > 0 {
		componentGraph = componentGraph.CreateCustomComponents(
			component.NewNamespaces(conf.CustomComponents),
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
		if err := componentGraph.ExtendComponentIDs(conf.ExtendComponentIDs); err != nil {
			return nil, err
		}
	}

	return componentGraph, nil
}
