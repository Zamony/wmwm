// Package config parses command line arguments
// and provides access to them
package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

// ColorLimit is a maximum color value in RGB palette
const ColorLimit = 1 << 24

// ColorFlag is a type used to represent a color as a CLI argument
type ColorFlag uint32

// String returns string representation of the color
func (c *ColorFlag) String() string {
	return fmt.Sprintf("0x%x", *c)
}

// Set parses and sets color
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

// NonNegativeFlag is a type used to represent
// non-negative number as a CLI argument
type NonNegativeFlag int

// String returns string representation of the non-negative number
func (u *NonNegativeFlag) String() string {
	return fmt.Sprint(*u)
}

// Set sets non negative number value
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

// StringsFlag is a type used to represent multiple string parameters
type StringsFlag struct {
	Value []string
}

// String returns string representation of multiple string parameters
func (s *StringsFlag) String() string {
	return fmt.Sprint(s.Value)
}

// Set appends string value to the end of current array of strings
func (s *StringsFlag) Set(v string) error {
	s.Value = append(s.Value, v)
	return nil
}

// ParseArgs parses CLI arguments
func ParseArgs() {
	flag.Var(&color, "color", "Background and border color")
	flag.Var(&paddingTop, "padding-top", "Value of top padding")
	flag.Var(&paddingBottom, "padding-bottom", "Value of bottom padding")
	flag.Var(&borderWidth, "border-width", "Border width of focused window")
	flag.Var(&nameLimit, "name-limit", "Maximum length of workspace name")
	flag.Var(&commands, "exec", "Commands to execute at startup")
	flag.StringVar(&terminal, "term", "xterm", "A command to launch terminal emulator")
	flag.StringVar(&launcher, "launcher", "rofi -show run", "A command to show application launcher")
	flag.StringVar(&locker, "lock", "slock", "A command to lock screen")
	flag.BoolVar(
		&debug, "debug", false,
		"Outputs debug information to Stderr",
	)
	flag.Parse()
}
