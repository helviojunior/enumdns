# threat-analysis

[🏠 Home](Home) · [🇺🇸 English](Threat-Analysis-EN)

O `enumdns threat-analysis` descobre **domínios parecidos (look-alike)
maliciosos** que poderiam ser usados em ataques contra a sua organização. Ele
gera variações de domínio com diversas técnicas e então testa cada uma para ver
quais realmente existem.

## Uso

```bash
enumdns threat-analysis -d <domínio> [técnicas] [writers] [flags globais]
enumdns threat-analysis -L <dominios.txt> [técnicas] [writers] [flags globais]
```

Se você não selecionar nenhuma técnica, **todas** são habilitadas
automaticamente.

## Técnicas de detecção

| Técnica | Flag | O que detecta |
|---------|------|---------------|
| Typosquatting | `--typosquatting` | Erros de adjacência de teclado (`goggle.com` por `google.com`) |
| Bitsquatting | `--bitsquatting` | Erros de inversão de um único bit |
| Homográfica | `--homographic` | Caracteres Unicode parecidos (`рaypal.com` com `р` cirílico) |
| Inserção / remoção | *(todas)* | Caracteres adicionados ou removidos |
| Transposição | *(todas)* | Troca de caracteres adjacentes |
| Variação de TLD | *(todas)* | Trocas suspeitas de TLD (`.tk`, `.ml`, `.com.br`, …) |
| Falsificação de sufixo | *(todas)* | Impersonar um sufixo multi-label (ex.: `g0v.br`) |
| Padrões de subdomínio / marca | `--brand-combo` | Padrões do tipo `secure-`, `login-`, `verify-` |

| Flag | Descrição |
|------|-----------|
| `--all-techniques` | Habilita todas as técnicas de uma vez. |

## Flags de escopo & ajuste

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-d, --domain` | — | Domínio-alvo (mutuamente exclusivo com `-L`). |
| `-L, --domain-list` | — | Arquivo de domínios, um por linha. |
| `--max-variations` | `1000` | Máximo de variações por domínio (10–50000). Avisa se muito alto/baixo. |
| `--target-tlds` | `com,net,org,co,io,com.br,net.br,org.br` | TLDs a testar nas trocas de TLD (separados por vírgula). |
| `--focus-suffix` | — | Enfatiza técnicas específicas de sufixo (ex.: `gov.br`). |
| `--span-last3` | off | Opera sobre os últimos 3 labels: muta o 3º da direita, mantém os 2 últimos como sufixo. |
| `--brand-combo` | off | Adiciona padrões de prefixo/sufixo de marca às variações de eTLD+1. |
| `--emit-candidates` | off | Grava **todos** os candidatos gerados antes de testar (inclui domínios NX). |
| `--warn-limits` | on | Avisa quando o limite de variações por domínio é atingido. |
| `-Q, --quick` | off | Consulta apenas registros A (mais rápido). |
| `-I, --ignore-nonexistent` | off | Ignora domínios inexistentes durante a análise. |

### Como o escopo é escolhido

As variações ocorrem no **domínio registrável** (conforme a Public Suffix List);
os subdomínios à esquerda são preservados.

- `microsoft.com` → varia `microsoft.*`
- `recife.pe.gov.br` → varia `pe.gov.br` e o sufixo `gov.br` (nunca toca em
  `recife`)

As trocas de TLD usam a união do sufixo real com `--target-tlds`. Use
`--span-last3` para casos difíceis de três labels, como `updates.microsoft.com`.

## Exemplos

```bash
# Todas as técnicas, relatório texto
enumdns threat-analysis -d example.com --all-techniques -o threats.txt

# Apenas typosquatting + homográfica, para SQLite
enumdns threat-analysis -d example.com --typosquatting --homographic --write-db

# Vários domínios de um arquivo, volume maior, JSON Lines
enumdns threat-analysis -L domains.txt --all-techniques --max-variations 5000 --write-jsonl

# Foco em gov.br, emitindo todos os candidatos (inclui NX)
enumdns threat-analysis -d recife.pe.gov.br --all-techniques --focus-suffix=gov.br --emit-candidates -o gov-br.txt

# com.br com um conjunto mais amplo de trocas de TLD
enumdns threat-analysis -d yeslinux.com.br --all-techniques --target-tlds com,net,org,co,info,io,com.br,net.br,org.br

# Host de três labels, mutando o 3º label da direita
enumdns threat-analysis -d updates.microsoft.com --all-techniques --span-last3
```

## Saída

Cada domínio look-alike existente é gravado nos [writers](Output-and-Writers-PT)
configurados. O writer de texto deduplica as linhas; use `--emit-candidates`
para também registrar os candidatos gerados (útil para auditar a cobertura,
mesmo para domínios que não resolvem). Domínios já analisados são ignorados via
banco de controle — use `--force` para reanalisar.

## Relacionados

- [Opções Globais](Global-Options-PT) — proxy, threads, timeout, resolver.
- [Saídas & Writers](Output-and-Writers-PT) — para onde os resultados vão.
