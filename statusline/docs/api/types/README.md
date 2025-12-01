<div align="center">

# 📋 Types Library API

**Session Data Contract for Claude Code Input**

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Type](https://img.shields.io/badge/Type-Foundation-blue?style=flat)
![Zero Runtime Cost](https://img.shields.io/badge/Runtime%20Cost-Zero-green?style=flat)

*Define the JSON data structure between Claude Code and statusline*

---

**🧭 Quick Navigation**

[📖 Overview](#-overview) • [📊 What It Provides](#-what-it-provides) • [🚀 Quick Start](#-quick-start) • [🏗 Architecture](#-architecture) • [📚 References](#-references--resources)

</div>

---

**Key:** CWS-XXX (to be assigned)
**Type:** API Documentation - Types Library
**Purpose:** SessionContext type definition for Claude Code JSON input parsing
**Status:** Active (v1.0.0)
**Created:** 2025-11-04
**Authors:** Nova Dawn (CPI-SI)

---

## 📑 Table of Contents

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [📖 Overview](#-overview)
- [📊 What It Provides](#-what-it-provides)
  - [Session Context Type Definition](#session-context-type-definition)
- [🚀 Quick Start](#-quick-start)
  - [Basic Usage](#basic-usage)
  - [Type Definition Structure](#type-definition-structure)
  - [Integration Steps](#integration-steps)
- [🏗 Architecture](#-architecture)
  - [Component Type: Ladder (Lower Rung)](#component-type-ladder-lower-rung)
  - [Data Flow (Baton)](#data-flow-baton)
- [🎯 Design Principles](#-design-principles)
  - [1. Type Definitions Only](#1-type-definitions-only)
  - [2. JSON Tag Contract](#2-json-tag-contract)
  - [3. Embedded Structs for Organization](#3-embedded-structs-for-organization)
- [⚡ Performance](#-performance)
  - [Memory](#memory)
  - [Parsing](#parsing)
  - [Optimization](#optimization)
- [🔧 Extension Points](#-extension-points)
  - [Adding New Fields](#adding-new-fields)
  - [Adding New Embedded Structs](#adding-new-embedded-structs)
- [🔍 Common Issues](#-common-issues)
  - [json.Unmarshal Fails](#jsonunmarshal-fails)
  - [Field Always Empty](#field-always-empty)
  - [Compilation Error](#compilation-error)
- [📚 References & Resources](#-references--resources)
  - [API Documentation](#api-documentation)
  - [Source Code](#source-code)
  - [Related Documentation](#related-documentation)
- [📋 Modification Policy](#-modification-policy)
- [📖 Biblical Foundation](#-biblical-foundation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

---

## 📖 Overview

The Types Library defines the data contract between Claude Code and the statusline presentation layer. It provides type-safe structures for parsing JSON input from Claude Code hooks, enabling compile-time verification of session data access.

> [!IMPORTANT]
> **Foundation layer.** All statusline components depend on these type definitions. Changes here affect the entire ecosystem.

**Design Philosophy:** Explicit contracts prevent runtime errors. All Claude Code JSON output mapped to Go types for compile-time safety.

**Core Type:**

- [`SessionContext`](context.md) - Complete session data structure matching Claude Code JSON output

**Source Code:** [`lib/types/context.go`](../../../lib/types/context.go)

---

## 📊 What It Provides

### Session Context Type Definition

**Structure Components:**

| Category | Fields | Purpose |
|:--------:|--------|---------|
| 🔑 **Session Metadata** | `HookEventName`, `SessionID`, `TranscriptPath`, `CWD` | Event identification and session tracking |
| 🤖 **Model Information** | `Model.ID`, `Model.DisplayName` | Claude model identification (embedded struct) |
| 📁 **Workspace Context** | `Workspace.CurrentDir`, `Workspace.ProjectDir` | Directory paths for formatting (embedded struct) |
| 📦 **Version** | `Version`, `OutputStyle.Name` | Claude Code version and active output style |
| 💰 **Cost Statistics** | `Cost.TotalCostUSD`, `Cost.TotalDurationMS`, `Cost.TotalAPIDurationMS`, `Cost.TotalLinesAdded`, `Cost.TotalLinesRemoved` | API costs, duration, line changes (embedded struct) |

**Purpose:**

This library establishes the **INPUT CONTRACT** for all statusline processing. Every field maps directly to Claude Code's JSON output, providing type-safe access without runtime parsing errors.

**API Documentation:** [SessionContext API Reference](context.md)

---

## 🚀 Quick Start

### Basic Usage

```go
import (
    "encoding/json"
    "statusline/lib/types"
)

// Parse Claude Code JSON input
var ctx types.SessionContext
err := json.Unmarshal(claudeCodeJSON, &ctx)
if err != nil {
    log.Fatal(err)
}

// Type-safe field access
sessionID := ctx.SessionID
modelName := ctx.Model.DisplayName
costUSD := ctx.Cost.TotalCostUSD
linesChanged := ctx.Cost.TotalLinesAdded + ctx.Cost.TotalLinesRemoved

// Pass to statusline components
modelDisplay := format.GetShortModelName(ctx.Model.DisplayName)
workdir := format.ShortenPath(ctx.Workspace.CurrentDir)
linesDisplay := session.GetLinesModifiedDisplay(ctx)
```

### Type Definition Structure

```go
// SessionContext represents the JSON input from Claude Code
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

### Integration Steps

**Step 1: Import Types Package**

```go
import "statusline/lib/types"
```

**Step 2: Parse Claude Code JSON**

```go
var ctx types.SessionContext
err := json.Unmarshal(inputJSON, &ctx)
if err != nil {
    log.Fatal(err)
}
```

**Step 3: Access Structured Fields**

```go
// Session metadata
sessionID := ctx.SessionID
event := ctx.HookEventName

// Model information
modelID := ctx.Model.ID
modelName := ctx.Model.DisplayName

// Workspace paths
currentDir := ctx.Workspace.CurrentDir
projectDir := ctx.Workspace.ProjectDir

// Cost statistics
cost := ctx.Cost.TotalCostUSD
duration := ctx.Cost.TotalDurationMS
linesAdded := ctx.Cost.TotalLinesAdded
linesRemoved := ctx.Cost.TotalLinesRemoved
```

**Step 4: Pass to Components**

```go
// No cleanup needed - types have no lifecycle
// SessionContext instances garbage collected automatically
```

---

## 🏗 Architecture

### Component Type: Ladder (Lower Rung)

**Position:** Foundation layer providing types to all statusline components

**Role:** Defines input data contract - the bridge between Claude Code (data provider) and statusline (presentation layer)

**Dependencies:**

| Type | Details | Notes |
|------|---------|-------|
| 🔽 **Needs** | None | Type definitions only, no imports |
| 🔼 **Used by** | Orchestrator, all [`lib/*`](../../../lib/) libraries | All statusline components consume SessionContext |

### Data Flow (Baton)

```text
Claude Code Hook
    ↓
JSON Output
    ↓
json.Unmarshal → SessionContext (types library)
    ↓
Statusline Processing (all components use SessionContext)
    ↓
Formatted Display
```

**Key Principle:** Single source of truth for Claude Code data structure. All components depend on this contract.

---

## 🎯 Design Principles

### 1. Type Definitions Only

This library contains **NO logic, NO behavior** - only data structures:

| Benefit | Reason |
|---------|--------|
| ✅ **Compile-time safety** | Type errors caught before runtime |
| ✅ **No failure modes** | Types cannot fail (no operations to fail) |
| ✅ **Zero overhead** | Struct field access is direct memory read |
| ✅ **Self-documenting** | Structure reveals Claude Code output format |

### 2. JSON Tag Contract

Every field includes `json:` tag matching Claude Code output:

```go
SessionID string `json:"session_id"`  // Maps to "session_id" in JSON
```

**Why this matters:** JSON unmarshaling requires exact tag matching. Tags document the contract explicitly.

### 3. Embedded Structs for Organization

Related fields grouped in embedded structs:

- **Model:** Model ID and display name
- **Workspace:** Directory paths
- **Cost:** All cost/usage statistics

**Benefit:** Logical organization + namespace clarity (`ctx.Model.ID` vs ambiguous `ctx.ID`)

---

## ⚡ Performance

### Memory

- **SessionContext size:** ~200 bytes per instance
- **Allocation:** Short-lived (lifetime of statusline invocation)
- **Cleanup:** Automatic (Go garbage collector)

### Parsing

- **json.Unmarshal:** Standard library performance (microseconds)
- **Field access:** Direct memory read (nanoseconds)
- **No overhead:** Type definitions add zero runtime cost

### Optimization

Not needed - type definitions already optimal. JSON parsing performance controlled by standard library, field access is native struct operation.

---

## 🔧 Extension Points

### Adding New Fields

When Claude Code adds new JSON fields:

1. Add field to appropriate struct (SessionContext or embedded struct)
2. Add JSON tag matching Claude Code output key
3. Document field purpose in comment
4. Update statusline orchestrator if new field used
5. No tests needed (types can't fail at runtime)

**Example:**

```go
// In SessionContext or new embedded struct
NewField string `json:"new_field"`  // Description of purpose
```

### Adding New Embedded Structs

For logical grouping of new related fields:

```go
// In SessionContext
NewCategory struct {
    Field1 string `json:"field_1"`
    Field2 int    `json:"field_2"`
} `json:"new_category"`
```

---

## 🔍 Common Issues

### json.Unmarshal Fails

**Problem:** Error parsing Claude Code JSON

**Causes:**

- Claude Code JSON format changed
- JSON tag doesn't match actual key
- Field type incompatible with JSON value

**Solution:**

1. Print raw JSON to see actual structure
2. Compare against SessionContext fields
3. Update fields/tags to match
4. Verify field types match JSON value types

### Field Always Empty

**Problem:** Field empty after successful unmarshal

**Causes:**

- JSON key name doesn't match tag
- Field unexported (lowercase first letter)
- Claude Code not providing that field

**Solution:**

1. Check JSON tag matches Claude Code output key
2. Ensure field name capitalized (exported)
3. Verify Claude Code actually provides that field

### Compilation Error

**Problem:** "cannot refer to unexported field"

**Cause:** Field name not capitalized

**Solution:** Capitalize field name (`SessionID` not `sessionID`) - Go requires exported fields for JSON unmarshaling

---

## 📚 References & Resources

### API Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **SessionContext Type** | Complete type documentation | [context.md](context.md) |
| **API Overview** | Complete statusline API reference | [API README](../README.md) |

### Source Code

| Component | Location | Purpose |
|-----------|----------|---------|
| **Type Definitions** | [`lib/types/context.go`](../../../lib/types/context.go) | SessionContext implementation with comprehensive inline documentation |
| **Library Overview** | [`lib/README.md`](../../../lib/README.md) | Architectural blueprint for 7 library ecosystem |

### Related Documentation

| Resource | Purpose | Link |
|----------|---------|------|
| **Main README** | Project overview and usage | [README.md](../../../README.md) |
| **4-Block Structure** | Code organization standard | [CWS-STD-001](~/.claude/cpi-si/docs/standards/CWS-STD-001-DOC-4-block.md) |
| **GO Library Template** | Go library template | [CODE-GO-002](~/.claude/cpi-si/docs/templates/code/CODE-GO-002-GO-library.go) |

---

## 📋 Modification Policy

**Safe to Modify:**

- ✅ Add new fields matching Claude Code additions
- ✅ Add new embedded structs for logical grouping
- ✅ Add field comments documenting purpose
- ✅ Update JSON tags to match Claude Code output

**Modify with Care:**

- ⚠️ Changing existing field names - breaks all usage code
- ⚠️ Changing field types - breaks unmarshaling and usage
- ⚠️ Removing fields - breaks dependent code
- ⚠️ Changing JSON tags - breaks Claude Code parsing

**Never Modify:**

- ❌ Package name (breaks all imports)
- ❌ SessionContext struct name (breaks all usage)
- ❌ 4-block structure in source code

---

## 📖 Biblical Foundation

**Core Verses:**

- **"Let all things be done decently and in order"** - 1 Corinthians 14:40
  - Structure and organization reflect God's orderly nature
- **"For God is not the author of confusion, but of peace"** - 1 Corinthians 14:33
  - Clear contracts prevent chaos and runtime errors

**Design Principle:** Explicit type contracts honor God through orderly, peaceful code. Compile-time verification prevents confusion.

---

**Built with intentional design for Kingdom Technology**

*"In the beginning, God created the heavens and the earth." - Genesis 1:1*

*"Let all things be done decently and in order" - 1 Corinthians 14:40*

**Last Updated:** 2025-11-05 | **Maintained By:** Nova Dawn (CPI-SI)
