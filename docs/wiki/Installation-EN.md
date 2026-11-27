# Installation

[🏠 Home](Home) · [🇧🇷 Português](Installation-PT)

EnumDNS is a single, statically-linked Go binary. You can download a pre-built
release or build it from source.

## Download a release (recommended)

Pre-built binaries for **Linux**, **Windows** and **macOS** (386 / amd64 /
arm64) are published on every release:

- Releases page: <https://github.com/helviojunior/enumdns/releases>
- OS-specific install steps: [INSTALL.md](https://github.com/helviojunior/enumdns/blob/main/INSTALL.md)

```bash
# Example (Linux amd64) — adjust the version/arch to your platform
VERSION=0.1.36
curl -sL -o enumdns.tar.gz \
  "https://github.com/helviojunior/enumdns/releases/download/v${VERSION}/enumdns-v${VERSION}-linux-amd64.tar.gz"
tar -xzf enumdns.tar.gz
sudo mv enumdns /usr/local/bin/
enumdns version
```

Each archive ships with a `.md5` checksum so you can verify the download.

## Build from source

Requirements: **Go 1.24+** and `git`.

```bash
git clone https://github.com/helviojunior/enumdns.git
cd enumdns
make local        # builds ./build/enumdns_<os>_<arch> for your current platform
```

Other make targets:

| Target | Result |
|--------|--------|
| `make local` | Build for the current OS/arch (fastest) |
| `make linux` / `make windows` / `make mac` | Cross-compile for a single OS (all archs) |
| `make all` | Cross-compile for every OS/arch |
| `make fmt` | Format all Go code with `gofmt -s -w` |
| `make clean` | Remove the `build/` folder |

The version string is embedded at build time. A development build reports a
`-dev` suffix (e.g. `0.1.36-dev-<githash>`); official release binaries report a
clean version.

## Verify the install

```bash
enumdns version          # print version / build info
enumdns --help           # list all commands and global flags
enumdns recon --help     # per-command help, with usage examples
```

## Next steps

- [Global Options](Global-Options-EN) — flags shared by every command (DNS
  server, proxy, timeouts, output).
- [recon](Recon-EN) — your first enumeration.

> ⚠️ EnumDNS is intended for educational purposes or for environments where you
> have explicit, legal authorization.
