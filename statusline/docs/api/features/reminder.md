<div align="center">

# ShouldShowReminder Function Reference

**Deterministic Session-Based Reminder Timing**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Function](https://img.shields.io/badge/Type-Function-blue?style=flat)
![Frequency](https://img.shields.io/badge/Frequency-20%25-orange?style=flat)

*Determines if Kingdom Technology reminder should display in the current session*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Signature](#-signature) • [⚙ Parameters](#-parameters) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - ShouldShowReminder Function
**Purpose:** Timing logic for Kingdom Technology reminder display
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [📝 Signature](#-signature)
- [⚙ Parameters](#-parameters)
  - [sessionID (string)](#sessionid-string)
- [🔄 Returns](#-returns)
  - [bool](#bool)
- [🎯 Behavior](#-behavior)
  - [Algorithm](#algorithm)
  - [Determinism](#determinism)
  - [Distribution](#distribution)
- [💻 Usage](#-usage)
  - [Basic Usage](#basic-usage)
  - [Testing Determinism](#testing-determinism)
  - [Testing Frequency](#testing-frequency)
- [🏥 Health Scoring](#-health-scoring)
- [⚡ Performance](#-performance)
- [🔧 Edge Cases](#-edge-cases)
  - [Empty String](#empty-string)
  - [Very Long Strings](#very-long-strings)
  - [Special Characters](#special-characters)
  - [Typical UUID](#typical-uuid)
- [⚙ Configuration](#-configuration)
  - [ReminderModulo Constant](#remindermodulo-constant)
- [🔍 Troubleshooting](#-troubleshooting)
  - [Reminder appears too frequently/rarely](#reminder-appears-too-frequentlyrarely)
  - [Reminder flickers on/off during session](#reminder-flickers-onoff-during-session)
- [🔗 Related Functions](#-related-functions)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**ShouldShowReminder** determines if Kingdom Technology reminder should display in the current session based on session ID hash.

> [!NOTE]
> **Pure deterministic function.** Same session ID always returns same result. No state, no side effects, no failures.

**Design:** Deterministic session-based timing ensuring consistent behavior across session restarts while maintaining ~20% display frequency.

**Library:** [`statusline/lib/features`](../../../lib/features/)

**Source Code:** [`lib/features/reminder.go`](../../../lib/features/reminder.go)

**Library README:** [Features Library Documentation](README.md)

---

## 📝 Signature

```go
func ShouldShowReminder(sessionID string) bool
```

---

## ⚙ Parameters

### sessionID (string)

Unique identifier for current session.

**Expected Format:** Any string (typically UUID format)

**Examples:**

| Input | Type | Valid |
|-------|:----:|:-----:|
| `"550e8400-e29b-41d4-a716-446655440000"` | Standard UUID | ✅ |
| `"test-session-123"` | Test session | ✅ |
| `""` | Empty string | ✅ (hash = 0) |
| `"session-™-∆-🌟"` | Special characters | ✅ (all runes contribute) |

**Validation:** None required - all string inputs valid

---

## 🔄 Returns

### bool

Indicates whether Kingdom Technology reminder should display.

**Values:**

| Value | Meaning | Frequency |
|:-----:|---------|:---------:|
| `true` | Show reminder in statusline | ~20% |
| `false` | Do not show reminder | ~80% |

**Guarantee:** Same session ID always returns same result (deterministic)

---

## 🎯 Behavior

### Algorithm

1. **Calculate hash:** Sum of all byte values in session ID string
2. **Evaluate modulo:** `hash % 5`
3. **Return decision:** `true` if result equals 0, `false` otherwise

**Display Frequency:** ~20% (1 in 5 sessions)

### Determinism

Same input always produces same output:

```go
sessionID := "test-123"
result1 := features.ShouldShowReminder(sessionID)
result2 := features.ShouldShowReminder(sessionID)
// result1 == result2 (guaranteed)

// Even after restart
result3 := features.ShouldShowReminder(sessionID)
// result3 == result1 (no flicker)
```

**Why this matters:** Reminder doesn't flicker on/off during session restarts.

### Distribution

Across different session IDs, approximately 20% will return true:

```go
showCount := 0
for i := 0; i < 1000; i++ {
    if features.ShouldShowReminder(fmt.Sprintf("session-%d", i)) {
        showCount++
    }
}
// showCount ≈ 200 (typically 150-250 due to hash distribution)
```

---

## 💻 Usage

### Basic Usage

```go
import "statusline/lib/features"

func buildStatusline(ctx SessionContext) string {
    var parts []string

    // ... add other statusline parts ...

    // Check if reminder should display
    if features.ShouldShowReminder(ctx.SessionID) {
        parts = append(parts, "⛪ Kingdom Technology")
    }

    return strings.Join(parts, " | ")
}
```

### Testing Determinism

```go
func TestDeterminism(t *testing.T) {
    sessionID := "test-session-abc-123"

    result1 := features.ShouldShowReminder(sessionID)
    result2 := features.ShouldShowReminder(sessionID)

    if result1 != result2 {
        t.Error("ShouldShowReminder not deterministic")
    }
}
```

### Testing Frequency

```go
func TestFrequency(t *testing.T) {
    showCount := 0
    totalTests := 1000

    for i := 0; i < totalTests; i++ {
        sessionID := fmt.Sprintf("session-%d", i)
        if features.ShouldShowReminder(sessionID) {
            showCount++
        }
    }

    frequency := float64(showCount) / float64(totalTests)
    if frequency < 0.15 || frequency > 0.25 {
        t.Errorf("Expected ~20%% frequency, got %.1f%%", frequency*100)
    }
}
```

---

## 🏥 Health Scoring

**Base100 Scale:** 100 points total per call

**Breakdown:**

| Operation | Points | Notes |
|-----------|:------:|-------|
| Hash calculation | +50 | Always succeeds |
| Threshold evaluation | +50 | Always succeeds |
| **Total** | **+100** | Per successful call |

> [!NOTE]
> **This function cannot fail.** All operations guaranteed to succeed. Health tracking demonstrates operation completion rather than error detection.

---

## ⚡ Performance

**Time Complexity:** O(n) where n is session ID length

**Typical Performance:**

| Metric | Value | Notes |
|--------|:-----:|-------|
| **Session ID length** | ~36 chars | UUID format |
| **Execution time** | <1 μs | Microseconds |
| **Memory** | ~8 bytes | Single int allocation for hash sum |

**Memory:**

- Single int allocation for hash sum (~8 bytes)
- No heap allocations beyond input string
- No persistent state

**Optimization:** Not needed - pure logic is already optimal. This is one of the lightest operations in the entire statusline.

---

## 🔧 Edge Cases

### Empty String

```go
features.ShouldShowReminder("")
// Valid input - hash = 0, may or may not trigger display
```

**Behavior:** Hash calculation treats empty string as sum of zero bytes (hash = 0)

### Very Long Strings

```go
longID := strings.Repeat("session-id-", 1000)
features.ShouldShowReminder(longID)
// Valid input - O(n) performance, still fast
```

**Behavior:** Longer strings take proportionally longer to hash, but still sub-millisecond

### Special Characters

```go
features.ShouldShowReminder("session-™-∆-🌟")
// Valid input - all runes contribute to hash
```

**Behavior:** Unicode characters contribute their byte values to hash sum

### Typical UUID

```go
features.ShouldShowReminder("550e8400-e29b-41d4-a716-446655440000")
// Most common format - optimal performance
```

**Behavior:** Standard case - 36 characters, deterministic hash, <1μs execution

---

## ⚙ Configuration

### ReminderModulo Constant

```go
const ReminderModulo = 5  // 20% display frequency
```

**Location:** [`lib/features/reminder.go`](../../../lib/features/reminder.go) (SETUP section)

**To Adjust Frequency:**

| Modulo Value | Frequency | Usage |
|:------------:|:---------:|-------|
| 2 | 50% | Every other session |
| 3 | 33% | One in three sessions |
| 4 | 25% | One in four sessions |
| **5** | **20%** | **One in five sessions** ← Current |
| 10 | 10% | One in ten sessions |

> [!WARNING]
> **Changing this value affects all users globally.** Consider adding user preferences for runtime configuration in future versions.

---

## 🔍 Troubleshooting

### Reminder appears too frequently/rarely

**Check:**

1. Verify `ReminderModulo` constant value (should be 5 for 20%)
2. Confirm session IDs are unique across sessions
3. Test frequency distribution with actual session IDs

**Expected:**

- Across 100 different sessions: ~20 should show reminder
- Within single session: result always consistent

**Diagnosis:**

```go
// Test frequency with your actual session IDs
showCount := 0
sessionIDs := []string{/* your actual session IDs */}
for _, id := range sessionIDs {
    if features.ShouldShowReminder(id) {
        showCount++
    }
}
fmt.Printf("Frequency: %.1f%%\n", float64(showCount)/float64(len(sessionIDs))*100)
```

### Reminder flickers on/off during session

**Problem:** Reminder appears and disappears unexpectedly during same session

**Cause:** Session ID changing during session (shouldn't happen)

**Solution:** Verify session ID remains constant for session duration

**Diagnosis:**

```go
// Log session ID at start and during checks
log.Printf("Session ID: %s, Show: %v",
    sessionID, features.ShouldShowReminder(sessionID))
```

**Expected:** Same session ID should always produce same result throughout session

---

## 🔗 Related Functions

**Same Library:**

- None currently (single-function library)

**Future Functions:**

| Function | Purpose | Status |
|----------|---------|:------:|
| `ShouldTimeOfDayReminder()` | Time-aware timing strategy | ⏳ Planned |
| `ShouldUserConfigReminder()` | User preference-based timing | ⏳ Planned |
| `ShouldContextReminder()` | Context-sensitive frequency | 🔬 Research |

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Features Library README** | Library overview and integration guide | [README.md](README.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Function Implementation** | [`lib/features/reminder.go`](../../../lib/features/reminder.go) | ShouldShowReminder with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 📖 Biblical Foundation

**Core Verse:**

- **"To every thing there is a season, and a time to every purpose under the heaven"** - Ecclesiastes 3:1
  - Wisdom in timing - knowing when to speak and when to be silent, when to remind and when to let work proceed

**Supporting Verse:**

- **"A word fitly spoken is like apples of gold in pictures of silver"** - Proverbs 25:11
  - Concise, well-timed communication

**Application:** 20% frequency balances witness with work - steady presence without overwhelming. Like salt preserving and flavoring without dominating, reminder appears with intentional rhythm that honors rather than interrupts work.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
