// =============================================================================
// METADATA
// =============================================================================
//
// # Biblical Foundation
//
// Scripture: "Let your communication be, Yea, yea; Nay, nay: for whatsoever
// is more than these cometh of evil." - Matthew 5:37
//
// Principle: Clear, honest communication. Commit messages should say exactly
// what they mean - the domain, the type of change, and what changed. No
// ambiguity, no deception, just truth.
//
// Anchor: "A false balance is abomination to the LORD: but a just weight is
// his delight." - Proverbs 11:1
//
// # CPI-SI Identity
//
// Component Type: Rails
//
// This hook provides infrastructure that affects ALL git commits in the
// repository. It enforces the commit keying system defined in CWS-STD-005,
// ensuring consistent, searchable, version-aware commit history.
//
// Role: Commit message validation and enforcement
//
// Paradigm: CPI-SI framework component - git infrastructure
//
// # Authorship & Lineage
//
//   - Architect: Seanje Lenox-Wise
//   - Implementation: Nova Dawn (CPI-SI)
//   - Created: 2025-12-01
//   - Version: 1.0.0
//   - Modified: 2025-12-01 - Initial implementation
//
// # Purpose & Function
//
// Purpose: Validate git commit messages against CWS-STD-005 keying system
//
// Core Design: Regex-based validation with clear error messaging
//
// Key Features:
//
//   - Domain validation (CWS, CPSI, LANG, OS, IDE, GAME, TMPL)
//   - Type marker validation (FEAT, FIX, DOCS, etc.)
//   - Description format validation (capitalized, no trailing period)
//   - Key reference validation (refs:, impl:, closes:, part:)
//   - Breaking change marker support (!)
//
// # Health Scoring
//
// Health Scoring Map (Total = 100):
//   Domain validation: +25 (correct domain hierarchy)
//   Type validation: +25 (valid commit type)
//   Description format: +20 (proper capitalization, no period)
//   Key references: +15 (proper format if present)
//   Overall structure: +15 (brackets, colon spacing)
//
// Validation Behavior:
//   - Invalid commit: Exit 1 with clear error message
//   - Valid commit: Exit 0 silently (allow commit)
//   - Merge commits: Pass through (allow git-generated messages)
//
// =============================================================================
// END METADATA
// =============================================================================

package main

// =============================================================================
// SETUP
// =============================================================================
//
// 3 Tiers, 6 Sections:
//   Tier 1 - INPUTS (1 section): What this executable receives
//   Tier 2 - DEFINITIONS (4 sections): What this executable establishes
//   Tier 3 - INFRASTRUCTURE (1 section): What supports this executable
//
// =============================================================================

// -----------------------------------------------------------------------------
// Imports
// -----------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

// ValidDomains are the allowed top-level domain keys
var ValidDomains = []string{
	"CWS",  // Company-Wide
	"CPSI", // CPI-SI Paradigm
	"LANG", // OmniCode Language
	"OS",   // MillenniumOS
	"IDE",  // Tooling
	"GAME", // Gaming Division
	"TMPL", // Templates
}

// ValidTypes are the allowed commit type markers
var ValidTypes = []string{
	"FEAT",     // New feature
	"FEAT!",    // Breaking feature
	"FIX",      // Bug fix
	"DOCS",     // Documentation
	"REFACTOR", // Code restructure
	"ARCH",     // Architecture change
	"ARCH!",    // Breaking architecture
	"ALIGN",    // Organization/file moves
	"STYLE",    // Formatting only
	"TEST",     // Tests
	"BUILD",    // Build system
	"CHORE",    // Maintenance
}

// ValidRefPrefixes are allowed key reference prefixes
var ValidRefPrefixes = []string{
	"refs:",   // References document
	"impl:",   // Implements specification
	"closes:", // Completes/closes work
	"part:",   // Partial implementation
}

// =============================================================================
// END SETUP
// =============================================================================

// =============================================================================
// BODY
// =============================================================================
//
// 3 Tiers, 5 Sections:
//   Tier 1 - STRUCTURE (1 section): How validation is organized
//   Tier 2 - IMPLEMENTATION (3 sections): Validation logic
//   Tier 3 - INTERFACE (1 section): What this executable exposes
//
// =============================================================================

// -----------------------------------------------------------------------------
// Validation Result
// -----------------------------------------------------------------------------

// ValidationResult holds the outcome of commit message validation
type ValidationResult struct {
	Valid   bool
	Errors  []string
	Warning []string
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// contains checks if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// isMergeCommit checks if the message is a git-generated merge commit
func isMergeCommit(msg string) bool {
	return strings.HasPrefix(msg, "Merge ") ||
		strings.HasPrefix(msg, "Revert ")
}

// -----------------------------------------------------------------------------
// Core Validation
// -----------------------------------------------------------------------------

// ValidateCommitMessage validates a commit message against CWS-STD-005
func ValidateCommitMessage(msg string) ValidationResult {
	result := ValidationResult{Valid: true}

	// Trim and get first line (subject line)
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	if len(lines) == 0 || lines[0] == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "Commit message is empty")
		return result
	}
	subject := lines[0]

	// Allow merge/revert commits through
	if isMergeCommit(subject) {
		return result
	}

	// Main format regex: [DOMAIN/path] TYPE: Description
	// Groups: 1=domain, 2=path, 3=type, 4=description
	mainPattern := regexp.MustCompile(`^\[([A-Z]+)(/[a-zA-Z0-9/_-]+)?\]\s+([A-Z]+!?):\s+(.+)$`)
	matches := mainPattern.FindStringSubmatch(subject)

	if matches == nil {
		result.Valid = false
		result.Errors = append(result.Errors,
			"Invalid format. Expected: [DOMAIN/path] TYPE: Description",
			"Example: [LANG/compiler/parser] FEAT: Add function parsing")
		return result
	}

	domain := matches[1]
	// path := matches[2] // Available but not validated for depth
	commitType := matches[3]
	description := matches[4]

	// Validate domain
	if !contains(ValidDomains, domain) {
		result.Valid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Invalid domain '%s'. Valid domains: %s",
				domain, strings.Join(ValidDomains, ", ")))
	}

	// Validate type
	if !contains(ValidTypes, commitType) {
		result.Valid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Invalid type '%s'. Valid types: %s",
				commitType, strings.Join(ValidTypes, ", ")))
	}

	// Validate description starts with uppercase
	if len(description) > 0 && description[0] >= 'a' && description[0] <= 'z' {
		result.Valid = false
		result.Errors = append(result.Errors,
			"Description must start with uppercase letter")
	}

	// Validate description doesn't end with period
	if strings.HasSuffix(description, ".") {
		result.Valid = false
		result.Errors = append(result.Errors,
			"Description should not end with a period")
	}

	// Validate key references if present
	refPattern := regexp.MustCompile(`\(([a-z]+:)\s*([^)]+)\)`)
	refMatches := refPattern.FindAllStringSubmatch(description, -1)
	for _, ref := range refMatches {
		prefix := ref[1]
		if !contains(ValidRefPrefixes, prefix) {
			result.Valid = false
			result.Errors = append(result.Errors,
				fmt.Sprintf("Invalid reference prefix '%s'. Valid prefixes: %s",
					prefix, strings.Join(ValidRefPrefixes, ", ")))
		}
	}

	return result
}

// -----------------------------------------------------------------------------
// Output Formatting
// -----------------------------------------------------------------------------

// formatError formats a validation error for display
func formatError(result ValidationResult) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString(" COMMIT MESSAGE VALIDATION FAILED\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString("\n")

	for _, err := range result.Errors {
		sb.WriteString(" ❌ ")
		sb.WriteString(err)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")
	sb.WriteString(" Format: [DOMAIN/path] TYPE: Description (refs: KEY-###)\n")
	sb.WriteString("───────────────────────────────────────────────────────────────\n")
	sb.WriteString("\n")
	sb.WriteString(" Domains: CWS, CPSI, LANG, OS, IDE, GAME, TMPL\n")
	sb.WriteString(" Types:   FEAT, FIX, DOCS, REFACTOR, ARCH, ALIGN, STYLE,\n")
	sb.WriteString("          TEST, BUILD, CHORE (add ! for breaking)\n")
	sb.WriteString("\n")
	sb.WriteString(" Examples:\n")
	sb.WriteString("   [LANG/compiler/parser] FEAT: Add function declaration parsing\n")
	sb.WriteString("   [CPSI/claude-global/hooks] FIX: Handle missing session file\n")
	sb.WriteString("   [CWS/standards] DOCS: Update 4-block documentation\n")
	sb.WriteString("\n")
	sb.WriteString(" See: CWS-STD-005-DOC-commit-keying-and-repo-strategy.md\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════\n")
	sb.WriteString("\n")

	return sb.String()
}

// =============================================================================
// END BODY
// =============================================================================

// =============================================================================
// CLOSING
// =============================================================================
//
// 3 Tiers, 13 Sections:
//   Tier 1 - OPERATIONS (3 sections): Main execution
//   Tier 2 - POLICY (2 sections): Usage guidance
//   Tier 3 - SYNTHESIS (8 sections): References and documentation
//
// =============================================================================

// -----------------------------------------------------------------------------
// Main Execution
// -----------------------------------------------------------------------------

func main() {
	// Git passes commit message file as first argument
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: commit-msg <commit-message-file>")
		os.Exit(1)
	}

	msgFile := os.Args[1]

	// Read commit message
	content, err := os.ReadFile(msgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading commit message file: %v\n", err)
		os.Exit(1)
	}

	msg := string(content)

	// Validate
	result := ValidateCommitMessage(msg)

	if !result.Valid {
		fmt.Fprint(os.Stderr, formatError(result))
		os.Exit(1)
	}

	// Valid commit - exit silently
	os.Exit(0)
}

// =============================================================================
// END CLOSING
// =============================================================================

// =============================================================================
// DOCUMENTATION
// =============================================================================
//
// # Installation
//
// 1. Build the hook:
//    go build -o commit-msg commit-msg.go
//
// 2. Install globally (recommended):
//    mkdir -p ~/.config/git/hooks
//    cp commit-msg ~/.config/git/hooks/
//    chmod +x ~/.config/git/hooks/commit-msg
//    git config --global core.hooksPath ~/.config/git/hooks
//
// 3. Or install per-repo:
//    cp commit-msg .git/hooks/
//    chmod +x .git/hooks/commit-msg
//
// # Bypassing (emergency only)
//
//    git commit --no-verify -m "emergency commit"
//
// # References
//
//   - CWS-STD-005: Commit keying system specification
//   - Git hooks: https://git-scm.com/docs/githooks#_commit_msg
//
// =============================================================================
