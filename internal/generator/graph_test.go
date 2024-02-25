package generator_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"
	"github.com/ilya2049/gocomponent/internal/config"
	"github.com/ilya2049/gocomponent/internal/generator"
	"github.com/ilya2049/gocomponent/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateGraph(t *testing.T) {
	// Given
	conf := config.Config{
		IncludeThirdPartyComponents: true,
		IncludeParentComponents: []string{
			"/internal",
			"net/http",
		},
		IncludeChildComponents: []string{
			"/internal",
			"net/http",
		},
		ExcludeParentComponents: []string{
			"/internal/postgresql",
		},
		ExcludeChildComponents: []string{
			"/internal/postgresql",
		},
		CustomComponents: []string{
			"user",
		},
	}

	cmdMain := component.New(component.NewNamespace("/cmd/main"))
	internalPostgresql := component.New(component.NewNamespace("/internal/postgresql"))
	domainUser := component.New(component.NewNamespace("/internal/domain/user"))
	domainProduct := component.New(component.NewNamespace("/internal/domain/product"))
	appUser := component.New(component.NewNamespace("/internal/app/user"))
	appProduct := component.New(component.NewNamespace("/internal/app/product"))
	internalPkg := component.New(component.NewNamespace("/internal/pkg"))
	netHttp := component.New(component.NewNamespace("net/http"))

	g := component.NewGraph(component.Imports{
		component.NewImport(cmdMain, appUser),
		component.NewImport(cmdMain, appProduct),
		component.NewImport(internalPostgresql, domainUser),
		component.NewImport(appUser, domainUser),
		component.NewImport(internalPostgresql, domainProduct),
		component.NewImport(appProduct, domainProduct),
		component.NewImport(domainUser, internalPkg),
		component.NewImport(domainProduct, internalPkg),
		component.NewImport(internalPkg, netHttp),
	})

	fsWalker := testutil.NewFsWalkerStub(g)

	// When
	generatedComponentGraph, err := generator.GenerateGraph(&conf, fsWalker)
	require.NoError(t, err)

	// Then
	wantGeneratedComponentGraphString := testutil.BuildGraphString(
		"/internal/app/product -> /internal/domain/product",
		"/internal/domain/product -> /internal/pkg",
		"/internal/pkg -> net/http",
		"user -> /internal/pkg",
	)

	assert.Equal(t, wantGeneratedComponentGraphString, generatedComponentGraph.String())
}
