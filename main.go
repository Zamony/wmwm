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
			log.Println(event)
			err := handleKeyPress(conn, event.(xproto.KeyPressEvent), keymap, manager)
			if err != nil {
				break eventloop
			}
		case xproto.ConfigureRequestEvent:
			log.Println(event)
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
			log.Println(event)
			wattr, err := xproto.GetWindowAttributes(conn, e.Window).Reply()
			if err != nil || !wattr.OverrideRedirect {
				log.Println("HAS NO OVERRIDEREDIRECT", e.Window)
				win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
				win.Reattach(manager.Curr())
				win.Activate(manager.Curr())
			} else {
				log.Println("HAS OVERRIDEREDIRECT", e.Window)
			}
		case xproto.DestroyNotifyEvent:
			log.Println(event)
			win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
			win.Remove()
		case xproto.ButtonPressEvent:
			log.Println(event)
			if e.Child > 0 {
				win := NewWindow(uint32(e.Child), manager.Mailbox(), conn)
				win.FocusHere()
			}
			xproto.AllowEventsChecked(conn, xproto.AllowReplayPointer, e.Time)
			xproto.AllowEventsChecked(conn, xproto.AllowReplayKeyboard, e.Time)
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

	case kbrd.XK_t:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			cmd := exec.Command("xterm")
			err := cmd.Start()
			go func() { cmd.Wait() }()
			if err != nil {
				return errors.New("Terminal launch failed")
			}
		}
		return nil
	case kbrd.XK_F1:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 1 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(1)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(1)
			manager.SetCurr(1)
		}
	case kbrd.XK_F2:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 2 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(2)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(2)
			manager.SetCurr(2)
		}
	case kbrd.XK_F3:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 3 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(3)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(3)
			manager.SetCurr(3)
		}
	case kbrd.XK_F4:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 4 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(4)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(4)
			manager.SetCurr(4)
		}
	case kbrd.XK_F5:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 5 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(5)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(5)
			manager.SetCurr(5)
		}
	case kbrd.XK_F6:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 6 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(6)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(6)
			manager.SetCurr(6)
		}
	case kbrd.XK_F7:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 7 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(7)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(7)
			manager.SetCurr(7)
		}
	case kbrd.XK_F8:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 8 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(8)
		} else if !winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.Deactivate(manager.Curr())
			win.Activate(8)
			manager.SetCurr(8)
		}
	case kbrd.XK_F9:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != 9 {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Attach(9)
		} else if !winActive {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.Activate(9)
			manager.SetCurr(9)
		}
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
	case kbrd.XK_grave:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			cmd := exec.Command("rofi", "-show", "run")
			err := cmd.Start()
			go func() { cmd.Wait() }()
			if err != nil {
				return errors.New("Application launcher failed")
			}
		}
		return nil
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
	cursor, err := xutil.CreateCursor(conn)
	if err != nil {
		log.Fatal(err)
	}
	if err := xproto.ChangeWindowAttributesChecked(
		conn, root.Root, xproto.CwBackPixel|xproto.CwCursor,
		[]uint32{root.BlackPixel, uint32(cursor)},
	).Check(); err != nil {
		log.Fatal(err)
	}

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

	xutil.GrabMouse(conn, root)
	xutil.GrabShortcuts(conn, root, keymap)
	manager := NewWorkspaceManager(mainscr, auxscr)

	processEvents(conn, keymap, manager)
}
