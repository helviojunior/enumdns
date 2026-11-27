# Saídas & Writers

[🏠 Home](Home) · [🇺🇸 English](Output-and-Writers-EN)

Todo comando de varredura (`recon`, `brute`, `resolve`, `threat-analysis`, `ad`)
envia seus resultados para um ou mais **writers**. Você pode habilitar quantos
quiser ao mesmo tempo — a mesma execução pode gravar um arquivo texto, um banco
SQLite e JSON Lines simultaneamente.

## Writers / formatos de saída

| Flag | Arquivo padrão | Descrição |
|------|----------------|-----------|
| `-o, --write-text-file` | — | Relatório texto legível. |
| `--write-jsonl` | `enumdns.jsonl` | JSON Lines (um objeto JSON por registro). |
| `--write-jsonl-file` | — | Sobrescreve o nome do arquivo JSONL. |
| `--write-csv` | `enumdns.csv` | CSV (conjunto limitado de colunas). |
| `--write-csv-file` | — | Sobrescreve o nome do arquivo CSV. |
| `--write-db` | `sqlite://enumdns.sqlite3` | Banco de dados relacional. |
| `--write-db-uri` | — | URI do banco — SQLite, Postgres ou MySQL (ex.: `postgres://user:pass@host:port/db`). |
| `--write-db-enable-debug` | — | Loga queries do banco (verboso). |
| `--write-elastic` | `http://localhost:9200/enumdns` | Grava no Elasticsearch ao vivo. |
| `--write-elasticsearch-uri` | — | URI do Elasticsearch (`http://user:senha@host:9200/index`). |
| `--write-none` | — | Writer vazio, para silenciar o aviso de "nenhum writer configurado". |

```bash
# Vários writers de uma vez
enumdns recon -d test.com -o out.txt --write-jsonl --write-db

# Para Postgres
enumdns brute -d test.com -w wl.txt --write-db --write-db-uri 'postgres://user:pass@db:5432/enumdns'

# Elasticsearch ao vivo
enumdns recon -d test.com --write-elastic --write-elasticsearch-uri http://localhost:9200/enumdns
```

Se nenhum writer for configurado, o EnumDNS exibe um aviso e os resultados só
aparecem na tela. Para persisti-los, adicione uma ou mais flags `--write-*`. Para
converter ou enviar resultados existentes depois, veja [report](Reports-PT).

## Vínculo de SOA

Em todos os modos de enumeração, o EnumDNS resolve o registro **SOA (Start Of
Authority)** da zona a que cada host pertence e vincula cada registro a ele.

O SOA é armazenado como um objeto dedicado (tabela `soa` / índice Elasticsearch
`<índice>-soa`) contendo o **apex da zona**, o nameserver primário (MNAME), a
caixa postal responsável (RNAME), o serial e os timers de
refresh/retry/expire/min-TTL. Cada resultado carrega uma coluna `soa` apontando
para o apex da sua zona, de modo que o vínculo é preservado em todos os formatos
de saída (DB, CSV, JSON, Elastic).

### Busca com cache

O SOA de uma zona é consultado **no máximo uma vez**, não importa quantos hosts
daquela zona você enumere:

- Para `a.b.c.example.com`, o cache é checado do mais específico ao menos
  específico (`a.b.c.example.com` → `b.c.example.com` → `c.example.com` →
  `example.com`), parando antes do TLD puro.
- Em um **hit**, o SOA cacheado é reutilizado sem consulta DNS.
- Em um **miss**, uma única consulta SOA revela o apex real da zona, que então é
  cacheado para todos os hosts seguintes daquela zona.

Isso mantém varreduras grandes rápidas, ainda atribuindo cada registro à sua
zona autoritativa. A flag `--allow-parent-soa` controla o comportamento para
nomes que não são o apex de uma zona (veja
[Opções Globais](Global-Options-PT#escopo-de-zona--public-suffix)).

## Banco de controle

O EnumDNS mantém um **banco de controle** em `~/.enumdns.db` que registra quais
hosts já foram testados. Reexecutar um comando então **ignora** os hosts já
vistos, permitindo retomar ou estender uma varredura de forma barata.

- `--force` reverifica tudo, ignorando o banco de controle.
- `--disable-control-db` o desativa por completo (sem dedup entre execuções); um
  banco temporário é usado na execução.
- `--local-workspace` armazena o banco de controle (`enumdns_ctrl.db`) no
  diretório atual em vez do `$HOME`.

O banco de controle é separado dos writers de resultado — ele é controle
interno, não saída.

## Relacionados

- [report](Reports-PT) — converte entre formatos e sincroniza com o
  Elasticsearch.
- [Opções Globais](Global-Options-PT) — as flags acima são compartilhadas por
  todos os comandos.
