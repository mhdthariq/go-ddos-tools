package proxy

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProxyString tests the String method
func TestProxyString(t *testing.T) {
	tests := []struct {
		name  string
		proxy Proxy
		want  string
	}{
		{
			name:  "HTTP proxy",
			proxy: Proxy{Type: HTTP, Host: "192.168.1.1", Port: "8080"},
			want:  "192.168.1.1:8080",
		},
		{
			name:  "SOCKS5 proxy",
			proxy: Proxy{Type: SOCKS5, Host: "10.0.0.1", Port: "1080"},
			want:  "10.0.0.1:1080",
		},
		{
			name:  "Custom port",
			proxy: Proxy{Type: HTTP, Host: "proxy.example.com", Port: "3128"},
			want:  "proxy.example.com:3128",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.proxy.String()
			if got != tt.want {
				t.Errorf("Proxy.String() = %s; want %s", got, tt.want)
			}
		})
	}
}

// TestProxyURL tests the URL method
func TestProxyURL(t *testing.T) {
	tests := []struct {
		name  string
		proxy Proxy
		want  string
	}{
		{
			name:  "HTTP proxy URL",
			proxy: Proxy{Type: HTTP, Host: "192.168.1.1", Port: "8080"},
			want:  "http://192.168.1.1:8080",
		},
		{
			name:  "SOCKS4 proxy URL",
			proxy: Proxy{Type: SOCKS4, Host: "10.0.0.1", Port: "1080"},
			want:  "socks4://10.0.0.1:1080",
		},
		{
			name:  "SOCKS5 proxy URL",
			proxy: Proxy{Type: SOCKS5, Host: "10.0.0.2", Port: "1080"},
			want:  "socks5://10.0.0.2:1080",
		},
		{
			name:  "Unknown proxy type",
			proxy: Proxy{Type: ProxyType(99), Host: "proxy.example.com", Port: "3128"},
			want:  "proxy.example.com:3128",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.proxy.URL()
			if got != tt.want {
				t.Errorf("Proxy.URL() = %s; want %s", got, tt.want)
			}
		})
	}
}

// TestParseProxy tests proxy string parsing
func TestParseProxy(t *testing.T) {
	tests := []struct {
		name      string
		proxyStr  string
		proxyType int
		wantHost  string
		wantPort  string
		wantType  ProxyType
		wantErr   bool
	}{
		{
			name:      "Simple HTTP format",
			proxyStr:  "192.168.1.1:8080",
			proxyType: 1,
			wantHost:  "192.168.1.1",
			wantPort:  "8080",
			wantType:  HTTP,
			wantErr:   false,
		},
		{
			name:      "SOCKS5 host:port format",
			proxyStr:  "10.0.0.1:1080",
			proxyType: 5,
			wantHost:  "10.0.0.1",
			wantPort:  "1080",
			wantType:  SOCKS5,
			wantErr:   false,
		},
		{
			name:      "HTTP URL format",
			proxyStr:  "http://proxy.example.com:8080",
			proxyType: 0,
			wantHost:  "proxy.example.com",
			wantPort:  "8080",
			wantType:  HTTP,
			wantErr:   false,
		},
		{
			name:      "SOCKS5 URL format",
			proxyStr:  "socks5://10.0.0.1:1080",
			proxyType: 0,
			wantHost:  "10.0.0.1",
			wantPort:  "1080",
			wantType:  SOCKS5,
			wantErr:   false,
		},
		{
			name:      "HTTPS URL format",
			proxyStr:  "https://proxy.example.com:443",
			proxyType: 0,
			wantHost:  "proxy.example.com",
			wantPort:  "443",
			wantType:  HTTP,
			wantErr:   false,
		},
		{
			name:      "SOCKS4 URL format",
			proxyStr:  "socks4://10.0.0.2:1080",
			proxyType: 0,
			wantHost:  "10.0.0.2",
			wantPort:  "1080",
			wantType:  SOCKS4,
			wantErr:   false,
		},
		{
			name:      "URL without port (default to 80)",
			proxyStr:  "http://proxy.example.com",
			proxyType: 0,
			wantHost:  "proxy.example.com",
			wantPort:  "80",
			wantType:  HTTP,
			wantErr:   false,
		},
		{
			name:      "Invalid format - no colon",
			proxyStr:  "192.168.1.1",
			proxyType: 1,
			wantErr:   true,
		},
		{
			name:      "Invalid format - too many parts",
			proxyStr:  "192.168.1.1:8080:extra",
			proxyType: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProxy(tt.proxyStr, tt.proxyType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseProxy() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseProxy() unexpected error: %v", err)
				return
			}

			if got.Host != tt.wantHost {
				t.Errorf("ParseProxy().Host = %s; want %s", got.Host, tt.wantHost)
			}
			if got.Port != tt.wantPort {
				t.Errorf("ParseProxy().Port = %s; want %s", got.Port, tt.wantPort)
			}
			if got.Type != tt.wantType {
				t.Errorf("ParseProxy().Type = %v; want %v", got.Type, tt.wantType)
			}
		})
	}
}

// TestLoadProxies tests loading proxies from a file
func TestLoadProxies(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "proxies.txt")

	// Test case 1: Valid proxy file
	content := `192.168.1.1:8080
10.0.0.1:1080
# This is a comment
proxy.example.com:3128

invalid-proxy
192.168.1.2:8081
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	proxies, err := LoadProxies(testFile, 1)
	if err != nil {
		t.Errorf("LoadProxies() error = %v; want nil", err)
	}

	// Should have loaded 4 valid proxies (skipping comment, empty line, and invalid entry)
	// Actually loads: 192.168.1.1:8080, 10.0.0.1:1080, proxy.example.com:3128, 192.168.1.2:8081
	expectedCount := 4
	if len(proxies) != expectedCount {
		t.Errorf("LoadProxies() loaded %d proxies; want %d", len(proxies), expectedCount)
	}

	// Verify first proxy
	if len(proxies) > 0 {
		if proxies[0].Host != "192.168.1.1" {
			t.Errorf("First proxy host = %s; want 192.168.1.1", proxies[0].Host)
		}
		if proxies[0].Port != "8080" {
			t.Errorf("First proxy port = %s; want 8080", proxies[0].Port)
		}
	}

	// Test case 2: Non-existent file
	_, err = LoadProxies(filepath.Join(tmpDir, "nonexistent.txt"), 1)
	if err == nil {
		t.Error("LoadProxies() on non-existent file should return error")
	}

	// Test case 3: Empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	proxies, err = LoadProxies(emptyFile, 1)
	if err != nil {
		t.Errorf("LoadProxies() on empty file error = %v; want nil", err)
	}
	if len(proxies) != 0 {
		t.Errorf("LoadProxies() on empty file returned %d proxies; want 0", len(proxies))
	}

	// Test case 4: File with only comments
	commentsFile := filepath.Join(tmpDir, "comments.txt")
	if err := os.WriteFile(commentsFile, []byte("# Comment 1\n# Comment 2\n"), 0644); err != nil {
		t.Fatalf("Failed to create comments test file: %v", err)
	}

	proxies, err = LoadProxies(commentsFile, 1)
	if err != nil {
		t.Errorf("LoadProxies() on comments file error = %v; want nil", err)
	}
	if len(proxies) != 0 {
		t.Errorf("LoadProxies() on comments file returned %d proxies; want 0", len(proxies))
	}
}

// TestProxyTypes tests proxy type constants
func TestProxyTypes(t *testing.T) {
	tests := []struct {
		name string
		pt   ProxyType
		val  int
	}{
		{"ALL", ALL, 0},
		{"HTTP", HTTP, 1},
		{"SOCKS4", SOCKS4, 4},
		{"SOCKS5", SOCKS5, 5},
		{"RANDOM", RANDOM, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.pt) != tt.val {
				t.Errorf("ProxyType %s = %d; want %d", tt.name, int(tt.pt), tt.val)
			}
		})
	}
}

// BenchmarkParseProxy benchmarks proxy parsing
func BenchmarkParseProxy(b *testing.B) {
	for b.Loop() {
		ParseProxy("192.168.1.1:8080", 1)
	}
}

// BenchmarkProxyURL benchmarks URL generation
func BenchmarkProxyURL(b *testing.B) {
	p := Proxy{Type: HTTP, Host: "192.168.1.1", Port: "8080"}
	for b.Loop() {
		p.URL()
	}
}
