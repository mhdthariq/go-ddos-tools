package tools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
			fmt.Println("DSTAT not implemented yet")
		case "CFIP":
			fmt.Println("CFIP not implemented yet")
		case "DNS":
			fmt.Println("DNS not implemented yet")
		case "CHECK":
			fmt.Println("CHECK not implemented yet")
		case "INFO":
			fmt.Println("INFO not implemented yet")
		case "TSSRV":
			fmt.Println("TSSRV not implemented yet")
		case "PING":
			fmt.Println("PING not implemented yet")
		default:
			fmt.Printf("%s command not found\n", command)
		}
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
