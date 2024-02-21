package component

import "strings"

const SectionSeparator = "/"

type Namespace string

func NewNamespace(value string) Namespace {
	return Namespace(value)
}

func (ns Namespace) LastSection() string {
	sections := strings.Split(string(ns), SectionSeparator)

	if len(sections) == 1 {
		return sections[0]
	}

	return sections[len(sections)-1]
}

func (ns Namespace) ExtendComponentID(componentIDSections string) string {
	if string(ns) == componentIDSections {
		return componentIDSections
	}

	namespaceSectionsWithoutComponentIDSections := strings.TrimSuffix(string(ns), SectionSeparator+componentIDSections)

	sectionToExtend := Namespace(namespaceSectionsWithoutComponentIDSections).LastSection()

	if componentIDSections != "" {
		sectionToExtend += SectionSeparator + componentIDSections
	}

	return sectionToExtend
}

func (ns Namespace) Contains(another Namespace) bool {
	return strings.Contains(string(ns), string(another))
}

func (ns Namespace) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(ns), prefix)
}

func (ns Namespace) TrimPrefix(prefix string) Namespace {
	return Namespace(strings.TrimPrefix(string(ns), prefix))
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
