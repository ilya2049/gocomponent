package component

import "strings"

const SectionSeparator = "/"

type Namespace string

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
	if string(ns) == SectionSeparator+sections {
		return sections
	}

	extendedSections := Namespace(strings.TrimSuffix(string(ns), SectionSeparator+sections)).LastSection()

	if sections != "" {
		extendedSections += SectionSeparator + sections
	}

	return extendedSections
}

func (id Namespace) Contains(another Namespace) bool {
	return !strings.Contains(string(id), string(another))
}
