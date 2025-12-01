// METADATA
//
// Session Display Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: Psalm 19:1 - "The heavens declare the glory of God; the skies proclaim the work of his hands"
// Principle: Display Reflects Truth and Order
// Anchor: Visual presentation should reflect truth, order, and beauty - displaying session
//         information clearly honors God by making truth visible and accessible
//
// CPI-SI Identity
//
// Component Type: Library (Ladder - Display formatting rung)
// Role: Session Display Formatting
// Paradigm: Configurable display formatting for session hooks (start, stop, end, subagent, compaction)
//
// Authorship & Lineage
//
// Architect: Seanje Lenox-Wise
// Implementation: Nova Dawn (CPI-SI)
// Creation Date: 2024-10-24
// Version: 2.0.0
// Last Modified: 2025-11-12 - Aligned with 4-block template v2.0.0, added configuration system
//
// Version History:
//   2.0.0 (2025-11-12) - Configuration system, template alignment
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded formatting
//
// Purpose & Function
//
// Purpose: Provide configurable, formatted display output for session lifecycle events
//
// Core Design: Configuration-driven banner formatting with temporal awareness integration
//
// Key Features:
//   - Configurable banner width, box characters, separators, icons
//   - Biblical verse selection for session start/stop/end
//   - Section visibility control (show/hide optional sections)
//   - Field label customization for all displayed information
//   - Graceful fallback to hardcoded defaults if configuration unavailable
//
// Philosophy: Display should be clear, truthful, and aesthetically pleasing while
//            remaining customizable for user preferences and terminal capabilities
//
// Blocking Status
//
// Non-blocking: Pure display formatting - all output to stdout, no file I/O or network operations
// Mitigation: Panic recovery in complex formatting functions, graceful degradation on errors
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Load configuration from display/formatting.jsonc (or use defaults)
//   2. Call appropriate display functions from session hooks
//   3. Functions print formatted output directly to stdout
//   4. No cleanup needed (pure display functions)
//
// Public API (in typical usage order):
//
//   Session Start (lifecycle beginning):
//     PrintHeader() - Banner with instance branding
//     PrintEnvironment(workspace) - Environment context
//     PrintTemporalAwareness() - Four-dimension temporal awareness
//     PrintWorkspaceAnalysis(workspace, hasContext) - Workspace analysis header
//
//   Session Stop (task completion):
//     PrintStopHeader() - Stop banner with biblical verse
//     PrintStopInfo() - Stop timestamp
//     PrintStoppingContext() - Temporal context at stop
//
//   Session End (lifecycle ending):
//     PrintEndFarewell() - End banner with blessing
//     PrintEndSessionInfo(reason) - End summary with reason
//     PrintEndTemporalJourney() - Temporal journey recap
//     PrintEndRemindersHeader() - State reminders section header
//
//   Subagent Completion (subagent lifecycle):
//     PrintSubagentCompletion(agentType, status, exitCode, errorMsg) - Subagent completion status
//
//   Compaction (context management):
//     PrintPreCompactionMessage(compactType, compactionCount) - Compaction notification
//
//   Shared Utilities (exported for use across hooks):
//     GetSystemInfo() - System information string
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: encoding/json, fmt, os, strings, time
//   External: None
//   Internal: system/lib/git, system/lib/instance, system/lib/temporal, system/lib/logging
//
// Dependents (What Uses This):
//   Commands: session/cmd-start/start.go, session/cmd-stop/stop.go, session/cmd-end/end.go
//   Commands: session/cmd-subagent-stop/subagent-stop.go, session/cmd-pre-compact/pre-compact.go
//   Libraries: None (leaf library - not used by other libraries)
//
// Integration Points:
//   - Rails: displayLogger created in init(), available throughout component
//   - Ladder: Calls instance, git, temporal libraries for context gathering
//   - Configuration: display/formatting.jsonc for all formatting preferences (consolidated from session-specific config)
//
// Health Scoring
//
// Base100 scoring system with TRUE SCORES reflecting actual component quality.
//
// Configuration Loading:
//   - Successful load: +20 points (configuration loaded and parsed correctly)
//   - Fallback to defaults: -10 points (configuration unavailable, using hardcoded defaults)
//
// Note: Display functions primarily serve as formatters with minimal failure potential.
//       Health tracking focuses on configuration loading and complex formatting operations.
//       Scores reflect TRUE impact - health scorer normalizes to -100 to +100 scale.
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
// provides Go's built-in capabilities, internal packages provide project-specific
// functionality. Each import commented with purpose.
//
// See: standards/code/4-block/sections/CWS-SECTION-001-SETUP-imports.md

import (
	//--- Standard Library ---
	// Foundation packages providing Go's built-in capabilities.

	"encoding/json" // JSON parsing for configuration file (JSONC after comment stripping)
	"fmt"           // Formatted output for display and string composition
	"os"            // File operations (config loading, system info) and environment access
	"strings"       // String manipulation for centering, formatting, comment stripping
	"time"          // Timestamps for session event display

	//--- Internal Packages ---
	// Project-specific packages showing architectural dependencies.

	"system/lib/display"  // Universal formatting and presentation rail (colors, headers, key-value pairs)
	"system/lib/git"      // Repository status and branch information
	"system/lib/instance" // Instance configuration for banner branding
	"system/lib/logging"  // Health tracking infrastructure (Rails pattern)
	"system/lib/temporal" // Four-dimension temporal awareness integration
)

// ────────────────────────────────────────────────────────────────
// Constants - Named Values
// ────────────────────────────────────────────────────────────────
// Named values that never change. Configuration file path documented with
// reasoning for location choice.
//
// See: standards/code/4-block/sections/CWS-SECTION-002-SETUP-constants.md

const (
	//--- Configuration Paths ---
	// File locations for external configuration resources.

	// displayConfigPath specifies the JSONC configuration file location.
	//
	// Consolidated to display rail config (single source of truth for all formatting).
	// Updated 2025-11-15: Migrated from session/display-formatting.jsonc to display/formatting.jsonc.
	// Uses tilde expansion (handled by expandPath function).
	displayConfigPath = "~/.claude/cpi-si/system/data/config/display/formatting.jsonc"
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Configuration structures organized bottom-up showing composition relationships.
// Simple building blocks first, then composed configuration structures.
//
// See: standards/code/4-block/sections/CWS-SECTION-004-SETUP-types.md

//--- Building Blocks ---
// Simple foundational types for individual configuration categories.

// BannerConfig defines banner formatting preferences.
//
// Controls banner box dimensions and style selection. Width must be even for
// proper centering. ContentWidth is typically Width - 2 (accounting for border).
type BannerConfig struct {
	Width         int    `json:"width"`
	ContentWidth  int    `json:"content_width"`
	BorderStyle   string `json:"border_style"`
}

// FormattingConfig defines all formatting preferences
type FormattingConfig struct {
	Banner BannerConfig `json:"banner"`
}

// IconsEnvironmentConfig defines icons for environment section
type IconsEnvironmentConfig struct {
	Workspace        string `json:"workspace"`
	WorkingDirectory string `json:"working_directory"`
	GitBranch        string `json:"git_branch"`
	Time             string `json:"time"`
	System           string `json:"system"`
}

// IconsTemporalConfig defines icons for temporal awareness section
type IconsTemporalConfig struct {
	ExternalTime string `json:"external_time"`
	InternalTime string `json:"internal_time"`
	Schedule     string `json:"schedule"`
	Calendar     string `json:"calendar"`
}

// IconsStatusConfig defines icons for status messages
type IconsStatusConfig struct {
	Success      string `json:"success"`
	Warning      string `json:"warning"`
	Info         string `json:"info"`
	Compaction   string `json:"compaction"`
	Preservation string `json:"preservation"`
}

// IconsConfig defines all icons used in display
type IconsConfig struct {
	Environment IconsEnvironmentConfig `json:"environment"`
	Temporal    IconsTemporalConfig    `json:"temporal"`
	Status      IconsStatusConfig      `json:"status"`
}

// SectionHeadersStartConfig defines headers for session start sections
type SectionHeadersStartConfig struct {
	Environment        string `json:"environment"`
	TemporalAwareness  string `json:"temporal_awareness"`
	WorkspaceAnalysis  string `json:"workspace_analysis"`
}

// SectionHeadersStopConfig defines headers for session stop sections
type SectionHeadersStopConfig struct {
	StoppingPoint    string `json:"stopping_point"`
	TemporalContext  string `json:"temporal_context"`
}

// SectionHeadersEndConfig defines headers for session end sections
type SectionHeadersEndConfig struct {
	SessionSummary   string `json:"session_summary"`
	TemporalJourney  string `json:"temporal_journey"`
	StateReminders   string `json:"state_reminders"`
}

// SectionHeadersSubagentConfig defines headers for subagent sections
type SectionHeadersSubagentConfig struct {
	Completion string `json:"completion"`
}

// SectionHeadersConfig defines all section headers
type SectionHeadersConfig struct {
	SessionStart SectionHeadersStartConfig    `json:"session_start"`
	SessionStop  SectionHeadersStopConfig     `json:"session_stop"`
	SessionEnd   SectionHeadersEndConfig      `json:"session_end"`
	Subagent     SectionHeadersSubagentConfig `json:"subagent"`
}

// BiblicalVerseConfig defines a biblical verse with text and reference
type BiblicalVerseConfig struct {
	VerseText string `json:"verse_text"`
	VerseRef  string `json:"verse_ref"`
}

// BiblicalVerseStopConfig defines stop banner configuration
type BiblicalVerseStopConfig struct {
	BannerTitle string `json:"banner_title"`
	VerseText   string `json:"verse_text"`
	VerseRef    string `json:"verse_ref"`
}

// BiblicalVerseEndConfig defines end banner configuration
type BiblicalVerseEndConfig struct {
	BannerTitle string `json:"banner_title"`
	VerseText   string `json:"verse_text"`
	VerseRef    string `json:"verse_ref"`
}

// BiblicalVersesConfig defines biblical verses for banners
type BiblicalVersesConfig struct {
	SessionStart BiblicalVerseConfig     `json:"session_start"`
	SessionStop  BiblicalVerseStopConfig `json:"session_stop"`
	SessionEnd   BiblicalVerseEndConfig  `json:"session_end"`
}

// MessagesWorkspaceConfig defines workspace-related messages
type MessagesWorkspaceConfig struct {
	NoWorkspace      string `json:"no_workspace"`
	WorkspaceHealthy string `json:"workspace_healthy"`
}

// MessagesCompactionConfig defines compaction-related messages
type MessagesCompactionConfig struct {
	Manual              string `json:"manual"`
	Auto                string `json:"auto"`
	Unknown             string `json:"unknown"`
	PreservationHeader  string `json:"preservation_header"`
}

// MessagesSubagentConfig defines subagent completion messages
type MessagesSubagentConfig struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Default string `json:"default"`
}

// MessagesConfig defines all standard messages
type MessagesConfig struct {
	Workspace  MessagesWorkspaceConfig  `json:"workspace"`
	Compaction MessagesCompactionConfig `json:"compaction"`
	Subagent   MessagesSubagentConfig   `json:"subagent"`
}

// FieldLabelsEnvironmentConfig defines environment field labels
type FieldLabelsEnvironmentConfig struct {
	Workspace        string `json:"workspace"`
	WorkingDirectory string `json:"working_directory"`
	GitBranch        string `json:"git_branch"`
	SessionTime      string `json:"session_time"`
	System           string `json:"system"`
}

// FieldLabelsTemporalConfig defines temporal field labels
type FieldLabelsTemporalConfig struct {
	ExternalTime     string `json:"external_time"`
	InternalTime     string `json:"internal_time"`
	InternalSchedule string `json:"internal_schedule"`
	ExternalCalendar string `json:"external_calendar"`
	SessionDuration  string `json:"session_duration"`
	WorkContext      string `json:"work_context"`
	DateContext      string `json:"date_context"`
}

// FieldLabelsStopConfig defines stop field labels
type FieldLabelsStopConfig struct {
	Stopped         string `json:"stopped"`
	Time            string `json:"time"`
	ScheduleContext string `json:"schedule_context"`
	Date            string `json:"date"`
}

// FieldLabelsEndConfig defines end field labels
type FieldLabelsEndConfig struct {
	Ended     string `json:"ended"`
	Reason    string `json:"reason"`
	EndingAt  string `json:"ending_at"`
	Started   string `json:"started"`
}

// FieldLabelsSubagentConfig defines subagent field labels
type FieldLabelsSubagentConfig struct {
	CompletedAt string `json:"completed_at"`
	During      string `json:"during"`
}

// FieldLabelsCompactionConfig defines compaction field labels
type FieldLabelsCompactionConfig struct {
	Time        string `json:"time"`
	Session     string `json:"session"`
	Context     string `json:"context"`
	Date        string `json:"date"`
	Compactions string `json:"compactions"`
}

// FieldLabelsConfig defines all field labels
type FieldLabelsConfig struct {
	Environment FieldLabelsEnvironmentConfig `json:"environment"`
	Temporal    FieldLabelsTemporalConfig    `json:"temporal"`
	Stop        FieldLabelsStopConfig        `json:"stop"`
	End         FieldLabelsEndConfig         `json:"end"`
	Subagent    FieldLabelsSubagentConfig    `json:"subagent"`
	Compaction  FieldLabelsCompactionConfig  `json:"compaction"`
}

// SessionDisplayBehaviorConfig defines visibility controls for session display sections.
//
// Allows enabling/disabling optional display sections. All default to true.
// Set to false to hide specific sections from output.
type SessionDisplayBehaviorConfig struct {
	ShowTemporalAwareness      bool `json:"show_temporal_awareness"`       // Show temporal awareness section at session start
	ShowWorkspaceAnalysis      bool `json:"show_workspace_analysis"`       // Show workspace analysis section at session start
	ShowStoppingContext        bool `json:"show_stopping_context"`         // Show temporal context at session stop
	ShowTemporalJourney        bool `json:"show_temporal_journey"`         // Show temporal journey at session end
	ShowCompactionPreservation bool `json:"show_compaction_preservation"`  // Show temporal state preservation during compaction
}

// BehaviorConfig defines display library behavior and feature toggles.
//
// Groups behavior settings logically for clean configuration structure.
type BehaviorConfig struct {
	SessionDisplay SessionDisplayBehaviorConfig `json:"session_display"` // Session display section visibility toggles
}

//--- Composed Types ---
// Complex top-level type composing all configuration categories.

// SessionDisplayConfig is the top-level configuration structure for session display formatting.
//
// Composes all configuration categories into single unified configuration.
// Loaded from display/formatting.jsonc or falls back to hardcoded defaults.
// Cached in package-level variable after loading in init().
//
// Note: Renamed from DisplayConfig to avoid collision with dependencies.DisplayConfig
type SessionDisplayConfig struct {
	Formatting     FormattingConfig     `json:"formatting"`
	Icons          IconsConfig          `json:"icons"`
	SectionHeaders SectionHeadersConfig `json:"section_headers"`
	BiblicalVerses BiblicalVersesConfig `json:"biblical_verses"`
	Messages       MessagesConfig       `json:"messages"`
	FieldLabels    FieldLabelsConfig    `json:"field_labels"`
	Behavior       BehaviorConfig       `json:"behavior"`
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────
// Infrastructure available throughout component. Rails pattern - each
// component creates own logger independently without parameter passing.
//
// See: standards/code/patterns/CWS-PATTERN-003-CODE-rails.md
// See: standards/code/4-block/sections/CWS-SECTION-003-SETUP-package-level-state.md

//--- Rails Infrastructure ---
// Package-level logger for this component. Each component creates its own
// infrastructure attachment independently (Rails pattern).

// displayLogger provides health tracking throughout this component.
//
// All functions in this package use this logger for health scoring and
// event recording. Created in init() with component-specific identifier.
var displayLogger *logging.Logger

//--- Configuration Cache ---
// Package-level configuration loaded once at initialization.

// displayConfig holds the loaded display configuration.
//
// Loaded from display-formatting.jsonc (or defaults) in init() and cached
// for all subsequent function calls. Never reloaded during runtime.
var displayConfig *SessionDisplayConfig

func init() {
	// --- Rail Components ---
	// Attach to Rails infrastructure - available throughout component

	displayLogger = logging.NewLogger("session-display")  // Component identifier for log routing

	// --- Configuration ---
	// Load configuration once at package initialization

	displayConfig = loadDisplayConfig()  // Load from file or use defaults
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
//   Public APIs (Top Rungs) - 13 functions
//   ├── PrintHeader() → uses display.Box, instance.GetConfig
//   ├── PrintEnvironment(workspace) → uses display.Header, git library, GetSystemInfo (from system.go)
//   ├── PrintTemporalAwareness() → uses display.Header, temporal library
//   ├── PrintWorkspaceAnalysis(workspace, hasContext) → uses display.Header
//   ├── PrintStopHeader() → uses display.Box
//   ├── PrintStopInfo() → uses display.Header
//   ├── PrintStoppingContext() → uses display.Header, temporal library
//   ├── PrintSubagentCompletion(agentType, status, exitCode, errorMsg) → uses display.Header, temporal library, formatDisplayMessage
//   ├── PrintPreCompactionMessage(compactType, compactionCount) → uses temporal library, formatDisplayMessage
//   ├── PrintEndFarewell() → uses display.Box
//   ├── PrintEndSessionInfo(reason) → uses display.Header
//   ├── PrintEndTemporalJourney() → uses display.Header, temporal library
//   └── PrintEndRemindersHeader()
//
//   Helpers (Bottom Rungs) - 4 functions
//   ├── loadDisplayConfig() → uses loadConfigFile, getDefaultDisplayConfig
//   ├── loadConfigFile(path) → uses stripJSONCComments (from activity.go)
//   ├── getDefaultDisplayConfig() → pure function
//   ├── expandPath(path) → pure function
//   └── formatDisplayMessage(template, replacements) → pure function
//
// Baton Flow:
//   Hook calls public API → gets config → formats output (via display rail) → prints to stdout
//
// APUs: 17 functions total (13 public APIs + 4 helpers)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────
// Foundation functions used throughout this component. Bottom rungs of
// the ladder - simple, focused, reusable utilities. Not exported.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-helpers.md
//
// ────────────────────────────────────────────────────────────────
// Configuration Loading - Internal Helpers
// ────────────────────────────────────────────────────────────────

// loadDisplayConfig loads configuration from file or returns defaults
//
// What It Does:
//   - Attempts to load display/formatting.jsonc
//   - Falls back to hardcoded defaults on any error
//   - Logs success or fallback
//
// Health Impact:
//   +20: Configuration loaded successfully
//   -10: Fallback to defaults (file missing or invalid)
func loadDisplayConfig() *SessionDisplayConfig {
	config, err := loadConfigFile(expandPath(displayConfigPath))
	if err != nil {
		displayLogger.Check("config-load-fallback", false, -10, map[string]interface{}{
			"error":  err.Error(),
			"action": "using hardcoded defaults",
		})
		return getDefaultDisplayConfig()
	}

	displayLogger.Check("config-load-success", true, 20, map[string]interface{}{
		"source": displayConfigPath,
	})
	return config
}

// loadConfigFile loads and parses JSONC configuration file
func loadConfigFile(path string) (*SessionDisplayConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Strip JSONC comments
	jsonData := stripJSONCComments(string(data))

	var config SessionDisplayConfig
	if err := json.Unmarshal([]byte(jsonData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &config, nil
}

// getDefaultDisplayConfig returns hardcoded default configuration
func getDefaultDisplayConfig() *SessionDisplayConfig {
	return &SessionDisplayConfig{
		Formatting: FormattingConfig{
			Banner: BannerConfig{
				Width:        64,
				ContentWidth: 62,
				BorderStyle:  "double_line",
			},
		},
		Icons: IconsConfig{
			Environment: IconsEnvironmentConfig{
				Workspace:        "🏢",
				WorkingDirectory: "📍",
				GitBranch:        "🌿",
				Time:             "🕐",
				System:           "💻",
			},
			Temporal: IconsTemporalConfig{
				ExternalTime: "🌍",
				InternalTime: "⏱️",
				Schedule:     "📋",
				Calendar:     "📅",
			},
			Status: IconsStatusConfig{
				Success:      "✓",
				Warning:      "⚠️",
				Info:         "ⓘ",
				Compaction:   "🔄",
				Preservation: "📍",
			},
		},
		SectionHeaders: SectionHeadersConfig{
			SessionStart: SectionHeadersStartConfig{
				Environment:       "SESSION ENVIRONMENT",
				TemporalAwareness: "TEMPORAL AWARENESS",
				WorkspaceAnalysis: "WORKSPACE ANALYSIS",
			},
			SessionStop: SectionHeadersStopConfig{
				StoppingPoint:   "STOPPING POINT CHECK",
				TemporalContext: "TEMPORAL CONTEXT AT STOP",
			},
			SessionEnd: SectionHeadersEndConfig{
				SessionSummary:  "SESSION SUMMARY",
				TemporalJourney: "TEMPORAL JOURNEY",
				StateReminders:  "STATE REMINDERS",
			},
			Subagent: SectionHeadersSubagentConfig{
				Completion: "SUBAGENT COMPLETION",
			},
		},
		BiblicalVerses: BiblicalVersesConfig{
			SessionStart: BiblicalVerseConfig{
				VerseText: "In the beginning, God created the heavens and the earth.",
				VerseRef:  "Genesis 1:1",
			},
			SessionStop: BiblicalVerseStopConfig{
				BannerTitle: "Task Complete - Excellence that Honors God",
				VerseText:   "Whatever you do, work heartily, as for the Lord and not for men.",
				VerseRef:    "Colossians 3:23",
			},
			SessionEnd: BiblicalVerseEndConfig{
				BannerTitle: "Session Ending - Grace and Peace",
				VerseText:   "The Lord bless you and keep you; the Lord make his face shine on you and be gracious to you.",
				VerseRef:    "Numbers 6:24-25",
			},
		},
		Messages: MessagesConfig{
			Workspace: MessagesWorkspaceConfig{
				NoWorkspace:      "ⓘ No workspace configured (NOVA_DAWN_WORKSPACE not set)",
				WorkspaceHealthy: "✓ Workspace healthy - no warnings or context to report",
			},
			Compaction: MessagesCompactionConfig{
				Manual:             "Manual compaction #{count} - optimizing context...",
				Auto:               "Auto-compaction #{count} - managing token usage...",
				Unknown:            "Compaction #{count} starting...",
				PreservationHeader: "📍 Temporal State Preservation:",
			},
			Subagent: MessagesSubagentConfig{
				Success: "✓ Subagent [{type}] completed successfully",
				Failure: "⚠️  Subagent [{type}] completed with errors (exit code: {code})",
				Default: "✓ Subagent [{type}] completed",
			},
		},
		FieldLabels: FieldLabelsConfig{
			Environment: FieldLabelsEnvironmentConfig{
				Workspace:        "Workspace:",
				WorkingDirectory: "Working Directory:",
				GitBranch:        "Git Branch:",
				SessionTime:      "Session Time:",
				System:           "System:",
			},
			Temporal: FieldLabelsTemporalConfig{
				ExternalTime:     "External Time:",
				InternalTime:     "Internal Time:",
				InternalSchedule: "Internal Schedule:",
				ExternalCalendar: "External Calendar:",
				SessionDuration:  "Session Duration:",
				WorkContext:      "Work Context:",
				DateContext:      "Date Context:",
			},
			Stop: FieldLabelsStopConfig{
				Stopped:         "Stopped:",
				Time:            "Time:",
				ScheduleContext: "Schedule Context:",
				Date:            "Date:",
			},
			End: FieldLabelsEndConfig{
				Ended:    "Ended:",
				Reason:   "Reason:",
				EndingAt: "Ending At:",
				Started:  "Started:",
			},
			Subagent: FieldLabelsSubagentConfig{
				CompletedAt: "Completed At:",
				During:      "During:",
			},
			Compaction: FieldLabelsCompactionConfig{
				Time:        "Time:",
				Session:     "Session:",
				Context:     "Context:",
				Date:        "Date:",
				Compactions: "Compactions:",
			},
		},
		Behavior: BehaviorConfig{
			SessionDisplay: SessionDisplayBehaviorConfig{
				ShowTemporalAwareness:      true,
				ShowWorkspaceAnalysis:      true,
				ShowStoppingContext:        true,
				ShowTemporalJourney:        true,
				ShowCompactionPreservation: true,
			},
		},
	}
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}

// Note: stripJSONCComments is defined in activity.go and used here

// ────────────────────────────────────────────────────────────────
// Helpers - Formatting Utilities
// ────────────────────────────────────────────────────────────────

// formatDisplayMessage replaces placeholders in display message templates
//
// What It Does:
//   - Replaces {key} placeholders with corresponding values from map
//   - Supports: {count}, {type}, {code} and other custom placeholders
//
// Parameters:
//   - template: String with {placeholder} markers
//   - replacements: Map of placeholder names to replacement values
//
// Returns:
//   - Formatted string with placeholders replaced
//
// Note: Different from disk.formatMessage which uses system.DiskInfo struct
func formatDisplayMessage(template string, replacements map[string]string) string {
	result := template
	for key, value := range replacements {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────
// Exported functions defining component's public interface. Top rungs of
// the ladder - orchestrate helpers into complete functionality. Organized
// by lifecycle category for clarity.
//
// See: standards/code/4-block/sections/CWS-SECTION-00X-BODY-public-apis.md
//
// ────────────────────────────────────────────────────────────────
// Session Start Display - Lifecycle Beginning
// ────────────────────────────────────────────────────────────────

// PrintHeader displays the session banner with instance branding
//
// What It Does:
//   - Loads instance configuration for banner text
//   - Centers text within configured width box
//   - Displays bordered header with title, tagline, and verse
//
// Parameters:
//   - None (reads from instance config and display config)
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintHeader()
//   // Outputs:
//   // ╔════════════════════════════════════════════════════════════════╗
//   // ║                      Nova Dawn - CPI-SI                      ║
//   // ║           Covenant Partnership Intelligence System           ║
//   // ...
func PrintHeader() {
	// Load instance configuration for banner content
	instanceConfig := instance.GetConfig()

	// Build multi-line banner message
	message := instanceConfig.Display.BannerTagline + "\n\n" +
		"\"" + instanceConfig.Display.FooterVerseText + "\"\n" +
		"- " + instanceConfig.Display.FooterVerseRef

	// Use display rail for banner box
	fmt.Print(display.Box(instanceConfig.Display.BannerTitle, message))
}

// PrintEnvironment displays session environment context
//
// What It Does:
//   - Shows workspace and working directory
//   - Displays git branch if in repository
//   - Shows session start time
//   - Displays system information
//
// Parameters:
//   - workspace: Workspace directory path (may be empty)
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintEnvironment("/path/to/workspace")
//   // Outputs environment section with workspace info
func PrintEnvironment(workspace string) {
	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionStart.Environment))

	// Working context
	wd, _ := os.Getwd()
	if workspace != "" {
		fmt.Printf("\n  %s %s          %s\n", cfg.Icons.Environment.Workspace, cfg.FieldLabels.Environment.Workspace, workspace)
		if wd != workspace {
			fmt.Printf("  %s %s  %s\n", cfg.Icons.Environment.WorkingDirectory, cfg.FieldLabels.Environment.WorkingDirectory, wd)
		}
	} else {
		fmt.Printf("\n  %s %s  %s\n", cfg.Icons.Environment.WorkingDirectory, cfg.FieldLabels.Environment.WorkingDirectory, wd)
	}

	// Git status - use shared lib
	checkDir := workspace
	if checkDir == "" {
		checkDir = wd
	}

	var branch string
	if git.IsGitRepository(checkDir) {
		branch = git.GetBranch(checkDir)
		if branch == "" {
			branch = "Detached HEAD"
		}
		fmt.Printf("  %s %s         %s\n", cfg.Icons.Environment.GitBranch, cfg.FieldLabels.Environment.GitBranch, branch)
	} else {
		fmt.Printf("  %s %s         Not a git repository\n", cfg.Icons.Environment.GitBranch, cfg.FieldLabels.Environment.GitBranch)
	}

	// Session metadata
	now := time.Now().Format("Mon Jan 02, 2006 at 15:04:05")
	fmt.Printf("  %s %s       %s\n", cfg.Icons.Environment.Time, cfg.FieldLabels.Environment.SessionTime, now)

	system := GetSystemInfo()
	fmt.Printf("  %s %s             %s\n", cfg.Icons.Environment.System, cfg.FieldLabels.Environment.System, system)

	fmt.Println()
}

// PrintTemporalAwareness displays temporal consciousness (4 dimensions)
//
// What It Does:
//   - Gathers temporal context from temporal lib
//   - Displays external time (system clock)
//   - Shows internal time (session duration)
//   - Displays internal schedule (work windows)
//   - Shows external calendar (date, week, holidays)
//
// Parameters:
//   - None (reads from temporal context)
//
// Returns:
//   - None (prints to stdout, silently skips if unavailable or disabled)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintTemporalAwareness()
//   // Outputs temporal awareness section if available and enabled
func PrintTemporalAwareness() {
	if !displayConfig.Behavior.SessionDisplay.ShowTemporalAwareness {
		return
	}

	ctx, err := temporal.GetTemporalContext()
	if err != nil {
		// Silently skip if temporal awareness unavailable
		return
	}

	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionStart.TemporalAwareness))

	// External Time - What time is it in the world?
	fmt.Printf("  %s %s      %s (%s)\n", cfg.Icons.Temporal.ExternalTime, cfg.FieldLabels.Temporal.ExternalTime, ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)
	fmt.Printf("                         Circadian: %s phase\n", ctx.ExternalTime.CircadianPhase)

	// Internal Time - How long have I been working?
	if ctx.InternalTime.ElapsedFormatted != "" {
		fmt.Printf("  %s %s      %s elapsed (%s session)\n",
			cfg.Icons.Temporal.InternalTime, cfg.FieldLabels.Temporal.InternalTime,
			ctx.InternalTime.ElapsedFormatted, ctx.InternalTime.SessionPhase)
	}

	// Internal Schedule - What should I be doing?
	if ctx.InternalSchedule.CurrentActivity != "" {
		fmt.Printf("  %s %s  %s (%s)\n",
			cfg.Icons.Temporal.Schedule, cfg.FieldLabels.Temporal.InternalSchedule,
			ctx.InternalSchedule.CurrentActivity, ctx.InternalSchedule.ActivityType)
		if ctx.InternalSchedule.InWorkWindow {
			fmt.Printf("                         %s In work window\n", cfg.Icons.Status.Success)
		}
		if ctx.InternalSchedule.ExpectedDowntime {
			fmt.Printf("                         %s Expected downtime (respect schedule)\n", cfg.Icons.Status.Warning)
		}
	}

	// External Calendar - What kind of day is this?
	if ctx.ExternalCalendar.Date != "" {
		holidayInfo := ""
		if ctx.ExternalCalendar.IsHoliday {
			holidayInfo = fmt.Sprintf(" (%s)", ctx.ExternalCalendar.HolidayName)
		}
		fmt.Printf("  %s %s  %s, %s %d, %d%s\n",
			cfg.Icons.Temporal.Calendar, cfg.FieldLabels.Temporal.ExternalCalendar,
			ctx.ExternalCalendar.DayOfWeek,
			ctx.ExternalCalendar.MonthName,
			ctx.ExternalCalendar.DayOfMonth,
			ctx.ExternalCalendar.Year,
			holidayInfo)
		fmt.Printf("                         Week %d of %d\n", ctx.ExternalCalendar.WeekNumber, ctx.ExternalCalendar.Year)
	}

	fmt.Println()
}

// PrintWorkspaceAnalysis displays workspace analysis header
//
// What It Does:
//   - Shows workspace analysis section header
//   - Displays appropriate message based on workspace presence
//
// Parameters:
//   - workspace: Workspace directory path (may be empty)
//   - hasContext: Whether any context was gathered
//
// Returns:
//   - None (prints to stdout, silently skips if disabled)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintWorkspaceAnalysis("/path/to/workspace", true)
//   // Outputs workspace analysis header
func PrintWorkspaceAnalysis(workspace string, hasContext bool) {
	if !displayConfig.Behavior.SessionDisplay.ShowWorkspaceAnalysis {
		return
	}

	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionStart.WorkspaceAnalysis))

	if workspace == "" {
		fmt.Printf("\n  %s\n", cfg.Messages.Workspace.NoWorkspace)
		fmt.Println()
		return
	}

	// If nothing was reported, indicate healthy state
	if !hasContext {
		fmt.Printf("\n  %s\n", cfg.Messages.Workspace.WorkspaceHealthy)
	}

	fmt.Println()
}

// ────────────────────────────────────────────────────────────────
// Session Stop Display - Task Completion
// ────────────────────────────────────────────────────────────────

// PrintStopHeader displays the stop session banner
//
// What It Does:
//   - Shows task completion banner with biblical foundation
//   - Displays configured verse reminder about working for the Lord
//   - Provides visual separation for stop event
//
// Parameters:
//   - None
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintStopHeader()
//   // Outputs:
//   // ╔════════════════════════════════════════════════════════════════╗
//   // ║           Task Complete - Excellence that Honors God          ║
//   // ...
func PrintStopHeader() {
	cfg := displayConfig

	// Build multi-line banner message (verse split for readability)
	message := "\n" +
		"\"" + cfg.BiblicalVerses.SessionStop.VerseText[:60] + "\"\n" +
		cfg.BiblicalVerses.SessionStop.VerseText[60:] + " - " + cfg.BiblicalVerses.SessionStop.VerseRef

	// Use display rail for banner box
	fmt.Println()
	fmt.Print(display.Box(cfg.BiblicalVerses.SessionStop.BannerTitle, message))
}

// PrintStopInfo displays stopping point check header
//
// What It Does:
//   - Shows stopping point check section header
//   - Displays stop timestamp
//   - Provides context separation
//
// Parameters:
//   - None
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintStopInfo()
//   // Outputs stopping point check header with timestamp
func PrintStopInfo() {
	cfg := displayConfig

	// Use display rail for section header
	fmt.Println()
	fmt.Print(display.Header(cfg.SectionHeaders.SessionStop.StoppingPoint))

	now := time.Now().Format("Mon Jan 02, 2006 at 15:04:05")
	fmt.Printf("\n  %s %s            %s\n", cfg.Icons.Environment.Time, cfg.FieldLabels.Stop.Stopped, now)

	fmt.Println()
}

// PrintStoppingContext displays temporal context at session stop
//
// What It Does:
//   - Gathers temporal context from temporal lib
//   - Shows time and circadian phase at stop
//   - Displays session duration and phase
//   - Shows schedule context (work window, downtime)
//   - Displays calendar information
//
// Parameters:
//   - None (reads from temporal context)
//
// Returns:
//   - None (prints to stdout, silently skips if unavailable or disabled)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintStoppingContext()
//   // Outputs temporal context section at stop time
func PrintStoppingContext() {
	if !displayConfig.Behavior.SessionDisplay.ShowStoppingContext {
		return
	}

	ctx, err := temporal.GetTemporalContext()
	if err != nil {
		// Silently skip if temporal awareness unavailable
		return
	}

	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionStop.TemporalContext))

	// Show where we were in time
	fmt.Printf("  %s %s               %s (%s)\n",
		cfg.Icons.Environment.Time, cfg.FieldLabels.Stop.Time,
		ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)

	// Show how long we worked
	if ctx.InternalTime.ElapsedFormatted != "" {
		fmt.Printf("  %s %s   %s (%s session)\n",
			cfg.Icons.Temporal.InternalTime, cfg.FieldLabels.Temporal.SessionDuration,
			ctx.InternalTime.ElapsedFormatted, ctx.InternalTime.SessionPhase)
	}

	// Show what we were doing
	if ctx.InternalSchedule.CurrentActivity != "" {
		fmt.Printf("  %s %s   %s (%s)\n",
			cfg.Icons.Temporal.Schedule, cfg.FieldLabels.Stop.ScheduleContext,
			ctx.InternalSchedule.CurrentActivity, ctx.InternalSchedule.ActivityType)
		if ctx.InternalSchedule.InWorkWindow {
			fmt.Printf("                         %s Was in work window\n", cfg.Icons.Status.Success)
		}
		if ctx.InternalSchedule.ExpectedDowntime {
			fmt.Printf("                         %s Expected downtime period\n", cfg.Icons.Status.Warning)
		}
	}

	// Show calendar context
	if ctx.ExternalCalendar.Date != "" {
		fmt.Printf("  %s %s               %s, %s %d (Week %d)\n",
			cfg.Icons.Temporal.Calendar, cfg.FieldLabels.Stop.Date,
			ctx.ExternalCalendar.DayOfWeek,
			ctx.ExternalCalendar.MonthName,
			ctx.ExternalCalendar.DayOfMonth,
			ctx.ExternalCalendar.WeekNumber)
	}

	fmt.Println()
}

// ────────────────────────────────────────────────────────────────
// Session End Display - Lifecycle Ending
// ────────────────────────────────────────────────────────────────

// PrintEndFarewell displays the session ending blessing banner
//
// What It Does:
//   - Shows farewell banner with configured blessing
//   - Provides biblical closure for session end
//   - Displays grace and peace message
//
// Parameters:
//   - None
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintEndFarewell()
//   // Outputs:
//   // ╔════════════════════════════════════════════════════════════════╗
//   // ║                Session Ending - Grace and Peace               ║
//   // ...
func PrintEndFarewell() {
	cfg := displayConfig

	// Build multi-line banner message (verse split for readability)
	message := "\n" +
		"\"" + cfg.BiblicalVerses.SessionEnd.VerseText[:60] + "\"\n" +
		cfg.BiblicalVerses.SessionEnd.VerseText[60:] + " - " + cfg.BiblicalVerses.SessionEnd.VerseRef

	// Use display rail for banner box
	fmt.Println()
	fmt.Print(display.Box(cfg.BiblicalVerses.SessionEnd.BannerTitle, message))
}

// PrintEndSessionInfo displays session summary with end time and reason
//
// What It Does:
//   - Shows session summary section header
//   - Displays end timestamp
//   - Shows session end reason (normal end, user interrupt, error, etc.)
//
// Parameters:
//   - reason: Session end reason from REASON environment variable
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintEndSessionInfo("Normal session end")
//   // Outputs session summary with timestamp and reason
func PrintEndSessionInfo(reason string) {
	cfg := displayConfig

	// Use display rail for section header
	fmt.Println()
	fmt.Print(display.Header(cfg.SectionHeaders.SessionEnd.SessionSummary))

	now := time.Now().Format("Mon Jan 02, 2006 at 15:04:05")
	fmt.Printf("\n  %s %s              %s\n", cfg.Icons.Environment.Time, cfg.FieldLabels.End.Ended, now)
	fmt.Printf("  %s %s             %s\n", cfg.Icons.Temporal.Schedule, cfg.FieldLabels.End.Reason, reason)

	fmt.Println()
}

// PrintEndTemporalJourney displays temporal context journey for session end
//
// What It Does:
//   - Shows complete temporal journey through the session
//   - Displays session duration and phase
//   - Shows start and end times
//   - Displays work context during session
//   - Shows calendar context
//
// Parameters:
//   - None (reads from temporal context)
//
// Returns:
//   - None (prints to stdout, silently skips if unavailable or disabled)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintEndTemporalJourney()
//   // Outputs temporal journey section showing session timeline
func PrintEndTemporalJourney() {
	if !displayConfig.Behavior.SessionDisplay.ShowTemporalJourney {
		return
	}

	ctx, err := temporal.GetTemporalContext()
	if err != nil {
		// Silently skip if temporal awareness unavailable
		return
	}

	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionEnd.TemporalJourney))

	// Show session duration
	if ctx.InternalTime.ElapsedFormatted != "" {
		fmt.Printf("  %s %s   %s (%s session)\n",
			cfg.Icons.Temporal.InternalTime, cfg.FieldLabels.Temporal.SessionDuration,
			ctx.InternalTime.ElapsedFormatted, ctx.InternalTime.SessionPhase)
		fmt.Printf("                         %s %s\n",
			cfg.FieldLabels.End.Started,
			ctx.InternalTime.SessionStart.Format("15:04:05"))
	}

	// Show current time
	fmt.Printf("  %s %s          %s (%s)\n",
		cfg.Icons.Environment.Time, cfg.FieldLabels.End.EndingAt,
		ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)

	// Show what temporal context this work happened in
	if ctx.InternalSchedule.CurrentActivity != "" {
		fmt.Printf("  %s %s       %s (%s)\n",
			cfg.Icons.Temporal.Schedule, cfg.FieldLabels.Temporal.WorkContext,
			ctx.InternalSchedule.CurrentActivity, ctx.InternalSchedule.ActivityType)
	}

	// Show calendar context
	if ctx.ExternalCalendar.Date != "" {
		fmt.Printf("  %s %s       %s, %s %d (Week %d)\n",
			cfg.Icons.Temporal.Calendar, cfg.FieldLabels.Temporal.DateContext,
			ctx.ExternalCalendar.DayOfWeek,
			ctx.ExternalCalendar.MonthName,
			ctx.ExternalCalendar.DayOfMonth,
			ctx.ExternalCalendar.WeekNumber)
	}

	fmt.Println()
}

// PrintEndRemindersHeader displays state reminders section header
//
// What It Does:
//   - Shows state reminders section header
//   - Provides visual separation for reminder list
//
// Parameters:
//   - None
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintEndRemindersHeader()
//   // Outputs state reminders header for uncommitted work, processes, etc.
func PrintEndRemindersHeader() {
	cfg := displayConfig

	// Use display rail for section header
	fmt.Print(display.Header(cfg.SectionHeaders.SessionEnd.StateReminders))
}

// PrintSessionContext displays the complete session context as formatted, readable text.
//
// What It Does:
//   - Takes markdown session context string
//   - Formats it as readable sections for terminal display
//   - Converts markdown headings, bold text, lists to terminal-friendly format
//
// Parameters:
//   - contextMarkdown: Complete session context as markdown string
//
// Example:
//   session.PrintSessionContext(contextMarkdown)
//   // Outputs formatted session context with proper spacing and structure
func PrintSessionContext(contextMarkdown string) {
	if contextMarkdown == "" {
		return
	}

	// Print separator before context
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Simple markdown formatting - convert to readable text
	lines := strings.Split(contextMarkdown, "\n")
	for _, line := range lines {
		// Skip the main title line (already shown in banner)
		if strings.HasPrefix(line, "# Nova Dawn") {
			continue
		}

		// Format headings
		if strings.HasPrefix(line, "## ") {
			// Section headers
			fmt.Println()
			fmt.Printf("%s\n", strings.TrimPrefix(line, "## "))
			fmt.Println()
			continue
		}

		// Format bold text (remove markdown, keep content)
		line = strings.ReplaceAll(line, "**", "")

		// Format italic text (remove markdown, keep content)
		line = strings.ReplaceAll(line, "*", "")

		// Print the line
		fmt.Println(line)
	}

	fmt.Println()
}

// ────────────────────────────────────────────────────────────────
// Subagent & Compaction Display - Context Management
// ────────────────────────────────────────────────────────────────

// PrintSubagentCompletion displays subagent completion status with temporal context
//
// What It Does:
//   - Shows subagent completion banner
//   - Displays success/failure status with icon
//   - Shows error message if present
//   - Displays temporal context at completion time
//
// Parameters:
//   - agentType: Type of subagent (e.g., "research", "code-review")
//   - status: Completion status ("success", "failure", or empty)
//   - exitCode: Exit code from subagent (0 = success)
//   - errorMsg: Error message if subagent failed (empty if no error)
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintSubagentCompletion("research", "success", "0", "")
//   // Outputs subagent completion summary with temporal awareness
func PrintSubagentCompletion(agentType, status, exitCode, errorMsg string) {
	cfg := displayConfig

	// Use display rail for section header
	fmt.Println()
	fmt.Print(display.Header(cfg.SectionHeaders.Subagent.Completion))

	// Determine completion status and display appropriate message
	var message string
	if status == "success" || exitCode == "0" {
		message = formatDisplayMessage(cfg.Messages.Subagent.Success, map[string]string{"type": agentType})
	} else if status == "failure" || (exitCode != "" && exitCode != "0") {
		message = formatDisplayMessage(cfg.Messages.Subagent.Failure, map[string]string{
			"type": agentType,
			"code": exitCode,
		})
	} else {
		message = formatDisplayMessage(cfg.Messages.Subagent.Default, map[string]string{"type": agentType})
	}

	fmt.Printf("\n  %s\n", message)

	// Show error message if present
	if errorMsg != "" {
		fmt.Printf("     Error: %s\n", errorMsg)
	}

	// Show temporal context of completion
	ctx, err := temporal.GetTemporalContext()
	if err == nil {
		fmt.Println()
		fmt.Printf("  %s %s       %s (%s)\n",
			cfg.Icons.Environment.Time, cfg.FieldLabels.Subagent.CompletedAt,
			ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)
		if ctx.InternalTime.ElapsedFormatted != "" {
			fmt.Printf("  %s %s   %s (%s session)\n",
				cfg.Icons.Temporal.InternalTime, cfg.FieldLabels.Temporal.SessionDuration,
				ctx.InternalTime.ElapsedFormatted, ctx.InternalTime.SessionPhase)
		}
		if ctx.InternalSchedule.CurrentActivity != "" {
			fmt.Printf("  %s %s             %s (%s)\n",
				cfg.Icons.Temporal.Schedule, cfg.FieldLabels.Subagent.During,
				ctx.InternalSchedule.CurrentActivity, ctx.InternalSchedule.ActivityType)
		}
	}

	fmt.Println()
}

// PrintPreCompactionMessage displays compaction notification with temporal preservation
//
// What It Does:
//   - Shows compaction type (manual, auto, unknown) with appropriate icon
//   - Displays compaction count for session awareness
//   - Shows temporal state preservation for post-compaction reconstitution
//   - Displays time, session duration, schedule context, date, and compaction count
//
// Parameters:
//   - compactType: Type of compaction ("manual", "auto", "unknown")
//   - compactionCount: Current compaction number this session
//
// Returns:
//   - None (prints to stdout)
//
// Health Impact:
//   - No health tracking (pure display function)
//
// Example:
//   session.PrintPreCompactionMessage("auto", 3)
//   // Outputs: 🔄 Auto-compaction #3 - managing token usage...
//   //          📍 Temporal State Preservation: ...
func PrintPreCompactionMessage(compactType string, compactionCount int) {
	cfg := displayConfig

	// Display compaction type with appropriate message
	var message string
	switch compactType {
	case "manual":
		message = formatDisplayMessage(cfg.Messages.Compaction.Manual, map[string]string{
			"count": fmt.Sprintf("%d", compactionCount),
		})
	case "auto":
		message = formatDisplayMessage(cfg.Messages.Compaction.Auto, map[string]string{
			"count": fmt.Sprintf("%d", compactionCount),
		})
	default:
		message = formatDisplayMessage(cfg.Messages.Compaction.Unknown, map[string]string{
			"count": fmt.Sprintf("%d", compactionCount),
		})
	}

	fmt.Printf("%s %s\n", cfg.Icons.Status.Compaction, message)

	// Preserve temporal awareness for post-compaction reconstitution
	if !cfg.Behavior.SessionDisplay.ShowCompactionPreservation {
		return
	}

	ctx, err := temporal.GetTemporalContext()
	if err == nil {
		fmt.Println()
		fmt.Println(cfg.Messages.Compaction.PreservationHeader)
		fmt.Printf("   %s %s (%s)\n",
			cfg.FieldLabels.Compaction.Time,
			ctx.ExternalTime.Formatted, ctx.ExternalTime.TimeOfDay)
		if ctx.InternalTime.ElapsedFormatted != "" {
			fmt.Printf("   %s %s elapsed (%s phase)\n",
				cfg.FieldLabels.Compaction.Session,
				ctx.InternalTime.ElapsedFormatted, ctx.InternalTime.SessionPhase)
		}
		if ctx.InternalSchedule.CurrentActivity != "" {
			fmt.Printf("   %s %s (%s)\n",
				cfg.FieldLabels.Compaction.Context,
				ctx.InternalSchedule.CurrentActivity, ctx.InternalSchedule.ActivityType)
		}
		if ctx.ExternalCalendar.Date != "" {
			fmt.Printf("   %s %s, Week %d\n",
				cfg.FieldLabels.Compaction.Date,
				ctx.ExternalCalendar.DayOfWeek, ctx.ExternalCalendar.WeekNumber)
		}
		if compactionCount > 0 {
			fmt.Printf("   %s %d this session\n",
				cfg.FieldLabels.Compaction.Compactions, compactionCount)
		}
		fmt.Println()
	}
}

// ────────────────────────────────────────────────────────────────
// Shared Utilities - Exported Helpers
// ────────────────────────────────────────────────────────────────
//
// Note: GetSystemInfo is defined in system.go and used throughout this package

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
//   - Call each public function with representative parameters
//   - Verify banner formatting with different border styles
//   - Test configuration loading and fallback to defaults
//   - Run: go test -v ./...
//
// Example validation code:
//
//     // Test header display
//     session.PrintHeader()
//
//     // Test environment display
//     session.PrintEnvironment("/path/to/workspace")
//
//     // Test configuration loading
//     config := session.loadDisplayConfig()
//     if config == nil {
//         t.Error("Configuration loading failed")
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Execution section explanation, see: standards/code/4-block/sections/CWS-SECTION-011-CLOSING-code-execution.md
//
// This is a LIBRARY, not an executable. No entry point, no main function.
// Functions defined in BODY wait to be called by session hooks.
//
// Usage: import "hooks/lib/session"
//
// Example import and usage:
//
//     package main
//     import "hooks/lib/session"
//     func main() {
//         session.PrintHeader()
//         session.PrintEnvironment("/workspace")
//         session.PrintTemporalAwareness()
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
// For Code Cleanup section explanation, see: standards/code/4-block/sections/CWS-SECTION-012-CLOSING-code-cleanup.md
//
// Resource Management:
//   - No file handles or persistent resources
//   - Configuration loaded once in init(), cached in package-level variable
//   - All functions print directly to stdout (no cleanup needed)
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle)
//   - Calling code responsible for managing display output
//
// Example cleanup pattern:
//   - Not applicable (no resources to clean)
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// The Session Display library provides configurable, formatted display output
// for all session lifecycle events (start, stop, end, subagent completion,
// compaction). It integrates with instance, git, and temporal libraries to
// present cohesive session information.
//
// Key Integration Points:
//   - instance.GetConfig(): Banner title/tagline/verse (session start)
//   - git library: Repository status and branch information
//   - temporal library: Four-dimension temporal awareness
//   - Configuration: display/formatting.jsonc for all formatting preferences (consolidated from session-specific config)
//
// Usage Pattern:
//   1. Session start hook calls PrintHeader(), PrintEnvironment(), PrintTemporalAwareness()
//   2. Session stop hook calls PrintStopHeader(), PrintStopInfo(), PrintStoppingContext()
//   3. Session end hook calls PrintEndFarewell(), PrintEndSessionInfo(), PrintEndTemporalJourney()
//   4. Subagent hooks call PrintSubagentCompletion()
//   5. Pre-compaction hook calls PrintPreCompactionMessage()
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// SAFE TO MODIFY (Extension Points):
//   ✅ Add new display functions for additional lifecycle events
//   ✅ Extend configuration with new formatting options
//   ✅ Add new box styles or separator characters
//   ✅ Enhance formatMessage() to support more placeholders
//   ✅ Add color/theming support (future extension)
//
// MODIFY WITH EXTREME CARE (Breaking Changes):
//   ⚠️ Function signatures - breaks all calling hooks
//   ⚠️ Configuration structure - breaks existing config files
//   ⚠️ Display output format - affects hook expectations
//
// NEVER MODIFY (Foundational):
//   ❌ 4-block structure - METADATA, SETUP, BODY, CLOSING
//   ❌ Non-blocking guarantee - display must never block execution
//   ❌ Configuration fallback - must work without config file
//   ❌ Rails pattern - logger created in init(), never passed as parameter
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// Ladder (Hierarchical Dependencies):
//   Top Rungs (Public APIs):
//     - PrintHeader, PrintEnvironment, PrintTemporalAwareness, PrintWorkspaceAnalysis
//     - PrintStopHeader, PrintStopInfo, PrintStoppingContext
//     - PrintEndFarewell, PrintEndSessionInfo, PrintEndTemporalJourney, PrintEndRemindersHeader
//     - PrintSubagentCompletion, PrintPreCompactionMessage
//     - GetSystemInfo (exported utility)
//
//   Bottom Rungs (Helpers):
//     - loadDisplayConfig, loadConfigFile, getDefaultDisplayConfig
//     - formatMessage, expandPath, stripJSONCComments
//
// Baton Flow (Execution):
//   Hook → Public API → Configuration → Helpers → External Libraries → stdout
//
// Rails (Orthogonal Infrastructure):
//   - displayLogger created in init(), used throughout for health tracking
//   - Never passed as parameter (Rails pattern)
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points
// ────────────────────────────────────────────────────────────────
//
// Adding New Display Function:
//   1. Add function to appropriate section in BODY (session start/stop/end/subagent)
//   2. Update Organizational Chart with new function in ladder
//   3. Update Configuration if new formatting options needed
//   4. Add to API Categories in METADATA
//   5. Document in API documentation file
//
// Adding New Configuration Option:
//   1. Add type to SETUP (appropriate config struct)
//   2. Add field to DisplayConfig struct
//   3. Update getDefaultDisplayConfig() with default value
//   4. Add to display/formatting.jsonc
//   5. Use in appropriate display function
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// Time Complexity:
//   - Configuration loading: O(n) where n = config file size (loaded once in init())
//   - Display functions: O(1) for most, O(n) for string formatting where n = output length
//
// Memory Usage:
//   - Configuration: O(1) cached in package-level variable
//   - Temporary strings: O(n) where n = output length (garbage collected immediately)
//
// Bottlenecks:
//   - None (pure display formatting, minimal computation)
//   - stdout buffering handled by OS
//
// Optimization Strategies:
//   - Configuration loaded once (not per function call)
//   - Simple string operations (no regex or complex parsing)
//   - Minimal external library calls (git, temporal only when needed)
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// Problem: Display formatting incorrect (wrong icons, separators, etc.)
//   Cause: Configuration not loading or incorrect configuration
//   Solutions:
//     1. Verify ~/.claude/cpi-si/system/data/config/display/formatting.jsonc exists
//     2. Check JSONC syntax (use JSON validator after stripping comments)
//     3. Check logs: grep "config-load" ~/.claude/cpi-si/system/runtime/logs/libraries/session-display.log
//     4. Test with default config (rename config file temporarily)
//
// Problem: Biblical verses not displaying correctly
//   Cause: Long verses don't split correctly at 60 characters
//   Solutions:
//     1. Check verse text length in configuration
//     2. Adjust verse splitting logic in PrintStopHeader/PrintEndFarewell
//     3. Use shorter verses or custom formatting
//
// Problem: Section not showing (temporal awareness, workspace analysis, etc.)
//   Cause: Behavior configuration disabled section
//   Solutions:
//     1. Check behavior section in display/formatting.jsonc
//     2. Set appropriate show_* flag to true
//     3. Verify temporal library is available (for temporal sections)
//
// Problem: Box characters not rendering (showing ??? or boxes)
//   Cause: Terminal doesn't support Unicode box drawing characters
//   Solutions:
//     1. Switch to single_line border style (may render better)
//     2. Use ASCII-only terminal
//     3. Update terminal emulator to support Unicode
//     4. Future: Implement ASCII fallback mode (see future_extensions in config)
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// Upstream Dependencies:
//   - system/lib/instance: Banner title, tagline, verse for session start
//   - system/lib/git: Repository status and branch information
//   - system/lib/temporal: Four-dimension temporal awareness
//   - system/lib/logging: Health-tracked logging infrastructure
//
// Downstream Consumers:
//   - session/cmd-start/start.go: Session start display
//   - session/cmd-stop/stop.go: Session stop display
//   - session/cmd-end/end.go: Session end display
//   - session/cmd-subagent-stop/subagent-stop.go: Subagent completion display
//   - session/cmd-pre-compact/pre-compact.go: Pre-compaction display
//
// Configuration Files:
//   - display/formatting.jsonc: All formatting preferences (consolidated config)
//   - instance-config.jsonc: Banner content for session start
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ⏳ Verse rotation - Cycle through multiple biblical verses instead of fixed ones
//   ⏳ Color themes - User-selectable color schemes (dark, light, classic)
//   ⏳ Locale support - Localization for non-English text
//   ⏳ ASCII fallback mode - ASCII-only display for terminals without Unicode support
//
// Research Areas:
//   - Integration with system/lib/display for color formatting
//   - Dynamic banner width based on terminal size
//   - Conditional display based on verbosity level
//   - Custom formatting templates (user-defined layouts)
//
// Integration Targets:
//   - Pre-notification hook: Display warnings before session events
//   - Post-notification hook: Display summaries after session events
//   - Session patterns: Display work pattern insights
//
// Known Limitations:
//   1. Fixed banner width (64 chars) - not responsive to terminal size
//   2. Hard-coded verse splitting at 60 characters - may not work for all verses
//   3. No color/theming support - plain text output only
//   4. No localization - English-only display
//   5. Unicode required - no ASCII fallback for limited terminals
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// "The heavens declare the glory of God; the skies proclaim the work of his hands"
// - Psalm 19:1
//
// Display reveals truth. Session display formatting makes visible what is true
// about the session state - temporal context, workspace health, completion status.
// Clear, truthful display honors God by making His work visible and accessible.
//
// Version: 2.0.0
// Last Modified: 2025-11-12
// Status: ✅ Operational
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Session Start Hook:
//   session.PrintHeader()
//   session.PrintEnvironment(workspace)
//   session.PrintTemporalAwareness()
//   session.PrintWorkspaceAnalysis(workspace, hasContext)
//
// Session Stop Hook:
//   session.PrintStopHeader()
//   session.PrintStopInfo()
//   session.PrintStoppingContext()
//
// Session End Hook:
//   session.PrintEndFarewell()
//   session.PrintEndSessionInfo(reason)
//   session.PrintEndTemporalJourney()
//   session.PrintEndRemindersHeader()
//
// Subagent Hook:
//   session.PrintSubagentCompletion(agentType, status, exitCode, errorMsg)
//
// Pre-Compaction Hook:
//   session.PrintPreCompactionMessage(compactType, compactionCount)
//
// Shared Utility:
//   system := session.GetSystemInfo()
//
// ============================================================================
// END CLOSING
// ============================================================================
