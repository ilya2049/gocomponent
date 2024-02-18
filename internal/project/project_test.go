package project_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/project"

	"github.com/stretchr/testify/assert"
)

func TestProject_MakeUniqueComponentIDs(t *testing.T) {
	project := project.New()

	_ = project.GetOrAddComponent(component.NewNamespace("postgresql/repository/user/edit"))
	_ = project.GetOrAddComponent(component.NewNamespace("domain/user/edit"))
	_ = project.GetOrAddComponent(component.NewNamespace("domain/product/edit"))
	_ = project.GetOrAddComponent(component.NewNamespace("pkg"))

	project.MakeUniqueComponentIDs()

	var uniqueIDs []string

	for _, c := range project.Components() {
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
