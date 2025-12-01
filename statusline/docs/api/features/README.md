<div align="center">

# ⚙ Features Library API

**Display Timing Logic for Kingdom Technology Reminder**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Function-1-blue?style=flat)
![Frequency](https://img.shields.io/badge/Frequency-20%25-orange?style=flat)

*Determine WHEN to show statusline elements with deterministic session-based timing*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [🔧 Public API](#-public-api) • [🚀 Quick Start](#-quick-start) • [🎯 Design](#-design-philosophy) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Features Library
**Purpose:** Kingdom Technology reminder display timing logic
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [🔧 Public API](#-public-api)
  - [ShouldShowReminder](#shouldshowreminder)
- [🚀 Quick Start](#-quick-start)
  - [Basic Integration](#basic-integration)
- [🎯 Design Philosophy](#-design-philosophy)
  - [Timing Strategy](#timing-strategy)
  - [Determinism Guarantee](#determinism-guarantee)
  - [Stateless Design](#stateless-design)
- [🧪 Demo: How It Works](#-demo-how-it-works)
  - [Deterministic Behavior](#deterministic-behavior)
  - [Frequency Distribution in Practice](#frequency-distribution-in-practice)
  - [Handles All Input Types](#handles-all-input-types)
- [⚡ Performance](#-performance)
- [🔧 Configuration](#-configuration)
  - [ReminderModulo Constant](#remindermodulo-constant)
- [🔧 Extension Points](#-extension-points)
  - [Adding New Timing Strategies](#adding-new-timing-strategies)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Reminder appears too frequently/rarely](#reminder-appears-too-frequentlyrarely)
  - [Reminder flickers on/off during session](#reminder-flickers-onoff-during-session)
- [📚 References & Resources](#-references--resources)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [🗺 Future Roadmap](#-future-roadmap)
- [📋 Modification Policy](#-modification-policy)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The Features library provides timing logic for Kingdom Technology reminder display in the statusline. It determines **WHEN** the reminder should appear (statusline orchestrator handles **HOW** to display it).

> [!IMPORTANT]
> **Timing-only library.** This library makes timing decisions. Display formatting and assembly handled by statusline orchestrator.

**Core Principle:** Deterministic session-based timing ensures consistent behavior across session restarts while maintaining ~20% display frequency.

**Biblical Foundation:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Wisdom in timing - knowing when to speak and when to be silent
- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, well-timed communication

**Source Code:** [`lib/features/reminder.go`](../../../lib/features/reminder.go)

---

## 🔧 Public API

### ShouldShowReminder

```go
func ShouldShowReminder(sessionID string) bool
```

Determines if Kingdom Technology reminder should display in the current session.

**Parameters:**

- `sessionID` (`string`): Unique identifier for current session

**Returns:**

- `bool`: `true` if reminder should display, `false` otherwise

**Behavior:**

- Uses simple hash modulo approach for deterministic distribution
- Same session ID always returns same result (deterministic guarantee)
- Achieves ~20% display frequency across different sessions
- Cannot fail - all string inputs valid, operation always succeeds

**Example Usage:**

```go
import "statusline/lib/features"

// In statusline orchestrator
if features.ShouldShowReminder(ctx.SessionID) {
    parts = append(parts, buildKingdomTechReminder())
}
```

**Health Impact:**

- Hash calculation: +50 points (always succeeds)
- Threshold evaluation: +50 points (always succeeds)
- **Total:** +100 points per call

**API Documentation:** [ShouldShowReminder Function Reference](reminder.md)

---

## 🚀 Quick Start

### Basic Integration

**Step 1: Import**

```go
import "statusline/lib/features"
```

**Step 2: Call During Statusline Assembly**

```go
func buildStatusline(ctx SessionContext) string {
    var parts []string

    // ... build other parts ...

    // Check if Kingdom Technology reminder should display
    if features.ShouldShowReminder(ctx.SessionID) {
        parts = append(parts, kingdomTechReminder)
    }

    return strings.Join(parts, " | ")
}
```

**Step 3: No Cleanup Required**

Stateless library requires no cleanup or lifecycle management.

---

## 🎯 Design Philosophy

### Timing Strategy

**Frequency:** 20% (1 in 5 sessions)

**Reasoning:**

| Frequency | Result | Why Not |
|:---------:|--------|---------|
| 50%+ | Too frequent | Becomes noise, loses impact |
| 20% | **Balanced** | Steady presence without overwhelming |
| 5% | Too rare | Loses visibility, fails to remind |

**20% strikes the balance:** Witness with work, presence without pressure.

### Determinism Guarantee

Same session ID always produces same result. This ensures:

- Reminder doesn't flicker on/off during session restarts
- Consistent user experience within a session
- Predictable testing and validation

### Stateless Design

| Property | Implementation |
|----------|----------------|
| ✅ **No package-level state** | Pure function, no globals |
| ✅ **No side effects** | Same input → same output |
| ✅ **Zero dependencies** | No imports required |
| ✅ **Lightweight** | <1 microsecond typical execution |

---

## 🧪 Demo: How It Works

### Deterministic Behavior

The same session ID always returns the same result - no randomness, no flicker:

```go
sessionID := "test-session-123"
result1 := features.ShouldShowReminder(sessionID)
result2 := features.ShouldShowReminder(sessionID)
// result1 == result2 (guaranteed every time)

// Restart the session with same ID
result3 := features.ShouldShowReminder(sessionID)
// result3 == result1 (still the same - no flicker across restarts)
```

### Frequency Distribution in Practice

Across different sessions, approximately 20% will show the reminder:

```go
showCount := 0
totalSessions := 1000

for i := 0; i < totalSessions; i++ {
    sessionID := fmt.Sprintf("session-%d", i)
    if features.ShouldShowReminder(sessionID) {
        showCount++
    }
}

frequency := float64(showCount) / float64(totalSessions)
// frequency ≈ 0.20 (20%)
// Typically ranges between 15-25% due to hash distribution
```

### Handles All Input Types

The function works with any string input - no special requirements:

```go
// Empty session ID - valid
result := features.ShouldShowReminder("")

// Very long session ID - valid (just takes longer to hash)
longID := strings.Repeat("session-id-", 1000)
result = features.ShouldShowReminder(longID)

// Special characters and Unicode - valid
result = features.ShouldShowReminder("session-™-∆-🌟")

// Typical UUID format - valid (most common)
result = features.ShouldShowReminder("550e8400-e29b-41d4-a716-446655440000")
```

---

## ⚡ Performance

**Time Complexity:** O(n) where n is session ID length

**Typical Session ID:** ~36 characters

**Typical Execution:** <1 microsecond

**Memory Usage:**

- Single int allocation for hash sum (~8 bytes)
- No heap allocations beyond input string
- No large buffers or persistent state

**Optimization:** Not needed - pure logic is already optimal. This is one of the lightest operations in the entire statusline.

---

## 🔧 Configuration

### ReminderModulo Constant

```go
const ReminderModulo = 5  // 20% display frequency
```

**To adjust frequency:**

| Modulo Value | Frequency | Usage |
|:------------:|:---------:|-------|
| 2 | 50% | Every other session |
| 3 | 33% | One in three sessions |
| 4 | 25% | One in four sessions |
| **5** | **20%** | **One in five sessions** ← Current |
| 10 | 10% | One in ten sessions |

> [!NOTE]
> **Changing this value affects all users.** Consider adding user preferences for runtime configuration in future versions.

---

## 🔧 Extension Points

### Adding New Timing Strategies

**Pattern:**

```go
// Should[Strategy]Reminder follows standard naming pattern
func ShouldTimeOfDayReminder(sessionID string, hour int) bool {
    // Implement deterministic logic
    // Same input = same output (no randomness)
    hash := calculateHash(sessionID)

    // Example: Show more during evening hours
    modulo := 5
    if hour >= 18 && hour <= 22 {
        modulo = 3  // ~33% during evening
    }

    return hash%modulo == 0
}
```

**Steps:**

1. Create function with signature: `func Should[Strategy]Reminder(...) bool`
2. Implement deterministic logic (same input = same output)
3. Document display frequency and reasoning
4. Update tests to verify expected frequency distribution
5. Update statusline orchestrator to use new strategy

---

## 🔍 Troubleshooting

### Reminder appears too frequently/rarely

**Check:**

1. Verify `ReminderModulo` constant value (should be 5 for 20%)
2. Confirm session IDs are actually unique across sessions
3. Test frequency distribution with actual session IDs

**Expected behavior:**

- Across 100 different sessions, ~20 should show reminder
- Within single session, result is always consistent

### Reminder flickers on/off during session

**Cause:** Session ID changing during session (shouldn't happen)

**Solution:** Verify session ID remains constant for duration of session

---

## 📚 References & Resources

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **ShouldShowReminder Function** | Complete function documentation | [reminder.md](reminder.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Timing Logic** | [`lib/features/reminder.go`](../../../lib/features/reminder.go) | ShouldShowReminder implementation with comprehensive inline documentation |
| **Library Overview** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 🗺 Future Roadmap

**Planned Features:**

- ⏳ Time-of-day aware timing (show more during certain hours)
- ⏳ User-configurable display frequency
- ⏳ Multiple reminder types with different frequencies

**Research Areas:**

- Circadian-aware reminder timing (align with natural work rhythms)
- Context-sensitive frequency (more frequent during onboarding, less after)
- Biblical calendar integration (special reminders on Sabbath, etc.)

**Integration Targets:**

- User preference system for reminder frequency control
- [`hooks/lib/temporal`](~/.claude/hooks/lib/temporal/) for time-of-day logic
- Session patterns for learned optimal timing

---

## 📋 Modification Policy

**Safe to Modify (Extension Points):**

- ✅ Add new timing strategies (follow `Should[Strategy]Reminder` pattern)
- ✅ Add helper functions for new hash algorithms
- ✅ Adjust `ReminderModulo` constant (changes display frequency)
- ✅ Add additional constants for new timing strategies

**Modify with Extreme Care (Breaking Changes):**

- ⚠️ `ShouldShowReminder` signature - breaks statusline orchestrator
- ⚠️ Return value semantics (true/false meaning) - breaks calling logic
- ⚠️ Determinism guarantee - breaks session consistency

**NEVER Modify (Foundational):**

- ❌ Pure function guarantee (no side effects, no state)
- ❌ Stateless design principle
- ❌ 4-block structure in source code

---

## 📖 Biblical Foundation

**Core Verses:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Wisdom in timing - knowing when to speak and when to be silent, when to remind and when to let work proceed
- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, well-timed communication

**Design Principle:** Kingdom Technology reminder serves as gentle persistent witness - not demanding constant attention, but maintaining steady presence. Like salt preserving and flavoring without dominating, reminder appears with intentional rhythm that honors rather than interrupts work.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
