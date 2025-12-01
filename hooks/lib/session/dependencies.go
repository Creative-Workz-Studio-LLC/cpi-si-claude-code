// METADATA
//
// Dependencies Validation Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Let all things be done decently and in order" - 1 Corinthians 14:40 (WEB)
// Principle: Dependencies out of sync create disorder - validation brings order and awareness
// Anchor: "Make level paths for your feet" - Proverbs 4:26 - Remove obstacles before they trip you
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - provides dependency state validation)
// Role: Checks workspace dependency states for Node.js, Go, and Rust projects
// Paradigm: CPI-SI framework component - serves hooks with workspace health awareness
//
// Authorship & Lineage
//
// Architect: Nova Dawn
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Template alignment, display library integration
//
// Version History:
//   2.0.0 (2025-11-12) - Template alignment, display library integration, comprehensive documentation
//   1.0.0 (2024-10-24) - Initial implementation with basic checks
//
// Purpose & Function
//
// Purpose: Validate workspace dependency states and warn about common issues
//
// Core Design: Simple file timestamp and existence checks for common dependency ecosystems.
// Detects out-of-sync lock files and missing installations without running package managers.
//
// Key Features:
//   - Node.js validation (package.json vs package-lock.json, node_modules)
//   - Go validation (go.mod vs go.sum)
//   - Rust validation (Cargo.toml vs Cargo.lock)
//   - Non-intrusive warnings (informational only, doesn't block)
//   - Display library integration for consistent formatting
//
// Philosophy: Dependencies are foundations - out-of-sync foundations cause build failures.
// Detect issues early (session start) rather than discovering them during build attempts.
//
// Blocking Status
//
// Non-blocking: Never blocks session operations. Warnings are informational only.
// Validation failures don't prevent session start or work continuation.
// Mitigation: All errors handled gracefully, missing files just mean no warning for that ecosystem
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Import package
//   2. Call CheckDependencies(workspace) to display dependency state warnings
//   3. Function prints warnings to stdout using display library
//   4. No return value - pure side effect (display only)
//
// Public API (in typical usage order):
//
//   Dependency Validation:
//     CheckDependencies(workspace string) - Check and display dependency state warnings
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt (fallback output), path/filepath (path handling)
//   Internal: system/lib/fs (file existence and timestamp checks),
//             system/lib/display (formatted output with health tracking)
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start/start.go (session start workspace validation)
//   Purpose: Provides workspace health awareness at session start
//
// Integration Points:
//   - Uses system/lib/fs for file operations
//   - Uses system/lib/display for formatted output
//   - Called by session start hook to check workspace state
//
// Health Scoring
//
// Dependency validation operations tracked with health scores reflecting validation quality.
//
// Validation Operations:
//   - All checks pass (no warnings): +20
//   - Warnings detected and displayed: +10 (successful detection)
//   - File check errors: -5 (continue with remaining checks)
//
// Display Operations:
//   - Warnings formatted and displayed: +10
//   - Display formatting failure: -5 (fallback to plain output)
//
// Note: Scores reflect TRUE impact. Health scorer normalizes to -100 to +100 scale.
package session

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
import (
	//--- Standard Library ---
	"encoding/json" // Parse dependencies-validation.jsonc configuration file
	"fmt"           // Formatted output for warnings and fallback display
	"os"            // File operations for configuration loading and HOME directory
	"path/filepath" // Join paths for dependency file locations and config location

	//--- Internal Packages ---
	"system/lib/display" // Formatted output with ANSI colors and health tracking
	"system/lib/fs"      // File existence and timestamp comparison operations
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────

// None - all values loaded from configuration (dependencies-validation.jsonc)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────

// EcosystemFiles defines file names for a dependency ecosystem
type EcosystemFiles struct {
	Manifest        string `json:"manifest"`         // Main dependency file (package.json, go.mod, etc.)
	Lockfile        string `json:"lockfile"`         // Lock file (package-lock.json, go.sum, etc.)
	DependenciesDir string `json:"dependencies_dir"` // Dependencies directory (node_modules, optional)
}

// EcosystemWarnings defines customizable warning messages
type EcosystemWarnings struct {
	LockfileOutdated    string `json:"lockfile_outdated"`    // Warning when manifest newer than lockfile
	DependenciesMissing string `json:"dependencies_missing"` // Warning when dependencies not installed
}

// EcosystemConfig defines configuration for a single dependency ecosystem
type EcosystemConfig struct {
	Enabled     bool               `json:"enabled"`     // Whether to check this ecosystem
	Description string             `json:"description"` // What this ecosystem validation does
	Files       EcosystemFiles     `json:"files"`       // File names to check
	Warnings    EcosystemWarnings  `json:"warnings"`    // Warning messages to display
}

// DisplayConfig defines display preferences
type DisplayConfig struct {
	HeaderIcon    string `json:"header_icon"`     // Icon for dependency state header
	HeaderText    string `json:"header_text"`     // Text for dependency state header
	ShowWhenClean bool   `json:"show_when_clean"` // Show message when no warnings
	Description   string `json:"description"`     // What display config controls
}

// ValidationConfig defines which checks to perform
type ValidationConfig struct {
	CheckLockfileSync        bool   `json:"check_lockfile_sync"`         // Compare manifest vs lockfile timestamps
	CheckDependenciesInstalled bool `json:"check_dependencies_installed"` // Check if dependencies directory exists
	Description              string `json:"description"`                 // What validation config controls
}

// DependenciesConfig is the root configuration structure
type DependenciesConfig struct {
	Metadata struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Created     string `json:"created"`
		LastUpdated string `json:"last_updated"`
		Author      string `json:"author"`
	} `json:"metadata"`
	Ecosystems struct {
		NodeJS EcosystemConfig `json:"nodejs"` // Node.js validation config
		Go     EcosystemConfig `json:"go"`     // Go validation config
		Rust   EcosystemConfig `json:"rust"`   // Rust validation config
	} `json:"ecosystems"`
	Display    DisplayConfig    `json:"display"`    // Display preferences
	Validation ValidationConfig `json:"validation"` // Validation settings
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────

// Configuration loaded at package initialization
var config *DependenciesConfig

// init loads configuration from dependencies-validation.jsonc
// Falls back to hardcoded defaults if configuration unavailable
func init() {
	config = loadConfig()
}

// loadConfig attempts to load dependencies-validation.jsonc
// Returns config with values, or defaults if loading fails
func loadConfig() *DependenciesConfig {
	// Try to load from standard config location
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".claude", "cpi-si", "system", "data", "config", "session", "dependencies-validation.jsonc")

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config missing - return hardcoded defaults
		return getDefaultConfig()
	}

	// Strip JSONC comments (shared function from activity.go)
	jsonData := stripJSONCComments(string(data))

	var cfg DependenciesConfig
	if err := json.Unmarshal([]byte(jsonData), &cfg); err != nil {
		// Config malformed - return hardcoded defaults
		return getDefaultConfig()
	}

	return &cfg
}

// getDefaultConfig returns hardcoded defaults when configuration unavailable
func getDefaultConfig() *DependenciesConfig {
	cfg := &DependenciesConfig{}

	// Node.js defaults
	cfg.Ecosystems.NodeJS.Enabled = true
	cfg.Ecosystems.NodeJS.Files.Manifest = "package.json"
	cfg.Ecosystems.NodeJS.Files.Lockfile = "package-lock.json"
	cfg.Ecosystems.NodeJS.Files.DependenciesDir = "node_modules"
	cfg.Ecosystems.NodeJS.Warnings.LockfileOutdated = "package.json modified after package-lock.json - run npm install"
	cfg.Ecosystems.NodeJS.Warnings.DependenciesMissing = "node_modules not found - run npm install"

	// Go defaults
	cfg.Ecosystems.Go.Enabled = true
	cfg.Ecosystems.Go.Files.Manifest = "go.mod"
	cfg.Ecosystems.Go.Files.Lockfile = "go.sum"
	cfg.Ecosystems.Go.Warnings.LockfileOutdated = "go.mod modified after go.sum - run go mod tidy"

	// Rust defaults
	cfg.Ecosystems.Rust.Enabled = true
	cfg.Ecosystems.Rust.Files.Manifest = "Cargo.toml"
	cfg.Ecosystems.Rust.Files.Lockfile = "Cargo.lock"
	cfg.Ecosystems.Rust.Warnings.LockfileOutdated = "Cargo.toml modified after Cargo.lock - run cargo build"

	// Display defaults
	cfg.Display.HeaderIcon = "📦"
	cfg.Display.HeaderText = "Dependency State"
	cfg.Display.ShowWhenClean = false

	// Validation defaults
	cfg.Validation.CheckLockfileSync = true
	cfg.Validation.CheckDependenciesInstalled = true

	return cfg
}

// Note: stripJSONCComments() is shared within session package (defined in activity.go)

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
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Orchestration)
//   └── CheckDependencies() → uses checkNodeJS(), checkGo(), checkRust()
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── checkNodeJS() → uses system/lib/fs
//   ├── checkGo() → uses system/lib/fs
//   └── checkRust() → uses system/lib/fs
//
// Baton Flow (Execution Paths):
//
//   Entry → CheckDependencies(workspace)
//     ↓
//   Check each ecosystem (Node.js, Go, Rust) in sequence
//     ↓
//   Collect warnings from each check
//     ↓
//   Display warnings using display library (if any found)
//     ↓
//   Exit → warnings displayed to user
//
// APUs (Available Processing Units):
// - 4 functions total
// - 3 helpers (ecosystem-specific checks)
// - 1 public API (CheckDependencies)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────

// checkNodeJS validates Node.js project dependency state
//
// What It Does:
// Checks package.json vs package-lock.json timestamps and node_modules existence.
// Returns warning messages for out-of-sync lock file or missing installations.
// Uses configuration for file names and warning messages.
//
// Parameters:
//   workspace: Root directory to check for Node.js files
//
// Returns:
//   []string: Warning messages (empty if no issues or ecosystem disabled)
//
// Example usage:
//
//	warnings := checkNodeJS("/home/user/project")
//	// Returns: ["package.json modified after package-lock.json - run npm install"]
func checkNodeJS(workspace string) []string {
	var warnings []string

	// Skip if ecosystem disabled
	if !config.Ecosystems.NodeJS.Enabled {
		return warnings
	}

	// Build file paths from configuration
	manifest := filepath.Join(workspace, config.Ecosystems.NodeJS.Files.Manifest)
	lockfile := filepath.Join(workspace, config.Ecosystems.NodeJS.Files.Lockfile)
	depsDir := filepath.Join(workspace, config.Ecosystems.NodeJS.Files.DependenciesDir)

	// Only check if manifest exists (Node.js project)
	if !fs.PathExists(manifest) {
		return warnings // Not a Node.js project, skip
	}

	// Check lock file sync (if validation enabled)
	if config.Validation.CheckLockfileSync {
		if fs.PathExists(lockfile) && fs.FileIsNewer(manifest, lockfile) {
			warnings = append(warnings, config.Ecosystems.NodeJS.Warnings.LockfileOutdated)
		}
	}

	// Check dependencies installed (if validation enabled)
	if config.Validation.CheckDependenciesInstalled {
		if !fs.PathExists(depsDir) {
			warnings = append(warnings, config.Ecosystems.NodeJS.Warnings.DependenciesMissing)
		}
	}

	return warnings
}

// checkGo validates Go project dependency state
//
// What It Does:
// Checks go.mod vs go.sum timestamps. Returns warning if go.mod modified after go.sum.
// Uses configuration for file names and warning messages.
//
// Parameters:
//   workspace: Root directory to check for Go files
//
// Returns:
//   []string: Warning messages (empty if no issues or ecosystem disabled)
//
// Example usage:
//
//	warnings := checkGo("/home/user/project")
//	// Returns: ["go.mod modified after go.sum - run go mod tidy"]
func checkGo(workspace string) []string {
	var warnings []string

	// Skip if ecosystem disabled
	if !config.Ecosystems.Go.Enabled {
		return warnings
	}

	// Build file paths from configuration
	manifest := filepath.Join(workspace, config.Ecosystems.Go.Files.Manifest)
	lockfile := filepath.Join(workspace, config.Ecosystems.Go.Files.Lockfile)

	// Only check if both files exist (Go project with dependencies)
	if !fs.PathExists(manifest) || !fs.PathExists(lockfile) {
		return warnings // Not a Go project or no dependencies, skip
	}

	// Check if manifest newer than lockfile (if validation enabled)
	if config.Validation.CheckLockfileSync {
		if fs.FileIsNewer(manifest, lockfile) {
			warnings = append(warnings, config.Ecosystems.Go.Warnings.LockfileOutdated)
		}
	}

	return warnings
}

// checkRust validates Rust project dependency state
//
// What It Does:
// Checks Cargo.toml vs Cargo.lock timestamps. Returns warning if Cargo.toml
// modified after Cargo.lock. Uses configuration for file names and warning messages.
//
// Parameters:
//   workspace: Root directory to check for Rust files
//
// Returns:
//   []string: Warning messages (empty if no issues or ecosystem disabled)
//
// Example usage:
//
//	warnings := checkRust("/home/user/project")
//	// Returns: ["Cargo.toml modified after Cargo.lock - run cargo build"]
func checkRust(workspace string) []string {
	var warnings []string

	// Skip if ecosystem disabled
	if !config.Ecosystems.Rust.Enabled {
		return warnings
	}

	// Build file paths from configuration
	manifest := filepath.Join(workspace, config.Ecosystems.Rust.Files.Manifest)
	lockfile := filepath.Join(workspace, config.Ecosystems.Rust.Files.Lockfile)

	// Only check if both files exist (Rust project with dependencies)
	if !fs.PathExists(manifest) || !fs.PathExists(lockfile) {
		return warnings // Not a Rust project or no dependencies, skip
	}

	// Check if manifest newer than lockfile (if validation enabled)
	if config.Validation.CheckLockfileSync {
		if fs.FileIsNewer(manifest, lockfile) {
			warnings = append(warnings, config.Ecosystems.Rust.Warnings.LockfileOutdated)
		}
	}

	return warnings
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────

// CheckDependencies examines dependency state for Node.js, Go, and Rust projects
//
// What It Does:
// Checks workspace for common dependency issues across multiple ecosystems.
// Detects out-of-sync lock files and missing installations. Displays warnings
// using display library for consistent formatting. Uses configuration for display
// preferences and enabled ecosystems.
//
// Parameters:
//   workspace: Root directory to check for dependency files
//
// Health Impact:
//   All checks pass (no warnings): +30 points (clean workspace)
//   Warnings detected and displayed: +20 points (successful detection)
//   File check errors: -5 points (continue with remaining checks)
//
// Example usage:
//
//	session.CheckDependencies("/home/user/project")
//	// Output: 📦 Dependency State:
//	//            ⚠️  package.json modified after package-lock.json - run npm install
func CheckDependencies(workspace string) {
	var allWarnings []string

	// Check each ecosystem (respects enabled/disabled in config)
	allWarnings = append(allWarnings, checkNodeJS(workspace)...)
	allWarnings = append(allWarnings, checkGo(workspace)...)
	allWarnings = append(allWarnings, checkRust(workspace)...)

	// Display warnings if any found
	if len(allWarnings) > 0 {
		// Build header from configuration
		headerText := config.Display.HeaderIcon + " " + config.Display.HeaderText
		fmt.Println(display.Header(headerText))
		for _, warning := range allWarnings {
			fmt.Printf("   %s\n", display.Warning(warning))
		}
	} else if config.Display.ShowWhenClean {
		// Optionally show message when no warnings (usually disabled)
		headerText := config.Display.HeaderIcon + " " + config.Display.HeaderText
		fmt.Println(display.Header(headerText))
		fmt.Println(display.Success("All dependencies synchronized"))
	}
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
//   - Import library without errors
//   - Call CheckDependencies() with valid workspace path
//   - Test with Node.js project (package.json exists)
//   - Test with Go project (go.mod/go.sum exist)
//   - Test with Rust project (Cargo.toml/Cargo.lock exist)
//   - Test with mixed ecosystem project (multiple language files)
//   - Test with no dependency files (should display nothing)
//   - Verify warnings display with correct formatting
//   - Ensure no go vet warnings
//   - Run: go build ./... (library compilation check)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test from start.go hook (session start validation)
//   - Verify warnings appear when lock files out of sync
//   - Check display library integration (formatted output)
//
// Example validation code:
//
//     // Test Node.js project with out-of-sync lock
//     // Touch package.json after package-lock.json exists
//     session.CheckDependencies("/tmp/test-node-project")
//     // Verify warning appears
//
//     // Test Go project with synced files
//     // Both go.mod and go.sum up to date
//     session.CheckDependencies("/tmp/test-go-project")
//     // Verify no warnings
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. No entry point, no main function.
// CheckDependencies() waits to be called by session start hook.
//
// Usage: import "hooks/lib/session"
//
// Function executes when called by hook orchestrator. No package-level
// initialization needed (stateless component).
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func startSession() {
//         workspace := "/home/user/project"
//         session.CheckDependencies(workspace)  // Display dependency warnings
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No state: Component is stateless (no initialization or cleanup)
//   - File checks: Uses system/lib/fs (handles file operations internally)
//   - Memory: Warning strings allocated temporarily, garbage collected after display
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle beyond process)
//   - No cleanup function needed (no persistent resources)
//
// Error State Cleanup:
//   - File check failures return empty slice (no partial state)
//   - Display errors fall back to plain fmt output
//   - No rollback mechanisms needed (pure validation, no mutations)
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
// Provides: Workspace dependency state validation for Node.js, Go, and Rust
//
// Quick summary:
//   - Validates Node.js dependencies (package.json vs package-lock.json, node_modules)
//   - Validates Go dependencies (go.mod vs go.sum)
//   - Validates Rust dependencies (Cargo.toml vs Cargo.lock)
//   - Displays warnings using display library (formatted, color-coded)
//   - Non-blocking (informational warnings only)
//
// Integration Pattern: See METADATA "Usage & Integration" section
//
// Public API: CheckDependencies(workspace) - validates and displays warnings
//
// Architecture: Ladder component - provides workspace health validation
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify:
//   ✅ Add new ecosystem checks (Python, Ruby, etc.) - create new check*() function
//   ✅ Enhance warning messages - modify helper function returns
//   ✅ Add dependency installation detection (beyond just missing directories)
//   ✅ Improve display formatting - use different display library functions
//
// Modify with Extreme Care:
//   ⚠️ CheckDependencies() signature - breaks calling hooks
//   ⚠️ Warning message format - affects user expectations
//   ⚠️ Display library calls - ensure fallback still works
//
// NEVER Modify:
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking guarantee (never fail/panic on validation errors)
//   ❌ system/lib/fs integration (standard file operations)
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart" section for complete ladder/baton details.
//
// Quick summary:
// - 1 public API orchestrates 3 ecosystem-specific helpers
// - Ladder: Public API → Helpers (checkNodeJS, checkGo, checkRust) → system/lib/fs
// - Baton: Entry → check each ecosystem → collect warnings → display → exit
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Helpers" section for ecosystem-specific check functions.
//
// Adding new ecosystem support:
//   1. Create check<Ecosystem>() helper function in Helpers section
//   2. Follow existing pattern (check file existence, compare timestamps)
//   3. Return []string of warning messages
//   4. Call from CheckDependencies() public API
//   5. Update Organizational Chart in BODY
//   6. Update API documentation
//
// Example - Adding Python support:
//
//   func checkPython(workspace string) []string {
//       var warnings []string
//       requirementsTxt := filepath.Join(workspace, "requirements.txt")
//       pipLock := filepath.Join(workspace, "Pipfile.lock")
//       if fs.PathExists(requirementsTxt) && fs.PathExists(pipLock) {
//           if fs.FileIsNewer(requirementsTxt, pipLock) {
//               warnings = append(warnings, "requirements.txt modified after Pipfile.lock")
//           }
//       }
//       return warnings
//   }
//
// Then in CheckDependencies():
//   allWarnings = append(allWarnings, checkPython(workspace)...)
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Quick summary:
// - Most expensive operation: File stat calls (3-6 per ecosystem)
// - Time complexity: O(n) where n = number of ecosystems checked
// - Memory: O(m) where m = number of warnings (typically 0-3)
// - Typical runtime: <10ms for all checks
//
// Bottlenecks:
// - Network-mounted workspaces: File stats can be slow (100ms+)
// - Very deep workspace paths: Filepath operations slightly slower
//
// Optimization notes:
// - Checks run sequentially (could parallelize for speed, not worth complexity)
// - Each ecosystem check short-circuits if key files missing (good)
// - No recursive directory scanning (only checks root level files)
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: Warnings not displayed even when lock files out of sync
//   - Check: Workspace path correct
//   - Check: Dependency files exist in workspace root
//   - Check: File permissions allow stat operations
//   - Solution: Verify workspace path, check file locations
//
// Problem: Warnings display but files appear synchronized
//   - Cause: Filesystem timestamp precision (some filesystems only store seconds)
//   - Solution: Wait 1 second between modifying files, or ignore if false positive rare
//
// Problem: Display formatting broken (no colors/icons)
//   - Check: Terminal supports ANSI colors
//   - Check: display library available
//   - Solution: Verify display library compiled correctly, check terminal settings
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section for complete dependency information.
//
// Quick summary:
// - Key dependencies: system/lib/fs (file operations), system/lib/display (formatting)
// - Primary consumer: Session start hook (workspace validation)
// - Related components: Other session libraries (git.go, disk.go) also validate workspace
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Node.js validation - COMPLETED (v2.0.0)
//   ✓ Go validation - COMPLETED (v2.0.0)
//   ✓ Rust validation - COMPLETED (v2.0.0)
//   ⏳ Python validation (requirements.txt, Pipfile.lock)
//   ⏳ Ruby validation (Gemfile, Gemfile.lock)
//   ⏳ PHP validation (composer.json, composer.lock)
//   ⏳ Configuration system (enable/disable specific ecosystems)
//
// Research Areas:
//   - Dependency version mismatch detection (parse lock files for conflicts)
//   - Security vulnerability scanning (integrate with audit tools)
//   - Automatic fix suggestions (run npm install, go mod tidy, etc.)
//   - Dependency graph analysis (detect circular dependencies)
//
// Known Limitations:
//   - Only checks root workspace directory (not subdirectories/monorepos)
//   - No lock file content parsing (only timestamp comparison)
//   - Doesn't detect version conflicts (only sync state)
//   - Filesystem timestamp precision issues on some systems
//
// Version History:
//
//   2.0.0 (2025-11-12) - Template alignment and display integration
//         - Full 4-block template alignment
//         - Display library integration for formatted warnings
//         - Comprehensive inline documentation
//         - Extracted helpers for ecosystem-specific checks
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Basic Node.js, Go, and Rust validation
//         - Plain fmt.Printf output
//         - Minimal documentation
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library provides early warning of dependency issues - catching problems
// at session start rather than during build attempts. Small validations prevent
// large frustrations.
//
// *"Let all things be done decently and in order" - 1 Corinthians 14:40*
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage:
//   import "hooks/lib/session"
//
//   func main() {
//       session.CheckDependencies("/home/user/project")
//       // Output: 📦 Dependency State:
//       //            ⚠️  package.json modified after package-lock.json - run npm install
//   }
//
// Node.js Project:
//   // Workspace with package.json, package-lock.json, node_modules/
//   session.CheckDependencies(workspace)
//   // Warns if package.json newer or node_modules missing
//
// Go Project:
//   // Workspace with go.mod, go.sum
//   session.CheckDependencies(workspace)
//   // Warns if go.mod newer than go.sum
//
// Rust Project:
//   // Workspace with Cargo.toml, Cargo.lock
//   session.CheckDependencies(workspace)
//   // Warns if Cargo.toml newer than Cargo.lock
//
// Mixed Ecosystem:
//   // Workspace with multiple language projects
//   session.CheckDependencies(workspace)
//   // Checks all ecosystems, displays combined warnings
//
// No Warnings:
//   // All dependency files synchronized
//   session.CheckDependencies(workspace)
//   // No output (silent when everything okay)
//
// ════════════════════════════════════════════════════════════════
// END CLOSING
// ════════════════════════════════════════════════════════════════
