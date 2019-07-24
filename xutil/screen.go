// Package xutil provides high-level abstraction for the XGB functions
package xutil

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/Zamony/wm/config"
)

// Screen represents logical screen
type Screen struct {
	width         int
	height        int
	xoffset       int
	paddingTop    int
	paddingBottom int
}

// NewScreen returns instance of Screen
func NewScreen(width, height, xoffset, pTop, pBot int) Screen {
	return Screen{width, height, xoffset, pTop, pBot}
}

// Width returns screen width
func (screen *Screen) Width() int {
	return screen.width
}

// Height returns screen height
func (screen *Screen) Height() int {
	return screen.height
}

// XOffset returns screen's ofsset on x-axis.
// For primary monitor XOffset is always zero.
func (screen *Screen) XOffset() int {
	return screen.xoffset
}

// PaddingTop returns top padding of the screen.
// One can use non-zero padding to left free space
// needed for placing a status bar
func (screen *Screen) PaddingTop() int {
	return screen.paddingTop
}

// PaddingBottom returns bottom padding of the screen.
// One can use non-zero padding to left free space
// needed for placing a status bar
func (screen *Screen) PaddingBottom() int {
	return screen.paddingBottom
}

// MonitorsInfo holds information about connected screens.
// Note that it's possible to connect only one external monitor.
type MonitorsInfo struct {
	primary   Screen
	secondary Screen
	dual      bool
}

// ReadMonitorsInfo returns information about connected monitors
// Note that it's possible to set padding only on primary screen
func ReadMonitorsInfo(conn *xgb.Conn) (MonitorsInfo, error) {
	var info MonitorsInfo
	r, err := xinerama.QueryScreens(conn).Reply()
	if err != nil {
		return info, err
	}

	nscreen := len(r.ScreenInfo)
	if nscreen < 1 {
		return info, errors.New("No screen info available")
	}

	if nscreen > 2 {
		return info, errors.New("Only 2 monitors setup is supported")
	}

	info.primary = Screen{
		int(r.ScreenInfo[0].Width),
		int(r.ScreenInfo[0].Height),
		int(r.ScreenInfo[0].XOrg),
		config.PaddingTop(),
		config.PaddingBottom(),
	}

	if nscreen == 2 {
		info.dual = true
		info.secondary = Screen{
			int(r.ScreenInfo[1].Width),
			int(r.ScreenInfo[1].Height),
			int(r.ScreenInfo[1].XOrg),
			0,
			0,
		}
	}

	return info, nil
}

// Primary returns information about primary screen
func (m MonitorsInfo) Primary() Screen {
	return m.primary
}

// Secondary returns information about secondary screen
func (m MonitorsInfo) Secondary() Screen {
	return m.secondary
}

// InPrimaryRegion checks whether point with specified
// x-coordinate is in the primary screen
func (m MonitorsInfo) InPrimaryRegion(x int) bool {
	return x <= m.primary.Width()
}

// IsDualSetup returns true if external monitor is being used
func (m MonitorsInfo) IsDualSetup() bool {
	return m.dual
}

// CommonWidth returns sum of primary and secondary monitors widths
func (m MonitorsInfo) CommonWidth() int {
	return m.primary.Width() + m.secondary.Width()
}
