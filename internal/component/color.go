package component

type Color string

func NewColor(value string) Color {
	return Color(value)
}

func (c Color) String() string {
	return string(c)
}
