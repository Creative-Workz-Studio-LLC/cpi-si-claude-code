// METADATA
//
// Git Display Library - CPI-SI Statusline
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Let all things be done decently and in order" - 1 Corinthians 14:40
// Principle: Clarity and order in communication - presenting information in ways that serve understanding
// Anchor: "A word fitly spoken is like apples of gold in pictures of silver" - Proverbs 25:11
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - Middle Rung)
// Role: Presentation layer for git repository status in statusline display
// Paradigm: CPI-SI framework component - formats assessment data for user visibility
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
//   0.1.0 (2025-10-24) - Initial implementation with basic git display formatting
//
// Purpose & Function
//
// Purpose: Transform git repository assessment data into formatted display strings
// suitable for statusline space constraints
//
// Core Design: Pure presentation layer - receives assessment from system/lib/git,
// outputs formatted display strings with visual indicators (colors, icons, symbols)
//
// Key Features:
//   - Branch name display with visual indicators for dirty status
//   - Ahead/behind tracking with directional arrows (↑↓)
//   - Color coding and icon selection for visual clarity
//   - Graceful handling of non-repository directories
//   - Zero-state awareness (returns HasInfo: false for non-git directories)
//
// Philosophy: Show repository context without overwhelming - branch, changes,
// sync status at a glance. Assessment happens in system lib, presentation happens here.
//
// Blocking Status
//
// Non-blocking: Returns display structures immediately, never blocks statusline
// Mitigation: Uses system/lib/git for all potentially blocking operations (git commands)
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/git"
//
// Integration Pattern:
//   1. Call GetGitDisplay(workdir) with working directory path
//   2. Check HasInfo field to determine if git information available
//   3. Use DisplayString, Color, Icon fields for statusline formatting
//   4. No cleanup needed - stateless library
//
// Public API (in typical usage order):
//
//   Display Formatting (presentation):
//     GetGitDisplay(workdir string) GitDisplay
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt
//   External: None
//   Internal: system/lib/git (assessment), system/lib/display (color constants)
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Ladder: Uses system/lib/git for assessment (lower rung)
//   - Baton: Receives workdir, returns GitDisplay structure
//   - Rails: N/A (pure function library, no logging infrastructure)
//
// Health Scoring
//
// Pure presentation library - no health scoring infrastructure needed.
// Functions are guaranteed to succeed (graceful degradation on missing data).
//
// Display String Generation:
//   - Branch formatting: Always succeeds (empty string if no branch)
//   - Status indicators: Always succeeds (omitted if not applicable)
//   - Color/icon assignment: Always succeeds (defaults provided)
//
// Note: This library cannot fail - all operations return valid display structures.
// Health tracking would measure "successfully did nothing" which provides no value.
package git

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

	"fmt" // String formatting for ahead/behind indicators

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	gitlib "system/lib/git"   // Assessment data (branch, dirty, ahead/behind)
	"system/lib/display"      // Color constants for terminal output
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

// GitDisplay represents formatted git repository information for statusline display.
//
// Contains all presentation elements needed to render git status in statusline -
// the formatted string, visual styling (color/icon), and availability flag.
//
// Zero value represents "no git info available" state (HasInfo: false, empty strings).
// This allows calling code to distinguish between "not a repo" vs "repo with no changes".
//
// Example usage:
//
//     gitDisplay := git.GetGitDisplay("/home/user/project")
//     if gitDisplay.HasInfo {
//         fmt.Printf("%s %s %s\n", gitDisplay.Icon, gitDisplay.DisplayString, gitDisplay.Color)
//     }
type GitDisplay struct {
	DisplayString string // Formatted git status (e.g., "main*↑2↓1")
	Color         string // Terminal color code for display
	Icon          string // Visual icon representing git (e.g., "🌿")
	HasInfo       bool   // True if git information available, false for non-repos
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
//   └── GetGitDisplay() → uses system/lib/git.GetInfo() [external assessment]
//
//   Helpers/Core Ops: None (single-function library)
//
// Baton Flow (Execution Paths):
//
//   Entry → GetGitDisplay(workdir)
//     ↓
//   system/lib/git.GetInfo(workdir) [assessment]
//     ↓
//   Format display string (branch, dirty, ahead/behind)
//     ↓
//   Assign color and icon
//     ↓
//   Exit → return GitDisplay
//
// APUs (Available Processing Units):
// - 1 function total
// - 0 helpers (all logic inline)
// - 0 core operations (presentation only)
// - 1 public API (exported interface)

// ────────────────────────────────────────────────────────────────
// Display Formatting - Presentation Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Transform git assessment data (from system/lib/git) into formatted display
// strings with visual indicators suitable for statusline constraints.
//
// Why Separated:
// Single-function library - all logic contained in GetGitDisplay(). No need
// for helpers or internal operations. Assessment handled by system lib.
//
// Extension Point:
// To add new git display indicators (stashes, conflicts, etc.):
//   1. Ensure system/lib/git provides the assessment data
//   2. Add formatting logic to GetGitDisplay() after ahead/behind section
//   3. Follow same pattern: append symbol + value if applicable
//   4. Update GitDisplay docstring with new indicator examples
//   5. Create API documentation showing new indicators

// GetGitDisplay returns formatted git repository information for statusline display.
//
// What It Does:
// Retrieves git assessment data from system/lib/git and formats it into a display
// string with visual indicators. Shows branch name, dirty status (*), and ahead/behind
// tracking (↑↓) in a compact format optimized for statusline space constraints.
//
// Parameters:
//   workdir: Absolute path to working directory to check for git repository
//
// Returns:
//   GitDisplay: Formatted display structure with all presentation elements
//
// Display Format Examples:
//   "main"         → Clean branch, no changes, synced with remote
//   "main*"        → Branch with uncommitted changes
//   "main↑2"       → Branch 2 commits ahead of remote
//   "main↓1"       → Branch 1 commit behind remote
//   "main*↑2↓1"    → Dirty branch, ahead 2, behind 1
//   ""             → Not a git repository (HasInfo: false)
//
// Behavior:
//   - Non-repo directories: Returns GitDisplay{HasInfo: false} with empty strings
//   - No remote tracking: Shows branch and dirty status, omits ahead/behind
//   - Detached HEAD: Shows short commit hash (provided by system/lib/git)
//
// Example usage:
//
//	gitDisplay := GetGitDisplay("/home/user/project")
//	if gitDisplay.HasInfo {
//	    fmt.Printf("%s %s\n", gitDisplay.Icon, gitDisplay.DisplayString)
//	}
func GetGitDisplay(workdir string) GitDisplay {
	// Get assessment data from system lib (branch, dirty, ahead/behind)
	gitInfo := gitlib.GetInfo(workdir)

	// No git information available - return zero-state display
	if gitInfo.Branch == "" {
		return GitDisplay{HasInfo: false}
	}

	// Start with branch name
	gitDisplayStr := gitInfo.Branch

	// Append dirty indicator if uncommitted changes exist
	if gitInfo.Dirty {
		gitDisplayStr += "*"
	}

	// Append ahead indicator if commits ahead of remote
	if gitInfo.Ahead > 0 {
		gitDisplayStr += fmt.Sprintf("↑%d", gitInfo.Ahead)
	}

	// Append behind indicator if commits behind remote
	if gitInfo.Behind > 0 {
		gitDisplayStr += fmt.Sprintf("↓%d", gitInfo.Behind)
	}

	// Return complete display structure with visual elements
	return GitDisplay{
		DisplayString: gitDisplayStr,
		Color:         display.Green, // Green indicates git context available
		Icon:          "🌿",           // Branch/nature icon represents git
		HasInfo:       true,
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
//   - Call GetGitDisplay() with git repository directory
//   - Call GetGitDisplay() with non-repository directory
//   - Verify DisplayString format matches expected patterns
//   - Confirm HasInfo correctly distinguishes repo vs non-repo
//   - Check Color and Icon fields populated for valid repos
//   - Ensure no panics or errors for edge cases (empty paths, invalid dirs)
//   - Run: go vet ./... (no warnings)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - gofmt -l . (no formatting issues)
//
// Integration Testing:
//   - Test with actual statusline orchestrator
//   - Verify display in terminal with different git states
//   - Check visual indicators render correctly
//   - Validate space usage in full statusline context
//
// Example validation code:
//
//     // Test with git repository
//     display := git.GetGitDisplay("/home/user/project")
//     if !display.HasInfo {
//         t.Error("Expected HasInfo true for git repository")
//     }
//     if display.DisplayString == "" {
//         t.Error("Expected non-empty DisplayString")
//     }
//
//     // Test with non-repository
//     display = git.GetGitDisplay("/tmp")
//     if display.HasInfo {
//         t.Error("Expected HasInfo false for non-repository")
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
// Usage: import "statusline/lib/git"
//
// The library is imported into the calling package (typically statusline orchestrator),
// making GetGitDisplay() available. No code executes during import - function is
// defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "statusline/lib/git"
//
//     func main() {
//         // Get git display for current directory
//         gitDisplay := git.GetGitDisplay(".")
//
//         if gitDisplay.HasInfo {
//             fmt.Printf("%s %s\n", gitDisplay.Icon, gitDisplay.DisplayString)
//         }
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Memory: Single GitDisplay struct allocation per call (~50 bytes)
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
//   - Small allocations (~50 bytes per call) don't require manual management
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
//   - Transforms git assessment data into formatted statusline display strings
//   - Shows branch, dirty status, ahead/behind tracking with visual indicators
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
//   ✅ Add new visual indicators (stashes, conflicts, etc.)
//   ✅ Adjust color or icon choices for different git states
//   ✅ Extend DisplayString format with additional symbols
//   ✅ Add helper functions for complex formatting logic
//   ✅ Create additional display structure types for different views
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ GitDisplay struct fields - breaks statusline orchestrator
//   ⚠️ GetGitDisplay() signature - breaks all calling code
//   ⚠️ DisplayString format - breaks parsing/expectation in consumers
//   ⚠️ HasInfo semantics - breaks zero-state handling
//   ⚠️ Color/Icon field types - breaks display rendering
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Pure function guarantee (stateless, no side effects)
//   ❌ Graceful degradation (always return valid display)
//   ❌ Assessment vs Presentation separation (system lib does assessment)
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
// - 1 public API (GetGitDisplay) orchestrates presentation logic
// - Ladder: Uses system/lib/git (lower rung) for assessment data
// - Baton: workdir → assessment → formatting → GitDisplay structure
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// Adding New Git Indicators:
//   1. Ensure system/lib/git.Info struct provides the data field
//   2. In GetGitDisplay(), after ahead/behind section, add new indicator logic
//   3. Follow pattern: if gitInfo.[Field] [condition] { gitDisplayStr += "[symbol][value]" }
//   4. Update GitDisplay docstring "Display Format Examples" with new patterns
//   5. Create API documentation showing new indicator behavior
//   6. Update tests to verify new indicator appears correctly
//
// Example pattern for adding stash indicator:
//
//     // After ahead/behind section in GetGitDisplay()
//     if gitInfo.Stashes > 0 {
//         gitDisplayStr += fmt.Sprintf("$%d", gitInfo.Stashes)
//     }
//     // Result: "main*↑2$3" (dirty, 2 ahead, 3 stashes)
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Types: GitDisplay struct (~50 bytes - minimal memory footprint)
//
// Quick summary:
// - GetGitDisplay(): O(1) formatting, ~50 byte allocation, <1 microsecond
// - Memory: Single struct allocation per call, garbage collected automatically
// - No I/O operations (assessment done by system/lib/git)
// - No caching needed - formatting is trivial cost
//
// Key optimization: This library needs no optimization. String formatting
// and struct allocation are already optimal for this use case.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// This library has no common failure modes - all operations guaranteed to succeed.
//
// Expected Behaviors:
//   - Non-repository returns HasInfo: false - This is correct, not an error
//   - No remote tracking omits ahead/behind - This is correct, not an error
//   - Empty DisplayString when HasInfo: false - This is correct zero-state
//
// If GetGitDisplay() returns unexpected results:
//   Problem: HasInfo true but DisplayString empty
//     - Cause: system/lib/git.GetInfo() returned Info with Branch set but empty
//     - Solution: Check system/lib/git implementation - should return Branch: ""
//
//   Problem: Ahead/behind not showing when expected
//     - Cause: Repository has no remote tracking branch
//     - Solution: Expected behavior - ahead/behind requires remote tracking
//     - Note: system/lib/git returns Ahead: 0, Behind: 0 when no remote
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Key dependencies: system/lib/git (assessment), system/lib/display (colors)
// - Primary consumer: statusline orchestrator (main command)
//
// Parallel Implementation:
//   - Go version: This file (statusline/lib/git/git.go)
//   - Shared philosophy: Assessment in system lib, presentation in statusline lib
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Branch name display - COMPLETED
//   ✓ Dirty status indicator - COMPLETED
//   ✓ Ahead/behind tracking - COMPLETED
//   ⏳ Stash indicator (if system/lib/git provides data)
//   ⏳ Conflict indicator (if system/lib/git provides data)
//   ⏳ Multiple remote tracking support
//   ⏳ Custom indicator symbols via configuration
//
// Research Areas:
//   - Color coding based on git state (dirty=yellow, conflicts=red)
//   - Abbreviated branch names for very long branch names
//   - User-configurable display format strings
//   - Icon variations based on git state (different icons for different states)
//
// Known Limitations to Address:
//   - Only shows single remote tracking branch (most repos have one remote)
//   - No visual distinction between different dirty states (staged vs unstaged)
//   - Fixed color (green) regardless of git state
//   - No configuration options (format, symbols, colors all hardcoded)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//         - Applied full GO library template (1014 lines → comprehensive docs)
//         - Added complete METADATA section (Biblical foundation, CPI-SI identity)
//         - Expanded SETUP with detailed type documentation
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 13 sections
//         - Established architectural pattern: assessment vs presentation separation
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - Basic git display formatting (branch, dirty, ahead/behind)
//         - GitDisplay struct with DisplayString, Color, Icon, HasInfo
//         - GetGitDisplay() function using system/lib/git for assessment
//         - Graceful degradation for non-repository directories
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component (middle rung) - provides presentation
// layer for git repository status. Sits between system/lib/git (assessment)
// and statusline orchestrator (display). Part of the architectural pattern
// where assessment logic lives in system libs, presentation logic lives in
// statusline libs.
//
// Modify thoughtfully - changes here affect statusline display. The separation
// between assessment (system/lib/git) and presentation (this library) must be
// maintained. Never duplicate assessment logic here - always use system lib.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with real git repositories before committing
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain assessment vs presentation separation
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
//     import "statusline/lib/git"
//
//     // Get git display for working directory
//     gitDisplay := git.GetGitDisplay("/home/user/project")
//
// Check if Git Info Available:
//
//     if gitDisplay.HasInfo {
//         fmt.Printf("%s %s\n", gitDisplay.Icon, gitDisplay.DisplayString)
//     } else {
//         fmt.Println("Not a git repository")
//     }
//
// Full Statusline Integration:
//
//     func buildStatusline(workdir string) string {
//         gitDisplay := git.GetGitDisplay(workdir)
//
//         var parts []string
//         if gitDisplay.HasInfo {
//             // Format: 🌿 main*↑2↓1
//             gitPart := fmt.Sprintf("%s %s", gitDisplay.Icon, gitDisplay.DisplayString)
//             parts = append(parts, gitPart)
//         }
//
//         // ... add other statusline parts ...
//
//         return strings.Join(parts, " | ")
//     }
//
// Handling Different Git States:
//
//     gitDisplay := git.GetGitDisplay(workdir)
//     if !gitDisplay.HasInfo {
//         return "no git"  // Not a repository
//     }
//
//     // Display will be formatted based on state:
//     // Clean: "main"
//     // Dirty: "main*"
//     // Ahead: "main↑2"
//     // Behind: "main↓1"
//     // Complex: "main*↑2↓1"
//
// ============================================================================
// END CLOSING
// ============================================================================
