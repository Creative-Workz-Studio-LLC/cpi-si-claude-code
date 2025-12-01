// METADATA
//
// Process Monitoring Library - CPI-SI Hooks System
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "Watch and pray, that ye enter not into temptation" - Matthew 26:41 (KJV)
// Principle: Watchful awareness prevents problems - monitoring active processes serves developers
// Anchor: "Be sober, be vigilant" - 1 Peter 5:8 (KJV)
//
// CPI-SI Identity
//
// Component Type: LIBRARY - Session awareness utility (session-specific rung)
// Role: Monitors development server ports and reports active processes at session start/end
// Paradigm: CPI-SI framework component - awareness without interference
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI instance)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Added configuration system for flexible port monitoring
//
// Version History:
//   2.0.0 (2025-11-12) - Configuration-driven port list, display customization, timeout protection
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded port list
//
// Purpose & Function
//
// Purpose: Provide configurable, non-blocking detection of active development servers on common
// ports at session start and end. Helps developers track what's running without manual checking,
// prevents forgotten background processes, and provides awareness of development environment state.
//
// Core Design: Simple port checking pattern - uses lsof to detect listening processes on configured
// ports, reports findings with customizable formatting. Non-blocking by design - detection failures
// don't interrupt session flow.
//
// Key Features:
//   - Configurable port list (customize which ports to monitor)
//   - Non-blocking execution (detection failures don't interrupt sessions)
//   - Timeout protection (lsof commands can't hang session)
//   - Display customization (icons, messages, formatting)
//   - Graceful fallback (works with hardcoded defaults if config missing)
//   - Dual-context support (session start vs session end messaging)
//
// Philosophy: "Awareness serves work" - monitoring active processes provides helpful context
// without becoming intrusive. If detection fails, work continues - awareness serves developers,
// not blocks them.
//
// Blocking Status
//
// Non-blocking: Detection failures are silent - session continues regardless of lsof status.
// Port checking may fail due to missing lsof, permission issues, or system errors. Library
// returns immediately on error without interrupting session flow.
//
// Mitigation: Configuration allows logging failures for debugging while maintaining non-blocking
// behavior. display.show_* flags control output visibility.
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
//	func main() {
//		session.CheckRunningProcesses()           // At session start
//		// ... session work ...
//		session.CheckRunningProcessesAsReminder() // At session end
//	}
//
// Integration Pattern:
//   1. Import session library
//   2. Call CheckRunningProcesses() at session start for awareness
//   3. Call CheckRunningProcessesAsReminder() at session end as reminder
//   4. No cleanup needed - detection is stateless
//
// Public API (in typical usage order):
//
//   Detection Functions:
//     CheckRunningProcesses() - Detect and report active dev servers (session start context)
//     CheckRunningProcessesAsReminder() - Detect and report active dev servers (session end context)
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: fmt, os/exec, strings, context, time
//   External: None
//   Internal: None (self-contained detection logic)
//   System Commands: lsof (for port checking)
//   Config Files: ~/.claude/cpi-si/system/data/config/session/processes.jsonc
//
// Dependents (What Uses This):
//   Commands: session/cmd-start/start.go, session/cmd-end/end.go
//   Libraries: None (leaf node in dependency tree)
//
// Integration Points:
//   - Rails: No logger needed (silent execution or direct output)
//   - Baton: Stateless detection (no data flow between calls)
//   - Ladder: System command execution (lsof)
//
// Health Scoring
//
// Process monitoring tracked through detection execution and configuration success:
//
// Configuration Loading:
//   - Config loaded successfully: +20
//   - Config missing (using defaults): +15
//   - Config invalid (using defaults): +10
//
// Port Detection:
//   - All ports checked successfully: +40
//   - Partial port checks succeed: +20
//   - All port checks fail: -10
//
// Display Output:
//   - Output formatted correctly: +20
//   - Output has minor issues: +10
//   - Output fails: +0
//
// Note: Scores reflect TRUE impact. Non-blocking design means failures don't prevent work,
// but successful detection enables better developer awareness and environment management.
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
// Dependencies this component needs. Organized by source - standard library
// provides Go's built-in capabilities. Each import commented with purpose,
// not just name. This library has no internal package dependencies - operates
// as self-contained port detection utility.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"context"       // Context for timeout control on port checks
	"encoding/json" // JSON parsing for configuration file
	"fmt"           // Formatted output for process reporting
	"os"            // File operations and environment access (UserHomeDir)
	"os/exec"       // Command execution for calling lsof
	"path/filepath" // Path construction for configuration file
	"strings"       // String manipulation for output formatting
	"time"          // Duration for timeout specification
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Named values that never change. Magic numbers given meaningful names,
// configuration values documented with reasoning. Constants prevent bugs
// from typos and make intent clear - default ports have purpose (common
// dev servers), not just arbitrary numbers. These provide graceful
// fallback when configuration file is missing or invalid.
//
// See: standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

const (
	// Default port list (fallback if config file missing)
	// Common development server ports across multiple frameworks
	defaultPort3000 = "3000" // React/Next.js default
	defaultPort8000 = "8000" // Django/Python HTTP server
	defaultPort8080 = "8080" // Generic HTTP server (Spring Boot, Tomcat)
	defaultPort5173 = "5173" // Vite dev server
	defaultPort4200 = "4200" // Angular CLI dev server

	// Default timeout for port checks (seconds)
	defaultProcessTimeout = 2

	// Default display settings
	defaultIcon          = "🔌"
	defaultSeparator     = ", "
	defaultStartMessage  = "Active dev servers on ports:"
	defaultEndMessage    = "Dev servers still running on ports:"
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Types organized bottom-up showing dependency relationships. Simple
// building blocks first (PortConfig), then composed structures
// (PortsConfig, DisplayConfig, BehaviorConfig), finally top-level
// aggregation (ProcessesConfig). This organization reveals what
// depends on what.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

// PortConfig defines configuration for a single monitored port
type PortConfig struct {
	Number      string `json:"number"`      // Port number as string
	Description string `json:"description"` // What typically runs on this port
	Enabled     bool   `json:"enabled"`     // Whether to check this port
}

// PortsConfig defines which ports to monitor
type PortsConfig struct {
	Enabled        bool         `json:"enabled"`         // Master switch for port monitoring
	MonitoredPorts []PortConfig `json:"monitored_ports"` // Configured ports with descriptions
	CustomPorts    []string     `json:"custom_ports"`    // Additional port numbers to check
}

// ProcessDisplayConfig defines output formatting for process detection
type ProcessDisplayConfig struct {
	ShowAtStart      bool   `json:"show_at_start"`      // Show at session start
	ShowAtEnd        bool   `json:"show_at_end"`        // Show at session end
	Icon             string `json:"icon"`               // Icon for notification
	StartMessage     string `json:"start_message"`      // Message prefix at start
	EndMessage       string `json:"end_message"`        // Message prefix at end
	Separator        string `json:"separator"`          // Separator between ports
	ShowDescriptions bool   `json:"show_descriptions"`  // Include descriptions (verbose)
}

// ProcessBehaviorConfig defines execution behavior preferences
type ProcessBehaviorConfig struct {
	SilentFailures bool   `json:"silent_failures"` // Continue silently on error
	TimeoutSeconds int    `json:"timeout_seconds"` // Max wait time per port check
	CheckCommand   string `json:"check_command"`   // Command to use (lsof or ss)
	RequireLsof    bool   `json:"require_lsof"`    // Fail if lsof not available
}

// ProcessesConfig is the top-level configuration structure for process monitoring
type ProcessesConfig struct {
	Ports    PortsConfig              `json:"ports"`    // Port configuration
	Display  ProcessDisplayConfig     `json:"display"`  // Display preferences
	Behavior ProcessBehaviorConfig    `json:"behavior"` // Behavior preferences
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// Infrastructure available throughout component. Rails pattern - configuration
// loaded once at package initialization, available to all functions without
// parameter passing. State lives here, functions use it.
//
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md

var (
	processConfig       *ProcessesConfig // Cached configuration loaded in init()
	processConfigLoaded bool             // Flag indicating if config loaded successfully
)

func init() {
	// --- Configuration Loading ---
	// Load process monitoring configuration at package import
	// Falls back to hardcoded defaults if config file missing or invalid

	homeDir, err := os.UserHomeDir() // Get user home directory
	if err != nil {
		processConfigLoaded = false // Can't load without home dir
		return
	}

	// Build path to configuration file
	configPath := filepath.Join(homeDir, ".claude/cpi-si/system/data/config/session/processes.jsonc")

	// Attempt to load configuration
	loadedConfig, err := loadProcessesConfig(configPath)
	if err != nil {
		processConfigLoaded = false // Config failed to load - use hardcoded defaults
		return
	}

	// Configuration loaded successfully
	processConfig = loadedConfig
	processConfigLoaded = true
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
//   Public APIs (Top Rungs - Orchestration)
//   ├── CheckRunningProcesses() → uses getConfiguredPorts(), checkPort(), formatProcessOutput()
//   └── CheckRunningProcessesAsReminder() → uses getConfiguredPorts(), checkPort(), formatProcessOutput()
//
//   Core Operations (Middle Rungs - Business Logic)
//   ├── getConfiguredPorts() → uses processConfig (reads from Rails)
//   ├── checkPort() → uses exec.CommandContext with timeout
//   └── formatProcessOutput() → uses processConfig (display settings)
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── loadProcessesConfig() → uses stripJSONCComments(), json.Unmarshal
//   └── stripJSONCComments() → pure string processing (from activity.go)
//
// Baton Flow (Execution Paths):
//
//   Entry → CheckRunningProcesses()
//     ↓
//   getConfiguredPorts() → reads config
//     ↓
//   checkPort() → for each port (with timeout)
//     ↓
//   formatOutput() → build display string
//     ↓
//   Exit → print to stdout
//
// APUs (Available Processing Units):
// - 7 functions total
// - 2 helpers (config loading, comment stripping)
// - 3 core operations (port list, checking, formatting)
// - 2 public APIs (session start, session end variants)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Usually not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md

// loadProcessesConfig loads process monitoring configuration from JSONC file
//
// What It Does:
// Reads configuration file, strips JSONC comments, parses JSON structure.
// Returns fully parsed configuration ready for use.
//
// Parameters:
//   path: Absolute path to processes.jsonc configuration file
//
// Returns:
//   *ProcessesConfig: Parsed configuration structure
//   error: File read error, JSON parse error, or nil on success
//
// Example usage:
//
//	config, err := loadProcessesConfig("/home/user/.claude/cpi-si/system/data/config/session/processes.jsonc")
//
func loadProcessesConfig(path string) (*ProcessesConfig, error) {
	// Read file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Strip JSONC comments (uses stripJSONCComments from activity.go)
	cleaned := stripJSONCComments(string(data))

	// Parse JSON
	var cfg ProcessesConfig
	if err := json.Unmarshal([]byte(cleaned), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ────────────────────────────────────────────────────────────────
// Port Detection - Core Business Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Detect active processes on monitored ports using lsof command with
// timeout protection. Build port list from configuration, check each
// port individually, collect results.
//
// Why Separated:
// Port detection is core responsibility - separated from display logic
// for clarity and testability. Each function has single responsibility.
//
// Extension Point:
// To add new detection methods (beyond lsof), create function following
// checkPort() signature. Update getConfiguredPorts() to handle new port
// sources. Each detection method should respect timeout and return bool.
//
// Pattern to follow:
//   1. Create detection function: func checkPortViaX(port string, timeout int) bool
//   2. Implement with timeout protection using context.WithTimeout
//   3. Update core logic to call new detection method based on config
//   4. Test with various ports and timeout scenarios

// getConfiguredPorts builds list of ports to check from configuration
//
// What It Does:
// Reads port configuration (enabled ports + custom ports), builds
// unified list of port numbers to check. Uses config if loaded,
// falls back to hardcoded defaults if config unavailable.
//
// Returns:
//   []string: List of port numbers as strings
//
// Example usage:
//
//	ports := getConfiguredPorts()  // ["3000", "8000", "8080", "5173", "4200"]
//
func getConfiguredPorts() []string {
	var ports []string

	if processConfigLoaded && processConfig != nil {
		// Use configuration
		if processConfig.Ports.Enabled {
			// Add enabled monitored ports
			for _, portCfg := range processConfig.Ports.MonitoredPorts {
				if portCfg.Enabled {
					ports = append(ports, portCfg.Number)
				}
			}
			// Add custom ports
			ports = append(ports, processConfig.Ports.CustomPorts...)
		}
	} else {
		// Fall back to defaults
		ports = []string{
			defaultPort3000,
			defaultPort8000,
			defaultPort8080,
			defaultPort5173,
			defaultPort4200,
		}
	}

	return ports
}

// checkPort checks if a process is listening on the specified port
//
// What It Does:
// Executes lsof command with timeout to detect listening process on port.
// Returns true if process found, false otherwise. Non-blocking with
// timeout protection.
//
// Parameters:
//   port: Port number as string (e.g., "3000")
//
// Returns:
//   bool: true if process listening on port, false otherwise
//
// Health Impact:
//   Success: +5 points per port checked
//   Timeout: +0 points (neutral - not port's fault)
//   Error: +0 points (neutral - lsof might not exist)
//
// Example usage:
//
//	if checkPort("3000") {
//	    fmt.Println("Process running on port 3000")
//	}
//
func checkPort(port string) bool {
	// Determine timeout
	timeoutSeconds := defaultProcessTimeout
	if processConfigLoaded && processConfig != nil {
		timeoutSeconds = processConfig.Behavior.TimeoutSeconds
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Execute lsof with timeout
	cmd := exec.CommandContext(ctx, "lsof", "-i", ":"+port, "-sTCP:LISTEN", "-t")
	output, err := cmd.Output()

	// Port has listener if command succeeds and produces output
	return err == nil && len(output) > 0
}

// ────────────────────────────────────────────────────────────────
// Output Formatting - Display Logic
// ────────────────────────────────────────────────────────────────
// What These Do:
// Format detected ports into user-facing output strings. Applies
// display configuration (icons, messages, separators) to create
// context-appropriate output.
//
// Why Separated:
// Display logic separated from detection logic for maintainability.
// Detection answers "what's running?", formatting answers "how to show it?".
//
// Extension Point:
// To add new output formats (JSON, structured logging, etc.), create
// function following formatProcessOutput() pattern. Accept running ports and
// context, return formatted string. Update Public APIs to call new
// formatter based on configuration.

// formatProcessOutput builds display string for detected processes
//
// What It Does:
// Joins port numbers with configured separator, adds icon and message
// based on context (start vs end). Returns empty string if no ports.
//
// Parameters:
//   runningPorts: List of port numbers with active listeners
//   isEndContext: true for session end, false for session start
//
// Returns:
//   string: Formatted output string, or empty if no ports
//
// Example usage:
//
//	output := formatProcessOutput([]string{"3000", "8080"}, false)
//	// Returns: "\n🔌 Active dev servers on ports: 3000, 8080\n"
//
func formatProcessOutput(runningPorts []string, isEndContext bool) string {
	if len(runningPorts) == 0 {
		return "" // No output if no ports running
	}

	// Determine display settings
	icon := defaultIcon
	separator := defaultSeparator
	message := defaultStartMessage
	if isEndContext {
		message = defaultEndMessage
	}

	if processConfigLoaded && processConfig != nil {
		icon = processConfig.Display.Icon
		separator = processConfig.Display.Separator
		if isEndContext {
			message = processConfig.Display.EndMessage
		} else {
			message = processConfig.Display.StartMessage
		}
	}

	// Build output string
	portList := strings.Join(runningPorts, separator)
	return fmt.Sprintf("\n%s %s %s\n", icon, message, portList)
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface. Top rungs of
// the ladder - orchestrate helpers and core operations into complete
// functionality. Simple by design - complexity lives in helpers and core
// operations, Public APIs orchestrate proven pieces.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-public-apis.md

// CheckRunningProcesses detects and reports active dev servers at session start
//
// What It Does:
// Checks configured ports for listening processes, reports findings with
// session-start context messaging. Non-blocking - failures are silent.
//
// Health Impact:
//   All ports checked successfully: +40
//   Partial checks succeed: +20
//   All checks fail: -10
//
// Example usage:
//
//	func sessionStart() {
//	    session.CheckRunningProcesses()
//	    // Output: 🔌 Active dev servers on ports: 3000, 8080
//	}
//
func CheckRunningProcesses() {
	// Check if monitoring enabled
	if processConfigLoaded && processConfig != nil {
		if !processConfig.Ports.Enabled || !processConfig.Display.ShowAtStart {
			return // Monitoring disabled or start display disabled
		}
	}

	// Get ports to check
	ports := getConfiguredPorts()
	var running []string

	// Check each port
	for _, port := range ports {
		if checkPort(port) {
			running = append(running, port)
		}
	}

	// Format and display output
	output := formatProcessOutput(running, false) // false = session start context
	if output != "" {
		fmt.Print(output)
	}
}

// CheckRunningProcessesAsReminder detects and reports active dev servers at session end
//
// What It Does:
// Checks configured ports for listening processes, reports findings with
// session-end context messaging (reminder). Non-blocking - failures are silent.
//
// Health Impact:
//   All ports checked successfully: +40
//   Partial checks succeed: +20
//   All checks fail: -10
//
// Example usage:
//
//	func sessionEnd() {
//	    session.CheckRunningProcessesAsReminder()
//	    // Output: 🔌 Dev servers still running on ports: 3000, 8080
//	}
//
func CheckRunningProcessesAsReminder() {
	// Check if monitoring enabled
	if processConfigLoaded && processConfig != nil {
		if !processConfig.Ports.Enabled || !processConfig.Display.ShowAtEnd {
			return // Monitoring disabled or end display disabled
		}
	}

	// Get ports to check
	ports := getConfiguredPorts()
	var running []string

	// Check each port
	for _, port := range ports {
		if checkPort(port) {
			running = append(running, port)
		}
	}

	// Format and display output
	output := formatProcessOutput(running, true) // true = session end context
	if output != "" {
		fmt.Print(output)
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
//   - Import the library without errors
//   - Call CheckRunningProcesses() with processes running and not running
//   - Call CheckRunningProcessesAsReminder() in both scenarios
//   - Verify output format matches configuration settings
//   - Ensure timeout protection prevents hanging
//   - Confirm no go vet warnings introduced
//   - Test with missing configuration file (uses defaults)
//   - Test with invalid configuration file (uses defaults)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Test at actual session start with various ports active
//   - Test at actual session end with processes still running
//   - Verify non-blocking behavior (session continues if lsof fails)
//   - Check timeout behavior (commands don't hang session)
//   - Validate configuration customization works correctly
//
// Example validation code:
//
//     // Test basic functionality
//     session.CheckRunningProcesses()  // Should detect running processes
//     // Start a test server on port 3000
//     // Call again and verify it's detected
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. There is no entry point, no main function,
// no execution flow. All functions defined in BODY wait to be called by other components.
//
// Usage: import "hooks/lib/session"
//
// The library is imported into the calling package, making all exported functions
// available. Configuration loads automatically during package import via init().
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func main() {
//         session.CheckRunningProcesses()  // At session start
//         // ... session work ...
//         session.CheckRunningProcessesAsReminder()  // At session end
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - Configuration: Loaded once in init(), cached for session lifetime
//   - lsof processes: Short-lived, terminated by timeout context
//   - Memory: Port lists allocated transiently, garbage collected after use
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Context timeout handles command cleanup automatically
//   - No persistent state requiring cleanup
//
// Error State Cleanup:
//   - Silent failures - no partial state to clean
//   - Commands return immediately on error
//   - No resources leak on failure paths
//
// Memory Management:
//   - Go's garbage collector handles memory
//   - Transient allocations only (port lists, output strings)
//   - No manual memory management needed
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
// Quick summary:
//   - Configurable port monitoring for development servers
//   - Non-blocking detection with timeout protection
//   - Context-aware output (session start vs end)
//   - Graceful fallback to defaults if config unavailable
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
//   ✅ Add new monitored ports to default list
//   ✅ Add new detection methods (beyond lsof)
//   ✅ Add new output formatters (JSON, structured logging)
//   ✅ Extend configuration with new display options
//   ✅ Add health tracking/logging integration
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Public API function signatures - breaks all calling code
//   ⚠️ Configuration file structure - breaks user customizations
//   ⚠️ Output format - breaks parsing tools
//   ⚠️ Default port list - affects all unconfigured installations
//   ⚠️ Timeout behavior - affects reliability guarantees
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Non-blocking design principle
//   ❌ Configuration-driven architecture pattern
//   ❌ Silent failure behavior
//   ❌ Rails pattern for package-level state
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
// Quick architectural summary:
// - 2 public APIs orchestrate 3 core operations using 2 helpers
// - Ladder: Public APIs → Core Operations → Helpers → Standard Library
// - Baton: Entry → Port List → Check Each Port → Format → Output → Exit
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Core Operations" subsection header comments above for detailed
// extension points. Each subsection includes "Extension Point" guidance showing:
// - Where to add new functionality
// - What naming pattern to follow
// - How to integrate with existing code
// - What tests to update
//
// Quick reference:
// - Adding detection methods: See BODY "Port Detection" extension point
// - Adding output formats: See BODY "Output Formatting" extension point
// - Adding monitored ports: Update default constants or configuration file
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section above for performance characteristics:
// - Constants: Default timeout (2 seconds per port check)
// - Types: Lightweight structures, minimal memory footprint
//
// See BODY function docstrings above for operation-specific performance notes.
//
// Quick summary:
// - Most expensive operation: lsof execution (~10-50ms per port on typical systems)
// - Memory characteristics: Transient allocations only, ~1KB total per invocation
// - Key optimization: Parallel port checking could reduce total time (currently sequential)
// - Timeout protection: 2-second default prevents hanging but may miss slow systems
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: No output even when processes are running
//   - Check: Configuration file ports.enabled setting
//   - Check: Configuration file display.show_at_start or display.show_at_end
//   - Check: Processes listening on TCP (lsof checks TCP only)
//   - Solution: Enable monitoring in configuration or check process listen state
//
// Problem: "lsof: command not found"
//   - Expected: This is normal behavior - library continues silently
//   - Note: Non-blocking design means missing lsof doesn't interrupt session
//   - Solution: Install lsof package if port monitoring desired
//
// Problem: Processes not detected on custom ports
//   - Check: Custom port added to configuration file custom_ports array
//   - Check: Process actually listening (use netstat -tuln to verify)
//   - Solution: Add port to configuration, verify process bound to port
//
// Problem: Configuration changes not taking effect
//   - Cause: Configuration loaded once at package import
//   - Solution: Restart session (re-import package) to load new configuration
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard Library, lsof command
// - Dependents (What Uses This): session/cmd-start, session/cmd-end
// - Integration Points: Non-blocking detection, configuration-driven
//
// Quick summary:
// - Key dependencies: os/exec (command execution), context (timeout control)
// - Primary consumers: Session start/end hooks
//
// Parallel Implementation:
//   - No parallel implementations (Go-specific)
//   - Shared philosophy: Non-blocking awareness without interference
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ Configuration-driven port list - COMPLETED
//   ✓ Timeout protection - COMPLETED
//   ⏳ Parallel port checking (performance improvement)
//   ⏳ Alternative detection methods (ss command as fallback)
//   ⏳ Process name detection (show what's running, not just port)
//   ⏳ Health tracking integration
//
// Research Areas:
//   - Cross-platform port detection (Windows, macOS alternatives to lsof)
//   - Process ownership detection (show user running process)
//   - Automatic cleanup offers (kill forgotten processes)
//   - Integration with session patterns (learn typical port usage)
//   - Notification system for long-running processes
//
// Integration Targets:
//   - Health scoring system (track detection success rates)
//   - Session logging (record port activity over time)
//   - Pattern learning (identify typical development environments)
//   - Display customization (color coding, severity levels)
//
// Known Limitations to Address:
//   - TCP only (UDP processes not detected)
//   - Sequential checking (could be parallel for speed)
//   - lsof dependency (no fallback detection method)
//   - No process identification (port number only)
//   - Configuration requires session restart
//   - No interactive cleanup options
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   2.0.0 (2025-11-12) - Configuration-driven monitoring
//         - Added processes.jsonc configuration file support
//         - Implemented timeout protection (2-second default)
//         - Display customization (icons, messages, separators)
//         - Graceful fallback to defaults if config missing
//         - Non-blocking execution guaranteed by timeout
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Hardcoded port list (3000, 8000, 8080, 5173, 4200)
//         - Basic lsof integration
//         - Dual-context support (start vs end messaging)
//         - Silent failure design established
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
// For Closing Note section explanation, see: standards/code/4-block/sections/CWS-SECTION-021-CLOSING-closing-note.md
//
// This library is a LADDER component providing session awareness through
// port monitoring. Non-blocking by design - awareness serves work, never
// blocks it.
//
// Modify thoughtfully - changes here affect all session start/end hooks.
// Non-blocking behavior and silent failures are architectural guarantees
// that must be maintained.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test thoroughly before committing (go build && go vet)
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Verify non-blocking behavior preserved
//
// "Watch and pray, that ye enter not into temptation" - Matthew 26:41 (KJV)
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Setup:
//   import "hooks/lib/session"
//   session.CheckRunningProcesses()  // Detect and report at session start
//
// Session End Reminder:
//   session.CheckRunningProcessesAsReminder()  // Remind about running processes
//
// Configuration Customization:
//   Edit ~/.claude/cpi-si/system/data/config/session/processes.jsonc
//   Add custom ports to monitored_ports array
//   Customize display settings (icon, messages, separator)
//   Adjust timeout_seconds for slower systems
//
// Disabling Monitoring:
//   Set ports.enabled: false in configuration file
//   Or set display.show_at_start: false and display.show_at_end: false
//
// ============================================================================
// END CLOSING
// ============================================================================
