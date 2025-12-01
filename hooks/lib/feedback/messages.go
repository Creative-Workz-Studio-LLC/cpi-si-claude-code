// METADATA
//
// Feedback Messages Library - CPI-SI Hook Support System
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Let your speech always be with grace, seasoned with salt, that you may know how you ought to answer each one." - Colossians 4:6 (WEB)
// Principle: Communication should edify, encourage, and guide toward excellence
// Anchor: "Therefore exhort one another, and build each other up" - 1 Thessalonians 5:11 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Hook support (mid-rung on ladder)
// Role: Provides contextual feedback and encouragement after tool operations
// Paradigm: Kingdom-honoring communication through gracious, purposeful feedback
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-11 - Refactored to use system display library
//
// Version History:
//   2.0.0 (2025-11-11) - Refactored to use system/lib/display for formatting
//   1.0.0 (2024-10-24) - Initial implementation with raw formatting
//
// Purpose & Function
//
// Purpose: Provide contextual guidance, encouragement, and quality reminders after tool operations.
// Helps maintain quality standards and Kingdom-honoring development practices through timely feedback.
//
// Core Design: Event-driven feedback system that responds to specific tool operations
// (commits, builds, dependency installs) with encouragement and quality checklists.
//
// Key Features:
//   - Post-commit encouragement (meaningful commit messages)
//   - Dependency install reminders (version verification)
//   - Build completion quality checklist (warnings, testing, excellence)
//   - Command failure warnings (non-zero exit codes)
//   - Build command detection (pattern matching)
//
// Philosophy: Communication should serve others - provide guidance that helps maintain
// quality and reminds of Kingdom principles without being burdensome or preachy.
//
// Blocking Status
//
// Non-blocking: All feedback operations print to stdout only - failures don't interrupt workflow.
// Mitigation: Uses system display library with panic recovery and health tracking.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/feedback"
//
//	// After git commit
//	feedback.OnCommitCreated()
//
//	// After dependency install
//	feedback.OnDependenciesInstalled()
//
//	// After build completes
//	feedback.OnBuildComplete()
//
//	// On command failure
//	feedback.OnCommandFailure("1")
//
// Integration Pattern:
//   1. Import feedback library
//   2. Call appropriate feedback function from hooks (tool/post-use, prompt/submit, etc.)
//   3. No cleanup needed - feedback is fire-and-forget
//
// Public API (in typical usage order):
//
//   Feedback Messages (contextual guidance):
//     OnCommitCreated() - Encourages meaningful commit messages
//     OnDependenciesInstalled() - Reminds to verify versions
//     OnBuildComplete() - Shows quality checklist
//     OnCommandFailure(exitCode string) - Warns about failures
//
//   Utility Functions (pattern detection):
//     IsBuildCommand(cmd string) bool - Detects build operations
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, strings
//   External: None
//   Internal: system/lib/display (formatting and color output)
//
// Dependents (What Uses This):
//   Hooks: tool/post-use, prompt/submit (contextual feedback after operations)
//
// Integration Points:
//   - Ladder: Mid-rung - uses display lib (lower), used by hooks (higher)
//   - Baton: Receives event notifications, outputs formatted feedback
//   - Rails: Uses system display library's logging infrastructure
//
// Health Scoring
//
// Health Scoring Map (Base100):
//
//   Feedback Operations (Total = 100):
//     All operations delegated to system/lib/display which tracks:
//       - Display formatting success: +100 (from display library)
//       - Panic recovery: -10 (from display library)
//
// Note: This library itself doesn't track health - it relies on system/lib/display's
// health tracking. Feedback operations are simple pass-through to display formatters.
package feedback

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
// functionality. Each import commented with purpose, not just name.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"encoding/json" // JSON parsing for configuration files
	"fmt"           // String formatting for composed messages
	"os"            // File operations for config loading
	"path/filepath" // Path manipulation for config file locations
	"strings"       // String operations for command pattern matching

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display" // ANSI color formatting and message display
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Named values that never change. Magic numbers given meaningful names,
// configuration values documented with reasoning. Constants prevent bugs
// from typos and make intent clear - timeout duration has purpose, not
// just "30" scattered through code.
//
// See: standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

// None needed - feedback messages are context-specific without hardcoded values

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Configuration Types ---
// Structures for loading feedback configuration from data files

// BuildPatternsConfig holds build command detection patterns
type BuildPatternsConfig struct {
	Patterns map[string]struct {
		Patterns    []string `json:"patterns"`     // List of command patterns to match
		Description string   `json:"description"`  // What these patterns detect
	} `json:"patterns"`  // Organized by language/tool category
}

// MessagesConfig holds feedback message templates
type MessagesConfig struct {
	Commit struct {
		Success struct {
			Message     string `json:"message"`       // Message text
			DisplayType string `json:"display_type"`  // success/warning/info/failure
		} `json:"success"`
	} `json:"commit"`

	Dependencies struct {
		Installed struct {
			Message     string `json:"message"`
			DisplayType string `json:"display_type"`
		} `json:"installed"`
	} `json:"dependencies"`

	Build struct {
		Complete struct {
			Header      string `json:"header"`
			DisplayType string `json:"display_type"`
		} `json:"complete"`
	} `json:"build"`

	CommandFailure struct {
		Generic struct {
			Message     string `json:"message"`
			DisplayType string `json:"display_type"`
		} `json:"generic"`
	} `json:"command_failure"`
}

// QualityChecklistConfig holds quality checklist items
type QualityChecklistConfig struct {
	Checklist struct {
		Technical []struct {
			Item        string `json:"item"`         // Checklist item text
			Category    string `json:"category"`     // Item category
			Priority    string `json:"priority"`     // high/medium/low/essential
			Explanation string `json:"explanation"`  // Detailed explanation
		} `json:"technical"`

		Kingdom []struct {
			Item        string `json:"item"`
			Category    string `json:"category"`
			Priority    string `json:"priority"`
			Explanation string `json:"explanation"`
		} `json:"kingdom"`
	} `json:"checklist"`

	Filtering struct {
		Modes struct {
			Default struct {
				IncludeCategories []string `json:"include_categories"`  // Categories to show
			} `json:"default"`
			Comprehensive struct {
				IncludeCategories []string `json:"include_categories"`  // Categories for comprehensive mode
			} `json:"comprehensive"`
			Quick struct {
				IncludeCategories []string `json:"include_categories"`  // Categories for quick mode
			} `json:"quick"`
			KingdomFocused struct {
				IncludeCategories []string `json:"include_categories"`  // Categories for kingdom_focused mode
			} `json:"kingdom_focused"`
		} `json:"modes"`
		ActiveMode string `json:"active_mode"`  // Which mode is active
	} `json:"filtering"`
}

// ────────────────────────────────────────────────────────────────
// Package-Level State
// ────────────────────────────────────────────────────────────────
// Package-level variables and initialization. Rails pattern: each component
// creates its own logger/inspector independently - never passed as parameters.
//
// See: standards/code/4-block/sections/CWS-SECTION-005-SETUP-package-state.md

// Configuration loaded from system data files
var (
	buildPatterns    *BuildPatternsConfig    // Build command detection patterns
	messages         *MessagesConfig         // Feedback message templates
	qualityChecklist *QualityChecklistConfig // Quality checklist items
	configLoaded     bool                    // Whether configs loaded successfully
)

// init loads configuration files at package initialization
// Falls back to hardcoded defaults if configs unavailable
func init() {
	// Get home directory for config file paths - needed to locate system data
	home := os.Getenv("HOME")  // Read HOME environment variable
	if home == "" {  // Check if HOME is empty - need fallback
		var err error
		home, err = os.UserHomeDir()  // Try system API as fallback
		if err != nil {  // Check if fallback failed
			configLoaded = false  // Mark configs as unavailable - will use hardcoded fallbacks
			return  // Exit init early - can't load configs without home directory
		}
	}

	// Build paths to configuration files - all in system/data/config/feedback/
	configDir := filepath.Join(home, ".claude/cpi-si/system/data/config/feedback")
	patternsPath := filepath.Join(configDir, "build-patterns.jsonc")  // Build command patterns
	messagesPath := filepath.Join(configDir, "messages.jsonc")        // Feedback messages
	checklistPath := filepath.Join(configDir, "quality-checklist.jsonc")  // Quality checklist

	// Load build patterns config - used by IsBuildCommand()
	buildPatterns = loadBuildPatterns(patternsPath)  // Parse JSON from file

	// Load messages config - used by feedback functions
	messages = loadMessages(messagesPath)  // Parse JSON from file

	// Load quality checklist config - used by OnBuildComplete()
	qualityChecklist = loadQualityChecklist(checklistPath)  // Parse JSON from file

	// Mark configs as loaded if all succeeded - enables config-driven behavior
	configLoaded = (buildPatterns != nil && messages != nil && qualityChecklist != nil)
}

// ============================================================================
// BODY
// ============================================================================
//
// For BODY structure explanation, see: standards/code/4-block/CWS-STD-007-CODE-body-block.md

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Ladder Structure
// ────────────────────────────────────────────────────────────────
// Shows component hierarchy and dependencies (who uses whom).
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Orchestration)
//   ├── OnCommitCreated() → uses display.Success()
//   ├── OnDependenciesInstalled() → uses display.Success()
//   ├── OnBuildComplete() → uses display.Info()
//   └── OnCommandFailure() → uses display.Warning()
//
//   Utility Functions (Middle Rungs - Business Logic)
//   └── IsBuildCommand() → uses strings.Contains()
//
// APUs (Available Processing Units):
// - 5 functions total
// - 4 public feedback APIs (exported interface)
// - 1 utility function (pattern detection)
//
// All operations delegate to system/lib/display for formatting - this library
// provides contextual feedback wrappers around display primitives.
//
// See: standards/code/4-block/sections/CWS-SECTION-006-BODY-organizational-chart.md

// ────────────────────────────────────────────────────────────────
// Helpers - Configuration Loading
// ────────────────────────────────────────────────────────────────
// Functions for loading configuration from JSON files with graceful fallback

// loadBuildPatterns loads build command patterns from config file
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadBuildPatterns(path string) *BuildPatternsConfig {
	data, err := os.ReadFile(path)  // Read config file content
	if err != nil {  // Check if read failed - file missing or unreadable
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments by removing lines starting with //
	lines := strings.Split(string(data), "\n")  // Split into individual lines
	var jsonLines []string  // Collect non-comment lines
	for _, line := range lines {  // Process each line
		trimmed := strings.TrimSpace(line)  // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {  // Check if not a comment line
			jsonLines = append(jsonLines, line)  // Keep this line
		}
	}
	jsonData := strings.Join(jsonLines, "\n")  // Rejoin into valid JSON

	var config BuildPatternsConfig  // Create config structure
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {  // Parse JSON into structure
		return nil  // Return nil if parsing failed - malformed JSON
	}

	return &config  // Return successfully loaded config
}

// loadMessages loads feedback message templates from config file
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadMessages(path string) *MessagesConfig {
	data, err := os.ReadFile(path)  // Read config file content
	if err != nil {  // Check if read failed
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments
	lines := strings.Split(string(data), "\n")
	var jsonLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			jsonLines = append(jsonLines, line)
		}
	}
	jsonData := strings.Join(jsonLines, "\n")

	var config MessagesConfig
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
		return nil
	}

	return &config
}

// loadQualityChecklist loads quality checklist items from config file
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadQualityChecklist(path string) *QualityChecklistConfig {
	data, err := os.ReadFile(path)  // Read config file content
	if err != nil {  // Check if read failed
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments
	lines := strings.Split(string(data), "\n")
	var jsonLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			jsonLines = append(jsonLines, line)
		}
	}
	jsonData := strings.Join(jsonLines, "\n")

	var config QualityChecklistConfig
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
		return nil
	}

	return &config
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Feedback Messages
// ────────────────────────────────────────────────────────────────
// Exported functions providing contextual feedback after tool operations.
// Each function delegates to system/lib/display for actual formatting.
//
// See: standards/code/4-block/sections/CWS-SECTION-008-BODY-public-apis.md

// OnCommitCreated provides feedback after git commit
// Encourages meaningful commit messages that tell a story
func OnCommitCreated() {
	// Try to use config message, fall back to hardcoded if unavailable
	message := "Commit created. Commits tell a story - make it meaningful."  // Default fallback
	if configLoaded && messages != nil {  // Check if config loaded successfully
		message = messages.Commit.Success.Message  // Use message from config file
	}

	// Delegate to display.Success for formatted output - adds checkmark icon and green color
	fmt.Println(display.Success(message))
}

// OnDependenciesInstalled reminds to verify versions after package install
// Encourages version verification for reproducible builds
func OnDependenciesInstalled() {
	// Try to use config message, fall back to hardcoded if unavailable
	message := "Dependencies installed. Consider verifying versions."  // Default fallback
	if configLoaded && messages != nil {  // Check if config loaded successfully
		message = messages.Dependencies.Installed.Message  // Use message from config file
	}

	// Delegate to display.Success for formatted output - adds checkmark icon and green color
	fmt.Println(display.Success(message))
}

// OnBuildComplete shows quality checklist
// Presents quality checklist (configurable or default 3-item checklist)
func OnBuildComplete() {
	// Try to use config header, fall back to hardcoded if unavailable
	header := "Build complete. Quality check:"  // Default fallback
	if configLoaded && messages != nil {  // Check if config loaded successfully
		header = messages.Build.Complete.Header  // Use header from config file
	}

	// Display header with info formatting - adds info icon and cyan color
	fmt.Println(display.Info(header))

	// Try to use config checklist, fall back to hardcoded if unavailable
	if configLoaded && qualityChecklist != nil {  // Check if config loaded successfully
		// Load checklist items from config filtered by active mode
		displayChecklistFromConfig()  // Use config-driven checklist
	} else {
		// Fallback to hardcoded 3-item checklist - matches original behavior
		fmt.Println("   • Compiles without warnings?")       // Technical quality
		fmt.Println("   • Edge cases tested?")               // Testing coverage
		fmt.Println("   • Honors God through excellence?")   // Kingdom principle
	}
}

// displayChecklistFromConfig renders quality checklist from loaded configuration
// Filters items based on active mode's include_categories
func displayChecklistFromConfig() {
	// Get active filtering mode - determines which categories to show
	activeMode := qualityChecklist.Filtering.ActiveMode  // e.g., "default", "comprehensive"
	var includeCategories []string  // Categories to display

	// Load category filter for active mode - switch on mode name to load correct config
	switch activeMode {
	case "comprehensive":  // Extended checklist for thorough review
		includeCategories = qualityChecklist.Filtering.Modes.Comprehensive.IncludeCategories
	case "quick":  // Minimal checklist for rapid iteration
		includeCategories = qualityChecklist.Filtering.Modes.Quick.IncludeCategories
	case "kingdom_focused":  // Focus on Kingdom principles only
		includeCategories = qualityChecklist.Filtering.Modes.KingdomFocused.IncludeCategories
	case "default", "":  // Standard 3-item checklist (or empty mode string)
		includeCategories = qualityChecklist.Filtering.Modes.Default.IncludeCategories
	default:  // Unknown mode - fall back to default
		includeCategories = qualityChecklist.Filtering.Modes.Default.IncludeCategories
	}

	// Create set of included categories for fast lookup
	categorySet := make(map[string]bool)  // Map for O(1) category checking
	for _, cat := range includeCategories {  // Add each category to set
		categorySet[cat] = true
	}

	// Display technical items that match included categories
	for _, item := range qualityChecklist.Checklist.Technical {  // Iterate technical checklist
		if categorySet[item.Category] {  // Check if this item's category is included
			fmt.Printf("   • %s\n", item.Item)  // Print checklist item with bullet
		}
	}

	// Display kingdom items that match included categories
	for _, item := range qualityChecklist.Checklist.Kingdom {  // Iterate kingdom checklist
		if categorySet[item.Category] {  // Check if this item's category is included
			fmt.Printf("   • %s\n", item.Item)  // Print checklist item with bullet
		}
	}
}

// OnCommandFailure warns about non-zero exit code
// Alerts to command failures that might otherwise go unnoticed
func OnCommandFailure(exitCode string) {
	if exitCode != "" && exitCode != "0" {  // Check if exit code indicates failure - empty or "0" means success
		// Try to use config message template, fall back to hardcoded if unavailable
		messageTemplate := "Command exited with code {exit_code}"  // Default fallback with placeholder
		if configLoaded && messages != nil {  // Check if config loaded successfully
			messageTemplate = messages.CommandFailure.Generic.Message  // Use template from config file
		}

		// Replace {exit_code} placeholder with actual exit code value
		message := strings.Replace(messageTemplate, "{exit_code}", exitCode, -1)  // Replace all occurrences

		// Format warning message - display.Warning adds warning icon and yellow color
		fmt.Println(display.Warning(message))
	}
	// Note: If exitCode is "0" or empty, do nothing - command succeeded, no warning needed
}

// ────────────────────────────────────────────────────────────────
// Utility Functions - Pattern Detection
// ────────────────────────────────────────────────────────────────
// Helper functions for detecting specific command patterns.
//
// See: standards/code/4-block/sections/CWS-SECTION-007-BODY-helpers.md

// IsBuildCommand checks if command string is a build operation
// Used to trigger build-specific feedback (quality checklist)
// Supports configuration-driven patterns or hardcoded fallback
func IsBuildCommand(cmd string) bool {
	var patterns []string  // Collect all patterns to check

	// Try to use config patterns, fall back to hardcoded if unavailable
	if configLoaded && buildPatterns != nil {  // Check if config loaded successfully
		// Collect patterns from all language/tool categories in config
		for _, langConfig := range buildPatterns.Patterns {  // Iterate through each language/tool category
			patterns = append(patterns, langConfig.Patterns...)  // Add all patterns from this category
		}
	} else {
		// Fallback to hardcoded patterns - matches original behavior
		patterns = []string{
			"cargo build",    // Rust build system
			"go build",       // Go compiler
			"npm run build",  // Node.js build script
			"make build",     // Make build target
			"gcc",            // GNU C compiler
			"g++",            // GNU C++ compiler
			"clang",          // LLVM C/C++ compiler
		}
	}

	// Check if command contains any build pattern - simple substring matching
	for _, pattern := range patterns {  // Iterate through each pattern - testing against command string
		if strings.Contains(cmd, pattern) {  // Check if pattern appears in command - case-sensitive match
			return true  // Found a match - this is a build command, return immediately
		}
	}

	return false  // No patterns matched - not a recognized build command
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
//   - Call feedback functions and verify output format
//   - Test IsBuildCommand() with various command strings
//   - Verify integration with system/lib/display functions
//   - Check hooks calling feedback functions correctly
//
// Example validation code:
//
//     // Test import
//     import "hooks/lib/feedback"
//
//     // Test feedback functions
//     feedback.OnCommitCreated()
//     // Expected: Green checkmark + "Commit created. Commits tell a story - make it meaningful."
//
//     feedback.OnBuildComplete()
//     // Expected: Cyan info icon + "Build complete. Quality check:"
//     //           Followed by three checklist items
//
//     // Test build detection
//     if !feedback.IsBuildCommand("go build main.go") {
//         t.Error("IsBuildCommand should detect 'go build'")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by hook components.
//
// Usage: import "hooks/lib/feedback"
//
// The library is imported into hook packages (tool/post-use, prompt/submit), making
// all exported functions available. No code executes during import - functions are
// defined and ready for hooks to call after tool operations.
//
// Example usage from hooks:
//
//     package toolpostuse
//
//     import "hooks/lib/feedback"
//
//     func AfterToolUse(toolName string) {
//         // Check if build command, show quality checklist
//         if feedback.IsBuildCommand(toolName) {
//             feedback.OnBuildComplete()
//         }
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No file handles: All operations use display library (handles cleanup internally)
//   - No network connections: Pure local output formatting
//   - No database connections: No persistence layer
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code (hooks) responsible for managing their own lifecycle
//   - This library is stateless - no cleanup function needed
//
// Error State Cleanup:
//   - Display library handles panic recovery (see system/lib/display format.go)
//   - No partial state to corrupt (stateless operations)
//   - Failures simply result in no output (graceful degradation)
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - No large allocations (simple string formatting)
//   - Build pattern array is small constant data (7 strings)

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// For Library Overview section explanation, see: standards/code/4-block/sections/CWS-SECTION-013-CLOSING-library-overview.md
//
// Purpose: See METADATA "Purpose & Function" section above - contextual guidance
// and quality reminders after tool operations
//
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Gracious feedback after commits, builds, dependency installs
//   - Quality checklists emphasizing Kingdom-honoring excellence
//   - Build command detection for triggering appropriate feedback
//   - Non-blocking operation (failures don't interrupt workflow)
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide. Simple import and call pattern.
//
// Public API: See METADATA "Public API" section above for complete function list
// organized by purpose (feedback messages, utility functions)
//
// Architecture: See METADATA "CPI-SI Identity" section above - this is a mid-rung
// library using display (lower) and used by hooks (higher)
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new feedback functions (OnTestComplete, OnDeploySuccess, etc.)
//   ✅ Add new build patterns to IsBuildCommand (cmake, bazel, etc.)
//   ✅ Extend quality checklist items in OnBuildComplete
//   ✅ Add new utility functions for pattern detection
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Function signatures - all hook code depends on these interfaces
//   ⚠️ Message content that other systems parse (if any)
//   ⚠️ Display library integration approach
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Delegation to system/lib/display (architectural pattern)
//   ❌ Non-blocking behavior (feedback must never crash workflow)
//   ❌ Kingdom-honoring communication principles
//
// Validation After Modifications:
//   See "Code Validation" section above for testing requirements and
//   example validation code
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart - Ladder Structure" section above for
// complete ladder structure (dependencies) and baton flow (function calls).
//
// Quick architectural summary:
// - 4 public feedback APIs orchestrate display library formatters
// - 1 utility function provides build command detection
// - Ladder: hooks (top) → feedback (mid) → display (bottom)
// - Baton: hooks call feedback functions → feedback calls display formatters → output to stdout
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Public APIs" and "Utility Functions" section headers above for
// detailed extension points.
//
// Quick reference:
// - Adding feedback functions: Add to "Public APIs - Feedback Messages" section,
//   follow existing pattern (delegate to display library formatters)
// - Adding command patterns: Extend buildPatterns array in IsBuildCommand()
// - Adding utility functions: Add to "Utility Functions - Pattern Detection" section
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// All operations are lightweight string formatting with minimal overhead.
//
// - Most expensive operation: IsBuildCommand() - O(n) where n = pattern count (currently 7)
// - Memory characteristics: Minimal (small constant strings, no allocations)
// - Output cost: Display library handles formatting with panic recovery
//
// Note: Performance is not a concern for this library - it runs after tool operations
// where feedback latency is imperceptible to users.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: No feedback appears after tool operations
//   Check: Hook integration - verify hooks import and call feedback functions
//   Check: Display library working - test system/lib/display directly
//   Solution: Review hook implementation, ensure feedback functions called
//
// Problem: Build feedback not appearing after build commands
//   Check: IsBuildCommand() detects your build command - test pattern matching
//   Check: Hook calls IsBuildCommand() before conditionally calling OnBuildComplete()
//   Solution: Add new build pattern to IsBuildCommand() or fix hook logic
//
// Problem: Formatting not correct (no colors/icons)
//   Expected: Display library provides ANSI formatting
//   Note: If terminal doesn't support ANSI, formatting may not display correctly
//   Solution: Verify display library working, check terminal ANSI support
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependency: system/lib/display (provides all formatting primitives)
// - Primary consumers: hooks/tool/post-use, hooks/prompt/submit (contextual feedback)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Refactored to use system display library - COMPLETED (v2.0.0)
//   ⏳ Add feedback for test operations (OnTestComplete, OnTestFailure)
//   ⏳ Add feedback for deployment operations (OnDeploySuccess, OnDeployFailure)
//   ⏳ Pattern-based feedback (detect specific tools and provide targeted guidance)
//
// Research Areas:
//   - Context-aware feedback (analyze command output for specific issues)
//   - Adaptive encouragement (vary messages to avoid monotony)
//   - Integration with activity logging (track which feedback resonates)
//
// Integration Targets:
//   - Session quality tracking (correlate feedback with outcome quality)
//   - Pattern learning (discover what feedback helps most)
//
// Known Limitations to Address:
//   - Build detection limited to hardcoded patterns (could be more flexible)
//   - No feedback personalization (same messages for everyone)
//   - No context from command output (only command string itself)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history:
//
//   2.0.0 (2025-11-11) - Architectural Refactoring
//         - Refactored to delegate all formatting to system/lib/display
//         - Removed raw fmt.Println with manual formatting
//         - Now uses display.Success(), display.Info(), display.Warning()
//         - Maintains same external API (hooks see no breaking changes)
//         - Comprehensive template alignment (METADATA, SETUP, BODY, CLOSING)
//
//   1.0.0 (2024-10-24) - Initial Implementation
//         - Basic feedback functions with manual formatting
//         - IsBuildCommand() pattern detection
//         - Post-commit, post-build, post-dependency feedback
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a MID-RUNG on the ladder - it uses display library (lower rung)
// and is used by hooks (higher rung). It provides contextual feedback that
// encourages quality and Kingdom-honoring development practices.
//
// Modify thoughtfully - hook code depends on these interfaces. Maintain the
// gracious, encouraging tone that edifies rather than condemns.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test integration with hooks after changes
//   - Maintain Kingdom-honoring communication principles
//
// "Let your speech always be with grace, seasoned with salt, that you may know
// how you ought to answer each one." - Colossians 4:6 (WEB)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage from Hooks:
//
//     import "hooks/lib/feedback"
//
//     // After git commit
//     feedback.OnCommitCreated()
//
//     // After build completes
//     feedback.OnBuildComplete()
//
//     // After dependency install
//     feedback.OnDependenciesInstalled()
//
//     // On command failure
//     feedback.OnCommandFailure("1")  // exit code as string
//
// Conditional Feedback Based on Command:
//
//     cmd := "go build main.go"
//     if feedback.IsBuildCommand(cmd) {
//         feedback.OnBuildComplete()
//     }
//
// Pattern Detection:
//
//     commands := []string{
//         "cargo build --release",  // true - Rust build
//         "go build",                // true - Go build
//         "npm run build",           // true - Node.js build
//         "make clean",              // false - not a build command
//     }
//
//     for _, cmd := range commands {
//         if feedback.IsBuildCommand(cmd) {
//             fmt.Printf("%s is a build command\n", cmd)
//         }
//     }

// ============================================================================
// END CLOSING
// ============================================================================
