#!/usr/bin/python3
# -*- coding: utf-8 -*-

'''
Author : Helvio Junior (M4v3r1cK)

Sample file:

<registry>|<cc>|<type>|<start>|<value>|<date>|<status>
lacnic|DO|ipv4|5.183.80.0|1024|20240711|allocated
lacnic|BR|ipv4|24.152.0.0|1024|20200310|allocated
lacnic|BR|ipv4|24.152.4.0|1024|20200312|allocated
lacnic|BR|ipv4|24.152.8.0|1024|20200309|allocated
lacnic|BR|ipv4|24.152.12.0|1024|20200309|allocated

| Column       | Description                                                                       |
| ------------ | --------------------------------------------------------------------------------- |
| **registry** | RIR name (`arin`, `lacnic`, `ripe`, `apnic`, `afrinic`)                           |
| **cc**       | ISO 3166-1 country code (e.g., `BR`, `US`)                                        |
| **type**     | Resource type: `ipv4`, `ipv6`, or `asn`                                           |
| **start**    | Starting IP address, ASN number, or IPv6 prefix                                   |
| **value**    | For IPs: number of addresses<br>For ASN: count (usually 1)<br>IPv6: prefix length |
| **date**     | Date of allocation or assignment (YYYYMMDD)                                       |
| **status**   | `allocated`, `assigned`, `available`, etc.                                        |


https://github.com/sapics/ip-location-db/blob/main/asn/asn-ipv4.csv
ip_range_start, ip_range_end, autonomous_system_number, autonomous_system_organization


'''

import requests
import re, sys
import math
import ipaddress

URLS = [
    "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-latest",
    "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
    "https://ftp.apnic.net/stats/apnic/delegated-apnic-latest",
    "https://ftp.afrinic.net/stats/afrinic/delegated-afrinic-latest",
    "https://ftp.lacnic.net/pub/stats/ripencc/delegated-ripencc-latest",
]

URL_IPASN = "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/asn/asn-ipv4.csv"

GO_FILE = "pkg/models/asn.go"
RIR_NAMES = [
    "lacnic", "arin", "apnic", "afrinic", "ripencc"
]

def calculate_subnet_mask(num_hosts):
    """Calculates the subnet mask for a given number of hosts.

    Args:
        num_hosts: The number of hosts required in the subnet.

    Returns:
         The subnet mask in CIDR notation (e.g., /24).
    """
    if num_hosts <= 0:
        raise ValueError("Number of hosts must be positive.")

    # Calculate the required number of host bits
    host_bits = math.ceil(math.log2(num_hosts))

    # Calculate the subnet mask in CIDR notation
    subnet_mask_cidr = 32 - host_bits

    return subnet_mask_cidr

def fetch_file(url):
    print(f"[*] Downloading {url}")
    response = requests.get(url)
    response.raise_for_status()
    return response.text.splitlines()

def parse_ipasn(lines):
    print("[*] Parsing asn list...")
    asn1 = {}
    asn2 = {}
    for l in lines:
        lp = l.lower().split(",")
        if len(lp) == 4:
            ip_range_start = lp[0]
            asn = lp[2].strip()
            org = lp[3].strip()

            org = re.sub(r'[^a-z0-9 .()-]', '', org, flags=re.IGNORECASE)

            asn1[ip_range_start] = {
                "ip_range_start": ip_range_start,
                "autonomous_system_number": asn,
                "autonomous_system_organization": org,
            }
            if asn.strip() != "":
                if asn not in asn2:
                    asn2[asn] = {
                        "autonomous_system_number": asn,
                        "autonomous_system_organization": org,
                    }

    return asn1, asn2

def parse_services(lines, subnet_list, asn_list):
    print("[*] Parsing services...")
    services = []
    for line in lines:
        line = line.strip()
        if not line or line.startswith('#'):
            continue

        parts = line.lower().split("|")
        if len(parts) < 7:
            continue

        rir_name = parts[0]
        cc = parts[1].upper()
        rtype = parts[2]
        start = parts[3]
        value = parts[4]
        date = parts[5]
        status = parts[6]

        # Extract port and protocol
        if rir_name not in RIR_NAMES:
            continue

        asn = None
        subnet = ""
        if rtype in ["ipv4", "ipv6"]:
            asn = subnet_list.get(start, None)
            m = calculate_subnet_mask(int(value))
            subnet = f"{start}/{m}"
        elif rtype == "asn":
            asn = asn_list.get(start, None)

        services.append({
            "rir_name": rir_name,
            "cc": cc,
            "type": rtype,
            "start": start,
            "value": value,
            "date": date,
            "status": status,
            "subnet": subnet,
            "asn": int(asn["autonomous_system_number"]) if asn is not None else 0,
            "asn_org": asn["autonomous_system_organization"] if asn is not None else ""
        })

    return services

def generate_go_file(services):
    print("[*] Generating go file...")
    with open(GO_FILE, 'w') as f:
        f.write("package models\n\n")

        #f.write(f'import (\n')
        #f.write(f'    "math/big"\n')
        #f.write(f'    "net"\n')
        #f.write(f')\n\n')

        f.write("var AsnList = []ASN{\n")
        for svc in services:
            if svc["type"] == "asn" and svc["asn"] != 0:
                f.write("    {\n")
                f.write(f'        Number: {svc["asn"]},\n')
                f.write(f'        RIRName: "{svc["rir_name"]}",\n')
                f.write(f'        CountryCode: "{svc["cc"]}",\n')
                f.write(f'        Org: "{svc["asn_org"]}",\n')
                f.write("    },\n")
        f.write("}\n\n")

        f.write("var AsnDelagated = []ASNIpDelegate{\n")
        for svc in services:
            if svc["type"] in ["ipv4", "ipv6"]:
                f.write("    {\n")
                f.write(f'        RIRName: "{svc["rir_name"]}",\n')
                f.write(f'        CountryCode: "{svc["cc"]}",\n')
                f.write(f'        Subnet: "{svc["subnet"]}",\n')
                

                if svc["type"] == "ipv4":
                    int_ip = int(ipaddress.IPv4Address(svc["subnet"].split('/')[0]))
                    f.write(f'        IntIPv4: {int_ip},\n')
                elif svc["type"] == "ipv6":
                    ip = str(ipaddress.IPv6Address(svc["subnet"].split('/')[0]))
                    #f.write(f'        IntIPv6: *(new(big.Int).SetBytes(net.ParseIP("{ip}").To16())),\n')

                f.write(f'        Addresses: {svc["value"]},\n')
                f.write(f'        Date: "{svc["date"]}",\n')
                f.write(f'        ASN: {svc["asn"]},\n')
                f.write(f'        Status: "{svc["status"]}",\n')
                f.write("    },\n")
        f.write("}\n")

    print(f"[+] Generated {GO_FILE} with {len(services)} services.")

'''
type ASNIpDelegate struct {
    ID uint `json:"id" gorm:"primarykey"`

    Hash                  string    `gorm:"column:hash;index:,unique;" json:"hash"`
    RIRName               string    `gorm:"column:rir_name" json:"rir_name"`
    CountryCode           string    `gorm:"column:country_code" json:"country_code"`
    Type                  string    `gorm:"column:type" json:"type"` 
    Start                 string    `gorm:"column:start" json:"start"`
    Value                 string    `gorm:"column:value" json:"value"`
    Date                  string    `gorm:"column:date" json:"date"`
    Status                string    `gorm:"column:status" json:"status"`
}

'''

if __name__ == "__main__":
    services = []

    lines = fetch_file(URL_IPASN)
    subnet_list, asn_list = parse_ipasn(lines)
    if len(subnet_list) < 100:
        print("[!] Error: Fail to get IP/ASN csv")
        sys.exit(2)

    for url in URLS:
        lines = fetch_file(url)
        svc = parse_services(lines, subnet_list, asn_list)
        if len(svc) < 100:
            print("[!] Error: Fail to update tables")
        print(f"[*] Registers: {len(svc)}")
        services += svc
    generate_go_file(services)
