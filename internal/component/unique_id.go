package component

type ExtendableID interface {
	ExtendID()
	ID() string
}

func MakeUniqueIDs(extendableIDs []ExtendableID) {
	for len(extendableIDs) > 0 {
		firstExtendableID := extendableIDs[0]
		isIDUnique := true

		for i := 1; i < len(extendableIDs); i++ {
			if extendableIDs[i].ID() == firstExtendableID.ID() {
				isIDUnique = false
				extendableIDs[i].ExtendID()
			}
		}

		if isIDUnique {
			extendableIDs = extendableIDs[1:]
		} else {
			firstExtendableID.ExtendID()
		}
	}
}
