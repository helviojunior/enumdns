package cmd

import (
	"errors"
	"log/slog"
	"os"

	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/helviojunior/enumdns/pkg/runner"
	"github.com/helviojunior/enumdns/pkg/writers"
	"github.com/spf13/cobra"
)

var adRunner *runner.AdRunner
var adWriters = []writers.Writer{}
var adCmd = &cobra.Command{
	Use:   "ad",
	Short: "Dump AD-integrated DNS zones and records (ADIDNS)",
	Long: ascii.LogoHelp(ascii.Markdown(`
# ad

Dump every DNS zone and record stored inside Active Directory (ADIDNS) by
querying a domain controller over LDAP, enumerating the DomainDnsZones /
ForestDnsZones partitions and decoding the binary **dnsRecord** attribute of
each node.

Authentication supports three methods:

1. Username and password (` + "`-u/-p`" + ` plus ` + "`-d`" + ` domain)
2. NTLM hash / pass-the-hash (` + "`-H`" + `, as ` + "`NT`" + ` or ` + "`LM:NT`" + `)
3. Kerberos ticket cache (` + "`-k`" + ` with ` + "`--krb5-ccache`" + ` or $KRB5CCNAME)
`)),
	Example: `
   - enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' -o adidns.txt
   - enumdns ad -s 10.0.0.1 --ldaps -d corp.local -u jdoe -p 'Passw0rd!' --write-jsonl
   - enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -H aad3b435b51404eeaad3b435b51404ee:<nthash> --write-db
   - enumdns ad --target-ip 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' --write-db
   - KRB5CCNAME=/tmp/jdoe.ccache enumdns ad -s 10.0.0.1 -d corp.local -k --write-db`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		// Annoying quirk, but because I'm overriding PersistentPreRun
		// here which overrides the parent it seems.
		// So we need to explicitly call the parent's one now.
		if err = rootCmd.PersistentPreRunE(cmd, args); err != nil {
			return err
		}

		// An slog-capable logger to use with the runner.
		logger := slog.New(log.Logger)

		if err := ConfigureWriters(&adWriters); err != nil {
			return err
		}

		// Print every record to stdout (mirrors recon's behaviour).
		if !opts.Logging.Silence && len(adWriters) > 1 {
			if stdoutWriter, ok := adWriters[1].(*writers.StdoutWriter); ok {
				stdoutWriter.WriteAll = true
			}
		}

		adRunner, err = runner.NewAD(logger, *opts, adWriters)
		if err != nil {
			return err
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Without a forced target, the domain controllers are discovered from
		// DNS, which requires the AD domain name.
		if opts.AD.TargetIP == "" && opts.AD.Domain == "" {
			return errors.New("specify the AD domain (-d) to discover the global catalogs, or force a server with --target-ip")
		}

		if !opts.AD.Kerberos && opts.AD.NtlmHash == "" && opts.AD.Username == "" {
			log.Warn("no credentials provided; attempting an anonymous LDAP bind (usually rejected by AD)")
		}

		if opts.AD.Kerberos && opts.AD.CCachePath != "" && !tools.FileExists(opts.AD.CCachePath) {
			return errors.New("kerberos ccache file is not readable")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		zones, records, err := adRunner.Run()
		if err != nil {
			log.Error("error dumping AD DNS", "err", err)
			adRunner.Close()
			os.Exit(2)
		}

		adRunner.Close()

		log.Infof("Dumped %s record(s) across %s zone(s)",
			tools.FormatInt(records), tools.FormatInt(zones))
	},
}

func init() {
	rootCmd.AddCommand(adCmd)

	adCmd.Flags().StringVar(&opts.AD.TargetIP, "target-ip", "", "Force a specific LDAP server (host/IP), skipping Global Catalog discovery")
	adCmd.Flags().IntVar(&opts.AD.Port, "ldap-port", 0, "LDAP port (default 389, or 636 with --ldaps)")
	adCmd.Flags().BoolVar(&opts.AD.LDAPS, "ldaps", false, "Use LDAPS (TLS) instead of plaintext LDAP")

	adCmd.Flags().StringVarP(&opts.AD.Username, "username", "u", "", "Username for LDAP authentication")
	adCmd.Flags().StringVarP(&opts.AD.Password, "password", "p", "", "Password for LDAP authentication")
	adCmd.Flags().StringVarP(&opts.AD.Domain, "domain", "d", "", "Active Directory domain / realm")

	adCmd.Flags().StringVarP(&opts.AD.NtlmHash, "hashes", "H", "", "NTLM hash for pass-the-hash (format: NT or LM:NT)")

	adCmd.Flags().BoolVarP(&opts.AD.Kerberos, "kerberos", "k", false, "Authenticate using a Kerberos ticket cache (ccache)")
	adCmd.Flags().StringVar(&opts.AD.CCachePath, "krb5-ccache", "", "Path to the Kerberos ccache file (defaults to $KRB5CCNAME)")
	adCmd.Flags().StringVar(&opts.AD.Krb5Conf, "krb5-config", "", "Path to krb5.conf (defaults to /etc/krb5.conf)")
	adCmd.Flags().StringVar(&opts.AD.SPN, "spn", "", "Service principal name for Kerberos (defaults to ldap/<dc>)")

	adCmd.Flags().StringVar(&opts.AD.BaseDN, "base-dn", "", "Restrict the dump to a single LDAP container (skips autodetection)")
}
