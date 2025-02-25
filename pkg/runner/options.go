package runner

// Options are global github.com/helviojunior/enumdnsenumdns options
type Options struct {
    // Logging is logging options
    Logging Logging
    // DNS Over HTTPs related options
    DnsOverHttps DnsOverHttps
    // Writer is output options
    Writer Writer
    // Scan is typically Scan options
    Scan Scan
    //
    DnsSufix string

    //
    DnsServer string
    DnsPort int
    DnsProtocol string
}

// Logging is log related options
type Logging struct {
    // Debug display debug level logging
    Debug bool
    // LogScanErrors log errors related to scanning
    LogScanErrors bool
    // Silence all logging
    Silence bool
}

// Writer options
type Writer struct {
    UserPath  string
    Db        bool
    DbURI     string
    DbDebug   bool // enables verbose database logs
    Csv       bool
    CsvFile   string
    Jsonl     bool
    JsonlFile string
    ELastic   bool
    ELasticURI string
    Text      bool
    TextFile  string
    Stdout    bool
    None      bool
}

// DNS Over HTTPs related options
type DnsOverHttps struct {

    // Don't write HTML response content
    SkipSSLCheck bool

    // Proxy server to use
    Proxy string

    ProxyUser string
    ProxyPassword string

    // UserAgent is the user-agent string to set for Chrome
    UserAgent string
    // Headers to add to every request
    Headers []string
    
}

// Scan is scanning related options
type Scan struct {
    // Threads (not really) are the number of goroutines to use.
    // More soecifically, its the go-rod page pool well use.
    Threads int
    // Timeout is the maximum time to wait for a page load before timing out.
    Timeout int
    // Number of seconds of delay between navigation and screenshotting
    Delay int
}

// NewDefaultOptions returns Options with some default values
func NewDefaultOptions() *Options {
    return &Options{
        Scan: Scan{
            Threads:          6,
            Timeout:          60,
        },
        Logging: Logging{
            Debug:         true,
            LogScanErrors: true,
        },
        DnsSufix: "",
    }
}