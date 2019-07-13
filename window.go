package main

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/proto"
	"github.com/Zamony/wm/xutil"
)

type Window struct {
	height  int
	width   int
	x       int
	y       int
	mailbox chan proto.Message
	id      uint32
	conn    *xgb.Conn
}

func NewWindow(id uint32, c chan proto.Message, xc *xgb.Conn) *Window {
	return &Window{0, 0, 0, 0, c, id, xc}
}

func (window *Window) Id() uint32 {
	return window.id
}

func (window *Window) Attach(to uint32) {
	msg := proto.Message{window.id, to, proto.Attach, window.conn}
	window.mailbox <- msg
}

func (window *Window) Detach(to uint32) {
	msg := proto.Message{window.id, to, proto.Detach, window.conn}
	window.mailbox <- msg
}

func (window *Window) Reattach(to uint32) {
	msg := proto.Message{window.id, to, proto.Reattach, window.conn}
	window.mailbox <- msg
}

func (window *Window) Deactivate(to uint32) {
	msg := proto.Message{window.id, to, proto.Deactivate, window.conn}
	window.mailbox <- msg
}

func (window *Window) Activate(id uint32) {
	msg := proto.Message{window.id, id, proto.Activate, window.conn}
	window.mailbox <- msg
}

func (window *Window) Remove() {
	msg := proto.Message{window.id, 0, proto.Remove, window.conn}
	window.mailbox <- msg
}

func (window *Window) MoveLeft() {
	msg := proto.Message{window.id, 0, proto.MoveLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) MoveRight() {
	msg := proto.Message{window.id, 0, proto.MoveRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) MoveUp() {
	msg := proto.Message{window.id, 0, proto.MoveUp, window.conn}
	window.mailbox <- msg
}

func (window *Window) MoveDown() {
	msg := proto.Message{window.id, 0, proto.MoveDown, window.conn}
	window.mailbox <- msg
}

func (window *Window) Close(id uint32) {
	msg := proto.Message{window.id, id, proto.Close, window.conn}
	window.mailbox <- msg
}

func (window *Window) Exit() {
	msg := proto.Message{window.id, 0, proto.Exit, window.conn}
	window.mailbox <- msg
}

func (window *Window) FocusHere() {
	msg := proto.Message{window.id, 0, proto.FocusHere, window.conn}
	window.mailbox <- msg
}

func (window *Window) FocusLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) FocusRight(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) FocusTop(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusTop, window.conn}
	window.mailbox <- msg
}

func (window *Window) FocusBottom(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusBottom, window.conn}
	window.mailbox <- msg
}

func (window *Window) ResizeLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) ResizeRight(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) SetX(x int) error {
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowX, []uint32{uint32(x)},
	).Check()
	window.x = x
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (window *Window) SetY(y int) error {
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowY, []uint32{uint32(y)},
	).Check()
	window.y = y
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (window *Window) SetWidth(w int) error {
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowWidth, []uint32{uint32(w)},
	).Check()
	window.width = w
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (window *Window) SetHeight(h int) error {
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowHeight, []uint32{uint32(h)},
	).Check()
	window.height = h
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (window *Window) MapW() error {
	err := xproto.MapWindowChecked(
		window.conn, xproto.Window(window.id),
	).Check()
	if err != nil {
		return err
	}
	err = xproto.ChangeWindowAttributesChecked(
		window.conn, xproto.Window(window.id),
		xproto.CwEventMask, []uint32{
			xproto.EventMaskStructureNotify | xproto.EventMaskEnterWindow,
		},
	).Check()
	return err
}

func (window *Window) UnmapW() error {
	err := xproto.UnmapWindowChecked(
		window.conn, xproto.Window(window.id),
	).Check()
	return err
}

func (window *Window) CloseW() error {
	// if !xutil.HasAtomDefined("WM_DELETE_WINDOW", window.id, window.conn) {
	// 	err :=
	// 	return err
	// }

	return xutil.SendClientEvent(
		"WM_DELETE_WINDOW",
		xproto.TimeCurrentTime,
		window.id,
		window.conn,
	)
}

func (window *Window) DestroyW() error {
	return xproto.DestroyWindowChecked(window.conn, xproto.Window(window.id)).Check()
}

func (window *Window) CouldBeManaged() bool {
	x := window.x
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowX, []uint32{uint32(x)},
	).Check()

	if err != nil {
		return false
	}

	return true
}

func (window *Window) Defocus() error {
	if window == nil {
		return nil
	}

	return xproto.ChangeWindowAttributesChecked(
		window.conn, xproto.Window(window.id),
		xproto.CwBorderPixel, []uint32{0x000000},
	).Check()
}

func (window *Window) TakeFocus() error {
	if window == nil {
		return nil
	}
	timepoint := uint32(xproto.TimeCurrentTime)
	if timepoint != 0 {
		timepoint--
	}
	if xutil.HasAtomDefined("WM_TAKE_FOCUS", window.id, window.conn) {
		println("HAS WM_TAKE_FOCUS")
		err := xutil.SendClientEvent(
			"WM_TAKE_FOCUS",
			timepoint,
			window.id,
			window.conn,
		)
		if err != nil {
			return err
		}
	}

	err := xproto.ChangeWindowAttributesChecked(
		window.conn, xproto.Window(window.id),
		xproto.CwBorderPixel, []uint32{0x00ff00},
	).Check()
	if err != nil {
		return err
	}

	err = xproto.SetInputFocusChecked(
		window.conn, xproto.InputFocusPointerRoot,
		xproto.Window(window.id), xproto.TimeCurrentTime,
	).Check()
	return err
}

func (window Window) PrintStatus() {
	println("X:", window.x, "Y:", window.y, "W:", window.width, "H:", window.height, "ID:", window.id)
}
