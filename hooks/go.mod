// ============================================================================
// METADATA
// ============================================================================
// Go Module Definition - Hooks package configuration
// Defines module structure for Nova Dawn hook system

// ============================================================================
// SETUP
// ============================================================================
module hooks

go 1.24.4

require (
	system/lib/calendar v0.0.0 // indirect
	system/lib/config v0.0.0 // indirect
	system/lib/display v0.0.0 // indirect
	system/lib/fs v0.0.0 // indirect
	system/lib/git v0.0.0
	system/lib/instance v0.0.0 // indirect
	system/lib/jsonc v0.0.0 // indirect
	system/lib/logging v0.0.0 // indirect
	system/lib/paths v0.0.0 // indirect
	system/lib/planner v0.0.0 // indirect
	system/lib/privacy v0.0.0
	system/lib/sessiontime v0.0.0 // indirect
	system/lib/system v0.0.0 // indirect
	system/lib/temporal v0.0.0
	system/lib/validation v0.0.0
)

require hooks/lib v0.0.0-00010101000000-000000000000

require github.com/BurntSushi/toml v1.5.0 // indirect

// Local module references
replace hooks/lib => ./lib

replace system/lib => ../cpi-si/system/runtime/lib

replace system/lib/calendar => ../cpi-si/system/runtime/lib/calendar

replace system/lib/config => ../cpi-si/system/runtime/lib/config

replace system/lib/display => ../cpi-si/system/runtime/lib/display

replace system/lib/fs => ../cpi-si/system/runtime/lib/fs

replace system/lib/git => ../cpi-si/system/runtime/lib/git

replace system/lib/instance => ../cpi-si/system/runtime/lib/instance

replace system/lib/jsonc => ../cpi-si/system/runtime/lib/jsonc

replace system/lib/logging => ../cpi-si/system/runtime/lib/logging

replace system/lib/paths => ../cpi-si/system/runtime/lib/paths

replace system/lib/planner => ../cpi-si/system/runtime/lib/planner

replace system/lib/privacy => ../cpi-si/system/runtime/lib/privacy

replace system/lib/sessiontime => ../cpi-si/system/runtime/lib/sessiontime

replace system/lib/system => ../cpi-si/system/runtime/lib/system

replace system/lib/temporal => ../cpi-si/system/runtime/lib/temporal

replace system/lib/validation => ../cpi-si/system/runtime/lib/validation

// ============================================================================
// CLOSING
// ============================================================================
// Module is ready for use by all hook files
