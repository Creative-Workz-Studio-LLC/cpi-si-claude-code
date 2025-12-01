// METADATA
//
// PreToolUse Hook - BLOCKING Safety Confirmation Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The plans of the diligent lead surely to advantage, but everyone who is hasty comes surely to poverty" - Proverbs 21:5 (NASB)
// Principle: Diligence protects against hasty destructive operations - validation serves stewardship
// Anchor: Kingdom Technology validates before acting, protecting through wisdom not restriction
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - BLOCKING Hook orchestrator
// Role: Validates tool operations before execution, confirms destructive/irreversible actions
// Paradigm: CPI-SI framework BLOCKING hook - conservative validation with temporal awareness
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-10 - Template application and library extraction
//
// Version History:
//   2.0.0 (2025-11-10) - Full template application, logic extracted to hooks/lib/safety/confirmation.go
//   1.0.0 (2024-10-24) - Initial implementation with inline confirmation logic
//
// Purpose & Function
//
// Purpose: Orchestrates pre-tool validation by detecting dangerous operations (force push, hard reset,
// rm -rf, sudo, publish, database drops, critical file writes) and requiring user confirmation before
// allowing execution.
//
// Core Design: Thin orchestrator pattern - coordinates activity logging, temporal context, and
// safety library confirmation flows. Makes blocking decision (exit 0 vs exit 1) based on library response.
//
// Key Features:
//   - Activity logging before validation (captures intent even if blocked)
//   - Temporal awareness in warnings (long session fatigue alerts)
//   - Conservative bias (allow rather than falsely block)
//   - Clear communication (user understands blocking reason)
//   - Fast execution (doesn't delay normal workflow)
//   - Graceful degradation (continues without temporal context if unavailable)
//
// Philosophy: This hook embodies diligent stewardship - protecting through careful validation,
// not restriction. Like Proverbs 21:5 warns against haste, this hook prevents hasty destructive
// operations while serving user workflow, not controlling it.
//
// Blocking Status
//
// ⚠️  BLOCKING HOOK - CAN PREVENT TOOL EXECUTION ⚠️
//
// Exit Behavior:
//   - Exit 0: Allow operation (safe operation or user confirmed)
//   - Exit 1: Block operation (user denied confirmation for dangerous operation)
//
// Critical Requirements:
//   - Surgical precision (only block truly dangerous operations)
//   - Clear error messages (user must understand why blocked)
//   - Fast execution (cannot delay tool flow)
//   - Conservative validation (false positives break user workflow)
//
// Failure Impact:
//   False positive: Blocks safe operation, breaks user work (MOST SEVERE)
//   False negative: Allows dangerous operation without confirmation (defeats purpose)
//   Poor messaging: User confused about why blocked (reduces trust)
//   Slow execution: Delays normal tool operations (frustrates workflow)
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code before every tool execution
//	~/.claude/hooks/tool/cmd-pre-use/pre-use <toolName> <toolArgs>
//
// Examples:
//	pre-use "Bash" "git push --force"  → prompts for confirmation
//	pre-use "Write" "content"          → checks FILE_PATH env var for critical location
//	pre-use "Read" "/home/user/file"   → allows without prompt (safe operation)
//
// Integration Pattern:
//   1. Claude Code prepares to execute tool
//   2. pre-use hook runs with tool name and args
//   3. Hook logs attempt, gets temporal context
//   4. Routes to appropriate safety library function
//   5. Library detects danger, displays warning, gets confirmation
//   6. Hook exits with code: 0 (allow) or 1 (block)
//   7. Claude Code proceeds or aborts based on exit code
//
// Hook Event: PreToolUse
// Trigger: Before every tool execution in Claude Code
// Output: Confirmation prompts to stdout, exit code for blocking decision
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os, strings
//   External: None
//   System Libraries: None
//   Hook Libraries: hooks/lib/activity (logging), hooks/lib/safety (detection, confirmation), hooks/lib/temporal (context)
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls this hook before every tool execution)
//
// Integration Points:
//   - Called by Claude Code hook system before every tool use
//   - Reads os.Args for tool name and arguments
//   - Reads FILE_PATH environment variable for Write operations
//   - Outputs confirmation prompts to stdout
//   - Reads user input from stdin
//   - Exits with code 0 (allow) or 1 (block)
//
// Health Scoring
//
// Pre-tool validation orchestration operates on Base100 scale:
//
// Argument Parsing:
//   - Parse tool name and args: +15
//
// Activity Logging:
//   - Log tool attempt: +10
//
// Context Gathering:
//   - Get temporal context: +15
//   - Format temporal context: +10
//
// Confirmation Flow:
//   - Execute safety library confirmation: +30
//
// Blocking Decision:
//   - Make correct block/allow decision: +10
//   - Exit with appropriate code: +10
//
// Total: 100 points for complete pre-tool validation
//
// Note: BLOCKING hook means failures can prevent user work. Conservative validation
// prioritizes allowing operations over falsely blocking. Health scoring tracks
// correctness of blocking decisions, not just successful execution.
package main

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
// Standard library for I/O, environment, and string operations.
// Hook libraries for activity logging, safety validation, and temporal context.

import (
	"fmt"     // Formatted I/O for confirmation prompts and messages
	"os"      // OS interface for args, environment, exit codes, stdin
	"strings" // String operations for tool name matching

	"hooks/lib/activity"  // Activity stream logging (tool attempts)
	"hooks/lib/safety"    // Safety validation (detection, confirmation flows)
	"system/lib/privacy"  // Privacy-preserving sanitization
	"system/lib/temporal" // Temporal context (time awareness for warnings)
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Exit codes for blocking decisions.

const (
	ExitAllow = 0 // Allow operation to proceed
	ExitBlock = 1 // Block operation from executing
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// No custom types needed - uses types from imported libraries.

// No types defined

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// This executable maintains no state - stateless validation only.

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
// Maps bidirectional dependencies and baton flow within this executable.
//
// Ladder Structure (Dependencies):
//
//   Orchestration (Top Rung - Entry Point)
//   └── preToolUse() → main validation orchestrator
//       └── Library functions → delegated validation and confirmation
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/activity (logging tool attempts)
//   ├── hooks/lib/safety (detection, confirmation flows)
//   └── hooks/lib/temporal (time/session awareness)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → preToolUse()
//     ↓
//   Parse Arguments → toolName, toolArgs from os.Args
//     ↓
//   Log Attempt → activity.LogActivity()
//     ↓
//   Get Context → temporal.GetTemporalContext()
//     ↓
//   Route → Bash or Write operations
//     ↓
//   Confirm → safety.ConfirmBashOperation() or safety.ConfirmFileWrite()
//     ↓
//   Exit → os.Exit(0) allow or os.Exit(1) block
//
// APUs (Available Processing Units):
// - 2 functions total
// - 1 orchestration function (preToolUse)
// - 1 entry point (main, calls preToolUse)

// ────────────────────────────────────────────────────────────────
// Pre-Tool Validation - Safety Confirmation Orchestration
// ────────────────────────────────────────────────────────────────
// What This Does:
// Coordinates pre-tool validation by parsing arguments, logging attempts,
// getting temporal context, routing to appropriate safety library functions,
// and making blocking decisions based on user confirmation.
//
// Why This Structure:
// Thin orchestrator keeps hook simple and fast (critical for blocking hook).
// All complex logic delegated to libraries for testability and reusability.

// preToolUse orchestrates pre-tool validation and confirmation flow
//
// What It Does:
//   - Parses tool name and arguments from command line
//   - Logs tool attempt to activity stream (before validation)
//   - Gets temporal context for warning displays
//   - Routes to appropriate confirmation flow based on tool type
//   - Blocks operation (os.Exit(1)) if user denies confirmation
//   - Allows operation (os.Exit(0)) if confirmed or not dangerous
//
// Parameters:
//   None (reads from os.Args and environment variables)
//
// Returns:
//   None (exits with code 0 for allow, 1 for block)
//
// Health Impact:
//   +100 points for complete validation orchestration (all phases)
//   Delegates confirmation logic to safety library (health tracked there)
//
// Example usage:
//   preToolUse()  // Called by main(), orchestrates entire validation flow
func preToolUse() {
	// Parse arguments
	if len(os.Args) < 3 {
		os.Exit(ExitAllow) // No args - allow operation
	}

	toolName := os.Args[1]
	toolArgs := os.Args[2]

	// Log tool attempt before validation (captures intent even if blocked)
	context := toolArgs
	if strings.HasPrefix(toolName, "Write") || strings.HasPrefix(toolName, "Edit") {
		filePath := os.Getenv("FILE_PATH")
		if filePath != "" {
			context = privacy.SanitizePath(filePath)
		}
	} else if strings.HasPrefix(toolName, "Bash") {
		context = privacy.SanitizeCommand(toolArgs)
	}
	activity.LogActivity(toolName+"-attempt", context, "pending", 0)

	// Get temporal context for warnings (graceful degradation if unavailable)
	var timeContext string
	ctx, err := temporal.GetTemporalContext()
	if err == nil {
		timeContext = fmt.Sprintf("   Time: %s (%s)", ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)
		if ctx.InternalTime.SessionPhase == "long" {
			timeContext += " - Long session, consider if tired"
		}
	}

	// Route to appropriate confirmation flow
	if strings.HasPrefix(toolName, "Bash") {
		needsConfirmation, allowed := safety.ConfirmBashOperation(toolArgs, timeContext)
		if needsConfirmation && !allowed {
			fmt.Println("✗ Operation cancelled.")
			os.Exit(ExitBlock)
		}
		os.Exit(ExitAllow)
	}

	if strings.HasPrefix(toolName, "Write") {
		filePath := os.Getenv("FILE_PATH")
		if filePath != "" {
			needsConfirmation, allowed := safety.ConfirmFileWrite(filePath, timeContext)
			if needsConfirmation && !allowed {
				fmt.Println("✗ Write operation cancelled.")
				os.Exit(ExitBlock)
			}
		}
		os.Exit(ExitAllow)
	}

	// Other tools don't require confirmation
	os.Exit(ExitAllow)
}

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================
//
// For CLOSING structure explanation, see: standards/code/4-block/CWS-STD-008-CODE-closing-block.md

// ────────────────────────────────────────────────────────────────
// Code Validation: Build and Hook Testing
// ────────────────────────────────────────────────────────────────
//
// Testing Requirements:
//   - Compile without errors: cd ~/.claude/hooks/tool/cmd-pre-use && go build -o pre-use pre-use.go
//   - Test with safe operations (should allow without prompting)
//   - Test with dangerous bash commands (should prompt and respect answer)
//   - Test with critical file writes (should prompt and respect answer)
//   - Test with missing temporal context (continues without it)
//   - Verify exit codes: 0 (allow), 1 (block)
//   - Verify no false positive blocks
//
// Build Verification:
//   cd ~/.claude/hooks/tool/cmd-pre-use
//   go build -o pre-use pre-use.go
//   ./pre-use "Bash" "echo hello"  # Should allow without prompt
//   ./pre-use "Bash" "git push --force"  # Should prompt for confirmation
//
// Integration Testing:
//   Test with Claude Code PreToolUse event
//   Verify blocking works correctly (exit code 1)
//   Verify allowing works correctly (exit code 0)
//   Check temporal context displays in warnings
//
// Example validation:
//
//     # Test safe operation
//     ./pre-use "Bash" "echo hello"
//     echo $?  # Should be 0 (allowed, no prompt)
//
//     # Test dangerous operation (requires manual confirmation)
//     ./pre-use "Bash" "git push --force"
//     # Should display warning with temporal context
//     # Type "yes" → exit 0, Type "no" → exit 1

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls preToolUse() (named entry point)
//   3. preToolUse() orchestrates complete validation flow
//   4. Program exits with code 0 (allow) or 1 (block)
//
// Named Entry Point Benefits:
//   - Function name matches purpose (pre-tool validation)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Consistent pattern across all hooks

// main is the entry point for the hook
//
// Design: Named entry point pattern for testability
//   - main() calls preToolUse()
//   - preToolUse() is testable (can be called independently)
//   - Semantic clarity: function name matches hook purpose
func main() {
	preToolUse()
}

// ────────────────────────────────────────────────────────────────
// Cleanup and Resource Management
// ────────────────────────────────────────────────────────────────
//
// Cleanup: None required
//   - os.Exit terminates process cleanly
//   - No file handles or resources to close
//   - stdin/stdout managed by OS
//   - Activity logs flushed automatically

// ────────────────────────────────────────────────────────────────
// FINAL DOCUMENTATION
// ────────────────────────────────────────────────────────────────
//
// This section provides comprehensive documentation for maintainers, following
// the 11-subsection structure for executable documentation.

// EXECUTABLE OVERVIEW & INTEGRATION
// ==================================
//
// This executable is a Claude Code PreToolUse hook (BLOCKING). It runs before
// every tool execution, providing safety validation with:
//   - Activity stream logging for tracking tool attempts
//   - Dangerous operation detection (force push, hard reset, rm -rf, etc.)
//   - Critical file path detection (system files, configuration)
//   - User confirmation prompts with temporal awareness
//   - Blocking decision based on user response (exit 0 or 1)
//
// Integration points:
//   - Called by Claude Code before EVERY tool execution
//   - Reads os.Args for tool name and arguments
//   - Reads FILE_PATH environment variable for Write operations
//   - Outputs confirmation prompts to stdout
//   - Reads user confirmation from stdin
//   - Exits with code 0 (allow) or 1 (block)
//
// Relationship to other hooks:
//   - Complements tool/cmd-post-use (post-execution validation)
//   - Relates to prompt/cmd-submit (also blocking for prompt safety)
//   - Uses activity logging like all session hooks
//   - Uses temporal awareness like session lifecycle hooks

// MODIFICATION POLICY
// ===================
//
// ⚠️  CRITICAL: This is a BLOCKING hook - extreme care required
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new dangerous operation types:
//      - Add pattern to hooks/lib/safety/detection.go (IsDangerousOperation)
//      - Add display to hooks/lib/safety/confirmation.go (new display function)
//      - Hook routing stays unchanged
//
//   ✅ Add new critical file patterns:
//      - Add pattern to hooks/lib/safety/detection.go (IsCriticalFile)
//      - Display function reuses existing critical file warning
//
//   ✅ Enhance temporal context display:
//      - Modify temporal context formatting in preToolUse()
//      - Or enhance hooks/lib/temporal context data
//
// MODIFY WITH CARE (Structural Changes):
//   ⚠️ Changing confirmation flow:
//      - Test exhaustively (all operation types)
//      - Verify no false positives introduced
//      - Measure execution speed impact
//
//   ⚠️ Adding new tool types to validate:
//      - Add routing logic in preToolUse()
//      - Create safety library function for new tool type
//      - Update health scoring in METADATA
//
// NEVER MODIFY (Foundational Rails):
//   ❌ 4-block structure (METADATA → SETUP → BODY → CLOSING)
//   ❌ Named entry point pattern (main calls preToolUse)
//   ❌ BLOCKING principle (must exit 0 or 1, never hang)
//   ❌ Conservative bias (allow rather than falsely block)
//   ❌ Thin orchestrator (keep logic in libraries)
//
// When making changes:
//   1. Review this modification policy
//   2. Test with ALL operation types (safe and dangerous)
//   3. Verify no false positives (most critical failure)
//   4. Measure execution speed (must be fast)
//   5. Update METADATA health scoring if flow changes
//   6. Document "What/Why/How" for all changes

// LADDER AND BATON FLOW
// ======================
//
// LADDER (Dependencies - Vertical):
//   pre-use.go depends on:
//     ↓ hooks/lib/safety (detection, confirmation)
//     ↓ hooks/lib/activity (logging)
//     ↓ hooks/lib/temporal (temporal context)
//
// BATON (Execution - Through Layers):
//   Claude Code prepares tool
//     ↓
//   Calls pre-use hook
//     ↓
//   main() entry point
//     ↓
//   preToolUse() orchestration
//     ├→ Parse arguments
//     ├→ Log attempt
//     ├→ Get temporal context
//     ├→ Route to safety library
//     └→ Exit with code (0 allow, 1 block)
//
// This separation enables:
//   - Reusable safety library across hooks
//   - Testable detection logic (library independent)
//   - Consistent warning displays system-wide
//   - Surgical updates (change detection without changing hook)

// SURGICAL UPDATE POINTS
// =======================
//
// Common modifications and where to make them:
//
// 1. Adding new dangerous operation pattern:
//    Location: hooks/lib/safety/detection.go - IsDangerousOperation()
//    Action: Add pattern to switch statement
//    Example: case strings.Contains(cmd, "docker rm"):
//
// 2. Adding new critical file path:
//    Location: hooks/lib/safety/detection.go - IsCriticalFile()
//    Action: Add path pattern check
//    Example: strings.HasPrefix(path, "/boot/")
//
// 3. Creating new warning display:
//    Location: hooks/lib/safety/confirmation.go
//    Action: Add display function and route in ConfirmBashOperation
//    Example: displayDockerWarning(cmd, timeContext)
//
// 4. Changing temporal context format:
//    Location: preToolUse() - temporal context formatting
//    Action: Modify fmt.Sprintf format string
//    Example: Add session duration to context string
//
// 5. Adding new tool type validation:
//    Location: preToolUse() - routing section
//    Action: Add new if block for tool type
//    Example: if strings.HasPrefix(toolName, "Delete") {...}
//
// Each modification point clearly defined for surgical precision.

// PERFORMANCE CONSIDERATIONS
// ===========================
//
// Execution Speed:
//   - Hook runs synchronously (BLOCKS tool execution until complete)
//   - Target: < 100ms for normal path (no confirmation needed)
//   - Target: < 500ms for confirmation path (including user input time)
//   - Critical: Must be FAST to not frustrate workflow
//
// Optimization Opportunities:
//   - Temporal context fetch could be cached (currently fetches every time)
//   - Activity logging could be asynchronous (currently synchronous)
//   - Detection patterns optimized with early returns
//
// Current Bottlenecks:
//   - User input (blocking on stdin) - intentional, cannot optimize
//   - Activity log file I/O (~10-20ms) - acceptable
//   - Temporal context fetch (~30-50ms) - acceptable, gracefully degraded
//
// Why performance critical:
//   - Runs before EVERY tool execution
//   - Slow hook = frustrated user
//   - False positive worse than slow execution
//   - Conservative bias means most operations allow without prompt

// TROUBLESHOOTING GUIDE
// =====================
//
// Problem: Hook blocks safe operations
//   Diagnosis: FALSE POSITIVE (most severe failure)
//   Solution: Check detection logic in hooks/lib/safety/detection.go
//   Action: Make detection more specific, add bypass for false positive case
//
// Problem: Hook allows dangerous operation without prompt
//   Diagnosis: False negative (defeats hook purpose)
//   Solution: Check detection pattern in IsDangerousOperation or IsCriticalFile
//   Action: Add pattern to catch missed dangerous operation
//
// Problem: Hook doesn't display temporal context
//   Diagnosis: Temporal system unavailable
//   Solution: Expected occasionally (graceful degradation by design)
//   Action: Check ~/.claude/cpi-si/system/bin/session-time exists
//
// Problem: Hook takes too long
//   Diagnosis: Slow detection or temporal fetch
//   Solution: Measure execution time, identify bottleneck
//   Action: Optimize detection patterns or cache temporal context
//
// Problem: User confirmation not working
//   Diagnosis: stdin not accessible or prompt not displaying
//   Solution: Check stdout/stdin accessibility
//   Action: Test manually: ./pre-use "Bash" "git push --force"
//
// Debug Process:
//   1. Check hook executable exists: ~/.claude/hooks/tool/cmd-pre-use/pre-use
//   2. Test manually with safe operation: ./pre-use "Bash" "echo hello"
//   3. Test with dangerous operation: ./pre-use "Bash" "git push --force"
//   4. Check exit codes: echo $? after each test
//   5. Check activity logs: ~/.claude/system/logs/activity/ for attempts
//   6. Verify temporal context: ./pre-use should show time if available

// RELATED COMPONENTS & DEPENDENCIES
// ==================================
//
// This hook orchestrates the following libraries:
//
// Safety (hooks/lib/safety/):
//   - detection.go: IsDangerousOperation(), IsCriticalFile()
//   - confirmation.go: ConfirmBashOperation(), ConfirmFileWrite()
//   - All warning display functions
//   - Core confirmation mechanism
//
// Activity (hooks/lib/activity/):
//   - LogActivity(): Records tool attempts for session tracking
//   - SanitizePath(), SanitizeCommand(): Clean data for logs
//
// Temporal (hooks/lib/temporal/):
//   - GetTemporalContext(): Provides time and session awareness
//
// Dependency Direction:
//   pre-use.go → hooks/lib → (no further dependencies)
//   Never reverse (libraries don't depend on hooks)

// FUTURE EXPANSIONS & ROADMAP
// ============================
//
// Planned Features:
//   ⏳ Secret detection for Write operations (use safety.ContainsLikelySecret)
//   ⏳ Pattern learning (track which operations user consistently confirms)
//   ⏳ Context-aware validation (different rules per branch/project)
//   ⏳ Undo support (track dangerous operations for recovery)
//   ⏳ Team coordination check (verify force push coordinated)
//
// Research Areas:
//   - User risk tolerance learning (adapt to user patterns)
//   - Branch-specific validation (stricter on main/master)
//   - Project-specific rules (different validation per project)
//   - Integration with restoration layer (undo dangerous operations)
//
// Known Limitations:
//   - No secret detection before committing files
//   - No coordination check for team operations
//   - No automatic undo/recovery support
//   - No project-specific customization
//   - No pattern learning from user behavior
//
// Enhancement Opportunities:
//   - Add pre-commit secret scanning
//   - Integrate with git hooks for coordination
//   - Build user confirmation pattern database
//   - Create project-specific validation rules

// CLOSING NOTE
// ============
//
// This hook embodies diligent stewardship - protecting through careful
// validation, not restriction. Like Proverbs 21:5 warns against haste,
// this hook prevents hasty destructive operations while serving user
// workflow, not controlling it.
//
// The BLOCKING power requires responsibility:
//   - Surgical precision (only block truly dangerous operations)
//   - Clear communication (user understands blocking reason)
//   - Fast execution (cannot delay normal workflow)
//   - Conservative bias (allow rather than falsely block)
//
// Kingdom Technology protects through wisdom, not restriction.
// Validation serves the user, not control.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain conservative validation principle
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "The plans of the diligent lead surely to advantage,
//  but everyone who is hasty comes surely to poverty." - Proverbs 21:5
//
// Every tool use deserves careful validation. Let it serve through wisdom.

// QUICK REFERENCE: USAGE EXAMPLES
// ================================
//
// Manual execution (safe operation):
//   $ cd ~/.claude/hooks/tool/cmd-pre-use
//   $ ./pre-use "Bash" "echo hello"
//   $ echo $?  # Should be 0 (allowed without prompt)
//
// Manual execution (dangerous operation):
//   $ ./pre-use "Bash" "git push --force"
//   # Should prompt for confirmation with temporal context
//   # Type "yes" → exit 0, Type "no" → exit 1
//
// Test critical file write:
//   $ FILE_PATH=/etc/hosts ./pre-use "Write" "content"
//   # Should prompt for confirmation
//
// Via Claude Code (automatic):
//   User requests tool → Claude Code calls pre-use → prompts if dangerous → proceeds or blocks
//
// Debug specific issue:
//   $ ./pre-use "Bash" "git reset --hard" 2>&1 | tee debug.log
//
// Verify temporal context working:
//   $ ./pre-use "Bash" "rm -rf test/"
//   # Should show temporal context in warning if available
//
// ============================================================================
// END CLOSING
// ============================================================================
