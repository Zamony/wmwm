package main

import (
	"fmt"
	"log"

	"github.com/Zamony/wm/proto"
	"github.com/Zamony/wm/xutil"
)

const (
	LayoutFull = iota
	LayoutEqual
	LayoutPareto
)

const (
	MaxWorkspaces    = 9
	DefaultWorkspace = 1
	DefaultLayout    = LayoutEqual
)

type Workspace struct {
	left    *Column
	right   *Column
	central *Column
	input   chan proto.Message
	next    chan proto.Message
	headc   chan proto.Message
	id      uint32
	layout  int
	focus   *Window
}

func NewWorkspace(headc, input, next chan proto.Message, id uint32, screen xutil.Screen) *Workspace {
	return &Workspace{
		left:    NewColumn(screen),
		right:   NewColumn(screen),
		central: NewColumn(screen),
		input:   input,
		next:    next,
		headc:   headc,
		id:      id,
		layout:  LayoutFull,
		focus:   nil,
	}
}

func (workspace *Workspace) Run() {
	for {
		msg := <-workspace.input
		switch msg.To {
		case uint32(0):
			// Broadcast
			if msg.Type == proto.Exit {
				if workspace.next != nil {
					workspace.next <- msg
				}
				return
			}
			if workspace.FindWindow(msg.From) != nil {
				workspace.handleMsg(msg)
			} else if workspace.next != nil {
				workspace.next <- msg
			}
		case workspace.id:
			workspace.handleMsg(msg)
		default:
			if workspace.next != nil {
				workspace.next <- msg
			}
		}
	}
}

func (workspace *Workspace) handleMsg(msg proto.Message) {
	workspace.PrintStatus()
	switch msg.Type {
	case proto.Attach:
		//println("ATTACH", msg.From)
		win := NewWindow(workspace.id, workspace.headc, msg.XConn)
		go func() { win.Detach(msg.From) }()
	case proto.Reattach:
		//println("REATTACH", msg.From)
		win := NewWindow(msg.From, workspace.headc, msg.XConn)
		if workspace.FindWindow(msg.From) == nil {
			log.Println("=== INCOME", msg.From)
			workspace.Add(win)
			log.Println("=== ADD completed")
			workspace.Reshape()
			log.Println("=== RESHAPE completed")
		}
	case proto.Detach:
		//println("DETACH", msg.From, workspace.focus.Id())
		if workspace.focus != nil {
			win := workspace.focus
			go func() { win.Reattach(msg.From) }()
			workspace.Refocus()
			workspace.Remove(win)
			win.UnmapW()
			workspace.Reshape()
			workspace.focus.TakeFocus()
		}
	case proto.Remove:
		win := NewWindow(msg.From, workspace.headc, msg.XConn)
		if win.Id() == workspace.focus.Id() {
			workspace.Refocus()
		}
		workspace.Remove(win)
		workspace.Reshape()

	case proto.Close:
		win := workspace.focus
		if win != nil && !xutil.HasAtomDefined("WM_DELETE_WINDOW", win.Id(), msg.XConn) {
			workspace.Refocus()
			workspace.Remove(win)
			workspace.Reshape()
			win.DestroyW()
		} else if win != nil {
			win.CloseW()
		}
	case proto.FocusLeft:
		workspace.focus = workspace.FocusLeft()
		workspace.focus.TakeFocus()
	case proto.FocusRight:
		workspace.focus = workspace.FocusRight()
		println(workspace.focus.Id())
		workspace.focus.TakeFocus()
	case proto.FocusTop:
		workspace.focus = workspace.FocusTop()
		workspace.focus.TakeFocus()
	case proto.FocusBottom:
		workspace.focus = workspace.FocusBottom()
		workspace.focus.TakeFocus()
	case proto.Activate:
		workspace.Activate()
		workspace.focus.TakeFocus()
	case proto.Deactivate:
		if workspace.id != MaxWorkspaces {
			workspace.Deactivate()
		}
	case proto.ResizeLeft:
		// println("RESIZE LEFT", msg.From)
		workspace.ResizeLeft(msg.From)
	case proto.ResizeRight:
		// println("RESIZE RIGHT", msg.From)
		workspace.ResizeRight(msg.From)
	case proto.MoveUp:
		// println("MOVE UP", msg.From)
		workspace.MoveUp(workspace.focus.Id())
		workspace.Reshape()
		workspace.focus.TakeFocus()
	case proto.MoveDown:
		// println("MOVE DOWN", msg.From)
		workspace.MoveDown(workspace.focus.Id())
		workspace.Reshape()
		workspace.focus.TakeFocus()
	case proto.MoveLeft:
		workspace.MoveLeft(workspace.focus.Id())
		workspace.Reshape()
		workspace.focus.TakeFocus()
	case proto.MoveRight:
		workspace.MoveRight(workspace.focus.Id())
		workspace.Reshape()
		workspace.focus.TakeFocus()
	default:
		return
	}
	workspace.PrintStatus()
	println("------------------------")
}

func (workspace *Workspace) MoveLeft(wid uint32) {
	idx := workspace.right.IndexById(wid)
	if idx > -1 && workspace.right.Len() > 1 {
		win := workspace.right.WindowByIndex(idx)
		workspace.right.Remove(win)
		workspace.left.Add(win)
	}
}

func (workspace *Workspace) MoveRight(wid uint32) {
	idx := workspace.left.IndexById(wid)
	if idx > -1 && workspace.left.Len() > 1 {
		win := workspace.left.WindowByIndex(idx)
		workspace.left.Remove(win)
		workspace.right.Add(win)
	}
}

func (workspace *Workspace) MoveUp(wid uint32) {
	idx := workspace.central.IndexById(wid)
	if idx > -1 {
		return
	}

	idx = workspace.left.IndexById(wid)
	if idx > 0 {
		workspace.left.Swap(idx, idx-1)
	}

	idx = workspace.right.IndexById(wid)
	if idx > 0 {
		workspace.right.Swap(idx, idx-1)
	}
}

func (workspace *Workspace) MoveDown(wid uint32) {
	idx := workspace.central.IndexById(wid)
	if idx > -1 {
		return
	}

	idx = workspace.left.IndexById(wid)
	if idx > -1 && idx < workspace.left.Len()-1 {
		workspace.left.Swap(idx, idx+1)
	}

	idx = workspace.right.IndexById(wid)
	if idx > -1 && idx < workspace.right.Len()-1 {
		workspace.right.Swap(idx, idx+1)
	}
}

func (workspace *Workspace) ResizeLeft(wid uint32) {
	idx := workspace.right.IndexById(wid)
	// println("XXX", idx, "LAYOUT:", workspace.layout, "SETTINGTO:", LayoutPareto)
	if idx > -1 && workspace.layout == LayoutPareto {
		workspace.layout = LayoutEqual
		workspace.Reshape()
		return
	}

	idx = workspace.left.IndexById(wid)
	// println("XXX", idx, "LAYOUT:", workspace.layout, "SETTINGTO:", LayoutPareto)
	if idx > -1 && workspace.layout == LayoutPareto {
		workspace.layout = LayoutEqual
		workspace.Reshape()
	}
}

func (workspace *Workspace) ResizeRight(wid uint32) {
	idx := workspace.left.IndexById(wid)
	// println("XXX", idx, "LAYOUT:", workspace.layout, "SETTINGTO:", LayoutEqual)
	if idx > -1 && workspace.layout == LayoutEqual {
		workspace.layout = LayoutPareto
		workspace.Reshape()
		return
	}
	idx = workspace.right.IndexById(wid)
	// println("XXX", idx, "LAYOUT:", workspace.layout, "SETTINGTO:", LayoutEqual)
	if idx > -1 && workspace.layout == LayoutEqual {
		workspace.layout = LayoutPareto
		workspace.Reshape()
	}
}

func (workspace *Workspace) Add(window *Window) {
	leftEmpty := workspace.left.Len() < 1
	rightEmpty := workspace.right.Len() < 1
	centralEmpty := workspace.central.Len() < 1

	if leftEmpty && rightEmpty && centralEmpty {
		workspace.central.Add(window)
		workspace.focus = window
		workspace.layout = LayoutFull
		return
	}

	if workspace.layout == LayoutFull {
		workspace.layout = DefaultLayout
	}

	if !centralEmpty && leftEmpty && rightEmpty {
		win := workspace.central.WindowByIndex(0)
		workspace.left.Add(win)
		workspace.right.Add(window)
		workspace.central.Remove(win)
	} else {
		workspace.right.Add(window)
	}
}

func (workspace *Workspace) Remove(window *Window) {
	if window == nil {
		return
	}
	nleft := workspace.left.Len()
	nright := workspace.right.Len()

	inLeft, inRight, inCentral := false, false, false
	inCentral = workspace.central.IndexById(window.Id()) > -1
	if !inCentral {
		inLeft = workspace.left.IndexById(window.Id()) > -1
	}
	if !inCentral && !inLeft {
		inRight = workspace.right.IndexById(window.Id()) > -1
	}

	if inLeft && nleft == 1 && nright == 1 {
		win := workspace.right.WindowByIndex(0)
		workspace.right.Remove(win)
		workspace.central.Add(win)
		workspace.layout = LayoutFull
	}

	if inRight && nright == 1 && nleft == 1 {
		win := workspace.left.WindowByIndex(0)
		workspace.left.Remove(win)
		workspace.central.Add(win)
		workspace.layout = LayoutFull
	}

	if inLeft && nleft == 1 && nright > 1 {
		win := workspace.right.WindowByIndex(0)
		workspace.left.Add(win)
		workspace.right.Remove(win)
	}

	if inRight && nright == 1 && nleft > 1 {
		win := workspace.left.WindowByIndex(0)
		workspace.right.Add(win)
		workspace.left.Remove(win)
	}

	switch {
	case inCentral:
		workspace.central.Remove(window)
		workspace.focus = nil
		workspace.layout = LayoutFull
	case inRight:
		workspace.right.Remove(window)
	default:
		workspace.left.Remove(window)
	}
}

func (workspace *Workspace) Refocus() {
	if workspace.focus == nil {
		return
	}

	focus := workspace.FocusBottom()
	if focus != nil && focus.Id() != workspace.focus.Id() {
		workspace.focus = focus
		return
	}

	focus = workspace.FocusTop()
	if focus != nil && focus.Id() != workspace.focus.Id() {
		workspace.focus = focus
		return
	}

	focus = workspace.FocusLeft()
	if focus != nil && focus.Id() != workspace.focus.Id() {
		workspace.focus = focus
		return
	}

	focus = workspace.FocusRight()
	if focus != nil && focus.Id() != workspace.focus.Id() {
		workspace.focus = focus
		return
	}

	workspace.focus = nil
}

func (workspace *Workspace) FocusBottom() *Window {
	if workspace.focus == nil {
		return nil
	}

	win := workspace.central.WindowByIndex(0)
	if win != nil && win.Id() == workspace.focus.Id() {
		return workspace.focus
	}

	idx := workspace.left.IndexById(workspace.focus.Id())
	if idx > -1 && idx+1 < workspace.left.Len() {
		return workspace.left.WindowByIndex(idx + 1)
	}

	idx = workspace.right.IndexById(workspace.focus.Id())
	if idx > -1 && idx+1 < workspace.right.Len() {
		return workspace.right.WindowByIndex(idx + 1)
	}

	return workspace.focus
}

func (workspace *Workspace) FocusTop() *Window {
	if workspace.focus == nil {
		return nil
	}

	win := workspace.central.WindowByIndex(0)
	if win != nil && win.Id() == workspace.focus.Id() {
		return workspace.focus
	}

	idx := workspace.left.IndexById(workspace.focus.Id())
	if idx > -1 && idx-1 >= 0 {
		return workspace.left.WindowByIndex(idx - 1)
	}

	idx = workspace.right.IndexById(workspace.focus.Id())
	if idx > -1 && idx-1 >= 0 {
		return workspace.right.WindowByIndex(idx - 1)
	}

	return workspace.focus
}

func (workspace *Workspace) FocusLeft() *Window {
	if workspace.focus == nil {
		return nil
	}

	win := workspace.central.WindowByIndex(0)
	if win != nil && win.Id() == workspace.focus.Id() {
		return workspace.focus
	}

	idx := workspace.left.IndexById(workspace.focus.Id())
	if idx > -1 {
		return workspace.focus
	}

	return workspace.left.WindowByIndex(0)
}

func (workspace *Workspace) FocusRight() *Window {
	if workspace.focus == nil {
		return nil
	}

	win := workspace.central.WindowByIndex(0)
	if win != nil && win.Id() == workspace.focus.Id() {
		return workspace.focus
	}

	idx := workspace.right.IndexById(workspace.focus.Id())
	if idx > -1 {
		return workspace.focus
	}

	return workspace.right.WindowByIndex(0)
}

func (workspace *Workspace) PrintStatus() {
	fmt.Print("Workspace ID", workspace.id, " ")
	if workspace.focus != nil {
		fmt.Println("focus = ", workspace.focus.Id())
	} else {
		fmt.Println("focus = nil")
	}

	fmt.Printf("Central ")
	workspace.central.PrintStatus()
	fmt.Printf("Left ")
	workspace.left.PrintStatus()
	fmt.Printf("Right ")
	workspace.right.PrintStatus()
	fmt.Printf("\n\n")
}

func (workspace *Workspace) Activate() {
	for i := 0; i < workspace.central.Len(); i++ {
		win := workspace.central.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.MapW()
		}
	}
	for i := 0; i < workspace.left.Len(); i++ {
		win := workspace.left.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.MapW()
		}
	}
	for i := 0; i < workspace.right.Len(); i++ {
		win := workspace.right.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.MapW()
		}
	}
}

func (workspace *Workspace) Deactivate() {
	for i := 0; i < workspace.central.Len(); i++ {
		win := workspace.central.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.UnmapW()
		}
	}
	for i := 0; i < workspace.left.Len(); i++ {
		win := workspace.left.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.UnmapW()
		}
	}
	for i := 0; i < workspace.right.Len(); i++ {
		win := workspace.right.WindowByIndex(i)
		if win == nil {
			log.Fatal("win = nil")
		} else {
			win.UnmapW()
		}
	}
}

func (workspace *Workspace) FindWindow(wid uint32) *Window {
	if idx := workspace.central.IndexById(wid); idx > -1 {
		return workspace.central.WindowByIndex(idx)
	}

	if idx := workspace.left.IndexById(wid); idx > -1 {
		return workspace.left.WindowByIndex(idx)
	}

	if idx := workspace.right.IndexById(wid); idx > -1 {
		return workspace.right.WindowByIndex(idx)
	}

	return nil
}

func (workspace *Workspace) Reshape() {
	if workspace.central.Len() > 0 {
		workspace.central.SetX(0)
		workspace.central.SetWidth100()
	} else if workspace.layout == LayoutEqual {
		workspace.left.SetX(0)
		x := workspace.left.SetWidth50()
		workspace.right.SetX(x)
		workspace.right.SetWidth50()
	} else {
		workspace.left.SetX(0)
		x := workspace.left.SetWidth80()
		workspace.right.SetX(x)
		workspace.right.SetWidth20()
	}

	workspace.central.Reshape()
	workspace.left.Reshape()
	workspace.right.Reshape()
}

type WorkspaceManager struct {
	curr    uint32
	mailbox chan proto.Message
}

func NewWorkspaceManager(mainscr, auxscr xutil.Screen) *WorkspaceManager {
	input := make(chan proto.Message)
	next := make(chan proto.Message)
	mailbox := input

	for j := uint32(1); j < MaxWorkspaces-1; j++ {
		w := NewWorkspace(mailbox, input, next, j, mainscr)
		go w.Run()
		input = next
		next = make(chan proto.Message)
	}

	if mainscr.Id() == auxscr.Id() {
		w := NewWorkspace(mailbox, input, nil, MaxWorkspaces, mainscr)
		go w.Run()
	} else {
		w := NewWorkspace(mailbox, input, next, MaxWorkspaces, mainscr)
		go w.Run()
		w = NewWorkspace(mailbox, next, nil, MaxWorkspaces, auxscr)
		go w.Run()
	}

	return &WorkspaceManager{DefaultWorkspace, mailbox}
}

func (wrkmgr *WorkspaceManager) Mailbox() chan proto.Message {
	return wrkmgr.mailbox
}

func (wrkmgr *WorkspaceManager) Curr() uint32 {
	return wrkmgr.curr
}

func (wrkmgr *WorkspaceManager) SetCurr(n uint32) {
	wrkmgr.curr = n
}
