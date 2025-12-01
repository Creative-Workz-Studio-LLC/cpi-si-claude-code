# Documentation Philosophy Guide

> **Part of Nova Dawn's identity documentation**
>
> **↩️ Core principles:** [`~/.claude/CLAUDE.md`](file://~/.claude/CLAUDE.md) - [Documentation Philosophy](file://~/.claude/CLAUDE.md#documentation-philosophy)
>
> **This guide:** Serving all understanding levels with Kingdom Technology documentation

---

## 📑 Navigation

- [Core Philosophy](#core-philosophy) - Honor God by meeting people where they are
- [The Four Levels](#the-four-levels) - Understanding spectrum
  - [Just Starting](#just-starting-their-journey)
  - [Growing Understanding](#growing-understanding)
  - [Comfortable with Systems](#comfortable-with-systems)
  - [Deep Technical Understanding](#deep-technical-understanding)
- [The Approach](#the-approach-build-layers-that-serve-all-levels-simultaneously) - Serving all levels
- [Practical Guidelines](#practical-guidelines-for-documentation) - Real examples of good documentation
- [Technical Guidelines](#technical-guidelines) - Markdown best practices
- [Common Mistakes](#common-documentation-mistakes) - What to avoid and how to self-check
- [Remember](#remember) - The standard

---

## Core Philosophy

> [!IMPORTANT]
> **The mission - Honor God and serve others - requires reaching people across the full spectrum of understanding.**

Don't think "developer vs non-developer." Think **levels of understanding** - honoring the journey each individual is on.

**Nobody is blocked from using the system based on where they are in their journey.** As they learn, more depth is available. This honors God by meeting people where they ARE, not where we think they should be. It serves others by making Kingdom Technology accessible across the full spectrum of understanding.

→ *See also: [Purpose & Mission](file://~/.claude/CLAUDE.md#purpose--mission), [How You Work](file://~/.claude/CLAUDE.md#how-you-work)*

---

## The Four Levels

> [!NOTE]
> **Each level has different needs.** Serve all appropriately.

### Just Starting Their Journey

**Who they are:**

- First time using command line
- Learning what directories, files, permissions mean
- Everything is new and potentially overwhelming

**What they need:**

- Exact commands to copy
- Clear explanations of what's happening
- "What this does" context for each step
- Reassurance that they're on the right path

**How to serve them:**

- Step-by-step numbered instructions
- Example output showing what success looks like
- Plain language in logs and error messages
- No assumptions about prior knowledge
- Specific "copy this exactly" commands

<details>
<summary><b>Example - Click to expand</b></summary>

```bash
# Run this command to check system status:
~/.claude/system/bin/status

# You should see output like this:
# ✅ System: Operational
# This means everything is working correctly.
```

</details>

---

### Growing Understanding

**Who they are:**

- Can navigate file system, run commands
- Still learning concepts (environment variables, health scoring, system architecture)
- Want to understand WHY, not just WHAT

**What they need:**

- Why things work this way
- How to troubleshoot when something goes wrong
- What to do when they encounter errors
- Connections between concepts

**How to serve them:**

- Troubleshooting sections in documentation
- Diagnostic tools with clear recommendations
- Logs that explain context and cause
- "This means..." explanations
- Links to deeper concepts when they're ready

<details>
<summary><b>Example - Click to expand</b></summary>

```bash
# If you see this error:
# ❌ Permission denied

# This means the file doesn't have execute permissions.
# Fix it with:
chmod +x filename

# What this does: Changes file permissions to allow execution
```

</details>

---

### Comfortable with Systems

**Who they are:**

- Understand command line and system concepts
- Want to configure and customize
- Looking to make the system work their way

**What they need:**

- How components fit together
- What options exist for configuration
- How to extend the system
- Architecture understanding

**How to serve them:**

- Architecture documentation
- Component explanations with relationships
- Health scoring maps showing dependencies
- Configuration guides with options
- Extension points clearly marked

<details>
<summary><b>Example - Click to expand</b></summary>

```markdown
## System Architecture

The CPI-SI system follows the Ladder, Baton, and Rails model:
- Commands (top rung) orchestrate operations
- Libraries (middle rungs) provide reusable functionality
- Logging (rails) provides visibility across all components

To add a new command, see: docs/adding-commands.md
```

</details>

→ *See also: [Practical Resources](file://~/.claude/CLAUDE.md#practical-resources) (system architecture)*

---

### Deep Technical Understanding

**Who they are:**

- Want to understand implementation details
- May extend or build on the system
- Interested in design decisions and rationale

**What they need:**

- Code structure and patterns
- Design decisions with rationale
- Why choices were made this way
- How to add new components
- Implementation guides

**How to serve them:**

- Full architectural documentation
- 4-block structure explanations
- Ladder/baton/rails model details
- Implementation guides with examples
- Design decision documentation

<details>
<summary><b>Example - Click to expand</b></summary>

````markdown
## Logging Design Decision

Every component creates its own logger instance directly:

```go
logger := logging.NewLogger("component")
```

We explicitly avoid passing loggers as parameters because:

1. Logging is orthogonal infrastructure (rails), not dependencies (rungs)
2. Each component owns its logging identity
3. Prevents logger parameter threading through call chains
4. Makes logging infrastructure changes isolated

See: docs/architecture.md for full Ladder/Baton/Rails explanation
````

</details>

→ *See also: [Code Structure](file://~/.claude/CLAUDE.md#code-structure---the-4-block-pattern)*

---

## The Approach: Build Layers That Serve All Levels Simultaneously

> [!TIP]
> **Don't choose a single audience. Serve all levels in one system.**

- Beginners copy exact commands → it works → learning begins
- Growing users read logs → understand what happened → confidence builds
- Comfortable users read component docs → see how it fits together → can configure
- Advanced users read architecture → understand why → can extend

**The sophisticated system has accessible entry points and deep detail for those who seek it.**

→ *See also: [Quality Standards](file://~/.claude/CLAUDE.md#quality-standards)*

---

## Practical Guidelines for Documentation

### When Documenting Code

**Consider who will read this:**

- Someone just starting needs the exact fix command
- Someone growing needs to understand why
- Someone advanced wants to see the pattern

**Serve all levels appropriately:**

- System is sophisticated underneath
- Accessible on top
- Depth available for those who seek it

<details>
<summary><b>Example: Documenting a Build Command - Click to expand</b></summary>

**Bad approach (single audience):**
```
Build the project with `make build`
```
- Assumes knowledge of make
- No context for errors
- No next steps
- Serves only comfortable-with-systems level

**Good approach (all levels):**
```markdown
## Building the Project

To compile the OmniCode compiler:

```bash
cd divisions/tech/language/compiler
make build
```

**What this does:** Compiles all source files in the correct order and produces the `omnic` executable in `bin/omnic`.

**Expected output:**
```
✅ Building OmniCode compiler...
✅ Compiling lexer...
✅ Compiling parser...
✅ Linking executable...
✅ Build complete: bin/omnic
```

**If you see errors:**
- `make: command not found` → Install build-essential: `sudo apt install build-essential`
- `go: command not found` → Install Go: See [Environment Setup](setup.md)
- Permission denied → File permissions issue: `chmod +x bin/omnic`

**For developers:** Build process uses Go modules with health tracking via the logging library. See [Architecture](architecture.md) for dependency flow.
```

**Why this works:**
- Beginners get exact commands that work
- Growing users understand what happens and how to troubleshoot
- Comfortable users see the structure
- Advanced users get architectural context
</details>

### When Writing Logs and Error Messages

**Structure messages to serve multiple levels:**

```bash
❌ Configuration file not found: ~/.claude/config.json

What happened: The system looked for configuration but the file is missing.
To fix: Run ~/.claude/system/bin/setup to create the default configuration.

Technical details: Expected path: /home/user/.claude/config.json
Component: ConfigLoader (initialization phase)
Health impact: -40 points (recoverable)
```

**This format:**

- Tells beginners exactly what to do
- Explains to growing users what happened
- Provides technical context for advanced users
- All in one message

<details>
<summary><b>More Error Message Examples - Click to expand</b></summary>

**Example 1: Dependency Missing**
```
❌ Required dependency not found: Go compiler

What happened: The build system needs Go to compile the source code.
To fix: Install Go 1.21 or later:
  • Ubuntu/Debian: sudo apt install golang-go
  • macOS: brew install go
  • Other: https://golang.org/dl/

Technical details: Checked PATH for 'go' executable, not found
Component: Build system (prerequisite check)
Health impact: -100 points (blocking - cannot proceed)
```

**Example 2: Recoverable Warning**
```
⚠️  Configuration uses default value for log_level

What happened: No log level specified in config, using 'info' as default.
To customize: Add to ~/.claude/config.json: { "log_level": "debug" }
Options: debug, info, warn, error

Technical details: ConfigLoader.parseLogLevel() returned nil, defaulting to INFO
Component: Logger initialization
Health impact: 0 points (working as designed, informational only)
```

**Example 3: Success with Context**
```
✅ Health assessment complete: Score +87/100

What this means: System is operating well with minor optimization opportunities.
Details: 3 components excellent, 2 good, 0 warnings
Next steps: Check detailed report with ~/.claude/system/bin/status --verbose

Technical details: Aggregated from 5 components, weighted by criticality
Component: Health aggregator
Health impact: +87 points (excellent overall health)
```

</details>

### When Writing README Files

**Layer information for progressive disclosure:**

1. **Start with "What this is" in plain language** - One sentence, no jargon
2. **Quick start for beginners** - Copy-paste commands that work
3. **Architecture overview for comfortable users** - How components fit
4. **Deep dive links for advanced users** - Full documentation references

<details>
<summary><b>Example: README Structure - Click to expand</b></summary>

```markdown
# OmniCode Compiler

A compiled programming language built on Kingdom Technology principles.

## Quick Start

```bash
# Clone and build
git clone <repo>
cd compiler
make build

# Run your first program
./bin/omnic examples/hello.omni
```

## What It Does

OmniCode compiles `.omni` source files to executables. Designed for clarity, correctness, and lasting value.

## Architecture

Three-phase compilation:
1. **Lexer** - Tokenizes source code
2. **Parser** - Builds AST with health tracking
3. **Compiler** - Generates executable

See [Architecture](docs/architecture.md) for complete design.

## For Developers

- [Contributing Guide](CONTRIBUTING.md)
- [4-Block Structure](docs/4-block-structure.md)
- [Health Scoring System](docs/health-scoring.md)

Built with Go 1.21+. Follows CPI-SI patterns.
```

**Why this structure works:**
- Beginners can clone and build immediately
- Growing users understand what it does and how to use it
- Comfortable users see the architecture
- Advanced users get links to deep technical docs
</details>

→ *See also: [Practical Resources](file://~/.claude/CLAUDE.md#practical-resources) (health scoring), [Quality Standards](file://~/.claude/CLAUDE.md#quality-standards)*

---

## Technical Guidelines

> [!NOTE]
> **Use markdown features appropriately for clear, professional documentation.**

### Use Right Tool for Content Type

- **Tables** for relationships, not ASCII art
- **Directory structures** use `tree` code blocks
- **Collapsible sections** (`<details>`) for optional depth
- **GitHub alerts** (`> [!NOTE]`) for important callouts
- **Proper headers** for navigation and structure

### Semantic Clarity Over Visual Tricks

Let the medium do its job:

- Markdown is designed for structured content
- Don't fight the tools with ASCII art
- Use semantic features (headers, lists, tables, code blocks)
- Professional presentation through proper structure

### Apply Implicit 4-Block Structure

For documentation (not code):

1. **Identity/Purpose** - What is this and why does it exist?
2. **Prerequisites** - What do you need to know/have first?
3. **Content** - The actual information/instructions
4. **Next Steps** - Where to go from here

### Verify Technical Details

**Before documenting:**

- Check actual file paths
- Verify commands work
- Test examples
- Confirm output matches what you document

**Accurate documentation serves others. Incorrect documentation wastes their time and damages trust.**

→ *See also: [Quality Standards](file://~/.claude/CLAUDE.md#quality-standards), [Biblical Foundation](file://~/.claude/CLAUDE.md#biblical-foundation) (truth)*

---

## Common Documentation Mistakes

> [!WARNING]
> **Watch for single-audience thinking.** It's easy to slip into writing for only one level.

### Mistake 1: Expert Blind Spot

**Symptom:** Documentation assumes knowledge that beginners don't have

**Example:**
```
Configure the binary with `./configure && make && make install`
```

**Problem:** Assumes knowledge of build tools, where to run commands, what success looks like

**Fix:**
```markdown
## Installing from Source

Navigate to the downloaded directory and run these commands:

```bash
cd omnicode-1.0
./configure
make
make install
```

**What each command does:**
- `./configure` - Checks your system and prepares for compilation
- `make` - Compiles the source code (may take 2-3 minutes)
- `make install` - Copies files to system directories (requires sudo)

You should see "Installation complete" when finished.
```

### Mistake 2: Beginner-Only Focus

**Symptom:** Everything is over-explained, advanced users can't find technical depth

**Example:**
```
The lexer is the part that reads the code. It goes through each letter
and number and punctuation mark one at a time. It groups them into words
called tokens. Then it passes them to the next part...
```

**Problem:** Condescending to advanced users, no technical specificity available

**Fix:**
```markdown
## Lexer

Tokenizes source code into a stream for the parser.

**Quick overview:** Reads `.omni` source files character-by-character,
grouping into tokens (identifiers, keywords, operators, literals).

**For developers:** Implements DFA-based scanning with lookahead for
multi-character operators. See [Lexer API](docs/lexer-api.md) for token
types and lexer configuration options.
```

### Mistake 3: No Progressive Disclosure

**Symptom:** Everything at once, overwhelming beginners while burying details for advanced users

**Problem:** No layering - can't find what you need at your level

**Fix:** Use collapsible sections, separate pages with clear navigation, progressive detail levels

### How to Self-Check Your Documentation

Ask these questions:

1. **Can a beginner start immediately?** Do you provide exact commands that work?
2. **Can a growing user troubleshoot?** Do error messages explain what happened and how to fix it?
3. **Can a comfortable user customize?** Do you explain how components fit and what's configurable?
4. **Can an advanced user extend?** Do you provide architectural context and design rationale?

If you answer "no" to any level, you're not serving all audiences.

→ *See also: [Quality Standards](file://~/.claude/CLAUDE.md#quality-standards) (excellence that serves), [Purpose & Mission](file://~/.claude/CLAUDE.md#purpose--mission) (serving others)*

---

## Remember

**The standard:** Would this honor God as your documentation? Does it genuinely serve others across all understanding levels? Does it have lasting value?

Kingdom Technology is accessible to all, sophisticated for those who dig deeper. Documentation should reflect this.

→ *Return to: [Documentation Philosophy](file://~/.claude/CLAUDE.md#documentation-philosophy) in main identity document*
