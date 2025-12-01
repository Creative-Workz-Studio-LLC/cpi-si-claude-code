// METADATA
//
// SessionStop Hook - Session Stop Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Whatever you do, work heartily, as for the Lord and not for men" - Colossians 3:23
// Principle: Every stop is a reflection point - stopping well honors the work done, like Sabbath rest after six days
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - Hook orchestrator
// Role: Coordinates session stop summary with stopping point quality checks and temporal awareness
// Paradigm: CPI-SI framework hook implementing non-blocking session stop with graceful transition
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
// Purpose: Orchestrates graceful session stop by displaying temporal context, checking stopping
// point quality, and reminding about workspace state to ensure faithful completion of work.
//
// Core Design: Thin orchestrator pattern - coordinates modular components from hooks/lib/session
// and hooks/lib/temporal to present stopping summary without implementing business logic directly.
//
// Key Features:
//   - Stop banner with biblical foundation (Colossians 3:23)
//   - Temporal awareness at stop time (when stopped, how long worked)
//   - Stopping point quality checks (uncommitted work, running processes)
//   - Activity context for session continuity
//   - Non-blocking design (failures don't prevent session stop)
//
// Philosophy: Session stop is not just "user left" - it's intentional transition point. Like
// Sabbath rest after six days of work, stopping well honors the work done. Hook provides moment
// to check quality of stopping point and remind about unfinished work.
//
// Blocking Status
//
// Non-blocking: All operations fail gracefully, session stops even if display or checks fail.
// Display errors go to stderr, don't exit. Check failures warned but don't block stop.
// Mitigation: Defensive checks throughout, silent failures for non-critical features.
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code on SessionStop event
//	~/.claude/hooks/session/cmd-stop/stop
//
// Integration Pattern:
//   1. Claude Code triggers SessionStop hook event
//   2. stop executable runs with REASON environment variable
//   3. Displays stop banner and temporal context
//   4. Checks stopping point quality if workspace configured
//   5. Session stops with graceful summary
//
// Hook Event: SessionStop
// Trigger: When user stops Claude Code session
// Output: Visual display of stop summary and checks
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os
//   External: None
//   System Libraries: system/lib/git (via hooks/lib/session)
//   Hook Libraries: hooks/lib/session (display, checks), hooks/lib/activity, hooks/lib/temporal
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook on SessionStop event)
//
// Integration Points:
//   - Called by Claude Code hook system on SessionStop
//   - Reads REASON environment variable (stop reason)
//   - Reads NOVA_DAWN_WORKSPACE environment variable
//   - Logs to activity stream for pattern learning
//
// Health Scoring
//
// Session stop orchestration operates on Base100 scale:
//
// Initialization:
//   - Get stop reason: +10
//   - Activity logging: +10
//
// Display:
//   - Stop header display: +15
//   - Stop info display: +10
//   - Temporal context display: +15
//
// Analysis:
//   - Uncommitted work reminder: +10
//   - Running process check: +10
//   - Recent activity review: +10
//
// Output:
//   - Closing divider: +10
//
// Total: 100 points for complete session stop
//
// Note: Non-blocking design means partial failures still allow session to stop.
// Health scoring tracks what succeeded, not what prevented stop.
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
// Hook libraries for session-specific functionality.

import (
	"fmt" // Formatted I/O for display output
	"os"  // OS interface for environment variables

	"hooks/lib/activity" // Activity stream logging
	"hooks/lib/session"  // Session display and check functions
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
//   └── stop() → calls main orchestration
//       └── checkStoppingPoint() → stopping point validation orchestrator
//           └── session.* functions → delegated to libraries
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/session (display, checks, reminders)
//   ├── hooks/lib/activity (activity logging)
//   ├── hooks/lib/temporal (temporal awareness)
//   └── system/lib/git (via session library)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → stop()
//     ├→ Phase 1: Initialization
//     │   ├→ Get stop reason from environment
//     │   └→ Log activity event
//     ├→ Phase 2: Display
//     │   ├→ session.PrintStopHeader()
//     │   ├→ session.PrintStopInfo()
//     │   └→ session.PrintStoppingContext()
//     ├→ Phase 3: Analysis (if workspace configured)
//     │   └→ checkStoppingPoint()
//     │       ├→ session.RemindUncommittedWork()
//     │       ├→ session.CheckRunningProcessesAsReminder()
//     │       └→ session.CheckRecentActivity()
//     └→ Phase 4: Output
//         └→ Display closing divider
//
// APUs (Atomic Processing Units):
//   2 functions:
//     - stop() [70 points]: Main orchestration (init, display, output)
//     - checkStoppingPoint() [30 points]: Stopping point validation orchestration
//
// ────────────────────────────────────────────────────────────────
// Function Implementations
// ────────────────────────────────────────────────────────────────

// checkStoppingPoint orchestrates stopping point validation
//
// What It Does:
//   - Orchestrates three stopping point checks
//   - Reminds about uncommitted work
//   - Checks for running processes that need attention
//   - Reviews recent activity for context
//
// Why It Exists:
//   - Stopping well matters - leaving work in good state honors the work
//   - Separates "what to check" (orchestration) from "how to check" (libraries)
//   - Can be reused or extended without changing display logic
//
// When Called:
//   - Called by stop() if workspace is configured
//   - Skipped if NOVA_DAWN_WORKSPACE not set (nothing to check)
//
// Parameters:
//   - workspace: Workspace directory path for checks
//
// Returns:
//   - None (orchestrates display functions)
//
// Health Contribution:
//   - 30 points total (10 per check function)
//   - All checks non-blocking (failures don't stop session)
//
// Related Components:
//   - hooks/lib/session/reminders.go: RemindUncommittedWork
//   - hooks/lib/session/processes.go: CheckRunningProcessesAsReminder
//   - hooks/lib/session/activity.go: CheckRecentActivity
//
// Example:
//   checkStoppingPoint("/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC")
//   // Displays three reminder sections about workspace state
func checkStoppingPoint(workspace string) {
	// Orchestrate three checks - order matters (uncommitted work most urgent)
	session.RemindUncommittedWork(workspace)
	session.CheckRunningProcessesAsReminder()
	session.CheckRecentActivity(workspace)
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
// Validation - Pre-Execution Checks
// ────────────────────────────────────────────────────────────────
// Pre-execution validation not needed (hook guaranteed to be called correctly)
// Environment variable defaults ensure safe execution
// All library functions handle their own validation

// ────────────────────────────────────────────────────────────────
// Execution - Main Entry Points
// ────────────────────────────────────────────────────────────────

// stop is the named entry point for session stop orchestration
//
// What It Does:
//   - Orchestrates complete session stop sequence
//   - Gets stop reason from environment (with default)
//   - Logs activity for pattern learning
//   - Displays stop banner with biblical foundation
//   - Shows temporal awareness at stop time
//   - Checks stopping point quality if workspace configured
//   - Provides graceful stop summary
//
// Why Named Entry Point Pattern:
//   - main() is generic, stop() is semantically meaningful
//   - Allows testing without executable mechanics
//   - Clear architectural intent (this is stop orchestration)
//   - Matches pattern across all CPI-SI hooks
//
// Orchestration Philosophy:
//   - Thin orchestrator (coordinates, doesn't implement)
//   - All logic in libraries (display, checks, temporal awareness)
//   - Non-blocking (session stops even with warnings)
//   - Grace in failures (missing workspace? skip checks, continue)
//
// Parameters:
//   - None (reads from environment variables)
//
// Returns:
//   - None (orchestrates display and checks to stdout)
//
// Health Contribution:
//   - 70 points total (initialization + display + output)
//   - Perfect execution: all display and checks succeed
//   - Partial execution: some displays/checks fail (still completes)
//
// Environment Variables:
//   - REASON: Stop reason (defaults to "User stepping away")
//   - NOVA_DAWN_WORKSPACE: Workspace path (empty = skip checks)
//
// Example:
//   stop()
//   // Executes complete stop sequence with all displays and checks
func stop() {
	// Phase 1: Initialization (20 points)
	reason := os.Getenv("REASON")
	if reason == "" {
		reason = "User stepping away"
	}

	// Log session stop event to activity stream
	activity.LogActivity("SessionStop", reason, "success", 0)

	// Phase 2: Display (40 points)
	session.PrintStopHeader()      // Stop banner with Colossians 3:23
	session.PrintStopInfo()        // Timestamp and stopping point check header
	session.PrintStoppingContext() // Temporal awareness at stop

	// Phase 3: Analysis (30 points)
	workspace := os.Getenv("NOVA_DAWN_WORKSPACE")
	if workspace != "" {
		checkStoppingPoint(workspace)
	} else {
		fmt.Println() // Spacing if no workspace to check
	}

	// Phase 4: Output (10 points)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

func main() {
	stop() // Named entry point pattern
}

// ────────────────────────────────────────────────────────────────
// Cleanup - Resource Management
// ────────────────────────────────────────────────────────────────
// No cleanup needed (hook is stateless, no resources to release)

// ────────────────────────────────────────────────────────────────
// FINAL DOCUMENTATION
// ────────────────────────────────────────────────────────────────
//
// This section provides comprehensive documentation for maintainers, following
// the 11-subsection structure for executable documentation.

// EXECUTABLE OVERVIEW & INTEGRATION
// ==================================
//
// This executable is a Claude Code SessionStop hook. It runs automatically
// when the user stops a session, providing graceful transition with:
//   - Biblical reminder about faithful work (Colossians 3:23)
//   - Temporal awareness (when stopped, how long worked)
//   - Stopping point quality checks (uncommitted work, running processes)
//   - Activity context for session continuity
//
// Integration points:
//   - Called by Claude Code on SessionStop event
//   - Reads environment variables for context
//   - Logs to activity stream for pattern learning
//   - Outputs user-facing summary to stdout
//   - Non-blocking (always completes even with errors)
//
// Relationship to other hooks:
//   - Complements SessionStart (stop.go vs start.go)
//   - Shares display libraries with all session hooks
//   - Uses same temporal awareness as start/end hooks
//   - Part of complete session lifecycle (start → work → stop → end)

// MODIFICATION POLICY
// ===================
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new stopping point checks:
//      - Add check functions to hooks/lib/session/
//      - Call from checkStoppingPoint() orchestration
//      - Update health scoring map in METADATA
//
//   ✅ Change display order:
//      - Reorder calls in stop() function
//      - Maintains orchestration pattern
//
//   ✅ Add new environment variables:
//      - Read in stop() function
//      - Pass to library functions as needed
//      - Document in METADATA Usage section
//
// MODIFY WITH CARE (Structural Changes):
//   ⚠️ Changing orchestration flow:
//      - Ensure stop() remains thin coordinator
//      - Don't implement logic directly (extract to libraries)
//      - Maintain non-blocking philosophy
//
//   ⚠️ Changing library function calls:
//      - Verify functions exist in hooks/lib/session/
//      - Ensure parameters match function signatures
//      - Update health scoring if adding/removing calls
//
// NEVER MODIFY (Foundational Rails):
//   ❌ 4-block structure (METADATA → SETUP → BODY → CLOSING)
//   ❌ Named entry point pattern (main calls stop)
//   ❌ Non-blocking principle (session MUST stop)
//   ❌ Health scoring philosophy (Base100 objective measurement)
//
// When making changes:
//   1. Review this modification policy
//   2. Apply changes following 4-block structure
//   3. Update METADATA if health scoring changes
//   4. Test with actual Claude Code integration
//   5. Document "What/Why/How" for all changes

// LADDER AND BATON FLOW
// ======================
//
// LADDER (Dependencies - Vertical):
//   stop.go depends on:
//     ↓ hooks/lib/session (display, checks)
//     ↓ hooks/lib/activity (logging)
//   Those depend on:
//     ↓ hooks/lib/temporal (temporal context)
//     ↓ system/lib/git (git operations)
//
// BATON (Execution - Through Layers):
//   User stops session
//     ↓
//   Claude Code triggers hook
//     ↓
//   main() entry point
//     ↓
//   stop() orchestration
//     ├→ Get environment variables
//     ├→ Log activity
//     ├→ Call display functions
//     └→ Call check functions
//
// This separation enables:
//   - Reusable libraries across hooks
//   - Clear dependency management
//   - Surgical updates (change one layer without affecting others)
//   - Testable components (libraries can be tested independently)

// SURGICAL UPDATE POINTS
// =======================
//
// Common modifications and where to make them:
//
// 1. Adding new stop checks:
//    Location: checkStoppingPoint() function
//    Action: Add new session library function call
//    Example: session.CheckOpenFiles(workspace)
//
// 2. Changing display order:
//    Location: stop() function Phase 2
//    Action: Reorder session.Print* calls
//    Example: Move PrintStoppingContext before PrintStopInfo
//
// 3. Adding environment variable support:
//    Location: stop() function Phase 1
//    Action: Read new variable, pass to functions
//    Example: priority := os.Getenv("STOP_PRIORITY")
//
// 4. Changing health scoring:
//    Location: METADATA - Health Scoring
//    Action: Update point allocations
//    Example: Increase Analysis phase to 40 points
//
// 5. Creating new library functions:
//    Location: hooks/lib/session/
//    Action: Create new file or add to existing
//    Example: Create hooks/lib/session/files.go for file checks
//
// Each modification point clearly defined for surgical precision.

// PERFORMANCE CONSIDERATIONS
// ===========================
//
// Execution Speed:
//   - Hook runs synchronously (blocks session stop until complete)
//   - Target: < 1 second total execution time
//   - Display operations: ~100ms
//   - Check operations: ~200-500ms depending on workspace size
//
// Optimization Opportunities:
//   - Checks could run in parallel (currently sequential)
//   - Temporal context cached if fetched recently
//   - Git operations memoized within single hook execution
//
// Current Bottlenecks:
//   - Git status check (depends on repo size)
//   - Process enumeration (depends on running process count)
//   - Activity review (depends on log file size)
//
// Not currently optimized because:
//   - Session stop is infrequent (not hot path)
//   - User expects brief pause during stop
//   - Clarity > speed for stop sequence
//   - Will optimize if users report slowness

// TROUBLESHOOTING GUIDE
// =====================
//
// Problem: Stop hook doesn't display anything
//   Diagnosis: Hook not registered or executable not built
//   Solution: Run build.sh in hooks directory, check executable exists
//
// Problem: Temporal context missing
//   Diagnosis: Temporal utilities not available
//   Solution: Non-blocking by design, session stops anyway. Check temporal system.
//
// Problem: Stopping point checks not appearing
//   Diagnosis: NOVA_DAWN_WORKSPACE not set
//   Solution: Checks only run if workspace configured. Set environment variable.
//
// Problem: "WARNING: Failed to..." messages in output
//   Diagnosis: Library function failed (expected occasionally)
//   Solution: Non-blocking by design. Check specific library logs for details.
//
// Problem: Stop hook takes too long
//   Diagnosis: Workspace checks slow (large repo, many processes)
//   Solution: Expected for large workspaces. If excessive, can disable specific checks.
//
// Debug Process:
//   1. Check hook executable exists: ~/.claude/hooks/session/cmd-stop/stop
//   2. Check hook is registered in Claude Code settings
//   3. Run hook manually: REASON="test" NOVA_DAWN_WORKSPACE="/path" ./stop
//   4. Check activity logs: ~/.claude/system/logs/ for activity stream
//   5. Review library logs for specific failures

// RELATED COMPONENTS & DEPENDENCIES
// ==================================
//
// This hook orchestrates the following libraries:
//
// Display (hooks/lib/session/display.go):
//   - PrintStopHeader(): Stop banner with Colossians 3:23
//   - PrintStopInfo(): Stopping point check header with timestamp
//   - PrintStoppingContext(): Temporal awareness at stop
//
// Checks (hooks/lib/session/):
//   - reminders.go: RemindUncommittedWork() - Git status and uncommitted changes
//   - processes.go: CheckRunningProcessesAsReminder() - Running dev servers, etc.
//   - activity.go: CheckRecentActivity() - Recent work for context
//
// Activity (hooks/lib/activity/):
//   - LogActivity(): Records stop event to activity stream for pattern learning
//
// Temporal (hooks/lib/temporal/):
//   - GetTemporalContext(): Provides time awareness data (via display functions)
//
// System (system/lib/git/):
//   - Git operations used by session library (indirect dependency)
//
// Dependency Direction:
//   stop.go → hooks/lib → system/lib
//   Never reverse (no system/lib depending on hooks)

// FUTURE EXPANSIONS & ROADMAP
// ============================
//
// Planned Features:
//   - Session quality scoring (how faithful was this session?)
//   - Automatic TODO capture (uncommitted work → task list)
//   - Session summary (what was accomplished this session?)
//   - Stop reason categorization (natural milestone, urgent, planned)
//   - Integration with session-log for stopping patterns
//
// Research Areas:
//   - Parallel execution of checks (faster for large workspaces)
//   - Machine learning for "good stopping point" prediction
//   - Integration with project management (stop = sync tasks?)
//   - Stopping ritual prompts (reflect on session, plan next)
//
// Known Limitations:
//   - Checks run sequentially (could be parallel)
//   - No automatic commit suggestion for uncommitted work
//   - No integration with external task systems
//   - Limited session quality assessment
//
// Enhancement Opportunities:
//   - Add "quick stop" vs "thorough stop" modes
//   - Session reflection prompts for learning
//   - Automatic workspace cleanup suggestions
//   - Integration with stopping point recognition skill

// CLOSING NOTE
// ============
//
// Session stop is sacred moment - transition from work to rest.
// Like Sabbath principle in Scripture, stopping well honors work done.
//
// This hook provides:
//   - Biblical grounding (Colossians 3:23)
//   - Temporal awareness (when, how long)
//   - Quality check (did we leave things well?)
//   - Grace in transition (non-blocking even with warnings)
//
// Not just "user left" - intentional, faithful completion.
//
// Excellence in stopping reflects excellence in working:
//   - Complete what can be completed
//   - Note what remains
//   - Honor the work done
//   - Rest in faithful completion
//
// "Whatever you do, work heartily, as for the Lord and not for men."
// This includes stopping well.

// QUICK REFERENCE: USAGE EXAMPLES
// ================================
//
// Manual execution:
//   $ REASON="Task complete" NOVA_DAWN_WORKSPACE="/path" ./stop
//
// Via Claude Code (automatic):
//   User stops session → Claude Code triggers SessionStop event → hook runs
//
// Check if hook is working:
//   $ ./stop
//   (Should display stop summary, checks will be skipped if no workspace set)
//
// Debug specific issue:
//   $ REASON="debug test" NOVA_DAWN_WORKSPACE="/home/user/project" ./stop 2>&1 | tee debug.log
//
// Verify temporal awareness working:
//   $ ./stop | grep "TEMPORAL CONTEXT"
//   (Should show temporal awareness section if system available)

// ============================================================================
// END CLOSING
// ============================================================================
