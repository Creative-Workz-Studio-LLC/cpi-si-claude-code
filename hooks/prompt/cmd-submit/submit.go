// METADATA
//
// UserPromptSubmit Hook - NON-BLOCKING Awareness Orchestrator
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Set a guard, O Lord, over my mouth; keep watch over the door of my lips" - Psalm 141:3
// Principle: Awareness before speech - check words before releasing them
// Anchor: Kingdom Technology promotes careful communication, protecting from hasty revelation
//
// CPI-SI Identity
//
// Component Type: EXECUTABLE - NON-BLOCKING Hook orchestrator
// Role: Provides awareness before prompt submission (secret detection warning)
// Paradigm: CPI-SI framework hook - warns but never blocks (user workflow priority)
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
//   2.0.0 (2025-11-10) - Full template application, logic extracted to hooks/lib/safety
//   1.0.0 (2024-10-24) - Initial implementation with inline warning logic
//
// Purpose & Function
//
// Purpose: Orchestrates pre-prompt awareness by detecting likely secrets in user prompts and
// warning before submission. Logs prompt activity for engagement tracking (privacy-preserving
// length only) and archives full prompts for session history.
//
// Core Design: Thin orchestrator pattern - coordinates activity logging, secret detection,
// warning display, and monitoring logging. NEVER blocks (always exits 0) - warnings only.
//
// Key Features:
//   - Privacy-preserving activity logging (prompt length only)
//   - Quick secret detection (API keys, passwords, tokens)
//   - Non-blocking warnings (awareness, not prevention)
//   - Async monitoring logging (full prompt for session history)
//   - VERY fast execution (runs before EVERY prompt)
//
// Philosophy: This hook embodies careful communication - awareness before speech, but not
// censorship. Like Psalm 141:3 requests guard over words, this hook provides awareness while
// respecting user autonomy. Warns but never blocks.
//
// Blocking Status
//
// NON-BLOCKING: Always exits 0 (never prevents prompt submission)
//
// Exit Behavior:
//   - Always Exit 0: Prompt submitted regardless of warnings
//
// Design Reasoning:
//   - Prompts run VERY frequently (every user interaction)
//   - False positives would frustrate workflow constantly
//   - User awareness more valuable than blocking
//   - Secret detection helpful but imperfect
//   - Workflow priority over protection
//
// Performance Critical:
//   - Target: < 50ms total execution (faster than pre-use)
//   - Secret detection must be quick
//   - Async logging (doesn't block prompt)
//   - Minimal display (brief warning only)
//
// Usage & Integration
//
// Usage:
//
//	# Called by Claude Code before every prompt submission
//	~/.claude/hooks/prompt/cmd-submit/submit
//
// Environment Variables:
//	PROMPT="user prompt text"  → Full prompt text for detection and logging
//
// Integration Pattern:
//   1. User types prompt, presses enter
//   2. Claude Code triggers UserPromptSubmit event
//   3. Hook runs with PROMPT environment variable
//   4. Hook logs activity (length only)
//   5. Hook checks for secrets, warns if detected
//   6. Hook logs full prompt async (monitoring)
//   7. Hook exits 0 (always allows)
//   8. Claude Code processes prompt normally
//
// Hook Event: UserPromptSubmit
// Trigger: Before every prompt submission
// Output: Warning to stdout if secrets detected, always exits 0
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os, strconv
//   External: None
//   System Libraries: None
//   Hook Libraries: hooks/lib/activity (logging), hooks/lib/safety (detection, display), hooks/lib/monitoring (async logging)
//
// Dependents (What Uses This):
//   Commands: None (top-level hook, not called by other executables)
//   Libraries: None
//   Tools: Claude Code (calls before every prompt submission)
//
// Integration Points:
//   - Called by Claude Code hook system on UserPromptSubmit event
//   - Reads PROMPT environment variable for full prompt text
//   - Outputs warning to stdout if secrets detected
//   - Always exits with code 0 (never blocks)
//
// Health Scoring
//
// Prompt submission awareness operates on Base100 scale:
//
// Activity Logging:
//   - Log prompt engagement (length only): +20
//
// Secret Detection:
//   - Check prompt for secrets: +30
//
// Warning Display:
//   - Display warning if secrets found: +20
//
// Monitoring Logging:
//   - Async log full prompt: +30
//
// Total: 100 points for complete awareness orchestration
//
// Note: NON-BLOCKING design means no false positive penalty - warnings acceptable.
// Health scoring tracks awareness quality, not blocking correctness.
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
// Standard library for I/O, environment, and string conversion.
// Hook libraries for activity logging, safety detection/display, and monitoring.

import (
	"os"      // OS interface for environment variables and exit codes
	"strconv" // String conversion for prompt length

	"hooks/lib/activity"   // Activity stream logging (engagement tracking)
	"hooks/lib/monitoring" // Monitoring logging (session history)
	"hooks/lib/safety"     // Safety detection and warning display
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Exit code for non-blocking operation.

const (
	ExitAllow = 0 // Always allow (never blocks)
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// No custom types needed - uses types from imported libraries.

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// This executable maintains no state - stateless awareness only.

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
//   └── userPromptSubmit() → main awareness orchestrator
//       └── Library functions → delegated detection and logging
//
//   Libraries (Bottom Rungs - Foundation)
//   ├── hooks/lib/activity (engagement logging)
//   ├── hooks/lib/safety (secret detection, warning display)
//   └── hooks/lib/monitoring (session history logging)
//
// Baton Flow (Execution Paths):
//
//   Entry → main()
//     ↓
//   Named Entry Point → userPromptSubmit()
//     ↓
//   Parse Environment → PROMPT variable
//     ↓
//   Log Activity → activity.LogActivity() with length
//     ↓
//   Detect Secrets → safety.ContainsLikelySecret()
//     ↓
//   Warn If Found → safety.DisplaySecretWarning()
//     ↓
//   Async Log → monitoring.LogPrompt() (goroutine)
//     ↓
//   Exit → os.Exit(0) always allow
//
// APUs (Available Processing Units):
// - 2 functions total
// - 1 orchestration function (userPromptSubmit)
// - 1 entry point (main, calls userPromptSubmit)

// ────────────────────────────────────────────────────────────────
// Prompt Submission Awareness - Secret Detection Orchestration
// ────────────────────────────────────────────────────────────────
// What This Does:
// Coordinates pre-prompt awareness by logging engagement, checking for secrets,
// warning if detected, and archiving prompt for session history. Never blocks.
//
// Why This Structure:
// Thin orchestrator keeps hook simple and VERY fast (critical for frequent execution).
// All logic delegated to libraries for testability and reusability.

// userPromptSubmit orchestrates pre-prompt awareness and logging
//
// What It Does:
//   - Gets prompt from PROMPT environment variable
//   - Logs engagement to activity stream (privacy-preserving length only)
//   - Checks prompt for likely secrets (API keys, passwords, tokens)
//   - Displays warning if secrets detected (non-blocking)
//   - Logs full prompt async to monitoring for session history
//   - Always exits 0 (never blocks prompt submission)
//
// Parameters:
//   None (reads from PROMPT environment variable)
//
// Returns:
//   None (exits with code 0 always)
//
// Health Impact:
//   +100 points for complete awareness orchestration (all phases)
//   Delegates detection and display to safety library
//
// Example usage:
//   userPromptSubmit()  // Called by main(), orchestrates entire awareness flow
func userPromptSubmit() {
	// Get prompt from environment
	prompt := os.Getenv("PROMPT")

	// Log engagement to activity stream (privacy-preserving - length only)
	promptLength := strconv.Itoa(len(prompt))
	activity.LogActivity("PromptSubmit", "length:"+promptLength, "success", 0)

	// Quick secret detection and warning (non-blocking)
	if prompt != "" && safety.ContainsLikelySecret(prompt) {
		safety.DisplaySecretWarning()
	}

	// Log full prompt async for session history (non-blocking)
	go monitoring.LogPrompt(prompt)

	// Always allow - warns but never blocks
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
//   - Compile without errors: cd ~/.claude/hooks/prompt/cmd-submit && go build -o submit submit.go
//   - Test with normal prompt (should not warn)
//   - Test with secret-like prompt (should warn but allow)
//   - Verify always exits 0 (never blocks)
//   - Verify async logging doesn't block
//   - Verify execution speed < 50ms
//
// Build Verification:
//   cd ~/.claude/hooks/prompt/cmd-submit
//   go build -o submit submit.go
//   PROMPT="hello world" ./submit  # Should not warn
//   PROMPT="sk-1234567890abcdef" ./submit  # Should warn
//
// Integration Testing:
//   Test with Claude Code UserPromptSubmit event
//   Verify warning displays correctly
//   Verify prompt still submits (never blocked)
//   Check activity logs for engagement tracking
//
// Example validation:
//
//     # Test normal prompt
//     PROMPT="explain recursion" ./submit
//     echo $?  # Should be 0 (allowed)
//
//     # Test secret-like prompt
//     PROMPT="my api key is sk-abc123" ./submit
//     # Should display warning, but still exit 0

// ────────────────────────────────────────────────────────────────
// Code Execution: Named Entry Point Pattern
// ────────────────────────────────────────────────────────────────
//
// Execution Flow:
//   1. main() called by Go runtime
//   2. main() calls userPromptSubmit() (named entry point)
//   3. userPromptSubmit() orchestrates complete awareness flow
//   4. Program exits with code 0 (always allows)
//
// Named Entry Point Benefits:
//   - Function name matches purpose (user prompt submission)
//   - Clear architectural intent (not generic "main")
//   - Enables testing without running full executable
//   - Consistent pattern across all hooks

// main is the entry point for the hook
//
// Design: Named entry point pattern for testability
//   - main() calls userPromptSubmit()
//   - userPromptSubmit() is testable (can be called independently)
//   - Semantic clarity: function name matches hook purpose
func main() {
	userPromptSubmit()
}

// ────────────────────────────────────────────────────────────────
// Cleanup and Resource Management
// ────────────────────────────────────────────────────────────────
//
// Cleanup: None required
//   - os.Exit terminates process cleanly
//   - No file handles or resources to close
//   - Async goroutine managed by runtime
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
// This executable is a Claude Code UserPromptSubmit hook (NON-BLOCKING). It runs before
// every prompt submission, providing awareness with:
//   - Privacy-preserving activity logging (engagement frequency tracking)
//   - Quick secret detection (API keys, passwords, tokens)
//   - Non-blocking warnings (user awareness without prevention)
//   - Async session history logging (full prompt archived)
//   - VERY fast execution (< 50ms target)
//
// Integration points:
//   - Called by Claude Code before EVERY prompt submission
//   - Reads PROMPT environment variable for full prompt text
//   - Outputs warning to stdout if secrets detected
//   - Always exits with code 0 (never blocks)
//   - Logs async to monitoring (doesn't block submission)
//
// Relationship to other hooks:
//   - Complements tool/cmd-pre-use (tool-level safety)
//   - Part of engagement tracking system (activity logging)
//   - Feeds session history (monitoring logs)
//   - Most frequently executed hook (every prompt)

// MODIFICATION POLICY
// ===================
//
// ⚠️  CRITICAL: This hook runs before EVERY prompt - must be VERY fast
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new secret detection patterns:
//      - Add pattern to hooks/lib/safety/detection.go (ContainsLikelySecret)
//      - Hook stays unchanged (already calls library)
//
//   ✅ Enhance warning display:
//      - Modify safety.DisplaySecretWarning() in confirmation.go
//      - Hook stays unchanged
//
//   ✅ Add new logging destinations:
//      - Add logging function call in userPromptSubmit()
//      - Keep async (don't block submission)
//
// MODIFY WITH CARE (Structural Changes):
//   ⚠️ Changing detection flow:
//      - Measure execution time (must stay < 50ms)
//      - Test false positive impact (prompts run frequently)
//      - Verify async logging doesn't block
//
//   ⚠️ Adding blocking behavior:
//      - DON'T. This hook must never block.
//      - User workflow priority over protection
//      - Awareness, not prevention
//
// NEVER MODIFY (Foundational Rails):
//   ❌ 4-block structure (METADATA → SETUP → BODY → CLOSING)
//   ❌ Named entry point pattern (main calls userPromptSubmit)
//   ❌ NON-BLOCKING principle (must ALWAYS exit 0)
//   ❌ Async logging (monitoring.LogPrompt in goroutine)
//   ❌ Privacy-preserving activity logging (length only)
//
// When making changes:
//   1. Review this modification policy
//   2. Measure execution time (must be VERY fast)
//   3. Test with frequent prompts (user workflow impact)
//   4. Update METADATA health scoring if flow changes
//   5. Document "What/Why/How" for all changes

// LADDER AND BATON FLOW
// ======================
//
// LADDER (Dependencies - Vertical):
//   submit.go depends on:
//     ↓ hooks/lib/safety (detection, display)
//     ↓ hooks/lib/activity (logging)
//     ↓ hooks/lib/monitoring (async logging)
//
// BATON (Execution - Through Layers):
//   User types prompt
//     ↓
//   Claude Code triggers hook
//     ↓
//   main() entry point
//     ↓
//   userPromptSubmit() orchestration
//     ├→ Parse PROMPT environment
//     ├→ Log activity (length)
//     ├→ Detect secrets
//     ├→ Warn if found
//     ├→ Log async (goroutine)
//     └→ Exit 0 (allow)
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
// 1. Adding new secret pattern:
//    Location: hooks/lib/safety/detection.go - ContainsLikelySecret()
//    Action: Add pattern to detection list
//    Example: Add "ghp_" prefix for GitHub tokens
//
// 2. Enhancing warning display:
//    Location: hooks/lib/safety/confirmation.go - DisplaySecretWarning()
//    Action: Modify display text or format
//    Example: Add specific secret type in warning
//
// 3. Adding new logging destination:
//    Location: userPromptSubmit() - after secret detection
//    Action: Add new async logging call
//    Example: go telemetry.LogPromptMetrics(promptLength)
//
// 4. Changing activity logging format:
//    Location: userPromptSubmit() - activity.LogActivity call
//    Action: Modify context string format
//    Example: Add additional metadata
//
// 5. Performance optimization:
//    Location: Secret detection or logging calls
//    Action: Optimize patterns, caching, or async execution
//    Example: Cache recent detection results
//
// Each modification point clearly defined for surgical precision.

// PERFORMANCE CONSIDERATIONS
// ===========================
//
// Execution Speed:
//   - Hook runs before EVERY prompt (highest frequency)
//   - Target: < 50ms total execution time
//   - Critical: Must be FASTER than pre-use hook
//   - User perceives delays > 100ms as "slow"
//
// Optimization Opportunities:
//   - Secret detection patterns optimized (early returns)
//   - Monitoring logging async (doesn't block)
//   - Activity logging lightweight (length only)
//   - No temporal context fetch (not needed)
//
// Current Bottlenecks:
//   - Secret detection regex (< 10ms typically)
//   - Activity log file I/O (~5-10ms)
//   - Display output (< 5ms)
//
// Why performance critical:
//   - Runs before EVERY prompt
//   - User workflow interrupted if slow
//   - False positives acceptable (warns frequently)
//   - Speed more important than perfect detection

// TROUBLESHOOTING GUIDE
// =====================
//
// Problem: Hook warns on normal prompts
//   Diagnosis: False positive (acceptable for this hook)
//   Solution: If frequent, check detection pattern specificity
//   Action: Review hooks/lib/safety/detection.go patterns
//
// Problem: Hook doesn't warn on obvious secrets
//   Diagnosis: False negative (detection pattern missed it)
//   Solution: Add pattern to ContainsLikelySecret()
//   Action: Update detection.go with new pattern
//
// Problem: Hook takes too long
//   Diagnosis: Slow detection or logging
//   Solution: Profile execution, identify bottleneck
//   Action: Optimize detection patterns or async logging
//
// Problem: Prompts not archived
//   Diagnosis: Monitoring logging failing
//   Solution: Check monitoring system availability
//   Action: Hook still works (async failure acceptable)
//
// Problem: Activity logs empty
//   Diagnosis: Activity logging system issue
//   Solution: Check ~/.claude/system/logs/activity/
//   Action: Non-blocking by design, warning only
//
// Debug Process:
//   1. Check hook executable exists: ~/.claude/hooks/prompt/cmd-submit/submit
//   2. Test manually: PROMPT="test" ./submit
//   3. Test with secret: PROMPT="sk-abc123" ./submit
//   4. Check exit code: echo $? (should always be 0)
//   5. Check activity logs for engagement tracking
//   6. Measure execution time: time ./submit

// RELATED COMPONENTS & DEPENDENCIES
// ==================================
//
// This hook orchestrates the following libraries:
//
// Safety (hooks/lib/safety/):
//   - detection.go: ContainsLikelySecret() - Pattern-based secret detection
//   - confirmation.go: DisplaySecretWarning() - Warning display
//
// Activity (hooks/lib/activity/):
//   - LogActivity(): Records prompt engagement (length only, privacy-preserving)
//
// Monitoring (hooks/lib/monitoring/):
//   - LogPrompt(): Archives full prompt for session history (async)
//
// Dependency Direction:
//   submit.go → hooks/lib → (no further dependencies)
//   Never reverse (libraries don't depend on hooks)

// FUTURE EXPANSIONS & ROADMAP
// ============================
//
// Planned Features:
//   ⏳ Secret type identification (API key vs password vs token)
//   ⏳ Contextual secret detection (only warn in certain contexts)
//   ⏳ User pattern learning (frequent false positives ignored)
//   ⏳ Prompt sanitization suggestions (how to remove secret)
//
// Research Areas:
//   - Machine learning for secret detection (beyond regex)
//   - Context-aware detection (code blocks vs prose)
//   - User-specific detection patterns
//   - Integration with secret management tools
//
// Known Limitations:
//   - Regex-based detection (false positives and negatives)
//   - No secret type identification (generic warning)
//   - No context awareness (warns same for all prompts)
//   - No sanitization suggestions (just warns)
//
// Enhancement Opportunities:
//   - Add secret type to warning ("API key detected")
//   - Suggest sanitization ("Use environment variable instead")
//   - Learn user patterns (ignore frequent false positives)
//   - Integrate with secret vaults (detect leaked vault secrets)

// CLOSING NOTE
// ============
//
// This hook embodies careful communication - awareness before speech, but not
// censorship. Like Psalm 141:3 requests guard over words, this hook provides
// awareness while respecting user autonomy.
//
// The NON-BLOCKING design prioritizes workflow:
//   - Awareness over prevention (warns but never blocks)
//   - Speed over perfection (false positives acceptable)
//   - Privacy-preserving logging (length only in activity)
//   - Async archival (doesn't block submission)
//
// Kingdom Technology promotes careful communication through awareness.
// User autonomy respected - warnings serve, don't control.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Maintain non-blocking principle (ALWAYS exit 0)
//   - Test with actual Claude Code integration
//   - Document all changes comprehensively (What/Why/How pattern)
//
// "Set a guard, O Lord, over my mouth; keep watch over the door of my lips." - Psalm 141:3
//
// Every prompt deserves awareness. Let it serve through wisdom.

// QUICK REFERENCE: USAGE EXAMPLES
// ================================
//
// Manual execution (normal prompt):
//   $ cd ~/.claude/hooks/prompt/cmd-submit
//   $ PROMPT="explain recursion" ./submit
//   $ echo $?  # Should be 0 (allowed, no warning)
//
// Manual execution (secret-like prompt):
//   $ PROMPT="my api key is sk-1234567890abcdef" ./submit
//   # Should display warning, but still exit 0
//
// Via Claude Code (automatic):
//   User types prompt → presses enter → Claude Code calls submit → warns if secrets → always allows
//
// Debug specific issue:
//   $ PROMPT="test secret detection: ghp_abc123" ./submit 2>&1 | tee debug.log
//
// Verify execution speed:
//   $ time PROMPT="normal prompt" ./submit
//   # Should be < 50ms

// ============================================================================
// END CLOSING
// ============================================================================
