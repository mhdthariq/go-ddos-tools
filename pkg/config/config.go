package config

import (
	"encoding/json"
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
