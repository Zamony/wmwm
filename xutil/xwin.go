// Package xutil provides high-level abstraction for the XGB functions
package xutil

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// The purpose of these declarations
// is to be monkey patched during testing
var (
	ConfigureWindowChecked        = xproto.ConfigureWindowChecked
	MapWindowChecked              = xproto.MapWindowChecked
	UnmapWindowChecked            = xproto.UnmapWindowChecked
	ChangeWindowAttributesChecked = xproto.ChangeWindowAttributesChecked
	DestroyWindowChecked          = xproto.DestroyWindowChecked
	SetInputFocusChecked          = xproto.SetInputFocusChecked
)

// SetWindowX sets window's x-coordinate value
func SetWindowX(x int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowX, []uint32{uint32(x)},
	).Check()
}

// SetWindowY sets window's y-coordinate value
func SetWindowY(y int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY, []uint32{uint32(y)},
	).Check()
}

// SetWindowWidth sets width of the window
func SetWindowWidth(w int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowWidth, []uint32{uint32(w)},
	).Check()
}

// SetWindowHeight sets height of the window
func SetWindowHeight(h int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowHeight, []uint32{uint32(h)},
	).Check()
}

// MapWindow maps window on the screen
func MapWindow(wid uint32, conn *xgb.Conn) error {
	return MapWindowChecked(conn, xproto.Window(wid)).Check()
}

// UnmapWindow unmaps window from screen
func UnmapWindow(wid uint32, conn *xgb.Conn) error {
	return UnmapWindowChecked(conn, xproto.Window(wid)).Check()
}

// RemoveWindowBorder removes window's border
func RemoveWindowBorder(wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowBorderWidth, []uint32{uint32(0)},
	).Check()
}

// WatchWindowEvents subscribes window manager to the
// MaskStructureNotify and MaskEnterWindow events
func WatchWindowEvents(wid uint32, conn *xgb.Conn) error {
	return ChangeWindowAttributesChecked(
		conn, xproto.Window(wid),
		xproto.CwEventMask, []uint32{
			xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow,
		},
	).Check()
}

// DestroyWindow destroys the window
func DestroyWindow(wid uint32, conn *xgb.Conn) error {
	return DestroyWindowChecked(conn, xproto.Window(wid)).Check()
}

// FocusWindow changes focus to window
func FocusWindow(wid uint32, conn *xgb.Conn) error {
	return SetInputFocusChecked(
		conn, xproto.InputFocusPointerRoot,
		xproto.Window(wid), xproto.TimeCurrentTime,
	).Check()
}

// RemovePaddingFromWindow removes padding from the window
func RemovePaddingFromWindow(y, height, bwidth int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{uint32(y - bwidth), uint32(height + bwidth)},
	).Check()
}

// AddPaddingToWindow adds padding to the window
func AddPaddingToWindow(y, height, bwidth int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{uint32(y + bwidth), uint32(height - bwidth)},
	).Check()
}
