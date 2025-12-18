# Config-Based Dictionary Implementation Summary

## Overview

The dictionary attack feature now uses a **config-based architecture** for managing dictionary sources, similar to how proxy providers work. All dictionary URLs are centralized in `config.json`, making the system more maintainable, flexible, and consistent with the project's existing architecture.

## What Changed

### Before (Hardcoded URLs)
```go
// URLs hardcoded in pkg/utils/dictionary.go
const (
    Top500PasswordsURL = "https://raw.githubusercontent.com/..."
    Best1kPasswordsURL = "https://raw.githubusercontent.com/..."
    CommonPasswordsURL = "https://raw.githubusercontent.com/..."
)
```

### After (Config-Based)
```json
// Configured in config.json
{
  "dictionary-providers": [
    {
      "name": "passwords-500",
      "type": "password",
      "url": "https://raw.githubusercontent.com/danielmiessler/SecLists/master/...",
      "filename": "files/passwords-500.txt"
    }
  ]
}
```

## Architecture Benefits

### ✅ Centralized Management
- All dictionary URLs in one place (`config.json`)
- Easy to view and modify all sources
- Consistent with proxy provider architecture

### ✅ No Code Changes Required
- Update URLs without recompiling
- Add new dictionaries without touching code
- Instant changes - just edit config.json

### ✅ Team Collaboration
- Version control friendly
- Everyone uses same sources
- Easy to share custom dictionaries

### ✅ Extensibility
- Add unlimited custom dictionaries
- Support multiple sources for same type
- Easy to maintain and update

## Configuration Structure

### config.json Layout

```json
{
  "MCBOT": "MHDDoS_",
  "MINECRAFT_DEFAULT_PROTOCOL": 47,
  "proxy-providers": [...],
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

### Dictionary Provider Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `name` | string | Unique identifier | `"passwords-1k"` |
| `type` | string | Dictionary category | `"password"` or `"username"` |
| `url` | string | Download source | `"https://..."` |
| `filename` | string | Local save path | `"files/passwords-1k.txt"` |

## Implementation Details

### New Functions (pkg/utils/dictionary.go)

```go
// Get dictionary provider by name from config
GetDictionaryProvider(cfg *config.Config, name string) *config.DictionaryProvider

// Download dictionary using config provider
DownloadDictionaryFromConfig(cfg *config.Config, dictName string) error

// Ensure dictionary exists, download from config if not
EnsureDictionaryFromConfig(cfg *config.Config, dictName string) error

// Get all configured dictionary providers
GetAllDictionaryProviders(cfg *config.Config) []config.DictionaryProvider
```

### Config Structure (pkg/config/config.go)

```go
type Config struct {
    MCBot               string
    MinecraftProtocol   int
    ProxyProviders      []ProxyProvider
    DictionaryProviders []DictionaryProvider  // NEW
}

type DictionaryProvider struct {
    Name     string `json:"name"`
    Type     string `json:"type"`
    URL      string `json:"url"`
    Filename string `json:"filename"`
}
```

### Command Integration (pkg/cmd/run.go)

The `dict` command now:
1. Loads `config.json` automatically
2. Reads available providers from config
3. Downloads using config URLs
4. Lists all configured dictionaries

## Usage Examples

### View Available Dictionaries

```bash
./ddos-tools dict list
```

Output:
```
[INFO] Available dictionaries:

[OK] ✓ passwords-500 (499 entries) - files/passwords-500.txt
[OK] ✓ passwords-1k (1049 entries) - files/passwords-1k.txt
[OK] ✓ passwords-10k (10000 entries) - files/passwords-10k.txt
[OK] ✓ usernames (17 entries) - files/usernames.txt
```

### Download Dictionary

```bash
./ddos-tools dict download passwords-1k
```

The tool:
1. Reads `config.json`
2. Finds the `passwords-1k` provider
3. Downloads from the configured URL
4. Saves to the configured filename

### Use in Attack

```bash
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-10k.txt \
  --threads 100
```

With auto-download:
```bash
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-10k.txt \
  --download-dict \
  --threads 100
```

## Adding Custom Dictionaries

### Step 1: Edit config.json

Add your custom dictionary provider:

```json
{
  "dictionary-providers": [
    {
      "name": "custom-passwords",
      "type": "password",
      "url": "https://mycompany.com/password-list.txt",
      "filename": "files/custom-passwords.txt"
    },
    {
      "name": "company-users",
      "type": "username",
      "url": "https://internal.company.com/usernames.txt",
      "filename": "files/company-users.txt"
    }
  ]
}
```

### Step 2: Download

```bash
./ddos-tools dict download custom-passwords
./ddos-tools dict download company-users
```

### Step 3: Use

```bash
./ddos-tools LOGIN https://target.com/login \
  --usernames files/company-users.txt \
  --passwords files/custom-passwords.txt \
  --threads 200
```

## Comparison with Proxy Providers

The dictionary provider implementation follows the **exact same pattern** as proxy providers:

### Proxy Providers
```json
{
  "proxy-providers": [
    {
      "type": 5,
      "url": "https://...",
      "timeout": 5
    }
  ]
}
```

### Dictionary Providers
```json
{
  "dictionary-providers": [
    {
      "name": "passwords-1k",
      "type": "password",
      "url": "https://...",
      "filename": "files/..."
    }
  ]
}
```

**Benefits of Consistency**:
- Same mental model for users
- Familiar pattern for developers
- Easier to maintain and extend
- Professional, polished architecture

## Migration Path

### For Users
**No migration needed!** The feature is backward compatible:
- Existing credential files still work
- `--data` flag unchanged
- No breaking changes

### For Developers
**Easy to extend:**
1. Edit `config.json` (no code changes)
2. Add new dictionary provider
3. Changes take effect immediately

## Default Configuration

The default `config.json` includes four dictionary providers from **SecLists**:

1. **passwords-500** (500 worst passwords)
   - Fast testing
   - 3.5 KB
   - 499 entries

2. **passwords-1k** (best 1,050 passwords)
   - Recommended for most use cases
   - 7.8 KB
   - 1,049 entries

3. **passwords-10k** (10,000 common passwords)
   - Comprehensive testing
   - 72 KB
   - 10,000 entries

4. **usernames** (common usernames)
   - System and service accounts
   - 112 bytes
   - 17 entries

All sourced from: https://github.com/danielmiessler/SecLists

## Error Handling

### Config Not Found
```
[ERROR] Failed to load config: open config.json: no such file or directory
[WARN] Cannot access dictionary providers without config.json
```

**Solution**: Ensure `config.json` exists in project root

### Provider Not Found
```
[ERROR] Failed to download dictionary: dictionary provider 'xyz' not found in config
[INFO] Available types:
  - passwords-500
  - passwords-1k
  - passwords-10k
  - usernames
```

**Solution**: Use a configured provider name or add it to config.json

### Download Failed
```
[ERROR] Failed to download dictionary: failed to download dictionary: HTTP 404
```

**Solution**: Update the URL in `config.json`

## Best Practices

### 1. Version Control
Commit `config.json` to track dictionary sources:
```bash
git add config.json
git commit -m "Add custom dictionary providers"
```

### 2. Documentation
Document custom dictionaries in your config:
```markdown
## Custom Dictionary: company-passwords
- Source: Internal password audit 2024
- Last Updated: 2024-01-15
- Maintained by: security-team@company.com
```

### 3. Environment-Specific Configs
```bash
# Development
cp config.json config-dev.json

# Production
cp config.json config-prod.json

# Use specific config
./ddos-tools LOGIN https://target.com --config config-prod.json
```

### 4. Backup URLs
Add multiple sources for redundancy:
```json
{
  "dictionary-providers": [
    {
      "name": "passwords-primary",
      "type": "password",
      "url": "https://primary-source.com/passwords.txt",
      "filename": "files/passwords.txt"
    },
    {
      "name": "passwords-backup",
      "type": "password",
      "url": "https://backup-source.com/passwords.txt",
      "filename": "files/passwords-backup.txt"
    }
  ]
}
```

## Security Considerations

### Trusted Sources Only
- ✅ SecLists (verified, community-maintained)
- ✅ Internal company sources
- ⚠️ Carefully vet third-party sources
- ❌ Avoid random internet URLs

### Config File Privacy
If using sensitive/internal URLs:
```bash
# Create local config
cp config.json config-local.json

# Add to .gitignore
echo "config-local.json" >> .gitignore
```

### HTTPS Only
Always use HTTPS URLs in config.json:
```json
// Good
"url": "https://secure-source.com/list.txt"

// Bad
"url": "http://insecure-source.com/list.txt"
```

## Testing

### Validate Config
```bash
# Check JSON syntax
cat config.json | jq .

# Check dictionary providers
cat config.json | jq '.["dictionary-providers"]'
```

### Test Download
```bash
# Test each provider
./ddos-tools dict download passwords-500
./ddos-tools dict download passwords-1k
./ddos-tools dict download passwords-10k
./ddos-tools dict download usernames

# Verify downloads
./ddos-tools dict list
```

### Test Usage
```bash
# Quick test with smallest dictionary
./ddos-tools LOGIN https://httpbin.org/post \
  --passwords files/passwords-500.txt \
  --threads 10 \
  --duration 10
```

## Documentation Files

- **[CONFIG.md](docs/CONFIG.md)** - Complete config.json documentation
- **[DICTIONARY_ATTACK.md](docs/DICTIONARY_ATTACK.md)** - Dictionary attack guide
- **[DICTIONARY_QUICKSTART.md](DICTIONARY_QUICKSTART.md)** - Quick start guide
- **[DICTIONARY_IMPLEMENTATION.md](DICTIONARY_IMPLEMENTATION.md)** - Implementation details

## Summary

The config-based dictionary implementation provides:

✅ **Maintainability** - URLs in config, not code
✅ **Flexibility** - Add custom sources easily
✅ **Consistency** - Matches proxy provider pattern
✅ **Simplicity** - No code changes for new dictionaries
✅ **Professional** - Clean, organized architecture
✅ **Extensible** - Unlimited custom providers
✅ **Team-Friendly** - Version control and sharing

This architectural improvement makes DDoS-Tools more maintainable, professional, and easier to customize while maintaining backward compatibility with existing functionality.

## Quick Reference

```bash
# View available dictionaries (from config.json)
./ddos-tools dict list

# Download dictionary (reads config.json)
./ddos-tools dict download passwords-1k

# Use in attack
./ddos-tools LOGIN https://target.com/login \
  --passwords files/passwords-1k.txt \
  --threads 100

# Add custom dictionary to config.json, then:
./ddos-tools dict download my-custom-dict
```

**Config Location**: `ddos-tools/config.json`
**Dictionary Section**: `dictionary-providers` array
**No Rebuild Required**: Edit config.json and run immediately