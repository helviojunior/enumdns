# Instalação

[🏠 Home](Home) · [🇺🇸 English](Installation-EN)

O EnumDNS é um único binário Go com link estático. Você pode baixar um release
pré-compilado ou compilar a partir do código-fonte.

## Baixar um release (recomendado)

Binários pré-compilados para **Linux**, **Windows** e **macOS** (386 / amd64 /
arm64) são publicados a cada release:

- Página de releases: <https://github.com/helviojunior/enumdns/releases>
- Passos de instalação por sistema: [INSTALL.md](https://github.com/helviojunior/enumdns/blob/main/INSTALL.md)

```bash
# Exemplo (Linux amd64) — ajuste a versão/arquitetura ao seu sistema
VERSION=0.1.36
curl -sL -o enumdns.tar.gz \
  "https://github.com/helviojunior/enumdns/releases/download/v${VERSION}/enumdns-v${VERSION}-linux-amd64.tar.gz"
tar -xzf enumdns.tar.gz
sudo mv enumdns /usr/local/bin/
enumdns version
```

Cada arquivo acompanha um checksum `.md5` para verificar o download.

## Compilar do código-fonte

Requisitos: **Go 1.24+** e `git`.

```bash
git clone https://github.com/helviojunior/enumdns.git
cd enumdns
make local        # gera ./build/enumdns_<os>_<arch> para a plataforma atual
```

Outros alvos do make:

| Alvo | Resultado |
|------|-----------|
| `make local` | Compila para o SO/arquitetura atual (mais rápido) |
| `make linux` / `make windows` / `make mac` | Cross-compila para um SO (todas as arquiteturas) |
| `make all` | Cross-compila para todos os SO/arquiteturas |
| `make fmt` | Formata todo o código Go com `gofmt -s -w` |
| `make clean` | Remove a pasta `build/` |

A versão é embutida em tempo de compilação. Um build de desenvolvimento reporta o
sufixo `-dev` (ex.: `0.1.36-dev-<githash>`); binários oficiais de release
reportam uma versão limpa.

## Verificar a instalação

```bash
enumdns version          # exibe versão / informações de build
enumdns --help           # lista todos os comandos e flags globais
enumdns recon --help     # ajuda por comando, com exemplos de uso
```

## Próximos passos

- [Opções Globais](Global-Options-PT) — flags compartilhadas por todos os
  comandos (servidor DNS, proxy, timeouts, saída).
- [recon](Recon-PT) — sua primeira enumeração.

> ⚠️ O EnumDNS é destinado a fins educacionais ou a ambientes onde você possui
> autorização explícita e legal.
