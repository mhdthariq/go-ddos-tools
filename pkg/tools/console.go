package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-ddos-tools/pkg/ui"
	"github.com/go-ddos-tools/pkg/utils"
)

// NetworkStats holds network statistics
type NetworkStats struct {
	BytesSent     int64
	BytesReceived int64
	PacketsSent   int64
	PacketsRecv   int64
}

// IPInfo holds IP information from API
type IPInfo struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
	Region  string `json:"region"`
	ISP     string `json:"isp"`
	Org     string `json:"org"`
}

var (
	lastNetStats  NetworkStats
	currentStats  NetworkStats
	statsRecorded bool
)

// RunConsole runs the interactive console
func RunConsole() {
	hostname, _ := os.Hostname()
	prompt := ui.Info(hostname + "@DDoSTools:~# ")

	fmt.Println()
	fmt.Println(ui.Header("  DDoSTools Interactive Tools Console"))
	fmt.Println(ui.Info("  ─────────────────────────────────────"))
	fmt.Println("  Type 'help' for available commands")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		cmd := strings.TrimSpace(scanner.Text())
		if cmd == "" {
			continue
		}

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "HELP":
			printHelp()
		case "CLEAR":
			fmt.Print("\033c")
		case "EXIT", "QUIT", "Q", "E", "LOGOUT", "CLOSE":
			ui.PrintInfo("Goodbye!")
			return
		case "DSTAT":
			runDstat()
		case "CFIP":
			runCFIP()
		case "DNS":
			runDNS(hostname, scanner)
		case "CHECK":
			runCheck(hostname, scanner)
		case "INFO":
			runInfo(hostname, scanner)
		case "TSSRV":
			runTSSRV(hostname, scanner)
		case "PING":
			runPing(hostname, scanner)
		default:
			ui.PrintError("Command '%s' not found. Type 'help' for available commands.", command)
		}
	}
}

// runDstat displays network and system statistics
func runDstat() {
	fmt.Println("Press 'q' to stop DSTAT and return to menu")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle keyboard input in separate goroutine
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Set stdin to non-blocking mode for character reading
				char, _, err := reader.ReadRune()
				if err != nil {
					continue
				}
				if char == 'q' || char == 'Q' {
					cancel()
					return
				}
			}
		}
	}()

	// Initialize stats
	if !statsRecorded {
		updateNetworkStats()
		statsRecorded = true
		time.Sleep(1 * time.Second)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nDSTAT stopped.")
			return
		case <-ticker.C:
			oldStats := currentStats
			updateNetworkStats()

			// Calculate deltas
			bytesSent := currentStats.BytesSent - oldStats.BytesSent
			bytesRecv := currentStats.BytesReceived - oldStats.BytesReceived
			packetsSent := currentStats.PacketsSent - oldStats.PacketsSent
			packetsRecv := currentStats.PacketsRecv - oldStats.PacketsRecv

			// Get memory stats
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			memUsedMB := m.Alloc / 1024 / 1024
			memTotalMB := m.Sys / 1024 / 1024
			memPercent := float64(m.Alloc) / float64(m.Sys) * 100

			fmt.Printf("\n--- Network & System Statistics ---\n")
			fmt.Printf("Bytes Sent:        %s/s\n", utils.HumanBytes(bytesSent))
			fmt.Printf("Bytes Received:    %s/s\n", utils.HumanBytes(bytesRecv))
			fmt.Printf("Packets Sent:      %s/s\n", utils.HumanFormat(packetsSent))
			fmt.Printf("Packets Received:  %s/s\n", utils.HumanFormat(packetsRecv))
			fmt.Printf("Memory Usage:      %d MB / %d MB (%.2f%%)\n", memUsedMB, memTotalMB, memPercent)
			fmt.Printf("Goroutines:        %d\n", runtime.NumGoroutine())
			fmt.Println("-----------------------------------")
		}
	}
}

// updateNetworkStats updates network statistics
// Platform-specific implementations are in console_linux.go, console_windows.go, and console_darwin.go

// runCheck checks if a website is online
func runCheck(hostname string, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s@DDoSTools:~# give-me-ipaddress# ", hostname)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		cmd := strings.ToUpper(domain)
		if cmd == "BACK" {
			return
		}
		if cmd == "CLEAR" {
			fmt.Print("\033c")
			continue
		}
		if cmd == "EXIT" || cmd == "QUIT" || cmd == "Q" || cmd == "E" || cmd == "LOGOUT" || cmd == "CLOSE" {
			os.Exit(0)
		}

		if !strings.Contains(domain, "/") {
			continue
		}

		fmt.Println("Please wait...")

		client := &http.Client{
			Timeout: 20 * time.Second,
		}

		resp, err := client.Get(domain)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		status := "ONLINE"
		if resp.StatusCode >= 500 {
			status = "OFFLINE"
		}

		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		fmt.Printf("Status: %s\n\n", status)
	}
}

// runInfo gets IP information
func runInfo(hostname string, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s@MHTools:~# give-me-ipaddress# ", hostname)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		cmd := strings.ToUpper(domain)
		if cmd == "BACK" {
			return
		}
		if cmd == "CLEAR" {
			fmt.Print("\033c")
			continue
		}
		if cmd == "EXIT" || cmd == "QUIT" || cmd == "Q" || cmd == "E" || cmd == "LOGOUT" || cmd == "CLOSE" {
			os.Exit(0)
		}

		// Clean domain
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.ReplaceAll(domain, "http://", "")
		if strings.Contains(domain, "/") {
			domain = strings.Split(domain, "/")[0]
		}

		fmt.Print("Please wait...\r")

		// Call IP info API
		url := fmt.Sprintf("https://ipwhois.app/json/%s", domain)
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error fetching info!")
			continue
		}
		defer resp.Body.Close()

		var info IPInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			fmt.Println("Error parsing response!")
			continue
		}

		if !info.Success {
			fmt.Println("Error: Failed to get info!")
			continue
		}

		fmt.Printf("Country: %s\n", info.Country)
		fmt.Printf("City: %s\n", info.City)
		fmt.Printf("Org: %s\n", info.Org)
		fmt.Printf("ISP: %s\n", info.ISP)
		fmt.Printf("Region: %s\n\n", info.Region)
	}
}

// runTSSRV resolves TeamSpeak SRV records
func runTSSRV(hostname string, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s@MHTools:~# give-me-domain# ", hostname)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		cmd := strings.ToUpper(domain)
		if cmd == "BACK" {
			return
		}
		if cmd == "CLEAR" {
			fmt.Print("\033c")
			continue
		}
		if cmd == "EXIT" || cmd == "QUIT" || cmd == "Q" || cmd == "E" || cmd == "LOGOUT" || cmd == "CLOSE" {
			os.Exit(0)
		}

		// Clean domain
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.ReplaceAll(domain, "http://", "")
		if strings.Contains(domain, "/") {
			domain = strings.Split(domain, "/")[0]
		}

		fmt.Print("Please wait...\r")

		// Lookup SRV records
		records := []string{"_ts3._udp.", "_tsdns._tcp."}

		for _, rec := range records {
			_, addrs, err := net.LookupSRV("", "", rec+domain)
			if err != nil {
				fmt.Printf("%s%s: Not found\n", rec, domain)
				continue
			}

			if len(addrs) > 0 {
				srv := addrs[0]
				target := strings.TrimSuffix(srv.Target, ".")
				fmt.Printf("%s%s: %s:%d\n", rec, domain, target, srv.Port)
			} else {
				fmt.Printf("%s%s: Not found\n", rec, domain)
			}
		}
		fmt.Println()
	}
}

// runPing pings a server
func runPing(hostname string, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s@MHTools:~# give-me-ipaddress# ", hostname)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		cmd := strings.ToUpper(domain)
		if cmd == "BACK" {
			return
		}
		if cmd == "CLEAR" {
			fmt.Print("\033c")
			continue
		}
		if cmd == "EXIT" || cmd == "QUIT" || cmd == "Q" || cmd == "E" || cmd == "LOGOUT" || cmd == "CLOSE" {
			os.Exit(0)
		}

		// Clean domain
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.ReplaceAll(domain, "http://", "")
		if strings.Contains(domain, "/") {
			domain = strings.Split(domain, "/")[0]
		}

		fmt.Println("Please wait...")

		// Resolve the domain
		ips, err := net.LookupIP(domain)
		if err != nil {
			fmt.Printf("Error resolving domain: %v\n", err)
			continue
		}

		if len(ips) == 0 {
			fmt.Println("No IP addresses found")
			continue
		}

		ip := ips[0].String()
		fmt.Printf("Address: %s\n", ip)

		// Perform TCP ping (simplified version)
		count := 5
		received := 0
		var totalRTT time.Duration

		for i := range count {
			start := time.Now()
			conn, err := net.DialTimeout("tcp", ip+":80", 2*time.Second)
			rtt := time.Since(start)

			if err == nil {
				conn.Close()
				received++
				totalRTT += rtt
				fmt.Printf("Reply from %s: time=%.2fms\n", ip, float64(rtt.Microseconds())/1000.0)
			} else {
				fmt.Printf("Request timeout\n")
			}

			if i < count-1 {
				time.Sleep(200 * time.Millisecond)
			}
		}

		avgRTT := float64(0)
		if received > 0 {
			avgRTT = float64(totalRTT.Microseconds()) / float64(received) / 1000.0
		}

		status := "ONLINE"
		if received == 0 {
			status = "OFFLINE"
		}

		fmt.Printf("\nPing Statistics:\n")
		fmt.Printf("Packets: Sent = %d, Received = %d, Lost = %d\n", count, received, count-received)
		fmt.Printf("Average RTT: %.2fms\n", avgRTT)
		fmt.Printf("Status: %s\n\n", status)
	}
}

// runCFIP fetches and displays CloudFlare IP ranges
func runCFIP() {
	fmt.Println("Fetching CloudFlare IP ranges...")
	fmt.Println()

	// CloudFlare publishes their IP ranges at these URLs
	ipv4URL := "https://www.cloudflare.com/ips-v4"
	ipv6URL := "https://www.cloudflare.com/ips-v6"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fetch IPv4 ranges
	fmt.Println(ui.Header("CloudFlare IPv4 Ranges:"))
	resp, err := client.Get(ipv4URL)
	if err != nil {
		ui.PrintError("Failed to fetch IPv4 ranges: %v", err)
	} else {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		count := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				fmt.Printf("  %s\n", line)
				count++
			}
		}
		fmt.Printf("\nTotal IPv4 ranges: %d\n", count)
	}

	fmt.Println()

	// Fetch IPv6 ranges
	fmt.Println(ui.Header("CloudFlare IPv6 Ranges:"))
	resp, err = client.Get(ipv6URL)
	if err != nil {
		ui.PrintError("Failed to fetch IPv6 ranges: %v", err)
	} else {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		count := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				fmt.Printf("  %s\n", line)
				count++
			}
		}
		fmt.Printf("\nTotal IPv6 ranges: %d\n", count)
	}

	fmt.Println()
	ui.PrintInfo("You can save these ranges to a file for firewall rules or analysis.")
	fmt.Println()
}

// runDNS performs advanced DNS lookups
func runDNS(hostname string, scanner *bufio.Scanner) {
	prompt := fmt.Sprintf("%s@DDoSTools:~# give-me-domain# ", hostname)

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			break
		}

		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		cmd := strings.ToUpper(domain)
		if cmd == "BACK" {
			return
		}
		if cmd == "CLEAR" {
			fmt.Print("\033c")
			continue
		}
		if cmd == "EXIT" || cmd == "QUIT" || cmd == "Q" || cmd == "E" || cmd == "LOGOUT" || cmd == "CLOSE" {
			os.Exit(0)
		}

		// Clean domain
		domain = strings.ReplaceAll(domain, "https://", "")
		domain = strings.ReplaceAll(domain, "http://", "")
		if strings.Contains(domain, "/") {
			domain = strings.Split(domain, "/")[0]
		}

		fmt.Println()
		fmt.Println(ui.Header(fmt.Sprintf("DNS Lookup Results for: %s", domain)))
		fmt.Println()

		// A Records (IPv4)
		fmt.Println(ui.Info("A Records (IPv4):"))
		ips, err := net.LookupIP(domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			foundIPv4 := false
			for _, ip := range ips {
				if ip.To4() != nil {
					fmt.Printf("  %s\n", ip.String())
					foundIPv4 = true
				}
			}
			if !foundIPv4 {
				fmt.Println("  No A records found")
			}
		}
		fmt.Println()

		// AAAA Records (IPv6)
		fmt.Println(ui.Info("AAAA Records (IPv6):"))
		foundIPv6 := false
		if err == nil {
			for _, ip := range ips {
				if ip.To4() == nil && ip.To16() != nil {
					fmt.Printf("  %s\n", ip.String())
					foundIPv6 = true
				}
			}
		}
		if !foundIPv6 {
			fmt.Println("  No AAAA records found")
		}
		fmt.Println()

		// CNAME Records
		fmt.Println(ui.Info("CNAME Record:"))
		cname, err := net.LookupCNAME(domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else if cname != domain && cname != domain+"." {
			fmt.Printf("  %s\n", cname)
		} else {
			fmt.Println("  No CNAME record found")
		}
		fmt.Println()

		// MX Records
		fmt.Println(ui.Info("MX Records (Mail Servers):"))
		mxRecords, err := net.LookupMX(domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else if len(mxRecords) > 0 {
			for _, mx := range mxRecords {
				fmt.Printf("  Priority: %d, Host: %s\n", mx.Pref, mx.Host)
			}
		} else {
			fmt.Println("  No MX records found")
		}
		fmt.Println()

		// NS Records
		fmt.Println(ui.Info("NS Records (Name Servers):"))
		nsRecords, err := net.LookupNS(domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else if len(nsRecords) > 0 {
			for _, ns := range nsRecords {
				fmt.Printf("  %s\n", ns.Host)
			}
		} else {
			fmt.Println("  No NS records found")
		}
		fmt.Println()

		// TXT Records
		fmt.Println(ui.Info("TXT Records:"))
		txtRecords, err := net.LookupTXT(domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else if len(txtRecords) > 0 {
			for _, txt := range txtRecords {
				// Truncate very long TXT records
				if len(txt) > 100 {
					fmt.Printf("  %s...\n", txt[:100])
				} else {
					fmt.Printf("  %s\n", txt)
				}
			}
		} else {
			fmt.Println("  No TXT records found")
		}
		fmt.Println()
	}
}

func printHelp() {
	fmt.Println()
	fmt.Println(ui.Header("Available Tools:"))
	fmt.Println()
	fmt.Printf("  %s	- Network and system statistics monitor\n", ui.Highlight("DSTAT"))
	fmt.Printf("  %s    - Check if a website is online\n", ui.Highlight("CHECK"))
	fmt.Printf("  %s	- Get IP address information\n", ui.Highlight("INFO"))
	fmt.Printf("  %s    - TeamSpeak SRV record lookup\n", ui.Highlight("TSSRV"))
	fmt.Printf("  %s	- Ping a server\n", ui.Highlight("PING"))
	fmt.Printf("  %s	- CloudFlare IP range finder\n", ui.Highlight("CFIP"))
	fmt.Printf("  %s	- Advanced DNS lookup tool\n", ui.Highlight("DNS"))
	fmt.Println()
	fmt.Println(ui.Header("Commands:"))
	fmt.Println()
	fmt.Printf("  %s	- Show this help message\n", ui.Highlight("HELP"))
	fmt.Printf("  %s	- Clear the screen\n", ui.Highlight("CLEAR"))
	fmt.Printf("  %s	- Exit the console\n", ui.Highlight("EXIT"))
	fmt.Println()
}

// StopAllAttacks stops all running attacks
func StopAllAttacks() {
	ui.PrintInfo("Stopping all attacks...")
	// In Go, we would signal all goroutines to stop
	// This is handled through context/channels in the main attack logic
	os.Exit(0)
}
