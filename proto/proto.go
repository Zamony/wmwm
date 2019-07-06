package proto

import (
	"github.com/BurntSushi/xgb"
)

const (
	Attach = iota
	Detach
	Activate
	Deactivate
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
