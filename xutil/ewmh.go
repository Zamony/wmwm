package xutil

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/logging"
)

func getRoot(conn *xgb.Conn) (xproto.Window, error) {
	coninfo := xproto.Setup(conn)
	if coninfo == nil {
		err := errors.New("Coudn't parse X connection info")
		return xproto.Window(0), err
	}
	return coninfo.Roots[0].Root, nil
}

func SetSupported(conn *xgb.Conn) error {
	atoms := []xproto.Atom{
		GetAtom("_NET_SUPPORTED", conn),
		GetAtom("_NET_NUMBER_OF_DESKTOPS", conn),
		GetAtom("_NET_DESKTOP_NAMES", conn),
		GetAtom("_NET_CURRENT_DESKTOP", conn),
	}
	buf := make([]byte, len(atoms)*4)
	for i, atom := range atoms {
		xgb.Put32(buf[i*4:], uint32(atom))
	}
	root, err := getRoot(conn)
	if err != nil {
		return err
	}
	err = xproto.ChangePropertyChecked(
		conn, xproto.PropModeReplace, root,
		GetAtom("_NET_SUPPORTED", conn),
		xproto.AtomAtom, 32, uint32(len(atoms)), buf,
	).Check()
	return err
}

func SetNumberOfDesktops(n uint32, conn *xgb.Conn) error {
	root, err := getRoot(conn)
	if err != nil {
		return err
	}
	buf := make([]byte, 4)
	xgb.Put32(buf, uint32(n))
	err = xproto.ChangePropertyChecked(
		conn, xproto.PropModeReplace, root,
		GetAtom("_NET_NUMBER_OF_DESKTOPS", conn),
		xproto.AtomCardinal, 32, 1, buf,
	).Check()
	return err
}

func SetCurrentDesktop(n uint32, conn *xgb.Conn) error {
	root, err := getRoot(conn)
	if err != nil {
		return err
	}
	buf := make([]byte, 4)
	xgb.Put32(buf, uint32(n-1))
	err = xproto.ChangePropertyChecked(
		conn, xproto.PropModeReplace, root,
		GetAtom("_NET_CURRENT_DESKTOP", conn),
		xproto.AtomCardinal, 32, 1, buf,
	).Check()
	return err
}

func SetDesktopNames(names []string, conn *xgb.Conn) error {
	nullterm := make([]byte, 0)
	for _, name := range names {
		for i := 0; i < len(name); i++ {
			if name[i] != 0 {
				nullterm = append(nullterm, name[i])
			}
		}
		nullterm = append(nullterm, 0)
	}

	root, err := getRoot(conn)
	if err != nil {
		return err
	}

	err = xproto.ChangePropertyChecked(
		conn, xproto.PropModeReplace, root,
		GetAtom("_NET_DESKTOP_NAMES", conn),
		GetAtom("UTF8_STRING", conn), 8, uint32(len(nullterm)), nullterm,
	).Check()

	return err
}

func GetDesktopNames(conn *xgb.Conn) ([]string, error) {
	root, err := getRoot(conn)
	if err != nil {
		return nil, err
	}
	reply, err := xproto.GetProperty(
		conn, false, root, GetAtom("_NET_DESKTOP_NAMES", conn),
		xproto.GetPropertyTypeAny, 0, (1<<32)-1,
	).Reply()

	if err != nil {
		return nil, err
	}

	if reply.Format != 8 {
		return nil, errors.New("Error in getting property, not a string")
	}

	names := make([]string, 0)
	start := 0
	logging.Println("REPLYVAL", reply.Value)
	for i, c := range reply.Value {
		if c == 0 {
			// if string(reply.Value[start:i]) == string(make([]byte, 0)) {
			// 	logging.Println("SKIPPING")
			// 	start = i + 1
			// 	continue
			// }
			names = append(names, string(reply.Value[start:i]))
			start = i + 1
		}
	}
	if start < int(reply.ValueLen) {
		names = append(names, string(reply.Value[start:]))
	}
	return names, nil
}

func GetWMName(wid uint32, conn *xgb.Conn) (string, error) {
	reply, err := xproto.GetProperty(
		conn, false, xproto.Window(wid), GetAtom("_NET_WM_NAME", conn),
		xproto.GetPropertyTypeAny, 0, (1<<32)-1,
	).Reply()

	if err != nil || reply.Format != 8 {
		reply, err := xproto.GetProperty(
			conn, false, xproto.Window(wid), GetAtom("WM_NAME", conn),
			xproto.GetPropertyTypeAny, 0, (1<<32)-1,
		).Reply()

		if err != nil {
			return "", errors.New("Error in getting property WM_NAME")
		}

		if reply.Format != 8 {
			return "", errors.New("Error in getting property, not a string")
		}
		logging.Println("WMNAME", string(reply.Value))
		return string(reply.Value), nil
	}

	logging.Println("NETWMNAME", string(reply.Value))
	return string(reply.Value), nil
}

func IsDock(wid uint32, conn *xgb.Conn) bool {
	reply, err := xproto.GetProperty(
		conn, false, xproto.Window(wid), GetAtom("_NET_WM_WINDOW_TYPE", conn),
		xproto.GetPropertyTypeAny, 0, (1<<32)-1,
	).Reply()

	if err != nil {
		return false
	}
	if reply.Format != 32 {
		return false
	}

	dockAtom := GetAtom("_NET_WM_WINDOW_TYPE_DOCK", conn)
	values := reply.Value
	for len(values) >= 4 {
		atom := xproto.Atom(xgb.Get32(values))
		if atom == dockAtom {
			return true
		}

		values = values[4:]
	}

	return false
}
