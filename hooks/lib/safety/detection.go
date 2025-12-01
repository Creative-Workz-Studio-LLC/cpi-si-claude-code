// METADATA
//
// Safety Detection Library - CPI-SI Hook Support System
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The prudent see danger and take refuge, but the simple keep going and pay the penalty." - Proverbs 27:12 (WEB)
// Principle: Wisdom recognizes danger before harm occurs - protection through early detection
// Anchor: "A prudent man sees danger and takes refuge; but the simple pass on, and suffer for it." - Proverbs 22:3 (WEB)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Hook support (mid-rung on ladder)
// Role: Detects potentially dangerous operations requiring user confirmation
// Paradigm: Conservative protection - better to allow than falsely block
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-11 - Added configuration loading and comprehensive documentation
//
// Version History:
//   2.0.0 (2025-11-11) - Added configuration loading with graceful fallbacks
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded patterns
//
// Purpose & Function
//
// Purpose: Detect potentially dangerous operations (destructive commands, critical file
// modifications, likely secrets) to enable BLOCKING hooks to request user confirmation
// before proceeding. Protects against accidental data loss, credential exposure, and
// system damage.
//
// Core Design: Pattern-matching detection library with three primary functions:
// 1. Dangerous operation detection (force push, hard reset, rm -rf, sudo, etc.)
// 2. Critical path detection (system files, authentication, version control)
// 3. Secret detection (API keys, tokens, private keys)
//
// Key Features:
//   - Configurable patterns (load from system/data/config/safety/)
//   - Graceful fallback to hardcoded defaults if config unavailable
//   - Simple substring matching (fast, reliable)
//   - Conservative bias (when uncertain, allow - don't falsely block)
//   - Privacy-preserving secret detection (never logs actual content)
//
// Philosophy: Detection serves protection, not paranoia. Conservative bias means
// better to allow questionable operations than falsely block legitimate work.
// User's workflow continues unless clear danger detected.
//
// Blocking Status
//
// Non-blocking: Detection itself never blocks - it reports findings to BLOCKING hooks
// which make final decision about confirmation requirements.
// Mitigation: All detection failures gracefully fall back to hardcoded safe defaults.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/safety"
//
//	// Check if command is dangerous
//	if safety.IsDangerousOperation("git push --force") {
//	    // Request user confirmation
//	}
//
//	// Check if file path is critical
//	if safety.IsCriticalFile("/etc/passwd") {
//	    // Request user confirmation
//	}
//
//	// Check if text contains likely secrets
//	if safety.ContainsLikelySecret("sk-1234567890abcdef") {
//	    // Warn user about potential secret
//	}
//
// Integration Pattern:
//   1. Import safety library into BLOCKING hooks (tool/pre-use, prompt/submit)
//   2. Pass command string or file path to detection functions
//   3. Use boolean result to decide whether to request confirmation
//   4. No cleanup needed - detection is stateless
//
// Public API (in typical usage order):
//
//   Dangerous Operation Detection (command validation):
//     IsDangerousOperation(cmd string) bool - Detects dangerous command patterns
//
//   Critical Path Detection (file path validation):
//     IsCriticalFile(filePath string) bool - Detects critical file locations
//
//   Secret Detection (content scanning):
//     ContainsLikelySecret(text string) bool - Detects likely secret patterns
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: strings, encoding/json, os, path/filepath
//   External: None
//   Internal: None (pure detection logic)
//
// Dependents (What Uses This):
//   Hooks: tool/pre-use (BLOCKING pre-tool validation)
//   Hooks: prompt/submit (NON-BLOCKING secret detection)
//
// Integration Points:
//   - Ladder: Mid-rung - pure library, used by hooks (higher rungs)
//   - Baton: Receives command/path/text, returns boolean detection result
//   - Rails: No logging infrastructure (stateless detection)
//
// Health Scoring
//
// Health Scoring Map (Base100):
//
//   Configuration Loading (Total = 30):
//     - Dangerous patterns loaded: +10
//     - Critical paths loaded: +10
//     - Secret patterns loaded: +10
//     - Fallback to defaults: 0 (neutral - still functional)
//     - Config load failure: -10 (per failed config)
//
//   Detection Operations (Total = 70):
//     - Successful pattern match: +35
//     - Successful no-match: +35
//     - Detection error (panic): -70
//
// Note: Detection is simple substring matching - errors extremely rare.
// Health tracking focuses on configuration loading success.
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

	"encoding/json" // JSON parsing for configuration files
	"os"            // File operations for config loading
	"path/filepath" // Path manipulation for config file locations
	"strings"       // String operations for pattern matching
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// No constants needed - patterns loaded from configuration or fallback arrays.

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first, then composed structures. This organization
// reveals what depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Configuration Types ---
// Structures matching JSONC configuration file formats.

// PatternCategory represents a category of patterns with metadata.
//
// Maps to dangerous-patterns.jsonc structure where patterns are organized
// by category (git_destructive, filesystem_destructive, etc.). Each category
// contains patterns array, description, severity, and reason.
//
// Example from config:
//
//	"git_destructive": {
//	  "patterns": ["git push --force", "git reset --hard"],
//	  "description": "Git operations that can destroy work",
//	  "severity": "high",
//	  "reason": "Can permanently lose work"
//	}
type PatternCategory struct {
	Patterns    []string `json:"patterns"`    // Array of pattern strings to match
	Description string   `json:"description"` // Human-readable category description
	Severity    string   `json:"severity"`    // Severity level: critical/high/medium/low
	Reason      string   `json:"reason"`      // Why this category is dangerous
}

// DangerousPatternsConfig represents the complete dangerous-patterns.jsonc structure.
//
// Top-level configuration object containing metadata and pattern categories.
// Loaded during init() with graceful fallback to hardcoded defaults.
//
// File location: $HOME/.claude/cpi-si/system/data/config/safety/dangerous-patterns.jsonc
type DangerousPatternsConfig struct {
	Metadata struct {
		Name        string `json:"name"`         // Config name
		Description string `json:"description"`  // Config purpose
		Version     string `json:"version"`      // Semantic version
		LastUpdated string `json:"last_updated"` // ISO date
		Author      string `json:"author"`       // Who created this
	} `json:"metadata"` // Configuration metadata
	Patterns map[string]PatternCategory `json:"patterns"` // Category name -> patterns
}

// PathCategory represents a category of file paths with metadata.
//
// Maps to critical-paths.jsonc structure where paths are organized
// by category (system_config, authentication, version_control, etc.).
//
// Example from config:
//
//	"authentication": {
//	  "patterns": ["/.ssh/", "/.gnupg/", "/etc/shadow"],
//	  "description": "Authentication and security configuration",
//	  "severity": "critical",
//	  "reason": "Compromise can lock you out or expose credentials"
//	}
type PathCategory struct {
	Patterns    []string `json:"patterns"`    // Array of path patterns to match
	Description string   `json:"description"` // Human-readable category description
	Severity    string   `json:"severity"`    // Severity level: critical/high/medium/low
	Reason      string   `json:"reason"`      // Why this category is critical
	Platform    string   `json:"platform"`    // Target platform: linux/unix/windows/all
}

// CriticalPathsConfig represents the complete critical-paths.jsonc structure.
//
// Top-level configuration object containing metadata and path categories.
// Loaded during init() with graceful fallback to hardcoded defaults.
//
// File location: $HOME/.claude/cpi-si/system/data/config/safety/critical-paths.jsonc
type CriticalPathsConfig struct {
	Metadata struct {
		Name        string `json:"name"`         // Config name
		Description string `json:"description"`  // Config purpose
		Version     string `json:"version"`      // Semantic version
		LastUpdated string `json:"last_updated"` // ISO date
		Author      string `json:"author"`       // Who created this
	} `json:"metadata"` // Configuration metadata
	Paths map[string]PathCategory `json:"paths"` // Category name -> paths
}

// SecretCategory represents a category of secret patterns with metadata.
//
// Maps to secret-patterns.jsonc structure where patterns are organized
// by category (api_keys, github_tokens, private_keys, etc.).
//
// Example from config:
//
//	"api_keys": {
//	  "patterns": ["sk-", "sk_live_", "AKIA"],
//	  "description": "Common API key prefixes",
//	  "severity": "high",
//	  "reason": "API keys grant programmatic access",
//	  "services": ["OpenAI", "Stripe", "AWS"]
//	}
type SecretCategory struct {
	Patterns    []string `json:"patterns"`    // Array of secret patterns to match
	Description string   `json:"description"` // Human-readable category description
	Severity    string   `json:"severity"`    // Severity level: critical/high/medium
	Reason      string   `json:"reason"`      // Why this category is sensitive
	Services    []string `json:"services"`    // Affected services/platforms
}

// SecretPatternsConfig represents the complete secret-patterns.jsonc structure.
//
// Top-level configuration object containing metadata and secret pattern categories.
// Loaded during init() with graceful fallback to hardcoded defaults.
//
// File location: $HOME/.claude/cpi-si/system/data/config/safety/secret-patterns.jsonc
type SecretPatternsConfig struct {
	Metadata struct {
		Name        string `json:"name"`         // Config name
		Description string `json:"description"`  // Config purpose
		Version     string `json:"version"`      // Semantic version
		LastUpdated string `json:"last_updated"` // ISO date
		Author      string `json:"author"`       // Who created this
	} `json:"metadata"` // Configuration metadata
	Patterns map[string]SecretCategory `json:"patterns"` // Category name -> secret patterns
}

// ────────────────────────────────────────────────────────────────
// Package-Level State - Initialization
// ────────────────────────────────────────────────────────────────
// Package-level state initialized once at import time. Configuration loading
// happens in init(), making patterns immediately available to all functions.
//
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md

//--- Configuration State ---
// Loaded configurations available to all detection functions.

var (
	dangerousConfig  *DangerousPatternsConfig // Loaded dangerous patterns config
	criticalConfig   *CriticalPathsConfig     // Loaded critical paths config
	secretConfig     *SecretPatternsConfig    // Loaded secret patterns config
	configLoaded     bool                     // True if all configs loaded successfully
)

// Hardcoded fallback patterns - used if config loading fails.
// These are the original patterns from v1.0.0, ensuring library
// continues functioning even if configuration files are missing.

var (
	// fallbackDangerousPatterns contains minimal set of dangerous operations.
	// Used only if dangerous-patterns.jsonc cannot be loaded.
	fallbackDangerousPatterns = []string{
		"git push --force",
		"git push -f",
		"git reset --hard",
		"rm -rf",
		"rm -r",
		"sudo",
		"npm publish",
		"cargo publish",
		"DROP DATABASE",
		"DROP TABLE",
	}

	// fallbackCriticalPaths contains minimal set of critical locations.
	// Used only if critical-paths.jsonc cannot be loaded.
	fallbackCriticalPaths = []string{
		"/etc/",
		"/.git/config",
		"/.ssh/",
		"/boot/",
	}

	// fallbackSecretPatterns contains minimal set of secret indicators.
	// Used only if secret-patterns.jsonc cannot be loaded.
	fallbackSecretPatterns = []string{
		"sk-",           // OpenAI keys
		"ghp_",          // GitHub tokens
		"xox",           // Slack tokens
		"AKIA",          // AWS keys
		"BEGIN PRIVATE", // Private keys
	}
)

func init() {
	// Load configuration files from standard location.
	// Gracefully falls back to hardcoded patterns if loading fails.

	homeDir := os.Getenv("HOME") // Get user home directory for config path
	if homeDir == "" {
		homeDir = "/home/seanje-lenox-wise" // Fallback to known home if $HOME unset
	}

	configBase := filepath.Join(homeDir, ".claude/cpi-si/system/data/config/safety") // Base path for safety configs

	// Load each configuration file independently - partial success is acceptable
	dangerousConfig = loadDangerousPatterns(filepath.Join(configBase, "dangerous-patterns.jsonc"))
	criticalConfig = loadCriticalPaths(filepath.Join(configBase, "critical-paths.jsonc"))
	secretConfig = loadSecretPatterns(filepath.Join(configBase, "secret-patterns.jsonc"))

	// Set configLoaded flag if all three configs loaded successfully
	configLoaded = (dangerousConfig != nil && criticalConfig != nil && secretConfig != nil)
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
//   Public APIs (Top Rungs - Exported Interface)
//   ├── IsDangerousOperation() → uses matchesAnyPattern() + config/fallback
//   ├── IsCriticalFile() → uses matchesAnyPattern() + config/fallback
//   └── ContainsLikelySecret() → uses matchesAnyPattern() + config/fallback
//
//   Helpers (Bottom Rungs - Internal Utilities)
//   ├── loadDangerousPatterns() → reads config, parses JSON
//   ├── loadCriticalPaths() → reads config, parses JSON
//   ├── loadSecretPatterns() → reads config, parses JSON
//   ├── stripJSONCComments() → removes // comments from JSONC
//   └── matchesAnyPattern() → pure string matching function
//
// Baton Flow (Execution Paths):
//
//   Entry → IsDangerousOperation("git push --force")
//     ↓
//   Check configLoaded → get patterns from config or fallback
//     ↓
//   matchesAnyPattern(cmd, patterns)
//     ↓
//   Exit → return true (matches) or false (no match)
//
// APUs (Available Processing Units):
// - 8 functions total
// - 5 helpers (config loading, JSON comment stripping, pattern matching)
// - 3 public APIs (dangerous operation, critical file, secret detection)

// ────────────────────────────────────────────────────────────────
// Helpers - Configuration Loading
// ────────────────────────────────────────────────────────────────
// Internal functions for loading and parsing configuration files.
// Handle JSONC comment stripping and graceful error handling.

// loadDangerousPatterns loads dangerous-patterns.jsonc configuration.
//
// What It Does:
// Reads dangerous-patterns.jsonc from disk, strips JSONC comments,
// parses JSON into DangerousPatternsConfig struct. Returns nil on
// any error (file missing, JSON invalid, etc.).
//
// Parameters:
//   path: Absolute path to dangerous-patterns.jsonc file
//
// Returns:
//   *DangerousPatternsConfig if successful, nil if any error occurs
//
// Error Handling:
// All errors result in nil return - caller uses fallback patterns.
// Silent failure is intentional - detection continues with defaults.
func loadDangerousPatterns(path string) *DangerousPatternsConfig {
	data, err := os.ReadFile(path) // Read entire file into memory - configs are small (<10KB)
	if err != nil {
		return nil // File doesn't exist or can't be read - use fallback
	}

	cleanJSON := stripJSONCComments(string(data)) // Remove // comments to make valid JSON

	var config DangerousPatternsConfig                   // Allocate config struct
	if err := json.Unmarshal([]byte(cleanJSON), &config); err != nil {
		return nil // JSON parsing failed - use fallback
	}

	return &config // Successfully loaded and parsed
}

// loadCriticalPaths loads critical-paths.jsonc configuration.
//
// What It Does:
// Reads critical-paths.jsonc from disk, strips JSONC comments,
// parses JSON into CriticalPathsConfig struct. Returns nil on
// any error (file missing, JSON invalid, etc.).
//
// Parameters:
//   path: Absolute path to critical-paths.jsonc file
//
// Returns:
//   *CriticalPathsConfig if successful, nil if any error occurs
//
// Error Handling:
// All errors result in nil return - caller uses fallback patterns.
// Silent failure is intentional - detection continues with defaults.
func loadCriticalPaths(path string) *CriticalPathsConfig {
	data, err := os.ReadFile(path) // Read entire file into memory - configs are small (<10KB)
	if err != nil {
		return nil // File doesn't exist or can't be read - use fallback
	}

	cleanJSON := stripJSONCComments(string(data)) // Remove // comments to make valid JSON

	var config CriticalPathsConfig                       // Allocate config struct
	if err := json.Unmarshal([]byte(cleanJSON), &config); err != nil {
		return nil // JSON parsing failed - use fallback
	}

	return &config // Successfully loaded and parsed
}

// loadSecretPatterns loads secret-patterns.jsonc configuration.
//
// What It Does:
// Reads secret-patterns.jsonc from disk, strips JSONC comments,
// parses JSON into SecretPatternsConfig struct. Returns nil on
// any error (file missing, JSON invalid, etc.).
//
// Parameters:
//   path: Absolute path to secret-patterns.jsonc file
//
// Returns:
//   *SecretPatternsConfig if successful, nil if any error occurs
//
// Error Handling:
// All errors result in nil return - caller uses fallback patterns.
// Silent failure is intentional - detection continues with defaults.
func loadSecretPatterns(path string) *SecretPatternsConfig {
	data, err := os.ReadFile(path) // Read entire file into memory - configs are small (<10KB)
	if err != nil {
		return nil // File doesn't exist or can't be read - use fallback
	}

	cleanJSON := stripJSONCComments(string(data)) // Remove // comments to make valid JSON

	var config SecretPatternsConfig                      // Allocate config struct
	if err := json.Unmarshal([]byte(cleanJSON), &config); err != nil {
		return nil // JSON parsing failed - use fallback
	}

	return &config // Successfully loaded and parsed
}

// stripJSONCComments removes // comments from JSONC text.
//
// What It Does:
// Splits text by newlines, filters out lines starting with //, rejoins.
// Simple approach sufficient for our config files which use only // comments.
//
// Parameters:
//   jsonc: JSONC text with // comments
//
// Returns:
//   Valid JSON text with comments removed
//
// Limitations:
// Does not handle /* */ block comments or // inside strings.
// Our config files only use // comments at start of lines.
func stripJSONCComments(jsonc string) string {
	lines := strings.Split(jsonc, "\n") // Split into individual lines
	var cleaned []string                // Accumulate non-comment lines

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)               // Remove leading/trailing whitespace
		if !strings.HasPrefix(trimmed, "//") {           // Keep if not a comment line
			cleaned = append(cleaned, line)              // Preserve original spacing
		}
	}

	return strings.Join(cleaned, "\n") // Rejoin with newlines
}

// matchesAnyPattern checks if text contains any pattern from list.
//
// What It Does:
// Iterates through pattern list, returns true if text contains any pattern.
// Uses simple substring matching via strings.Contains().
//
// Parameters:
//   text: Text to search (command, file path, prompt content)
//   patterns: Array of pattern strings to search for
//
// Returns:
//   true if any pattern found in text, false if no matches
//
// Performance:
// O(n*m) where n=pattern count, m=average pattern length.
// Fast enough for small pattern lists (typically <50 patterns).
func matchesAnyPattern(text string, patterns []string) bool {
	for _, pattern := range patterns { // Iterate through all patterns
		if strings.Contains(text, pattern) { // Substring match - simple and reliable
			return true // Found a match - no need to continue
		}
	}
	return false // No patterns matched
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Detection Functions
// ────────────────────────────────────────────────────────────────
// Exported functions providing core detection capabilities.
// Used by BLOCKING hooks to identify dangerous operations.

// IsDangerousOperation checks if command contains dangerous patterns.
//
// What It Does:
// Searches command string for patterns indicating potentially dangerous
// operations (force push, hard reset, rm -rf, sudo, etc.). Uses configured
// patterns if available, falls back to hardcoded defaults if config missing.
//
// Parameters:
//   cmd: Command string to check (e.g., "git push --force origin main")
//
// Returns:
//   true if command matches dangerous pattern, false otherwise
//
// Conservative Bias:
// When uncertain, returns false (allow). False positives disrupt workflow,
// so detection errs on side of permissiveness.
//
// Example Usage:
//
//	if IsDangerousOperation("git push --force") {
//	    // Request user confirmation
//	}
func IsDangerousOperation(cmd string) bool {
	var patterns []string // Pattern list to search

	if configLoaded && dangerousConfig != nil {
		// Use configured patterns - flatten all categories into single list
		for _, category := range dangerousConfig.Patterns {
			patterns = append(patterns, category.Patterns...) // Spread category patterns into combined list
		}
	} else {
		// Config unavailable - use hardcoded fallback
		patterns = fallbackDangerousPatterns
	}

	return matchesAnyPattern(cmd, patterns) // Delegate to pattern matcher
}

// IsCriticalFile checks if file path is in a critical location.
//
// What It Does:
// Searches file path for patterns indicating critical system locations
// (/etc/, /.ssh/, /boot/, git config, etc.). Uses configured paths if
// available, falls back to hardcoded defaults if config missing.
//
// Parameters:
//   filePath: File path to check (e.g., "/etc/passwd" or "~/.ssh/id_rsa")
//
// Returns:
//   true if path matches critical location, false otherwise
//
// Conservative Bias:
// When uncertain, returns false (allow). False positives disrupt workflow,
// so detection errs on side of permissiveness.
//
// Example Usage:
//
//	if IsCriticalFile("/etc/shadow") {
//	    // Request user confirmation
//	}
func IsCriticalFile(filePath string) bool {
	var paths []string // Path pattern list to search

	if configLoaded && criticalConfig != nil {
		// Use configured paths - flatten all categories into single list
		for _, category := range criticalConfig.Paths {
			paths = append(paths, category.Patterns...) // Spread category patterns into combined list
		}
	} else {
		// Config unavailable - use hardcoded fallback
		paths = fallbackCriticalPaths
	}

	return matchesAnyPattern(filePath, paths) // Delegate to pattern matcher
}

// ContainsLikelySecret performs basic check for obvious secret patterns.
//
// What It Does:
// Searches text for patterns indicating likely secrets (API keys, tokens,
// private keys, etc.). Uses configured patterns if available, falls back
// to hardcoded defaults if config missing.
//
// Privacy Preservation:
// This function only detects - it NEVER logs or stores actual secret content.
// Calling code should log detection events with length only, never content.
//
// Parameters:
//   text: Text to scan for secrets (prompt content, file content, etc.)
//
// Returns:
//   true if text likely contains secrets, false otherwise
//
// False Positives:
// Some patterns (like "eyJ" for JWT) are common in non-secret base64 data.
// Conservative bias: better to warn unnecessarily than miss actual secrets.
//
// Example Usage:
//
//	if ContainsLikelySecret(promptText) {
//	    // Warn user (log length only, NEVER content)
//	}
func ContainsLikelySecret(text string) bool {
	var patterns []string // Secret pattern list to search

	if configLoaded && secretConfig != nil {
		// Use configured patterns - flatten all categories into single list
		for _, category := range secretConfig.Patterns {
			patterns = append(patterns, category.Patterns...) // Spread category patterns into combined list
		}
	} else {
		// Config unavailable - use hardcoded fallback
		patterns = fallbackSecretPatterns
	}

	return matchesAnyPattern(text, patterns) // Delegate to pattern matcher
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
//   - Test IsDangerousOperation() with known dangerous commands
//   - Test IsCriticalFile() with known critical paths
//   - Test ContainsLikelySecret() with known secret patterns
//   - Verify false positives are minimal (conservative bias)
//   - Test graceful fallback when configs missing
//   - Confirm no panics on malformed input
//   - Run: go test -v ./... (when tests exist)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - No circular imports (safety has no internal dependencies)
//
// Integration Testing:
//   - Test with tool/pre-use hook (BLOCKING dangerous operations)
//   - Test with prompt/submit hook (NON-BLOCKING secret detection)
//   - Verify conservative bias (false negatives acceptable, false positives disruptive)
//   - Test config loading from actual config files
//   - Verify fallback patterns work when configs missing
//
// Example validation code:
//
//     // Test dangerous operation detection
//     if !IsDangerousOperation("git push --force") {
//         t.Error("Failed to detect force push")
//     }
//     if IsDangerousOperation("git status") {
//         t.Error("False positive on safe command")
//     }
//
//     // Test critical file detection
//     if !IsCriticalFile("/etc/passwd") {
//         t.Error("Failed to detect critical system file")
//     }
//     if IsCriticalFile("/home/user/document.txt") {
//         t.Error("False positive on normal file")
//     }
//
//     // Test secret detection
//     if !ContainsLikelySecret("sk-1234567890abcdef") {
//         t.Error("Failed to detect OpenAI key")
//     }
//     if ContainsLikelySecret("normal text content") {
//         t.Error("False positive on normal text")
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
// The library is imported into BLOCKING hook packages (tool/pre-use, prompt/submit),
// making all exported detection functions available. Configuration loading happens
// automatically during init() - no setup required from calling code.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/safety"
//
//     func preUseHook(toolName, toolInput string) {
//         if toolName == "Bash" {
//             // Check if command is dangerous
//             if safety.IsDangerousOperation(toolInput) {
//                 // Request user confirmation
//                 confirmed := requestConfirmation("Dangerous operation detected. Proceed?")
//                 if !confirmed {
//                     os.Exit(1) // Block execution
//                 }
//             }
//         }
//         // Allow execution
//         os.Exit(0)
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Configuration structs: Allocated once in init(), never freed (package lifetime)
//   - Pattern arrays: Immutable after init(), no cleanup needed
//   - Temporary strings: Garbage collected automatically
//
// Graceful Shutdown:
// No shutdown needed - library is stateless with no persistent resources.
// Configuration loaded once at import time, detection functions are pure
// operations on immutable data.
//
// Memory Profile:
//   - Config files: ~10-15KB each (3 files = ~30-45KB total)
//   - Parsed structs: ~50-100KB in memory (negligible)
//   - No goroutines, no file handles, no network connections
//   - Total memory footprint: <100KB persistent
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
//   - Pattern-based detection library for dangerous operations, critical paths, and secrets
//   - Configuration-driven with graceful fallback to hardcoded defaults
//   - Conservative bias (better to allow than falsely block)
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
//   ✅ Add new patterns to configuration files (dangerous-patterns.jsonc, critical-paths.jsonc, secret-patterns.jsonc)
//   ✅ Add new categories to existing config files (follow existing pattern structure)
//   ✅ Extend fallback pattern arrays (add patterns, update existing)
//   ✅ Improve pattern matching logic (enhance matchesAnyPattern function)
//   ✅ Add new detection functions (follow existing naming pattern)
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Public API function signatures - breaks all calling hooks (tool/pre-use, prompt/submit)
//   ⚠️ Config struct definitions - must match JSONC schema exactly
//   ⚠️ Conservative bias principle - making detection too aggressive disrupts workflow
//   ⚠️ Privacy preservation - never log secret content, only metadata
//   ⚠️ Boolean return types - hooks depend on true/false semantics
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Graceful fallback behavior (library must work without configs)
//   ❌ Stateless nature (no side effects, thread-safe guarantee)
//   ❌ Package name or import path (hooks depend on "hooks/lib/safety")
//   ❌ Init-time configuration loading (no runtime reload without rearchitecture)
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
// - 3 public APIs orchestrate pattern matching using config loading helpers
// - Ladder: Pure detection logic, no dependencies, used by BLOCKING hooks
// - Baton: Hook → detection function → pattern match → boolean result
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// Quick reference:
// - Adding dangerous pattern category: Edit dangerous-patterns.jsonc (no code changes)
// - Adding critical path category: Edit critical-paths.jsonc (no code changes)
// - Adding secret pattern category: Edit secret-patterns.jsonc (no code changes)
// - Improving pattern matching: See matchesAnyPattern() function in BODY
// - Adding detection function: Follow IsDangerousOperation() pattern
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Types: Config struct memory usage (~50-100KB total)
// - Package-Level State: Init-time config loading (<10ms)
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
// - Config loading: Once per process at init time (<10ms total)
// - Pattern matching: O(n*m) where n=patterns, m=pattern length (<100μs per call)
// - Key optimization: Current performance excellent for hook usage, no optimization needed unless patterns grow to hundreds/thousands
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// See BODY function docstrings above for operation-specific troubleshooting.
//
// Quick reference (details in BODY function docstrings above):
// - Detection not working: Check config files exist, JSONC syntax valid, $HOME set
// - Too many false positives: Pattern too broad, make more specific (conservative bias)
// - Missing dangerous operations: Add pattern to config file, test with detection function
// - Config changes not taking effect: Restart process (config loaded once at init time)
// - Privacy concern: Calling code must NEVER log secret content, only length/metadata
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
// - Key dependencies: None (pure standard library usage)
// - Primary consumers: hooks/tool/pre-use (BLOCKING), hooks/prompt/submit (NON-BLOCKING)
// - Related components: hooks/lib/safety/confirmation.go (uses detection results)
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration-driven pattern detection - COMPLETED (v2.0.0)
//   ✓ Graceful fallback to hardcoded defaults - COMPLETED (v2.0.0)
//   ⏳ Regex pattern support for complex matching
//   ⏳ Platform-aware filtering (Linux vs Windows vs macOS)
//   ⏳ Severity-based filtering (only block critical/high)
//   ⏳ Runtime config reload capability
//
// Research Areas:
//   - Context-aware detection (reduce false positives)
//   - Entropy analysis for secret detection (catch random-looking strings)
//   - Pattern confidence scores (not just boolean match/no-match)
//   - Whitelist support (allow specific known-safe patterns)
//   - Machine learning pattern detection (train on actual dangerous operations)
//
// Integration Targets:
//   - System audit logs (correlation with actual harm)
//   - Detection telemetry (learn from false positives/negatives)
//   - Team-wide pattern sharing (learn from collective experience)
//   - Custom detection plugins (user-defined detection functions)
//
// Known Limitations to Address:
//   - Config loaded once at init (requires process restart for changes)
//   - Simple substring matching (no regex or context awareness)
//   - No confidence scoring (binary true/false only)
//   - No platform filtering (all patterns apply to all platforms)
//   - No severity-based filtering (all matches treated equally)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-11) - Configuration-driven detection with graceful fallbacks
//         - Added JSONC configuration loading for all pattern types
//         - Created dangerous-patterns.jsonc, critical-paths.jsonc, secret-patterns.jsonc
//         - Implemented graceful fallback to hardcoded defaults if config unavailable
//         - Comprehensive 4-block template alignment with full documentation
//         - Design principle: Configuration enables adaptation without code changes
//
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded patterns
//         - Dangerous operation detection (force push, hard reset, rm -rf, sudo, etc.)
//         - Critical path detection (system files, authentication, version control)
//         - Secret detection (API keys, tokens, private keys)
//         - Conservative bias design principle established
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a mid-rung LADDER component providing pure detection logic.
// Used by BLOCKING hooks (tool/pre-use) and NON-BLOCKING secret detection (prompt/submit)
// to protect against accidental data loss, credential exposure, and system damage.
//
// Modify thoughtfully - changes here affect all BLOCKING validation flows. Conservative
// bias must be maintained (false positives disrupt workflow, false negatives are acceptable).
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go test -v ./...)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Maintain conservative bias principle (better to allow than falsely block)
//
// "The prudent see danger and take refuge, but the simple keep going and pay the penalty." - Proverbs 27:12 (WEB)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic dangerous operation check:
//
//     import "hooks/lib/safety"
//
//     if safety.IsDangerousOperation("git push --force origin main") {
//         // Request confirmation from user
//     }
//
// Basic critical file check:
//
//     if safety.IsCriticalFile("/etc/passwd") {
//         // Request confirmation from user
//     }
//
// Basic secret detection (privacy-preserving):
//
//     promptText := getUserPrompt()
//     if safety.ContainsLikelySecret(promptText) {
//         // Warn user (log length only, NEVER content)
//         fmt.Printf("Warning: Prompt may contain secrets (length: %d)\n", len(promptText))
//     }
//
// Using in BLOCKING hook context:
//
//     func preUseHook(toolName, toolInput string) {
//         if toolName == "Bash" {
//             if safety.IsDangerousOperation(toolInput) {
//                 confirmed := requestConfirmation("Dangerous command detected")
//                 if !confirmed {
//                     os.Exit(1) // Block execution
//                 }
//             }
//         }
//         os.Exit(0) // Allow execution
//     }
//
// ============================================================================
// END CLOSING
// ============================================================================
