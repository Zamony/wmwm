package config

var (
	color         ColorFlag
	paddingTop    NonNegativeFlag
	paddingBottom NonNegativeFlag
	borderWidth   NonNegativeFlag
	nameLimit     NonNegativeFlag
	commands      StringsFlag
	debug         bool
)

func Color() uint32 {
	return uint32(color)
}

func PaddingTop() int {
	return int(paddingTop)
}

func PaddingBottom() int {
	return int(paddingBottom)
}

func BorderWidth() int {
	return int(borderWidth)
}

func NameLimit() int {
	if nameLimit < 1 {
		return 1
	}
	return int(nameLimit)
}

func Commands() []string {
	return commands.Value
}

func Debug() bool {
	return debug
}
