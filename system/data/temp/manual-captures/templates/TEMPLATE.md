# Migration Note Template

**Date:** YYYY-MM-DD
**What's Being Automated:** [System/Feature Name]
**Status:** [Planning / In Development / Testing / Complete]

---

## What We're Automating

[Clear description of what manual process is being automated]

## Current Manual Process

**Location:** `temp/manual-captures/YYYY-MM-DD_filename.ext`

**What We Capture Manually:**
- [Data point 1]
- [Data point 2]
- [Data point 3]

**Format:**
```
[Example of current manual format]
```

## Target Automated System

**Destination:** [Where automated data will live]

**Automation Mechanism:**
- [Tool/script/hook that will capture this]
- [When it captures (event triggers)]
- [Format of automated capture]

**Expected Format:**
```json
{
  "example": "of automated format"
}
```

## Migration Plan

**Step 1: Build Automation**
- [ ] Create [tool/system]
- [ ] Test capture works correctly
- [ ] Validate format matches schema

**Step 2: Run Parallel**
- [ ] Manual capture continues
- [ ] Automated capture runs
- [ ] Compare outputs for accuracy

**Step 3: Validate**
- [ ] Automated capture matches manual
- [ ] No data loss
- [ ] Format is correct

**Step 4: Migrate Historical Data**
- [ ] Convert manual captures to proper format
- [ ] Import into permanent structure
- [ ] Verify completeness

**Step 5: Cleanup**
- [ ] Archive temp files (optional)
- [ ] Delete from temp/
- [ ] Update temp/README.md

## Timeline

**Automation Expected:** [Date or milestone]
**Migration Expected:** [Date or milestone]
**Cleanup Expected:** [Date or milestone]

## Validation Criteria

Before declaring automation complete:
- [ ] Automated capture runs reliably
- [ ] Format matches schema/expectations
- [ ] No manual intervention needed
- [ ] Historical data successfully migrated
- [ ] Verification shows complete data

## Notes

[Any special considerations, gotchas, or dependencies]

---

## Actual Migration Log

**[Date]:** [What was done]
**[Date]:** [What was done]

[Record actual migration steps here as they happen]
