# EnumDNS

EnumDNS is a modular DNS reconnaissance tool capable of resolving hosts from various sources, including wordlists, BloodHound files, and Active Directory environments.

Available modules:

1. Brute-force
2. Enumerate DNS registers (CNAME, A, AAAA, NS and so on)
3. Resolve DNS hosts from txt file
4. Resolve DNS hosts from BloodHound file (.zip or .json)
5. **Threat Analysis** - Advanced domain security analysis for typosquatting, homographic attacks, and malicious domain detection


## Main features

- [x] Perform brute-force DNS enumeration to discover hostnames  
- [x] Support for custom DNS suffix lists  
- [x] Automatically identify cloud provider services  
- [x] Retrieve multiple DNS record types (e.g., CNAME, A, AAAA)  
- [x] Enumerate all domain controllers names and IPs (in a Active Directory environment)
- [x] Support to SOCKS (socks4/socks5) proxy
- [x] **Threat analysis** with 8 detection techniques (typosquatting, bitsquatting, homographic attacks, etc.)
- [x] **Comprehensive test coverage** (98.4% on threat analysis module)
- [x] Additional advanced features and enhancements  


## Get last release

Check how to get last release by your Operational Systems procedures here [INSTALL.md](https://github.com/helviojunior/enumdns/blob/main/INSTALL.md)


# Utilization

```
$ enumdns -h


    ______                      ____  _   _______
   / ____/___  __  ______ ___  / __ \/ | / / ___/
  / __/ / __ \/ / / / __ '__ \/ / / /  |/ /\__ \
 / /___/ / / / /_/ / / / / / / /_/ / /|  /___/ /
/_____/_/ /_/\__,_/_/ /_/ /_/_____/_/ |_//____/

Usage:
  enumdns [command]

Examples:

   - enumdns recon -d test.com -o enumdns.txt
   - enumdns recon -d test.com --write-jsonl
   - enumdns recon -L domains.txt --write-db

   - enumdns brute -d test.com -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d test.com -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -L domains.txt -w /tmp/wordlist.txt --write-db

   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
   - enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db

   - enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
   - enumdns resolve file -L /tmp/host_list.txt --write-jsonl
   - enumdns resolve file -L /tmp/host_list.txt --write-db

   - enumdns threat-analysis -d example.com --all-techniques -o threats.txt
   - enumdns threat-analysis -d example.com --typosquatting --homographic --write-db
   - enumdns threat-analysis -L domains.txt --all-techniques --max-variations 5000 --write-jsonl

Available Commands:
  brute           Perform brute-force enumeration
  help            Help about any command
  recon           Perform recon enumeration
  report          Work with enumdns reports
  threat-analysis Advanced domain threat analysis for typosquatting and malicious domains
  version         Get the enumdns version

Flags:
  -D, --debug-log                Enable debug logging
  -h, --help                     help for enumdns
  -X, --proxy string             Proxy to pass traffic through: <scheme://ip:port> (e.g., socks4://user:pass@proxy_host:1080
  -q, --quiet                    Silence (almost all) logging
  -o, --write-text-file string   The file to write Text lines to

Use "enumdns [command] --help" for more information about a command.

```


## Disclaimer

This tool is intended for educational purpose or for use in environments where you have been given explicit/legal authorization to do so.
## Threat Analysis Module

The `threat-analysis` module provides advanced domain security analysis to detect malicious domains that could be used in attacks against your organization. This module implements multiple techniques for identifying suspicious domains:

### Available Techniques

- **Typosquatting**: Detects domains with keyboard adjacency errors (e.g., `goggle.com` for `google.com`)
- **Bitsquatting**: Identifies domains created through single bit-flip errors  
- **Homographic Attacks**: Detects Unicode characters that look similar to ASCII (e.g., `рaypal.com` with Cyrillic 'р')
- **Character Insertion/Deletion**: Finds domains with added or removed characters
- **Character Transposition**: Detects swapped adjacent characters
- **TLD Variations**: Analyzes suspicious TLDs (.tk, .ml, .ga, etc.)
- **Subdomain Patterns**: Identifies phishing patterns like "secure-", "login-", "verify-"

### Scope & Flags

- Scope: Variations occur on the registrable domain (PSL). Subdomains to the left are preserved.
  - `microsoft.com` → vary `microsoft.*`
  - `recife.pe.gov.br` → vary `pe.gov.br` and suffix `gov.br` (no changes to `recife`).
- Suffix focus: `gov.br` includes suffix impersonation (e.g., `g0v.br`, homoglyphs) without touching subdomains.
- TLD swaps: Uses union of real suffix + `--target-tlds` (default includes `com.br, net.br, org.br`).
- Deduplicated output: Text writer avoids duplicated lines; use `--emit-candidates` to print generated candidates (including NX).

New/advanced flags:
- `--span-last3`: operate over last 3 labels (mutate 3rd-from-right, keep last 2 as suffix) for tricky cases.
- `--focus-suffix=<suffix>`: emphasize suffix-specific techniques (e.g., `gov.br`).
- `--emit-candidates`: write all generated candidates to outputs before probing.
- `--brand-combo`: add brand prefix/suffix patterns.

### Quick Examples

```bash
# Basic threat analysis with all techniques
enumdns threat-analysis -d yourcompany.com --all-techniques

# Specific techniques only
enumdns threat-analysis -d yourcompany.com --typosquatting --homographic

# Analyze multiple domains from file
enumdns threat-analysis -L company-domains.txt --all-techniques --write-db

# High-volume analysis with custom limits
enumdns threat-analysis -d example.com --all-techniques --max-variations 10000

# Output to different formats
enumdns threat-analysis -d example.com --all-techniques --write-jsonl --write-csv

# Focus on gov.br with candidates (includes NX)
enumdns threat-analysis -d recife.pe.gov.br --all-techniques --focus-suffix=gov.br --emit-candidates -o gov-br.txt

# com.br with broader TLD swaps
enumdns threat-analysis -d yeslinux.com.br --all-techniques --target-tlds com,net,org,co,info,io,com.br,net.br,org.br
```

### Security Features

- **Risk Scoring**: Each domain receives a threat score (0.0-1.0) based on multiple indicators
- **Threat Indicators**: Automatic identification of suspicious patterns
- **Rate Limiting**: Configurable limits to prevent overwhelming DNS servers
- **Proxy Support**: Works with SOCKS proxies for discrete analysis

For detailed documentation, see [documentation.md](documentation.md#análise-de-ameaças-threat-analysis---guia-detalhado).


## Disclaimer

This tool is intended for educational purpose or for use in environments where you have been given explicit/legal authorization to do so.
