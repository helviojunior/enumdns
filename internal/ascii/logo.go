package ascii

// Logo returns the enumdns ascii logo
func Logo() string {
	return `                   
                                                     
    ______                      ____  _   _______
   / ____/___  __  ______ ___  / __ \/ | / / ___/
  / __/ / __ \/ / / / __ '__ \/ / / /  |/ /\__ \ 
 / /___/ / / / /_/ / / / / / / /_/ / /|  /___/ / 
/_____/_/ /_/\__,_/_/ /_/ /_/_____/_/ |_//____/  
                                                 
                                                                                                                                                                                                                                                                                                                                                                       
`
}

// LogoHelp returns the logo, with help
func LogoHelp(s string) string {
	return Logo() + "\n\n" + s
}
