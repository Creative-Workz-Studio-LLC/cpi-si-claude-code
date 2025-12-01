<div align="center">

# 🌐 Environment Configuration

**Non-interactive defaults for seamless terminal operations**

![Bash](https://img.shields.io/badge/Bash-4EAA25?style=flat&logo=gnu-bash&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)

*Eliminate interactive prompts across all development tools*

[Integration](#integration) • [Variables](#configuration-categories) • [Testing](#testing)

</div>

---

## Table of Contents

- [🌐 Environment Configuration](#-environment-configuration)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Files](#files)
  - [Integration](#integration)
    - [Automatic Integration](#automatic-integration)
    - [Manual Integration](#manual-integration)
  - [Configuration Categories](#configuration-categories)
    - [Package Management](#package-management)
    - [Development Tools](#development-tools)
    - [CPI-SI Framework](#cpi-si-framework)
  - [Testing](#testing)
  - [Usage](#usage)
  - [Integration with settings.json](#integration-with-settingsjson)
  - [Troubleshooting](#troubleshooting)
    - [Environment Variables Not Applied](#environment-variables-not-applied)
    - [Still Seeing Prompts](#still-seeing-prompts)

---

## Overview

> [!IMPORTANT]
> **Architecture Change:** This folder is transitioning to a generated configuration model.
>
> - Environment defaults now defined in `../system.toml` and `../user.toml`
> - `non-interactive.conf` will be generated into `../runtime/` from composed configs
> - This README documents the environment configuration system
>
> Current state: Transition in progress

**What is this?** Settings that make command-line tools stop asking questions and just do their work.

**Why it matters:**

- No more "Do you want to continue? [Y/n]" prompts
- No waiting for confirmations during package installs
- No interruptions when services restart
- Everything runs smoothly in automated workflows

**What you get:**

| Tool | Without this | With this |
|------|--------------|-----------|
| apt/dpkg | Asks "Continue? [Y/n]" | Just installs ✓ |
| pip | Prompts for confirmations | Silent install ✓ |
| npm | Asks "Is this OK? (yes)" | Auto-yes ✓ |
| System services | Interrupts to ask about restarts | Restarts automatically ✓ |

> [!NOTE]
> These settings only affect prompts - they don't force dangerous operations. Your system stays safe.

---

## Current Architecture

**Environment configuration is now generated from composed configs:**

| Source | Purpose | Location |
|--------|---------|----------|
| System defaults | Universal tool settings | `../system.toml` → `[environment.*]` |
| User preferences | Personal customizations | `../user.toml` → `[environment.*]` |
| Instance identity | CPI-SI workspace paths | `~/.claude/cpi-si/config/instance/*/config.json` |
| **Generated file** | **Composed environment** | **`../runtime/environment.conf`** |

## Deprecated Files

| File | Status | Migration |
|------|--------|-----------|
| `non-interactive.conf.deprecated` | Archived | Values moved to system.toml and user.toml |

**Migration complete:** All values extracted to appropriate config layers.

---

## Integration

### Automatic Integration

**Run the integration script (recommended):**

```bash
# Future: when configurator exists
~/.claude/cpi-si/system/bin/configure --integrate-environment
source ~/.bashrc
```

**What it does:**

1. ✓ Regenerates environment.conf from configs
2. ✓ Adds source line to `.bashrc` automatically
3. ✓ Validates composed configuration
4. ✓ Provides verification instructions
5. ✓ Safe to run multiple times (idempotent)

### Manual Integration

**Add to `~/.bashrc` manually:**

```bash
# CPI-SI Runtime Environment - Generated from composed configs
if [ -f "$HOME/.claude/cpi-si/system/config/runtime/environment.conf" ]; then
    source "$HOME/.claude/cpi-si/system/config/runtime/environment.conf"
fi
```

**Then reload:**

```bash
source ~/.bashrc
```

**Note:** This sources the **generated** file, not a static config.

> [!IMPORTANT]
> Logout and login again (or source ~/.bashrc) for changes to take effect in all shells.

---

## Configuration Categories

### Package Management

**Eliminate installation prompts:**

| Variable | Value | Effect |
|----------|-------|--------|
| `DEBIAN_FRONTEND` | `noninteractive` | No apt configuration prompts |
| `NEEDRESTART_MODE` | `a` | Automatic restart decisions after updates |
| `PIP_NO_INPUT` | `1` | No pip confirmation prompts |
| `NPM_CONFIG_YES` | `true` | Auto-yes for npm operations |

**Examples:**

```bash
# All work without interactive prompts
sudo apt install python3-pip       # No debconf prompts
pip install requests               # No confirmation needed
npm init                           # Auto-accepts defaults
```

### Development Tools

**Language and tool optimization:**

| Category | Variables | Purpose |
|----------|-----------|---------|
| **Python** | `PYTHONUNBUFFERED=1`<br>`PYTHONDONTWRITEBYTECODE=1` | Real-time output<br>No .pyc files |
| **Node.js** | `NODE_ENV=development`<br>`NPM_CONFIG_COLOR=true` | Dev mode defaults<br>Colored output |
| **Rust** | `RUST_BACKTRACE=1`<br>`CARGO_TERM_COLOR=always` | Full error traces<br>Colored output |
| **Go** | `GOPATH=~/go` | Go workspace location |
| **Git** | `GIT_EDITOR=nano`<br>`GIT_PAGER=less` | Default editor<br>Pager for diff/log |
| **Sudo** | `SUDO_ASKPASS=~/.local/bin/sudo-askpass` | GUI password prompt helper<br>(Ready for sudo-rs support) |

### CPI-SI Framework

**Framework-specific paths and configuration:**

| Variable | Value | Purpose |
|----------|-------|---------|
| `NOVA_DAWN_WORKSPACE` | Project path | Instance-specific workspace location |
| `CPI_SI_BIBLICAL_PATH` | Biblical texts | Location of Scripture for framework |
| `CPI_SI_SKILLS_PATH` | Skills directory | Framework skills location |
| `CPI_SI_SYSTEM_PATH` | `~/.claude/system` | System tools and configuration |

**Why these matter:**

- Enable framework components to find resources
- Provide consistent paths across tools
- Support instance-specific workspace management
- Give direct access to biblical foundation

---

## Testing

**Verify integration:**

```bash
# Check package management variables
env | grep DEBIAN_FRONTEND   # Should show: noninteractive
env | grep PIP_NO_INPUT      # Should show: 1
env | grep NPM_CONFIG_YES    # Should show: true

# Check CPI-SI framework variables
env | grep CPI_SI            # Should show all framework paths
env | grep NOVA_DAWN         # Should show workspace path
```

**Test actual behavior:**

```bash
# Should complete without prompts
sudo apt install tree
pip install requests
npm init -y

# All should run non-interactively
```

**Run validation:**

```bash
cd ~/.claude/system
./bin/validate    # Validates environment integration
```

---

## Usage

These environment variables are automatically used by:

| Tool | Behavior | Without Variables | With Variables |
|------|----------|-------------------|----------------|
| `apt install` | Package installation | Prompts for configuration | Auto-installs with defaults |
| `pip install` | Python packages | May prompt for confirmations | Silent install |
| `npm init` | Project initialization | Asks multiple questions | Auto-accepts defaults |
| Development tools | Output formatting | May lack color/clarity | Optimized output |

**Real-world examples:**

```bash
# Before: Requires manual responses
sudo apt install postgresql
# Prompts: Do you want to continue? [Y/n]
# Prompts: Configure postgresql? [yes/no]

# After: Autonomous execution
sudo apt install postgresql
# Installs automatically with sensible defaults
```

> [!TIP]
> These variables enable Claude Code to perform installations and configurations autonomously while maintaining safety through sudoers boundaries.

---

## Integration with settings.json

Key variables are also defined in `~/.claude/settings.json` for Claude Code session environment:

```json
{
  "env": {
    "DEBIAN_FRONTEND": "noninteractive",
    "PIP_NO_INPUT": "1",
    "NPM_CONFIG_YES": "true",
    "NOVA_DAWN_WORKSPACE": "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC",
    "CPI_SI_BIBLICAL_PATH": "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC/foundation/biblical",
    "CPI_SI_SKILLS_PATH": "/home/seanje-lenox-wise/.claude/skills",
    "CPI_SI_SYSTEM_PATH": "/home/seanje-lenox-wise/.claude/system"
  }
}
```

**Why dual definition?**

| Location | Used By | Scope |
|----------|---------|-------|
| `non-interactive.conf` | Interactive shells | Your terminal sessions, bash scripts |
| `settings.json` | Claude Code tools | Tool executions (Bash, Write, etc.) |

This ensures variables are available in **both** interactive work and Claude Code autonomous operations.

---

## Troubleshooting

### Environment Variables Not Applied

**Symptoms:** Tools still prompting for input, CPI-SI paths not found

**Solutions:**

1. **Verify integration exists:**

   ```bash
   grep "CPI-SI" ~/.bashrc
   ```

   Should show the source line for `non-interactive.conf`

2. **Reload shell configuration:**

   ```bash
   source ~/.bashrc
   ```

3. **Restart Claude Code session:**
   - Settings changes require session restart
   - Use `/reset` command or restart

4. **Check shell type:**

   ```bash
   echo $SHELL
   ```

   If using zsh, add to `~/.zshrc` instead of `~/.bashrc`

### Still Seeing Prompts

**Possible causes:**

1. **Variable not exported:**

   ```bash
   echo $DEBIAN_FRONTEND
   ```

   Should show `noninteractive`, not empty

2. **Tool-specific configuration overrides:**
   - Some tools have their own config files (e.g., `~/.npmrc`)
   - These may override environment variables
   - Check tool documentation for precedence

3. **Session context:**
   - Variables set in current shell may not propagate to sudo
   - Sudoers `env_keep` preserves specific variables through sudo
   - Check: `sudo env | grep DEBIAN_FRONTEND`

**Run diagnostics:**

```bash
cd ~/.claude/system
./bin/diagnose    # Comprehensive environment diagnostics
```

---

<div align="center">

**Built with intentional design for Kingdom Technology**

[Back to System Documentation](../README.md)

</div>
