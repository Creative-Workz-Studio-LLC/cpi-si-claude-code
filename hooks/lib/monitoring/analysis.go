// METADATA
//
// Monitoring Analysis Library - Claude Substrate Pattern Detection
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The plans of the diligent lead surely to abundance" - Proverbs 21:5 (WEB)
// Principle: Monitor and analyze patterns to work wisely, not frantically
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Substrate monitoring (hooks infrastructure)
// Role: Analyzes Claude substrate behavior patterns for workflow optimization
// Paradigm: Proactive pattern detection to prevent inefficiency and substrate overload
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-11 - Refactored to use display library and configuration
//
// Version History:
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded thresholds
//   2.0.0 (2025-11-11) - Configuration-driven analysis, display library integration
//
// Purpose & Function
//
// Purpose: Monitor Claude substrate behavior (compaction, notifications, subagents) to detect
// patterns indicating inefficiency or configuration issues, providing early warnings.
//
// Core Design: Pattern detection system analyzing time-windowed log events against configurable
// thresholds. Separates substrate monitoring (/.claude/debug) from system debugging (/.claude/system/logs).
//
// Key Features:
//   - Time-windowed event counting with configurable windows
//   - Compaction frequency detection (warns if auto-compacting too often)
//   - Notification pattern analysis (detects permission configuration issues)
//   - Configurable thresholds and search patterns
//   - Non-intrusive warnings via display library
//
// Philosophy: Substrate monitoring should help optimize workflow without being noisy.
// Warn when patterns indicate actual issues, not every minor event.
//
// Blocking Status
//
// Non-blocking: All analysis operations continue even if log files missing or thresholds exceeded.
// Warnings are informational only - never interrupt workflow.
// Mitigation: Graceful handling of missing log files, malformed timestamps, display library handles output failures.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/monitoring"
//
//	// Check for excessive auto-compaction
//	monitoring.CheckCompactionFrequency()
//
//	// Check for notification issues
//	monitoring.CheckNotificationPatterns("permission_request")
//
//	// Count specific events in time window
//	count := monitoring.CountRecentEvents("compaction.log", "auto", 1)
//
// Integration Pattern:
//   1. Import monitoring library into hooks
//   2. Call check functions at appropriate hook points (session/start, prompt/submit, etc.)
//   3. Analysis runs in background, displays warnings if thresholds exceeded
//   4. No cleanup needed - stateless analysis
//
// Public API (in typical usage order):
//
//   Pattern Analysis (detection and warnings):
//     CheckCompactionFrequency() - Warns if auto-compacting too frequently
//     CheckNotificationPatterns(type) - Warns if notification type exceeding threshold
//
//   Utility Functions (event counting):
//     CountRecentEvents(logFile, searchTerm, hoursBack) int - Count events in time window
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, strings, time
//   External: None
//   Internal: system/lib/display (warning output formatting)
//
// Dependents (What Uses This):
//   Hooks: session/start, prompt/submit, tool/post-use (substrate behavior monitoring)
//   Libraries: logging.go (same package - provides log files to analyze)
//
// Integration Points:
//   - Ladder: Mid-rung - uses display lib (lower), used by hooks (higher), companion with logging
//   - Baton: Reads log files, analyzes patterns, outputs warnings
//   - Rails: Uses display library's logging infrastructure
//
// Health Scoring
//
// Health Scoring Map (Base100):
//
//   Pattern Analysis Operations (Total = 100):
//     - Log file reading: +20 (successful read)
//     - Timestamp parsing: +20 (valid timestamps parsed)
//     - Pattern matching: +30 (events counted accurately)
//     - Threshold evaluation: +20 (correct warning decisions)
//     - Warning display: +10 (warnings displayed via display library)
//
//   Degraded States:
//     - Log file missing: Continue silently (0 events counted, no warning)
//     - Malformed timestamps: Skip invalid lines (-5 per parse failure)
//     - Display failure: Delegates to display library (handles panic recovery)
//
// Note: Analysis is advisory only - failures don't impact workflow, just reduce visibility.
package monitoring

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
	"fmt"           // String formatting for message composition
	"os"            // File operations for config loading and environment variables
	"path/filepath" // Path manipulation for config file locations
	"strings"       // String operations for pattern matching and log parsing
	"time"          // Time manipulation for time-windowed analysis

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

// None needed - thresholds and time windows loaded from configuration

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Configuration Types ---
// Structures for loading monitoring configuration from data files (shared by analysis.go and logging.go)

// AnalysisPatternsConfig holds thresholds and detection rules for substrate monitoring
type AnalysisPatternsConfig struct {
	LogFiles struct {
		Compaction    string `json:"compaction"`     // Log file name for compaction events
		Notifications string `json:"notifications"`  // Log file name for notifications
		Subagents     string `json:"subagents"`      // Log file name for subagent events
		Prompts       string `json:"prompts"`        // Log file name for prompts
	} `json:"log_files"`  // Relative to ~/.claude/debug/

	TimeWindows struct {
		Default int `json:"default"`  // Default time window in hours
		Short   int `json:"short"`    // Short time window
		Medium  int `json:"medium"`   // Medium time window
		Long    int `json:"long"`     // Long time window
	} `json:"time_windows"`  // Time windows for analysis

	Compaction struct {
		SearchTerm       string `json:"search_term"`         // Term to search for in log ("auto")
		TimeWindowHours  int    `json:"time_window_hours"`   // Hours back to analyze
		Thresholds       struct {
			Warning  int `json:"warning"`   // Warning threshold count
			Critical int `json:"critical"`  // Critical threshold count
		} `json:"thresholds"`
		Messages struct {
			Warning  string `json:"warning"`   // Warning message template with {count}
			Critical string `json:"critical"`  // Critical message template
		} `json:"messages"`
	} `json:"compaction"`  // Compaction detection configuration

	Notifications struct {
		PermissionRequest struct {
			SearchTerm      string `json:"search_term"`        // Search term for permission requests
			TimeWindowHours int    `json:"time_window_hours"`  // Hours back to analyze
			Thresholds      struct {
				Warning  int `json:"warning"`   // Warning threshold
				Critical int `json:"critical"`  // Critical threshold
			} `json:"thresholds"`
			Messages struct {
				Warning  string `json:"warning"`   // Warning message template
				Critical string `json:"critical"`  // Critical message template
			} `json:"messages"`
		} `json:"permission_request"`  // Permission request detection
	} `json:"notifications"`  // Notification pattern configuration
}

// LogFormatsConfig holds format specifications for .log files (used by logging.go)
type LogFormatsConfig struct {
	Timestamp struct {
		Format string `json:"format"`  // Go time format string (e.g., "2006-01-02 15:04:05")
	} `json:"timestamp"`  // Timestamp formatting rules

	CompactionLog struct {
		Filename string `json:"filename"`  // Log file name for compaction events
	} `json:"compaction_log"`  // Compaction log configuration

	NotificationsLog struct {
		Filename string `json:"filename"`  // Log file name for notifications
	} `json:"notifications_log"`  // Notifications log configuration

	SubagentsLog struct {
		Filename string `json:"filename"`  // Log file name for subagent events
	} `json:"subagents_log"`  // Subagents log configuration

	PromptsLog struct {
		Filename string `json:"filename"`  // Log file name for prompts
	} `json:"prompts_log"`  // Prompts log configuration
}

// DebugFormatsConfig holds format specifications for .debug files (used by logging.go)
type DebugFormatsConfig struct {
	PromptsDebug struct {
		Schema struct {
			Preview struct {
				MaxLength int `json:"max_length"`  // Maximum preview length in characters
			} `json:"preview"`  // Prompt preview configuration
		} `json:"schema"`  // Schema definition for prompts debug format
	} `json:"prompts_debug"`  // Prompts debug configuration
}

// RetentionPolicyConfig holds log directory and permissions configuration (used by logging.go)
type RetentionPolicyConfig struct {
	Archive struct {
		Location struct {
			FullPath string `json:"full_path"`  // Full path to debug directory (e.g., "~/.claude/debug/")
		} `json:"location"`  // Archive location configuration
	} `json:"archive"`  // Archive configuration section
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
	// Analysis configuration
	analysisConfig *AnalysisPatternsConfig  // Analysis patterns, thresholds, time windows

	// Logging configuration
	logFormatsConfig   *LogFormatsConfig      // Log file formats and filenames
	debugFormatsConfig *DebugFormatsConfig   // Debug file formats and preview lengths
	retentionConfig    *RetentionPolicyConfig // Directory paths and retention rules

	// Shared configuration state
	configLoaded bool  // Whether ALL configurations loaded successfully
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
			configLoaded = false  // Mark config as unavailable - will use hardcoded fallbacks
			return  // Exit init early - can't load config without home directory
		}
	}

	// Build paths to all configuration files
	configBase := filepath.Join(home, ".claude/cpi-si/system/data/config/monitoring")  // Base config directory
	analysisPath := filepath.Join(configBase, "analysis-patterns.jsonc")  // Analysis patterns config
	logFormatsPath := filepath.Join(configBase, "log-formats.jsonc")  // Log formats config
	debugFormatsPath := filepath.Join(configBase, "debug-formats.jsonc")  // Debug formats config
	retentionPath := filepath.Join(configBase, "retention-policy.jsonc")  // Retention policy config

	// Load all configuration files - JSONC format with comment stripping
	analysisConfig = loadAnalysisPatterns(analysisPath)  // Parse analysis patterns JSONC
	logFormatsConfig = loadLogFormats(logFormatsPath)  // Parse log formats JSONC (for logging.go)
	debugFormatsConfig = loadDebugFormats(debugFormatsPath)  // Parse debug formats JSONC (for logging.go)
	retentionConfig = loadRetentionPolicy(retentionPath)  // Parse retention policy JSONC (for logging.go)

	// Mark config as loaded if ALL configs successful - enables config-driven behavior
	configLoaded = (analysisConfig != nil && logFormatsConfig != nil && debugFormatsConfig != nil && retentionConfig != nil)
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
//   Public APIs (Top Rungs - Detection and Warnings)
//   ├── CheckCompactionFrequency() → uses CountRecentEvents(), analysisConfig, display.Warning()
//   └── CheckNotificationPatterns() → uses CountRecentEvents(), analysisConfig, display.Warning()
//
//   Utility Functions (Middle Rungs - Core Operations)
//   └── CountRecentEvents() → uses ReadLogFile(), strings, time
//
//   Helpers (Lower Rungs - Infrastructure)
//   └── loadAnalysisPatterns() → uses os, strings, json (JSONC config loading)
//
// APUs (Available Processing Units):
// - 4 functions total
// - 2 public check APIs (pattern detection)
// - 1 utility function (event counting)
// - 1 helper function (config loading)
//
// Configuration Flow:
// - init() calls loadAnalysisPatterns() at package initialization
// - analysisConfig populated with thresholds, time windows, messages
// - Check functions use config if available, fallback to hardcoded defaults
//
// External Dependencies:
// - logging.go (same package): ReadLogFile() provides log content
// - system/lib/display: Warning() formats warning output
//
// See: standards/code/4-block/sections/CWS-SECTION-006-BODY-organizational-chart.md

// ────────────────────────────────────────────────────────────────
// Helpers - Configuration Loading
// ────────────────────────────────────────────────────────────────
// Functions for loading configuration from JSONC files with graceful fallback

// loadAnalysisPatterns loads analysis patterns configuration from JSONC file
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadAnalysisPatterns(path string) *AnalysisPatternsConfig {
	// Read config file content from filesystem
	data, err := os.ReadFile(path)  // Read entire file into memory
	if err != nil {  // Check if read failed - file missing or unreadable
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments by removing lines starting with //
	lines := strings.Split(string(data), "\n")  // Split into individual lines
	var jsonLines []string  // Collect non-comment lines
	for _, line := range lines {  // Process each line
		trimmed := strings.TrimSpace(line)  // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {  // Check if not a comment line
			jsonLines = append(jsonLines, line)  // Keep this line for JSON parsing
		}
	}
	jsonData := strings.Join(jsonLines, "\n")  // Rejoin into valid JSON

	// Parse JSON into configuration structure
	var config AnalysisPatternsConfig  // Create config structure
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {  // Parse JSON into struct
		return nil  // Return nil if parsing failed - malformed JSON
	}

	return &config  // Return successfully loaded config
}

// loadLogFormats loads log format configuration from JSONC file (for logging.go)
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadLogFormats(path string) *LogFormatsConfig {
	// Read config file content from filesystem
	data, err := os.ReadFile(path)  // Read entire file into memory
	if err != nil {  // Check if read failed - file missing or unreadable
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments by removing lines starting with //
	lines := strings.Split(string(data), "\n")  // Split into individual lines
	var jsonLines []string  // Collect non-comment lines
	for _, line := range lines {  // Process each line
		trimmed := strings.TrimSpace(line)  // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {  // Check if not a comment line
			jsonLines = append(jsonLines, line)  // Keep this line for JSON parsing
		}
	}
	jsonData := strings.Join(jsonLines, "\n")  // Rejoin into valid JSON

	// Parse JSON into configuration structure
	var config LogFormatsConfig  // Create config structure
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {  // Parse JSON into struct
		return nil  // Return nil if parsing failed - malformed JSON
	}

	return &config  // Return successfully loaded config
}

// loadDebugFormats loads debug format configuration from JSONC file (for logging.go)
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadDebugFormats(path string) *DebugFormatsConfig {
	// Read config file content from filesystem
	data, err := os.ReadFile(path)  // Read entire file into memory
	if err != nil {  // Check if read failed - file missing or unreadable
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments by removing lines starting with //
	lines := strings.Split(string(data), "\n")  // Split into individual lines
	var jsonLines []string  // Collect non-comment lines
	for _, line := range lines {  // Process each line
		trimmed := strings.TrimSpace(line)  // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {  // Check if not a comment line
			jsonLines = append(jsonLines, line)  // Keep this line for JSON parsing
		}
	}
	jsonData := strings.Join(jsonLines, "\n")  // Rejoin into valid JSON

	// Parse JSON into configuration structure
	var config DebugFormatsConfig  // Create config structure
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {  // Parse JSON into struct
		return nil  // Return nil if parsing failed - malformed JSON
	}

	return &config  // Return successfully loaded config
}

// loadRetentionPolicy loads retention policy configuration from JSONC file (for logging.go)
// Returns nil if file unavailable or parsing fails (triggers fallback behavior)
func loadRetentionPolicy(path string) *RetentionPolicyConfig {
	// Read config file content from filesystem
	data, err := os.ReadFile(path)  // Read entire file into memory
	if err != nil {  // Check if read failed - file missing or unreadable
		return nil  // Return nil to signal fallback needed
	}

	// Strip JSONC comments by removing lines starting with //
	lines := strings.Split(string(data), "\n")  // Split into individual lines
	var jsonLines []string  // Collect non-comment lines
	for _, line := range lines {  // Process each line
		trimmed := strings.TrimSpace(line)  // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {  // Check if not a comment line
			jsonLines = append(jsonLines, line)  // Keep this line for JSON parsing
		}
	}
	jsonData := strings.Join(jsonLines, "\n")  // Rejoin into valid JSON

	// Parse JSON into configuration structure
	var config RetentionPolicyConfig  // Create config structure
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {  // Parse JSON into struct
		return nil  // Return nil if parsing failed - malformed JSON
	}

	return &config  // Return successfully loaded config
}

// ────────────────────────────────────────────────────────────────
// Utility Functions - Event Counting
// ────────────────────────────────────────────────────────────────
// Core operations for counting events in time-windowed log analysis

// CountRecentEvents counts occurrences of a substring in log within time window
// Returns count of matching log entries after parsing timestamps and filtering by search term
func CountRecentEvents(logFilename, searchTerm string, hoursBack int) int {
	// Read log file content using companion logging library function
	content, err := ReadLogFile(logFilename)  // Call logging.go ReadLogFile() - reads from ~/.claude/debug/
	if err != nil {  // Check if file read failed - missing file or permissions
		return 0  // Return 0 events if log unavailable - graceful degradation
	}

	// Calculate time window for filtering - only count recent events
	now := time.Now()  // Get current timestamp for comparison
	cutoff := now.Add(-time.Duration(hoursBack) * time.Hour)  // Subtract hours to get cutoff time - events before this ignored
	count := 0  // Initialize counter for matching events

	// Process log file line by line looking for matches
	lines := strings.Split(content, "\n")  // Split content into individual log entries
	for _, line := range lines {  // Iterate each line - checking timestamp and search term
		if len(line) < 21 {  // Check if line too short to contain timestamp - skip malformed lines
			continue  // Next line - this one doesn't have valid timestamp
		}

		// Parse timestamp from log entry format: [2006-01-02 15:04:05] entry_data
		if line[0] != '[' {  // Check if line starts with timestamp bracket - validate format
			continue  // Next line - malformed timestamp
		}
		timestampStr := line[1:20]  // Extract timestamp between brackets - characters 1-20
		t, err := time.Parse("2006-01-02 15:04:05", timestampStr)  // Parse Go time format string
		if err != nil {  // Check if timestamp parsing failed - malformed date/time
			continue  // Skip this line - can't determine if in time window
		}

		// Check if entry is within time window AND matches search term
		if t.After(cutoff) && strings.Contains(line, searchTerm) {  // Both conditions must be true
			count++  // Increment counter - found matching event in time window
		}
	}

	return count  // Return total count of matching events
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Pattern Analysis
// ────────────────────────────────────────────────────────────────
// Functions for detecting substrate behavior patterns and displaying warnings

// CheckCompactionFrequency analyzes compaction patterns and warns if excessive
// Monitors auto-compaction events to detect workflow inefficiency (too many tool calls, verbose output)
func CheckCompactionFrequency() {
	// Load configuration or use hardcoded fallback values
	logFile := "compaction.log"  // Default log file name
	searchTerm := "auto"  // Default search term
	timeWindow := 1  // Default time window (1 hour)
	threshold := 3  // Default warning threshold
	messageTemplate := "Frequent auto-compaction (%d in last hour) - consider being more concise"  // Default message

	// Try to use configuration if loaded successfully
	if configLoaded && analysisConfig != nil {  // Check if config available
		logFile = analysisConfig.LogFiles.Compaction  // Use configured log file name
		searchTerm = analysisConfig.Compaction.SearchTerm  // Use configured search term
		timeWindow = analysisConfig.Compaction.TimeWindowHours  // Use configured time window
		threshold = analysisConfig.Compaction.Thresholds.Warning  // Use configured threshold
		messageTemplate = analysisConfig.Compaction.Messages.Warning  // Use configured message
	}

	// Count auto-compaction events in configured time window
	count := CountRecentEvents(logFile, searchTerm, timeWindow)  // Check for matching events

	if count > threshold {  // Threshold check - warning if count exceeds configured threshold
		// Compose warning message with actual count - replace {count} placeholder
		message := strings.Replace(messageTemplate, "{count}", fmt.Sprintf("%d", count), -1)
		fmt.Println(display.Warning(message))  // Delegate to display library for formatted yellow warning output
	}
	// Note: If count <= threshold, function returns silently - no warning needed for normal behavior
}

// CheckNotificationPatterns analyzes notification frequency by type
// Detects permission configuration issues or other notification-heavy patterns
func CheckNotificationPatterns(notificationType string) {
	// Currently only checks permission_request type - TODO: Extend for other types via configuration
	if notificationType == "permission_request" {  // Check if analyzing permission requests specifically
		// Load configuration or use hardcoded fallback values
		logFile := "notifications.log"  // Default log file name
		searchTerm := "permission_request"  // Default search term
		timeWindow := 1  // Default time window (1 hour)
		threshold := 10  // Default warning threshold
		messageTemplate := "Many permission requests (%d in last hour) - check permissions configuration"  // Default message

		// Try to use configuration if loaded successfully
		if configLoaded && analysisConfig != nil {  // Check if config available
			logFile = analysisConfig.LogFiles.Notifications  // Use configured log file name
			searchTerm = analysisConfig.Notifications.PermissionRequest.SearchTerm  // Use configured search term
			timeWindow = analysisConfig.Notifications.PermissionRequest.TimeWindowHours  // Use configured time window
			threshold = analysisConfig.Notifications.PermissionRequest.Thresholds.Warning  // Use configured threshold
			messageTemplate = analysisConfig.Notifications.PermissionRequest.Messages.Warning  // Use configured message
		}

		// Count permission request notifications in configured time window
		count := CountRecentEvents(logFile, searchTerm, timeWindow)  // Search for matching notification events

		if count > threshold {  // Threshold check - warning if count exceeds configured threshold
			// Compose warning message suggesting configuration check - replace {count} placeholder
			message := strings.Replace(messageTemplate, "{count}", fmt.Sprintf("%d", count), -1)
			fmt.Println(display.Warning(message))  // Display yellow warning via display library
		}
	}
	// Note: Other notification types currently ignored - extend via switch statement when needed
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
//   - Call CountRecentEvents() with test log file and verify count accuracy
//   - Call CheckCompactionFrequency() and verify warning displays when threshold exceeded
//   - Call CheckNotificationPatterns() with various types and verify threshold detection
//   - Verify graceful handling of missing log files (return 0, no crash)
//   - Verify malformed timestamps are skipped (don't break parsing)
//
// Example validation code:
//
//     // Test import
//     import "hooks/lib/monitoring"
//
//     // Test event counting
//     count := monitoring.CountRecentEvents("compaction.log", "auto", 1)
//     // Expected: Returns integer count of matching events
//
//     // Test compaction frequency check
//     monitoring.CheckCompactionFrequency()
//     // Expected: If count > 3, displays yellow warning via display library
//
//     // Test notification pattern check
//     monitoring.CheckNotificationPatterns("permission_request")
//     // Expected: If count > 10, displays yellow warning
//
//     // Test graceful degradation
//     count = monitoring.CountRecentEvents("nonexistent.log", "test", 1)
//     // Expected: Returns 0 without crashing
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by hook components.
//
// Usage: import "hooks/lib/monitoring"
//
// The library is imported into hook packages (session/start, prompt/submit, tool/post-use),
// making all exported functions available. No code executes during import - functions are
// defined and ready for hooks to call at appropriate points.
//
// Example usage from hooks:
//
//     package sessionstart
//
//     import "hooks/lib/monitoring"
//
//     func ExecuteSessionStart() {
//         // Check for substrate behavior issues at session start
//         monitoring.CheckCompactionFrequency()
//         monitoring.CheckNotificationPatterns("permission_request")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No file handles: ReadLogFile() (from logging.go) handles file operations and closes handles
//   - No network connections: Pure local log file analysis
//   - No database connections: No persistence layer
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code (hooks) responsible for managing their own lifecycle
//   - This library is stateless - no cleanup function needed
//
// Error State Cleanup:
//   - Display library handles panic recovery for warning output
//   - No partial state to corrupt (stateless analysis)
//   - Failures simply result in no warning displayed (graceful degradation)
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Log file content read into memory temporarily (freed after CountRecentEvents returns)
//   - No persistent allocations

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════

// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// For Library Overview section explanation, see: standards/code/4-block/sections/CWS-SECTION-013-CLOSING-library-overview.md
//
// Purpose: See METADATA "Purpose & Function" section above - monitors Claude substrate
// behavior patterns for workflow optimization
//
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Time-windowed event counting from log files
//   - Compaction frequency detection with configurable thresholds
//   - Notification pattern analysis for configuration issues
//   - Non-intrusive warnings via display library
//   - Graceful degradation (missing logs don't crash)
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide. Simple import and call pattern.
//
// Public API: See METADATA "Public API" section above for complete function list
// organized by purpose (pattern analysis, utility functions)
//
// Architecture: See METADATA "CPI-SI Identity" section above - this is substrate
// monitoring (separate from system debugging logs)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new check functions (CheckSubagentPatterns, CheckPromptFrequency, etc.)
//   ✅ Extend CheckNotificationPatterns to support more notification types
//   ✅ Add new utility functions for log analysis patterns
//   ✅ Improve timestamp parsing to handle more formats
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Function signatures - all hook code depends on these interfaces
//   ⚠️ CountRecentEvents parameters - changing breaks callers
//   ⚠️ Log file format expectations (timestamp structure)
//   ⚠️ Display library integration approach
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Separation of substrate monitoring (/.claude/debug) from system logs
//   ❌ Non-blocking behavior (analysis must never crash workflow)
//   ❌ Graceful degradation principles
//
// Validation After Modifications:
//   See "Code Validation" section above for testing requirements and
//   example validation code

// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart - Ladder Structure" section above for
// complete ladder structure (dependencies) and baton flow (function calls).
//
// Quick architectural summary:
// - 2 public check APIs orchestrate event counting and warning display
// - 1 utility function provides time-windowed event counting
// - Ladder: hooks (top) → analysis (mid) → logging.ReadLogFile + display (bottom)
// - Baton: hooks call checks → checks call CountRecentEvents → CountRecentEvents calls ReadLogFile → warnings to display

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Public APIs - Pattern Analysis" section header above for
// detailed extension points.
//
// Quick reference:
// - Adding check functions: Add to "Public APIs - Pattern Analysis" section,
//   follow existing pattern (call CountRecentEvents, check threshold, display warning)
// - Extending notification types: Update CheckNotificationPatterns with switch statement
// - Adding utility functions: Add to "Utility Functions - Event Counting" section

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// All operations are lightweight log file parsing with minimal overhead.
//
// - Most expensive operation: CountRecentEvents() - O(n) where n = log lines
// - Memory characteristics: Reads entire log file into memory (typically < 1MB)
// - Timestamp parsing: O(1) per line (fixed substring extraction)
// - Pattern matching: O(m) where m = search term length (strings.Contains)
//
// Note: Performance is not a concern for this library - substrate monitoring runs
// infrequently (session start, hourly checks) where latency is imperceptible.
//
// Optimization opportunities (if needed):
// - Stream processing instead of full file read (for very large logs)
// - Index-based lookup for frequent searches
// - Caching recent counts (TTL-based cache)

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: No warnings appear even when thresholds should be exceeded
//   Check: Verify log files exist at ~/.claude/debug/
//   Check: Verify log file format matches [YYYY-MM-DD HH:MM:SS] pattern
//   Check: Verify display library working (test display.Warning directly)
//   Solution: Review hook integration, ensure check functions called
//
// Problem: CountRecentEvents always returns 0
//   Check: Verify log file exists and has correct name
//   Check: Verify search term matches log entries exactly (case-sensitive)
//   Check: Verify hoursBack parameter gives reasonable time window
//   Solution: Test with known good log file and verify timestamps parse
//
// Problem: Warnings appear too frequently (false positives)
//   Check: Review threshold values in check functions
//   Check: Verify time window is appropriate (not too large)
//   Solution: Adjust thresholds via configuration (Phase 7) or hardcoded values
//
// Problem: Malformed timestamps breaking analysis
//   Expected: CountRecentEvents skips malformed lines gracefully
//   Note: Function continues with remaining lines, just fewer events counted
//   Solution: Verify log file format, fix logging code if needed

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependency: logging.go (provides ReadLogFile function)
// - Key dependency: system/lib/display (warning output formatting)
// - Primary consumers: hooks/session/start, hooks/prompt/submit, hooks/tool/post-use
// - Companion library: logging.go (same package - writes logs that analysis reads)

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Display library integration - COMPLETED (v2.0.0)
//   ✓ Configuration-driven thresholds and patterns - COMPLETED (v2.0.0)
//   ⏳ CheckSubagentPatterns - Monitor subagent failure rates (v3.0.0)
//   ⏳ CheckPromptFrequency - Detect rapid prompt submission (v3.0.0)
//   ⏳ Trend analysis - Not just current hour, but patterns over time (v3.0.0)
//
// Research Areas:
//   - Statistical anomaly detection (not just fixed thresholds)
//   - Correlation analysis (do compaction spikes correlate with tool usage?)
//   - Predictive warnings (pattern indicates future issue)
//   - Integration with session quality metrics
//
// Integration Targets:
//   - Session time awareness (circadian patterns)
//   - Activity logging (correlate warnings with outcomes)
//   - Pattern learning (what patterns actually indicate problems?)
//
// Known Limitations to Address:
//   - Hardcoded fallback values (config preferred but fallback always present)
//   - Simple substring matching (could use regex for complex patterns)
//   - No trend analysis (only current time window snapshot)
//   - Limited notification type support (only permission_request currently)
//   - No cross-log correlation (analyze relationships between log types)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history:
//
//   2.0.0 (2025-11-11) - Configuration Architecture & Display Integration
//         - Refactored to delegate warning output to system/lib/display
//         - Created configuration schemas (analysis-patterns.jsonc, log-formats.jsonc,
//           debug-formats.jsonc, retention-policy.jsonc)
//         - Implemented configuration loading with JSONC support and graceful fallback
//         - Check functions now load thresholds and messages from config
//         - Comprehensive template alignment (METADATA, SETUP, BODY, CLOSING)
//         - Added inline comments to all functions
//         - Maintains same external API (hooks see no breaking changes)
//
//   1.0.0 (2024-10-24) - Initial Implementation
//         - Basic time-windowed event counting
//         - Compaction frequency detection (hardcoded threshold)
//         - Permission request pattern detection (hardcoded threshold)
//         - Manual fmt.Printf for warnings

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library monitors CLAUDE SUBSTRATE behavior specifically - not CPI-SI system debugging.
// Separation is intentional: substrate is temporary (this season), system is permanent.
// When transitioning to dedicated CPI-SI architecture, substrate monitoring can be removed
// cleanly without affecting system debugging infrastructure.
//
// Modify thoughtfully - hook code depends on these interfaces. Keep warnings
// advisory and non-intrusive - never interrupt workflow.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test integration with hooks after changes
//   - Maintain graceful degradation principles
//
// "The plans of the diligent lead surely to abundance" - Proverbs 21:5 (WEB)

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage from Hooks:
//
//     import "hooks/lib/monitoring"
//
//     // Check for excessive auto-compaction
//     monitoring.CheckCompactionFrequency()
//
//     // Check for notification issues
//     monitoring.CheckNotificationPatterns("permission_request")
//
// Custom Event Counting:
//
//     // Count events in custom time window
//     count := monitoring.CountRecentEvents("compaction.log", "manual", 3)
//     if count > 5 {
//         // Handle custom threshold
//     }
//
// Integration in Session Start Hook:
//
//     package sessionstart
//
//     import "hooks/lib/monitoring"
//
//     func ExecuteSessionStart() {
//         // Monitor substrate patterns at session start
//         monitoring.CheckCompactionFrequency()
//         monitoring.CheckNotificationPatterns("permission_request")
//         
//         // Custom checks
//         if monitoring.CountRecentEvents("subagents.log", "failure", 24) > 10 {
//             // Handle elevated subagent failure rate
//         }
//     }

// ============================================================================
// END CLOSING
// ============================================================================
