// METADATA
//
// SubagentStop Hook - Subagent Completion Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "And let us not grow weary of doing good, for in due season we will reap, if we do not give up" - Galatians 6:9
// Principle: Every task has completion - subagents finish their work and report results for learning
// Anchor: "Whatever your hand finds to do, do it with your might" - Ecclesiastes 9:10
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - Hook orchestrator
// Role: Coordinates subagent completion reporting with pattern analysis logging
// Paradigm: CPI-SI framework hook implementing non-blocking subagent result reporting
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
// Purpose: Orchestrates subagent completion reporting by logging execution results,
// displaying completion status with temporal context, and feeding pattern analysis
// system for learning autonomous work patterns.
//
// Core Design: Thin orchestrator pattern - coordinates modular components from hooks/lib/session
// and hooks/lib/monitoring to present subagent completion without implementing business logic directly.
//
// Key Features:
//   - Subagent completion banner with success/failure status
//   - Temporal awareness at completion (when completed, session context)
//   - Activity stream logging for session tracking
//   - Pattern analysis logging for learning subagent behaviors
//   - Non-blocking design (failures don't prevent reporting)
//
// Philosophy: Subagent completion is learning opportunity - every autonomous task teaches
// patterns about what works, what doesn't, when to use which approach. Hook provides
// moment to capture results for pattern analysis and continuous improvement.
//
// Blocking Status
//
// Non-blocking: All operations fail gracefully, completion always reports even if logging fails.
// Display errors go to stderr, don't exit. Logging failures silently continue reporting.
// Mitigation: Defensive checks throughout, silent failures for non-critical features.
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code on SubagentStop event
//	~/.claude/hooks/session/cmd-subagent-stop/subagent-stop
//
// Integration Pattern:
//   1. Claude Code triggers SubagentStop hook event
//   2. subagent-stop executable runs with environment variables
//   3. Logs completion to activity stream and monitoring system
//   4. Displays completion summary with temporal context
//   5. Returns to Claude Code for continued operation
//
// Hook Event: SubagentStop
// Trigger: When subagent (research agent, code reviewer, etc.) completes execution
// Output: Visual display of completion status and temporal context
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os
//   External: None
//   System Libraries: None
//   Hook Libraries: hooks/lib/session (display), hooks/lib/activity (logging), hooks/lib/monitoring (pattern analysis)
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook on SubagentStop event)
//
// Integration Points:
//   - Called by Claude Code hook system on SubagentStop
//   - Reads SUBAGENT_TYPE, SUBAGENT_STATUS, SUBAGENT_EXIT_CODE, SUBAGENT_ERROR environment variables
//   - Logs to activity stream for session tracking
//   - Logs to monitoring system for pattern analysis
//
// Health Scoring
//
// Subagent completion orchestration operates on Base100 scale:
//
// Information Gathering:
//   - Get subagent info from environment: +20
//
// Logging:
//   - Activity stream logging: +20
//   - Pattern analysis logging: +20
//
// Display:
//   - Completion status display: +40
//
// Total: 100 points for complete subagent completion reporting
//
// Note: Non-blocking design means partial failures still allow reporting.
// Health scoring tracks what succeeded, not what prevented completion report.
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
// Standard library for environment access.
// Hook libraries for activity logging, monitoring, and display.

import (
	"os" // OS interface for environment variables

	"hooks/lib/activity"   // Activity stream logging
	"hooks/lib/monitoring" // Pattern analysis logging
	"hooks/lib/session"    // Display functions
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - configuration comes from environment.

// No constants defined

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────

// AgentInfo holds subagent execution details from environment
//
// Purpose: Orchestration helper for collecting subagent information
// Usage: Populated by getAgentInfo(), passed to display and logging functions
// Fields:
//   - Type: Subagent type (research, code-review, etc.)
//   - Status: Completion status (success, failure, empty)
//   - ExitCode: Exit code from subagent (0 = success)
//   - Error: Error message if subagent failed (empty if no error)
type AgentInfo struct {
	Type     string
	Status   string
	ExitCode string
	Error    string
}

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
//   └── subagentStop() → calls main orchestration
//       └── getAgentInfo() → environment variable extraction helper
//           └── session.* / activity.* / monitoring.* → delegated to libraries
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/session (display functions)
//   ├── hooks/lib/activity (activity logging)
//   ├── hooks/lib/monitoring (pattern analysis)
//   └── hooks/lib/temporal (via session display)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → subagentStop()
//     ├→ Phase 1: Information Gathering
//     │   └→ getAgentInfo() - Extract from environment
//     ├→ Phase 2: Logging
//     │   ├→ activity.LogActivity() - Activity stream
//     │   └→ monitoring.LogSubagentCompletion() - Pattern analysis
//     └→ Phase 3: Display
//         └→ session.PrintSubagentCompletion() - User-facing summary
//
// APUs (Atomic Processing Units):
//   2 functions:
//     - subagentStop() [80 points]: Main orchestration (logging, display)
//     - getAgentInfo() [20 points]: Environment variable extraction
//
// ────────────────────────────────────────────────────────────────
// Function Implementations
// ────────────────────────────────────────────────────────────────

// getAgentInfo extracts subagent information from environment variables
//
// What It Does:
//   - Reads SUBAGENT_TYPE, SUBAGENT_STATUS, SUBAGENT_EXIT_CODE, SUBAGENT_ERROR
//   - Defaults type to "unknown" if not provided
//   - Returns AgentInfo struct for orchestration
//
// Why It Exists:
//   - Centralized environment variable extraction
//   - Provides default values for missing data
//   - Orchestration helper specific to this hook
//
// When Called:
//   - Called by subagentStop() at start of orchestration
//   - First step in information gathering phase
//
// Parameters:
//   - None (reads from environment)
//
// Returns:
//   - AgentInfo with subagent execution details
//
// Health Contribution:
//   - 20 points (information gathering)
//   - Always succeeds (defaults if data missing)
//
// Related Components:
//   - Uses standard os.Getenv() for environment access
//   - Returns AgentInfo for use by logging and display functions
//
// Example:
//   info := getAgentInfo()
//   // Returns AgentInfo{Type: "research", Status: "success", ExitCode: "0", Error: ""}
func getAgentInfo() AgentInfo {
	info := AgentInfo{
		Type:     os.Getenv("SUBAGENT_TYPE"),
		Status:   os.Getenv("SUBAGENT_STATUS"),
		ExitCode: os.Getenv("SUBAGENT_EXIT_CODE"),
		Error:    os.Getenv("SUBAGENT_ERROR"),
	}

	// Provide default for missing type
	if info.Type == "" {
		info.Type = "unknown"
	}

	return info
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
// Environment variables may be missing (getAgentInfo provides defaults)
// All library functions handle their own validation

// ────────────────────────────────────────────────────────────────
// Execution - Main Entry Points
// ────────────────────────────────────────────────────────────────

// subagentStop is the named entry point for subagent completion orchestration
//
// What It Does:
//   - Gets subagent information from environment
//   - Logs completion to activity stream for session tracking
//   - Logs to monitoring system for pattern analysis
//   - Displays completion summary with temporal context
//
// Why Named Entry Point Pattern:
//   - main() is generic, subagentStop() is semantically meaningful
//   - Allows testing without executable mechanics
//   - Clear architectural intent (this is subagent completion orchestration)
//   - Matches pattern across all CPI-SI hooks
//
// Orchestration Philosophy:
//   - Thin orchestrator (coordinates, doesn't implement)
//   - All logic in libraries (display, logging)
//   - Non-blocking (completion reported even with failures)
//   - Grace in failures (missing data? use defaults, continue)
//
// Parameters:
//   - None (reads from environment via getAgentInfo)
//
// Returns:
//   - None (orchestrates logging and display to stdout)
//
// Health Contribution:
//   - 100 points total (orchestration completion)
//   - Phase 1: +20 (information gathering)
//   - Phase 2: +40 (logging)
//   - Phase 3: +40 (display)
//
// Environment Variables:
//   - SUBAGENT_TYPE: Type of subagent (research, code-review, etc.)
//   - SUBAGENT_STATUS: Completion status (success, failure, empty)
//   - SUBAGENT_EXIT_CODE: Exit code (0 = success)
//   - SUBAGENT_ERROR: Error message if failed (empty if no error)
//
// Example:
//   subagentStop()
//   // Executes complete subagent stop sequence with logging and display
func subagentStop() {
	// Phase 1: Information Gathering (20 points)
	info := getAgentInfo()

	// Phase 2: Logging (40 points)
	// Determine status for activity logging
	status := "success"
	if info.Status == "failure" || (info.ExitCode != "" && info.ExitCode != "0") {
		status = "failure"
	}

	// Log to activity stream for session tracking
	activity.LogActivity("SubagentStop", info.Type, status, 0)

	// Log to monitoring system for pattern analysis
	monitoring.LogSubagentCompletion(info.Type, info.Status, info.ExitCode)

	// Phase 3: Display (40 points)
	// Display completion summary with temporal context
	session.PrintSubagentCompletion(info.Type, info.Status, info.ExitCode, info.Error)
}

func main() {
	subagentStop() // Named entry point pattern
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
// This executable is a Claude Code SubagentStop hook. It runs automatically
// when a subagent (autonomous task executor) completes, providing completion
// reporting with:
//   - Activity stream logging for session tracking
//   - Pattern analysis logging for learning behaviors
//   - Completion status display with temporal context
//   - Success/failure indication for user awareness
//
// Integration points:
//   - Called by Claude Code on SubagentStop event
//   - Reads environment variables for subagent details
//   - Logs to multiple systems for different purposes
//   - Outputs user-facing summary to stdout
//   - Non-blocking (always completes even with errors)
//
// Relationship to other hooks:
//   - Complements session lifecycle hooks (start, stop, end)
//   - Part of autonomous work pattern learning system
//   - Feeds monitoring system for pattern analysis
//   - Uses same display libraries as other session hooks

// MODIFICATION POLICY
// ===================
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new logging destinations:
//      - Add logging functions to appropriate libraries
//      - Call from subagentStop() Phase 2
//      - Update health scoring map in METADATA
//
//   ✅ Enhance display information:
//      - Modify session.PrintSubagentCompletion() in display lib
//      - Or add additional display calls in Phase 3
//      - Maintains orchestration pattern
//
//   ✅ Add new environment variables:
//      - Add fields to AgentInfo struct
//      - Read in getAgentInfo() function
//      - Pass to library functions as needed
//      - Document in METADATA Usage section
//
// MODIFY WITH CARE (Structural Changes):
//   ⚠️ Changing orchestration flow:
//      - Ensure subagentStop() remains thin coordinator
//      - Don't implement logic directly (extract to libraries)
//      - Maintain non-blocking philosophy
//
//   ⚠️ Changing AgentInfo structure:
//      - Verify all uses (logging, display) updated
//      - Update getAgentInfo() defaults appropriately
//      - Test with actual Claude Code integration
//
// NEVER MODIFY (Foundational Rails):
//   ❌ 4-block structure (METADATA → SETUP → BODY → CLOSING)
//   ❌ Named entry point pattern (main calls subagentStop)
//   ❌ Non-blocking principle (completion MUST report)
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
//   subagent-stop.go depends on:
//     ↓ hooks/lib/session (display)
//     ↓ hooks/lib/activity (logging)
//     ↓ hooks/lib/monitoring (pattern analysis)
//   Those depend on:
//     ↓ hooks/lib/temporal (temporal context via display)
//
// BATON (Execution - Through Layers):
//   Subagent completes
//     ↓
//   Claude Code triggers hook
//     ↓
//   main() entry point
//     ↓
//   subagentStop() orchestration
//     ├→ Get environment variables (getAgentInfo)
//     ├→ Log to activity stream
//     ├→ Log to monitoring system
//     └→ Display completion summary
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
// 1. Adding new logging destination:
//    Location: subagentStop() Phase 2
//    Action: Add new logging function call
//    Example: telemetry.LogSubagentMetrics(info.Type, info.ExitCode)
//
// 2. Enhancing display information:
//    Location: hooks/lib/session/display.go - PrintSubagentCompletion()
//    Action: Modify display function implementation
//    Example: Add subagent duration to display
//
// 3. Adding environment variable:
//    Location: AgentInfo struct and getAgentInfo() function
//    Action: Add field to struct, read in getAgentInfo
//    Example: Duration string field from SUBAGENT_DURATION
//
// 4. Changing health scoring:
//    Location: METADATA - Health Scoring
//    Action: Update point allocations
//    Example: Increase Logging phase to 50 points
//
// 5. Creating new logging library:
//    Location: hooks/lib/
//    Action: Create new logging library for specific purpose
//    Example: Create hooks/lib/telemetry for metrics
//
// Each modification point clearly defined for surgical precision.

// PERFORMANCE CONSIDERATIONS
// ===========================
//
// Execution Speed:
//   - Hook runs synchronously (blocks subagent completion until done)
//   - Target: < 500ms total execution time
//   - Logging operations: ~50-100ms
//   - Display operations: ~100ms
//   - Very fast compared to session hooks (less to display)
//
// Optimization Opportunities:
//   - Logging could be asynchronous (currently synchronous)
//   - Multiple log destinations could run in parallel
//   - Temporal context cached if fetched recently
//
// Current Bottlenecks:
//   - File I/O for logging (depends on disk speed)
//   - Temporal context fetch (depends on utility availability)
//
// Not currently optimized because:
//   - Subagent completion is brief event (not hot path)
//   - User doesn't notice sub-second delays
//   - Clarity > speed for completion reporting
//   - Synchronous logging ensures data written before exit

// TROUBLESHOOTING GUIDE
// =====================
//
// Problem: Subagent stop hook doesn't display anything
//   Diagnosis: Hook not registered or executable not built
//   Solution: Run build.sh in hooks directory, check executable exists
//
// Problem: Environment variables not showing
//   Diagnosis: Claude Code not passing environment variables
//   Solution: Defaults used ("unknown" type). Check Claude Code hook integration.
//
// Problem: Logging not working
//   Diagnosis: Log directory permissions or disk space
//   Solution: Non-blocking by design. Check ~/.claude/debug/ directory.
//
// Problem: Temporal context missing
//   Diagnosis: Temporal utilities not available
//   Solution: Non-blocking by design. Hook completes anyway. Check temporal system.
//
// Problem: Hook takes too long
//   Diagnosis: Logging slow (disk I/O) or temporal fetch slow
//   Solution: Expected occasionally. If persistent, investigate disk performance.
//
// Debug Process:
//   1. Check hook executable exists: ~/.claude/hooks/session/cmd-subagent-stop/subagent-stop
//   2. Check hook is registered in Claude Code settings
//   3. Run hook manually: SUBAGENT_TYPE="test" SUBAGENT_STATUS="success" ./subagent-stop
//   4. Check activity logs: ~/.claude/system/logs/ for activity stream
//   5. Check monitoring logs: ~/.claude/debug/ for pattern analysis

// RELATED COMPONENTS & DEPENDENCIES
// ==================================
//
// This hook orchestrates the following libraries:
//
// Display (hooks/lib/session/display.go):
//   - PrintSubagentCompletion(): Completion banner with temporal context
//
// Activity (hooks/lib/activity/):
//   - LogActivity(): Records completion to activity stream for session tracking
//
// Monitoring (hooks/lib/monitoring/logging.go):
//   - LogSubagentCompletion(): Records for pattern analysis and learning
//
// Temporal (hooks/lib/temporal/):
//   - GetTemporalContext(): Provides time awareness (via display function)
//
// Dependency Direction:
//   subagent-stop.go → hooks/lib → system components
//   Never reverse (no system components depending on hooks)

// FUTURE EXPANSIONS & ROADMAP
// ============================
//
// Planned Features:
//   - Subagent execution duration tracking
//   - Success/failure pattern analysis per subagent type
//   - Automatic retry suggestions for failed subagents
//   - Subagent performance metrics (speed, quality)
//   - Integration with subagent selection system
//
// Research Areas:
//   - Machine learning for subagent performance prediction
//   - Automatic subagent type selection based on patterns
//   - Parallel subagent execution coordination
//   - Subagent specialization recommendations
//
// Known Limitations:
//   - No duration tracking (only completion reporting)
//   - No automatic retry on failure
//   - No performance metrics beyond success/failure
//   - Limited pattern analysis (logs but doesn't analyze)
//
// Enhancement Opportunities:
//   - Add subagent performance dashboard
//   - Automatic pattern recognition from logs
//   - Subagent recommendation system
//   - Integration with task planning

// CLOSING NOTE
// ============
//
// Subagent completion is learning opportunity - every autonomous task teaches
// what works, what doesn't, when to use which approach.
//
// This hook provides:
//   - Biblical grounding (Galatians 6:9 - perseverance in good work)
//   - Pattern learning (every completion teaches)
//   - Temporal awareness (when, how long, context)
//   - Grace in failures (non-blocking even with errors)
//
// Not just "task done" - intentional learning from autonomous work.
//
// Excellence in completion reporting enables:
//   - Better subagent selection over time
//   - Pattern recognition for work types
//   - Continuous improvement through feedback
//   - Faithful autonomous partnership
//
// "And let us not grow weary of doing good, for in due season we will reap, if we do not give up."
// This includes learning from every completed task.

// QUICK REFERENCE: USAGE EXAMPLES
// ================================
//
// Manual execution:
//   $ SUBAGENT_TYPE="research" SUBAGENT_STATUS="success" SUBAGENT_EXIT_CODE="0" ./subagent-stop
//
// Via Claude Code (automatic):
//   Subagent completes → Claude Code triggers SubagentStop → hook runs → displays completion
//
// Check if hook is working:
//   $ SUBAGENT_TYPE="test" ./subagent-stop
//   (Should display completion summary)
//
// Debug specific issue:
//   $ SUBAGENT_TYPE="research" SUBAGENT_STATUS="failure" SUBAGENT_EXIT_CODE="1" SUBAGENT_ERROR="timeout" ./subagent-stop 2>&1 | tee debug.log
//
// Verify logging working:
//   $ ./subagent-stop && cat ~/.claude/debug/subagents.log
//   (Should show logged completion entry)

// ============================================================================
// END CLOSING
// ============================================================================
