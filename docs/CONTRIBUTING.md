# Contributing to ddos-tools Documentation

Thank you for your interest in improving the ddos-tools documentation! This guide will help you contribute effectively to our documentation.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Documentation Standards](#documentation-standards)
- [Types of Contributions](#types-of-contributions)
- [Submission Process](#submission-process)
- [Style Guide](#style-guide)
- [Review Process](#review-process)

## 🚀 Getting Started

### Prerequisites

- Git installed on your system
- GitHub account
- Text editor or IDE with Markdown support
- Basic understanding of Markdown syntax

### Setting Up

1. **Fork the repository**
   ```bash
   # Click "Fork" on GitHub, then clone your fork
   git clone https://github.com/YOUR-USERNAME/ddos-tools.git
   cd ddos-tools
   ```

2. **Create a documentation branch**
   ```bash
   git checkout -b docs/your-improvement-name
   ```

3. **Make your changes**
   - Edit existing documentation in `docs/`
   - Add new documentation files as needed
   - Update the `docs/README.md` index if adding new files

4. **Preview your changes**
   - Use a Markdown previewer (VS Code, GitHub preview, etc.)
   - Check all links work correctly
   - Verify code examples are accurate

## 📝 Documentation Standards

### File Organization

```
docs/
├── README.md              # Documentation index (update when adding files)
├── USAGE.md               # User guide
├── LEGAL.md               # Legal information
├── LEGAL-QUICK-REF.md     # Legal quick reference
├── USER-AGENTS.md         # Technical documentation
├── CHANGELOG.md           # Version history
└── CONTRIBUTING.md        # This file
```

### Naming Conventions

- Use UPPERCASE for major documents (README.md, USAGE.md, etc.)
- Use kebab-case for technical docs (user-agents.md, configuration.md)
- Use descriptive names that clearly indicate content
- Add `.md` extension to all Markdown files

### File Structure

Each documentation file should include:

1. **Title** (H1 heading)
2. **Overview** - Brief description of document purpose
3. **Table of Contents** - For documents longer than 3 sections
4. **Main Content** - Organized with clear headings
5. **Examples** - Practical code examples where applicable
6. **Footer** - Last updated date and maintainer info

## 🎯 Types of Contributions

### 1. Fixing Errors

**Examples:**
- Typos and grammar mistakes
- Incorrect code examples
- Broken links
- Outdated information

**Process:**
1. Identify the error
2. Make the correction
3. Verify the fix is accurate
4. Submit a pull request

### 2. Improving Clarity

**Examples:**
- Rephrasing confusing sections
- Adding missing explanations
- Breaking up long paragraphs
- Improving organization

**Process:**
1. Identify unclear content
2. Rewrite for better clarity
3. Ensure technical accuracy
4. Submit a pull request

### 3. Adding New Content

**Examples:**
- New tutorials or guides
- Additional examples
- FAQ sections
- Troubleshooting guides

**Process:**
1. Check if content already exists
2. Discuss in an issue first (for major additions)
3. Write the new content
4. Update the documentation index
5. Submit a pull request

### 4. Updating Examples

**Examples:**
- Modernizing code examples
- Adding new use cases
- Fixing deprecated methods
- Improving example clarity

**Process:**
1. Test the new examples
2. Ensure compatibility with current version
3. Update relevant documentation
4. Submit a pull request

## 📋 Submission Process

### Step-by-Step Guide

1. **Create an Issue (Optional but Recommended)**
   - For major changes, create an issue first
   - Describe what you want to improve
   - Get feedback from maintainers

2. **Make Your Changes**
   ```bash
   # Create a branch
   git checkout -b docs/descriptive-name
   
   # Edit files
   nano docs/USAGE.md
   
   # Stage changes
   git add docs/USAGE.md
   
   # Commit with descriptive message
   git commit -m "docs: improve USAGE.md proxy configuration examples"
   ```

3. **Test Your Changes**
   - Preview Markdown rendering
   - Test all code examples
   - Check all links work
   - Verify formatting is correct

4. **Push to Your Fork**
   ```bash
   git push origin docs/descriptive-name
   ```

5. **Create Pull Request**
   - Go to GitHub and create a pull request
   - Use a descriptive title (e.g., "docs: add troubleshooting section to USER-AGENTS.md")
   - Describe your changes in the PR description
   - Reference any related issues

### Commit Message Format

Use conventional commits format:

```
docs: <short description>

<optional longer description>

<optional footer>
```

**Examples:**
```
docs: fix typos in USAGE.md

docs: add proxy rotation examples to USAGE.md

docs: update USER-AGENTS.md with mobile examples

Added examples for iOS and Android user agents to improve
mobile testing coverage.

Closes #123
```

**Prefixes:**
- `docs:` - Documentation changes
- `fix:` - Fixing errors in docs
- `feat:` - Adding new documentation
- `refactor:` - Restructuring existing docs

## ✍️ Style Guide

### Markdown Formatting

#### Headings

```markdown
# H1 - Document Title (once per file)

## H2 - Major Sections

### H3 - Subsections

#### H4 - Minor Subsections
```

#### Code Blocks

Always specify the language:

````markdown
```bash
./ddos-tools GET http://example.com
```

```go
package main

func main() {
    // Go code here
}
```

```json
{
    "config": "value"
}
```
````

#### Lists

**Unordered:**
```markdown
- First item
- Second item
  - Nested item
  - Another nested item
- Third item
```

**Ordered:**
```markdown
1. First step
2. Second step
3. Third step
```

#### Links

```markdown
[Link Text](URL)
[Internal Link](../README.md)
[Section Link](#section-name)
```

#### Emphasis

```markdown
**Bold** for important terms
*Italic* for emphasis
`Code` for inline code, commands, or filenames
```

#### Tables

```markdown
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
| Data 4   | Data 5   | Data 6   |
```

### Writing Style

#### Be Clear and Concise

✅ **Good:**
```markdown
To start a Layer 7 attack, use the GET method:
```

❌ **Avoid:**
```markdown
If you want to perhaps initiate what we call a Layer 7 attack, you might want to consider using the GET method:
```

#### Use Active Voice

✅ **Good:**
```markdown
Run the following command to test connectivity:
```

❌ **Avoid:**
```markdown
The following command should be run to test connectivity:
```

#### Be Specific

✅ **Good:**
```markdown
Set the thread count to 100 for optimal performance on systems with 8+ CPU cores.
```

❌ **Avoid:**
```markdown
Set the thread count to something reasonable based on your system.
```

#### Include Examples

Always provide practical examples:

```markdown
### Basic Usage

To perform a GET request:

```bash
./ddos-tools GET http://example.com 5 100 proxies.txt 100 60
```

This command:
- Uses the GET method
- Targets http://example.com
- Uses proxy type 5 (SOCKS5)
- Runs 100 threads
- Loads proxies from proxies.txt
- Sends 100 requests per connection
- Runs for 60 seconds
```

### Legal and Safety Notices

Always include appropriate warnings:

```markdown
⚠️ **WARNING**: Only test systems you own or have explicit written authorization to test.

⚠️ **IMPORTANT**: This feature requires root/administrator privileges.

ℹ️ **NOTE**: Performance may vary based on network conditions.

✅ **TIP**: Use proxy rotation to avoid rate limiting.
```

## 🔍 Review Process

### What Reviewers Look For

1. **Technical Accuracy**
   - Code examples work correctly
   - Information is up-to-date
   - No misleading statements

2. **Clarity**
   - Easy to understand
   - Well-organized
   - Appropriate level of detail

3. **Formatting**
   - Proper Markdown syntax
   - Consistent style
   - Working links

4. **Completeness**
   - All necessary information included
   - Examples provided where helpful
   - Related docs updated if needed

### Review Timeline

- Simple fixes: 1-3 days
- Medium changes: 3-7 days
- Major additions: 1-2 weeks

### Addressing Feedback

1. **Read all feedback carefully**
2. **Ask questions if unclear**
3. **Make requested changes**
4. **Push updates to the same PR**
5. **Respond to comments**

## 🏆 Best Practices

### Documentation Checklist

Before submitting, verify:

- [ ] All code examples are tested and working
- [ ] All links are valid and point to correct locations
- [ ] Spelling and grammar are correct
- [ ] Formatting is consistent with existing docs
- [ ] New files are added to docs/README.md index
- [ ] Legal warnings are included where appropriate
- [ ] Examples include clear explanations
- [ ] Technical accuracy is verified
- [ ] Commit messages follow conventional format

### Common Mistakes to Avoid

❌ **Don't:**
- Copy-paste from external sources without attribution
- Include outdated version information
- Use ambiguous pronouns (it, this, that) without clear references
- Assume reader's technical level
- Include untested code examples
- Forget to update related documentation

✅ **Do:**
- Test all examples before submitting
- Link to related documentation
- Use clear, descriptive headings
- Include practical use cases
- Update the changelog when appropriate
- Add version information when relevant

## 📞 Getting Help

### Resources

- **Markdown Guide**: https://www.markdownguide.org/
- **GitHub Markdown**: https://guides.github.com/features/mastering-markdown/
- **Conventional Commits**: https://www.conventionalcommits.org/

### Contact

- **Questions**: Open a GitHub Discussion
- **Issues**: Create a GitHub Issue with `documentation` label
- **Quick Help**: Comment on your pull request

## 🙏 Recognition

All contributors are valued! Your contributions will be:

- Acknowledged in commit history
- Appreciated by the community
- Helping others use the tool effectively

Thank you for helping improve ddos-tools documentation! 🎉

---

**Maintained By**: Muhammad Thariq  
**Last Updated**: November 2025  
**Version**: 2.4 SNAPSHOT  
**Copyright**: © 2025 Muhammad Thariq