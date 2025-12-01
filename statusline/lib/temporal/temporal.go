// METADATA
//
// Temporal Display Library - CPI-SI Statusline
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1
// Principle: Wisdom in timing - knowing what time it is, what phase we're in, what season of work
// Anchor: "Redeeming the time, because the days are evil" - Ephesians 5:16
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - Middle Rung)
// Role: Presentation layer for temporal awareness in statusline display
// Paradigm: CPI-SI framework component - formats time/schedule data for user visibility
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
//   0.1.0 (2025-10-24) - Initial implementation with temporal display functions
//
// Purpose & Function
//
// Purpose: Transform temporal awareness data (time of day, session phase, schedule status, calendar)
// into formatted display structures with contextual icons and colors suitable for statusline
//
// Core Design: Pure presentation layer - receives temporal data from hooks/lib/temporal,
// outputs formatted display structures with visual indicators (icons, colors, labels)
//
// Key Features:
//   - Time of day display with contextual icons (🌅 morning, ☀️ afternoon, 🌆 evening, 🌙 night)
//   - Session phase display with duration and color coding
//   - Schedule/downtime status with activity indicators
//   - Calendar display with day of week and week number
//   - Zero-state handling (gracefully omits unavailable data)
//
// Philosophy: Show temporal context without overwhelming - what time it is, how long working,
// what should be happening. Data collection happens in hooks/lib/temporal, presentation here.
//
// Blocking Status
//
// Non-blocking: Returns display structures immediately, never blocks statusline
// Mitigation: Uses hooks/lib/temporal which provides in-memory temporal context
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/temporal"
//	import "hooks/lib/temporal" // For TemporalContext data
//
// Integration Pattern:
//   1. Get TemporalContext from hooks/lib/temporal
//   2. Pass context to display functions
//   3. Check return values for empty/zero states
//   4. Use display fields for statusline formatting
//   5. Apply Icon and Color for visual styling
//   6. No cleanup needed - stateless library
//
// Public API (in typical usage order):
//
//   Temporal Display (presentation):
//     GetTimeOfDayDisplay(tempCtx *temporal.TemporalContext) TimeOfDayDisplay
//     GetSessionPhaseDisplay(tempCtx *temporal.TemporalContext) SessionPhaseDisplay
//     GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay
//     GetCalendarDisplay(tempCtx *temporal.TemporalContext) string
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt
//   External: None
//   Internal: hooks/lib/temporal (temporal context data), system/lib/display (color constants)
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Ladder: Uses hooks/lib/temporal for data collection (lower rung)
//   - Baton: Receives TemporalContext, returns display structures
//   - Rails: N/A (pure function library, no logging infrastructure)
//
// Health Scoring
//
// Pure presentation library - no health scoring infrastructure needed.
// Functions are guaranteed to succeed (graceful degradation on missing data).
//
// Display Generation:
//   - Time of day: Always succeeds (returns default icon if data missing)
//   - Session phase: Returns empty display if no duration tracked
//   - Schedule: Returns empty display if no schedule data
//   - Calendar: Returns empty string if no calendar data
//
// Note: This library cannot fail - all operations return valid display structures or empty values.
package temporal

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

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"fmt" // String formatting for calendar display

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/temporal"  // Temporal context data (time, schedule, calendar)
	"system/lib/display"   // Color constants for terminal output
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────

//--- Display Types ---
// Presentation structures for statusline rendering.

// TimeOfDayDisplay represents formatted time of day with contextual icon.
//
// Contains time of day label (morning/afternoon/evening/night), icon representing
// that time period, and color for display. Used to show current time context.
//
// Zero value represents default state (empty values).
//
// Example usage:
//
//     todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
//     fmt.Printf("%s %s\n", todDisplay.Icon, todDisplay.Label)
//     // → "🌅 morning"
type TimeOfDayDisplay struct {
	Icon  string // Contextual icon (🌅 morning, ☀️ afternoon, 🌆 evening, 🌙 night)
	Label string // Time of day label ("morning", "afternoon", "evening", "night")
	Color string // Terminal color code for display
}

// SessionPhaseDisplay represents formatted session phase with duration.
//
// Contains session duration (human-readable), phase label (fresh/active/long),
// and color coding based on phase. Used to show how long session has been running.
//
// Zero value represents "no session data" state (empty values).
//
// Example usage:
//
//     phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
//     if phaseDisplay.Duration != "" {
//         fmt.Printf("%s (%s)\n", phaseDisplay.Duration, phaseDisplay.Phase)
//         // → "1h23m (active)"
//     }
type SessionPhaseDisplay struct {
	Duration string // Human-readable duration (e.g., "1h23m", "45s")
	Phase    string // Session phase ("fresh", "active", "long")
	Color    string // Terminal color code (green for fresh/active, yellow for long)
}

// ScheduleDisplay represents formatted schedule/downtime status.
//
// Contains current activity, color coding, and icon. Used to show whether in
// scheduled work window or expected downtime.
//
// Zero value represents "no schedule data" state (empty values).
//
// Example usage:
//
//     schedDisplay := temporal.GetScheduleDisplay(tempCtx)
//     if schedDisplay.Activity != "" {
//         fmt.Printf("%s %s\n", schedDisplay.Icon, schedDisplay.Activity)
//         // → "📋 Deep Work" or "⏸️ Downtime"
//     }
type ScheduleDisplay struct {
	Activity string // Current activity name or "Downtime"
	Color    string // Terminal color code (green for work, yellow for downtime)
	Icon     string // Activity icon (📋 for work, ⏸️ for downtime)
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
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Orchestration)
//   ├── GetTimeOfDayDisplay(tempCtx) → icon mapping + color
//   ├── GetSessionPhaseDisplay(tempCtx) → duration + phase color
//   ├── GetScheduleDisplay(tempCtx) → activity/downtime formatting
//   └── GetCalendarDisplay(tempCtx) → calendar string formatting
//
//   Icon/Color Mappings (Presentation Logic):
//   ├── Time of Day: morning=🌅, afternoon=☀️, evening=🌆, night=🌙
//   ├── Session Phase: fresh/active=green, long=yellow
//   └── Schedule: work=📋+green, downtime=⏸️+yellow
//
// Baton Flow (Execution Paths):
//
//   Entry → GetTimeOfDayDisplay(tempCtx)
//     ↓
//   Extract tempCtx.ExternalTime.TimeOfDay
//     ↓
//   Map to icon (switch statement)
//     ↓
//   Exit → return TimeOfDayDisplay
//
//   Entry → GetSessionPhaseDisplay(tempCtx)
//     ↓
//   Check tempCtx.InternalTime.ElapsedFormatted
//     ↓
//   Determine phase color
//     ↓
//   Exit → return SessionPhaseDisplay (or empty if no duration)
//
//   Entry → GetScheduleDisplay(tempCtx)
//     ↓
//   Check InWorkWindow or ExpectedDowntime
//     ↓
//   Format activity with icon/color
//     ↓
//   Exit → return ScheduleDisplay (or empty if no schedule)
//
//   Entry → GetCalendarDisplay(tempCtx)
//     ↓
//   Check tempCtx.ExternalCalendar.DayOfWeek
//     ↓
//   Format calendar string
//     ↓
//   Exit → return string (or empty if no calendar)
//
// APUs (Available Processing Units):
// - 4 functions total
// - 0 helpers (all presentation logic inline)
// - 4 public APIs (exported display functions)

// ────────────────────────────────────────────────────────────────
// Temporal Display Functions - Presentation Logic
// ────────────────────────────────────────────────────────────────

// GetTimeOfDayDisplay returns formatted time of day with contextual icon.
//
// What It Does:
// Maps time of day label (morning/afternoon/evening/night) to contextual icon
// for visual time awareness in statusline.
//
// Parameters:
//   tempCtx: TemporalContext from hooks/lib/temporal containing ExternalTime
//
// Returns:
//   TimeOfDayDisplay: Icon, label, and color for statusline display
//
// Icon Mapping:
//   morning   → 🌅 (sunrise)
//   afternoon → ☀️ (sun)
//   evening   → 🌆 (cityscape at dusk)
//   night     → 🌙 (crescent moon)
//
// Example usage:
//
//	todDisplay := GetTimeOfDayDisplay(tempCtx)
//	fmt.Printf("%s %s\n", todDisplay.Icon, todDisplay.Label)
//	// → "🌙 night"
func GetTimeOfDayDisplay(tempCtx *temporal.TemporalContext) TimeOfDayDisplay {
	timeOfDayIcon := "☀️" // Default to afternoon icon

	switch tempCtx.ExternalTime.TimeOfDay {
	case "morning":
		timeOfDayIcon = "🌅"
	case "afternoon":
		timeOfDayIcon = "☀️"
	case "evening":
		timeOfDayIcon = "🌆"
	case "night":
		timeOfDayIcon = "🌙"
	}

	return TimeOfDayDisplay{
		Icon:  timeOfDayIcon,
		Label: tempCtx.ExternalTime.TimeOfDay,
		Color: display.Yellow,
	}
}

// GetSessionPhaseDisplay returns formatted session phase with duration.
//
// What It Does:
// Formats session duration and phase with color coding based on how long
// session has been running (fresh/active=green, long=yellow).
//
// Parameters:
//   tempCtx: TemporalContext from hooks/lib/temporal containing InternalTime
//
// Returns:
//   SessionPhaseDisplay: Duration, phase, and color (or empty if no session)
//
// Phase Color Mapping:
//   fresh/active → Green (healthy session length)
//   long         → Yellow (consider taking break)
//
// Zero-State:
//   No duration → Returns empty SessionPhaseDisplay{}
//
// Example usage:
//
//	phaseDisplay := GetSessionPhaseDisplay(tempCtx)
//	if phaseDisplay.Duration != "" {
//	    fmt.Printf("%s (%s)\n", phaseDisplay.Duration, phaseDisplay.Phase)
//	    // → "1h23m (active)"
//	}
func GetSessionPhaseDisplay(tempCtx *temporal.TemporalContext) SessionPhaseDisplay {
	if tempCtx.InternalTime.ElapsedFormatted == "" {
		return SessionPhaseDisplay{} // No session data - return empty
	}

	// Determine phase color
	phaseColor := display.Green
	if tempCtx.InternalTime.SessionPhase == "long" {
		phaseColor = display.Yellow
	}

	return SessionPhaseDisplay{
		Duration: tempCtx.InternalTime.ElapsedFormatted,
		Phase:    tempCtx.InternalTime.SessionPhase,
		Color:    phaseColor,
	}
}

// GetScheduleDisplay returns formatted schedule/downtime status.
//
// What It Does:
// Formats current schedule activity or downtime status with appropriate icon
// and color coding for statusline display.
//
// Parameters:
//   tempCtx: TemporalContext from hooks/lib/temporal containing InternalSchedule
//
// Returns:
//   ScheduleDisplay: Activity, icon, and color (or empty if no schedule)
//
// Display Logic:
//   In work window      → Activity name with 📋 (green)
//   Expected downtime   → "Downtime" with ⏸️ (yellow)
//   No schedule data    → Empty ScheduleDisplay{}
//
// Example usage:
//
//	schedDisplay := GetScheduleDisplay(tempCtx)
//	if schedDisplay.Activity != "" {
//	    fmt.Printf("%s %s\n", schedDisplay.Icon, schedDisplay.Activity)
//	    // → "📋 Deep Work"
//	}
func GetScheduleDisplay(tempCtx *temporal.TemporalContext) ScheduleDisplay {
	// Check if in work window with activity
	if tempCtx.InternalSchedule.CurrentActivity != "" && tempCtx.InternalSchedule.InWorkWindow {
		return ScheduleDisplay{
			Activity: tempCtx.InternalSchedule.CurrentActivity,
			Color:    display.Green,
			Icon:     "📋",
		}
	}

	// Check if expected downtime (sleep, meal, break)
	if tempCtx.InternalSchedule.ExpectedDowntime {
		return ScheduleDisplay{
			Activity: "Downtime",
			Color:    display.Yellow,
			Icon:     "⏸️",
		}
	}

	// No schedule data
	return ScheduleDisplay{}
}

// GetCalendarDisplay returns formatted calendar information.
//
// What It Does:
// Formats calendar data (day of week, week number) into colored string for
// statusline display.
//
// Parameters:
//   tempCtx: TemporalContext from hooks/lib/temporal containing ExternalCalendar
//
// Returns:
//   string: Formatted calendar display (or empty string if no calendar data)
//
// Format:
//   "📅 Tuesday, Week 45" (in cyan)
//
// Zero-State:
//   No calendar data → Returns empty string ""
//
// Example usage:
//
//	calDisplay := GetCalendarDisplay(tempCtx)
//	if calDisplay != "" {
//	    fmt.Println(calDisplay)
//	    // → "📅 Tuesday, Week 45" (in cyan)
//	}
func GetCalendarDisplay(tempCtx *temporal.TemporalContext) string {
	if tempCtx.ExternalCalendar.DayOfWeek == "" {
		return "" // No calendar data
	}

	return fmt.Sprintf("%s📅 %s, Week %d%s",
		display.Cyan,
		tempCtx.ExternalCalendar.DayOfWeek,
		tempCtx.ExternalCalendar.WeekNumber,
		display.Reset)
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
//
// Testing Requirements:
//   - Import the library without errors
//   - Create mock TemporalContext with various time states
//   - Call each display function with mock contexts
//   - Verify icon mapping correct for each time of day
//   - Check session phase color coding (green vs yellow)
//   - Test schedule/downtime display logic
//   - Verify calendar formatting
//   - Ensure no panics for nil or empty contexts
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
//
// This is a LIBRARY, not an executable. No entry point, no main function.
// All functions wait to be called by statusline orchestrator.
//
// Usage: import "statusline/lib/temporal"
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - Memory: Four display struct allocations per statusline render (~200 bytes total)
//   - No file handles, network connections, or persistent resources
//   - Go's garbage collector handles memory automatically
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// Purpose: Transform temporal awareness data into formatted statusline displays
//
// Provides:
//   - Time of day display with contextual icons
//   - Session phase display with duration and color
//   - Schedule/downtime status indicators
//   - Calendar information formatting
//
// Architecture: Ladder - Middle Rung presentation layer
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify:
//   ✅ Add new temporal displays (circadian phase, next activity, etc.)
//   ✅ Adjust icon choices for different times/phases
//   ✅ Change color coding schemes
//   ✅ Extend display structures with additional fields
//
// Modify with Care:
//   ⚠️ Display struct fields - breaks statusline orchestrator
//   ⚠️ Function signatures - breaks all calling code
//   ⚠️ Icon/color semantics - changes visual meaning
//
// Never Modify:
//   ❌ Pure function guarantee (stateless, no side effects)
//   ❌ Graceful degradation (always return valid values)
//   ❌ Data vs Presentation separation (hooks/lib/temporal provides data)
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart" section above for complete structure.
//
// The Organizational Chart in BODY provides the detailed map showing:
// - All functions and their dependencies (ladder)
// - Complete execution flow paths (baton)
// - APU count (Available Processing Units)
//
// Quick summary:
// - 4 public APIs (GetTimeOfDayDisplay, GetSessionPhaseDisplay, GetScheduleDisplay, GetCalendarDisplay)
// - 0 helpers (all logic inline)
// - Ladder: Uses hooks/lib/temporal (lower rung) for data
// - Baton: TemporalContext → formatting → display structures
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Core Operations" subsection header comments above for detailed
// extension points. Each subsection includes "Extension Point" guidance showing:
// - Where to add new functionality
// - What naming pattern to follow
// - How to integrate with existing code
// - What tests to update
//
// Quick reference (details in BODY subsection comments):
//
// Adding New Temporal Displays:
//   1. Check if hooks/lib/temporal provides the data (if not, extend temporal lib first)
//   2. Define new display structure in SETUP Types section
//   3. Create Get[Aspect]Display() function in BODY
//   4. Follow existing patterns: extract from TemporalContext, format, return display struct
//   5. Update API documentation with new display
//   6. Add tests for new display function
//
// Extending Existing Displays:
//   1. Locate display function in BODY
//   2. Add new fields to display structure in SETUP Types
//   3. Extract additional data from TemporalContext
//   4. Update API documentation with new fields
//   5. Test with various temporal states
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Types: Display structs (~50 bytes each)
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
//
// Function Performance:
// - GetTimeOfDayDisplay(): O(1) switch statement, ~50 byte allocation
// - GetSessionPhaseDisplay(): O(1) check + assignment, ~50 byte allocation
// - GetScheduleDisplay(): O(1) conditional logic, ~50 byte allocation
// - GetCalendarDisplay(): O(1) sprintf, ~50 byte allocation
//
// Memory: Four struct allocations per statusline render (~200 bytes total)
//
// Key optimization: Already optimal - no I/O, minimal allocations
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// See BODY function docstrings above for operation-specific troubleshooting.
// Functions that commonly have issues include "Troubleshooting" sections in
// their docstrings with problem/check/solution patterns.
//
// Quick reference (details in BODY function docstrings above):
//
// No common failure modes - all operations guaranteed to succeed.
//
// Expected Behaviors:
//   - Empty session phase when no duration → Correct zero-state
//   - Empty schedule when no activity → Correct zero-state
//   - Empty calendar string when no data → Correct zero-state
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependency: hooks/lib/temporal (data collection layer)
// - Primary consumer: statusline orchestrator (main command)
// - Parallel libraries: lib/format, lib/session, lib/system, lib/git (other presentation layers)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Time of day display - COMPLETED
//   ✓ Session phase display - COMPLETED
//   ✓ Schedule display - COMPLETED
//   ✓ Calendar display - COMPLETED
//   ⏳ Circadian phase display
//   ⏳ Next activity display
//   ⏳ Time until next scheduled item
//
// Research Areas:
//   - Color coding based on time of day (morning/afternoon/evening/night)
//   - Schedule conflict indicators
//   - Timezone-aware displays
//   - Holiday-specific formatting
//
// Integration Targets:
//   - Calendar systems (iCal, Google Calendar integration)
//   - Time tracking tools (historical session patterns)
//   - Productivity dashboards
//   - Circadian rhythm monitoring
//
// Known Limitations to Address:
//   - No timezone display (assumes local time)
//   - Schedule display limited to single line
//   - Calendar display without holiday context
//   - No color variation based on time of day
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//         - Applied full GO library template (130 lines → comprehensive docs)
//         - Added complete METADATA section (Biblical foundation, CPI-SI identity)
//         - Expanded SETUP with detailed type documentation
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 11 sections
//         - Established architectural pattern: data vs presentation separation
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - Basic temporal display formatting (time, phase, schedule, calendar)
//         - Four display structures (TimeOfDayDisplay, SessionPhaseDisplay, etc.)
//         - Icon selection based on time of day
//         - Graceful degradation for missing data
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This library is a LADDER component (middle rung) - provides presentation
// layer for temporal awareness. Receives data from hooks/lib/temporal (lower rung),
// formats for statusline display. Part of the pattern where data collection
// happens in hooks libs, presentation happens in statusline libs.
//
// "To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Setup:
//
//     import "statusline/lib/temporal"
//     import "hooks/lib/temporal"
//
//     tempCtx, _ := temporal.GetTemporalContext()
//
//     // Get temporal displays
//     todDisplay := temporal.GetTimeOfDayDisplay(tempCtx)
//     phaseDisplay := temporal.GetSessionPhaseDisplay(tempCtx)
//     schedDisplay := temporal.GetScheduleDisplay(tempCtx)
//     calDisplay := temporal.GetCalendarDisplay(tempCtx)
//
// ============================================================================
// END CLOSING
// ============================================================================
