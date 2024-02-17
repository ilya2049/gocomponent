package component

type Import struct {
	from *Component
	to   *Component
}

func NewImport(from, to *Component) *Import {
	return &Import{
		from: from,
		to:   to,
	}
}

func (i *Import) From() *Component {
	return i.from
}

func (i *Import) To() *Component {
	return i.to
}

type Imports []*Import
