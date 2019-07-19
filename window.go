package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wm/logging"
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

func (window *Window) SendAttach(to uint32) {
	msg := proto.Message{window.id, to, proto.Attach, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendDetach(to uint32) {
	msg := proto.Message{window.id, to, proto.Detach, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendReattach(to uint32) {
	msg := proto.Message{window.id, to, proto.Reattach, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendDeactivate(to uint32) {
	msg := proto.Message{window.id, to, proto.Deactivate, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendActivate(id uint32) {
	msg := proto.Message{window.id, id, proto.Activate, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendRemove() {
	msg := proto.Message{window.id, 0, proto.Remove, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendMoveLeft() {
	msg := proto.Message{window.id, 0, proto.MoveLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendMoveRight() {
	msg := proto.Message{window.id, 0, proto.MoveRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendMoveUp() {
	msg := proto.Message{window.id, 0, proto.MoveUp, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendMoveDown() {
	msg := proto.Message{window.id, 0, proto.MoveDown, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendClose(id uint32) {
	msg := proto.Message{window.id, id, proto.Close, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendExit() {
	msg := proto.Message{window.id, 0, proto.Exit, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendFocusHere() {
	msg := proto.Message{window.id, 0, proto.FocusHere, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendFocusLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendFocusRight(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendFocusUp(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusUp, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendFocusDown(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusDown, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendResizeLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeLeft, window.conn}
	window.mailbox <- msg
}

func (window *Window) SendResizeRight(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeRight, window.conn}
	window.mailbox <- msg
}

func (window *Window) SetX(x int) error {
	window.x = x
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowX, []uint32{uint32(x)},
	).Check()
	return err
}

func (window *Window) SetY(y int) error {
	window.y = y
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowY, []uint32{uint32(y)},
	).Check()
	return err
}

func (window *Window) SetWidth(w int) error {
	window.width = w
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowWidth, []uint32{uint32(w)},
	).Check()
	return err
}

func (window *Window) SetHeight(h int) error {
	window.height = h
	err := xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowHeight, []uint32{uint32(h)},
	).Check()
	return err
}

func (window *Window) Map() error {
	err := xproto.MapWindowChecked(
		window.conn, xproto.Window(window.id),
	).Check()
	if err != nil {
		return err
	}

	err = xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowBorderWidth, []uint32{uint32(0)},
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

func (window *Window) Unmap() error {
	err := xproto.UnmapWindowChecked(
		window.conn, xproto.Window(window.id),
	).Check()
	return err
}

func (window *Window) Close() error {
	return xutil.SendClientEvent(
		"WM_DELETE_WINDOW",
		xproto.TimeCurrentTime,
		window.id,
		window.conn,
	)
}

func (window *Window) Destroy() error {
	return xproto.DestroyWindowChecked(window.conn, xproto.Window(window.id)).Check()
}

func (window *Window) CouldBeDestroyed() bool {
	return !xutil.HasAtomDefined("WM_DELETE_WINDOW", window.id, window.conn)
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

	return window.UnsetBorder()
}

func (window *Window) UnsetBorder() error {
	return xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{uint32(window.y - 2), uint32(window.height + 2)},
	).Check()
}

func (window *Window) SetBorder() error {
	return xproto.ConfigureWindowChecked(
		window.conn, xproto.Window(window.id),
		xproto.ConfigWindowY|xproto.ConfigWindowHeight,
		[]uint32{
			uint32(window.y + 2),
			uint32(window.height - 2),
		},
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

	if err := window.SetBorder(); err != nil {
		return err
	}

	return xproto.SetInputFocusChecked(
		window.conn, xproto.InputFocusPointerRoot,
		xproto.Window(window.id), xproto.TimeCurrentTime,
	).Check()
}

func (window Window) IsDock() bool {
	return xutil.IsDock(window.Id(), window.conn)
}

func (window Window) LogStatus() {
	logging.Println(
		"X:", window.x, "Y:", window.y,
		"W:", window.width, "H:", window.height, "ID:", window.id,
	)
}
