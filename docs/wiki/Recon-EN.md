# recon

[🏠 Home](Home) · [🇧🇷 Português](Recon-PT)

`enumdns recon` performs **DNS record enumeration** for one or more domains. For
each name it queries the common record types — **A, AAAA, CNAME, NS, MX, TXT,
SOA, SRV, PTR** — identifies cloud-provider services, resolves the zone **SOA**
and links every record back to its zone.

Recon also follows newly discovered domains: when enumeration surfaces a new
domain (e.g. via CNAME/NS targets), that domain is queued and enumerated too,
until no new domains appear.

## Usage

```bash
enumdns recon -d <domain> [writers] [global flags]
enumdns recon -L <domains.txt> [writers] [global flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-d, --dns-name` | A single DNS name/suffix to enumerate (e.g. `test.com`). |
| `-L, --dns-list` | A file with one DNS name per line. |
| `-I, --IgnoreNonexistent` | Ignore non-existent entries from `--dns-list` (skip names with no SOA instead of aborting). |

Combine with any [global flag](Global-Options-EN): `-s` (resolver), `-X`
(proxy), `-t/-T` (threads/timeout), and the output flags below.

## Examples

```bash
# Enumerate a single domain, write a text report
enumdns recon -d test.com -o enumdns.txt

# Enumerate and store as JSON Lines
enumdns recon -d test.com --write-jsonl

# Enumerate a list of domains into a SQLite database
enumdns recon -L domains.txt --write-db

# Internal / AD zone through a specific domain controller
enumdns recon -d corp.local -s 10.0.0.10 --write-db

# A non-apex name — resolve the real parent SOA instead of failing
enumdns recon -d www.test.com --allow-parent-soa -o out.txt
```

## How it works

1. The target name is validated and its **zone SOA** is resolved (cached per
   zone, so each zone is queried at most once — see
   [Output & Writers → SOA linking](Output-and-Writers-EN#soa-linking)).
2. Record types are queried for the name; results are handed to every configured
   [writer](Output-and-Writers-EN).
3. Any new domain discovered during enumeration is queued and the loop repeats
   until no new domains remain.

By default a public suffix / TLD is **refused** (use `--allow-tld` to override),
and a non-apex name requires `--allow-parent-soa` to resolve its parent zone.

## Related

- [Global Options](Global-Options-EN) — resolver, proxy, scoping, performance.
- [brute](Brute-Force-EN) — add hostname brute-forcing on top of recon.
- [resolve](Resolve-EN) — resolve hosts from files, BloodHound or crt.sh.
- [Output & Writers](Output-and-Writers-EN) — choose your output formats.
