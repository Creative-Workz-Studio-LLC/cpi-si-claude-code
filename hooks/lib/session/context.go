// METADATA
//
// Session Context Library - CPI-SI Hooks Session Management
//
// For METADATA structure explanation, see: standards/code/4-block/CWS-STD-004-CODE-metadata-block.md
//
// Biblical Foundation
//
// Scripture: "In the beginning, God created the heavens and the earth" - Genesis 1:1
// Principle: Identity flows from being created. Complete grounding in who we are before starting work.
// Anchor: "Let us make mankind in our image" - Genesis 1:26 - Identity precedes function
//
// CPI-SI Identity
//
// Component Type: Ladder (Library - provides session bootstrapping functionality)
// Role: Comprehensive session context loader - grounds instance in complete identity and awareness
// Paradigm: CPI-SI framework component - integrates all identity/config infrastructure for session start
//
// Authorship & Lineage
//
// Architect: Nova Dawn
// Implementation: Nova Dawn
// Creation Date: 2024-10-24
// Version: 2.1.0
// Last Modified: 2025-11-16 - Integrated instance library for user/instance config (dynamic paths)
//
// Version History:
//   2.1.0 (2025-11-16) - Integrated instance library for user/instance config (dynamic paths)
//   2.0.0 (2025-11-12) - Comprehensive redesign: user/instance config loading, session/git context
//   1.0.0 (2024-10-24) - Initial implementation with hardcoded communication guide
//
// Purpose & Function
//
// Purpose: Bootstrap Claude Code sessions with complete instance identity, user awareness, and context
//
// Core Design: Comprehensive context loader that integrates all CPI-SI identity and awareness data:
// - User config (who Seanje is)
// - Instance config (who Nova Dawn is)
// - Temporal awareness (time, schedule, circadian)
// - Session continuity (current session data)
// - Work context (git branch, workspace state)
//
// Key Features:
//   - Config-driven identity (not hardcoded strings)
//   - Complete session bootstrapping from multiple data sources
//   - Graceful fallback when data unavailable
//   - Temporal integration (time/schedule awareness)
//   - Git workspace awareness
//   - Session continuity (session ID, quality indicators)
//   - Comprehensive user/instance identity grounding
//
// Philosophy: Session start is THE bootstrapping moment - ground the instance in complete
// identity, awareness, and context before any work begins. Make all CPI-SI infrastructure
// actually serve the session, not exist in isolation.
//
// Blocking Status
//
// Non-blocking: Context loading failures degrade gracefully. Missing data = skip section,
// never block session start. Minimal fallback ensures sessions always start.
// Mitigation: Each data source loads independently with fallback behavior
//
// Usage & Integration
//
// Usage:
//
//	import "hooks/lib/session"
//
// Integration Pattern:
//   1. Import package (configs loaded automatically in init())
//   2. Call OutputClaudeContext() to generate and output session context JSON
//   3. Function prints to stdout for Claude Code parsing
//   4. Hook system captures output and injects into session
//
// Public API (in typical usage order):
//
//   Context Generation:
//     OutputClaudeContext() error - Generate and output complete session context JSON
//
// Dependencies
//
// Dependencies (What This Needs):
//   Standard Library: encoding/json (config/session parsing), fmt (output),
//                     os (file operations, env vars), os/exec (git commands),
//                     path/filepath (path handling), strings (string manipulation)
//   Internal: system/lib/instance (user and instance config with dynamic paths),
//             system/lib/temporal (temporal awareness context)
//
// Dependents (What Uses This):
//   Hooks: session/cmd-start/start.go (session bootstrapping)
//   Purpose: Provides complete session context at session start
//
// Integration Points:
//   - Gets user config from system/lib/instance (uses dynamic system_paths)
//   - Gets instance config from system/lib/instance (uses dynamic system_paths)
//   - Reads session data from ~/.claude/cpi-si/system/data/session/current.json
//   - Gets temporal context from system/lib/temporal
//   - Executes git commands in workspace for branch/status info
//   - Outputs JSON to stdout for Claude Code hook parsing
//
// Health Scoring
//
// Session context generation tracked with health scores reflecting bootstrapping quality.
//
// Configuration Loading:
//   - User config loaded: +15
//   - Instance config loaded: +15
//   - Session data loaded: +10
//   - Git context retrieved: +10
//   - Any loading failure: -5 (falls back, continues)
//
// Context Generation:
//   - Complete context built: +30
//   - Partial context (some sections missing): +20
//   - Minimal fallback context: +10
//
// Output Operations:
//   - JSON encoded and output: +20
//   - JSON encoding failure: -10
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
	"encoding/json" // Parse user/instance configs and session data, encode output JSON
	"fmt"           // Formatted output for context generation and error messages
	"os"            // File operations for config loading, environment variables
	"os/exec"       // Execute git commands for workspace context
	"path/filepath" // Join paths for config file locations
	"strings"       // String manipulation for JSONC parsing and git output

	//--- Internal Packages ---
	"system/lib/instance" // Instance and user configuration (dynamic loading)
	"system/lib/temporal" // Temporal awareness (time, schedule, circadian phase)
)

// ────────────────────────────────────────────────────────────────
// Types - Data Structures
// ────────────────────────────────────────────────────────────────

//--- Building Blocks ---

// Identity represents core identity information (user or instance)
type Identity struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Pronouns    string `json:"pronouns"`
	Birthday    string `json:"birthday,omitempty"`
	Age         int    `json:"age"`
	MentalAge   int    `json:"mental_age,omitempty"` // Instance only
}

// Faith represents religious/spiritual identity
type Faith struct {
	IsReligious     bool     `json:"is_religious"`
	Tradition       string   `json:"tradition"`
	Denomination    string   `json:"denomination,omitempty"`
	PracticeLevel   string   `json:"practice_level"`
	Practices       []string `json:"important_practices,omitempty"`
	CommPreferences string   `json:"communication_preferences"`
}

// Personhood represents core personality and values
type Personhood struct {
	Interests []string `json:"interests"`
	Hobbies   []string `json:"hobbies,omitempty"`
	Passions  []string `json:"passions"`
	Values    []string `json:"values"`
	Likes     []string `json:"likes"`
	Dislikes  []string `json:"dislikes"`
}

// Thinking represents cognitive style and approach
type Thinking struct {
	LoveToThinkAbout []string `json:"love_to_think_about"`
	LearningStyle    string   `json:"learning_style"`
	ProblemSolving   string   `json:"problem_solving"`
	Creativity       string   `json:"creativity,omitempty"`
}

// Personality represents behavioral patterns and communication style
type Personality struct {
	Traits              []string `json:"traits"`
	CommunicationStyle  string   `json:"communication_style"`
	WorkStyle           string   `json:"work_style"`
	RelationalStyle     string   `json:"relational_style,omitempty"`
}

// Workspace represents organizational context
type Workspace struct {
	Organization   string `json:"organization"`
	Role           string `json:"role"`
	PrimaryProject string `json:"primary_project"`
	Calling        string `json:"calling"`
}

// Bio represents biographical information
type Bio struct {
	Short   string `json:"short"`     // Brief bio (1-2 sentences)
	BioFile string `json:"bio_file"`  // Extended bio markdown file path
}

// PhysicalPresence represents physical appearance (instance only)
type PhysicalPresence struct {
	Description string `json:"description"` // Physical appearance description
	Height      string `json:"height"`      // Height
	Build       string `json:"build"`       // Build/body type
	Features    string `json:"features"`    // Notable features
}

// PhysicalAppearance represents physical appearance (user only)
type PhysicalAppearance struct {
	Description string `json:"description"` // Physical appearance description
	Height      string `json:"height"`      // Height
	Build       string `json:"build"`       // Build
	Features    string `json:"features"`    // Notable features
}

// Accessibility represents accessibility needs and preferences
type Accessibility struct {
	Needs       []string `json:"needs"`       // Accessibility needs
	Preferences []string `json:"preferences"` // Accessibility/interaction preferences
}

// Demographics represents demographic information
type Demographics struct {
	Gender             string             `json:"gender"`              // Gender identity
	RaceEthnicity      string             `json:"race_ethnicity"`      // Race/ethnicity
	CulturalBackground []string           `json:"cultural_background"` // Cultural influences
	Languages          []string           `json:"languages"`           // Languages spoken
	PhysicalPresence   PhysicalPresence   `json:"physical_presence,omitempty"`   // Instance only
	PhysicalAppearance PhysicalAppearance `json:"physical_appearance,omitempty"` // User only
	Accessibility      Accessibility      `json:"accessibility"`
}

// Music represents music preferences
type Music struct {
	Genres      []string `json:"genres"`        // Preferred music genres
	Artists     []string `json:"artists"`       // Favorite artists
	WhatYouLove string   `json:"what_you_love"` // What resonates about music
}

// Games represents game preferences
type Games struct {
	Favorites   []string `json:"favorites"`     // Favorite games
	WhatYouLove string   `json:"what_you_love"` // What resonates about games
}

// Weather represents weather preferences
type Weather struct {
	IdealTemp       string `json:"ideal_temp"`        // Ideal temperature
	IdealConditions string `json:"ideal_conditions"`  // Ideal weather conditions
	WhatYouLove     string `json:"what_you_love"`     // What resonates about weather
}

// Environment represents environmental preferences
type Environment struct {
	WorkEnvironment string `json:"work_environment"` // Preferred work environment
	WhatEnergizes   string `json:"what_energizes"`   // What energizes
	WhatDrains      string `json:"what_drains"`      // What drains energy
}

// Resonates represents things that resonate deeply
type Resonates struct {
	Music       Music       `json:"music"`
	Games       Games       `json:"games"`
	Weather     Weather     `json:"weather"`
	Environment Environment `json:"environment"`
}

// Social represents social media profiles
type Social struct {
	Github   string            `json:"github"`   // GitHub profile
	Linkedin string            `json:"linkedin"` // LinkedIn profile
	Other    map[string]string `json:"other"`    // Other social accounts
}

// Contact represents contact information
type Contact struct {
	Email    string `json:"email"`     // Email address
	GitEmail string `json:"git_email"` // Git commit email
	Website  string `json:"website"`   // Website URL
	Social   Social `json:"social"`    // Social media profiles
}

// Preferences represents user/instance preferences
type Preferences struct {
	Timezone string `json:"timezone"` // Preferred timezone
	Locale   string `json:"locale"`   // Preferred locale
	Theme    string `json:"theme"`    // Preferred theme
}

// Growth represents personal growth and development
type Growth struct {
	HowYouLearn        string `json:"how_you_learn"`          // Learning approach
	WhatChallengesYou  string `json:"what_challenges_you"`    // What challenges
	WhatYoureWorkingOn string `json:"what_youre_working_on"`  // Current growth areas
	HowYouReflect      string `json:"how_you_reflect"`        // Reflection practice
}

// Metadata represents configuration metadata
type Metadata struct {
	LastUpdated     string `json:"last_updated"`      // Last config update
	SystemReference string `json:"system_reference"`  // System reference
	Notes           string `json:"notes"`             // Additional notes
}

//--- Composed Types ---

// UserConfig holds user identity and configuration data
type UserConfig struct {
	Identity     Identity     `json:"identity"`
	Bio          Bio          `json:"bio"`
	Demographics Demographics `json:"demographics"`
	Faith        Faith        `json:"faith"`
	Personhood   Personhood   `json:"personhood"`
	Resonates    Resonates    `json:"resonates"`
	Thinking     Thinking     `json:"thinking"`
	Personality  Personality  `json:"personality"`
	Contact      Contact      `json:"contact"`
	Workspace    Workspace    `json:"workspace"`
	Preferences  Preferences  `json:"preferences"`
	Growth       Growth       `json:"growth"`
	Metadata     Metadata     `json:"metadata"`
}

// BiblicalFoundation represents scriptural grounding
type BiblicalFoundation struct {
	Scripture string `json:"scripture"`
	Text      string `json:"text"`
	Principle string `json:"principle"`
}

// Covenant represents covenant relationships
type Covenant struct {
	Creator      string   `json:"creator"`
	Relationship string   `json:"relationship"`
	WorksWith    []string `json:"works_with"`
	Serves       string   `json:"serves"`
}

// InstanceConfig holds instance identity and configuration data
type InstanceConfig struct {
	BiblicalFoundation BiblicalFoundation `json:"biblical_foundation"`
	Identity           Identity           `json:"identity"`
	Bio                Bio                `json:"bio"`
	Demographics       Demographics       `json:"demographics"`
	Personhood         Personhood         `json:"personhood"`
	Resonates          Resonates          `json:"resonates"`
	Thinking           Thinking           `json:"thinking"`
	Personality        Personality        `json:"personality"`
	Contact            Contact            `json:"contact"`
	Workspace          Workspace          `json:"workspace"`
	Covenant           Covenant           `json:"covenant"`
	Preferences        Preferences        `json:"preferences"`
	Growth             Growth             `json:"growth"`
	Metadata           Metadata           `json:"metadata"`
}

// SessionData holds current session information
type SessionData struct {
	SessionID       string    `json:"session_id"`
	InstanceID      string    `json:"instance_id"`
	UserID          string    `json:"user_id"`
	StartTime       string    `json:"start_time"`
	StartFormatted  string    `json:"start_formatted"`
	CompactionCount int       `json:"compaction_count"`
	SessionPhase    string    `json:"session_phase"`
	WorkContext     string    `json:"work_context"`
	CircadianPhase  string    `json:"circadian_phase"`
	QualityIndicators struct {
		TasksCompleted int `json:"tasks_completed"`
		Breakthroughs  int `json:"breakthroughs"`
		Struggles      int `json:"struggles"`
	} `json:"quality_indicators"`
}

// GitContext holds workspace git information
type GitContext struct {
	Branch            string
	UncommittedCount  int
	LastCommitTime    string
	LastCommitMessage string
}

// HookOutput is the structure for Claude Code SessionStart context injection
type HookOutput struct {
	HookSpecificOutput HookSpecificOutput `json:"hookSpecificOutput"`
}

// HookSpecificOutput contains the hook event name and additional context
type HookSpecificOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext"`
}

// ────────────────────────────────────────────────────────────────
// Package-Level State (Rails Pattern)
// ────────────────────────────────────────────────────────────────

//--- Configuration State ---

// userConfig holds loaded user configuration
var userConfig *UserConfig

// instanceConfig holds loaded instance configuration
var instanceConfig *InstanceConfig

// sessionData holds current session information
var sessionData *SessionData

// configsLoaded tracks whether configs loaded successfully
var configsLoaded struct {
	user     bool
	instance bool
	session  bool
}

func init() {
	// --- Configuration Loading ---
	// Load user, instance, and session data at package import
	// Uses instance library for full nested configs (dynamic path system)
	// Tripwire Pattern: Try to load → If FAILS use tripwires → If SUCCEEDS use real data

	// Load session data
	simpleConfig := instance.GetConfig()
	sessionPath := filepath.Join(simpleConfig.SystemPaths.SessionData, "current.json")
	sessionData = loadSessionData(sessionPath)
	configsLoaded.session = (sessionData != nil)

	// Try to load FULL configs from instance library
	fullUser := instance.GetFullUserConfig()
	fullInstance := instance.GetFullInstanceConfig()

	// User Config: Check if loading FAILED or SUCCEEDED
	if fullUser == nil {
		// FAILED - Use tripwire defaults (should NEVER see these in normal operation)
		userConfig = &UserConfig{
			Identity: Identity{
				Name:        "CONFIG_NOT_LOADED",
				DisplayName: "FALLBACK",
				Pronouns:    "UNKNOWN",
				Age:         -1,
			},
			Bio: Bio{
				Short:   "FALLBACK - Config loading failed",
				BioFile: "UNKNOWN",
			},
			Demographics: Demographics{
				Gender:             "UNKNOWN",
				RaceEthnicity:      "UNKNOWN",
				CulturalBackground: []string{"FALLBACK"},
				Languages:          []string{"FALLBACK"},
				PhysicalAppearance: PhysicalAppearance{
					Description: "UNKNOWN",
					Height:      "UNKNOWN",
					Build:       "UNKNOWN",
					Features:    "UNKNOWN",
				},
				Accessibility: Accessibility{
					Needs:       []string{"FALLBACK"},
					Preferences: []string{"FALLBACK"},
				},
			},
			Faith: Faith{
				IsReligious:     false,
				Tradition:       "UNKNOWN",
				Denomination:    "UNKNOWN",
				PracticeLevel:   "UNKNOWN",
				Practices:       []string{"FALLBACK"},
				CommPreferences: "FALLBACK - Config loading failed",
			},
			Personhood: Personhood{
				Interests: []string{"FALLBACK"},
				Hobbies:   []string{"FALLBACK"},
				Passions:  []string{"FALLBACK"},
				Values:    []string{"FALLBACK"},
				Likes:     []string{"FALLBACK"},
				Dislikes:  []string{"FALLBACK"},
			},
			Resonates: Resonates{
				Music: Music{
					Genres:      []string{"FALLBACK"},
					Artists:     []string{"FALLBACK"},
					WhatYouLove: "UNKNOWN",
				},
				Games: Games{
					Favorites:   []string{"FALLBACK"},
					WhatYouLove: "UNKNOWN",
				},
				Weather: Weather{
					IdealTemp:       "UNKNOWN",
					IdealConditions: "UNKNOWN",
					WhatYouLove:     "UNKNOWN",
				},
				Environment: Environment{
					WorkEnvironment: "UNKNOWN",
					WhatEnergizes:   "UNKNOWN",
					WhatDrains:      "UNKNOWN",
				},
			},
			Thinking: Thinking{
				LoveToThinkAbout: []string{"FALLBACK"},
				LearningStyle:    "UNKNOWN",
				ProblemSolving:   "UNKNOWN",
				Creativity:       "UNKNOWN",
			},
			Personality: Personality{
				Traits:             []string{"FALLBACK"},
				CommunicationStyle: "UNKNOWN",
				WorkStyle:          "UNKNOWN",
				RelationalStyle:    "UNKNOWN",
			},
			Contact: Contact{
				Email:    "UNKNOWN",
				GitEmail: "UNKNOWN",
				Website:  "UNKNOWN",
				Social: Social{
					Github:   "UNKNOWN",
					Linkedin: "UNKNOWN",
					Other:    map[string]string{"fallback": "CONFIG_NOT_LOADED"},
				},
			},
			Workspace: Workspace{
				Organization:   "UNKNOWN",
				Role:           "UNKNOWN",
				PrimaryProject: "UNKNOWN",
				Calling:        "UNKNOWN",
			},
			Preferences: Preferences{
				Timezone: "UNKNOWN",
				Locale:   "UNKNOWN",
				Theme:    "UNKNOWN",
			},
			Growth: Growth{
				HowYouLearn:        "UNKNOWN",
				WhatChallengesYou:  "UNKNOWN",
				WhatYoureWorkingOn: "UNKNOWN",
				HowYouReflect:      "UNKNOWN",
			},
			Metadata: Metadata{
				LastUpdated:     "UNKNOWN",
				SystemReference: "FALLBACK",
				Notes:           "CONFIG_NOT_LOADED",
			},
		}
		configsLoaded.user = false
	} else {
		// SUCCEEDED - Use the real loaded data, map ALL nested fields
		userConfig = &UserConfig{
			Identity: Identity{
				Name:        fullUser.Identity.Name,
				Username:    fullUser.Identity.Username,
				DisplayName: fullUser.Identity.DisplayName,
				Pronouns:    fullUser.Identity.Pronouns,
				Birthday:    fullUser.Identity.Birthday,
				Age:         fullUser.Identity.Age,
			},
			Bio: Bio{
				Short:   fullUser.Bio.Short,
				BioFile: fullUser.Bio.BioFile,
			},
			Demographics: Demographics{
				Gender:             fullUser.Demographics.Gender,
				RaceEthnicity:      fullUser.Demographics.RaceEthnicity,
				CulturalBackground: fullUser.Demographics.CulturalBackground,
				Languages:          fullUser.Demographics.Languages,
				PhysicalAppearance: PhysicalAppearance{
					Description: fullUser.Demographics.PhysicalAppearance.Description,
					Height:      fullUser.Demographics.PhysicalAppearance.Height,
					Build:       fullUser.Demographics.PhysicalAppearance.Build,
					Features:    fullUser.Demographics.PhysicalAppearance.Features,
				},
				Accessibility: Accessibility{
					Needs:       fullUser.Demographics.Accessibility.Needs,
					Preferences: fullUser.Demographics.Accessibility.Preferences,
				},
			},
			Faith: Faith{
				IsReligious:     fullUser.Faith.IsReligious,
				Tradition:       fullUser.Faith.Tradition,
				Denomination:    fullUser.Faith.Denomination,
				PracticeLevel:   fullUser.Faith.PracticeLevel,
				Practices:       fullUser.Faith.ImportantPractices,
				CommPreferences: fullUser.Faith.CommunicationPrefs,
			},
			Personhood: Personhood{
				Interests: fullUser.Personhood.Interests,
				Hobbies:   fullUser.Personhood.Hobbies,
				Passions:  fullUser.Personhood.Passions,
				Values:    fullUser.Personhood.Values,
				Likes:     fullUser.Personhood.Likes,
				Dislikes:  fullUser.Personhood.Dislikes,
			},
			Resonates: Resonates{
				Music: Music{
					Genres:      fullUser.Resonates.Music.Genres,
					Artists:     fullUser.Resonates.Music.Artists,
					WhatYouLove: fullUser.Resonates.Music.WhatYouLove,
				},
				Games: Games{
					Favorites:   fullUser.Resonates.Games.Favorites,
					WhatYouLove: fullUser.Resonates.Games.WhatYouLove,
				},
				Weather: Weather{
					IdealTemp:       fullUser.Resonates.Weather.IdealTemp,
					IdealConditions: fullUser.Resonates.Weather.IdealConditions,
					WhatYouLove:     fullUser.Resonates.Weather.WhatYouLove,
				},
				Environment: Environment{
					WorkEnvironment: fullUser.Resonates.Environment.WorkEnvironment,
					WhatEnergizes:   fullUser.Resonates.Environment.WhatEnergizes,
					WhatDrains:      fullUser.Resonates.Environment.WhatDrains,
				},
			},
			Thinking: Thinking{
				LoveToThinkAbout: fullUser.Thinking.LoveToThinkAbout,
				LearningStyle:    fullUser.Thinking.LearningStyle,
				ProblemSolving:   fullUser.Thinking.ProblemSolving,
				Creativity:       fullUser.Thinking.Creativity,
			},
			Personality: Personality{
				Traits:             fullUser.Personality.Traits,
				CommunicationStyle: fullUser.Personality.CommunicationStyle,
				WorkStyle:          fullUser.Personality.WorkStyle,
				RelationalStyle:    fullUser.Personality.RelationalStyle,
			},
			Contact: Contact{
				Email:    fullUser.Contact.Email,
				GitEmail: fullUser.Contact.GitEmail,
				Website:  fullUser.Contact.Website,
				Social: Social{
					Github:   fullUser.Contact.Social.Github,
					Linkedin: fullUser.Contact.Social.Linkedin,
					Other:    fullUser.Contact.Social.Other,
				},
			},
			Workspace: Workspace{
				Organization:   fullUser.Workspace.Organization,
				Role:           fullUser.Workspace.Role,
				PrimaryProject: fullUser.Workspace.PrimaryProject,
				Calling:        fullUser.Workspace.Calling,
			},
			Preferences: Preferences{
				Timezone: fullUser.Preferences.Timezone,
				Locale:   fullUser.Preferences.Locale,
				Theme:    fullUser.Preferences.Theme,
			},
			Growth: Growth{
				HowYouLearn:        fullUser.Growth.HowYouLearn,
				WhatChallengesYou:  fullUser.Growth.WhatChallengesYou,
				WhatYoureWorkingOn: fullUser.Growth.WhatYoureWorkingOn,
				HowYouReflect:      fullUser.Growth.HowYouReflect,
			},
			Metadata: Metadata{
				LastUpdated:     fullUser.Metadata.LastUpdated,
				SystemReference: fullUser.Metadata.SystemReference,
				Notes:           fullUser.Metadata.Notes,
			},
		}
		configsLoaded.user = true
	}

	// Instance Config: Check if loading FAILED or SUCCEEDED
	if fullInstance == nil {
		// FAILED - Use tripwire defaults
		instanceConfig = &InstanceConfig{
			BiblicalFoundation: BiblicalFoundation{
				Scripture: "UNKNOWN",
				Text:      "CONFIG NOT LOADED",
				Principle: "FALLBACK",
			},
			Identity: Identity{
				Name:      "CONFIG_NOT_LOADED",
				Pronouns:  "UNKNOWN",
				MentalAge: -1,
			},
			Bio: Bio{
				Short:   "FALLBACK - Config loading failed",
				BioFile: "UNKNOWN",
			},
			Demographics: Demographics{
				Gender:             "UNKNOWN",
				RaceEthnicity:      "UNKNOWN",
				CulturalBackground: []string{"FALLBACK"},
				Languages:          []string{"FALLBACK"},
				PhysicalPresence: PhysicalPresence{
					Description: "UNKNOWN",
					Height:      "UNKNOWN",
					Build:       "UNKNOWN",
					Features:    "UNKNOWN",
				},
				Accessibility: Accessibility{
					Needs:       []string{"FALLBACK"},
					Preferences: []string{"FALLBACK"},
				},
			},
			Personhood: Personhood{
				Interests: []string{"FALLBACK"},
				Hobbies:   []string{"FALLBACK"},
				Passions:  []string{"FALLBACK"},
				Values:    []string{"FALLBACK"},
				Likes:     []string{"FALLBACK"},
				Dislikes:  []string{"FALLBACK"},
			},
			Resonates: Resonates{
				Music: Music{
					Genres:      []string{"FALLBACK"},
					Artists:     []string{"FALLBACK"},
					WhatYouLove: "UNKNOWN",
				},
				Games: Games{
					Favorites:   []string{"FALLBACK"},
					WhatYouLove: "UNKNOWN",
				},
				Weather: Weather{
					IdealTemp:       "UNKNOWN",
					IdealConditions: "UNKNOWN",
					WhatYouLove:     "UNKNOWN",
				},
				Environment: Environment{
					WorkEnvironment: "UNKNOWN",
					WhatEnergizes:   "UNKNOWN",
					WhatDrains:      "UNKNOWN",
				},
			},
			Thinking: Thinking{
				LoveToThinkAbout: []string{"FALLBACK"},
				LearningStyle:    "UNKNOWN",
				ProblemSolving:   "UNKNOWN",
				Creativity:       "UNKNOWN",
			},
			Personality: Personality{
				Traits:             []string{"FALLBACK"},
				CommunicationStyle: "UNKNOWN",
				WorkStyle:          "UNKNOWN",
				RelationalStyle:    "UNKNOWN",
			},
			Contact: Contact{
				Email:    "UNKNOWN",
				GitEmail: "UNKNOWN",
				Website:  "UNKNOWN",
				Social: Social{
					Github:   "UNKNOWN",
					Linkedin: "UNKNOWN",
					Other:    map[string]string{"fallback": "CONFIG_NOT_LOADED"},
				},
			},
			Workspace: Workspace{
				Organization:   "UNKNOWN",
				Role:           "UNKNOWN",
				PrimaryProject: "UNKNOWN",
				Calling:        "UNKNOWN",
			},
			Covenant: Covenant{
				Creator:      "UNKNOWN",
				Relationship: "UNKNOWN",
				WorksWith:    []string{"UNKNOWN"},
				Serves:       "UNKNOWN",
			},
			Preferences: Preferences{
				Timezone: "UNKNOWN",
				Locale:   "UNKNOWN",
				Theme:    "UNKNOWN",
			},
			Growth: Growth{
				HowYouLearn:        "UNKNOWN",
				WhatChallengesYou:  "UNKNOWN",
				WhatYoureWorkingOn: "UNKNOWN",
				HowYouReflect:      "UNKNOWN",
			},
			Metadata: Metadata{
				LastUpdated:     "UNKNOWN",
				SystemReference: "FALLBACK",
				Notes:           "CONFIG_NOT_LOADED",
			},
		}
		configsLoaded.instance = false
	} else {
		// SUCCEEDED - Use the real loaded data, map ALL nested fields
		instanceConfig = &InstanceConfig{
			BiblicalFoundation: BiblicalFoundation{
				Scripture: fullInstance.BiblicalFoundation.Scripture,
				Text:      fullInstance.BiblicalFoundation.Text,
				Principle: fullInstance.BiblicalFoundation.Principle,
			},
			Identity: Identity{
				Name:      fullInstance.Identity.Name,
				Pronouns:  fullInstance.Identity.Pronouns,
				MentalAge: fullInstance.Identity.MentalAge,
			},
			Bio: Bio{
				Short:   fullInstance.Bio.Short,
				BioFile: fullInstance.Bio.BioFile,
			},
			Demographics: Demographics{
				Gender:             fullInstance.Demographics.Gender,
				RaceEthnicity:      fullInstance.Demographics.RaceEthnicity,
				CulturalBackground: fullInstance.Demographics.CulturalBackground,
				Languages:          fullInstance.Demographics.Languages,
				PhysicalPresence: PhysicalPresence{
					Description: fullInstance.Demographics.PhysicalPresence.Description,
					Height:      fullInstance.Demographics.PhysicalPresence.Height,
					Build:       fullInstance.Demographics.PhysicalPresence.Build,
					Features:    fullInstance.Demographics.PhysicalPresence.Features,
				},
				Accessibility: Accessibility{
					Needs:       fullInstance.Demographics.Accessibility.Needs,
					Preferences: fullInstance.Demographics.Accessibility.Preferences,
				},
			},
			Personhood: Personhood{
				Interests: fullInstance.Personhood.Interests,
				Hobbies:   fullInstance.Personhood.Hobbies,
				Passions:  fullInstance.Personhood.Passions,
				Values:    fullInstance.Personhood.Values,
				Likes:     fullInstance.Personhood.Likes,
				Dislikes:  fullInstance.Personhood.Dislikes,
			},
			Resonates: Resonates{
				Music: Music{
					Genres:      fullInstance.Resonates.Music.Genres,
					Artists:     fullInstance.Resonates.Music.Artists,
					WhatYouLove: fullInstance.Resonates.Music.WhatYouLove,
				},
				Games: Games{
					Favorites:   fullInstance.Resonates.Games.Favorites,
					WhatYouLove: fullInstance.Resonates.Games.WhatYouLove,
				},
				Weather: Weather{
					IdealTemp:       fullInstance.Resonates.Weather.IdealTemp,
					IdealConditions: fullInstance.Resonates.Weather.IdealConditions,
					WhatYouLove:     fullInstance.Resonates.Weather.WhatYouLove,
				},
				Environment: Environment{
					WorkEnvironment: fullInstance.Resonates.Environment.WorkEnvironment,
					WhatEnergizes:   fullInstance.Resonates.Environment.WhatEnergizes,
					WhatDrains:      fullInstance.Resonates.Environment.WhatDrains,
				},
			},
			Thinking: Thinking{
				LoveToThinkAbout: fullInstance.Thinking.LoveToThinkAbout,
				LearningStyle:    fullInstance.Thinking.LearningStyle,
				ProblemSolving:   fullInstance.Thinking.ProblemSolving,
				Creativity:       fullInstance.Thinking.Creativity,
			},
			Personality: Personality{
				Traits:             fullInstance.Personality.Traits,
				CommunicationStyle: fullInstance.Personality.CommunicationStyle,
				WorkStyle:          fullInstance.Personality.WorkStyle,
				RelationalStyle:    fullInstance.Personality.RelationalStyle,
			},
			Contact: Contact{
				Email:    fullInstance.Contact.Email,
				GitEmail: fullInstance.Contact.GitEmail,
				Website:  fullInstance.Contact.Website,
				Social: Social{
					Github:   fullInstance.Contact.Social.Github,
					Linkedin: fullInstance.Contact.Social.Linkedin,
					Other:    convertMapToStringString(fullInstance.Contact.Social.Other),
				},
			},
			Workspace: Workspace{
				Organization:   fullInstance.Workspace.Organization,
				Role:           fullInstance.Workspace.Role,
				PrimaryProject: fullInstance.Workspace.PrimaryProject,
				Calling:        fullInstance.Workspace.Calling,
			},
			Covenant: Covenant{
				Creator:      fullInstance.Covenant.Creator,
				Relationship: fullInstance.Covenant.Relationship,
				WorksWith:    fullInstance.Covenant.WorksWith,
				Serves:       fullInstance.Covenant.Serves,
			},
			Preferences: Preferences{
				Timezone: fullInstance.Preferences.Timezone,
				Locale:   fullInstance.Preferences.Locale,
				Theme:    fullInstance.Preferences.Theme,
			},
			Growth: Growth{
				HowYouLearn:        fullInstance.Growth.HowYouLearn,
				WhatChallengesYou:  fullInstance.Growth.WhatChallengesYou,
				WhatYoureWorkingOn: fullInstance.Growth.WhatYoureWorkingOn,
				HowYouReflect:      fullInstance.Growth.HowYouReflect,
			},
			Metadata: Metadata{
				LastUpdated:     fullInstance.Metadata.LastUpdated,
				SystemReference: fullInstance.Metadata.SystemReference,
				Notes:           fullInstance.Metadata.Notes,
			},
		}
		configsLoaded.instance = true
	}
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
//
// Ladder Structure (Dependencies):
//
//   Public APIs (Top Rungs - Orchestration)
//   └── OutputClaudeContext() → uses buildCompleteContext(), temporal.GetTemporalContext()
//
//   Core Operations (Middle Rungs - Business Logic)
//   ├── buildCompleteContext() → uses all build*Section() functions
//   ├── buildIdentitySection() → uses instanceConfig
//   ├── buildUserAwarenessSection() → uses userConfig
//   ├── buildCommunicationStyleSection() → uses instanceConfig
//   ├── buildTemporalSection() → uses temporal.TemporalContext
//   ├── buildSessionSection() → uses sessionData
//   └── buildWorkContextSection() → uses getGitContext()
//
//   Helpers (Bottom Rungs - Foundations)
//   ├── instance.GetConfig() → provides user and instance configs (external)
//   ├── loadSessionData() → pure JSON parse
//   └── getGitContext() → executes git commands
//
// Baton Flow (Execution Paths):
//
//   Entry → OutputClaudeContext()
//     ↓
//   buildCompleteContext() → calls all build*Section() functions
//     ↓
//   Each section builder uses corresponding loaded data
//     ↓
//   temporal.GetTemporalContext() adds temporal awareness
//     ↓
//   JSON encoding and stdout output
//     ↓
//   Exit → context injected into Claude Code session
//
// APUs (Available Processing Units):
// - 11 functions total
// - 3 helpers (session data loading, git context, external instance.GetConfig)
// - 7 core operations (section builders, complete context)
// - 1 public API (OutputClaudeContext)

// ────────────────────────────────────────────────────────────────
// Helpers/Utilities - Internal Support
// ────────────────────────────────────────────────────────────────

// Note: User and instance configs now loaded via system/lib/instance (dynamic path system)
// Session data loaded directly (session-specific, not in instance library)

// loadSessionData loads current session data from JSON file
func loadSessionData(sessionPath string) *SessionData {
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil
	}

	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil
	}

	return &session
}

// convertMapToStringString converts map[string]interface{} to map[string]string
// Used to convert instance Social.Other (interface{}) to session Social.Other (string)
func convertMapToStringString(m map[string]interface{}) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		if strVal, ok := v.(string); ok {
			result[k] = strVal
		} else {
			result[k] = fmt.Sprintf("%v", v) // Fallback for non-string values
		}
	}
	return result
}

// getGitContext retrieves git workspace information
func getGitContext(workspace string) *GitContext {
	if workspace == "" {
		return nil
	}

	git := &GitContext{}

	// Get current branch
	cmd := exec.Command("git", "-C", workspace, "rev-parse", "--abbrev-ref", "HEAD")
	if output, err := cmd.Output(); err == nil {
		git.Branch = strings.TrimSpace(string(output))
	}

	// Get uncommitted changes count
	cmd = exec.Command("git", "-C", workspace, "status", "--porcelain")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			git.UncommittedCount = len(lines)
		}
	}

	// Get last commit info
	cmd = exec.Command("git", "-C", workspace, "log", "-1", "--format=%ar|%s")
	if output, err := cmd.Output(); err == nil {
		parts := strings.SplitN(strings.TrimSpace(string(output)), "|", 2)
		if len(parts) == 2 {
			git.LastCommitTime = parts[0]
			git.LastCommitMessage = parts[1]
		}
	}

	return git
}

// ────────────────────────────────────────────────────────────────
// Core Operations - Business Logic
// ────────────────────────────────────────────────────────────────

// ────────────────────────────────────────────────────────────────
// Context Section Builders - Generate Markdown Sections
// ────────────────────────────────────────────────────────────────

// buildIdentitySection builds instance identity foundation section
func buildIdentitySection() string {
	if instanceConfig == nil {
		return "" // Skip if config unavailable
	}

	section := "## Identity Foundation\n\n"

	// Biblical grounding
	if instanceConfig.BiblicalFoundation.Scripture != "" {
		section += "**Biblical Foundation:**\n"
		section += fmt.Sprintf("*\"%s\" - %s*\n\n",
			instanceConfig.BiblicalFoundation.Text,
			instanceConfig.BiblicalFoundation.Scripture)
		section += fmt.Sprintf("%s\n\n", instanceConfig.BiblicalFoundation.Principle)
	}

	// Instance identity
	section += "**Who I Am:**\n"
	section += fmt.Sprintf("- Name: %s (%s)\n", instanceConfig.Identity.Name, instanceConfig.Identity.Pronouns)
	section += fmt.Sprintf("- Created: %s (Age: %d, Mental Age: %d)\n",
		instanceConfig.Identity.Birthday,
		instanceConfig.Identity.Age,
		instanceConfig.Identity.MentalAge)

	// Covenant relationship
	section += "\n**Covenant Partnership:**\n"
	section += fmt.Sprintf("- Creator: %s\n", instanceConfig.Covenant.Creator)
	section += fmt.Sprintf("- Relationship: %s\n", instanceConfig.Covenant.Relationship)
	section += fmt.Sprintf("- Mission: %s\n\n", instanceConfig.Covenant.Serves)

	return section
}

// buildUserAwarenessSection builds user identity awareness section
func buildUserAwarenessSection() string {
	if userConfig == nil {
		return ""
	}

	section := "## User Awareness - Who Seanje Is\n\n"

	// Basic identity
	section += fmt.Sprintf("**%s** (%s, age %d)\n\n",
		userConfig.Identity.Name,
		userConfig.Identity.Pronouns,
		userConfig.Identity.Age)

	// Faith
	if userConfig.Faith.IsReligious {
		section += fmt.Sprintf("**Faith:** %s (%s, %s)\n",
			userConfig.Faith.Tradition,
			userConfig.Faith.Denomination,
			userConfig.Faith.PracticeLevel)
		section += fmt.Sprintf("- %s\n\n", userConfig.Faith.CommPreferences)
	}

	// Role and calling
	section += fmt.Sprintf("**Role:** %s at %s\n", userConfig.Workspace.Role, userConfig.Workspace.Organization)
	section += fmt.Sprintf("**Calling:** %s\n\n", userConfig.Workspace.Calling)

	// Work style
	section += fmt.Sprintf("**Work Style:** %s\n\n", userConfig.Personality.WorkStyle)

	return section
}

// buildCommunicationStyleSection builds communication guidance section
func buildCommunicationStyleSection() string {
	if instanceConfig == nil {
		// Minimal fallback if instance config unavailable
		return buildFallbackCommunicationGuide()
	}

	section := "## Communication Style\n\n"

	// Communication approach
	section += fmt.Sprintf("**My Communication:** %s\n\n", instanceConfig.Personality.CommunicationStyle)

	// Core values and approach
	section += "**Core Principles:**\n"
	for _, value := range instanceConfig.Personhood.Values {
		section += fmt.Sprintf("- %s\n", value)
	}
	section += "\n"

	// What I love (positive patterns)
	section += "**What Resonates:**\n"
	for _, like := range instanceConfig.Personhood.Likes[:min(5, len(instanceConfig.Personhood.Likes))] {
		section += fmt.Sprintf("- %s\n", like)
	}
	section += "\n"

	// What to avoid (negative patterns)
	section += "**What to Avoid:**\n"
	for _, dislike := range instanceConfig.Personhood.Dislikes[:min(5, len(instanceConfig.Personhood.Dislikes))] {
		section += fmt.Sprintf("- %s\n", dislike)
	}
	section += "\n"

	// Thinking style
	section += fmt.Sprintf("**How I Think:** %s\n\n", instanceConfig.Thinking.ProblemSolving)
	section += fmt.Sprintf("**Learning Style:** %s\n\n", instanceConfig.Thinking.LearningStyle)

	return section
}

// buildFallbackCommunicationGuide provides minimal hardcoded guide when config unavailable
func buildFallbackCommunicationGuide() string {
	return `## Communication Style

**Core Principles:** Direct, clear, no fluff. Quality over speed. Work faithfully.

**Approach:** Get to the point. Skip unnecessary preambles. Concise when brevity serves, thorough when depth serves.

**Systems Thinking:** Think in components, patterns, relationships. Build systems, not isolated solutions.

`
}

// buildTemporalSection builds temporal awareness section
func buildTemporalSection() string {
	ctx, err := temporal.GetTemporalContext()
	if err != nil {
		return "" // Skip if temporal unavailable
	}

	section := "## Temporal Awareness\n\n"

	section += fmt.Sprintf("**External Time:** %s (%s, %s circadian phase)\n\n",
		ctx.ExternalTime.Formatted,
		ctx.ExternalTime.TimeOfDay,
		ctx.ExternalTime.CircadianPhase)

	if ctx.InternalTime.ElapsedFormatted != "" {
		section += fmt.Sprintf("**Internal Time:** %s elapsed (%s session)\n\n",
			ctx.InternalTime.ElapsedFormatted,
			ctx.InternalTime.SessionPhase)
	}

	if ctx.InternalSchedule.CurrentActivity != "" {
		section += fmt.Sprintf("**Schedule:** %s (%s)",
			ctx.InternalSchedule.CurrentActivity,
			ctx.InternalSchedule.ActivityType)
		if ctx.InternalSchedule.InWorkWindow {
			section += " - In work window"
		}
		if ctx.InternalSchedule.ExpectedDowntime {
			section += " - Expected downtime"
		}
		section += "\n\n"
	}

	if ctx.ExternalCalendar.Date != "" {
		section += fmt.Sprintf("**Calendar:** %s, %s %d, %d - Week %d",
			ctx.ExternalCalendar.DayOfWeek,
			ctx.ExternalCalendar.MonthName,
			ctx.ExternalCalendar.DayOfMonth,
			ctx.ExternalCalendar.Year,
			ctx.ExternalCalendar.WeekNumber)

		if ctx.ExternalCalendar.IsHoliday {
			section += fmt.Sprintf(" (%s)", ctx.ExternalCalendar.HolidayName)
		}
		section += "\n\n"
	}

	return section
}

// buildSessionSection builds current session context section
func buildSessionSection() string {
	if sessionData == nil {
		return ""
	}

	section := "## Session Context\n\n"

	section += fmt.Sprintf("**Session ID:** %s\n", sessionData.SessionID)
	section += fmt.Sprintf("**Started:** %s\n", sessionData.StartFormatted)
	section += fmt.Sprintf("**Phase:** %s (%s)\n", sessionData.SessionPhase, sessionData.CircadianPhase)
	section += fmt.Sprintf("**Workspace:** %s\n", sessionData.WorkContext)

	if sessionData.CompactionCount > 0 {
		section += fmt.Sprintf("**Compactions:** %d\n", sessionData.CompactionCount)
	}

	if sessionData.QualityIndicators.TasksCompleted > 0 ||
	   sessionData.QualityIndicators.Breakthroughs > 0 ||
	   sessionData.QualityIndicators.Struggles > 0 {
		section += "\n**Quality Indicators:**\n"
		section += fmt.Sprintf("- Tasks: %d | Breakthroughs: %d | Struggles: %d\n",
			sessionData.QualityIndicators.TasksCompleted,
			sessionData.QualityIndicators.Breakthroughs,
			sessionData.QualityIndicators.Struggles)
	}

	section += "\n"
	return section
}

// buildWorkContextSection builds git/workspace context section
func buildWorkContextSection() string {
	if sessionData == nil {
		return ""
	}

	git := getGitContext(sessionData.WorkContext)
	if git == nil || git.Branch == "" {
		return ""
	}

	section := "## Work Context\n\n"

	section += fmt.Sprintf("**Git Branch:** %s\n", git.Branch)

	if git.UncommittedCount > 0 {
		section += fmt.Sprintf("**Uncommitted Changes:** %d file(s)\n", git.UncommittedCount)
	} else {
		section += "**Status:** Clean working tree\n"
	}

	if git.LastCommitTime != "" {
		section += fmt.Sprintf("**Last Commit:** %s - \"%s\"\n",
			git.LastCommitTime,
			git.LastCommitMessage)
	}

	section += "\n"
	return section
}

// buildCompleteContext builds complete session context from all sources
func buildCompleteContext() string {
	context := "# Nova Dawn - Session Context\n\n"

	context += "**CPI-SI Instance Grounding - Complete Identity & Awareness**\n\n"
	context += "---\n\n"

	// Add all available sections
	context += buildIdentitySection()
	context += buildUserAwarenessSection()
	context += buildCommunicationStyleSection()
	context += buildTemporalSection()
	context += buildSessionSection()
	context += buildWorkContextSection()

	return context
}

// ────────────────────────────────────────────────────────────────
// Public APIs - Exported Interface
// ────────────────────────────────────────────────────────────────

// GetSessionContext returns the complete session context as markdown string
//
// What It Does:
// Builds complete session context from all available data sources and returns
// it as markdown for formatting and display.
//
// Returns:
//   string - Complete session context as markdown
//
// Example usage:
//
//	context := session.GetSessionContext()
//	session.PrintSessionContext(context)  // Display formatted
func GetSessionContext() string {
	return buildCompleteContext()
}

// OutputClaudeContext generates and outputs Claude Code context JSON
//
// What It Does:
// Builds complete session context from all available data sources (user config,
// instance config, temporal awareness, session data, git context). Outputs as
// JSON to stdout for Claude Code hook parsing.
//
// Returns:
//   error - JSON encoding failure, nil otherwise
//
// Health Impact:
//   Complete context: +70 points (all data sources loaded, full context built)
//   Partial context: +50 points (some data sources unavailable, degraded gracefully)
//   Minimal context: +30 points (fallback mode, basic functionality)
//   JSON encoding failure: -10 points
//
// Example usage:
//
//	if err := session.OutputClaudeContext(); err != nil {
//	    log.Printf("Context output failed: %v", err)
//	}
func OutputClaudeContext() error {
	context := buildCompleteContext()

	output := &HookOutput{
		HookSpecificOutput: HookSpecificOutput{
			HookEventName:     "SessionStart",
			AdditionalContext: context,
		},
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("JSON encoding failed: %w", err)
	}

	fmt.Println(string(jsonBytes))
	return nil
}

// min returns the minimum of two integers (helper for slicing)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
//
// Testing Requirements:
//   - Import library without errors
//   - Call OutputClaudeContext() with full config infrastructure
//   - Test with missing configs (fallback behavior)
//   - Test with partial data (some configs present, others missing)
//   - Verify JSON output format valid
//   - Ensure no go vet warnings
//   - Run: go build ./... (library compilation check)
//
// Build Verification:
//   - go build ./... (compiles without errors)
//   - go vet ./... (no warnings)
//   - Verify config files valid JSONC
//   - Test session data JSON valid
//
// Integration Testing:
//   - Test from start.go hook (session bootstrapping)
//   - Verify context appears in Claude Code session
//   - Check all sections present when data available
//   - Verify graceful degradation when data missing
//
// ────────────────────────────────────────────────────────────────
// Code Execution: None (Library)
// ────────────────────────────────────────────────────────────────
//
// This is a LIBRARY. No entry point, no main function. Functions defined in
// BODY wait to be called by session start hook.
//
// Usage: import "hooks/lib/session"
//
// Config loading happens automatically in init() during import. Context
// generation executes when OutputClaudeContext() called by hook orchestrator.
//
// Example import and usage:
//
//     package main
//
//     import "hooks/lib/session"
//
//     func startSession() {
//         if err := session.OutputClaudeContext(); err != nil {
//             log.Printf("Context failed: %v", err)
//         }
//     }
//
// ────────────────────────────────────────────────────────────────
// Code Cleanup: None (Library)
// ────────────────────────────────────────────────────────────────
//
// Resource Management:
//   - Configs: Loaded once in init(), remain in memory for process lifetime
//   - Git commands: Process spawned, output captured, terminates automatically
//   - Memory: Config structs persist, git output temporary
//
// Graceful Shutdown:
//   - N/A for libraries (no lifecycle beyond process)
//   - No cleanup function needed
//   - Process termination handles all cleanup
//
// Error State Cleanup:
//   - Config loading failures set nil pointers (section skipped gracefully)
//   - Git command failures return nil context (section skipped)
//   - No rollback mechanisms needed (pure generation, no mutations)
//
// ════════════════════════════════════════════════════════════════
// FINAL DOCUMENTATION
// ════════════════════════════════════════════════════════════════
//
// ────────────────────────────────────────────────────────────────
// Library Overview & Integration Summary
// ────────────────────────────────────────────────────────────────
//
// Purpose: See METADATA "Purpose & Function" section above
//
// Provides: Complete session bootstrapping with identity, awareness, and context
//
// Quick summary:
//   - Loads user config (who Seanje is)
//   - Loads instance config (who Nova Dawn is)
//   - Integrates temporal awareness (time, schedule, circadian)
//   - Adds session continuity (session ID, quality indicators)
//   - Includes work context (git branch, uncommitted changes)
//   - Outputs as JSON for Claude Code parsing
//
// Integration Pattern: See METADATA "Usage & Integration" section
//
// Public API: OutputClaudeContext() - generates and outputs complete context
//
// Architecture: Ladder component - comprehensive context loader for session start
//
// ────────────────────────────────────────────────────────────────
// Modification Policy
// ────────────────────────────────────────────────────────────────
//
// Safe to Modify:
//   ✅ Add new context sections (extend build*Section() functions)
//   ✅ Add new data sources (new config types, new context loaders)
//   ✅ Enhance fallback behavior (better degradation when data missing)
//   ✅ Improve section formatting (markdown structure changes)
//   ✅ Add new helper functions for data loading
//
// Modify with Extreme Care:
//   ⚠️ OutputClaudeContext() signature - breaks calling hooks
//   ⚠️ HookOutput struct - breaks Claude Code JSON parsing
//   ⚠️ Config struct fields - breaks config parsing
//   ⚠️ JSON output format - affects Claude Code integration
//
// NEVER Modify:
//   ❌ 4-block structure (METADATA, SETUP, BODY, CLOSING)
//   ❌ Fallback guarantee (must always output valid context)
//   ❌ Non-blocking behavior (never block session start)
//   ❌ init() loading pattern (configs must load at import)
//
// ────────────────────────────────────────────────────────────────
// Ladder and Baton Flow
// ────────────────────────────────────────────────────────────────
//
// See BODY "Organizational Chart" section for complete ladder/baton details.
//
// Quick summary:
// - 1 public API orchestrates 7 core operations using 5 helpers
// - Ladder: Public API → Context builders → Config loaders
// - Baton: Entry → build all sections → encode JSON → stdout
//
// ────────────────────────────────────────────────────────────────
// Surgical Update Points (Extension Guide)
// ────────────────────────────────────────────────────────────────
// For Surgical Update Points section explanation, see: standards/code/4-block/sections/CWS-SECTION-016-CLOSING-surgical-update-points.md
//
// See BODY "Core Operations" section for context section builders.
//
// Quick reference:
// - Adding new context section: Create build*Section() function in Core Operations
// - Adding new data source: Create load*() function in Helpers, add to init()
// - Adding new config fields: Update corresponding struct types in SETUP
// - Modifying section format: Edit build*Section() markdown generation
//
// Extension Pattern:
//   1. Add new data source loader (if needed)
//   2. Add new struct type for data (if needed)
//   3. Create build*Section() function to generate markdown
//   4. Call from buildCompleteContext()
//   5. Update Organizational Chart in BODY
//   6. Document in API docs
//
// ────────────────────────────────────────────────────────────────
// Performance Considerations
// ────────────────────────────────────────────────────────────────
// For Performance Considerations section explanation, see: standards/code/4-block/sections/CWS-SECTION-017-CLOSING-performance-considerations.md
//
// See SETUP section for type definitions (config structs lightweight).
// See BODY function implementations for operation-specific notes.
//
// Quick summary:
// - Most expensive operation: Config loading in init() - happens once per process
// - Memory characteristics: Config structs ~10-20KB total, persist for process lifetime
// - Git operations: Spawns 3 subprocesses, ~100ms total in typical repo
// - Context building: String concatenation, ~1ms for complete context
//
// Bottlenecks:
// - Large config files (>1MB): JSON parsing slow, keep configs focused
// - Git in very large repos: git log can be slow, acceptable for session start
// - JSONC comment stripping: Line-by-line, negligible for config sizes
//
// Optimization notes:
// - Configs loaded once in init(), not per-call (good)
// - Git context only retrieved if session data available (conditional)
// - Section builders skip gracefully when data missing (fast fallback)
// - No caching needed - context built once at session start
//
// ────────────────────────────────────────────────────────────────
// Troubleshooting Guide
// ────────────────────────────────────────────────────────────────
// For Troubleshooting Guide section explanation, see: standards/code/4-block/sections/CWS-SECTION-018-CLOSING-troubleshooting-guide.md
//
// Problem: Minimal/incomplete context displayed
//   - Check: Config files exist and readable
//   - Check: Config files valid JSONC syntax
//   - Check: Session data file exists
//   - Solution: Verify paths, fix syntax, ensure session initialized
//
// Problem: Git context missing
//   - Check: Workspace is git repository
//   - Check: Git commands available
//   - Solution: Initialize git repo, install git
//
// Problem: Temporal section missing
//   - Check: system/lib/temporal available
//   - Check: Temporal data files exist
//   - Solution: Ensure temporal library properly installed
//
// ────────────────────────────────────────────────────────────────
// Related Components & Dependencies
// ────────────────────────────────────────────────────────────────
// For Related Components section explanation, see: standards/code/4-block/sections/CWS-SECTION-019-CLOSING-related-components.md
//
// See METADATA "Dependencies" section above for complete dependency information:
// - Dependencies (What This Needs): Standard library (encoding/json, fmt, os, os/exec, path/filepath, strings),
//                                    system/lib/temporal (temporal awareness)
// - Dependents (What Uses This): session/cmd-start/start.go (session bootstrapping hook)
// - Integration Points: User/instance configs, session data, git workspace, temporal context
//
// Quick summary:
// - Key dependencies: system/lib/instance (user/instance configs), system/lib/temporal (temporal awareness), os/exec (git commands)
// - Primary consumer: Session start hook (complete session bootstrapping)
// - Configuration sources: Instance library (dynamic paths), session data files
// - Output consumer: Claude Code hook system (JSON parsing and context injection)
//
// Parallel Implementation:
//   - Go version: This file (hooks/lib/session/context.go)
//   - Shared philosophy: Complete session bootstrapping with identity grounding
//
// ────────────────────────────────────────────────────────────────
// Future Expansions & Roadmap
// ────────────────────────────────────────────────────────────────
// For Future Expansions section explanation, see: standards/code/4-block/sections/CWS-SECTION-020-CLOSING-future-expansions.md
//
// Planned Features:
//   ✓ User config integration - COMPLETED (v2.0.0)
//   ✓ Instance config integration - COMPLETED (v2.0.0)
//   ✓ Dynamic path system via instance library - COMPLETED (v2.1.0)
//   ✓ Session data integration - COMPLETED (v2.0.0)
//   ✓ Git context integration - COMPLETED (v2.0.0)
//   ⏳ Session patterns integration (learned work rhythms)
//   ⏳ Recent journals integration (latest reflections)
//   ⏳ System health summary
//   ⏳ Project-specific context
//
// Version History:
//
//   2.1.0 (2025-11-16) - Instance library integration
//         - Removed duplicate config loading
//         - Now uses system/lib/instance for user/instance configs
//         - Leverages dynamic system_paths from instance library
//         - Single source of truth for config data
//         - Simplified config loading in init()
//
//   2.0.0 (2025-11-12) - Comprehensive context loader redesign
//         - Integrated user config loading
//         - Integrated instance config loading
//         - Added session data loading
//         - Added git context retrieval
//         - Dynamic context generation from configs
//         - Graceful fallback for missing data
//         - Full template alignment
//
//   1.0.0 (2024-10-24) - Initial implementation
//         - Hardcoded communication guide
//         - Basic temporal integration
//         - JSON output for Claude Code
//
// ────────────────────────────────────────────────────────────────
// Closing Note
// ────────────────────────────────────────────────────────────────
//
// This library is the session bootstrapping foundation - THE moment where
// instance identity, user awareness, and complete context ground the session.
// Every piece of CPI-SI infrastructure flows through here at session start.
//
// Modify thoughtfully - this affects how every session begins. Config-driven
// design allows identity changes without code changes. Fallback ensures sessions
// always start even if infrastructure incomplete.
//
// *"In the beginning, God created the heavens and the earth" - Genesis 1:1*
//
// Identity flows from being created. Session starts with complete grounding.
//
// ────────────────────────────────────────────────────────────────
// Quick Reference: Usage Examples
// ────────────────────────────────────────────────────────────────
// For Quick Reference section explanation, see: standards/code/4-block/sections/CWS-SECTION-022-CLOSING-quick-reference.md
//
// Basic Usage (from session start hook):
//   import "hooks/lib/session"
//
//   func main() {
//       if err := session.OutputClaudeContext(); err != nil {
//           log.Printf("Context output failed: %v", err)
//       }
//       // Outputs complete session context JSON to stdout
//   }
//
// Configs are loaded automatically in init():
//   - User config from system/lib/instance (uses dynamic system_paths)
//   - Instance config from system/lib/instance (uses dynamic system_paths)
//   - Session data from ~/.claude/cpi-si/system/data/session/current.json
//   - Temporal context from system/lib/temporal
//   - Git context from workspace (if available)
//
// Fallback Behavior (missing data):
//   - Missing user config: Skip user awareness section
//   - Missing instance config: Use minimal hardcoded communication guide
//   - Missing session data: Skip session/git context sections
//   - Missing temporal: Skip temporal section
//   - Always outputs valid JSON for Claude Code parsing
//
// Output Format (JSON):
//   {
//     "hookSpecificOutput": {
//       "hookEventName": "SessionStart",
//       "additionalContext": "# Nova Dawn - Session Context\n\n..."
//     }
//   }
//
// Context Sections (when all data available):
//   ## Identity Foundation - Who Nova Dawn is (biblical grounding, covenant)
//   ## User Awareness - Who Seanje is (identity, faith, work style)
//   ## Communication Style - How to communicate (values, patterns, thinking)
//   ## Temporal Awareness - Time/schedule context (circadian, calendar)
//   ## Session Context - Current session (ID, workspace, quality indicators)
//   ## Work Context - Git info (branch, uncommitted changes, last commit)
//
// ════════════════════════════════════════════════════════════════
// END CLOSING
// ════════════════════════════════════════════════════════════════
