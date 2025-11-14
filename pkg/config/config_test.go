package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig tests loading configuration from a JSON file
func TestLoadConfig(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test case 1: Valid config file with all fields
	validConfig := `{
	"MCBOT": "TestBot_",
	"MINECRAFT_DEFAULT_PROTOCOL": 340,
	"proxy-providers": [
		{
			"type": 1,
			"url": "http://example.com/proxies.txt",
			"timeout": 5
		},
		{
			"type": 5,
			"url": "http://example.com/socks5.txt",
			"timeout": 10
		}
	]
}`
	validFile := filepath.Join(tmpDir, "valid.json")
	if err := os.WriteFile(validFile, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to create valid config file: %v", err)
	}

	cfg, err := LoadConfig(validFile)
	if err != nil {
		t.Errorf("LoadConfig() on valid file error = %v; want nil", err)
	}

	if cfg.MCBot != "TestBot_" {
		t.Errorf("LoadConfig().MCBot = %s; want TestBot_", cfg.MCBot)
	}

	if cfg.MinecraftProtocol != 340 {
		t.Errorf("LoadConfig().MinecraftProtocol = %d; want 340", cfg.MinecraftProtocol)
	}

	if len(cfg.ProxyProviders) != 2 {
		t.Errorf("LoadConfig().ProxyProviders length = %d; want 2", len(cfg.ProxyProviders))
	}

	if len(cfg.ProxyProviders) > 0 {
		if cfg.ProxyProviders[0].Type != 1 {
			t.Errorf("First proxy provider type = %d; want 1", cfg.ProxyProviders[0].Type)
		}
		if cfg.ProxyProviders[0].URL != "http://example.com/proxies.txt" {
			t.Errorf("First proxy provider URL = %s; want http://example.com/proxies.txt", cfg.ProxyProviders[0].URL)
		}
		if cfg.ProxyProviders[0].Timeout != 5 {
			t.Errorf("First proxy provider timeout = %d; want 5", cfg.ProxyProviders[0].Timeout)
		}
	}

	// Test case 2: Config file with defaults (missing optional fields)
	defaultConfig := `{
	"proxy-providers": []
}`
	defaultFile := filepath.Join(tmpDir, "defaults.json")
	if err := os.WriteFile(defaultFile, []byte(defaultConfig), 0644); err != nil {
		t.Fatalf("Failed to create default config file: %v", err)
	}

	cfg, err = LoadConfig(defaultFile)
	if err != nil {
		t.Errorf("LoadConfig() on default file error = %v; want nil", err)
	}

	if cfg.MCBot != "MHDDoS_" {
		t.Errorf("LoadConfig().MCBot = %s; want MHDDoS_ (default)", cfg.MCBot)
	}

	if cfg.MinecraftProtocol != 47 {
		t.Errorf("LoadConfig().MinecraftProtocol = %d; want 47 (default)", cfg.MinecraftProtocol)
	}

	// Test case 3: Non-existent file
	_, err = LoadConfig(filepath.Join(tmpDir, "nonexistent.json"))
	if err == nil {
		t.Error("LoadConfig() on non-existent file should return error")
	}

	// Test case 4: Invalid JSON
	invalidJSON := `{
	"MCBOT": "TestBot_",
	"invalid": json
}`
	invalidFile := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	_, err = LoadConfig(invalidFile)
	if err == nil {
		t.Error("LoadConfig() on invalid JSON should return error")
	}

	// Test case 5: Empty JSON object
	emptyConfig := `{}`
	emptyFile := filepath.Join(tmpDir, "empty.json")
	if err := os.WriteFile(emptyFile, []byte(emptyConfig), 0644); err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}

	cfg, err = LoadConfig(emptyFile)
	if err != nil {
		t.Errorf("LoadConfig() on empty JSON error = %v; want nil", err)
	}

	// Should have default values
	if cfg.MCBot != "MHDDoS_" {
		t.Errorf("LoadConfig() empty config MCBot = %s; want MHDDoS_ (default)", cfg.MCBot)
	}

	if cfg.MinecraftProtocol != 47 {
		t.Errorf("LoadConfig() empty config MinecraftProtocol = %d; want 47 (default)", cfg.MinecraftProtocol)
	}

	// Test case 6: Partial config with only MCBot
	partialConfig := `{
	"MCBOT": "CustomBot_"
}`
	partialFile := filepath.Join(tmpDir, "partial.json")
	if err := os.WriteFile(partialFile, []byte(partialConfig), 0644); err != nil {
		t.Fatalf("Failed to create partial config file: %v", err)
	}

	cfg, err = LoadConfig(partialFile)
	if err != nil {
		t.Errorf("LoadConfig() on partial file error = %v; want nil", err)
	}

	if cfg.MCBot != "CustomBot_" {
		t.Errorf("LoadConfig().MCBot = %s; want CustomBot_", cfg.MCBot)
	}

	if cfg.MinecraftProtocol != 47 {
		t.Errorf("LoadConfig().MinecraftProtocol = %d; want 47 (default)", cfg.MinecraftProtocol)
	}
}

// TestConfigStruct tests the Config struct
func TestConfigStruct(t *testing.T) {
	cfg := Config{
		MCBot:             "TestBot_",
		MinecraftProtocol: 340,
		ProxyProviders: []ProxyProvider{
			{Type: 1, URL: "http://example.com", Timeout: 5},
		},
	}

	if cfg.MCBot != "TestBot_" {
		t.Errorf("Config.MCBot = %s; want TestBot_", cfg.MCBot)
	}

	if cfg.MinecraftProtocol != 340 {
		t.Errorf("Config.MinecraftProtocol = %d; want 340", cfg.MinecraftProtocol)
	}

	if len(cfg.ProxyProviders) != 1 {
		t.Errorf("Config.ProxyProviders length = %d; want 1", len(cfg.ProxyProviders))
	}
}

// TestProxyProvider tests the ProxyProvider struct
func TestProxyProvider(t *testing.T) {
	provider := ProxyProvider{
		Type:    1,
		URL:     "http://example.com/proxies.txt",
		Timeout: 10,
	}

	if provider.Type != 1 {
		t.Errorf("ProxyProvider.Type = %d; want 1", provider.Type)
	}

	if provider.URL != "http://example.com/proxies.txt" {
		t.Errorf("ProxyProvider.URL = %s; want http://example.com/proxies.txt", provider.URL)
	}

	if provider.Timeout != 10 {
		t.Errorf("ProxyProvider.Timeout = %d; want 10", provider.Timeout)
	}
}

// BenchmarkLoadConfig benchmarks config loading
func BenchmarkLoadConfig(b *testing.B) {
	// Create a temporary config file
	tmpDir := b.TempDir()
	configContent := `{
	"MCBOT": "TestBot_",
	"MINECRAFT_DEFAULT_PROTOCOL": 340,
	"proxy-providers": [
		{
			"type": 1,
			"url": "http://example.com/proxies.txt",
			"timeout": 5
		}
	]
}`
	configFile := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to create config file: %v", err)
	}

	b.ResetTimer()
	for b.Loop() {
		LoadConfig(configFile)
	}
}
