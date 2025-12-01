// METADATA
//
// Disk Space Monitoring Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "The prudent see danger and take refuge, but the simple keep going and pay the penalty" - Proverbs 27:12 (WEB)
// Principle: Proactive awareness prevents crises - monitor resources before they become critical
// Anchor: "Suppose one of you wants to build a tower. Won't you first sit down and estimate the cost to see if you have enough money to complete it?" - Luke 14:28 (NIV)
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - provides disk space monitoring functionality)
// Role: Monitors workspace disk usage and warns at configurable thresholds
// Paradigm: CPI-SI framework component - serves hooks with resource awareness
//
// Authorship & Lineage
//
// Architect: Nova Dawn
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Template alignment, configuration system, display library integration
//
// Version History:
//   2.0.0 (2025-11-12) - Template alignment, config-driven thresholds, display library integration
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded 80% threshold
//
// Purpose & Function
//
// Purpose: Provides workspace disk space monitoring with configurable warning thresholds
//
// Core Design: Simple disk usage percentage checking using system/lib/system.GetDiskUsage()
// with configurable warning and critical thresholds. Integrates with display library for
// consistent formatting across CPI-SI system.
//
// Key Features:
//   - Configurable warning threshold (default 80%)
//   - Configurable critical threshold (default 95%)
//   - Display library integration for consistent formatting
//   - Optional "healthy" message when disk space is good
//   - Graceful fallback to hardcoded defaults if configuration unavailable
//   - Non-intrusive monitoring (warns only when needed)
//
// Philosophy: Resource awareness prevents crises. Warning at appropriate thresholds gives
// time to clean up before critical situations. Default to sensible behavior even if
// configuration missing.
//
// Blocking Status
//
// Non-blocking: Never blocks session operations. If disk check fails, continues silently
// without disrupting workflow.
// Mitigation: All errors handled gracefully, failed checks just mean no warning displayed
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Import package (configuration loaded automatically in init())
//   2. Call CheckDiskSpace(workspace) to display disk usage warnings
//   3. Function prints to stdout using display library formatting
//   4. No return value - pure side effect (display only)
//
// Public API (in typical usage order):
//
//   Disk Monitoring:
//     CheckDiskSpace(workspace string) - Check and display disk space warnings
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: os (file operations), path/filepath (path handling),
//                     encoding/json (config parsing), fmt (fallback output)
//   Internal: system/lib/system (GetDiskUsage for disk information),
//             system/lib/display (formatted output with health tracking)
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start/start.go (session start disk check)
//   Purpose: Provides workspace resource awareness at session start
//
// Integration Points:
//   - Uses system/lib/system for disk usage information
//   - Uses system/lib/display for formatted output
//   - Reads configuration from system/data/config/session/disk-monitoring.jsonc
//
// Health Scoring
//
// Disk monitoring operations tracked with health scores reflecting monitoring quality.
//
// Configuration Loading:
//   - Config loaded successfully: +10
//   - Config missing/malformed: -5 (falls back to defaults)
//
// Disk Monitoring:
//   - Disk healthy (below thresholds): +15 (good state detected)
//   - Warning threshold exceeded: +10 (successful detection and display)
//   - Critical threshold exceeded: +10 (successful critical detection)
//   - Check failed: -10 (system call error)
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
	"encoding/json" // Parse disk-monitoring.jsonc configuration file
	"fmt"           // Formatted output for warnings and fallback display
	"os"            // File operations for configuration loading and HOME directory
	"path/filepath" // Join paths for configuration file location

	//--- Internal Packages ---
	"system/lib/display" // Formatted output with ANSI colors and health tracking
	"system/lib/system"  // GetDiskUsage for disk space information
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────

// None - all values loaded from configuration (disk-monitoring.jsonc)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────

// ThresholdsConfig defines percentage thresholds for warnings
type ThresholdsConfig struct {
	WarningPercent  float64 `json:"warning_percent"`  // Percentage for warning level
	CriticalPercent float64 `json:"critical_percent"` // Percentage for critical level
	Description     string  `json:"description"`      // What thresholds control
	Reasoning       struct {
		Warning  string `json:"warning"`  // Why this warning threshold
		Critical string `json:"critical"` // Why this critical threshold
	} `json:"reasoning"`
}

// DiskDisplayConfig defines display preferences for disk monitoring
type DiskDisplayConfig struct {
	HeaderIcon       string `json:"header_icon"`        // Icon for disk status header
	HeaderText       string `json:"header_text"`        // Text for disk status header
	ShowWhenHealthy  bool   `json:"show_when_healthy"`  // Show message when disk is healthy
	HealthyMessage   string `json:"healthy_message"`    // Message when disk healthy (with placeholders)
	Description      string `json:"description"`        // What display config controls
}

// DiskMessagesConfig defines customizable warning messages
type DiskMessagesConfig struct {
	Warning     string `json:"warning"`     // Warning level message
	Critical    string `json:"critical"`    // Critical level message
	Description string `json:"description"` // What message config controls
}

// DiskBehaviorConfig defines when and if disk monitoring runs
type DiskBehaviorConfig struct {
	Enabled             bool   `json:"enabled"`               // Master enable/disable switch
	CheckOnSessionStart bool   `json:"check_on_session_start"` // Check at session start
	Description         string `json:"description"`           // What behavior config controls
}

// DiskConfig is the root configuration structure
type DiskConfig struct {
	Metadata struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Created     string `json:"created"`
		LastUpdated string `json:"last_updated"`
		Author      string `json:"author"`
	} `json:"metadata"`
	Thresholds ThresholdsConfig      `json:"thresholds"` // Warning/critical thresholds
	Display    DiskDisplayConfig     `json:"display"`    // Display preferences
	Messages   DiskMessagesConfig    `json:"messages"`   // Warning messages
	Behavior   DiskBehaviorConfig    `json:"behavior"`   // Behavior settings
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────

// Configuration loaded at package initialization
var diskConfig *DiskConfig

// init loads configuration from disk-monitoring.jsonc
// Falls back to hardcoded defaults if configuration unavailable
func init() {
	diskConfig = loadDiskConfig()
}

// loadDiskConfig attempts to load disk-monitoring.jsonc
// Returns config with values, or defaults if loading fails
func loadDiskConfig() *DiskConfig {
	// Try to load from standard config location
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".claude", "cpi-si", "system", "data", "config", "session", "disk-monitoring.jsonc")

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config missing - return hardcoded defaults
		return getDefaultDiskConfig()
	}

	// Strip JSONC comments (shared function from activity.go)
	jsonData := stripJSONCComments(string(data))

	var cfg DiskConfig
	if err := json.Unmarshal([]byte(jsonData), &cfg); err != nil {
		// Config malformed - return hardcoded defaults
		return getDefaultDiskConfig()
	}

	return &cfg
}

// getDefaultDiskConfig returns hardcoded defaults when configuration unavailable
func getDefaultDiskConfig() *DiskConfig {
	cfg := &DiskConfig{}

	// Threshold defaults
	cfg.Thresholds.WarningPercent = 80
	cfg.Thresholds.CriticalPercent = 95

	// Display defaults
	cfg.Display.HeaderIcon = "💾"
	cfg.Display.HeaderText = "Disk Space Status"
	cfg.Display.ShowWhenHealthy = false
	cfg.Display.HealthyMessage = "Disk space: {percent}% used ({available} available)"

	// Message defaults
	cfg.Messages.Warning = "Disk space: {percent}% used ({available} available)"
	cfg.Messages.Critical = "⚠️  CRITICAL: Disk nearly full - {percent}% used (only {available} remaining)"

	// Behavior defaults
	cfg.Behavior.Enabled = true
	cfg.Behavior.CheckOnSessionStart = true

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
//   └── CheckDiskSpace() → uses formatMessage(), system.GetDiskUsage()
//
//   Helpers (Bottom Rungs - Foundations)
//   └── formatMessage() → replaces placeholders in message templates
//
// Baton Flow (Execution Paths):
//
//   Entry → CheckDiskSpace(workspace)
//     ↓
//   Get disk usage via system.GetDiskUsage()
//     ↓
//   Compare usage against thresholds (critical, warning, healthy)
//     ↓
//   Format appropriate message with placeholders replaced
//     ↓
//   Display using display library (or stay silent)
//     ↓
//   Exit → warnings/status displayed to user
//
// APUs (Available Processing Units):
// - 2 functions total
// - 1 helper (formatMessage)
// - 1 public API (CheckDiskSpace)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────

// formatMessage replaces placeholders in message template
//
// What It Does:
// Takes a message template with placeholders ({percent}, {available}, {used}, {total})
// and replaces them with actual disk usage values from system.GetDiskUsage().
//
// Parameters:
//   template: Message string with placeholders
//   diskInfo: Disk usage information from system.GetDiskUsage()
//
// Returns:
//   string: Message with placeholders replaced by actual values
//
// Example usage:
//
//	msg := formatMessage("Disk: {percent}% used", diskInfo)
//	// Returns: "Disk: 85% used"
func formatMessage(template string, diskInfo system.DiskInfo) string {
	// Replace placeholders with actual values
	// Note: Simple string replacement for now
	// Future: Could use text/template for more complex formatting
	msg := template

	// Replace {percent} with usage percentage
	percentStr := fmt.Sprintf("%.0f", diskInfo.UsagePercent)
	msg = replaceAll(msg, "{percent}", percentStr)

	// Replace {available} with available space
	msg = replaceAll(msg, "{available}", diskInfo.Available)

	// Replace {used} with used space
	msg = replaceAll(msg, "{used}", diskInfo.Used)

	// Replace {total} with total space
	msg = replaceAll(msg, "{total}", diskInfo.Total)

	return msg
}

// replaceAll is a simple helper to replace all occurrences of a substring
// Note: Using this instead of strings.ReplaceAll to avoid importing strings
func replaceAll(s, old, new string) string {
	result := ""
	for {
		i := indexOf(s, old)
		if i == -1 {
			result += s
			break
		}
		result += s[:i] + new
		s = s[i+len(old):]
	}
	return result
}

// indexOf finds the first occurrence of substr in s
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────

// CheckDiskSpace examines disk usage and displays warnings at configured thresholds
//
// What It Does:
// Checks workspace disk usage and displays warnings based on configured thresholds.
// Shows critical warning if usage >= critical_percent, warning if >= warning_percent,
// or optional healthy message if below thresholds. Uses configuration for all
// thresholds, messages, and display preferences.
//
// Parameters:
//   workspace: Root directory to check for disk usage
//
// Health Impact:
//   Disk healthy (below thresholds): +15 points (good state)
//   Warning threshold exceeded: +10 points (successful detection)
//   Critical threshold exceeded: +10 points (successful critical detection)
//   Check failed: -10 points (system call error)
//
// Example usage:
//
//	session.CheckDiskSpace("/home/user/project")
//	// Output: 💾 Disk Space Status:
//	//            ⚠️  Disk space: 85% used (150GB available)
func CheckDiskSpace(workspace string) {
	// Check if monitoring enabled
	if !diskConfig.Behavior.Enabled {
		return // Silent when disabled
	}

	// Get disk usage information
	diskInfo := system.GetDiskUsage(workspace)

	// Determine severity level and display appropriate message
	if diskInfo.UsagePercent >= diskConfig.Thresholds.CriticalPercent {
		// Critical level - display critical warning
		headerText := diskConfig.Display.HeaderIcon + " " + diskConfig.Display.HeaderText
		fmt.Println(display.Header(headerText))

		message := formatMessage(diskConfig.Messages.Critical, diskInfo)
		fmt.Printf("   %s\n", display.Failure(message))

	} else if diskInfo.UsagePercent >= diskConfig.Thresholds.WarningPercent {
		// Warning level - display warning
		headerText := diskConfig.Display.HeaderIcon + " " + diskConfig.Display.HeaderText
		fmt.Println(display.Header(headerText))

		message := formatMessage(diskConfig.Messages.Warning, diskInfo)
		fmt.Printf("   %s\n", display.Warning(message))

	} else if diskConfig.Display.ShowWhenHealthy {
		// Healthy level - optionally display success message
		headerText := diskConfig.Display.HeaderIcon + " " + diskConfig.Display.HeaderText
		fmt.Println(display.Header(headerText))

		message := formatMessage(diskConfig.Display.HealthyMessage, diskInfo)
		fmt.Println(display.Success(message))
	}
	// Otherwise: silent (healthy and show_when_healthy is false)
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
//   - Call CheckDiskSpace() with valid workspace path
//   - Test with workspace at various usage levels (below warning, at warning, at critical)
//   - Test with disabled monitoring (config.Behavior.Enabled = false)
//   - Test with show_when_healthy enabled and disabled
//   - Verify message formatting with different templates
//   - Ensure no go vet warnings
//   - Run: go build ./... (library compilation check)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test from start.go hook (session start monitoring)
//   - Verify warnings appear at configured thresholds
//   - Check display library integration (formatted output)
//   - Test configuration fallback (remove config file, verify defaults work)
//
// Example validation code:
//
//     // Test with workspace below warning threshold
//     session.CheckDiskSpace("/home/user/project")
//     // Verify no output (or healthy message if enabled)
//
//     // Test with workspace above warning threshold
//     // (Use test workspace or mock system.GetDiskUsage)
//     session.CheckDiskSpace("/full/workspace")
//     // Verify warning appears
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. No entry point, no main function.
// CheckDiskSpace() waits to be called by session start hook.
//
// Usage: import "hooks/lib/session"
//
// Function executes when called by hook orchestrator. Configuration loaded
// automatically via init() at package import time.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func startSession() {
//         workspace := "/home/user/project"
//         session.CheckDiskSpace(workspace)  // Display disk space warnings
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No state: Component loads config once, uses it throughout session
//   - Disk checks: Uses system.GetDiskUsage() (handles system calls internally)
//   - Memory: Message strings allocated temporarily, garbage collected after display
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle beyond process)
//   - No cleanup function needed (no persistent resources)
//
// Error State Cleanup:
//   - Disk check failures stay silent (no partial state)
//   - Display errors fall back to plain fmt output
//   - No rollback mechanisms needed (pure monitoring, no mutations)
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
// Provides: Workspace disk space monitoring with configurable thresholds
//
// Quick summary:
//   - Monitors disk usage percentage for workspace path
//   - Warns at configurable thresholds (default: 80% warning, 95% critical)
//   - Displays using display library (formatted, color-coded)
//   - Optional healthy message when below thresholds
//   - Non-blocking (informational warnings only)
//
// Integration Pattern: See METADATA "Usage & Integration" section
//
// Public API: CheckDiskSpace(workspace) - monitors and displays disk warnings
//
// Architecture: Ladder component - provides resource awareness validation
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
// For Modification Policy section explanation, see: standards/code/4-block/sections/CWS-SECTION-014-CLOSING-modification-policy.md
//
// Safe to Modify:
//   ✅ Add new threshold levels (info, severe) - extend ThresholdsConfig
//   ✅ Enhance message formatting - improve formatMessage() function
//   ✅ Add placeholder types (filesystem, mount point) - extend formatMessage()
//   ✅ Improve display formatting - use different display library functions
//
// Modify with Extreme Care:
//   ⚠️ CheckDiskSpace() signature - breaks calling hooks
//   ⚠️ Configuration structure - breaks existing config files
//   ⚠️ Threshold defaults - affects users without custom config
//   ⚠️ Display library calls - ensure fallback still works
//
// NEVER Modify:
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking guarantee (never fail/panic on check errors)
//   ❌ system/lib/system integration (standard disk usage API)
//   ❌ Configuration fallback behavior (must degrade gracefully)
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
// For Ladder and Baton Flow section explanation, see: standards/code/4-block/sections/CWS-SECTION-015-CLOSING-ladder-baton-flow.md
//
// See BODY "Organizational Chart" section for complete ladder/baton details.
//
// Quick summary:
// - 1 public API orchestrates disk checking and message formatting
// - Ladder: Public API → formatMessage() helper → system.GetDiskUsage()
// - Baton: Entry → get disk info → compare thresholds → format message → display → exit
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Helpers" section for message formatting function.
//
// Adding new threshold levels:
//   1. Add threshold field to ThresholdsConfig in SETUP
//   2. Update getDefaultConfig() with default value
//   3. Add check in CheckDiskSpace() in proper severity order
//   4. Add message template to MessagesConfig
//   5. Update configuration file schema
//   6. Update API documentation
//
// Example - Adding "info" threshold at 50%:
//
//   // In ThresholdsConfig:
//   InfoPercent float64 `json:"info_percent"`
//
//   // In getDefaultConfig():
//   cfg.Thresholds.InfoPercent = 50
//
//   // In MessagesConfig:
//   Info string `json:"info"`
//
//   // In CheckDiskSpace():
//   } else if diskInfo.UsagePercent >= config.Thresholds.InfoPercent {
//       message := formatMessage(config.Messages.Info, diskInfo)
//       fmt.Println(display.Info(message))
//   }
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// Quick summary:
// - Most expensive operation: system.GetDiskUsage() system call
// - Time complexity: O(1) - single system call per check
// - Memory: O(1) - fixed size message strings
// - Typical runtime: <10ms for disk stat
//
// Bottlenecks:
// - Network-mounted filesystems: Disk stats can be slow (100ms+)
// - Very busy systems: Disk I/O contention delays stats
//
// Optimization notes:
// - Configuration loaded once at init() (not per call)
// - Message formatting is simple string replacement (fast)
// - No caching (disk usage changes rapidly, caching counterproductive)
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: Warnings not displayed even when disk is full
//   - Check: config.Behavior.Enabled is true
//   - Check: Workspace path exists and is valid
//   - Check: Threshold percentages are reasonable (80, 95)
//   - Solution: Verify configuration loaded correctly, check file permissions
//
// Problem: Wrong threshold triggering
//   - Cause: Configuration thresholds set incorrectly
//   - Solution: Check disk-monitoring.jsonc values, ensure warning < critical
//
// Problem: Display formatting broken (no colors/icons)
//   - Check: Terminal supports ANSI colors
//   - Check: display library available
//   - Solution: Verify display library compiled correctly, check terminal settings
//
// Problem: Configuration not loading (using fallback defaults)
//   - Check: File exists at correct path
//   - Check: File permissions allow reading
//   - Check: JSONC syntax is valid
//   - Solution: Verify file location, validate JSON structure
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section for complete dependency information.
//
// Quick summary:
// - Key dependencies: system/lib/system (disk usage), system/lib/display (formatting)
// - Primary consumer: Session start hook (workspace monitoring)
// - Related components: Other session libraries (dependencies.go, git.go) also validate workspace
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Warning threshold - COMPLETED (v2.0.0)
//   ✓ Critical threshold - COMPLETED (v2.0.0)
//   ✓ Configuration system - COMPLETED (v2.0.0)
//   ⏳ Info threshold (50% for awareness)
//   ⏳ Absolute minimum space (warn if <10GB regardless of %)
//   ⏳ Unit preference control (GB, TB, auto)
//   ⏳ Multiple path monitoring (check several filesystems)
//
// Research Areas:
//   - Trend analysis (disk usage increasing over time)
//   - Predictive warnings (project when disk will be full)
//   - Cleanup suggestions (find large files, old caches)
//   - Integration with system cleanup tools
//
// Known Limitations:
//   - Single workspace path check (no multi-path support)
//   - No historical tracking (only current usage)
//   - No size-based warnings (only percentage)
//   - No cleanup automation (only warnings)
//
// Version History:
//
//   2.0.0 (2025-11-12) - Configuration system and template alignment
//         - Full 4-block template alignment
//         - Display library integration for formatted warnings
//         - Configuration system with warning/critical thresholds
//         - Customizable messages and display preferences
//         - Comprehensive inline documentation
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Basic 80% threshold check
//         - Plain fmt.Printf output
//         - Minimal documentation
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library provides early warning of disk space issues - catching low space
// at session start rather than during critical work. Proactive monitoring prevents
// crises and lost work.
//
// *"The prudent see danger and take refuge" - Proverbs 27:12*
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
//       session.CheckDiskSpace("/home/user/project")
//       // Output: 💾 Disk Space Status:
//       //            ⚠️  Disk space: 85% used (150GB available)
//   }
//
// At Warning Threshold (80%):
//   session.CheckDiskSpace(workspace)
//   // Shows warning with configured message
//
// At Critical Threshold (95%):
//   session.CheckDiskSpace(workspace)
//   // Shows critical warning with different formatting
//
// Below Thresholds (healthy):
//   session.CheckDiskSpace(workspace)
//   // Silent by default, or shows success message if show_when_healthy enabled
//
// With Custom Configuration:
//   // Edit ~/.claude/cpi-si/system/data/config/session/disk-monitoring.jsonc
//   // Change warning_percent to 70
//   session.CheckDiskSpace(workspace)
//   // Now warns at 70% instead of 80%
//
// Disabled Monitoring:
//   // Set "enabled": false in config
//   session.CheckDiskSpace(workspace)
//   // No output (monitoring disabled)
//
// ════════════════════════════════════════════════════════════════
// END CLOSING
// ════════════════════════════════════════════════════════════════
