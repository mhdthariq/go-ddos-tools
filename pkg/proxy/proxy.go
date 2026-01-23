package proxy

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-ddos-tools/pkg/config"
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
	switch p.Type {
	case HTTP:
		return p.dialHTTP(address)
	case SOCKS4:
		return p.dialSOCKS4(address)
	case SOCKS5:
		return p.dialSOCKS5(address)
	default:
		// Fallback to direct dial if proxy type is unknown
		return net.DialTimeout(network, address, 10*time.Second)
	}
}

// dialHTTP creates a connection through an HTTP proxy using CONNECT method
func (p *Proxy) dialHTTP(address string) (net.Conn, error) {
	// Connect to the proxy server
	conn, err := net.DialTimeout("tcp", p.String(), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to HTTP proxy: %w", err)
	}

	// Send CONNECT request
	connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", address, address)
	_, err = conn.Write([]byte(connectReq))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send CONNECT request: %w", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read CONNECT response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("HTTP proxy CONNECT failed with status: %d", resp.StatusCode)
	}

	return conn, nil
}

// dialSOCKS4 creates a connection through a SOCKS4 proxy
func (p *Proxy) dialSOCKS4(address string) (net.Conn, error) {
	// Connect to the proxy server
	conn, err := net.DialTimeout("tcp", p.String(), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SOCKS4 proxy: %w", err)
	}

	// Parse target address
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("invalid target address: %w", err)
	}

	// Resolve target IP
	ip := net.ParseIP(host)
	if ip == nil {
		// Resolve hostname to IP
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			conn.Close()
			return nil, fmt.Errorf("failed to resolve target host: %w", err)
		}
		// Use first IPv4 address
		for _, resolvedIP := range ips {
			if ipv4 := resolvedIP.To4(); ipv4 != nil {
				ip = ipv4
				break
			}
		}
		if ip == nil {
			conn.Close()
			return nil, fmt.Errorf("no IPv4 address found for host")
		}
	}
	ip = ip.To4()
	if ip == nil {
		conn.Close()
		return nil, fmt.Errorf("SOCKS4 only supports IPv4 addresses")
	}

	// Parse port
	var port uint16
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	// Build SOCKS4 connect request
	// VN (1) | CD (1) | DSTPORT (2) | DSTIP (4) | USERID (variable) | NULL (1)
	req := make([]byte, 9)
	req[0] = 0x04 // SOCKS version 4
	req[1] = 0x01 // CONNECT command
	binary.BigEndian.PutUint16(req[2:4], port)
	copy(req[4:8], ip)
	req[8] = 0x00 // NULL terminator for empty user ID

	_, err = conn.Write(req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send SOCKS4 request: %w", err)
	}

	// Read response
	// VN (1) | CD (1) | DSTPORT (2) | DSTIP (4)
	resp := make([]byte, 8)
	_, err = io.ReadFull(conn, resp)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read SOCKS4 response: %w", err)
	}

	// Check response code
	if resp[1] != 0x5A { // 0x5A = request granted
		conn.Close()
		return nil, fmt.Errorf("SOCKS4 request failed with code: %d", resp[1])
	}

	return conn, nil
}

// dialSOCKS5 creates a connection through a SOCKS5 proxy
func (p *Proxy) dialSOCKS5(address string) (net.Conn, error) {
	// Connect to the proxy server
	conn, err := net.DialTimeout("tcp", p.String(), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SOCKS5 proxy: %w", err)
	}

	// Step 1: Send greeting with supported authentication methods
	// VER (1) | NMETHODS (1) | METHODS (variable)
	greeting := []byte{0x05, 0x01, 0x00} // Version 5, 1 method, no auth
	_, err = conn.Write(greeting)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send SOCKS5 greeting: %w", err)
	}

	// Read server's chosen method
	// VER (1) | METHOD (1)
	authResp := make([]byte, 2)
	_, err = io.ReadFull(conn, authResp)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read SOCKS5 auth response: %w", err)
	}

	if authResp[0] != 0x05 {
		conn.Close()
		return nil, fmt.Errorf("invalid SOCKS5 version in response")
	}

	if authResp[1] == 0xFF {
		conn.Close()
		return nil, fmt.Errorf("SOCKS5 server requires authentication")
	}

	// Step 2: Send connect request
	// VER (1) | CMD (1) | RSV (1) | ATYP (1) | DST.ADDR (variable) | DST.PORT (2)
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("invalid target address: %w", err)
	}

	var port uint16
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	var req []byte
	ip := net.ParseIP(host)

	if ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			// IPv4 address
			req = make([]byte, 10)
			req[0] = 0x05 // Version 5
			req[1] = 0x01 // CONNECT command
			req[2] = 0x00 // Reserved
			req[3] = 0x01 // IPv4 address type
			copy(req[4:8], ipv4)
			binary.BigEndian.PutUint16(req[8:10], port)
		} else {
			// IPv6 address
			req = make([]byte, 22)
			req[0] = 0x05 // Version 5
			req[1] = 0x01 // CONNECT command
			req[2] = 0x00 // Reserved
			req[3] = 0x04 // IPv6 address type
			copy(req[4:20], ip.To16())
			binary.BigEndian.PutUint16(req[20:22], port)
		}
	} else {
		// Domain name
		if len(host) > 255 {
			conn.Close()
			return nil, fmt.Errorf("domain name too long")
		}
		req = make([]byte, 7+len(host))
		req[0] = 0x05            // Version 5
		req[1] = 0x01            // CONNECT command
		req[2] = 0x00            // Reserved
		req[3] = 0x03            // Domain name type
		req[4] = byte(len(host)) // Domain length
		copy(req[5:5+len(host)], host)
		binary.BigEndian.PutUint16(req[5+len(host):], port)
	}

	_, err = conn.Write(req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send SOCKS5 connect request: %w", err)
	}

	// Read response header
	// VER (1) | REP (1) | RSV (1) | ATYP (1)
	respHeader := make([]byte, 4)
	_, err = io.ReadFull(conn, respHeader)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read SOCKS5 response header: %w", err)
	}

	if respHeader[0] != 0x05 {
		conn.Close()
		return nil, fmt.Errorf("invalid SOCKS5 version in connect response")
	}

	if respHeader[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("SOCKS5 connect failed with code: %d", respHeader[1])
	}

	// Read the rest of the response based on address type
	var addrLen int
	switch respHeader[3] {
	case 0x01: // IPv4
		addrLen = 4 + 2 // 4 bytes IP + 2 bytes port
	case 0x03: // Domain
		// Read domain length first
		lenByte := make([]byte, 1)
		_, err = io.ReadFull(conn, lenByte)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to read domain length: %w", err)
		}
		addrLen = int(lenByte[0]) + 2 // domain + 2 bytes port
	case 0x04: // IPv6
		addrLen = 16 + 2 // 16 bytes IP + 2 bytes port
	default:
		conn.Close()
		return nil, fmt.Errorf("unknown address type in SOCKS5 response")
	}

	// Read and discard the bound address and port
	boundAddr := make([]byte, addrLen)
	_, err = io.ReadFull(conn, boundAddr)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read SOCKS5 bound address: %w", err)
	}

	return conn, nil
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

// DownloadFromProvider downloads proxies from a single provider
func DownloadFromProvider(provider config.ProxyProvider) ([]Proxy, error) {
	log.Printf("Downloading proxies from URL: %s, Type: %d, Timeout: %d",
		provider.URL, provider.Type, provider.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(provider.Timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", provider.URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Duration(provider.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download proxies: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse proxies from response
	return parseProxiesFromText(string(body), provider.Type)
}

// parseProxiesFromText parses proxies from text content
func parseProxiesFromText(text string, proxyType int) ([]Proxy, error) {
	var proxies []Proxy
	ipPortRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d+)`)

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := ipPortRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			proxy := &Proxy{
				Type: ProxyType(proxyType),
				Host: matches[1],
				Port: matches[2],
			}
			proxies = append(proxies, *proxy)
		}
	}

	return proxies, nil
}

// DownloadFromConfig downloads proxies from all providers in config
func DownloadFromConfig(cfg *config.Config, proxyType int) ([]Proxy, error) {
	var providers []config.ProxyProvider

	// Filter providers based on proxy type
	for _, provider := range cfg.ProxyProviders {
		if proxyType == 0 || provider.Type == proxyType {
			providers = append(providers, provider)
		}
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no proxy providers found for type %d", proxyType)
	}

	log.Printf("Downloading Proxies from %d Providers", len(providers))

	var allProxies []Proxy
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Download from all providers concurrently
	for _, provider := range providers {
		wg.Add(1)
		go func(p config.ProxyProvider) {
			defer wg.Done()

			proxies, err := DownloadFromProvider(p)
			if err != nil {
				log.Printf("Error downloading from provider %s: %v", p.URL, err)
				return
			}

			mu.Lock()
			allProxies = append(allProxies, proxies...)
			mu.Unlock()
		}(provider)
	}

	wg.Wait()

	// Remove duplicates
	uniqueProxies := make(map[string]Proxy)
	for _, proxy := range allProxies {
		key := proxy.String()
		uniqueProxies[key] = proxy
	}

	result := make([]Proxy, 0, len(uniqueProxies))
	for _, proxy := range uniqueProxies {
		result = append(result, proxy)
	}

	return result, nil
}

// CheckProxy checks if a proxy is working by attempting to dial through it
func CheckProxy(proxy Proxy, testURL string, timeout time.Duration) bool {
	// Parse the test URL to get host and port
	u, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		if u.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	// Set a deadline for the proxy connection test
	done := make(chan bool, 1)

	go func() {
		conn, err := proxy.Dial("tcp", host)
		if err != nil {
			done <- false
			return
		}
		conn.Close()
		done <- true
	}()

	select {
	case result := <-done:
		return result
	case <-time.After(timeout):
		return false
	}
}

// CheckAllProxies checks all proxies concurrently
func CheckAllProxies(proxies []Proxy, testURL string, timeout time.Duration, threads int) []Proxy {
	log.Printf("%d Proxies are getting checked, this may take awhile!", len(proxies))

	var validProxies []Proxy
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a higher thread count for proxy checking (at least 100, max 500)
	checkThreads := threads
	if checkThreads < 100 {
		checkThreads = 100
	}
	if checkThreads > 500 {
		checkThreads = 500
	}

	// Create a semaphore to limit concurrent checks
	sem := make(chan struct{}, checkThreads)

	for _, proxy := range proxies {
		wg.Add(1)
		go func(p Proxy) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			if CheckProxy(p, testURL, timeout) {
				mu.Lock()
				validProxies = append(validProxies, p)
				mu.Unlock()
			}
		}(proxy)
	}

	wg.Wait()
	return validProxies
}

// SaveProxies saves proxies to a file
func SaveProxies(proxies []Proxy, filename string) error {
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, proxy := range proxies {
		fmt.Fprintln(writer, proxy.String())
	}

	return writer.Flush()
}

// LoadOrDownloadProxies loads proxies from file or downloads them if file doesn't exist
func LoadOrDownloadProxies(filename string, proxyType int, cfg *config.Config, testURL string, threads int) ([]Proxy, error) {
	// Validate proxy type
	if proxyType != 0 && proxyType != 1 && proxyType != 4 && proxyType != 5 && proxyType != 6 {
		return nil, fmt.Errorf("invalid proxy type: must be one of [0, 1, 4, 5, 6]")
	}

	// Handle RANDOM type (6)
	if proxyType == 6 {
		randomTypes := []int{1, 4, 5}
		proxyType = randomTypes[time.Now().UnixNano()%int64(len(randomTypes))]
	}

	// Try to load from file first
	if _, err := os.Stat(filename); err == nil {
		proxies, err := LoadProxies(filename, proxyType)
		if err == nil && len(proxies) > 0 {
			log.Printf("Proxy Count: %d", len(proxies))
			return proxies, nil
		}
	}

	// File doesn't exist or is empty, download and check proxies
	log.Println("The file doesn't exist, creating files and downloading proxies.")

	// Download proxies from providers
	proxies, err := DownloadFromConfig(cfg, proxyType)
	if err != nil {
		return nil, fmt.Errorf("failed to download proxies: %w", err)
	}

	if len(proxies) == 0 {
		log.Println("Warning: No proxies downloaded")
		return nil, nil
	}

	// Check proxies
	if testURL == "" {
		testURL = "http://httpbin.org/get"
	}

	validProxies := CheckAllProxies(proxies, testURL, 5*time.Second, threads)

	if len(validProxies) == 0 {
		return nil, fmt.Errorf("proxy check failed: your network may be unreachable or the target may not be available")
	}

	// Save valid proxies to file
	if err := SaveProxies(validProxies, filename); err != nil {
		log.Printf("Warning: failed to save proxies to file: %v", err)
	}

	log.Printf("Proxy Count: %d", len(validProxies))
	return validProxies, nil
}
