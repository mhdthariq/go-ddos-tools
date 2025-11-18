# Proxy Configuration Guide

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.4 SNAPSHOT

---

## Overview

This guide explains how to configure and use proxies with ddos-tools. Proxies help distribute traffic, avoid rate limiting, and maintain anonymity during authorized security testing.

## Table of Contents

- [What Are Proxies?](#what-are-proxies)
- [Proxy Types](#proxy-types)
- [Creating proxies.txt](#creating-proxiestxt)
- [Proxy File Format](#proxy-file-format)
- [Where to Get Proxies](#where-to-get-proxies)
- [Testing Proxies](#testing-proxies)
- [Usage Examples](#usage-examples)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## What Are Proxies?

Proxies are intermediate servers that forward your requests to the target server. They provide several benefits:

### Benefits of Using Proxies

1. **Traffic Distribution**: Spread requests across multiple IP addresses
2. **Rate Limit Evasion**: Bypass IP-based rate limiting
3. **Geographic Diversity**: Route traffic through different regions
4. **Anonymity**: Hide your real IP address
5. **Load Balancing**: Distribute attack load across multiple sources

### How Proxies Work

```
Your Computer → Proxy Server → Target Server
     ↑                            ↓
     └──────── Response ←─────────┘
```

When using proxies:
1. Your tool connects to the proxy
2. Proxy forwards your request to the target
3. Target sees proxy IP, not yours
4. Response comes back through proxy

---

## Proxy Types

ddos-tools supports three types of proxies:

### 1. HTTP Proxies (Type 1)

**Description**: Standard HTTP/HTTPS proxies

**Format**:
```
ip:port
ip:port:username:password
```

**Example**:
```
192.168.1.100:8080
proxy.example.com:3128
10.0.0.5:8888:user:pass
```

**Best For**:
- Layer 7 (HTTP/HTTPS) attacks
- Web traffic
- GET, POST, HEAD methods

**Limitations**:
- HTTP/HTTPS traffic only
- May not support all methods

---

### 2. SOCKS4 Proxies (Type 4)

**Description**: SOCKS version 4 proxies

**Format**:
```
ip:port
```

**Example**:
```
192.168.1.101:1080
socks4.example.com:9050
```

**Best For**:
- TCP connections
- Layer 4 attacks
- General purpose

**Limitations**:
- No authentication support
- No UDP support
- IPv4 only

---

### 3. SOCKS5 Proxies (Type 5) - **RECOMMENDED**

**Description**: SOCKS version 5 proxies (most versatile)

**Format**:
```
ip:port
ip:port:username:password
```

**Example**:
```
192.168.1.102:1080
socks5.example.com:1080:user:pass
```

**Best For**:
- All attack types
- Both TCP and UDP
- Layer 4 and Layer 7
- IPv4 and IPv6

**Advantages**:
- ✅ Supports authentication
- ✅ Supports UDP
- ✅ Supports IPv6
- ✅ Most compatible
- ✅ Best performance

**Recommended**: Use SOCKS5 when available for maximum compatibility.

---

## Creating proxies.txt

### Step 1: Create the File

Create a text file named `proxies.txt` in your project directory or `files/proxies/` folder.

**Using Command Line:**

```bash
# Linux/macOS
touch proxies.txt

# Windows
type nul > proxies.txt
```

**Using Text Editor:**
- Open Notepad (Windows) or any text editor
- Save as `proxies.txt`

---

### Step 2: Add Proxy Entries

Add one proxy per line in the correct format:

**Example proxies.txt:**
```
# HTTP Proxies
192.168.1.100:8080
proxy1.example.com:3128
10.0.0.5:8888:user:pass

# SOCKS4 Proxies
192.168.1.101:1080
socks4.example.com:9050

# SOCKS5 Proxies (RECOMMENDED)
192.168.1.102:1080
socks5.example.com:1080:admin:secret123
```

---

### Step 3: Save the File

Save the file with these settings:
- **Encoding**: UTF-8 (no BOM)
- **Line Endings**: Unix (LF) or Windows (CRLF) - both work
- **Extension**: `.txt`

---

## Proxy File Format

### Basic Format (No Authentication)

```
IP:PORT
```

**Examples:**
```
192.168.1.100:8080
proxy.example.com:3128
10.0.0.5:1080
```

---

### Format with Authentication

```
IP:PORT:USERNAME:PASSWORD
```

**Examples:**
```
192.168.1.100:8080:admin:password123
proxy.example.com:3128:user:secretpass
10.0.0.5:1080:john:doe123
```

---

### Mixed Proxy File

You can mix different formats in the same file:

```
# Fast HTTP proxies
192.168.1.100:8080
192.168.1.101:8080:user:pass

# SOCKS5 proxies
10.0.0.5:1080
10.0.0.6:1080:admin:secret

# Domain-based proxies
proxy1.example.com:3128
proxy2.example.com:1080:user:pass
```

---

### Comments and Blank Lines

```
# This is a comment - lines starting with # are ignored

192.168.1.100:8080        # HTTP proxy

# Blank lines are ignored too

192.168.1.101:1080        # SOCKS5 proxy
```

---

## Where to Get Proxies

### Free Proxy Lists

⚠️ **Warning**: Free proxies are often unreliable, slow, and may log your traffic.

**Popular Sources:**
- Free Proxy List (free-proxy-list.net)
- ProxyScrape (proxyscrape.com/free-proxy-list)
- GeoNode (geonode.com/free-proxy-list)
- Spys.one (spys.one/en/)

**Free Proxy Caveats:**
- ❌ Often offline or slow
- ❌ May log your traffic
- ❌ Short lifespan
- ❌ Many are blacklisted
- ❌ Limited anonymity

---

### Paid Proxy Services (Recommended)

✅ **Advantages**: Reliable, fast, private, good uptime

**Popular Providers:**
- Bright Data (brightdata.com)
- Smartproxy (smartproxy.com)
- Oxylabs (oxylabs.io)
- IPRoyal (iproyal.com)
- Webshare (webshare.io)

**Paid Proxy Benefits:**
- ✅ High uptime (99%+)
- ✅ Fast speeds
- ✅ No logging
- ✅ Large IP pools
- ✅ Geographic targeting
- ✅ Support and documentation

---

### Self-Hosted Proxies

Run your own proxy servers for maximum control:

**Options:**
- Squid Proxy (HTTP)
- Shadowsocks (SOCKS5)
- 3proxy (Multi-protocol)
- Dante (SOCKS)
- TinyProxy (HTTP)

**Benefits:**
- ✅ Full control
- ✅ No third-party logging
- ✅ Customizable
- ✅ Cost-effective at scale

**Example: Install Squid on Ubuntu:**
```bash
sudo apt update
sudo apt install squid
sudo systemctl enable squid
sudo systemctl start squid
```

Default Squid config listens on port 3128.

---

### Scraping Your Own Proxies

Use proxy scraper tools to gather free proxies:

**Tools:**
- ProxyBroker (Python)
- Proxy-Scraper (Python)
- ProxyChains (Shell)

**Example with ProxyBroker:**
```bash
pip install proxybroker
proxybroker find --types SOCKS5 HTTP --limit 100 --outfile proxies.txt
```

---

## Testing Proxies

### Manual Testing

Test a single proxy using curl:

```bash
# HTTP proxy
curl -x http://192.168.1.100:8080 http://httpbin.org/ip

# SOCKS5 proxy
curl --socks5 192.168.1.102:1080 http://httpbin.org/ip

# With authentication
curl -x http://user:pass@192.168.1.100:8080 http://httpbin.org/ip
```

---

### Automated Testing Tools

**ProxyChains:**
```bash
sudo apt install proxychains
# Configure in /etc/proxychains.conf
proxychains curl http://httpbin.org/ip
```

**Custom Python Script:**
```python
import requests

def test_proxy(proxy):
    try:
        r = requests.get('http://httpbin.org/ip', 
                        proxies={'http': proxy, 'https': proxy}, 
                        timeout=5)
        print(f"✅ {proxy}: {r.json()['origin']}")
        return True
    except:
        print(f"❌ {proxy}: Failed")
        return False

# Test from file
with open('proxies.txt') as f:
    for line in f:
        line = line.strip()
        if line and not line.startswith('#'):
            proxy = f"http://{line}"
            test_proxy(proxy)
```

---

### Built-in Testing (Future Feature)

In a future version, ddos-tools will include:
```bash
# Coming soon
./ddos-tools CHECK-PROXIES proxies.txt
```

---

## Usage Examples

### Basic Usage

**Syntax:**
```bash
./ddos-tools <method> <target> <proxy_type> <threads> <proxy_file> <rpc> <duration>
```

---

### Layer 7 with Proxies

```bash
# GET flood with HTTP proxies (type 1)
./ddos-tools GET http://example.com 1 100 proxies.txt 100 60

# POST flood with SOCKS4 proxies (type 4)
./ddos-tools POST https://example.com/api 4 200 proxies.txt 50 120

# STRESS with SOCKS5 proxies (type 5) - RECOMMENDED
./ddos-tools STRESS http://example.com 5 100 proxies.txt 100 90

# Using all proxy types (type 0)
./ddos-tools GET http://example.com 0 100 proxies.txt 100 60

# Random proxy selection (type 6)
./ddos-tools POST http://example.com 6 100 proxies.txt 100 60
```

---

### Layer 4 with Proxies

```bash
# TCP flood with SOCKS5 proxies
./ddos-tools TCP 192.168.1.1:80 100 60 5 proxies.txt

# UDP flood with SOCKS5 proxies
./ddos-tools UDP 192.168.1.1:53 200 120 5 proxies.txt
```

**Note**: Only TCP and UDP Layer 4 methods support proxies.

---

### Proxy Type Selection

```bash
# Type 0: Use all available proxies
./ddos-tools GET http://example.com 0 100 proxies.txt 100 60

# Type 1: HTTP proxies only
./ddos-tools GET http://example.com 1 100 http-proxies.txt 100 60

# Type 4: SOCKS4 proxies only
./ddos-tools GET http://example.com 4 100 socks4-proxies.txt 100 60

# Type 5: SOCKS5 proxies only (RECOMMENDED)
./ddos-tools GET http://example.com 5 100 socks5-proxies.txt 100 60

# Type 6: Random selection from all types
./ddos-tools GET http://example.com 6 100 mixed-proxies.txt 100 60
```

---

### Multiple Proxy Files

Organize proxies by type or quality:

```bash
# High-quality paid proxies
./ddos-tools GET http://example.com 5 50 premium-proxies.txt 100 60

# Free proxies (backup)
./ddos-tools GET http://example.com 5 200 free-proxies.txt 100 60

# Geographic-specific
./ddos-tools GET http://example.com 5 100 us-proxies.txt 100 60
```

---

## Troubleshooting

### Proxies Not Working

**Problem**: Attack not starting or connection errors

**Solutions:**
1. Verify proxy file format
2. Test proxies manually
3. Check proxy type matches file content
4. Ensure proxies support the protocol (HTTP/SOCKS)
5. Try different proxy type

```bash
# Test with different types
./ddos-tools GET http://example.com 1 10 proxies.txt 10 30  # HTTP
./ddos-tools GET http://example.com 5 10 proxies.txt 10 30  # SOCKS5
```

---

### Slow Performance

**Problem**: Attack is very slow

**Solutions:**
1. Use faster proxies (paid > free)
2. Reduce thread count
3. Use SOCKS5 instead of HTTP
4. Filter slow proxies from list
5. Use geographically closer proxies

---

### Authentication Errors

**Problem**: Proxy requires authentication but fails

**Solutions:**
1. Verify username:password format
2. Check for special characters in password
3. URL-encode special characters
4. Test proxy manually with curl

**Example with special characters:**
```
# If password is: p@ss:word
# Encode as: p%40ss%3Aword
192.168.1.100:8080:user:p%40ss%3Aword
```

---

### Connection Refused

**Problem**: Proxy refuses connection

**Solutions:**
1. Check proxy is online (ping IP)
2. Verify port is correct
3. Check firewall rules
4. Try different proxy
5. Test proxy with curl

```bash
# Test if proxy is reachable
ping 192.168.1.100
nc -zv 192.168.1.100 8080
```

---

### High Failure Rate

**Problem**: Many proxies failing

**Solutions:**
1. Update proxy list (free proxies die quickly)
2. Use paid proxies for better reliability
3. Filter tested/working proxies
4. Increase timeout values
5. Reduce concurrent connections per proxy

---

## Best Practices

### Proxy List Management

**DO:**
- ✅ Keep proxies updated regularly
- ✅ Test proxies before use
- ✅ Remove dead proxies from list
- ✅ Organize by type and quality
- ✅ Use separate files for different purposes
- ✅ Document proxy sources

**DON'T:**
- ❌ Use old proxy lists
- ❌ Mix incompatible proxy types
- ❌ Overload individual proxies
- ❌ Share proxy lists publicly
- ❌ Use untrusted free proxies for sensitive tests

---

### Performance Optimization

**1. Proxy Quality:**
- Use paid proxies for important tests
- Filter by response time (<1s)
- Prefer SOCKS5 over HTTP
- Choose geographically close proxies

**2. Thread Distribution:**
```
Total Threads = Number of Proxies × Threads per Proxy
```

Example:
- 100 proxies × 2 threads each = 200 total threads
- Better than 10 proxies × 20 threads each

**3. File Organization:**
```
files/proxies/
├── premium-socks5.txt        # Paid SOCKS5 (best)
├── free-socks5.txt           # Free SOCKS5 (backup)
├── http-proxies.txt          # HTTP only
├── us-proxies.txt            # US-based
├── eu-proxies.txt            # EU-based
└── working-proxies.txt       # Recently tested
```

---

### Security Considerations

**1. Proxy Trust:**
- Free proxies may log traffic
- Use HTTPS for sensitive data
- Assume proxies can see your requests
- Never send credentials through untrusted proxies

**2. Anonymity Layers:**
```
Your IP → VPN → Proxy → Target
```
- Add VPN before proxies for extra anonymity
- Use Tor for maximum privacy
- Rotate proxies frequently

**3. Legal Compliance:**
- Only use proxies for authorized testing
- Ensure proxy usage is legal in your jurisdiction
- Don't use proxies to hide illegal activity
- Keep logs of authorized testing

---

### Proxy Rotation Strategy

**Manual Rotation:**
```bash
# Test 1: First proxy set
./ddos-tools GET http://example.com 5 100 proxies-set1.txt 100 60

# Test 2: Second proxy set
./ddos-tools GET http://example.com 5 100 proxies-set2.txt 100 60
```

**Automatic Rotation:**
- Tool automatically rotates through proxy list
- Each thread gets a random proxy
- Failed proxies are skipped
- Built-in load balancing

---

## Example Proxy Lists

### Small Test List (10 proxies)

**test-proxies.txt:**
```
# For quick testing
192.168.1.100:8080
192.168.1.101:8080
192.168.1.102:1080
192.168.1.103:1080
192.168.1.104:3128
192.168.1.105:3128
192.168.1.106:1080:user:pass
192.168.1.107:1080:admin:secret
192.168.1.108:8080:test:test123
192.168.1.109:3128:proxy:password
```

---

### Production List (100+ proxies)

**production-proxies.txt:**
```
# Premium SOCKS5 proxies - Updated: 2025-11-15

# US East Region
us-proxy-001.example.com:1080:user:pass
us-proxy-002.example.com:1080:user:pass
us-proxy-003.example.com:1080:user:pass
# ... (97 more)

# Backup HTTP proxies
http-backup-001.example.com:8080
http-backup-002.example.com:8080
# ... (more backups)
```

---

### Geographic-Specific Lists

**us-proxies.txt:**
```
# United States proxies only
us-ny-01.example.com:1080:user:pass
us-ca-01.example.com:1080:user:pass
us-tx-01.example.com:1080:user:pass
```

**eu-proxies.txt:**
```
# European Union proxies only
eu-uk-01.example.com:1080:user:pass
eu-de-01.example.com:1080:user:pass
eu-fr-01.example.com:1080:user:pass
```

---

## Advanced Topics

### Proxy Chaining

Route through multiple proxies:

```
Your PC → Proxy 1 → Proxy 2 → Proxy 3 → Target
```

**Setup with ProxyChains:**
```bash
# /etc/proxychains.conf
[ProxyList]
socks5  192.168.1.100  1080  user  pass
socks5  192.168.1.101  1080  user  pass
socks5  192.168.1.102  1080  user  pass
```

---

### Proxy Performance Testing

**Measure Response Time:**
```bash
time curl -x http://192.168.1.100:8080 http://example.com
```

**Benchmark Multiple Proxies:**
```bash
for proxy in $(cat proxies.txt); do
    echo "Testing $proxy"
    time curl -x http://$proxy http://httpbin.org/ip --max-time 5
done
```

---

### Creating a Proxy Validator Script

**proxy-validator.sh:**
```bash
#!/bin/bash

INPUT_FILE="proxies.txt"
OUTPUT_FILE="working-proxies.txt"
TIMEOUT=5

> $OUTPUT_FILE

while IFS= read -r proxy; do
    # Skip comments and empty lines
    [[ "$proxy" =~ ^#.*$ ]] && continue
    [[ -z "$proxy" ]] && continue
    
    echo "Testing: $proxy"
    if curl -x "http://$proxy" "http://httpbin.org/ip" \
       --max-time $TIMEOUT --silent > /dev/null 2>&1; then
        echo "✅ $proxy"
        echo "$proxy" >> $OUTPUT_FILE
    else
        echo "❌ $proxy"
    fi
done < "$INPUT_FILE"

echo "Working proxies saved to $OUTPUT_FILE"
```

---

## Legal Notice

⚠️ **IMPORTANT**: Proxy usage for security testing must comply with all applicable laws.

**Legal Requirements:**
- ✅ Only use proxies for authorized testing
- ✅ Ensure proxy provider allows your use case
- ✅ Don't use proxies to hide illegal activity
- ✅ Comply with proxy service terms of service
- ✅ Keep authorization documentation

**Prohibited Activities:**
- ❌ Using proxies for unauthorized attacks
- ❌ Violating proxy provider terms
- ❌ Evading law enforcement
- ❌ Accessing restricted content illegally
- ❌ Hiding malicious activity

See [LEGAL.md](LEGAL.md) for comprehensive legal guidelines.

---

## Additional Resources

- **[LAYER7-METHODS.md](LAYER7-METHODS.md)** - Layer 7 attack methods
- **[LAYER4-METHODS.md](LAYER4-METHODS.md)** - Layer 4 attack methods
- **[USAGE.md](USAGE.md)** - Complete usage guide

- **[README.md](../README.md)** - Project overview

### External Resources

- **Proxy Providers:**
  - Bright Data: https://brightdata.com
  - Smartproxy: https://smartproxy.com
  - Webshare: https://webshare.io

- **Free Proxy Lists:**
  - Free Proxy List: https://free-proxy-list.net
  - ProxyScrape: https://proxyscrape.com
  - GeoNode: https://geonode.com/free-proxy-list

- **Proxy Tools:**
  - ProxyChains: https://github.com/haad/proxychains
  - ProxyBroker: https://github.com/constverum/ProxyBroker
  - Squid Proxy: http://www.squid-cache.org

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Copyright**: © 2025 Muhammad Thariq  
**License**: MIT with Educational Use Terms

**Remember**: Proxies are powerful tools for legitimate security testing. Always use them responsibly, legally, and only with proper authorization.