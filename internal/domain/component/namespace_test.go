package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/domain/component"

	"github.com/stretchr/testify/assert"
)

func TestNamespace_LastSection(t *testing.T) {
	tests := []struct {
		name      string
		namespace component.Namespace
		want      string
	}{
		{
			name:      "The last section is not empty",
			namespace: "/internal/postgresql/connection",
			want:      "connection",
		},
		{
			name:      "Namespace is a empty string",
			namespace: "",
			want:      "",
		},
		{
			name:      "Namespace is the root",
			namespace: "/",
			want:      "",
		},
		{
			name:      "Namespace has only one section and the root",
			namespace: "/connection",
			want:      "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.namespace.LastSection())
		})
	}
}

func TestNamespace_ExtendComponentID(t *testing.T) {
	tests := []struct {
		name                string
		namespace           component.Namespace
		componentIDSections string
		want                string
	}{
		{
			name:                "There is on section, just the root",
			namespace:           "/",
			componentIDSections: "",
			want:                "/",
		},
		{
			name:                "There are no sections in a component id yet",
			namespace:           "/internal/postgresql/connection",
			componentIDSections: "",
			want:                "connection",
		},
		{
			name:                "There is only one section in a component id",
			namespace:           "/internal/postgresql/connection",
			componentIDSections: "connection",
			want:                "postgresql/connection",
		},
		{
			name:                "There are two sections in a component id",
			namespace:           "/internal/postgresql/connection",
			componentIDSections: "postgresql/connection",
			want:                "/internal/postgresql/connection",
		},
		{
			name:                "A component id cannot be extended",
			namespace:           "/internal/postgresql/connection",
			componentIDSections: "/internal/postgresql/connection",
			want:                "/internal/postgresql/connection",
		},
		{
			name:                "When a section is before the root, then the root is included in an id",
			namespace:           "/internal",
			componentIDSections: "",
			want:                "/internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.namespace.ExtendComponentID(tt.componentIDSections))
		})
	}
}
