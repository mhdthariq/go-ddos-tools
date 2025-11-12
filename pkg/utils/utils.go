package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
)

// Counter is a thread-safe counter
type Counter struct {
	value int64
}

// NewCounter creates a new counter
func NewCounter() *Counter {
	return &Counter{value: 0}
}

// Add increments the counter
func (c *Counter) Add(delta int64) {
	atomic.AddInt64(&c.value, delta)
}

// Get returns the current value
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set sets the counter value
func (c *Counter) Set(value int64) {
	atomic.StoreInt64(&c.value, value)
}

// HumanBytes converts bytes to human-readable format
func HumanBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// HumanFormat converts number to human-readable format
func HumanFormat(num int64) string {
	suffixes := []string{"", "k", "m", "g", "t", "p"}
	if num <= 999 {
		return fmt.Sprintf("%d", num)
	}

	exp := 0
	fnum := float64(num)
	for fnum >= 1000 && exp < len(suffixes)-1 {
		fnum /= 1000
		exp++
	}
	return fmt.Sprintf("%.2f%s", fnum, suffixes[exp])
}

// LoadLines loads lines from a file
func LoadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

// GetDefaultUserAgents returns default user agents
func GetDefaultUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101 Firefox/69.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
	}
}

// GetDefaultReferers returns default referers
func GetDefaultReferers() []string {
	return []string{
		"https://www.google.com/",
		"https://www.facebook.com/",
		"https://www.twitter.com/",
		"https://www.reddit.com/",
	}
}

// RandomBytes generates random bytes
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(RandInt(0, 255))
	}
	return b
}

// RandInt generates a random integer between min and max
func RandInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

// RandString generates a random string of length n
func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[RandInt(0, len(letters)-1)]
	}
	return string(b)
}

// RandIPv4 generates a random IPv4 address
func RandIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		RandInt(1, 254),
		RandInt(1, 254),
		RandInt(1, 254),
		RandInt(1, 254))
}
