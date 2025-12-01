// METADATA
//
// Path Formatting Library - Statusline Display Optimization
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Make the path straight" - Matthew 3:3
// Principle: Clarity in direction - presenting paths in ways that serve understanding and navigation
// Anchor: "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// CPI-SI Identity
//
// Component Type: LADDER - Presentation formatting layer
// Role: Converts verbose filesystem paths to concise display forms for statusline space constraints
// Paradigm: CPI-SI framework component implementing clarity through intelligent truncation
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 1.0.0
// Last Modified: 2025-11-04 - Full template application
//
// Version History:
//   1.0.0 (2025-11-04) - Full template application with comprehensive documentation
//   0.1.0 (2024-10-24) - Initial implementation
//
// Purpose & Function
//
// Purpose: Shortens verbose filesystem paths to readable display forms suitable
// for statusline space constraints while maintaining clear location identification.
//
// Core Design: Two-stage optimization - home directory replacement with ~ tilde,
// then basename extraction for paths exceeding length threshold.
//
// Key Features:
//   - Home directory shortening: /home/user/project → ~/project
//   - Length-based truncation: Long paths show basename only
//   - Configurable threshold: 40 character limit for statusline
//   - Graceful fallback: Returns original path if shortening not needed
//
// Philosophy: Filesystem paths in statusline show location context, not full
// hierarchy. Home replacement provides familiar shorthand, basename shows what
// matters (current location) when full path too verbose.
//
// Blocking Status
//
// Non-blocking: Minimal I/O (single UserHomeDir call), no network, no state.
// Degrades gracefully - if home directory unavailable, returns path unchanged.
// Mitigation: Error from UserHomeDir ignored, falls back to original path.
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/format"
//
// Integration Pattern:
//   1. Obtain working directory path from context
//   2. Call ShortenPath with full path
//   3. Display shortened path in statusline
//   4. User sees concise location identifier
//
// Public API (in typical usage order):
//
//   Path Shortening (display optimization):
//     ShortenPath(path string) string
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os (home directory), path/filepath (basename), strings (prefix operations)
//   External: None
//   Internal: None
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Called by statusline orchestrator during working directory display
//   - Working directory path provided via context
//   - Shortened result displayed in statusline output
//
// Health Scoring
//
// Path shortening operates on Base100 scale:
//
// Home Directory Replacement:
//   - Successful home directory retrieval: +25
//   - Successful tilde substitution: +25
//
// Length Evaluation:
//   - Length check and decision: +25
//   - Path return (original or basename): +25
//
// Total: 100 points for complete operation
//
// Note: Home directory retrieval can fail (no home), but operation continues
// with graceful degradation. Health reflects successful completion of operation.
package format

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
// Dependencies for path manipulation and home directory detection.

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"os"            // Home directory detection
	"path/filepath" // Basename extraction
	"strings"       // String prefix operations for home replacement
	"unicode/utf8"  // Rune-aware truncation for path ellipsis handling
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Configuration values defining path shortening behavior.

const (
	//--- Display Length Configuration ---
	// Controls when paths are truncated for statusline display.

	// MaxPathLength defines the character limit for full path display.
	//
	// Set to 40 characters to fit comfortably in statusline alongside
	// other information (git, session stats, system metrics).
	//
	// Reasoning: Statusline typically 80-120 characters total. Allocating
	// 40 to path leaves room for other essential information. Longer paths
	// would crowd out other context.
	//
	// Paths exceeding this get compacted using an ellipsis while preserving head/tail
	// context so multiple similarly named directories remain distinguishable.
	MaxPathLength = 40

	// PathPrefixLength controls how many characters from the start of the path are
	// preserved when compaction is necessary. Chosen to retain home/root indicators.
	PathPrefixLength = 15

	// PathSuffixLength controls how many characters from the end of the path are
	// preserved when compaction is necessary. Chosen to keep the immediate directory.
	PathSuffixLength = 20
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// This component uses no custom types - operates on primitives only.

// No types needed

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// This component maintains no state - pure stateless functions.

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
// Maps bidirectional dependencies and baton flow within this component.
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Single Rung - Self-Contained)
//   └── ShortenPath() → pure logic (no helpers needed)
//
// Baton Flow (Execution Paths):
//
//   Entry → ShortenPath(path)
//     ↓
//   Home directory replacement (~ substitution)
//     ↓
//   Length evaluation (exceeds MaxPathLength?)
//     ↓
//   Return shortened path OR basename
//     ↓
//   Exit → return string
//
// APUs (Available Processing Units):
// - 1 function total
// - 1 public API (exported interface)
// - 0 helpers (self-contained logic)

// ────────────────────────────────────────────────────────────────
// Path Shortening - Display Optimization
// ────────────────────────────────────────────────────────────────
// What These Do:
// Transform verbose filesystem paths into concise display forms optimized
// for statusline space constraints while maintaining location clarity.
//
// Why Separated:
// Presentation logic isolated - this determines HOW to show path (shortened),
// statusline orchestrator determines WHERE to show it.
//
// Extension Point:
// To add new shortening strategies, create function following pattern:
//   func Shorten[Strategy]Path(path string) string
// Each strategy should return optimized path form for specific context.
//
// Pattern to follow:
//   1. Identify display constraint (screen width, context, etc.)
//   2. Implement shortening logic serving that constraint
//   3. Document reasoning and threshold values
//   4. Update tests to verify shortening behavior

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface.

// ShortenPath makes paths more readable for statusline display
//
// What It Does:
// Applies two-stage optimization to filesystem paths for display in statusline:
// 1. Replaces home directory prefix with ~ tilde (familiar Unix convention)
// 2. If result still exceeds MaxPathLength, compacts with ellipsis while
//    preserving head/tail context
//
// Parameters:
//   path: Full filesystem path (e.g., "/home/user/project/subdir/file.txt")
//
// Returns:
//   string: Shortened path optimized for display
//
// Health Impact:
//   Home directory retrieval: +25 points (success or graceful degradation)
//   Tilde substitution: +25 points (if applicable)
//   Length check: +25 points (decision made)
//   Path return: +25 points (result produced)
//   Total: +100 points per call
//
// Example usage:
//
//	short := format.ShortenPath("/home/user/very/long/nested/project/directory")
//	// If home = "/home/user":
//	//   Step 1: "~/very/long/nested/project/directory" (42 chars > 40)
//	//   Step 2: "~/very/long/n...project/directory" (head/tail retained)
//
//	short = format.ShortenPath("/home/user/project")
//	// If home = "/home/user":
//	//   Step 1: "~/project" (9 chars < 40)
//	//   Result: "~/project" (full shortened path)
//
//	short = format.ShortenPath("/etc/config/app.conf")
//	// Not in home directory:
//	//   Step 1: "/etc/config/app.conf" (20 chars < 40)
//	//   Result: "/etc/config/app.conf" (unchanged)
//
// Graceful degradation: If home directory unavailable (rare), path unchanged.
// Length threshold configurable via MaxPathLength constant.
func ShortenPath(path string) string {
	// Health: +25 points for home directory operation (success or graceful fail)
	// Replace home directory with ~
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		// Health: +25 points for successful tilde substitution
		path = "~" + strings.TrimPrefix(path, home)
	}
	// If home unavailable, continue with original path (graceful degradation)

	// Health: +25 points for length check and decision
	// Compact path with ellipsis if it exceeds threshold while keeping context
	if utf8.RuneCountInString(path) > MaxPathLength {
		return compactPath(path)
	}

	// Health: +25 points for returning shortened/original path
	return path
}

// compactPath shortens long paths by preserving both the beginning (home/root) and
// the most relevant ending (current directory) with an ellipsis separator.
func compactPath(path string) string {
	runes := []rune(path)
	total := len(runes)

	// If the path barely exceeds the threshold, trimming to basename still helps
	// without losing all context. Preserve at least one directory level.
	if total <= PathPrefixLength+PathSuffixLength {
		return filepath.Base(path)
	}

	prefixLen := PathPrefixLength
	suffixLen := PathSuffixLength

	// Ensure we do not slice beyond the available runes.
	if prefixLen > total {
		prefixLen = total
	}
	if suffixLen > total-prefixLen {
		suffixLen = total - prefixLen
	}

	prefix := string(runes[:prefixLen])
	suffix := string(runes[total-suffixLen:])
	return prefix + "..." + suffix
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
//   - Call ShortenPath with various path formats
//   - Verify home directory replacement with ~
//   - Verify basename extraction for long paths
//   - Test edge cases (empty string, root path, no home directory)
//   - Ensure no panics on special characters or unusual paths
//   - Confirm no go vet warnings introduced
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test with actual working directory paths from session context
//   - Verify shortened paths display correctly in statusline
//   - Confirm proper handling when home directory unavailable
//
// Example validation code:
//
//     // Test home replacement (assuming home = "/home/user")
//     expectExact := map[string]string{
//         "/home/user/project": "~/project",
//         "/etc/config":        "/etc/config", // unchanged
//         "":                   "",
//     }
//
//     for input, expected := range expectExact {
//         result := format.ShortenPath(input)
//         if result != expected {
//             t.Errorf("Input %q: expected %q, got %q", input, expected, result)
//         }
//     }
//
//     // Long path should include ellipsis while retaining tail context
//     longPath := "/home/user/very/long/path/that/exceeds/threshold"
//     longResult := format.ShortenPath(longPath)
//     if !strings.Contains(longResult, "...") {
//         t.Errorf("Expected ellipsis in long path result, got %q", longResult)
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
// Usage: import "statusline/lib/format"
//
// The library is imported into statusline orchestrator, making ShortenPath
// available for working directory display optimization. No code executes during
// import - function is defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "statusline/lib/format"
//
//     func buildStatusline(workdir string) string {
//         // Shorten working directory for display
//         workdirShort := format.ShortenPath(workdir)
//         parts := []string{workdirShort}
//
//         // ... build other parts ...
//
//         return strings.Join(parts, " | ")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No resources to manage (stateless pure functions)
//   - Single UserHomeDir() call (minimal I/O)
//   - String operations use Go's built-in memory management
//
// Graceful Shutdown:
//   - N/A for stateless library
//   - No lifecycle management needed
//   - Calling code has no cleanup responsibilities
//
// Error State Cleanup:
//   - UserHomeDir errors ignored gracefully (path returned unchanged)
//   - No partial state or rollback needed
//   - Always returns valid path string
//
// Memory Management:
//   - Go's garbage collector handles all memory
//   - Minimal allocations (string operations, single basename call if needed)
//   - No large buffers or persistent state

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
// Quick summary:
//   - Shortens filesystem paths for statusline display
//   - Replaces home directory with ~ tilde
//   - Shows basename only for paths exceeding 40 characters
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (LADDER presentation layer)
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new shortening strategies (follow Shorten[Strategy]Path pattern)
//   ✅ Adjust MaxPathLength constant (changes display threshold)
//   ✅ Extend logic for different path formats or conventions
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ ShortenPath signature - breaks statusline orchestrator
//   ⚠️ Return value semantics (what shortened means) - breaks display expectations
//   ⚠️ MaxPathLength drastic changes - affects layout consistency
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Graceful degradation guarantee (always return valid path)
//   ❌ Stateless design principle
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing requirements.
//   Always verify home replacement and length-based truncation after changes.
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart - Internal Structure" section above for
// complete ladder structure (dependencies) and baton flow (execution paths).
//
// Quick architectural summary:
// - 1 public API (ShortenPath) with self-contained logic
// - Ladder: Single-rung (no dependencies or helpers)
// - Baton: Linear flow (input → home replace → length check → return)
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Path Shortening" subsection above for detailed extension point.
//
// Quick reference:
// - Adding new shortening strategy: See BODY "Path Shortening" extension point
//   - Create Shorten[Strategy]Path function
//   - Implement logic for specific display constraint
//   - Document threshold and reasoning
//   - Update tests for strategy verification
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Quick summary:
// - UserHomeDir() call: Once per invocation (minimal I/O, cached by OS)
// - String operations: O(n) where n=path length (typically <100 chars)
// - Basename extraction: O(1) (filepath.Base is efficient)
// - Typical execution: <1 microsecond for standard paths
//
// Performance is negligible. UserHomeDir() typically cached, string ops minimal.
// Optimization not needed.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Expected Behavior:
//   - Home replacement: /home/user/... → ~/...
//   - Long path truncation: Shows basename only if >40 chars
//   - No home directory: Path unchanged (graceful degradation)
//   - Empty string: Returns empty string (valid behavior)
//
// If home replacement not working:
//   - Verify path actually starts with home directory
//   - Check UserHomeDir() returns expected value
//   - Confirm path uses absolute form (not relative)
//
// If truncation threshold seems wrong:
//   - Check MaxPathLength constant value (should be 40)
//   - Verify statusline has room for 40-char paths
//   - Consider adjusting MaxPathLength if layout changes
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Standard library only (os, path/filepath, strings)
// - Primary consumer: statusline orchestrator (statusline.go)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Home directory replacement with ~ - COMPLETED
//   ✓ Length-based basename truncation - COMPLETED
//   ⏳ Smart truncation (preserve key path segments)
//   ⏳ User-configurable length threshold
//   ⏳ Git repository root awareness (show relative to repo)
//
// Research Areas:
//   - Intelligent path abbreviation (show first/last segments)
//   - Context-sensitive truncation based on available statusline space
//   - Integration with workspace/project concepts
//
// Integration Targets:
//   - User preference system for path display customization
//   - Git library for repository-relative paths
//   - Dynamic length threshold based on terminal width
//
// Known Limitations to Address:
//   - Fixed 40-char threshold (not adaptive to terminal width)
//   - Simple basename truncation (loses context)
//   - No awareness of important path segments (like git repo root)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
//
//   1.0.0 (2025-11-04) - Full template application
//         - Comprehensive documentation added
//         - MaxPathLength constant
//         - Extension point documentation
//         - Health scoring map
//
//   0.1.0 (2024-10-24) - Initial implementation
//         - Home directory ~ replacement
//         - Length-based basename truncation
//         - Graceful degradation for missing home
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing presentation formatting for
// filesystem paths in statusline. It demonstrates the principle that clarity
// comes through intelligent brevity - showing location without overwhelming
// with full hierarchy.
//
// Modify thoughtfully - changes here affect how users identify where they're
// working. The balance between context and conciseness matters.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test home replacement and truncation thoroughly
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain graceful degradation guarantee
//
// "Make the path straight" - Matthew 3:3
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "statusline/lib/format"
//
//     // In statusline orchestrator
//     workdirShort := format.ShortenPath(workdir)
//     parts = append(parts, workdirShort)
//
// Home Replacement:
//
//     // Assuming home = "/home/user"
//     format.ShortenPath("/home/user/project")
//     // Returns: "~/project"
//
//     format.ShortenPath("/home/user/go/src/github.com/project")
//     // Returns: "~/go/src/github.com/project" (if <40 chars)
//
// Length-Based Truncation:
//
//     // Long path (>40 chars after home replacement)
//     format.ShortenPath("/home/user/very/long/nested/project/directory")
//     // Step 1: "~/very/long/nested/project/directory" (42 chars)
//     // Step 2: Returns "directory" (basename only)
//
// Non-Home Paths:
//
//     format.ShortenPath("/etc/nginx/sites-available")
//     // Returns: "/etc/nginx/sites-available" (unchanged, <40 chars)
//
//     format.ShortenPath("/var/log/application/server/production/access.log")
//     // Returns: "access.log" (>40 chars, basename only)
//
// ============================================================================
// END CLOSING
// ============================================================================
