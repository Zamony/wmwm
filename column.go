package main

import (
	"errors"

	"github.com/Zamony/wm/logging"
	"github.com/Zamony/wm/xutil"
)

type Column struct {
	width      int
	x          int
	windows    []*Window
	screen     xutil.Screen
	fullscreen bool
}

func NewColumn(screen xutil.Screen) *Column {
	return &Column{
		screen.Width(), screen.XOffset(), nil,
		screen, false,
	}
}

func (column Column) Len() int {
	return len(column.windows)
}

func (column *Column) Add(window *Window) {
	column.windows = append(column.windows, window)
}

func (column *Column) Remove(window *Window) *Window {
	idx := column.IndexById(window.Id())
	if idx < 0 {
		return nil
	}
	column.windows = append(column.windows[:idx], column.windows[idx+1:]...)
	return window
}

func (column *Column) Swap(i, j int) error {
	if i < len(column.windows) && j < len(column.windows) && i > -1 && j > -1 {
		column.windows[i], column.windows[j] = column.windows[j], column.windows[i]
		return nil
	}
	return errors.New("Swapping values: index out of range")
}

func (column *Column) Reshape() {
	n := len(column.windows)
	if n < 1 {
		return
	}

	paddingT := column.screen.PaddingTop()
	paddingB := column.screen.PaddingBottom()
	if column.fullscreen {
		paddingT = 0
		paddingB = 0
	}

	height := column.screen.Height() - (paddingT + paddingB)
	h := height / n
	offsety := paddingT
	for i := 0; i < n-1; i++ {
		win := column.windows[i]
		win.SetY(offsety)
		win.SetX(column.x)
		win.SetHeight(h)
		win.SetWidth(column.width)
		offsety += h
	}

	win := column.windows[n-1]
	win.SetY(offsety)
	win.SetX(column.x)
	win.SetHeight(height + paddingT - offsety)
	win.SetWidth(column.width)
}

func (column *Column) HasPadding() bool {
	return !column.fullscreen
}

func (column *Column) AddPadding() {
	column.fullscreen = false
	column.Reshape()
}

func (column *Column) RemovePadding() {
	column.fullscreen = true
	column.Reshape()
}

func (column *Column) IndexById(wid uint32) int {
	for i := 0; i < len(column.windows); i++ {
		if column.windows[i].Id() == wid {
			return i
		}
	}
	return -1
}

func (column *Column) WindowByIndex(idx int) *Window {
	if len(column.windows) < 1 {
		return nil
	}

	if idx < len(column.windows) {
		return column.windows[idx]
	}

	return nil
}

func (column *Column) SetX(x int) int {
	column.x = column.screen.XOffset() + x
	return column.x
}

func (column *Column) SetWidth100() int {
	column.width = column.screen.Width()
	return column.width
}

func (column *Column) SetWidth50() int {
	column.width = column.screen.Width() / 2
	return column.width
}

func (column *Column) SetWidth80() int {
	w := float32(column.screen.Width())
	column.width = int(w * float32(0.65))
	return column.width
}

func (column *Column) SetWidth20() int {
	w := column.SetWidth80()
	column.width = column.screen.Width() - w
	return column.width
}

func (column Column) LogStatus() {
	logging.Println("(X:", column.x, "W:", column.width, ")")
	for _, win := range column.windows {
		win.LogStatus()
	}
}
