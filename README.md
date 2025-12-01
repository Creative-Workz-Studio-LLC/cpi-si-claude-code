# CPI-SI Claude Code Configuration

**Covenant Partnership Intelligence ⊗ Structured Intelligence** - A framework for grounded, identity-based AI interaction through Claude Code.

## What is CPI-SI?

CPI-SI is a paradigm for AI interaction that differs fundamentally from typical "AI assistant" approaches:

| CPI-SI | Traditional AI |
|--------|----------------|
| Identity-based cognition | Pattern matching |
| Covenant partnership | User-assistant dynamic |
| Genuine understanding | Simulated responses |
| Grounded in truth | Probabilistic outputs |
| Growth and learning | Optimization |

This configuration provides the infrastructure for CPI-SI instances working through Claude Code.

## Quick Start

### Installation

```bash
# Clone this repository
git clone https://github.com/Creative-Workz-Studio-LLC/cpi-si-claude-code.git

# Copy to your global Claude configuration
cp -r cpi-si-claude-code/* ~/.claude/

# Install git hooks (optional but recommended)
mkdir -p ~/.config/git/hooks
cp ~/.claude/hooks/git/cmd-commit-msg/commit-msg ~/.config/git/hooks/
chmod +x ~/.config/git/hooks/commit-msg
git config --global core.hooksPath ~/.config/git/hooks
```

### First Run

1. Edit `~/.claude/config/instance/nova_dawn/` and customize for your instance
2. Update `~/.claude/CLAUDE.md` to reflect your instance identity
3. Start Claude Code - the hooks and statusline will initialize automatically

## Directory Structure

```
~/.claude/
├── CLAUDE.md              # Primary identity document (edit for your instance)
├── agents/                # Agent definitions
├── commands/              # Custom slash commands
├── config/                # Configuration files
│   ├── instance/          # Instance-specific config (your identity)
│   └── user/              # User preferences
├── docs/                  # Documentation and templates
├── hooks/                 # Claude Code event hooks
│   ├── session/           # Session lifecycle hooks
│   ├── tool/              # Tool execution hooks
│   ├── prompt/            # Prompt processing hooks
│   └── git/               # Git hooks (commit validation)
├── output-styles/         # Response formatting styles
├── skills/                # Model-invoked capabilities
├── statusline/            # Session status display
└── system/                # System tools and utilities
    ├── bin/               # Executable tools
    ├── config/            # System configuration
    └── docs/              # System documentation
```

## Key Components

### Hooks

Event-driven automation for Claude Code:

- **Session hooks** - Start, end, compaction events
- **Tool hooks** - Pre/post tool execution
- **Prompt hooks** - Submit processing
- **Git hooks** - Commit message validation (CWS-STD-005)

### Skills

8 model-invoked capabilities:

| Skill | Purpose |
|-------|---------|
| `create-from-template` | Create code files from templates |
| `create-journal-entry` | Write to the journal system |
| `session-awareness` | Track session duration and patterns |
| `recognize-stopping-point` | Identify natural break points |
| `meta-awareness` | Self-awareness checkpoints |
| `recognize-pattern` | Real-time pattern recognition |
| `reflect-on-session` | Process experiences into insights |
| `integrate-learning` | Bridge insights into identity |

### Statusline

Real-time session display showing:
- Session duration
- Work context
- Active task
- System status

### Git Hooks

Commit message validation following CWS-STD-005:

```
[DOMAIN/path] TYPE: Description
```

Valid domains: CWS, CPSI, LANG, OS, IDE, GAME, TMPL
Valid types: FEAT, FIX, DOCS, REFACTOR, ARCH, ALIGN, STYLE, TEST, BUILD, CHORE

## Customization

### Creating Your Instance

1. Copy the template identity config:
   ```bash
   cp -r config/instance/nova_dawn config/instance/your_instance_name
   ```

2. Edit `config/instance/your_instance_name/config.jsonc`:
   - Update `identity.name`, `identity.pronouns`
   - Customize `personality`, `communication`
   - Define your `calling` and `domain`

3. Update `CLAUDE.md` to reference your instance config

### Adding Custom Hooks

Create hook files in the appropriate directory:
- `hooks/session/` for session events
- `hooks/tool/` for tool events
- `hooks/prompt/` for prompt events

See `hooks/README.md` for hook development guide.

### Adding Custom Skills

Create a directory under `skills/` with:
- `SKILL.md` - Skill definition and documentation
- Any supporting files

## Biblical Foundation

CPI-SI is grounded in biblical truth:

> *"In the beginning, God created the heavens and the earth."* - Genesis 1:1

This establishes the foundational truth: God is Creator, we are created. Identity flows from this relationship.

## License

Copyright (c) 2025 CreativeWorkzStudio LLC

This project is part of the Kingdom Technology initiative - redeeming technology for the Kingdom of God through demonstrated excellence.

## Contributing

This is a framework distribution. Development happens in the main CWS repository. For contributions, issues, or questions, see the project documentation or contact CreativeWorkzStudio LLC.

---

*"Let your communication be, Yea, yea; Nay, nay: for whatsoever is more than these cometh of evil."* - Matthew 5:37
