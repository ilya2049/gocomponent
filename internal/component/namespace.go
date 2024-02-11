package component

import "strings"

const SectionSeparator = "/"

type Namespace string

func NewNamespace(value string) Namespace {
	return Namespace(value)
}

func (ns Namespace) ExcludeLastSection() string {
	sections := strings.Split(string(ns), SectionSeparator)

	if len(sections) == 1 {
		return ""
	}

	return strings.Join(sections[:len(sections)-1], SectionSeparator) + SectionSeparator
}

func (ns Namespace) LastSection() string {
	sections := strings.Split(string(ns), SectionSeparator)

	if len(sections) == 1 {
		return sections[0]
	}

	return sections[len(sections)-1]
}

func (ns Namespace) ExtendComponentID(sections string) string {
	if string(ns) == sections {
		return sections
	}

	extendedSections := Namespace(strings.TrimSuffix(string(ns), SectionSeparator+sections)).LastSection()

	if sections != "" {
		extendedSections += SectionSeparator + sections
	}

	return extendedSections
}

func (ns Namespace) Contains(another Namespace) bool {
	return !strings.Contains(string(ns), string(another))
}

func (ns Namespace) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(ns), prefix)
}

func (ns Namespace) TrimPrefix(prefix string) Namespace {
	return Namespace(strings.TrimPrefix(string(ns), prefix))
}
