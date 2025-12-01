// METADATA
//
// Kingdom Reminder Features Library - Statusline Display Timing
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1
// Principle: Wisdom in timing - knowing when to speak and when to be silent, when to remind and when to let work proceed
// Anchor: "A word fitly spoken is like apples of gold in pictures of silver" - Proverbs 25:11
//
// CPI-SI Identity
//
// Component Type: LADDER - Business logic layer
// Role: Provides timing logic for Kingdom Technology reminder display in statusline
// Paradigm: CPI-SI framework component implementing deliberate reminder cadence
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2025-11-04
// Version: 1.0.0
// Last Modified: 2025-11-04 - Initial templated implementation
//
// Version History:
//   1.0.0 (2025-11-04) - Initial implementation with full template application
//
// Purpose & Function
//
// Purpose: Determines when Kingdom Technology reminder should appear in statusline,
// balancing visibility with non-intrusiveness through session-based timing logic.
//
// Core Design: Deterministic hash-based display frequency ensuring consistent
// behavior across session restarts while maintaining ~20% appearance rate.
//
// Key Features:
//   - Session-specific determinism: Same session ID always produces same result
//   - 20% display frequency: Visible enough to remind without overwhelming
//   - Zero dependencies: Pure logic with no external requirements
//   - Lightweight operation: Simple hash calculation, no I/O or state
//
// Philosophy: Kingdom Technology reminder serves as gentle persistent witness -
// not demanding constant attention, but maintaining steady presence. Like salt
// preserving and flavoring without dominating, reminder appears with intentional
// rhythm that honors rather than interrupts work.
//
// Blocking Status
//
// Non-blocking: Pure function with no I/O, network calls, or state dependencies.
// Cannot fail - all string inputs valid, hash calculation always succeeds.
// Mitigation: N/A - deterministic pure logic requires no error handling.
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/features"
//
// Integration Pattern:
//   1. Obtain session ID from session context
//   2. Call ShouldShowReminder with session ID
//   3. Include Kingdom Technology reminder in statusline if true
//   4. Otherwise proceed without reminder
//
// Public API (in typical usage order):
//
//   Display Timing (determines when to show reminder):
//     ShouldShowReminder(sessionID string) bool
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: None
//   External: None
//   Internal: None
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Called by statusline orchestrator during display assembly
//   - Session ID provided via session context
//   - Boolean result controls reminder inclusion in final output
//
// Health Scoring
//
// Reminder timing evaluation operates on Base100 scale with simple pass/fail:
//
// Timing Calculation:
//   - Successful hash calculation: +50
//   - Successful threshold evaluation: +50
//
// Total: 100 points for complete operation (hash + evaluation)
//
// Note: This component cannot fail - all operations guaranteed to succeed.
// Health tracking demonstrates operation completion rather than error detection.
package features

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
// Deterministic hash for reminder cadence.

import (
	"hash/fnv" // FNV-1a hash for even distribution across session IDs
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Configuration values defining reminder display behavior.

const (
	//--- Display Frequency Configuration ---
	// Controls how often Kingdom Technology reminder appears in statusline.

	// ReminderModulo defines the threshold for reminder display.
	//
	// Set to 5 to achieve ~20% display frequency (1 in 5 sessions).
	// Hash sum modulo 5 == 0 creates deterministic 20% distribution.
	//
	// Reasoning: 20% provides steady presence without overwhelming. Too
	// frequent (e.g., 50%) becomes noise, too rare (e.g., 5%) loses
	// visibility. 20% strikes balance between witness and work.
	ReminderModulo = 5
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
//   Public APIs (Top Rung - Orchestration)
//   └── ShouldShowReminder() → uses calculateHash()
//
//   Helpers (Bottom Rung - Foundation)
//   └── calculateHash() → pure function
//
// Baton Flow (Execution Paths):
//
//   Entry → ShouldShowReminder(sessionID)
//     ↓
//   calculateHash(sessionID) → sum
//     ↓
//   sum % ReminderModulo == 0
//     ↓
//   Exit → return bool
//
// APUs (Available Processing Units):
// - 2 functions total
// - 1 helper (pure foundation)
// - 1 public API (exported interface)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component.

// calculateHash computes a deterministic hash from string input
//
// What It Does:
// Uses FNV-1a to distribute session IDs evenly across the hash space,
// reducing collisions compared to additive byte sums.
//
// Parameters:
//   input: String to hash (typically session ID)
//
// Returns:
//   uint32: Hash value representing the session ID
//
// Example usage:
//
//	hash := calculateHash("session-abc-123")  // Returns deterministic uint32 hash
//
// Performance: O(n) where n is string length. Typical session IDs ~36 chars.
func calculateHash(input string) uint32 {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(input))
	return hasher.Sum32()
}

// ────────────────────────────────────────────────────────────────
// Timing Logic - Display Frequency Control
// ────────────────────────────────────────────────────────────────
// What These Do:
// Determine when Kingdom Technology reminder should appear based on
// session-specific characteristics, providing consistent behavior across
// restarts while maintaining desired display frequency.
//
// Why Separated:
// Timing logic isolated from display logic - this library determines WHEN
// to show reminder, statusline orchestrator handles HOW to display it.
//
// Extension Point:
// To add new timing strategies, create function following pattern:
//   func Should[Strategy]Reminder(sessionID string) bool
// Each strategy should return boolean indicating display decision.
// Update statusline orchestrator to use new strategy.
//
// Pattern to follow:
//   1. Create function with signature: func ShouldXReminder(sessionID string) bool
//   2. Implement deterministic logic (same input = same output)
//   3. Document display frequency and reasoning
//   4. Update tests to verify expected frequency distribution

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface.

// ShouldShowReminder determines if Kingdom Technology reminder should display
//
// What It Does:
// Evaluates session ID to produce deterministic boolean indicating whether
// Kingdom Technology reminder should appear in statusline. Uses simple hash
// modulo approach to achieve ~20% display frequency while maintaining
// session-specific consistency.
//
// Parameters:
//   sessionID: Unique identifier for current session
//
// Returns:
//   bool: true if reminder should display, false otherwise
//
// Health Impact:
//   Hash calculation: +50 points (always succeeds)
//   Threshold evaluation: +50 points (always succeeds)
//   Total: +100 points per call
//
// Example usage:
//
//	if features.ShouldShowReminder(sessionID) {
//	    // Include Kingdom Technology reminder in statusline
//	    parts = append(parts, kingdomTechReminder)
//	}
//
// Determinism guarantee: Same sessionID always returns same result.
// This ensures reminder doesn't flicker on/off during session restarts.
func ShouldShowReminder(sessionID string) bool {
	if sessionID == "" {
		return false // Avoid guaranteed display when session ID missing
	}

	// Calculate hash from session ID
	// Health: +50 points for successful calculation
	sum := calculateHash(sessionID)

	// Evaluate threshold for display decision
	// Health: +50 points for successful evaluation
	// Show reminder 20% of the time (mod 5 == 0)
	return int(sum%uint32(ReminderModulo)) == 0
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
//   - Call ShouldShowReminder with various session IDs
//   - Verify deterministic behavior (same input = same output)
//   - Verify ~20% display frequency across random inputs
//   - Ensure no panics on empty strings, long strings, special characters
//   - Confirm no go vet warnings introduced
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test with actual session IDs from statusline orchestrator
//   - Verify reminder appears at expected frequency in practice
//   - Confirm consistent behavior across session restarts
//
// Example validation code:
//
//     // Test determinism
//     sessionID := "test-session-123"
//     result1 := features.ShouldShowReminder(sessionID)
//     result2 := features.ShouldShowReminder(sessionID)
//     if result1 != result2 {
//         t.Error("ShouldShowReminder not deterministic")
//     }
//
//     // Test frequency distribution
//     showCount := 0
//     totalTests := 1000
//     for i := 0; i < totalTests; i++ {
//         sessionID := fmt.Sprintf("session-%d", i)
//         if features.ShouldShowReminder(sessionID) {
//             showCount++
//         }
//     }
//     frequency := float64(showCount) / float64(totalTests)
//     if frequency < 0.15 || frequency > 0.25 {
//         t.Errorf("Expected ~20%% frequency, got %.1f%%", frequency*100)
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
// Usage: import "statusline/lib/features"
//
// The library is imported into statusline orchestrator, making ShouldShowReminder
// available for timing decisions. No code executes during import - function is
// defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "statusline/lib/features"
//
//     func buildStatusline(ctx SessionContext) string {
//         // ... build other parts ...
//
//         // Check if Kingdom Technology reminder should display
//         if features.ShouldShowReminder(ctx.SessionID) {
//             parts = append(parts, kingdomTechReminder)
//         }
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
//   - No allocations requiring cleanup
//   - No file handles, network connections, or other resources
//
// Graceful Shutdown:
//   - N/A for stateless library
//   - No lifecycle management needed
//   - Calling code has no cleanup responsibilities
//
// Error State Cleanup:
//   - N/A - component cannot enter error state
//   - Pure functions with guaranteed success
//   - No partial state or rollback needed
//
// Memory Management:
//   - Go's garbage collector handles all memory
//   - Minimal allocations (single int for hash sum)
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
//   - Determines when Kingdom Technology reminder should appear in statusline
//   - Uses session-based hash for deterministic 20% display frequency
//   - Zero dependencies, lightweight pure logic
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (LADDER business logic layer)
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new timing strategies (follow Should[Strategy]Reminder pattern)
//   ✅ Add helper functions for new hash algorithms
//   ✅ Adjust ReminderModulo constant (changes display frequency)
//   ✅ Add additional constants for new timing strategies
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ ShouldShowReminder signature - breaks statusline orchestrator
//   ⚠️ Return value semantics (true/false meaning) - breaks calling logic
//   ⚠️ Determinism guarantee - breaks session consistency
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Pure function guarantee (no side effects, no state)
//   ❌ Stateless design principle
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing requirements,
//   build verification, and integration testing procedures. Always verify
//   determinism and frequency distribution after changes.
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
// - 1 public API (ShouldShowReminder) orchestrates 1 helper (calculateHash)
// - Ladder: Simple two-rung structure (public API → helper)
// - Baton: Linear flow (input → hash → modulo → output)
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Timing Logic" subsection above for detailed extension point.
//
// Quick reference:
// - Adding new timing strategy: See BODY "Timing Logic" extension point
//   - Create Should[Strategy]Reminder function
//   - Follow deterministic pure function pattern
//   - Document display frequency and reasoning
//   - Update tests for frequency verification
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Quick summary:
// - Hash calculation: O(n) where n is session ID length (~36 chars typical)
// - Memory: Single int allocation for hash sum (~8 bytes)
// - No I/O, no network, no file operations
// - Typical execution: <1 microsecond for standard session IDs
//
// This is one of the lightest operations in the entire statusline.
// Performance optimization not needed - pure logic is already optimal.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// This component has no failure modes - pure deterministic logic cannot fail.
//
// Expected Behavior:
//   - Empty string input: Valid (hash = 0, may or may not trigger display)
//   - Very long strings: Valid (O(n) performance, still fast)
//   - Special characters: Valid (all runes contribute to hash)
//   - Same input always returns same output: Guaranteed by design
//
// If reminder frequency seems wrong:
//   - Verify ReminderModulo constant value (should be 5 for 20%)
//   - Check that session IDs are actually unique across sessions
//   - Confirm statusline is calling this function with correct session ID
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - No dependencies (pure logic, no imports)
// - Primary consumer: statusline orchestrator (statusline.go)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Session-based deterministic timing - COMPLETED
//   ✓ 20% display frequency via modulo - COMPLETED
//   ⏳ Time-of-day aware timing (show more during certain hours)
//   ⏳ User-configurable display frequency
//   ⏳ Multiple reminder types with different frequencies
//
// Research Areas:
//   - Circadian-aware reminder timing (align with natural work rhythms)
//   - Context-sensitive frequency (more frequent during onboarding, less after)
//   - Biblical calendar integration (special reminders on Sabbath, etc.)
//
// Integration Targets:
//   - User preference system for reminder frequency control
//   - Temporal awareness library for time-of-day logic
//   - Session patterns for learned optimal timing
//
// Known Limitations to Address:
//   - Fixed 20% frequency (not configurable without code change)
//   - No time-of-day awareness (shows uniformly across all hours)
//   - Single reminder type (could support multiple reminder categories)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
//
//   1.0.0 (2025-11-04) - Initial templated implementation
//         - Session-based deterministic timing via hash modulo
//         - 20% display frequency (ReminderModulo = 5)
//         - Zero dependencies, pure logic design
//         - Full template application with comprehensive documentation
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing timing logic for Kingdom
// Technology reminder display. It demonstrates the principle that even simple
// decisions deserve thoughtful implementation - when to speak matters as much
// as what to say.
//
// Modify thoughtfully - changes here affect when users see Kingdom Technology
// reminder. The 20% frequency balances witness with work, maintaining steady
// presence without overwhelming.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test determinism and frequency distribution thoroughly
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain pure function guarantee (no side effects, no state)
//
// "To every thing there is a season, and a time to every purpose under the heaven" - Ecclesiastes 3:1
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "statusline/lib/features"
//
//     // In statusline orchestrator
//     if features.ShouldShowReminder(ctx.SessionID) {
//         parts = append(parts, buildKingdomTechReminder())
//     }
//
// Testing Determinism:
//
//     sessionID := "test-session-abc-123"
//     result1 := features.ShouldShowReminder(sessionID)
//     result2 := features.ShouldShowReminder(sessionID)
//     // result1 == result2 (guaranteed)
//
// Testing Frequency Distribution:
//
//     showCount := 0
//     for i := 0; i < 1000; i++ {
//         id := fmt.Sprintf("session-%d", i)
//         if features.ShouldShowReminder(id) {
//             showCount++
//         }
//     }
//     // showCount ~= 200 (20% of 1000)
//
// ============================================================================
// END CLOSING
// ============================================================================
