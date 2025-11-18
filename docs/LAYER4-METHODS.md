# Layer 4 Attack Methods - Quick Reference Guide

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.4 SNAPSHOT

---

## Overview

This document provides a comprehensive reference for all 14 Layer 4 (Transport Layer) attack methods implemented in ddos-tools. All methods are **100% implemented and ready to use**.

Layer 4 attacks target the transport layer, focusing on TCP/UDP protocols and specific game server implementations.

## Table of Contents

- [Basic Transport Layer Methods](#basic-transport-layer-methods)
- [Game Server Methods](#game-server-methods)
- [Specialized Methods](#specialized-methods)
- [Method Comparison Table](#method-comparison-table)
- [Usage Examples](#usage-examples)
- [Best Practices](#best-practices)
- [Privilege Requirements](#privilege-requirements)

---

## Basic Transport Layer Methods

### 1. TCP - TCP Connection Flood

**Description**: Establishes TCP connections and sends random data to exhaust server resources.

**Syntax**:
```bash
./ddos-tools TCP <ip:port> <threads> <duration>
```

**With Proxies**:
```bash
./ddos-tools TCP <ip:port> <threads> <duration> <proxy_type> <proxy_file>
```

**Example**:
```bash
./ddos-tools TCP 192.168.1.1:80 100 60
./ddos-tools TCP 192.168.1.1:80 100 60 5 proxies.txt
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Low
- **Bandwidth**: Medium
- **Detection**: Easy
- **Privileges**: None required
- **Proxy Support**: Yes

**Use Cases**:
- Connection pool exhaustion
- TCP stack stress testing
- Port availability testing
- Server capacity testing

**How It Works**:
- Establishes TCP connections
- Sends random data packets
- Maintains connections
- Exhausts server connection limits

---

### 2. UDP - UDP Packet Flood

**Description**: Sends UDP packets to the target port to saturate bandwidth and resources.

**Syntax**:
```bash
./ddos-tools UDP <ip:port> <threads> <duration>
```

**With Proxies**:
```bash
./ddos-tools UDP <ip:port> <threads> <duration> <proxy_type> <proxy_file>
```

**Example**:
```bash
./ddos-tools UDP 192.168.1.1:53 200 120
./ddos-tools UDP 192.168.1.1:53 200 120 5 proxies.txt
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: High
- **Detection**: Easy
- **Privileges**: None required
- **Proxy Support**: Yes

**Use Cases**:
- Bandwidth saturation
- UDP service testing
- DNS server stress testing
- Game server flooding

**How It Works**:
- Sends UDP packets with random data
- No connection establishment needed
- High packet rate
- Bypasses TCP handshake overhead

---

### 3. SYN - SYN Flood Attack

**Description**: TCP SYN flood attack that exploits the TCP three-way handshake.

**Syntax**:
```bash
sudo ./ddos-tools SYN <ip:port> <threads> <duration>
```

**Example**:
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Medium
- **Bandwidth**: Low
- **Detection**: Medium
- **Privileges**: **ROOT/ADMIN REQUIRED**
- **Proxy Support**: No

**Use Cases**:
- SYN queue exhaustion
- Connection table overflow
- Firewall rule testing
- TCP stack vulnerability testing

**How It Works**:
- Sends TCP SYN packets
- Does not complete handshake
- Fills server's SYN queue
- Prevents legitimate connections

⚠️ **Warning**: Requires root/administrator privileges for raw socket creation.

---

### 4. ICMP - ICMP Echo Flood

**Description**: ICMP echo request (ping) flood attack.

**Syntax**:
```bash
sudo ./ddos-tools ICMP <ip:port> <threads> <duration>
```

**Example**:
```bash
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Low
- **Bandwidth**: Medium
- **Detection**: Easy
- **Privileges**: **ROOT/ADMIN REQUIRED**
- **Proxy Support**: No

**Use Cases**:
- ICMP rate limiting testing
- Network saturation
- Firewall ICMP rule testing
- Basic connectivity flooding

**How It Works**:
- Sends ICMP echo request packets
- Can saturate network bandwidth
- Tests ICMP handling
- Bypasses some firewalls

⚠️ **Warning**: Requires root/administrator privileges.

---

### 5. OVH-UDP - OVH-Optimized UDP Flood

**Description**: OVH-specific UDP flood optimized for OVH infrastructure.

**Syntax**:
```bash
./ddos-tools OVH-UDP <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools OVH-UDP 192.168.1.1:80 200 90
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Medium
- **Bandwidth**: High
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Testing OVH DDoS protection
- OVH-hosted server testing
- Bypass OVH rate limiting

**How It Works**:
- Optimized UDP packet structure
- Targets OVH network characteristics
- High packet rate
- Specialized for OVH infrastructure

---

## Game Server Methods

### 6. MINECRAFT - Minecraft Server Flood

**Description**: Specialized attack for Minecraft Java Edition servers using server list ping.

**Syntax**:
```bash
./ddos-tools MINECRAFT <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools MINECRAFT mc.example.com:25565 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: High
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Minecraft server capacity testing
- Server list ping flood
- Testing anti-bot protection
- Query system testing

**How It Works**:
- Sends Minecraft protocol packets
- Uses server list ping format
- Exploits query response overhead
- Protocol version configurable via config.json

**Target**: Minecraft Java Edition servers (default port 25565)

---

### 7. MCBOT - Minecraft Bot Spam Attack

**Description**: Simulates Minecraft bot connections with login attempts.

**Syntax**:
```bash
./ddos-tools MCBOT <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools MCBOT mc.example.com:25565 50 120
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: High
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Bot protection testing
- Login rate limiting
- Authentication system stress
- Player slot exhaustion

**How It Works**:
- Initiates Minecraft login sequence
- Random player names
- Complete handshake process
- Tests authentication overhead

---

### 8. MCPE - Minecraft Pocket Edition Flood

**Description**: Targets Minecraft Bedrock/Pocket Edition servers.

**Syntax**:
```bash
./ddos-tools MCPE <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools MCPE pe.example.com:19132 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Bedrock Edition server testing
- Mobile server stress testing
- RakNet protocol testing

**How It Works**:
- Uses RakNet protocol
- Unconnected ping packets
- Bedrock-specific packet structure

**Target**: Minecraft Bedrock/Pocket Edition (default port 19132)

---

### 9. VSE - Valve Source Engine Query Flood

**Description**: Floods Valve Source Engine game servers with query packets.

**Syntax**:
```bash
./ddos-tools VSE <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools VSE 192.168.1.1:27015 100 60
```

**Characteristics**:
- **Speed**: Fast
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Source Engine server testing
- Query rate limiting
- Server browser flood
- A2S_INFO packet testing

**How It Works**:
- Sends A2S_INFO query packets
- Source Engine protocol
- Server info request flood
- Response amplification possible

**Target Games**: CS:GO, TF2, Garry's Mod, Half-Life 2, etc.

---

### 10. TS3 - TeamSpeak 3 Query Flood

**Description**: Floods TeamSpeak 3 servers with voice and query packets.

**Syntax**:
```bash
./ddos-tools TS3 <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools TS3 192.168.1.1:9987 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- TeamSpeak server capacity testing
- Voice packet handling
- Query system stress
- Anti-flood protection testing

**How It Works**:
- Sends TeamSpeak 3 protocol packets
- UDP-based voice packets
- Query request flood
- Tests voice server limits

**Default Port**: 9987 (voice), 10011 (query)

---

### 11. FIVEM - FiveM Server Info Flood

**Description**: Targets FiveM (GTA V multiplayer) servers with info requests.

**Syntax**:
```bash
./ddos-tools FIVEM <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools FIVEM 192.168.1.1:30120 100 60
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: Medium
- **Bandwidth**: Medium
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- FiveM server stress testing
- Server info query flood
- Player list request spam
- Server browser protection testing

**How It Works**:
- Sends getinfo packets
- FiveM-specific protocol
- Server status requests
- Response overhead exploitation

**Default Port**: 30120

---

### 12. FIVEM-TOKEN - FiveM Token Request Flood

**Description**: FiveM token/authentication request flood.

**Syntax**:
```bash
./ddos-tools FIVEM-TOKEN <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools FIVEM-TOKEN 192.168.1.1:30120 50 90
```

**Characteristics**:
- **Speed**: Medium
- **Complexity**: High
- **Bandwidth**: Low
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Authentication system testing
- Token generation stress
- Connection handshake flood

**How It Works**:
- Requests authentication tokens
- Connection establishment attempts
- Tests authentication overhead
- Token generation exhaustion

---

## Specialized Methods

### 13. CONNECTION - Connection Hold Attack

**Description**: Establishes and holds connections open to exhaust connection pools.

**Syntax**:
```bash
./ddos-tools CONNECTION <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools CONNECTION 192.168.1.1:80 200 300
```

**Characteristics**:
- **Speed**: Slow
- **Complexity**: Low
- **Bandwidth**: Low
- **Detection**: Medium
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Connection pool exhaustion
- Long-lived connection testing
- Keep-alive abuse
- Connection limit testing

**How It Works**:
- Establishes TCP connections
- Keeps connections alive
- Minimal data transfer
- Exhausts available connection slots

---

### 14. CPS - Connections Per Second Test

**Description**: Measures and stresses connections per second capacity.

**Syntax**:
```bash
./ddos-tools CPS <ip:port> <threads> <duration>
```

**Example**:
```bash
./ddos-tools CPS 192.168.1.1:80 100 60
```

**Characteristics**:
- **Speed**: Very Fast
- **Complexity**: Low
- **Bandwidth**: Low
- **Detection**: Easy
- **Privileges**: None required
- **Proxy Support**: No

**Use Cases**:
- Connection rate limiting testing
- Server accept() queue testing
- Connection establishment overhead
- Rate limiting validation

**How It Works**:
- Rapid connection establishment
- Immediate disconnection
- High connection rate
- Tests handshake processing

---

## Method Comparison Table

| Method | Speed | Complexity | Bandwidth | Detection | Privileges | Proxy Support |
|--------|-------|------------|-----------|-----------|------------|---------------|
| TCP | Fast | Low | Medium | Easy | None | Yes |
| UDP | Very Fast | Low | High | Easy | None | Yes |
| SYN | Very Fast | Medium | Low | Medium | ROOT | No |
| ICMP | Fast | Low | Medium | Easy | ROOT | No |
| OVH-UDP | Fast | Medium | High | Medium | None | No |
| MINECRAFT | Medium | High | Medium | Medium | None | No |
| MCBOT | Slow | High | Medium | Medium | None | No |
| MCPE | Medium | Medium | Medium | Medium | None | No |
| VSE | Fast | Medium | Medium | Medium | None | No |
| TS3 | Medium | Medium | Medium | Medium | None | No |
| FIVEM | Medium | Medium | Medium | Medium | None | No |
| FIVEM-TOKEN | Medium | High | Low | Medium | None | No |
| CONNECTION | Slow | Low | Low | Medium | None | No |
| CPS | Very Fast | Low | Low | Easy | None | No |

---

## Usage Examples

### Basic Attacks

```bash
# TCP flood
./ddos-tools TCP 192.168.1.1:80 100 60

# UDP flood on DNS server
./ddos-tools UDP 192.168.1.1:53 200 120

# ICMP flood (requires root)
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

### With Proxies

```bash
# TCP flood with SOCKS5 proxies
./ddos-tools TCP 192.168.1.1:80 100 60 5 proxies.txt

# UDP flood with HTTP proxies
./ddos-tools UDP 192.168.1.1:53 200 120 1 proxies.txt
```

### Game Server Attacks

```bash
# Minecraft Java Edition
./ddos-tools MINECRAFT mc.example.com:25565 100 60

# Minecraft Bedrock Edition
./ddos-tools MCPE pe.example.com:19132 100 60

# CS:GO server
./ddos-tools VSE 192.168.1.1:27015 100 60

# TeamSpeak 3
./ddos-tools TS3 192.168.1.1:9987 100 60

# FiveM server
./ddos-tools FIVEM 192.168.1.1:30120 100 60
```

### Advanced Attacks

```bash
# SYN flood (requires root)
sudo ./ddos-tools SYN 192.168.1.1:80 500 60

# Connection exhaustion
./ddos-tools CONNECTION 192.168.1.1:80 200 300

# Connections per second test
./ddos-tools CPS 192.168.1.1:80 100 60
```

---

## Best Practices

### Method Selection

**For Maximum Speed:**
- UDP, SYN, CPS

**For Stealth:**
- CONNECTION, MCBOT, FIVEM-TOKEN

**For Game Servers:**
- MINECRAFT (Java Edition)
- MCPE (Bedrock Edition)
- VSE (Source Engine games)
- TS3 (TeamSpeak)
- FIVEM (GTA V multiplayer)

**For Connection Exhaustion:**
- CONNECTION, TCP, SYN

**For Bandwidth Saturation:**
- UDP, OVH-UDP

### Parameter Recommendations

**Threads:**
- Light attack: 10-50 threads
- Medium attack: 50-200 threads
- Heavy attack: 200-500 threads
- SYN flood: 500-1000 threads

**Duration:**
- Quick test: 30-60 seconds
- Standard test: 60-300 seconds
- Extended test: 300-600 seconds
- Connection hold: 600+ seconds

**Target Ports:**
- Web servers: 80 (HTTP), 443 (HTTPS)
- DNS: 53
- Minecraft Java: 25565
- Minecraft Bedrock: 19132
- Source Engine: 27015
- TeamSpeak 3: 9987
- FiveM: 30120

### Proxy Usage

**Supported Methods:**
- TCP (all proxy types)
- UDP (all proxy types)

**Not Supported:**
- SYN (requires raw sockets)
- ICMP (requires raw sockets)
- Game-specific methods (protocol restrictions)

**Recommended Proxy Type:**
- SOCKS5 (type 5) for best compatibility

---

## Privilege Requirements

### Requires Root/Administrator

**SYN Flood:**
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
```

**ICMP Flood:**
```bash
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

**Reason**: Raw socket creation requires elevated privileges.

### No Privileges Required

All other methods work without special privileges:
- TCP, UDP, CONNECTION, CPS, OVH-UDP
- MINECRAFT, MCBOT, MCPE
- VSE, TS3, FIVEM, FIVEM-TOKEN

---

## Performance Optimization

### Thread Tuning

**Low-end systems (1-2 cores):**
- 10-50 threads

**Medium systems (4-8 cores):**
- 50-200 threads

**High-end systems (8+ cores):**
- 200-500+ threads

### Memory Considerations

**Low memory (<2GB):**
- Use simple methods: TCP, UDP, ICMP
- Limit threads to 50-100

**High memory (8GB+):**
- Can use all methods
- Higher thread counts
- Multiple concurrent attacks

### Network Bandwidth

**Limited bandwidth (<10Mbps):**
- Use SYN, CONNECTION, CPS
- Lower thread counts

**High bandwidth (100Mbps+):**
- Use UDP, OVH-UDP
- Maximum thread counts
- Bandwidth saturation attacks

---

## Troubleshooting

### Common Issues

**Issue**: "Permission denied" on SYN/ICMP
- **Solution**: Run with sudo/administrator privileges

**Issue**: Connection refused
- **Solution**: Verify target IP and port are correct

**Issue**: Low packet rate
- **Solution**: Increase thread count or use faster method (UDP, SYN)

**Issue**: Proxies not working
- **Solution**: Verify proxy file format and use supported method (TCP/UDP)

**Issue**: Game server method not working
- **Solution**: Verify correct port and server type

---

## Legal Notice

⚠️ **IMPORTANT**: All Layer 4 methods must only be used on systems you own or have explicit written authorization to test.

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
- ❌ Game servers without permission

**Special Warnings:**
- SYN floods can affect entire networks
- UDP floods can saturate bandwidth
- Game server attacks can impact legitimate players
- Always obtain explicit permission

Violation may result in criminal prosecution under computer fraud laws.

See [LEGAL.md](LEGAL.md) for comprehensive legal guidelines.

---

## Additional Resources

- **[LAYER7-METHODS.md](LAYER7-METHODS.md)** - Layer 7 attack methods reference
- **[USAGE.md](USAGE.md)** - Complete usage guide

- **[PROXIES.md](PROXIES.md)** - Proxy configuration guide
- **[README.md](../README.md)** - Project overview

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Copyright**: © 2025 Muhammad Thariq  
**License**: MIT with Educational Use Terms

**Remember**: Use responsibly and legally. All 14 Layer 4 methods are powerful tools that should only be used for legitimate security testing with proper authorization.