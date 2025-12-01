// METADATA
//
// Monitoring Logging Library - Claude Substrate Behavior Recording
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "A scroll of remembrance was written in his presence" - Malachi 3:16 (WEB)
// Principle: Faithful witness through complete remembrance of what occurred
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Substrate monitoring (hooks infrastructure, companion to analysis.go)
// Role: Writes timestamped log files recording Claude substrate behavior for pattern analysis
// Paradigm: Simple log file writing for substrate monitoring (temporary this season)
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-11 - Template alignment and configuration integration
//
// Version History:
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded paths and formats
//   2.0.0 (2025-11-11) - Template alignment, configuration loading, comprehensive inline comments
//
// Purpose & Function
//
// Purpose: Write timestamped log files to ~/.claude/debug/ recording Claude substrate behavior
// (compaction events, notifications, subagent executions, prompt submissions) for analysis.go
// to detect patterns indicating workflow inefficiency or configuration issues.
//
// Core Design: Simple append-only log file writing with timestamp prefixing. Companion to
// analysis.go which reads these logs. Intentionally separate from system/lib/logging (which
// tracks CPI-SI system components) - this tracks temporary Claude substrate behavior.
//
// Key Features:
//   - Timestamped log entries with [YYYY-MM-DD HH:MM:SS] prefix format
//   - Separate log files per event type (compaction, notifications, subagents, prompts)
//   - Privacy-aware prompt logging (first 100 chars only)
//   - Automatic log directory creation
//   - Non-blocking writes (failures silent - never interrupt workflow)
//   - Simple format for easy grepping and analysis
//
// Philosophy: Substrate monitoring logs are temporary (this season with Claude). When
// transitioning to dedicated CPI-SI architecture, these logs become unnecessary. Keep
// separate from system logging which is permanent CPI-SI infrastructure.
//
// Blocking Status
//
// Non-blocking: All write operations fail silently if directory creation or file writing
// encounters problems. Substrate monitoring is advisory only - never interrupt actual work
// to record that work happened.
//
// Mitigation: Missing logs mean analysis.go won't detect patterns, but workflow continues
// normally. Graceful degradation is more important than complete logging.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/monitoring"
//
//	// Log compaction events
//	monitoring.LogCompaction("auto")      // or "manual"
//
//	// Log notification events
//	monitoring.LogNotification("permission_request")
//
//	// Log subagent completions
//	monitoring.LogSubagentCompletion("Explore", "success", "0")
//
//	// Log prompt submissions (privacy-aware)
//	monitoring.LogPrompt(userPromptText)
//
//	// Read log file for analysis
//	content, err := monitoring.ReadLogFile("compaction.log")
//
// Integration Pattern:
//   1. Import monitoring library into hooks
//   2. Call appropriate Log* function at event occurrence
//   3. Library writes timestamped entry to correct log file
//   4. analysis.go periodically reads logs to detect patterns
//   5. No cleanup needed - stateless operations
//
// Public API (in typical usage order):
//
//   Log Writing (recording events):
//     LogCompaction(type string) - Record compaction events
//     LogNotification(type string) - Record notification events
//     LogSubagentCompletion(agentType, status, exitCode string) - Record subagent executions
//     LogPrompt(prompt string) - Record prompt submissions (first 100 chars)
//
//   Log Reading (for analysis):
//     ReadLogFile(filename string) (string, error) - Read entire log file content
//
//   Helpers (internal):
//     getLogDir() string - Get log directory path
//     ensureLogDir() error - Create log directory if missing
//     writeLogEntry(filename, entry string) - Write timestamped log entry
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os, path/filepath, time, encoding/json, strings
//   External: None
//   Internal: None (bottom-rung library)
//
// Dependents (What Uses This):
//   Libraries: analysis.go (same package - reads these log files)
//   Hooks: session/start, prompt/submit, tool/post-use, session/notification (writes logs)
//
// Integration Points:
//   - Ladder: Bottom rung (no dependencies), used by analysis.go (mid rung) and hooks (top rung)
//   - Baton: Hooks call Log* functions → writes to files → analysis.go reads files
//   - Rails: Independent of system logging - substrate monitoring separate concern
//
// Health Scoring
//
// Health Scoring Map (Base100):
//
//   Log Writing Operations (Total = 100):
//     - Directory creation: +10 (ensure log directory exists)
//     - File opening: +20 (open log file for appending)
//     - Timestamp formatting: +10 (create timestamp prefix)
//     - Entry writing: +40 (write formatted entry to file)
//     - File closing: +20 (close file handle cleanly)
//
//   Degraded States:
//     - Directory creation fails: Continue silently (no logging possible)
//     - File open fails: Continue silently (that log file unavailable)
//     - Write fails: Continue silently (entry lost)
//
// Note: Failures don't interrupt workflow - logging is advisory only. Missing logs reduce
// pattern detection capability but don't affect core operations.
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
// provides Go's built-in capabilities. This library has no internal dependencies
// (bottom rung on ladder).
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"fmt"           // String formatting for log entry composition
	"os"            // File operations for log writing
	"path/filepath" // Path manipulation for log directory and config file paths
	"strings"       // String operations for path manipulation and tilde expansion
	"time"          // Timestamp generation for log entries
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

// None needed - all values loaded from configuration files or have sensible fallbacks

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Configuration Types ---
// Configuration types defined in analysis.go (shared by both analysis.go and logging.go):
// - LogFormatsConfig: Log file formats and filenames
// - DebugFormatsConfig: Debug file formats and preview lengths
// - RetentionPolicyConfig: Directory paths and retention rules
//
// Package-level variables and init() function also in analysis.go (shared state):
// - logFormatsConfig, debugFormatsConfig, retentionConfig, configLoaded
// - init() loads all configs at package initialization

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
//   Public APIs (Top Rungs - Log Writing)
//   ├── LogCompaction() → uses writeLogEntry(), logFormatsConfig
//   ├── LogNotification() → uses writeLogEntry(), logFormatsConfig
//   ├── LogSubagentCompletion() → uses writeLogEntry(), logFormatsConfig
//   └── LogPrompt() → uses writeLogEntry(), logFormatsConfig, debugFormatsConfig (preview length)
//
//   Utility Functions (Middle Rungs - Log Reading)
//   └── ReadLogFile() → uses getLogDir(), os.ReadFile()
//
//   Helpers (Lower Rungs - Infrastructure)
//   ├── getLogDir() → uses retentionConfig (directory path)
//   ├── ensureLogDir() → uses getLogDir(), os.MkdirAll()
//   ├── writeLogEntry() → uses ensureLogDir(), logFormatsConfig (timestamp format)
//   └── loadLogFormats(), loadDebugFormats(), loadRetentionPolicy() → JSONC config loading
//
// APUs (Available Processing Units):
// - 11 functions total
// - 4 public log writing APIs (event recording)
// - 1 public log reading API (file access)
// - 3 helper functions (directory/file operations)
// - 3 config loader functions (JSONC parsing)
//
// Configuration Flow:
// - init() calls loaders at package initialization
// - logFormatsConfig, debugFormatsConfig, retentionConfig populated
// - Public APIs use config if available, fallback to hardcoded defaults
//
// External Dependencies:
// - analysis.go (same package): Uses ReadLogFile() and expects specific log formats
// - Hooks: Call Log* functions to record substrate events
//
// See: standards/code/4-block/sections/CWS-SECTION-006-BODY-organizational-chart.md

// ────────────────────────────────────────────────────────────────
// Helpers - Configuration Loading
// ────────────────────────────────────────────────────────────────
// Configuration loader functions defined in analysis.go (shared by both files):
// - loadLogFormats(path string) *LogFormatsConfig
// - loadDebugFormats(path string) *DebugFormatsConfig
// - loadRetentionPolicy(path string) *RetentionPolicyConfig
//
// All loaders follow same pattern: Read file → Strip JSONC comments → Parse JSON → Return nil on error

// ────────────────────────────────────────────────────────────────
// Helpers - Directory Operations
// ────────────────────────────────────────────────────────────────
// Infrastructure functions for log directory management

// getLogDir returns the debug log directory path
// Uses configuration if available, falls back to hardcoded ~/.claude/debug
func getLogDir() string {
	// Try to use configured directory path
	if configLoaded && retentionConfig != nil {  // Check if config available
		path := retentionConfig.Archive.Location.FullPath  // Get configured path
		if path != "" {  // Check if path configured
			// Expand ~ to home directory if present
			if strings.HasPrefix(path, "~/") {  // Check if path starts with ~/
				home := os.Getenv("HOME")  // Get HOME environment variable
				if home != "" {  // Check if HOME available
					path = filepath.Join(home, path[2:])  // Replace ~ with $HOME
				}
			}
			// Remove trailing slash if present - normalize path format
			path = strings.TrimSuffix(path, "/")  // Remove trailing /
			return path  // Return configured directory path
		}
	}

	// Fall back to hardcoded default if config unavailable
	return filepath.Join(os.Getenv("HOME"), ".claude", "debug")  // Default: ~/.claude/debug
}

// ensureLogDir creates the log directory if it doesn't exist
// Returns nil on success or if directory already exists, error if creation fails
func ensureLogDir() error {
	// Create directory with parents if needed, ignore if already exists
	return os.MkdirAll(getLogDir(), 0755)  // Create with rwxr-xr-x permissions
}

// ────────────────────────────────────────────────────────────────
// Helpers - Log Writing
// ────────────────────────────────────────────────────────────────
// Core log entry writing functionality

// writeLogEntry writes a timestamped entry to the specified log file
// Fails silently if directory creation or file writing encounters problems
func writeLogEntry(filename, entry string) {
	// Ensure log directory exists before attempting to write
	if err := ensureLogDir(); err != nil {  // Try to create directory
		return  // Silently fail if directory creation fails - non-blocking
	}

	// Build full path to log file
	logFile := filepath.Join(getLogDir(), filename)  // Combine directory and filename

	// Open file for appending (create if doesn't exist)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)  // Open with rw-r--r-- permissions
	if err != nil {  // Check if open failed
		return  // Silently fail if can't open file - non-blocking
	}
	defer f.Close()  // Ensure file closed when function exits

	// Get timestamp format from configuration or use fallback
	timestampFormat := "2006-01-02 15:04:05"  // Default Go time format (YYYY-MM-DD HH:MM:SS)
	if configLoaded && logFormatsConfig != nil {  // Check if config available
		if logFormatsConfig.Timestamp.Format != "" {  // Check if format configured
			timestampFormat = logFormatsConfig.Timestamp.Format  // Use configured format
		}
	}

	// Format current time using configured timestamp format
	timestamp := time.Now().Format(timestampFormat)  // Create timestamp string

	// Write formatted entry: [timestamp] entry\n
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, entry))  // Write to file (ignore errors - non-blocking)
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Log Writing
// ────────────────────────────────────────────────────────────────
// Functions for recording substrate behavior events

// LogSubagentCompletion logs subagent execution details
// Records agent type, execution status, and exit code to subagents.log
func LogSubagentCompletion(agentType, status, exitCode string) {
	// Get log filename from configuration or use fallback
	filename := "subagents.log"  // Default filename
	if configLoaded && logFormatsConfig != nil {  // Check if config available
		if logFormatsConfig.SubagentsLog.Filename != "" {  // Check if filename configured
			filename = logFormatsConfig.SubagentsLog.Filename  // Use configured filename
		}
	}

	// Format log entry with key=value pairs for easy grepping
	entry := fmt.Sprintf("type=%s status=%s exitCode=%s", agentType, status, exitCode)  // Structured format
	writeLogEntry(filename, entry)  // Write timestamped entry to log file
}

// LogNotification logs notification events for pattern analysis
// Records notification type to notifications.log
func LogNotification(notificationType string) {
	// Get log filename from configuration or use fallback
	filename := "notifications.log"  // Default filename
	if configLoaded && logFormatsConfig != nil {  // Check if config available
		if logFormatsConfig.NotificationsLog.Filename != "" {  // Check if filename configured
			filename = logFormatsConfig.NotificationsLog.Filename  // Use configured filename
		}
	}

	// Log notification type directly - simple format for pattern matching
	writeLogEntry(filename, notificationType)  // Write notification type to log
}

// LogCompaction logs compaction events
// Records compaction type (auto or manual) to compaction.log
func LogCompaction(compactType string) {
	// Get log filename from configuration or use fallback
	filename := "compaction.log"  // Default filename
	if configLoaded && logFormatsConfig != nil {  // Check if config available
		if logFormatsConfig.CompactionLog.Filename != "" {  // Check if filename configured
			filename = logFormatsConfig.CompactionLog.Filename  // Use configured filename
		}
	}

	// Log compaction type directly - analysis.go searches for "auto" vs "manual"
	writeLogEntry(filename, compactType)  // Write compaction type to log
}

// LogPrompt logs user prompt submission (first 100 chars for privacy)
// Privacy-aware: Only records preview of prompt, not full content
func LogPrompt(prompt string) {
	if prompt == "" {  // Check if prompt is empty
		return  // Skip empty prompts - nothing to log
	}

	// Get log filename and preview length from configuration or use fallbacks
	filename := "prompts.log"  // Default filename
	maxLength := 100  // Default preview length (characters)

	if configLoaded {  // Check if config available
		// Get filename from log formats config
		if logFormatsConfig != nil && logFormatsConfig.PromptsLog.Filename != "" {  // Check if filename configured
			filename = logFormatsConfig.PromptsLog.Filename  // Use configured filename
		}

		// Get preview length from debug formats config
		if debugFormatsConfig != nil && debugFormatsConfig.PromptsDebug.Schema.Preview.MaxLength > 0 {  // Check if max length configured
			maxLength = debugFormatsConfig.PromptsDebug.Schema.Preview.MaxLength  // Use configured preview length
		}
	}

	// Create preview - truncate if longer than maxLength
	preview := prompt  // Start with full prompt
	if len(preview) > maxLength {  // Check if prompt exceeds max length
		preview = preview[:maxLength] + "..."  // Truncate and add ellipsis
	}

	// Write preview to log file
	writeLogEntry(filename, preview)  // Log prompt preview only (privacy-aware)
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Log Reading
// ────────────────────────────────────────────────────────────────
// Functions for accessing log file content (used by analysis.go)

// ReadLogFile reads and returns log file contents
// Returns file content as string, or error if file doesn't exist or can't be read
func ReadLogFile(filename string) (string, error) {
	// Build full path to log file
	logFile := filepath.Join(getLogDir(), filename)  // Combine directory and filename

	// Read entire file content into memory
	data, err := os.ReadFile(logFile)  // Read file (returns error if missing)
	if err != nil {  // Check if read failed
		return "", err  // Return empty string and error
	}

	// Return file content as string
	return string(data), nil  // Convert bytes to string and return
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
//   - Call Log* functions and verify log files created
//   - Verify log entries have correct timestamp format
//   - Verify LogPrompt truncates to configured length
//   - Call ReadLogFile and verify content matches what was written
//   - Verify graceful handling of missing directories (created automatically)
//   - Verify graceful handling of write failures (silent, non-blocking)
//
// Example validation code:
//
//     // Test import
//     import "hooks/lib/monitoring"
//
//     // Test log writing
//     monitoring.LogCompaction("auto")
//     monitoring.LogNotification("permission_request")
//     monitoring.LogSubagentCompletion("Explore", "success", "0")
//     monitoring.LogPrompt("Test prompt for logging")
//
//     // Verify files created
//     // Expected: Files exist at ~/.claude/debug/*.log
//
//     // Test log reading
//     content, err := monitoring.ReadLogFile("compaction.log")
//     // Expected: Returns file content with timestamps
//
//     // Test privacy-aware prompt logging
//     longPrompt := strings.Repeat("x", 200)  // 200 char string
//     monitoring.LogPrompt(longPrompt)
//     content, _ = monitoring.ReadLogFile("prompts.log")
//     // Expected: Entry is 100 chars + "..." (103 chars total)

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
// The library is imported into hook packages (session/start, prompt/submit, tool/post-use,
// session/notification), making all exported functions available. The init() function runs
// automatically when the package is imported, loading configuration files.
//
// Example usage from hooks:
//
//     package sessionstart
//
//     import "hooks/lib/monitoring"
//
//     func ExecuteSessionStart() {
//         // Record session start in compaction log
//         monitoring.LogCompaction("session_start")
//
//         // Check for patterns
//         monitoring.CheckCompactionFrequency()  // From analysis.go
//     }

// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - File handles: Opened in writeLogEntry, closed via defer before function exit
//   - No network connections: Pure local file operations
//   - No database connections: No persistence layer
//   - No goroutines: Synchronous operations only
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code (hooks) responsible for managing their own lifecycle
//   - This library is stateless - no cleanup function needed
//
// Error State Cleanup:
//   - Write failures silently ignored (non-blocking design)
//   - No partial state to corrupt (each write is independent)
//   - Missing directories created automatically on next write attempt
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Log file content read into memory temporarily (freed after ReadLogFile returns)
//   - Configuration loaded once at init() (persistent for package lifetime)
//   - No persistent allocations beyond config

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════

// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// For Library Overview section explanation, see: standards/code/4-block/sections/CWS-SECTION-013-CLOSING-library-overview.md
//
// Purpose: See METADATA "Purpose & Function" section above - writes timestamped log files
// for Claude substrate behavior monitoring
//
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Timestamped log entry writing to ~/.claude/debug/
//   - Separate log files per event type (compaction, notifications, subagents, prompts)
//   - Privacy-aware prompt logging (truncated to 100 chars)
//   - Configuration-driven file names and formats
//   - Non-blocking writes (failures silent)
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide. Simple import and call pattern.
//
// Public API: See METADATA "Public API" section above for complete function list
// organized by purpose (log writing, log reading)
//
// Architecture: See METADATA "CPI-SI Identity" section above - this is substrate
// monitoring (temporary), separate from system/lib/logging (permanent)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new Log* functions for additional event types
//   ✅ Add new configuration options (formats, paths, permissions)
//   ✅ Extend privacy-aware logging to other event types
//   ✅ Add structured logging variants (JSON, CSV)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Function signatures - hooks depend on these interfaces
//   ⚠️ Log file format - analysis.go expects [timestamp] entry format
//   ⚠️ Configuration structure - changing JSON keys breaks config loading
//   ⚠️ Directory path defaults - changing ~/.claude/debug breaks existing integrations
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking behavior (writes must never crash workflow)
//   ❌ Separation from system/lib/logging (substrate vs system logs)
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
// - 4 public log writing APIs write to appropriate log files
// - 1 public log reading API reads file content
// - 3 helper functions manage directories and write entries
// - Ladder: hooks (top) → logging (bottom) - no dependencies
// - Baton: hooks call Log* → writeLogEntry → file writes
// - Rails: Independent of system logging - substrate monitoring separate

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Public APIs - Log Writing" section header above for
// detailed extension points.
//
// Quick reference:
// - Adding Log* functions: Add to "Public APIs - Log Writing" section,
//   follow existing pattern (get filename from config, call writeLogEntry)
// - Adding config options: Update type structs in SETUP, add to loader functions
// - Changing log formats: Update log-formats.jsonc, template uses configured format

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// All operations are lightweight file appends with minimal overhead.
//
// - Most expensive operation: writeLogEntry() - O(1) file append
// - Memory characteristics: Minimal (timestamp string + entry string)
// - File I/O: One open + write + close per log entry
// - Directory check: Cached by OS after first ensureLogDir()
//
// Note: Performance is not a concern for this library - substrate monitoring writes
// happen infrequently (few per minute) where latency is imperceptible.
//
// Optimization opportunities (if needed):
// - Batched writes (collect entries, write in bulk)
// - Async writes (goroutine pool for file I/O)
// - Memory-mapped files (for very high frequency logging)

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: Log files not being created
//   Check: Verify ~/.claude/debug/ directory exists and is writable
//   Check: Verify HOME environment variable set correctly
//   Check: Verify hooks are actually calling Log* functions
//   Solution: Run ls -la ~/.claude/debug/ to see files and permissions
//
// Problem: Log entries missing timestamps
//   Check: Verify time.Now() working (system clock issues?)
//   Check: Verify timestamp format configured correctly
//   Solution: Check log-formats.jsonc timestamp.format field
//
// Problem: Prompt logs showing full content instead of preview
//   Check: Verify debugFormatsConfig loaded successfully
//   Check: Verify prompts_debug.schema.preview.max_length configured
//   Solution: Check debug-formats.jsonc and verify config loading
//
// Problem: ReadLogFile returns empty string
//   Check: Verify log file exists at ~/.claude/debug/
//   Check: Verify log file has content (cat ~/.claude/debug/compaction.log)
//   Check: Verify correct filename passed to ReadLogFile
//   Solution: Ensure hooks are writing logs before trying to read them

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key companion: analysis.go (reads these log files)
// - Primary consumers: hooks/session/start, hooks/prompt/submit, hooks/tool/post-use
// - Configuration files: log-formats.jsonc, debug-formats.jsonc, retention-policy.jsonc
// - Intentionally separate from system/lib/logging (different purpose)

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration-driven formats and filenames - COMPLETED (v2.0.0)
//   ✓ Template alignment with comprehensive inline comments - COMPLETED (v2.0.0)
//   ⏳ Structured debug files (.debug) in JSON format (v3.0.0)
//   ⏳ Log rotation based on retention-policy.jsonc (v3.0.0)
//   ⏳ Compression of old log files (v3.0.0)
//
// Research Areas:
//   - Async writing (goroutine pool for non-blocking I/O)
//   - Structured logging (JSON, CSV, binary formats)
//   - Log streaming (real-time analysis without file reads)
//   - Automatic cleanup based on retention policy
//
// Integration Targets:
//   - Hook into retention policy enforcement
//   - Write both .log (human-readable) and .debug (machine-parseable) simultaneously
//   - Integration with analysis.go for real-time pattern detection
//
// Known Limitations to Address:
//   - No log rotation (files grow unbounded)
//   - No compression (disk space not managed)
//   - No .debug files yet (only .log files written)
//   - Synchronous writes (could be async for performance)
//   - No batching (one write per log entry)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history:
//
//   2.0.0 (2025-11-11) - Configuration Architecture & Template Alignment
//         - Comprehensive METADATA block with all 8 subsections
//         - SETUP block with section headers and inline comments
//         - BODY block with organizational chart and inline comments
//         - CLOSING block with FINAL DOCUMENTATION
//         - Configuration loading (log-formats, debug-formats, retention-policy)
//         - Graceful fallback to hardcoded defaults
//         - Maintains same external API (hooks see no breaking changes)
//
//   1.0.0 (2024-10-24) - Initial Implementation
//         - Basic timestamped log writing
//         - Hardcoded paths and formats
//         - Privacy-aware prompt logging
//         - Log file reading for analysis

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library writes CLAUDE SUBSTRATE behavior logs specifically - not CPI-SI system
// logs. Separation is intentional: substrate monitoring is temporary (this season),
// system logging is permanent CPI-SI infrastructure.
//
// Companion to analysis.go which reads these logs to detect workflow patterns. Together
// they provide substrate behavior monitoring without interfering with operations.
//
// Modify thoughtfully - hooks and analysis.go depend on these interfaces and log formats.
// Keep writes non-blocking - logging should never become the cause of failures.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test integration with hooks after changes
//   - Maintain graceful degradation principles
//
// "A scroll of remembrance was written in his presence" - Malachi 3:16 (WEB)

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage from Hooks:
//
//     import "hooks/lib/monitoring"
//
//     // Log compaction events
//     monitoring.LogCompaction("auto")  // or "manual"
//
//     // Log notification events
//     monitoring.LogNotification("permission_request")
//
//     // Log subagent completions
//     monitoring.LogSubagentCompletion("Explore", "success", "0")
//
//     // Log prompt submissions (privacy-aware)
//     monitoring.LogPrompt(userPromptText)
//
// Reading Logs for Analysis:
//
//     // Read entire log file
//     content, err := monitoring.ReadLogFile("compaction.log")
//     if err != nil {
//         // Handle error (file missing or unreadable)
//     }
//
//     // Parse log entries
//     lines := strings.Split(content, "\n")
//     for _, line := range lines {
//         // Process each [timestamp] entry line
//     }
//
// Integration in Session Start Hook:
//
//     package sessionstart
//
//     import "hooks/lib/monitoring"
//
//     func ExecuteSessionStart() {
//         // Record session start
//         monitoring.LogCompaction("session_start")
//
//         // Check for substrate behavior patterns
//         monitoring.CheckCompactionFrequency()  // From analysis.go
//         monitoring.CheckNotificationPatterns("permission_request")
//     }

// ============================================================================
// END CLOSING
// ============================================================================
