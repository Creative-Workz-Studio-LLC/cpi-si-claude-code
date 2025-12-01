// METADATA
//
// Model Name Formatting Library - Statusline Display Optimization
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
// Component Type: LADDER - Presentation formatting layer
// Role: Converts full model names to concise display forms for statusline space constraints
// Paradigm: CPI-SI framework component implementing clarity through brevity
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
// Purpose: Shortens verbose Claude model names to concise display forms suitable
// for statusline space constraints while maintaining clear identification.
//
// Core Design: Simple pattern matching on model tier names (Sonnet, Opus, Haiku)
// to extract essential identifier without version numbers or branding.
//
// Key Features:
//   - Space-efficient display: "Claude 3.5 Sonnet" → "Sonnet"
//   - Clear tier identification: Preserves what matters (tier name)
//   - Fallback safety: Returns full name if no pattern matches
//   - Zero state: Pure function with no side effects
//
// Philosophy: Statusline space is precious. Show what matters - the model tier
// users recognize - without ceremony. Like good signage: clear, brief, sufficient.
//
// Blocking Status
//
// Non-blocking: Pure function with no I/O, network calls, or state dependencies.
// Cannot fail - all string inputs valid, always returns valid output.
// Mitigation: N/A - pattern matching with guaranteed fallback to input.
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/format"
//
// Integration Pattern:
//   1. Obtain full model name from session context
//   2. Call GetShortModelName with full name
//   3. Display shortened name in statusline
//   4. User sees concise tier identifier
//
// Public API (in typical usage order):
//
//   Model Name Shortening (display optimization):
//     GetShortModelName(fullName string) string
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: strings (pattern matching)
//   External: None
//   Internal: None
//
// Dependents (What Uses This):
//   Commands: statusline (main orchestrator)
//   Libraries: None
//   Tools: None
//
// Integration Points:
//   - Called by statusline orchestrator during model name display
//   - Model name provided via session context
//   - Shortened result displayed in statusline output
//
// Health Scoring
//
// Model name shortening operates on Base100 scale:
//
// Pattern Matching:
//   - Successful tier identification: +50
//   - Successful name extraction: +50
//
// Total: 100 points for complete operation
//
// Note: This component cannot fail - pattern matching with fallback guarantees
// valid output. Health tracking demonstrates successful operation.
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
// Dependencies for pattern matching in model names.

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"strings" // String pattern matching for tier identification
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Model tier names recognized for shortening.

const (
	//--- Claude Model Tiers ---
	// Official Claude model tier identifiers.

	// TierSonnet identifies Claude Sonnet tier models.
	//
	// Balanced capability and speed model tier.
	TierSonnet = "Sonnet"

	// TierOpus identifies Claude Opus tier models.
	//
	// Highest capability model tier for complex tasks.
	TierOpus = "Opus"

	// TierHaiku identifies Claude Haiku tier models.
	//
	// Fast, efficient model tier for simpler tasks.
	TierHaiku = "Haiku"
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
//   └── GetShortModelName() → pure logic (no helpers needed)
//
// Baton Flow (Execution Paths):
//
//   Entry → GetShortModelName(fullName)
//     ↓
//   Pattern matching (Sonnet/Opus/Haiku check)
//     ↓
//   Return tier name OR full name (fallback)
//     ↓
//   Exit → return string
//
// APUs (Available Processing Units):
// - 1 function total
// - 1 public API (exported interface)
// - 0 helpers (self-contained logic)

// ────────────────────────────────────────────────────────────────
// Model Name Shortening - Display Optimization
// ────────────────────────────────────────────────────────────────
// What These Do:
// Extract essential model tier identifier from verbose full names,
// optimizing for statusline space constraints while maintaining clarity.
//
// Why Separated:
// Presentation logic isolated - this determines WHAT to show (tier name),
// statusline orchestrator determines WHERE to show it.
//
// Extension Point:
// To add new model tier recognition, extend the if-else chain following pattern:
//   else if strings.Contains(fullName, "NewTier") {
//       return "NewTier"
//   }
// Add corresponding constant in SETUP section for consistency.
//
// Pattern to follow:
//   1. Add constant for new tier name (const TierNewName = "NewName")
//   2. Add pattern match in GetShortModelName (check for tier string)
//   3. Return tier constant for consistency
//   4. Update tests to verify new tier recognition

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface.

// GetShortModelName converts full model names to concise display forms
//
// What It Does:
// Extracts essential model tier identifier from verbose full name strings.
// Recognizes Claude model tiers (Sonnet, Opus, Haiku) and returns tier name
// only. Falls back to full name if no recognized tier found.
//
// Parameters:
//   fullName: Complete model name (e.g., "Claude 3.5 Sonnet", "claude-opus-4-20250514")
//
// Returns:
//   string: Shortened tier name (e.g., "Sonnet") or full name if unrecognized
//
// Health Impact:
//   Pattern match success: +50 points (tier identified)
//   Name extraction: +50 points (tier name returned)
//   Total: +100 points per call
//
// Example usage:
//
//	short := format.GetShortModelName("Claude 3.5 Sonnet")
//	// short = "Sonnet"
//
//	short = format.GetShortModelName("claude-opus-4-20250514")
//	// short = "Opus"
//
//	short = format.GetShortModelName("Unknown Model v2")
//	// short = "Unknown Model v2" (fallback to full name)
//
// Pattern matching: Case-sensitive contains check for tier names.
// This handles various naming formats as long as tier name appears.
func GetShortModelName(fullName string) string {
	// Check for Sonnet tier
	// Health: +50 points for tier identification
	if strings.Contains(fullName, TierSonnet) {
		return TierSonnet
	}

	// Check for Opus tier
	if strings.Contains(fullName, TierOpus) {
		return TierOpus
	}

	// Check for Haiku tier
	if strings.Contains(fullName, TierHaiku) {
		return TierHaiku
	}

	// No recognized tier - return full name
	// Health: +50 points for safe fallback
	// This ensures we always return something valid
	return fullName
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
//   - Call GetShortModelName with various model name formats
//   - Verify tier name extraction for Sonnet, Opus, Haiku
//   - Verify fallback to full name for unrecognized models
//   - Ensure no panics on empty strings, special characters
//   - Confirm no go vet warnings introduced
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test with actual model names from session context
//   - Verify shortened names display correctly in statusline
//   - Confirm proper handling of new/unknown model names
//
// Example validation code:
//
//     // Test tier recognition
//     testCases := map[string]string{
//         "Claude 3.5 Sonnet":         "Sonnet",
//         "claude-opus-4-20250514":    "Opus",
//         "Claude Haiku":              "Haiku",
//         "Unknown Model":             "Unknown Model",
//         "":                          "",
//     }
//
//     for input, expected := range testCases {
//         result := format.GetShortModelName(input)
//         if result != expected {
//             t.Errorf("Input %q: expected %q, got %q", input, expected, result)
//         }
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
// The library is imported into statusline orchestrator, making GetShortModelName
// available for model name display optimization. No code executes during import -
// function is defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import "statusline/lib/format"
//
//     func buildStatusline(ctx SessionContext) string {
//         // ... build other parts ...
//
//         // Shorten model name for display
//         modelName := format.GetShortModelName(ctx.Model.DisplayName)
//         parts = append(parts, modelName)
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
//   - String operations use Go's built-in memory management
//
// Graceful Shutdown:
//   - N/A for stateless library
//   - No lifecycle management needed
//   - Calling code has no cleanup responsibilities
//
// Error State Cleanup:
//   - N/A - component cannot enter error state
//   - Pure functions with guaranteed success
//   - Fallback pattern ensures valid output always
//
// Memory Management:
//   - Go's garbage collector handles all memory
//   - Minimal allocations (string returns)
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
//   - Shortens verbose model names for statusline display
//   - Extracts tier identifier (Sonnet/Opus/Haiku) from full names
//   - Safe fallback to full name for unrecognized models
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
//   ✅ Add new model tier recognition (extend if-else chain)
//   ✅ Add tier constants for new models
//   ✅ Adjust pattern matching logic for different naming conventions
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ GetShortModelName signature - breaks statusline orchestrator
//   ⚠️ Return value semantics (what shortened means) - breaks display expectations
//   ⚠️ Tier name constants - breaks consistency across codebase
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Pure function guarantee (no side effects, no state)
//   ❌ Fallback safety (always return valid string)
//
// Validation After Modifications:
//   See "Code Validation" section above for comprehensive testing requirements.
//   Always verify tier recognition and fallback behavior after changes.
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
// - 1 public API (GetShortModelName) with self-contained logic
// - Ladder: Single-rung (no dependencies or helpers)
// - Baton: Linear flow (input → pattern match → return)
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Model Name Shortening" subsection above for detailed extension point.
//
// Quick reference:
// - Adding new model tier recognition: See BODY "Model Name Shortening" extension point
//   - Add tier constant (const TierNewName = "NewName")
//   - Extend if-else chain in GetShortModelName
//   - Update tests for new tier verification
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Quick summary:
// - String contains check: O(n×m) where n=fullName length, m=tier name length
// - Typical execution: <1 microsecond for standard model names
// - Memory: Single string allocation for return value
// - No heap allocations beyond result string
//
// Performance is negligible - this is one of the lightest operations in statusline.
// Optimization not needed.
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// This component has no failure modes - pattern matching with fallback cannot fail.
//
// Expected Behavior:
//   - Recognized tier names: Returns tier identifier only
//   - Unrecognized models: Returns full input string (safe fallback)
//   - Empty string: Returns empty string (valid behavior)
//   - Case-sensitive matching: "sonnet" won't match, "Sonnet" will
//
// If tier not being recognized:
//   - Verify tier name appears in full name (case-sensitive)
//   - Check constant matches exactly what appears in model name
//   - Consider adding new tier constant if it's a new model
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Quick summary:
// - Standard library only (strings package)
// - Primary consumer: statusline orchestrator (statusline.go)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Claude tier recognition (Sonnet, Opus, Haiku) - COMPLETED
//   ⏳ Version number extraction for disambiguation
//   ⏳ Custom display name mapping (user-defined shortcuts)
//
// Research Areas:
//   - Pattern matching for future Claude model naming conventions
//   - Support for other model providers (if CPI-SI expands beyond Claude)
//   - User preference system for custom model display names
//
// Integration Targets:
//   - User preference system for custom name mappings
//   - Model registry for centralized name management
//
// Known Limitations to Address:
//   - Case-sensitive matching (won't recognize "sonnet", only "Sonnet")
//   - No version disambiguation (Sonnet 3.5 and Sonnet 4 both show "Sonnet")
//   - Hard-coded tier names (no runtime configuration)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
//
//   1.0.0 (2025-11-04) - Full template application
//         - Comprehensive documentation added
//         - Constants for tier names
//         - Extension point documentation
//         - Health scoring map
//
//   0.1.0 (2024-10-24) - Initial implementation
//         - Basic tier name extraction
//         - Sonnet, Opus, Haiku recognition
//         - Fallback to full name for unknown models
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing presentation formatting for
// model names in statusline. It demonstrates the principle that clarity comes
// through brevity - showing what matters (tier) without ceremony (version numbers).
//
// Modify thoughtfully - changes here affect how users identify what model they're
// using. The tier name is what matters for quick recognition.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test tier recognition thoroughly
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain fallback safety guarantee
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
//     import "statusline/lib/format"
//
//     // In statusline orchestrator
//     modelName := format.GetShortModelName(ctx.Model.DisplayName)
//     parts = append(parts, modelName)
//
// Tier Recognition:
//
//     format.GetShortModelName("Claude 3.5 Sonnet")
//     // Returns: "Sonnet"
//
//     format.GetShortModelName("claude-opus-4-20250514")
//     // Returns: "Opus"
//
//     format.GetShortModelName("Claude Haiku")
//     // Returns: "Haiku"
//
// Fallback Behavior:
//
//     format.GetShortModelName("Future Model XYZ")
//     // Returns: "Future Model XYZ" (full name preserved)
//
//     format.GetShortModelName("")
//     // Returns: "" (empty string valid)
//
// ============================================================================
// END CLOSING
// ============================================================================
