package internal

import (
	"context"
	"net/url"
	"net"
	"errors"
	"golang.org/x/net/proxy"
	"github.com/miekg/dns"

)

type SocksClient struct {
	Client *dns.Client
}

// Exchange performs a synchronous query. It sends the message m to the address
// contained in a and waits for a reply. Basic use pattern with a *dns.Client:
//
//	c := new(dns.Client)
//	in, rtt, err := c.Exchange(message, "127.0.0.1:53")
//
// Exchange does not retry a failed query, nor will it fall back to TCP in
// case of truncation.
// It is up to the caller to create a message that allows for larger responses to be
// returned. Specifically this means adding an EDNS0 OPT RR that will advertise a larger
// buffer, see SetEdns0. Messages without an OPT RR will fallback to the historic limit
// of 512 bytes
// To specify a local address or a timeout, the caller has to set the `Client.Dialer`
// attribute appropriately
func (c *SocksClient) Exchange(m *dns.Msg, proxyUri *url.URL, address string) (*dns.Msg, error) {
	if proxyUri == nil {
		return dns.Exchange(m, address); 
	}

	c.Client = new(dns.Client)
	co, err := c.Dial(proxyUri, address)
	if err != nil {
		return nil, err
	}
	defer co.Close()
	r, _, err := c.Client.ExchangeWithConn(m, co)
	return r, err
}

// Dial connects to the address on the named network.
func (c *SocksClient) Dial(proxyUri *url.URL, address string) (conn *dns.Conn, err error) {
	return c.DialContext(context.Background(), proxyUri, address)
}

// DialContext connects to the address on the named network, with a context.Context.
func (c *SocksClient) DialContext(ctx context.Context, proxyUri *url.URL, address string) (conn *dns.Conn, err error) {
	d, err := FromURL(proxyUri, proxy.Direct)
	if err != nil {
	    return nil, errors.New("Error connecting to proxy: " +  err.Error())
	}

	conn = new(dns.Conn)
	conn.Conn, err = d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	conn.UDPSize = c.Client.UDPSize
	return conn, nil
}

func FromURL(u *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {

	var auth *proxy.Auth
	if u.User != nil {
		auth = new(proxy.Auth)
		auth.User = u.User.Username()
		if p, ok := u.User.Password(); ok {
			auth.Password = p
		}
	}

	switch u.Scheme {
	case "socks4", "socks5", "socks5h":
		addr := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "1080"
		}
		return proxy.SOCKS5("tcp", net.JoinHostPort(addr, port), auth, forward)
	}

	return nil, errors.New("proxy: unknown scheme: " + u.Scheme)
}