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

### Installing Go v1.23.5

```
wget https://go.dev/dl/go1.23.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.5.linux-amd64.tar.gz
rm -rf /usr/bin/go && ln -s /usr/local/go/bin/go /usr/bin/go
```


## Disclaimer

This tool is intended for educational purpose or for use in environments where you have been given explicit/legal authorization to do so.