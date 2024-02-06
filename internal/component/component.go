package component

import "github.com/ilya2049/gocomponent/internal/impset"

type Component struct {
	id      ID
	shortID string

	importSet impset.NamespaceImportSet
}

func New(id ID) *Component {
	return &Component{
		id:        id,
		shortID:   "",
		importSet: *impset.NewNamespaceImportSet(id.Namespace()),
	}
}

func (c *Component) ExtendShortID() {
	c.shortID = c.id.ExtendShortID(c.shortID)
}

func (c *Component) ShortID() string {
	return c.shortID
}

func (c *Component) ParseImportsOfGoFile(goFileName string) error {
	return c.importSet.ParseImportsOfGoFile(goFileName)
}
