# Output & Writers

[🏠 Home](Home) · [🇧🇷 Português](Output-and-Writers-PT)

Every scan command (`recon`, `brute`, `resolve`, `threat-analysis`, `ad`) sends
its results to one or more **writers**. You can enable as many as you like at
once — the same run can write a text file, a SQLite DB and JSON Lines
simultaneously.

## Writers / output formats

| Flag | Default file | Description |
|------|--------------|-------------|
| `-o, --write-text-file` | — | Human-readable text report. |
| `--write-jsonl` | `enumdns.jsonl` | JSON Lines (one JSON object per record). |
| `--write-jsonl-file` | — | Override the JSONL filename. |
| `--write-csv` | `enumdns.csv` | CSV (a limited set of columns). |
| `--write-csv-file` | — | Override the CSV filename. |
| `--write-db` | `sqlite://enumdns.sqlite3` | Relational database. |
| `--write-db-uri` | — | DB URI — SQLite, Postgres or MySQL (e.g. `postgres://user:pass@host:port/db`). |
| `--write-db-enable-debug` | — | Log DB queries (verbose). |
| `--write-elastic` | `http://localhost:9200/enumdns` | Write to Elasticsearch live. |
| `--write-elasticsearch-uri` | — | Elasticsearch URI (`http://user:pass@host:9200/index`). |
| `--write-none` | — | Empty writer, to silence "no writers configured" warnings. |

```bash
# Several writers at once
enumdns recon -d test.com -o out.txt --write-jsonl --write-db

# Into Postgres
enumdns brute -d test.com -w wl.txt --write-db --write-db-uri 'postgres://user:pass@db:5432/enumdns'

# Live Elasticsearch
enumdns recon -d test.com --write-elastic --write-elasticsearch-uri http://localhost:9200/enumdns
```

If no writer is configured, EnumDNS prints a warning and results are only shown
on screen. To persist them, add one or more `--write-*` flags. To convert or
ship existing results later, see [report](Reports-EN).

## SOA linking

In every enumeration mode, EnumDNS resolves the **Start Of Authority (SOA)**
record of the zone each host belongs to and links every record back to it.

The SOA is stored as a dedicated object (table `soa` / Elasticsearch index
`<index>-soa`) holding the **zone apex**, primary nameserver (MNAME),
responsible mailbox (RNAME), serial and the refresh/retry/expire/min-TTL timers.
Each result carries a `soa` column pointing to its zone apex, so the link is
preserved in every output format (DB, CSV, JSON, Elastic).

### Cache-aware lookup

A zone's SOA is queried **at most once**, no matter how many hosts of that zone
you enumerate:

- For `a.b.c.example.com`, the cache is checked from most specific to least
  specific (`a.b.c.example.com` → `b.c.example.com` → `c.example.com` →
  `example.com`), stopping before the bare TLD.
- On a **hit**, the cached SOA is reused with no DNS query.
- On a **miss**, a single SOA query reveals the real zone apex, which is then
  cached for every subsequent host of that zone.

This keeps large scans fast while still attributing every record to its
authoritative zone. The `--allow-parent-soa` flag controls behaviour for names
that are not a zone apex (see [Global Options](Global-Options-EN#zone--public-suffix-scoping)).

## Control database

EnumDNS keeps a **control database** at `~/.enumdns.db` that records which hosts
were already probed. Re-running a command then **skips** hosts it has already
seen, so you can resume or extend a scan cheaply.

- `--force` re-checks everything, ignoring the control DB.
- `--disable-control-db` turns it off entirely (no cross-run dedup); a temporary
  DB is used for the run.
- `--local-workspace` stores the control DB (`enumdns_ctrl.db`) in the current
  directory instead of `$HOME`.

The control DB is separate from your result writers — it is bookkeeping, not
output.

## Related

- [report](Reports-EN) — convert between formats and sync to Elasticsearch.
- [Global Options](Global-Options-EN) — the flags above are shared by all
  commands.
