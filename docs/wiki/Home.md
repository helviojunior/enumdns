# EnumDNS

**Modular DNS reconnaissance tool** — resolve hosts from wordlists, text files,
BloodHound exports and Active Directory, enumerate DNS records, hunt for
malicious look-alike domains, and dump AD-integrated DNS (ADIDNS) over LDAP.

🌐 **Idioma / Language:** **English** · [Português](#português)

---

## English

Usage guide, one page per command:

| Command / Topic | Description |
|-----------------|-------------|
| [Installation](Installation-EN) | Download a release or build from source |
| [Global Options](Global-Options-EN) | DNS server, proxy, timeouts, threads, SOA/TLD scoping |
| [recon](Recon-EN) | Enumerate DNS records (A, AAAA, CNAME, NS, MX, TXT…) for a domain |
| [brute](Brute-Force-EN) | Brute-force hostnames from a wordlist |
| [resolve](Resolve-EN) | Resolve hosts from a text file, a BloodHound export or crt.sh |
| [threat-analysis](Threat-Analysis-EN) | Typosquatting, homographic and other malicious-domain detection |
| [ad (ADIDNS)](AD-DNS-Dump-EN) | Dump AD-integrated DNS zones/records over LDAP |
| [wordlist](Wordlist-EN) | Build a clean DNS-label wordlist from arbitrary files |
| [report](Reports-EN) | Convert between SQLite/JSONL and sync results to Elasticsearch |
| [Output & Writers](Output-and-Writers-EN) | Output formats, SOA linking, the control database |

### Quick start

```bash
# Enumerate a domain and save a text report
enumdns recon -d example.com -o enumdns.txt

# Brute-force hostnames with a wordlist, store in SQLite
enumdns brute -d example.com -w /tmp/wordlist.txt --write-db

# See all commands
enumdns --help
```

---

## Português

Guia de uso, uma página por comando:

| Comando / Tópico | Descrição |
|------------------|-----------|
| [Instalação](Installation-PT) | Baixar um release ou compilar do código-fonte |
| [Opções Globais](Global-Options-PT) | Servidor DNS, proxy, timeouts, threads, escopo SOA/TLD |
| [recon](Recon-PT) | Enumerar registros DNS (A, AAAA, CNAME, NS, MX, TXT…) de um domínio |
| [brute](Brute-Force-PT) | Força-bruta de hostnames a partir de uma wordlist |
| [resolve](Resolve-PT) | Resolver hosts de um arquivo texto, export do BloodHound ou crt.sh |
| [threat-analysis](Threat-Analysis-PT) | Detecção de typosquatting, homográficos e domínios maliciosos |
| [ad (ADIDNS)](AD-DNS-Dump-PT) | Extrair zonas/registros DNS integrados ao AD via LDAP |
| [wordlist](Wordlist-PT) | Gerar uma wordlist limpa de labels DNS a partir de arquivos |
| [report](Reports-PT) | Converter entre SQLite/JSONL e sincronizar resultados com o Elasticsearch |
| [Saídas & Writers](Output-and-Writers-PT) | Formatos de saída, vínculo de SOA, banco de controle |

### Começando

```bash
# Enumerar um domínio e salvar um relatório texto
enumdns recon -d example.com -o enumdns.txt

# Força-bruta de hostnames com uma wordlist, gravando em SQLite
enumdns brute -d example.com -w /tmp/wordlist.txt --write-db

# Ver todos os comandos
enumdns --help
```

---

> ⚠️ **Disclaimer / Aviso:** EnumDNS is intended for educational purposes or for
> environments where you have explicit, legal authorization. / O EnumDNS é
> destinado a fins educacionais ou a ambientes onde você possui autorização
> explícita e legal.
