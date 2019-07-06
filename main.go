package main

import (
	"errors"
	"log"
	"os/exec"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/kbrd"
	"github.com/Zamony/wm/xutil"
)

func processEvents(conn *xgb.Conn, keymap [256][]xproto.Keysym, manager *WorkspaceManager) {
eventloop:
	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			log.Println(err)
			continue
		}

		switch e := event.(type) {
		case xproto.KeyPressEvent:
			err := handleKeyPress(conn, event.(xproto.KeyPressEvent), keymap, manager)
			if err != nil {
				break eventloop
			}
		case xproto.ConfigureRequestEvent:
			ev := xproto.ConfigureNotifyEvent{
				Event:            e.Window,
				Window:           e.Window,
				AboveSibling:     0,
				X:                e.X,
				Y:                e.Y,
				Width:            e.Width,
				Height:           e.Height,
				BorderWidth:      0,
				OverrideRedirect: false,
			}
			xproto.SendEventChecked(conn, false, e.Window, xproto.EventMaskStructureNotify, string(ev.Bytes()))
		case xproto.MapRequestEvent:
			win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
			win.Attach(manager.Curr())
			win.Activate(manager.Curr())
		default:
			log.Println(event)
		}
	}
}

func handleKeyPress(conn *xgb.Conn, key xproto.KeyPressEvent, keymap [256][]xproto.Keysym, manager *WorkspaceManager) error {
	keysym := keymap[key.Detail][0]

	switch keysym {
	case kbrd.XK_BackSpace:
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if ctrlActive && altActive {
			return errors.New("Error: quit event")
		}

		return nil

	case kbrd.XK_e:
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if ctrlActive && altActive {
			cmd := exec.Command("xterm")
			err := cmd.Start()
			go func() { cmd.Wait() }()
			if err != nil {
				return errors.New("Terminal launch failed")
			}
		}
		return nil
	case kbrd.XK_F1:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(1)
	case kbrd.XK_F2:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(2)
	case kbrd.XK_F3:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(3)
	case kbrd.XK_F4:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(4)
	case kbrd.XK_F5:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(5)
	case kbrd.XK_F6:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(6)
	case kbrd.XK_F7:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(7)
	case kbrd.XK_F8:
		win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
		win.Deactivate()
		win.Activate(8)
	case kbrd.XK_Left:
		winActive := (key.State & xproto.ModMask4) != 0
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && ctrlActive && !altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.ResizeLeft(manager.Curr())
		}

		if winActive && !ctrlActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.FocusLeft(manager.Curr())
		}

		if winActive && !ctrlActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.MoveLeft()
		}

	case kbrd.XK_Right:
		winActive := (key.State & xproto.ModMask4) != 0
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && ctrlActive && !altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.ResizeRight(manager.Curr())
		}

		if winActive && !ctrlActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.FocusRight(manager.Curr())
		}

		if winActive && !ctrlActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.MoveRight()
		}
	case kbrd.XK_Up:
		winActive := (key.State & xproto.ModMask4) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.FocusTop(manager.Curr())
		}

		if winActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.MoveUp()
		}

	case kbrd.XK_Down:
		winActive := (key.State & xproto.ModMask4) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.FocusBottom(manager.Curr())
		}
		if winActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.MoveDown()
		}
	case kbrd.XK_q:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Close(manager.Curr())
		}
	default:
		return nil
	}

	return nil
}

func main() {
	conn, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	coninfo := xproto.Setup(conn)
	if coninfo == nil {
		log.Fatal("Coudn't parse X connection info")
	}

	if len(coninfo.Roots) != 1 {
		log.Fatal("Number of roots > 1, Xinerama init failed")
	}
	root := coninfo.Roots[0]

	if err := xinerama.Init(conn); err != nil {
		log.Fatal(err)
	}

	mainscr, auxscr, err := xutil.GetScreens(conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := xutil.BecomeWM(conn, root); err != nil {
		log.Println(err)
		log.Fatal("Cannot take WM ownership")
	}

	keymap, err := kbrd.Mapping(conn)
	if err != nil {
		log.Fatal(err)
	}

	xutil.GrabShortcuts(conn, root, keymap)
	manager := NewWorkspaceManager(mainscr, auxscr)

	processEvents(conn, keymap, manager)
}
