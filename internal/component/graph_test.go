package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/component/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph_MakeUniqueComponentIDs(t *testing.T) {
	tests := []struct {
		name             string
		newGraph         func() *component.Graph
		wantComponentIDs []string
	}{
		{
			name: "Many components",
			newGraph: func() *component.Graph {
				component1 := component.New(component.NewNamespace("/postgresql/repository/user/edit"))
				component2 := component.New(component.NewNamespace("/domain/user/edit"))
				component3 := component.New(component.NewNamespace("/domain/product/edit"))
				component4 := component.New(component.NewNamespace("/internal")) // already unique

				return component.NewGraph(component.Imports{
					component.NewImport(component1, component2),
					component.NewImport(component2, component3),
					component.NewImport(component3, component4),
				})
			},
			wantComponentIDs: []string{
				"repository/user/edit",
				"/domain/user/edit",
				"product/edit",
				"/internal",
			},
		},
		{
			name: "When a section is before the root, then the root is included in an id",
			newGraph: func() *component.Graph {
				component1 := component.New(component.NewNamespace("/postgresql/repository/user"))
				component2 := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(component1, component2),
				})
			},
			wantComponentIDs: []string{
				"user",
				"/pkg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			g := tt.newGraph()

			// When
			g.MakeUniqueComponentIDs()

			// Then
			componentIDs := testutil.GetComponentIDs(g)

			assert.ElementsMatch(t, tt.wantComponentIDs, componentIDs)
		})
	}
}

func TestGraph_CreateCustomComponents(t *testing.T) {
	tests := []struct {
		name             string
		newGraph         func() *component.Graph
		customComponents []string
		wantGraphString  string
	}{
		{
			name: "Create a section-marker custom component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				pgUser := component.New(component.NewNamespace("/postgresql/user"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, pgUser),
					component.NewImport(main, domainUser),
					component.NewImport(pgUser, domainUser),
					component.NewImport(pgUser, pkg),
					component.NewImport(domainUser, pkg),
				})
			},
			customComponents: []string{"user"},
			wantGraphString: testutil.BuildGraphString(
				"user -> /pkg",
				"/cmd/main -> user",
			),
		},
		{
			name: "Create a root-based custom component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				pgUser := component.New(component.NewNamespace("/postgresql/user"))
				pgProduct := component.New(component.NewNamespace("/postgresql/product"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, pgUser),
					component.NewImport(main, pgProduct),
					component.NewImport(pgUser, pgProduct),
					component.NewImport(pgProduct, pkg),
					component.NewImport(pgUser, pkg),
				})
			},
			customComponents: []string{"/postgresql"},
			wantGraphString: testutil.BuildGraphString(
				"/postgresql -> /pkg",
				"/cmd/main -> /postgresql",
			),
		},
		{
			name: "There are no custom components found",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				pgUser := component.New(component.NewNamespace("/postgresql/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, pgUser),
					component.NewImport(pgUser, pkg),
				})
			},
			customComponents: []string{"/mongodb"},
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /postgresql/user",
				"/postgresql/user -> /pkg",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := tt.newGraph()

			graphWithCustomComponents := graph.CreateCustomComponents(
				component.NewNamespaces(tt.customComponents),
			)

			assert.Equal(t, tt.wantGraphString, graphWithCustomComponents.String())
		})
	}
}

func TestGraph_String(t *testing.T) {
	tests := []struct {
		name            string
		newGraph        func() *component.Graph
		wantGraphString string
	}{
		{
			name: "The graph is not empty",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				pgUser := component.New(component.NewNamespace("/postgresql/user"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, pgUser),
					component.NewImport(main, domainUser),
					component.NewImport(pgUser, domainUser),
					component.NewImport(pgUser, pkg),
					component.NewImport(domainUser, pkg),
				})
			},
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /postgresql/user",
				"/cmd/main -> /domain/user",
				"/postgresql/user -> /domain/user",
				"/postgresql/user -> /pkg",
				"/domain/user -> /pkg",
			),
		},
		{
			name: "The graph is empty",
			newGraph: func() *component.Graph {
				return component.NewGraph(component.Imports{})
			},
			wantGraphString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphString := tt.newGraph().String()

			assert.Equal(t, tt.wantGraphString, graphString)
		})
	}
}

func TestGraph_ExtendComponentIDs(t *testing.T) {
	tests := []struct {
		name                       string
		newGraph                   func() *component.Graph
		idRegexpPatternAndSections map[string]int
		wantComponentIDs           []string
	}{
		{
			name: "Extend version-component id",
			newGraph: func() *component.Graph {
				component1 := component.New(component.NewNamespace("/internal"))
				component2 := component.New(component.NewNamespace("github.com/user/lib/v5"))

				return component.NewGraph(component.Imports{
					component.NewImport(component1, component2),
				})
			},
			idRegexpPatternAndSections: map[string]int{
				`v\d+$`: 2, // add two extra sections in a unique component id 'v5'
			},
			wantComponentIDs: []string{
				"/internal",
				"user/lib/v5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			g := tt.newGraph()

			g.MakeUniqueComponentIDs()

			// When
			err := g.ExtendComponentIDs(tt.idRegexpPatternAndSections)
			require.NoError(t, err)

			// Then
			componentIDs := testutil.GetComponentIDs(g)

			assert.ElementsMatch(t, tt.wantComponentIDs, componentIDs)
		})
	}
}

func TestGraph_ExtendComponentIDs_ErrorCase(t *testing.T) {
	// Given
	component1 := component.New(component.NewNamespace("/internal"))
	component2 := component.New(component.NewNamespace("github.com/user/lib/v5"))

	g := component.NewGraph(component.Imports{
		component.NewImport(component1, component2),
	})

	g.MakeUniqueComponentIDs()

	// When
	err := g.ExtendComponentIDs(map[string]int{
		"+": 1, // invalid regexp
	})

	// Then
	assert.Error(t, err)
}

func TestGraph_ExcludeChildComponents(t *testing.T) {
	tests := []struct {
		name                string
		newGraph            func() *component.Graph
		componentsToExclude component.Namespaces
		wantGraphString     string
	}{
		{
			name: "Exclude a root-based child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToExclude: component.NewNamespaces([]string{
				"/domain/user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/domain/user -> /pkg",
			),
		},
		{
			name: "Exclude a section-marker child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToExclude: component.NewNamespaces([]string{
				"user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/domain/user -> /pkg",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.newGraph()

			g = g.ExcludeChildComponents(tt.componentsToExclude)

			assert.Equal(t, tt.wantGraphString, g.String())
		})
	}
}

func TestGraph_ExcludeParentComponents(t *testing.T) {
	tests := []struct {
		name                string
		newGraph            func() *component.Graph
		componentsToExclude component.Namespaces
		wantGraphString     string
	}{
		{
			name: "Exclude a root-based parent component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToExclude: component.NewNamespaces([]string{
				"/domain/user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /domain/user",
			),
		},
		{
			name: "Exclude a section-marker parent component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToExclude: component.NewNamespaces([]string{
				"user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /domain/user",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.newGraph()

			g = g.ExcludeParentComponents(tt.componentsToExclude)

			assert.Equal(t, tt.wantGraphString, g.String())
		})
	}
}

func TestGraph_IncludeChildComponents(t *testing.T) {
	tests := []struct {
		name                string
		newGraph            func() *component.Graph
		componentsToInclude component.Namespaces
		wantGraphString     string
	}{
		{
			name: "Include a root-based child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToInclude: component.NewNamespaces([]string{
				"/domain/user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /domain/user",
			),
		},
		{
			name: "Include a section-marker child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToInclude: component.NewNamespaces([]string{
				"user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/cmd/main -> /domain/user",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.newGraph()

			g = g.IncludeChildComponents(tt.componentsToInclude)

			assert.Equal(t, tt.wantGraphString, g.String())
		})
	}
}

func TestGraph_IncludeParentComponents(t *testing.T) {
	tests := []struct {
		name                string
		newGraph            func() *component.Graph
		componentsToInclude component.Namespaces
		wantGraphString     string
	}{
		{
			name: "Include a root-based child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToInclude: component.NewNamespaces([]string{
				"/domain/user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/domain/user -> /pkg",
			),
		},
		{
			name: "Include a section-marker child component",
			newGraph: func() *component.Graph {
				main := component.New(component.NewNamespace("/cmd/main"))
				domainUser := component.New(component.NewNamespace("/domain/user"))
				pkg := component.New(component.NewNamespace("/pkg"))

				return component.NewGraph(component.Imports{
					component.NewImport(main, domainUser),
					component.NewImport(domainUser, pkg),
				})
			},
			componentsToInclude: component.NewNamespaces([]string{
				"user",
			}),
			wantGraphString: testutil.BuildGraphString(
				"/domain/user -> /pkg",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.newGraph()

			g = g.IncludeParentComponents(tt.componentsToInclude)

			assert.Equal(t, tt.wantGraphString, g.String())
		})
	}
}

func TestGraph_RemoveThirdPartyComponents(t *testing.T) {
	// Given
	main := component.New(component.NewNamespace("/cmd/main"))
	domainUser := component.New(component.NewNamespace("/domain/user"))
	pkg := component.New(component.NewNamespace("/pkg"))
	thirdParty := component.New(component.NewNamespace("github.com/user/lib/v5"))
	thirdParty.MarkAsThirdParty()

	g := component.NewGraph(component.Imports{
		component.NewImport(main, domainUser),
		component.NewImport(domainUser, pkg),
		component.NewImport(pkg, thirdParty),
	})

	// When
	g = g.RemoveThirdPartyComponents()

	// Then
	wantGraphString := testutil.BuildGraphString(
		"/cmd/main -> /domain/user",
		"/domain/user -> /pkg",
	)

	assert.Equal(t, wantGraphString, g.String())
}
