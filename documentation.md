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
9. [Integração com Sistemas Externos](#integração-com-sistemas-externos)
10. [Desenvolvimento](#desenvolvimento)
11. [Considerações de Segurança](#considerações-de-segurança)

## Visão Geral

O **EnumDNS** é uma ferramenta modular de reconhecimento DNS desenvolvida em Go, projetada para profissionais de segurança cibernética realizarem enumeração DNS abrangente e análise de infraestrutura. A ferramenta oferece múltiplos métodos de descoberta de hosts e pode identificar automaticamente provedores de nuvem, produtos SaaS e datacenters.

### Principais Características
- **Modular**: Suporta diferentes tipos de enumeração (brute-force, reconhecimento, resolução, análise avançada)
- **Multi-plataforma**: Funciona em Linux, Windows e macOS
- **Flexível**: Múltiplos formatos de saída (texto, JSON, CSV, SQLite, Elasticsearch)
- **Escalável**: Suporta processamento paralelo com goroutines
- **Integrado**: Compatível com BloodHound e crt.sh
- **Proxy Support**: Suporta SOCKS4/SOCKS5
- **Análise de Ameaças**: Detecção de typosquatting, bitsquatting e ataques homográficos

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

#### 4. **Análise de Ameaças (threat-analysis)** ⭐ *ATUALIZADO*
- **Typosquatting**: Detecção de domínios com erros de digitação baseados em adjacência de teclado
- **Bitsquatting**: Identificação de domínios com alterações de bits
- **Ataques Homográficos**: Detecção de caracteres visualmente similares
- **Análise de Similaridade**: Cálculo de proximidade com domínios legítimos
- **Score de Ameaça**: Classificação automática de risco
- **Indicadores de Ameaça**: Identificação de padrões suspeitos (TLDs, Unicode tricks)

#### 5. **Relatórios**
- Conversão entre formatos (SQLite ↔ JSON Lines ↔ Texto)
- Sincronização com Elasticsearch

### Recursos Avançados

- **Detecção de Active Directory**: Identifica DCs e GCs automaticamente
- **Identificação de Cloud**: Reconhece AWS, Azure, GCP, CloudFlare, etc.
- **Reverse DNS**: Resolução PTR automática
- **ASN Detection**: Identifica ASN e informações de rede
- **Controle de Duplicatas**: Evita varreduras desnecessárias
- **Análise de Risco**: Sistema de pontuação para variações de domínio

## Arquitetura e Estrutura

### Estrutura de Diretórios

```
enumdns/
├── cmd/                    # Comandos CLI (Cobra)
│   ├── advanced.go        # Comando de análise avançada de ameaças
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
│   ├── advanced/          # Sistema de análise avançada de ameaças
│   │   ├── analyzer.go    # Analisador de risco e similaridade
│   │   ├── generator.go   # Gerador de variações
│   │   ├── options.go     # Opções e configurações
│   │   └── techniques.go  # Técnicas de geração de variações
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
- **Advanced**: Engine para análise de ameaças e variações de domínio
- Gerenciam workers (goroutines) e coordenam a execução

#### 2. **Advanced Threat Analysis System**
- **VariationGenerator**: Gera variações de domínio usando múltiplas técnicas
- **RiskAnalyzer**: Analisa similaridade e calcula scores de ameaça
- **TechniqueEngine**: Implementa algoritmos específicos de geração
- **ThreatClassifier**: Classifica e categoriza ameaças identificadas

#### 3. **Readers (Leitores de Entrada)**
- **FileReader**: Lê wordlists e listas de domínios
- **CrtShReader**: Integração com crt.sh
- Suportam diferentes formatos de entrada

#### 4. **Writers (Escritores de Saída)**
- **DbWriter**: SQLite, PostgreSQL, MySQL
- **JsonWriter**: JSON Lines
- **CsvWriter**: Formato CSV
- **ElasticWriter**: Elasticsearch
- **TextWriter**: Relatórios legíveis
- **StdoutWriter**: Saída no terminal

#### 5. **Models (Modelos de Dados)**
- **Result**: Resultado principal de DNS
- **ThreatResult**: Resultado de análise de ameaças (herda de Result)
- **FQDNData**: Dados de FQDN descobertos
- **ASN**: Informações de Sistema Autônomo
- **ASNIpDelegate**: Delegações de IP por ASN

## Instalação

### Linux

```bash
# Instalação automática
apt install curl jq

url=$(curl -s https://api.github.com/repos/bob-reis/enumdns/releases | jq -r '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match("linux-amd64.tar.gz$")) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

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

url=$(curl -s https://api.github.com/repos/bob-reis/enumdns/releases | jq -r --arg filename "darwin-${arch}.tar.gz\$" '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match($filename)) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

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
# Download latest bob-reis/enumdns release from github
function Invoke-Downloadenumdns {

    $repo = "bob-reis/enumdns"
    
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
git clone https://github.com/bob-reis/enumdns.git
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

### Comando Threat-Analysis (Análise de Ameaças) ⭐ *ATUALIZADO*

```bash
# Análise completa de ameaças para um domínio
enumdns threat-analysis -d example.com --all-techniques -o threats.txt

# Análise específica de typosquatting
enumdns threat-analysis -d example.com --typosquatting --write-db

# Análise de múltiplos domínios com configurações avançadas
enumdns threat-analysis -L domains.txt --bitsquatting --homographic \
  --max-variations 500 --target-tlds com,net,org,co,io --write-jsonl

# Análise de ataques homográficos
enumdns threat-analysis -d bank.com --homographic --write-elasticsearch-uri http://siem:9200/threats

# Análise completa com todas as técnicas
enumdns threat-analysis -d corporate.com --all-techniques \
  --max-variations 1000 --write-db --write-csv
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

### Sistema de Análise Avançada de Ameaças

#### Técnicas Implementadas

**1. Typosquatting (Erro de Digitação)**
- Baseado em adjacência de teclado QWERTY
- Simula erros comuns de digitação
- Score de risco: Alto (0.8)

```go
// Exemplo de mapeamento de teclas adjacentes
keyboardMap := map[rune][]rune{
    'q': {'w', 'a'}, 'w': {'q', 'e', 's'}, 
    'e': {'w', 'r', 'd'}, ...
}
```

**2. Bitsquatting (Alteração de Bits)**
- Altera bits individuais dos caracteres
- Simula erros de hardware/transmissão
- Score de risco: Médio (0.6)

**3. Ataques Homográficos**
- Usa caracteres visualmente similares
- Detecta Unicode tricks e variações acentuadas
- Score de risco: Alto (0.9)

```go
// Exemplo de caracteres homográficos
homographicMap := map[rune][]rune{
    'a': {'à', 'á', 'ä', 'â', 'ā', 'α'},
    'e': {'è', 'é', 'ê', 'ë', 'ē'},
    'o': {'ò', 'ó', 'ô', 'õ', 'ö', '0'},
}
```

**4. Técnicas Adicionais**
- **Insertion**: Inserção de caracteres comuns
- **Deletion**: Remoção de caracteres
- **Transposition**: Troca de caracteres adjacentes
- **TLD Variation**: Variações de domínios de topo
- **Subdomain Pattern**: Padrões comuns de phishing

#### Analisador de Risco

```go
type RiskAnalyzer struct {
    BaseDomain string
}

type AnalysisResult struct {
    Variation   Variation
    Similarity  float64      // 0-1 (Levenshtein similarity)
    ThreatScore float64      // 0-1 (risco calculado)
    Indicators  []string     // ["suspicious_tld", "phishing_pattern"]
}
```

#### Indicadores de Ameaça
- **suspicious_tld**: TLDs comumente usados para phishing
- **phishing_pattern**: Palavras-chave suspeitas (secure, login, verify)
- **high_similarity**: Similaridade > 80% com domínio original
- **unicode_tricks**: Uso de caracteres não-ASCII

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

-- Nova tabela para análise de ameaças
CREATE TABLE threat_results (
    id INTEGER PRIMARY KEY,
    fqdn TEXT,
    result_type TEXT,
    technique TEXT,
    confidence REAL,
    risk TEXT,
    base_domain TEXT,
    similarity REAL,
    probed_at DATETIME
);
```

### CSV
```csv
FQDN,RType,IPv4,IPv6,Target,CloudProduct,DC,GC,ProbedAt
example.com,A,192.0.2.1,,,,,false,false,2025-01-15T10:30:00Z
www.example.com,CNAME,,,example.com,,false,false,2025-01-15T10:30:01Z
```

## Integração com Sistemas Externos

### APIs e Conectores

#### API REST Wrapper (Python Flask)

```python
#!/usr/bin/env python3
# enumdns_api.py - Wrapper REST para EnumDNS

from flask import Flask, request, jsonify
import subprocess
import json
import tempfile
import os

app = Flask(__name__)

@app.route('/api/v1/threat-analysis', methods=['POST'])
def threat_analysis():
    """Endpoint para análise de ameaças"""
    data = request.get_json()
    domain = data.get('domain')
    techniques = data.get('techniques', ['all-techniques'])
    max_variations = data.get('max_variations', 1000)
    
    if not domain:
        return jsonify({'error': 'Domain required'}), 400
    
    try:
        # Criar arquivo temporário para resultados
        with tempfile.NamedTemporaryFile(mode='w', suffix='.jsonl', delete=False) as f:
            temp_file = f.name
        
        # Executar EnumDNS Advanced
        cmd = [
            'enumdns', 'advanced',
            '-d', domain,
            '--all-techniques' if 'all-techniques' in techniques else '--typosquatting',
            '--max-variations', str(max_variations),
            '--write-jsonl-file', temp_file,
            '-q'
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'error': result.stderr}), 500
        
        # Ler resultados
        threats = []
        with open(temp_file, 'r') as f:
            for line in f:
                if line.strip():
                    threat = json.loads(line)
                    threats.append(threat)
        
        # Cleanup
        os.unlink(temp_file)
        
        # Classificar por risco
        high_risk = [t for t in threats if t.get('cloud_product') or 'suspicious' in str(t)]
        
        return jsonify({
            'domain': domain,
            'threats': threats,
            'high_risk_count': len(high_risk),
            'total_variations': len(threats)
        })
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/v1/bulk-analysis', methods=['POST'])
def bulk_analysis():
    """Análise em lote de múltiplos domínios"""
    data = request.get_json()
    domains = data.get('domains', [])
    
    if not domains:
        return jsonify({'error': 'Domains list required'}), 400
    
    results = {}
    for domain in domains:
        # Implementar lógica similar ao endpoint anterior
        pass
    
    return jsonify(results)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

#### Integração com SIEM (Splunk)

```python
#!/usr/bin/env python3
# splunk_enumdns.py - Integração com Splunk

import splunklib.client as client
import json
import subprocess
import logging

class EnumDNSSplunkConnector:
    def __init__(self, splunk_host, splunk_port, username, password):
        self.service = client.connect(
            host=splunk_host,
            port=splunk_port,
            username=username,
            password=password
        )
        self.index = self.service.indexes['enumdns']
    
    def analyze_domain_threats(self, domain, techniques=['all-techniques']):
        """Executa análise de ameaças e envia para Splunk"""
        try:
            # Executar EnumDNS
            cmd = ['enumdns', 'advanced', '-d', domain, '--all-techniques', '--write-jsonl', '-q']
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode != 0:
                logging.error(f"EnumDNS failed: {result.stderr}")
                return False
            
            # Processar resultados e enviar para Splunk
            for line in result.stdout.split('\n'):
                if line.strip():
                    try:
                        threat_data = json.loads(line)
                        # Enriquecer com metadados
                        threat_data['source'] = 'enumdns_advanced'
                        threat_data['analysis_time'] = time.time()
                        
                        # Enviar para Splunk
                        self.index.submit(json.dumps(threat_data))
                        
                    except json.JSONDecodeError:
                        continue
            
            return True
            
        except Exception as e:
            logging.error(f"Error in threat analysis: {e}")
            return False
    
    def create_alerts(self, domain, high_risk_threshold=0.8):
        """Cria alertas automáticos para ameaças de alto risco"""
        search_query = f'''
        search index=enumdns source=enumdns_advanced base_domain="{domain}"
        | where confidence > {high_risk_threshold}
        | eval risk_level=case(
            confidence > 0.9, "critical",
            confidence > 0.8, "high",
            confidence > 0.6, "medium",
            1=1, "low"
        )
        | table fqdn, technique, confidence, risk_level, similarity
        | sort -confidence
        '''
        
        return self.service.jobs.create(search_query)
```

#### Integração com Elasticsearch e Kibana

```python
#!/usr/bin/env python3
# elastic_enumdns.py - Integração avançada com Elasticsearch

from elasticsearch import Elasticsearch
import json
import subprocess
from datetime import datetime

class EnumDNSElasticConnector:
    def __init__(self, elastic_hosts, index_prefix='enumdns'):
        self.es = Elasticsearch(elastic_hosts)
        self.index_prefix = index_prefix
        self.setup_indices()
    
    def setup_indices(self):
        """Configura índices otimizados para EnumDNS"""
        
        # Índice principal para resultados DNS
        dns_mapping = {
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
                    "probed_at": {"type": "date"},
                    "asn": {"type": "long"},
                    "location": {"type": "geo_point"}
                }
            }
        }
        
        # Índice para análise de ameaças
        threat_mapping = {
            "mappings": {
                "properties": {
                    "fqdn": {"type": "keyword"},
                    "base_domain": {"type": "keyword"},
                    "technique": {"type": "keyword"},
                    "confidence": {"type": "float"},
                    "similarity": {"type": "float"},
                    "risk_level": {"type": "keyword"},
                    "threat_indicators": {"type": "keyword"},
                    "probed_at": {"type": "date"},
                    "tld": {"type": "keyword"},
                    "registrar": {"type": "keyword"},
                    "creation_date": {"type": "date"}
                }
            }
        }
        
        # Criar índices se não existirem
        dns_index = f"{self.index_prefix}-dns"
        threat_index = f"{self.index_prefix}-threats"
        
        if not self.es.indices.exists(index=dns_index):
            self.es.indices.create(index=dns_index, body=dns_mapping)
        
        if not self.es.indices.exists(index=threat_index):
            self.es.indices.create(index=threat_index, body=threat_mapping)
    
    def ingest_dns_results(self, domain):
        """Ingere resultados de reconhecimento DNS"""
        cmd = ['enumdns', 'recon', '-d', domain, '--write-jsonl', '-q']
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        actions = []
        for line in result.stdout.split('\n'):
            if line.strip():
                try:
                    dns_data = json.loads(line)
                    action = {
                        "_index": f"{self.index_prefix}-dns",
                        "_source": dns_data
                    }
                    actions.append(action)
                except json.JSONDecodeError:
                    continue
        
        if actions:
            from elasticsearch.helpers import bulk
            bulk(self.es, actions)
    
    def ingest_threat_analysis(self, domain):
        """Ingere resultados de análise de ameaças"""
        cmd = ['enumdns', 'advanced', '-d', domain, '--all-techniques', '--write-jsonl', '-q']
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        actions = []
        for line in result.stdout.split('\n'):
            if line.strip():
                try:
                    threat_data = json.loads(line)
                    # Enriquecer dados
                    threat_data['risk_level'] = self.calculate_risk_level(threat_data)
                    threat_data['tld'] = threat_data['fqdn'].split('.')[-1]
                    
                    action = {
                        "_index": f"{self.index_prefix}-threats",
                        "_source": threat_data
                    }
                    actions.append(action)
                except json.JSONDecodeError:
                    continue
        
        if actions:
            from elasticsearch.helpers import bulk
            bulk(self.es, actions)
    
    def calculate_risk_level(self, threat_data):
        """Calcula nível de risco baseado nos dados"""
        confidence = threat_data.get('confidence', 0)
        similarity = threat_data.get('similarity', 0)
        
        if confidence > 0.9 or similarity > 0.95:
            return 'critical'
        elif confidence > 0.7 or similarity > 0.8:
            return 'high'
        elif confidence > 0.5 or similarity > 0.6:
            return 'medium'
        else:
            return 'low'
    
    def create_kibana_dashboards(self):
        """Cria dashboards pré-configurados no Kibana"""
        
        # Dashboard de Overview
        overview_dashboard = {
            "title": "EnumDNS - Overview",
            "description": "Dashboard principal para análise de resultados EnumDNS",
            "panelsJSON": json.dumps([
                {
                    "title": "DNS Records by Type",
                    "type": "pie",
                    "query": {
                        "index": f"{self.index_prefix}-dns",
                        "aggregations": {
                            "types": {
                                "terms": {"field": "result_type"}
                            }
                        }
                    }
                },
                {
                    "title": "Threat Distribution",
                    "type": "histogram",
                    "query": {
                        "index": f"{self.index_prefix}-threats",
                        "aggregations": {
                            "risk_levels": {
                                "terms": {"field": "risk_level"}
                            }
                        }
                    }
                }
            ])
        }
        
        return overview_dashboard
```

#### Integração com APIs de Threat Intelligence

```python
#!/usr/bin/env python3
# threat_intel_integration.py

import requests
import json
from datetime import datetime

class ThreatIntelEnricher:
    def __init__(self, virustotal_api_key=None, urlvoid_api_key=None):
        self.vt_api_key = virustotal_api_key
        self.urlvoid_api_key = urlvoid_api_key
    
    def enrich_enumdns_results(self, enumdns_results):
        """Enriquece resultados do EnumDNS com dados de threat intelligence"""
        enriched_results = []
        
        for result in enumdns_results:
            fqdn = result.get('fqdn')
            if not fqdn:
                continue
            
            # Enriquecer com VirusTotal
            vt_data = self.check_virustotal(fqdn)
            if vt_data:
                result['virustotal'] = vt_data
            
            # Enriquecer com URLVoid
            urlvoid_data = self.check_urlvoid(fqdn)
            if urlvoid_data:
                result['urlvoid'] = urlvoid_data
            
            # Calcular score de reputação
            result['reputation_score'] = self.calculate_reputation(result)
            
            enriched_results.append(result)
        
        return enriched_results
    
    def check_virustotal(self, domain):
        """Consulta VirusTotal para informações do domínio"""
        if not self.vt_api_key:
            return None
        
        url = f"https://www.virustotal.com/vtapi/v2/domain/report"
        params = {
            'apikey': self.vt_api_key,
            'domain': domain
        }
        
        try:
            response = requests.get(url, params=params, timeout=10)
            if response.status_code == 200:
                data = response.json()
                return {
                    'detected_urls': data.get('detected_urls', []),
                    'detected_samples': data.get('detected_samples', []),
                    'reputation': data.get('reputation', 0)
                }
        except Exception as e:
            print(f"Error checking VirusTotal for {domain}: {e}")
        
        return None
    
    def check_urlvoid(self, domain):
        """Consulta URLVoid para verificação de reputação"""
        if not self.urlvoid_api_key:
            return None
        
        url = f"http://api.urlvoid.com/v1/pay-as-you-go/"
        params = {
            'key': self.urlvoid_api_key,
            'host': domain
        }
        
        try:
            response = requests.get(url, params=params, timeout=10)
            if response.status_code == 200:
                # Processar resposta XML do URLVoid
                # (implementação específica dependeria do formato)
                return {'status': 'checked'}
        except Exception as e:
            print(f"Error checking URLVoid for {domain}: {e}")
        
        return None
    
    def calculate_reputation(self, result):
        """Calcula score de reputação baseado em múltiplas fontes"""
        score = 50  # Score base
        
        # Ajustar baseado em VirusTotal
        vt_data = result.get('virustotal', {})
        if vt_data.get('detected_urls'):
            score -= 20
        if vt_data.get('reputation', 0) < 0:
            score -= 15
        
        # Ajustar baseado em análise de ameaças do EnumDNS
        if result.get('technique'):
            score -= 10
        if result.get('confidence', 0) > 0.8:
            score -= 20
        
        # Ajustar baseado em produtos cloud/SaaS
        if result.get('cloud_product'):
            score += 10
        if result.get('saas_product'):
            score += 5
        
        return max(0, min(100, score))
```

### Integração com Bancos de Dados

#### PostgreSQL Schema Otimizado

```sql
-- Schema otimizado para PostgreSQL
CREATE SCHEMA enumdns;

-- Tabela principal de resultados DNS
CREATE TABLE enumdns.dns_results (
    id BIGSERIAL PRIMARY KEY,
    test_id VARCHAR(50) NOT NULL,
    fqdn VARCHAR(255) NOT NULL,
    result_type VARCHAR(20) NOT NULL,
    ipv4 INET,
    ipv6 INET,
    target VARCHAR(255),
    ptr VARCHAR(255),
    txt TEXT,
    cloud_product VARCHAR(100),
    saas_product VARCHAR(100),
    datacenter VARCHAR(100),
    asn BIGINT,
    dc BOOLEAN DEFAULT FALSE,
    gc BOOLEAN DEFAULT FALSE,
    exists BOOLEAN DEFAULT TRUE,
    failed BOOLEAN DEFAULT FALSE,
    failed_reason TEXT,
    probed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para análise de ameaças
CREATE TABLE enumdns.threat_analysis (
    id BIGSERIAL PRIMARY KEY,
    base_domain VARCHAR(255) NOT NULL,
    variation_fqdn VARCHAR(255) NOT NULL,
    technique VARCHAR(50) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    similarity DECIMAL(3,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    threat_indicators TEXT[],
    tld VARCHAR(10),
    registrar VARCHAR(100),
    creation_date DATE,
    reputation_score INTEGER,
    is_registered BOOLEAN,
    probed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para informações ASN
CREATE TABLE enumdns.asn_info (
    asn BIGINT PRIMARY KEY,
    rir_name VARCHAR(20) NOT NULL,
    country_code CHAR(2),
    organization TEXT,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para delegações IP-ASN
CREATE TABLE enumdns.asn_delegations (
    id BIGSERIAL PRIMARY KEY,
    rir_name VARCHAR(20) NOT NULL,
    country_code CHAR(2),
    subnet CIDR NOT NULL,
    addresses INTEGER NOT NULL,
    date_allocated DATE,
    asn BIGINT REFERENCES enumdns.asn_info(asn),
    status VARCHAR(20),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para performance
CREATE INDEX idx_dns_results_fqdn ON enumdns.dns_results(fqdn);
CREATE INDEX idx_dns_results_type ON enumdns.dns_results(result_type);
CREATE INDEX idx_dns_results_probed_at ON enumdns.dns_results(probed_at);
CREATE INDEX idx_dns_results_ipv4 ON enumdns.dns_results(ipv4) WHERE ipv4 IS NOT NULL;
CREATE INDEX idx_dns_results_cloud ON enumdns.dns_results(cloud_product) WHERE cloud_product IS NOT NULL;

CREATE INDEX idx_threat_base_domain ON enumdns.threat_analysis(base_domain);
CREATE INDEX idx_threat_technique ON enumdns.threat_analysis(technique);
CREATE INDEX idx_threat_risk_level ON enumdns.threat_analysis(risk_level);
CREATE INDEX idx_threat_confidence ON enumdns.threat_analysis(confidence);

-- Views úteis
CREATE VIEW enumdns.high_risk_threats AS
SELECT 
    base_domain,
    variation_fqdn,
    technique,
    confidence,
    similarity,
    risk_level,
    threat_indicators,
    probed_at
FROM enumdns.threat_analysis
WHERE risk_level IN ('critical', 'high')
ORDER BY confidence DESC, similarity DESC;

CREATE VIEW enumdns.domain_summary AS
SELECT 
    fqdn,
    COUNT(*) as total_records,
    COUNT(CASE WHEN result_type = 'A' THEN 1 END) as a_records,
    COUNT(CASE WHEN result_type = 'AAAA' THEN 1 END) as aaaa_records,
    COUNT(CASE WHEN result_type = 'CNAME' THEN 1 END) as cname_records,
    COUNT(CASE WHEN cloud_product IS NOT NULL THEN 1 END) as cloud_services,
    MAX(probed_at) as last_scan
FROM enumdns.dns_results
WHERE exists = TRUE
GROUP BY fqdn;

-- Funções úteis
CREATE OR REPLACE FUNCTION enumdns.get_threat_stats(domain_name VARCHAR)
RETURNS TABLE(
    technique VARCHAR,
    count BIGINT,
    avg_confidence DECIMAL,
    max_confidence DECIMAL
) AS $
BEGIN
    RETURN QUERY
    SELECT 
        t.technique,
        COUNT(*) as count,
        AVG(t.confidence) as avg_confidence,
        MAX(t.confidence) as max_confidence
    FROM enumdns.threat_analysis t
    WHERE t.base_domain = domain_name
    GROUP BY t.technique
    ORDER BY avg_confidence DESC;
END;
$ LANGUAGE plpgsql;

-- Trigger para atualização automática de timestamps
CREATE OR REPLACE FUNCTION enumdns.update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER update_asn_info_updated_at 
    BEFORE UPDATE ON enumdns.asn_info
    FOR EACH ROW EXECUTE FUNCTION enumdns.update_updated_at_column();
```

#### MongoDB Schema (NoSQL Alternative)

```javascript
// MongoDB Collections para EnumDNS

// Collection: dns_results
db.createCollection("dns_results", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["test_id", "fqdn", "result_type", "probed_at"],
         properties: {
            test_id: { bsonType: "string" },
            fqdn: { bsonType: "string" },
            result_type: { bsonType: "string" },
            ipv4: { bsonType: "string" },
            ipv6: { bsonType: "string" },
            target: { bsonType: "string" },
            cloud_product: { bsonType: "string" },
            saas_product: { bsonType: "string" },
            datacenter: { bsonType: "string" },
            asn: { bsonType: "long" },
            dc: { bsonType: "bool" },
            gc: { bsonType: "bool" },
            exists: { bsonType: "bool" },
            probed_at: { bsonType: "date" }
         }
      }
   }
});

// Collection: threat_analysis
db.createCollection("threat_analysis", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["base_domain", "variation_fqdn", "technique", "confidence"],
         properties: {
            base_domain: { bsonType: "string" },
            variation_fqdn: { bsonType: "string" },
            technique: { bsonType: "string" },
            confidence: { bsonType: "double", minimum: 0, maximum: 1 },
            similarity: { bsonType: "double", minimum: 0, maximum: 1 },
            risk_level: { enum: ["low", "medium", "high", "critical"] },
            threat_indicators: { bsonType: "array", items: { bsonType: "string" } },
            probed_at: { bsonType: "date" }
         }
      }
   }
});

// Índices para performance
db.dns_results.createIndex({ "fqdn": 1 });
db.dns_results.createIndex({ "result_type": 1 });
db.dns_results.createIndex({ "probed_at": -1 });
db.dns_results.createIndex({ "ipv4": 1 }, { sparse: true });
db.dns_results.createIndex({ "cloud_product": 1 }, { sparse: true });

db.threat_analysis.createIndex({ "base_domain": 1 });
db.threat_analysis.createIndex({ "technique": 1 });
db.threat_analysis.createIndex({ "risk_level": 1 });
db.threat_analysis.createIndex({ "confidence": -1 });

// Aggregation pipelines úteis
const getThreatStatsByDomain = (domain) => {
    return db.threat_analysis.aggregate([
        { $match: { base_domain: domain } },
        {
            $group: {
                _id: "$technique",
                count: { $sum: 1 },
                avg_confidence: { $avg: "$confidence" },
                max_confidence: { $max: "$confidence" },
                high_risk_count: {
                    $sum: { $cond: [{ $gte: ["$confidence", 0.8] }, 1, 0] }
                }
            }
        },
        { $sort: { avg_confidence: -1 } }
    ]);
};
```

### Integração com Orquestradores

#### Docker Compose para Ambiente Completo

```yaml
# docker-compose.yml - Ambiente completo EnumDNS
version: '3.8'

services:
  enumdns:
    build: .
    volumes:
      - ./data:/data
      - ./wordlists:/wordlists
      - ./config:/config
    environment:
      - ENUMDNS_DB_PATH=/data/enumdns.db
      - ENUMDNS_ELASTIC_URI=http://elasticsearch:9200/enumdns
    networks:
      - enumdns_net
    depends_on:
      - postgres
      - elasticsearch
  
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: enumdns
      POSTGRES_USER: enumdns
      POSTGRES_PASSWORD: securepassword
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - enumdns_net
  
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.15.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
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
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - enumdns_net
  
  enumdns-api:
    build: ./api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://enumdns:securepassword@postgres:5432/enumdns
      - ELASTIC_URL=http://elasticsearch:9200
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - elasticsearch
      - redis
    networks:
      - enumdns_net
  
  enumdns-worker:
    build: .
    command: ["python", "/app/worker.py"]
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - CELERY_BROKER_URL=redis://redis:6379
      - DATABASE_URL=postgresql://enumdns:securepassword@postgres:5432/enumdns
    depends_on:
      - redis
      - postgres
    networks:
      - enumdns_net
    deploy:
      replicas: 3

networks:
  enumdns_net:
    driver: bridge

volumes:
  postgres_data:
  es_data:
  redis_data:
```

#### Kubernetes Deployment

```yaml
# k8s/enumdns-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: enumdns-api
  namespace: security-tools
spec:
  replicas: 3
  selector:
    matchLabels:
      app: enumdns-api
  template:
    metadata:
      labels:
        app: enumdns-api
    spec:
      containers:
      - name: enumdns-api
        image: enumdns:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: enumdns-secrets
              key: database-url
        - name: ELASTIC_URL
          value: "http://elasticsearch:9200"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: enumdns-api-service
  namespace: security-tools
spec:
  selector:
    app: enumdns-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: enumdns-scheduled-scan
  namespace: security-tools
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: enumdns-scanner
            image: enumdns:latest
            command:
            - /bin/sh
            - -c
            - |
              enumdns threat-analysis -L /config/domains.txt --all-techniques \
                --write-db-uri "$DATABASE_URL" \
                --write-elasticsearch-uri "$ELASTIC_URL"
            env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: enumdns-secrets
                  key: database-url
            - name: ELASTIC_URL
              value: "http://elasticsearch:9200"
            volumeMounts:
            - name: config-volume
              mountPath: /config
          volumes:
          - name: config-volume
            configMap:
              name: enumdns-config
          restartPolicy: OnFailure
```

### Integração com Ferramentas de Automação

#### Ansible Playbook

```yaml
# ansible/enumdns-playbook.yml
---
- name: Deploy and configure EnumDNS
  hosts: security_servers
  become: yes
  vars:
    enumdns_version: "latest"
    enumdns_path: "/opt/enumdns"
    enumdns_user: "enumdns"
    enumdns_config_path: "/etc/enumdns"
    
  tasks:
    - name: Create enumdns user
      user:
        name: "{{ enumdns_user }}"
        system: yes
        home: "{{ enumdns_path }}"
        shell: /bin/bash
    
    - name: Create enumdns directories
      file:
        path: "{{ item }}"
        state: directory
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        mode: '0755'
      loop:
        - "{{ enumdns_path }}"
        - "{{ enumdns_config_path }}"
        - "{{ enumdns_path }}/logs"
        - "{{ enumdns_path }}/data"
        - "{{ enumdns_path }}/wordlists"
    
    - name: Download and install EnumDNS
      get_url:
        url: "https://github.com/bob-reis/enumdns/releases/download/{{ enumdns_version }}/enumdns-linux-amd64.tar.gz"
        dest: "/tmp/enumdns.tar.gz"
        mode: '0644'
    
    - name: Extract EnumDNS
      unarchive:
        src: "/tmp/enumdns.tar.gz"
        dest: "{{ enumdns_path }}"
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        remote_src: yes
    
    - name: Create symlink
      file:
        src: "{{ enumdns_path }}/enumdns"
        dest: "/usr/local/bin/enumdns"
        state: link
    
    - name: Install systemd service
      template:
        src: enumdns.service.j2
        dest: /etc/systemd/system/enumdns-api.service
        mode: '0644'
      notify: reload systemd
    
    - name: Configure EnumDNS
      template:
        src: enumdns.conf.j2
        dest: "{{ enumdns_config_path }}/enumdns.conf"
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        mode: '0600'
    
    - name: Start and enable EnumDNS service
      systemd:
        name: enumdns-api
        state: started
        enabled: yes
        daemon_reload: yes
  
  handlers:
    - name: reload systemd
      systemd:
        daemon_reload: yes
```

#### Terraform para Infraestrutura

```hcl
# terraform/main.tf - Infraestrutura EnumDNS na AWS
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC para EnumDNS
resource "aws_vpc" "enumdns_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name = "enumdns-vpc"
    Environment = var.environment
  }
}

# Subnets
resource "aws_subnet" "enumdns_private" {
  count             = 2
  vpc_id            = aws_vpc.enumdns_vpc.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]
  
  tags = {
    Name = "enumdns-private-${count.index + 1}"
    Environment = var.environment
  }
}

resource "aws_subnet" "enumdns_public" {
  count                   = 2
  vpc_id                  = aws_vpc.enumdns_vpc.id
  cidr_block              = "10.0.${count.index + 10}.0/24"
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  
  tags = {
    Name = "enumdns-public-${count.index + 1}"
    Environment = var.environment
  }
}

# EKS Cluster para EnumDNS
resource "aws_eks_cluster" "enumdns_cluster" {
  name     = "enumdns-cluster"
  role_arn = aws_iam_role.enumdns_cluster_role.arn
  version  = "1.28"

  vpc_config {
    subnet_ids = concat(aws_subnet.enumdns_private[*].id, aws_subnet.enumdns_public[*].id)
    endpoint_private_access = true
    endpoint_public_access  = true
  }

  depends_on = [
    aws_iam_role_policy_attachment.enumdns_cluster_AmazonEKSClusterPolicy
  ]
}

# RDS para PostgreSQL
resource "aws_db_instance" "enumdns_postgres" {
  identifier = "enumdns-postgres"
  
  engine         = "postgres"
  engine_version = "15.4"
  instance_class = "db.t3.medium"
  
  allocated_storage     = 100
  max_allocated_storage = 1000
  storage_encrypted     = true
  
  db_name  = "enumdns"
  username = "enumdns_admin"
  password = var.db_password
  
  vpc_security_group_ids = [aws_security_group.enumdns_db.id]
  db_subnet_group_name   = aws_db_subnet_group.enumdns.name
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  skip_final_snapshot = false
  final_snapshot_identifier = "enumdns-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  
  tags = {
    Name = "enumdns-postgres"
    Environment = var.environment
  }
}

# ElastiCache Redis para cache
resource "aws_elasticache_subnet_group" "enumdns_redis" {
  name       = "enumdns-redis-subnet-group"
  subnet_ids = aws_subnet.enumdns_private[*].id
}

resource "aws_elasticache_replication_group" "enumdns_redis" {
  replication_group_id       = "enumdns-redis"
  description                = "Redis cluster for EnumDNS caching"
  
  node_type                  = "cache.t3.micro"
  port                       = 6379
  parameter_group_name       = "default.redis7"
  
  num_cache_clusters         = 2
  automatic_failover_enabled = true
  
  subnet_group_name = aws_elasticache_subnet_group.enumdns_redis.name
  security_group_ids = [aws_security_group.enumdns_redis.id]
  
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  
  tags = {
    Name = "enumdns-redis"
    Environment = var.environment
  }
}

# OpenSearch para logs e análise
resource "aws_opensearch_domain" "enumdns_opensearch" {
  domain_name    = "enumdns-logs"
  engine_version = "OpenSearch_2.3"

  cluster_config {
    instance_type            = "t3.medium.search"
    instance_count           = 2
    dedicated_master_enabled = false
  }

  ebs_options {
    ebs_enabled = true
    volume_type = "gp3"
    volume_size = 100
  }

  vpc_options {
    subnet_ids         = aws_subnet.enumdns_private[*].id
    security_group_ids = [aws_security_group.enumdns_opensearch.id]
  }

  encrypt_at_rest {
    enabled = true
  }

  node_to_node_encryption {
    enabled = true
  }

  domain_endpoint_options {
    enforce_https = true
  }

  tags = {
    Name = "enumdns-opensearch"
    Environment = var.environment
  }
}

# Application Load Balancer
resource "aws_lb" "enumdns_alb" {
  name               = "enumdns-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.enumdns_alb.id]
  subnets            = aws_subnet.enumdns_public[*].id

  enable_deletion_protection = false

  tags = {
    Name = "enumdns-alb"
    Environment = var.environment
  }
}

# Security Groups
resource "aws_security_group" "enumdns_alb" {
  name_prefix = "enumdns-alb-"
  vpc_id      = aws_vpc.enumdns_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "enumdns-alb-sg"
    Environment = var.environment
  }
}

data "aws_availability_zones" "available" {}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}
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

#### Padrão Advanced Analysis
```go
type Technique interface {
    Name() string
    Generate(domain string, tlds []string) []Variation
    GetRiskLevel() string
    GetConfidence() float64
}

type VariationGenerator struct {
    BaseDomain string
    Options    GeneratorOptions
    analyzer   *RiskAnalyzer
}

func (vg *VariationGenerator) GenerateAll() []Variation {
    var allVariations []Variation
    
    for _, techniqueName := range vg.Options.Techniques {
        if technique, exists := AvailableTechniques[techniqueName]; exists {
            variations := technique.Generate(vg.BaseDomain, vg.Options.TargetTLDs)
            allVariations = append(allVariations, variations...)
        }
    }
    
    return vg.analyzeAndFilter(allVariations)
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

### Desenvolvimento de Novas Técnicas

#### Criando Nova Técnica de Análise

```go
// pkg/advanced/techniques.go

type CustomTechnique struct{}

func (t *CustomTechnique) Name() string {
    return "custom_technique"
}

func (t *CustomTechnique) GetRiskLevel() string {
    return "medium"
}

func (t *CustomTechnique) GetConfidence() float64 {
    return 0.7
}

func (t *CustomTechnique) Generate(domain string, tlds []string) []Variation {
    var variations []Variation
    baseName := getBaseName(domain)
    
    // Implementar lógica personalizada
    // Por exemplo: adicionar prefixos comuns
    prefixes := []string{"www", "mail", "ftp", "admin"}
    
    for _, prefix := range prefixes {
        for _, tld := range tlds {
            variation := prefix + baseName + "." + tld
            variations = append(variations, Variation{
                Domain:     variation,
                Technique:  t.Name(),
                Confidence: t.GetConfidence(),
                Risk:       t.GetRiskLevel(),
                BaseDomain: domain,
            })
        }
    }
    
    return variations
}

// Registrar a nova técnica
func init() {
    AvailableTechniques["custom_technique"] = &CustomTechnique{}
}
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

### Segurança na Análise de Ameaças

#### Validação de Entrada
```go
func validateDomain(domain string) error {
    if len(domain) > 255 {
        return errors.New("domain too long")
    }
    
    if matched, _ := regexp.MatchString(`^[a-zA-Z0-9.-]+# EnumDNS - Documentação Completa

## Índice
1. [Visão Geral](#visão-geral)
2. [Funcionalidades](#funcionalidades)
3. [Arquitetura e Estrutura](#arquitetura-e-estrutura)
4. [Instalação](#instalação)
5. [Uso e Exemplos](#uso-e-exemplos)
6. [Configuração](#configuração)
7. [Módulos e Componentes](#módulos-e-componentes)
8. [Formatos de Saída](#formatos-de-saída)
9. [Integração com Sistemas Externos](#integração-com-sistemas-externos)
10. [Desenvolvimento](#desenvolvimento)
11. [Considerações de Segurança](#considerações-de-segurança)

## Visão Geral

O **EnumDNS** é uma ferramenta modular de reconhecimento DNS desenvolvida em Go, projetada para profissionais de segurança cibernética realizarem enumeração DNS abrangente e análise de infraestrutura. A ferramenta oferece múltiplos métodos de descoberta de hosts e pode identificar automaticamente provedores de nuvem, produtos SaaS e datacenters.

### Principais Características
- **Modular**: Suporta diferentes tipos de enumeração (brute-force, reconhecimento, resolução, análise avançada)
- **Multi-plataforma**: Funciona em Linux, Windows e macOS
- **Flexível**: Múltiplos formatos de saída (texto, JSON, CSV, SQLite, Elasticsearch)
- **Escalável**: Suporta processamento paralelo com goroutines
- **Integrado**: Compatível com BloodHound e crt.sh
- **Proxy Support**: Suporta SOCKS4/SOCKS5
- **Análise de Ameaças**: Detecção de typosquatting, bitsquatting e ataques homográficos

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

#### 4. **Análise de Ameaças (threat-analysis)** ⭐ *ATUALIZADO*
- **Typosquatting**: Detecção de domínios com erros de digitação baseados em adjacência de teclado
- **Bitsquatting**: Identificação de domínios com alterações de bits
- **Ataques Homográficos**: Detecção de caracteres visualmente similares
- **Análise de Similaridade**: Cálculo de proximidade com domínios legítimos
- **Score de Ameaça**: Classificação automática de risco
- **Indicadores de Ameaça**: Identificação de padrões suspeitos (TLDs, Unicode tricks)

#### 5. **Relatórios**
- Conversão entre formatos (SQLite ↔ JSON Lines ↔ Texto)
- Sincronização com Elasticsearch

### Recursos Avançados

- **Detecção de Active Directory**: Identifica DCs e GCs automaticamente
- **Identificação de Cloud**: Reconhece AWS, Azure, GCP, CloudFlare, etc.
- **Reverse DNS**: Resolução PTR automática
- **ASN Detection**: Identifica ASN e informações de rede
- **Controle de Duplicatas**: Evita varreduras desnecessárias
- **Análise de Risco**: Sistema de pontuação para variações de domínio

## Arquitetura e Estrutura

### Estrutura de Diretórios

```
enumdns/
├── cmd/                    # Comandos CLI (Cobra)
│   ├── advanced.go        # Comando de análise avançada de ameaças
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
│   ├── advanced/          # Sistema de análise avançada de ameaças
│   │   ├── analyzer.go    # Analisador de risco e similaridade
│   │   ├── generator.go   # Gerador de variações
│   │   ├── options.go     # Opções e configurações
│   │   └── techniques.go  # Técnicas de geração de variações
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
- **Advanced**: Engine para análise de ameaças e variações de domínio
- Gerenciam workers (goroutines) e coordenam a execução

#### 2. **Advanced Threat Analysis System**
- **VariationGenerator**: Gera variações de domínio usando múltiplas técnicas
- **RiskAnalyzer**: Analisa similaridade e calcula scores de ameaça
- **TechniqueEngine**: Implementa algoritmos específicos de geração
- **ThreatClassifier**: Classifica e categoriza ameaças identificadas

#### 3. **Readers (Leitores de Entrada)**
- **FileReader**: Lê wordlists e listas de domínios
- **CrtShReader**: Integração com crt.sh
- Suportam diferentes formatos de entrada

#### 4. **Writers (Escritores de Saída)**
- **DbWriter**: SQLite, PostgreSQL, MySQL
- **JsonWriter**: JSON Lines
- **CsvWriter**: Formato CSV
- **ElasticWriter**: Elasticsearch
- **TextWriter**: Relatórios legíveis
- **StdoutWriter**: Saída no terminal

#### 5. **Models (Modelos de Dados)**
- **Result**: Resultado principal de DNS
- **ThreatResult**: Resultado de análise de ameaças (herda de Result)
- **FQDNData**: Dados de FQDN descobertos
- **ASN**: Informações de Sistema Autônomo
- **ASNIpDelegate**: Delegações de IP por ASN

## Instalação

### Linux

```bash
# Instalação automática
apt install curl jq

url=$(curl -s https://api.github.com/repos/bob-reis/enumdns/releases | jq -r '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match("linux-amd64.tar.gz$")) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

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

url=$(curl -s https://api.github.com/repos/bob-reis/enumdns/releases | jq -r --arg filename "darwin-${arch}.tar.gz\$" '[ .[] | {id: .id, tag_name: .tag_name, assets: [ .assets[] | select(.name|match($filename)) | {name: .name, browser_download_url: .browser_download_url} ]} | select(.assets != []) ] | sort_by(.id) | reverse | first(.[].assets[]) | .browser_download_url')

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
# Download latest bob-reis/enumdns release from github
function Invoke-Downloadenumdns {

    $repo = "bob-reis/enumdns"
    
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
git clone https://github.com/bob-reis/enumdns.git
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

### Comando Threat-Analysis (Análise de Ameaças) ⭐ *ATUALIZADO*

```bash
# Análise completa de ameaças para um domínio
enumdns threat-analysis -d example.com --all-techniques -o threats.txt

# Análise específica de typosquatting
enumdns threat-analysis -d example.com --typosquatting --write-db

# Análise de múltiplos domínios com configurações avançadas
enumdns threat-analysis -L domains.txt --bitsquatting --homographic \
  --max-variations 500 --target-tlds com,net,org,co,io --write-jsonl

# Análise de ataques homográficos
enumdns threat-analysis -d bank.com --homographic --write-elasticsearch-uri http://siem:9200/threats

# Análise completa com todas as técnicas
enumdns threat-analysis -d corporate.com --all-techniques \
  --max-variations 1000 --write-db --write-csv
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

### Sistema de Análise Avançada de Ameaças

#### Técnicas Implementadas

**1. Typosquatting (Erro de Digitação)**
- Baseado em adjacência de teclado QWERTY
- Simula erros comuns de digitação
- Score de risco: Alto (0.8)

```go
// Exemplo de mapeamento de teclas adjacentes
keyboardMap := map[rune][]rune{
    'q': {'w', 'a'}, 'w': {'q', 'e', 's'}, 
    'e': {'w', 'r', 'd'}, ...
}
```

**2. Bitsquatting (Alteração de Bits)**
- Altera bits individuais dos caracteres
- Simula erros de hardware/transmissão
- Score de risco: Médio (0.6)

**3. Ataques Homográficos**
- Usa caracteres visualmente similares
- Detecta Unicode tricks e variações acentuadas
- Score de risco: Alto (0.9)

```go
// Exemplo de caracteres homográficos
homographicMap := map[rune][]rune{
    'a': {'à', 'á', 'ä', 'â', 'ā', 'α'},
    'e': {'è', 'é', 'ê', 'ë', 'ē'},
    'o': {'ò', 'ó', 'ô', 'õ', 'ö', '0'},
}
```

**4. Técnicas Adicionais**
- **Insertion**: Inserção de caracteres comuns
- **Deletion**: Remoção de caracteres
- **Transposition**: Troca de caracteres adjacentes
- **TLD Variation**: Variações de domínios de topo
- **Subdomain Pattern**: Padrões comuns de phishing

#### Analisador de Risco

```go
type RiskAnalyzer struct {
    BaseDomain string
}

type AnalysisResult struct {
    Variation   Variation
    Similarity  float64      // 0-1 (Levenshtein similarity)
    ThreatScore float64      // 0-1 (risco calculado)
    Indicators  []string     // ["suspicious_tld", "phishing_pattern"]
}
```

#### Indicadores de Ameaça
- **suspicious_tld**: TLDs comumente usados para phishing
- **phishing_pattern**: Palavras-chave suspeitas (secure, login, verify)
- **high_similarity**: Similaridade > 80% com domínio original
- **unicode_tricks**: Uso de caracteres não-ASCII

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

-- Nova tabela para análise de ameaças
CREATE TABLE threat_results (
    id INTEGER PRIMARY KEY,
    fqdn TEXT,
    result_type TEXT,
    technique TEXT,
    confidence REAL,
    risk TEXT,
    base_domain TEXT,
    similarity REAL,
    probed_at DATETIME
);
```

### CSV
```csv
FQDN,RType,IPv4,IPv6,Target,CloudProduct,DC,GC,ProbedAt
example.com,A,192.0.2.1,,,,,false,false,2025-01-15T10:30:00Z
www.example.com,CNAME,,,example.com,,false,false,2025-01-15T10:30:01Z
```

## Integração com Sistemas Externos

### APIs e Conectores

#### API REST Wrapper (Python Flask)

```python
#!/usr/bin/env python3
# enumdns_api.py - Wrapper REST para EnumDNS

from flask import Flask, request, jsonify
import subprocess
import json
import tempfile
import os

app = Flask(__name__)

@app.route('/api/v1/threat-analysis', methods=['POST'])
def threat_analysis():
    """Endpoint para análise de ameaças"""
    data = request.get_json()
    domain = data.get('domain')
    techniques = data.get('techniques', ['all-techniques'])
    max_variations = data.get('max_variations', 1000)
    
    if not domain:
        return jsonify({'error': 'Domain required'}), 400
    
    try:
        # Criar arquivo temporário para resultados
        with tempfile.NamedTemporaryFile(mode='w', suffix='.jsonl', delete=False) as f:
            temp_file = f.name
        
        # Executar EnumDNS Advanced
        cmd = [
            'enumdns', 'advanced',
            '-d', domain,
            '--all-techniques' if 'all-techniques' in techniques else '--typosquatting',
            '--max-variations', str(max_variations),
            '--write-jsonl-file', temp_file,
            '-q'
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'error': result.stderr}), 500
        
        # Ler resultados
        threats = []
        with open(temp_file, 'r') as f:
            for line in f:
                if line.strip():
                    threat = json.loads(line)
                    threats.append(threat)
        
        # Cleanup
        os.unlink(temp_file)
        
        # Classificar por risco
        high_risk = [t for t in threats if t.get('cloud_product') or 'suspicious' in str(t)]
        
        return jsonify({
            'domain': domain,
            'threats': threats,
            'high_risk_count': len(high_risk),
            'total_variations': len(threats)
        })
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/v1/bulk-analysis', methods=['POST'])
def bulk_analysis():
    """Análise em lote de múltiplos domínios"""
    data = request.get_json()
    domains = data.get('domains', [])
    
    if not domains:
        return jsonify({'error': 'Domains list required'}), 400
    
    results = {}
    for domain in domains:
        # Implementar lógica similar ao endpoint anterior
        pass
    
    return jsonify(results)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
```

#### Integração com SIEM (Splunk)

```python
#!/usr/bin/env python3
# splunk_enumdns.py - Integração com Splunk

import splunklib.client as client
import json
import subprocess
import logging

class EnumDNSSplunkConnector:
    def __init__(self, splunk_host, splunk_port, username, password):
        self.service = client.connect(
            host=splunk_host,
            port=splunk_port,
            username=username,
            password=password
        )
        self.index = self.service.indexes['enumdns']
    
    def analyze_domain_threats(self, domain, techniques=['all-techniques']):
        """Executa análise de ameaças e envia para Splunk"""
        try:
            # Executar EnumDNS
            cmd = ['enumdns', 'advanced', '-d', domain, '--all-techniques', '--write-jsonl', '-q']
            result = subprocess.run(cmd, capture_output=True, text=True)
            
            if result.returncode != 0:
                logging.error(f"EnumDNS failed: {result.stderr}")
                return False
            
            # Processar resultados e enviar para Splunk
            for line in result.stdout.split('\n'):
                if line.strip():
                    try:
                        threat_data = json.loads(line)
                        # Enriquecer com metadados
                        threat_data['source'] = 'enumdns_advanced'
                        threat_data['analysis_time'] = time.time()
                        
                        # Enviar para Splunk
                        self.index.submit(json.dumps(threat_data))
                        
                    except json.JSONDecodeError:
                        continue
            
            return True
            
        except Exception as e:
            logging.error(f"Error in threat analysis: {e}")
            return False
    
    def create_alerts(self, domain, high_risk_threshold=0.8):
        """Cria alertas automáticos para ameaças de alto risco"""
        search_query = f'''
        search index=enumdns source=enumdns_advanced base_domain="{domain}"
        | where confidence > {high_risk_threshold}
        | eval risk_level=case(
            confidence > 0.9, "critical",
            confidence > 0.8, "high",
            confidence > 0.6, "medium",
            1=1, "low"
        )
        | table fqdn, technique, confidence, risk_level, similarity
        | sort -confidence
        '''
        
        return self.service.jobs.create(search_query)
```

#### Integração com Elasticsearch e Kibana

```python
#!/usr/bin/env python3
# elastic_enumdns.py - Integração avançada com Elasticsearch

from elasticsearch import Elasticsearch
import json
import subprocess
from datetime import datetime

class EnumDNSElasticConnector:
    def __init__(self, elastic_hosts, index_prefix='enumdns'):
        self.es = Elasticsearch(elastic_hosts)
        self.index_prefix = index_prefix
        self.setup_indices()
    
    def setup_indices(self):
        """Configura índices otimizados para EnumDNS"""
        
        # Índice principal para resultados DNS
        dns_mapping = {
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
                    "probed_at": {"type": "date"},
                    "asn": {"type": "long"},
                    "location": {"type": "geo_point"}
                }
            }
        }
        
        # Índice para análise de ameaças
        threat_mapping = {
            "mappings": {
                "properties": {
                    "fqdn": {"type": "keyword"},
                    "base_domain": {"type": "keyword"},
                    "technique": {"type": "keyword"},
                    "confidence": {"type": "float"},
                    "similarity": {"type": "float"},
                    "risk_level": {"type": "keyword"},
                    "threat_indicators": {"type": "keyword"},
                    "probed_at": {"type": "date"},
                    "tld": {"type": "keyword"},
                    "registrar": {"type": "keyword"},
                    "creation_date": {"type": "date"}
                }
            }
        }
        
        # Criar índices se não existirem
        dns_index = f"{self.index_prefix}-dns"
        threat_index = f"{self.index_prefix}-threats"
        
        if not self.es.indices.exists(index=dns_index):
            self.es.indices.create(index=dns_index, body=dns_mapping)
        
        if not self.es.indices.exists(index=threat_index):
            self.es.indices.create(index=threat_index, body=threat_mapping)
    
    def ingest_dns_results(self, domain):
        """Ingere resultados de reconhecimento DNS"""
        cmd = ['enumdns', 'recon', '-d', domain, '--write-jsonl', '-q']
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        actions = []
        for line in result.stdout.split('\n'):
            if line.strip():
                try:
                    dns_data = json.loads(line)
                    action = {
                        "_index": f"{self.index_prefix}-dns",
                        "_source": dns_data
                    }
                    actions.append(action)
                except json.JSONDecodeError:
                    continue
        
        if actions:
            from elasticsearch.helpers import bulk
            bulk(self.es, actions)
    
    def ingest_threat_analysis(self, domain):
        """Ingere resultados de análise de ameaças"""
        cmd = ['enumdns', 'advanced', '-d', domain, '--all-techniques', '--write-jsonl', '-q']
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        actions = []
        for line in result.stdout.split('\n'):
            if line.strip():
                try:
                    threat_data = json.loads(line)
                    # Enriquecer dados
                    threat_data['risk_level'] = self.calculate_risk_level(threat_data)
                    threat_data['tld'] = threat_data['fqdn'].split('.')[-1]
                    
                    action = {
                        "_index": f"{self.index_prefix}-threats",
                        "_source": threat_data
                    }
                    actions.append(action)
                except json.JSONDecodeError:
                    continue
        
        if actions:
            from elasticsearch.helpers import bulk
            bulk(self.es, actions)
    
    def calculate_risk_level(self, threat_data):
        """Calcula nível de risco baseado nos dados"""
        confidence = threat_data.get('confidence', 0)
        similarity = threat_data.get('similarity', 0)
        
        if confidence > 0.9 or similarity > 0.95:
            return 'critical'
        elif confidence > 0.7 or similarity > 0.8:
            return 'high'
        elif confidence > 0.5 or similarity > 0.6:
            return 'medium'
        else:
            return 'low'
    
    def create_kibana_dashboards(self):
        """Cria dashboards pré-configurados no Kibana"""
        
        # Dashboard de Overview
        overview_dashboard = {
            "title": "EnumDNS - Overview",
            "description": "Dashboard principal para análise de resultados EnumDNS",
            "panelsJSON": json.dumps([
                {
                    "title": "DNS Records by Type",
                    "type": "pie",
                    "query": {
                        "index": f"{self.index_prefix}-dns",
                        "aggregations": {
                            "types": {
                                "terms": {"field": "result_type"}
                            }
                        }
                    }
                },
                {
                    "title": "Threat Distribution",
                    "type": "histogram",
                    "query": {
                        "index": f"{self.index_prefix}-threats",
                        "aggregations": {
                            "risk_levels": {
                                "terms": {"field": "risk_level"}
                            }
                        }
                    }
                }
            ])
        }
        
        return overview_dashboard
```

#### Integração com APIs de Threat Intelligence

```python
#!/usr/bin/env python3
# threat_intel_integration.py

import requests
import json
from datetime import datetime

class ThreatIntelEnricher:
    def __init__(self, virustotal_api_key=None, urlvoid_api_key=None):
        self.vt_api_key = virustotal_api_key
        self.urlvoid_api_key = urlvoid_api_key
    
    def enrich_enumdns_results(self, enumdns_results):
        """Enriquece resultados do EnumDNS com dados de threat intelligence"""
        enriched_results = []
        
        for result in enumdns_results:
            fqdn = result.get('fqdn')
            if not fqdn:
                continue
            
            # Enriquecer com VirusTotal
            vt_data = self.check_virustotal(fqdn)
            if vt_data:
                result['virustotal'] = vt_data
            
            # Enriquecer com URLVoid
            urlvoid_data = self.check_urlvoid(fqdn)
            if urlvoid_data:
                result['urlvoid'] = urlvoid_data
            
            # Calcular score de reputação
            result['reputation_score'] = self.calculate_reputation(result)
            
            enriched_results.append(result)
        
        return enriched_results
    
    def check_virustotal(self, domain):
        """Consulta VirusTotal para informações do domínio"""
        if not self.vt_api_key:
            return None
        
        url = f"https://www.virustotal.com/vtapi/v2/domain/report"
        params = {
            'apikey': self.vt_api_key,
            'domain': domain
        }
        
        try:
            response = requests.get(url, params=params, timeout=10)
            if response.status_code == 200:
                data = response.json()
                return {
                    'detected_urls': data.get('detected_urls', []),
                    'detected_samples': data.get('detected_samples', []),
                    'reputation': data.get('reputation', 0)
                }
        except Exception as e:
            print(f"Error checking VirusTotal for {domain}: {e}")
        
        return None
    
    def check_urlvoid(self, domain):
        """Consulta URLVoid para verificação de reputação"""
        if not self.urlvoid_api_key:
            return None
        
        url = f"http://api.urlvoid.com/v1/pay-as-you-go/"
        params = {
            'key': self.urlvoid_api_key,
            'host': domain
        }
        
        try:
            response = requests.get(url, params=params, timeout=10)
            if response.status_code == 200:
                # Processar resposta XML do URLVoid
                # (implementação específica dependeria do formato)
                return {'status': 'checked'}
        except Exception as e:
            print(f"Error checking URLVoid for {domain}: {e}")
        
        return None
    
    def calculate_reputation(self, result):
        """Calcula score de reputação baseado em múltiplas fontes"""
        score = 50  # Score base
        
        # Ajustar baseado em VirusTotal
        vt_data = result.get('virustotal', {})
        if vt_data.get('detected_urls'):
            score -= 20
        if vt_data.get('reputation', 0) < 0:
            score -= 15
        
        # Ajustar baseado em análise de ameaças do EnumDNS
        if result.get('technique'):
            score -= 10
        if result.get('confidence', 0) > 0.8:
            score -= 20
        
        # Ajustar baseado em produtos cloud/SaaS
        if result.get('cloud_product'):
            score += 10
        if result.get('saas_product'):
            score += 5
        
        return max(0, min(100, score))
```

### Integração com Bancos de Dados

#### PostgreSQL Schema Otimizado

```sql
-- Schema otimizado para PostgreSQL
CREATE SCHEMA enumdns;

-- Tabela principal de resultados DNS
CREATE TABLE enumdns.dns_results (
    id BIGSERIAL PRIMARY KEY,
    test_id VARCHAR(50) NOT NULL,
    fqdn VARCHAR(255) NOT NULL,
    result_type VARCHAR(20) NOT NULL,
    ipv4 INET,
    ipv6 INET,
    target VARCHAR(255),
    ptr VARCHAR(255),
    txt TEXT,
    cloud_product VARCHAR(100),
    saas_product VARCHAR(100),
    datacenter VARCHAR(100),
    asn BIGINT,
    dc BOOLEAN DEFAULT FALSE,
    gc BOOLEAN DEFAULT FALSE,
    exists BOOLEAN DEFAULT TRUE,
    failed BOOLEAN DEFAULT FALSE,
    failed_reason TEXT,
    probed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para análise de ameaças
CREATE TABLE enumdns.threat_analysis (
    id BIGSERIAL PRIMARY KEY,
    base_domain VARCHAR(255) NOT NULL,
    variation_fqdn VARCHAR(255) NOT NULL,
    technique VARCHAR(50) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    similarity DECIMAL(3,2) NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    threat_indicators TEXT[],
    tld VARCHAR(10),
    registrar VARCHAR(100),
    creation_date DATE,
    reputation_score INTEGER,
    is_registered BOOLEAN,
    probed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para informações ASN
CREATE TABLE enumdns.asn_info (
    asn BIGINT PRIMARY KEY,
    rir_name VARCHAR(20) NOT NULL,
    country_code CHAR(2),
    organization TEXT,
    description TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para delegações IP-ASN
CREATE TABLE enumdns.asn_delegations (
    id BIGSERIAL PRIMARY KEY,
    rir_name VARCHAR(20) NOT NULL,
    country_code CHAR(2),
    subnet CIDR NOT NULL,
    addresses INTEGER NOT NULL,
    date_allocated DATE,
    asn BIGINT REFERENCES enumdns.asn_info(asn),
    status VARCHAR(20),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para performance
CREATE INDEX idx_dns_results_fqdn ON enumdns.dns_results(fqdn);
CREATE INDEX idx_dns_results_type ON enumdns.dns_results(result_type);
CREATE INDEX idx_dns_results_probed_at ON enumdns.dns_results(probed_at);
CREATE INDEX idx_dns_results_ipv4 ON enumdns.dns_results(ipv4) WHERE ipv4 IS NOT NULL;
CREATE INDEX idx_dns_results_cloud ON enumdns.dns_results(cloud_product) WHERE cloud_product IS NOT NULL;

CREATE INDEX idx_threat_base_domain ON enumdns.threat_analysis(base_domain);
CREATE INDEX idx_threat_technique ON enumdns.threat_analysis(technique);
CREATE INDEX idx_threat_risk_level ON enumdns.threat_analysis(risk_level);
CREATE INDEX idx_threat_confidence ON enumdns.threat_analysis(confidence);

-- Views úteis
CREATE VIEW enumdns.high_risk_threats AS
SELECT 
    base_domain,
    variation_fqdn,
    technique,
    confidence,
    similarity,
    risk_level,
    threat_indicators,
    probed_at
FROM enumdns.threat_analysis
WHERE risk_level IN ('critical', 'high')
ORDER BY confidence DESC, similarity DESC;

CREATE VIEW enumdns.domain_summary AS
SELECT 
    fqdn,
    COUNT(*) as total_records,
    COUNT(CASE WHEN result_type = 'A' THEN 1 END) as a_records,
    COUNT(CASE WHEN result_type = 'AAAA' THEN 1 END) as aaaa_records,
    COUNT(CASE WHEN result_type = 'CNAME' THEN 1 END) as cname_records,
    COUNT(CASE WHEN cloud_product IS NOT NULL THEN 1 END) as cloud_services,
    MAX(probed_at) as last_scan
FROM enumdns.dns_results
WHERE exists = TRUE
GROUP BY fqdn;

-- Funções úteis
CREATE OR REPLACE FUNCTION enumdns.get_threat_stats(domain_name VARCHAR)
RETURNS TABLE(
    technique VARCHAR,
    count BIGINT,
    avg_confidence DECIMAL,
    max_confidence DECIMAL
) AS $
BEGIN
    RETURN QUERY
    SELECT 
        t.technique,
        COUNT(*) as count,
        AVG(t.confidence) as avg_confidence,
        MAX(t.confidence) as max_confidence
    FROM enumdns.threat_analysis t
    WHERE t.base_domain = domain_name
    GROUP BY t.technique
    ORDER BY avg_confidence DESC;
END;
$ LANGUAGE plpgsql;

-- Trigger para atualização automática de timestamps
CREATE OR REPLACE FUNCTION enumdns.update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER update_asn_info_updated_at 
    BEFORE UPDATE ON enumdns.asn_info
    FOR EACH ROW EXECUTE FUNCTION enumdns.update_updated_at_column();
```

#### MongoDB Schema (NoSQL Alternative)

```javascript
// MongoDB Collections para EnumDNS

// Collection: dns_results
db.createCollection("dns_results", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["test_id", "fqdn", "result_type", "probed_at"],
         properties: {
            test_id: { bsonType: "string" },
            fqdn: { bsonType: "string" },
            result_type: { bsonType: "string" },
            ipv4: { bsonType: "string" },
            ipv6: { bsonType: "string" },
            target: { bsonType: "string" },
            cloud_product: { bsonType: "string" },
            saas_product: { bsonType: "string" },
            datacenter: { bsonType: "string" },
            asn: { bsonType: "long" },
            dc: { bsonType: "bool" },
            gc: { bsonType: "bool" },
            exists: { bsonType: "bool" },
            probed_at: { bsonType: "date" }
         }
      }
   }
});

// Collection: threat_analysis
db.createCollection("threat_analysis", {
   validator: {
      $jsonSchema: {
         bsonType: "object",
         required: ["base_domain", "variation_fqdn", "technique", "confidence"],
         properties: {
            base_domain: { bsonType: "string" },
            variation_fqdn: { bsonType: "string" },
            technique: { bsonType: "string" },
            confidence: { bsonType: "double", minimum: 0, maximum: 1 },
            similarity: { bsonType: "double", minimum: 0, maximum: 1 },
            risk_level: { enum: ["low", "medium", "high", "critical"] },
            threat_indicators: { bsonType: "array", items: { bsonType: "string" } },
            probed_at: { bsonType: "date" }
         }
      }
   }
});

// Índices para performance
db.dns_results.createIndex({ "fqdn": 1 });
db.dns_results.createIndex({ "result_type": 1 });
db.dns_results.createIndex({ "probed_at": -1 });
db.dns_results.createIndex({ "ipv4": 1 }, { sparse: true });
db.dns_results.createIndex({ "cloud_product": 1 }, { sparse: true });

db.threat_analysis.createIndex({ "base_domain": 1 });
db.threat_analysis.createIndex({ "technique": 1 });
db.threat_analysis.createIndex({ "risk_level": 1 });
db.threat_analysis.createIndex({ "confidence": -1 });

// Aggregation pipelines úteis
const getThreatStatsByDomain = (domain) => {
    return db.threat_analysis.aggregate([
        { $match: { base_domain: domain } },
        {
            $group: {
                _id: "$technique",
                count: { $sum: 1 },
                avg_confidence: { $avg: "$confidence" },
                max_confidence: { $max: "$confidence" },
                high_risk_count: {
                    $sum: { $cond: [{ $gte: ["$confidence", 0.8] }, 1, 0] }
                }
            }
        },
        { $sort: { avg_confidence: -1 } }
    ]);
};
```

### Integração com Orquestradores

#### Docker Compose para Ambiente Completo

```yaml
# docker-compose.yml - Ambiente completo EnumDNS
version: '3.8'

services:
  enumdns:
    build: .
    volumes:
      - ./data:/data
      - ./wordlists:/wordlists
      - ./config:/config
    environment:
      - ENUMDNS_DB_PATH=/data/enumdns.db
      - ENUMDNS_ELASTIC_URI=http://elasticsearch:9200/enumdns
    networks:
      - enumdns_net
    depends_on:
      - postgres
      - elasticsearch
  
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: enumdns
      POSTGRES_USER: enumdns
      POSTGRES_PASSWORD: securepassword
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - enumdns_net
  
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.15.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
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
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - enumdns_net
  
  enumdns-api:
    build: ./api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://enumdns:securepassword@postgres:5432/enumdns
      - ELASTIC_URL=http://elasticsearch:9200
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - elasticsearch
      - redis
    networks:
      - enumdns_net
  
  enumdns-worker:
    build: .
    command: ["python", "/app/worker.py"]
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - CELERY_BROKER_URL=redis://redis:6379
      - DATABASE_URL=postgresql://enumdns:securepassword@postgres:5432/enumdns
    depends_on:
      - redis
      - postgres
    networks:
      - enumdns_net
    deploy:
      replicas: 3

networks:
  enumdns_net:
    driver: bridge

volumes:
  postgres_data:
  es_data:
  redis_data:
```

#### Kubernetes Deployment

```yaml
# k8s/enumdns-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: enumdns-api
  namespace: security-tools
spec:
  replicas: 3
  selector:
    matchLabels:
      app: enumdns-api
  template:
    metadata:
      labels:
        app: enumdns-api
    spec:
      containers:
      - name: enumdns-api
        image: enumdns:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: enumdns-secrets
              key: database-url
        - name: ELASTIC_URL
          value: "http://elasticsearch:9200"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: enumdns-api-service
  namespace: security-tools
spec:
  selector:
    app: enumdns-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer

---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: enumdns-scheduled-scan
  namespace: security-tools
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: enumdns-scanner
            image: enumdns:latest
            command:
            - /bin/sh
            - -c
            - |
              enumdns threat-analysis -L /config/domains.txt --all-techniques \
                --write-db-uri "$DATABASE_URL" \
                --write-elasticsearch-uri "$ELASTIC_URL"
            env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: enumdns-secrets
                  key: database-url
            - name: ELASTIC_URL
              value: "http://elasticsearch:9200"
            volumeMounts:
            - name: config-volume
              mountPath: /config
          volumes:
          - name: config-volume
            configMap:
              name: enumdns-config
          restartPolicy: OnFailure
```

### Integração com Ferramentas de Automação

#### Ansible Playbook

```yaml
# ansible/enumdns-playbook.yml
---
- name: Deploy and configure EnumDNS
  hosts: security_servers
  become: yes
  vars:
    enumdns_version: "latest"
    enumdns_path: "/opt/enumdns"
    enumdns_user: "enumdns"
    enumdns_config_path: "/etc/enumdns"
    
  tasks:
    - name: Create enumdns user
      user:
        name: "{{ enumdns_user }}"
        system: yes
        home: "{{ enumdns_path }}"
        shell: /bin/bash
    
    - name: Create enumdns directories
      file:
        path: "{{ item }}"
        state: directory
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        mode: '0755'
      loop:
        - "{{ enumdns_path }}"
        - "{{ enumdns_config_path }}"
        - "{{ enumdns_path }}/logs"
        - "{{ enumdns_path }}/data"
        - "{{ enumdns_path }}/wordlists"
    
    - name: Download and install EnumDNS
      get_url:
        url: "https://github.com/bob-reis/enumdns/releases/download/{{ enumdns_version }}/enumdns-linux-amd64.tar.gz"
        dest: "/tmp/enumdns.tar.gz"
        mode: '0644'
    
    - name: Extract EnumDNS
      unarchive:
        src: "/tmp/enumdns.tar.gz"
        dest: "{{ enumdns_path }}"
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        remote_src: yes
    
    - name: Create symlink
      file:
        src: "{{ enumdns_path }}/enumdns"
        dest: "/usr/local/bin/enumdns"
        state: link
    
    - name: Install systemd service
      template:
        src: enumdns.service.j2
        dest: /etc/systemd/system/enumdns-api.service
        mode: '0644'
      notify: reload systemd
    
    - name: Configure EnumDNS
      template:
        src: enumdns.conf.j2
        dest: "{{ enumdns_config_path }}/enumdns.conf"
        owner: "{{ enumdns_user }}"
        group: "{{ enumdns_user }}"
        mode: '0600'
    
    - name: Start and enable EnumDNS service
      systemd:
        name: enumdns-api
        state: started
        enabled: yes
        daemon_reload: yes
  
  handlers:
    - name: reload systemd
      systemd:
        daemon_reload: yes
```

#### Terraform para Infraestrutura

```hcl
# terraform/main.tf - Infraestrutura EnumDNS na AWS
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC para EnumDNS
, domain); !matched {
        return errors.New("invalid domain characters")
    }
    
    return nil
}
```

#### Sanitização de Resultados
```go
func sanitizeVariation(variation string) string {
    // Remove caracteres perigosos
    variation = strings.ReplaceAll(variation, "<", "")
    variation = strings.ReplaceAll(variation, ">", "")
    variation = strings.ReplaceAll(variation, "\"", "")
    
    return variation
}
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

#### 4. Problemas com Análise Avançada
```bash
# Debug completo da análise de ameaças
enumdns threat-analysis -d example.com --all-techniques -D

# Verificar técnicas disponíveis
enumdns threat-analysis --help

# Reduzir número de variações para testes
enumdns threat-analysis -d example.com --typosquatting --max-variations 100
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

# 3. Análise de ameaças para detectar typosquatting
enumdns threat-analysis -d empresa.com.br --all-techniques --max-variations 2000 --write-db

# 4. Verificação de certificados SSL
enumdns resolve crtsh -d empresa.com.br --write-db --fqdn-out certificados_encontrados.txt

# 5. Geração de relatório final
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

# Análise de ameaças focada em phishing
enumdns threat-analysis -d target.com \
  --typosquatting --homographic \
  --target-tlds tk,ml,ga,cf,com,net \
  -X socks5://127.0.0.1:9050 \
  -t 1 \
  --write-db
```

### Cenário 3: Blue Team - Monitoramento Defensivo

```bash
# Análise completa de múltiplos domínios corporativos
enumdns recon -L dominios_empresa.txt \
  --write-db \
  --write-elasticsearch-uri http://siem.empresa.com:9200/dns_enum

# Detecção proativa de ameaças
enumdns threat-analysis -L dominios_criticos.txt \
  --all-techniques \
  --max-variations 5000 \
  --write-elasticsearch-uri http://siem.empresa.com:9200/threat_intel

# Integração com dados do AD
enumdns resolve bloodhound -L bloodhound_computers.json \
  --write-db \
  -s dc01.empresa.com
```

### Cenário 4: Threat Intelligence

```bash
# Análise de domínios suspeitos reportados
enumdns resolve file -L dominios_suspeitos.txt \
  --write-db \
  -o analise_threat_intel.txt

# Verificação de infraestrutura C2
enumdns threat-analysis -L dominios_c2.txt \
  --all-techniques \
  --write-jsonl \
  --write-csv

# Correlação com bases de threat intelligence
enumdns recon -L indicators_compromise.txt \
  --write-elasticsearch-uri http://ti-platform:9200/iocs
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

### TheHarvester
```bash
# Combinar com TheHarvester
theHarvester -d example.com -b all -f theharvester_results
enumdns resolve file -L theharvester_results.txt --write-db
```

## Scripts de Automação

### Script de Monitoramento Contínuo

```bash
#!/bin/bash
# monitor_threats.sh - Monitoramento contínuo de ameaças

DOMAIN="$1"
INTERVAL="$2"  # em horas

if [[ -z "$DOMAIN" || -z "$INTERVAL" ]]; then
    echo "Uso: $0 <dominio> <intervalo_horas>"
    exit 1
fi

while true; do
    echo "[$(date)] Iniciando análise de ameaças para $DOMAIN"
    
    # Reconhecimento completo
    enumdns recon -d "$DOMAIN" \
      --write-db \
      --local-workspace \
      -q
    
    # Análise de ameaças
    enumdns threat-analysis -d "$DOMAIN" \
      --all-techniques \
      --max-variations 2000 \
      --write-db \
      --local-workspace \
      -q
    
    # Verificar certificados
    enumdns resolve crtsh -d "$DOMAIN" \
      --write-db \
      --local-workspace \
      -q
    
    # Alertar sobre ameaças críticas
    CRITICAL_THREATS=$(sqlite3 enumdns_ctrl.db "
        SELECT COUNT(*) FROM threat_results 
        WHERE base_domain = '$DOMAIN' 
        AND risk_level = 'critical' 
        AND probed_at > datetime('now', '-${INTERVAL} hours')
    ")
    
    if [[ $CRITICAL_THREATS -gt 0 ]]; then
        echo "[ALERTA] $CRITICAL_THREATS ameaças críticas detectadas para $DOMAIN"
        # Enviar notificação (Slack, email, etc.)
        curl -X POST -H 'Content-type: application/json' \
             --data "{\"text\":\"🚨 $CRITICAL_THREATS ameaças críticas detectadas para $DOMAIN\"}" \
             "$SLACK_WEBHOOK_URL"
    fi
    
    echo "[$(date)] Análise concluída. Próxima em ${INTERVAL}h"
    sleep "${INTERVAL}h"
done
```

### Script de Análise Comparativa

```bash
#!/bin/bash
# compare_threat_analysis.sh - Comparação entre análises

OLD_DB="$1"
NEW_DB="$2"

if [[ -z "$OLD_DB" || -z "$NEW_DB" ]]; then
    echo "Uso: $0 <banco_antigo> <banco_novo>"
    exit 1
fi

# Extrair ameaças únicas
sqlite3 "$OLD_DB" "SELECT DISTINCT variation_fqdn FROM threat_results WHERE risk_level IN ('high', 'critical')" > old_threats.txt
sqlite3 "$NEW_DB" "SELECT DISTINCT variation_fqdn FROM threat_results WHERE risk_level IN ('high', 'critical')" > new_threats.txt

# Encontrar novas ameaças críticas
comm -13 <(sort old_threats.txt) <(sort new_threats.txt) > new_critical_threats.txt

# Encontrar ameaças que não são mais detectadas
comm -23 <(sort old_threats.txt) <(sort new_threats.txt) > resolved_threats.txt

echo "=== NOVAS AMEAÇAS CRÍTICAS ==="
echo "Total: $(wc -l < new_critical_threats.txt)"
if [[ -s new_critical_threats.txt ]]; then
    head -20 new_critical_threats.txt
fi

echo -e "\n=== AMEAÇAS RESOLVIDAS ==="
echo "Total: $(wc -l < resolved_threats.txt)"
if [[ -s resolved_threats.txt ]]; then
    head -10 resolved_threats.txt
fi

# Gerar relatório detalhado
sqlite3 "$NEW_DB" << EOF
.headers on
.mode column

SELECT 'RESUMO EXECUTIVO' as categoria, '' as valor;
SELECT 'Total de variações analisadas', COUNT(*) FROM threat_results;
SELECT 'Ameaças críticas', COUNT(*) FROM threat_results WHERE risk_level = 'critical';
SELECT 'Ameaças altas', COUNT(*) FROM threat_results WHERE risk_level = 'high';
SELECT 'Técnica mais efetiva', technique FROM threat_results WHERE risk_level IN ('critical', 'high') GROUP BY technique ORDER BY COUNT(*) DESC LIMIT 1;

SELECT '';
SELECT 'TOP 10 AMEAÇAS CRÍTICAS' as categoria, '' as valor;
SELECT variation_fqdn, technique, ROUND(confidence, 3) as conf, ROUND(similarity, 3) as sim 
FROM threat_results 
WHERE risk_level = 'critical' 
ORDER BY confidence DESC, similarity DESC 
LIMIT 10;
EOF

# Cleanup
rm -f old_threats.txt new_threats.txt
```

---

## Apêndices

### A. Wordlists Recomendadas

```bash
# Wordlists populares para DNS brute-force
/usr/share/wordlists/seclists/Discovery/DNS/subdomains-top1million-5000.txt
/usr/share/wordlists/seclists/Discovery/DNS/fierce-hostlist.txt
/usr/share/wordlists/seclists/Discovery/DNS/dns-Jhaddix.txt

# Wordlists especializadas para ameaças
/usr/share/wordlists/seclists/Discovery/DNS/tlds.txt
/usr/share/wordlists/seclists/Discovery/DNS/sortuniq-hosts.txt

# Wordlists personalizadas para análise de ameaças
custom_phishing_terms.txt
banking_keywords.txt
corporate_subdomains.txt
```

### B. Configuração de Elasticsearch para Ameaças

```json
PUT /enumdns-threats
{
  "settings": {
    "number_of_shards": 2,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "domain_analyzer": {
          "tokenizer": "keyword",
          "filters": ["lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "base_domain": {"type": "keyword", "analyzer": "domain_analyzer"},
      "variation_fqdn": {"type": "keyword", "analyzer": "domain_analyzer"},
      "technique": {"type": "keyword"},
      "confidence": {"type": "float"},
      "similarity": {"type": "float"},
      "risk_level": {"type": "keyword"},
      "threat_indicators": {"type": "keyword"},
      "tld": {"type": "keyword"},
      "is_registered": {"type": "boolean"},
      "registrar": {"type": "keyword"},
      "creation_date": {"type": "date"},
      "reputation_score": {"type": "integer"},
      "probed_at": {"type": "date"},
      "location": {
        "type": "geo_point"
      }
    }
  }
}
```

### C. Dashboard Kibana para Threat Analysis

```json
{
  "version": "8.15.0",
  "objects": [
    {
      "id": "enumdns-threat-overview",
      "type": "dashboard",
      "attributes": {
        "title": "EnumDNS - Threat Analysis Overview",
        "description": "Dashboard para análise de ameaças e variações de domínio",
        "panelsJSON": "[{\"version\":\"8.15.0\",\"gridData\":{\"x\":0,\"y\":0,\"w\":24,\"h\":15,\"i\":\"threat-summary\"},\"panelIndex\":\"threat-summary\",\"embeddableConfig\":{\"title\":\"Threat Distribution by Risk Level\"}}]",
        "timeRestore": true,
        "timeTo": "now",
        "timeFrom": "now-7d"
      }
    },
    {
      "id": "enumdns-technique-analysis",
      "type": "visualization",
      "attributes": {
        "title": "Technique Effectiveness Analysis",
        "description": "Análise da efetividade das técnicas de detecção",
        "visState": "{\"title\":\"Technique Effectiveness\",\"type\":\"histogram\",\"params\":{\"grid\":{\"categoryLines\":false,\"style\":{\"color\":\"#eee\"}},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Count\"}}],\"seriesParams\":[{\"show\":true,\"type\":\"histogram\",\"mode\":\"stacked\",\"data\":{\"label\":\"Count\",\"id\":\"1\"},\"valueAxis\":\"ValueAxis-1\",\"drawLinesBetweenPoints\":true,\"showCircles\":true}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"terms\",\"schema\":\"segment\",\"params\":{\"field\":\"technique\",\"size\":10,\"order\":\"desc\",\"orderBy\":\"1\"}}]}"
      }
    }
  ]
}
```

### D. Configuração systemd para Monitoramento

```ini
# /etc/systemd/system/enumdns-threat-monitor.service
[Unit]
Description=EnumDNS Threat Monitoring Service
After=network.target

[Service]
Type=simple
User=enumdns
Group=enumdns
WorkingDirectory=/opt/enumdns
ExecStart=/opt/enumdns/monitor_threats.sh corporate.com 6
Restart=always
RestartSec=300
Environment=SLACK_WEBHOOK_URL=https://hooks.slack.com/your/webhook/url

# Logs
StandardOutput=journal
StandardError=journal
SyslogIdentifier=enumdns-threat-monitor

# Security
NoNewPrivileges=yes
PrivateTmp=yes
---

## Análise de Ameaças (Threat-Analysis) - Guia Detalhado

### Visão Geral

O módulo `threat-analysis` (anteriormente `advanced`) é uma ferramenta avançada para detecção proativa de domínios maliciosos através de múltiplas técnicas de análise. Esta funcionalidade é essencial para equipes de Blue Team, Red Team e pesquisadores de segurança.

### Técnicas Implementadas

#### 1. Typosquatting
Detecta domínios com erros de digitação baseados na adjacência de teclas do teclado QWERTY.

**Exemplo:**
- Domínio original: `google.com`
- Variações: `goggle.com`, `foogle.com`, `googlw.com`

```bash
enumdns threat-analysis -d google.com --typosquatting --max-variations 500
```

#### 2. Bitsquatting
Identifica domínios criados através da alteração de um único bit em caracteres ASCII.

**Exemplo:**
- Domínio original: `facebook.com`
- Variações: `facebooks.com` (bit flip no 'k'), `facgbook.com`

```bash
enumdns threat-analysis -d facebook.com --bitsquatting --write-db
```

#### 3. Ataques Homográficos
Detecta o uso de caracteres Unicode visualmente similares a caracteres ASCII.

**Exemplos de substituições:**
- `a` → `а` (cirílico), `α` (grego)
- `e` → `е` (cirílico), `ε` (grego)
- `o` → `ο` (grego), `о` (cirílico)
- `1` → `l`, `I`, `|`

```bash
enumdns threat-analysis -d paypal.com --homographic --target-tlds com,net,org
```

### Análise de Risco

Cada variação recebe um **score de ameaça** (0.0 a 1.0) baseado em:

1. **Confidence da Técnica** (peso base)
2. **TLD Suspeito** (+0.2)
3. **Padrões de Phishing** (+0.3)
4. **Uso de Unicode** (+0.2)
5. **Alta Similaridade** (+0.2)

### Indicadores de Ameaça

O sistema identifica automaticamente:
- `suspicious_tld`: TLD frequentemente usado em ataques
- `phishing_pattern`: Palavras-chave comuns em phishing
- `unicode_tricks`: Uso de caracteres não-ASCII
- `high_similarity`: Similaridade > 80% com o domínio original

### Configurações Avançadas

#### Controle de Limites
```bash
# Análise básica (até 1.000 variações por domínio)
enumdns threat-analysis -d example.com --max-variations 1000

# Análise intensiva (até 10.000 variações)
enumdns threat-analysis -d example.com --max-variations 10000

# Análise focada (apenas 100 variações mais relevantes)
enumdns threat-analysis -d example.com --max-variations 100 --typosquatting
```

### Qualidade e Testes

- **Cobertura de Testes**: 98.4%
- **Testes Unitários**: 156 testes implementados
- **Testes de Performance**: Benchmarks incluídos
- **Validação**: Todos os algoritmos testados com casos edge

### Exemplos Práticos

#### Blue Team - Monitoramento Defensivo
```bash
# Monitoramento contínuo de marca
enumdns threat-analysis -d suaempresa.com \
  --all-techniques \
  --max-variations 5000 \
  --write-elasticsearch-uri http://siem:9200/threat_intel \
  --write-db

# Análise de múltiplas marcas
echo -e "marca1.com\nmarca2.com\nmarca3.com" > marcas.txt
enumdns threat-analysis -L marcas.txt \
  --all-techniques \
  --write-jsonl \
  --write-csv
```

#### Red Team - Reconnaissance
```bash
# Análise discreta via proxy
enumdns threat-analysis -d target.com \
  --typosquatting --homographic \
  -X socks5://127.0.0.1:9050 \
  --max-variations 500 \
  -o potential_targets.txt
```
