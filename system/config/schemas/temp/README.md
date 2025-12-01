# Temporary Data Schemas

**Purpose:** Schemas for manual capture data during bootstrap phase while automation is being built.

## Schema: manual-capture.schema.json

**Purpose:** Validate structured manual captures in JSONC format.

**Intentionally Flexible:** The `data` object has `additionalProperties: true` because temp captures vary widely based on what's being recorded. The schema enforces basic metadata structure but allows any data structure within the `data` field.

**Required Fields:**
- `capture_metadata.date` - When captured (YYYY-MM-DD)
- `capture_metadata.type` - What type of capture
- `capture_metadata.description` - Brief description
- `data` - The actual captured information

**Optional but Recommended:**
- `capture_metadata.captured_by` - Who captured this
- `capture_metadata.automation_target` - Where it will be migrated
- `migration` - Migration planning details
- `notes` - Additional context

## Usage

The schema is referenced in manual capture templates:

```jsonc
{
  "$schema": "../../../../system/config/schemas/temp/manual-capture.schema.json",
  "capture_metadata": { ... },
  "data": { ... }
}
```

## Why Separate Temp Schema?

Temp captures are fundamentally different from permanent data:
- Structure varies by capture type
- Fields evolve as we learn what to capture
- Flexibility more important than strict validation
- Will be migrated to proper schemas once automation exists

The temp schema validates enough to be useful without being restrictive during bootstrap phase.

---

*Once automation is built and manual data migrated, these schemas may become obsolete. That's expected - they serve the bootstrap phase.*
