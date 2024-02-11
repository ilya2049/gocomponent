package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_MakeUniqueComponentIDs(t *testing.T) {
	registry := component.NewRegistry(false)

	_ = registry.GetOrAddComponent(component.NewNamespace("postgresql/repository/user/edit"))
	_ = registry.GetOrAddComponent(component.NewNamespace("domain/user/edit"))
	_ = registry.GetOrAddComponent(component.NewNamespace("domain/product/edit"))
	_ = registry.GetOrAddComponent(component.NewNamespace("pkg"))

	registry.MakeUniqueComponentIDs()

	var uniqueIDs []string

	for _, c := range registry.Components() {
		uniqueIDs = append(uniqueIDs, c.ID())
	}

	want := []string{
		"repository/user/edit",
		"domain/user/edit",
		"product/edit",
		"pkg",
	}

	assert.ElementsMatch(t, want, uniqueIDs)
}
