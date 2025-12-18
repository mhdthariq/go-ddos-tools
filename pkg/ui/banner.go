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
	fmt.Printf("  %s %s %s [%s]\n\n", os.Args[0], Color(BrightYellow, "<method>"), Color(BrightMagenta, "<target>"), Color(Dim, "flags"))

	fmt.Println(Header("EXAMPLES:"))
	fmt.Printf("  %s %s %s -threads 1000 -duration 60 -proxy-file http.txt\n", os.Args[0], Color(BrightYellow, "CFB"), Color(BrightMagenta, "https://example.com"))
	fmt.Printf("  %s %s %s -threads 100 -duration 60\n\n", os.Args[0], Color(BrightYellow, "TCP"), Color(BrightMagenta, "1.1.1.1:80"))

	fmt.Println(Header("FLAGS:"))
	fmt.Println("  -threads <int>         Number of threads (default: 100)")
	fmt.Println("  -duration <int>        Duration in seconds (default: 60)")
	fmt.Println("  -rpc <int>             Requests per connection (L7) (default: 100)")
	fmt.Println("  -proxy-file <string>   File with proxies (default: proxies.txt)")
	fmt.Println("  -proxy-type <int>      SOCKS proxy type (4 or 5) (default: 5)")
	fmt.Println("  -config <string>       Configuration file (default: config.json)")
	fmt.Println("  -user-agents <string>  File with user agents (default: files/useragent.txt)")
	fmt.Println("  -referers <string>     File with referers (default: files/referers.txt)")
	fmt.Println("  -reflectors <string>   File with reflectors for amplification attacks")
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
			end := min(i+4, len(methods))

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
