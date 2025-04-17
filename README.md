# EnumDNS

EnumDNS is a modular DNS recon tool! 

Available modules:

1. Brute-force
2. Enumerate DNS registers (CNAME, A, AAAA, NS and so on)


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


## Get Linux last release
```
apt install curl jq

url=$(curl -s https://api.github.com/repos/helviojunior/enumdns/releases | jq -r '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match("linux-amd64.tar.gz$")) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

cd /opt
rm -rf enumdns-latest.tar.gz enumdns
wget -nv -O enumdns-latest.tar.gz "$url"
tar -xzf enumdns-latest.tar.gz

rsync -av enumdns /usr/local/sbin/
chmod +x /usr/local/sbin/enumdns

enumdns version
```

# Utilization

```
$ enumdns -h


    ______                      ____  _   _______
   / ____/___  __  ______ ___  / __ \/ | / / ___/
  / __/ / __ \/ / / / __ '__ \/ / / /  |/ /\__ \
 / /___/ / / / /_/ / / / / / / /_/ / /|  /___/ /
/_____/_/ /_/\__,_/_/ /_/ /_/_____/_/ |_//____/
                                  Ver: dev-dev

Usage:
  enumdns [command]

Examples:

   - enumdns recon -d helviojunior.com.br -o enumdns.txt
   - enumdns recon -d helviojunior.com.br --write-jsonl
   - enumdns recon -D domains.txt --write-db

   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -D domains.txt -w /tmp/wordlist.txt --write-db

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

### Installing Go v1.23.5

```
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
rm -rf /usr/bin/go && ln -s /usr/local/go/bin/go /usr/bin/go
```


## Disclaimer

This tool is intended for educational purpose or for use in environments where you have been given explicit/legal authorization to do so.