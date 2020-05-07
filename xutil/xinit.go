// Package xutil provides high-level abstraction for the XGB functions
package xutil

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wmwm/kbrd"
)

// Shortcut represents keys combination
type Shortcut struct {
	Modifiers uint16
	Keycode   xproto.Keycode
}

// BecomeWM asks X to grant it permission to manage windows
func BecomeWM(conn *xgb.Conn, xroot xproto.ScreenInfo) error {
	mask := []uint32{
		xproto.EventMaskKeyPress |
			xproto.EventMaskKeyRelease |
			xproto.EventMaskButtonPress |
			xproto.EventMaskButtonRelease |
			xproto.EventMaskStructureNotify |
			xproto.EventMaskSubstructureRedirect,
	}

	changed := xproto.ChangeWindowAttributesChecked(
		conn, xroot.Root, xproto.CwEventMask, mask,
	)

	return changed.Check()
}

// GrabShortcuts tells X that it should send
// specified keys combinations directly to the WM
func GrabShortcuts(conn *xgb.Conn, xroot xproto.ScreenInfo, keymap [256][]xproto.Keysym) error {
	sym2code := make(map[xproto.Keysym]xproto.Keycode)
	needed := map[xproto.Keysym]uint8{
		kbrd.XK_BackSpace: 0, kbrd.XK_F1: 0,
		kbrd.XK_F2: 0, kbrd.XK_F3: 0, kbrd.XK_F4: 0, kbrd.XK_F5: 0,
		kbrd.XK_F6: 0, kbrd.XK_F7: 0, kbrd.XK_F8: 0, kbrd.XK_F9: 0,
		kbrd.XK_Left: 0, kbrd.XK_Right: 0, kbrd.XK_Up: 0, kbrd.XK_Down: 0,
		kbrd.XK_q: 0, kbrd.XK_Return: 0, kbrd.XK_grave: 0, kbrd.XK_t: 0,
		kbrd.XK_f: 0, kbrd.XK_l: 0,
	}
	for i, syms := range keymap {
		for _, sym := range syms {
			if _, ok := needed[sym]; ok {
				sym2code[sym] = xproto.Keycode(i)
			}
		}
	}

	shortcuts := []Shortcut{
		{xproto.ModMask4, sym2code[kbrd.XK_t]},
		{xproto.ModMask4, sym2code[kbrd.XK_q]},
		{xproto.ModMask4, sym2code[kbrd.XK_grave]},
		{xproto.ModMask4, sym2code[kbrd.XK_f]},
		{xproto.ModMask4, sym2code[kbrd.XK_l]},
		{uint16(0), sym2code[kbrd.XK_F1]},
		{uint16(0), sym2code[kbrd.XK_F2]}, {uint16(0), sym2code[kbrd.XK_F3]},
		{uint16(0), sym2code[kbrd.XK_F4]}, {uint16(0), sym2code[kbrd.XK_F5]},
		{uint16(0), sym2code[kbrd.XK_F6]}, {uint16(0), sym2code[kbrd.XK_F7]},
		{uint16(0), sym2code[kbrd.XK_F8]}, {uint16(0), sym2code[kbrd.XK_F9]},
		{xproto.ModMask4, sym2code[kbrd.XK_Left]},
		{xproto.ModMask4, sym2code[kbrd.XK_Right]},
		{xproto.ModMask4, sym2code[kbrd.XK_Up]},
		{xproto.ModMask4, sym2code[kbrd.XK_Down]},
		{xproto.ModMaskControl | xproto.ModMask4, sym2code[kbrd.XK_Left]},
		{xproto.ModMaskControl | xproto.ModMask4, sym2code[kbrd.XK_Right]},
		{xproto.ModMask4 | xproto.ModMask1, sym2code[kbrd.XK_Up]},
		{xproto.ModMask4 | xproto.ModMask1, sym2code[kbrd.XK_Down]},
		{xproto.ModMask4 | xproto.ModMask1, sym2code[kbrd.XK_Left]},
		{xproto.ModMask4 | xproto.ModMask1, sym2code[kbrd.XK_Right]},
		{xproto.ModMask4, sym2code[kbrd.XK_F1]},
		{xproto.ModMask4, sym2code[kbrd.XK_F2]},
		{xproto.ModMask4, sym2code[kbrd.XK_F3]},
		{xproto.ModMask4, sym2code[kbrd.XK_F4]},
		{xproto.ModMask4, sym2code[kbrd.XK_F5]},
		{xproto.ModMask4, sym2code[kbrd.XK_F6]},
		{xproto.ModMask4, sym2code[kbrd.XK_F7]},
		{xproto.ModMask4, sym2code[kbrd.XK_F8]},
		{xproto.ModMask4, sym2code[kbrd.XK_F9]},
	}

	for _, shortcut := range shortcuts {
		err := xproto.GrabKeyChecked(
			conn, false, xroot.Root, shortcut.Modifiers,
			shortcut.Keycode, xproto.GrabModeAsync, xproto.GrabModeAsync,
		).Check()
		if err != nil {
			return err
		}
	}

	return nil
}

// GrabMouse tells X that it should send left mouse
// button click directly to the WM
func GrabMouse(conn *xgb.Conn, xroot xproto.ScreenInfo) error {
	return xproto.GrabButtonChecked(
		conn, true, xroot.Root, xproto.EventMaskButtonPress,
		xproto.GrabModeSync, xproto.GrabModeSync, xproto.WindowNone,
		xproto.CursorNone, xproto.ButtonIndex1, xproto.ModMaskAny,
	).Check()
}

// CreateCursor creates X cursor (XC_left_ptr)
func CreateCursor(conn *xgb.Conn) (xproto.Cursor, error) {
	cursor, err := xproto.NewCursorId(conn)
	if err != nil {
		return cursor, err
	}
	font, err := xproto.NewFontId(conn)
	if err != nil {
		return cursor, err
	}
	err = xproto.OpenFontChecked(conn, font, uint16(len("cursor")), "cursor").Check()
	if err != nil {
		return cursor, err
	}
	const xcLeftPtr = 68 // XC_left_ptr from cursorfont.h.
	err = xproto.CreateGlyphCursorChecked(
		conn, cursor, font, font, xcLeftPtr, xcLeftPtr+1,
		0xffff, 0xffff, 0xffff, 0, 0, 0,
	).Check()
	if err != nil {
		return cursor, err
	}
	err = xproto.CloseFontChecked(conn, font).Check()
	if err != nil {
		return cursor, err
	}
	return cursor, nil
}
