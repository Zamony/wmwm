package xutil

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/Zamony/wm/config"
)

type Screen struct {
	width         int
	height        int
	xoffset       int
	paddingTop    int
	paddingBottom int
}

func NewScreen(width, height, xoffset, pTop, pBot int) Screen {
	return Screen{width, height, xoffset, pTop, pBot}
}

func (screen *Screen) Width() int {
	return screen.width
}

func (screen *Screen) Height() int {
	return screen.height
}

func (screen *Screen) XOffset() int {
	return screen.xoffset
}

func (screen *Screen) PaddingTop() int {
	return screen.paddingTop
}

func (screen *Screen) PaddingBottom() int {
	return screen.paddingBottom
}

type MonitorsInfo struct {
	primary   Screen
	secondary Screen
	dual      bool
}

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

func (m MonitorsInfo) Primary() Screen {
	return m.primary
}

func (m MonitorsInfo) Secondary() Screen {
	return m.secondary
}

func (m MonitorsInfo) InPrimaryRegion(x int) bool {
	return x <= m.primary.Width()
}

func (m MonitorsInfo) IsDualSetup() bool {
	return m.dual
}

func (m MonitorsInfo) CommonWidth() int {
	return m.primary.Width() + m.secondary.Width()
}
