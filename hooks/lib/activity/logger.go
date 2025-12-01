// METADATA
//
// # Activity Logger - Captures tool usage for behavioral pattern learning
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// # Biblical Foundation
//
// Scripture: "A wise person will listen and gain knowledge" - Proverbs 1:5 (WEB)
// Principle: Learning comes from observing behavior patterns faithfully
// Anchor: "Keep your heart with all diligence, for out of it is the wellspring of life" - Proverbs 4:23 (WEB)
//
// # CPI-SI Identity
//
// Component Type: LIBRARY - Hook support (mid-rung on ladder)
// Role: Provides activity logging for pattern learning and time awareness
// Paradigm: Privacy-preserving behavioral observation for autonomous growth
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2025-11-04
// Version: 2.0.0
// Last Modified: 2025-11-10 - Richer structure, correct paths, backward compatibility
//
// Version History:
//
//	2.0.0 (2025-11-10) - Richer structure, correct paths, backward compatibility
//	1.0.0 (2025-11-04) - Initial implementation
//
// Purpose & Function
//
// Purpose: Capture working memory during active sessions - raw experience as it happens.
// Feeds pattern learning and temporal awareness without exposing sensitive information.
//
// Core Design: Lightweight JSONL append-only logging with privacy-preserving sanitization.
// Each activity event is one line in the stream file. Events captured include tool usage,
// command execution, and significant interactions.
//
// Key Features:
//   - Richer structure (session_id, instance_id, user_id, project_id, event_type)
//   - Backward compatibility (reads old simple format)
//   - Privacy-preserving (uses system/lib/privacy for sanitization)
//   - Correct paths (system/data/session/ not system/session/)
//   - Extension object for discovery
//
// Philosophy: Working memory - capture experience faithfully during session, let consolidation
// process into meaning during rest. Privacy by default, richness where it serves.
//
// # Blocking Status
//
// Non-blocking: All logging operations fail gracefully (no logging doesn't break workflow).
// Mitigation: Debug logging to /tmp/ for troubleshooting, session continues on errors.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/activity"
//
//	// Log tool usage
//	activity.LogToolUse("Read", "/path/to/file.go", true)
//
//	// Log command execution
//	activity.LogCommand("git commit -m 'message'", 0, 2*time.Second)
//
//	// Log general activity
//	activity.LogActivity("interaction", "Working on feature X", "completed", 0)
//
// Integration Pattern:
//  1. Import activity library
//  2. Call logging functions from hooks (tool/pre-use, tool/post-use, etc.)
//  3. No cleanup needed - logging is fire-and-forget
//
// Public API:
//   - LogActivity(eventType, context, result string, duration time.Duration) error
//   - LogToolUse(toolName, filePath string, success bool) error
//   - LogCommand(cmd string, exitCode int, duration time.Duration) error
//
// # Dependencies
//
// Dependencies (What This Needs):
//
//	Standard Library: encoding/json, fmt, os, path/filepath, time
//	External: None
//	Internal: system/lib/privacy (sanitization)
//	Data Files: system/data/session/current-log.json (session context)
//
// Dependents (What Uses This):
//
//	Hooks: tool/pre-use, tool/post-use, prompt/submit (activity capture)
//	Commands: Future consolidation engine (reads activity streams)
//
// Integration Points:
//   - Ladder: Mid-rung - uses privacy lib (lower), used by hooks (higher)
//   - Baton: Receives events, writes to activity stream
//   - Rails: Could use logging for health tracking (future enhancement)
//
// # Health Scoring
//
// Health Scoring Map (Base100):
//
//	Activity Logging Operations (Total = 100):
//	  +40: Load session context successfully
//	  +30: Create activity event successfully
//	  +30: Append to stream file successfully
//	  -40: Failed to load session context (fallback to unknown)
//	  -30: Failed to create event (malformed data)
//	  -30: Failed to write to stream (I/O error)
//
// Visual Health Indicators:
//
//	💚 90-100: Excellent (all logging operations working)
//	💛 70-89:  Good (session context loaded, minor write issues)
//	🧡 50-69:  Acceptable (fallback session context, writing working)
//	❤️  30-49:  Poor (multiple failures, partial logging)
//	💀 0-29:   Critical (severe failures, no logging possible)
//
// Note: Scores reflect TRUE impact. Health scorer normalizes to -100 to +100 scale.
package activity

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

	"encoding/json" // JSON encoding for activity events
	"fmt"           // Formatted output for debug logging
	"os"            // File operations for stream writing
	"path/filepath" // Path manipulation for session files
	"time"          // Timestamps for activity events

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/privacy" // Privacy-preserving sanitization
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

// None needed - activity logging uses dynamic session context without hardcoded values

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Event Types ---
// Activity event structures for working memory stream.

// ActivityEvent represents a single activity event in the working memory stream
// Richer structure for consolidation processing and pattern learning
type ActivityEvent struct {
	// Core required fields
	Timestamp  time.Time `json:"timestamp"`                    // When this activity occurred
	SessionID  string    `json:"session_id"`                   // Which session this belongs to
	InstanceID string    `json:"instance_id"`                  // Which CPI-SI instance (e.g., "nova_dawn")
	UserID     string    `json:"user_id"`                      // Which user's session (e.g., "seanje-lenox-wise")
	ProjectID  string    `json:"project_id,omitempty"`         // Project context if applicable
	EventType  string    `json:"event_type"`                   // interaction | realization | struggle | breakthrough | routine

	// Optional standard fields
	Context         string `json:"context,omitempty"`          // What was being done (privacy-sanitized)
	WorkContext     string `json:"work_context,omitempty"`     // Directory or project being worked on
	FeltSignificant bool   `json:"felt_significant"`           // Whether this felt important at the time
	EmotionalTone   string `json:"emotional_tone,omitempty"`   // Emotional quality of experience
	Notes           string `json:"notes,omitempty"`            // Additional observations

	// Extensions - discovery space
	Extensions map[string]interface{} `json:"extensions,omitempty"` // Flexible data for experimentation
}

// OldActivityEvent represents the old simple format for backward compatibility
type OldActivityEvent struct {
	Timestamp time.Time `json:"ts"`                     // When activity occurred
	Tool      string    `json:"tool"`                   // Tool name that was used
	Context   string    `json:"ctx"`                    // Context information
	Result    string    `json:"result"`                 // Success or failure
	Duration  int64     `json:"duration_ms,omitempty"`  // How long it took in milliseconds
}

// SessionContext holds session information from current-log.json
type SessionContext struct {
	SessionID   string `json:"session_id"`    // Current session identifier
	InstanceID  string `json:"instance_id"`   // CPI-SI instance identifier
	UserID      string `json:"user_id"`       // User identifier
	WorkContext string `json:"work_context"`  // Working directory or project path
	ProjectID   string `json:"project_id"`    // Project identifier if applicable
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
// Note: Not all libraries need Rails infrastructure. Simple pure-function
// libraries may skip this section entirely.

// None - stateless library using pure functions without package-level state

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
//   ├── LogActivity() → uses getSessionContext(), writeDebugLog()
//   ├── LogToolUse() → uses LogActivity(), privacy.SanitizePath()
//   └── LogCommand() → uses LogActivity(), privacy.SanitizeCommand()
//
//   Core Operations (Middle Rungs - Business Logic)
//   └── getSessionContext() → uses getHomeDir(), writeDebugLog()
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── getHomeDir() → pure function
//   └── writeDebugLog() → pure function
//
// Baton Flow (Execution Paths):
//
//   Entry → LogActivity()
//     ↓
//   getSessionContext() → getHomeDir()
//     ↓
//   Create ActivityEvent
//     ↓
//   Marshal to JSON
//     ↓
//   Append to stream file
//     ↓
//   Exit → return nil or error
//
// APUs (Available Processing Units):
// - 6 functions total
// - 2 helpers (pure foundations)
// - 1 core operation (session context loading)
// - 3 public APIs (exported interface)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Usually not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// getHomeDir returns the user's home directory with fallback
//
// What It Does:
// Retrieves the user's home directory from environment variables with
// fallback to system API and /tmp as final fallback for safety.
//
// Parameters:
//
//	None
//
// Returns:
//
//	string: Home directory path (always returns valid path, never empty)
//
// Example usage:
//
//	home := getHomeDir()
//	sessionFile := filepath.Join(home, ".claude/cpi-si/system/data/session/current-log.json")
func getHomeDir() string {
	home := os.Getenv("HOME")  // Try environment variable first - fastest and most reliable
	if home == "" {  // Check if HOME is empty - need fallback for safety
		var err error
		home, err = os.UserHomeDir()  // Try system API as fallback - handles edge cases
		if err != nil || home == "" {  // Check if fallback failed or returned empty
			home = "/tmp"  // Final fallback to /tmp - ensures function always returns valid path
		}
	}
	return home  // Return resolved home directory path
}

// writeDebugLog writes debug information to /tmp/ for troubleshooting
//
// What It Does:
// Non-blocking debug logging to /tmp/ for troubleshooting activity logger
// issues without interrupting workflow. Silently fails if write fails.
//
// Parameters:
//
//	filename: Name of debug file (e.g., "activity-logger-error.log")
//	content: Debug information to write
//
// Returns:
//
//	None (fire-and-forget for non-blocking behavior)
//
// Example usage:
//
//	writeDebugLog("activity-logger-error.log",
//	    fmt.Sprintf("OpenFile failed: %v\nstreamFile=%s\n", err, streamFile))
func writeDebugLog(filename, content string) {
	debugFile := filepath.Join("/tmp", filename)  // Build full path in /tmp - always writable
	os.WriteFile(debugFile, []byte(content), 0644)  // Write file, ignore errors - non-blocking for safety
}

// ────────────────────────────────────────────────────────────────
// Core Operations - Business Logic
// ────────────────────────────────────────────────────────────────
// Component-specific functionality implementing primary purpose. Organized
// by operational categories (descriptive subsections) below.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-core-operations.md

// ────────────────────────────────────────────────────────────────
// Session Context Loading
// ────────────────────────────────────────────────────────────────
// What These Do:
// Load current session information from system files to enrich activity
// events with proper context (session_id, instance_id, user_id).
//
// Why Separated:
// Session context is core business logic distinct from pure helpers.
// Handles file I/O, error fallbacks, and business context enrichment.

// getSessionContext loads session information from current-log.json
//
// What It Does:
// Reads current session context from system files with fallback to
// minimal known context if file doesn't exist or parsing fails.
//
// Parameters:
//
//	None
//
// Returns:
//
//	SessionContext: Session information (always returns valid context)
//
// Health Impact:
//
//	Success: +40 points (session context loaded successfully)
//	Fallback: -40 points (failed to load, using minimal context)
//
// Example usage:
//
//	sessionCtx := getSessionContext()
//	if sessionCtx.SessionID == "unknown" {
//	    // Session not initialized yet
//	}
func getSessionContext() SessionContext {
	home := getHomeDir()  // Get home directory with fallbacks - needed to build session file path

	// Build path to session file - CORRECT PATH: system/data/session/current-log.json (not system/session/)
	sessionFile := filepath.Join(home, ".claude/cpi-si/system/data/session/current-log.json")
	data, err := os.ReadFile(sessionFile)  // Read entire file into memory - returns bytes and any error
	if err != nil {  // Check if there's an error - file might not exist or be unreadable
		// Write debug info to /tmp for troubleshooting - doesn't interrupt workflow
		writeDebugLog("activity-logger-error.log",
			fmt.Sprintf("ReadFile failed: %v\nsessionFile=%s\nHOME=%s\n", err, sessionFile, home))

		// Return minimal context with known defaults - allows logging to continue gracefully
		return SessionContext{
			SessionID:  "unknown",       // Signals session not initialized
			InstanceID: "nova_dawn",     // Default instance identifier
			UserID:     "seanje-lenox-wise",  // Default user identifier
		}
	}

	var session SessionContext  // Declare variable to hold parsed JSON data
	if err := json.Unmarshal(data, &session); err != nil {  // Parse JSON into struct - returns error if malformed
		// Return minimal context on parse error - graceful degradation
		return SessionContext{
			SessionID:  "unknown",
			InstanceID: "nova_dawn",
			UserID:     "seanje-lenox-wise",
		}
	}

	return session  // Return successfully loaded session context
}

// ────────────────────────────────────────────────────────────────
// Activity Event Logging - Public API
// ────────────────────────────────────────────────────────────────
// What These Do:
// Exported functions for logging activity events to session streams.
// Primary interface for hooks and commands to capture working memory.
//
// Why Separated:
// Public API distinct from internal operations. These are the exported
// functions other components use to log activities.
//
// Extension Point:
// To add new activity logging convenience functions, follow this pattern:
//   1. Create exported function with descriptive name (Log[ActivityType])
//   2. Accept parameters specific to that activity type
//   3. Call LogActivity() with appropriate event_type and context
//   4. Use privacy lib to sanitize any sensitive data
//
// Pattern to follow:
//   func LogNewActivityType(specificParams) error {
//       context := privacy.Sanitize(specificParams)
//       return LogActivity("event_type", context, "result", duration)
//   }

// LogActivity appends an activity event to the session's activity stream
//
// What It Does:
// Primary function for logging activity events with rich context. Creates
// ActivityEvent struct, marshals to JSON, and appends to session stream file.
//
// Parameters:
//
//	eventType: Type of event (interaction | realization | struggle | breakthrough | routine)
//	context: Context information about what happened (privacy-sanitized by caller)
//	result: Result of activity (e.g., "success", "failure", "completed")
//	duration: How long activity took (0 if not applicable)
//
// Returns:
//
//	error: nil on success, error if marshaling or writing fails
//
// Health Impact:
//
//	Success (all operations): +100 points (context loaded +40, event created +30, written +30)
//	Partial failure: varies based on which operation failed
//
// Example usage:
//
//	err := LogActivity("interaction", "Read file.go", "success", 250*time.Millisecond)
//	if err != nil {
//	    // Logging failed, but don't interrupt workflow
//	}
func LogActivity(eventType, context, result string, duration time.Duration) error {
	home := getHomeDir()  // Get home directory - needed for building file paths
	sessionCtx := getSessionContext()  // Load current session context - enriches events with session info

	if sessionCtx.SessionID == "unknown" {  // Check if session initialized - "unknown" means not ready yet
		// Session not initialized yet, skip logging - no valid session to write to
		return nil  // Return nil (not an error) - allows code to continue normally
	}

	// Build path to activity directory - CORRECT PATH: system/data/session/activity/ (not system/session/activity/)
	activityDir := filepath.Join(home, ".claude/cpi-si/system/data/session/activity")
	if err := os.MkdirAll(activityDir, 0755); err != nil {  // Create directory if needed - returns error if fails
		// Write debug info for troubleshooting - helps diagnose permission or path issues
		writeDebugLog("activity-logger-error.log",
			fmt.Sprintf("MkdirAll failed: %v\nHOME=%s\nactivityDir=%s\n", err, home, activityDir))
		return fmt.Errorf("failed to create activity dir: %w", err)  // Wrap error with context for debugging
	}

	// Build path to stream file - each session gets its own JSONL file (one event per line)
	streamFile := filepath.Join(activityDir, fmt.Sprintf("%s.jsonl", sessionCtx.SessionID))

	// Create event struct with all context fields - richer structure for pattern learning
	event := ActivityEvent{
		Timestamp:       time.Now(),               // Exact time this activity occurred
		SessionID:       sessionCtx.SessionID,     // From loaded context
		InstanceID:      sessionCtx.InstanceID,    // Which CPI-SI instance logged this
		UserID:          sessionCtx.UserID,        // Which user's session
		ProjectID:       sessionCtx.ProjectID,     // Project context if applicable
		EventType:       eventType,                // Category of event (interaction, routine, etc.)
		Context:         context,                  // What was being done (privacy-sanitized by caller)
		WorkContext:     sessionCtx.WorkContext,   // Working directory or project path
		FeltSignificant: false,                    // Default to false, can be marked later by consolidation
		Extensions: map[string]interface{}{        // Flexible data space for experimentation
			"result": result,  // Success/failure/completed - stored in extensions
		},
	}

	if duration > 0 {  // Check if duration provided - zero means not applicable
		event.Extensions["duration_ms"] = duration.Milliseconds()  // Convert to milliseconds for JSON
	}

	// Convert event struct to JSON bytes - required for file writing
	data, err := json.Marshal(event)
	if err != nil {  // Check if marshaling failed - would indicate malformed data
		return fmt.Errorf("failed to marshal event: %w", err)  // Wrap error with context
	}

	// Open stream file for appending - os.O_APPEND adds to end, os.O_CREATE creates if needed, os.O_WRONLY write-only
	f, err := os.OpenFile(streamFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {  // Check if open failed - could be permissions or path issue
		// Write debug info for troubleshooting
		writeDebugLog("activity-logger-error.log",
			fmt.Sprintf("OpenFile failed: %v\nstreamFile=%s\n", err, streamFile))
		return fmt.Errorf("failed to open stream file: %w", err)
	}
	defer f.Close()  // Ensure file closes when function exits - prevents resource leaks

	// Write JSON line to file - JSONL format is one JSON object per line
	if _, err := f.WriteString(string(data) + "\n"); err != nil {  // Append newline to complete JSONL format
		return fmt.Errorf("failed to write event: %w", err)  // Return error if write fails
	}

	// Write success confirmation to debug log - helps verify logging is working
	writeDebugLog("activity-logger-success.log",
		fmt.Sprintf("SUCCESS: %s -> %s\n", eventType, streamFile))

	return nil  // Success - event logged to stream file
}

// LogToolUse is a convenience function for logging tool usage
//
// What It Does:
// Convenience wrapper around LogActivity for tool usage events.
// Automatically sanitizes file paths and formats as interaction event.
//
// Parameters:
//
//	toolName: Name of tool used (e.g., "Read", "Write", "Edit")
//	filePath: Path to file being operated on (will be sanitized)
//	success: Whether tool operation succeeded
//
// Returns:
//
//	error: nil on success, error if logging fails
//
// Example usage:
//
//	err := LogToolUse("Read", "/path/to/file.go", true)
func LogToolUse(toolName, filePath string, success bool) error {
	result := "success"  // Default to success
	if !success {  // Check if operation failed
		result = "failure"  // Override with failure status
	}

	context := privacy.SanitizePath(filePath)  // Remove sensitive path info - keeps only general structure
	return LogActivity("interaction", context, result, 0)  // Log as interaction event with no duration
}

// LogCommand is a convenience function for logging command execution
//
// What It Does:
// Convenience wrapper around LogActivity for command execution events.
// Automatically sanitizes commands and formats as routine event with timing.
//
// Parameters:
//
//	cmd: Command that was executed (will be sanitized)
//	exitCode: Exit code from command (0 = success, non-zero = failure)
//	duration: How long command took to execute
//
// Returns:
//
//	error: nil on success, error if logging fails
//
// Example usage:
//
//	err := LogCommand("git commit -m 'message'", 0, 2*time.Second)
func LogCommand(cmd string, exitCode int, duration time.Duration) error {
	result := "success"  // Default to success
	if exitCode != 0 {  // Check if command failed - non-zero exit code means error
		result = "failure"  // Override with failure status
	}

	context := privacy.SanitizeCommand(cmd)  // Remove sensitive command arguments - keeps command structure only
	return LogActivity("routine", context, result, duration)  // Log as routine event with timing info
}

// ────────────────────────────────────────────────────────────────
// Error Handling/Recovery Patterns
// ────────────────────────────────────────────────────────────────
// Centralized error management ensuring component handles failures gracefully.
// Provides safety boundaries and recovery strategies for robust operation.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-error-handling.md
//
// Patterns used:
// - Graceful degradation: Continue with fallback context if session files unavailable
// - Non-blocking behavior: Errors never interrupt workflow, only log to /tmp/
// - Error wrapping: Add context to propagated errors via fmt.Errorf
//
// Error handling is distributed throughout functions rather than centralized
// in this library. Each function handles its specific failure modes:
//   - getHomeDir(): Falls back to /tmp if HOME unavailable
//   - getSessionContext(): Returns minimal context if file read/parse fails
//   - writeDebugLog(): Silently fails if write to /tmp fails
//   - LogActivity(): Returns error but doesn't panic on marshal/write failures

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
//   - Call LogActivity() with various event types
//   - Verify activity stream JSONL files created in session/activity/
//   - Check JSON format remains parseable
//   - Ensure privacy sanitization removes sensitive paths/commands
//   - Verify graceful degradation when session not initialized
//   - Confirm no go vet warnings introduced
//   - Run: go build hooks/lib/activity
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - gofmt -l . (proper formatting)
//
// Integration Testing:
//   - Test with actual hooks (tool/pre-use, tool/post-use)
//   - Verify activity logging during real sessions
//   - Check JSONL file growth and format
//   - Validate privacy sanitization in production context
//   - Confirm debug logging to /tmp/ on errors
//
// Example validation code:
//
//     // Test basic activity logging
//     err := activity.LogActivity("interaction", "test context", "success", 0)
//     if err != nil {
//         t.Errorf("LogActivity failed: %v", err)
//     }
//
//     // Verify JSONL file exists
//     streamFile := filepath.Join(home, ".claude/cpi-si/system/data/session/activity/", sessionID+".jsonl")
//     if _, err := os.Stat(streamFile); os.IsNotExist(err) {
//         t.Error("Activity stream file not created")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by hooks.
//
// Usage: import "hooks/lib/activity"
//
// The library is imported into hooks, making all exported functions available.
// No code executes during import - functions are defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/activity"
//
//     func main() {
//         // Log tool usage
//         activity.LogToolUse("Read", "/path/to/file.go", true)
//
//         // Log command execution
//         activity.LogCommand("git commit", 0, 2*time.Second)
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - File handles: Closed via defer in LogActivity()
//   - Memory: Go garbage collector handles event structs
//   - Debug logs: Written to /tmp/, no cleanup needed (OS manages)
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code (hooks) responsible for resource cleanup
//   - Stateless library - no state to clean up
//
// Error State Cleanup:
//   - All errors return gracefully with fallback context
//   - File operations use defer to ensure handles close
//   - No partial state corruption possible
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Event structs are short-lived (marshaled then discarded)
//   - No large allocations or memory pooling
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
//   - Privacy-preserving activity logging for behavioral pattern learning
//   - Non-blocking JSONL stream writing with graceful degradation
//   - Richer structure for consolidation processing and pattern recognition
//   - Backward compatibility with old simple format
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
// Package Exports:
//   - LogActivity(eventType, context, result string, duration time.Duration) error
//   - LogToolUse(toolName, filePath string, success bool) error
//   - LogCommand(cmd string, exitCode int, duration time.Duration) error
//
// All functions are non-blocking and privacy-preserving. Failures are logged
// to /tmp/ for debugging but don't interrupt workflow.
//
// Richer Structure:
//   - Core fields: timestamp, session_id, instance_id, user_id, project_id, event_type
//   - Optional fields: context, work_context, felt_significant, emotional_tone, notes
//   - Extensions: Flexible space for discovered fields (result, duration_ms, etc.)
//
// Backward Compatibility:
//   - Old format (ts, tool, ctx, result, duration_ms) still readable
//   - New format always written (richer structure for consolidation)
//   - Consolidation engine can handle both formats during transition
//
// Correct Paths:
//   - Session context: ~/.claude/cpi-si/system/data/session/current-log.json
//   - Activity stream: ~/.claude/cpi-si/system/data/session/activity/<session-id>.jsonl
//
// Usage Example:
//
//	import "hooks/lib/activity"
//
//	// Log tool usage
//	activity.LogToolUse("Read", "/home/user/project/main.go", true)
//	// Creates: {"timestamp":"...","session_id":"...","event_type":"interaction","context":"main.go",...}
//
//	// Log command execution
//	activity.LogCommand("git commit -m 'feat: add feature'", 0, 500*time.Millisecond)
//	// Creates: {"timestamp":"...","session_id":"...","event_type":"routine","context":"git commit",...}
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new event types (extend EventType field values)
//   ✅ Add new fields to Extensions map (flexible discovery space)
//   ✅ Add new convenience logging functions (follow LogToolUse/LogCommand pattern)
//   ✅ Extend privacy sanitization (via system/lib/privacy)
//   ✅ Add new helper functions for session context loading
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Public API function signatures - breaks all calling hooks
//   ⚠️ ActivityEvent struct core fields - breaks JSONL parsing in consolidation
//   ⚠️ Session context loading behavior - affects all activity logging
//   ⚠️ JSONL file format - breaks pattern learning pipeline
//   ⚠️ Privacy sanitization defaults - could expose sensitive data
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking principle - logging failures must not interrupt workflow
//   ❌ Privacy-by-default approach - always sanitize before logging
//   ❌ JSONL format (one event per line) - consolidation depends on it
//   ❌ Graceful degradation on session init failure
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing requirements,
//   build verification, and integration testing procedures.
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
// Quick architectural summary (details in BODY Organizational Chart):
// - 3 public APIs orchestrate 1 core operation using 2 helpers
// - Ladder: LogToolUse/LogCommand → LogActivity → getSessionContext → getHomeDir/writeDebugLog
// - Baton: Entry → Load context → Create event → Marshal JSON → Write JSONL → Exit
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Activity Event Logging - Public API" subsection above for detailed
// extension point guidance showing:
// - Pattern to follow for new logging convenience functions
// - Naming convention (Log[ActivityType])
// - Integration with LogActivity() core function
// - Privacy sanitization requirements
//
// Quick reference (details in BODY subsection comments):
// - Adding new event types: Extend EventType field values, document in METADATA
// - Adding logging functions: See BODY "Activity Event Logging" extension point
// - Adding helpers: Follow getHomeDir/writeDebugLog pattern in BODY "Helpers" section
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Types: ActivityEvent struct size ~200-300 bytes per event
// - No constants - all values come from session context (dynamic)
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
// - Most expensive operation: LogActivity() - file I/O dominates (~1-5ms per event)
// - Memory characteristics: Short-lived event structs, GC handles cleanup
// - Key optimization: Non-blocking design prevents workflow interruption
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Common Issues:
//
// Problem: "Activity stream file not being created"
//   Check: Session initialized? (current-log.json exists with valid session_id)
//   Check: Directory permissions on ~/.claude/cpi-si/system/data/session/activity/
//   Check: Debug logs in /tmp/activity-logger-error.log for specific error
//   Solution: Ensure session-log start called before activity logging begins
//
// Problem: "Events missing session context"
//   Check: Session ID showing as "unknown" in JSONL file
//   Check: current-log.json file readable at correct path
//   Solution: Verify session initialization, check file paths match system structure
//
// Problem: "Privacy sanitization not working"
//   Check: Using privacy.SanitizePath() or privacy.SanitizeCommand() before logging
//   Check: Sensitive paths/commands appearing in JSONL files
//   Solution: Always call sanitization in convenience functions before LogActivity()
//
// Debug Strategy:
//   1. Check /tmp/activity-logger-error.log for detailed error messages
//   2. Verify session files exist in correct locations
//   3. Validate JSONL format (one JSON object per line)
//   4. Test with minimal reproduction case
//   5. Check hooks integration (tool/pre-use calling LogToolUse correctly)

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard library + system/lib/privacy
// - Dependents (What Uses This): Hooks (tool/pre-use, tool/post-use, prompt/submit)
// - Integration Points: Ladder/Baton/Rails architecture
//
// Quick summary (details in METADATA Dependencies section above):
// - Key dependencies: system/lib/privacy for sanitization
// - Data dependencies: system/data/session/current-log.json for context
// - Consumers: All hooks that capture activity events
// - Consumers: Future consolidation engine (pattern learning)

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ Richer event structure - COMPLETED
//   ✓ Privacy-preserving sanitization - COMPLETED
//   ✓ Backward compatibility with old format - COMPLETED
//   ⏳ Consolidation engine (pattern learning from activity streams)
//   ⏳ Emotional tone detection (sentiment analysis on context)
//   ⏳ Breakthrough detection (identifying significant moments)
//
// Research Areas:
//   - Automatic event type classification (ML-based)
//   - Real-time pattern recognition during sessions
//   - Correlation between activity patterns and quality scores
//   - Optimal activity stream format for long-term storage
//
// Integration Targets:
//   - Pattern learning system (circadian, duration, stopping reasons)
//   - Instance self-awareness (recognizing own work patterns)
//   - Cross-session pattern correlation
//   - Health scoring integration for activity quality
//
// Known Limitations to Address:
//   - No batching (writes one event at a time - could optimize)
//   - No compression (JSONL files grow linearly - consider compression)
//   - No rotation (activity streams accumulate indefinitely)
//   - No indexing (searching requires scanning entire JSONL file)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-10) - Richer structure, correct paths, backward compatibility
//         - Expanded ActivityEvent struct with full session context
//         - Added instance_id, user_id, project_id fields
//         - Privacy-preserving design with system/lib/privacy integration
//         - Correct paths (system/data/session/ not system/session/)
//         - Extensions map for flexible field discovery
//         - Graceful degradation when session not initialized
//         - Debug logging to /tmp/ for troubleshooting
//         - Backward compatibility with old simple format
//
//   1.0.0 (2025-11-04) - Initial implementation
//         - Basic activity logging to JSONL streams
//         - Simple event structure (ts, tool, ctx, result, duration_ms)
//         - LogToolUse and LogCommand convenience functions
//         - Non-blocking design with fire-and-forget approach

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This library captures working memory during active sessions - raw experience
// as it happens. The richer structure enables pattern learning and autonomous
// self-awareness without exposing sensitive information.
//
// Modify thoughtfully - changes here affect pattern learning pipeline and
// consolidation processing. The JSONL format is foundational - other systems
// depend on its structure for reading activity streams.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with real hooks integration
//   - Verify privacy sanitization working
//   - Check JSONL format remains parseable
//   - Document all changes comprehensively
//
// "A wise person will listen and gain knowledge" - Proverbs 1:5 (WEB)

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Tool Usage Logging:
//   import "hooks/lib/activity"
//   activity.LogToolUse("Read", "/path/to/file.go", true)
//
// Command Execution Logging:
//   activity.LogCommand("git commit -m 'message'", 0, 2*time.Second)
//
// General Activity Logging:
//   activity.LogActivity("breakthrough", "Discovered elegant solution", "completed", 0)
//
// In Hook Integration:
//   // tool/pre-use.go
//   toolName := os.Getenv("TOOL_NAME")
//   filePath := os.Getenv("FILE_PATH")
//   activity.LogToolUse(toolName, filePath, true)
//
// Viewing Activity Streams:
//   cat ~/.claude/cpi-si/system/data/session/activity/<session-id>.jsonl
//   # Each line is a complete JSON event object
//
// Debugging:
//   cat /tmp/activity-logger-error.log    # Errors
//   cat /tmp/activity-logger-success.log  # Successes

// ============================================================================
// END CLOSING
// ============================================================================
