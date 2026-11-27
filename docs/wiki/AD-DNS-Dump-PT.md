# ad — Dump de DNS do AD (ADIDNS)

[🏠 Home](Home) · [🇺🇸 English](AD-DNS-Dump-EN)

O `enumdns ad` extrai **todas as zonas e registros DNS armazenados dentro do
Active Directory** (ADIDNS). Ele consulta um domain controller via **LDAP**,
enumera as partições `DomainDnsZones` / `ForestDnsZones` e decodifica o atributo
binário `dnsRecord` de cada nó — **A, AAAA, NS, CNAME, SOA, PTR, MX, TXT, SRV**.
Os resultados vão para os mesmos [writers](Output-and-Writers-PT) dos demais
comandos.

## Uso

```bash
enumdns ad -d <domínio> <autenticação> [writers] [flags globais]
enumdns ad --target-ip <dc> -d <domínio> <autenticação> [writers] [flags globais]
```

## Autenticação

Escolha **um** método:

| Método | Flags |
|--------|-------|
| Usuário + senha | `-u <user> -p <senha> -d <domínio>` |
| NTLM pass-the-hash | `-u <user> -H <NT>` ou `-H <LM:NT>` |
| Kerberos (ccache) | `-k` com `--krb5-ccache <arquivo>` ou `$KRB5CCNAME` |

Sem credenciais, uma conexão LDAP **anônima** é tentada (normalmente rejeitada
pelo AD) e um aviso é exibido.

## Flags

| Flag | Descrição |
|------|-----------|
| `-d, --domain` | Domínio / realm do Active Directory (usado para descobrir DCs). |
| `-u, --username` | Usuário para autenticação LDAP. |
| `-p, --password` | Senha para autenticação LDAP. |
| `-H, --hashes` | Hash NTLM para pass-the-hash (`NT` ou `LM:NT`). |
| `-k, --kerberos` | Autentica com um ticket cache Kerberos (ccache). |
| `--krb5-ccache` | Caminho do arquivo ccache (padrão `$KRB5CCNAME`). |
| `--krb5-config` | Caminho do `krb5.conf` (padrão `/etc/krb5.conf`). |
| `--spn` | Service principal name para Kerberos (padrão `ldap/<dc>`). |
| `--target-ip` | Força um servidor LDAP específico (host/IP), pulando a descoberta via Global Catalog. |
| `--ldap-port` | Porta LDAP (padrão 389, ou 636 com `--ldaps`). |
| `--ldaps` | Usa LDAPS (TLS) em vez de LDAP em texto puro. |
| `--base-dn` | Restringe o dump a um único container LDAP (pula a autodetecção). |

A descoberta de servidor reaproveita o módulo de recon: com `-d` + `-s <dns>` os
domain controllers / Global Catalogs são encontrados via DNS; `--target-ip`
conecta direto a um servidor. Proxy SOCKS (`-X`) e LDAPS são suportados.

## Exemplos

```bash
# Usuário/senha via LDAP, relatório texto
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' -o adidns.txt

# Via LDAPS (TLS), JSON Lines
enumdns ad -s 10.0.0.1 --ldaps -d corp.local -u jdoe -p 'Passw0rd!' --write-jsonl

# Pass-the-hash, para SQLite
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -H aad3b435b51404eeaad3b435b51404ee:<nthash> --write-db

# Conectar direto a um DC específico (sem descoberta)
enumdns ad --target-ip 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' --write-db

# Restringir a um único container LDAP
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' --base-dn 'DC=corp,DC=local' -o adidns.txt

# Kerberos via ccache
KRB5CCNAME=/tmp/jdoe.ccache enumdns ad -s 10.0.0.1 -d corp.local -k --write-db
```

## Como funciona

1. **Descobre / conecta** a um DC via DNS (`-d` + `-s`) ou diretamente
   (`--target-ip`), autenticando com o método escolhido.
2. **Localiza as zonas**: primeiro os locais padrão do MicrosoftDNS; se vierem
   vazios, entra em ação um fallback extenso (container Partitions da floresta,
   subárvore da partição de domínio, sondagem de filhos).
3. **Decodifica os registros**: o atributo binário `dnsRecord` de cada nó é
   parseado em registros tipados e gravado.

O acesso aos atributos do AD é case-insensitive. O comando reporta claramente
quando um domínio **não é integrado ao AD** (servido por Infoblox/BIND) versus
partições presentes mas ilegíveis com as credenciais atuais.

No final, imprime um resumo: `Dumped N record(s) across M zone(s)`.

## Relacionados

- [resolve bloodhound](Resolve-PT#resolve-bloodhound) — resolve contas de
  computador do AD a partir de um export do BloodHound.
- [recon](Recon-PT) — enumera DCs e registros a partir do DNS.
- [Opções Globais](Global-Options-PT) · [Saídas & Writers](Output-and-Writers-PT)
