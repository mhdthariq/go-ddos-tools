# DDoS Tools - Network Stress Testing Suite

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey)](https://github.com/go-ddos-tools)

A powerful network stress testing and analysis toolkit written in Go. This project is a modern reimplementation of MHDDoS with enhanced performance, cross-platform support, and additional security testing capabilities.

## ⚠️ Legal Disclaimer

**THIS SOFTWARE IS FOR EDUCATIONAL AND AUTHORIZED TESTING PURPOSES ONLY.**

- Only use this tool on networks and systems you own or have explicit written permission to test
- Unauthorized access to computer systems is illegal and punishable by law
- The authors assume no liability and are not responsible for any misuse or damage caused by this program
- By using this software, you agree to use it responsibly and legally

**📄 IMPORTANT**: Read [LICENSE](LICENSE) and [LEGAL.md](LEGAL.md) before using this software.

## 📋 Table of Contents

- [Features](#features)
- [Supported Attack Methods](#supported-attack-methods)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Requirements](#requirements)
- [License](#license)
- [Contributing](#contributing)

## ✨ Features

### Core Capabilities
- **Multi-Protocol Support**: Layer 4 (TCP/UDP) and Layer 7 (HTTP/HTTPS) attack methods
- **Proxy Integration**: Support for HTTP, SOCKS4, and SOCKS5 proxies
- **High Performance**: Written in Go for efficient concurrent operations
- **Cross-Platform**: Runs on Linux, Windows, and macOS
- **Real-time Monitoring**: Live statistics for PPS (Packets Per Second) and bandwidth usage
- **Interactive Console**: Built-in tools for network diagnostics and analysis

### Advanced Features
- **Amplification Attacks**: DNS, NTP, Memcached, and other reflection methods
- **Game Server Testing**: Specialized methods for Minecraft, FiveM, TeamSpeak, and Source Engine
- **Custom User Agents**: Randomized or custom user agent rotation
- **Flexible Configuration**: JSON-based configuration system
- **Thread Safety**: Atomic counters and safe concurrent operations

### Built-in Tools
- **DSTAT**: Real-time network and system statistics monitoring
- **CHECK**: Website availability checker
- **INFO**: IP geolocation and ISP information lookup
- **TSSRV**: TeamSpeak SRV record resolver
- **PING**: Advanced TCP ping utility
- **DNS**: DNS lookup and resolution tools (Coming Soon)
- **CFIP**: CloudFlare IP range finder (Coming Soon)

## 🎯 Supported Attack Methods

### Layer 7 (Application Layer) - 27 Methods

**✅ Implemented (8):**
```
GET, POST, HEAD, STRESS, SLOW, NULL, COOKIE, PPS
```

**🚧 Coming Soon (19):**
```
CFB, BYPASS, OVH, DYN, EVEN, GSB, DGB, AVB, CFBUAM, APACHE,
XMLRPC, BOT, BOMB, DOWNLOADER, KILLER, TOR, RHEX, STOMP
```

### Layer 4 (Transport Layer) - 14 Methods

**✅ Implemented (14):**
```
TCP, UDP, SYN, VSE, MINECRAFT, MCBOT, CONNECTION, CPS,
FIVEM, FIVEM-TOKEN, TS3, MCPE, ICMP, OVH-UDP
```

### Amplification Methods - 7 Methods

**✅ Implemented (7):**
```
MEM, NTP, DNS, ARD, CLDAP, CHAR, RDP
```

## 🚀 Installation

### Prerequisites
- Go 1.22 or higher
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools

# Build the project
go build -o ddos-tools main.go

# Or install directly
go install
```

### Quick Install (Linux/macOS)
```bash
curl -sSL https://raw.githubusercontent.com/go-ddos-tools/ddos-tools/main/install.sh | bash
```

## 🏃 Quick Start

### Basic Layer 7 Attack
```bash
# HTTP GET flood with proxies
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
#            └─┬─┘ └──────┬───────┘ ┬ └─┬┘ └────┬────┘ └─┬┘ └┬┘
#          Method    Target URL    │ Threads  Proxy   RPC Duration
#                              Proxy Type      File         (seconds)
```

### Basic Layer 4 Attack
```bash
# UDP flood
./ddos-tools UDP 192.168.1.1:80 100 60
#            └┬┘ └──────┬───────┘ └─┬┘ └┬┘
#          Method   Target:Port  Threads Duration
```

### Interactive Tools
```bash
# Launch tools console
./ddos-tools TOOLS

# Available commands in console:
# - DSTAT: Network statistics
# - CHECK: Website checker
# - INFO: IP information
# - PING: TCP ping utility
# - HELP: Show help
# - CLEAR: Clear screen
# - EXIT: Exit console
```

## 📁 Project Structure

```
ddos-tools/
├── main.go                      # Entry point
├── config.json                  # Configuration file
├── go.mod                       # Go module definition
├── README.md                    # This file
├── USAGE.md                     # Detailed usage guide
│
├── pkg/                         # Core packages
│   ├── attacks/                 # Attack implementations
│   │   ├── layer4.go           # Layer 4 attack methods
│   │   └── layer7.go           # Layer 7 attack methods
│   │
│   ├── config/                  # Configuration management
│   │   ├── config.go           # Config loader
│   │   └── config_test.go      # Tests
│   │
│   ├── methods/                 # Method definitions
│   │   ├── methods.go          # Method validation
│   │   └── methods_test.go     # Tests
│   │
│   ├── minecraft/               # Minecraft protocol
│   │   └── minecraft.go        # Packet builders
│   │
│   ├── proxy/                   # Proxy support
│   │   ├── proxy.go            # Proxy handler
│   │   └── proxy_test.go       # Tests
│   │
│   ├── tools/                   # Interactive tools
│   │   ├── console.go          # Console interface
│   │   ├── console_linux.go    # Linux-specific stats
│   │   ├── console_darwin.go   # macOS-specific stats
│   │   └── console_windows.go  # Windows-specific stats
│   │
│   └── utils/                   # Utilities
│       ├── utils.go            # Helper functions
│       └── utils_test.go       # Tests
│
└── files/                       # Resource files
    ├── proxies/                 # Proxy lists
    ├── useragent.txt           # User agent list
    └── referers.txt            # Referer list
```

## 📚 Documentation

For detailed usage instructions, examples, and advanced configurations, see:
- **[USAGE.md](USAGE.md)** - Comprehensive usage guide
- **[Configuration Guide](docs/CONFIGURATION.md)** - Config file setup (Coming Soon)
- **[API Documentation](docs/API.md)** - Developer API reference (Coming Soon)

## 🔧 Requirements

### Runtime Requirements
- Go 1.22+ (for range over int support)
- Network connectivity
- Sufficient system resources (RAM, CPU) for concurrent operations

### Optional Requirements
- Proxy lists (for proxy-based attacks)
- Root/Administrator privileges (for raw socket operations like SYN flood)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## 🛠️ Development

### Code Style
This project follows standard Go conventions:
- `gofmt` for code formatting
- Go 1.22+ modern idioms (range over int, slices.Contains, etc.)
- Comprehensive test coverage
- Platform-specific code using build tags

### Building for Different Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ddos-tools-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o ddos-tools.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o ddos-tools-mac
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines
- Write tests for new features
- Follow Go best practices
- Update documentation
- Ensure all tests pass

## 📊 Performance

- **Concurrent Operations**: Supports thousands of concurrent threads
- **Memory Efficient**: Optimized for low memory footprint
- **Cross-Platform**: Native performance on all supported platforms
- **Real-time Stats**: Minimal overhead monitoring

## 🔒 Security Notes

- Some methods require elevated privileges (root/admin)
- Raw socket operations may be restricted by OS security policies
- Use responsibly and only on authorized systems
- Always obtain proper authorization before testing

## 📝 Version History

- **v2.4 SNAPSHOT (Go)** - Current development version
  - Complete Go rewrite
  - Enhanced performance
  - Cross-platform support
  - Modern Go idioms (Go 1.22+)
  - **Layer 7**: 8/27 methods implemented (30%)
  - **Layer 4**: 14/14 methods implemented (100%)
  - **Amplification**: 7/7 methods implemented (100%)
  - **Tools**: 5/7 tools implemented (71%)

## 📄 License

This project is licensed under the **MIT License** with additional terms for responsible use.

**Key Points:**
- ✅ Free to use, modify, and distribute
- ✅ Commercial use allowed with proper authorization
- ✅ Attribution required
- ⚠️ Must comply with educational/authorized testing terms
- ⚠️ No warranty provided
- ⚠️ Authors not liable for misuse

**Required Reading:**
- [LICENSE](LICENSE) - Full MIT License text and legal disclaimer
- [LEGAL.md](LEGAL.md) - Comprehensive legal guidelines and responsible use

**Summary**: You are free to use this software for legitimate security testing, education, and research. However, you **MUST** obtain proper authorization before testing any systems you don't own. Unauthorized use is illegal and violations will be prosecuted.

## 🙏 Acknowledgments

- Original MHDDoS project for inspiration
- Go community for excellent tooling and libraries
- Contributors and testers

## 📞 Contact & Support

- **Issues**: [GitHub Issues](https://github.com/go-ddos-tools/ddos-tools/issues)
- **Discussions**: [GitHub Discussions](https://github.com/go-ddos-tools/ddos-tools/discussions)

---

**Remember**: Use this tool responsibly and legally. Unauthorized attacks are illegal and unethical.