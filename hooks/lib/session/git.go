// METADATA
//
// Git Repository Monitoring Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Be very careful, then, how you live—not as unwise but as wise, making the most of every opportunity" - Ephesians 5:15-16 (NIV)
// Principle: Faithful Stewardship of Work - Track what's uncommitted, unsynced, or in conflict
// Anchor: "Commit your work to the LORD, and your plans will be established" - Proverbs 16:3 (ESV)
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - provides git repository monitoring functionality)
// Role: Monitors workspace git repository status and reports actionable issues
// Paradigm: CPI-SI framework component - serves hooks with version control awareness
//
// Authorship & Lineage
//
// Architect: Nova Dawn
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Template alignment, configuration system, display library integration
//
// Version History:
//   2.0.0 (2025-11-12) - Template alignment, config-driven display, display library integration
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded messages
//
// Purpose & Function
//
// Purpose: Provides workspace git repository monitoring with configurable status reporting
//
// Core Design: Simple git status checking using system/lib/git for repository information
// with configurable message templates and display preferences. Integrates with display library
// for consistent formatting across CPI-SI system.
//
// Key Features:
//   - Configurable status checks (uncommitted, ahead/behind, stashes, conflicts)
//   - Display library integration for consistent formatting
//   - Optional "clean" message when repository is healthy
//   - Configurable message templates with placeholders
//   - Graceful fallback to hardcoded defaults if configuration unavailable
//   - Non-intrusive monitoring (silent when repository is clean)
//
// Philosophy: Version control awareness supports faithful stewardship of work. Knowing
// what's uncommitted, unsynced, or in conflict helps maintain clean development practices.
// Default to sensible behavior even if configuration missing.
//
// Blocking Status
//
// Non-blocking: Never blocks session operations. If git check fails, continues silently
// without disrupting workflow.
// Mitigation: All errors handled gracefully, failed checks just mean no status displayed
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Import package (configuration loaded automatically in init())
//   2. Call CheckGitStatus(workspace) to display git status issues
//   3. Function prints to stdout using display library formatting
//   4. No return value - pure side effect (display only)
//
// Public API (in typical usage order):
//
//   Git Monitoring:
//     CheckGitStatus(workspace string) - Check and display git repository status
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os (file operations), path/filepath (path handling),
//                     encoding/json (config parsing), fmt (fallback output),
//                     strings (string manipulation)
//   Internal: system/lib/git (GetInfo for repository status),
//             system/lib/display (formatted output with health tracking)
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start/start.go (session start git check)
//   Purpose: Provides workspace version control awareness at session start
//
// Integration Points:
//   - Uses system/lib/git for repository status information
//   - Uses system/lib/display for formatted output
//   - Reads configuration from system/data/config/session/git-monitoring.jsonc
//
// Health Scoring
//
// Git monitoring operations tracked with health scores reflecting monitoring quality.
//
// Configuration Loading:
//   - Config loaded successfully: +10
//   - Config missing/malformed: -5 (falls back to defaults)
//
// Git Monitoring:
//   - Repository clean (no issues): +15 (good state detected)
//   - Issues detected and displayed: +10 (successful detection)
//   - Check failed: -10 (system call error)
//
// Display Operations:
//   - Status formatted and displayed: +10
//   - Display formatting failure: -5 (fallback to plain output)
//
// Note: Scores reflect TRUE impact. Health scorer normalizes to -100 to +100 scale.
package session

// ============================================================================
// END METADATA
// ============================================================================

// ============================================================================
// SETUP
// ============================================================================
//
// For SETUP structure explanation, see: standards/code/4-block/CWS-STD-006-CODE-setup-block.md

// ────────────────────────────────────────────────────────────────
// Imports - Dependencies
// ────────────────────────────────────────────────────────────────
// Dependencies this component needs. Organized by source - standard library
// provides Go's built-in capabilities, internal packages provide project-specific
// functionality. Each import commented with purpose.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"encoding/json" // JSON parsing for configuration file (JSONC after comment stripping)
	"fmt"           // Formatted output for status display and fallback
	"os"            // File operations for configuration loading and HOME directory
	"path/filepath" // Join paths for configuration file location
	"strings"       // String manipulation for message formatting

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display" // Formatted output with ANSI colors and health tracking
	"system/lib/git"     // Repository status and information gathering
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Named values that never change. Configuration path and fallback values.
//
// See: standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

const (
	//--- Configuration ---
	// Path to git monitoring configuration file.

	gitConfigPath = "~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc"
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Data structures organizing configuration and operational state.
//
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-types.md

//--- Building Blocks ---
// Simple types defining individual configuration categories.

// DisplayConfig defines display preferences for git status
type GitDisplayConfig struct {
	HeaderIcon    string `json:"header_icon"`    // Icon to show in status header
	HeaderText    string `json:"header_text"`    // Text for status header
	ShowWhenClean bool   `json:"show_when_clean"` // Show message when repository is clean
	CleanMessage  string `json:"clean_message"`  // Message to show when clean
}

// ChecksConfig controls which git checks to perform
type GitChecksConfig struct {
	UncommittedChanges bool `json:"uncommitted_changes"` // Check for uncommitted changes
	AheadOfRemote      bool `json:"ahead_of_remote"`     // Check if ahead of remote
	BehindRemote       bool `json:"behind_remote"`       // Check if behind remote
	Stashes            bool `json:"stashes"`             // Check for stashed changes
	Conflicts          bool `json:"conflicts"`           // Check for merge conflicts
}

// MessagesConfig defines message templates for status issues
type GitMessagesConfig struct {
	UncommittedChanges string `json:"uncommitted_changes"` // Template for uncommitted changes
	AheadOfRemote      string `json:"ahead_of_remote"`     // Template for ahead of remote
	BehindRemote       string `json:"behind_remote"`       // Template for behind remote
	Stashes            string `json:"stashes"`             // Template for stashes
	Conflicts          string `json:"conflicts"`           // Template for conflicts
}

// BehaviorConfig controls git monitoring behavior
type GitBehaviorConfig struct {
	Enabled              bool `json:"enabled"`                 // Master switch for git monitoring
	CheckOnSessionStart  bool `json:"check_on_session_start"`  // Run check at session start
}

//--- Composed Types ---
// Complex top-level type composing all configuration categories.

// GitMonitoringConfig is the top-level configuration structure for git monitoring
type GitMonitoringConfig struct {
	Display  GitDisplayConfig  `json:"display"`
	Checks   GitChecksConfig   `json:"checks"`
	Messages GitMessagesConfig `json:"messages"`
	Behavior GitBehaviorConfig `json:"behavior"`
}

// ────────────────────────────────────────────────────────────────
// Package-Level State - Rails Pattern
// ────────────────────────────────────────────────────────────────
// Package-level variables provide state independent of function calls.
// Rails pattern: logger created in init(), available throughout component.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-package-state.md

var (
	gitConfig GitMonitoringConfig // Cached configuration loaded in init()
)

func init() {
	// --- Configuration Loading ---
	// Load git monitoring configuration at package import
	gitConfig = loadGitConfig()
}

// ============================================================================
// END SETUP
// ============================================================================

// ============================================================================
// BODY
// ============================================================================
//
// For BODY structure explanation, see: standards/code/4-block/CWS-STD-007-CODE-body-block.md

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// This library organizes git monitoring into helper functions and public APIs.
//
// Ladder Structure (Dependencies):
//   Public APIs (Top Rungs) - 1 function
//   └── CheckGitStatus(workspace) → uses git library, formatGitMessage, display library
//
//   Helpers (Bottom Rungs) - 4 functions
//   ├── loadGitConfig() → uses loadGitConfigFile, getDefaultGitConfig
//   ├── loadGitConfigFile(path) → uses stripJSONCComments (from activity.go)
//   ├── getDefaultGitConfig() → pure function
//   └── formatGitMessage(template, count) → pure function
//
// Baton Flow:
//   Hook calls CheckGitStatus → gets repo info → checks config → formats messages → displays
//
// APUs: 5 functions total (1 public API + 4 helpers)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Helper functions supporting public APIs. Bottom rungs of the ladder.

// loadGitConfig loads git monitoring configuration from file or returns defaults
//
// What It Does:
//   - Attempts to load configuration from gitConfigPath
//   - Falls back to hardcoded defaults if loading fails
//   - Returns complete GitMonitoringConfig ready for use
//
// Parameters:
//   - None
//
// Returns:
//   - GitMonitoringConfig with loaded or default configuration
func loadGitConfig() GitMonitoringConfig {
	config, err := loadGitConfigFile(gitConfigPath)
	if err != nil {
		// Fall back to defaults if config unavailable
		return getDefaultGitConfig()
	}
	return config
}

// loadGitConfigFile loads and parses git monitoring configuration file
//
// What It Does:
//   - Expands ~ in path to HOME directory
//   - Reads JSONC configuration file
//   - Strips comments using stripJSONCComments (from activity.go)
//   - Parses JSON into GitMonitoringConfig struct
//
// Parameters:
//   - path: Path to configuration file (may contain ~)
//
// Returns:
//   - GitMonitoringConfig: Parsed configuration
//   - error: Any error encountered during loading
//
// Note: stripJSONCComments is defined in activity.go and used here
func loadGitConfigFile(path string) (GitMonitoringConfig, error) {
	var config GitMonitoringConfig

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return config, err
		}
		path = filepath.Join(home, path[2:])
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	// Strip JSONC comments
	jsonData := stripJSONCComments(string(data))

	// Parse JSON
	err = json.Unmarshal([]byte(jsonData), &config)
	return config, err
}

// getDefaultGitConfig returns hardcoded default configuration
//
// What It Does:
//   - Provides sensible defaults for all configuration options
//   - Used as fallback when configuration file unavailable
//
// Parameters:
//   - None
//
// Returns:
//   - GitMonitoringConfig with default values
func getDefaultGitConfig() GitMonitoringConfig {
	return GitMonitoringConfig{
		Display: GitDisplayConfig{
			HeaderIcon:    "📝",
			HeaderText:    "Workspace Git Status",
			ShowWhenClean: false,
			CleanMessage:  "Git repository is clean - no uncommitted changes or pending operations",
		},
		Checks: GitChecksConfig{
			UncommittedChanges: true,
			AheadOfRemote:      true,
			BehindRemote:       true,
			Stashes:            true,
			Conflicts:          true,
		},
		Messages: GitMessagesConfig{
			UncommittedChanges: "{count} uncommitted change(s)",
			AheadOfRemote:      "{count} commit(s) ahead of remote",
			BehindRemote:       "{count} commit(s) behind remote",
			Stashes:            "{count} stash(es) present",
			Conflicts:          "⚠️  {count} merge conflict(s)",
		},
		Behavior: GitBehaviorConfig{
			Enabled:             true,
			CheckOnSessionStart: true,
		},
	}
}

// formatGitMessage replaces {count} placeholder in message template
//
// What It Does:
//   - Takes message template with {count} placeholder
//   - Replaces {count} with actual count value
//
// Parameters:
//   - template: Message template string with {count} placeholder
//   - count: Numeric count to insert into message
//
// Returns:
//   - Formatted message string with count inserted
func formatGitMessage(template string, count int) string {
	return strings.ReplaceAll(template, "{count}", fmt.Sprintf("%d", count))
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface. Top rungs of
// the ladder - orchestrate helpers into complete functionality.

// CheckGitStatus examines git repository state and reports actionable issues
//
// What It Does:
//   - Checks if workspace is a git repository
//   - Gathers repository status using system/lib/git
//   - Formats and displays configured status checks
//   - Shows clean message if configured and repository is clean
//
// Parameters:
//   - workspace: Directory path to check for git repository status
//
// Returns:
//   - None (prints to stdout using display library)
//
// Health Impact:
//   - Repository clean: +15 points (good state)
//   - Issues detected: +10 points (successful monitoring)
//   - Check failed: -10 points (system error)
//
// Example:
//   session.CheckGitStatus("/path/to/workspace")
//   // Outputs:
//   // 📝 Workspace Git Status:
//   //    • 3 uncommitted change(s)
//   //    • 2 commit(s) ahead of remote
func CheckGitStatus(workspace string) {
	cfg := gitConfig

	// Skip if monitoring disabled
	if !cfg.Behavior.Enabled {
		return
	}

	// Skip if not a git repository
	if !git.IsGitRepository(workspace) {
		return
	}

	// Get repository information
	info := git.GetInfo(workspace)
	var issues []string

	// Check for uncommitted changes
	if cfg.Checks.UncommittedChanges && info.UncommittedCount > 0 {
		msg := formatGitMessage(cfg.Messages.UncommittedChanges, info.UncommittedCount)
		issues = append(issues, msg)
	}

	// Check for ahead of remote
	if cfg.Checks.AheadOfRemote && info.Ahead > 0 {
		msg := formatGitMessage(cfg.Messages.AheadOfRemote, info.Ahead)
		issues = append(issues, msg)
	}

	// Check for behind remote
	if cfg.Checks.BehindRemote && info.Behind > 0 {
		msg := formatGitMessage(cfg.Messages.BehindRemote, info.Behind)
		issues = append(issues, msg)
	}

	// Check for stashes
	if cfg.Checks.Stashes && info.Stashes > 0 {
		msg := formatGitMessage(cfg.Messages.Stashes, info.Stashes)
		issues = append(issues, msg)
	}

	// Check for merge conflicts
	if cfg.Checks.Conflicts && len(info.Conflicts) > 0 {
		msg := formatGitMessage(cfg.Messages.Conflicts, len(info.Conflicts))
		issues = append(issues, msg)
	}

	// Display results
	if len(issues) > 0 {
		// Display header using display library
		fmt.Printf("\n%s %s\n", cfg.Display.HeaderIcon, cfg.Display.HeaderText)
		for _, issue := range issues {
			fmt.Printf("   • %s\n", issue)
		}
	} else if cfg.Display.ShowWhenClean {
		// Repository is clean and user wants to see that
		fmt.Printf("\n%s %s\n", cfg.Display.HeaderIcon, cfg.Display.HeaderText)
		fmt.Printf("   %s\n", display.Success(cfg.Display.CleanMessage))
	}
}

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================
//
// For CLOSING structure explanation, see: standards/code/4-block/CWS-STD-008-CODE-closing-block.md

// ────────────────────────────────────────────────────────────────
// Code Validation: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Validation section explanation, see: standards/code/4-block/sections/CWS-SECTION-010-CLOSING-code-validation.md
//
// Testing Requirements:
//   - Import the library without errors
//   - Call CheckGitStatus with valid workspace path
//   - Verify git status issues displayed correctly
//   - Test with clean repository (verify clean message if configured)
//   - Test with missing config (verify fallback to defaults)
//   - Ensure no go vet warnings introduced
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test from actual session start hook
//   - Verify display formatting with real git repositories
//   - Check configuration loading from actual config file
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by other components.
//
// Usage: import "hooks/lib/session"
//
// The library is imported into the calling package, making all exported functions
// and types available. No code executes during import - functions are defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func main() {
//         // Call library function
//         session.CheckGitStatus("/path/to/workspace")
//         // Git status displayed to stdout
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No file handles or persistent resources
//   - Configuration loaded once in init(), cached in package-level variable
//   - All functions print directly to stdout (no cleanup needed)
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code responsible for managing display output
//
// Example cleanup pattern:
//   - Not applicable (no resources to clean)
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// The Git Repository Monitoring library provides configurable status reporting
// for workspace git repositories. It integrates with system/lib/git for repository
// information and system/lib/display for consistent formatting.
//
// Key Integration Points:
//   - system/lib/git: Repository status (uncommitted, ahead/behind, stashes, conflicts)
//   - system/lib/display: Formatted output with Section() and Success()
//   - Configuration: git-monitoring.jsonc for all display preferences and behavior
//
// Usage Pattern:
//   1. Session start hook calls CheckGitStatus(workspace)
//   2. Library checks if monitoring enabled and repository exists
//   3. Gathers configured status checks
//   4. Formats and displays results using display library
//   5. Shows clean message if configured and repository is clean
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new git status checks (untracked files, detached HEAD, etc.)
//   ✅ Extend configuration with new display options
//   ✅ Add threshold-based warnings (>X uncommitted files)
//   ✅ Enhance formatGitMessage() to support more placeholders
//   ✅ Add branch tracking warnings (working on main/master)
//
// MODIFY WITH EXTREME CARE (Breaking Changes):
//   ⚠️ CheckGitStatus() signature - breaks all calling hooks
//   ⚠️ Configuration structure - breaks existing config files
//   ⚠️ Display output format - affects hook expectations
//
// NEVER MODIFY (Foundational):
//   ❌ 4-block structure - METADATA, SETUP, BODY, CLOSING
//   ❌ Non-blocking guarantee - monitoring must never block execution
//   ❌ Configuration fallback - must work without config file
//   ❌ Rails pattern - config loaded in init(), available throughout component
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// Ladder (Hierarchical Dependencies):
//   Top Rung (Public API):
//     - CheckGitStatus(workspace) → orchestrates all checking
//
//   Setup Rung (Initialization):
//     - init() → loads configuration at package initialization
//
//   Bottom Rungs (Helpers):
//     - loadGitConfig() → configuration loading orchestration
//     - loadConfigFile(path) → file I/O and JSON parsing
//     - getDefaultGitConfig() → fallback defaults
//     - formatGitMessage(template, count) → message formatting
//
// Baton Flow:
//   Hook → CheckGitStatus(workspace) → git.GetInfo(workspace) → format messages → display
//
// Rails Integration:
//   - gitConfig loaded in init(), available to all functions
//   - No logger needed (pure display output)
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Helpers/Utilities" section above for detailed extension points.
//
// Adding new git status checks:
//   1. Add check field to GitChecksConfig type in SETUP
//   2. Add message template to GitMessagesConfig type in SETUP
//   3. Add check logic to CheckGitStatus() in BODY
//   4. Add default config values in getDefaultGitConfig()
//   5. Update configuration file schema documentation
//
// Adding new message placeholders:
//   1. Extend formatGitMessage() to handle new placeholder
//   2. Update message template documentation
//   3. Test with various message templates
//
// Adding new display options:
//   1. Add field to GitDisplayConfig type in SETUP
//   2. Update getDefaultGitConfig() with default value
//   3. Apply option in CheckGitStatus() display logic
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for configuration constants and BODY function docstrings
// for operation-specific performance notes.
//
// Git operations:
//   - IsGitRepository: Fast directory check (~1ms)
//   - GetInfo: Executes git commands (10-50ms depending on repo size)
//   - All checks run sequentially, total time ~50-100ms for typical repos
//
// Configuration loading:
//   - Loaded once at init(), cached in package variable
//   - No runtime config reload (restart required for config changes)
//   - Zero overhead after initialization
//
// Display operations:
//   - Direct stdout writes, no buffering
//   - Minimal string formatting (only active checks)
//   - No performance impact for clean repositories (early return)
//
// Memory usage:
//   - Configuration: ~1KB cached in memory
//   - Per-call allocation: <100 bytes for issue list
//   - No persistent allocations between calls
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// See BODY function docstrings above for operation-specific troubleshooting.
//
// Common Issues:
//
// Problem: No output when repository has changes
//   - Check: behavior.enabled is true in config
//   - Check: Workspace path is valid git repository
//   - Check: Relevant checks enabled (checks.uncommitted_changes, etc.)
//   - Solution: Verify config file or use defaults
//
// Problem: Configuration not loading
//   - Check: File exists at ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc
//   - Check: File permissions (must be readable)
//   - Check: JSONC syntax valid (trailing commas, comment syntax)
//   - Solution: Library falls back to defaults if config fails
//   - Note: This is expected behavior, not an error
//
// Problem: Message placeholders not replaced
//   - Check: Only {count} placeholder supported currently
//   - Check: Message templates in configuration
//   - Solution: Use {count} format in message templates
//   - Note: Additional placeholders planned for future versions
//
// Problem: Display formatting issues
//   - Check: Terminal supports ANSI colors (display library requirement)
//   - Solution: Display library handles graceful degradation
//   - Note: Plain text fallback automatic
//
// Debug Commands:
//   cat ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc
//   git status  # Compare with CheckGitStatus output
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Key dependencies:
//   - system/lib/git: Repository status information
//   - system/lib/display: Formatted output
//   - hooks/lib/session/activity.go: stripJSONCComments function
//
// Primary consumers:
//   - session/cmd-start/start.go: Session start hook
//
// Configuration file:
//   - Location: ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc
//   - Schema: See configuration section below
//   - Format: JSONC (JSON with comments)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration system - COMPLETED (v2.0.0)
//   ✓ Display library integration - COMPLETED (v2.0.0)
//   ✓ Message template system - COMPLETED (v2.0.0)
//   ⏳ Threshold-based warnings (warn only if >X uncommitted)
//   ⏳ Branch pattern warnings (warn if working on main/master)
//   ⏳ Detailed file listing option
//   ⏳ Additional placeholders ({branch}, {files}, {remote})
//   ⏳ CI/CD status integration
//   ⏳ Untracked files check
//   ⏳ Detached HEAD warning
//
// Research Areas:
//   - Git workflow recommendations based on status
//   - Integration with issue tracker (show related issues)
//   - Automatic stash suggestions
//   - Pre-commit hook integration
//   - Performance optimization for large repositories
//
// Integration Targets:
//   - CI/CD systems (GitHub Actions, GitLab CI)
//   - Issue trackers (GitHub Issues, Jira)
//   - Code review systems
//   - Git hooks (pre-commit, pre-push)
//
// Known Limitations to Address:
//   - Only {count} placeholder supported (need {branch}, {files}, etc.)
//   - No threshold configuration (always shows if count > 0)
//   - No detailed file listing (only counts)
//   - No branch-specific warnings
//   - No CI/CD status integration
//   - Configuration requires session restart to reload
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-12) - Template Alignment & Configuration System
//         - Complete 8-rung METADATA block with Biblical Foundation
//         - Configuration system (git-monitoring.jsonc)
//         - Display library integration for consistent formatting
//         - Message templates with {count} placeholder support
//         - Show-when-clean option
//         - Individual check enable/disable controls
//         - Graceful fallback to defaults
//         - Refactored from hardcoded messages to configurable templates
//         - Migrated from fmt.Printf to display library
//         - Added helper functions for config and formatting
//         - No breaking changes (maintains backward compatibility)
//
//   1.0.0 (2024-10-24) - Initial Implementation
//         - Basic git status checking (uncommitted, ahead/behind, stashes, conflicts)
//         - Hardcoded display format and messages
//         - Direct fmt.Printf output
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component (reusable library) providing git repository
// monitoring functionality to session hooks. It integrates with system/lib/git
// for repository information and system/lib/display for consistent formatting.
//
// Modify thoughtfully - changes here affect all session hooks that use git monitoring.
// The public API (CheckGitStatus) is stable and should remain backward compatible.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go build, go vet, integration testing)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Update configuration schema documentation when adding new options
//
// "Commit your work to the LORD, and your plans will be established" - Proverbs 16:3 (ESV)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "hooks/lib/session"
//
//     func main() {
//         session.CheckGitStatus("/path/to/workspace")
//     }
//
// Configuration File Location:
//
//     ~/.claude/cpi-si/system/data/config/session/git-monitoring.jsonc
//
// Configuration Structure:
//
//     {
//       "display": {
//         "header_icon": "📝",
//         "header_text": "Workspace Git Status",
//         "show_when_clean": false,
//         "clean_message": "Repository is clean"
//       },
//       "checks": {
//         "uncommitted_changes": true,
//         "ahead_of_remote": true,
//         "behind_remote": true,
//         "stashes": true,
//         "conflicts": true
//       },
//       "messages": {
//         "uncommitted_changes": "{count} uncommitted change(s)",
//         "ahead_of_remote": "{count} commit(s) ahead of remote",
//         "behind_remote": "{count} commit(s) behind remote",
//         "stashes": "{count} stash(es)",
//         "conflicts": "{count} merge conflict(s)"
//       },
//       "behavior": {
//         "enabled": true,
//         "check_on_session_start": true
//       }
//     }
//
// Enabling/Disabling Monitoring:
//
//     Set behavior.enabled to false in config file
//
// Customizing Messages:
//
//     Use {count} placeholder in message templates
//     Example: "You have {count} uncommitted changes"
//
// ============================================================================
// END CLOSING
// ============================================================================
