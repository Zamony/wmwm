package proto

import (
	"github.com/BurntSushi/xgb"
)

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
	FocusLeft
	FocusRight
	FocusTop
	FocusBottom
	SetLayoutFull
	SetLayoutEqual
	SetLayoutPareto
	ResizeLeft
	ResizeRight
	Close
	Exit
)

type Message struct {
	From  uint32
	To    uint32
	Type  uint
	XConn *xgb.Conn
}
