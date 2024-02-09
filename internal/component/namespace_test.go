package component_test

import (
	"testing"

	"github.com/ilya2049/gocomponent/internal/component"

	"github.com/stretchr/testify/assert"
)

func TestNamespace_ExcludeLastSection(t *testing.T) {
	tests := []struct {
		name      string
		namespace component.Namespace
		want      string
	}{
		{
			name:      "The last section is not empty",
			namespace: "internal/postgresql/connection",
			want:      "internal/postgresql/",
		},
		{
			name:      "Namespace is a empty string",
			namespace: "",
			want:      "",
		},
		{
			name:      "Namespace has only one section",
			namespace: "connection",
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.namespace.ExcludeLastSection())
		})
	}
}

func TestNamespace_LastSection(t *testing.T) {
	tests := []struct {
		name      string
		namespace component.Namespace
		want      string
	}{
		{
			name:      "The last section is not empty",
			namespace: "internal/postgresql/connection",
			want:      "connection",
		},
		{
			name:      "Namespace is a empty string",
			namespace: "",
			want:      "",
		},
		{
			name:      "Namespace has only one section",
			namespace: "connection",
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
		name        string
		namespace   component.Namespace
		componentID string
		want        string
	}{
		{
			name:        "There are not sections in a component id yet",
			namespace:   "internal/postgresql/connection",
			componentID: "",
			want:        "connection",
		},
		{
			name:        "There is only one section in a component id",
			namespace:   "internal/postgresql/connection",
			componentID: "connection",
			want:        "postgresql/connection",
		},
		{
			name:        "There are two sections in a component id",
			namespace:   "internal/postgresql/connection",
			componentID: "postgresql/connection",
			want:        "internal/postgresql/connection",
		},
		{
			name:        "A component id cannot be extended",
			namespace:   "internal/postgresql/connection",
			componentID: "internal/postgresql/connection",
			want:        "internal/postgresql/connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.namespace.ExtendComponentID(tt.componentID))
		})
	}
}
