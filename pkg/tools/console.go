package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

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
	prompt := fmt.Sprintf("%s@DDoSTools:~# ", hostname)

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
			return
		case "DSTAT":
			runDstat()
		case "CFIP":
			fmt.Println("CFIP: CloudFlare IP finder - Coming soon")
		case "DNS":
			fmt.Println("DNS: DNS lookup tool - Coming soon")
		case "CHECK":
			runCheck(hostname, scanner)
		case "INFO":
			runInfo(hostname, scanner)
		case "TSSRV":
			runTSSRV(hostname, scanner)
		case "PING":
			runPing(hostname, scanner)
		default:
			fmt.Printf("%s command not found\n", command)
		}
	}
}

// runDstat displays network and system statistics
func runDstat() {
	fmt.Println("Press Ctrl+C to stop DSTAT")

	// Initialize stats
	if !statsRecorded {
		updateNetworkStats()
		statsRecorded = true
		time.Sleep(1 * time.Second)
	}

	for {
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
		fmt.Printf("Bytes Sent:      %s/s\n", utils.HumanBytes(bytesSent))
		fmt.Printf("Bytes Received:  %s/s\n", utils.HumanBytes(bytesRecv))
		fmt.Printf("Packets Sent:    %s/s\n", utils.HumanFormat(packetsSent))
		fmt.Printf("Packets Received: %s/s\n", utils.HumanFormat(packetsRecv))
		fmt.Printf("Memory Usage:    %d MB / %d MB (%.2f%%)\n", memUsedMB, memTotalMB, memPercent)
		fmt.Printf("Goroutines:      %d\n", runtime.NumGoroutine())
		fmt.Println("-----------------------------------")

		time.Sleep(1 * time.Second)
	}
}

// updateNetworkStats updates network statistics
// Platform-specific implementations are in console_linux.go, console_windows.go, and console_darwin.go

// runCheck checks if a website is online
func runCheck(hostname string, scanner *bufio.Scanner) {
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

func printHelp() {
	fmt.Println("Tools: DSTAT, CFIP, DNS, CHECK, INFO, TSSRV, PING")
	fmt.Println("Commands: HELP, CLEAR, EXIT")
}

// StopAllAttacks stops all running attacks
func StopAllAttacks() {
	log.Println("Stopping all attacks...")
	// In Go, we would signal all goroutines to stop
	// This is handled through context/channels in the main attack logic
	os.Exit(0)
}
