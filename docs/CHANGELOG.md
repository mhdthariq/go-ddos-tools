# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Code Modernization (November 2025)
- **Modernized to Go 1.21+ syntax**
  - Removed custom `min()` and `max()` functions in favor of built-in functions (Go 1.21+)
  - Updated `pkg/ui/colors.go` ProgressBar to use `min()` and `max()` built-ins
  - Updated `pkg/ui/validation.go` to use built-in `min()` and `max()`
- **Modernized to Go 1.22+ range-over-int syntax**
  - Converted 19 traditional for-loops (`for i := 0; i < n; i++`) to modern `for range n` syntax
  - Updated all Layer 7 attack methods in `pkg/attacks/layer7.go` (17 methods)
  - Updated `pkg/ui/color_test.go` TestSpinner function
  - Updated `pkg/ui/validation.go` levenshteinClose function
- **Improved error handling**
  - Changed error wrapping from `%v` to `%w` for proper error chains (Go 1.13+ best practice)
  - Updated error handling in `main.go` and `pkg/proxy/proxy.go`
- **Zero compiler warnings or errors**
- **All tests passing**

### Added
- **All 26 Layer 7 attack methods now implemented (100% complete!)**
  - CFB - Cloudflare Bypass
  - BYPASS - Generic bypass method
  - OVH - OVH-specific attack
  - DYN - Dynamic host header attack
  - EVEN - Event-based attack
  - GSB - Google Search Bot simulation
  - DGB - DDoS Guard bypass
  - AVB - Anti-bot bypass
  - CFBUAM - Cloudflare Under Attack Mode bypass
  - APACHE - Apache Range header attack (CVE-2011-3192)
  - XMLRPC - XML-RPC pingback attack
  - BOT - Bot simulation (robots.txt/sitemap.xml)
  - BOMB - High-volume bombardment
  - DOWNLOADER - Download and read attack
  - KILLER - Multi-spawn GET attack
  - TOR - Tor network bypass via tor2web gateways
  - RHEX - Random hex path attack
  - STOMP - Special hex pattern attack
- Linux user agents to `pkg/utils/utils.go` for better platform diversity
  - Added Chrome on Linux user agent
  - Added Firefox on Ubuntu user agent
- Created `docs/` folder for better documentation organization
- Moved all documentation files into `docs/` folder (except README.md and LICENSE)

### Changed
- **Layer 7 implementation status: 26/26 methods (100%)**
- Reorganized project documentation structure
- Updated user agent list to include 6 default user agents covering Windows, macOS, and Linux platforms
- Enhanced attack method registry in `pkg/methods/methods.go`

### Improved
- Better cross-platform coverage in HTTP request headers
- Improved documentation accessibility and organization

## [Previous Updates]

### Previous Code Modernization
- Modernized codebase to Go 1.24+ features
  - Replaced manual slice lookups with `slices.Contains`
  - Replaced `strings.Replace(..., -1)` with `strings.ReplaceAll`
  - Modernized benchmark loops to `b.Loop()`
  - Used `unsafe.Add` for safer pointer arithmetic
  - Used `net.JoinHostPort` for IPv6-safe address formatting
  - Replaced string concatenation with `strings.Builder`

### Features
- **Implemented all 26 Layer 7 attack methods (100%)**
  - Basic: GET, POST, HEAD, STRESS, SLOW, NULL, COOKIE, PPS
  - Bypass: CFB, BYPASS, DYN, EVEN, DGB, AVB, CFBUAM
  - Advanced: OVH, GSB, APACHE, XMLRPC, BOT, BOMB, DOWNLOADER, KILLER, TOR, RHEX, STOMP
- Implemented 14 Layer 4 attack methods (100%)
- Implemented 7 Amplification attack methods (100%)
- Added comprehensive configuration via JSON
- Added proxy support (HTTP/SOCKS4/SOCKS5)
- Added rate limiting and concurrency controls

### Documentation
- Added MIT LICENSE with educational disclaimer
- Created comprehensive USAGE.md guide
- Created LEGAL.md with detailed legal guidance
- Created LEGAL-QUICK-REF.md for quick legal reference

### Testing & Quality
- All existing tests passing
- Resolved all linter warnings and diagnostics
- Fixed IPv6 address handling issues
- Fixed unsafe pointer usage issues
- Removed redundant code and checks

---

## Notes

This project is intended for **educational and authorized security testing purposes only**.
Unauthorized use of this tool against systems you do not own or have explicit permission to test is illegal.
See [LEGAL.md](LEGAL.md) for detailed legal information.

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025