package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// PrintBanner prints the application banner
func PrintBanner() {
	banner := `
   ____    ____         _____           _
  |  _ \  |  _ \   ___ |_   _|__   ___ | |___
  | | | | | | | | / _ \  | |/ _ \ / _ \| / __|
  | |_| | | |_| || (_) | | | (_) | (_) | \__ \
  |____/  |____/  \___/  |_|\___/ \___/|_|___/
`
	if ColorsEnabled {
		fmt.Println(BrightCyan + banner + Reset)
		fmt.Println(Dim + "        Made with " + Red + "<3" + Dim + " by " + BrightBlue + "MHDDoS" + Reset)
		fmt.Println()
	} else {
		fmt.Println(banner)
		fmt.Println("        Made with <3 by MHDDoS")
		fmt.Println()
	}
}

// PrintUsage prints the usage instructions
func PrintUsage() {
	fmt.Println(Header("USAGE:"))
	fmt.Printf("  %s %s %s [%s]\n\n", os.Args[0], Color(BrightYellow, "<method>"), Color(BrightMagenta, "<url>"), Color(Dim, "flags"))

	fmt.Println(Header("EXAMPLES:"))
	fmt.Printf("  %s %s %s 5 1000 proxies.txt 100 60\n", os.Args[0], Color(BrightYellow, "CFB"), Color(BrightMagenta, "https://example.com"))
	fmt.Printf("  %s %s %s 100 60\n\n", os.Args[0], Color(BrightYellow, "TCP"), Color(BrightMagenta, "1.1.1.1:80"))

	fmt.Println(Header("LAYER 7:"))
	fmt.Printf("  %s <url> <socks_type> <threads> <proxies> <rpc> <duration>\n", Color(BrightCyan, "General"))
	fmt.Printf("  %s: 4=SOCKS4, 5=SOCKS5, 1=HTTP\n\n", Color(Dim, "socks_type"))

	fmt.Println(Header("LAYER 4:"))
	fmt.Printf("  %s <ip:port> <threads> <duration>\n", Color(BrightCyan, "General"))
	fmt.Println()
}

// PrintMethods prints the available attack methods
func PrintMethods(l7Methods, l4Methods, ampMethods []string) {
	sort.Strings(l7Methods)
	sort.Strings(l4Methods)
	sort.Strings(ampMethods)

	fmt.Println(Header("AVAILABLE METHODS"))
	fmt.Println()

	// Helper to print a customized box/table like structure
	printCategory := func(name string, methods []string, color string) {
		fmt.Printf(" %s %s\n", Color(color, "::"), Header(name))

		// Group in columns of 4
		for i := 0; i < len(methods); i += 4 {
			end := i + 4
			if end > len(methods) {
				end = len(methods)
			}

			row := methods[i:end]
			formattedRow := make([]string, len(row))
			for j, m := range row {
				formattedRow[j] = Color(Reset, m)
			}
			fmt.Printf("    %s\n", strings.Join(formattedRow, "\t"))
		}
		fmt.Println()
	}

	printCategory("LAYER 7 (HTTP/S)", l7Methods, BrightMagenta)
	printCategory("LAYER 4 (TCP/UDP)", l4Methods, BrightBlue)
	printCategory("AMPLIFICATION", ampMethods, BrightYellow)
}
