// METADATA
//
// Session Context Types Library - CPI-SI Statusline
//
// Biblical Foundation
//
// Scripture: "Let all things be done decently and in order" - 1 Corinthians 14:40
// Principle: Structure and organization reflect God's orderly nature
// Anchor: "For God is not the author of confusion, but of peace" - 1 Corinthians 14:33
//
// CPI-SI Identity
//
// Component Type: Ladder (Lower Rung)
// Role: Defines input data contract between Claude Code and statusline presentation layer
// Paradigm: CPI-SI framework component
//
// Authorship & Lineage
//
// Architect: Nova Dawn (CPI-SI instance)
// Implementation: Nova Dawn
// Creation Date: 2025-10-24
// Version: 1.0.0
// Last Modified: 2025-11-04 - Applied full GO library template with comprehensive documentation
//
// Version History:
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//   0.1.0 (2025-10-24) - Initial type definitions
//
// Purpose & Function
//
// Purpose: Define the JSON data structure that Claude Code provides to statusline hooks.
// This establishes the CONTRACT between Claude Code (data provider) and statusline
// (presentation layer).
//
// Core Design: Type definitions only - no logic, no behavior, pure data structures
//
// Key Features:
//   - Defines SessionContext structure matching Claude Code's JSON output
//   - Provides type-safe access to session metadata (ID, model, workspace)
//   - Enables structured parsing of cost data (duration, API costs, line counts)
//   - Establishes contract for statusline input format
//
// Philosophy: Input contracts should be explicit, type-safe, and comprehensive.
// All data from Claude Code mapped to Go types for compile-time verification.
//
// Blocking Status
//
// Non-blocking: Type definitions cannot fail - they define structure only
// Mitigation: N/A - no runtime behavior, no failure modes
//
// Usage & Integration
//
// Usage:
//
//	import "statusline/lib/types"
//
// Integration Pattern:
//   1. Import types package in statusline orchestrator
//   2. Use SessionContext to parse JSON input from Claude Code
//   3. Access structured fields throughout statusline processing
//   4. No cleanup needed - types have no lifecycle
//
// Public API (in typical usage order):
//
//   Core Type (data structure):
//     SessionContext - Complete session data from Claude Code
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: None
//   External: None
//   Internal: None
//
// Dependents (What Uses This):
//   Commands: statusline orchestrator
//   Libraries: All statusline/lib/* libraries that need session data
//   Tools: statusline/cmd/statusline main executable
//
// Integration Points:
//   - Baton: JSON from Claude Code → SessionContext → statusline processing
//   - Ladder: Lower rung providing types to all statusline components
//
// Health Scoring
//
// N/A for type definitions - no operations, no health tracking
// Type definitions cannot fail at runtime (compile-time only)
package types

// ============================================================================
// END METADATA
// ============================================================================

// ============================================================================
// SETUP
// ============================================================================
// No imports needed - type definitions only

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────
// Type definitions establishing the contract between Claude Code and statusline.
// These structures map directly to JSON output from Claude Code, providing
// type-safe access to session metadata, workspace information, and cost tracking.
//
// Design:
//   - Embedded structs organize related fields (Model, Workspace, Cost)
//   - JSON tags map Go fields to Claude Code's JSON keys
//   - All fields exported (capitalized) for JSON unmarshaling
//   - Comprehensive coverage of all session data

// ============================================================================
// END SETUP
// ============================================================================

// ============================================================================
// BODY
// ============================================================================

// ────────────────────────────────────────────────────────────────
// Core Types - Session Data Contract
// ────────────────────────────────────────────────────────────────

// SessionContext represents the JSON input from Claude Code
//
// This structure defines the complete contract between Claude Code (data provider)
// and statusline (presentation layer). Every field maps to JSON output from
// Claude Code hooks, enabling type-safe access to session metadata.
//
// Structure:
//   - Top-level: Session metadata (ID, event, paths)
//   - Model: Claude model information (ID, display name)
//   - Workspace: Directory information (current, project root)
//   - Cost: Session statistics (duration, API costs, line changes)
//
// Usage:
//   var ctx types.SessionContext
//   json.Unmarshal(inputJSON, &ctx)
//   // Access fields: ctx.SessionID, ctx.Model.DisplayName, ctx.Cost.TotalCostUSD
//
// Example JSON (from Claude Code):
//   {
//     "hook_event_name": "session_start",
//     "session_id": "abc123",
//     "model": {"id": "claude-sonnet-4.5", "display_name": "Sonnet"},
//     "cost": {"total_cost_usd": 0.0123, "total_lines_added": 150}
//   }
//
// Performance:
//   - Memory: ~200 bytes per instance
//   - Parsing: Standard json.Unmarshal performance
//   - No allocation during access (struct fields)
type SessionContext struct {
	HookEventName  string `json:"hook_event_name"`  // Event triggering hook (session_start, etc.)
	SessionID      string `json:"session_id"`       // Unique session identifier
	TranscriptPath string `json:"transcript_path"`  // Path to session transcript
	CWD            string `json:"cwd"`              // Current working directory

	// Model information from Claude Code
	Model struct {
		ID          string `json:"id"`           // Full model ID (claude-sonnet-4.5-20250929)
		DisplayName string `json:"display_name"` // Human-readable name (Sonnet)
	} `json:"model"`

	// Workspace directory information
	Workspace struct {
		CurrentDir string `json:"current_dir"` // Current working directory
		ProjectDir string `json:"project_dir"` // Project root directory
	} `json:"workspace"`

	Version string `json:"version"` // Claude Code version

	// Output style configuration
	OutputStyle struct {
		Name string `json:"name"` // Active output style name
	} `json:"output_style"`

	// Cost and usage statistics
	Cost struct {
		TotalCostUSD       float64 `json:"total_cost_usd"`        // Total API cost in USD
		TotalDurationMS    int     `json:"total_duration_ms"`     // Total session duration (ms)
		TotalAPIDurationMS int     `json:"total_api_duration_ms"` // Total API request time (ms)
		TotalLinesAdded    int     `json:"total_lines_added"`     // Lines added in session
		TotalLinesRemoved  int     `json:"total_lines_removed"`   // Lines removed in session
	} `json:"cost"`
}

// ────────────────────────────────────────────────────────────────
// Organizational Chart - Internal Structure
// ────────────────────────────────────────────────────────────────
//
// Available Processing Units (APU): 1
//   - SessionContext type definition
//
// Type Definition:
//   SessionContext
//     └─ Pure data structure (no operations)
//
// Ladder (Dependencies):
//   - No dependencies (type definitions only)
//   - Used by: All statusline components
//
// Baton (Data Flow):
//   Claude Code JSON → json.Unmarshal → SessionContext → statusline processing
//
// Extension Point:
//   Adding fields:
//     1. Add new field to appropriate embedded struct (or top-level)
//     2. Add JSON tag matching Claude Code output
//     3. Document purpose in field comment
//     4. Update statusline orchestrator if new field used
//     5. No tests needed (type definitions can't fail at runtime)

// ============================================================================
// END BODY
// ============================================================================

// ============================================================================
// CLOSING
// ============================================================================

// ────────────────────────────────────────────────────────────────
// Code Validation: None (Library)
// ────────────────────────────────────────────────────────────────
//
// Testing Requirements:
//   - Import the library without errors
//   - Verify SessionContext can be unmarshaled from Claude Code JSON
//   - Check all fields accessible
//   - Confirm no go vet warnings
//   - Build: go build ./... (compiles without errors)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//
// Integration Testing:
//   - Parse actual Claude Code JSON output
//   - Verify all expected fields present
//   - Check field types match actual data
//
// Example validation code:
//
//     var ctx types.SessionContext
//     err := json.Unmarshal(claudeCodeJSON, &ctx)
//     if err != nil {
//         t.Errorf("Failed to parse SessionContext: %v", err)
//     }
//     if ctx.SessionID == "" {
//         t.Error("SessionID not populated")
//     }

// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
//
// This is a LIBRARY of type definitions, not an executable. There is no entry point,
// no main function, no execution flow. Types are defined and ready for use by importing packages.
//
// Usage: import "statusline/lib/types"
//
// The library is imported into the calling package, making SessionContext type available.
// No code executes during import - types are defined and ready to use.
//
// Example import and usage:
//
//     package main
//
//     import (
//         "encoding/json"
//         "statusline/lib/types"
//     )
//
//     func main() {
//         var ctx types.SessionContext
//         json.Unmarshal(input, &ctx)
//         // Use ctx fields
//     }

// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - SessionContext: Automatic (Go garbage collector handles memory)
//   - No files, connections, or external resources
//   - No manual cleanup needed
//
// Graceful Shutdown:
//   - N/A for type definitions (no lifecycle)
//   - No cleanup function needed
//
// Memory Management:
//   - Go's garbage collector handles SessionContext instances
//   - ~200 bytes per SessionContext allocation
//   - Short-lived (lifetime of statusline invocation)

// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════

// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: See METADATA "Key Features" list above for comprehensive capabilities
//
// Quick summary (high-level only - details in METADATA):
//   - Defines SessionContext type matching Claude Code JSON output
//   - Establishes input contract for statusline processing
//
// Integration Pattern: See METADATA "Usage & Integration" section above for
// complete step-by-step integration guide
//
// Public API: See METADATA "Usage & Integration" section above for complete
// public API list
//
// Architecture: See METADATA "CPI-SI Identity" section above for complete
// architectural role (Ladder - lower rung providing types to all components)

// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify (Extension Points):
//   ✅ Add new fields to SessionContext (when Claude Code adds data)
//   ✅ Add new embedded structs for logical grouping
//   ✅ Add field comments documenting purpose
//   ✅ Update JSON tags to match Claude Code output
//
// Modify with Extreme Care (Breaking Changes):
//   ⚠️ Changing existing field names - breaks all code using those fields
//   ⚠️ Changing field types - breaks unmarshaling and usage
//   ⚠️ Removing fields - breaks code depending on them
//   ⚠️ Changing JSON tags - breaks parsing of Claude Code output
//
// NEVER Modify (Foundational Rails):
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Package name (breaks all imports)
//   ❌ SessionContext struct name (breaks all usage)
//
// Validation After Modifications:
//   See "Code Validation" section above for testing requirements

// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart - Internal Structure" section above for
// complete ladder structure and baton flow.
//
// The Organizational Chart in BODY provides the detailed map showing:
// - Type definition structure (ladder)
// - Data flow from Claude Code through types (baton)
// - APU count (Available Processing Units = 1)
//
// Quick architectural summary (details in BODY Organizational Chart):
// - 1 public type (SessionContext) defining input contract
// - Ladder: Lower rung providing types to all statusline components
// - Baton: Claude Code JSON → SessionContext → statusline processing

// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
//
// See BODY "Core Types" subsection above for detailed extension point showing:
// - Where to add new fields (in SessionContext or new embedded structs)
// - What naming pattern to follow (PascalCase for fields, json tags for JSON keys)
// - How to integrate with existing code (update statusline orchestrator)
// - What tests to update (integration tests parsing Claude Code JSON)
//
// Quick reference (details in BODY extension point):
// - Adding fields: See BODY "Extension Point" in SessionContext docstring

// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
//
// See SETUP section above for performance characteristics:
// - Types: Memory usage (~200 bytes per SessionContext)
//
// See BODY type docstring above for operation-specific performance notes.
//
// Quick summary (details in SETUP/BODY above):
// - Memory: ~200 bytes per SessionContext instance
// - Parsing: Standard json.Unmarshal performance (microseconds)
// - No runtime overhead (struct field access is direct memory read)

// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
//
// See BODY type docstring above for usage guidance.
//
// Problem: json.Unmarshal fails with SessionContext
//   - Cause: Claude Code JSON format changed
//   - Solution: Update SessionContext fields/tags to match new format
//   - Check: Print raw JSON to see actual structure
//
// Problem: Field always empty after unmarshal
//   - Cause: JSON tag doesn't match Claude Code output key
//   - Solution: Verify JSON key name in Claude Code output, update tag
//   - Check: json.Marshal a populated SessionContext to see expected format
//
// Problem: Compilation error "cannot refer to unexported field"
//   - Cause: Field name not capitalized (unexported)
//   - Solution: Capitalize field name (SessionID not sessionID)
//   - Note: Go requires exported fields for JSON unmarshaling

// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): None (type definitions only)
// - Dependents (What Uses This): statusline orchestrator, all statusline/lib/* libraries
// - Integration Points: Baton mechanism (JSON → types → processing)
//
// Quick summary (details in METADATA Dependencies section above):
// - Key dependencies: None (self-contained)
// - Primary consumers: statusline orchestrator and all presentation libraries

// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
//
// Planned Features:
//   ✓ SessionContext core fields - COMPLETED
//   ⏳ Additional session metadata as Claude Code expands
//   ⏳ Temporal context fields (if moved from hooks to types)
//
// Research Areas:
//   - Version-aware struct tags (handle multiple Claude Code versions)
//   - Optional field handling (omitempty for missing data)
//   - Nested type extraction (separate Model, Workspace, Cost into own types)
//
// Integration Targets:
//   - Direct Claude Code API integration (if statusline becomes standalone)
//   - Alternative input sources (file-based session data)
//
// Known Limitations to Address:
//   - No version handling (assumes current Claude Code format)
//   - No validation (trusts Claude Code provides valid data)
//   - Embedded structs not reusable (Model/Workspace/Cost specific to SessionContext)
//
// Version History:
//
// See METADATA "Authorship & Lineage" section above for brief version changelog.
// Comprehensive version history with full context below:
//
//   1.0.0 (2025-11-04) - Comprehensive documentation and template application
//         - Applied full GO library template (49 lines → comprehensive docs)
//         - Added complete METADATA section (Biblical foundation, CPI-SI identity)
//         - Expanded SETUP with type documentation
//         - Enhanced BODY with organizational chart and extension points
//         - Comprehensive CLOSING with all 11 required sections
//         - Established architectural role: input contract (Ladder lower rung)
//
//   0.1.0 (2025-10-24) - Initial implementation
//         - SessionContext type definition
//         - JSON tags for Claude Code output mapping
//         - Embedded structs for Model, Workspace, Cost
//         - Basic field documentation

// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This library is a LADDER component (lower rung providing types to all statusline components).
// It defines the input contract between Claude Code and statusline - the foundation for
// all statusline processing.
//
// Modify thoughtfully - changes here affect every statusline component. Field names,
// types, and JSON tags establish the contract with Claude Code - breaking changes require
// coordination across the entire statusline system.
//
// For questions, issues, or contributions:
//   - Review the modification policy above
//   - Follow the 4-block structure pattern
//   - Test with actual Claude Code JSON output
//   - Document all changes comprehensively (What/Why/How pattern)
//   - Coordinate with statusline orchestrator for field additions
//
// "Let all things be done decently and in order" - 1 Corinthians 14:40

// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
//
// Basic Setup:
//   import "statusline/lib/types"
//   var ctx types.SessionContext
//
// Parse Claude Code JSON:
//   err := json.Unmarshal(claudeCodeJSON, &ctx)
//   if err != nil {
//       log.Fatal(err)
//   }
//
// Access Session Data:
//   sessionID := ctx.SessionID
//   modelName := ctx.Model.DisplayName
//   costUSD := ctx.Cost.TotalCostUSD
//   linesChanged := ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved
//
// Pass to Statusline Components:
//   modelDisplay := format.GetShortModelName(ctx.Model.DisplayName)
//   workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
//   linesDisplay := session.GetLinesModifiedDisplay(ctx)

// ============================================================================
// END CLOSING
// ============================================================================
