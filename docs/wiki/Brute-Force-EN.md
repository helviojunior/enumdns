# brute

[🏠 Home](Home) · [🇧🇷 Português](Brute-Force-PT)

`enumdns brute` performs **brute-force hostname discovery**: it combines every
label of a wordlist with one or more DNS suffixes (`<word>.<suffix>`) and
resolves each candidate. Results are linked to the zone SOA just like
[recon](Recon-EN).

## Usage

```bash
enumdns brute -d <suffix> -w <wordlist.txt> [writers] [global flags]
enumdns brute -L <suffixes.txt> -w <wordlist.txt> [writers] [global flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-d, --dns-suffix` | A single DNS suffix to brute-force (e.g. `test.com`). |
| `-L, --dns-list` | A file with one DNS suffix per line. |
| `-w, --word-list` | **Required.** Wordlist file, one label per line. |
| `-Q, --quick` | Request **only A records** (faster). |
| `-I, --IgnoreNonexistent` | Ignore non-existent suffixes from `--dns-list`. |

Plus any [global flag](Global-Options-EN) — notably `-t`/`-T` to tune speed, and
the output flags below.

## Examples

```bash
# Brute-force one domain with a wordlist, text output
enumdns brute -d test.com -w /tmp/wordlist.txt -o enumdns.txt

# Store results as JSON Lines
enumdns brute -d test.com -w /tmp/wordlist.txt --write-jsonl

# Multiple suffixes from a file, into SQLite
enumdns brute -L domains.txt -w /tmp/wordlist.txt --write-db

# Fast A-only scan, more workers
enumdns brute -d test.com -w big-wordlist.txt -Q -t 30 -T 5 --write-db
```

## Behaviour & scoping

- The total work is `len(suffixes) × len(wordlist)`. Already-scanned hosts are
  **skipped** using the [control database](Output-and-Writers-EN#control-database);
  pass `--force` to re-scan everything.
- Public suffixes / TLDs are **refused** by default — use `--allow-tld` to allow
  them.
- With `--allow-parent-soa`, the real zone apex is resolved and the wordlist is
  brute-forced against both the apex and the name you supplied.

## Building a good wordlist

Use the [wordlist](Wordlist-EN) command to turn previous results, scraped data
or arbitrary files into a clean, deduplicated DNS-label wordlist:

```bash
enumdns wordlist -i '*.txt' --exclude cloud,mail -o custom_wl.txt
enumdns brute -d test.com -w custom_wl.txt --write-db
```

## Related

- [wordlist](Wordlist-EN) — generate the `-w` wordlist.
- [recon](Recon-EN) — record enumeration without brute-forcing.
- [Global Options](Global-Options-EN) · [Output & Writers](Output-and-Writers-EN)
