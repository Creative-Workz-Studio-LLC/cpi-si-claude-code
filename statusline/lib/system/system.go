// METADATA
//
// System Health Display Library - CPI-SI Statusline
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Let all things be done decently and in order" - 1 Corinthians 14:40
// Principle: Order and clarity in communication - showing system health clearly and truthfully
// Anchor: "A faithful witness will not lie: but a false witness will utter lies" - Proverbs 14:5
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - Middle Rung)
// Role: Presentation layer for system health metrics in statusline display
// Paradigm: CPI-SI framework component - formats system data for user visibility
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
//   0.1.0 (2025-10-24) - Initial implementation with system health display functions
//
// Purpose & Function
//
// Purpose: Transform system health metrics (load average, memory, disk) into formatted
// display structures with health-based color coding suitable for statusline constraints
//
// Core Design: Pure presentation layer - receives system data from system/lib/system,
// outputs formatted display structures with health thresholds and visual indicators
//
// Key Features:
//   - Load average display with CPU-relative health coloring
//   - Memory usage display with usage-based health coloring
//   - Disk usage display with capacity-based health coloring
//   - Zero-state handling (gracefully omits unavailable data)
//   - Visual clarity through color coding and icons
//   - Health thresholds for quick status recognition
//
// Philosophy: Show system health truthfully without alarm - what's happening, is it healthy.
// Data collection happens in system lib, presentation happens here.
//
// Blocking Status
//
// Non-blocking: Returns display structures immediately, never blocks statusline
// Mitigation: Uses system/lib/system which reads /proc files directly (fast operations)
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/system"
//
// Integration Pattern:
//   1. Call display functions to get formatted health data
//   2. Check HasInfo field to determine if data available
//   3. Use display fields (LoadAvg, Percent, etc.) for statusline formatting
//   4. Apply Color and Icon for visual styling
//   5. No cleanup needed - stateless library
//
// Public API (in typical usage order):
//
//   System Health (presentation):
//     GetLoadDisplay() LoadDisplay
//     GetMemoryDisplay() MemoryDisplay
//     GetDiskDisplay(path string) DiskDisplay
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: None
//   External: None
//   Internal: system/lib/system (system data collection), system/lib/display (color constants)
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Ladder: Uses system/lib/system for data collection (lower rung)
//   - Baton: Receives system data, returns display structures
//   - Rails: N/A (pure function library, no logging infrastructure)
//
// Health Scoring
//
// Pure presentation library - no health scoring infrastructure needed.
// Functions are guaranteed to succeed (graceful degradation on missing data).
//
// Display Generation:
//   - Load display: Always succeeds (zero-state if no load data)
//   - Memory display: Always succeeds (zero-state if no memory data)
//   - Disk display: Always succeeds (zero-state if path invalid)
//
// Note: This library cannot fail - all operations return valid display structures.
// Health tracking would measure "successfully did nothing" which provides no value.
package system

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
	// Synchronization for shared system snapshots per render.
	"sync"

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	syslib "system/lib/system" // System data collection (load, memory, disk)
	"system/lib/display"       // Color constants for terminal output
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

// LoadDisplay represents formatted load average with health coloring.
//
// Contains 1-minute load average, health-based color (relative to CPU count),
// and availability flag. Used to show system load in statusline.
//
// Zero value represents "no load data" state (HasInfo: false, LoadAvg: 0).
//
// Example usage:
//
//     loadDisplay := system.GetLoadDisplay()
//     if loadDisplay.HasInfo {
//         fmt.Printf("%s %.2f\n", loadDisplay.Icon, loadDisplay.LoadAvg)
//     }
type LoadDisplay struct {
	LoadAvg float64 // 1-minute load average
	Color   string  // Terminal color code based on load/CPU ratio
	Icon    string  // Visual icon representing load (e.g., "⚡")
	HasInfo bool    // True if load data available, false if no data
}

// MemoryDisplay represents formatted memory usage with health coloring.
//
// Contains used/total memory in GB, usage percentage, health-based color,
// and availability flag. Used to show memory status in statusline.
//
// Zero value represents "no memory data" state (HasInfo: false, all values: 0).
//
// Example usage:
//
//     memDisplay := system.GetMemoryDisplay()
//     if memDisplay.HasInfo {
//         fmt.Printf("%s %.1fGB/%.1fGB (%.0f%%)\n",
//             memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB, memDisplay.Percent)
//     }
type MemoryDisplay struct {
	UsedGB  float64 // Memory used in GB
	TotalGB float64 // Total memory in GB
	Percent float64 // Usage percentage (0-100)
	Color   string  // Terminal color code based on usage percent
	Icon    string  // Visual icon representing memory (e.g., "💾")
	HasInfo bool    // True if memory data available, false if no data
}

// DiskDisplay represents formatted disk usage with health coloring.
//
// Contains usage percentage, health-based color, and availability flag.
// Used to show disk capacity status in statusline.
//
// Zero value represents "no disk data" state (HasInfo: false, Percent: 0).
//
// Example usage:
//
//     diskDisplay := system.GetDiskDisplay("/")
//     if diskDisplay.HasInfo {
//         fmt.Printf("%s %.0f%% used\n", diskDisplay.Icon, diskDisplay.Percent)
//     }
type DiskDisplay struct {
	Percent float64 // Disk usage percentage (0-100)
	Color   string  // Terminal color code based on usage percent
	Icon    string  // Visual icon representing disk (e.g., "💿")
	HasInfo bool    // True if disk data available, false if no data
}

var (
	systemInfoOnce   sync.Once // Ensures single info capture per render
	cachedSystemInfo syslib.Info
)

// getSystemInfo retrieves and caches system info so multiple display components
// share a single snapshot (single df/proc read per statusline render).
func getSystemInfo() syslib.Info {
	systemInfoOnce.Do(func() {
		cachedSystemInfo = syslib.GetInfo()
	})
	return cachedSystemInfo
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
//   ├── GetLoadDisplay() → uses syslib.GetInfo() + health thresholds
//   ├── GetMemoryDisplay() → uses syslib.GetInfo() + health thresholds
//   └── GetDiskDisplay(path) → uses syslib.GetDiskUsage() + health thresholds
//
//   Health Thresholds (Presentation Logic):
//   ├── Load: >80% CPU = red, >50% = yellow, else green
//   ├── Memory: >80% = red, >60% = yellow, else green
//   └── Disk: >90% = red, >75% = yellow, else green
//
// Baton Flow (Execution Paths):
//
//   Entry → GetLoadDisplay()
//     ↓
//   syslib.GetInfo() [external - data collection]
//     ↓
//   Calculate load/CPU ratio
//     ↓
//   Apply health thresholds → assign color
//     ↓
//   Exit → return LoadDisplay
//
//   Entry → GetMemoryDisplay()
//     ↓
//   syslib.GetInfo() [external - data collection]
//     ↓
//   Calculate memory usage percent
//     ↓
//   Apply health thresholds → assign color
//     ↓
//   Exit → return MemoryDisplay
//
//   Entry → GetDiskDisplay(path)
//     ↓
//   syslib.GetDiskUsage(path) [external - data collection]
//     ↓
//   Extract usage percent
//     ↓
//   Apply health thresholds → assign color
//     ↓
//   Exit → return DiskDisplay
//
// APUs (Available Processing Units):
// - 3 functions total
// - 0 helpers (all presentation logic inline)
// - 3 public APIs (exported display functions)

// ────────────────────────────────────────────────────────────────
// Health Display Functions - Presentation Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Transform system health metrics (from system/lib/system) into formatted display
// structures with health-based color coding suitable for statusline constraints.
//
// Why Separated:
// Three distinct types of system health (load, memory, disk) each need
// independent health thresholds and formatting. Separation allows statusline
// orchestrator to choose which metrics to display based on available space.
//
// Extension Point:
// To add new system health metrics (network, CPU temp, etc.):
//   1. Define new display structure type in SETUP (follow existing pattern)
//   2. Create Get[Metric]Display() function following naming pattern
//   3. Get data from system/lib/system (or extend if metric not available)
//   4. Apply health thresholds for color selection
//   5. Update API documentation with new metric display
//   6. Add tests for new display function

// GetLoadDisplay returns formatted load average with health-based color coding.
//
// What It Does:
// Retrieves 1-minute load average from system/lib/system and applies health thresholds
// relative to CPU count. Shows if system is under load stress.
//
// Returns:
//   LoadDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   Load 1.5 on 4 CPUs → Green ⚡ (37.5% load - healthy)
//   Load 3.0 on 4 CPUs → Yellow ⚡ (75% load - elevated)
//   Load 4.5 on 4 CPUs → Red ⚡ (112.5% load - stressed)
//   No load data → HasInfo: false
//
// Health Thresholds:
//   - Load/CPU > 80%: Red (system stressed)
//   - Load/CPU > 50%: Yellow (elevated load)
//   - Load/CPU ≤ 50%: Green (healthy)
//
// Behavior:
//   - No load data: Returns HasInfo: false
//   - Any load: Calculates load/CPU ratio and applies thresholds
//
// Example usage:
//
//	loadDisplay := GetLoadDisplay()
//	if loadDisplay.HasInfo {
//	    fmt.Printf("%s %.2f load\n", loadDisplay.Icon, loadDisplay.LoadAvg)
//	}
func GetLoadDisplay() LoadDisplay {
	sysInfo := getSystemInfo()

	if sysInfo.LoadAvg1 <= 0 {
		return LoadDisplay{HasInfo: false}
	}

	// Color code based on load vs CPU count
	loadColor := display.Green
	loadPercent := (sysInfo.LoadAvg1 / float64(sysInfo.CPUCount)) * 100

	if loadPercent > 80 {
		loadColor = display.Red
	} else if loadPercent > 50 {
		loadColor = display.Yellow
	}

	return LoadDisplay{
		LoadAvg: sysInfo.LoadAvg1,
		Color:   loadColor,
		Icon:    "⚡",
		HasInfo: true,
	}
}

// GetMemoryDisplay returns formatted memory usage with health-based color coding.
//
// What It Does:
// Retrieves memory statistics from system/lib/system and applies health thresholds
// based on usage percentage. Shows current memory pressure.
//
// Returns:
//   MemoryDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   5.2GB/16GB (32%) → Green 💾 (healthy)
//   10.5GB/16GB (65%) → Yellow 💾 (elevated usage)
//   14.0GB/16GB (87%) → Red 💾 (high memory pressure)
//   No memory data → HasInfo: false
//
// Health Thresholds:
//   - Usage > 80%: Red (high memory pressure)
//   - Usage > 60%: Yellow (elevated memory usage)
//   - Usage ≤ 60%: Green (healthy)
//
// Behavior:
//   - No memory data: Returns HasInfo: false
//   - Any memory data: Calculates usage percent and applies thresholds
//
// Example usage:
//
//	memDisplay := GetMemoryDisplay()
//	if memDisplay.HasInfo {
//	    fmt.Printf("%s %.1fGB/%.1fGB\n", memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB)
//	}
func GetMemoryDisplay() MemoryDisplay {
	sysInfo := getSystemInfo()

	if sysInfo.MemTotalGB <= 0 {
		return MemoryDisplay{HasInfo: false}
	}

	memPercent := (sysInfo.MemUsedGB / sysInfo.MemTotalGB) * 100
	memColor := display.Green

	if memPercent > 80 {
		memColor = display.Red
	} else if memPercent > 60 {
		memColor = display.Yellow
	}

	return MemoryDisplay{
		UsedGB:  sysInfo.MemUsedGB,
		TotalGB: sysInfo.MemTotalGB,
		Percent: memPercent,
		Color:   memColor,
		Icon:    "💾",
		HasInfo: true,
	}
}

// GetDiskDisplay returns formatted disk usage with health-based color coding.
//
// What It Does:
// Retrieves disk usage for specified path from system/lib/system and applies
// health thresholds based on capacity percentage. Shows disk space availability.
//
// Parameters:
//   path: Filesystem path to check (e.g., "/", "/home")
//
// Returns:
//   DiskDisplay: Formatted display structure with all presentation elements
//
// Display Examples:
//   45% used → Green 💿 (plenty of space)
//   78% used → Yellow 💿 (running low)
//   92% used → Red 💿 (critically low)
//   Invalid path → HasInfo: false
//
// Health Thresholds:
//   - Usage > 90%: Red (critically low space)
//   - Usage > 75%: Yellow (running low)
//   - Usage ≤ 75%: Green (healthy)
//
// Behavior:
//   - Invalid path or no data: Returns HasInfo: false
//   - Any disk data: Uses existing percent and applies thresholds
//
// Example usage:
//
//	diskDisplay := GetDiskDisplay("/")
//	if diskDisplay.HasInfo {
//	    fmt.Printf("%s %.0f%% disk used\n", diskDisplay.Icon, diskDisplay.Percent)
//	}
func GetDiskDisplay(path string) DiskDisplay {
	diskInfo := syslib.GetDiskUsage(path)

	if diskInfo.UsagePercent < 0 {
		return DiskDisplay{HasInfo: false}
	}

	diskColor := display.Green
	if diskInfo.UsagePercent > 90 {
		diskColor = display.Red
	} else if diskInfo.UsagePercent > 75 {
		diskColor = display.Yellow
	}

	return DiskDisplay{
		Percent: diskInfo.UsagePercent,
		Color:   diskColor,
		Icon:    "💿",
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
//   - Call each display function
//   - Verify HasInfo correctly distinguishes data vs no-data states
//   - Check display fields populated correctly (LoadAvg, UsedGB/TotalGB/Percent)
//   - Confirm Color thresholds working (green/yellow/red based on metrics)
//   - Test with various system states (low/medium/high usage)
//   - Verify Icon fields set appropriately
//   - Ensure no panics or errors for edge cases
//   - Run: go vet ./... (no warnings)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - gofmt -l . (no formatting issues)
//
// Integration Testing:
//   - Test with actual statusline orchestrator
//   - Verify display in terminal with different system loads
//   - Check visual indicators render correctly
//   - Validate health color coding matches thresholds
//   - Test on different systems (various CPU counts, memory sizes)
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by other components.
//
// Usage: import "statusline/lib/system"
//
// The library is imported into the calling package (typically statusline orchestrator),
// making all display functions available. No code executes during import - functions
// are defined and ready to use.
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
//   - Transforms system health metrics into formatted statusline display structures
//   - Shows load average, memory usage, and disk capacity with health coloring
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
//   ✅ Add new system health metrics (network, CPU temp, etc.)
//   ✅ Adjust health threshold values (current: 50/80 for load, 60/80 for memory, 75/90 for disk)
//   ✅ Extend display structures with additional fields
//   ✅ Change color or icon choices for different metrics
//   ✅ Add alternative threshold strategies (e.g., warning levels)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Display struct fields - breaks statusline orchestrator
//   ⚠️ Function signatures - breaks all calling code
//   ⚠️ HasInfo semantics - breaks zero-state handling
//   ⚠️ Color/Icon field types - breaks display rendering
//   ⚠️ Health threshold semantics - changes what "healthy" means
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Pure function guarantee (stateless, no side effects)
//   ❌ Graceful degradation (always return valid display)
//   ❌ Data vs Presentation separation (system/lib/system provides data)
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
// The Organizational Chart in BODY provides the detailed map showing:
// - All functions and their dependencies (ladder)
// - Complete execution flow paths (baton)
// - APU count (Available Processing Units)
//
// Quick architectural summary (details in BODY Organizational Chart):
// - 3 public APIs (GetLoadDisplay, GetMemoryDisplay, GetDiskDisplay)
// - 0 helpers (all logic inline)
// - Ladder: Uses system/lib/system (lower rung) for data collection
// - Baton: System data → health thresholds → display structures
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
// Adding New System Health Metrics:
//   1. Check if system/lib/system provides the data (if not, extend system lib first)
//   2. Define new display structure in SETUP Types section
//   3. Create Get[Metric]Display() function in BODY
//   4. Apply health thresholds for color selection
//   5. Update API documentation with new metric
//   6. Add tests for new display function
//
// Adjusting Health Thresholds:
//   1. Locate threshold values in Get[Metric]Display() functions
//   2. Modify percentage thresholds (current values documented in function comments)
//   3. Test with various system states to verify new thresholds make sense
//   4. Update documentation with new threshold values
//   5. Consider: Do new thresholds align with industry standards?
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
// - GetLoadDisplay(): O(1) /proc read + calculation, ~50 byte allocation, <100 microseconds
// - GetMemoryDisplay(): O(1) /proc read + calculation, ~50 byte allocation, <100 microseconds
// - GetDiskDisplay(): O(1) df command + parsing, ~50 byte allocation, <500 microseconds
//
// Memory: Three struct allocations per statusline render (~150 bytes total)
//
// Key optimization: This library needs no optimization. System data retrieval and
// struct allocation are already optimal for this use case. The system/lib/system
// reads /proc directly which is fast.
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
// This library has no common failure modes - all operations guaranteed to succeed.
//
// Expected Behaviors:
//   - Zero load returns HasInfo: false - This is correct, not an error
//   - Zero memory returns HasInfo: false - This is correct, not an error
//   - Invalid disk path returns HasInfo: false - This is correct, not an error
//   - Empty display data when HasInfo: false - This is correct zero-state
//
// If display functions return unexpected results:
//   Problem: HasInfo true but color seems wrong
//     - Cause: Health thresholds may not match expectations
//     - Solution: Verify threshold values in function code
//     - Check: Load (>80%/>50%), Memory (>80%/>60%), Disk (>90%/>75%)
//
//   Problem: Load percentage calculation seems off
//     - Cause: Load average relative to CPU count
//     - Solution: Load 4.0 on 4 CPUs = 100% (expected behavior)
//     - Formula: (LoadAvg1 / CPUCount) * 100
//
//   Problem: HasInfo false when system clearly has data
//     - Cause: system/lib/system not returning data
//     - Solution: Check system/lib/system implementation
//     - Debug: Call syslib.GetInfo() directly to verify
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependency: system/lib/system (data collection layer)
// - Primary consumer: statusline orchestrator (main command)
// - Parallel libraries: lib/format, lib/session, lib/git (other presentation layers)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Load average display - COMPLETED
//   ✓ Memory usage display - COMPLETED
//   ✓ Disk usage display - COMPLETED
//   ⏳ Network throughput display
//   ⏳ CPU temperature display
//   ⏳ Process count display
//   ⏳ Configurable health thresholds
//
// Research Areas:
//   - Adaptive thresholds based on historical patterns
//   - Multi-level health indicators (warning/critical separation)
//   - Trend indicators (↑↓ for increasing/decreasing metrics)
//   - Threshold configuration via external file
//
// Integration Targets:
//   - Monitoring systems (Prometheus, Grafana integration for historical tracking)
//   - Alert systems (threshold breach notifications)
//   - Performance profiling tools
//   - System dashboards (centralized health monitoring)
//
// Known Limitations to Address:
//   - Health thresholds hardcoded (not configurable)
//   - Binary health states (green/yellow/red - no gradients)
//   - No historical context (can't show "getting worse")
//   - Disk display only shows single path (not multiple mounts)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//         - Applied full GO library template (136 lines → comprehensive docs)
//         - Added complete METADATA section (Biblical foundation, CPI-SI identity)
//         - Expanded SETUP with detailed type documentation
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 11 sections
//         - Established architectural pattern: data vs presentation separation
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - Basic system health formatting (load, memory, disk)
//         - Three display structures (LoadDisplay, MemoryDisplay, DiskDisplay)
//         - Health threshold coloring (green/yellow/red)
//         - Graceful degradation for missing data
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component (middle rung) - provides presentation
// layer for system health metrics. Receives data from system/lib/system (lower rung),
// formats for statusline display. Part of the architectural pattern where data
// collection happens in system libs, presentation happens in statusline libs.
//
// Modify thoughtfully - changes here affect statusline health display. The separation
// between data collection (system/lib/system) and presentation (this library) must
// be maintained. Never duplicate data collection logic here - always receive data
// via system/lib/system.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with various system states before committing
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain data vs presentation separation
//
// "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "statusline/lib/system"
//
//     // Get system health displays
//     loadDisplay := system.GetLoadDisplay()
//     memDisplay := system.GetMemoryDisplay()
//     diskDisplay := system.GetDiskDisplay("/")
//
// Check if Data Available:
//
//     if loadDisplay.HasInfo {
//         fmt.Printf("%s %.2f load\n", loadDisplay.Icon, loadDisplay.LoadAvg)
//     }
//
//     if memDisplay.HasInfo {
//         fmt.Printf("%s %.1f/%.1fGB\n", memDisplay.Icon, memDisplay.UsedGB, memDisplay.TotalGB)
//     }
//
//     if diskDisplay.HasInfo {
//         fmt.Printf("%s %.0f%% used\n", diskDisplay.Icon, diskDisplay.Percent)
//     }
//
// Full Statusline Integration:
//
//     func buildStatusline() string {
//         var parts []string
//
//         // System health
//         loadDisplay := system.GetLoadDisplay()
//         if loadDisplay.HasInfo {
//             parts = append(parts, fmt.Sprintf("%.2f", loadDisplay.LoadAvg))
//         }
//
//         memDisplay := system.GetMemoryDisplay()
//         if memDisplay.HasInfo {
//             parts = append(parts, fmt.Sprintf("%.0f%% mem", memDisplay.Percent))
//         }
//
//         diskDisplay := system.GetDiskDisplay("/")
//         if diskDisplay.HasInfo {
//             parts = append(parts, fmt.Sprintf("%.0f%% disk", diskDisplay.Percent))
//         }
//
//         return strings.Join(parts, " | ")
//         // → "1.25 | 65% mem | 45% disk"
//     }
//
// ============================================================================
// END CLOSING
// ============================================================================
