package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestMakeUniqueIDs(t *testing.T) {
	components := []*component.Package{
		component.NewPackage("/postgresql/repository/user/edit"),
		component.NewPackage("/domain/user/edit"),
		component.NewPackage("/domain/product/edit"),
		component.NewPackage("/pkg"),
	}

	var shortIDExtenders []component.ExtendableID

	for _, c := range components {
		shortIDExtenders = append(shortIDExtenders, c)
	}

	shortIDExtenders = append(shortIDExtenders, component.New("/domain/user/edit"))

	component.MakeUniqueIDs(shortIDExtenders)

	var shortIDs []string

	for _, c := range components {
		shortIDs = append(shortIDs, c.ID())
	}

	want := []string{
		"repository/user/edit",
		"domain/user/edit",
		"domain/user/edit",
		"product/edit",
		"pkg",
	}

	assert.Equal(t, want, shortIDs)
}
