package xutil

type Screen struct {
	width  int
	height int
	id     uint32
}

func (screen *Screen) Width() int {
	return screen.width
}

func (screen *Screen) Height() int {
	return screen.height
}

func (screen *Screen) Id() uint32 {
	return screen.id
}
