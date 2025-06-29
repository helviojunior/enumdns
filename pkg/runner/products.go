package runner

var products = map[string][]string{
    "CloudFront": []string{ "cloudfront.net", "cloudfront", },
    "CloudFlare": []string{ "cloudflare.com", "cloudflare", "cc-ecdn.net" },
    "Akamai": []string{ "akamaitechnologies.com", "edgekey.net", "akamaiedge.net", "akam.net" },
    "Imperva": []string{ "incapsula.com" },
    "Sucuri": []string{ "sucuri.net" },
    "Bunny": []string{ "bunnyinfra.net" },
    "KeyCDN": []string{ "proinity.net" },
    "CDN77": []string{ "cdn77.com" },
    "AWS Global Accelerator": []string{ "awsglobalaccelerator.com", },
    "AWS": []string{ "amazonaws.com", "awsdns", },
    "Microsoft Office 365": []string{ "lync.com", "office.com", "outlook.com" },
    "Microsoft Sharepoint": []string{ "sharepointonline.com" },
    "Azure": []string{ "azure-dns.com", "azure-dns.net", "azure-dns.org", "azure-dns.info", "azurewebsites.net", "cloudapp.net" },
    "Oracle Cloud": []string{ "oraclecloud.net" },
    "GCP": []string{ "googleusercontent.com" },
    "Registro.BR": []string { "dns.br" },
}

var saas_products = map[string][]string{
    "CloudFront": []string{ "cloudfront.net", "cloudfront", },
    "CloudFlare": []string{ "cloudflare.com", "cloudflare", "cc-ecdn.net" },
    "Akamai": []string{ "akamaitechnologies.com", "edgekey.net", "akamaiedge.net", "akam.net" },
    "Imperva": []string{ "incapsula.com" },
    "Sucuri": []string{ "sucuri.net" },
    "Bunny": []string{ "bunnyinfra.net" },
    "KeyCDN": []string{ "proinity.net" },
    "CDN77": []string{ "cdn77.com" },
    "AWS Global Accelerator": []string{ "awsglobalaccelerator.com", },
    "Microsoft Office 365": []string{ "lync.com", "office.com", "outlook.com" },
    "Microsoft Sharepoint": []string{ "sharepointonline.com" },
    "Azure": []string{ "azure-dns.com", "azure-dns.net", "azure-dns.org", "azure-dns.info", "azurewebsites.net" },
    "Heroku": []string{ "herokuapp.com", "herokudns.com" },
    "Registro.BR": []string { "dns.br" },
    "Trend Micro Email Security": []string { "tmes.trendmicro.com" },
    "Wix": []string{ "wixsite.com", "wixdns.net" },
    "Github": []string{ "github.io", "github.com" },
    "SalesForce": []string{ "exacttarget.com" },
    "Shopify": []string{ "myshopify.com" },
}

var datacenter = map[string][]string{
    "ALog": []string{ "alog.com.br", },
    "Toweb": []string{ "datacenter1.com.br", },
    "Uni5": []string{ "uni5.net", },
    "Hosting Service": []string{ "hostingservice.com", },
    "Locaweb": []string{ "locaweb.com.br" },
    "Equinix": []string{ "equinix.com" },
    "Telefonica": []string{ "tdatabrasil.net.br" },
    "UOL": []string{ "uoldiveo.com.br", "compasso.com.br" },
    "HostGator": []string{ "hostgator.com.br" },
    "Datacom": []string{ "dialhost.com.br", "stackpath.net" },
    "DialHost": []string{ "brascloud.com.br", "dialhost.com.br" },
}
