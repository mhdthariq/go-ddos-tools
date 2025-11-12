package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

// ProxyType represents different proxy types
type ProxyType int

const (
	HTTP   ProxyType = 1
	SOCKS4 ProxyType = 4
	SOCKS5 ProxyType = 5
	ALL    ProxyType = 0
	RANDOM ProxyType = 6
)

// Proxy represents a proxy server
type Proxy struct {
	Type ProxyType
	Host string
	Port string
}

// String returns the proxy as a string
func (p *Proxy) String() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}

// URL returns the proxy as a URL string
func (p *Proxy) URL() string {
	switch p.Type {
	case HTTP:
		return fmt.Sprintf("http://%s:%s", p.Host, p.Port)
	case SOCKS4:
		return fmt.Sprintf("socks4://%s:%s", p.Host, p.Port)
	case SOCKS5:
		return fmt.Sprintf("socks5://%s:%s", p.Host, p.Port)
	default:
		return fmt.Sprintf("%s:%s", p.Host, p.Port)
	}
}

// Dial creates a connection through the proxy
func (p *Proxy) Dial(network, address string) (net.Conn, error) {
	// For now, direct dial (proxy support can be added later with proper libraries)
	return net.DialTimeout(network, address, 10*time.Second)
}

// LoadProxies loads proxies from a file
func LoadProxies(filename string, proxyType int) ([]Proxy, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []Proxy
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		proxy, err := ParseProxy(line, proxyType)
		if err != nil {
			continue
		}
		proxies = append(proxies, *proxy)
	}

	return proxies, scanner.Err()
}

// ParseProxy parses a proxy string
func ParseProxy(proxyStr string, proxyType int) (*Proxy, error) {
	// Handle URL format
	if strings.Contains(proxyStr, "://") {
		u, err := url.Parse(proxyStr)
		if err != nil {
			return nil, err
		}

		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "80"
		}

		var pType ProxyType
		switch u.Scheme {
		case "http", "https":
			pType = HTTP
		case "socks4":
			pType = SOCKS4
		case "socks5":
			pType = SOCKS5
		default:
			pType = HTTP
		}

		return &Proxy{
			Type: pType,
			Host: host,
			Port: port,
		}, nil
	}

	// Handle host:port format
	parts := strings.Split(proxyStr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy format: %s", proxyStr)
	}

	var pType ProxyType
	switch proxyType {
	case 1:
		pType = HTTP
	case 4:
		pType = SOCKS4
	case 5:
		pType = SOCKS5
	default:
		pType = HTTP
	}

	return &Proxy{
		Type: pType,
		Host: parts[0],
		Port: parts[1],
	}, nil
}
