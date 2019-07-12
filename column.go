package main

import (
	"errors"
	"fmt"

	"github.com/Zamony/wm/xutil"
)

type Column struct {
	width   int
	x       int
	windows []*Window
	screen  *xutil.Screen
}

func NewColumn(screen xutil.Screen) *Column {
	return &Column{
		screen.Width(), screen.XOffset(), nil,
		xutil.NewScreen(
			screen.Width(),
			screen.Height(),
			screen.XOffset(),
			screen.YOffset(),
			screen.Id(),
		),
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

	h := column.screen.Height() / n
	offsety := 0
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
	win.SetHeight(column.screen.Height() - offsety)
	win.SetWidth(column.width)
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
	fmt.Println(column.screen)
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

func (column Column) PrintStatus() {
	println("(X:", column.x, "W:", column.width, ")")
	for _, win := range column.windows {
		win.PrintStatus()
	}
}
