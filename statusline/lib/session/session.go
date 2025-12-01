// METADATA
//
// Session Statistics Display Library - CPI-SI Statusline
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "So teach us to number our days, that we may apply our hearts unto wisdom" - Psalm 90:12
// Principle: Awareness of time and stewardship - tracking work to understand patterns and use time wisely
// Anchor: "Redeeming the time, because the days are evil" - Ephesians 5:16
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - Middle Rung)
// Role: Presentation layer for session statistics in statusline display
// Paradigm: CPI-SI framework component - formats session data for user visibility
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise (CreativeWorkzStudio LLC)
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2025-10-24
// Version: 1.0.0
// Last Modified: 2025-11-04 - Applied comprehensive GO library template
//
// Version History:
//   1.0.0 (2025-11-04) - Applied full template, comprehensive documentation
//   0.2.0 (2025-11-04) - Updated to use system/lib/sessiontime.FormatDuration
//   0.1.0 (2025-10-24) - Initial implementation with session statistics formatting
//
// Purpose & Function
//
// Purpose: Transform session statistics (lines modified, duration, cost) into formatted
// display strings suitable for statusline space constraints
//
// Core Design: Pure presentation layer - receives session data from SessionContext,
// outputs formatted display structures with visual indicators (colors, icons)
//
// Key Features:
//   - Lines modified display (added + removed)
//   - Session duration formatting (human-readable time)
//   - Cost tracking display (USD with 4 decimal precision)
//   - Zero-state handling (gracefully omits unavailable data)
//   - Visual clarity through color coding and icons
//
// Philosophy: Show session progress without overwhelming - what changed, how long,
// what it cost. Data collection happens elsewhere, presentation happens here.
//
// Blocking Status
//
// Non-blocking: Returns display structures immediately, never blocks statusline
// Mitigation: No I/O operations, pure formatting of in-memory data
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/session"
//
// Integration Pattern:
//   1. Pass SessionContext to display functions
//   2. Check HasInfo field to determine if data available
//   3. Use display fields (TotalLines, Duration, Cost) for statusline formatting
//   4. Apply Color and Icon for visual styling
//   5. No cleanup needed - stateless library
//
// Public API (in typical usage order):
//
//   Session Statistics (presentation):
//     GetLinesModifiedDisplay(ctx SessionContext) LinesModifiedDisplay
//     GetDurationDisplay(ctx SessionContext) DurationDisplay
//     GetCostDisplay(ctx SessionContext) CostDisplay
//     GetFormattedCost(cost float64) string
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, time
//   External: None
//   Internal: system/lib/display (color constants), system/lib/sessiontime (duration formatting), statusline/lib/types (SessionContext)
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Ladder: Uses system/lib/sessiontime for duration formatting (lower rung)
//   - Baton: Receives SessionContext, returns display structures
//   - Rails: N/A (pure function library, no logging infrastructure)
//
// Health Scoring
//
// Pure presentation library - no health scoring infrastructure needed.
// Functions are guaranteed to succeed (graceful degradation on missing data).
//
// Display Generation:
//   - Lines modified: Always succeeds (zero-state if no changes)
//   - Duration formatting: Always succeeds (zero-state if no duration)
//   - Cost display: Always succeeds (zero-state if no cost)
//
// Note: This library cannot fail - all operations return valid display structures.
// Health tracking would measure "successfully did nothing" which provides no value.
package session

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

	"fmt"  // String formatting for cost display
	"time" // Duration types for time conversion

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display"     // Color constants for terminal output
	"system/lib/sessiontime" // Duration formatting (human-readable time)
	"statusline/lib/types"   // SessionContext struct (data source)
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Display Types ---
// Presentation structures for statusline rendering.

// LinesModifiedDisplay represents formatted lines modified information for statusline display.
//
// Contains total lines changed (added + removed), visual styling (color/icon), and
// availability flag. Used to show code editing activity in session.
//
// Zero value represents "no changes" state (HasInfo: false, TotalLines: 0).
//
// Example usage:
//
//     linesDisplay := session.GetLinesModifiedDisplay(ctx)
//     if linesDisplay.HasInfo {
//         fmt.Printf("%s %d lines\n", linesDisplay.Icon, linesDisplay.TotalLines)
//     }
type LinesModifiedDisplay struct {
	TotalLines int    // Total lines added + removed
	Color      string // Terminal color code for display
	Icon       string // Visual icon representing editing (e.g., "📝")
	HasInfo    bool   // True if changes tracked, false if no data
}

// DurationDisplay represents formatted session duration for statusline display.
//
// Contains human-readable duration string (e.g., "1h 23m"), visual styling, and
// availability flag. Used to show how long session has been active.
//
// Zero value represents "no duration" state (HasInfo: false, empty Duration).
//
// Example usage:
//
//     durationDisplay := session.GetDurationDisplay(ctx)
//     if durationDisplay.HasInfo {
//         fmt.Printf("%s %s\n", durationDisplay.Icon, durationDisplay.Duration)
//     }
type DurationDisplay struct {
	Duration string // Human-readable duration (e.g., "1h 23m", "45s")
	Color    string // Terminal color code for display
	Icon     string // Visual icon representing time (e.g., "⏱")
	HasInfo  bool   // True if duration tracked, false if no data
}

// CostDisplay represents formatted cost tracking for statusline display.
//
// Contains USD cost value, visual styling, and availability flag. Used to show
// session cost (API usage, token costs, etc.).
//
// Zero value represents "no cost" state (HasInfo: false, Cost: 0).
//
// Example usage:
//
//     costDisplay := session.GetCostDisplay(ctx)
//     if costDisplay.HasInfo {
//         fmt.Printf("%s %s\n", costDisplay.Icon, session.GetFormattedCost(costDisplay.Cost))
//     }
type CostDisplay struct {
	Cost    float64 // Cost in USD
	Color   string  // Terminal color code for display
	Icon    string  // Visual icon representing cost (e.g., "💰")
	HasInfo bool    // True if cost tracked, false if no data
}

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
//   ├── GetLinesModifiedDisplay() → uses SessionContext.Cost fields
//   ├── GetDurationDisplay() → uses SessionContext.Cost + sessiontime.FormatDuration()
//   ├── GetCostDisplay() → uses SessionContext.Cost fields
//   └── GetFormattedCost() → uses fmt.Sprintf() [helper]
//
//   Helpers: GetFormattedCost() (pure formatting helper)
//
// Baton Flow (Execution Paths):
//
//   Entry → GetLinesModifiedDisplay(ctx)
//     ↓
//   Calculate total (added + removed)
//     ↓
//   Assign color/icon
//     ↓
//   Exit → return LinesModifiedDisplay
//
//   Entry → GetDurationDisplay(ctx)
//     ↓
//   system/lib/sessiontime.FormatDuration() [external]
//     ↓
//   Assign color/icon
//     ↓
//   Exit → return DurationDisplay
//
//   Entry → GetCostDisplay(ctx)
//     ↓
//   Extract cost value
//     ↓
//   Assign color/icon
//     ↓
//   Exit → return CostDisplay
//
//   Entry → GetFormattedCost(cost)
//     ↓
//   fmt.Sprintf("$%.4f", cost)
//     ↓
//   Exit → return string
//
// APUs (Available Processing Units):
// - 4 functions total
// - 1 helper (GetFormattedCost)
// - 0 core operations (all presentation)
// - 3 public APIs (exported display functions)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Usually not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// GetFormattedCost formats cost value as USD string with 4 decimal precision.
//
// What It Does:
// Converts float64 cost to standardized string format "$0.0000" for consistent
// display across statusline. Four decimal places provide precision for token costs.
//
// Parameters:
//   cost: Float64 value representing USD cost
//
// Returns:
//   string: Formatted cost string (e.g., "$0.0123", "$1.5000")
//
// Example usage:
//
//	costStr := GetFormattedCost(0.0123)  // → "$0.0123"
//	costStr := GetFormattedCost(1.5)     // → "$1.5000"
func GetFormattedCost(cost float64) string {
	return fmt.Sprintf("$%.4f", cost)
}

// ────────────────────────────────────────────────────────────────
// Display Formatting - Presentation Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Transform session statistics (from SessionContext) into formatted display
// structures with visual indicators suitable for statusline constraints.
//
// Why Separated:
// Three distinct types of session data (lines, duration, cost) each need
// independent formatting. Separation allows statusline orchestrator to choose
// which statistics to display based on available space and user preferences.
//
// Extension Point:
// To add new session statistics (commits, files changed, etc.):
//   1. Define new display structure type in SETUP (follow existing pattern)
//   2. Create Get[Statistic]Display() function following naming pattern
//   3. Extract data from SessionContext (or add new field if needed)
//   4. Apply color/icon selection for visual consistency
//   5. Update API documentation with new statistic display
//   6. Add tests for new display function

// GetLinesModifiedDisplay returns formatted lines modified information for statusline display.
//
// What It Does:
// Calculates total lines changed (added + removed) from SessionContext and formats
// into display structure with visual indicators. Shows code editing activity.
//
// Parameters:
//   ctx: SessionContext containing Cost.TotalLinesAdded and Cost.TotalLinesRemoved
//
// Returns:
//   LinesModifiedDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   153 lines → Yellow 📝 (153 total changes)
//   0 lines   → HasInfo: false (no changes to display)
//
// Behavior:
//   - No changes (0 lines): Returns HasInfo: false
//   - Any changes: Returns TotalLines with yellow color and edit icon
//
// Example usage:
//
//	linesDisplay := GetLinesModifiedDisplay(ctx)
//	if linesDisplay.HasInfo {
//	    fmt.Printf("%s %d lines modified\n", linesDisplay.Icon, linesDisplay.TotalLines)
//	}
func GetLinesModifiedDisplay(ctx types.SessionContext) LinesModifiedDisplay {
	// Calculate total lines changed (added + removed)
	totalLinesModified := ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved

	// No changes - return zero-state
	if totalLinesModified <= 0 {
		return LinesModifiedDisplay{HasInfo: false}
	}

	// Return formatted display with visual elements
	return LinesModifiedDisplay{
		TotalLines: totalLinesModified,
		Color:      display.Yellow, // Yellow indicates active editing
		Icon:       "📝",           // Edit icon represents code changes
		HasInfo:    true,
	}
}

// GetDurationDisplay returns formatted session duration for statusline display.
//
// What It Does:
// Converts session duration from milliseconds to human-readable format (e.g., "1h 23m")
// using system/lib/sessiontime and formats into display structure with visual indicators.
//
// Parameters:
//   ctx: SessionContext containing Cost.TotalDurationMS (duration in milliseconds)
//
// Returns:
//   DurationDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   "45s"      → Gray ⏱ (45 seconds)
//   "1h 23m"   → Gray ⏱ (1 hour 23 minutes)
//   ""         → HasInfo: false (no duration tracked)
//
// Behavior:
//   - No duration (0 ms): Returns HasInfo: false
//   - Any duration: Formats using sessiontime.FormatDuration() with gray color
//
// Example usage:
//
//	durationDisplay := GetDurationDisplay(ctx)
//	if durationDisplay.HasInfo {
//	    fmt.Printf("%s %s elapsed\n", durationDisplay.Icon, durationDisplay.Duration)
//	}
func GetDurationDisplay(ctx types.SessionContext) DurationDisplay {
	// No duration tracked - return zero-state
	if ctx.Cost.TotalDurationMS <= 0 {
		return DurationDisplay{HasInfo: false}
	}

	// Convert milliseconds to time.Duration and format for humans
	duration := sessiontime.FormatDuration(time.Duration(ctx.Cost.TotalDurationMS) * time.Millisecond)

	// Return formatted display with visual elements
	return DurationDisplay{
		Duration: duration,
		Color:    display.Gray, // Gray indicates passive time tracking
		Icon:     "⏱",          // Timer icon represents duration
		HasInfo:  true,
	}
}

// GetCostDisplay returns formatted cost tracking for statusline display.
//
// What It Does:
// Extracts cost value from SessionContext and formats into display structure
// with visual indicators. Shows session cost (API usage, token costs, etc.).
//
// Parameters:
//   ctx: SessionContext containing Cost.TotalCostUSD (cost in USD)
//
// Returns:
//   CostDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   $0.0123 → Yellow 💰 (cost value ready for formatting)
//   $0.0000 → HasInfo: false (no cost to display)
//
// Behavior:
//   - No cost (0 USD): Returns HasInfo: false
//   - Any cost: Returns Cost value with yellow color and money icon
//
// Note: Use GetFormattedCost(costDisplay.Cost) to convert to string for display
//
// Example usage:
//
//	costDisplay := GetCostDisplay(ctx)
//	if costDisplay.HasInfo {
//	    fmt.Printf("%s %s\n", costDisplay.Icon, GetFormattedCost(costDisplay.Cost))
//	}
func GetCostDisplay(ctx types.SessionContext) CostDisplay {
	// No cost tracked - return zero-state
	if ctx.Cost.TotalCostUSD <= 0 {
		return CostDisplay{HasInfo: false}
	}

	// Return formatted display with visual elements
	return CostDisplay{
		Cost:    ctx.Cost.TotalCostUSD,
		Color:   display.Yellow, // Yellow indicates cost awareness
		Icon:    "💰",           // Money icon represents cost
		HasInfo: true,
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
//
// ────────────────────────────────────────────────────────────────
// Code Validation: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Validation section explanation, see: standards/code/4-block/sections/CWS-SECTION-010-CLOSING-code-validation.md
//
// Testing Requirements:
//   - Import the library without errors
//   - Call each display function with SessionContext containing data
//   - Call each display function with SessionContext containing zero values
//   - Verify HasInfo correctly distinguishes data vs no-data states
//   - Check display fields populated correctly (TotalLines, Duration, Cost)
//   - Confirm Color and Icon fields set appropriately
//   - Test GetFormattedCost() with various cost values
//   - Ensure no panics or errors for edge cases (negative values, very large numbers)
//   - Run: go vet ./... (no warnings)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - gofmt -l . (no formatting issues)
//
// Integration Testing:
//   - Test with actual statusline orchestrator
//   - Verify display in terminal with different session states
//   - Check visual indicators render correctly
//   - Validate formatting of duration strings (45s, 1h 23m, etc.)
//   - Verify cost precision (4 decimal places)
//
// Example validation code:
//
//     // Test with session data
//     ctx := types.SessionContext{}
//     ctx.Cost.TotalLinesAdded = 100
//     ctx.Cost.TotalLinesRemoved = 50
//     ctx.Cost.TotalDurationMS = 5_400_000  // 1h 30m
//     ctx.Cost.TotalCostUSD = 0.0123
//
//     linesDisplay := session.GetLinesModifiedDisplay(ctx)
//     if !linesDisplay.HasInfo || linesDisplay.TotalLines != 150 {
//         t.Error("Expected 150 total lines")
//     }
//
//     durationDisplay := session.GetDurationDisplay(ctx)
//     if !durationDisplay.HasInfo || durationDisplay.Duration == "" {
//         t.Error("Expected non-empty duration")
//     }
//
//     costDisplay := session.GetCostDisplay(ctx)
//     if !costDisplay.HasInfo || costDisplay.Cost != 0.0123 {
//         t.Error("Expected cost 0.0123")
//     }
//
//     // Test zero-state
//     emptyCtx := types.SessionContext{}
//     if session.GetLinesModifiedDisplay(emptyCtx).HasInfo {
//         t.Error("Expected HasInfo false for empty context")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by other components.
//
// Usage: import "statusline/lib/session"
//
// The library is imported into the calling package (typically statusline orchestrator),
// making all display functions available. No code executes during import - functions
// are defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "statusline/lib/session"
//     import "statusline/lib/types"
//
//     func main() {
//         ctx := types.SessionContext{ /* ... populated ... */ }
//
//         // Get session statistics displays
//         linesDisplay := session.GetLinesModifiedDisplay(ctx)
//         durationDisplay := session.GetDurationDisplay(ctx)
//         costDisplay := session.GetCostDisplay(ctx)
//
//         // Build statusline with available data
//         var parts []string
//         if linesDisplay.HasInfo {
//             parts = append(parts, fmt.Sprintf("%d lines", linesDisplay.TotalLines))
//         }
//         if durationDisplay.HasInfo {
//             parts = append(parts, durationDisplay.Duration)
//         }
//         if costDisplay.HasInfo {
//             parts = append(parts, session.GetFormattedCost(costDisplay.Cost))
//         }
//
//         fmt.Println(strings.Join(parts, " | "))
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Memory: Three display struct allocations per orchestrator call (~150 bytes total)
//   - No file handles, network connections, or persistent resources
//   - Go's garbage collector handles memory automatically
//
// Graceful Shutdown:
//   - N/A for stateless library (no lifecycle)
//   - Calling code responsible for any display cleanup
//   - No state to persist or restore
//
// Error State Cleanup:
//   - No error states possible - all operations guaranteed to succeed
//   - Graceful degradation returns valid zero-state (HasInfo: false)
//   - No partial state or corruption possible
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Small allocations (~50 bytes per display struct) don't require manual management
//   - No large buffers or long-lived allocations
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
//   - Transforms session statistics into formatted statusline display strings
//   - Shows lines modified, session duration, and cost tracking with visual indicators
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Ladder - Middle Rung presentation layer)
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new session statistic displays (commits, files, etc.)
//   ✅ Adjust color or icon choices for different statistics
//   ✅ Extend display structures with additional fields
//   ✅ Add helper functions for complex formatting logic
//   ✅ Create alternative formatting functions (compact, verbose, etc.)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Display struct fields - breaks statusline orchestrator
//   ⚠️ Function signatures - breaks all calling code
//   ⚠️ HasInfo semantics - breaks zero-state handling
//   ⚠️ Color/Icon field types - breaks display rendering
//   ⚠️ GetFormattedCost() format - breaks cost display expectations
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Pure function guarantee (stateless, no side effects)
//   ❌ Graceful degradation (always return valid display)
//   ❌ Data vs Presentation separation (SessionContext provides data)
//   ❌ Non-blocking guarantee (no I/O operations)
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
// Quick architectural summary (details in BODY Organizational Chart):
// - 3 public APIs (GetLinesModifiedDisplay, GetDurationDisplay, GetCostDisplay)
// - 1 helper (GetFormattedCost)
// - Ladder: Uses system/lib/sessiontime (lower rung) for duration formatting
// - Baton: SessionContext → formatting → display structures
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// Adding New Session Statistics:
//   1. Define new display structure in SETUP Types section
//   2. Create Get[Statistic]Display() function in BODY
//   3. Extract data from SessionContext (or extend SessionContext if needed)
//   4. Follow existing patterns: check for zero-state, assign color/icon
//   5. Update API documentation with new statistic
//   6. Add tests for new display function
//
// Example pattern for adding commits display:
//
//     // In SETUP Types:
//     type CommitsDisplay struct {
//         CommitCount int
//         Color       string
//         Icon        string
//         HasInfo     bool
//     }
//
//     // In BODY:
//     func GetCommitsDisplay(ctx types.SessionContext) CommitsDisplay {
//         if ctx.Git.CommitCount <= 0 {
//             return CommitsDisplay{HasInfo: false}
//         }
//
//         return CommitsDisplay{
//             CommitCount: ctx.Git.CommitCount,
//             Color:       display.Green,
//             Icon:        "📦",
//             HasInfo:     true,
//         }
//     }
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Types: Display structs (~50 bytes each)
//
// Quick summary:
// - GetLinesModifiedDisplay(): O(1) arithmetic, ~50 byte allocation, <1 microsecond
// - GetDurationDisplay(): O(1) formatting, ~50 byte allocation, <10 microseconds
// - GetCostDisplay(): O(1) extraction, ~50 byte allocation, <1 microsecond
// - GetFormattedCost(): O(1) sprintf, ~20 byte allocation, <1 microsecond
// - Memory: Three struct allocations per statusline render (~150 bytes total)
// - No I/O operations (all in-memory formatting)
//
// Key optimization: This library needs no optimization. Struct allocation and
// string formatting are already optimal for this use case.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// This library has no common failure modes - all operations guaranteed to succeed.
//
// Expected Behaviors:
//   - Zero lines modified returns HasInfo: false - This is correct, not an error
//   - Zero duration returns HasInfo: false - This is correct, not an error
//   - Zero cost returns HasInfo: false - This is correct, not an error
//   - Empty duration string when HasInfo: false - This is correct zero-state
//
// If display functions return unexpected results:
//   Problem: HasInfo true but display data seems wrong
//     - Cause: SessionContext populated incorrectly
//     - Solution: Check where SessionContext gets populated
//     - Verify Cost.TotalLinesAdded, Cost.TotalLinesRemoved, Cost.TotalDurationMS, Cost.TotalCostUSD
//
//   Problem: Duration format looks wrong
//     - Cause: system/lib/sessiontime.FormatDuration() behavior changed
//     - Solution: Check system/lib/sessiontime implementation
//     - Expected formats: "45s", "1m 23s", "1h 23m", "2h 15m 30s"
//
//   Problem: Cost precision incorrect (not 4 decimals)
//     - Cause: GetFormattedCost() format string modified
//     - Solution: Verify fmt.Sprintf("$%.4f", cost) unchanged
//     - Expected: Always 4 decimal places ($0.0000, $0.0123, $1.5000)
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependencies: system/lib/sessiontime (duration formatting), statusline/lib/types (SessionContext)
// - Primary consumer: statusline orchestrator (main command)
//
// Parallel Implementation:
//   - Go version: This file (statusline/lib/session/session.go)
//   - Shared philosophy: Data collection elsewhere, presentation here
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Lines modified display - COMPLETED
//   ✓ Session duration display - COMPLETED
//   ✓ Cost tracking display - COMPLETED
//   ⏳ Commit count display
//   ⏳ Files changed display
//   ⏳ Token usage display (input/output breakdown)
//   ⏳ Custom formatting options (compact, verbose, detailed)
//
// Research Areas:
//   - Color coding based on thresholds (high cost = red, normal = yellow)
//   - Configurable cost precision (2 vs 4 decimal places)
//   - Human-readable large numbers (1.5K lines instead of 1500)
//   - Session comparison (current vs previous session stats)
//
// Known Limitations to Address:
//   - No breakdown of lines added vs removed (just total)
//   - Duration loses seconds for long sessions (shows "1h 23m", not "1h 23m 45s")
//   - Cost always 4 decimals even for large values ($100.0000 looks odd)
//   - No configuration options (colors, icons, formats all hardcoded)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//         - Applied full GO library template (114 lines → comprehensive docs)
//         - Added complete METADATA section (Biblical foundation, CPI-SI identity)
//         - Expanded SETUP with detailed type documentation
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 13 sections
//         - Established architectural pattern: data vs presentation separation
//
//   0.2.0 (2025-11-04) - Duration formatting refactor
//         - Updated GetDurationDisplay() to use system/lib/sessiontime.FormatDuration()
//         - Removed duplicate duration formatting logic
//         - Improved consistency with system-wide duration display
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - Basic session statistics formatting (lines, duration, cost)
//         - Three display structures (LinesModifiedDisplay, DurationDisplay, CostDisplay)
//         - GetFormattedCost() helper for USD formatting
//         - Graceful degradation for missing data
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component (middle rung) - provides presentation
// layer for session statistics. Receives data from SessionContext (populated
// elsewhere), formats for statusline display. Part of the architectural pattern
// where data collection happens upstream, presentation happens here.
//
// Modify thoughtfully - changes here affect statusline display. The separation
// between data collection (SessionContext population) and presentation (this library)
// must be maintained. Never duplicate data collection logic here - always receive
// data via SessionContext.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with various session states before committing
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain data vs presentation separation
//
// "So teach us to number our days, that we may apply our hearts unto wisdom" - Psalm 90:12
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "statusline/lib/session"
//     import "statusline/lib/types"
//
//     ctx := types.SessionContext{ /* ... populated ... */ }
//
//     // Get session statistics displays
//     linesDisplay := session.GetLinesModifiedDisplay(ctx)
//     durationDisplay := session.GetDurationDisplay(ctx)
//     costDisplay := session.GetCostDisplay(ctx)
//
// Check if Data Available:
//
//     if linesDisplay.HasInfo {
//         fmt.Printf("%s %d lines\n", linesDisplay.Icon, linesDisplay.TotalLines)
//     }
//
//     if durationDisplay.HasInfo {
//         fmt.Printf("%s %s\n", durationDisplay.Icon, durationDisplay.Duration)
//     }
//
//     if costDisplay.HasInfo {
//         fmt.Printf("%s %s\n", costDisplay.Icon, session.GetFormattedCost(costDisplay.Cost))
//     }
//
// Full Statusline Integration:
//
//     func buildStatusline(ctx types.SessionContext) string {
//         var parts []string
//
//         // Session statistics
//         linesDisplay := session.GetLinesModifiedDisplay(ctx)
//         if linesDisplay.HasInfo {
//             parts = append(parts, fmt.Sprintf("%d lines", linesDisplay.TotalLines))
//         }
//
//         durationDisplay := session.GetDurationDisplay(ctx)
//         if durationDisplay.HasInfo {
//             parts = append(parts, durationDisplay.Duration)
//         }
//
//         costDisplay := session.GetCostDisplay(ctx)
//         if costDisplay.HasInfo {
//             parts = append(parts, session.GetFormattedCost(costDisplay.Cost))
//         }
//
//         // ... add other statusline parts ...
//
//         return strings.Join(parts, " | ")
//         // → "153 lines | 1h 23m | $0.0123"
//     }
//
// Handling Different Session States:
//
//     // New session (no activity yet)
//     emptyCtx := types.SessionContext{}
//     linesDisplay := session.GetLinesModifiedDisplay(emptyCtx)
//     // → HasInfo: false (nothing to display)
//
//     // Active session with changes
//     activeCtx := types.SessionContext{
//         Cost: types.CostInfo{
//             TotalLinesAdded: 100,
//             TotalLinesRemoved: 50,
//             TotalDurationMS: 5400000,  // 1h 30m
//             TotalCostUSD: 0.0123,
//         },
//     }
//     linesDisplay = session.GetLinesModifiedDisplay(activeCtx)
//     // → TotalLines: 150, HasInfo: true
//
// ============================================================================
// END CLOSING
// ============================================================================
