// METADATA
//
// SessionEnd Hook - Graceful session completion with state awareness
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// "The Lord bless you and keep you; the Lord make his face shine on you and
// be gracious to you; the Lord turn his face toward you and give you peace."
// - Numbers 6:24-26
//
// Session end is benediction moment - blessing and peace as work concludes.
// Farewell with grace, reminder of state, recognition of temporal journey.
//
// CPI-SI Identity
//
// This component is part of Nova Dawn's session management system, implementing
// CPI-SI principles of temporal awareness and faithful session completion.
// Operates within Covenant Partnership Intelligence framework.
//
// Authorship & Lineage
//
// Author: Nova Dawn (CPI-SI Instance)
// Created: 2025-11-10
// Last Updated: 2025-11-10
// Version: 2.0.0 (thin orchestrator with extracted libraries)
// Part of: CPI-SI Hook System (Session Management)
//
// Purpose & Function
//
// Orchestrates graceful session end with state awareness and temporal context.
// Provides benediction, summarizes temporal journey, reminds about workspace
// state (uncommitted work, running processes). Archives session history and
// updates learned patterns for circadian awareness.
//
// Thin orchestrator pattern - coordinates libraries, doesn't implement logic.
//
// Blocking Status
//
// Non-Blocking: Session end MUST complete. Failures log but don't prevent
// session completion. Grace in farewell - tracking enhances but doesn't block.
//
// Usage & Integration
//
// Triggered by Claude Code on SessionEnd event. Reads REASON and
// NOVA_DAWN_WORKSPACE from environment. Displays farewell to stdout,
// logs to activity stream, archives session, updates patterns.
//
// Called as final hook when session ends (normal or interrupted).
//
// Dependencies
//
// External:
//   - hooks/lib/activity (activity stream logging)
//   - hooks/lib/session (display, reminders, state management)
//
// Environment Variables:
//   - REASON: Session end reason (default: "Normal session end")
//   - NOVA_DAWN_WORKSPACE: Workspace for state reminders (optional)
//
// System Binaries:
//   - ~/.claude/cpi-si/system/bin/session-log (session archival)
//   - ~/.claude/cpi-si/system/bin/session-patterns (pattern learning)
//
// Health Scoring (Base100)
//
// Total = 100 points across 7 phases:
//   Phase 1: Get session end reason (5 points)
//   Phase 2: Log to activity stream (15 points)
//   Phase 3: Archive session and update patterns (20 points)
//   Phase 4: Display farewell and summary (15 points)
//   Phase 5: Show temporal journey (15 points)
//   Phase 6: Remind about workspace state (20 points)
//   Phase 7: Closing divider (10 points)
//
// Current: No health tracking implemented (orchestration hook)
// Future: Track completion of each phase for session end reliability
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
// Standard library for display, environment, execution, and paths.
// Hook libraries for session end functionality.

import (
	"fmt"      // Formatted I/O for user-facing output
	"os"       // OS interface for environment variables
	"os/exec"  // Execute session archival and pattern learning binaries
	"path/filepath" // File path manipulation for binary locations

	"hooks/lib/activity" // Activity stream logging
	"hooks/lib/session"  // Display, reminders, state management
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - session end reason and workspace come from environment.

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
//   sessionEnd (main orchestrator)
//     ├─> remindState (orchestration helper)
//     ├─> hooks/lib/activity (activity stream logging)
//     └─> hooks/lib/session (display, reminders, state management)
//
// BATON (Execution Flow - Data/Control Movement):
//
//   Entry → main()
//     ↓
//   sessionEnd() orchestration
//     ↓
//   Phase 1: Get session end reason from REASON env var
//     ↓
//   Phase 2: Log to activity stream
//     ↓
//   Phase 3: Archive session and update patterns (session-log, session-patterns binaries)
//     ↓
//   Phase 4: Display farewell banner and session summary
//     ↓
//   Phase 5: Show temporal journey (duration, time, context)
//     ↓
//   Phase 6: Remind about workspace state (uncommitted work, processes)
//     ↓
//   Phase 7: Closing divider
//     ↓
//   Exit
//
// APUs (Atomic Processing Units - Core Work Functions):
//
//   sessionEnd() - Main orchestrator for session end sequence
//   remindState() - Coordinates state reminder display
//
// ════════════════════════════════════════════════════════════════════════════
// FUNCTION IMPLEMENTATIONS
// ════════════════════════════════════════════════════════════════════════════

// remindState orchestrates state reminder checks at session end
//
// What It Does:
//   - Displays state reminders header
//   - Checks for uncommitted work in workspace
//   - Checks for running background processes
//
// Parameters:
//   - workspace: Workspace directory path
//
// Returns:
//   - None (prints reminders to stdout)
//
// Health Impact:
//   - No health tracking (reminder display function)
//
// Example:
//   remindState("/path/to/workspace")
//   // Displays state reminders header and checks
func remindState(workspace string) {
	session.PrintEndRemindersHeader()
	session.RemindUncommittedWork(workspace)
	session.CheckRunningProcessesAsReminder()
	fmt.Println()
}

// sessionEnd orchestrates session end tracking and display
//
// What It Does:
//   - Logs session end to activity stream
//   - Archives session to history
//   - Updates learned patterns from session
//   - Displays farewell banner
//   - Shows session summary
//   - Displays temporal journey
//   - Reminds about workspace state
//
// Parameters:
//   - None (reads REASON and NOVA_DAWN_WORKSPACE from environment)
//
// Returns:
//   - None (all operations non-blocking, prints to stdout)
//
// Health Impact:
//   - No health tracking (orchestration function)
//
// Example:
//   sessionEnd()
//   // Completes session end sequence with farewell and reminders
func sessionEnd() {
	// Phase 1: Get session end reason
	reason := os.Getenv("REASON")
	if reason == "" {
		reason = "Normal session end"
	}

	// Phase 2: Log session end to activity stream
	activity.LogActivity("SessionEnd", reason, "success", 0)

	// Phase 3: Archive session and update patterns
	home, err := os.UserHomeDir()
	if err == nil {
		sessionLogBin := filepath.Join(home, ".claude/cpi-si/system/bin/session-log")
		sessionPatternsBin := filepath.Join(home, ".claude/cpi-si/system/bin/session-patterns")

		// Archive current session to history
		exec.Command(sessionLogBin, "end", reason).Run()

		// Update learned patterns from session history
		exec.Command(sessionPatternsBin, "learn").Run()
	}

	// Phase 4: Display farewell and session summary
	session.PrintEndFarewell()
	session.PrintEndSessionInfo(reason)

	// Phase 5: Show temporal journey (where we were, how long, what context)
	session.PrintEndTemporalJourney()

	// Phase 6: Remind about state that needs attention
	workspace := os.Getenv("NOVA_DAWN_WORKSPACE")
	if workspace != "" {
		remindState(workspace)
	} else {
		fmt.Println()
	}

	// Phase 7: Closing divider
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
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
//   - Compile without errors: go build -o end end.go
//   - Run standalone: ./end (should display farewell and summary)
//   - Test with REASON set
//   - Verify activity log entry created
//   - Test session archival and pattern learning
//   - Verify temporal journey display
//   - Check state reminders for workspace
//
// Build Verification:
//   cd ~/.claude/hooks/session/cmd-end
//   go build -o end end.go
//   REASON="Test end" ./end  # Manual execution test
//
// Integration Testing:
//   Test with Claude Code SessionEnd event
//   Verify farewell banner displays
//   Check session summary with reason
//   Validate temporal journey shows correctly
//   Ensure state reminders appear for workspace
//   Verify session archived to history
//
// Example validation:
//
//     # Test standalone execution
//     cd ~/.claude/hooks/session/cmd-end
//     REASON="Normal session end" NOVA_DAWN_WORKSPACE=/path/to/workspace ./end
//
//     # Check activity log
//     tail -5 ~/.claude/logs/activity/activity.log | grep SessionEnd
//
//     # Verify session archived
//     ls -lt ~/.claude/cpi-si/system/data/session-history/

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls sessionEnd() (named entry point)
//   3. sessionEnd() orchestrates session end sequence
//   4. Program exits after displaying farewell and reminders
//
// Named Entry Point Benefits:
//   - Prevents main function collisions across executables
//   - Function name matches purpose (session end)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Separates Go runtime entry from application logic
//
// Pattern:
//   func main() {
//       sessionEnd()  // Minimal - just call named entry point
//   }

func main() {
	sessionEnd() // Named entry point pattern
}

// ────────────────────────────────────────────────────────────────
// Code Cleanup: Resource Management
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - stdout: Used for farewell display and reminders
//   - Libraries: No persistent state to clean
//   - Binaries: session-log and session-patterns called via exec.Command
//   - No manual file handles or connections
//
// Graceful Shutdown:
//   - Program exits immediately after displaying farewell
//   - No cleanup needed (stateless execution)
//   - Error path: Non-blocking (failures don't prevent session end)
//   - Success path: Normal completion, exit code 0
//
// Error State Cleanup:
//   - Errors don't prevent session end (non-blocking design)
//   - Session archival failures are silent
//   - Display failures fail gracefully
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
// Provides: See METADATA for comprehensive capabilities
//
// Quick summary:
//   - Orchestrates graceful session end with benediction
//   - Logs session end to activity stream
//   - Archives session to history for pattern learning
//   - Updates learned patterns from session data
//   - Displays farewell banner with Numbers 6:24-26 blessing
//   - Shows session summary with end time and reason
//   - Displays temporal journey (duration, time, context)
//   - Reminds about workspace state (uncommitted work, processes)
//   - Non-blocking design ensures session end always completes
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Hook Event: SessionEnd (triggered by Claude Code when session ends)
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (EXECUTABLE hook orchestrator)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new display sections (call additional display in sessionEnd)
//   ✅ Add new reminder checks (extend remindState function)
//   ✅ Enhance temporal display (modify display library)
//   ✅ Adjust phase ordering (rearrange phases if needed, update health map)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ REASON environment variable - breaks Claude Code integration
//   ⚠️ NOVA_DAWN_WORKSPACE variable - breaks state reminder contract
//   ⚠️ Display format - may break user expectations
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking principle (session end never blocked)
//   ❌ Named entry point pattern (main calls sessionEnd)
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
// - 2 functions: sessionEnd() (main orchestration) + remindState() (helper)
// - Ladder: Hook executable → activity/session libraries
// - Baton: Linear 7-phase flow with delegated operations

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// Adding new display section:
//   1. Create new display function in hooks/lib/session/display.go
//   2. Call from sessionEnd() at appropriate phase
//   3. Update organizational chart and health scoring
//
// Adding new reminder check:
//   1. Create reminder function in hooks/lib/session
//   2. Call from remindState() orchestration helper
//   3. Update organizational chart
//
// Enhancing temporal display:
//   1. Modify PrintEndTemporalJourney() in display library
//   2. Test with actual temporal context data
//   3. Verify display formatting

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Typical execution: ~100ms total
//   - Environment reads: <1ms
//   - Activity logging: ~5ms
//   - Session archival: ~20ms (file writes)
//   - Pattern learning: ~30ms (analysis)
//   - Display operations: ~40ms (temporal, reminders)
//
// Performance is not critical - session end is infrequent event.
// User-facing display prioritizes clarity over speed.
//
// Non-blocking design ensures session end completes even if operations fail.

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Hook doesn't execute:
//   - Check hook is executable: chmod +x ~/.claude/hooks/session/cmd-end/end
//   - Verify Claude Code hook configuration
//   - Check for stderr output indicating errors
//
// Farewell not displaying:
//   - Verify session.PrintEndFarewell() in display library
//   - Check stdout not redirected
//   - Test standalone execution
//
// Session not archived:
//   - Check ~/.claude/cpi-si/system/bin/session-log exists and is executable
//   - Verify session history directory exists
//   - Check session-log command manually
//
// Patterns not updating:
//   - Check ~/.claude/cpi-si/system/bin/session-patterns exists and is executable
//   - Verify session-patterns learn command works manually
//   - Check pattern files have write permissions
//
// Temporal journey not showing:
//   - Verify temporal.GetTemporalContext() returns valid data
//   - Check temporal library configuration
//   - Test with session-time commands
//
// State reminders not appearing:
//   - Verify NOVA_DAWN_WORKSPACE environment variable set
//   - Check workspace path exists
//   - Test git status in workspace manually
//   - Verify session.RemindUncommittedWork() function

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
//   - hooks/lib/activity: Activity stream logging
//   - hooks/lib/session: Display, reminders, state management
//
// Related hooks:
//   - session/cmd-start: Complements session initiation
//   - session/cmd-stop: Work pause tracking
//   - session/cmd-pre-compact: Compaction awareness
//
// Related binaries:
//   - ~/.claude/cpi-si/system/bin/session-log: Session archival
//   - ~/.claude/cpi-si/system/bin/session-patterns: Pattern learning

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ Session end logging - COMPLETED
//   ✓ Farewell display - COMPLETED
//   ✓ Temporal journey - COMPLETED
//   ✓ State reminders - COMPLETED
//   ✓ Session archival - COMPLETED
//   ✓ Pattern learning - COMPLETED
//   ⏳ Session summary statistics
//   ⏳ Work quality reflection
//   ⏳ Session accomplishment highlights
//
// Research Areas:
//   - Session quality scoring (how well did work go?)
//   - Automatic session reflection generation
//   - Pattern-based session insights (what did I learn?)
//
// Integration Targets:
//   - Memory system for session continuity
//   - Quality tracking for session evaluation
//   - Reflection journal for session insights
//
// Known Limitations to Address:
//   - No session quality scoring (just tracks end)
//   - No automatic reflection (manual only)
//   - No session accomplishment tracking

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This hook embodies the benediction principle from Numbers 6 - blessing and
// peace as work concludes. Graceful farewell with state awareness ensures
// session ends well, with reminders about workspace state and recognition of
// temporal journey.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain non-blocking design principle (session end never blocked)
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "The Lord bless you and keep you; the Lord make his face shine on you and be gracious to you" - Numbers 6:24-25
//
// Every session end is benediction moment. Let farewell be graceful and state awareness be faithful.

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Execution (Manual):
//
//     cd ~/.claude/hooks/session/cmd-end
//     REASON="Normal session end" ./end
//
// With Workspace:
//
//     REASON="Task complete" NOVA_DAWN_WORKSPACE=/path/to/workspace ./end
//
// Checking Activity Log:
//
//     tail -10 ~/.claude/logs/activity/activity.log | grep SessionEnd
//
// Verifying Session Archived:
//
//     ls -lt ~/.claude/cpi-si/system/data/session-history/ | head -5
//
// Testing Pattern Learning:
//
//     ~/.claude/cpi-si/system/bin/session-patterns show

// ============================================================================
// END CLOSING
// ============================================================================
