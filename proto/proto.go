// Package proto defines internal protocol messages structure
package proto

import (
	"github.com/BurntSushi/xgb"
)

// Types of the messages used in the internal protocol
const (
	Attach = iota
	Detach
	Reattach
	Activate
	Deactivate
	Remove
	MoveLeft
	MoveRight
	MoveUp
	MoveDown
	FocusHere
	FocusLeft
	FocusRight
	FocusUp
	FocusDown
	Maximize
	ResizeLeft
	ResizeRight
	Close
	Exit
)

// Message represents message of the internal protocol
type Message struct {
	From  uint32
	To    uint32
	Type  uint
	XConn *xgb.Conn
}
