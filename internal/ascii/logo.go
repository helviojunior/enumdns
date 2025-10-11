package ascii

import (
	"fmt"
	"strings"

	"github.com/helviojunior/enumdns/internal/version"
)

// Logo returns the enumdns ascii logo
func Logo() string {
	txt := `                   
                                                     
{G}    ______                    {O}  ____  _   _______
{G}   / ____/___  __  ______ ___ {O} / __ \/ | / / ___/
{G}  / __/ / __ \/ / / / __ '__ \{O}/ / / /  |/ /\__ \ 
{G} / /___/ / / / /_/ / / / / / /{O} /_/ / /|  /___/ / 
{G}/_____/_/ /_/\__,_/_/ /_/ /_/{O}_____/_/ |_//____/  {B}
`

	v := fmt.Sprintf("Ver: %s-%s", version.Version, version.GitHash)
	r := 46 - len(v)
	if r < 0 {
		r = 0
	}
	txt += strings.Repeat(" ", r)
	txt += v + "{W}"
	txt = strings.Replace(txt, "{G}", "\033[32m", -1)
	txt = strings.Replace(txt, "{B}", "\033[36m", -1)
	txt = strings.Replace(txt, "{O}", "\033[33m", -1)
	txt = strings.Replace(txt, "{W}", "\033[0m", -1)
	return fmt.Sprintln(txt)
}

// LogoHelp returns the logo, with help
func LogoHelp(s string) string {
	return Logo() + "\n\n" + s
}
