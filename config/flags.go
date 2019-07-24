// Package config parses command line arguments
// and provides access to them
package config

var (
	color         ColorFlag
	paddingTop    NonNegativeFlag
	paddingBottom NonNegativeFlag
	borderWidth   NonNegativeFlag
	nameLimit     NonNegativeFlag
	commands      StringsFlag
	terminal      string
	launcher      string
	locker        string
	debug         bool
)

// Color returns --color command line argument value
func Color() uint32 {
	return uint32(color)
}

// PaddingTop returns --padding-top command line argument value
func PaddingTop() int {
	return int(paddingTop)
}

// PaddingBottom returns --padding-bottom command line argument value
func PaddingBottom() int {
	return int(paddingBottom)
}

// BorderWidth returns --border-width command line argument value
func BorderWidth() int {
	return int(borderWidth)
}

// NameLimit returns --name-limit command line argument value
func NameLimit() int {
	if nameLimit < 1 {
		return 1
	}
	return int(nameLimit)
}

// Commands returns values of --exec command line arguments
func Commands() []string {
	return commands.Value
}

// TerminalCmd returns value of --term command line argument
func TerminalCmd() string {
	return terminal
}

// LockerCmd returns value of --lock command line argument
func LockerCmd() string {
	return locker
}

// LauncherCmd returns value of --launcher command line argument
func LauncherCmd() string {
	return launcher
}

// Debug returns value of --debug command line argument
func Debug() bool {
	return debug
}
