package proxy

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-ddos-tools/pkg/config"
)

func TestLoadOrDownloadProxies_ExistingFile(t *testing.T) {
	// Create a temporary proxy file
	tmpDir := t.TempDir()
	proxyFile := filepath.Join(tmpDir, "test_proxies.txt")

	// Write some test proxies
	content := "192.168.1.1:8080\n192.168.1.2:8080\n"
	if err := os.WriteFile(proxyFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a minimal config
	cfg := &config.Config{
		ProxyProviders: []config.ProxyProvider{},
	}

	// Test loading from existing file
	proxies, err := LoadOrDownloadProxies(proxyFile, 1, cfg, "", 10)
	if err != nil {
		t.Fatalf("LoadOrDownloadProxies failed: %v", err)
	}

	if len(proxies) != 2 {
		t.Errorf("Expected 2 proxies, got %d", len(proxies))
	}
}

func TestLoadOrDownloadProxies_InvalidProxyType(t *testing.T) {
	tmpDir := t.TempDir()
	proxyFile := filepath.Join(tmpDir, "test_proxies.txt")

	cfg := &config.Config{}

	// Test with invalid proxy type
	_, err := LoadOrDownloadProxies(proxyFile, 99, cfg, "", 10)
	if err == nil {
		t.Error("Expected error for invalid proxy type, got nil")
	}
}

func TestLoadOrDownloadProxies_RandomType(t *testing.T) {
	tmpDir := t.TempDir()
	proxyFile := filepath.Join(tmpDir, "test_proxies.txt")

	// Write some test proxies
	content := "192.168.1.1:8080\n"
	if err := os.WriteFile(proxyFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{}

	// Test with RANDOM type (6) - should convert to 1, 4, or 5
	proxies, err := LoadOrDownloadProxies(proxyFile, 6, cfg, "", 10)
	if err != nil {
		t.Fatalf("LoadOrDownloadProxies failed: %v", err)
	}

	if len(proxies) != 1 {
		t.Errorf("Expected 1 proxy, got %d", len(proxies))
	}
}

func TestCheckProxy(t *testing.T) {
	// Test with a proxy that should fail
	proxy := Proxy{
		Type: HTTP,
		Host: "192.0.2.1", // TEST-NET-1 (documentation/testing)
		Port: "9999",
	}

	result := CheckProxy(proxy, "", 1*time.Second)
	// This should fail as it's a non-routable address
	if result {
		t.Log("Proxy check unexpectedly succeeded (network might be configured differently)")
	}
}

func TestSaveProxies(t *testing.T) {
	tmpDir := t.TempDir()
	proxyFile := filepath.Join(tmpDir, "subdir", "proxies.txt")

	proxies := []Proxy{
		{Type: HTTP, Host: "192.168.1.1", Port: "8080"},
		{Type: SOCKS5, Host: "192.168.1.2", Port: "1080"},
	}

	err := SaveProxies(proxies, proxyFile)
	if err != nil {
		t.Fatalf("SaveProxies failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(proxyFile); os.IsNotExist(err) {
		t.Error("Proxy file was not created")
	}

	// Verify content
	content, err := os.ReadFile(proxyFile)
	if err != nil {
		t.Fatalf("Failed to read proxy file: %v", err)
	}

	expected := "192.168.1.1:8080\n192.168.1.2:1080\n"
	if string(content) != expected {
		t.Errorf("Expected content:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestParseProxiesFromText(t *testing.T) {
	text := `
# Comment line
192.168.1.1:8080
192.168.1.2:1080

invalid line
10.0.0.1:3128
`

	proxies, err := parseProxiesFromText(text, 1)
	if err != nil {
		t.Fatalf("parseProxiesFromText failed: %v", err)
	}

	if len(proxies) != 3 {
		t.Errorf("Expected 3 proxies, got %d", len(proxies))
	}

	// Verify first proxy
	if proxies[0].Host != "192.168.1.1" || proxies[0].Port != "8080" {
		t.Errorf("First proxy incorrect: %v", proxies[0])
	}
}
