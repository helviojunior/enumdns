# wordlist

[🏠 Home](Home) · [🇧🇷 Português](Wordlist-PT)

`enumdns wordlist` builds a **clean, deduplicated, sorted DNS-label wordlist**
from one or more input files (plain text, previous results, scraped data, …).
The result is ideal as the `-w` input for [brute](Brute-Force-EN).

## Usage

```bash
enumdns wordlist -i '<glob|file>' [-i ...] [filters] -o <out.txt>
```

Without `-o`, the wordlist is printed to stdout.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i, --input` | — | Input file or glob pattern (repeatable). **Quote globs** (`-i '*.txt'`) so enumdns expands them itself. |
| `--exclude` | — | Terms to drop: a file (one per line) or comma-separated literals (repeatable), e.g. `--exclude cloud,mail,onmicrosoft` or `--exclude deny.txt`. |
| `--min-length` | `4` | Minimum label length to keep. |
| `--max-length` | `63` | Maximum label length (DNS limit is 63). |
| `--keep-tld` | off | Keep the trailing TLD of dotted names (e.g. `com`, `br`) instead of removing it. |
| `--keep-numeric` | off | Keep all-numeric tokens (e.g. `144`, `3600`) instead of dropping them. |
| `-o, --write-text-file` | — | Output file (global flag). Always excluded from inputs. |

## How tokens are processed

Each line is split into candidate labels on anything that is not a letter,
digit, underscore or hyphen — so `host.example.com` yields `host` and `example`.
Each token is then:

- lower-cased and **deduplicated**;
- **validated** as a DNS label (flag-like junk such as `--allow-parent-soa` is
  dropped);
- dropped when **all-numeric** (IP octets, TTLs, SOA serials) unless
  `--keep-numeric`;
- filtered by **length** (`--min-length` / `--max-length`);
- filtered against your **exclusion list** (`--exclude`).

Only the trailing TLD of a *dotted* name is removed (the `com` of
`www.example.com`), so an internal label that merely shares a name with a gTLD
(e.g. `dev`) is kept. A bare token (no dot) is never treated as a TLD. Use
`--keep-tld` to keep the TLD too.

The output file (`-o`) is always removed from the inputs, so regenerating the
list never feeds back on itself.

## Examples

```bash
# All .txt files in the current dir -> a clean wordlist
enumdns wordlist -i '*.txt' -o custom_wl.txt

# Multiple inputs + inline exclusions
enumdns wordlist -i '/data/recon/*.txt' -i extra.txt --exclude cloud,mail,onmicrosoft -o custom_wl.txt

# One file, exclusion file, min length 3, keep the TLD
enumdns wordlist -i dump.txt --exclude deny.txt --min-length 3 --keep-tld -o custom_wl.txt

# Feed it straight into brute
enumdns wordlist -i '*.txt' -o wl.txt
enumdns brute -d test.com -w wl.txt --write-db
```

## Related

- [brute](Brute-Force-EN) — consume the wordlist as `-w`.
- [Global Options](Global-Options-EN)
