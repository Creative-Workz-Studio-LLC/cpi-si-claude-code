// METADATA
//
// Safety Confirmation Library - CPI-SI Hook Support System
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty." - Proverbs 21:5 (WEB)
// Principle: Diligence requires confirmation before destructive operations - pause prevents hasty loss
// Anchor: "A prudent person foresees danger and takes precautions. The simpleton goes blindly on and suffers the consequences." - Proverbs 22:3 (NLT)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Hook support (mid-rung on ladder)
// Role: Provides user confirmation flows and warning displays for BLOCKING safety validation
// Paradigm: Conservative confirmation - clear communication, respectful prompting, exact match validation
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.1
// Last Modified: 2025-11-11 - Fixed displayForcePushWarning to show command parameter
//
// Version History:
//   2.0.1 (2025-11-11) - Fixed displayForcePushWarning to show cmd parameter (consistency with other warnings)
//   2.0.0 (2025-11-11) - Refactored to use system/lib/display for consistent formatting
//   1.0.0 (2024-10-24) - Initial implementation with raw fmt.Println formatting
//
// Purpose & Function
//
// Purpose: Orchestrate user confirmation flows for dangerous operations detected by safety/detection.go.
// Display contextual warnings with temporal awareness, prompt for explicit confirmation, validate
// user responses for BLOCKING decisions in tool/pre-use hook.
//
// Core Design: Confirmation orchestration layer that routes different operation types to specific
// warning displays, prompts user via stdin, validates exact response match, returns boolean results
// for BLOCKING decision making.
//
// Key Features:
//   - Bash operation confirmation (force push, hard reset, rm -rf, sudo, publish, database)
//   - Critical file write confirmation (system files, authentication, version control)
//   - Temporal context integration (long session warnings, time of day awareness)
//   - Exact match validation ("yes" full word, "y" single character)
//   - Conservative bias (ambiguous responses default to denial for safety)
//   - Non-blocking secret warning display (awareness without blocking)
//
// Philosophy: Confirmation serves protection through informed decision, not restriction through
// fear. Clear warnings explain why dangerous. Temporal context provides additional awareness.
// Exact match prevents accidental confirmation. Respectful prompting trusts user's judgment.
//
// Blocking Status
//
// Non-blocking: Confirmation functions themselves never block - they return boolean results.
// BLOCKING hooks (tool/pre-use) make final blocking decision based on returned values.
// Mitigation: Simple stdin reading with exact match validation, minimal failure points.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/safety"
//
//	// Bash operation confirmation
//	needs, allowed := safety.ConfirmBashOperation("git push --force", timeContext)
//	if needs && !allowed {
//	    os.Exit(1) // Block execution
//	}
//
//	// File write confirmation
//	needs, allowed := safety.ConfirmFileWrite("/etc/passwd", timeContext)
//	if needs && !allowed {
//	    os.Exit(1) // Block execution
//	}
//
//	// Non-blocking secret warning
//	safety.DisplaySecretWarning() // Shows warning, doesn't block
//
// Integration Pattern:
//   1. BLOCKING hook (tool/pre-use) detects dangerous operation via safety/detection.go
//   2. Hook calls ConfirmBashOperation() or ConfirmFileWrite() with temporal context
//   3. Function displays warning, prompts user, validates response
//   4. Function returns (needsConfirmation bool, allowed bool)
//   5. Hook makes blocking decision: if needs && !allowed → exit 1 (block)
//   6. No cleanup needed - confirmation is stateless user interaction
//
// Public API (in typical usage order):
//
//   Bash Operation Confirmation (BLOCKING context):
//     ConfirmBashOperation(cmd, timeContext) (needsConfirmation, allowed bool)
//       - Detects operation type, displays specific warning, gets confirmation
//       - Returns (true, true) if dangerous and user confirmed
//       - Returns (true, false) if dangerous and user denied
//       - Returns (false, true) if operation not dangerous
//
//   File Write Confirmation (BLOCKING context):
//     ConfirmFileWrite(filePath, timeContext) (needsConfirmation, allowed bool)
//       - Checks if critical path, displays warning, gets confirmation
//       - Same return semantics as ConfirmBashOperation
//
//   Secret Warning Display (NON-BLOCKING context):
//     DisplaySecretWarning()
//       - Shows warning about potential secrets, returns immediately
//       - Used by prompt/submit hook (NON-BLOCKING design)
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: bufio (stdin reading), fmt (output), os (stdin handle), strings (response trimming)
//   External: None
//   Internal: safety/detection.go (IsDangerousOperation, IsCriticalFile), system/lib/display (Warning, Info formatters)
//
// Dependents (What Uses This):
//   Hooks: tool/pre-use (BLOCKING bash/write confirmation)
//   Hooks: prompt/submit (NON-BLOCKING secret warnings)
//
// Integration Points:
//   - Ladder: Mid-rung - uses detection.go (same rung) and display (lower rung), used by hooks (higher rungs)
//   - Baton: Hook → confirmation function → warning display → stdin prompt → validation → boolean result → hook blocking decision
//   - Rails: Uses system/lib/display logging infrastructure
//
// Health Scoring
//
// Health Scoring Map (Base100):
//
//   Confirmation Orchestration (Total = 100):
//     All operations delegated to detection.go and system/lib/display which track:
//       - Detection success: +35 (from detection.go)
//       - Display formatting: +10 (from system/lib/display)
//       - User interaction: +35 (stdin read, response validation)
//       - Boolean result: +20 (correct true/false return based on confirmation)
//
// Note: This library orchestrates confirmation flows. Health tracking happens in:
// - detection.go (pattern matching)
// - system/lib/display (warning formatting)
// - This library (user interaction and response validation)
package safety

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

	"bufio"  // Buffered I/O for reading user input from stdin line by line
	"fmt"    // Formatted I/O for displaying prompts and constructing messages
	"os"     // Operating system interface for accessing stdin handle
	"strings" // String manipulation for trimming whitespace from user responses

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display" // ANSI color formatting and consistent message display (lower rung)
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - confirmation prompts and messages composed dynamically.

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// No custom types needed - confirmation functions use built-in types (string, bool).

// ────────────────────────────────────────────────────────────────
// Package-Level State - Initialization
// ────────────────────────────────────────────────────────────────
// No package-level state - confirmation is stateless user interaction.
// Each confirmation flow is independent with no shared state.

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
//   ├── ConfirmBashOperation() → uses IsDangerousOperation() + displayWarning() + confirm()
//   ├── ConfirmFileWrite() → uses IsCriticalFile() + displayCriticalFileWarning() + confirm()
//   └── DisplaySecretWarning() → uses display.Warning()
//
//   Warning Display Functions (Middle Rungs - Presentation)
//   ├── displayForcePushWarning() → uses display.Warning() + fmt.Println
//   ├── displayHardResetWarning() → uses display.Warning() + fmt.Println
//   ├── displayRecursiveDeletionWarning() → uses display.Warning() + fmt.Println
//   ├── displaySudoWarning() → uses display.Info() + fmt.Println
//   ├── displayPublishWarning() → uses display.Info() + fmt.Println
//   ├── displayDatabaseWarning() → uses display.Warning() + fmt.Println
//   └── displayCriticalFileWarning() → uses display.Warning() + fmt.Println
//
//   Helpers (Bottom Rungs - Core I/O)
//   └── confirm() → uses bufio, stdin, string validation
//
// Baton Flow (Execution Paths):
//
//   Entry → ConfirmBashOperation("git push --force", timeContext)
//     ↓
//   IsDangerousOperation(cmd) → detect if dangerous
//     ↓
//   Switch on operation type → route to specific warning
//     ↓
//   displayForcePushWarning(cmd, timeContext) → show formatted warning
//     ↓
//   confirm("Continue with force push?", "yes") → prompt user, read stdin
//     ↓
//   Validate response == "yes" → exact match check
//     ↓
//   Exit → return (needsConfirmation=true, allowed=true/false)
//
// APUs (Available Processing Units):
// - 11 functions total
// - 1 helper (stdin confirmation)
// - 7 warning display functions (operation-specific formatting)
// - 3 public APIs (bash confirmation, file confirmation, secret warning)

// ────────────────────────────────────────────────────────────────
// Helpers - User Interaction
// ────────────────────────────────────────────────────────────────
// Core confirmation mechanism for user input validation.

// confirm prompts user for confirmation and validates response.
//
// What It Does:
// Displays confirmation prompt with expected response format, reads user input
// from stdin using buffered reader, validates response matches expected value
// exactly. Supports two confirmation levels: "yes" (full word for serious
// operations) and "y" (single character for medium severity).
//
// Parameters:
//   message: The confirmation question to display to user
//   expected: The exact string required for confirmation ("yes" or "y")
//
// Returns:
//   bool: true if user response matches expected value, false otherwise
//
// Conservative Design:
// Requires exact match for safety - ambiguous responses default to denial.
// "yes" requires full word (reduces accidental confirmation). "y" accepts
// single character or capital "Y" (convenience for lower severity operations).
//
// Example Usage:
//
//	if confirm("Continue with force push?", "yes") {
//	    // User typed exactly "yes" - proceed
//	} else {
//	    // User typed anything else - deny
//	}
func confirm(message, expected string) bool {
	reader := bufio.NewReader(os.Stdin) // Create buffered reader for stdin - enables ReadString('\n')

	// Construct prompt based on expected response format
	var prompt string                                        // Prompt string to display
	if expected == "yes" {                                   // Full word confirmation (serious operations)
		prompt = fmt.Sprintf("   %s (yes/NO): ", message)    // Make clear NO is default
	} else {                                                 // Single character confirmation (medium severity)
		prompt = fmt.Sprintf("   %s (y/N): ", message)       // Make clear N is default
	}

	fmt.Print(prompt)                              // Display prompt (no newline - user types on same line)
	response, _ := reader.ReadString('\n')         // Read until newline - blocks until user presses Enter
	response = strings.TrimSpace(response)         // Remove leading/trailing whitespace and newline

	// Validate exact match (case-sensitive for "yes", case-insensitive for "y")
	return response == expected || (expected == "y" && response == "Y") // Allow Y for convenience
}

// ────────────────────────────────────────────────────────────────
// Warning Display Functions - Operation-Specific Formatting
// ────────────────────────────────────────────────────────────────
// Contextual warning displays for different dangerous operation types.
// Use system/lib/display for consistent formatting across all output.

// displayForcePushWarning shows warning for git force push operations.
//
// What It Does:
// Displays prominent warning box explaining impact (rewriting remote history),
// shows command being executed, includes temporal context if provided. Uses
// display.Warning() for consistent formatting.
//
// Parameters:
//   cmd: The git force push command (e.g., "git push --force origin main")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Emphasis:
// Force push is team-coordination-critical - rewriting shared history affects
// all team members. Warning emphasizes coordination requirement.
func displayForcePushWarning(cmd string, timeContext string) {
	fmt.Println()                                                                // Blank line for separation
	fmt.Println(display.Warning("═══════════════════════════════════════════════════════════")) // Use display library for consistent warning color
	fmt.Println(display.Warning("WARNING: FORCE PUSH DETECTED"))                 // Clear danger statement
	fmt.Println("   This will rewrite history on the remote repository.")        // Explain impact
	fmt.Printf("   Command: %s\n", cmd)                                          // Show exact command for verification
	fmt.Println("   Ensure this is intentional and coordinated with your team.") // Emphasize coordination requirement
	if timeContext != "" {                                                       // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                       // Show time/session phase context
	}
	fmt.Println(display.Warning("═══════════════════════════════════════════════════════════")) // Closing separator
	fmt.Println()                                                                // Blank line before prompt
}

// displayHardResetWarning shows warning for git hard reset operations.
//
// What It Does:
// Displays destructive operation warning explaining impact (discarding uncommitted
// changes permanently), shows command being executed, includes temporal context.
// Uses display.Warning() for consistent formatting.
//
// Parameters:
//   cmd: The git hard reset command (e.g., "git reset --hard HEAD~1")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Emphasis:
// Hard reset is data-loss-critical - discards work permanently with no undo.
func displayHardResetWarning(cmd string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Warning("DESTRUCTIVE OPERATION: Hard reset"))                    // Clear danger type
	fmt.Println("   This will discard all uncommitted changes permanently.")             // Explain permanence
	fmt.Printf("   Command: %s\n", cmd)                                                  // Show exact command
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// displayRecursiveDeletionWarning shows warning for recursive file deletion.
//
// What It Does:
// Displays destructive operation warning for rm -rf commands, shows command
// being executed, includes temporal context. Uses display.Warning() for
// consistent formatting.
//
// Parameters:
//   cmd: The rm -rf command (e.g., "rm -rf build/")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Emphasis:
// Recursive deletion is file-loss-critical - removes entire directory trees
// unrecoverably. Command display lets user verify target path.
func displayRecursiveDeletionWarning(cmd string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Warning("DESTRUCTIVE OPERATION: Recursive removal"))             // Clear danger type
	fmt.Printf("   Command: %s\n", cmd)                                                  // Show exact command for verification
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// displaySudoWarning shows warning for elevated privilege operations.
//
// What It Does:
// Displays elevated privileges notification for sudo commands, shows command
// being executed, includes temporal context. Uses display.Info() (less severe
// than Warning) since sudo is common and sometimes necessary.
//
// Parameters:
//   cmd: The sudo command (e.g., "sudo apt install package")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Level:
// Medium severity - sudo is necessary for legitimate operations. Info-level
// display (🔐) rather than warning-level (⚠️).
func displaySudoWarning(cmd string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Info("ELEVATED PRIVILEGES REQUESTED"))                           // Info level (not severe warning)
	fmt.Printf("   Command: %s\n", cmd)                                                  // Show exact command
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// displayPublishWarning shows warning for package publishing operations.
//
// What It Does:
// Displays package publishing notification for npm/cargo publish commands,
// explains impact (difficult to undo once published), shows command being
// executed, includes temporal context. Uses display.Info() since publishing
// is intentional workflow step.
//
// Parameters:
//   cmd: The publish command (e.g., "npm publish" or "cargo publish")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Level:
// Medium severity - publishing is intentional but difficult to undo. Info-level
// display (📦) with explanation of impact.
func displayPublishWarning(cmd string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Info("PACKAGE PUBLISHING"))                                      // Info level (intentional action)
	fmt.Println("   Publishing to registry - difficult to undo.")                        // Explain permanence
	fmt.Printf("   Command: %s\n", cmd)                                                  // Show exact command
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// displayDatabaseWarning shows warning for database destructive operations.
//
// What It Does:
// Displays database destructive operation warning for DROP DATABASE/TABLE
// commands, shows command being executed, includes temporal context. Uses
// display.Warning() for high severity.
//
// Parameters:
//   cmd: The database command (e.g., "DROP DATABASE production")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Emphasis:
// Database drops are data-loss-critical - permanent and often unrecoverable.
func displayDatabaseWarning(cmd string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Warning("DATABASE DESTRUCTIVE OPERATION"))                       // Clear danger type
	fmt.Printf("   Command: %s\n", cmd)                                                  // Show exact command for verification
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// displayCriticalFileWarning shows warning for critical file write operations.
//
// What It Does:
// Displays critical file write notification when writing to system files,
// configuration, or authentication locations. Shows file path, explains
// criticality, includes temporal context. Uses display.Warning() for high severity.
//
// Parameters:
//   filePath: The critical file path (e.g., "/etc/passwd" or "~/.ssh/config")
//   timeContext: Formatted temporal context string or empty string if unavailable
//
// Returns:
//   None (prints to stdout)
//
// Warning Emphasis:
// Critical file writes can break system, lock out user, or corrupt configuration.
func displayCriticalFileWarning(filePath string, timeContext string) {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Warning("CRITICAL FILE WRITE"))                                  // Clear danger type
	fmt.Printf("   Path: %s\n", filePath)                                                // Show exact path for verification
	fmt.Println("   This file is in a critical system location.")                        // Explain criticality
	if timeContext != "" {                                                               // If temporal context available
		fmt.Printf("   %s\n", timeContext)                                               // Show time/session phase context
	}
	fmt.Println()                                                                        // Blank line before prompt
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Confirmation Orchestration
// ────────────────────────────────────────────────────────────────
// High-level confirmation flows for different operation types.
// Used by BLOCKING hooks to make blocking decisions.

// ConfirmBashOperation checks if bash command is dangerous and gets confirmation.
//
// What It Does:
// Orchestrates complete confirmation flow for bash commands: checks if command
// matches dangerous patterns via IsDangerousOperation(), identifies specific
// operation type (force push, hard reset, rm -rf, sudo, publish, database),
// displays appropriate warning with temporal context, prompts user for explicit
// confirmation, validates exact response match.
//
// Parameters:
//   cmd: The bash command to check (e.g., "git push --force origin main")
//   timeContext: Formatted temporal context string from temporal library (or empty string)
//
// Returns:
//   needsConfirmation: true if operation is dangerous and required confirmation
//   allowed: true if user confirmed (or operation not dangerous), false if user denied
//
// Return Semantics:
//   (false, true) - Operation not dangerous, allowed to proceed immediately
//   (true, true) - Operation dangerous, user confirmed, allowed to proceed
//   (true, false) - Operation dangerous, user denied, BLOCK execution
//
// Conservative Bias:
// If IsDangerousOperation() returns true but no specific pattern matches (edge case),
// returns (false, true) - allow rather than false positive block.
//
// Example Usage:
//
//	needs, allowed := safety.ConfirmBashOperation("git push --force", timeContext)
//	if needs && !allowed {
//	    fmt.Println("❌ Operation cancelled")
//	    os.Exit(1) // Block execution
//	}
//	// Continue with operation
func ConfirmBashOperation(cmd string, timeContext string) (needsConfirmation bool, allowed bool) {
	// Check if dangerous using detection library
	if !IsDangerousOperation(cmd) {                                                     // If command not dangerous
		return false, true                                                              // No confirmation needed, allowed
	}

	// Identify operation type and display appropriate warning with confirmation
	// Each case shows specific warning and gets confirmation with appropriate severity
	switch {
	case strings.Contains(cmd, "git push --force") || strings.Contains(cmd, "git push -f"): // Force push detection
		displayForcePushWarning(cmd, timeContext)                                       // Show force push warning (team coordination critical)
		return true, confirm("Continue with force push?", "yes")                        // Require full "yes" (high severity)

	case strings.Contains(cmd, "git reset --hard"):                                     // Hard reset detection
		displayHardResetWarning(cmd, timeContext)                                       // Show hard reset warning (data loss critical)
		return true, confirm("Confirm hard reset?", "yes")                              // Require full "yes" (high severity)

	case strings.Contains(cmd, "rm -rf") || strings.Contains(cmd, "rm -r"):            // Recursive deletion detection
		displayRecursiveDeletionWarning(cmd, timeContext)                               // Show deletion warning (file loss critical)
		return true, confirm("Confirm deletion?", "yes")                                // Require full "yes" (high severity)

	case strings.Contains(cmd, "sudo"):                                                 // Sudo detection
		displaySudoWarning(cmd, timeContext)                                            // Show sudo warning (medium severity)
		return true, confirm("Confirm sudo operation?", "y")                            // Accept single "y" (medium severity)

	case strings.Contains(cmd, "npm publish") || strings.Contains(cmd, "cargo publish"): // Package publish detection
		displayPublishWarning(cmd, timeContext)                                         // Show publish warning (medium severity)
		return true, confirm("Confirm publish?", "yes")                                 // Require full "yes" (hard to undo)

	case strings.Contains(cmd, "DROP DATABASE") || strings.Contains(cmd, "DROP TABLE"): // Database drop detection
		displayDatabaseWarning(cmd, timeContext)                                        // Show database warning (data loss critical)
		return true, confirm("Confirm database operation?", "yes")                      // Require full "yes" (high severity)
	}

	// Edge case: IsDangerousOperation returned true but no specific pattern matched
	// Conservative bias: allow operation rather than false positive block
	return false, true                                                                  // No confirmation needed (pattern match failed), allowed
}

// ConfirmFileWrite checks if file path is critical and gets confirmation.
//
// What It Does:
// Orchestrates complete confirmation flow for file write operations: checks if
// file path is in critical location via IsCriticalFile(), displays warning with
// file path and temporal context, prompts user for explicit confirmation,
// validates exact response match.
//
// Parameters:
//   filePath: The file path to check (e.g., "/etc/hosts" or "~/.ssh/config")
//   timeContext: Formatted temporal context string from temporal library (or empty string)
//
// Returns:
//   needsConfirmation: true if file is critical and required confirmation
//   allowed: true if user confirmed (or file not critical), false if user denied
//
// Return Semantics:
//   (false, true) - File not critical, allowed to proceed immediately
//   (true, true) - File critical, user confirmed, allowed to proceed
//   (true, false) - File critical, user denied, BLOCK execution
//
// Example Usage:
//
//	needs, allowed := safety.ConfirmFileWrite("/etc/passwd", timeContext)
//	if needs && !allowed {
//	    fmt.Println("❌ Operation cancelled")
//	    os.Exit(1) // Block execution
//	}
//	// Continue with write
func ConfirmFileWrite(filePath string, timeContext string) (needsConfirmation bool, allowed bool) {
	// Check if critical using detection library
	if !IsCriticalFile(filePath) {                                                     // If file not critical
		return false, true                                                              // No confirmation needed, allowed
	}

	// Display warning and get confirmation
	displayCriticalFileWarning(filePath, timeContext)                                   // Show critical file warning
	return true, confirm("Confirm write operation?", "yes")                             // Require full "yes" (system file modification)
}

// DisplaySecretWarning shows warning for likely secret detection in prompts.
//
// What It Does:
// Displays NON-BLOCKING warning that potential secret detected in user prompt.
// Suggests user review prompt before submitting. Returns immediately without
// prompting for confirmation - NON-BLOCKING design for prompt/submit hook.
//
// Parameters:
//   None (generic warning - doesn't show secret content for privacy)
//
// Returns:
//   None (prints to stdout and returns immediately)
//
// NON-BLOCKING Design:
// Used by prompt/submit hook which runs frequently and cannot block submission.
// False positives acceptable for secret detection - warning-only approach.
// User awareness more important than blocking submission.
//
// Privacy Preservation:
// Function does not receive or display secret content - only shows generic
// warning. Calling code responsible for privacy-preserving logging (length
// only, never content).
//
// Example Usage:
//
//	if safety.ContainsLikelySecret(promptText) {
//	    safety.DisplaySecretWarning() // Show warning, don't block
//	}
//	// Continue (NON-BLOCKING)
func DisplaySecretWarning() {
	fmt.Println()                                                                        // Blank line for separation
	fmt.Println(display.Warning("Potential secret detected in prompt"))                 // Clear warning statement
	fmt.Println("   Review prompt before submitting")                                   // Suggest review action
	fmt.Println()                                                                        // Blank line after warning
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
// Code Validation: Testing Requirements
// ────────────────────────────────────────────────────────────────
// For Code Validation section explanation, see: standards/code/4-block/sections/CWS-SECTION-010-CLOSING-code-validation.md
//
// Testing Requirements:
//   - Import the library without errors: import "hooks/lib/safety"
//   - Test ConfirmBashOperation() with dangerous commands (force push, hard reset, rm -rf, sudo)
//   - Test ConfirmBashOperation() with safe commands (should return false, true immediately)
//   - Test ConfirmFileWrite() with critical paths (/etc/, /.ssh/)
//   - Test ConfirmFileWrite() with normal paths (should return false, true immediately)
//   - Test DisplaySecretWarning() (should display and return without blocking)
//   - Verify warning displays use system/lib/display formatting consistently
//   - Test confirmation validation (exact match required, ambiguous responses denied)
//   - Verify temporal context displayed when provided
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - Verify import of system/lib/display works correctly
//
// Integration Testing:
//   - Test with tool/pre-use hook (complete BLOCKING flow)
//   - Test with prompt/submit hook (NON-BLOCKING secret warning)
//   - Verify user can confirm dangerous operations (should allow when "yes" typed)
//   - Verify user can deny dangerous operations (should block when anything else typed)
//   - Verify exact match validation ("yes" requires full word, not "y" or "YES")
//   - Test temporal context integration (long session warnings, time of day)
//
// Example validation code:
//
//     // Test dangerous operation confirmation
//     needs, allowed := ConfirmBashOperation("git push --force", "")
//     if !needs {
//         t.Error("Should require confirmation for force push")
//     }
//
//     // Test safe operation (no confirmation)
//     needs, allowed := ConfirmBashOperation("git status", "")
//     if needs {
//         t.Error("Should not require confirmation for safe command")
//     }
//     if !allowed {
//         t.Error("Safe command should be allowed")
//     }
//
//     // Test critical file confirmation
//     needs, allowed := ConfirmFileWrite("/etc/passwd", "")
//     if !needs {
//         t.Error("Should require confirmation for critical file")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by BLOCKING hooks.
//
// Usage: import "hooks/lib/safety"
//
// The library is imported into BLOCKING hook packages (tool/pre-use) and NON-BLOCKING
// hooks (prompt/submit), making all exported confirmation functions available. No
// initialization needed - confirmation flows are stateless user interactions.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/safety"
//
//     func preUseHook(toolName, toolInput string) {
//         if toolName == "Bash" {
//             // Get temporal context for warning display
//             timeContext := getTemporalContext()
//
//             // Check if dangerous and get confirmation
//             needs, allowed := safety.ConfirmBashOperation(toolInput, timeContext)
//
//             if needs && !allowed {
//                 // User denied confirmation
//                 fmt.Println("❌ Operation cancelled")
//                 os.Exit(1) // BLOCK execution
//             }
//         }
//         // Operation allowed
//         os.Exit(0)
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Stdin reader: Created per confirmation, garbage collected after use
//   - Display formatters: Stateless functions, no resources to clean
//   - Temporary strings: Garbage collected automatically
//
// Graceful Shutdown:
// No shutdown needed - confirmation is stateless user interaction with no
// persistent resources. Each confirmation creates buffered reader for stdin,
// reads response, validates, returns result. Reader garbage collected after
// function returns.
//
// Memory Profile:
//   - Per-confirmation: ~1KB for buffered reader + response string
//   - No persistent allocations
//   - No goroutines, no file handles (stdin provided by OS)
//   - Total memory footprint: Negligible (temporary allocations only)
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
//   - User confirmation orchestration for dangerous operations and critical file writes
//   - Operation-specific warning displays with temporal awareness integration
//   - Exact match validation with conservative bias (ambiguous responses default to denial)
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list organized by category in typical usage order
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Rails/Ladder/Baton) explanation
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify (Extension Points):
//   ✅ Add new operation types to ConfirmBashOperation switch (follow existing pattern)
//   ✅ Add new warning display functions (follow displayWarning naming pattern)
//   ✅ Customize warning messages (explain impact clearly, include temporal context)
//   ✅ Add new confirmation functions (follow Confirm naming pattern)
//   ✅ Enhance display formatting (use system/lib/display consistently)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ ConfirmBashOperation() signature - breaks tool/pre-use hook integration
//   ⚠️ ConfirmFileWrite() signature - breaks tool/pre-use hook integration
//   ⚠️ Return value semantics (needsConfirmation, allowed) - hooks depend on this contract
//   ⚠️ Exact match validation logic - changing to loose match reduces safety
//   ⚠️ Conservative bias principle - false positives break workflow trust
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Exact match requirement (safety depends on explicit confirmation)
//   ❌ Conservative bias (ambiguous responses must default to denial)
//   ❌ Stateless nature (each confirmation independent, no shared state)
//   ❌ Package name or import path (hooks depend on "hooks/lib/safety")
//
// Validation After Modifications:
//   See "Code Validation" section in GROUP 1: CODING above for comprehensive
//   testing requirements, build verification, and integration testing procedures.
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
// - 3 public APIs orchestrate warning display and confirmation validation
// - Ladder: Uses detection.go (same rung) and system/lib/display (lower rung)
// - Baton: Hook → confirmation function → detection → warning display → stdin prompt → validation → boolean result → hook blocking decision
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// Quick reference:
// - Adding new operation type: Add case to ConfirmBashOperation switch, create displayWarning function
// - Adding new confirmation level: Add case to confirm() prompt construction (e.g., "YES" for critical)
// - Customizing warnings: Edit display functions, maintain system/lib/display usage for consistency
// - Changing confirmation prompts: Edit message strings in public API function calls to confirm()
// - Adding temporal context: Already integrated, available in all warning displays
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - No package-level state (stateless confirmation)
// - Temporary allocations only (buffered reader, response string)
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
// - Detection: <100μs (pattern matching via detection.go)
// - Display: <1ms (ANSI formatting via system/lib/display)
// - User interaction: Blocks until user presses Enter (intentional - human decision time)
// - Key consideration: User interaction dominates performance (seconds), code execution negligible
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// See BODY function docstrings above for operation-specific troubleshooting.
//
// Quick reference (details in BODY function docstrings above):
// - Confirmation not accepting "yes": Check exact match requirement (case-sensitive, full word)
// - Warning not displaying colors: Terminal doesn't support ANSI colors (system/lib/display limitation)
// - False positives: Operation detected as dangerous but shouldn't be (see detection.go patterns)
// - Temporal context not showing: Empty timeContext string passed (check temporal library availability)
// - Stdin reading hangs: Buffered reader blocks until Enter pressed (intentional design)
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard Library, External, Internal
// - Dependents (What Uses This): Commands, Libraries, Tools that depend on this
// - Integration Points: How other systems connect and interact
//
// Quick summary (details in METADATA Dependencies section above):
// - Key dependencies: safety/detection.go (pattern detection), system/lib/display (consistent formatting)
// - Primary consumers: hooks/tool/pre-use (BLOCKING confirmation), hooks/prompt/submit (NON-BLOCKING warnings)
// - Related components: system/lib/temporal (provides timeContext), hooks/lib/activity (logs confirmation results)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ System/lib/display integration - COMPLETED (v2.0.0)
//   ⏳ Configurable warning messages (JSONC config files for customization)
//   ⏳ Localization support (multi-language warning displays)
//   ⏳ Confirmation history tracking (learn user patterns)
//   ⏳ Rich terminal UI (interactive selection instead of text prompts)
//
// Research Areas:
//   - Confirmation shortcuts (whitelist frequently-confirmed operations)
//   - Context-aware prompts (different prompts per branch, project, time of day)
//   - Team coordination integration (check if operation coordinated with team)
//   - Undo suggestions (provide recovery guidance if operation was mistake)
//   - Voice confirmation (accessibility feature for hands-free confirmation)
//
// Integration Targets:
//   - Confirmation telemetry (anonymous stats on confirmation patterns)
//   - Team-wide confirmation sharing (learn from collective experience)
//   - Git hook integration (coordinate with pre-push, pre-commit hooks)
//   - Restoration layer integration (automatic recovery tracking)
//
// Known Limitations to Address:
//   - Stdin blocking (no timeout - waits indefinitely for user response)
//   - No cancellation mechanism (user must type something, Ctrl+C only option)
//   - Terminal-only (no GUI confirmation dialog support)
//   - English-only messages (no localization support)
//   - Static warning messages (no customization without code changes)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-11) - System library integration and comprehensive documentation
//         - Refactored to use system/lib/display for consistent formatting
//         - Replaced raw fmt.Println with display.Warning() and display.Info()
//         - Comprehensive 4-block template alignment with full documentation
//         - Enhanced inline comments explaining syntax and architecture
//         - Design principle: Consistent formatting across all output
//
//   1.0.0 (2024-10-24) - Initial implementation with raw formatting
//         - Bash operation confirmation (force push, hard reset, rm -rf, sudo, publish, database)
//         - Critical file write confirmation (system files, authentication, version control)
//         - Temporal context integration (long session warnings, time of day)
//         - Exact match validation ("yes" full word, "y" single character)
//         - Conservative bias principle established
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a mid-rung LADDER component providing confirmation orchestration.
// Used by BLOCKING hooks (tool/pre-use) to make informed blocking decisions through
// clear warnings and explicit user confirmation. Temporal awareness integration
// provides additional context for user decisions.
//
// Modify thoughtfully - changes here affect all BLOCKING validation flows. Exact
// match validation must be maintained (safety depends on explicit confirmation).
// Conservative bias must be preserved (ambiguous responses default to denial).
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go test -v ./...)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain exact match validation (safety-critical)
//   - Preserve conservative bias (false positives break trust)
//
// "The plans of the diligent surely lead to profit; and everyone who is hasty surely rushes to poverty." - Proverbs 21:5 (WEB)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic bash confirmation (BLOCKING):
//
//     import "hooks/lib/safety"
//
//     timeContext := getTemporalContext() // From temporal library
//     needs, allowed := safety.ConfirmBashOperation("git push --force", timeContext)
//
//     if needs && !allowed {
//         os.Exit(1) // Block execution
//     }
//     // Continue
//
// Basic file write confirmation (BLOCKING):
//
//     needs, allowed := safety.ConfirmFileWrite("/etc/passwd", timeContext)
//
//     if needs && !allowed {
//         os.Exit(1) // Block execution
//     }
//     // Continue
//
// Secret warning display (NON-BLOCKING):
//
//     if safety.ContainsLikelySecret(promptText) {
//         safety.DisplaySecretWarning() // Show warning, don't block
//     }
//     // Continue (always)
//
// Complete BLOCKING flow:
//
//     func preUseHook() {
//         toolName := os.Args[1]
//         if toolName == "Bash" {
//             cmd := os.Args[2]
//             timeContext := getTemporalContext()
//
//             needs, allowed := safety.ConfirmBashOperation(cmd, timeContext)
//
//             if needs && !allowed {
//                 fmt.Println("❌ Operation cancelled")
//                 os.Exit(1)
//             }
//         }
//         os.Exit(0)
//     }
//
// ============================================================================
// END CLOSING
// ============================================================================
