package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/helviojunior/enumdns/internal/ascii"
	"github.com/helviojunior/enumdns/internal/tools"
	"github.com/helviojunior/enumdns/pkg/log"
	"github.com/spf13/cobra"
)

var wordlistOpts = struct {
	Inputs      []string
	Exclude     []string
	MinLength   int
	MaxLength   int
	KeepTLD     bool
	KeepNumeric bool
}{}

// fqdnRun matches a maximal run of DNS-label characters (letters, digits,
// underscore and hyphen) optionally joined by dots, so a name such as
// www.example.com is captured as a single run and can be split into its labels
// while remembering which label was the trailing TLD. A bare token (no dot) is
// matched too, and is never TLD-filtered.
var fqdnRun = regexp.MustCompile(`[a-z0-9_-]+(?:\.[a-z0-9_-]+)*`)

// validLabel reports whether a token is a syntactically valid DNS label: it must
// start and end with an alphanumeric or underscore and may contain hyphens only
// in the middle. This is what drops flag-like garbage such as "--allow-parent-soa"
// or "-foo"/"bar-", which is exactly how CLI flags leak into a scraped wordlist.
var validLabel = regexp.MustCompile(`^[a-z0-9_]([a-z0-9_-]*[a-z0-9_])?$`)

// allNumeric matches a token made entirely of digits. Such labels are valid DNS
// labels in a subdomain position (RFC 1123 §2.1), but in a scraped wordlist they
// are almost always noise (IP octets, TTLs, SOA serials) rather than useful
// brute-force candidates, so they are dropped unless --keep-numeric is set.
var allNumeric = regexp.MustCompile(`^[0-9]+$`)

var wordlistCmd = &cobra.Command{
	Use:   "wordlist",
	Short: "Build a clean DNS-label wordlist from arbitrary files",
	Long: ascii.LogoHelp(ascii.Markdown(`
# wordlist

Generate a deduplicated, sorted DNS-label wordlist from one or more input files
(plain text, previous results, scraped data, etc.).

Each line is split into candidate labels on anything that is not a letter, digit,
underscore or hyphen, so an FQDN like _host.example.com_ yields _host_ and
_example_. Each token is then:

  - lower-cased and deduplicated;
  - validated as a DNS label (flag-like junk such as _--allow-parent-soa_ is dropped);
  - dropped when all-numeric (IP octets, TTLs, SOA serials) unless _--keep-numeric_;
  - filtered by length (_--min-length_ / _--max-length_);
  - filtered against your own exclusion list (_--exclude_).

Only the trailing TLD of a _dotted_ name is removed (the _com_ of _www.example.com_,
the _br_ of _ed.dev.br_), so internal labels that merely share a name with a gTLD
(e.g. _dev_) are kept. A bare token (no dot) is never treated as a TLD. Use
_--keep-tld_ to keep the trailing TLD as well.

Inputs accept shell globs. Quote them (e.g. _-i '*.txt'_) so enumdns expands the
pattern itself. The output file (_-o_) is always removed from the inputs, so
regenerating the list never feeds back on itself.
`)),
	Example: `
   - enumdns wordlist -i '*.txt' -o custom_wl.txt
   - enumdns wordlist -i '/data/recon/*.txt' -i extra.txt --exclude cloud,mail,onmicrosoft -o wl.txt
   - enumdns wordlist -i dump.txt --exclude deny.txt --min-length 3 --keep-tld -o wl.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Treat positional args as additional inputs (handles a shell-expanded
		// glob that was passed unquoted).
		patterns := append(append([]string{}, wordlistOpts.Inputs...), args...)
		return runWordlist(patterns)
	},
}

func runWordlist(patterns []string) error {
	if len(patterns) == 0 {
		return errors.New("at least one input file or glob is required (-i)")
	}
	if wordlistOpts.MinLength < 1 {
		wordlistOpts.MinLength = 1
	}
	if wordlistOpts.MaxLength <= 0 || wordlistOpts.MaxLength > 63 {
		wordlistOpts.MaxLength = 63
	}
	if wordlistOpts.MinLength > wordlistOpts.MaxLength {
		return fmt.Errorf("--min-length (%d) cannot be greater than --max-length (%d)", wordlistOpts.MinLength, wordlistOpts.MaxLength)
	}

	// opts.Writer.TextFile was already resolved to an absolute path by the root
	// PersistentPreRunE. Keep it out of the inputs so a previous run's output can
	// never reseed the wordlist (the feedback loop behind leaked flag labels).
	outPath := strings.TrimSpace(opts.Writer.TextFile)

	files, err := expandWordlistInputs(patterns, outPath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no input files matched the given pattern(s)")
	}
	log.Infof("Reading %s input file(s)", tools.FormatInt(len(files)))

	exclude, err := buildExclusionSet(wordlistOpts.Exclude)
	if err != nil {
		return err
	}
	if len(exclude) > 0 {
		log.Infof("Loaded %s exclusion term(s)", tools.FormatInt(len(exclude)))
	}

	words := map[string]struct{}{}
	scanned := 0
	for _, f := range files {
		n, err := collectTokens(f, exclude, words)
		if err != nil {
			log.Warnf("skipping '%s': %s", f, err.Error())
			continue
		}
		scanned += n
		log.Debugf("%s: %s candidate token(s)", f, tools.FormatInt(n))
	}

	out := make([]string, 0, len(words))
	for w := range words {
		out = append(out, w)
	}
	sort.Strings(out)

	if err := writeWordlist(outPath, out); err != nil {
		return err
	}

	log.Infof("Wordlist generated: %s unique label(s) (from %s scanned token(s))",
		tools.FormatInt(len(out)), tools.FormatInt(scanned))
	if outPath != "" {
		log.Infof("Saved to %s", outPath)
	}
	return nil
}

// expandWordlistInputs expands each pattern with filepath.Glob, falling back to a
// literal file check, and returns the deduplicated list of regular files. The
// resolved output path (if any) is skipped.
func expandWordlistInputs(patterns []string, outPath string) ([]string, error) {
	outAbs := ""
	if outPath != "" {
		if a, err := filepath.Abs(outPath); err == nil {
			outAbs = a
		}
	}

	seen := map[string]struct{}{}
	var files []string

	add := func(p string) {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		if outAbs != "" && abs == outAbs {
			log.Debugf("excluding output file from inputs: %s", p)
			return
		}
		if _, ok := seen[abs]; ok {
			return
		}
		fi, err := os.Stat(abs)
		if err != nil || fi.IsDir() {
			return
		}
		seen[abs] = struct{}{}
		files = append(files, abs)
	}

	for _, pat := range patterns {
		pat = strings.TrimSpace(pat)
		if pat == "" {
			continue
		}
		matches, err := filepath.Glob(pat)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern '%s': %s", pat, err.Error())
		}
		if len(matches) == 0 {
			if tools.FileExists(pat) {
				add(pat)
			} else {
				log.Warnf("no files matched '%s'", pat)
			}
			continue
		}
		for _, m := range matches {
			add(m)
		}
	}

	return files, nil
}

// buildExclusionSet turns the --exclude values into a lookup set. Each value is
// either a path to an existing file (one term per line) or a comma-separated list
// of literal terms, so both "--exclude deny.txt" and "--exclude cloud,mail" work.
func buildExclusionSet(values []string) (map[string]struct{}, error) {
	set := map[string]struct{}{}

	addTerm := func(t string) {
		t = strings.ToLower(strings.Trim(t, ". \t\r\n"))
		if t != "" {
			set[t] = struct{}{}
		}
	}

	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if tools.FileExists(v) {
			f, err := os.Open(v)
			if err != nil {
				return nil, fmt.Errorf("cannot read exclude file '%s': %s", v, err.Error())
			}
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				addTerm(sc.Text())
			}
			cerr := sc.Err()
			f.Close()
			if cerr != nil {
				return nil, cerr
			}
			continue
		}
		for _, part := range strings.Split(v, ",") {
			addTerm(part)
		}
	}

	return set, nil
}

// collectTokens scans one file, tokenizes every line into DNS-label candidates,
// applies the filters and stores the survivors in out. It returns the number of
// candidate tokens scanned (before filtering/dedup).
func collectTokens(path string, exclude map[string]struct{}, out map[string]struct{}) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanned := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.ToLower(sc.Text())
		for _, run := range fqdnRun.FindAllString(line, -1) {
			labels := strings.Split(run, ".")
			// Only a dotted name has a TLD position. Drop the genuine trailing
			// TLD (the "br" of ed.dev.br, the "com" of www.example.com) but keep
			// every internal label, even one that merely shares its name with a
			// gTLD (e.g. "dev"). Bare tokens have no dot, so they are never
			// TLD-filtered. --keep-tld keeps the trailing TLD too.
			if !wordlistOpts.KeepTLD && len(labels) > 1 {
				if last := len(labels) - 1; tools.IsTLD(labels[last]) {
					labels = labels[:last]
				}
			}
			for _, tok := range labels {
				scanned++
				if acceptToken(tok, exclude) {
					out[tok] = struct{}{}
				}
			}
		}
	}
	return scanned, sc.Err()
}

// acceptToken applies every filter to a single token.
func acceptToken(tok string, exclude map[string]struct{}) bool {
	if len(tok) < wordlistOpts.MinLength || len(tok) > wordlistOpts.MaxLength {
		return false
	}
	if !validLabel.MatchString(tok) {
		return false // flag-like junk: "--allow-parent-soa", "-foo", "bar-"
	}
	if !wordlistOpts.KeepNumeric && allNumeric.MatchString(tok) {
		return false // IP octets, TTLs, SOA serials: 144, 226, 300, 3600, ...
	}
	if _, bad := exclude[tok]; bad {
		return false
	}
	return true
}

// writeWordlist writes the sorted words to path, or to stdout when path is empty.
func writeWordlist(path string, words []string) error {
	if path == "" {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		for _, word := range words {
			fmt.Fprintln(w, word)
		}
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, word := range words {
		if _, err := fmt.Fprintln(w, word); err != nil {
			return err
		}
	}
	return w.Flush()
}

func init() {
	rootCmd.AddCommand(wordlistCmd)

	wordlistCmd.Flags().StringArrayVarP(&wordlistOpts.Inputs, "input", "i", []string{}, "Input file or glob pattern (repeatable). Quote globs (e.g. -i '*.txt') to let enumdns expand them.")
	wordlistCmd.Flags().StringArrayVar(&wordlistOpts.Exclude, "exclude", []string{}, "Terms to exclude: a file (one term per line) or comma-separated literals (repeatable), e.g. --exclude cloud,mail,onmicrosoft or --exclude deny.txt")
	wordlistCmd.Flags().IntVar(&wordlistOpts.MinLength, "min-length", 4, "Minimum label length to keep")
	wordlistCmd.Flags().IntVar(&wordlistOpts.MaxLength, "max-length", 63, "Maximum label length to keep (DNS label max is 63)")
	wordlistCmd.Flags().BoolVar(&wordlistOpts.KeepTLD, "keep-tld", false, "Keep the trailing TLD of dotted names (e.g. com, br) instead of removing it")
	wordlistCmd.Flags().BoolVar(&wordlistOpts.KeepNumeric, "keep-numeric", false, "Keep all-numeric tokens (e.g. 144, 3600) instead of dropping them")
}
