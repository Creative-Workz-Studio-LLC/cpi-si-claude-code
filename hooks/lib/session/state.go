// ============================================================================
// METADATA
// ============================================================================
// Session State Library - CPI-SI Hooks System
//
// For METADATA structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Be watchful, stand firm in the faith" - 1 Corinthians 16:13 (WEB)
// Principle: Awareness of state enables faithful continuation - tracking session
//            context provides continuity and coherence across work sessions
// Anchor: "A wise man's heart discerns both time and judgment" - Ecclesiastes 8:5 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Orchestration wrapper (hooks-specific rung)
// Role: Provides hooks-layer access to session state operations while delegating
//       to system/lib/sessiontime (authoritative source). Maintains hooks-compatible
//       interface for backward compatibility.
// Paradigm: CPI-SI framework component - hooks orchestrate, system implements
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2025-11-10
// Version: 3.0.0
// Last Modified: 2025-11-12 - Removed duplication, imports from system/lib/sessiontime
//
// Version History:
//   3.0.0 (2025-11-12) - Architectural consolidation, removed duplicate SessionState
//   2.0.0 (2025-11-10) - Config inheritance, richer structure, correct paths
//   1.0.0 (2024-10-24) - Initial implementation with basic session state
//
// Purpose & Function
//
// Purpose: Provide hooks-layer orchestration wrapper around session state operations.
// Delegates to system/lib/sessiontime (authoritative implementation) while maintaining
// backward-compatible interface for existing hooks. Ensures hooks can access session
// state without tight coupling to system internals.
//
// Core Design: Thin wrapper pattern - re-exports types and delegates function calls
// to authoritative system library. No business logic duplication, only interface
// adaptation for hooks ecosystem compatibility.
//
// Key Features:
//   - Re-exports SessionState type from system/lib/sessiontime
//   - Delegates compaction count operations to system library
//   - Delegates session state retrieval to system library
//   - Maintains backward-compatible function signatures
//   - Zero duplication - all logic in system/lib/sessiontime
//
// Philosophy: "Hooks orchestrate, system implements" - this library provides the hooks
// interface layer while system provides implementation. Changes to SessionState structure
// happen in system/lib/sessiontime, this wrapper automatically inherits them.
//
// Blocking Status
//
// Non-blocking: All operations delegate to system/lib/sessiontime which handles errors.
// This wrapper propagates errors from system library without additional blocking behavior.
//
// Mitigation: Errors returned for caller handling. No file I/O or blocking operations
// in wrapper itself - all actual work happens in system library.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
//	func sessionHook() {
//		// Increment compaction count
//		count, err := session.IncrementCompactionCount()
//
//		// Get session state
//		state, err := session.GetSessionState()
//
//		// Access SessionState fields
//		fmt.Println(state.SessionID)
//		fmt.Println(state.CompactionCount)
//	}
//
// Integration Pattern:
//   1. Import hooks/lib/session
//   2. Call wrapper functions (IncrementCompactionCount, GetCompactionCount, GetSessionState)
//   3. Wrapper delegates to system/lib/sessiontime
//   4. No cleanup needed - stateless wrapper
//
// Public API (in typical usage order):
//
//   Compaction Operations:
//     IncrementCompactionCount() (int, error) - Increment and return compaction count
//     GetCompactionCount() (int, error) - Get current compaction count
//
//   State Access:
//     GetSessionState() (*SessionState, error) - Get complete session state
//
//   Types:
//     SessionState - Re-exported from system/lib/sessiontime
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: None (pure delegation wrapper)
//   External: None
//   Internal: system/lib/sessiontime (authoritative session state library)
//
// Dependents (What Uses This):
//   Commands: hooks/session/cmd-start, hooks/session/cmd-end, hooks/session/cmd-pre-compact
//   Libraries: None (leaf node for hooks ecosystem)
//
// Integration Points:
//   - Rails: No logger needed (wrapper delegates to system)
//   - Baton: Receives calls from hooks, passes to system library
//   - Ladder: Depends on system/lib/sessiontime (lower rung)
//
// Health Scoring
//
// This wrapper delegates to system/lib/sessiontime for all operations. Health tracking
// happens in system library - wrapper success measured by delegation correctness:
//
// Delegation Success:
//   - Successful call to system library: +80
//   - System library returns expected result: +20
//
// Note: Actual session state operation health tracked in system/lib/sessiontime.
// This wrapper's health reflects delegation correctness only.
package session

// ============================================================================
// END METADATA
// ============================================================================

// ============================================================================
// SETUP
// ============================================================================
//
// For SETUP structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-006-CODE-setup-block.md

// ────────────────────────────────────────────────────────────────
// Imports - Dependencies
// ────────────────────────────────────────────────────────────────
// Dependencies this component needs. This wrapper has minimal dependencies -
// only imports the authoritative system library it delegates to. No standard
// library imports needed for pure delegation wrapper.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// None needed - pure delegation wrapper

	//--- External Packages ---
	// None needed - pure delegation wrapper

	//--- Internal Packages ---
	// Authoritative session state library providing actual implementation.

	"system/lib/sessiontime" // Session state operations (authoritative source)
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// None needed for pure delegation wrapper. All constants live in system/lib/sessiontime.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types re-exported from system library for hooks compatibility.
// This maintains backward compatibility while delegating to authoritative source.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

// SessionState is re-exported from system/lib/sessiontime for convenience.
//
// This maintains backward compatibility for hooks code that imports
// "hooks/lib/session" and expects SessionState to be available. All actual
// type definition and implementation lives in system/lib/sessiontime.
//
// IMPORTANT: This is a type alias, not a duplicate definition. Changes to
// SessionState happen in system/lib/sessiontime and automatically propagate
// to this wrapper.
//
// Usage:
//
//     import "hooks/lib/session"
//
//     state, err := session.GetSessionState()
//     // state is *session.SessionState (which is *sessiontime.SessionState)
//     fmt.Println(state.SessionID)
//
type SessionState = sessiontime.SessionState

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// None needed for pure delegation wrapper. This library has no logger or
// inspector because it performs no operations itself - all work delegated
// to system/lib/sessiontime which has its own Rails infrastructure.
//
// See: ~/.claude/cpi-si/docs/standards/code/patterns/CWS-PATTERN-003-CODE-rails.md
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md
//
// Note: Pure delegation wrappers skip Rails infrastructure. The authoritative
// implementation (system/lib/sessiontime) handles health tracking and debugging.

// ============================================================================
// END SETUP
// ============================================================================

// ============================================================================
// BODY
// ============================================================================
//
// For BODY structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-007-CODE-body-block.md

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Internal Structure
// ────────────────────────────────────────────────────────────────
// Maps bidirectional dependencies and baton flow within this component.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-organizational-chart.md
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Pure Delegation)
//   ├── IncrementCompactionCount() → delegates to sessiontime.IncrementCompactionCount()
//   ├── GetCompactionCount() → delegates to sessiontime.GetCompactionCount()
//   └── GetSessionState() → delegates to sessiontime.ReadSession()
//
//   Core Operations: None (pure delegation wrapper)
//   Helpers: None (pure delegation wrapper)
//
// Baton Flow (Execution Paths):
//
//   Entry → [Public API Function]
//     ↓
//   Delegate to system/lib/sessiontime
//     ↓
//   Return result from system library
//     ↓
//   Exit → return to caller
//
// APUs (Available Processing Units):
// - 3 functions total
// - 0 helpers (pure delegation)
// - 0 core operations (pure delegation)
// - 3 public APIs (exported delegation wrappers)
//
// Note: All actual processing happens in system/lib/sessiontime (authoritative source).
// This wrapper provides hooks-compatible interface only.

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// None needed for pure delegation wrapper. All helper functions live in
// system/lib/sessiontime (authoritative implementation).
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// ────────────────────────────────────────────────────────────────
// Core Operations - Business Logic
// ────────────────────────────────────────────────────────────────
// None in this wrapper. All business logic lives in system/lib/sessiontime
// (authoritative implementation). This wrapper delegates calls without
// performing operations itself.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-core-operations.md

// ────────────────────────────────────────────────────────────────
// Error Handling/Recovery Patterns
// ────────────────────────────────────────────────────────────────
// None needed for pure delegation wrapper. Errors from system/lib/sessiontime
// are propagated unchanged to caller. No panic recovery, error wrapping, or
// retry logic in wrapper - all error handling happens in authoritative implementation.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-error-handling.md

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining hooks-compatible interface. All functions
// delegate to system/lib/sessiontime for actual implementation. Simple by
// design - complexity lives in system library, wrapper provides clean interface.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-public-apis.md

// IncrementCompactionCount increments compaction count and returns new value.
//
// What It Does:
// Delegates to system/lib/sessiontime.IncrementCompactionCount() which reads
// current session state, increments CompactionCount field, writes updated state
// back to file, and returns new count.
//
// Parameters: None
//
// Returns:
//   int: New compaction count after increment
//   error: Error from system library (file read/write/JSON errors)
//
// Health Impact:
//   Delegation success: +80 points (call succeeded)
//   Expected result: +20 points (got valid count)
//   System library failure: -50 points (propagated error)
//
// Example usage:
//
//	count, err := session.IncrementCompactionCount()
//	if err != nil {
//	    log.Printf("Failed to increment compaction count: %v", err)
//	    return
//	}
//	fmt.Printf("Compaction count: %d\n", count)
//
func IncrementCompactionCount() (int, error) {
	return sessiontime.IncrementCompactionCount()
}

// GetCompactionCount returns current compaction count from session state.
//
// What It Does:
// Delegates to system/lib/sessiontime.GetCompactionCount() which reads current
// session state and extracts CompactionCount field.
//
// Parameters: None
//
// Returns:
//   int: Current compaction count
//   error: Error from system library (file read/JSON parse errors)
//
// Health Impact:
//   Delegation success: +80 points (call succeeded)
//   Expected result: +20 points (got valid count)
//   System library failure: -50 points (propagated error)
//
// Example usage:
//
//	count, err := session.GetCompactionCount()
//	if err != nil {
//	    log.Printf("Failed to get compaction count: %v", err)
//	    return
//	}
//	fmt.Printf("Current compaction count: %d\n", count)
//
func GetCompactionCount() (int, error) {
	return sessiontime.GetCompactionCount()
}

// GetSessionState returns the current session state.
//
// What It Does:
// Delegates to system/lib/sessiontime.ReadSession() which reads session state
// file, parses JSON, and returns complete SessionState struct.
//
// Parameters: None
//
// Returns:
//   *SessionState: Pointer to current session state with all fields
//   error: Error from system library (file read/JSON parse errors)
//
// Health Impact:
//   Delegation success: +80 points (call succeeded)
//   Expected result: +20 points (got valid state)
//   System library failure: -50 points (propagated error)
//
// Example usage:
//
//	state, err := session.GetSessionState()
//	if err != nil {
//	    log.Printf("Failed to get session state: %v", err)
//	    return
//	}
//	fmt.Printf("Session ID: %s\n", state.SessionID)
//	fmt.Printf("Started: %s\n", state.StartFormatted)
//	fmt.Printf("Compactions: %d\n", state.CompactionCount)
//
func GetSessionState() (*SessionState, error) {
	return sessiontime.ReadSession()
}

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================
//
// For CLOSING structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-008-CODE-closing-block.md
//
// ────────────────────────────────────────────────────────────────
// Code Validation: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Validation section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-010-CLOSING-code-validation.md
//
// Testing Requirements:
//   - Import the library without errors
//   - Call each wrapper function with no parameters
//   - Verify delegation to system/lib/sessiontime works
//   - Check error propagation from system library
//   - Ensure SessionState type alias resolves correctly
//   - Confirm no go vet warnings introduced
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - cd ~/.claude/hooks/lib/session && go build .
//   - go vet . (no warnings)
//   - Verify system/lib/sessiontime import resolves
//
// Integration Testing:
//   - Test with actual hooks (cmd-start, cmd-end, cmd-pre-compact)
//   - Verify compaction count increments correctly
//   - Check session state retrieval returns expected data
//   - Validate backward compatibility with existing hooks code
//
// Example validation code:
//
//     // Test delegation works
//     count, err := session.IncrementCompactionCount()
//     if err != nil {
//         t.Errorf("IncrementCompactionCount failed: %v", err)
//     }
//
//     state, err := session.GetSessionState()
//     if err != nil {
//         t.Errorf("GetSessionState failed: %v", err)
//     }
//     if state.CompactionCount != count {
//         t.Errorf("Compaction count mismatch")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. No entry point, no main function.
// All functions wait to be called by hooks orchestration code.
//
// Usage: import "hooks/lib/session"
//
// The library is imported into hooks code, making wrapper functions available.
// No code executes during import - functions are defined and ready to delegate.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func main() {
//         // Call wrapper functions
//         state, err := session.GetSessionState()
//         if err != nil {
//             log.Fatal(err)
//         }
//         fmt.Println(state.SessionID)
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No resources managed by wrapper (pure delegation)
//   - System library handles file I/O and cleanup
//   - Wrapper is stateless - no cleanup needed
//
// Graceful Shutdown:
//   - N/A for stateless wrapper
//   - No lifecycle to manage
//   - System library handles resource cleanup
//
// Error State Cleanup:
//   - Errors propagated from system library unchanged
//   - No partial state in wrapper
//   - No rollback needed (stateless)
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - No manual allocations in wrapper
//   - SessionState memory managed by system library
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// For Library Overview section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-013-CLOSING-library-overview.md
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: Hooks-layer orchestration wrapper around session state operations
//
// Quick summary:
//   - Thin wrapper delegating to system/lib/sessiontime (authoritative source)
//   - Re-exports SessionState type for backward compatibility
//   - Zero duplication - all logic in system library
//
// Integration Pattern: See METADATA "Usage & Integration" section above
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list organized by category
//
// Architecture: Hooks orchestration wrapper - delegates to system implementation
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new wrapper functions (follow delegation pattern)
//   ✅ Add convenience methods that delegate to system library
//   ✅ Extend Public API with hooks-specific helpers
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Function signatures - breaks all calling hooks
//   ⚠️ SessionState type alias - breaks hooks accessing state
//   ⚠️ Delegation targets - affects all consumers
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Delegation pattern - wrapper MUST delegate to system library
//   ❌ Zero duplication principle - no business logic in wrapper
//   ❌ SessionState definition location - lives in system/lib/sessiontime only
//
// Validation After Modifications:
//   See "Code Validation" section above for testing requirements
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart" section above for structure.
//
// Quick architectural summary:
// - 3 public wrapper APIs delegate to system/lib/sessiontime
// - Ladder: hooks/lib/session → system/lib/sessiontime
// - Baton: Call flows through wrapper to system library, result flows back
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// Adding New Wrapper Functions:
//   1. Add function in BODY "Public APIs" section
//   2. Follow delegation pattern: return systemtime.FunctionName(...)
//   3. Document with full docstring (What/Parameters/Returns/Health/Example)
//   4. Update BODY Organizational Chart
//   5. Update METADATA Public API list
//
// Example new wrapper:
//
//     // NewWrapperFunction does X by delegating to system library.
//     func NewWrapperFunction(param type) (result, error) {
//         return sessiontime.SystemFunction(param)
//     }
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Wrapper Overhead:
//   - Function call overhead: Negligible (single delegation)
//   - Memory overhead: None (stateless wrapper)
//   - Type alias overhead: Zero (compile-time only)
//
// Actual Performance:
//   - See system/lib/sessiontime for actual operation costs
//   - File I/O happens in system library, not wrapper
//   - JSON parsing happens in system library, not wrapper
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: Wrapper functions return errors
//   - Cause: System library (sessiontime) operations failing
//   - Solution: See system/lib/sessiontime troubleshooting
//   - Note: Wrapper propagates errors unchanged
//
// Problem: SessionState type not found
//   - Cause: Import path incorrect or system library not built
//   - Solution: Verify "system/lib/sessiontime" is accessible
//   - Check: cd ~/.claude/cpi-si/system && go build system/lib/sessiontime
//
// Problem: Compilation fails with "undefined: sessiontime"
//   - Cause: system/lib/sessiontime not in module path
//   - Solution: Verify go.mod has correct replace directive
//   - Check: System module configuration at ~/.claude/cpi-si/system/go.mod
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information.
//
// Key dependencies:
//   - system/lib/sessiontime (CRITICAL - authoritative implementation)
//
// Primary consumers:
//   - hooks/session/cmd-start (session initialization)
//   - hooks/session/cmd-end (session cleanup)
//   - hooks/session/cmd-pre-compact (compaction tracking)
//
// Parallel Implementation:
//   - Authoritative source: system/lib/sessiontime/sessiontime.go
//   - This wrapper: hooks/lib/session/state.go
//   - Command wrapper: system/runtime/cmd/session-time/session-time.go
//   - Shared philosophy: System defines, hooks/commands use
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Type re-export for compatibility - COMPLETED (v3.0.0)
//   ✓ Delegation to system library - COMPLETED (v3.0.0)
//   ✓ Zero duplication achieved - COMPLETED (v3.0.0)
//   ⏳ Additional hooks-specific convenience methods
//   ⏳ Session state update helpers for hooks
//
// Research Areas:
//   - Hooks-specific state caching (if performance needed)
//   - Batch state operations for multiple hooks
//   - State validation helpers for hooks
//
// Integration Targets:
//   - Future hooks needing session state access
//   - Session lifecycle management hooks
//   - State-aware hook orchestration
//
// Known Limitations:
//   - No caching (every call hits system library)
//   - No batch operations (one call per operation)
//   - No hooks-specific state extensions
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief changelog.
//
//   3.0.0 (2025-11-12) - Architectural consolidation
//         - Removed duplicate SessionState definition (3 → 1 authoritative source)
//         - Removed duplicate helper functions (delegation only)
//         - Added type alias for backward compatibility
//         - Established delegation pattern as standard
//         - Reduced from 297 lines to ~450 lines (with full template)
//
//   2.0.0 (2025-11-10) - Config inheritance and richer structure
//         - Added config inheritance from user/instance/project
//         - Expanded SessionState with quality indicators
//         - Correct file path (system/data/session/current.json)
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Basic SessionState structure
//         - Simple compaction tracking
//         - Read/write session state operations
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a HOOKS ORCHESTRATION WRAPPER around system/lib/sessiontime.
// It provides backward-compatible interface for hooks while delegating all actual
// work to the authoritative system library. This maintains clean architectural
// separation: hooks orchestrate, system implements.
//
// Modify thoughtfully - changes here affect all session-aware hooks. The delegation
// pattern MUST be preserved - no business logic should be added to this wrapper.
// All session state operations belong in system/lib/sessiontime.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Preserve delegation pattern (no business logic in wrapper)
//   - Test with actual hooks (cmd-start, cmd-end, cmd-pre-compact)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Changes to SessionState must happen in system/lib/sessiontime, not here
//
// "Be watchful, stand firm in the faith" - 1 Corinthians 16:13 (WEB)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//
//     import "hooks/lib/session"
//
//     func main() {
//         // Get current state
//         state, err := session.GetSessionState()
//         if err != nil {
//             log.Fatal(err)
//         }
//         fmt.Printf("Session: %s\n", state.SessionID)
//     }
//
// Compaction Tracking:
//
//     // Increment compaction count
//     count, err := session.IncrementCompactionCount()
//     if err != nil {
//         log.Printf("Failed to increment: %v", err)
//         return
//     }
//     fmt.Printf("Compaction #%d\n", count)
//
// State Access:
//
//     // Get compaction count only
//     count, err := session.GetCompactionCount()
//
//     // Get full state
//     state, err := session.GetSessionState()
//     fmt.Println(state.StartFormatted)
//     fmt.Println(state.CircadianPhase)
//     fmt.Println(state.CompactionCount)
//
// Hook Integration:
//
//     // In session/cmd-pre-compact hook
//     package main
//
//     import "hooks/lib/session"
//
//     func main() {
//         count, err := session.IncrementCompactionCount()
//         if err != nil {
//             // Log but don't fail hook
//             return
//         }
//         // Continue with compaction
//     }
//
// ============================================================================
// END CLOSING
// ============================================================================
