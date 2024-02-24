package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/component/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGraph_MakeUniqueComponentIDs(t *testing.T) {
	component1 := component.New(component.NewNamespace("/postgresql/repository/user/edit"))
	component2 := component.New(component.NewNamespace("/domain/user/edit"))
	component3 := component.New(component.NewNamespace("/domain/product/edit"))
	component4 := component.New(component.NewNamespace("/pkg"))      // is not a section-marker
	component5 := component.New(component.NewNamespace("pkg"))       // a section-marker
	component6 := component.New(component.NewNamespace("/internal")) // already unique

	g := component.NewGraph(component.Imports{
		component.NewImport(component1, component2),
		component.NewImport(component2, component3),
		component.NewImport(component3, component4),
		component.NewImport(component4, component5),
		component.NewImport(component5, component6),
	})

	g.MakeUniqueComponentIDs()

	var uniqueIDs []string

	for _, c := range g.Components() {
		uniqueIDs = append(uniqueIDs, c.ID())
	}

	want := []string{
		"repository/user/edit",
		"domain/user/edit",
		"product/edit",
		"/pkg",
		"pkg",
		"internal",
	}

	assert.ElementsMatch(t, want, uniqueIDs)
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
