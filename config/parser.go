package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

const ColorLimit = 1 << 24

type ColorFlag uint32

func (c *ColorFlag) String() string {
	return fmt.Sprintf("0x%x", *c)
}

func (c *ColorFlag) Set(v string) error {
	col, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		return err
	}
	if col < ColorLimit {
		*c = ColorFlag(col)
		return nil
	}

	return errors.New("Not a color in RGB palette")
}

type NonNegativeFlag int

func (u *NonNegativeFlag) String() string {
	return fmt.Sprint(*u)
}

func (u *NonNegativeFlag) Set(v string) error {
	x, err := strconv.ParseInt(v, 0, strconv.IntSize)
	if x < 0 {
		return errors.New("Non-negative number required")
	}
	if err == nil {
		*u = NonNegativeFlag(x)
	}
	return err
}

type StringsFlag struct {
	Value []string
}

func (s *StringsFlag) String() string {
	return fmt.Sprint(s.Value)
}

func (s *StringsFlag) Set(v string) error {
	s.Value = append(s.Value, v)
	return nil
}

func ParseArgs() {
	flag.Var(&color, "color", "Background and border color")
	flag.Var(&paddingTop, "padding-top", "Value of top padding")
	flag.Var(&paddingBottom, "padding-bottom", "Value of bottom padding")
	flag.Var(&borderWidth, "border-width", "Border width of focused window")
	flag.Var(&nameLimit, "name-limit", "Maximum length of workspace name")
	flag.Var(&commands, "exec", "Commands to execute at startup")
	flag.BoolVar(
		&debug, "debug", false,
		"Outputs debug information to Stderr",
	)
	flag.Parse()
}
