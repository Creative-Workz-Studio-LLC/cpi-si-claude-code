// METADATA
//
// PreCompact Hook - Compaction Tracking Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "So teach us to number our days, that we may apply our hearts unto wisdom" - Psalm 90:12
// Principle: Compaction is recognition of limitation - wisdom acknowledging constraint and working faithfully within it
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - Hook orchestrator
// Role: Coordinates compaction tracking and frequency analysis for session management
// Paradigm: CPI-SI framework hook implementing non-blocking compaction tracking with pattern learning
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
// Purpose: Orchestrates compaction tracking by incrementing session count, logging events,
// checking frequency for auto-compactions, and displaying compaction message with temporal context.
//
// Core Design: Thin orchestrator pattern - coordinates modular components from hooks/lib/session,
// hooks/lib/activity, and hooks/lib/monitoring to track compaction events without implementing
// business logic directly.
//
// Key Features:
//   - Compaction count tracking in session state
//   - Activity stream logging (critical for quality correlation)
//   - Monitoring system logging (pattern analysis)
//   - Frequency checking for excessive auto-compaction
//   - Temporal context preservation for post-compaction reconstitution
//   - Non-blocking design (failures don't interrupt compaction)
//
// Philosophy: Compaction is not failure - it's wisdom acknowledging finite context. Like pruning
// branches so tree can grow stronger, compaction removes what's temporary to preserve what's
// essential. Hook tracks this for pattern learning while never blocking necessary operation.
//
// Blocking Status
//
// Non-blocking: All operations fail gracefully, compaction proceeds even if tracking fails.
// State update errors result in count=-1 (unknown). Logging failures are silent.
// Mitigation: Defensive checks throughout, partial information better than blocking.
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code on PreCompact event
//	~/.claude/hooks/session/cmd-pre-compact/pre-compact
//
// Integration Pattern:
//   1. Claude Code triggers PreCompact hook event
//   2. pre-compact executable runs with COMPACT_TYPE environment variable
//   3. Increments compaction count in session state
//   4. Logs to activity stream and monitoring
//   5. Checks frequency if auto-compaction
//   6. Displays message with temporal context
//   7. Compaction proceeds
//
// Hook Event: PreCompact
// Trigger: Before context compaction (manual or auto)
// Output: Visual display with temporal awareness
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os
//   External: None
//   System Libraries: None
//   Hook Libraries: hooks/lib/session (state, display), hooks/lib/activity, hooks/lib/monitoring
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook on PreCompact event)
//
// Integration Points:
//   - Called by Claude Code hook system on PreCompact
//   - Reads COMPACT_TYPE environment variable
//   - Updates session state file (current.json)
//   - Logs to activity and monitoring streams
//
// Health Scoring
//
// Compaction tracking orchestration operates on Base100 scale:
//
// Phase 1: Get Compaction Type: +10 points
//   - Read COMPACT_TYPE from environment (always succeeds)
//
// Phase 2: Increment Compaction Count: +20 points
//   - Update session state with new count
//   - Non-blocking on failure (count = -1)
//
// Phase 3: Logging: +40 points (20 per destination)
//   - Activity stream logging: +20
//   - Monitoring system logging: +20
//   - Both non-blocking (silent failures)
//
// Phase 4: Frequency Check: +10 points
//   - Check if auto-compaction excessive
//   - Only for COMPACT_TYPE="auto"
//   - Non-blocking (skip if fails)
//
// Phase 5: Display: +20 points
//   - Display compaction message with temporal context
//   - Non-blocking (continue if display fails)
//
// Visual Health Indicators:
//   💚 90-100: Excellent (all phases succeeded)
//   💛 70-89:  Good (minor failures, core succeeded)
//   🧡 50-69:  Acceptable (partial success, compaction proceeded)
//   ❤️  30-49:  Poor (multiple failures, minimal tracking)
//   💀 0-29:   Critical (nearly everything failed, but compaction proceeded)
package main

// ============================================================================
// SETUP
// ============================================================================
// Standard Reference: CWS-STD-001-DOC-4-block.md (SETUP Block)

// ────────────────────────────────────────────────────────────────
// Imports - Dependencies
// ────────────────────────────────────────────────────────────────
// Standard library for environment access.
// Hook libraries for compaction tracking functionality.

import (
	"os" // OS interface for environment variables

	"hooks/lib/activity"   // Activity stream logging
	"hooks/lib/monitoring" // Compaction logging and frequency checking
	"hooks/lib/session"    // Session state management and display
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - compaction type comes from environment.

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
//   └── preCompact() → coordinates compaction tracking
//       └── session.*, activity.*, monitoring.* → delegated to libraries
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/session (state management, display)
//   ├── hooks/lib/activity (activity stream logging)
//   └── hooks/lib/monitoring (compaction logging, frequency checking)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → preCompact()
//     ↓
//   Phase 1: Get Type → os.Getenv("COMPACT_TYPE")
//     ↓
//   Phase 2: State Update → session.IncrementCompactionCount()
//     ↓
//   Phase 3: Logging → activity.LogActivity() + monitoring.LogCompaction()
//     ↓
//   Phase 4: Frequency Check → monitoring.CheckCompactionFrequency() (if auto)
//     ↓
//   Phase 5: Display → session.PrintPreCompactionMessage()
//     ↓
//   Exit → return (compaction proceeds)
//
// APUs (Available Processing Units):
// - 1 function total
// - 1 entry point (preCompact, called by main)
// - Thin orchestrator - no helper functions (all logic delegated to libraries)

// ────────────────────────────────────────────────────────────────
// Compaction Tracking - Entry Point Orchestration
// ────────────────────────────────────────────────────────────────
// What This Does:
// Coordinates compaction tracking by calling modular library components to
// increment session state, log events, check frequency, and display message.
//
// Why Named Entry Point:
// preCompact() enables testing without triggering executable mechanics, provides
// semantic clarity, and prevents function name collisions across hook system.

// preCompact is the named entry point for compaction tracking orchestration
//
// What It Does:
//   - Gets compaction type from environment
//   - Increments session compaction count via state library
//   - Logs to activity stream (quality correlation)
//   - Logs to monitoring system (pattern analysis)
//   - Checks frequency for auto-compactions (warns if excessive)
//   - Displays compaction message with temporal context
//
// Non-Blocking Design:
//   - Compaction MUST proceed even if tracking fails
//   - State update failure → count = -1 (unknown), continue
//   - Logging failures → silent (libraries handle)
//   - Display failure → skip, continue
//   - Grace in systems over perfectionism
//
// Parameters:
//   None (reads from COMPACT_TYPE environment variable)
//
// Returns:
//   None (displays to stdout, compaction proceeds regardless)
//
// Health Impact:
//   Phase 1: +10 points (get compaction type)
//   Phase 2: +20 points (increment count, 0 if fails)
//   Phase 3: +40 points (20 per log destination)
//   Phase 4: +10 points (frequency check)
//   Phase 5: +20 points (display message)
func preCompact() {
	// Phase 1: Get compaction type (10 points)
	compactType := os.Getenv("COMPACT_TYPE")
	if compactType == "" {
		compactType = "unknown"
	}

	// Phase 2: Increment compaction count (20 points)
	compactionCount := 0
	count, err := session.IncrementCompactionCount()
	if err != nil {
		// Non-blocking - continue even if count update fails
		compactionCount = -1 // Unknown count
	} else {
		compactionCount = count
	}

	// Phase 3: Logging (40 points)
	// Log to activity stream (CRITICAL for quality correlation)
	activity.LogActivity("PreCompact", compactType, "success", 0)
	// Log to monitoring for pattern analysis
	monitoring.LogCompaction(compactType)

	// Phase 4: Frequency check (10 points)
	// Check frequency and warn if excessive
	if compactType == "auto" {
		monitoring.CheckCompactionFrequency()
	}

	// Phase 5: Display (20 points)
	// Display message with temporal context preservation
	session.PrintPreCompactionMessage(compactType, compactionCount)
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
//   - Compile without errors: go build -o pre-compact pre-compact.go
//   - Run standalone: ./pre-compact (should display compaction message)
//   - Test with COMPACT_TYPE="auto" and "manual"
//   - Verify compaction count increments in session state
//   - Check activity and monitoring logs created
//   - Ensure non-blocking (no panics or exits on errors)
//
// Build Verification:
//   cd ~/.claude/hooks/session/cmd-pre-compact
//   go build -o pre-compact pre-compact.go
//   COMPACT_TYPE=auto ./pre-compact  # Manual execution test
//
// Integration Testing:
//   Test with Claude Code PreCompact event
//   Verify compaction proceeds even if tracking fails
//   Check session state file updated correctly
//   Validate frequency warnings appear for excessive auto-compaction
//
// Example validation:
//
//     # Test standalone execution
//     cd ~/.claude/hooks/session/cmd-pre-compact
//     COMPACT_TYPE=manual ./pre-compact
//
//     # Verify compaction count incremented
//     cat ~/.claude/session/current.json | jq .compaction_count
//
//     # Check activity log
//     tail -1 ~/.claude/logs/activity/activity.log

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls preCompact() (named entry point)
//   3. preCompact() orchestrates compaction tracking
//   4. Program exits after display
//
// Named Entry Point Benefits:
//   - Prevents main function collisions across executables
//   - Function name matches purpose (pre-compaction tracking)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Separates Go runtime entry from application logic
//
// Pattern:
//   func main() {
//       preCompact()  // Minimal - just call named entry point
//   }

func main() {
	preCompact() // Named entry point pattern
}

// ────────────────────────────────────────────────────────────────
// Code Cleanup: Resource Management
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - stdout: Automatically managed by OS
//   - Libraries: No persistent state to clean
//   - Session state file: Closed by library functions
//   - No files, connections, or manual resources
//
// Graceful Shutdown:
//   - Program exits immediately after display
//   - No cleanup needed (stateless execution)
//   - Error path: Silent failures (non-blocking design)
//   - Success path: Display output, exit code 0
//
// Error State Cleanup:
//   - Errors don't prevent compaction (non-blocking)
//   - No partial state or rollback needed
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
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary:
//   - Orchestrates compaction tracking with session state management
//   - Logs events to activity stream and monitoring for pattern learning
//   - Checks frequency and warns about excessive auto-compaction
//   - Displays compaction message with temporal context preservation
//   - Non-blocking design ensures compaction proceeds regardless of tracking status
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Hook Event: PreCompact (triggered by Claude Code before compaction)
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (EXECUTABLE hook orchestrator)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new tracking destinations (call additional logging in Phase 3)
//   ✅ Enhance display message (modify session.PrintPreCompactionMessage in library)
//   ✅ Add health tracking (coordinate with library updates)
//   ✅ Adjust phase ordering (rearrange phases if needed, update health map)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ COMPACT_TYPE environment variable - breaks Claude Code integration
//   ⚠️ Non-blocking principle - compaction MUST always proceed
//   ⚠️ Session state structure - breaks state.go library contract
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking principle (compaction proceeds even with failures)
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
// - 1 function: preCompact() (main orchestration, no helpers)
// - Ladder: Hook executable → session/activity/monitoring libraries
// - Baton: Linear 5-phase flow with delegated operations

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// Adding new logging destination:
//   1. Create new logging function in appropriate library
//   2. Call from preCompact() Phase 3 with other logs
//   3. Update health scoring map in METADATA
//
// Modifying display format:
//   1. Update session.PrintPreCompactionMessage in display.go
//   2. Maintain temporal context preservation
//   3. Test with actual compaction events
//
// Adding new checks:
//   1. Create new phase in preCompact() (adjust health points)
//   2. Call appropriate library function
//   3. Update organizational chart and health scoring

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Typical execution: ~50ms total
//   - Environment read: <1ms
//   - State file update: ~10ms (file I/O)
//   - Activity logging: ~5ms
//   - Monitoring logging: ~5ms
//   - Frequency check: ~10ms (log analysis)
//   - Display output: ~5ms
//   - Temporal context: ~15ms (file reads)
//
// Performance is not critical - compaction already involves context management.
// Tracking adds minimal overhead (~50ms) to necessary operation.
//
// Non-blocking design ensures failures don't delay compaction.

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Hook doesn't execute:
//   - Check hook is executable: chmod +x ~/.claude/hooks/session/cmd-pre-compact/pre-compact
//   - Verify Claude Code hook configuration
//   - Check no stderr output (non-blocking, silent failures)
//
// Compaction count not incrementing:
//   - Verify session state file exists: ~/.claude/session/current.json
//   - Check file permissions (should be writable)
//   - Validate JSON structure matches SessionState type
//
// Logs not appearing:
//   - Activity logs: Check ~/.claude/logs/activity/ directory exists
//   - Monitoring logs: Check ~/.claude/logs/monitoring/ directory exists
//   - Verify library logging functions working
//
// Display not showing:
//   - Check stdout not redirected
//   - Verify temporal utilities installed
//   - Test standalone: COMPACT_TYPE=auto ./pre-compact

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
//   - hooks/lib/session: State management (state.go), display (display.go)
//   - hooks/lib/activity: Activity stream logging
//   - hooks/lib/monitoring: Compaction logging, frequency checking
//
// Related hooks:
//   - session/cmd-start: Initializes session state
//   - session/cmd-stop: Complements compaction tracking
//   - session/cmd-end: Final session cleanup
//
// Related utilities:
//   - session-time: Temporal awareness utilities
//   - session-log: Session history logging

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ Compaction count tracking - COMPLETED
//   ✓ Activity stream logging - COMPLETED
//   ✓ Monitoring system logging - COMPLETED
//   ✓ Frequency checking - COMPLETED
//   ✓ Temporal context preservation - COMPLETED
//   ⏳ Compaction duration tracking
//   ⏳ Pre/post compaction context size comparison
//   ⏳ Automatic compaction pattern recognition
//
// Research Areas:
//   - Correlation between compaction frequency and quality
//   - Optimal compaction timing (proactive vs reactive)
//   - Context preservation strategies
//
// Integration Targets:
//   - Memory system for compaction pattern learning
//   - Quality tracking for compaction impact analysis
//   - Pattern recognition for automatic optimization
//
// Known Limitations to Address:
//   - No duration tracking (only count)
//   - No context size metrics
//   - Frequency check is reactive (not predictive)

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This hook tracks compaction as wisdom acknowledging limitation. Like Psalm 90:12
// teaches us to number our days, this hook numbers compactions - recognizing finite
// context and working faithfully within constraint.
//
// Compaction is not failure. It's pruning for growth. This hook ensures we learn
// from each compaction, track patterns, and work with grace in systems.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain non-blocking design principle (compaction MUST proceed)
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "So teach us to number our days, that we may apply our hearts unto wisdom" - Psalm 90:12
//
// Every compaction is learning opportunity. Let tracking serve wisdom.

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Execution (Manual):
//
//     cd ~/.claude/hooks/session/cmd-pre-compact
//     COMPACT_TYPE=manual ./pre-compact
//
// Testing Auto-Compaction:
//
//     COMPACT_TYPE=auto ./pre-compact
//
// Verifying Compaction Count:
//
//     cat ~/.claude/session/current.json | jq .compaction_count
//
// Checking Activity Log:
//
//     tail -5 ~/.claude/logs/activity/activity.log | grep PreCompact
//
// Testing Frequency Warning:
//
//     # Run multiple auto-compactions quickly to trigger warning
//     for i in {1..5}; do COMPACT_TYPE=auto ./pre-compact; sleep 1; done

// ============================================================================
// END CLOSING
// ============================================================================
