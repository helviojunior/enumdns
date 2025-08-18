//go:build windows

package ascii

func GetNextSpinner(spin string) string {
	switch spin {
	case "[=====]":
		return "[ ====]"
	case "[ ====]":
		return "[  ===]"
	case "[  ===]":
		return "[=  ==]"
	case "[=  ==]":
		return "[==  =]"
	case "[==  =]":
		return "[===  ]"
	case "[===  ]":
		return "[==== ]"
	default:
		return "[=====]"
	}
}

func ColoredSpin(spin string) string {
	return spin
}
