// Package kbrd provides constants representing physical keys
// and the mapping from physical keys to logical keys
package kbrd

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Mapping returns mapping from keycode to keysym
func Mapping(conn *xgb.Conn) ([256][]xproto.Keysym, error) {
	const low = 8
	const high = 255
	var mapping [256][]xproto.Keysym

	keymap, err := xproto.GetKeyboardMapping(conn, low, high-low+1).Reply()
	if err != nil {
		return mapping, err
	}

	if keymap == nil {
		return mapping, errors.New("Error getting keyboard mapping")
	}

	n := int(keymap.KeysymsPerKeycode)
	for i, k := low, 0; i <= high; i, k = i+1, k+n {
		mapping[i] = keymap.Keysyms[k : k+n]
	}

	return mapping, nil
}
