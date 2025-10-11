package cmd

import (
	//"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"

	"github.com/bob-reis/enumdns/internal"
	"github.com/bob-reis/enumdns/internal/ascii"
	"github.com/bob-reis/enumdns/internal/tools"
	"github.com/bob-reis/enumdns/pkg/log"
	"github.com/bob-reis/enumdns/pkg/readers"
	"github.com/bob-reis/enumdns/pkg/runner"
	"github.com/bob-reis/enumdns/pkg/writers"
	resolver "github.com/helviojunior/gopathresolver"
	"github.com/spf13/cobra"
)

const sqlitePrefix = "sqlite:///"

var (
	opts        = &runner.Options{}
	fileOptions = &readers.FileReaderOptions{}
	tProxy      = ""
	forceCheck  = false
	tempFolder  = ""
)

var rootCmd = &cobra.Command{
	Use:   "enumdns",
	Short: "enumdns is a modular DNS recon tool",
	Long:  ascii.Logo(),
	Example: `
       - enumdns recon -d test.com -o enumdns.txt
       - enumdns recon -d test.com --write-jsonl
       - enumdns recon -L domains.txt --write-db

   - enumdns brute -d test.com -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d test.com -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -L domains.txt -w /tmp/wordlist.txt --write-db

       - enumdns threat-analysis -d test.com --all-techniques -o threats.txt
       - enumdns threat-analysis -d test.com --typosquatting --write-jsonl
       - enumdns threat-analysis -L domains.txt --bitsquatting --write-db
       - enumdns threat-analysis -d recife.pe.gov.br --all-techniques --focus-suffix=gov.br --emit-candidates -o gov-br.txt
       - enumdns threat-analysis -d yeslinux.com.br --all-techniques --target-tlds com,net,org,co,info,io,com.br,net.br,org.br
       - enumdns threat-analysis -d updates.microsoft.com --all-techniques --span-last3

   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
   - enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db

       - enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
       - enumdns resolve file -L /tmp/host_list.txt --write-jsonl
       - enumdns resolve file -L /tmp/host_list.txt --write-db`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if cmd.CalledAs() != "version" && !opts.Logging.Silence {
			fmt.Println(ascii.Logo())
		}

		if opts.Logging.Silence {
			log.EnableSilence()
		}

		if opts.Logging.Debug && !opts.Logging.Silence {
			log.EnableDebug()
			log.Debug("debug logging enabled")
		}

		usr, err := user.Current()
		if err != nil {
			return err
		}

		opts.Writer.UserPath = usr.HomeDir

		opts.Writer.CtrlDbURI = sqlitePrefix + opts.Writer.UserPath + "/.enumdns.db"

		basePath := ""
		if opts.StoreTempAsWorkspace {
			basePath = "./"
		}

		if tempFolder, err = tools.CreateDir(tools.TempFileName(basePath, "enumdns_", "")); err != nil {
			log.Error("error creatting temp folder", "err", err)
			os.Exit(2)
		}

		if opts.LocalWorkspace {
			tmp, err := resolver.ResolveFullPath("./enumdns_ctrl.db")
			if err != nil {
				return err
			}

			opts.Writer.CtrlDbURI = sqlitePrefix + tmp
			log.Info("Control DB", "Path", tmp)
		}

		if opts.Writer.NoControlDb {
			opts.Writer.CtrlDbURI = sqlitePrefix + tools.TempFileName(tempFolder, "enumdns_", ".db")
		}

		if opts.Writer.TextFile != "" {

			opts.Writer.TextFile, err = resolver.ResolveFullPath(opts.Writer.TextFile)
			if err != nil {
				return err
			}

			opts.Writer.Text = true
		}

		if opts.DnsServer == "" {
			opts.DnsServer = tools.GetDefaultDnsServer("")
		}
		opts.PrivateDns = tools.IsPrivateIP(opts.DnsServer)

		fileOptions.DnsServer = opts.DnsServer + ":" + fmt.Sprintf("%d", opts.DnsPort)

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
		} else {
			opts.Proxy = nil
		}

		return nil
	},
}

func Execute() {

	if err := ascii.SetConsoleColors(); err != nil {
		log.Error("failed to set console colors", "err", err)
	}

	c := make(chan os.Signal, 1)
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
		if err := tools.RemoveFolder(tempFolder); err != nil {
			log.Error("failed to remove temp folder", "err", err)
		}
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
	time.Sleep(time.Second / 4)
	if err := tools.RemoveFolder(tempFolder); err != nil {
		log.Error("failed to remove temp folder", "err", err)
	}
	ascii.ShowCursor()
	fmt.Printf("\n")
}

func init() {
	// Disable Certificate Validation (Globally)
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	rootCmd.PersistentFlags().StringVarP(&opts.DnsServer, "server", "s", "", "DNS Server")
	rootCmd.PersistentFlags().IntVar(&opts.DnsPort, "port", 53, "DNS Server Port")
	rootCmd.PersistentFlags().StringVarP(&opts.DnsProtocol, "protocol", "", "UDP", "DNS Server protocol (TCP/UDP)")

	rootCmd.PersistentFlags().BoolVarP(&opts.Logging.Debug, "debug-log", "D", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&opts.Logging.Silence, "quiet", "q", false, "Silence (almost all) logging")
	rootCmd.PersistentFlags().BoolVarP(&forceCheck, "force", "F", false, "Force to check all hosts again.")

	// Logging control for subcommands
	rootCmd.PersistentFlags().BoolVar(&opts.Logging.LogScanErrors, "log-scan-errors", false, "Log scan errors (timeouts, DNS errors, etc.) to stderr (warning: can be verbose!)")

	rootCmd.PersistentFlags().StringVarP(&opts.Writer.TextFile, "write-text-file", "o", "", "The file to write Text lines to")

	//rootCmd.PersistentFlags().BoolVarP(&opts.DnsOverHttps.SkipSSLCheck, "ssl-insecure", "K", true, "SSL Insecure")
	rootCmd.PersistentFlags().StringVarP(&tProxy, "proxy", "X", "", "Proxy to pass traffic through: <scheme://ip:port> (e.g., socks4://user:pass@proxy_host:1080")
	//rootCmd.PersistentFlags().StringVarP(&opts.DnsOverHttps.ProxyUser, "proxy-user", "", "", "Proxy User")
	//rootCmd.PersistentFlags().StringVarP(&opts.DnsOverHttps.ProxyPassword, "proxy-pass", "", "", "Proxy Password")

	// "Threads" & other
	rootCmd.PersistentFlags().IntVarP(&opts.Scan.Threads, "threads", "t", 6, "Number of concurrent threads (goroutines) to use")
	rootCmd.PersistentFlags().IntVarP(&opts.Scan.Timeout, "timeout", "T", 60, "Number of seconds before considering a page timed out")

	rootCmd.PersistentFlags().BoolVar(&opts.Writer.NoControlDb, "disable-control-db", false, "Disable utilization of database ~/.enumdns.db.")
	rootCmd.PersistentFlags().BoolVar(&opts.StoreTempAsWorkspace, "local-temp", false, "Use execution path to store temp files")
	rootCmd.PersistentFlags().BoolVar(&opts.LocalWorkspace, "local-workspace", false, "Use execution path to store .enumdns.db intead of user home")

	// Write options for scan subcommands
	rootCmd.PersistentFlags().BoolVar(&opts.Writer.Db, "write-db", false, "Write results to a SQLite database")
	rootCmd.PersistentFlags().StringVar(&opts.Writer.DbURI, "write-db-uri", "sqlite://enumdns.sqlite3", "The database URI to use. Supports SQLite, Postgres, and MySQL (e.g., postgres://user:pass@host:port/db)")
	rootCmd.PersistentFlags().BoolVar(&opts.Writer.DbDebug, "write-db-enable-debug", false, "Enable database query debug logging (warning: verbose!)")
	rootCmd.PersistentFlags().BoolVar(&opts.Writer.Csv, "write-csv", false, "Write results as CSV (has limited columns)")
	rootCmd.PersistentFlags().StringVar(&opts.Writer.CsvFile, "write-csv-file", "enumdns.csv", "The file to write CSV rows to")
	rootCmd.PersistentFlags().BoolVar(&opts.Writer.Jsonl, "write-jsonl", false, "Write results as JSON lines")
	rootCmd.PersistentFlags().StringVar(&opts.Writer.JsonlFile, "write-jsonl-file", "enumdns.jsonl", "The file to write JSON lines to")
	rootCmd.PersistentFlags().BoolVar(&opts.Writer.None, "write-none", false, "Use an empty writer to silence warnings")

	rootCmd.PersistentFlags().BoolVar(&opts.Writer.ELastic, "write-elastic", false, "Write results to a SQLite database")
	rootCmd.PersistentFlags().StringVar(&opts.Writer.ELasticURI, "write-elasticsearch-uri", "http://localhost:9200/enumdns", "The elastic search URI to use. (e.g., http://user:pass@host:9200/index)")

}

// ConfigureWriters configura os writers baseado nas opções do usuário
func ConfigureWriters(writersList *[]writers.Writer) error {
	// Writer de controle (obrigatório)
	w, err := writers.NewDbWriter(opts.Writer.CtrlDbURI, opts.Writer.DbDebug)
	if err != nil {
		return err
	}
	*writersList = append(*writersList, w)

	// Writer de stdout
	if !opts.Logging.Silence {
		w, err := writers.NewStdoutWriter()
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	// Writers opcionais
	if opts.Writer.Text {
		w, err := writers.NewTextWriter(opts.Writer.TextFile)
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if opts.Writer.Jsonl {
		w, err := writers.NewJsonWriter(opts.Writer.JsonlFile)
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if opts.Writer.Db {
		w, err := writers.NewDbWriter(opts.Writer.DbURI, opts.Writer.DbDebug)
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if opts.Writer.Csv {
		w, err := writers.NewCsvWriter(opts.Writer.CsvFile)
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if opts.Writer.ELastic {
		w, err := writers.NewElasticWriter(opts.Writer.ELasticURI)
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if opts.Writer.None {
		w, err := writers.NewNoneWriter()
		if err != nil {
			return err
		}
		*writersList = append(*writersList, w)
	}

	if len(*writersList) == 0 {
		log.Warn("no writers have been configured. to persist probe results, add writers using --write-* flags")
	}

	return nil
}
