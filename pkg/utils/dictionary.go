package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ddos-tools/pkg/config"
)

// GetDictionaryProvider finds a dictionary provider by name from config
func GetDictionaryProvider(cfg *config.Config, name string) *config.DictionaryProvider {
	if cfg == nil || cfg.DictionaryProviders == nil {
		return nil
	}

	for _, provider := range cfg.DictionaryProviders {
		if provider.Name == name {
			return &provider
		}
	}
	return nil
}

// DownloadDictionary downloads a dictionary file from a URL and saves it to the specified path
func DownloadDictionary(url, savePath string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file already exists
	if _, err := os.Stat(savePath); err == nil {
		fmt.Printf("Dictionary already exists at %s\n", savePath)
		return nil
	}

	fmt.Printf("Downloading dictionary from %s...\n", url)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download dictionary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download dictionary: HTTP %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to file
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save dictionary: %w", err)
	}

	fmt.Printf("Successfully downloaded dictionary (%d bytes) to %s\n", written, savePath)
	return nil
}

// LoadDictionary loads a dictionary file and returns its contents as a slice of strings
func LoadDictionary(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dictionary file: %w", err)
	}

	return lines, nil
}

// GenerateCredentials creates credential pairs from usernames and passwords
func GenerateCredentials(usernames, passwords []string) []string {
	var credentials []string

	// If no usernames provided, use common defaults
	if len(usernames) == 0 {
		usernames = []string{"admin", "user", "root", "administrator", "test"}
	}

	// If no passwords provided, use common defaults
	if len(passwords) == 0 {
		passwords = []string{"password", "123456", "admin", "12345678", "password123"}
	}

	// Generate all combinations of username:password
	for _, username := range usernames {
		for _, password := range passwords {
			credentials = append(credentials, fmt.Sprintf("%s:%s", username, password))
		}
	}

	return credentials
}

// LoadCredentialsFromFile loads credentials in username:password format from a file
func LoadCredentialsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file: %w", err)
	}
	defer file.Close()

	var credentials []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			credentials = append(credentials, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading credentials file: %w", err)
	}

	return credentials, nil
}

// EnsureDictionary checks if a dictionary exists, downloads it if not
func EnsureDictionary(name, url, savePath string) error {
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		fmt.Printf("Dictionary '%s' not found, downloading...\n", name)
		return DownloadDictionary(url, savePath)
	}
	return nil
}

// EnsureDictionaryFromConfig checks if a dictionary exists, downloads it from config if not
func EnsureDictionaryFromConfig(cfg *config.Config, dictName string) error {
	provider := GetDictionaryProvider(cfg, dictName)
	if provider == nil {
		return fmt.Errorf("dictionary provider '%s' not found in config", dictName)
	}

	if _, err := os.Stat(provider.Filename); os.IsNotExist(err) {
		fmt.Printf("Dictionary '%s' not found, downloading...\n", provider.Name)
		return DownloadDictionary(provider.URL, provider.Filename)
	}
	return nil
}

// DownloadDictionaryFromConfig downloads a dictionary using config provider
func DownloadDictionaryFromConfig(cfg *config.Config, dictName string) error {
	provider := GetDictionaryProvider(cfg, dictName)
	if provider == nil {
		return fmt.Errorf("dictionary provider '%s' not found in config", dictName)
	}

	return DownloadDictionary(provider.URL, provider.Filename)
}

// GetAllDictionaryProviders returns all dictionary providers from config
func GetAllDictionaryProviders(cfg *config.Config) []config.DictionaryProvider {
	if cfg == nil {
		return []config.DictionaryProvider{}
	}
	return cfg.DictionaryProviders
}
