# EnumDNS - Documentação Completa

## Índice
1. [Visão Geral](#visão-geral)
2. [Funcionalidades](#funcionalidades)
3. [Arquitetura e Estrutura](#arquitetura-e-estrutura)
4. [Instalação](#instalação)
5. [Uso e Exemplos](#uso-e-exemplos)
6. [Configuração](#configuração)
7. [Módulos e Componentes](#módulos-e-componentes)
8. [Formatos de Saída](#formatos-de-saída)
9. [Desenvolvimento](#desenvolvimento)
10. [Considerações de Segurança](#considerações-de-segurança)

## Visão Geral

O **EnumDNS** é uma ferramenta modular de reconhecimento DNS desenvolvida em Go, projetada para profissionais de segurança cibernética realizarem enumeração DNS abrangente e análise de infraestrutura. A ferramenta oferece múltiplos métodos de descoberta de hosts e pode identificar automaticamente provedores de nuvem, produtos SaaS e datacenters.

### Principais Características
- **Modular**: Suporta diferentes tipos de enumeração (brute-force, reconhecimento, resolução)
- **Multi-plataforma**: Funciona em Linux, Windows e macOS
- **Flexível**: Múltiplos formatos de saída (texto, JSON, CSV, SQLite, Elasticsearch)
- **Escalável**: Suporta processamento paralelo com goroutines
- **Integrado**: Compatível com BloodHound e crt.sh
- **Proxy Support**: Suporta SOCKS4/SOCKS5

## Funcionalidades

### Módulos Principais

#### 1. **Brute-force DNS**
- Enumeração por força bruta usando wordlists
- Suporte a múltiplos sufixos DNS
- Detecção automática de produtos cloud/SaaS

#### 2. **Reconhecimento DNS**
- Enumeração completa de registros DNS
- Detecção automática de controladores de domínio
- Descoberta de novos domínios através de registros encontrados

#### 3. **Resolução de Hosts**
- Resolução a partir de arquivos de texto
- Integração com BloodHound (arquivos .zip e .json)
- Integração com crt.sh para descoberta de certificados

#### 4. **Relatórios**
- Conversão entre formatos (SQLite ↔ JSON Lines ↔ Texto)
- Sincronização com Elasticsearch

### Recursos Avançados

- **Detecção de Active Directory**: Identifica DCs e GCs automaticamente
- **Identificação de Cloud**: Reconhece AWS, Azure, GCP, CloudFlare, etc.
- **Reverse DNS**: Resolução PTR automática
- **ASN Detection**: Identifica ASN e informações de rede
- **Controle de Duplicatas**: Evita varreduras desnecessárias

## Arquitetura e Estrutura

### Estrutura de Diretórios

```
enumdns/
├── cmd/                    # Comandos CLI (Cobra)
│   ├── brute.go           # Comando de brute-force
│   ├── recon.go           # Comando de reconhecimento  
│   ├── resolve*.go        # Comandos de resolução
│   ├── report*.go         # Comandos de relatório
│   └── root.go            # Comando raiz e configurações
├── internal/              # Código interno
│   ├── ascii/             # Interface de usuário
│   ├── disk/              # Utilitários de disco
│   ├── tools/             # Ferramentas auxiliares
│   └── version/           # Informações de versão
├── pkg/                   # Pacotes públicos
│   ├── database/          # Gerenciamento de banco de dados
│   ├── log/               # Sistema de logging
│   ├── models/            # Modelos de dados
│   ├── readers/           # Leitores de entrada
│   ├── runner/            # Engines de execução
│   └── writers/           # Writers de saída
└── main.go               # Ponto de entrada
```

### Componentes Arquiteturais

#### 1. **Runners (Engines de Execução)**
- **Runner**: Engine principal para brute-force e resolução
- **Recon**: Engine especializado para reconhecimento completo
- Gerenciam workers (goroutines) e coordenam a execução

#### 2. **Readers (Leitores de Entrada)**
- **FileReader**: Lê wordlists e listas de domínios
- **CrtShReader**: Integração com crt.sh
- Suportam diferentes formatos de entrada

#### 3. **Writers (Escritores de Saída)**
- **DbWriter**: SQLite, PostgreSQL, MySQL
- **JsonWriter**: JSON Lines
- **CsvWriter**: Formato CSV
- **ElasticWriter**: Elasticsearch
- **TextWriter**: Relatórios legíveis
- **StdoutWriter**: Saída no terminal

#### 4. **Models (Modelos de Dados)**
- **Result**: Resultado principal de DNS
- **FQDNData**: Dados de FQDN descobertos
- **ASN**: Informações de Sistema Autônomo
- **ASNIpDelegate**: Delegações de IP por ASN

## Instalação

### Linux

```bash
# Instalação automática
apt install curl jq

url=$(curl -s https://api.github.com/repos/helviojunior/enumdns/releases | jq -r '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match("linux-amd64.tar.gz$")) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

cd /tmp
rm -rf enumdns-latest.tar.gz enumdns
wget -nv -O enumdns-latest.tar.gz "$url"
tar -xzf enumdns-latest.tar.gz

rsync -av enumdns /usr/local/sbin/
chmod +x /usr/local/sbin/enumdns

enumdns version
```

### macOS

```bash
# Instalar Homebrew (se necessário)
/bin/bash -c "$(curl -fsSL raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Instalar EnumDNS
brew install curl jq

arch=$(if [[ "$(uname -m)" -eq "x86_64" ]]; then echo "amd64"; else echo "arm64"; fi)

url=$(curl -s https://api.github.com/repos/helviojunior/enumdns/releases | jq -r --arg filename "darwin-${arch}.tar.gz\$" '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match($filename)) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

cd /tmp
rm -rf enumdns-latest.tar.gz enumdns
curl -sS -L -o enumdns-latest.tar.gz "$url"
tar -xzf enumdns-latest.tar.gz

rsync -av enumdns /usr/local/sbin/
chmod +x /usr/local/sbin/enumdns

enumdns version
```

### Windows (PowerShell)

```powershell
# Download latest helviojunior/enumdns release from github
function Invoke-Downloadenumdns {

    $repo = "helviojunior/enumdns"
    
    # Determine OS and Architecture
    $osPlatform = [System.Runtime.InteropServices.RuntimeInformation]::OSDescription
    $architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture

    if ($architecture -eq $null -or $architecture -eq "") {
        $architecture = $Env:PROCESSOR_ARCHITECTURE
    }
        
    if ($osPlatform -eq $null -or $osPlatform -eq "") {
        $osPlatform = $Env:OS
    }

    # Adjust the platform and architecture for the API call
    $platform = switch -Wildcard ($osPlatform) {
        "*Windows*" { "windows" }
        "*Linux*"   { "linux" }
        "*Darwin*"  { "darwin" } # MacOS is identified as Darwin
        Default     { "unknown" }
    }
    $arch = switch ($architecture) {
        "X64"  { "amd64" }
        "AMD64"  { "amd64" }
        "X86"  { "386" }
        "Arm"  { "arm" }
        "Arm64" { "arm64" }
        Default { "unknown" }
    }

    if ($platform -eq "unknown" -or $arch -eq "unknown") {
        Write-Error "Cannot get OS Platform and Architecture"
        Return
    }

    Write-Host Getting release list
    $releases = "https://api.github.com/repos/$repo/releases"

    $asset = Invoke-WebRequest $releases | ConvertFrom-Json | Sort-Object -Descending -Property "Id" | ForEach-Object -Process { Get-AssetData -Release $_ -OSPlatform $platform -OSArchitecture $arch } | Select-Object -First 1

    if ($asset -eq $null -or $asset.browser_download_url -eq $null){
        Write-Error "Cannot find a valid URL"
        Return
    }

    $tmpPath = $Env:Temp
    if ($tmpPath -eq $null -or $tmpPath -eq "") {
        $tmpPath = $Env:TMPDIR
    }
    if ($tmpPath -eq $null -or $tmpPath -eq "") {
        $tmpPath = switch ($platform) {
            "windows" { "c:\windows\temp\" }
            "linux"   { "/tmp" }
            "darwin"  { "/tmp" }
        }
    }

    $extension = switch ($platform) {
        "windows" { ".zip" }
        "linux"   { ".tar.gz" }
        "darwin"  { ".tar.gz" } # MacOS is identified as Darwin
        Default     { "unknown" }
    }

    $file = "enumdns-latest$extension"

    Write-Host Dowloading latest release
    $zip = Join-Path -Path $tmpPath -ChildPath $file
    Remove-Item $zip -Force -ErrorAction SilentlyContinue 
    Invoke-WebRequest $asset.browser_download_url -Out $zip

    Write-Host Extracting release files
    if ($extension -eq ".zip") {
        Expand-Archive $zip -Force -DestinationPath $tmpPath
    }else{
        . tar -xzf "$zip" -C "$tmpPath"
    }

    $exeFilename = switch ($platform) {
        "windows" { "enumdns.exe" }
        "linux"   { "enumdns" }
        "darwin"  { "enumdns" } 
    }

    try {
        $dstPath = (New-Object -ComObject Shell.Application).NameSpace('shell:Downloads').Self.Path
    } catch {
        $dstPath = switch ($platform) {
            "windows" { "~\Downloads\" }
            "linux"   { "/usr/local/sbin/" }
            "darwin"  { "/usr/local/sbin/" } 
        }
    }

    try {
        $name = Join-Path -Path $dstPath -ChildPath $exeFilename

        # Cleaning up target dir
        Remove-Item $name -Recurse -Force -ErrorAction SilentlyContinue 

        # Moving from temp dir to target dir
        Move-Item $(Join-Path -Path $tmpPath -ChildPath $exeFilename) -Destination $name -Force

        # Removing temp files
        Remove-Item $zip -Force
    } catch {
        $name = Join-Path -Path $tmpPath -ChildPath $exeFilename
    }
    
    Write-Host "EnumDNS saved at $name" -ForegroundColor DarkYellow

    Write-Host "Getting enumdns version banner"
    . $name version 
}

Function Get-AssetData {
    [CmdletBinding(SupportsShouldProcess = $False)]
    [OutputType([object])]
    Param (
        [Parameter(Mandatory = $True, Position = 0)]
        [object]$Release,
        [Parameter(Mandatory = $True, Position = 1)]
        [string]$OSPlatform,
        [Parameter(Mandatory = $True, Position = 2)]
        [string]$OSArchitecture
    )

    if($Release -is [system.array]){
        $Release = $Release[0]
    }
    
    if (Get-Member -inputobject $Release -name "assets" -Membertype Properties) {
        
        $extension = switch ($OSPlatform) {
            "windows" { ".zip" }
            "linux"   { ".tar.gz" }
            "darwin"  { ".tar.gz" } # MacOS is identified as Darwin
            Default     { "unknown" }
        }

        foreach ($asset in $Release.assets)
        {
            If ($asset.name.Contains("enumdns-") -and $asset.name.Contains("$OSPlatform-$OSArchitecture$extension")) { Return $asset }
        }

    }
    Return $null
} 

Invoke-Downloadenumdns
```

### Compilação Manual

```bash
# Pré-requisitos: Go 1.23.0+
git clone https://github.com/helviojunior/enumdns.git
cd enumdns
go get ./...
go build

# Instalação system-wide
go install .
ln -s /root/go/bin/enumdns /usr/bin/enumdns
```

## Uso e Exemplos

### Sintaxe Básica

```bash
enumdns [comando] [flags]
```

### Comando Recon (Reconhecimento)

```bash
# Reconhecimento de domínio único
enumdns recon -d example.com -o results.txt

# Múltiplos domínios
enumdns recon -L domains.txt --write-jsonl

# Com banco de dados
enumdns recon -d example.com --write-db
```

### Comando Brute (Força Bruta)

```bash
# Brute-force básico
enumdns brute -d example.com -w wordlist.txt -o results.txt

# Múltiplos domínios com wordlist
enumdns brute -L domains.txt -w wordlist.txt --write-db

# Modo rápido (apenas registros A)
enumdns brute -d example.com -w wordlist.txt -Q --write-jsonl
```

### Comando Resolve

#### Resolução de Arquivo
```bash
# Lista de hosts
enumdns resolve file -L hosts.txt -o results.txt

# Com saída JSON
enumdns resolve file -L hosts.txt --write-jsonl
```

#### Integração BloodHound
```bash
# Arquivo JSON de computadores
enumdns resolve bloodhound -L computers.json -o results.txt

# Arquivo ZIP completo  
enumdns resolve bloodhound -L bloodhound_data.zip --write-db
```

#### Integração crt.sh
```bash
# Descoberta via certificados
enumdns resolve crtsh -d example.com --write-db

# Salvar FQDNs encontrados
enumdns resolve crtsh -d example.com --fqdn-out discovered_hosts.txt
```

### Comando Report

#### Conversão de Formatos
```bash
# SQLite para JSON Lines
enumdns report convert --from-file results.sqlite3 --to-file results.jsonl

# JSON Lines para texto
enumdns report convert --from-file results.jsonl --to-file results.txt

# SQLite para texto
enumdns report convert --from-file results.sqlite3 --to-file results.txt
```

#### Sincronização Elasticsearch
```bash
# Enviar dados para Elasticsearch
enumdns report elastic --from-file results.sqlite3 --elasticsearch-uri http://localhost:9200/enumdns
```

## Configuração

### Flags Globais

```bash
# Servidor DNS customizado
enumdns recon -d example.com -s 8.8.8.8 --port 53

# Protocolo DNS (TCP/UDP)
enumdns recon -d example.com --protocol TCP

# Proxy SOCKS
enumdns recon -d example.com -X socks5://127.0.0.1:1080

# Controle de threads
enumdns recon -d example.com -t 10

# Timeout personalizado
enumdns recon -d example.com -T 120

# Debug logging
enumdns recon -d example.com -D

# Modo silencioso
enumdns recon -d example.com -q

# Forçar re-verificação
enumdns recon -d example.com -F
```

### Opções de Saída

```bash
# Múltiplas saídas simultâneas
enumdns recon -d example.com \
  -o results.txt \
  --write-jsonl \
  --write-db \
  --write-csv

# Banco de dados customizado
enumdns recon -d example.com --write-db-uri "postgres://user:pass@host:5432/db"

# Elasticsearch customizado
enumdns recon -d example.com --write-elasticsearch-uri "http://user:pass@host:9200/index"

# Workspace local
enumdns recon -d example.com --local-workspace
```

## Módulos e Componentes

### Sistema de DNS

#### Cliente SOCKS
- Suporte nativo a SOCKS4/SOCKS5
- Integração transparente com proxies
- Fallback automático para DNS direto

#### Resolução Multi-Tipo
- A, AAAA, CNAME, MX, NS, SOA, SRV, TXT, PTR
- Resolução reversa automática
- Detecção de registros especiais (_ldap._tcp, _gc._tcp)

### Detecção de Produtos

#### Cloud Providers
- AWS (amazonaws.com, awsdns, etc.)
- Azure (azure-dns.*, azurewebsites.net, etc.)
- GCP (googleusercontent.com)
- CloudFlare (cloudflare.com, cloudflare)
- Akamai (edgekey.net, akamaiedge.net)

#### SaaS Products
- Office 365 (lync.com, office.com, outlook.com)
- SharePoint (sharepointonline.com)
- Heroku (herokuapp.com, herokudns.com)
- GitHub (github.io, github.com)
- Salesforce (exacttarget.com)

#### Datacenters
- Locaweb, Equinix, UOL, HostGator
- Identificação baseada em padrões de DNS

### Active Directory Detection

#### Indicadores de AD
```go
// Busca por registros LDAP
_ldap._tcp.domain.com SRV

// Identificação de Global Catalogs
_gc._tcp.domain.com SRV

// Marcação automática DC/GC
result.DC = true
result.GC = true
```

### Sistema ASN

#### Base de Dados ASN
- Dados de RIRs (ARIN, LACNIC, RIPE, APNIC, AFRINIC)
- Mapeamento IP → ASN automático
- Informações de país e organização

## Formatos de Saída

### Texto (.txt)
```
FQDN                                                                   Type       Value
====================================================================== ========== ==================================================
example.com                                                           A          192.0.2.1
www.example.com                                                       CNAME      example.com
mail.example.com                                                      A          192.0.2.2 (Cloud = AWS)
```

### JSON Lines (.jsonl)
```json
{"fqdn":"example.com","result_type":"A","ipv4":"192.0.2.1","probed_at":"2025-01-15T10:30:00Z"}
{"fqdn":"www.example.com","result_type":"CNAME","target":"example.com","probed_at":"2025-01-15T10:30:01Z"}
```

### SQLite/PostgreSQL/MySQL
```sql
CREATE TABLE results (
    id INTEGER PRIMARY KEY,
    fqdn TEXT,
    result_type TEXT,
    ipv4 TEXT,
    ipv6 TEXT,
    target TEXT,
    cloud_product TEXT,
    dc BOOLEAN,
    gc BOOLEAN,
    probed_at DATETIME
);
```

### CSV
```csv
FQDN,RType,IPv4,IPv6,Target,CloudProduct,DC,GC,ProbedAt
example.com,A,192.0.2.1,,,,,false,false,2025-01-15T10:30:00Z
www.example.com,CNAME,,,example.com,,false,false,2025-01-15T10:30:01Z
```

## Desenvolvimento

### Tecnologias Utilizadas

#### Linguagem e Framework
- **Go 1.23.0+**: Linguagem principal
- **Cobra**: Framework CLI
- **GORM**: ORM para banco de dados
- **Miekg/DNS**: Biblioteca DNS nativa

#### Dependências Principais
```go
// DNS e rede
github.com/miekg/dns v1.1.63
golang.org/x/net v0.31.0

// CLI e UI
github.com/spf13/cobra v1.8.1
github.com/charmbracelet/glamour v0.8.0
github.com/charmbracelet/lipgloss v0.12.1

// Banco de dados
gorm.io/gorm v1.30.0
gorm.io/driver/mysql v1.5.7
gorm.io/driver/postgres v1.5.11
github.com/glebarez/sqlite v1.11.0

// Elasticsearch
github.com/elastic/go-elasticsearch/v8 v8.17.0
```

### Padrões de Código

#### Estrutura de Comandos
```go
var exampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Example command",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Configuração comum
        return nil
    },
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Validações específicas
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Lógica principal
    },
}
```

#### Padrão Runner
```go
type Runner struct {
    Targets chan string
    ctx     context.Context
    cancel  context.CancelFunc
    writers []writers.Writer
    options Options
}

func (r *Runner) Run(total int) {
    wg := sync.WaitGroup{}
    
    // Spawn workers
    for w := 0; w < r.options.Scan.Threads; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Worker logic
        }()
    }
    
    wg.Wait()
}
```

#### Padrão Writer
```go
type Writer interface {
    Write(*models.Result) error
    WriteFqdn(*models.FQDNData) error
    Finish() error
}
```

### Build e Release

#### Makefile Targets
```bash
# Build local
make build

# Build multi-plataforma
make build-all

# Testes
make test

# Linting
make lint

# Release
make release
```

#### Cross-compilation
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o enumdns-linux-amd64

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o enumdns-windows-amd64.exe

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o enumdns-darwin-arm64
```

## Considerações de Segurança

### Uso Responsável
```bash
# IMPORTANTE: Use apenas em ambientes autorizados
# Esta ferramenta é destinada para:
# - Testes de penetração autorizados
# - Auditorias de segurança internas
# - Pesquisa educacional
```

### Rate Limiting
- Delays automáticos entre requisições
- Configuração de timeouts
- Suporte a servidores DNS privados

### Controle de Dados
- Banco de dados local por padrão (~/.enumdns.db)
- Opção de workspace local
- Controle de dados temporários

### Proxies e Anonimato
```bash
# Uso com Tor
enumdns recon -d example.com -X socks5://127.0.0.1:9050

# Uso com proxy corporativo
enumdns recon -d example.com -X socks5://proxy.corp.com:1080
```

## Troubleshooting

### Problemas Comuns

#### 1. Erro de Conectividade DNS
```bash
# Verificar conectividade
enumdns recon -d google.com -s 8.8.8.8

# Usar servidor DNS alternativo
enumdns recon -d example.com -s 1.1.1.1
```

#### 2. Erro de Proxy
```bash
# Testar proxy
curl --socks5 127.0.0.1:1080 http://httpbin.org/ip

# Verificar formato da URL
enumdns recon -d example.com -X socks5://127.0.0.1:1080
```

#### 3. Erro de Banco de Dados
```bash
# Recriar banco de controle
enumdns recon -d example.com --disable-control-db

# Usar workspace local
enumdns recon -d example.com --local-workspace
```

### Debug e Logging
```bash
# Debug completo
enumdns recon -d example.com -D

# Log de erros de scan
enumdns recon -d example.com --log-scan-errors

# Modo verboso com debug
enumdns recon -d example.com -D --log-scan-errors
```

## Exemplos Práticos

### Cenário 1: Auditoria de Domínio Corporativo

```bash
# 1. Reconhecimento inicial
enumdns recon -d empresa.com.br --write-db -o recon_inicial.txt

# 2. Brute-force com wordlist comum
enumdns brute -d empresa.com.br -w /usr/share/wordlists/subdomains-top1million-5000.txt --write-db

# 3. Verificação de certificados SSL
enumdns resolve crtsh -d empresa.com.br --write-db --fqdn-out certificados_encontrados.txt

# 4. Geração de relatório final
enumdns report convert --from-file ~/.enumdns.db --to-file relatorio_final.txt
```

### Cenário 2: Red Team - Enumeração Stealthy

```bash
# Usando Tor para anonimato
enumdns recon -d target.com \
  -X socks5://127.0.0.1:9050 \
  -t 2 \
  -T 180 \
  --write-jsonl

# Brute-force com delay customizado
enumdns brute -d target.com \
  -w small_wordlist.txt \
  -X socks5://127.0.0.1:9050 \
  -t 1 \
  --write-db
```

### Cenário 3: Blue Team - Análise de Infraestrutura

```bash
# Análise completa de múltiplos domínios
enumdns recon -L dominios_empresa.txt \
  --write-db \
  --write-elasticsearch-uri http://siem.empresa.com:9200/dns_enum

# Integração com dados do AD
enumdns resolve bloodhound -L bloodhound_computers.json \
  --write-db \
  -s dc01.empresa.com
```

### Cenário 4: Análise Forense

```bash
# Resolução de lista de hosts suspeitos
enumdns resolve file -L ips_suspeitos.txt \
  --write-db \
  -o analise_forense.txt

# Verificação de infraestrutura C2
enumdns recon -L dominios_c2.txt \
  --write-jsonl \
  --write-csv
```

## Integração com Outras Ferramentas

### Amass
```bash
# Usar resultados do Amass como entrada
amass enum -d example.com -o amass_results.txt
enumdns resolve file -L amass_results.txt --write-db
```

### Subfinder
```bash
# Combinar com Subfinder
subfinder -d example.com -o subfinder_results.txt
enumdns resolve file -L subfinder_results.txt --write-db
```

### Nmap
```bash
# Usar IPs descobertos no Nmap
enumdns recon -d example.com --write-jsonl
# Extrair IPs e usar no Nmap
cat results.jsonl | jq -r '.ipv4' | sort -u > ips_discovered.txt
nmap -iL ips_discovered.txt -sS -p- --open
```

### BloodHound
```bash
# Análise pós-coleta BloodHound
enumdns resolve bloodhound -L 20250115_computers.json \
  --write-db \
  -s 192.168.1.10
```

## Scripts de Automação

### Script de Monitoramento Contínuo

```bash
#!/bin/bash
# monitor_dns.sh

DOMAIN="$1"
INTERVAL="$2"  # em horas

if [[ -z "$DOMAIN" || -z "$INTERVAL" ]]; then
    echo "Uso: $0 <dominio> <intervalo_horas>"
    exit 1
fi

while true; do
    echo "[$(date)] Iniciando varredura de $DOMAIN"
    
    # Reconhecimento completo
    enumdns recon -d "$DOMAIN" \
      --write-db \
      --local-workspace \
      -q
    
    # Verificar certificados
    enumdns resolve crtsh -d "$DOMAIN" \
      --write-db \
      --local-workspace \
      -q
    
    echo "[$(date)] Varredura concluída. Próxima em ${INTERVAL}h"
    sleep "${INTERVAL}h"
done
```

### Script de Análise Comparativa

```bash
#!/bin/bash
# compare_scans.sh

OLD_DB="$1"
NEW_DB="$2"

if [[ -z "$OLD_DB" || -z "$NEW_DB" ]]; then
    echo "Uso: $0 <banco_antigo> <banco_novo>"
    exit 1
fi

# Extrair hosts únicos
sqlite3 "$OLD_DB" "SELECT DISTINCT fqdn FROM results WHERE exists=1" > old_hosts.txt
sqlite3 "$NEW_DB" "SELECT DISTINCT fqdn FROM results WHERE exists=1" > new_hosts.txt

# Encontrar novos hosts
comm -13 <(sort old_hosts.txt) <(sort new_hosts.txt) > new_hosts_found.txt

# Encontrar hosts removidos
comm -23 <(sort old_hosts.txt) <(sort new_hosts.txt) > hosts_removed.txt

echo "Novos hosts encontrados: $(wc -l < new_hosts_found.txt)"
echo "Hosts removidos: $(wc -l < hosts_removed.txt)"

if [[ -s new_hosts_found.txt ]]; then
    echo "=== NOVOS HOSTS ==="
    cat new_hosts_found.txt
fi

if [[ -s hosts_removed.txt ]]; then
    echo "=== HOSTS REMOVIDOS ==="
    cat hosts_removed.txt
fi

# Cleanup
rm -f old_hosts.txt new_hosts.txt
```

## APIs e Integrações

### API REST Wrapper (Exemplo)

```python
#!/usr/bin/env python3
# enumdns_api.py - Wrapper REST para EnumDNS

from flask import Flask, request, jsonify
import subprocess
import json
import tempfile
import os

app = Flask(__name__)

@app.route('/api/v1/recon', methods=['POST'])
def recon_domain():
    data = request.get_json()
    domain = data.get('domain')
    
    if not domain:
        return jsonify({'error': 'Domain required'}), 400
    
    try:
        # Criar arquivo temporário para resultados
        with tempfile.NamedTemporaryFile(mode='w', suffix='.jsonl', delete=False) as f:
            temp_file = f.name
        
        # Executar EnumDNS
        cmd = [
            'enumdns', 'recon',
            '-d', domain,
            '--write-jsonl-file', temp_file,
            '-q'
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'error': result.stderr}), 500
        
        # Ler resultados
        results = []
        with open(temp_file, 'r') as f:
            for line in f:
                if line.strip():
                    results.append(json.loads(line))
        
        # Cleanup
        os.unlink(temp_file)
        
        return jsonify({
            'domain': domain,
            'results': results,
            'count': len(results)
        })
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/v1/brute', methods=['POST'])
def brute_domain():
    data = request.get_json()
    domain = data.get('domain')
    wordlist = data.get('wordlist', [])
    
    if not domain:
        return jsonify({'error': 'Domain required'}), 400
    
    try:
        # Criar wordlist temporária
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as f:
            for word in wordlist:
                f.write(f"{word}\n")
            wordlist_file = f.name
        
        # Criar arquivo temporário para resultados
        with tempfile.NamedTemporaryFile(mode='w', suffix='.jsonl', delete=False) as f:
            temp_file = f.name
        
        # Executar EnumDNS
        cmd = [
            'enumdns', 'brute',
            '-d', domain,
            '-w', wordlist_file,
            '--write-jsonl-file', temp_file,
            '-q'
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'error': result.stderr}), 500
        
        # Ler resultados
        results = []
        with open(temp_file, 'r') as f:
            for line in f:
                if line.strip():
                    results.append(json.loads(line))
        
        # Cleanup
        os.unlink(temp_file)
        os.unlink(wordlist_file)
        
        return jsonify({
            'domain': domain,
            'results': results,
            'count': len(results)
        })
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

### Integração com Docker

```dockerfile
# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o enumdns

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/enumdns .
COPY --from=builder /app/wordlists ./wordlists/

CMD ["./enumdns"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  enumdns:
    build: .
    volumes:
      - ./data:/data
      - ./wordlists:/wordlists
    environment:
      - ENUMDNS_DB_PATH=/data/enumdns.db
    networks:
      - enumdns_net
  
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.15.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data
    networks:
      - enumdns_net
  
  kibana:
    image: docker.elastic.co/kibana/kibana:8.15.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch
    networks:
      - enumdns_net

networks:
  enumdns_net:
    driver: bridge

volumes:
  es_data:
    driver: local
```

## Análise de Desempenho

### Benchmarks Típicos

| Configuração | Threads | Hosts/min | Uso CPU | Uso RAM |
|-------------|---------|-----------|---------|---------|
| Brute-force Local | 6 | ~300 | 15% | 50MB |
| Brute-force Proxy | 3 | ~100 | 10% | 45MB |
| Reconhecimento | 6 | ~500 | 20% | 60MB |
| Resolução Arquivo | 10 | ~800 | 25% | 70MB |

### Otimizações de Performance

```bash
# Para redes rápidas
enumdns brute -d example.com -w wordlist.txt -t 20 -T 30

# Para redes lentas/instáveis  
enumdns brute -d example.com -w wordlist.txt -t 3 -T 180

# Para preservar banda
enumdns brute -d example.com -w wordlist.txt -t 1 -T 300
```

### Monitoramento de Recursos

```bash
#!/bin/bash
# monitor_resources.sh

PID=$(pgrep enumdns)
if [[ -z "$PID" ]]; then
    echo "EnumDNS não está executando"
    exit 1
fi

echo "Monitorando EnumDNS (PID: $PID)"
echo "Time,CPU%,MEM%,VIRT,RES"

while kill -0 "$PID" 2>/dev/null; do
    ps -p "$PID" -o %cpu,%mem,vsize,rss --no-headers | \
    awk -v time="$(date +%H:%M:%S)" '{printf "%s,%.1f,%.1f,%s,%s\n", time, $1, $2, $3, $4}'
    sleep 5
done
```

## Casos de Uso Avançados

### 1. Descoberta de Shadow IT

```bash
#!/bin/bash
# shadow_it_discovery.sh

COMPANY_DOMAIN="$1"
OUTPUT_DIR="shadow_it_$(date +%Y%m%d)"

mkdir -p "$OUTPUT_DIR"

echo "[+] Descobrindo infraestrutura shadow IT para $COMPANY_DOMAIN"

# 1. Buscar certificados SSL
enumdns resolve crtsh -d "$COMPANY_DOMAIN" \
  --fqdn-out "$OUTPUT_DIR/ssl_certificates.txt" \
  --write-db --local-workspace

# 2. Extrair domínios únicos de terceiros
sqlite3 enumdns_ctrl.db "
SELECT DISTINCT target 
FROM results 
WHERE target LIKE '%.com' 
AND target NOT LIKE '%$COMPANY_DOMAIN%'
AND result_type IN ('CNAME', 'MX', 'NS')
" > "$OUTPUT_DIR/third_party_services.txt"

# 3. Categorizar por tipo de serviço
grep -i "office365\|outlook\|microsoft" "$OUTPUT_DIR/third_party_services.txt" > "$OUTPUT_DIR/microsoft_services.txt"
grep -i "google\|gmail\|ghs" "$OUTPUT_DIR/third_party_services.txt" > "$OUTPUT_DIR/google_services.txt"
grep -i "aws\|amazon\|cloudfront" "$OUTPUT_DIR/third_party_services.txt" > "$OUTPUT_DIR/aws_services.txt"
grep -i "azure\|windows" "$OUTPUT_DIR/third_party_services.txt" > "$OUTPUT_DIR/azure_services.txt"

echo "[+] Análise completa em $OUTPUT_DIR/"
```

### 2. Detecção de Phishing e Typosquatting

```bash
#!/bin/bash
# typosquatting_check.sh

DOMAIN="$1"
WORDLIST="/usr/share/wordlists/typos.txt"

# Gerar variações do domínio
python3 << EOF
import sys
domain = "$DOMAIN"
base = domain.split('.')[0]
tld = '.'.join(domain.split('.')[1:])

variations = []

# Substituições comuns
subs = {'o': '0', 'i': '1', 'l': '1', 'e': '3', 's': '5'}
for char, replacement in subs.items():
    if char in base:
        variations.append(base.replace(char, replacement) + '.' + tld)

# Adições comuns
additions = ['app', 'secure', 'login', 'mail', 'www']
for add in additions:
    variations.append(add + base + '.' + tld)
    variations.append(base + add + '.' + tld)

for var in set(variations):
    print(var)
EOF > typo_domains.txt

# Verificar quais existem
enumdns resolve file -L typo_domains.txt \
  --write-db \
  -o potential_typosquatting.txt

echo "[!] Verificar potential_typosquatting.txt para domínios suspeitos"
```

### 3. Análise de Takedown/Shutdown

```bash
#!/bin/bash
# infrastructure_monitoring.sh

DOMAIN="$1"
BASELINE_DB="baseline_$(echo $DOMAIN | tr '.' '_').db"
CURRENT_DB="current_$(date +%Y%m%d)_$(echo $DOMAIN | tr '.' '_').db"

# Criar baseline se não existir
if [[ ! -f "$BASELINE_DB" ]]; then
    echo "[+] Criando baseline para $DOMAIN"
    enumdns recon -d "$DOMAIN" \
      --write-db-uri "sqlite:///$BASELINE_DB" \
      --disable-control-db
    exit 0
fi

# Scan atual
echo "[+] Realizando scan atual de $DOMAIN"
enumdns recon -d "$DOMAIN" \
  --write-db-uri "sqlite:///$CURRENT_DB" \
  --disable-control-db

# Comparar resultados
echo "[+] Comparando com baseline..."

# Hosts que não respondem mais
sqlite3 << EOF
ATTACH DATABASE '$BASELINE_DB' AS baseline;
ATTACH DATABASE '$CURRENT_DB' AS current;

.headers on
.mode column

SELECT 'HOSTS_DOWN' as status, baseline.fqdn, baseline.ipv4, baseline.result_type
FROM baseline.results baseline
LEFT JOIN current.results current ON baseline.fqdn = current.fqdn AND baseline.result_type = current.result_type
WHERE baseline.exists = 1 
AND current.fqdn IS NULL;

SELECT 'NEW_HOSTS' as status, current.fqdn, current.ipv4, current.result_type  
FROM current.results current
LEFT JOIN baseline.results baseline ON current.fqdn = baseline.fqdn AND current.result_type = baseline.result_type
WHERE current.exists = 1
AND baseline.fqdn IS NULL;
EOF
```

## Troubleshooting Avançado

### Análise de Logs

```bash
# Habilitar logging completo
enumdns recon -d example.com -D --log-scan-errors 2>&1 | tee enumdns.log

# Analisar padrões de erro
grep "Error" enumdns.log | sort | uniq -c | sort -nr

# Verificar timeouts
grep -i "timeout" enumdns.log | wc -l

# Analisar desempenho DNS
grep "DNS request" enumdns.log | awk '{print $NF}' | sort -n | tail -10
```

### Problemas de Conectividade

```bash
# Testar conectividade DNS básica
dig @8.8.8.8 google.com

# Testar DNS via TCP
dig @8.8.8.8 +tcp google.com

# Verificar proxy SOCKS
curl --socks5 127.0.0.1:1080 http://ifconfig.me

# Testar resolução específica
enumdns recon -d google.com -s 8.8.8.8 -D
```

### Depuração de Banco de Dados

```sql
-- Verificar integridade do banco
PRAGMA integrity_check;

-- Estatísticas da base
SELECT 
    result_type,
    COUNT(*) as count,
    COUNT(CASE WHEN exists = 1 THEN 1 END) as existing,
    COUNT(CASE WHEN failed = 1 THEN 1 END) as failed
FROM results 
GROUP BY result_type;

-- Hosts com mais registros
SELECT fqdn, COUNT(*) as records 
FROM results 
WHERE exists = 1 
GROUP BY fqdn 
ORDER BY records DESC 
LIMIT 20;

-- Produtos cloud mais comuns
SELECT cloud_product, COUNT(*) as count 
FROM results 
WHERE cloud_product IS NOT NULL 
GROUP BY cloud_product 
ORDER BY count DESC;
```

## Contribuição

### Como Contribuir
1. Fork o repositório
2. Crie uma branch para sua feature: `git checkout -b feature/nova-funcionalidade`
3. Implemente mudanças com testes adequados
4. Faça commit seguindo convenções: `git commit -m "feat: adiciona nova funcionalidade"`
5. Push para sua branch: `git push origin feature/nova-funcionalidade`
6. Abra Pull Request com descrição detalhada

### Padrões de Código
- Go fmt obrigatório antes de commit
- Documentação GoDoc para funções públicas
- Testes unitários para lógica crítica
- Error handling adequado e consistente
- Logging estruturado com níveis apropriados

### Estrutura de Commits
```
feat: nova funcionalidade
fix: correção de bug
docs: atualização de documentação
style: formatação de código
refactor: refatoração sem mudança de funcionalidade
test: adição ou correção de testes
chore: tarefas de manutenção
```

### Testes

```bash
# Executar todos os testes
go test ./...

# Testes com coverage
go test -v -cover ./...

# Testes de integração
go test -tags=integration ./...

# Benchmark
go test -bench=. ./...
```

---

## Apêndices

### A. Wordlists Recomendadas

```bash
# Wordlists populares para DNS brute-force
/usr/share/wordlists/seclists/Discovery/DNS/subdomains-top1million-5000.txt
/usr/share/wordlists/seclists/Discovery/DNS/fierce-hostlist.txt
/usr/share/wordlists/seclists/Discovery/DNS/dns-Jhaddix.txt

# Wordlists especializadas
/usr/share/wordlists/seclists/Discovery/DNS/tlds.txt
/usr/share/wordlists/seclists/Discovery/DNS/sortuniq-hosts.txt
```

### B. Configuração de Elasticsearch

```json
PUT /enumdns
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "fqdn": {"type": "keyword"},
      "result_type": {"type": "keyword"},
      "ipv4": {"type": "ip"},
      "ipv6": {"type": "ip"},
      "target": {"type": "keyword"},
      "cloud_product": {"type": "keyword"},
      "saas_product": {"type": "keyword"},
      "datacenter": {"type": "keyword"},
      "dc": {"type": "boolean"},
      "gc": {"type": "boolean"},
      "probed_at": {"type": "date"}
    }
  }
}
```

### C. Dashboard Kibana

```json
{
  "version": "8.15.0",
  "objects": [
    {
      "id": "enumdns-overview",
      "type": "dashboard",
      "attributes": {
        "title": "EnumDNS Overview",
        "description": "Dashboard principal para análise de resultados EnumDNS",
        "panelsJSON": "[{\"version\":\"8.15.0\",\"gridData\":{\"x\":0,\"y\":0,\"w\":24,\"h\":15,\"i\":\"1\"},\"panelIndex\":\"1\",\"embeddableConfig\":{},\"panelRefName\":\"panel_1\"}]"
      }
    }
  ]
}
```

### D. Configuração systemd

```ini
# /etc/systemd/system/enumdns-monitor.service
[Unit]
Description=EnumDNS Continuous Monitor
After=network.target

[Service]
Type=simple
User=enumdns
Group=enumdns
WorkingDirectory=/opt/enumdns
ExecStart=/opt/enumdns/monitor_dns.sh example.com 24
Restart=always
RestartSec=300

[Install]
WantedBy=multi-user.target
```

### E. Configuração Nginx (Proxy para API)

```nginx
server {
    listen 80;
    server_name enumdns-api.example.com;
    
    location / {
        proxy_pass http://127.0.0.1:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/m;
    limit_req zone=api burst=5 nodelay;
}
```

---

**Autor**: Helvio Junior (M4v3r1ck)  
**Repositório**: https://github.com/helviojunior/enumdns  
**Versão da Documentação**: 1.0  
**Data**: Janeiro 2025

**Disclaimer**: Esta ferramenta é destinada exclusivamente para uso educacional e em ambientes onde você possui autorização explícita/legal para realizar testes de segurança. O uso inadequado desta ferramenta em sistemas sem autorização pode violar leis locais e internacionais. O autor não se responsabiliza pelo uso indevido da ferramenta.

**Licença**: Consulte o arquivo LICENSE no repositório do projeto para informações sobre licenciamento e termos de uso.