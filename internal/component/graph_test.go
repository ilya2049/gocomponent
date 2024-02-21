package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/stretchr/testify/assert"
)

func TestGraph_MakeUniqueComponentIDs(t *testing.T) {
	component1 := component.New(component.NewNamespace("/postgresql/repository/user/edit"))
	component2 := component.New(component.NewNamespace("/domain/user/edit"))
	component3 := component.New(component.NewNamespace("/domain/product/edit"))
	component4 := component.New(component.NewNamespace("/pkg"))      // is not a section marker
	component5 := component.New(component.NewNamespace("pkg"))       // a section marker
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
