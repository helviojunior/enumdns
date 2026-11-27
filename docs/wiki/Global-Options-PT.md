# Opções Globais

[🏠 Home](Home) · [🇺🇸 English](Global-Options-EN)

Estas flags são **persistentes** — funcionam com todos os comandos (`recon`,
`brute`, `resolve`, `threat-analysis`, `ad`, `wordlist`, `report`). As flags de
saída são detalhadas em [Saídas & Writers](Output-and-Writers-PT).

## Resolver DNS

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-s, --server` | resolver do sistema | Servidor DNS a consultar. Em cenários de AD/interno, aponte para um domain controller. |
| `--port` | `53` | Porta do servidor DNS. |
| `--protocol` | `UDP` | Protocolo de transporte: `UDP` ou `TCP`. |

Quando `--server` é omitido, o EnumDNS usa o resolver padrão do host. Um servidor
privado (RFC 1918) é detectado automaticamente e registrado como `(private)`.

```bash
# Consultar um resolver específico via TCP
enumdns recon -d example.com -s 1.1.1.1 --protocol TCP

# Usar um DNS interno de AD (domain controller)
enumdns recon -d corp.local -s 10.0.0.10
```

> Cada execução faz uma verificação de conectividade (best-effort) contra o
> `google.com`. Se o seu resolver for apenas interno (não recursa para a
> internet), a verificação falha com um **aviso** e a execução continua — isso é
> esperado.

## Proxy (SOCKS / HTTP)

| Flag | Descrição |
|------|-----------|
| `-X, --proxy` | Encaminha o tráfego DNS por um proxy: `<esquema://[user:senha@]host:porta>`. Suporta `socks4`, `socks5` e `http`. |

```bash
enumdns recon -d example.com -X socks5://127.0.0.1:1080
enumdns recon -d example.com -X socks5://user:senha@127.0.0.1:1080
```

Com um proxy configurado, as consultas DNS usam o `--timeout` tanto na conexão
com o proxy quanto na consulta, de modo que um proxy inacessível falha rápido.

## Desempenho & confiabilidade

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-t, --threads` | `6` | Workers concorrentes (goroutines). |
| `-T, --timeout` | `60` | Segundos de espera por cada consulta DNS (direta ou via proxy) antes de expirar. |
| `-F, --force` | off | Reverifica hosts mesmo que o banco de controle indique que já foram escaneados. |

```bash
# Varredura mais rápida, com mais workers e timeout menor
enumdns brute -d example.com -w wl.txt -t 20 -T 10
```

## Escopo de zona / public-suffix

| Flag | Descrição |
|------|-----------|
| `--allow-parent-soa` | Quando um nome não é o apex da zona (ex.: `www.test.com`), resolve o **SOA** real da zona pai em vez de falhar. |
| `--allow-tld` | Permite explicitamente enumerar um public suffix / TLD (ex.: `com.br`, `co.uk`). Recusado por padrão, como proteção. |

Veja [Saídas & Writers → Vínculo de SOA](Output-and-Writers-PT#vínculo-de-soa)
para entender como o SOA é resolvido e anexado a cada registro.

## Logging

| Flag | Descrição |
|------|-----------|
| `-D, --debug-log` | Habilita log de debug detalhado. |
| `-q, --quiet` | Silencia (quase) todo o log — útil para scripts. |
| `--log-scan-errors` | Loga erros de consulta (timeouts, erros DNS) no stderr. Pode ser muito verboso. |

## Workspace / banco de controle

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `--disable-control-db` | off | Não usa o banco de controle `~/.enumdns.db` (sem dedup entre execuções). |
| `--local-workspace` | off | Armazena o banco de controle (`enumdns_ctrl.db`) no diretório atual em vez do `$HOME`. |
| `--local-temp` | off | Armazena arquivos temporários no caminho de execução em vez do temp do sistema. |

O **banco de controle** (`~/.enumdns.db`) lembra quais hosts já foram testados,
de modo que reexecutar um comando os ignora; use `--force` para reescanear. Veja
[Saídas & Writers](Output-and-Writers-PT#banco-de-controle).

## Saída (resumo)

`-o/--write-text-file`, `--write-db`, `--write-jsonl`, `--write-csv`,
`--write-elastic` e `--write-none` definem para onde os resultados vão. Detalhes
completos em [Saídas & Writers](Output-and-Writers-PT).

## Relacionados

- [Saídas & Writers](Output-and-Writers-PT) — cada formato de saída em detalhe.
- [recon](Recon-PT) · [brute](Brute-Force-PT) · [resolve](Resolve-PT)
