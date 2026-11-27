# threat-analysis

[🏠 Home](Home) · [🇧🇷 Português](Threat-Analysis-PT)

`enumdns threat-analysis` discovers **malicious look-alike domains** that could
be used to attack your organization. It generates domain variations with several
techniques, then probes each one to see which actually exist.

## Usage

```bash
enumdns threat-analysis -d <domain> [techniques] [writers] [global flags]
enumdns threat-analysis -L <domains.txt> [techniques] [writers] [global flags]
```

If you select no technique, **all techniques** are enabled automatically.

## Detection techniques

| Technique | Flag | What it catches |
|-----------|------|-----------------|
| Typosquatting | `--typosquatting` | Keyboard-adjacency typos (`goggle.com` for `google.com`) |
| Bitsquatting | `--bitsquatting` | Single bit-flip errors |
| Homographic | `--homographic` | Unicode look-alikes (`рaypal.com` with Cyrillic `р`) |
| Insertion / deletion | *(all)* | Added or removed characters |
| Transposition | *(all)* | Swapped adjacent characters |
| TLD variation | *(all)* | Suspicious TLD swaps (`.tk`, `.ml`, `.com.br`, …) |
| Suffix impersonation | *(all)* | Impersonating a multi-label suffix (e.g. `g0v.br`) |
| Subdomain / brand patterns | `--brand-combo` | `secure-`, `login-`, `verify-` style patterns |

| Flag | Description |
|------|-------------|
| `--all-techniques` | Enable every technique at once. |

## Scope & tuning flags

| Flag | Default | Description |
|------|---------|-------------|
| `-d, --domain` | — | Target domain (mutually exclusive with `-L`). |
| `-L, --domain-list` | — | File of domains, one per line. |
| `--max-variations` | `1000` | Max variations per domain (10–50000). Warns if very high/low. |
| `--target-tlds` | `com,net,org,co,io,com.br,net.br,org.br` | TLDs to try for TLD-swap variations (comma-separated). |
| `--focus-suffix` | — | Emphasize suffix-specific techniques (e.g. `gov.br`). |
| `--span-last3` | off | Operate over the last 3 labels: mutate the 3rd-from-right, keep the last 2 as suffix. |
| `--brand-combo` | off | Add brand prefix/suffix patterns on top of eTLD+1 variations. |
| `--emit-candidates` | off | Write **all** generated candidates to outputs before probing (includes NX domains). |
| `--warn-limits` | on | Warn when the per-domain variation limit is reached. |
| `-Q, --quick` | off | Request only A records (faster). |
| `-I, --ignore-nonexistent` | off | Ignore nonexistent domains during analysis. |

### How the scope is chosen

Variations happen on the **registrable domain** (per the Public Suffix List);
subdomains to the left are preserved.

- `microsoft.com` → vary `microsoft.*`
- `recife.pe.gov.br` → vary `pe.gov.br` and the suffix `gov.br` (never touches
  `recife`)

TLD swaps use the union of the real suffix plus `--target-tlds`. Use
`--span-last3` for tricky three-label cases such as `updates.microsoft.com`.

## Examples

```bash
# All techniques, text report
enumdns threat-analysis -d example.com --all-techniques -o threats.txt

# Only typosquatting + homographic, into SQLite
enumdns threat-analysis -d example.com --typosquatting --homographic --write-db

# Many domains from a file, higher volume, JSON Lines
enumdns threat-analysis -L domains.txt --all-techniques --max-variations 5000 --write-jsonl

# gov.br focus, emit every candidate (including NX)
enumdns threat-analysis -d recife.pe.gov.br --all-techniques --focus-suffix=gov.br --emit-candidates -o gov-br.txt

# com.br with a broader TLD-swap set
enumdns threat-analysis -d yeslinux.com.br --all-techniques --target-tlds com,net,org,co,info,io,com.br,net.br,org.br

# Three-label host, mutate the 3rd-from-right label
enumdns threat-analysis -d updates.microsoft.com --all-techniques --span-last3
```

## Output

Each existing look-alike domain is written to your configured
[writers](Output-and-Writers-EN). The text writer deduplicates lines; use
`--emit-candidates` to also record the generated candidates (useful to audit
coverage, even for domains that do not resolve). Already-analyzed domains are
skipped via the control database — use `--force` to re-analyze.

## Related

- [Global Options](Global-Options-EN) — proxy, threads, timeout, resolver.
- [Output & Writers](Output-and-Writers-EN) — where results go.
