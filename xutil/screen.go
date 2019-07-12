package xutil

type Screen struct {
	width   int
	height  int
	xoffset int
	yoffset int
	id      uint32
}

func NewScreen(width, height, xoffset, yoffset int, id uint32) *Screen {
	return &Screen{width, height, xoffset, yoffset, id}
}

func (screen *Screen) Width() int {
	return screen.width
}

func (screen *Screen) Height() int {
	return screen.height
}

func (screen *Screen) XOffset() int {
	return screen.xoffset
}

func (screen *Screen) YOffset() int {
	return screen.yoffset
}

func (screen *Screen) Id() uint32 {
	return screen.id
}
