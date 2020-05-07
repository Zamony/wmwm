// Package xutil provides high-level abstraction for the XGB functions
package xutil

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wmwm/logging"
)

// GetAtom returns atom according to the specified string
func GetAtom(name string, conn *xgb.Conn) xproto.Atom {
	r, err := xproto.InternAtom(conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		logging.Fatal(err)
	}
	if r == nil {
		return 0
	}

	return r.Atom
}

// HasAtomDefined returns true if window has specified atom defined
func HasAtomDefined(atom string, wid uint32, conn *xgb.Conn) bool {
	prop, err := xproto.GetProperty(
		conn, false, xproto.Window(wid), GetAtom("WM_PROTOCOLS", conn),
		xproto.GetPropertyTypeAny, 0, 64,
	).Reply()

	if err != nil {
		log.Fatal(err)
	}

	atomv := GetAtom(atom, conn)
	if prop != nil {
		for v := prop.Value; len(v) >= 4; v = v[4:] {
			val := xproto.Atom(uint32(v[0]) | uint32(v[1])<<8 | uint32(v[2])<<16 | uint32(v[3])<<24)
			if val == atomv {
				return true
			}
		}
	}

	return false
}

// SendClientEvent sends client event
func SendClientEvent(atom string, timepoint, wid uint32, conn *xgb.Conn) error {
	return xproto.SendEventChecked(
		conn, false, xproto.Window(wid), xproto.EventMaskNoEvent,
		string(
			xproto.ClientMessageEvent{
				Format: 32,
				Window: xproto.Window(wid),
				Type:   GetAtom("WM_PROTOCOLS", conn),
				Data: xproto.ClientMessageDataUnionData32New(
					[]uint32{
						uint32(GetAtom(atom, conn)), timepoint,
						0, 0, 0,
					},
				),
			}.Bytes(),
		),
	).Check()
}
