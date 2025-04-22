# EnumDNS

EnumDNS is a modular DNS reconnaissance tool capable of resolving hosts from various sources, including wordlists, BloodHound files, and Active Directory environments.

Available modules:

1. Brute-force
2. Enumerate DNS registers (CNAME, A, AAAA, NS and so on)
3. Resolve DNS hosts from txt file
4. Resolve DNS hosts from BloodHound file (.zip or .json)


## Main features

- [x] Perform brute-force DNS enumeration to discover hostnames  
- [x] Support for custom DNS suffix lists  
- [x] Automatically identify cloud provider services  
- [x] Retrieve multiple DNS record types (e.g., CNAME, A, AAAA)  
- [x] Enumerate all domain controllers names and IPs (in a Active Directory environment)
- [x] Support to SOCKS (socks4/socks5) proxy
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

   - enumdns recon -d helviojunior.com.br -o enumdns.txt
   - enumdns recon -d helviojunior.com.br --write-jsonl
   - enumdns recon -D domains.txt --write-db

   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -D domains.txt -w /tmp/wordlist.txt --write-db

   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json -o enumdns.txt
   - enumdns resolve bloodhound -L /tmp/bloodhound_files.zip --write-jsonl
   - enumdns resolve bloodhound -L /tmp/bloodhound_computers.json --write-db

   - enumdns resolve file -L /tmp/host_list.txt -o enumdns.txt
   - enumdns resolve file -L /tmp/host_list.txt --write-jsonl
   - enumdns resolve file -L /tmp/host_list.txt --write-db

Available Commands:
  brute       Perform brute-force enumeration
  help        Help about any command
  recon       Perform recon enumeration
  report      Work with enumdns reports
  version     Get the enumdns version

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