package main

import (
    "fmt"
	"github.com/BurntSushi/xgb"
    "github.com/Zamony/wm/proto"
    "github.com/Zamony/wm/xutil"
    "github.com/Zamony/wm/config"
    "github.com/Zamony/wm/logging"
)

const (
    LayoutFull = iota
    LayoutEqual
    LayoutLeftWide
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
    conn    *xgb.Conn
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
        conn:    nil,
    }
}

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
            workspace.CleanUp()
            if workspace.FindWindow(msg.From) != nil {
                workspace.handleMsg(msg)
            } else if workspace.next != nil {
                workspace.next <- msg
            }
        case workspace.id:
            workspace.CleanUp()
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
            workspace.Focus()// workspace.focus.TakeFocus()
            if workspace.id == MaxWorkspaces {
                workspace.Activate()
            }
        }
    case proto.Detach:
        if workspace.focus != nil {
            win := workspace.focus
            go func() {
                win.SendAttach(msg.From)
                if workspace.id == MaxWorkspaces {
                    win.SendActivate(msg.From)
                }
            }()
            win.Defocus()
            workspace.Refocus()
            workspace.Remove(win)
            win.Unmap()
            workspace.Reshape()
            workspace.Focus()// workspace.focus.TakeFocus()
            go func() {
                if workspace.focus != nil {
                    workspace.focus.SendFocusHere()
                }
            }()
        }
    case proto.Remove:
        win := NewWindow(msg.From, workspace.headc, msg.XConn)
        if win != nil && workspace.focus != nil {
            if win.Id() == workspace.focus.Id() {
                win.Defocus()
                workspace.Refocus()
            }
            workspace.Remove(win)
            workspace.Reshape()
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
                workspace.Focus()// workspace.focus.TakeFocus()
            }
        }
    case proto.FocusLeft:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusLeft()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.FocusRight:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusRight()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.FocusUp:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusUp()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.FocusDown:
        workspace.focus.Defocus()
        workspace.focus = workspace.FocusDown()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.Activate:
        workspace.Activate()
        workspace.Focus()// workspace.focus.TakeFocus()
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
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.MoveDown:
        workspace.MoveDown(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.MoveLeft:
        workspace.MoveLeft(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()// workspace.focus.TakeFocus()
    case proto.MoveRight:
        workspace.MoveRight(workspace.focus.Id())
        workspace.Reshape()
        workspace.Focus()// workspace.focus.TakeFocus()
    default:
        return
    }

    workspace.SetName()
    workspace.LogStatus()
}

func (workspace *Workspace) CleanUp() {
    cleaned := false
    for i := 0; i < workspace.central.Len(); i++ {
        win := workspace.central.WindowByIndex(i)
        if managed := win.CouldBeManaged(); !managed {
            if win.Id() == workspace.focus.Id() {
                workspace.Refocus()
            }
            workspace.Remove(win)
            cleaned = true
        }
    }
    for i := 0; i < workspace.left.Len(); i++ {
        win := workspace.left.WindowByIndex(i)
        if managed := win.CouldBeManaged(); !managed {
            if win.Id() == workspace.focus.Id() {
                workspace.Refocus()
            }
            workspace.Remove(win)
            cleaned = true
        }
    }
    for i := 0; i < workspace.right.Len(); i++ {
        win := workspace.right.WindowByIndex(i)
        if managed := win.CouldBeManaged(); !managed {
            if win.Id() == workspace.focus.Id() {
                workspace.Refocus()
            }
            workspace.Remove(win)
            cleaned = true
        }
    }

    if cleaned {
        workspace.Reshape()
        if workspace.focus != nil {
            workspace.Focus()// workspace.focus.TakeFocus()
        }
    }
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
        workspace.SetName()
    case inRight:
        workspace.right.Remove(window)
    default:
        workspace.left.Remove(window)
    }
}

func (workspace *Workspace) Focus() {
    if workspace.focus == nil {
        return
    }

    workspace.focus.TakeFocus()

    if workspace.central.Len() > 0 {
        workspace.focus.UnsetBorder()
    }
}

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

func (workspace *Workspace) Deactivate() {
    for i := 0; i < workspace.central.Len(); i++ {
        win := workspace.central.WindowByIndex(i)
        win.Unmap()
    }
    for i := 0; i < workspace.left.Len(); i++ {
        win := workspace.left.WindowByIndex(i)
        win.Unmap()
    }
    for i := 0; i < workspace.right.Len(); i++ {
        win := workspace.right.WindowByIndex(i)
        win.Unmap()
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

func (workspace *Workspace) SetName() {
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

type WorkspaceManager struct {
    prev    uint32
    curr    uint32
    mailbox chan proto.Message
}

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

func (wrkmgr *WorkspaceManager) Mailbox() chan proto.Message {
    return wrkmgr.mailbox
}

func (wrkmgr WorkspaceManager) Prev() uint32 {
    return wrkmgr.prev
}

func (wrkmgr *WorkspaceManager) Curr() uint32 {
    return wrkmgr.curr
}

func (wrkmgr *WorkspaceManager) SetCurr(n uint32) {
    if wrkmgr.curr != n {
        wrkmgr.prev = wrkmgr.curr
    }
    wrkmgr.curr = n
}

func (wrkmgr WorkspaceManager) SpecialWorkspace() uint32 {
    return MaxWorkspaces
}
