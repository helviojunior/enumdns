# Global Options

[🏠 Home](Home) · [🇧🇷 Português](Global-Options-PT)

These flags are **persistent** — they work with every command (`recon`, `brute`,
`resolve`, `threat-analysis`, `ad`, `wordlist`, `report`). Output-related flags
are documented in detail on [Output & Writers](Output-and-Writers-EN).

## DNS resolver

| Flag | Default | Description |
|------|---------|-------------|
| `-s, --server` | system default | DNS server to query. In AD/internal scenarios, point this at a domain controller. |
| `--port` | `53` | DNS server port. |
| `--protocol` | `UDP` | Transport protocol: `UDP` or `TCP`. |

When `--server` is omitted, EnumDNS uses the host's default resolver. A private
(RFC 1918) server is detected automatically and logged as `(private)`.

```bash
# Query a specific resolver over TCP
enumdns recon -d example.com -s 1.1.1.1 --protocol TCP

# Use an internal AD DNS (domain controller)
enumdns recon -d corp.local -s 10.0.0.10
```

> Each run does a best-effort connectivity probe against `google.com`. If your
> resolver is internal-only (does not recurse to the internet), the probe fails
> with a **warning** and the run continues — this is expected.

## Proxy (SOCKS / HTTP)

| Flag | Description |
|------|-------------|
| `-X, --proxy` | Route DNS traffic through a proxy: `<scheme://[user:pass@]host:port>`. Supports `socks4`, `socks5` and `http`. |

```bash
enumdns recon -d example.com -X socks5://127.0.0.1:1080
enumdns recon -d example.com -X socks5://user:pass@127.0.0.1:1080
```

When a proxy is set, DNS queries use the configured `--timeout` for both the
proxy dial and the query, so an unreachable proxy fails fast.

## Performance & reliability

| Flag | Default | Description |
|------|---------|-------------|
| `-t, --threads` | `6` | Concurrent workers (goroutines). |
| `-T, --timeout` | `60` | Seconds to wait for each DNS query (direct or over proxy) before timing out. |
| `-F, --force` | off | Re-check hosts even if the control DB says they were already scanned. |

```bash
# Faster scan with more workers and a tighter timeout
enumdns brute -d example.com -w wl.txt -t 20 -T 10
```

## Zone / public-suffix scoping

| Flag | Description |
|------|-------------|
| `--allow-parent-soa` | When a name is not a zone apex (e.g. `www.test.com`), resolve its parent zone's real **SOA** instead of failing. |
| `--allow-tld` | Explicitly allow enumerating a public suffix / TLD (e.g. `com.br`, `co.uk`). Refused by default as a safety guard. |

See [Output & Writers → SOA linking](Output-and-Writers-EN#soa-linking) for how
the SOA is resolved and attached to every record.

## Logging

| Flag | Description |
|------|-------------|
| `-D, --debug-log` | Enable verbose debug logging. |
| `-q, --quiet` | Silence (almost) all logging — useful for scripting. |
| `--log-scan-errors` | Log per-query scan errors (timeouts, DNS errors) to stderr. Can be very verbose. |

## Workspace / control database

| Flag | Default | Description |
|------|---------|-------------|
| `--disable-control-db` | off | Do not use the `~/.enumdns.db` control database (no cross-run dedup). |
| `--local-workspace` | off | Store the control DB (`enumdns_ctrl.db`) in the current directory instead of `$HOME`. |
| `--local-temp` | off | Store temp files under the execution path instead of the system temp. |

The **control database** (`~/.enumdns.db`) remembers which hosts were already
probed so re-running a command skips them; use `--force` to re-scan. See
[Output & Writers](Output-and-Writers-EN#control-database).

## Output (summary)

`-o/--write-text-file`, `--write-db`, `--write-jsonl`, `--write-csv`,
`--write-elastic` and `--write-none` select where results go. Full details on
[Output & Writers](Output-and-Writers-EN).

## Related

- [Output & Writers](Output-and-Writers-EN) — every output format in depth.
- [recon](Recon-EN) · [brute](Brute-Force-EN) · [resolve](Resolve-EN)
