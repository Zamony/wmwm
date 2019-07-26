// Package main implements logic of the window manager
package main

import (
    "fmt"
    "sync"
    "github.com/BurntSushi/xgb"
    "github.com/Zamony/wm/proto"
    "github.com/Zamony/wm/xutil"
    "github.com/Zamony/wm/config"
    "github.com/Zamony/wm/logging"
)

// Layouts affect width of the columns
const (
    LayoutFull = iota
    LayoutEqual
    LayoutLeftWide
)

const (
    // MaxWorkspaces sets the number of workspaces available
    MaxWorkspaces    = 9
    // DefaultWorkspace sets default active workspace
    DefaultWorkspace = 1
    // DefaultLayout sets default column layout
    DefaultLayout    = LayoutEqual
)

var (
    unmapLock  ReattachLock
    attachLock ReattachLock
)

// ReattachLock represents lock that can be harmlessly unlocked
// even if there are not any locked goroutines
type ReattachLock struct {
    wg sync.WaitGroup
    locked bool
}

// Lock makes the caller wait till unlock happens
func (m *ReattachLock) Lock() {
    m.wg.Add(1)
    m.locked = true
    m.wg.Wait()
}

// Unlock unblocks all awaiting goroutines if there are any
func (m *ReattachLock) Unlock() {
    if m.locked {
        m.locked = false
        m.wg.Done()
    }
}

// Workspace represents a group of related windows
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
    conn    *xgb.Conn
}

// NewWorkspace creates instance of Workspace
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
        conn:    nil,
    }
}

// Run makes workspace listen for incoming messages from windows
func (workspace *Workspace) Run() {
    for {
        msg := <-workspace.input
        if workspace.conn == nil {
            workspace.conn = msg.XConn
        }
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
    workspace.LogStatus()
    switch msg.Type {
    case proto.Reattach:
        win := NewWindow(workspace.id, workspace.headc, msg.XConn)
        go func() { win.SendDetach(msg.From) }()
    case proto.Attach:
        win := NewWindow(msg.From, workspace.headc, msg.XConn)
        if win.IsDock() {
            win.Map()
            break
        }
        if workspace.FindWindow(msg.From) == nil {
            workspace.Add(win)
            workspace.Reshape()
            if workspace.id == MaxWorkspaces {
                workspace.Activate()
            }
        }
    case proto.Detach:
        if workspace.focus != nil {
            win := workspace.focus
            win.Unmap()
            go func() {
                unmapLock.Lock()
                win.SendAttach(msg.From)
                attachLock.Unlock()
            }()
        }
    case proto.Remove:
        win := workspace.FindWindow(msg.From)
        if win != nil && workspace.focus != nil {
            if !win.IsRemovalAllowed() {
                win.AllowRemoval()
                break
            }
            if win.Id() == workspace.focus.Id() {
                win.Defocus()
                workspace.Refocus()
            }
            workspace.Remove(win)
            workspace.Reshape()
            workspace.Focus()
            unmapLock.Unlock()
        }

    case proto.Close:
        win := workspace.focus
        if win == nil {
            break
        }
        if win.CouldBeDestroyed() {
            win.Defocus()
            workspace.Refocus()
            workspace.Remove(win)
            workspace.Reshape()
            win.Destroy()
        } else {
            win.Close()
        }
    case proto.FocusHere:
        if workspace.focus.Id() != msg.From {
            win := workspace.FindWindow(msg.From)
            if win != nil {
                workspace.focus.Defocus()
                workspace.focus = win
                workspace.Focus()
            }
        }
    case proto.FocusLeft:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusLeft()
        workspace.Focus()
    case proto.FocusRight:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusRight()
        workspace.Focus()
    case proto.FocusUp:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusUp()
        workspace.Focus()
    case proto.FocusDown:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusDown()
        workspace.Focus()
    case proto.Maximize:
        if workspace.central.HasPadding() {
            workspace.central.RemovePadding()
        } else {
            workspace.central.AddPadding()
        }
        workspace.Activate()
    case proto.Activate:
        workspace.Reshape()
        workspace.Activate()
        workspace.Focus()
        xutil.SetCurrentDesktop(workspace.id, msg.XConn)
    case proto.Deactivate:
        if workspace.id != MaxWorkspaces {
            workspace.Deactivate()
        }
    case proto.ResizeLeft:
        workspace.ResizeLeft(msg.From)
        workspace.Focus()
    case proto.ResizeRight:
        workspace.ResizeRight(msg.From)
        workspace.Focus()
    case proto.MoveUp:
        workspace.MoveUp(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()
    case proto.MoveDown:
        workspace.MoveDown(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()
    case proto.MoveLeft:
        workspace.MoveLeft(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()
    case proto.MoveRight:
        workspace.MoveRight(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()
    default:
        return
    }

    workspace.ChangeName()
    workspace.LogStatus()
}

// MoveLeft moves window to the left column
func (workspace *Workspace) MoveLeft(wid uint32) {
    idx := workspace.right.IndexById(wid)
    if idx > -1 && workspace.right.Len() > 1 {
        win := workspace.right.WindowByIndex(idx)
        workspace.right.Remove(win)
        workspace.left.Add(win)
    }
}

// MoveRight moves window to the right column
func (workspace *Workspace) MoveRight(wid uint32) {
    idx := workspace.left.IndexById(wid)
    if idx > -1 && workspace.left.Len() > 1 {
        win := workspace.left.WindowByIndex(idx)
        workspace.left.Remove(win)
        workspace.right.Add(win)
    }
}

// MoveUp moves window upward
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

// MoveDown moves window downward
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

// ResizeLeft resizes current window to the left
func (workspace *Workspace) ResizeLeft(wid uint32) {
    idx := workspace.right.IndexById(wid)
    if idx > -1 && workspace.layout == LayoutLeftWide {
        workspace.layout = LayoutEqual
        workspace.Reshape()
        return
    }

    idx = workspace.left.IndexById(wid)
    if idx > -1 && workspace.layout == LayoutLeftWide {
        workspace.layout = LayoutEqual
        workspace.Reshape()
    }
}

// ResizeRight resizes current window to the right
func (workspace *Workspace) ResizeRight(wid uint32) {
    idx := workspace.left.IndexById(wid)
    if idx > -1 && workspace.layout == LayoutEqual {
        workspace.layout = LayoutLeftWide
        workspace.Reshape()
        return
    }
    idx = workspace.right.IndexById(wid)
    if idx > -1 && workspace.layout == LayoutEqual {
        workspace.layout = LayoutLeftWide
        workspace.Reshape()
    }
}

// Add adds new window to the workspace
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

// Remove removes window from the workspace
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
        workspace.ChangeName()
    case inRight:
        workspace.right.Remove(window)
    default:
        workspace.left.Remove(window)
    }
}

// Focus changes focus to current focus window
func (workspace *Workspace) Focus() {
    if workspace.focus == nil {
        return
    }

    workspace.focus.TakeFocus()

    if workspace.central.Len() > 0 {
        workspace.focus.UnsetBorder()
    }
}

// Refocus finds new focus window
func (workspace *Workspace) Refocus() {
    if workspace.focus == nil {
        return
    }

    focus := workspace.FocusDown()
    if focus != nil && focus.Id() != workspace.focus.Id() {
        workspace.focus = focus
        return
    }

    focus = workspace.FocusUp()
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

// FocusDown changes focus downward
func (workspace *Workspace) FocusDown() *Window {
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

// FocusUp changes focus upward
func (workspace *Workspace) FocusUp() *Window {
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

// FocusLeft changes focus to the left
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

// FocusRight changes focus to the right
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

// Activate makes workspace active, making all its windows visible
func (workspace *Workspace) Activate() {
    for i := 0; i < workspace.central.Len(); i++ {
        win := workspace.central.WindowByIndex(i)
        win.Map()
    }
    for i := 0; i < workspace.left.Len(); i++ {
        win := workspace.left.WindowByIndex(i)
        win.Map()
    }
    for i := 0; i < workspace.right.Len(); i++ {
        win := workspace.right.WindowByIndex(i)
        win.Map()
    }
}

// Deactivate makes workspace active, making all its windows invisible
func (workspace *Workspace) Deactivate() {
    for i := 0; i < workspace.central.Len(); i++ {
        win := workspace.central.WindowByIndex(i)
        win.DenyRemoval()
        win.Unmap()
    }
    for i := 0; i < workspace.left.Len(); i++ {
        win := workspace.left.WindowByIndex(i)
        win.DenyRemoval()
        win.Unmap()
    }
    for i := 0; i < workspace.right.Len(); i++ {
        win := workspace.right.WindowByIndex(i)
        win.DenyRemoval()
        win.Unmap()
    }
}

// FindWindow searches window by its identifier
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

// Reshape changes window sizes according to current layout
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
        x := workspace.left.SetWidth65()
        workspace.right.SetX(x)
        workspace.right.SetWidth35()
    }

    workspace.central.Reshape()
    workspace.left.Reshape()
    workspace.right.Reshape()
}

// ChangeName changes name of the workspace according
// to current focused window name
func (workspace *Workspace) ChangeName() {
    if workspace.conn == nil {
        return
    }

    repr := fmt.Sprintf("%d", workspace.id)
    if workspace.focus != nil {
        n := workspace.left.Len() + workspace.right.Len()
        n += workspace.central.Len()
        name, err := xutil.GetWMName(
            workspace.focus.Id(), workspace.conn,
        )
        if err == nil || n > 0 {
            name = string([]rune(name)[:config.NameLimit()])
            if n == 1 {
                repr = fmt.Sprintf("%d:%s", workspace.id, name)
            } else {
                repr = fmt.Sprintf("%d:%s(%d)", workspace.id, name, n)
            }
        }
    }
    names, err := xutil.GetDesktopNames(workspace.conn)
    if err != nil {
        names = make([]string, MaxWorkspaces)
        for i := 0; i < MaxWorkspaces; i++ {
            names[i] = fmt.Sprintf("%d", i+1)
        }
        xutil.SetDesktopNames(names, workspace.conn)
    }
    names[workspace.id - 1] = repr
    xutil.SetDesktopNames(names, workspace.conn)
}

// LogStatus logs workspace's information for debugging purposes
func (workspace *Workspace) LogStatus() {
    logging.Print("Workspace ID", workspace.id, " ")
    if workspace.focus != nil {
        logging.Println("focus = ", workspace.focus.Id())
    } else {
        logging.Println("focus = nil")
    }

    logging.Println("Central ")
    workspace.central.LogStatus()
    logging.Println("Left ")
    workspace.left.LogStatus()
    logging.Println("Right ")
    workspace.right.LogStatus()
    logging.Print("\n\n")
}

// WorkspaceManager represents a logical bridge
// between windows and workspaces
type WorkspaceManager struct {
    prev    uint32
    curr    uint32
    mailbox chan proto.Message
}

// NewWorkspaceManager creates instance of WorkspaceManager
func NewWorkspaceManager(monitors xutil.MonitorsInfo) *WorkspaceManager {
    input := make(chan proto.Message)
    next := make(chan proto.Message)
    mailbox := input

    for j := uint32(1); j < MaxWorkspaces-1; j++ {
        w := NewWorkspace(mailbox, input, next, j, monitors.Primary())
        go w.Run()
        input = next
        next = make(chan proto.Message)
    }

    if !monitors.IsDualSetup() {
        w := NewWorkspace(mailbox, input, nil, MaxWorkspaces-1, monitors.Primary())
        go w.Run()
    } else {
        w := NewWorkspace(mailbox, input, next, MaxWorkspaces-1, monitors.Primary())
        go w.Run()
        w = NewWorkspace(mailbox, next, nil, MaxWorkspaces, monitors.Secondary())
        go w.Run()
    }

    return &WorkspaceManager{DefaultWorkspace, DefaultWorkspace, mailbox}
}

// Mailbox returns channel, that is used for passing messages to workspaces
func (wrkmgr *WorkspaceManager) Mailbox() chan proto.Message {
    return wrkmgr.mailbox
}

// Prev returns id of previously active workspace
func (wrkmgr WorkspaceManager) Prev() uint32 {
    return wrkmgr.prev
}

// Curr return id of currently active workspace
func (wrkmgr *WorkspaceManager) Curr() uint32 {
    return wrkmgr.curr
}

// SetCurr sets current active workspace
func (wrkmgr *WorkspaceManager) SetCurr(n uint32) {
    if wrkmgr.curr != n {
        wrkmgr.prev = wrkmgr.curr
    }
    wrkmgr.curr = n
}

// SpecialWorkspace returns id of special workspace, used for external monitor
func (wrkmgr WorkspaceManager) SpecialWorkspace() uint32 {
    return MaxWorkspaces
}
