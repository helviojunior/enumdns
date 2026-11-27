# wordlist

[🏠 Home](Home) · [🇺🇸 English](Wordlist-EN)

O `enumdns wordlist` gera uma **wordlist de labels DNS limpa, deduplicada e
ordenada** a partir de um ou mais arquivos de entrada (texto puro, resultados
anteriores, dados coletados, …). O resultado é ideal como entrada `-w` do
[brute](Brute-Force-PT).

## Uso

```bash
enumdns wordlist -i '<glob|arquivo>' [-i ...] [filtros] -o <saida.txt>
```

Sem `-o`, a wordlist é impressa no stdout.

## Flags

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-i, --input` | — | Arquivo ou padrão glob de entrada (repetível). **Coloque globs entre aspas** (`-i '*.txt'`) para o enumdns expandi-los. |
| `--exclude` | — | Termos a descartar: um arquivo (um por linha) ou literais separados por vírgula (repetível), ex.: `--exclude cloud,mail,onmicrosoft` ou `--exclude deny.txt`. |
| `--min-length` | `4` | Comprimento mínimo de label a manter. |
| `--max-length` | `63` | Comprimento máximo de label (o limite DNS é 63). |
| `--keep-tld` | off | Mantém o TLD final de nomes com ponto (ex.: `com`, `br`) em vez de removê-lo. |
| `--keep-numeric` | off | Mantém tokens totalmente numéricos (ex.: `144`, `3600`) em vez de descartá-los. |
| `-o, --write-text-file` | — | Arquivo de saída (flag global). Sempre excluído das entradas. |

## Como os tokens são processados

Cada linha é dividida em labels candidatos em tudo que não for letra, dígito,
underscore ou hífen — então `host.example.com` gera `host` e `example`. Cada
token então é:

- convertido para minúsculas e **deduplicado**;
- **validado** como label DNS (lixo com cara de flag, como `--allow-parent-soa`,
  é descartado);
- descartado quando **totalmente numérico** (octetos de IP, TTLs, seriais de
  SOA), a menos que `--keep-numeric`;
- filtrado por **comprimento** (`--min-length` / `--max-length`);
- filtrado contra a sua **lista de exclusão** (`--exclude`).

Apenas o TLD final de um nome *com ponto* é removido (o `com` de
`www.example.com`), de modo que um label interno que apenas compartilha o nome
com um gTLD (ex.: `dev`) é mantido. Um token sem ponto nunca é tratado como TLD.
Use `--keep-tld` para manter o TLD também.

O arquivo de saída (`-o`) é sempre removido das entradas, então regenerar a lista
nunca realimenta a si mesma.

## Exemplos

```bash
# Todos os .txt do diretório atual -> uma wordlist limpa
enumdns wordlist -i '*.txt' -o custom_wl.txt

# Múltiplas entradas + exclusões inline
enumdns wordlist -i '/data/recon/*.txt' -i extra.txt --exclude cloud,mail,onmicrosoft -o custom_wl.txt

# Um arquivo, arquivo de exclusão, comprimento mínimo 3, mantendo o TLD
enumdns wordlist -i dump.txt --exclude deny.txt --min-length 3 --keep-tld -o custom_wl.txt

# Alimentando direto o brute
enumdns wordlist -i '*.txt' -o wl.txt
enumdns brute -d test.com -w wl.txt --write-db
```

## Relacionados

- [brute](Brute-Force-PT) — consome a wordlist como `-w`.
- [Opções Globais](Global-Options-PT)
