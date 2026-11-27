# ad — AD DNS dump (ADIDNS)

[🏠 Home](Home) · [🇧🇷 Português](AD-DNS-Dump-PT)

`enumdns ad` dumps **every DNS zone and record stored inside Active Directory**
(ADIDNS). It queries a domain controller over **LDAP**, enumerates the
`DomainDnsZones` / `ForestDnsZones` partitions and decodes the binary
`dnsRecord` attribute of each node — **A, AAAA, NS, CNAME, SOA, PTR, MX, TXT,
SRV**. Results go to the same [writers](Output-and-Writers-EN) as every other
command.

## Usage

```bash
enumdns ad -d <domain> <auth> [writers] [global flags]
enumdns ad --target-ip <dc> -d <domain> <auth> [writers] [global flags]
```

## Authentication

Pick **one** method:

| Method | Flags |
|--------|-------|
| Username + password | `-u <user> -p <pass> -d <domain>` |
| NTLM pass-the-hash | `-u <user> -H <NT>` or `-H <LM:NT>` |
| Kerberos (ccache) | `-k` with `--krb5-ccache <file>` or `$KRB5CCNAME` |

Without credentials, an **anonymous** LDAP bind is attempted (usually rejected
by AD) and a warning is printed.

## Flags

| Flag | Description |
|------|-------------|
| `-d, --domain` | Active Directory domain / realm (used to discover DCs). |
| `-u, --username` | Username for LDAP authentication. |
| `-p, --password` | Password for LDAP authentication. |
| `-H, --hashes` | NTLM hash for pass-the-hash (`NT` or `LM:NT`). |
| `-k, --kerberos` | Authenticate with a Kerberos ticket cache (ccache). |
| `--krb5-ccache` | Path to the ccache file (defaults to `$KRB5CCNAME`). |
| `--krb5-config` | Path to `krb5.conf` (defaults to `/etc/krb5.conf`). |
| `--spn` | Service principal name for Kerberos (defaults to `ldap/<dc>`). |
| `--target-ip` | Force a specific LDAP server (host/IP), skipping Global Catalog discovery. |
| `--ldap-port` | LDAP port (default 389, or 636 with `--ldaps`). |
| `--ldaps` | Use LDAPS (TLS) instead of plaintext LDAP. |
| `--base-dn` | Restrict the dump to a single LDAP container (skips autodetection). |

Server discovery reuses the recon module: with `-d` + `-s <dns>` the domain
controllers / Global Catalogs are found via DNS; `--target-ip` connects to one
server directly. SOCKS proxy (`-X`) and LDAPS are supported.

## Examples

```bash
# User/password over LDAP, text report
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' -o adidns.txt

# Over LDAPS (TLS), JSON Lines
enumdns ad -s 10.0.0.1 --ldaps -d corp.local -u jdoe -p 'Passw0rd!' --write-jsonl

# Pass-the-hash, into SQLite
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -H aad3b435b51404eeaad3b435b51404ee:<nthash> --write-db

# Connect straight to a specific DC (no discovery)
enumdns ad --target-ip 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' --write-db

# Restrict to a single LDAP container
enumdns ad -s 10.0.0.1 -d corp.local -u jdoe -p 'Passw0rd!' --base-dn 'DC=corp,DC=local' -o adidns.txt

# Kerberos via ccache
KRB5CCNAME=/tmp/jdoe.ccache enumdns ad -s 10.0.0.1 -d corp.local -k --write-db
```

## How it works

1. **Discover / connect** to a DC via DNS (`-d` + `-s`) or directly
   (`--target-ip`), binding with the chosen auth method.
2. **Find the zones**: the standard MicrosoftDNS locations first; if those come
   up empty, an extensive fallback kicks in (forest Partitions container,
   domain-partition subtree, child probing).
3. **Decode records**: each node's binary `dnsRecord` attribute is parsed into
   typed records and written out.

AD attribute access is case-insensitive. The command reports clearly when a
domain is **not AD-integrated** (served by Infoblox/BIND instead) versus
partitions that are present but unreadable with the current credentials.

At the end it prints a summary: `Dumped N record(s) across M zone(s)`.

## Related

- [resolve bloodhound](Resolve-EN#resolve-bloodhound) — resolve AD computer
  accounts from a BloodHound export.
- [recon](Recon-EN) — enumerate DCs and records from DNS.
- [Global Options](Global-Options-EN) · [Output & Writers](Output-and-Writers-EN)
