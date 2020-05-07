// Package main implements logic of the window manager
package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wmwm/config"
	"github.com/Zamony/wmwm/logging"
	"github.com/Zamony/wmwm/proto"
	"github.com/Zamony/wmwm/xutil"
)

// Window is the basic structure representing X window
type Window struct {
	height         int
	width          int
	x              int
	y              int
	mailbox        chan proto.Message
	id             uint32
	conn           *xgb.Conn
	removalAllowed bool
}

// NewWindow creates instance of Window
func NewWindow(id uint32, c chan proto.Message, xc *xgb.Conn) *Window {
	return &Window{0, 0, 0, 0, c, id, xc, true}
}

// Id returns identifier of window
func (window *Window) Id() uint32 {
	return window.id
}

// SendAttach sends attach request to the specified workspace
func (window *Window) SendAttach(to uint32) {
	msg := proto.Message{window.id, to, proto.Attach, window.conn}
	window.mailbox <- msg
}

// SendDetach sends detach request to the specified workspace
func (window *Window) SendDetach(to uint32) {
	msg := proto.Message{window.id, to, proto.Detach, window.conn}
	window.mailbox <- msg
}

// SendReattach sends reattach request to the specified workspace
func (window *Window) SendReattach(to uint32) {
	msg := proto.Message{window.id, to, proto.Reattach, window.conn}
	window.mailbox <- msg
}

// SendDeactivate sends request to deactivate specified workspace
func (window *Window) SendDeactivate(to uint32) {
	msg := proto.Message{window.id, to, proto.Deactivate, window.conn}
	window.mailbox <- msg
}

// SendActivate sends request to activate specified workspace
func (window *Window) SendActivate(id uint32) {
	msg := proto.Message{window.id, id, proto.Activate, window.conn}
	window.mailbox <- msg
}

// SendRemove sends "remove me" request to workspace which it belongs to
func (window *Window) SendRemove() {
	msg := proto.Message{window.id, 0, proto.Remove, window.conn}
	window.mailbox <- msg
}

// SendMoveLeft sends "move me to the left" request
// to workspace which it belongs to
func (window *Window) SendMoveLeft() {
	msg := proto.Message{window.id, 0, proto.MoveLeft, window.conn}
	window.mailbox <- msg
}

// SendMoveRight sends "move me to the right" request
// to workspace which it belongs to
func (window *Window) SendMoveRight() {
	msg := proto.Message{window.id, 0, proto.MoveRight, window.conn}
	window.mailbox <- msg
}

// SendMoveUp sends "move me to the up" request
// to workspace which it belongs to
func (window *Window) SendMoveUp() {
	msg := proto.Message{window.id, 0, proto.MoveUp, window.conn}
	window.mailbox <- msg
}

// SendMoveDown sends "move me to the down" request
// to workspace which it belongs to
func (window *Window) SendMoveDown() {
	msg := proto.Message{window.id, 0, proto.MoveDown, window.conn}
	window.mailbox <- msg
}

// SendReattach sends close request to the specified workspace
func (window *Window) SendClose(id uint32) {
	msg := proto.Message{window.id, id, proto.Close, window.conn}
	window.mailbox <- msg
}

// SendExit broadcasts exit message
func (window *Window) SendExit() {
	msg := proto.Message{window.id, 0, proto.Exit, window.conn}
	window.mailbox <- msg
}

// SendFocusHere sends "focus on me" request
// to workspace which it belongs to
func (window *Window) SendFocusHere() {
	msg := proto.Message{window.id, 0, proto.FocusHere, window.conn}
	window.mailbox <- msg
}

// SendFocusLeft sends "focus on the window on the left from the current focus"
// request to workspace which it belongs to
func (window *Window) SendFocusLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusLeft, window.conn}
	window.mailbox <- msg
}

// SendFocusRight sends "focus on the window on the right from the current focus"
// request to workspace which it belongs to
func (window *Window) SendFocusRight(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusRight, window.conn}
	window.mailbox <- msg
}

// SendFocusUp sends "focus on the window which is above current focus"
// request to workspace which it belongs to
func (window *Window) SendFocusUp(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusUp, window.conn}
	window.mailbox <- msg
}

// SendFocusDown sends "focus on the window which is under current focus"
// request to workspace which it belongs to
func (window *Window) SendFocusDown(id uint32) {
	msg := proto.Message{window.id, id, proto.FocusDown, window.conn}
	window.mailbox <- msg
}

// SendMaximize sends request to the specified workspace,
// which makes central column full in size
func (window *Window) SendMaximize(id uint32) {
	msg := proto.Message{window.id, id, proto.Maximize, window.conn}
	window.mailbox <- msg
}

// SendResizeLeft sends request to resize current window to the left
func (window *Window) SendResizeLeft(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeLeft, window.conn}
	window.mailbox <- msg
}

// SendResizeRight sends request to resize current window to the right
func (window *Window) SendResizeRight(id uint32) {
	msg := proto.Message{window.id, id, proto.ResizeRight, window.conn}
	window.mailbox <- msg
}

// SetX sets window's x-coordinate value
func (window *Window) SetX(x int) error {
	window.x = x
	return xutil.SetWindowX(x, window.id, window.conn)
}

// SetY sets window's y-coordinate value
func (window *Window) SetY(y int) error {
	window.y = y
	return xutil.SetWindowY(y, window.id, window.conn)
}

// SetWidth sets width of the window
func (window *Window) SetWidth(w int) error {
	window.width = w
	return xutil.SetWindowWidth(w, window.id, window.conn)
}

// SetHeight sets height of the window
func (window *Window) SetHeight(h int) error {
	window.height = h
	return xutil.SetWindowHeight(h, window.id, window.conn)
}

// Map makes window visible on the screen
func (window *Window) Map() error {
	if err := xutil.MapWindow(window.id, window.conn); err != nil {
		return err
	}

	if err := xutil.RemoveWindowBorder(window.id, window.conn); err != nil {
		return err
	}

	return xutil.WatchWindowEvents(window.id, window.conn)
}

// Unmap hides the window
func (window *Window) Unmap() error {
	return xutil.UnmapWindow(window.id, window.conn)
}

// Close closes the window
func (window *Window) Close() error {
	return xutil.SendClientEvent(
		"WM_DELETE_WINDOW", xproto.TimeCurrentTime, window.id, window.conn,
	)
}

// Destroy kills the window
func (window *Window) Destroy() error {
	return xutil.DestroyWindow(window.id, window.conn)
}

// CouldBeDestroyed checks whether window could be destroyed
func (window *Window) CouldBeDestroyed() bool {
	return !xutil.HasAtomDefined("WM_DELETE_WINDOW", window.id, window.conn)
}

// Defocus makes window appear like the unfocused one
func (window *Window) Defocus() error {
	if window == nil {
		return nil
	}

	return window.UnsetBorder()
}

// UnsetBorder removes padding to make the window
// looks like it doesn't has border
func (window *Window) UnsetBorder() error {
	return xutil.RemovePaddingFromWindow(
		window.y, window.height, config.BorderWidth(),
		window.id, window.conn,
	)
}

// SetBorder adds padding to make the window looks like it has border
func (window *Window) SetBorder() error {
	return xutil.AddPaddingToWindow(
		window.y, window.height, config.BorderWidth(),
		window.id, window.conn,
	)
}

// TakeFocus makes window manager take focus on the window
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

	return xutil.FocusWindow(window.id, window.conn)
}

// DenyRemoval sets a flag that will forbid window removal
func (window *Window) DenyRemoval() {
	window.removalAllowed = false
}

// AllowRemoval sets a flag that permits window removal
func (window *Window) AllowRemoval() {
	window.removalAllowed = true
}

// IsRemovalAllowed checks whether it is allowed to remove the window
func (window Window) IsRemovalAllowed() bool {
	return window.removalAllowed
}

// IsDock performs check whether this window is dock or not
func (window Window) IsDock() bool {
	return xutil.IsDock(window.Id(), window.conn)
}

// LogStatus logs window's information for debugging purposes
func (window Window) LogStatus() {
	logging.Println(
		"X:", window.x, "Y:", window.y,
		"W:", window.width, "H:", window.height, "ID:", window.id,
	)
}
