# resolve

[🏠 Home](Home) · [🇧🇷 Português](Resolve-PT)

`enumdns resolve` resolves a set of hosts that come from an external **source**
instead of being brute-forced. It has three subcommands:

| Subcommand | Source |
|------------|--------|
| [`resolve file`](#resolve-file) | A plain text file, one host per line |
| [`resolve bloodhound`](#resolve-bloodhound) | A BloodHound export (`.zip` or `_computers.json`) |
| [`resolve crtsh`](#resolve-crtsh) | Certificate Transparency logs via crt.sh |

All three resolve records and link them to the zone SOA, writing to every
configured [writer](Output-and-Writers-EN). The shared flag
`-I/--IgnoreNonexistent` skips non-resolvable entries instead of aborting.

---

## resolve file

Resolve every host listed in a text file.

```bash
enumdns resolve file -L <host_list.txt> [writers] [global flags]
```

| Flag | Description |
|------|-------------|
| `-L, --host-list` | File with one DNS host per line. |

```bash
enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
enumdns resolve file -L /tmp/host_list.txt --write-jsonl
enumdns resolve file -L /tmp/host_list.txt --write-db
```

Already-resolved hosts are skipped via the control database; use `--force` to
re-resolve.

---

## resolve bloodhound

Resolve the computer accounts from a **BloodHound** collection. Accepts either a
raw `_computers.json` file or the collection `.zip` (the `_computers.json` is
extracted automatically). Disabled computer accounts are ignored.

```bash
enumdns resolve bloodhound -L <file.zip|_computers.json> [writers] [global flags]
```

| Flag | Description |
|------|-------------|
| `-L, --bloodhound-file` | BloodHound output: a `.zip` or a `_computers.json`. |

```bash
enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db
```

The command first checks DNS connectivity for each domain found in the export.
In an AD environment, point `-s` at a domain controller so the internal names
resolve:

```bash
enumdns resolve bloodhound -L export.zip -s 10.0.0.10 --write-db
```

---

## resolve crtsh

Crawl **Certificate Transparency** logs (crt.sh) for a domain, collect every
subdomain seen in issued certificates, then resolve them.

```bash
enumdns resolve crtsh -d <domain> [writers] [global flags]
enumdns resolve crtsh -L <domains.txt> [writers] [global flags]
```

| Flag | Description |
|------|-------------|
| `-d, --dns-suffix` | A single domain to crawl on crt.sh. |
| `-L, --dns-list` | File with one domain per line. |
| `--fqdn-out` | Also append every discovered FQDN to this output file. |

```bash
enumdns resolve crtsh -d test.com -o enumdns.txt
enumdns resolve crtsh -L domains.txt --write-jsonl
enumdns resolve crtsh -d test.com --fqdn-out fqdns.txt --write-db
```

Like `brute`/`recon`, crt.sh crawling respects `--allow-parent-soa` and
`--allow-tld`, and works through a `-X` proxy.

## Related

- [recon](Recon-EN) — enumerate records for a domain directly.
- [brute](Brute-Force-EN) — discover hosts via wordlist.
- [Global Options](Global-Options-EN) · [Output & Writers](Output-and-Writers-EN)
