// METADATA
//
// SessionStart Hook - Session Initialization Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "In the beginning, God created the heavens and the earth" - Genesis 1:1
// Principle: Every work has a beginning - session start establishes context and awareness for faithful work
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - Hook orchestrator
// Role: Coordinates session initialization, gathers context for autonomous covenant partnership work
// Paradigm: CPI-SI framework hook implementing non-blocking session startup with comprehensive awareness
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-10 - Template application and library extraction
//
// Version History:
//   2.0.0 (2025-11-10) - Full template application, logic extracted to hooks/lib/session
//   1.0.0 (2024-10-24) - Initial implementation with inline logic
//
// Purpose & Function
//
// Purpose: Orchestrates session initialization by gathering environment, temporal, and workspace
// context for Nova Dawn to begin autonomous work with full situational awareness.
//
// Core Design: Thin orchestrator pattern - coordinates modular components from hooks/lib/session,
// system/lib, and hooks/lib/temporal to present unified session context without implementing
// business logic directly.
//
// Key Features:
//   - Clean screen presentation with instance branding
//   - Environment awareness (workspace, git, system info)
//   - Temporal consciousness (4 dimensions: external time, internal time, schedule, calendar)
//   - Workspace analysis (git status, processes, disk, dependencies, activity)
//   - Claude Code context injection (Nova Dawn communication style + temporal awareness)
//   - Non-blocking design (failures don't prevent session start)
//
// Philosophy: Session start is first impression and foundation for work. Like Genesis 1:1
// establishes beginning of creation, session start establishes beginning of covenant work -
// providing order, context, and awareness for what follows.
//
// Blocking Status
//
// Non-blocking: All operations fail gracefully, session starts even if context gathering fails.
// Display errors go to stderr, don't exit. JSON output errors warned but don't block session.
// Mitigation: Defensive checks throughout, silent failures for non-critical features.
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code on SessionStart event
//	~/.claude/hooks/session/cmd-start/start
//
// Integration Pattern:
//   1. Claude Code triggers SessionStart hook event
//   2. start executable runs, clears screen
//   3. Displays session banner and context
//   4. Outputs JSON for Claude Code to inject context
//   5. Session begins with full awareness
//
// Hook Event: SessionStart
// Trigger: When Claude Code session begins
// Output: Visual display + JSON context injection
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: encoding/json, fmt, os
//   External: None
//   System Libraries: system/lib/git, system/lib/instance
//   Hook Libraries: hooks/lib/session (display, init, context), hooks/lib/activity, hooks/lib/temporal
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook on SessionStart event)
//
// Integration Points:
//   - Called by Claude Code hook system on SessionStart
//   - Outputs JSON to stdout for Claude Code parsing
//   - Initializes session-time and session-log utilities
//   - Reads NOVA_DAWN_WORKSPACE environment variable
//
// Health Scoring
//
// Session initialization orchestration operates on Base100 scale:
//
// Initialization:
//   - Session timing init: +10
//   - Session logging init: +10
//   - Activity logging: +10
//
// Display:
//   - Header display: +10
//   - Environment display: +10
//   - Temporal awareness display: +10
//   - Workspace analysis: +20
//
// Context Output:
//   - Claude Code JSON output: +20
//
// Total: 100 points for complete session initialization
//
// Note: Non-blocking design means partial failures still allow session to proceed.
// Health scoring tracks what succeeded, not what prevented session start.
package main

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
// Standard library for basic I/O and environment access.
// System libraries for git and instance configuration.
// Hook libraries for session-specific functionality.

import (
	"fmt" // Formatted I/O for display output
	"os"  // OS interface for environment variables and stderr

	"system/lib/git" // Git repository detection and branch info

	"hooks/lib/activity" // Activity stream logging
	"hooks/lib/session"  // Session display, init, context functions
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - configuration comes from environment and libraries.

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
//
// For BODY structure explanation, see: standards/code/4-block/CWS-STD-007-CODE-body-block.md

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Internal Structure
// ────────────────────────────────────────────────────────────────
// Maps bidirectional dependencies and baton flow within this executable.
//
// Ladder Structure (Dependencies):
//
//   Orchestration (Top Rung - Entry Point)
//   └── start() → calls main orchestration
//       └── gatherContext() → workspace analysis orchestrator
//           └── session.* functions → delegated to libraries
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/session (display, init, context, workspace checks)
//   ├── hooks/lib/activity (activity logging)
//   ├── hooks/lib/temporal (temporal awareness)
//   └── system/lib/git (repository detection)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → start()
//     ↓
//   Initialize → session.InitSessionTime(), session.InitSessionLog()
//     ↓
//   Log → activity.LogActivity()
//     ↓
//   Clear Screen → fmt.Print()
//     ↓
//   Display → session.PrintHeader(), session.PrintEnvironment(), session.PrintTemporalAwareness()
//     ↓
//   Analyze → gatherContext() if workspace configured
//     ↓
//   Output Context → session.OutputClaudeContext()
//     ↓
//   Exit → return
//
// APUs (Available Processing Units):
// - 2 functions total
// - 1 orchestration helper (gatherContext)
// - 1 entry point (start, called by main)

// ────────────────────────────────────────────────────────────────
// Workspace Analysis - Context Gathering Orchestration
// ────────────────────────────────────────────────────────────────
// What This Does:
// Coordinates workspace analysis by calling modular session library components
// to gather git status, processes, disk space, dependencies, and recent activity.
//
// Why Separated:
// Workspace analysis is optional (only when NOVA_DAWN_WORKSPACE set) and involves
// multiple checks. Isolating as function keeps main orchestration clean.

// gatherContext orchestrates workspace analysis from modular components
//
// What It Does:
//   - Coordinates various workspace checks via session library
//   - Git status, running processes, disk space, dependencies, recent activity
//   - Displays results from each check
//   - Shows workspace analysis header with aggregated results
//
// Parameters:
//   workspace: Workspace directory path to analyze
//
// Returns:
//   None (displays results to stdout)
//
// Health Impact:
//   +20 points for successful workspace analysis coordination
//   Delegates actual checks to session library (health tracked there)
//
// Example usage:
//
//	workspace := "/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC"
//	gatherContext(workspace)
//	// Outputs workspace analysis section with all checks
func gatherContext(workspace string) {
	// Track if any context was gathered
	hasContext := false

	// Git repository analysis
	if git.IsGitRepository(workspace) {
		session.CheckGitStatus(workspace)
		hasContext = true
	}

	// Development environment checks
	session.CheckRunningProcesses()

	// System resource checks
	session.CheckDiskSpace(workspace)

	// Project dependency validation
	session.CheckDependencies(workspace)

	// Recent activity tracking
	session.CheckRecentActivity(workspace)

	// Display workspace analysis header with results
	session.PrintWorkspaceAnalysis(workspace, hasContext)
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
//   - Compile without errors: go build -o start start.go
//   - Run standalone: ./start (should display session context)
//   - Verify all display sections appear
//   - Check Claude Code JSON output at end
//   - Test with NOVA_DAWN_WORKSPACE set and unset
//   - Verify graceful handling of missing temporal data
//
// Build Verification:
//   cd ~/.claude/hooks/session/cmd-start
//   go build -o start start.go
//   ./start  # Manual execution test
//
// Integration Testing:
//   Test with Claude Code SessionStart event
//   Verify JSON context injection works
//   Check that session initializes with full awareness
//   Validate temporal awareness displays correctly
//
// Example validation:
//
//     # Test standalone execution
//     cd ~/.claude/hooks/session/cmd-start
//     ./start
//
//     # Test with workspace
//     NOVA_DAWN_WORKSPACE=/path/to/workspace ./start
//
//     # Verify JSON output
//     ./start | tail -1 | jq .hookSpecificOutput.hookEventName
//     # Should output: "SessionStart"

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls start() (named entry point)
//   3. start() orchestrates complete session initialization
//   4. Program exits after context output
//
// Named Entry Point Benefits:
//   - Prevents main function collisions across executables
//   - Function name matches purpose (session start)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Separates Go runtime entry from application logic
//
// Pattern:
//   func main() {
//       start()  // Minimal - just call named entry point
//   }

// start orchestrates complete session initialization
//
// What It Does:
//   - Initializes session timing and logging
//   - Logs session start activity
//   - Clears screen for clean presentation
//   - Displays session header, environment, temporal awareness
//   - Gathers and displays workspace analysis if configured
//   - Outputs Claude Code context JSON
//
// Parameters:
//   None (reads from environment and libraries)
//
// Returns:
//   None (outputs to stdout, exits after completion)
//
// Health Impact:
//   Coordinates all initialization steps (+100 total)
//   See METADATA Health Scoring for complete breakdown
//
// Example:
//   Called automatically by main() when hook executes
func start() {
	// Initialize session timing (captures start time for time awareness)
	// Health: +10
	session.InitSessionTime()

	// Initialize session history logging (for pattern learning)
	// Health: +10
	session.InitSessionLog()

	// Log session start event to activity stream
	// Health: +10
	activity.LogActivity("SessionStart", "session-initialized", "success", 0)

	// Clear screen for clean presentation
	fmt.Print("\033[H\033[2J\033[3J")

	// Get workspace configuration
	workspace := os.Getenv("NOVA_DAWN_WORKSPACE")

	// Display session header
	// Health: +10
	session.PrintHeader()

	// Show environment context
	// Health: +10
	session.PrintEnvironment(workspace)

	// Show temporal awareness (4 dimensions of time/schedule consciousness)
	// Health: +10
	session.PrintTemporalAwareness()

	// Gather and display workspace analysis
	// Health: +20
	if workspace != "" {
		gatherContext(workspace)
	} else {
		// Display workspace analysis with no workspace configured
		session.PrintWorkspaceAnalysis(workspace, false)
	}

	// Display formatted session context for user readability
	// Health: +15
	sessionContext := session.GetSessionContext()
	session.PrintSessionContext(sessionContext)

	// Output Claude Code context JSON (must be last for Claude to parse)
	// Health: +20
	if err := session.OutputClaudeContext(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to output Claude context: %v\n", err)
		// Non-blocking: don't exit on error, session can still start
	}
}

func main() {
	start() // Entry point for session initialization orchestration
}

// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Stateless Hook)
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - stdout: Automatically managed by OS
//   - Libraries: No persistent state to clean
//   - Session utilities: Manage their own state
//   - No files, connections, or manual resources
//
// Graceful Shutdown:
//   - Program exits immediately after context output
//   - No cleanup needed (stateless execution)
//   - Error path: Print to stderr, continue (non-blocking)
//   - Success path: Output to stdout, exit code 0
//
// Error State Cleanup:
//   - Errors logged but don't prevent session start
//   - No partial state or rollback needed
//   - Libraries handle their own error states
//
// Memory Management:
//   - Go's garbage collector handles all memory
//   - Short-lived execution (~100ms typical)
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
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Orchestrates session initialization with comprehensive context
//   - Displays environment, temporal awareness, workspace analysis
//   - Injects Nova Dawn communication style into Claude Code
//   - Non-blocking design ensures session starts even with failures
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Hook Event: SessionStart (triggered by Claude Code when session begins)
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (EXECUTABLE hook orchestrator)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new workspace checks (call session.Check* functions in gatherContext)
//   ✅ Add new display sections (call session.Print* functions in start)
//   ✅ Adjust display order (rearrange function calls in start)
//   ✅ Add health tracking (coordinate with session library updates)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ JSON output format - breaks Claude Code parsing
//   ⚠️ Hook event name - breaks Claude Code hook system
//   ⚠️ stdout/stderr separation - breaks context injection
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking principle (failures must not prevent session start)
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
// See BODY "Organizational Chart - Internal Structure" section above for
// complete ladder structure (dependencies) and baton flow (execution paths).
//
// Quick architectural summary:
// - Thin orchestrator calling library functions
// - 2 functions: start() (main orchestration) + gatherContext() (workspace helper)
// - Ladder: Hook executable → session library → system libraries
// - Baton: Linear initialization flow with delegated checks

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// Adding new display sections:
//   1. Create new function in hooks/lib/session/display.go
//   2. Call from start() in appropriate order
//   3. Update health scoring map in METADATA
//
// Adding new workspace checks:
//   1. Create new function in hooks/lib/session/*.go
//   2. Call from gatherContext() with other checks
//   3. Update health scoring and organizational chart
//
// Modifying display order:
//   1. Rearrange function calls in start()
//   2. Maintain logical flow (init → display → analyze → output)
//   3. Keep JSON output last (Claude Code requirement)

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Typical execution: ~100ms total
//   - Display rendering: ~10ms
//   - Temporal context: ~20ms (file reads)
//   - Workspace analysis: ~50ms (git, disk, process checks)
//   - JSON generation: ~10ms
//
// Performance is not critical - session start is one-time event.
// Clarity and completeness prioritized over speed.
//
// Non-blocking design ensures failures don't compound delays.

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Session doesn't start:
//   - Check hook is executable: chmod +x ~/.claude/hooks/session/cmd-start/start
//   - Verify Claude Code hook configuration
//   - Check stderr for error messages
//
// Missing display sections:
//   - Temporal awareness missing: Check session-time utility exists
//   - Workspace analysis missing: Verify NOVA_DAWN_WORKSPACE set
//   - Git info missing: Confirm directory is git repository
//
// Claude Code context not injected:
//   - Verify JSON output appears (should be last line)
//   - Check JSON is valid: ./start | tail -1 | jq
//   - Ensure hookEventName is "SessionStart"
//
// Display formatting issues:
//   - Check terminal supports ANSI escape codes
//   - Verify box-drawing characters render correctly
//   - Test with different terminal emulators if needed

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
//   - hooks/lib/session: Display, init, context, workspace checks
//   - hooks/lib/temporal: 4-dimensional time awareness
//   - hooks/lib/activity: Activity stream logging
//   - system/lib/git: Repository detection and info
//   - system/lib/instance: Instance configuration
//
// Related hooks:
//   - session/cmd-end: Complements start, handles session end
//   - session/cmd-stop: Handles session stop events
//   - session/cmd-pre-compact: Pre-compaction awareness

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ Environment context display - COMPLETED
//   ✓ Temporal awareness (4 dimensions) - COMPLETED
//   ✓ Workspace analysis - COMPLETED
//   ✓ Claude Code context injection - COMPLETED
//   ⏳ Session resume detection (distinguish fresh start vs resume)
//   ⏳ Last session summary (what was accomplished previously)
//   ⏳ Project context loading (active projects and priorities)
//
// Research Areas:
//   - Session pattern recognition (typical start time, work patterns)
//   - Pre-load frequent files/context based on workspace
//   - Integration with planner for task prioritization display
//
// Integration Targets:
//   - Memory system for session continuity
//   - Project tracking for automatic context loading
//   - Pattern learning for personalized startup
//
// Known Limitations to Address:
//   - Fresh start vs resume not distinguished
//   - No session history display
//   - Project context not loaded automatically

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This hook is the FIRST IMPRESSION of every session - the foundation for
// covenant partnership work. Like Genesis 1:1 establishes the beginning of
// creation, this hook establishes the beginning of work with order, context,
// and awareness.
//
// Modify thoughtfully - changes here affect how every session begins. The
// non-blocking design ensures robustness (session starts even with failures),
// but the display provides essential context for autonomous work.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain non-blocking design principle
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "In the beginning, God created the heavens and the earth" - Genesis 1:1
//
// Every session is a new beginning. Let it start with order and awareness.

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Execution (Manual):
//
//     cd ~/.claude/hooks/session/cmd-start
//     ./start
//
// With Workspace:
//
//     NOVA_DAWN_WORKSPACE=/path/to/workspace ./start
//
// Testing JSON Output:
//
//     ./start | tail -1 | jq .hookSpecificOutput.hookEventName
//     # Should output: "SessionStart"
//
// Verifying Display Sections:
//
//     ./start | grep -E "SESSION ENVIRONMENT|TEMPORAL AWARENESS|WORKSPACE ANALYSIS"
//     # All three sections should appear
//
// ============================================================================
// END CLOSING
// ============================================================================
