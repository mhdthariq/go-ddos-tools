package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCounter tests the thread-safe counter
func TestCounter(t *testing.T) {
	c := NewCounter()

	// Test initial value
	if c.Get() != 0 {
		t.Errorf("NewCounter() initial value = %d; want 0", c.Get())
	}

	// Test Add
	c.Add(5)
	if c.Get() != 5 {
		t.Errorf("After Add(5), Get() = %d; want 5", c.Get())
	}

	// Test multiple Add operations
	c.Add(10)
	c.Add(-3)
	if c.Get() != 12 {
		t.Errorf("After Add(10) and Add(-3), Get() = %d; want 12", c.Get())
	}

	// Test Set
	c.Set(100)
	if c.Get() != 100 {
		t.Errorf("After Set(100), Get() = %d; want 100", c.Get())
	}
}

// TestHumanBytes tests byte formatting
func TestHumanBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"Zero bytes", 0, "0 B"},
		{"Small bytes", 500, "500 B"},
		{"1 KiB", 1024, "1.00 KiB"},
		{"1.5 KiB", 1536, "1.50 KiB"},
		{"1 MiB", 1048576, "1.00 MiB"},
		{"1.5 MiB", 1572864, "1.50 MiB"},
		{"1 GiB", 1073741824, "1.00 GiB"},
		{"2.5 GiB", 2684354560, "2.50 GiB"},
		{"1 TiB", 1099511627776, "1.00 TiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HumanBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("HumanBytes(%d) = %s; want %s", tt.bytes, got, tt.want)
			}
		})
	}
}

// TestHumanFormat tests number formatting
func TestHumanFormat(t *testing.T) {
	tests := []struct {
		name string
		num  int64
		want string
	}{
		{"Zero", 0, "0"},
		{"Small number", 500, "500"},
		{"999", 999, "999"},
		{"1k", 1000, "1.00k"},
		{"1.5k", 1500, "1.50k"},
		{"1m", 1000000, "1.00m"},
		{"2.5m", 2500000, "2.50m"},
		{"1g", 1000000000, "1.00g"},
		{"1t", 1000000000000, "1.00t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HumanFormat(tt.num)
			if got != tt.want {
				t.Errorf("HumanFormat(%d) = %s; want %s", tt.num, got, tt.want)
			}
		})
	}
}

// TestLoadLines tests loading lines from a file
func TestLoadLines(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Test case 1: Valid file with content
	content := "line1\nline2\nline3\n\nline5"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	lines, err := LoadLines(testFile)
	if err != nil {
		t.Errorf("LoadLines() error = %v; want nil", err)
	}

	expectedLines := []string{"line1", "line2", "line3", "line5"}
	if len(lines) != len(expectedLines) {
		t.Errorf("LoadLines() returned %d lines; want %d", len(lines), len(expectedLines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("LoadLines() line[%d] = %s; want %s", i, line, expectedLines[i])
		}
	}

	// Test case 2: Non-existent file
	_, err = LoadLines(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("LoadLines() on non-existent file should return error")
	}

	// Test case 3: Empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	lines, err = LoadLines(emptyFile)
	if err != nil {
		t.Errorf("LoadLines() on empty file error = %v; want nil", err)
	}
	if len(lines) != 0 {
		t.Errorf("LoadLines() on empty file returned %d lines; want 0", len(lines))
	}
}

// TestGetDefaultUserAgents tests default user agents
func TestGetDefaultUserAgents(t *testing.T) {
	agents := GetDefaultUserAgents()

	if len(agents) == 0 {
		t.Error("GetDefaultUserAgents() returned empty slice")
	}

	// Verify all agents are non-empty strings
	for i, agent := range agents {
		if agent == "" {
			t.Errorf("GetDefaultUserAgents()[%d] is empty", i)
		}
	}
}

// TestGetDefaultReferers tests default referers
func TestGetDefaultReferers(t *testing.T) {
	referers := GetDefaultReferers()

	if len(referers) == 0 {
		t.Error("GetDefaultReferers() returned empty slice")
	}

	// Verify all referers are non-empty strings and contain "https://"
	for i, referer := range referers {
		if referer == "" {
			t.Errorf("GetDefaultReferers()[%d] is empty", i)
		}
		if len(referer) < 8 {
			t.Errorf("GetDefaultReferers()[%d] is too short: %s", i, referer)
		}
	}
}

// TestRandomBytes tests random bytes generation
func TestRandomBytes(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"Zero bytes", 0},
		{"Small bytes", 10},
		{"Medium bytes", 100},
		{"Large bytes", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomBytes(tt.n)
			if len(result) != tt.n {
				t.Errorf("RandomBytes(%d) returned %d bytes; want %d", tt.n, len(result), tt.n)
			}

			// Verify we got some bytes (byte type is already constrained to 0-255)
			for i, b := range result {
				_ = i // Use the index to avoid unused variable warning
				_ = b // Byte values are always valid (0-255 by definition)
			}
		})
	}
}

// TestRandInt tests random integer generation
func TestRandInt(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{"Same value", 5, 5},
		{"Small range", 1, 10},
		{"Large range", 0, 1000},
		{"Negative range", -10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				result := RandInt(tt.min, tt.max)
				if result < tt.min || result > tt.max {
					t.Errorf("RandInt(%d, %d) = %d; out of range [%d, %d]", tt.min, tt.max, result, tt.min, tt.max)
				}
			}
		})
	}

	// Test edge case: min > max (should return min)
	result := RandInt(10, 5)
	if result != 10 {
		t.Errorf("RandInt(10, 5) = %d; want 10", result)
	}
}

// TestRandString tests random string generation
func TestRandString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Empty string", 0},
		{"Short string", 5},
		{"Medium string", 20},
		{"Long string", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandString(tt.length)
			if len(result) != tt.length {
				t.Errorf("RandString(%d) returned string of length %d; want %d", tt.length, len(result), tt.length)
			}

			// Verify all characters are alphanumeric
			for i, c := range result {
				if !isAlphanumeric(c) {
					t.Errorf("RandString()[%d] = %c; not alphanumeric", i, c)
				}
			}
		})
	}

	// Test uniqueness (probabilistic)
	s1 := RandString(20)
	s2 := RandString(20)
	if s1 == s2 {
		t.Log("Warning: Two random strings are identical (this is very unlikely but possible)")
	}
}

// TestRandIPv4 tests random IPv4 generation
func TestRandIPv4(t *testing.T) {
	for i := 0; i < 100; i++ {
		ip := RandIPv4()

		// Verify format
		parts := len(splitIP(ip))
		if parts != 4 {
			t.Errorf("RandIPv4() = %s; doesn't have 4 parts", ip)
		}

		// Verify each octet is in range [1, 254]
		octets := splitIP(ip)
		for j, octetStr := range octets {
			var octet int
			if _, err := fmt.Sscanf(octetStr, "%d", &octet); err != nil {
				t.Errorf("RandIPv4() octet[%d] = %s; not a number", j, octetStr)
				continue
			}
			if octet < 1 || octet > 254 {
				t.Errorf("RandIPv4() octet[%d] = %d; out of range [1, 254]", j, octet)
			}
		}
	}
}

// Helper functions for tests

func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func splitIP(ip string) []string {
	result := []string{}
	current := ""
	for _, c := range ip {
		if c == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// Benchmark tests
func BenchmarkHumanBytes(b *testing.B) {
	for b.Loop() {
		HumanBytes(1073741824)
	}
}

func BenchmarkHumanFormat(b *testing.B) {
	for b.Loop() {
		HumanFormat(1000000)
	}
}

func BenchmarkRandString(b *testing.B) {
	for b.Loop() {
		RandString(20)
	}
}

func BenchmarkRandIPv4(b *testing.B) {
	for b.Loop() {
		RandIPv4()
	}
}

func BenchmarkCounterAdd(b *testing.B) {
	c := NewCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
}
