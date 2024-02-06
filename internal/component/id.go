package component

import "strings"

const NameSpaceSeparator = "/"

type ID string

func (id ID) Namespace() string {
	idParts := strings.Split(string(id), NameSpaceSeparator)

	if len(idParts) == 1 {
		return ""
	}

	return strings.Join(idParts[:len(idParts)-1], NameSpaceSeparator) + NameSpaceSeparator
}

func (id ID) ComponentName() string {
	idParts := strings.Split(string(id), NameSpaceSeparator)

	if len(idParts) == 1 {
		return idParts[0]
	}

	return idParts[len(idParts)-1]
}

func (id ID) ExtendShortID(shortID string) string {
	if string(id) == NameSpaceSeparator+shortID {
		return shortID
	}

	extendedShortID := ID(strings.TrimSuffix(string(id), NameSpaceSeparator+shortID)).ComponentName()

	if shortID != "" {
		extendedShortID += NameSpaceSeparator + shortID
	}

	return extendedShortID
}
