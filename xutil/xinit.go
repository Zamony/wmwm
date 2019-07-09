package xutil

import (
	"errors"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/kbrd"
)

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

func GrabShortcuts(conn *xgb.Conn, xroot xproto.ScreenInfo, keymap [256][]xproto.Keysym) {
	sym2code := make(map[xproto.Keysym]xproto.Keycode)
	needed := map[xproto.Keysym]uint8{
		kbrd.XK_BackSpace: 0, kbrd.XK_F1: 0,
		kbrd.XK_F2: 0, kbrd.XK_F3: 0, kbrd.XK_F4: 0, kbrd.XK_F5: 0,
		kbrd.XK_F6: 0, kbrd.XK_F7: 0, kbrd.XK_F8: 0, kbrd.XK_F9: 0,
		kbrd.XK_Left: 0, kbrd.XK_Right: 0, kbrd.XK_Up: 0, kbrd.XK_Down: 0,
		kbrd.XK_q: 0, kbrd.XK_Return: 0, kbrd.XK_grave: 0, kbrd.XK_t: 0,
	}
	for i, syms := range keymap {
		for _, sym := range syms {
			if _, ok := needed[sym]; ok {
				sym2code[sym] = xproto.Keycode(i)
			}
		}
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_t], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F1], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F2], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F3], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F4], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F5], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F6], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F7], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		uint16(0),
		sym2code[kbrd.XK_F8], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_Left], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_Right], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_Up], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_Down], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_q], xproto.GrabModeAsync,
		xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMaskControl|xproto.ModMask4,
		sym2code[kbrd.XK_Left], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMaskControl|xproto.ModMask4,
		sym2code[kbrd.XK_Right], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4|xproto.ModMask1,
		sym2code[kbrd.XK_Up], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4|xproto.ModMask1,
		sym2code[kbrd.XK_Down], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4|xproto.ModMask1,
		sym2code[kbrd.XK_Left], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4|xproto.ModMask1,
		sym2code[kbrd.XK_Right], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F1], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F2], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F3], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F4], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F5], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F6], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F7], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F8], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_F9], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}

	if err := xproto.GrabKeyChecked(
		conn, false, xroot.Root,
		xproto.ModMask4,
		sym2code[kbrd.XK_grave], xproto.GrabModeAsync, xproto.GrabModeAsync,
	).Check(); err != nil {
		log.Print(err)
	}
}

func GrabMouse(conn *xgb.Conn, xroot xproto.ScreenInfo) error {
	return xproto.GrabButtonChecked(
		conn, true, xroot.Root, xproto.EventMaskButtonPress,
		xproto.GrabModeSync, xproto.GrabModeSync, xproto.WindowNone,
		xproto.CursorNone, xproto.ButtonIndex1, xproto.ModMaskAny,
	).Check()
}

func GetScreens(conn *xgb.Conn) (main Screen, aux Screen, err error) {
	r, err := xinerama.QueryScreens(conn).Reply()
	if err != nil {
		return main, aux, err
	}

	nscreen := len(r.ScreenInfo)
	if nscreen < 1 {
		return main, aux, errors.New("No screen info available")
	}

	if nscreen == 1 {
		main = Screen{int(r.ScreenInfo[0].Width), int(r.ScreenInfo[0].Height), 0}
		aux = main
	} else {
		main = Screen{int(r.ScreenInfo[0].Width), int(r.ScreenInfo[0].Height), 0}
		aux = Screen{int(r.ScreenInfo[nscreen-1].Width), int(r.ScreenInfo[nscreen-1].Height), 1}
	}

	return main, aux, nil
}

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
