package cmd

import (
	//"crypto/tls"
	"net/url"
	"os/user"
	"os"
	"fmt"
	"errors"
	"os/signal"
    "syscall"
    "time"

	"github.com/helviojunior/enumdns/internal"
	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/runner"
	"github.com/helviojunior/enumdns/pkg/readers"
    resolver "github.com/helviojunior/gopathresolver"
	"github.com/spf13/cobra"
)

var (
	opts = &runner.Options{}
	fileOptions = &readers.FileReaderOptions{}
	tProxy = ""
	forceCheck = false
)

var rootCmd = &cobra.Command{
	Use:   "enumdns",
	Short: "enumdns is a modular DNS recon tool",
	Long:  ascii.Logo(),
	Example: `
   - enumdns recon -d helviojunior.com.br -o enumdns.txt
   - enumdns recon -d helviojunior.com.br --write-jsonl
   - enumdns recon -D domains.txt --write-db   

   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -D domains.txt -w /tmp/wordlist.txt --write-db`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		
		usr, err := user.Current()
	    if err != nil {
	       return err
	    }

	    opts.Writer.UserPath = usr.HomeDir

	    if cmd.CalledAs() != "version" {
			fmt.Println(ascii.Logo())
		}

		if opts.Logging.Silence {
			log.EnableSilence()
		}

		if opts.Logging.Debug && !opts.Logging.Silence {
			log.EnableDebug()
			log.Debug("debug logging enabled")
		}

        if opts.Writer.TextFile != "" {

        	opts.Writer.TextFile, err = resolver.ResolveFullPath(opts.Writer.TextFile)
	        if err != nil {
	            return err
	        }

            opts.Writer.Text = true
        }

        //Check Proxy config
        if tProxy != "" {
        	u, err := url.Parse(tProxy)
        	if err != nil {
	        	return errors.New("Error parsing URL: " + err.Error())
	        }

        	_, err = internal.FromURL(u, nil)
        	if err != nil {
	        	return errors.New("Error parsing URL: " + err.Error())
	        }
	        opts.Proxy = u
	        fileOptions.ProxyUri = opts.Proxy

			port := u.Port()
			if port == "" {
				port = "1080"
			}
	        log.Warn("Setting proxy to " + u.Scheme + "://" + u.Hostname() + ":" + port)
        }else{
        	opts.Proxy = nil
        }
        
		return nil
	},
}

func Execute() {
	
	ascii.SetConsoleColors()

	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        ascii.ClearLine()
        fmt.Fprintf(os.Stderr, "\r\n")
        ascii.ClearLine()
        ascii.ShowCursor()
        log.Warn("interrupted, shutting down...                            ")
        ascii.ClearLine()
        fmt.Printf("\n")
        os.Exit(2)
    }()

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true
	err := rootCmd.Execute()
	if err != nil {
		var cmd string
		c, _, cerr := rootCmd.Find(os.Args[1:])
		if cerr == nil {
			cmd = c.Name()
		}

		v := "\n"

		if cmd != "" {
			v += fmt.Sprintf("An error occured running the `%s` command\n", cmd)
		} else {
			v += "An error has occured. "
		}

		v += "The error was:\n\n" + fmt.Sprintf("```%s```", err)
		fmt.Println(ascii.Markdown(v))

		os.Exit(1)
	}

	//Time to wait the logger flush
	time.Sleep(time.Second/4)
    ascii.ShowCursor()
    fmt.Printf("\n")
}

func init() {
	// Disable Certificate Validation (Globally)
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	rootCmd.PersistentFlags().BoolVarP(&opts.Logging.Debug, "debug-log", "D", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&opts.Logging.Silence, "quiet", "q", false, "Silence (almost all) logging")
	rootCmd.PersistentFlags().BoolVarP(&forceCheck, "force", "F", false, "Force to check all hosts again.")

	rootCmd.PersistentFlags().StringVarP(&opts.Writer.TextFile, "write-text-file", "o", "", "The file to write Text lines to")
    

	//rootCmd.PersistentFlags().BoolVarP(&opts.DnsOverHttps.SkipSSLCheck, "ssl-insecure", "K", true, "SSL Insecure")
	rootCmd.PersistentFlags().StringVarP(&tProxy, "proxy", "X", "", "Proxy to pass traffic through: <scheme://ip:port> (e.g., socks4://user:pass@proxy_host:1080")
	//rootCmd.PersistentFlags().StringVarP(&opts.DnsOverHttps.ProxyUser, "proxy-user", "", "", "Proxy User")
	//rootCmd.PersistentFlags().StringVarP(&opts.DnsOverHttps.ProxyPassword, "proxy-pass", "", "", "Proxy Password")

}
