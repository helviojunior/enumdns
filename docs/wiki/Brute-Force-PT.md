# brute

[🏠 Home](Home) · [🇺🇸 English](Brute-Force-EN)

O `enumdns brute` faz **descoberta de hostnames por força-bruta**: combina cada
label de uma wordlist com um ou mais sufixos DNS (`<palavra>.<sufixo>`) e resolve
cada candidato. Os resultados são vinculados ao SOA da zona, assim como no
[recon](Recon-PT).

## Uso

```bash
enumdns brute -d <sufixo> -w <wordlist.txt> [writers] [flags globais]
enumdns brute -L <sufixos.txt> -w <wordlist.txt> [writers] [flags globais]
```

## Flags

| Flag | Descrição |
|------|-----------|
| `-d, --dns-suffix` | Um único sufixo DNS para a força-bruta (ex.: `test.com`). |
| `-L, --dns-list` | Arquivo com um sufixo DNS por linha. |
| `-w, --word-list` | **Obrigatório.** Arquivo de wordlist, um label por linha. |
| `-Q, --quick` | Consulta **apenas registros A** (mais rápido). |
| `-I, --IgnoreNonexistent` | Ignora sufixos inexistentes de `--dns-list`. |

Mais qualquer [flag global](Global-Options-PT) — em especial `-t`/`-T` para
ajustar a velocidade, e as flags de saída abaixo.

## Exemplos

```bash
# Força-bruta de um domínio com uma wordlist, saída texto
enumdns brute -d test.com -w /tmp/wordlist.txt -o enumdns.txt

# Armazenar resultados como JSON Lines
enumdns brute -d test.com -w /tmp/wordlist.txt --write-jsonl

# Múltiplos sufixos de um arquivo, para SQLite
enumdns brute -L domains.txt -w /tmp/wordlist.txt --write-db

# Varredura rápida (só A), mais workers
enumdns brute -d test.com -w big-wordlist.txt -Q -t 30 -T 5 --write-db
```

## Comportamento & escopo

- O trabalho total é `nº de sufixos × nº de palavras`. Hosts já escaneados são
  **ignorados** usando o [banco de controle](Output-and-Writers-PT#banco-de-controle);
  passe `--force` para reescanear tudo.
- Public suffixes / TLDs são **recusados** por padrão — use `--allow-tld` para
  permiti-los.
- Com `--allow-parent-soa`, o apex real da zona é resolvido e a wordlist é
  aplicada tanto ao apex quanto ao nome informado.

## Montando uma boa wordlist

Use o comando [wordlist](Wordlist-PT) para transformar resultados anteriores,
dados coletados ou arquivos diversos em uma wordlist de labels DNS limpa e
deduplicada:

```bash
enumdns wordlist -i '*.txt' --exclude cloud,mail -o custom_wl.txt
enumdns brute -d test.com -w custom_wl.txt --write-db
```

## Relacionados

- [wordlist](Wordlist-PT) — gera a wordlist do `-w`.
- [recon](Recon-PT) — enumeração de registros sem força-bruta.
- [Opções Globais](Global-Options-PT) · [Saídas & Writers](Output-and-Writers-PT)
