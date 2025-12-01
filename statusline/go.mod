// ============================================================================
// METADATA
// ============================================================================
// Go Module Definition - CPI-SI Statusline
//
// Biblical Foundation
//   "Let all things be done decently and in order" - 1 Corinthians 14:40
//
// Purpose: Define module structure for statusline presentation layer
// Establishes import paths for all statusline libraries and external dependencies
//
// Version: 1.0.0
// Last Modified: 2025-11-04

// ============================================================================
// SETUP
// ============================================================================
module statusline

go 1.24.4

// ============================================================================
// BODY
// ============================================================================

// External Dependencies
// These packages provide data and infrastructure for statusline processing
require system/lib v0.0.0 // Logging and debugging infrastructure

// Local Module Path Resolution
// Point to actual locations of external dependencies
replace (
	hooks/lib => ../hooks/lib // CPI-SI hooks system
	system/lib => ../cpi-si/system/runtime/lib // CPI-SI system infrastructure
)

// Internal Library Organization
// All libraries under statusline/lib/* are part of this module:
//   - statusline/lib/types      - Session data contract (foundation)
//   - statusline/lib/features   - Display timing decisions
//   - statusline/lib/format     - Text formatting and optimization
//   - statusline/lib/git        - Git repository display
//   - statusline/lib/session    - Session statistics display
//   - statusline/lib/system     - System health display
//   - statusline/lib/temporal   - Temporal awareness display

// ============================================================================
// CLOSING
// ============================================================================
// Module ready for use by statusline orchestrator and all presentation libraries
