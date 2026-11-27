# report

[🏠 Home](Home) · [🇺🇸 English](Reports-EN)

O `enumdns report` trabalha com os arquivos de resultado produzidos pelos
comandos de varredura. Possui dois subcomandos: **convert** (entre formatos) e
**elastic** (sincronização com o Elasticsearch).

---

## report convert

Converte resultados entre **SQLite** e **JSON Lines** (e exporta para texto). As
extensões dos arquivos definem a direção e o alvo da conversão.

```bash
enumdns report convert --from-file <entrada> --to-file <saida>
```

| Flag | Descrição |
|------|-----------|
| `--from-file` | Arquivo de origem a converter. |
| `--to-file` | Arquivo de destino. Use `.sqlite3` para gerar SQLite, `.jsonl` para JSON Lines, `.txt` para texto. |

```bash
# SQLite -> relatório texto
enumdns report convert --from-file enumdns.sqlite3 --to-file enumdns.txt

# SQLite -> JSON Lines
enumdns report convert --from-file enumdns.sqlite3 --to-file data.jsonl

# JSON Lines -> SQLite
enumdns report convert --from-file enumdns.jsonl --to-file db.sqlite3
```

---

## report elastic

Sincroniza um arquivo local existente de resultados em **SQLite** ou **JSON
Lines** para um índice **Elasticsearch** (útil para dashboards/busca após uma
varredura offline).

```bash
enumdns report elastic --from-file <entrada> --elasticsearch-uri <uri>
```

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `--from-file` | `~/.enumdns.db` | Arquivo de origem. Use `.sqlite3` (SQLite) ou `.jsonl` (JSON Lines). |
| `--elasticsearch-uri` | `http://localhost:9200/enumdns` | URI do Elasticsearch (ex.: `http://user:senha@host:9200/index`). |

```bash
# A partir de SQLite
enumdns report elastic --from-file enumdns.sqlite3 --elasticsearch-uri http://localhost:9200/enumdns

# A partir de JSON Lines
enumdns report elastic --from-file enumdns.jsonl --elasticsearch-uri http://localhost:9200/enumdns
```

> Você também pode gravar no Elasticsearch **ao vivo** durante uma varredura com
> `--write-elastic` / `--write-elasticsearch-uri` — veja
> [Saídas & Writers](Output-and-Writers-PT). Os objetos SOA são armazenados em um
> índice complementar (`<índice>-soa`).

## Relacionados

- [Saídas & Writers](Output-and-Writers-PT) — os formatos que estes comandos leem
  e escrevem.
- [recon](Recon-PT) · [brute](Brute-Force-PT) · [resolve](Resolve-PT)
