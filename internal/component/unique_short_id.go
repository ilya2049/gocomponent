package component

func MakeUniqueShortIDs(components []*Component) []*Component {
	var componentsWithUniqueShortIDs []*Component

	for len(components) > 0 {
		firstComponent := components[0]
		isShortIDUnique := true

		for i := 1; i < len(components); i++ {
			if components[i].shortID == firstComponent.shortID {
				isShortIDUnique = false
				components[i].ExtendShortID()
			}
		}

		if isShortIDUnique {
			componentsWithUniqueShortIDs = append(componentsWithUniqueShortIDs, firstComponent)
			components = components[1:]
		} else {
			firstComponent.ExtendShortID()
		}
	}

	return componentsWithUniqueShortIDs
}
