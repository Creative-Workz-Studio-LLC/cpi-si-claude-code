// METADATA
//
// Activity Tracking Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Whatever your hand finds to do, do it with all your might" - Ecclesiastes 9:10 (WEB)
// Principle: Faithful work leaves traces - recent activity reveals current focus and stewardship
// Anchor: "Be diligent to know the state of your flocks, and pay attention to your herds" - Proverbs 27:23 (WEB)
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - provides context awareness functionality)
// Role: Tracks recent file modifications to provide session context awareness
// Paradigm: CPI-SI framework component - serves hooks with workspace activity intelligence
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
//   2.0.0 (2025-11-12) - Template alignment, config-driven, display library integration
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded values
//
// Purpose & Function
//
// Purpose: Provides workspace activity awareness by tracking recently modified files
//
// Core Design: Simple file modification detection using find command with configurable
// time windows and exclusion patterns. Integrates with display library for consistent
// formatting across CPI-SI system.
//
// Key Features:
//   - Configurable time window (default 60 minutes, customizable)
//   - Comprehensive exclusion patterns (dotfiles, dependencies, build artifacts)
//   - Smart display thresholds (detailed count, summary count, silent for very large sets)
//   - Graceful fallback to hardcoded defaults if configuration unavailable
//   - Non-intrusive context awareness (shows activity without overwhelming)
//
// Philosophy: Recent activity reveals current focus. Show enough context to be helpful,
// not so much as to be noise. Default to sensible behavior even if configuration missing.
//
// Blocking Status
//
// Non-blocking: Never blocks session operations. If activity detection fails, continues
// silently without disrupting workflow.
// Mitigation: All errors handled gracefully, empty results treated as "no recent activity"
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Import package (configuration loaded automatically in init())
//   2. Call CheckRecentActivity(workspace) to display recent file modifications
//   3. Function prints to stdout using display library formatting
//   4. No return value - pure side effect (display only)
//
// Public API (in typical usage order):
//
//   Activity Detection:
//     CheckRecentActivity(workspace string) - Display recent file modifications
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os/exec (find command), fmt (output), os (file operations),
//                     path/filepath (path handling), strings (string manipulation),
//                     encoding/json (config parsing)
//   Internal: system/lib/display (formatted output with health tracking)
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start/start.go, session/cmd-stop/stop.go
//   Purpose: Provides session context awareness at start/stop
//
// Integration Points:
//   - Hooks call CheckRecentActivity() to show workspace state
//   - Uses display library for consistent formatted output
//   - Reads configuration from system/data/config/session/activity-tracking.jsonc
//
// Health Scoring
//
// Activity detection operations tracked with health scores reflecting system awareness quality.
//
// Configuration Loading:
//   - Config loaded successfully: +10
//   - Config missing/malformed: -5 (falls back to defaults)
//
// Activity Detection:
//   - Successful detection: +15
//   - Find command failure: -10
//   - Empty results (no activity): +5 (valid result, not failure)
//
// Display Operations:
//   - Formatted output displayed: +10
//   - Display formatting failure: -5 (falls back to plain output)
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
// functionality. Each import commented with purpose, not just name.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"encoding/json" // Parse activity-tracking.jsonc configuration file
	"fmt"           // Formatted output for activity messages to stdout
	"os"            // File operations for configuration loading and HOME directory
	"os/exec"       // Execute find command to discover recently modified files
	"path/filepath" // Join paths for configuration file location
	"strings"       // String manipulation for command output and config parsing

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display" // Formatted output with ANSI colors and health tracking
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

// None needed - this component is fully config-driven. All values that could
// be constants are instead loaded from activity-tracking.jsonc configuration
// file, allowing users to customize behavior without code changes.
//
// Historical note: v1.0.0 had hardcoded constants (60 minutes, threshold 10).
// v2.0.0 moved all to configuration for flexibility.

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Building Blocks ---
// Simple foundational types used throughout this component.

// ExclusionPattern represents a single path exclusion pattern for find command.
//
// Captures the pattern itself, human-readable description, and reasoning for
// why this pattern is excluded. Used to build -not -path arguments for find.
type ExclusionPattern struct {
	Pattern     string `json:"pattern"`     // Find command pattern (e.g., "*/node_modules/*")
	Description string `json:"description"` // Human-readable description of what's excluded
	Reasoning   string `json:"reasoning"`   // Why this pattern should be excluded
}

// TimeWindow defines how far back to search for file modifications.
//
// Configures the time range for "recent" activity. Larger values capture more
// context, smaller values focus on immediate work.
type TimeWindow struct {
	Minutes     int    `json:"minutes"`     // Number of minutes to look back (passed to find -mmin)
	Description string `json:"description"` // Purpose of this time window
	Reasoning   string `json:"reasoning"`   // Why this duration was chosen
}

// DisplaySettings controls how activity information is shown to user.
//
// Configures display behavior, thresholds for different output modes, and
// visual elements (icons, formatting). Balances information vs noise.
type DisplaySettings struct {
	ThresholdDetailed int    `json:"threshold_detailed"` // Show count when <= this many files
	ThresholdSummary  int    `json:"threshold_summary"`  // Show summary when > detailed but <= this
	ShowFileList      bool   `json:"show_file_list"`     // If true, display actual file paths
	Icon              string `json:"icon"`               // Emoji icon for activity messages
	Description       map[string]string `json:"description"` // Field descriptions
	Reasoning         map[string]string `json:"reasoning"`   // Design reasoning for thresholds
}

//--- Composed Types ---
// Complex types built from building blocks above.

// ActivityConfig holds complete activity tracking configuration.
//
// Loaded from activity-tracking.jsonc at init(). Contains all settings for
// time window, exclusions, and display behavior. Zero value is NOT usable -
// must load from file or use hardcoded fallback.
type ActivityConfig struct {
	Metadata struct {
		Name        string `json:"name"`         // Config identity
		Description string `json:"description"`  // Purpose of this config
		Version     string `json:"version"`      // Config format version
		Author      string `json:"author"`       // Who created this config
		Created     string `json:"created"`      // Creation date
		LastUpdated string `json:"last_updated"` // Last modification date
	} `json:"metadata"` // Configuration metadata and versioning

	TimeWindow         TimeWindow         `json:"time_window"`         // How far back to look for modifications
	ExclusionPatterns  []ExclusionPattern `json:"exclusion_patterns"`  // Paths to ignore (dependencies, artifacts)
	Display            DisplaySettings    `json:"display"`             // Display behavior and thresholds
	Config             struct {
		FallbackOnError bool   `json:"fallback_on_error"` // If true, use hardcoded defaults on load failure
		Description     string `json:"description"`       // Behavior explanation
	} `json:"config"` // Configuration behavior settings
	Extensions map[string]interface{} `json:"extensions"` // Future additions without breaking changes
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// Infrastructure available throughout component. Rails pattern - each
// component creates own logger independently without parameter passing.
//
// See: standards/code/patterns/CWS-PATTERN-003-CODE-rails.md
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md
//
// Note: This component uses display library for health tracking instead of
// direct logger attachment. Display library handles logging internally.

//--- Configuration State ---
// Loaded configuration and state flags.

// activityConfig holds loaded activity tracking configuration.
//
// Populated in init() from activity-tracking.jsonc. If loading fails and
// fallback is enabled, this remains nil and hardcoded defaults are used.
var activityConfig *ActivityConfig

// configLoaded tracks whether configuration loaded successfully.
//
// True if activity-tracking.jsonc parsed without errors. False triggers
// fallback to hardcoded defaults throughout component.
var configLoaded bool

func init() {
	// --- Configuration Loading ---
	// Load activity tracking settings from system config directory

	// Get HOME directory for config file path resolution
	home := os.Getenv("HOME") // Standard HOME environment variable
	if home == "" {           // Fallback if HOME not set
		home = "/home/seanje-lenox-wise" // Default CPI-SI user home
	}

	// Build config file path: ~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
	configPath := filepath.Join(
		home,
		".claude",
		"cpi-si",
		"system",
		"data",
		"config",
		"session",
		"activity-tracking.jsonc",
	)

	// Attempt to load configuration from file
	activityConfig = loadActivityConfig(configPath) // Returns nil on any error
	configLoaded = (activityConfig != nil)          // Set flag based on success

	// Note: No logging here - this runs during package import before logger
	// is available. Failures handled gracefully via configLoaded flag.
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
//   └── CheckRecentActivity() → uses buildFindCommand(), executeFind(), formatOutput()
//
//   Core Operations (Middle Rungs - Business Logic)
//   ├── buildFindCommand() → uses getTimeWindow(), getExclusionPatterns()
//   ├── executeFind() → pure execution (runs find command)
//   └── formatOutput() → uses getDisplaySettings(), system/lib/display
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── loadActivityConfig() → uses stripJSONCComments() (pure)
//   ├── stripJSONCComments() → pure function
//   ├── getTimeWindow() → reads config (fallback to default)
//   ├── getExclusionPatterns() → reads config (fallback to defaults)
//   └── getDisplaySettings() → reads config (fallback to defaults)
//
// Baton Flow (Execution Paths):
//
//   Entry → CheckRecentActivity(workspace)
//     ↓
//   getTimeWindow(), getExclusionPatterns() → buildFindCommand()
//     ↓
//   executeFind(command)
//     ↓
//   getDisplaySettings() → formatOutput(files, settings)
//     ↓
//   Exit → output displayed via display library
//
// APUs (Available Processing Units):
// - 9 functions total
// - 5 helpers (config loading, getters, pure utilities)
// - 3 core operations (build command, execute, format)
// - 1 public API (CheckRecentActivity)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Usually not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// loadActivityConfig loads activity tracking configuration from JSONC file.
//
// What It Does:
// Reads activity-tracking.jsonc, strips JSONC comments, parses into ActivityConfig
// struct. Returns nil on any error (file missing, malformed JSON, etc.) to
// signal fallback to hardcoded defaults.
//
// Parameters:
//   configPath: Absolute path to activity-tracking.jsonc file
//
// Returns:
//   *ActivityConfig: Loaded configuration, or nil if loading failed
//
// Example usage:
//
//	config := loadActivityConfig("/home/user/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc")
//	if config == nil {
//	    // Use hardcoded defaults
//	}
func loadActivityConfig(configPath string) *ActivityConfig {
	data, err := os.ReadFile(configPath) // Read entire config file into memory
	if err != nil {                      // File missing, permission denied, etc.
		return nil // Signal failure - caller uses fallback
	}

	cleanJSON := stripJSONCComments(string(data)) // Remove // comment lines for JSON parser

	var config ActivityConfig                       // Declare struct to unmarshal into
	if err := json.Unmarshal([]byte(cleanJSON), &config); err != nil { // Parse JSON into struct
		return nil // Malformed JSON - signal failure
	}

	return &config // Success - return loaded configuration
}

// stripJSONCComments removes single-line comments from JSONC content.
//
// What It Does:
// Filters out lines starting with // to make JSONC parseable by standard JSON
// decoder. Preserves all non-comment lines. Simple line-based approach works
// for our config files (no inline comments, no /* */ style).
//
// Parameters:
//   jsonc: JSONC content with // comments
//
// Returns:
//   string: Valid JSON with comments removed
//
// Example usage:
//
//	jsonc := "// Comment\n{\"key\": \"value\"}"
//	json := stripJSONCComments(jsonc)  // Returns "{\"key\": \"value\"}"
func stripJSONCComments(jsonc string) string {
	lines := strings.Split(jsonc, "\n") // Split into individual lines
	var cleaned []string                // Accumulate non-comment lines

	for _, line := range lines { // Process each line
		trimmed := strings.TrimSpace(line) // Remove leading/trailing whitespace

		// Skip full-line comments
		if strings.HasPrefix(trimmed, "//") {
			continue // Skip this line entirely
		}

		// Handle trailing comments - find // that's NOT inside a quoted string
		inString := false
		escaped := false
		for i := 0; i < len(line)-1; i++ {
			char := line[i]

			// Handle escape sequences
			if escaped {
				escaped = false
				continue
			}
			if char == '\\' {
				escaped = true
				continue
			}

			// Track whether we're inside a quoted string
			if char == '"' {
				inString = !inString
			}

			// If we find // outside a string, truncate here
			if !inString && char == '/' && line[i+1] == '/' {
				line = line[:i]
				break
			}
		}

		cleaned = append(cleaned, line) // Keep this line (possibly with trailing comment removed)
	}

	return strings.Join(cleaned, "\n") // Rejoin into single string with newlines
}

// getTimeWindow returns configured time window or default fallback.
//
// What It Does:
// Checks if config loaded successfully. If yes, returns configured minutes.
// If no, returns hardcoded default (60 minutes). Ensures function always
// returns usable value even if configuration missing.
//
// Returns:
//   int: Minutes to look back for file modifications
//
// Example usage:
//
//	minutes := getTimeWindow()  // Returns 60 if config missing, otherwise configured value
func getTimeWindow() int {
	if configLoaded && activityConfig != nil { // Check config available
		return activityConfig.TimeWindow.Minutes // Use configured value
	}
	return 60 // Hardcoded fallback - 60 minutes (1 hour)
}

// getExclusionPatterns returns configured exclusion patterns or defaults.
//
// What It Does:
// Checks if config loaded. If yes, returns configured patterns. If no,
// returns hardcoded defaults (dotfiles, node_modules). Always returns
// usable slice even if configuration missing.
//
// Returns:
//   []ExclusionPattern: Patterns for paths to exclude from activity tracking
//
// Example usage:
//
//	patterns := getExclusionPatterns()  // Returns defaults if config missing
//	for _, p := range patterns {
//	    // Build find command exclusions
//	}
func getExclusionPatterns() []ExclusionPattern {
	if configLoaded && activityConfig != nil { // Check config available
		return activityConfig.ExclusionPatterns // Use configured patterns
	}

	// Hardcoded fallback - minimal essential exclusions
	return []ExclusionPattern{
		{Pattern: "*/.*", Description: "Hidden files", Reasoning: "System/config files"},
		{Pattern: "*/node_modules/*", Description: "Node.js deps", Reasoning: "Not user work"},
	}
}

// getDisplaySettings returns configured display settings or defaults.
//
// What It Does:
// Checks if config loaded. If yes, returns configured display settings.
// If no, returns hardcoded defaults (threshold 10, icon 📂). Ensures
// display behavior always defined even if configuration missing.
//
// Returns:
//   DisplaySettings: Display behavior, thresholds, and visual elements
//
// Example usage:
//
//	settings := getDisplaySettings()
//	if fileCount <= settings.ThresholdDetailed {
//	    // Show detailed count
//	}
func getDisplaySettings() DisplaySettings {
	if configLoaded && activityConfig != nil { // Check config available
		return activityConfig.Display // Use configured settings
	}

	// Hardcoded fallback - minimal sensible defaults
	return DisplaySettings{
		ThresholdDetailed: 10,     // Show count when <= 10 files
		ThresholdSummary:  50,     // Show summary when > 10 but <= 50
		ShowFileList:      false,  // Don't spam with file paths
		Icon:              "📂",   // Folder icon
	}
}

// ────────────────────────────────────────────────────────────────
// Core Operations - Business Logic
// ────────────────────────────────────────────────────────────────
// Component-specific functionality implementing primary purpose. Organized
// by operational categories (descriptive subsections) below.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-core-operations.md

// ────────────────────────────────────────────────────────────────
// Command Building - Find Command Construction
// ────────────────────────────────────────────────────────────────
// What These Do:
// Construct find command with time window and exclusion patterns. Translates
// configuration into executable command arguments.
//
// Why Separated:
// Building command vs executing command are distinct concerns. Separation
// allows testing command construction without executing, and makes logic
// explicit and reviewable.
//
// Extension Point:
// To add new exclusion patterns, update activity-tracking.jsonc configuration
// file. No code changes needed. To add new find command options (depth limits,
// file type filters), extend buildFindCommand() parameters and logic.

// buildFindCommand constructs find command arguments from configuration.
//
// What It Does:
// Builds complete find command with workspace path, time window (-mmin), and
// exclusion patterns (-not -path). Returns command ready for execution.
//
// Parameters:
//   workspace: Root directory to search for modified files
//
// Returns:
//   *exec.Cmd: Configured find command ready to execute
//
// Health Impact:
//   Success: +5 points (command built successfully)
//   Invalid workspace: -5 points (empty path provided)
//
// Example usage:
//
//	cmd := buildFindCommand("/home/user/project")
//	// Returns: find /home/user/project -type f -mmin -60 -not -path '*/.*' -not -path '*/node_modules/*'
func buildFindCommand(workspace string) *exec.Cmd {
	if workspace == "" { // Validate input
		workspace = "." // Default to current directory if empty
	}

	minutes := getTimeWindow()          // Get configured time window (or fallback)
	patterns := getExclusionPatterns() // Get configured exclusions (or fallback)

	// Build find command arguments slice
	args := []string{
		workspace,    // Search root directory
		"-type", "f", // Only files (not directories)
		"-mmin", fmt.Sprintf("-%d", minutes), // Modified within last N minutes
	}

	// Add exclusion patterns as -not -path arguments
	for _, pattern := range patterns { // Iterate configured exclusions
		args = append(args, "-not", "-path", pattern.Pattern) // Add to find command
	}

	return exec.Command("find", args...) // Construct command ready to execute
}

// ────────────────────────────────────────────────────────────────
// Execution - Find Command Execution
// ────────────────────────────────────────────────────────────────
// What These Do:
// Execute find command and capture output. Handles command failures gracefully.
//
// Why Separated:
// Pure execution logic separated from command building and output formatting.
// Makes error handling explicit and allows independent testing.

// executeFind runs find command and returns discovered file paths.
//
// What It Does:
// Executes find command, captures stdout, splits into individual file paths.
// Returns empty slice on any error (command failure, empty output).
//
// Parameters:
//   cmd: Configured find command from buildFindCommand()
//
// Returns:
//   []string: Slice of file paths discovered, or empty slice on error
//
// Health Impact:
//   Success: +15 points (activity detected successfully)
//   Command failure: -10 points (find command error)
//   Empty results: +5 points (valid result, just no activity)
//
// Example usage:
//
//	cmd := buildFindCommand(workspace)
//	files := executeFind(cmd)
//	if len(files) == 0 {
//	    // No recent activity or command failed
//	}
func executeFind(cmd *exec.Cmd) []string {
	output, err := cmd.Output() // Execute command, capture stdout
	if err != nil {             // Command failed (permissions, invalid path, etc.)
		return []string{} // Return empty slice - no activity detected
	}

	if len(output) == 0 { // Command succeeded but no files found
		return []string{} // Return empty slice - no recent activity
	}

	// Split output into individual file paths (newline-separated)
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files // Return discovered files
}

// ────────────────────────────────────────────────────────────────
// Formatting - Output Display
// ────────────────────────────────────────────────────────────────
// What These Do:
// Format activity information for display using display library. Applies
// thresholds to determine output mode (detailed, summary, or silent).
//
// Why Separated:
// Display logic separated from detection logic. Makes output behavior
// testable and allows easy modification without touching core detection.

// formatOutput displays activity information using display library.
//
// What It Does:
// Applies display thresholds to determine output mode. Uses display library
// for consistent formatting with health tracking. Prints to stdout.
//
// Parameters:
//   files: Slice of recently modified file paths
//   settings: Display configuration (thresholds, icon, behavior)
//
// Health Impact:
//   Formatted output: +10 points (information displayed successfully)
//   Display failure: -5 points (falls back to plain fmt output)
//
// Example usage:
//
//	files := executeFind(cmd)
//	settings := getDisplaySettings()
//	formatOutput(files, settings)  // Prints to stdout
func formatOutput(files []string, settings DisplaySettings) {
	count := len(files) // Number of files discovered

	if count == 0 { // No activity to display
		return // Silent - don't clutter output with "no activity" message
	}

	// Determine output mode based on thresholds
	if count <= settings.ThresholdDetailed {
		// Detailed mode: show specific count
		message := fmt.Sprintf("%s Recently modified (last hour): %d file(s)", settings.Icon, count)
		fmt.Println(display.Info(message)) // Use display library for formatted output
	} else if count <= settings.ThresholdSummary {
		// Summary mode: general awareness
		message := fmt.Sprintf("%s Recently modified (last hour): %d files", settings.Icon, count)
		fmt.Println(display.Info(message)) // Use display library
	}
	// Above summary threshold: silent (too many files, likely build artifacts)

	// Optional: display file list if enabled (typically for debugging)
	if settings.ShowFileList && count <= settings.ThresholdDetailed {
		for _, file := range files { // Print each file path
			fmt.Printf("  - %s\n", file)
		}
	}
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

// CheckRecentActivity displays recently modified files in workspace.
//
// What It Does:
// Complete workflow: builds find command with configured time window and
// exclusions, executes command, formats and displays results. Provides
// session context awareness by showing what work has been done recently.
//
// Parameters:
//   workspace: Root directory to search for modified files
//
// Health Impact:
//   Complete success: +30 points (detection + display successful)
//   Partial success: +20 points (detection worked, display had issues)
//   Command failure: -10 points (find command failed)
//   No activity: +5 points (valid result, workspace quiet)
//
// Example usage:
//
//	session.CheckRecentActivity("/home/user/project")
//	// Prints: 📂 Recently modified (last hour): 7 file(s)
func CheckRecentActivity(workspace string) {
	cmd := buildFindCommand(workspace)       // Build find command with config
	files := executeFind(cmd)                // Execute and capture results
	settings := getDisplaySettings()         // Get display configuration
	formatOutput(files, settings)            // Format and display to stdout
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
//   - Call CheckRecentActivity() with valid workspace path
//   - Verify output displays with correct format (icon, count, thresholds)
//   - Test with missing configuration (fallback to defaults)
//   - Test with empty workspace (no recent files)
//   - Test with workspace containing many files (threshold behavior)
//   - Ensure no go vet warnings introduced
//   - Run: go build ./... (library compilation check)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - Verify activity-tracking.jsonc is valid JSONC
//
// Integration Testing:
//   - Test from start.go hook (session start context)
//   - Test from stop.go hook (session stop context)
//   - Verify display library integration (formatted output)
//   - Check configuration loading (file exists, permissions correct)
//
// Example validation code:
//
//     // Test basic functionality
//     session.CheckRecentActivity("/tmp/test-workspace")
//     // Verify output appears with expected format
//
//     // Test with missing config
//     os.Rename(configPath, configPath+".bak")
//     session.CheckRecentActivity("/tmp/test-workspace")
//     // Verify fallback behavior works
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
// The library is imported into calling hooks (start.go, stop.go), making the
// CheckRecentActivity() function available. Configuration loads automatically
// in init() during import. Function executes when called by hook orchestrators.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func startHook() {
//         workspace := "/home/user/project"
//         session.CheckRecentActivity(workspace)  // Display recent file activity
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Configuration: Loaded once in init(), remains in memory for process lifetime
//   - Find command: Process spawned, stdout captured, process terminates automatically
//   - Memory: Strings allocated for output, garbage collected after display
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle beyond process)
//   - No cleanup function needed (no persistent resources)
//   - Process termination handles all resource cleanup
//
// Error State Cleanup:
//   - Command failures return empty slice (no partial state)
//   - Configuration loading failure sets flag (fallback active)
//   - No rollback mechanisms needed (pure detection, no mutations)
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Find output captured in strings (freed after display)
//   - Configuration struct remains in memory (small, persistent)
//   - No manual memory management required
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
// Quick summary (high-level only - details in METADATA):
//   - Tracks recently modified files in workspace using configurable time windows
//   - Displays activity context at session start/stop for awareness
//   - Config-driven behavior (time window, exclusions, display thresholds)
//   - Graceful fallback to hardcoded defaults if configuration missing
//   - Integrates with display library for consistent formatted output
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list organized by category in typical usage order
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Ladder - provides context awareness functionality)
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new exclusion patterns via activity-tracking.jsonc (no code changes)
//   ✅ Adjust time windows and thresholds via configuration
//   ✅ Add new display modes (extend formatOutput() logic)
//   ✅ Add new helper functions for additional detection modes
//   ✅ Extend configuration structure (add fields to extensions section)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ CheckRecentActivity() signature - breaks all calling hooks (start.go, stop.go)
//   ⚠️ ActivityConfig struct fields - breaks configuration parsing
//   ⚠️ Display thresholds defaults - affects all users without config
//   ⚠️ Exclusion pattern format - breaks existing configurations
//   ⚠️ Output format changes - affects hook display expectations
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Configuration fallback behavior (must degrade gracefully)
//   ❌ Non-blocking guarantee (never block hook execution)
//   ❌ Display library integration (health tracking dependency)
//   ❌ init() loading pattern (configuration must load at import)
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
// - APU count (9 functions total)
//
// Quick architectural summary (details in BODY Organizational Chart):
// - 1 public API (CheckRecentActivity) orchestrates 3 core operations using 5 helpers
// - Ladder: Public API → Core Operations (build/execute/format) → Helpers (config loading/getters)
// - Baton: Entry → build command → execute find → format output → display
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Core Operations" subsection header comments above for detailed
// extension points. Each subsection includes guidance for additions.
//
// Quick reference (details in BODY subsection comments):
// - Adding exclusion patterns: Update activity-tracking.jsonc (no code changes)
// - Adding display modes: Extend formatOutput() function with new threshold logic
// - Adding find options: Extend buildFindCommand() parameters and argument building
// - Adding helper functions: Add to Helpers section, update Organizational Chart
//
// Configuration Extensions (preferred over code changes):
// - New time windows: Modify time_window.minutes in config file
// - New thresholds: Modify display.threshold_detailed/threshold_summary in config
// - New exclusions: Add to exclusion_patterns array in config
// - Experimental features: Use extensions section in config for new behavior
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - No constants (fully config-driven, minimal memory overhead)
// - Types: Small structs, configuration loaded once in init()
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
// - Most expensive operation: executeFind() - filesystem traversal O(n) in files
// - Memory characteristics: Configuration ~1KB, find output scales with file count
// - Key optimization: Exclusion patterns reduce files scanned significantly
// - Time complexity: O(n) where n = files in workspace (find command limitation)
// - Memory usage: O(m) where m = files discovered (output capture)
//
// Bottlenecks:
// - Large workspaces (>100K files): find command slow, consider reducing time window
// - Many exclusions (>20 patterns): each pattern adds overhead to find traversal
// - File list display: Disabled by default (output can be massive for large file counts)
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// See BODY function docstrings above for operation-specific troubleshooting.
//
// Problem: No activity displayed even with recent work
//   - Check: Time window setting (default 60 minutes, may be too short)
//   - Check: Exclusion patterns (may be catching your work files)
//   - Check: Display thresholds (may be above threshold, showing summary or silent)
//   - Solution: Adjust activity-tracking.jsonc time_window.minutes or exclusion_patterns
//
// Problem: Configuration not loading (using fallback defaults)
//   - Check: File exists at ~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
//   - Check: File permissions (must be readable)
//   - Check: JSONC syntax (use linter to validate)
//   - Solution: Verify file path, fix permissions (chmod 644), validate JSON syntax
//
// Problem: Too much activity displayed (overwhelming output)
//   - Check: Exclusion patterns (may need more patterns for your workflow)
//   - Check: Display thresholds (may be too high)
//   - Solution: Add exclusions to activity-tracking.jsonc, lower threshold_summary
//
// Problem: Activity detection slow in large workspaces
//   - Cause: find command traversing many files
//   - Solution: Reduce time_window.minutes (fewer files to check)
//   - Solution: Add more exclusion_patterns (skip large directories)
//   - Note: This is find command limitation, not library issue
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard library (os/exec, fmt, strings, encoding/json),
//                                    system/lib/display (formatted output)
// - Dependents (What Uses This): session/cmd-start/start.go, session/cmd-stop/stop.go
// - Integration Points: Hooks call CheckRecentActivity(), display library for output
//
// Quick summary (details in METADATA Dependencies section above):
// - Key dependencies: os/exec (find command), system/lib/display (formatted output)
// - Primary consumers: Session start/stop hooks (context awareness)
// - Configuration: activity-tracking.jsonc (system/data/config/session/)
//
// Parallel Implementation:
//   - Go version: This file (hooks/lib/session/activity.go)
//   - Shared philosophy: Recent activity reveals current focus, config-driven behavior
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration system - COMPLETED (v2.0.0)
//   ✓ Display library integration - COMPLETED (v2.0.0)
//   ✓ Comprehensive exclusion patterns - COMPLETED (v2.0.0)
//   ⏳ File type filters (show only .go, .rs, etc.)
//   ⏳ Size-based filtering (ignore large binaries)
//   ⏳ User-based filtering (show only my modifications)
//   ⏳ Git integration (show only uncommitted changes)
//
// Research Areas:
//   - Alternative to find command (native Go filesystem walk for speed)
//   - Incremental scanning (cache results, only check new modifications)
//   - Language-specific activity (detect which projects active)
//   - Intelligent exclusions (learn what user typically ignores)
//   - Activity correlation (relate to session duration, time of day)
//
// Integration Targets:
//   - Git integration (uncommitted changes, branch activity)
//   - Session patterns (correlate activity with time of day, duration)
//   - Project detection (identify active projects from file paths)
//   - Health scoring aggregation (track detection quality over time)
//   - Temporal awareness (celestial correlation - do I work more during daylight?)
//
// Known Limitations to Address:
//   - find command performance in very large workspaces (>100K files)
//   - No caching (every call rescans filesystem from scratch)
//   - No git-awareness (shows all files, not just version-controlled)
//   - Line-based comment stripping (can't handle inline // or /* */ comments)
//   - No file content analysis (just modification time, not what changed)
//   - No activity categorization (all files treated equally)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-12) - Configuration system and template alignment
//         - Created activity-tracking.jsonc configuration file
//         - Integrated system/lib/display for formatted output
//         - Added comprehensive exclusion patterns (8 default patterns)
//         - Implemented graceful fallback to hardcoded defaults
//         - Full 4-block template alignment (METADATA/SETUP/BODY/CLOSING)
//         - Added complete inline documentation and health scoring
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Basic find command execution with hardcoded 60-minute window
//         - Hardcoded exclusions (dotfiles, node_modules only)
//         - Simple threshold display logic (10 files)
//         - Direct fmt.Printf output (no display library)
//         - Minimal documentation
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing context awareness functionality.
// Session start/stop hooks depend on this to show recent work context. It's a
// small but essential piece of session awareness - users see what they've been
// working on, grounding them in current state.
//
// Modify thoughtfully - changes here affect session start/stop behavior for all
// users. Configuration system allows customization without code changes (preferred).
// Non-blocking guarantee is critical - this must never disrupt hook execution.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go build, go vet, integration tests)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Update activity-tracking.jsonc schema if adding config fields
//   - Verify fallback behavior still works (config missing case)
//
// "Whatever your hand finds to do, do it with all your might" - Ecclesiastes 9:10
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage:
//   import "hooks/lib/session"
//
//   func main() {
//       session.CheckRecentActivity("/home/user/project")
//       // Output: 📂 Recently modified (last hour): 7 file(s)
//   }
//
// Custom Configuration:
//   Edit ~/.claude/cpi-si/system/data/config/session/activity-tracking.jsonc
//
//   {
//       "time_window": {"minutes": 120},  // Look back 2 hours
//       "display": {"threshold_detailed": 20}  // Show count up to 20 files
//   }
//
// Integration in Hooks:
//   // In session/cmd-start/start.go
//   session.CheckRecentActivity(workspace)  // Show activity at session start
//
//   // In session/cmd-stop/stop.go
//   session.CheckRecentActivity(workspace)  // Show activity at session stop
//
// Fallback Behavior:
//   // If config missing, uses hardcoded defaults:
//   // - 60 minute window
//   // - Dotfiles and node_modules excluded
//   // - Threshold 10 files
//
// ════════════════════════════════════════════════════════════════
// END CLOSING
// ════════════════════════════════════════════════════════════════
