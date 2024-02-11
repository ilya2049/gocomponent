package component

type Registry struct {
	components map[Namespace]*Component
}

func NewRegistry() *Registry {
	return &Registry{
		components: make(map[Namespace]*Component),
	}
}

func (r *Registry) GetOrAddComponent(namespace Namespace) *Component {
	existingComponent, ok := r.components[namespace]
	if ok {
		return existingComponent
	}

	newComponent := New(namespace)

	r.components[namespace] = newComponent

	return newComponent
}

func (r *Registry) MakeUniqueComponentIDs() {
	components := r.Components()

	for len(components) > 0 {
		firstComponent := components[0]
		isComponentIDUnique := true

		for i := 1; i < len(components); i++ {
			if components[i].ID() == firstComponent.ID() {
				isComponentIDUnique = false
				components[i].ExtendID()
			}
		}

		if isComponentIDUnique {
			components = components[1:]
		} else {
			firstComponent.ExtendID()
		}
	}
}

func (r *Registry) Components() []*Component {
	var components []*Component

	for _, component := range r.components {
		components = append(components, component)
	}

	return components
}
