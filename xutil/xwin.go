package xutil

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

var (
	ConfigureWindowChecked        = xproto.ConfigureWindowChecked
	MapWindowChecked              = xproto.MapWindowChecked
	UnmapWindowChecked            = xproto.UnmapWindowChecked
	ChangeWindowAttributesChecked = xproto.ChangeWindowAttributesChecked
	DestroyWindowChecked          = xproto.DestroyWindowChecked
	SetInputFocusChecked          = xproto.SetInputFocusChecked
)

func SetWindowX(x int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowX, []uint32{uint32(x)},
	).Check()
}

func SetWindowY(y int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY, []uint32{uint32(y)},
	).Check()
}

func SetWindowWidth(w int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowWidth, []uint32{uint32(w)},
	).Check()
}

func SetWindowHeight(h int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowHeight, []uint32{uint32(h)},
	).Check()
}

func MapWindow(wid uint32, conn *xgb.Conn) error {
	return MapWindowChecked(conn, xproto.Window(wid)).Check()
}

func UnmapWindow(wid uint32, conn *xgb.Conn) error {
	return UnmapWindowChecked(conn, xproto.Window(wid)).Check()
}

func RemoveWindowBorder(wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowBorderWidth, []uint32{uint32(0)},
	).Check()
}

func WatchWindowEvents(wid uint32, conn *xgb.Conn) error {
	return ChangeWindowAttributesChecked(
		conn, xproto.Window(wid),
		xproto.CwEventMask, []uint32{
			xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow,
		},
	).Check()
}

func DestroyWindow(wid uint32, conn *xgb.Conn) error {
	return DestroyWindowChecked(conn, xproto.Window(wid)).Check()
}

func FocusWindow(wid uint32, conn *xgb.Conn) error {
	return SetInputFocusChecked(
		conn, xproto.InputFocusPointerRoot,
		xproto.Window(wid), xproto.TimeCurrentTime,
	).Check()
}

func RemovePaddingFromWindow(y, height, bwidth int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{uint32(y - bwidth), uint32(height + bwidth)},
	).Check()
}

func AddPaddingToWindow(y, height, bwidth int, wid uint32, conn *xgb.Conn) error {
	return ConfigureWindowChecked(
		conn, xproto.Window(wid),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{uint32(y + bwidth), uint32(height - bwidth)},
	).Check()
}
