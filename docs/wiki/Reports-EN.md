# report

[🏠 Home](Home) · [🇧🇷 Português](Reports-PT)

`enumdns report` works with the result files produced by the scan commands. It
has two subcommands: **convert** (between formats) and **elastic** (sync to
Elasticsearch).

---

## report convert

Convert results between **SQLite** and **JSON Lines** (and export to text). The
file extensions decide the conversion direction and target.

```bash
enumdns report convert --from-file <in> --to-file <out>
```

| Flag | Description |
|------|-------------|
| `--from-file` | Source file to convert from. |
| `--to-file` | Destination file. Use `.sqlite3` to produce SQLite, `.jsonl` for JSON Lines, `.txt` for text. |

```bash
# SQLite -> text report
enumdns report convert --from-file enumdns.sqlite3 --to-file enumdns.txt

# SQLite -> JSON Lines
enumdns report convert --from-file enumdns.sqlite3 --to-file data.jsonl

# JSON Lines -> SQLite
enumdns report convert --from-file enumdns.jsonl --to-file db.sqlite3
```

---

## report elastic

Sync an existing local **SQLite** or **JSON Lines** result file into an
**Elasticsearch** index (useful for dashboards/searching after an offline scan).

```bash
enumdns report elastic --from-file <in> --elasticsearch-uri <uri>
```

| Flag | Default | Description |
|------|---------|-------------|
| `--from-file` | `~/.enumdns.db` | Source file. Use `.sqlite3` (SQLite) or `.jsonl` (JSON Lines). |
| `--elasticsearch-uri` | `http://localhost:9200/enumdns` | Elasticsearch URI (e.g. `http://user:pass@host:9200/index`). |

```bash
# From SQLite
enumdns report elastic --from-file enumdns.sqlite3 --elasticsearch-uri http://localhost:9200/enumdns

# From JSON Lines
enumdns report elastic --from-file enumdns.jsonl --elasticsearch-uri http://localhost:9200/enumdns
```

> You can also write to Elasticsearch **live** during a scan with
> `--write-elastic` / `--write-elasticsearch-uri` — see
> [Output & Writers](Output-and-Writers-EN). The SOA objects are stored in a
> companion index (`<index>-soa`).

## Related

- [Output & Writers](Output-and-Writers-EN) — the formats these commands read
  and write.
- [recon](Recon-EN) · [brute](Brute-Force-EN) · [resolve](Resolve-EN)
