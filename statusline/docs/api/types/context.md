<div align="center">

# SessionContext Type Reference

**Complete Session Data Contract for Claude Code**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Type](https://img.shields.io/badge/Type-Struct-blue?style=flat)
![Status](https://img.shields.io/badge/Status-Stable-green?style=flat)

*Type-safe access to all session metadata, model information, workspace context, and cost statistics*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📝 Type Definition](#-type-definition) • [🏗 Fields](#-fields) • [💻 Usage](#-usage) • [📚 References](#-references)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - SessionContext Type
**Purpose:** Complete type reference for Claude Code session data structure
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [📝 Type Definition](#-type-definition)
- [🏗 Fields](#-fields)
  - [Top-Level Fields](#top-level-fields)
  - [Model (Embedded Struct)](#model-embedded-struct)
  - [Workspace (Embedded Struct)](#workspace-embedded-struct)
  - [OutputStyle (Embedded Struct)](#outputstyle-embedded-struct)
  - [Cost (Embedded Struct)](#cost-embedded-struct)
- [💻 Usage](#-usage)
  - [Import](#import)
  - [Parse Claude Code JSON](#parse-claude-code-json)
  - [Access Fields](#access-fields)
  - [Pass to Components](#pass-to-components)
- [📊 Example JSON](#-example-json)
  - [Complete Session Start](#complete-session-start)
  - [Minimal Session](#minimal-session)
- [🎯 Understanding SessionContext](#-understanding-sessioncontext)
  - [What It Represents](#what-it-represents)
  - [Why Embedded Structs](#why-embedded-structs)
  - [JSON Tag Contract](#json-tag-contract)
- [🔧 Edge Cases](#-edge-cases)
  - [Empty Cost Statistics](#empty-cost-statistics)
  - [Missing Model Information](#missing-model-information)
  - [Nil Pointer (Invalid)](#nil-pointer-invalid)
- [⚡ Performance](#-performance)
- [🔍 Troubleshooting](#-troubleshooting)
  - [json.Unmarshal Fails](#jsonunmarshal-fails)
  - [Field Always Empty After Parse](#field-always-empty-after-parse)
  - [Cannot Access Field](#cannot-access-field)
  - [Type Mismatch Error](#type-mismatch-error)
- [🔧 Extension Points](#-extension-points)
  - [Adding New Top-Level Fields](#adding-new-top-level-fields)
  - [Adding New Embedded Struct](#adding-new-embedded-struct)
  - [Handling Optional Fields](#handling-optional-fields)
- [📚 References](#-references)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

**SessionContext** defines the complete data contract between Claude Code and statusline, providing type-safe access to all session metadata, model information, workspace context, and cost statistics.

> [!NOTE]
> **Foundation type.** Every statusline component receives SessionContext as input. Changes to this structure affect the entire ecosystem.

**Design:** Direct mapping from Claude Code JSON output to Go struct with embedded structures for logical organization and compile-time safety.

**Library:** [`statusline/lib/types`](../../../lib/types/)

**Source Code:** [`lib/types/context.go`](../../../lib/types/context.go)

**Library README:** [Types Library Documentation](README.md)

---

## 📝 Type Definition

```go
type SessionContext struct {
    HookEventName  string `json:"hook_event_name"`
    SessionID      string `json:"session_id"`
    TranscriptPath string `json:"transcript_path"`
    CWD            string `json:"cwd"`

    Model struct {
        ID          string `json:"id"`
        DisplayName string `json:"display_name"`
    } `json:"model"`

    Workspace struct {
        CurrentDir string `json:"current_dir"`
        ProjectDir string `json:"project_dir"`
    } `json:"workspace"`

    Version string `json:"version"`

    OutputStyle struct {
        Name string `json:"name"`
    } `json:"output_style"`

    Cost struct {
        TotalCostUSD       float64 `json:"total_cost_usd"`
        TotalDurationMS    int     `json:"total_duration_ms"`
        TotalAPIDurationMS int     `json:"total_api_duration_ms"`
        TotalLinesAdded    int     `json:"total_lines_added"`
        TotalLinesRemoved  int     `json:"total_lines_removed"`
    } `json:"cost"`
}
```

---

## 🏗 Fields

### Top-Level Fields

| Field | Type | JSON Tag | Purpose |
|-------|:----:|----------|---------|
| **HookEventName** | `string` | `hook_event_name` | Hook event triggering statusline update ("session_start", "tool_use", etc.) |
| **SessionID** | `string` | `session_id` | Unique identifier for this Claude Code session (UUID format) |
| **TranscriptPath** | `string` | `transcript_path` | Filesystem path to session transcript file |
| **CWD** | `string` | `cwd` | Current working directory (absolute path) |
| **Version** | `string` | `version` | Claude Code version string (semantic version) |

---

### Model (Embedded Struct)

Model information from Claude Code.

| Field | Type | JSON Tag | Purpose |
|-------|:----:|----------|---------|
| **Model.ID** | `string` | `id` | Full Claude model identifier (e.g., "claude-sonnet-4-5-20250929") |
| **Model.DisplayName** | `string` | `display_name` | Human-readable model name (e.g., "Sonnet", "Opus", "Haiku") |

**Access:** `ctx.Model.ID`, `ctx.Model.DisplayName`

---

### Workspace (Embedded Struct)

Workspace directory information.

| Field | Type | JSON Tag | Purpose |
|-------|:----:|----------|---------|
| **Workspace.CurrentDir** | `string` | `current_dir` | Current working directory (absolute path) |
| **Workspace.ProjectDir** | `string` | `project_dir` | Project root directory (absolute path) |

**Access:** `ctx.Workspace.CurrentDir`, `ctx.Workspace.ProjectDir`

---

### OutputStyle (Embedded Struct)

Output style configuration.

| Field | Type | JSON Tag | Purpose |
|-------|:----:|----------|---------|
| **OutputStyle.Name** | `string` | `name` | Active output style name |

**Access:** `ctx.OutputStyle.Name`

---

### Cost (Embedded Struct)

Session cost and usage statistics.

| Field | Type | JSON Tag | Purpose |
|-------|:----:|----------|---------|
| **Cost.TotalCostUSD** | `float64` | `total_cost_usd` | Total API cost for session in US dollars (e.g., 0.0123) |
| **Cost.TotalDurationMS** | `int` | `total_duration_ms` | Total session duration in milliseconds |
| **Cost.TotalAPIDurationMS** | `int` | `total_api_duration_ms` | Total API request time in milliseconds |
| **Cost.TotalLinesAdded** | `int` | `total_lines_added` | Total lines added in session |
| **Cost.TotalLinesRemoved** | `int` | `total_lines_removed` | Total lines removed in session |

**Access:** `ctx.Cost.TotalCostUSD`, `ctx.Cost.TotalDurationMS`, etc.

---

## 💻 Usage

### Import

```go
import (
    "encoding/json"
    "statusline/lib/types"
)
```

### Parse Claude Code JSON

```go
// Read JSON from stdin or file
inputJSON := []byte(`{
    "hook_event_name": "session_start",
    "session_id": "abc123",
    "model": {
        "id": "claude-sonnet-4-5-20250929",
        "display_name": "Sonnet"
    },
    "cost": {
        "total_cost_usd": 0.0123,
        "total_lines_added": 150,
        "total_lines_removed": 25
    }
}`)

var ctx types.SessionContext
err := json.Unmarshal(inputJSON, &ctx)
if err != nil {
    log.Fatal(err)
}
```

### Access Fields

```go
// Session metadata
fmt.Println("Event:", ctx.HookEventName)
fmt.Println("Session:", ctx.SessionID)

// Model information
fmt.Println("Model:", ctx.Model.DisplayName)
fmt.Println("Full ID:", ctx.Model.ID)

// Workspace context
fmt.Println("Working dir:", ctx.Workspace.CurrentDir)
fmt.Println("Project:", ctx.Workspace.ProjectDir)

// Cost statistics
fmt.Printf("Cost: $%.4f\n", ctx.Cost.TotalCostUSD)
fmt.Printf("Duration: %dms\n", ctx.Cost.TotalDurationMS)
fmt.Printf("Lines changed: %d\n",
    ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved)
```

### Pass to Components

```go
// Use with format library
modelDisplay := format.GetShortModelName(ctx.Model.DisplayName)
workdir := format.ShortenPath(ctx.Workspace.CurrentDir)

// Use with session library
linesDisplay := session.GetLinesModifiedDisplay(ctx)
costDisplay := session.GetCostDisplay(ctx)
durationDisplay := session.GetDurationDisplay(ctx)

// Use with features library
showReminder := features.ShouldShowReminder(ctx.SessionID)
```

---

## 📊 Example JSON

### Complete Session Start

```json
{
  "hook_event_name": "session_start",
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "transcript_path": "/home/user/.claude/transcripts/session.json",
  "cwd": "/home/user/project",
  "model": {
    "id": "claude-sonnet-4-5-20250929",
    "display_name": "Sonnet"
  },
  "workspace": {
    "current_dir": "/home/user/project/src",
    "project_dir": "/home/user/project"
  },
  "version": "1.0.0",
  "output_style": {
    "name": "default"
  },
  "cost": {
    "total_cost_usd": 0.0123,
    "total_duration_ms": 45000,
    "total_api_duration_ms": 3500,
    "total_lines_added": 150,
    "total_lines_removed": 25
  }
}
```

### Minimal Session

```json
{
  "hook_event_name": "session_start",
  "session_id": "abc123",
  "cwd": "/home/user/project",
  "model": {
    "id": "claude-sonnet-4-5",
    "display_name": "Sonnet"
  },
  "workspace": {
    "current_dir": "/home/user/project",
    "project_dir": "/home/user/project"
  },
  "cost": {
    "total_cost_usd": 0.0,
    "total_duration_ms": 0,
    "total_api_duration_ms": 0,
    "total_lines_added": 0,
    "total_lines_removed": 0
  }
}
```

---

## 🎯 Understanding SessionContext

### What It Represents

**SessionContext** = Complete snapshot of Claude Code session state at the moment statusline hook fires

**Components:**

| Category | Contains | Purpose |
|:--------:|----------|---------|
| 🔑 **Session Identity** | ID, event, paths | Who/when/where tracking |
| 🤖 **Model Context** | ID, display name | What Claude model is running |
| 📁 **Workspace Context** | Current dir, project dir | Where work is happening |
| 💰 **Cost Tracking** | Time, lines, cost | How much work has been done |

**Purpose:** Single unified structure containing all data statusline needs to format display

### Why Embedded Structs

**Logical organization:**

- Model data grouped together (ID + DisplayName)
- Workspace data grouped together (CurrentDir + ProjectDir)
- Cost data grouped together (all statistics)

**Benefits:**

| Benefit | Example |
|---------|---------|
| **Clear namespace** | `ctx.Model.ID` vs ambiguous `ctx.ModelID` |
| **Matches JSON structure** | Nested objects map to embedded structs |
| **Extensibility** | Add fields to existing groups without cluttering top level |

### JSON Tag Contract

Every field has `json:"..."` tag:

```go
SessionID string `json:"session_id"`
```

**Why essential:**

- Maps Go field name to JSON key (`SessionID` ← `session_id`)
- Enables `json.Unmarshal` to populate fields correctly
- Documents Claude Code output format explicitly
- **Breaking tags breaks parsing**

---

## 🔧 Edge Cases

### Empty Cost Statistics

```go
// Cost fields default to zero values if not in JSON
ctx.Cost.TotalCostUSD       // → 0.0
ctx.Cost.TotalDurationMS    // → 0
ctx.Cost.TotalLinesAdded    // → 0
```

**Behavior:** Go unmarshaling sets unspecified numeric fields to zero

### Missing Model Information

```go
// If model object missing in JSON
ctx.Model.ID          // → ""
ctx.Model.DisplayName // → ""
```

**Behavior:** Embedded struct fields default to zero values (empty strings)

### Nil Pointer (Invalid)

```go
var ctx *types.SessionContext  // nil pointer
ctx.SessionID                  // PANIC: nil pointer dereference
```

**Solution:** Always use value (not pointer) or check for nil:

```go
var ctx types.SessionContext   // ✅ Value - always safe
// OR
var ctx *types.SessionContext  // Pointer
if ctx != nil {                 // ✅ Check before access
    id := ctx.SessionID
}
```

---

## ⚡ Performance

**Memory:** ~200 bytes per SessionContext instance

**Allocation:**

- **Stack:** If declared in function scope and not escaped
- **Heap:** If returned from function or stored beyond scope
- **Typical:** Short-lived (lifetime of statusline invocation)

**Parsing:**

- **json.Unmarshal:** Standard library performance (microseconds for typical session data)
- **Field access:** Direct memory read (nanoseconds)

**Optimization:** Not needed - type definitions add zero runtime overhead

---

## 🔍 Troubleshooting

### json.Unmarshal Fails

**Problem:** `json.Unmarshal` returns error

**Causes:**

- Invalid JSON syntax
- Field type mismatch (string in JSON, int in struct)
- Unexpected JSON structure

**Solution:**

1. Validate JSON syntax: `json.Unmarshal(data, &map[string]interface{}{})`
2. Print raw JSON to inspect structure
3. Compare JSON keys against struct tags
4. Check field types match JSON value types

**Example:**

```go
// Debug: Print raw JSON structure
var raw map[string]interface{}
json.Unmarshal(data, &raw)
fmt.Printf("Raw JSON: %+v\n", raw)

// Then parse into SessionContext
var ctx types.SessionContext
err := json.Unmarshal(data, &ctx)
```

### Field Always Empty After Parse

**Problem:** Field empty despite being in JSON

**Causes:**

- JSON tag doesn't match key name
- Field unexported (lowercase first letter)
- JSON value is null

**Solution:**

1. Check JSON tag matches exact key: `json:"session_id"` for "session_id"
2. Ensure field capitalized: `SessionID` not `sessionID`
3. Check JSON value not null: `"session_id": null` → empty string

### Cannot Access Field

**Problem:** Compilation error "cannot refer to unexported field"

**Cause:** Field name not capitalized (unexported in Go)

**Solution:** Capitalize field name:

- ❌ `sessionID string` (unexported)
- ✅ `SessionID string` (exported, JSON can access)

### Type Mismatch Error

**Problem:** `json.Unmarshal` fails with type error

**Cause:** JSON value type doesn't match struct field type

**Example:**

```go
// Struct expects int
TotalLinesAdded int `json:"total_lines_added"`

// JSON provides string
"total_lines_added": "150"  // ❌ String, not number
```

**Solution:** Match types or use `json.Number` for flexible parsing

---

## 🔧 Extension Points

### Adding New Top-Level Fields

```go
type SessionContext struct {
    // ... existing fields ...

    NewField string `json:"new_field"`  // Add here with JSON tag
}
```

### Adding New Embedded Struct

```go
type SessionContext struct {
    // ... existing fields ...

    NewCategory struct {
        Field1 string `json:"field_1"`
        Field2 int    `json:"field_2"`
    } `json:"new_category"`
}
```

### Handling Optional Fields

```go
// Use pointers for truly optional fields
OptionalField *string `json:"optional_field,omitempty"`

// Access safely
if ctx.OptionalField != nil {
    value := *ctx.OptionalField
}
```

---

## 📚 References

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Types Library README** | Library overview and integration guide | [README.md](README.md) |
| **API Overview** | Complete statusline API documentation | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Type Definition** | [`lib/types/context.go`](../../../lib/types/context.go) | SessionContext implementation with comprehensive inline documentation |
| **Library Architecture** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 📖 Biblical Foundation

**Core Verses:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Structure and organization reflect God's orderly nature
- **"For God is not the author of confusion, but of peace"** - 1 Corinthians 14:33
  - Clear contracts prevent chaos and runtime errors

**Application:** SessionContext brings order to data flow - clear structure, explicit mapping, type safety. Reflects the orderly nature of God who establishes covenants and keeps them. Not chaotic runtime parsing, but structured compile-time verification.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
