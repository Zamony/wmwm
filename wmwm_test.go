package main

import (
	"testing"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/Zamony/wmwm/proto"
	"github.com/Zamony/wmwm/xutil"
)

func TestColumnIndexById(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	var conn interface{}
	xconn, _ := conn.(*xgb.Conn)
	w1 := NewWindow(1, ch, xconn)
	w2 := NewWindow(2, ch, xconn)
	w3 := NewWindow(3, ch, xconn)
	c.Add(w1)
	c.Add(w2)
	c.Add(w3)
	if i := c.IndexById(2); i != 1 {
		t.Error("Window2 index != 1")
	}
}

func TestColumnRemoveExistent(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	var conn interface{}
	xconn, _ := conn.(*xgb.Conn)
	w1 := NewWindow(1, ch, xconn)
	w2 := NewWindow(2, ch, xconn)
	w3 := NewWindow(3, ch, xconn)
	c.Add(w1)
	c.Add(w2)
	c.Add(w3)
	c.Remove(w2)
	if c.Len() != 2 {
		t.Error("Column length != 2")
	}
	if c.IndexById(1) < 0 || c.IndexById(3) < 0 {
		t.Error("Removed wrong window")
	}
}

func TestColumnRemoveNonExistent(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	var conn interface{}
	xconn, _ := conn.(*xgb.Conn)
	w1 := NewWindow(1, ch, xconn)
	w2 := NewWindow(2, ch, xconn)
	c.Add(w1)
	if rem := c.Remove(w2); rem != nil {
		t.Error("Failed to remove non-existent window")
	}
}

func TestColumnWindowByIndex(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	var conn interface{}
	xconn, _ := conn.(*xgb.Conn)
	w1 := NewWindow(1, ch, xconn)
	w2 := NewWindow(2, ch, xconn)
	c.Add(w1)
	c.Add(w2)
	if w := c.WindowByIndex(3); w != nil {
		t.Error("Found non-existent window")
	}
	if w := c.WindowByIndex(1); w == nil {
		t.Error("Existent window not found at index 1")
	}
}

func TestColumnSwap(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	var conn interface{}
	xconn, _ := conn.(*xgb.Conn)
	w1 := NewWindow(1, ch, xconn)
	w2 := NewWindow(2, ch, xconn)
	c.Add(w1)
	c.Add(w2)
	if err := c.Swap(0, 2); err == nil {
		t.Error("Swapping non existent column with index 2")
	}

	if err := c.Swap(-1, 1); err == nil {
		t.Error("Swapping non existent column with index -1")
	}

	if err := c.Swap(0, 1); err != nil {
		t.Error("Swapping failed")
	}

	w2 = c.WindowByIndex(0)
	w1 = c.WindowByIndex(1)
	if w1 == nil || w2 == nil {
		t.Error("Swapping caused one of the windows to disappear")
	}

	if w1.Id() != 1 || w2.Id() != 2 {
		t.Error("Swapped windows not in the specified positions")
	}
}

func TestColumnReshape(t *testing.T) {
	screen := xutil.NewScreen(4, 3, 0, 0, 0)
	c := NewColumn(screen)
	ch := make(chan proto.Message)
	w1 := NewWindow(1, ch, nil)
	w2 := NewWindow(2, ch, nil)
	w3 := NewWindow(3, ch, nil)
	c.Add(w1)
	c.Add(w2)
	c.Add(w3)
	oldConfWC := xutil.ConfigureWindowChecked
	cookie := &xgb.Cookie{}
	xutil.ConfigureWindowChecked = func(
		c *xgb.Conn, window xproto.Window, ValueMask uint16, ValueList []uint32,
	) xproto.ConfigureWindowCookie {
		return xproto.ConfigureWindowCookie{cookie}
	}
	c.Reshape()
	if w1.y != 0 || w1.height != 1 {
		t.Error("W1 has invalid geometry", w1.y, w1.height)
	}
	if w2.y != 1 || w2.height != 1 {
		t.Error("W2 has invalid geometry", w2.y, w2.height)
	}
	if w3.y != 2 || w3.height != 1 {
		t.Error("W3 has invalid geometry", w3.y, w3.height)
	}

	xutil.ConfigureWindowChecked = oldConfWC
}

func TestWorkspaceAdd(t *testing.T) {
	c := make(chan proto.Message)
	screen := xutil.NewScreen(8, 6, 0, 0, 0)
	wr := NewWorkspace(c, c, nil, 1, screen)
	w1 := NewWindow(1, c, nil)
	wr.Add(w1)
	if wr.central.Len() != 1 {
		t.Error("Win1: wr.central.Len() != 1")
	}
	w2 := NewWindow(2, c, nil)
	wr.Add(w2)
	switch {
	case wr.central.Len() != 0:
		t.Error("Win2: wr.central.Len() != 0")
	case wr.left.Len() != 1:
		t.Error("Win2: wr.left.Len() != 1")
	case wr.right.Len() != 1:
		t.Error("Win2: wr.right.Len() != 1")
	}

	w3 := NewWindow(3, c, nil)
	wr.Add(w3)
	switch {
	case wr.central.Len() != 0:
		t.Error("Win3: wr.central.Len() != 0")
	case wr.left.Len() != 1:
		t.Error("Win3: wr.left.Len() != 1")
	case wr.right.Len() != 2:
		t.Error("Win3: wr.right.Len() != 2")
	}
}

func TestWorkspaceRemoveCentral(t *testing.T) {
	c := make(chan proto.Message)
	screen := xutil.NewScreen(8, 6, 0, 0, 0)
	wr := NewWorkspace(c, c, nil, 1, screen)
	w1 := NewWindow(1, c, nil)
	wr.Add(w1)
	wr.Remove(w1)
	if wr.central.Len() != 0 {
		t.Error("wr.central.Len() != 0")
	}
}

func TestWorkspaceRemoveLastInColumn(t *testing.T) {
	c := make(chan proto.Message)
	screen := xutil.NewScreen(8, 6, 0, 0, 0)
	wr := NewWorkspace(c, c, nil, 1, screen)
	w1 := NewWindow(1, c, nil)
	w2 := NewWindow(2, c, nil)
	wr.Add(w1)
	wr.Add(w2)
	wr.Remove(w1)
	if wr.central.Len() != 1 {
		t.Error("wr.central.Len() != 0")
	}
	wr.Add(w1)
	wr.Remove(w1)
	if wr.central.Len() != 1 {
		t.Error("wr.central.Len() != 0")
	}
}

func TestWorkspaceRemove(t *testing.T) {
	c := make(chan proto.Message)
	screen := xutil.NewScreen(8, 6, 0, 0, 0)
	wr := NewWorkspace(c, c, nil, 1, screen)
	w1 := NewWindow(1, c, nil)
	w2 := NewWindow(2, c, nil)
	w3 := NewWindow(3, c, nil)
	wr.Add(w1)
	wr.Add(w2)
	wr.Add(w3)
	wr.Remove(w2)
	if wr.central.Len() != 0 {
		t.Error("wr.central.Len() != 0")
	}
	if wr.left.Len() != 1 {
		t.Error("wr.left.Len() != 1")
	}
	if wr.right.Len() != 1 {
		t.Error("wr.right.Len() != 1")
	}
}
