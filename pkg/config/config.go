package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

// Config holds the application configuration
type Config struct {
	MCBot             string          `json:"MCBOT"`
	MinecraftProtocol int             `json:"MINECRAFT_DEFAULT_PROTOCOL"`
	ProxyProviders    []ProxyProvider `json:"proxy-providers"`
	UserAgentFile     string          `json:"useragent_file"`
	RefererFile       string          `json:"referer_file"`
}

// ProxyProvider represents a proxy provider configuration
type ProxyProvider struct {
	Type    int    `json:"type"`
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Set defaults if not provided
	if config.MinecraftProtocol == 0 {
		config.MinecraftProtocol = 47
	}
	if config.MCBot == "" {
		config.MCBot = "MHDDoS_"
	}
	if config.UserAgentFile == "" {
		config.UserAgentFile = "files/useragent.txt"
	}
	if config.RefererFile == "" {
		config.RefererFile = "files/referers.txt"
	}

	return &config, nil
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate proxy providers
	for i, provider := range c.ProxyProviders {
		if provider.Timeout <= 0 {
			return fmt.Errorf("proxy provider %d: timeout must be positive, got %d", i, provider.Timeout)
		}
		if provider.URL == "" {
			return fmt.Errorf("proxy provider %d: URL cannot be empty", i)
		}
		if _, err := url.Parse(provider.URL); err != nil {
			return fmt.Errorf("proxy provider %d: invalid URL %q: %w", i, provider.URL, err)
		}
		// Validate proxy type (0=all, 1=HTTP, 4=SOCKS4, 5=SOCKS5, 6=RANDOM)
		validTypes := map[int]bool{0: true, 1: true, 4: true, 5: true, 6: true}
		if !validTypes[provider.Type] {
			return fmt.Errorf("proxy provider %d: invalid type %d, must be one of [0, 1, 4, 5, 6]", i, provider.Type)
		}
	}

	// Validate Minecraft protocol (if set)
	if c.MinecraftProtocol < 0 {
		return fmt.Errorf("minecraft protocol must be non-negative, got %d", c.MinecraftProtocol)
	}

	// Validate file paths are not empty if specified
	if c.UserAgentFile != "" {
		// Path validation - just ensure it's not obviously invalid
		if len(c.UserAgentFile) > 4096 {
			return fmt.Errorf("useragent file path too long")
		}
	}

	if c.RefererFile != "" {
		if len(c.RefererFile) > 4096 {
			return fmt.Errorf("referer file path too long")
		}
	}

	return nil
}
