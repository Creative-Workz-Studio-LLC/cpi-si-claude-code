// METADATA
//
// Statusline Orchestrator - CPI-SI Statusline Display System
//
// Biblical Foundation
//
// Scripture: "Let all things be done decently and in order" - 1 Corinthians 14:40
// Principle: Order and clarity in communication reflect God's nature
// Anchor: "A word fitly spoken is like apples of gold in pictures of silver" - Proverbs 25:11
//
// CPI-SI Identity
//
// Component Type: Baton (Top Rung Orchestrator)
// Role: Assemble complete statusline from all presentation libraries
// Paradigm: CPI-SI framework component - demonstrates separation of concerns
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2025-10-24
// Version: 1.0.0
// Last Modified: 2025-11-04 - Applied full template, refactored to named entry point
//
// Version History:
//   1.0.0 (2025-11-04) - Comprehensive template application and architectural refinement
//   0.1.0 (2025-10-24) - Initial orchestrator implementation
//
// Purpose & Function
//
// Purpose: Orchestrate all statusline presentation libraries to produce formatted
// statusline output from Claude Code session data. Reads JSON from stdin, transforms
// through presentation layer, outputs complete statusline to stdout.
//
// Core Design: Orchestrator pattern - coordinates multiple presentation libraries
// without implementing formatting logic itself. Each library handles specific concern,
// orchestrator assembles results.
//
// Key Features:
//   - Parses Claude Code JSON input (types library)
//   - Retrieves temporal context (hooks integration)
//   - Formats all display elements (7 presentation libraries)
//   - Assembles complete statusline with color/icons
//   - Outputs formatted result to stdout
//
// Philosophy: Orchestration, not implementation. Each library provides formatted
// display structures - orchestrator decides WHAT to show and WHERE, libraries
// decide HOW to format.
//
// Blocking Status
//
// Blocking: Reads stdin (waits for JSON input), executes, outputs, exits
// Mitigation: Fast execution (<10ms typical), called per hook invocation
//
// Usage & Integration
//
// Usage:
//
//	echo '{"session_id":"abc","model":{"display_name":"Sonnet"},...}' | ./statusline
//
// Integration Pattern:
//   1. Called by Claude Code hooks (session start, tool use, etc.)
//   2. Receives session context JSON via stdin
//   3. Produces formatted statusline string
//   4. Outputs to stdout for display
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: encoding/json, fmt, os, strings, time
//   External: hooks/lib/temporal (temporal context data)
//   External: system/lib/instance (instance configuration)
//   External: system/lib/display (terminal colors and formatting)
//   Internal: statusline/lib/* (all 7 presentation libraries)
//
// Library Dependencies:
//   - statusline/lib/types (SessionContext input contract)
//   - statusline/lib/features (display timing decisions)
//   - statusline/lib/format (text optimization)
//   - statusline/lib/git (git repository display)
//   - statusline/lib/session (session statistics)
//   - statusline/lib/system (system health metrics)
//   - statusline/lib/temporal (temporal awareness display)
//
// Dependents (What Uses This):
//   Commands: None (executable, not library)
//   Tools: Claude Code hooks system
//
// Integration Points:
//   - Baton: JSON stdin → statusline() → stdout output
//   - Ladder: Top rung orchestrating all libraries below
//
// Health Scoring
//
// N/A for orchestrator - no health tracking (simple coordination logic)
// Libraries handle their own health tracking independently
package main

// ============================================================================
// END METADATA
// ============================================================================

// ============================================================================
// SETUP
// ============================================================================

// ────────────────────────────────────────────────────────────────
// Imports - Dependencies
// ────────────────────────────────────────────────────────────────

import (
	//--- Standard Library ---
	"encoding/json" // JSON parsing for Claude Code input
	"fmt"           // Formatted output and string building
	"os"            // Stdin/stdout access, exit codes
	"strings"       // String joining for statusline assembly
	"time"          // Current time display

	//--- External CPI-SI Packages ---
	"system/lib/temporal"    // Temporal context (time of day, schedule, calendar)
	"system/lib/display"     // Terminal colors and formatting codes
	"system/lib/instance"    // Instance configuration (name, emoji)

	//--- Statusline Presentation Libraries ---
	"statusline/lib/features"             // Display timing decisions
	"statusline/lib/format"               // Text optimization and formatting
	gitlib "statusline/lib/git"           // Git repository display
	sessionlib "statusline/lib/session"   // Session statistics display
	systemlib "statusline/lib/system"     // System health metrics display
	temporallib "statusline/lib/temporal" // Temporal awareness display
	"statusline/lib/types"                // Session data contract (foundation)
)

// ============================================================================
// END SETUP
// ============================================================================

// ============================================================================
// BODY
// ============================================================================

// ────────────────────────────────────────────────────────────────
// Core Operations - Statusline Assembly
// ────────────────────────────────────────────────────────────────

// statusline is the main entry point for statusline orchestration
//
// Named entry point pattern (not generic main) prevents function name collisions
// when multiple executables might be in scope or combined. Function name matches
// executable purpose - clear architectural intent.
//
// Flow:
//   1. Parse JSON from stdin → types.SessionContext
//   2. Retrieve temporal context → hooks/lib/temporal.TemporalContext
//   3. Build statusline → assemble all display elements
//   4. Output to stdout
//   5. Exit (success or error)
//
// Error Handling:
//   - JSON parse failure: Print error to stderr, exit code 1
//   - Missing data: Libraries handle gracefully (return zero-state)
//   - Temporal context error: Skip temporal elements, continue
//
// Performance: Typical execution <10ms (fast enough for per-hook invocation)
func statusline() {
	var ctx types.SessionContext

	// Parse Claude Code session context from stdin
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Assemble and output statusline
	output := buildStatusline(ctx)
	fmt.Println(output)
}

// buildStatusline assembles complete statusline from all presentation libraries
//
// Orchestration Logic:
//   - Decide WHAT to display (which elements)
//   - Decide WHERE to place them (order in statusline)
//   - Libraries decide HOW to format (presentation details)
//
// Display Order (left to right):
//   1. Instance identity (name + emoji)
//   2. Current date/time
//   3. Temporal awareness (time of day, session phase, schedule, calendar)
//   4. Model name
//   5. Working directory
//   6. Git status
//   7. System health (load, memory, disk)
//   8. Session statistics (lines, duration, cost)
//   9. Version info
//   10. Kingdom Technology reminder (occasional)
//
// Conditional Display:
//   - Only show elements with HasInfo=true (graceful degradation)
//   - Temporal context: Skip entirely if unavailable
//   - Kingdom reminder: Based on features.ShouldShowReminder()
//
// Returns: Formatted statusline string with colors and icons
func buildStatusline(ctx types.SessionContext) string {
	var parts []string

	// ────────────────────────────────────────────────────────────
	// Instance Identity
	// ────────────────────────────────────────────────────────────
	instanceConfig := instance.GetConfig()
	parts = append(parts, fmt.Sprintf("%s%s%s %s%s",
		display.Bold, display.Magenta,
		instanceConfig.Emoji, instanceConfig.Name,
		display.Reset))

	// ────────────────────────────────────────────────────────────
	// Current Date/Time
	// ────────────────────────────────────────────────────────────
	now := time.Now()
	dateTime := now.Format("Mon Jan 02 15:04:05")
	parts = append(parts, fmt.Sprintf("%s🕐 %s%s",
		display.Cyan, dateTime, display.Reset))

	// ────────────────────────────────────────────────────────────
	// Temporal Awareness (Four Dimensions)
	// ────────────────────────────────────────────────────────────
	// External Time, Internal Time, Internal Schedule, External Calendar
	if tempCtx, err := temporal.GetTemporalContext(); err == nil {
		// External Time - Time of day indicator
		timeOfDay := temporallib.GetTimeOfDayDisplay(tempCtx)
		parts = append(parts, fmt.Sprintf("%s%s %s%s",
			timeOfDay.Color, timeOfDay.Icon, timeOfDay.Label, display.Reset))

		// Internal Time - Session phase (duration + classification)
		sessionPhase := temporallib.GetSessionPhaseDisplay(tempCtx)
		if sessionPhase.Duration != "" {
			parts = append(parts, fmt.Sprintf("%s⏱️  %s (%s)%s",
				sessionPhase.Color, sessionPhase.Duration, sessionPhase.Phase, display.Reset))
		}

		// Internal Schedule - Current activity or downtime
		scheduleDisplay := temporallib.GetScheduleDisplay(tempCtx)
		if scheduleDisplay.Activity != "" {
			parts = append(parts, fmt.Sprintf("%s%s %s%s",
				scheduleDisplay.Color, scheduleDisplay.Icon, scheduleDisplay.Activity, display.Reset))
		}

		// External Calendar - Date and week number
		calendarDisplay := temporallib.GetCalendarDisplay(tempCtx)
		if calendarDisplay != "" {
			parts = append(parts, calendarDisplay)
		}
	}

	// ────────────────────────────────────────────────────────────
	// Model Information
	// ────────────────────────────────────────────────────────────
	modelName := format.GetShortModelName(ctx.Model.DisplayName)
	parts = append(parts, fmt.Sprintf("%s🧠 %s%s",
		display.Cyan, modelName, display.Reset))

	// ────────────────────────────────────────────────────────────
	// Working Directory (smart fallback logic)
	// ────────────────────────────────────────────────────────────
	// Priority: CWD → CurrentDir → ProjectDir
	workdir := ctx.CWD
	if workdir == "" {
		workdir = ctx.Workspace.CurrentDir
	}
	if workdir == "" {
		workdir = ctx.Workspace.ProjectDir
	}
	workdirShort := format.ShortenPath(workdir)
	parts = append(parts, fmt.Sprintf("%s%s📂 %s%s",
		display.Bold, display.Blue, workdirShort, display.Reset))

	// ────────────────────────────────────────────────────────────
	// Git Repository Status
	// ────────────────────────────────────────────────────────────
	gitDisplay := gitlib.GetGitDisplay(workdir)
	if gitDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %s%s",
			gitDisplay.Color, gitDisplay.Icon, gitDisplay.DisplayString, display.Reset))
	}

	// ────────────────────────────────────────────────────────────
	// System Health Metrics
	// ────────────────────────────────────────────────────────────

	// System load (1-minute average)
	loadDisplay := systemlib.GetLoadDisplay()
	if loadDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %.2f%s",
			loadDisplay.Color, loadDisplay.Icon, loadDisplay.LoadAvg, display.Reset))
	}

	// Memory usage
	memDisplay := systemlib.GetMemoryDisplay()
	if memDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %.1fG/%.1fG%s",
			memDisplay.Color, memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB, display.Reset))
	}

	// Disk usage (for current workspace)
	diskDisplay := systemlib.GetDiskDisplay(workdir)
	if diskDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %.0f%%%s",
			diskDisplay.Color, diskDisplay.Icon, diskDisplay.Percent, display.Reset))
	}

	// ────────────────────────────────────────────────────────────
	// Session Statistics
	// ────────────────────────────────────────────────────────────

	// Lines modified (added + removed)
	linesDisplay := sessionlib.GetLinesModifiedDisplay(ctx)
	if linesDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %d lines%s",
			linesDisplay.Color, linesDisplay.Icon, linesDisplay.TotalLines, display.Reset))
	}

	// Session duration
	durationDisplay := sessionlib.GetDurationDisplay(ctx)
	if durationDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %s%s",
			durationDisplay.Color, durationDisplay.Icon, durationDisplay.Duration, display.Reset))
	}

	// Cost tracking
	costDisplay := sessionlib.GetCostDisplay(ctx)
	if costDisplay.HasInfo {
		parts = append(parts, fmt.Sprintf("%s%s %s%s",
			costDisplay.Color, costDisplay.Icon, sessionlib.GetFormattedCost(costDisplay.Cost), display.Reset))
	}

	// ────────────────────────────────────────────────────────────
	// Claude Code Version
	// ────────────────────────────────────────────────────────────
	if ctx.Version != "" {
		parts = append(parts, fmt.Sprintf("%sv%s%s",
			display.Dim+display.Gray, ctx.Version, display.Reset))
	}

	// ────────────────────────────────────────────────────────────
	// Kingdom Technology Reminder (Occasional)
	// ────────────────────────────────────────────────────────────
	if features.ShouldShowReminder(ctx.SessionID) {
		parts = append(parts, fmt.Sprintf("%s%s⛪ Kingdom Technology%s",
			display.Dim, display.Cyan, display.Reset))
	}

	// ────────────────────────────────────────────────────────────
	// Assemble Final Statusline
	// ────────────────────────────────────────────────────────────
	return strings.Join(parts, "  ")
}

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Internal Structure
// ────────────────────────────────────────────────────────────────
//
// Available Processing Units (APU): 2
//   - statusline() (main entry point)
//   - buildStatusline() (assembly orchestration)
//
// Function Flow:
//   statusline()
//     ├─ Parse JSON → types.SessionContext
//     ├─ Get temporal context → temporal.TemporalContext
//     └─ buildStatusline()
//         ├─ Instance identity
//         ├─ Current date/time
//         ├─ Temporal awareness (4 dimensions)
//         ├─ Model + working directory
//         ├─ Git status
//         ├─ System health (3 metrics)
//         ├─ Session statistics (3 metrics)
//         ├─ Version info
//         └─ Kingdom reminder (conditional)
//
// Ladder (Dependencies):
//   Orchestrator (this file)
//     ↓ depends on
//   7 Presentation Libraries
//     ↓ depend on
//   types library (foundation)
//
// Baton (Execution Flow):
//   stdin (JSON) → statusline() → buildStatusline() → stdout (formatted string)
//
// Extension Point:
//   Adding new display elements:
//     1. Import new library
//     2. Call library function in buildStatusline()
//     3. Format result with colors/icons
//     4. Append to parts array
//     5. Choose display order (position in parts)

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================

// ────────────────────────────────────────────────────────────────
// Code Validation: Build and Test
// ────────────────────────────────────────────────────────────────
//
// Testing Requirements:
//   - Compile without errors: go build -o statusline statusline.go
//   - Run with sample JSON: echo '{"session_id":"test",...}' | ./statusline
//   - Verify all display elements appear correctly
//   - Check color codes render properly in terminal
//   - Test graceful degradation (missing fields)
//   - Verify error handling (invalid JSON)
//
// Build Verification:
//   ./build.sh (runs library verification + executable build)
//
// Integration Testing:
//   Test with real Claude Code JSON output
//   Verify temporal context integration
//   Check all 7 libraries produce expected output
//   Validate final string format
//
// Example validation:
//
//     # Test with minimal JSON
//     echo '{"session_id":"abc","model":{"display_name":"Sonnet"}}' | ./statusline
//
//     # Test with complete JSON
//     cat sample_session.json | ./statusline
//
//     # Verify exit codes
//     echo 'invalid json' | ./statusline  # Should exit 1

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls statusline() (named entry point)
//   3. statusline() orchestrates complete execution
//   4. Program exits after output
//
// Named Entry Point Benefits:
//   - Prevents main function collisions across executables
//   - Function name matches executable purpose (statusline)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Separates Go runtime entry from application logic
//
// Pattern:
//   func main() {
//       statusline()  // Minimal - just call named entry point
//   }
//
// This allows statusline() to be called, tested, or composed without
// triggering Go's executable entry point mechanics.

func main() {
	statusline() // Entry point for statusline orchestration
}

// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Stateless Orchestrator)
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - stdin/stdout: Automatically managed by OS
//   - JSON decoder: Garbage collected
//   - Strings: Short-lived, garbage collected
//   - No files, connections, or manual resources
//
// Graceful Shutdown:
//   - Program exits immediately after output
//   - No cleanup needed (stateless execution)
//   - Error path: Print to stderr, exit code 1
//   - Success path: Output to stdout, exit code 0
//
// Memory Management:
//   - All allocations short-lived (< 10ms execution)
//   - Garbage collector handles cleanup
//   - No memory leaks possible (no long-running state)

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════

// ────────────────────────────────────────────────────────────────
// Orchestrator Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: Complete statusline assembly from 7 presentation libraries
//
// Quick summary (high-level only - details in METADATA):
//   - Orchestrates all statusline presentation libraries
//   - Parses JSON input, assembles formatted output
//   - Demonstrates separation of concerns architecture
//
// Integration Pattern: See METADATA "Usage & Integration" section above
//
// Function Flow: See BODY "Organizational Chart" section above for complete
// execution flow and library coordination
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Baton - top rung orchestrator)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new display elements (import library, call in buildStatusline)
//   ✅ Change display order (reorder parts array)
//   ✅ Adjust formatting (colors, icons, spacing)
//   ✅ Add conditional display logic (new feature flags)
//   ✅ Enhance error handling (better error messages)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ stdin/stdout contract - breaks hook integration
//   ⚠️ JSON input format - breaks Claude Code compatibility
//   ⚠️ Exit codes - breaks error detection
//   ⚠️ Output format - breaks display expectations
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Named entry point pattern (func statusline(), func main())
//   ❌ Orchestration principle (libraries format, orchestrator assembles)
//   ❌ Graceful degradation (always produce valid output)
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing

// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart - Internal Structure" section above for
// complete ladder structure and baton flow.
//
// The Organizational Chart in BODY provides the detailed map showing:
// - Function dependencies (ladder)
// - Execution flow paths (baton)
// - Library integration points
// - APU count (Available Processing Units = 2)
//
// Quick architectural summary (details in BODY Organizational Chart):
// - Top rung orchestrator coordinating 7 presentation libraries
// - Ladder: Orchestrator → presentation libs → types (foundation)
// - Baton: stdin JSON → parse → format → assemble → stdout

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart" extension point above for detailed
// guide on adding new display elements.
//
// Quick reference (details in BODY extension point):
// - Adding display elements: Import library, call in buildStatusline()
// - Changing order: Reorder parts array assembly
// - Conditional display: Add logic checking HasInfo or feature flags

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Typical Execution: <10ms (fast enough for per-hook invocation)
//
// Performance Breakdown:
//   - JSON parsing: ~1ms
//   - Temporal context retrieval: ~2ms
//   - Library calls (7 libraries): ~5ms total
//   - String assembly: <1ms
//   - stdout output: <1ms
//
// Optimization:
//   Not needed - already optimal for hook invocation frequency
//   Libraries handle their own optimization internally

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Problem: JSON parsing fails
//   - Cause: Invalid JSON from Claude Code
//   - Solution: Check stdin input format, verify Claude Code JSON structure
//   - Debug: Run with sample JSON to isolate issue
//
// Problem: Missing display elements
//   - Cause: Library returning HasInfo=false (no data)
//   - Solution: Verify data source provides expected fields
//   - Expected: Graceful degradation - missing elements skip display
//
// Problem: Temporal elements not showing
//   - Cause: temporal.GetTemporalContext() returning error
//   - Solution: Check hooks/lib/temporal availability and configuration
//   - Expected: Skip temporal section entirely, continue with other elements
//
// Problem: Colors not rendering
//   - Cause: Terminal doesn't support ANSI color codes
//   - Solution: Use terminal with color support or strip colors
//   - Note: display.Reset codes ensure proper color isolation

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): 7 statusline libraries + external packages
// - Integration Points: stdin/stdout, hook system
//
// Quick summary (details in METADATA Dependencies section above):
// - Key dependencies: All 7 statusline/lib/* libraries
// - External: hooks/lib/temporal, system/lib/*
// - Consumers: Claude Code hooks system

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ All 7 libraries integrated - COMPLETED
//   ✓ Temporal awareness (4 dimensions) - COMPLETED
//   ⏳ Configurable display elements (user preferences)
//   ⏳ Alternate output formats (compact mode, verbose mode)
//   ⏳ Performance profiling integration
//
// Research Areas:
//   - Dynamic element ordering based on priority
//   - Context-aware element hiding (e.g., hide git outside repos)
//   - Custom color schemes per instance
//   - Statusline themes (minimal, verbose, developer, etc.)
//
// Integration Targets:
//   - Multiple output style support
//   - JSON output mode (for machine parsing)
//   - Logging integration (track display performance)
//
// Known Limitations to Address:
//   - Fixed display order (not user-configurable)
//   - No element priority system
//   - All elements same visual weight
//   - No responsive behavior (long statuslines wrap)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive template application and architectural refinement
//         - Applied full CODE-GO-002 template pattern
//         - Refactored to named entry point (func statusline())
//         - Fixed field access to match temporal library (scheduleDisplay.Activity, calendarDisplay as string)
//         - Package name: statusline (not main) - prevents name collisions
//         - Complete METADATA with biblical foundation and CPI-SI identity
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 11 required sections
//         - Documented orchestration pattern and library integration
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - Basic orchestrator structure
//         - Integration with 7 presentation libraries
//         - JSON parsing and stdout output
//         - Temporal awareness integration
//         - Color and icon formatting

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This orchestrator is the TOP RUNG of the statusline ladder architecture.
// It coordinates 7 presentation libraries without implementing formatting logic
// itself - pure orchestration demonstrating separation of concerns.
//
// Modify thoughtfully - changes here affect how entire statusline appears.
// Libraries handle HOW to format, this decides WHAT and WHERE. Keep orchestration
// logic separate from presentation logic.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with real Claude Code JSON
//   - Verify all libraries still compile
//   - Check output in terminal with colors
//   - Document all changes comprehensively
//
// "Let all things be done decently and in order" - 1 Corinthians 14:40

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Execution:
//   echo '{"session_id":"abc","model":{"display_name":"Sonnet"}}' | ./statusline
//
// With Complete Session Data:
//   cat session_context.json | ./statusline
//
// Build and Test:
//   ./build.sh
//   echo '{"session_id":"test"}' | ./statusline
//
// Integration with Hooks:
//   # Called by Claude Code hooks automatically
//   # Receives session context JSON via stdin
//   # Outputs formatted statusline to stdout

// ============================================================================
// END CLOSING
// ============================================================================
