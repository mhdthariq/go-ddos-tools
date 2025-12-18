# Configuration File Documentation

## Overview

The `config.json` file is the central configuration file for DDoS-Tools. It contains settings for Minecraft-specific features, proxy providers, and dictionary providers. This configuration-based approach allows for easy updates and customization without code changes.

## File Location

```
ddos-tools/config.json
```

## Complete Structure

```json
{
  "MCBOT": "MHDDoS_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [
    {
      "type": 4,
      "url": "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/refs/heads/master/socks4.txt",
      "timeout": 5
    },
    {
      "type": 5,
      "url": "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/refs/heads/master/socks5.txt",
      "timeout": 5
    },
    {
      "type": 1,
      "url": "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
      "timeout": 5
    }
  ],
  "dictionary-providers": [
    {
      "name": "passwords-500",
      "type": "password",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/500-worst-passwords.txt",
      "filename": "files/passwords-500.txt"
    },
    {
      "name": "passwords-1k",
      "type": "password",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/best1050.txt",
      "filename": "files/passwords-1k.txt"
    },
    {
      "name": "passwords-10k",
      "type": "password",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10k-most-common.txt",
      "filename": "files/passwords-10k.txt"
    },
    {
      "name": "usernames",
      "type": "username",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Usernames/top-usernames-shortlist.txt",
      "filename": "files/usernames.txt"
    }
  ]
}
```

## Configuration Sections

### 1. Minecraft Settings

#### MCBOT
- **Type**: String
- **Default**: `"MHDDoS_"`
- **Purpose**: Bot name prefix for Minecraft attacks
- **Example**: `"MHDDoS_"`

#### MINECRAFT_DEFAULT_PROTOCOL
- **Type**: Integer
- **Default**: `47`
- **Purpose**: Default Minecraft protocol version
- **Example**: `47` (Minecraft 1.8.x)
- **Common Values**:
  - `47` - Minecraft 1.8.x
  - `210` - Minecraft 1.10.x
  - `335` - Minecraft 1.12.x
  - `404` - Minecraft 1.13.x
  - `498` - Minecraft 1.14.x

### 2. Proxy Providers

The `proxy-providers` array contains configurations for downloading proxy lists.

#### Proxy Provider Object

```json
{
  "type": 5,
  "url": "https://example.com/proxies.txt",
  "timeout": 5
}
```

#### Fields

**type** (integer)
- `1` - HTTP/HTTPS proxy
- `4` - SOCKS4 proxy
- `5` - SOCKS5 proxy

**url** (string)
- URL to download the proxy list from
- Must return a text file with one proxy per line
- Format: `ip:port`

**timeout** (integer)
- Connection timeout in seconds
- Recommended: `5` seconds

#### Default Providers

The default configuration includes three proxy providers:
1. SOCKS4 proxies from TheSpeedX/PROXY-List
2. SOCKS5 proxies from TheSpeedX/PROXY-List
3. HTTP proxies from TheSpeedX/PROXY-List

#### Adding Custom Proxy Providers

```json
{
  "proxy-providers": [
    {
      "type": 5,
      "url": "https://your-proxy-source.com/socks5.txt",
      "timeout": 10
    }
  ]
}
```

### 3. Dictionary Providers

The `dictionary-providers` array contains configurations for downloading password and username dictionaries.

#### Dictionary Provider Object

```json
{
  "name": "passwords-1k",
  "type": "password",
  "url": "https://example.com/passwords.txt",
  "filename": "files/passwords-1k.txt"
}
```

#### Fields

**name** (string)
- Unique identifier for the dictionary
- Used with `dict download <name>` command
- Examples: `"passwords-1k"`, `"usernames"`, `"custom-list"`

**type** (string)
- `"password"` - Password dictionary
- `"username"` - Username dictionary
- Used for categorization and documentation

**url** (string)
- URL to download the dictionary from
- Must return a text file with one entry per line
- SecLists URLs are recommended for security testing

**filename** (string)
- Local path where the dictionary will be saved
- Relative to project root
- Recommended: `"files/dictionary-name.txt"`

#### Default Providers

The default configuration includes four dictionary providers:
1. **passwords-500** - Top 500 worst passwords (fast)
2. **passwords-1k** - Top 1,000 best passwords (recommended)
3. **passwords-10k** - Top 10,000 common passwords (comprehensive)
4. **usernames** - Common usernames

All default dictionaries are sourced from [SecLists](https://github.com/danielmiessler/SecLists).

#### Adding Custom Dictionary Providers

Example: Adding a custom password list

```json
{
  "dictionary-providers": [
    {
      "name": "custom-passwords",
      "type": "password",
      "url": "https://example.com/my-passwords.txt",
      "filename": "files/custom-passwords.txt"
    },
    {
      "name": "company-usernames",
      "type": "username",
      "url": "https://internal.company.com/usernames.txt",
      "filename": "files/company-users.txt"
    }
  ]
}
```

Then download with:
```bash
./ddos-tools dict download custom-passwords
./ddos-tools dict download company-usernames
```

## Usage Examples

### Loading Configuration

The configuration file is loaded automatically when running commands:

```bash
# Config is loaded from config.json by default
./ddos-tools CFB https://example.com --proxy-file http.txt

# Specify custom config file
./ddos-tools CFB https://example.com --config custom-config.json
```

### Using Proxy Providers

```bash
# Download proxies from configured providers
./ddos-tools CFB https://example.com \
  --proxy-file http.txt \
  --threads 500
```

The tool will:
1. Check for `files/proxies/http.txt`
2. If not found, download from the HTTP proxy provider in config
3. Use the downloaded proxies for the attack

### Using Dictionary Providers

```bash
# List available dictionaries from config
./ddos-tools dict list

# Download a dictionary
./ddos-tools dict download passwords-10k

# Use in login attack
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-10k.txt \
  --threads 100
```

## Configuration Best Practices

### 1. Version Control

**DO**: Commit `config.json` to version control
```bash
git add config.json
git commit -m "Update dictionary providers"
```

**Benefits**:
- Team consistency
- Track changes over time
- Easy rollback if needed

### 2. Custom Sources

Create organization-specific configurations:

```json
{
  "dictionary-providers": [
    {
      "name": "company-breach-2024",
      "type": "password",
      "url": "https://internal.company.com/breach-passwords.txt",
      "filename": "files/company-breach.txt"
    }
  ]
}
```

### 3. Multiple Environments

Create different config files for different purposes:

- `config.json` - Default/production
- `config-dev.json` - Development with test sources
- `config-pentest.json` - Penetration testing specific

Use with:
```bash
./ddos-tools LOGIN https://target.com --config config-pentest.json
```

### 4. Backup URLs

Add multiple providers for redundancy:

```json
{
  "proxy-providers": [
    {
      "type": 5,
      "url": "https://primary-source.com/socks5.txt",
      "timeout": 5
    },
    {
      "type": 5,
      "url": "https://backup-source.com/socks5.txt",
      "timeout": 5
    }
  ]
}
```

### 5. Documentation Comments

While JSON doesn't support comments, document your config in a separate file:

```markdown
# Config Documentation

## Custom Dictionary: company-passwords
- Source: Internal password audit
- Updated: 2024-01-15
- Contact: security@company.com
```

## Validation

### Valid Configuration

```json
{
  "MCBOT": "TestBot_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [],
  "dictionary-providers": []
}
```

### Invalid Configurations

**Missing required fields:**
```json
{
  "MCBOT": "Test_"
  // Missing other required fields
}
```

**Invalid JSON syntax:**
```json
{
  "MCBOT": "Test_",  // Trailing comma on last item - INVALID
}
```

**Invalid proxy type:**
```json
{
  "proxy-providers": [
    {
      "type": 999,  // Invalid type - INVALID
      "url": "https://example.com/proxies.txt",
      "timeout": 5
    }
  ]
}
```

## Troubleshooting

### Config File Not Found

**Error**: `Could not load config.json, using defaults`

**Solution**:
```bash
# Ensure config.json exists in project root
ls -la config.json

# Or specify path
./ddos-tools CFB https://target.com --config /path/to/config.json
```

### Invalid JSON Format

**Error**: `Failed to load config: invalid character...`

**Solution**:
```bash
# Validate JSON syntax
cat config.json | jq .

# Or use online validator
# https://jsonlint.com/
```

### Dictionary Provider Not Found

**Error**: `dictionary provider 'xyz' not found in config`

**Solution**:
```bash
# Check available providers
./ddos-tools dict list

# Verify config.json contains the provider
cat config.json | jq '.["dictionary-providers"]'
```

### Download Failures

**Error**: `Failed to download dictionary: HTTP 404`

**Solution**:
1. Verify the URL is accessible:
   ```bash
   curl -I "URL_FROM_CONFIG"
   ```
2. Update the URL in config.json
3. No rebuild needed - changes take effect immediately

## Security Considerations

### 1. Sensitive URLs

If your config contains sensitive URLs:

```bash
# Use environment-specific configs
cp config.json config-local.json

# Add to .gitignore
echo "config-local.json" >> .gitignore
```

### 2. Proxy Privacy

Be cautious with proxy provider URLs:
- Use trusted sources only
- Consider hosting your own proxy lists
- Avoid untrusted public proxy sources

### 3. Dictionary Sources

For dictionary providers:
- ✅ **Recommended**: SecLists (verified, maintained)
- ✅ **Safe**: Internal company sources
- ⚠️ **Caution**: Random internet sources
- ❌ **Avoid**: Untrusted/suspicious URLs

## Migration Guide

### From Hardcoded URLs

**Before** (hardcoded in code):
```go
const PasswordURL = "https://example.com/passwords.txt"
```

**After** (config-based):
```json
{
  "dictionary-providers": [
    {
      "name": "custom",
      "type": "password",
      "url": "https://example.com/passwords.txt",
      "filename": "files/custom.txt"
    }
  ]
}
```

**Benefits**:
- No code changes needed for URL updates
- Easy to add new sources
- Version control friendly

### Adding New Dictionary Types

1. **Edit config.json**:
   ```json
   {
     "name": "my-new-dict",
     "type": "password",
     "url": "https://source.com/list.txt",
     "filename": "files/my-dict.txt"
   }
   ```

2. **Download**:
   ```bash
   ./ddos-tools dict download my-new-dict
   ```

3. **Use**:
   ```bash
   ./ddos-tools LOGIN https://target.com --passwords files/my-dict.txt
   ```

## Reference

### Config Struct (Go)

```go
type Config struct {
    MCBot               string
    MinecraftProtocol   int
    ProxyProviders      []ProxyProvider
    DictionaryProviders []DictionaryProvider
}

type ProxyProvider struct {
    Type    int
    URL     string
    Timeout int
}

type DictionaryProvider struct {
    Name     string
    Type     string
    URL      string
    Filename string
}
```

### Default Values

If config.json is missing or incomplete:

- `MCBOT`: `"MHDDoS_"`
- `MINECRAFT_DEFAULT_PROTOCOL`: `47`
- `proxy-providers`: `[]` (empty)
- `dictionary-providers`: `[]` (empty)

## Additional Resources

- [Main Documentation](../README.md)
- [Dictionary Attack Guide](DICTIONARY_ATTACK.md)
- [Proxy Configuration](../README.md#proxy-support)
- [SecLists Repository](https://github.com/danielmiessler/SecLists)

## Support

For configuration issues:
1. Validate JSON syntax with `jq`
2. Check file permissions
3. Verify URLs are accessible
4. Review error messages for specific issues

## Examples

### Minimal Config

```json
{
  "MCBOT": "Bot_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [],
  "dictionary-providers": []
}
```

### Production Config

```json
{
  "MCBOT": "ProdBot_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [
    {
      "type": 5,
      "url": "https://company-proxies.internal/socks5.txt",
      "timeout": 10
    }
  ],
  "dictionary-providers": [
    {
      "name": "company-passwords",
      "type": "password",
      "url": "https://security.company.internal/common-passwords.txt",
      "filename": "files/company-passwords.txt"
    },
    {
      "name": "company-users",
      "type": "username",
      "url": "https://security.company.internal/common-usernames.txt",
      "filename": "files/company-users.txt"
    }
  ]
}
```

### Testing Config

```json
{
  "MCBOT": "TestBot_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [
    {
      "type": 1,
      "url": "https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
      "timeout": 3
    }
  ],
  "dictionary-providers": [
    {
      "name": "test-passwords",
      "type": "password",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/500-worst-passwords.txt",
      "filename": "files/test-passwords.txt"
    }
  ]
}
```
