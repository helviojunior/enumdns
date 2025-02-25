# EnumDNS

EnumDNS is a modular DNS recon tool! 

Available modules:

1. Brute-force

# Build

Clone the repository and build the project with Golang:

```
git clone https://github.com/helviojunior/enumdns.git
cd enumdns
go get ./...
go build
```

If you want to update go.sum file just run the command `go mod tidy`.

# Installing system wide

After build run the commands bellow

```
go install .
ln -s /root/go/bin/enumdns /usr/bin/enumdns
```

# Utilization

```
$ enumdns brute -h


    ______                      ____  _   _______
   / ____/___  __  ______ ___  / __ \/ | / / ___/
  / __/ / __ \/ / / / __ '__ \/ / / /  |/ /\__ \
 / /___/ / / / /_/ / / / / / / /_/ / /|  /___/ /
/_____/_/ /_/\__,_/_/ /_/ /_/_____/_/ |_//____/


Usage:
  enumdns brute [flags]

Examples:

   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt -o enumdns.txt
   - enumdns brute -d helviojunior.com.br -w /tmp/wordlist.txt --write-jsonl
   - enumdns brute -D domains.txt -w /tmp/wordlist.txt --write-db

Flags:
      --delay int                        Number of seconds delay between navigation and screenshotting (default 3)
  -L, --dns-list string                  File containing a list of DNS sufix
  -d, --dns-sufix string                 Single DNS sufix. (ex: helviojunior.com.br)
  -h, --help                             help for brute
      --log-scan-errors                  Log scan errors (timeouts, DNS errors, etc.) to stderr (warning: can be verbose!)
      --port int                         DNS Server Port (default 53)
      --protocol string                  DNS Server protocol (TCP/UDP) (default "UDP")
  -s, --server string                    DNS Server (default "8.8.8.8")
  -t, --threads int                      Number of concurrent threads (goroutines) to use (default 16)
  -T, --timeout int                      Number of seconds before considering a page timed out (default 60)
  -w, --word-list string                 File containing a list of DNS hosts
      --write-csv                        Write results as CSV (has limited columns)
      --write-csv-file string            The file to write CSV rows to (default "enumdns.csv")
      --write-db                         Write results to a SQLite database
      --write-db-enable-debug            Enable database query debug logging (warning: verbose!)
      --write-db-uri string              The database URI to use. Supports SQLite, Postgres, and MySQL (e.g., postgres://user:pass@host:port/db) (default "sqlite://enumdns.sqlite3")
      --write-elastic                    Write results to a SQLite database
      --write-elasticsearch-uri string   The elastic search URI to use. (e.g., http://user:pass@host:9200/index) (default "http://localhost:9200/intelparser")
      --write-jsonl                      Write results as JSON lines
      --write-jsonl-file string          The file to write JSON lines to (default "enumdns.jsonl")
      --write-none                       Use an empty writer to silence warnings

Global Flags:
  -D, --debug-log                Enable debug logging
  -q, --quiet                    Silence (almost all) logging
  -o, --write-text-file string   The file to write Text lines to

```

### Installing Go v1.23.5

```
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
rm -rf /usr/bin/go && ln -s /usr/local/go/bin/go /usr/bin/go
```


## Disclaimer

This tool is intended for educational purpose or for use in environments where you have been given explicit/legal authorization to do so.