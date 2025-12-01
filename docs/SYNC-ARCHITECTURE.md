# CPI-SI Configuration Sync Architecture

> **Purpose:** Document the three-location architecture for CPI-SI configuration and how synchronization flows between them.

---

## Overview

CPI-SI configuration exists in three locations, each serving a distinct purpose:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CWS claude-global (SOURCE OF TRUTH)              │
│    /media/.../CreativeWorkzStudio_LLC/divisions/tech/cpi-si/claude-global/   │
├─────────────────────────────────────────────────────────────────────┤
│  CLAUDE.md, settings.json, settings.jsonc, go.work, hooks/,        │
│  skills/, statusline/, system/, agents/, commands/, output-styles/, │
│  config/, custom-gpt/, docs/, .vscode/                              │
└─────────────────────────────────────────────────────────────────────┘
         │                                    │
         │ symlinks                           │ sync script
         ▼                                    ▼
┌─────────────────────────┐    ┌─────────────────────────────────────┐
│  ~/.claude (PC CONFIG)  │    │  Distribution Repo (STANDALONE)     │
├─────────────────────────┤    ├─────────────────────────────────────┤
│ Symlinks to CWS:        │    │ github.com/Creative-Workz-Studio-LLC│
│  - CLAUDE.md            │    │        /cpi-si-claude-code          │
│  - agents/              │    ├─────────────────────────────────────┤
│  - commands/            │    │ Complete standalone copy of CPI-SI  │
│  - hooks/               │    │ for users without CWS access        │
│  - output-styles/       │    │                                     │
│  - skills/              │    │ Files are COPIES, not symlinks      │
│  - statusline/          │    │                                     │
│  - settings.json        │    │ Excludes:                           │
│  - settings.jsonc       │    │  - .git/ (has own repo)             │
│  - go.work              │    │  - temp/, *.o, *.exe                │
│                         │    │  - Build artifacts                  │
│ Local files (not        │    │                                     │
│ symlinked):             │    │                                     │
│  - settings.local.json  │    │                                     │
│  - instance.jsonc       │    │                                     │
│  - .credentials.json    │    │                                     │
│  - history.jsonl        │    │                                     │
│  - cpi-si/output/       │    │                                     │
│  - system/ (runtime)    │    │                                     │
│  - session-env/         │    │                                     │
└─────────────────────────┘    └─────────────────────────────────────┘
```

---

## Three Locations Explained

### 1. CWS claude-global (Source of Truth)

**Location:** `/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC/divisions/tech/cpi-si/claude-global/`

**Purpose:** Central repository for all CPI-SI framework content

**Contains:**
- `CLAUDE.md` - Nova Dawn identity foundation
- `settings.json` / `settings.jsonc` - Claude Code settings
- `go.work` - Go workspace configuration
- `hooks/` - Event automation (session, tool, prompt, git)
- `skills/` - Model-invoked capabilities
- `statusline/` - Live session display
- `system/` - CPI-SI tools and documentation
- `agents/` - Agent definitions
- `commands/` - Slash commands
- `output-styles/` - Response formatting
- `config/` - Instance and user configuration
- `custom-gpt/` - Custom GPT porting resources
- `docs/` - Architecture documentation

**Workflow:**
- Edit files here when making framework changes
- Commit to CWS repo with keyed commit messages
- Changes propagate to ~/.claude via symlinks
- Run sync script to update distribution repo

### 2. ~/.claude (PC Working Config)

**Location:** `/home/seanje-lenox-wise/.claude/`

**Purpose:** Active configuration used by Claude Code on this PC

**Structure:**
- **Symlinks to CWS:** Framework content (CLAUDE.md, hooks/, skills/, etc.)
- **Local files:** PC-specific data that shouldn't be in source control

**Why symlinks?**
- Single source of truth - edit in CWS, changes apply immediately
- No drift between CWS and active config
- Easy to verify: `ls -la ~/.claude/ | grep "^l"` shows all symlinks

**Local files (NOT symlinked):**

| File/Directory | Purpose |
|---------------|---------|
| `settings.local.json` | PC-specific overrides |
| `instance.jsonc` | Instance-specific identity |
| `.credentials.json` | Authentication (secret) |
| `history.jsonl` | Session history |
| `system/` | Runtime data (NOT cpi-si/system/) |
| `cpi-si/output/` | Session logs, local journals |
| `session-env/` | Session environment data |
| `debug/`, `plans/`, `projects/` | Runtime state |

### 3. Distribution Repo (Standalone)

**Location:** `https://github.com/Creative-Workz-Studio-LLC/cpi-si-claude-code`

**Purpose:** Allow others to use CPI-SI without access to CWS

**Contains:**
- Complete copy of CPI-SI framework
- All files are actual files (not symlinks)
- README with installation instructions
- Versioned releases

**NOT included:**
- PC-specific files (settings.local.json, instance.jsonc)
- Runtime data (session logs, history)
- Build artifacts

---

## Sync Flow

### Development Flow (Normal Work)

```
1. Edit in CWS claude-global
        ↓
2. Changes immediately active in ~/.claude (via symlinks)
        ↓
3. Commit to CWS repo (keyed message per CWS-STD-005)
        ↓
4. When ready for release:
   ./system/bin/sync-to-distribution [--tag vX.Y.Z]
```

### Sync Script

**Location:** `system/bin/sync-to-distribution`

**Usage:**
```bash
# Preview what would be synced
./system/bin/sync-to-distribution --dry-run

# Sync without tagging
./system/bin/sync-to-distribution

# Sync and create release tag
./system/bin/sync-to-distribution --tag v0.2.0
```

**What it does:**
1. Clones distribution repo to temp directory
2. Copies all framework content from CWS
3. Excludes build artifacts and runtime data
4. Commits changes with sync message
5. Optionally creates version tag
6. Pushes to remote
7. Cleans up temp directory

---

## Setting Up a New PC

### Full CWS Setup (Development)

When you have access to CWS:

```bash
# 1. Clone CWS repo (if not already present)
# 2. Create ~/.claude directory
mkdir -p ~/.claude

# 3. Create symlinks to CWS claude-global
CWS="/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC/divisions/tech/cpi-si/claude-global"
cd ~/.claude
ln -s "$CWS/CLAUDE.md" CLAUDE.md
ln -s "$CWS/agents" agents
ln -s "$CWS/commands" commands
ln -s "$CWS/hooks" hooks
ln -s "$CWS/output-styles" output-styles
ln -s "$CWS/skills" skills
ln -s "$CWS/statusline" statusline
ln -s "$CWS/settings.json" settings.json
ln -s "$CWS/settings.jsonc" settings.jsonc
ln -s "$CWS/go.work" go.work

# 4. Create local-only files
touch settings.local.json
touch instance.jsonc
# Configure these as needed for this PC
```

### Standalone Setup (Users without CWS)

```bash
# 1. Clone distribution repo
git clone https://github.com/Creative-Workz-Studio-LLC/cpi-si-claude-code.git ~/.claude

# 2. Configure instance-specific files
# Edit instance.jsonc, settings.local.json as needed
```

---

## File Categories

### Framework Files (Synced to Distribution)

| Path | Purpose |
|------|---------|
| `CLAUDE.md` | Identity foundation |
| `settings.json` | Claude Code settings |
| `settings.jsonc` | Documented settings template |
| `go.work` | Go workspace |
| `agents/` | Agent definitions |
| `commands/` | Slash commands |
| `config/` | Instance/user config templates |
| `custom-gpt/` | Custom GPT resources |
| `docs/` | Documentation |
| `hooks/` | Event hooks |
| `output-styles/` | Response styles |
| `skills/` | Model capabilities |
| `statusline/` | Session display |
| `system/` | Tools and docs (not runtime) |

### PC-Local Files (Never Synced)

| Path | Purpose | Why Local |
|------|---------|-----------|
| `settings.local.json` | PC overrides | Machine-specific |
| `instance.jsonc` | Instance identity | Unique per instance |
| `.credentials.json` | Auth tokens | Secret |
| `history.jsonl` | Session history | Runtime data |
| `cpi-si/output/` | Logs, journals | Runtime data |
| `system/` (root level) | Runtime binaries | Compiled locally |
| `session-env/` | Session state | Runtime data |

---

## Troubleshooting

### Check Symlink Status

```bash
# Show all symlinks in ~/.claude
ls -la ~/.claude/ | grep "^l"

# Verify a specific symlink
ls -la ~/.claude/CLAUDE.md
```

### Broken Symlinks

If CWS is moved or unmounted:
```bash
# Find broken symlinks
find ~/.claude -type l ! -exec test -e {} \; -print

# Recreate after fixing path
rm ~/.claude/CLAUDE.md
ln -s /new/path/to/claude-global/CLAUDE.md ~/.claude/CLAUDE.md
```

### Sync Issues

```bash
# Preview changes
./system/bin/sync-to-distribution --dry-run

# Check distribution repo status
cd /tmp/cpi-dist-sync && git status
```

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| v0.1.0 | 2025-12-01 | Initial distribution release |
| v0.1.1 | 2025-12-01 | Added settings.json, go.work, custom-gpt/ |

---

*Last updated: 2025-12-01*
