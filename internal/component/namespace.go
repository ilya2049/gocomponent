package component

import "strings"

const Slash = "/"

type Namespace string

func NewNamespace(value string) Namespace {
	return Namespace(value)
}

func (ns Namespace) LastSection() string {
	sections := strings.Split(ns.String(), Slash)

	if len(sections) == 1 {
		return sections[0]
	}

	return sections[len(sections)-1]
}

func (ns Namespace) ExtendComponentID(componentIDSections string) string {
	if ns.String() == componentIDSections {
		return componentIDSections
	}

	if ns.String() == Slash {
		return Slash
	}

	namespaceSectionsWithoutComponentIDSections := strings.TrimSuffix(ns.String(), Slash+componentIDSections)

	sectionToExtend := Namespace(namespaceSectionsWithoutComponentIDSections).LastSection()

	if componentIDSections != "" {
		sectionToExtend += Slash + componentIDSections
	}

	return sectionToExtend
}

func (ns Namespace) Contains(another Namespace) bool {
	if strings.HasPrefix(another.String(), Slash) {
		return strings.HasPrefix(ns.String(), another.String())
	}

	return strings.Contains(ns.String(), another.String())
}

func (ns Namespace) HasPrefix(prefix string) bool {
	return strings.HasPrefix(ns.String(), prefix)
}

func (ns Namespace) String() string {
	return string(ns)
}

func (ns Namespace) TrimPrefix(prefix string) Namespace {
	return Namespace(strings.TrimPrefix(ns.String(), prefix))
}

type Namespaces []Namespace

func NewNamespaces(values []string) Namespaces {
	namespaces := make(Namespaces, 0, len(values))

	for _, value := range values {
		namespaces = append(namespaces, NewNamespace(value))
	}

	return namespaces
}

func NewNamespaceColorMap(values map[string]string) map[Namespace]Color {
	namespaceColorMap := make(map[Namespace]Color)

	for namespaceValue, colorValue := range values {
		namespaceColorMap[NewNamespace(namespaceValue)] = NewColor(colorValue)
	}

	return namespaceColorMap
}
