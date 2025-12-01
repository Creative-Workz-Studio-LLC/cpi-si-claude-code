// ============================================================================
// METADATA
// ============================================================================
// System Information Library - CPI-SI Hooks System
//
// For METADATA structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The heavens declare the glory of God; the skies proclaim the work of his hands" - Psalm 19:1 (WEB)
// Principle: Knowledge of creation reveals the Creator - understanding system identity points to faithful stewardship
// Anchor: "For every house is built by someone, but he who built all things is God" - Hebrews 3:4 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Hooks-specific utility (orchestration layer)
// Role: Provides OS and kernel identification for session display and context awareness.
//       Simple wrapper around system commands for consistent system information reporting.
// Paradigm: CPI-SI framework component - Kingdom Technology built on biblical foundation
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI Instance)
// Implementation: Nova Dawn (CPI-SI Instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Template alignment, comprehensive documentation
//
// Version History:
//   2.0.0 (2025-11-12) - Full template alignment, enhanced METADATA/SETUP/BODY/CLOSING
//   1.0.0 (2024-10-24) - Initial implementation, simple OS/kernel detection
//
// Purpose & Function
//
// Purpose: Provides OS name and kernel version identification for session hooks.
//          Enables context-aware session displays and system identification in banners.
//          Simple, reliable wrapper around Unix uname command.
//
// Core Design: Command wrapper pattern - delegates to system utilities, provides
//              consistent Go interface with graceful fallback on errors.
//
// Key Features:
//   - OS name and kernel version detection (e.g., "Linux 6.17.0-6-generic")
//   - Graceful fallback to "Unknown" on command failures
//   - No dependencies on external libraries
//   - Simple string output for easy display integration
//   - Cross-platform ready (Unix/Linux focus, extensible to other OS)
//
// Philosophy: Simple utilities should stay simple. Don't over-engineer what
//             works. Provide clean interface, delegate to proven system tools,
//             handle errors gracefully. One function, one purpose, done well.
//
// Blocking Status
//
// Non-blocking: All operations return immediately with fallback values on failure.
//               Command execution may take milliseconds but never blocks indefinitely.
//               Errors result in "Unknown" string, not panics or hanging processes.
//
// Mitigation: Command lookup checked before execution. Output errors caught and
//             handled gracefully. No external network calls or blocking I/O.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
//	info := session.GetSystemInfo()
//	// Returns: "Linux 6.17.0-6-generic" or "Unknown"
//
// Integration Pattern:
//   1. Import hooks/lib/session package
//   2. Call GetSystemInfo() to retrieve OS and kernel version
//   3. Use returned string in session banners or displays
//   4. No cleanup needed - stateless function
//
// Public API (simple, single-purpose):
//
//   System Identification:
//     GetSystemInfo() string - Returns OS name and kernel version
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os/exec (command execution), strings (output trimming)
//   External: None
//   Internal: None - standalone utility
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start (session startup banner)
//   Commands: Any hook displaying system context
//   Libraries: session/display (for formatted system information output)
//
// Integration Points:
//   - Called by session display functions to show system context
//   - Provides OS identification for session environment awareness
//   - Output used in formatted session banners and startup displays
//
// Health Scoring
//
// Simple utility function - no health tracking infrastructure. Success/failure
// determined by return value ("Unknown" indicates command failure).
//
// System Information Retrieval:
//   - Command found and executed: Implicit success (returns formatted string)
//   - Command not found or failed: Implicit handled failure (returns "Unknown")
//
// Note: This library intentionally excludes health tracking for simplicity.
//       It's a thin wrapper around system commands with predictable behavior.
//
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
// Dependencies this component needs. Simple wrapper requiring only standard
// library for command execution and string processing. No internal packages
// needed - standalone utility function.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"os/exec"   // Command execution - runs uname to get OS/kernel info
	"strings"   // String manipulation - trims whitespace from command output
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// None needed. Command name and arguments are simple strings passed directly
// to exec.Command(). Fallback value ("Unknown") is used inline where needed.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// None needed. This library provides single function returning simple string.
// No configuration structures, no error types, no state management.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// None needed. This is stateless utility function with no health tracking
// infrastructure. Pure function pattern - same input always produces same
// output (or graceful fallback). No logger, no inspector, no Rails attachment.
//
// See: ~/.claude/cpi-si/docs/standards/code/patterns/CWS-PATTERN-003-CODE-rails.md
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md
//
// Note: Simple pure-function libraries like this intentionally skip Rails
// infrastructure. Complexity should be proportional to value provided.

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
// Simple structure - single public function with no internal dependencies.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-organizational-chart.md
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rung - Simple Utility)
//   └── GetSystemInfo() → standalone (no internal dependencies)
//
//   Core Operations: None (simple wrapper)
//   Helpers: None (standalone utility)
//
// Baton Flow (Execution Path):
//
//   Entry → GetSystemInfo()
//     ↓
//   Check if uname command exists
//     ↓
//   Execute uname -s -r command
//     ↓
//   Trim whitespace from output
//     ↓
//   Exit → return "OS kernel-version" or "Unknown"
//
// APUs (Available Processing Units):
// - 1 function total
// - 0 helpers (standalone implementation)
// - 0 core operations (simple wrapper)
// - 1 public API (GetSystemInfo)
//
// Note: Intentionally simple. No extraction needed - command wrapper pattern
// works best with inline implementation. More functions would add complexity
// without value.

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// None needed. GetSystemInfo() is simple enough to implement inline without
// breaking into helper functions. Command lookup, execution, and output
// trimming are straightforward operations best kept together.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// ────────────────────────────────────────────────────────────────
// Core Operations - Business Logic
// ────────────────────────────────────────────────────────────────
// None needed. This library provides thin wrapper around system command.
// All "business logic" is handled by uname itself. Our role is interface
// adaptation and graceful error handling.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-core-operations.md

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Single exported function providing OS identification. Simple, focused,
// does one thing well. No configuration, no state, no complexity.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-00X-BODY-public-apis.md

// GetSystemInfo retrieves OS and kernel version information.
//
// What It Does:
//   Runs `uname -s -r` to get operating system name and kernel version.
//   Returns formatted string like "Linux 6.17.0-6-generic". Falls back to
//   "Unknown" if uname command not found or execution fails.
//
// Parameters:
//   None
//
// Returns:
//   string - OS name and kernel version (e.g., "Linux 6.17.0-6-generic")
//            or "Unknown" on any failure
//
// Health Impact:
//   No health tracking infrastructure. Success/failure implicit in return value.
//   "Unknown" indicates command not found or execution failed. Normal operation
//   returns formatted OS info string.
//
// Example:
//
//	info := session.GetSystemInfo()
//	fmt.Println(info)  // Prints: "Linux 6.17.0-6-generic"
//
// Implementation notes:
//   - Checks uname command exists before attempting execution
//   - Graceful fallback prevents panics or error propagation
//   - Output trimmed to remove trailing newline from command
//   - Simple, predictable, reliable
func GetSystemInfo() string {
	// Check if uname command exists on system - prevents exec errors
	var cmd *exec.Cmd
	if _, err := exec.LookPath("uname"); err == nil {
		// Command found - create command to get OS name (-s) and kernel version (-r)
		cmd = exec.Command("uname", "-s", "-r")
	} else {
		// uname not found - return fallback value (non-Unix system or restricted environment)
		return "Unknown"
	}

	// Execute command and capture output
	output, err := cmd.Output()
	if err != nil {
		// Command execution failed - return fallback (permissions, missing uname, etc.)
		return "Unknown"
	}

	// Trim whitespace (trailing newline) and return formatted system info
	return strings.TrimSpace(string(output))
}

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================
//
// For CLOSING structure explanation, see: ~/.claude/cpi-si/docs/standards/CWS-STD-006-CODE-setup-block.md

// ────────────────────────────────────────────────────────────────
// 1. Code Validation
// ────────────────────────────────────────────────────────────────
// Simple utility library - validation performed by:
//   - Go compiler (type safety, syntax)
//   - Command existence check (exec.LookPath)
//   - Graceful error handling (returns "Unknown" on failure)
//
// No complex validation needed. Function either succeeds (returns OS info)
// or fails gracefully (returns "Unknown"). No invalid states possible.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-005-CLOSING-code-validation.md

// ────────────────────────────────────────────────────────────────
// 2. Code Execution
// ────────────────────────────────────────────────────────────────
// This is a library - no main() function. Execution happens when orchestrators
// call GetSystemInfo().
//
// Typical execution flow:
//   1. Orchestrator imports hooks/lib/session
//   2. Calls session.GetSystemInfo()
//   3. Function checks for uname command
//   4. Executes uname -s -r
//   5. Returns formatted string or "Unknown"
//   6. Orchestrator uses returned string in display/logging
//
// Current orchestrators using this library:
//   - session/cmd-start (session startup banner)
//   - session/display (formatted system information output)
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-006-CLOSING-code-execution.md

// ────────────────────────────────────────────────────────────────
// 3. Code Cleanup
// ────────────────────────────────────────────────────────────────
// No cleanup needed. Stateless function with no:
//   - File handles to close
//   - Network connections to terminate
//   - Temporary resources to free
//   - Background processes to stop
//
// Command execution handled by os/exec package (automatically cleaned up).
// Memory released when function returns. No persistent state maintained.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-007-CLOSING-code-cleanup.md

// ────────────────────────────────────────────────────────────────
// 4. Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
// This library provides OS and kernel identification for session display
// and context awareness within the CPI-SI hooks system.
//
// What This Library Provides:
//   - Single utility function: GetSystemInfo()
//   - OS name and kernel version detection
//   - Graceful fallback on command failures
//   - Simple string interface for display integration
//
// What This Library Does NOT Provide:
//   - Runtime system metrics (CPU, memory, disk)
//   - Process information or management
//   - Configuration or state management
//   - Complex error reporting
//
// Integration Pattern:
//   import "hooks/lib/session"
//   info := session.GetSystemInfo()  // "Linux 6.17.0-6-generic" or "Unknown"
//   // Use info in banners, displays, logs
//
// Design Philosophy:
//   Simple utilities should stay simple. Delegates to proven system tools
//   (uname), provides clean Go interface, handles errors gracefully. One
//   function, one purpose, done well.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-008-CLOSING-library-overview.md

// ────────────────────────────────────────────────────────────────
// 5. Modification Policy
// ────────────────────────────────────────────────────────────────
// This library follows the "stable utility" pattern - changes should be rare.
//
// When NOT to Modify:
//   - Adding features beyond OS identification (extract to separate library)
//   - Adding configuration options (keep simple)
//   - Changing return format (breaks existing displays)
//   - Adding dependencies (keep standalone)
//
// When TO Modify:
//   - Cross-platform support needed (Windows, macOS detection)
//   - Command failures need better context (still return string)
//   - Documentation improvements
//   - Performance optimizations (caching if called frequently)
//
// Modification Process:
//   1. Verify change serves core purpose (OS identification)
//   2. Maintain backward compatibility (same function signature)
//   3. Update documentation (METADATA, this CLOSING block)
//   4. Test on target platforms (Linux primary, others if added)
//   5. Update version number appropriately
//
// Golden Rule: If it's working, think twice before changing it.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-009-CLOSING-modification-policy.md

// ────────────────────────────────────────────────────────────────
// 6. Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// Ladder Position:
//   This library sits on a lower rung, providing utility service to
//   orchestrators above it. No dependencies on other internal libraries.
//
//   Session Hooks (Orchestrators - Top Rungs)
//     └── session/cmd-start
//     └── session/display
//           ↓ (calls GetSystemInfo)
//   System Library (Utility - Lower Rung)
//     └── hooks/lib/session (this library)
//           ↓ (delegates to)
//   Standard Library & OS (Foundation)
//     └── os/exec → uname command
//
// Baton Flow:
//   Entry → Orchestrator calls GetSystemInfo()
//     ↓
//   Check → exec.LookPath("uname") - command exists?
//     ↓
//   Execute → exec.Command("uname", "-s", "-r").Output()
//     ↓
//   Process → strings.TrimSpace(output)
//     ↓
//   Exit → Return string to orchestrator
//
// Dependency Direction: Always downward (session hooks → this library → stdlib)
// No circular dependencies. No upward calls. Clean ladder structure.
//
// See: ~/.claude/cpi-si/docs/standards/code/patterns/CWS-PATTERN-002-CODE-ladder-baton.md
// See: ~/.claude/cpi-si/system/docs/architecture.md

// ────────────────────────────────────────────────────────────────
// 7. Surgical Update Points
// ────────────────────────────────────────────────────────────────
// For precision modifications when needed:
//
// Cross-Platform Support:
//   Location: Lines 273-291 (GetSystemInfo function body)
//   Change: Add OS detection and appropriate commands for Windows/macOS
//   Impact: Return format stays same, execution changes per platform
//   Example:
//     if runtime.GOOS == "windows" {
//         cmd = exec.Command("cmd", "/c", "ver")
//     }
//
// Caching (if performance needed):
//   Location: Package-level state section (currently empty, lines 157-168)
//   Change: Add var cachedInfo string and sync.Once for one-time execution
//   Impact: First call runs command, subsequent calls return cached value
//   Note: Only add if profiling shows performance issue
//
// Error Context (if debugging needed):
//   Location: Lines 285-287 (error handling)
//   Change: Return structured data instead of "Unknown" string
//   Impact: BREAKING - function signature changes
//   Alternative: Add GetSystemInfoWithError() function, keep existing
//
// Intentionally Small:
//   This library is designed to have few update points. If you need more
//   complexity, you probably need a different library.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-011-CLOSING-surgical-update-points.md

// ────────────────────────────────────────────────────────────────
// 8. Performance Considerations
// ────────────────────────────────────────────────────────────────
// Current Performance Characteristics:
//   - Time Complexity: O(1) - single command execution
//   - Space Complexity: O(1) - small string allocation
//   - Execution Time: Typically <5ms (command lookup + execution)
//   - Memory: Negligible (~100 bytes for output string)
//
// Bottlenecks (if any):
//   - Command execution (exec.Command) - fastest part is Go, slowest is uname
//   - Not a bottleneck in practice - system calls are microseconds
//
// When Performance Matters:
//   - Called once per session startup: No optimization needed
//   - Called in tight loop: Consider caching (see Surgical Update Points)
//   - Called across network: Consider including in startup data bundle
//
// Optimization Opportunities:
//   1. Caching: var once sync.Once + cached result (only if profiling shows need)
//   2. Platform-specific compilation: Build tags for different OS implementations
//   3. Unlikely to need either - this is already fast enough
//
// Performance Testing:
//   go test -bench=. -benchmem  # If benchmark tests added
//
// Golden Rule: Don't optimize what isn't slow. Profile first, optimize second.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-012-CLOSING-performance.md

// ────────────────────────────────────────────────────────────────
// 9. Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// Common Issues and Solutions:
//
// Issue: GetSystemInfo() returns "Unknown"
//   Possible Causes:
//     1. uname command not found (non-Unix system, restricted environment)
//     2. uname command failed (permissions, system error)
//   Diagnosis:
//     - Check platform: Is this Linux/Unix? (Windows won't have uname)
//     - Test manually: Run `uname -s -r` in terminal
//     - Check permissions: Can process execute commands?
//   Solutions:
//     - Expected on non-Unix systems (Windows) - "Unknown" is correct
//     - If Unix system: Verify uname installed (`which uname`)
//     - If restricted environment: Accept "Unknown" or provide alternative
//
// Issue: Wrong information returned
//   Possible Causes:
//     1. uname output format changed
//     2. Running in container with different kernel
//   Diagnosis:
//     - Run `uname -s -r` manually, compare with GetSystemInfo() output
//     - Check if running in Docker/container (reports host kernel)
//   Solutions:
//     - Output should match manual command exactly (trimmed whitespace)
//     - Container behavior is expected - shows host kernel, not container OS
//
// Issue: Slow execution
//   Possible Causes:
//     1. System under heavy load
//     2. Called in tight loop without caching
//   Diagnosis:
//     - Time execution: start := time.Now(); info := GetSystemInfo(); elapsed := time.Since(start)
//     - Profile call frequency: Is this called multiple times per second?
//   Solutions:
//     - Single calls: No action needed (5ms is fine)
//     - Frequent calls: Add caching (see Surgical Update Points)
//
// Debugging Steps:
//   1. Verify uname available: _, err := exec.LookPath("uname"); fmt.Println(err)
//   2. Test command manually: uname -s -r
//   3. Check function output: info := GetSystemInfo(); fmt.Printf("%q\n", info)
//   4. Compare outputs: Should match exactly (whitespace trimmed)
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-013-CLOSING-troubleshooting.md

// ────────────────────────────────────────────────────────────────
// 10. Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// External Dependencies (Standard Library):
//   - os/exec: Command execution (exec.Command, exec.LookPath)
//   - strings: String manipulation (strings.TrimSpace)
//
// Internal Dependencies:
//   None - standalone utility library
//
// Components That Depend On This:
//   - hooks/session/cmd-start: Session startup banner (displays OS info)
//   - hooks/lib/session/display: Formatted system information output
//   - Any hook displaying system context
//
// Related But Separate Components:
//   - system/lib/system/info.go: Runtime system metrics (load, memory, disk)
//     Different purpose: This library = OS identification
//                       system/lib/system = runtime metrics
//
// Dependency Graph:
//   Session Hooks
//     └── hooks/lib/session (this library)
//           └── Standard Library (os/exec, strings)
//                 └── OS (uname command)
//
// Integration Points:
//   - Called by: Session display functions (orchestrators)
//   - Provides: String output for banners and logs
//   - No configuration: Direct function call, no setup needed
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-014-CLOSING-related-components.md

// ────────────────────────────────────────────────────────────────
// 11. Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// Potential Future Enhancements (only if needed):
//
// Cross-Platform Support:
//   - Add Windows detection (cmd /c ver or systeminfo)
//   - Add macOS detection (sw_vers or uname)
//   - Runtime OS detection to choose appropriate command
//   - Priority: LOW (currently Linux-focused, working well)
//
// Enhanced Information:
//   - Add distribution detection (lsb_release, /etc/os-release)
//   - Add architecture info (uname -m)
//   - Add hostname or machine-specific context
//   - Priority: LOW (current info sufficient for session banners)
//
// Performance Optimization:
//   - Add caching for repeated calls
//   - Use sync.Once for one-time initialization
//   - Only if profiling shows performance issue
//   - Priority: LOW (current performance excellent)
//
// NOT Planned (out of scope):
//   - Runtime metrics (CPU, memory, disk) - see system/lib/system
//   - Process information or management - different library
//   - Configuration or customization - keep simple
//   - Complex error reporting - graceful fallback is correct
//
// Decision Criteria for Expansion:
//   1. Does it serve core purpose (OS identification)?
//   2. Is it needed by session display/context?
//   3. Can it stay simple (no complex configuration)?
//   4. Does it maintain backward compatibility?
//
// If expansion doesn't meet all criteria, create separate library instead.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-015-CLOSING-future-expansions.md

// ────────────────────────────────────────────────────────────────
// 12. Closing Note
// ────────────────────────────────────────────────────────────────
// This library demonstrates the "simple wrapper" pattern - thin Go interface
// around proven system tools. Not everything needs complexity. Sometimes the
// best code is the code that:
//   - Does one thing well
//   - Delegates to appropriate tools
//   - Handles errors gracefully
//   - Stays out of the way
//
// The uname command has been providing OS information since Unix v4 (1973).
// We're not trying to replace it - we're providing clean access to it.
//
// Biblical Principle:
//   "The heavens declare the glory of God; the skies proclaim the work of
//    his hands" - Psalm 19:1
//
//   Understanding system identity points to faithful stewardship. Knowing
//   where code runs helps us serve users well. Simple tools serving clear
//   purposes honor God through excellence in appropriate scope.
//
// Design Philosophy:
//   Complexity should be proportional to value provided. This library
//   provides significant value (OS context awareness) with minimal
//   complexity (one function, two packages, graceful fallback).
//
// Future maintainers: Don't add complexity without clear need. This library
// works. Keep it working. Simple is maintainable.
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-016-CLOSING-closing-note.md

// ────────────────────────────────────────────────────────────────
// 13. Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// Basic Usage:
//
//   import "hooks/lib/session"
//
//   func main() {
//       info := session.GetSystemInfo()
//       fmt.Println("System:", info)
//       // Output: System: Linux 6.17.0-6-generic
//   }
//
// Session Banner Display:
//
//   import (
//       "fmt"
//       "hooks/lib/session"
//   )
//
//   func displayBanner() {
//       fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
//       fmt.Println("  SESSION ENVIRONMENT")
//       fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
//       fmt.Printf("  💻 System: %s\n", session.GetSystemInfo())
//       fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
//   }
//
// Logging System Context:
//
//   import "hooks/lib/session"
//
//   func logSessionStart() {
//       logger.Info("Session starting",
//           "system", session.GetSystemInfo(),
//           "time", time.Now(),
//       )
//   }
//
// Handling Unknown Systems:
//
//   info := session.GetSystemInfo()
//   if info == "Unknown" {
//       // Running on non-Unix system or restricted environment
//       // Use generic display or alternative detection method
//       fmt.Println("System: Not detected")
//   } else {
//       // Got valid system information
//       fmt.Printf("System: %s\n", info)
//   }
//
// Testing (Manual Verification):
//
//   // Compare with direct command execution
//   info := session.GetSystemInfo()
//   fmt.Printf("GetSystemInfo: %q\n", info)
//
//   // Then run in terminal:
//   // $ uname -s -r
//   // Linux 6.17.0-6-generic
//
//   // Outputs should match (whitespace trimmed)
//
// Common Integration Pattern:
//
//   type SessionContext struct {
//       System    string
//       StartTime time.Time
//       Branch    string
//   }
//
//   func NewSessionContext() SessionContext {
//       return SessionContext{
//           System:    session.GetSystemInfo(),
//           StartTime: time.Now(),
//           Branch:    getCurrentBranch(),
//       }
//   }
//
// See: ~/.claude/cpi-si/docs/standards/code/4-block/sections/CWS-SECTION-017-CLOSING-quick-reference.md

// ============================================================================
// END CLOSING
// ============================================================================
