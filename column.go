// Package main implements logic of the window manager
package main

import (
	"errors"

	"github.com/Zamony/wm/logging"
	"github.com/Zamony/wm/xutil"
)

// Column represent group of windows having
// same height and position at x-axis
type Column struct {
	width      int
	x          int
	windows    []*Window
	screen     xutil.Screen
	fullscreen bool
}

// NewColumn creates instance of Column
func NewColumn(screen xutil.Screen) *Column {
	return &Column{
		screen.Width(), screen.XOffset(), nil,
		screen, false,
	}
}

// Len returns number of windows in the column
func (column Column) Len() int {
	return len(column.windows)
}

// Add adds window to the column
func (column *Column) Add(window *Window) {
	column.windows = append(column.windows, window)
}

// Remove removes window from column
func (column *Column) Remove(window *Window) *Window {
	idx := column.IndexById(window.Id())
	if idx < 0 {
		return nil
	}
	column.windows = append(column.windows[:idx], column.windows[idx+1:]...)
	return window
}

// Swap swaps windows in the column
func (column *Column) Swap(i, j int) error {
	if i < len(column.windows) && j < len(column.windows) && i > -1 && j > -1 {
		column.windows[i], column.windows[j] = column.windows[j], column.windows[i]
		return nil
	}
	return errors.New("Swapping values: index out of range")
}

// Reshape reshapes changes windows sizes in the column
// in such way that the have the same height and position on x-axis
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

// HasPadding checks whether column has padding
func (column *Column) HasPadding() bool {
	return !column.fullscreen
}

// AddPadding adds padding to the column
func (column *Column) AddPadding() {
	column.fullscreen = false
	column.Reshape()
}

// RemovePadding removes column's padding
func (column *Column) RemovePadding() {
	column.fullscreen = true
	column.Reshape()
}

// IndexById returns index of window by its id
func (column *Column) IndexById(wid uint32) int {
	for i := 0; i < len(column.windows); i++ {
		if column.windows[i].Id() == wid {
			return i
		}
	}
	return -1
}

// WindowByIndex return window by its index
func (column *Column) WindowByIndex(idx int) *Window {
	if len(column.windows) < 1 {
		return nil
	}

	if idx < len(column.windows) {
		return column.windows[idx]
	}

	return nil
}

// SetX sets column's x-coordinate
func (column *Column) SetX(x int) int {
	column.x = column.screen.XOffset() + x
	return column.x
}

// SetWidth100 sets column width to 100%
func (column *Column) SetWidth100() int {
	column.width = column.screen.Width()
	return column.width
}

// SetWidth50 sets column width to 50%
func (column *Column) SetWidth50() int {
	column.width = column.screen.Width() / 2
	return column.width
}

// SetWidth65 sets column width to 65%
func (column *Column) SetWidth65() int {
	w := float32(column.screen.Width())
	column.width = int(w * float32(0.65))
	return column.width
}

// SetWidth35 sets column width to 35%
func (column *Column) SetWidth35() int {
	w := column.SetWidth65()
	column.width = column.screen.Width() - w
	return column.width
}

// LogStatus logs column's information for debugging purposes
func (column Column) LogStatus() {
	logging.Println("(X:", column.x, "W:", column.width, ")")
	for _, win := range column.windows {
		win.LogStatus()
	}
}
