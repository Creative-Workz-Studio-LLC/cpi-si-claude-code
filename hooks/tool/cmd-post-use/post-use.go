// METADATA
//
// PostToolUse Hook - Post-tool validation and contextual feedback orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Test everything; hold fast to what is good" - 1 Thessalonians 5:21
// Principle: Post-tool validation is testing the work - verifying quality, providing feedback, ensuring excellence
// Anchor: "Whatever you do, work heartily, as for the Lord" - Colossians 3:23
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - Hook orchestrator
// Role: Coordinates post-tool validation, formatting, and contextual feedback
// Paradigm: CPI-SI framework hook implementing quality assurance and user guidance
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-10 - Template application and named entry point pattern
//
// Version History:
//   2.0.0 (2025-11-10) - Full template application, named entry point, removed debug code
//   1.0.0 (2024-10-24) - Initial implementation
//
// Purpose & Function
//
// Purpose: Orchestrates post-tool validation and feedback by formatting code files,
// validating syntax, logging tool usage, and providing contextual guidance based on
// tool type and command patterns.
//
// Core Design: Thin orchestrator pattern - routes tool usage to appropriate handlers
// (file editing, bash commands, read/search operations) without implementing validation
// or feedback logic directly.
//
// Key Features:
//   - Automatic code formatting after Write/Edit operations
//   - Syntax validation with immediate feedback
//   - Tool usage activity logging with temporal context
//   - Contextual feedback (commits, builds, dependencies)
//   - Command failure detection and reporting
//   - Non-blocking design (failures don't interrupt tool completion)
//
// Philosophy: Like testing and holding fast to good (1 Thessalonians 5:21), post-tool
// validation ensures quality through immediate verification and constructive feedback.
//
// Blocking Status
//
// Non-blocking: All operations fail gracefully, tool completion never blocked.
// Validation failures report but don't prevent tool use. Logging failures are silent.
// Mitigation: Defensive checks throughout, tool execution is priority.
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code after each tool use
//	~/.claude/hooks/tool/cmd-post-use/post-use <tool_name> <tool_args>
//
// Integration Pattern:
//   1. Claude Code executes tool (Write, Edit, Bash, etc.)
//   2. PostToolUse hook runs with tool name and args
//   3. Routes to appropriate handler based on tool type
//   4. Logs tool usage with temporal context
//   5. Validates, formats, or provides feedback
//   6. Exits silently or reports results
//
// Hook Event: PostToolUse
// Trigger: After every tool execution
// Output: Validation results, formatting reports, contextual feedback
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os, path/filepath, strconv, strings, time
//   External: None
//   System Libraries: None
//   Hook Libraries: hooks/lib/activity, hooks/lib/feedback, system/lib/temporal, system/lib/validation
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook after tool use)
//
// Integration Points:
//   - Called by Claude Code hook system after tool execution
//   - Reads tool name and args from os.Args
//   - Reads environment variables (FILE_PATH, BASH_EXIT_CODE, etc.)
//   - Logs to activity stream
//   - Displays validation/feedback to user
//
// Health Scoring
//
// Post-tool validation orchestration operates on Base100 scale:
//
// Phase 1: Argument Validation: +5 points
//   - Check os.Args length (need at least 3)
//   - Silent exit if insufficient args
//
// Phase 2: Tool Routing: +10 points
//   - Route to appropriate handler based on tool type
//   - Detect Write/Edit/Bash/Read/Grep/Glob
//
// Phase 3: Handler Execution: +60 points (varies by tool)
//   - File editing (Write/Edit): +60 (logging + formatting + validation)
//   - Bash commands: +60 (logging + feedback + failure check)
//   - Read operations: +60 (logging only)
//   - Search operations (Grep/Glob): +60 (logging only)
//
// Phase 4: Temporal Logging: +15 points
//   - Gather temporal metadata
//   - Log context for pattern recognition
//
// Phase 5: Completion: +10 points
//   - Clean exit
//   - No resource cleanup needed
//
// Total Possible: 100 points for complete post-tool validation
//
// Current: No health tracking implemented (orchestration hook)
// Future: Track validation success rates, formatting coverage, feedback effectiveness
//
// END METADATA

package main

// ============================================================================
// SETUP
// ============================================================================
// Standard Reference: CWS-STD-001-DOC-4-block.md (SETUP Block)

// ────────────────────────────────────────────────────────────────
// Imports - Dependencies
// ────────────────────────────────────────────────────────────────
// Standard library for formatting, OS, paths, parsing, strings, time.
// Hook libraries for post-tool validation functionality.

import (
	"fmt"            // Formatted I/O for string construction
	"os"             // OS interface for environment and arguments
	"path/filepath"  // File path extension parsing
	"strconv"        // String to integer conversion for exit codes
	"strings"        // String manipulation for tool name detection
	"time"           // Duration types for command timing

	"hooks/lib/activity"   // Activity stream logging
	"hooks/lib/feedback"   // Contextual user feedback
	"system/lib/temporal"  // Temporal context for pattern recognition
	"system/lib/validation" // File formatting and syntax validation (v2.0.0 config-driven)
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - tool routing based on runtime args/env.

// No constants defined

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// No custom types needed - uses types from imported libraries.

// No types defined

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// This executable maintains no state - stateless orchestration only.

// No package-level state needed

// ============================================================================
// END SETUP
// ============================================================================

// ============================================================================
// BODY
// ============================================================================
// Standard Reference: CWS-STD-001-DOC-4-block.md (BODY Block)

// ════════════════════════════════════════════════════════════════════════════
// ORGANIZATIONAL CHART
// ════════════════════════════════════════════════════════════════════════════
//
// LADDER (Hierarchical Dependencies - What's Built On What):
//
//   postToolUse (main orchestrator)
//     ├─> getTemporalMetadata (temporal context helper)
//     ├─> handleFileEdit (file editing orchestration)
//     ├─> handleBashCommand (bash command orchestration)
//     ├─> hooks/lib/activity (activity stream logging)
//     ├─> hooks/lib/feedback (contextual user feedback)
//     ├─> hooks/lib/temporal (temporal context)
//     └─> system/lib/validation (file formatting and validation - v2.0.0 config-driven)
//
// BATON (Execution Flow - Data/Control Movement):
//
//   Entry → main()
//     ↓
//   postToolUse() orchestration
//     ↓
//   Phase 1: Argument validation (check os.Args length)
//     ↓
//   Phase 2: Parse tool name and args from command line
//     ↓
//   Phase 3: Gather temporal metadata for pattern recognition
//     ↓
//   Phase 4: Route to handler based on tool type:
//     ├─> Write/Edit → handleFileEdit()
//     │     ├─> Log tool use
//     │     ├─> Format file (validation.FormatFile)
//     │     └─> Validate file (validation.ValidateFile)
//     ├─> Bash → handleBashCommand()
//     │     ├─> Log command with exit code/duration
//     │     ├─> Provide contextual feedback (commits, builds, installs)
//     │     └─> Check for command failures
//     └─> Read/Grep/Glob → Log activity directly
//     ↓
//   Phase 5: Log temporal context for pattern recognition
//     ↓
//   Exit
//
// APUs (Atomic Processing Units - Core Work Functions):
//
//   postToolUse() - Main orchestrator for post-tool validation
//   getTemporalMetadata() - Temporal context formatter
//   handleFileEdit() - File editing validation orchestrator
//   handleBashCommand() - Bash command feedback orchestrator
//
// ════════════════════════════════════════════════════════════════════════════
// FUNCTION IMPLEMENTATIONS
// ════════════════════════════════════════════════════════════════════════════

// getTemporalMetadata returns temporal context as string for activity logging
//
// What It Does:
//   - Gathers temporal context from temporal library
//   - Formats as pipe-delimited string for logging
//   - Returns time of day, session phase, activity type
//
// Parameters:
//   - None (reads from temporal context)
//
// Returns:
//   - Formatted temporal metadata string (empty if unavailable)
//
// Health Impact:
//   - No health tracking (helper function)
//
// Example:
//   meta := getTemporalMetadata()
//   // Returns: "evening|long|focused-work"
func getTemporalMetadata() string {
	ctx, err := temporal.GetTemporalContext()
	if err != nil {
		return ""
	}

	metadata := fmt.Sprintf("%s|%s|%s",
		ctx.ExternalTime.TimeOfDay,
		ctx.InternalTime.SessionPhase,
		ctx.InternalSchedule.ActivityType)

	return metadata
}

// handleFileEdit processes Write/Edit tool usage
//
// What It Does:
//   - Logs tool usage to activity stream
//   - Formats code file using validation library
//   - Validates file after formatting
//   - Reports results to user
//
// Parameters:
//   - toolName: Name of tool used (Write or Edit)
//   - pattern: Tool arguments/pattern
//
// Returns:
//   - None (logs and displays results)
//
// Health Impact:
//   - No health tracking (orchestration helper)
//
// Example:
//   handleFileEdit("Write", "file.go")
//   // Logs, formats, validates, reports
func handleFileEdit(toolName, pattern string) {
	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		return
	}

	ext := filepath.Ext(filePath)

	// Log activity (non-blocking, privacy-preserving)
	activity.LogToolUse(toolName, filePath, true)

	// Format and validate code files
	result := validation.FormatFile(filePath, ext)
	result.Report()

	// Validate after formatting
	validationResult := validation.ValidateFile(filePath, ext)
	validationResult.Report()
}

// handleBashCommand processes Bash tool usage
//
// What It Does:
//   - Logs bash command with exit code and duration
//   - Provides contextual feedback based on command type
//   - Detects git commits, builds, dependency installs
//   - Checks for command failures
//
// Parameters:
//   - cmdStr: Bash command string
//
// Returns:
//   - None (logs and displays feedback)
//
// Health Impact:
//   - No health tracking (orchestration helper)
//
// Example:
//   handleBashCommand("git commit -m 'message'")
//   // Logs command, displays commit feedback
func handleBashCommand(cmdStr string) {
	// Log bash activity with exit code and duration
	exitCodeStr := os.Getenv("BASH_EXIT_CODE")
	exitCode := 0
	if exitCodeStr != "" {
		exitCode, _ = strconv.Atoi(exitCodeStr)
	}

	durationStr := os.Getenv("BASH_DURATION_MS")
	duration := time.Duration(0)
	if durationStr != "" {
		if ms, err := strconv.ParseInt(durationStr, 10, 64); err == nil {
			duration = time.Duration(ms) * time.Millisecond
		}
	}

	activity.LogCommand(cmdStr, exitCode, duration)

	// Provide contextual feedback based on command type
	switch {
	case strings.Contains(cmdStr, "git commit"):
		feedback.OnCommitCreated()
	case feedback.IsBuildCommand(cmdStr):
		feedback.OnBuildComplete()
	case strings.Contains(cmdStr, "npm install") ||
		strings.Contains(cmdStr, "cargo install") ||
		strings.Contains(cmdStr, "go install"):
		feedback.OnDependenciesInstalled()
	}

	// Check for command failures
	feedback.OnCommandFailure(exitCodeStr)
}

// postToolUse orchestrates post-tool validation and feedback
//
// What It Does:
//   - Routes tool usage to appropriate handler
//   - Logs temporal context for pattern recognition
//   - Coordinates validation, formatting, and feedback
//   - Handles Write, Edit, Bash, Read, Grep, Glob tools
//
// Parameters:
//   - None (reads from os.Args and environment)
//
// Returns:
//   - None (all operations non-blocking)
//
// Health Impact:
//   - No health tracking (orchestration function)
//
// Example:
//   postToolUse()
//   // Routes to handleFileEdit, handleBashCommand, or logging
func postToolUse() {
	if len(os.Args) < 3 {
		return
	}

	toolName := os.Args[1]
	toolArgs := os.Args[2]

	// Get temporal context for tool use (helps understand patterns)
	temporalMeta := getTemporalMetadata()

	// Route to appropriate handler based on tool type
	switch {
	case strings.HasPrefix(toolName, "Write") || strings.HasPrefix(toolName, "Edit"):
		handleFileEdit(toolName, toolArgs)
	case strings.HasPrefix(toolName, "Bash"):
		handleBashCommand(toolArgs)
	case strings.HasPrefix(toolName, "Read"):
		// Log Read operations
		filePath := os.Getenv("FILE_PATH")
		activity.LogToolUse("Read", filePath, true)
	case strings.HasPrefix(toolName, "Grep"):
		// Log Grep operations (search activity)
		pattern := os.Getenv("GREP_PATTERN")
		if pattern == "" {
			pattern = "search"
		}
		activity.LogActivity("Grep", pattern, "success", 0)
	case strings.HasPrefix(toolName, "Glob"):
		// Log Glob operations (file discovery)
		pattern := os.Getenv("GLOB_PATTERN")
		if pattern == "" {
			pattern = "search"
		}
		activity.LogActivity("Glob", pattern, "success", 0)
	}

	// Log temporal context for this tool use (enables pattern recognition)
	if temporalMeta != "" {
		activity.LogActivity("ToolContext", temporalMeta, "info", 0)
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
// Code Validation: Build and Hook Testing
// ────────────────────────────────────────────────────────────────
//
// Testing Requirements:
//   - Compile without errors: go build -o post-use post-use.go
//   - Run standalone: ./post-use Write file.go (should format and validate)
//   - Test with various tool types (Write, Edit, Bash, Read, Grep, Glob)
//   - Verify activity log entries created
//   - Test validation feedback displays
//   - Ensure non-blocking (no panics or exits on errors)
//
// Build Verification:
//   cd ~/.claude/hooks/tool/cmd-post-use
//   go build -o post-use post-use.go
//   FILE_PATH=/tmp/test.go ./post-use Write test.go  # Manual execution test
//
// Integration Testing:
//   Test with Claude Code tool execution
//   Verify formatting occurs after Write/Edit
//   Check validation feedback appears
//   Validate contextual feedback for commits/builds
//   Ensure temporal logging works
//
// Example validation:
//
//     # Test file editing
//     cd ~/.claude/hooks/tool/cmd-post-use
//     FILE_PATH=/path/to/file.go ./post-use Write file.go
//
//     # Check activity log
//     tail -5 ~/.claude/logs/activity/activity.log | grep ToolUse
//
//     # Test bash feedback
//     BASH_EXIT_CODE=0 ./post-use Bash "git commit -m 'test'"

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls postToolUse() (named entry point)
//   3. postToolUse() orchestrates validation and feedback
//   4. Program exits after processing
//
// Named Entry Point Benefits:
//   - Prevents main function collisions across executables
//   - Function name matches purpose (post-tool use)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Separates Go runtime entry from application logic
//
// Pattern:
//   func main() {
//       postToolUse()  // Minimal - just call named entry point
//   }

func main() {
	postToolUse() // Named entry point pattern
}

// ────────────────────────────────────────────────────────────────
// Code Cleanup: Resource Management
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - stdout: Used for validation/feedback output
//   - Libraries: No persistent state to clean
//   - File handles: Managed by validation library
//   - No manual resources
//
// Graceful Shutdown:
//   - Program exits immediately after feedback
//   - No cleanup needed (stateless execution)
//   - Error path: Non-blocking (failures don't prevent completion)
//   - Success path: Normal completion, exit code 0
//
// Error State Cleanup:
//   - Errors don't prevent tool completion (non-blocking design)
//   - Validation failures report but don't block
//   - Libraries handle their own error states
//
// Memory Management:
//   - Go's garbage collector handles all memory
//   - Short-lived execution (~50ms typical)
//   - No large buffers or persistent state

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════

// ────────────────────────────────────────────────────────────────
// Executable Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: See METADATA for comprehensive capabilities
//
// Quick summary:
//   - Orchestrates post-tool validation and feedback
//   - Formats code files automatically after Write/Edit
//   - Validates syntax with immediate feedback
//   - Logs tool usage with temporal context
//   - Provides contextual feedback (commits, builds, dependencies)
//   - Detects command failures and reports
//   - Non-blocking design ensures tool completion never blocked
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Hook Event: PostToolUse (triggered by Claude Code after every tool execution)
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (EXECUTABLE hook orchestrator)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new tool type handlers (extend switch in postToolUse)
//   ✅ Add new contextual feedback patterns (extend handleBashCommand)
//   ✅ Enhance validation rules (modify validation library)
//   ✅ Add new feedback types (extend feedback library)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ os.Args structure - breaks Claude Code integration
//   ⚠️ Environment variable names - breaks tool routing
//   ⚠️ Handler signatures - breaks orchestration flow
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking principle (tool completion never blocked)
//   ❌ Named entry point pattern (main calls postToolUse)
//   ❌ Stateless design (no persistent state in hook executable)
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing requirements,
//   build verification, and integration testing procedures. Always verify
//   hook still works with Claude Code after changes.

// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart" section above for complete ladder structure
// (dependencies) and baton flow (execution paths).
//
// Quick architectural summary:
// - Thin orchestrator calling library functions
// - 4 functions: postToolUse() (main) + 3 helpers (getTemporalMetadata, handleFileEdit, handleBashCommand)
// - Ladder: Hook executable → activity/feedback/temporal/validation libraries
// - Baton: Argument validation → tool routing → handler execution → temporal logging

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// Adding new tool type handler:
//   1. Add case to switch statement in postToolUse()
//   2. Create handler function or inline logging
//   3. Update organizational chart and health scoring
//
// Adding new bash feedback pattern:
//   1. Add case to switch in handleBashCommand()
//   2. Call feedback function from feedback library
//   3. Test with actual bash commands
//
// Enhancing file validation:
//   1. Modify validation library functions
//   2. Test with various file types
//   3. Verify feedback display

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Typical execution: ~50ms total
//   - Argument parsing: <1ms
//   - Tool routing: <1ms
//   - File formatting: ~20ms (depends on file size)
//   - Syntax validation: ~10ms
//   - Feedback display: ~5ms
//   - Temporal logging: ~5ms
//
// Performance not critical - post-tool hooks run after tool completion.
// User-facing validation prioritizes accuracy over speed.
//
// Non-blocking design ensures tool completion proceeds regardless of hook timing.

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Hook doesn't execute:
//   - Check hook is executable: chmod +x ~/.claude/hooks/tool/cmd-post-use/post-use
//   - Verify Claude Code hook configuration
//   - Check for stderr output indicating errors
//
// Formatting not happening:
//   - Verify FILE_PATH environment variable set
//   - Check validation.FormatFile() working
//   - Test standalone: FILE_PATH=/path/to/file.go ./post-use Write file.go
//
// Validation not showing:
//   - Verify validation.ValidateFile() implementation
//   - Check file extension supported
//   - Test with known invalid syntax
//
// Feedback not appearing:
//   - Check feedback library functions working
//   - Verify command pattern detection (git commit, build, install)
//   - Test standalone with specific commands
//
// Activity logs missing:
//   - Check ~/.claude/logs/activity/ directory exists
//   - Verify activity.LogToolUse() and activity.LogActivity() working
//   - Check log file permissions

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
//   - hooks/lib/activity: Activity stream logging
//   - hooks/lib/feedback: Contextual user feedback
//   - hooks/lib/temporal: Temporal context
//   - system/lib/validation: File formatting and validation (v2.0.0 config-driven)
//
// Related hooks:
//   - tool/cmd-pre-use: Complements pre-tool preparation
//   - session/cmd-start: Session initialization
//   - session/cmd-end: Session completion
//
// Related utilities:
//   - None (standalone post-tool processing)

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ File formatting - COMPLETED
//   ✓ Syntax validation - COMPLETED
//   ✓ Tool usage logging - COMPLETED
//   ✓ Contextual feedback - COMPLETED
//   ✓ Temporal context - COMPLETED
//   ⏳ AI-powered code review suggestions
//   ⏳ Performance profiling integration
//   ⏳ Security scanning for common vulnerabilities
//
// Research Areas:
//   - Machine learning for code quality prediction
//   - Automated refactoring suggestions
//   - Pattern detection across tool usage
//
// Integration Targets:
//   - Code quality metrics dashboard
//   - Continuous improvement tracking
//   - Team-wide feedback aggregation
//
// Known Limitations to Address:
//   - No AI-powered review (manual validation only)
//   - No performance profiling (timing only)
//   - No security scanning (syntax validation only)

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This hook embodies the testing principle from 1 Thessalonians 5:21 - verifying
// quality through immediate validation and constructive feedback. Excellence in
// post-tool validation ensures work meets standards before proceeding.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain non-blocking design principle (tool completion never blocked)
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "Test everything; hold fast to what is good" - 1 Thessalonians 5:21
//
// Every tool use is opportunity for excellence. Let validation serve quality without obstruction.

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Execution (Manual):
//
//     cd ~/.claude/hooks/tool/cmd-post-use
//     FILE_PATH=/path/to/file.go ./post-use Write file.go
//
// With Bash Command:
//
//     BASH_EXIT_CODE=0 BASH_DURATION_MS=1234 ./post-use Bash "git commit -m 'message'"
//
// Checking Activity Log:
//
//     tail -10 ~/.claude/logs/activity/activity.log | grep ToolUse
//
// Testing Validation:
//
//     # Create test file with syntax error
//     echo "package main\nfunc main( {" > /tmp/test.go
//     FILE_PATH=/tmp/test.go ./post-use Write test.go
//     # Should show validation error

// ============================================================================
// END CLOSING
// ============================================================================
