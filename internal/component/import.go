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

func (i *Import) String() string {
	return i.From().Namespace().String() + " -> " + i.To().Namespace().String()
}

type Imports []*Import
