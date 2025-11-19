package proxy

import (
	"bufio"
	"context"
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

// CheckProxy checks if a proxy is working
func CheckProxy(proxy Proxy, testURL string, timeout time.Duration) bool {
	// For now, we'll do a basic connection test
	// In a full implementation, we'd actually test HTTP/SOCKS connectivity
	conn, err := net.DialTimeout("tcp", proxy.String(), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
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
		return nil, fmt.Errorf("socks Type Not Found [4, 5, 1, 0, 6]")
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
		return nil, fmt.Errorf("failed to download proxies: %v", err)
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
		return nil, fmt.Errorf("proxy Check failed, Your network may be the problem | The target may not be available")
	}

	// Save valid proxies to file
	if err := SaveProxies(validProxies, filename); err != nil {
		log.Printf("Warning: failed to save proxies to file: %v", err)
	}

	log.Printf("Proxy Count: %d", len(validProxies))
	return validProxies, nil
}
