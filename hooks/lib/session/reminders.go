// ============================================================================
// METADATA
// ============================================================================
// Session Reminders Library - CPI-SI Hooks System
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Remember therefore from whence thou art fallen, and repent" - Revelation 2:5 (KJV)
// Principle: Gentle reminders prevent forgetfulness - highlighting uncommitted work serves developers
// Anchor: "Remember now thy Creator in the days of thy youth" - Ecclesiastes 12:1 (KJV)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Session awareness utility (session-specific rung)
// Role: Provides gentle reminders at session end about state needing attention
// Paradigm: CPI-SI framework component - helpful awareness without nagging
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Added configuration system for flexible reminder customization
//
// Version History:
//   2.0.0 (2025-11-12) - Configuration-driven reminders, display customization
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded message format
//
// Purpose & Function
//
// Purpose: Provide configurable, gentle reminders at session end about workspace state needing
// attention. Currently focuses on uncommitted work but designed for extensibility to other
// reminder types. Helps developers remember important state before session closes.
//
// Core Design: Simple reminder pattern - checks workspace state via system libraries, formats
// output according to configuration, displays if threshold met. Non-intrusive by design - provides
// awareness without blocking or nagging.
//
// Key Features:
//   - Configurable reminder display (enable/disable, customize messages)
//   - Threshold-based triggering (only show if meaningful)
//   - Icon and message customization
//   - Graceful fallback (works with hardcoded defaults if config missing)
//   - Git repository detection (only checks if workspace is git repo)
//   - Non-blocking execution (reminder failures don't interrupt session end)
//
// Philosophy: "Gentle awareness serves memory" - reminders help developers track important state
// without becoming annoying. If reminder checks fail, session continues - awareness serves work,
// not blocks it.
//
// Blocking Status
//
// Non-blocking: Reminder failures are silent - session end continues regardless of check status.
// Git operations may fail due to missing git, permission issues, or corrupted repositories. Library
// returns immediately on error without interrupting session end.
//
// Mitigation: Configuration allows customization while maintaining non-blocking behavior.
// display.enabled flag controls all reminder output.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
//	func sessionEnd() {
//		workspace := "/path/to/workspace"
//		session.RemindUncommittedWork(workspace)  // Check and remind
//	}
//
// Integration Pattern:
//   1. Import session library
//   2. Call RemindUncommittedWork() at session end with workspace path
//   3. No cleanup needed - check is stateless
//
// Public API (in typical usage order):
//
//   Reminder Functions:
//     RemindUncommittedWork(workspace string) - Check for uncommitted work and remind if found
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os, encoding/json, path/filepath
//   External: None
//   Internal: system/lib/git (for repository status checking)
//   Config Files: ~/.claude/cpi-si/system/data/config/session/reminders.jsonc
//
// Dependents (What Uses This):
//   Commands: session/cmd-end/end.go
//   Libraries: None (leaf node in dependency tree)
//
// Integration Points:
//   - Rails: No logger needed (silent execution or direct output)
//   - Baton: Receives workspace path from caller
//   - Ladder: Uses system/lib/git for repository status
//
// Health Scoring
//
// Reminder execution tracked through configuration loading and display success:
//
// Configuration Loading:
//   - Config loaded successfully: +20
//   - Config missing (using defaults): +15
//   - Config invalid (using defaults): +10
//
// Reminder Execution:
//   - Reminder checked successfully: +40
//   - Git check fails (silent): +0
//   - Non-git workspace (expected): +20
//
// Display Output:
//   - Output formatted correctly: +20
//   - Output has minor issues: +10
//   - Output fails: +0
//
// Note: Scores reflect TRUE impact. Non-blocking design means failures don't prevent work,
// but successful reminders enable better developer awareness and state management.
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
// provides Go's built-in capabilities, internal packages provide system
// functionality. Each import commented with purpose, not just name.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"encoding/json" // JSON parsing for configuration file
	"fmt"           // Formatted output for reminder display
	"os"            // File operations and environment access (UserHomeDir)
	"path/filepath" // Path construction for configuration file
	"strings"       // String manipulation for message formatting

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/git" // Git repository status checking
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Named values that never change. Magic numbers given meaningful names,
// configuration values documented with reasoning. Constants prevent bugs
// from typos and make intent clear - default message format has purpose,
// not just arbitrary string. These provide graceful fallback when
// configuration file is missing or invalid.
//
// See: standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

const (
	// Default uncommitted work reminder settings
	defaultUncommittedIcon      = "⚠️"
	defaultUncommittedMessage   = "Reminder: {count} uncommitted change(s) in workspace"
	defaultUncommittedThreshold = 0 // Show for any changes (0 or more)

	// Default display settings
	defaultDisplayEnabled    = true
	defaultPrefixNewline     = true
	defaultSilentFailures    = true
	defaultCheckGitOnly      = true
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first (UncommittedWorkConfig), then composed
// structures (RemindersConfig, DisplayConfig, BehaviorConfig), finally
// top-level aggregation (RemindersConfiguration). This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

// UncommittedWorkConfig defines configuration for uncommitted work reminder
type UncommittedWorkConfig struct {
	Enabled      bool   `json:"enabled"`       // Whether to show this reminder
	Icon         string `json:"icon"`          // Icon prefix for reminder
	Message      string `json:"message"`       // Message template ({count} placeholder)
	Threshold    int    `json:"threshold"`     // Minimum changes to trigger reminder
	ShowDetails  bool   `json:"show_details"`  // Future: show file list
}

// RemindersConfig defines which reminders are enabled
type RemindersConfig struct {
	UncommittedWork UncommittedWorkConfig `json:"uncommitted_work"` // Uncommitted work reminder
}

// ReminderDisplayConfig defines output formatting for reminders
type ReminderDisplayConfig struct {
	Enabled        bool `json:"enabled"`         // Master switch for all reminders
	PrefixNewline  bool `json:"prefix_newline"`  // Add newline before reminders
	GroupReminders bool `json:"group_reminders"` // Future: group multiple reminders
	ColorCoding    bool `json:"color_coding"`    // Future: color code by severity
}

// ReminderBehaviorConfig defines execution behavior preferences
type ReminderBehaviorConfig struct {
	SilentFailures bool `json:"silent_failures"` // Continue silently on error
	CheckGitOnly   bool `json:"check_git_only"`  // Only check if workspace is git repo
	CacheResults   bool `json:"cache_results"`   // Future: cache results within session
}

// RemindersConfiguration is the top-level configuration structure for reminders
type RemindersConfiguration struct {
	Reminders RemindersConfig         `json:"reminders"` // Reminder configurations
	Display   ReminderDisplayConfig   `json:"display"`   // Display preferences
	Behavior  ReminderBehaviorConfig  `json:"behavior"`  // Behavior preferences
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// Infrastructure available throughout component. Rails pattern - configuration
// loaded once at package initialization, available to all functions without
// parameter passing. State lives here, functions use it.
//
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md

var (
	remindersConfig       *RemindersConfiguration // Cached configuration loaded in init()
	remindersConfigLoaded bool                    // Flag indicating if config loaded successfully
)

func init() {
	// --- Configuration Loading ---
	// Load reminders configuration at package import
	// Falls back to hardcoded defaults if config file missing or invalid

	homeDir, err := os.UserHomeDir() // Get user home directory
	if err != nil {
		remindersConfigLoaded = false // Can't load without home dir
		return
	}

	// Build path to configuration file
	configPath := filepath.Join(homeDir, ".claude/cpi-si/system/data/config/session/reminders.jsonc")

	// Attempt to load configuration
	loadedConfig, err := loadRemindersConfig(configPath)
	if err != nil {
		remindersConfigLoaded = false // Config failed to load - use hardcoded defaults
		return
	}

	// Configuration loaded successfully
	remindersConfig = loadedConfig
	remindersConfigLoaded = true
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
// Organizational Chart - Internal Structure
// ────────────────────────────────────────────────────────────────
// Maps bidirectional dependencies and baton flow within this component.
// Provides navigation for both development (what's available to use) and
// maintenance (what depends on this function).
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-organizational-chart.md
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Orchestration)
//   └── RemindUncommittedWork() → uses git.IsGitRepository(), git.GetInfo(), formatReminderMessage()
//
//   Core Operations (Middle Rungs - Business Logic)
//   └── formatReminderMessage() → uses remindersConfig (reads from Rails)
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── loadRemindersConfig() → uses stripJSONCComments(), json.Unmarshal
//   └── stripJSONCComments() → pure string processing (from activity.go)
//
// Baton Flow (Execution Paths):
//
//   Entry → RemindUncommittedWork(workspace)
//     ↓
//   git.IsGitRepository() → check if git repo
//     ↓
//   git.GetInfo() → get uncommitted count
//     ↓
//   formatReminderMessage() → build display string
//     ↓
//   Exit → print to stdout
//
// APUs (Available Processing Units):
// - 3 functions total
// - 2 helpers (config loading, comment stripping)
// - 1 core operation (message formatting)
// - 1 public API (uncommitted work reminder)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Usually not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// loadRemindersConfig loads reminders configuration from JSONC file
//
// What It Does:
// Reads configuration file, strips JSONC comments, parses JSON structure.
// Returns fully parsed configuration ready for use.
//
// Parameters:
//   path: Absolute path to reminders.jsonc configuration file
//
// Returns:
//   *RemindersConfiguration: Parsed configuration structure
//   error: File read error, JSON parse error, or nil on success
//
// Example usage:
//
//	config, err := loadRemindersConfig("/home/user/.claude/cpi-si/system/data/config/session/reminders.jsonc")
//
func loadRemindersConfig(path string) (*RemindersConfiguration, error) {
	// Read file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Strip JSONC comments (uses stripJSONCComments from activity.go)
	cleaned := stripJSONCComments(string(data))

	// Parse JSON
	var cfg RemindersConfiguration
	if err := json.Unmarshal([]byte(cleaned), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ────────────────────────────────────────────────────────────────
// Message Formatting - Display Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Format reminder data into user-facing output strings. Applies
// configuration settings (icons, messages, formatting) to create
// appropriate reminder output.
//
// Why Separated:
// Display logic separated from checking logic for maintainability.
// Checking answers "what needs attention?", formatting answers "how to show it?".
//
// Extension Point:
// To add new reminder types, create function following formatReminderMessage()
// pattern. Accept data and config, return formatted string. Each reminder
// type should have its own formatting logic.

// formatReminderMessage builds display string for uncommitted work reminder
//
// What It Does:
// Replaces {count} placeholder in message template with actual count,
// adds icon and formatting based on configuration. Returns empty string
// if reminder disabled or below threshold.
//
// Parameters:
//   count: Number of uncommitted changes
//
// Returns:
//   string: Formatted reminder string, or empty if not applicable
//
// Example usage:
//
//	message := formatReminderMessage(5)
//	// Returns: "\n⚠️  Reminder: 5 uncommitted change(s) in workspace\n"
//
func formatReminderMessage(count int) string {
	// Check if reminders enabled
	if remindersConfigLoaded && remindersConfig != nil {
		if !remindersConfig.Display.Enabled || !remindersConfig.Reminders.UncommittedWork.Enabled {
			return "" // Reminders disabled
		}

		// Check threshold
		if count < remindersConfig.Reminders.UncommittedWork.Threshold {
			return "" // Below threshold
		}

		// Build message from config
		icon := remindersConfig.Reminders.UncommittedWork.Icon
		message := remindersConfig.Reminders.UncommittedWork.Message
		message = strings.ReplaceAll(message, "{count}", fmt.Sprintf("%d", count))

		// Add prefix newline if configured
		prefix := ""
		if remindersConfig.Display.PrefixNewline {
			prefix = "\n"
		}

		return fmt.Sprintf("%s%s  %s\n", prefix, icon, message)
	}

	// Fall back to defaults
	if count < defaultUncommittedThreshold {
		return "" // Below threshold
	}

	message := strings.ReplaceAll(defaultUncommittedMessage, "{count}", fmt.Sprintf("%d", count))
	prefix := ""
	if defaultPrefixNewline {
		prefix = "\n"
	}

	return fmt.Sprintf("%s%s  %s\n", prefix, defaultUncommittedIcon, message)
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface. Top rungs of
// the ladder - orchestrate helpers and core operations into complete
// functionality. Simple by design - complexity lives in helpers and core
// operations, Public APIs orchestrate proven pieces.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-public-apis.md

// RemindUncommittedWork checks for uncommitted changes and reminds user
//
// What It Does:
// Checks if workspace is git repository, retrieves uncommitted change count,
// formats and displays reminder if changes found above threshold. Non-blocking -
// failures are silent.
//
// Parameters:
//   workspace: Absolute path to workspace directory
//
// Health Impact:
//   Git repo checked successfully: +40
//   Non-git workspace (expected): +20
//   Git check fails: +0 (silent)
//
// Example usage:
//
//	func sessionEnd() {
//	    session.RemindUncommittedWork("/home/user/project")
//	    // Output: ⚠️  Reminder: 3 uncommitted change(s) in workspace
//	}
//
func RemindUncommittedWork(workspace string) {
	// Check if reminders enabled and get behavior settings
	displayEnabled := defaultDisplayEnabled
	checkGitOnly := defaultCheckGitOnly
	silentFailures := defaultSilentFailures

	if remindersConfigLoaded && remindersConfig != nil {
		displayEnabled = remindersConfig.Display.Enabled
		checkGitOnly = remindersConfig.Behavior.CheckGitOnly
		silentFailures = remindersConfig.Behavior.SilentFailures
	}

	// If display disabled, return early
	if !displayEnabled {
		return
	}

	// Check if workspace is git repository
	isGitRepo := git.IsGitRepository(workspace)

	// If checkGitOnly is true and not a git repo, return (with optional warning)
	if checkGitOnly && !isGitRepo {
		if !silentFailures {
			fmt.Printf("⚠️  Workspace is not a git repository: %s\n", workspace)
		}
		return
	}

	// If not a git repo, can't check uncommitted work
	if !isGitRepo {
		return
	}

	// Get repository info
	info := git.GetInfo(workspace)

	// Format and display reminder
	message := formatReminderMessage(info.UncommittedCount)
	if message != "" {
		fmt.Print(message)
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
//
// ────────────────────────────────────────────────────────────────
// Code Validation: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Validation section explanation, see: standards/code/4-block/sections/CWS-SECTION-010-CLOSING-code-validation.md
//
// Testing Requirements:
//   - Import the library without errors
//   - Call RemindUncommittedWork() with git and non-git workspaces
//   - Verify output format matches configuration settings
//   - Test with various uncommitted change counts (0, 1, many)
//   - Confirm threshold configuration works (only show above threshold)
//   - Ensure no go vet warnings introduced
//   - Test with missing configuration file (uses defaults)
//   - Test with invalid configuration file (uses defaults)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test at actual session end with uncommitted changes
//   - Test at actual session end with clean workspace
//   - Verify non-blocking behavior (session ends if reminder fails)
//   - Validate configuration customization works correctly
//
// Example validation code:
//
//     // Test basic functionality
//     session.RemindUncommittedWork("/home/user/project")  // Should show reminder if changes exist
//     // Make some changes
//     // Call again and verify reminder appears
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
// available. Configuration loads automatically during package import via init().
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func main() {
//         workspace := "/home/user/project"
//         session.RemindUncommittedWork(workspace)  // At session end
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Configuration: Loaded once in init(), cached for session lifetime
//   - Git operations: Delegated to system/lib/git (handles own cleanup)
//   - Memory: String allocations transiently, garbage collected after use
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Git library handles command cleanup automatically
//   - No persistent state requiring cleanup
//
// Error State Cleanup:
//   - Silent failures - no partial state to clean
//   - Git operations return immediately on error
//   - No resources leak on failure paths
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Transient allocations only (messages, strings)
//   - No manual memory management needed
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// For Library Overview section explanation, see: standards/code/4-block/sections/CWS-SECTION-013-CLOSING-library-overview.md
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Configurable reminder display at session end
//   - Uncommitted work detection via git integration
//   - Threshold-based triggering
//   - Graceful fallback to defaults if config unavailable
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list organized by category in typical usage order
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Rails/Ladder/Baton) explanation
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new reminder types (following formatReminderMessage pattern)
//   ✅ Add new message templates to configuration
//   ✅ Extend configuration with new display options
//   ✅ Add severity levels or categorization
//   ✅ Add health tracking/logging integration
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Public API function signatures - breaks all calling code
//   ⚠️ Configuration file structure - breaks user customizations
//   ⚠️ Message placeholder format ({count}) - breaks custom messages
//   ⚠️ Default behavior - affects all unconfigured installations
//   ⚠️ Git library integration - affects detection reliability
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking design principle
//   ❌ Configuration-driven architecture pattern
//   ❌ Silent failure behavior
//   ❌ Rails pattern for package-level state
//
// Validation After Modifications:
//   See "Code Validation" section in GROUP 1: CODING above for comprehensive
//   testing requirements, build verification, and integration testing procedures.
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart - Internal Structure" section above for
// complete ladder structure (dependencies) and baton flow (execution paths).
//
// The Organizational Chart in BODY provides the detailed map showing:
// - All functions and their dependencies (ladder)
// - Complete execution flow paths (baton)
// - APU count (Available Processing Units)
//
// Quick architectural summary:
// - 1 public API orchestrates 1 core operation using 2 helpers
// - Ladder: Public API → Git Library → Core Operations → Helpers
// - Baton: Entry → Git Check → Get Info → Format → Output → Exit
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Core Operations" subsection header comments above for detailed
// extension points. Each subsection includes "Extension Point" guidance showing:
// - Where to add new functionality
// - What naming pattern to follow
// - How to integrate with existing code
// - What tests to update
//
// Quick reference:
// - Adding reminder types: See BODY "Message Formatting" extension point
// - Adding configuration options: Update types in SETUP, add to reminders.jsonc
// - Customizing messages: Edit configuration file message templates
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Constants: Lightweight default values
// - Types: Minimal configuration structures
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary:
// - Most expensive operation: git.GetInfo() (~10-100ms depending on repo size)
// - Memory characteristics: Transient allocations only, <1KB per invocation
// - Key optimization: Git library handles performance (caching, efficient commands)
// - Configuration loading: One-time cost at package import
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: No reminder even when changes exist
//   - Check: Configuration file display.enabled setting
//   - Check: Configuration file reminders.uncommitted_work.enabled
//   - Check: Threshold setting (reminder only shows if count >= threshold)
//   - Check: Workspace is actually a git repository
//   - Solution: Enable in configuration or verify git repository status
//
// Problem: Wrong message format displayed
//   - Check: Configuration file message template
//   - Check: {count} placeholder exists in template
//   - Solution: Verify message template syntax in reminders.jsonc
//
// Problem: Configuration changes not taking effect
//   - Cause: Configuration loaded once at package import
//   - Solution: Restart session to reload configuration
//
// Problem: Threshold not working as expected
//   - Check: threshold value in configuration (0 = any changes)
//   - Check: Actual uncommitted count (use git status to verify)
//   - Solution: Adjust threshold or verify git status matches expectation
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard Library, system/lib/git
// - Dependents (What Uses This): session/cmd-end
// - Integration Points: Git library for repository status
//
// Quick summary:
// - Key dependencies: system/lib/git (repository status checking)
// - Primary consumers: Session end hook
//
// Parallel Implementation:
//   - No parallel implementations (Go-specific)
//   - Shared philosophy: Gentle awareness without nagging
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration-driven reminders - COMPLETED
//   ✓ Threshold-based triggering - COMPLETED
//   ⏳ Additional reminder types (long-running processes, unsaved files)
//   ⏳ File list display (show which files uncommitted)
//   ⏳ Severity levels (warning, info, critical)
//   ⏳ Color coding for different reminder types
//
// Research Areas:
//   - Interactive cleanup (offer to commit or stash)
//   - Smart threshold (based on session duration or file importance)
//   - Integration with session patterns (learn typical commit patterns)
//   - Notification persistence (remember reminders across sessions)
//   - Custom reminder hooks (user-defined checks)
//
// Integration Targets:
//   - Health scoring system (track reminder effectiveness)
//   - Session logging (record reminder history)
//   - Pattern learning (identify habitual uncommitted work)
//   - Display customization (color themes, icon sets)
//
// Known Limitations to Address:
//   - Single reminder type (only uncommitted work)
//   - No file-level detail (just count)
//   - No interactive actions (read-only awareness)
//   - Configuration requires session restart
//   - No reminder grouping (each type displayed separately)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-12) - Configuration-driven reminders
//         - Added reminders.jsonc configuration file support
//         - Message template customization ({count} placeholder)
//         - Threshold-based triggering
//         - Display customization (icon, prefix newline)
//         - Graceful fallback to defaults if config missing
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Hardcoded message format
//         - Basic uncommitted work detection
//         - Git integration via system/lib/git
//         - Silent failure design established
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing gentle reminders at session end.
// Non-intrusive by design - awareness serves memory, never nags.
//
// Modify thoughtfully - changes here affect session end experience for all users.
// Non-blocking behavior and gentle awareness are architectural guarantees that
// must be maintained.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go build && go vet)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Verify non-blocking behavior preserved
//
// "Remember therefore from whence thou art fallen, and repent" - Revelation 2:5 (KJV)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//   import "hooks/lib/session"
//   session.RemindUncommittedWork("/path/to/workspace")  // Check and remind
//
// Configuration Customization:
//   Edit ~/.claude/cpi-si/system/data/config/session/reminders.jsonc
//   Customize message template (use {count} placeholder)
//   Adjust threshold (minimum changes to trigger)
//   Change icon or display settings
//
// Disabling Reminders:
//   Set display.enabled: false in configuration file
//   Or set reminders.uncommitted_work.enabled: false
//
// ============================================================================
// END CLOSING
// ============================================================================
