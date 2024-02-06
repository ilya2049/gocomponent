package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestMakeUniqueShortIDs(t *testing.T) {
	components := []*component.Component{
		component.New("/postgresql/repository/user/edit"),
		component.New("/domain/user/edit"),
		component.New("/domain/product/edit"),
		component.New("/pkg"),
	}

	componentsWithUniqueShortIDs := component.MakeUniqueShortIDs(components)

	var shortIDs []string

	for _, c := range componentsWithUniqueShortIDs {
		shortIDs = append(shortIDs, c.ShortID())
	}

	want := []string{
		"repository/user/edit",
		"domain/user/edit",
		"product/edit",
		"pkg",
	}

	assert.Equal(t, want, shortIDs)
}
