// Package main implements logic of the window manager
package main

import (
	"errors"
	"os/exec"
	"regexp"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/config"
	"github.com/Zamony/wm/kbrd"
	"github.com/Zamony/wm/logging"
	"github.com/Zamony/wm/xutil"
)

func processEvents(
	conn *xgb.Conn, keymap [256][]xproto.Keysym,
	monitors xutil.MonitorsInfo, manager *WorkspaceManager,
) {
eventloop:
	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			logging.Println(err)
			continue
		}

		switch e := event.(type) {
		case xproto.KeyPressEvent:
			err := handleKeyPress(
				conn, event.(xproto.KeyPressEvent), keymap, manager,
			)
			if err != nil {
				break eventloop
			}
		case xproto.ConfigureRequestEvent:
			logging.Println(event)
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
			xproto.SendEventChecked(
				conn, false, e.Window, xproto.EventMaskStructureNotify,
				string(ev.Bytes()),
			)
		case xproto.MapRequestEvent:
			logging.Println(event)
			wattr, err := xproto.GetWindowAttributes(conn, e.Window).Reply()
			if err != nil || !wattr.OverrideRedirect {
				win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
				win.SendAttach(manager.Curr())
				win.SendActivate(manager.Curr())
			}
		case xproto.UnmapNotifyEvent:
			logging.Println(event)
			win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
			win.SendRemove()
		case xproto.DestroyNotifyEvent:
			logging.Println(event)
			win := NewWindow(uint32(e.Window), manager.Mailbox(), conn)
			win.SendRemove()
		case xproto.ButtonPressEvent:
			logging.Println(event)
			if e.Child > 0 {
				win := NewWindow(uint32(e.Child), manager.Mailbox(), conn)
				if monitors.IsDualSetup() {
					secondMonitorClick := !monitors.InPrimaryRegion(int(e.RootX))
					secondMonitorActive := manager.Curr() == manager.SpecialWorkspace()
					switch {
					case !secondMonitorClick && secondMonitorActive:
						win.SendActivate(manager.Prev())
						manager.SetCurr(manager.Prev())
					case secondMonitorClick && !secondMonitorActive:
						win.SendActivate(manager.SpecialWorkspace())
						manager.SetCurr(manager.SpecialWorkspace())
					}
				}

				win.SendFocusHere()
			}
			xproto.AllowEventsChecked(conn, xproto.AllowReplayPointer, e.Time)
			xproto.AllowEventsChecked(conn, xproto.AllowReplayKeyboard, e.Time)
		default:
			logging.Println(event)
		}
	}
}

func handleKeyPress(conn *xgb.Conn, key xproto.KeyPressEvent, keymap [256][]xproto.Keysym, manager *WorkspaceManager) error {
	keysym := keymap[key.Detail][0]
	fkeys := map[xproto.Keysym]uint32{
		kbrd.XK_F1: 1, kbrd.XK_F2: 2, kbrd.XK_F3: 3, kbrd.XK_F4: 4,
		kbrd.XK_F5: 5, kbrd.XK_F6: 6, kbrd.XK_F7: 7, kbrd.XK_F8: 8, kbrd.XK_F9: 9,
	}
	switch keysym {
	case kbrd.XK_BackSpace:
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if ctrlActive && altActive {
			return errors.New("Error: quit event")
		}
	case kbrd.XK_t:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			if _, err := RunCommand(config.TerminalCmd()); err != nil {
				logging.Error("Terminal launch failed")
			}
		}
	case kbrd.XK_f:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendMaximize(manager.Curr())
		}
	case kbrd.XK_F1, kbrd.XK_F2, kbrd.XK_F3, kbrd.XK_F4, kbrd.XK_F5:
		fallthrough
	case kbrd.XK_F6, kbrd.XK_F7, kbrd.XK_F8, kbrd.XK_F9:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive && manager.Curr() != fkeys[keysym] {
			win := NewWindow(manager.Curr(), manager.Mailbox(), conn)
			win.SendReattach(fkeys[keysym])
			isSpecial := manager.Curr() == manager.SpecialWorkspace()
			if isSpecial && manager.Prev() == fkeys[keysym] {
				go func() {
					attachLock.Lock()
					win.SendActivate(manager.Prev())
					win.SendActivate(manager.Curr())
				}()
			}
		} else if !winActive && manager.Curr() != fkeys[keysym] {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			if manager.Curr() == manager.SpecialWorkspace() {
				win.SendDeactivate(manager.Prev())
			} else if fkeys[keysym] != manager.SpecialWorkspace() {
				win.SendDeactivate(manager.Curr())
			}
			win.SendActivate(fkeys[keysym])
			manager.SetCurr(fkeys[keysym])
		}
	case kbrd.XK_Left:
		winActive := (key.State & xproto.ModMask4) != 0
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && ctrlActive && !altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendResizeLeft(manager.Curr())
		}

		if winActive && !ctrlActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.SendFocusLeft(manager.Curr())
		}

		if winActive && !ctrlActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendMoveLeft()
		}

	case kbrd.XK_Right:
		winActive := (key.State & xproto.ModMask4) != 0
		ctrlActive := (key.State & xproto.ModMaskControl) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && ctrlActive && !altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendResizeRight(manager.Curr())
		}

		if winActive && !ctrlActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.SendFocusRight(manager.Curr())
		}

		if winActive && !ctrlActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendMoveRight()
		}
	case kbrd.XK_Up:
		winActive := (key.State & xproto.ModMask4) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.SendFocusUp(manager.Curr())
		}

		if winActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendMoveUp()
		}

	case kbrd.XK_Down:
		winActive := (key.State & xproto.ModMask4) != 0
		altActive := (key.State & xproto.ModMask1) != 0
		if winActive && !altActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.SendFocusDown(manager.Curr())
		}
		if winActive && altActive {
			win := NewWindow(uint32(key.Child), manager.Mailbox(), conn)
			win.SendMoveDown()
		}
	case kbrd.XK_q:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			win := NewWindow(uint32(key.Root), manager.Mailbox(), conn)
			win.SendClose(manager.Curr())
		}
	case kbrd.XK_grave:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			if _, err := RunCommand(config.LauncherCmd()); err != nil {
				logging.Error("Application launcher failed")
			}
		}
		return nil
	case kbrd.XK_l:
		winActive := (key.State & xproto.ModMask4) != 0
		if winActive {
			if _, err := RunCommand(config.LockerCmd()); err != nil {
				logging.Error("Locker launch failed")
			}
		}
	default:
		return nil
	}

	return nil
}

// RunCommand starts specified command in a separate goroutine
func RunCommand(c string) (*exec.Cmd, error) {
	args := regexp.MustCompile(" +").Split(c, -1)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Start()
	if err != nil {
		return cmd, err
	}
	go func() {
		cmd.Wait()
	}()
	return cmd, nil
}

func main() {
	config.ParseArgs()
	logging.Debug = config.Debug()

	conn, err := xgb.NewConn()
	if err != nil {
		logging.Fatal(err)
	}
	defer conn.Close()

	coninfo := xproto.Setup(conn)
	if coninfo == nil {
		logging.Fatal("Coudn't parse X connection info")
	}

	if len(coninfo.Roots) != 1 {
		logging.Fatal("Number of roots > 1, Xinerama init failed")
	}
	root := coninfo.Roots[0]
	cursor, err := xutil.CreateCursor(conn)
	if err != nil {
		logging.Fatal(err)
	}
	if err := xproto.ChangeWindowAttributesChecked(
		conn, root.Root, xproto.CwBackPixel|xproto.CwCursor,
		[]uint32{config.Color(), uint32(cursor)},
	).Check(); err != nil {
		logging.Fatal(err)
	}

	if err := xinerama.Init(conn); err != nil {
		logging.Fatal(err)
	}

	monitors, err := xutil.ReadMonitorsInfo(conn)
	if err != nil {
		logging.Fatal(err)
	} else if monitors.IsDualSetup() {
		xutil.SetNumberOfDesktops(MaxWorkspaces, conn)
	} else {
		xutil.SetNumberOfDesktops(MaxWorkspaces-1, conn)
	}

	if err := xutil.BecomeWM(conn, root); err != nil {
		logging.Println(err)
		logging.Fatal("Cannot take WM ownership")
	}

	keymap, err := kbrd.Mapping(conn)
	if err != nil {
		logging.Fatal(err)
	}

	if err := xutil.GrabMouse(conn, root); err != nil {
		logging.Fatal(err)
	}

	if err := xutil.GrabShortcuts(conn, root, keymap); err != nil {
		logging.Fatal(err)
	}

	xutil.SetSupported(conn) // Set EWMH supported atoms
	manager := NewWorkspaceManager(monitors)

	for _, cmd := range config.Commands() {
		c, _ := RunCommand(cmd)
		defer c.Process.Kill()
	}

	processEvents(conn, keymap, monitors, manager)
}
