package ascii

import (
	"regexp"
)

var ansiPattern *regexp.Regexp = regexp.MustCompile(`(\x9B|\x1B\[)[0-?]*[ -/]*[@-~]`)

func ScapeAnsi(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}
