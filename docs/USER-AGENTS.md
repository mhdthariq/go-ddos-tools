# User Agent Documentation

## Overview

This document describes the user agent functionality in the ddos-tools project. User agents are HTTP headers that identify the client software making requests to web servers. Proper user agent rotation is essential for realistic traffic simulation during authorized security testing.

## Implementation

User agents are managed in `pkg/utils/utils.go` through the `GetDefaultUserAgents()` function.

### Default User Agents

The tool includes 6 default user agents covering major platforms and browsers:

#### Windows Platform
1. **Chrome 74 on Windows 10**
   ```
   Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36
   ```

2. **Chrome 77 on Windows 10**
   ```
   Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36
   ```

3. **Firefox 69 on Windows 10**
   ```
   Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101 Firefox/69.0
   ```

#### macOS Platform
4. **Safari 14 on macOS Big Sur**
   ```
   Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15
   ```

#### Linux Platform
5. **Chrome 91 on Linux**
   ```
   Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36
   ```

6. **Firefox 89 on Ubuntu**
   ```
   Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0
   ```

## Platform Coverage

| Platform | Browser | Count | Percentage |
|----------|---------|-------|------------|
| Windows  | Chrome, Firefox | 3 | 50% |
| macOS    | Safari | 1 | 17% |
| Linux    | Chrome, Firefox | 2 | 33% |

## Usage

### Using Default User Agents

Default user agents are automatically used when no custom user agent file is specified:

```go
userAgents := utils.GetDefaultUserAgents()
randomUA := userAgents[rand.Intn(len(userAgents))]
```

### Custom User Agent Files

You can provide a custom list of user agents via a text file:

1. Create a text file with one user agent per line:
   ```
   Mozilla/5.0 (custom user agent 1)
   Mozilla/5.0 (custom user agent 2)
   Mozilla/5.0 (custom user agent 3)
   ```

2. Specify the file in your configuration:
   ```json
   {
     "useragent_file": "path/to/useragents.txt"
   }
   ```

3. The tool will load and randomly select from your custom list

### Command Line Usage

```bash
# Use default user agents (automatic)
./ddos-tools -target https://example.com -method GET -threads 10

# Use custom user agent file
./ddos-tools -target https://example.com -method GET -threads 10 -ua-file my-useragents.txt
```

## Best Practices

### For Security Testing

1. **Diversity**: Use a diverse set of user agents to simulate realistic traffic patterns
2. **Rotation**: Rotate user agents frequently to avoid detection
3. **Relevance**: Choose user agents that match your testing scenario
4. **Updates**: Keep user agents up-to-date with current browser versions

### For Production Testing

1. **Custom Lists**: Maintain your own user agent lists based on your actual traffic
2. **Analytics**: Use user agents that match your site's visitor analytics
3. **Mobile Support**: Include mobile user agents if testing mobile endpoints
4. **Bot Identification**: Clearly identify test traffic if required by policy

## Adding New User Agents

To add new default user agents to the codebase:

1. Edit `pkg/utils/utils.go`
2. Locate the `GetDefaultUserAgents()` function
3. Add new user agent strings to the returned slice:

```go
func GetDefaultUserAgents() []string {
	return []string{
		// Existing user agents...
		"Mozilla/5.0 (new user agent here)",
	}
}
```

4. Consider platform balance and version diversity
5. Test that the new user agents are valid and accepted by target servers

## Security Considerations

### Why User Agent Rotation Matters

1. **Realism**: Real traffic comes from diverse clients
2. **Detection Avoidance**: Single user agent can trigger rate limiting
3. **Testing Accuracy**: Proper headers ensure realistic test results
4. **WAF/CDN Testing**: Some protections specifically analyze user agents

### User Agent Fingerprinting

Modern web application firewalls (WAFs) and content delivery networks (CDNs) may:

- Validate user agent format and consistency
- Check user agent against TLS fingerprints
- Compare user agent with HTTP/2 settings
- Track user agent anomalies across requests

For accurate testing, ensure:

1. User agents match realistic browser capabilities
2. HTTP headers are consistent with the claimed browser
3. TLS settings align with the user agent version
4. JavaScript fingerprints match if testing client-side code

## Resources

### User Agent Databases

- [WhatIsMyBrowser User Agent Database](https://www.whatismybrowser.com/guides/the-latest-user-agent/)
- [User Agent String.com](http://useragentstring.com/)
- [MDN User-Agent Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent)

### Testing Tools

- [User Agent Parser](https://developers.whatismybrowser.com/useragents/parse/)
- [Chrome User Agent Switcher Extension](https://chrome.google.com/webstore/detail/user-agent-switcher/)
- [Firefox User Agent Switcher Extension](https://addons.mozilla.org/en-US/firefox/addon/user-agent-string-switcher/)

## Troubleshooting

### Common Issues

**Issue**: Requests blocked despite user agent rotation
- **Solution**: Check if other headers (Accept, Accept-Language, etc.) match the user agent
- **Solution**: Verify TLS fingerprint matches the claimed browser version

**Issue**: User agents appear outdated
- **Solution**: Update the user agent list with current browser versions
- **Solution**: Use a custom user agent file with fresh data

**Issue**: Mobile endpoints not accessible
- **Solution**: Add mobile user agents to your custom list
- **Solution**: Include appropriate mobile headers (Accept, Viewport, etc.)

## Examples

### Layer 7 GET Attack with User Agent Rotation

```go
package main

import (
	"ddos-tools/pkg/utils"
	"net/http"
	"math/rand"
)

func makeRequest(url string) {
	userAgents := utils.GetDefaultUserAgents()
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}
```

### Custom User Agent Pool

```go
package main

import (
	"ddos-tools/pkg/utils"
)

func main() {
	// Load custom user agents
	customUAs, err := utils.LoadLines("custom-ua.txt")
	if err != nil {
		// Fallback to defaults
		customUAs = utils.GetDefaultUserAgents()
	}
	
	// Use in your attack logic
	for i := 0; i < 100; i++ {
		ua := customUAs[i % len(customUAs)]
		// Make request with ua...
	}
}
```

## Legal Notice

User agent spoofing for testing purposes is legal only when:

1. You own the target system, or
2. You have explicit written authorization to test

See [LEGAL.md](LEGAL.md) for comprehensive legal guidance.

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Copyright**: © 2025 Muhammad Thariq