# recon

[🏠 Home](Home) · [🇺🇸 English](Recon-EN)

O `enumdns recon` realiza a **enumeração de registros DNS** de um ou mais
domínios. Para cada nome, consulta os tipos de registro comuns — **A, AAAA,
CNAME, NS, MX, TXT, SOA, SRV, PTR** —, identifica serviços de provedores de
nuvem, resolve o **SOA** da zona e vincula cada registro à sua zona.

O recon também segue domínios recém-descobertos: quando a enumeração revela um
novo domínio (ex.: via alvos de CNAME/NS), esse domínio entra na fila e também é
enumerado, até que nenhum novo domínio apareça.

## Uso

```bash
enumdns recon -d <domínio> [writers] [flags globais]
enumdns recon -L <dominios.txt> [writers] [flags globais]
```

## Flags

| Flag | Descrição |
|------|-----------|
| `-d, --dns-name` | Um único nome/sufixo DNS a enumerar (ex.: `test.com`). |
| `-L, --dns-list` | Arquivo com um nome DNS por linha. |
| `-I, --IgnoreNonexistent` | Ignora entradas inexistentes de `--dns-list` (pula nomes sem SOA em vez de abortar). |

Combine com qualquer [flag global](Global-Options-PT): `-s` (resolver), `-X`
(proxy), `-t/-T` (threads/timeout) e as flags de saída abaixo.

## Exemplos

```bash
# Enumerar um único domínio e gravar um relatório texto
enumdns recon -d test.com -o enumdns.txt

# Enumerar e armazenar como JSON Lines
enumdns recon -d test.com --write-jsonl

# Enumerar uma lista de domínios para um banco SQLite
enumdns recon -L domains.txt --write-db

# Zona interna / AD através de um domain controller específico
enumdns recon -d corp.local -s 10.0.0.10 --write-db

# Um nome que não é apex — resolve o SOA real da zona pai em vez de falhar
enumdns recon -d www.test.com --allow-parent-soa -o out.txt
```

## Como funciona

1. O nome-alvo é validado e o **SOA da zona** é resolvido (cacheado por zona, de
   modo que cada zona é consultada no máximo uma vez — veja
   [Saídas & Writers → Vínculo de SOA](Output-and-Writers-PT#vínculo-de-soa)).
2. Os tipos de registro são consultados para o nome; os resultados são entregues
   a todos os [writers](Output-and-Writers-PT) configurados.
3. Qualquer novo domínio descoberto durante a enumeração entra na fila e o laço
   se repete até não restarem novos domínios.

Por padrão, um public suffix / TLD é **recusado** (use `--allow-tld` para
sobrepor) e um nome que não seja apex exige `--allow-parent-soa` para resolver a
zona pai.

## Relacionados

- [Opções Globais](Global-Options-PT) — resolver, proxy, escopo, desempenho.
- [brute](Brute-Force-PT) — adiciona força-bruta de hostnames ao recon.
- [resolve](Resolve-PT) — resolve hosts de arquivos, BloodHound ou crt.sh.
- [Saídas & Writers](Output-and-Writers-PT) — escolha os formatos de saída.
