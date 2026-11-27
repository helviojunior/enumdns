# resolve

[🏠 Home](Home) · [🇺🇸 English](Resolve-EN)

O `enumdns resolve` resolve um conjunto de hosts vindos de uma **fonte** externa,
em vez de descobri-los por força-bruta. Possui três subcomandos:

| Subcomando | Fonte |
|------------|-------|
| [`resolve file`](#resolve-file) | Um arquivo texto, um host por linha |
| [`resolve bloodhound`](#resolve-bloodhound) | Um export do BloodHound (`.zip` ou `_computers.json`) |
| [`resolve crtsh`](#resolve-crtsh) | Logs de Certificate Transparency via crt.sh |

Os três resolvem registros e os vinculam ao SOA da zona, gravando em todos os
[writers](Output-and-Writers-PT) configurados. A flag compartilhada
`-I/--IgnoreNonexistent` pula entradas não resolvíveis em vez de abortar.

---

## resolve file

Resolve todos os hosts listados em um arquivo texto.

```bash
enumdns resolve file -L <host_list.txt> [writers] [flags globais]
```

| Flag | Descrição |
|------|-----------|
| `-L, --host-list` | Arquivo com um host DNS por linha. |

```bash
enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
enumdns resolve file -L /tmp/host_list.txt --write-jsonl
enumdns resolve file -L /tmp/host_list.txt --write-db
```

Hosts já resolvidos são ignorados via banco de controle; use `--force` para
resolver novamente.

---

## resolve bloodhound

Resolve as contas de computador de uma coleta do **BloodHound**. Aceita tanto um
arquivo `_computers.json` quanto o `.zip` da coleta (o `_computers.json` é
extraído automaticamente). Contas de computador desabilitadas são ignoradas.

```bash
enumdns resolve bloodhound -L <arquivo.zip|_computers.json> [writers] [flags globais]
```

| Flag | Descrição |
|------|-----------|
| `-L, --bloodhound-file` | Saída do BloodHound: um `.zip` ou um `_computers.json`. |

```bash
enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db
```

O comando primeiro verifica a conectividade DNS de cada domínio encontrado no
export. Em um ambiente de AD, aponte `-s` para um domain controller para que os
nomes internos resolvam:

```bash
enumdns resolve bloodhound -L export.zip -s 10.0.0.10 --write-db
```

---

## resolve crtsh

Varre os logs de **Certificate Transparency** (crt.sh) de um domínio, coleta cada
subdomínio visto em certificados emitidos e então os resolve.

```bash
enumdns resolve crtsh -d <domínio> [writers] [flags globais]
enumdns resolve crtsh -L <dominios.txt> [writers] [flags globais]
```

| Flag | Descrição |
|------|-----------|
| `-d, --dns-suffix` | Um único domínio a varrer no crt.sh. |
| `-L, --dns-list` | Arquivo com um domínio por linha. |
| `--fqdn-out` | Também anexa cada FQDN descoberto a este arquivo de saída. |

```bash
enumdns resolve crtsh -d test.com -o enumdns.txt
enumdns resolve crtsh -L domains.txt --write-jsonl
enumdns resolve crtsh -d test.com --fqdn-out fqdns.txt --write-db
```

Assim como `brute`/`recon`, a varredura do crt.sh respeita `--allow-parent-soa` e
`--allow-tld`, e funciona através de um proxy `-X`.

## Relacionados

- [recon](Recon-PT) — enumera registros de um domínio diretamente.
- [brute](Brute-Force-PT) — descobre hosts via wordlist.
- [Opções Globais](Global-Options-PT) · [Saídas & Writers](Output-and-Writers-PT)
