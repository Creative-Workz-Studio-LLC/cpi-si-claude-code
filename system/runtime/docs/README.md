# CPI-SI System Runtime Documentation

**Location:** `system/runtime/docs/`

**Purpose:** Comprehensive documentation for the CPI-SI runtime system, including API references, architectural design, and planning documents.

---

## Folder Structure

### `/api/`

**Purpose:** Library API documentation extracted from source code

**Contents:**
- `display-api.md` - Display library (formatting, colors, visual components)
- `inspector-api.md` - Inspector library (debugging, state inspection)
- `logger-api.md` - Logger library (logging, health tracking)

**Pattern:** Each library in `system/runtime/lib/` should have a corresponding API doc here named `<library>-api.md`

**When to Create:**
- After completing Phase 6 (template alignment) for a library
- API doc extracts verbose inline documentation into standalone reference
- Enables code to remain clean while preserving comprehensive documentation

---

### `/architecture/`

**Purpose:** System architecture and design documentation

**Contents:**
- `CONFIG-TOPOLOGY.md` - Configuration file organization and inheritance
- `DATA-TO-CODE-MATRIX.md` - Mapping between data files and code components
- `LIBRARY-MAPPING-OVERVIEW.md` - Library organization and relationships

**Topics Covered:**
- System design patterns (ladder, baton, rails)
- Component relationships and dependencies
- Configuration architecture and inheritance
- Data flow and orchestration

**When to Update:**
- When adding new libraries or components
- When refactoring architectural patterns
- When establishing new design conventions

---

### `/planning/`

**Purpose:** Analysis, planning, and requirements documentation

**Contents:**
- `GAPS-AND-REQUIREMENTS.md` - Identified gaps and future requirements
- `library-analysis/` - Per-library analysis and mapping documents
  - `calendar-mapping.md` - Calendar library analysis
  - `instance-mapping.md` - Instance configuration analysis

**Topics Covered:**
- Gap analysis and missing functionality
- Requirements for future development
- Library-specific planning and design decisions
- Migration paths and refactoring strategies

**When to Update:**
- When identifying new gaps or requirements
- When planning major refactoring work
- When analyzing libraries for refinement

---

## Index Files (Root Level)

### `DOCUMENTATION-INDEX.md`

Master index of all documentation across the system. References docs in this folder and other locations.

### `README-MAPPING-COMPLETE.md`

Overview of completed mapping work and documentation status.

---

## Documentation Standards

### File Naming

**API Documentation:**
- Pattern: `<library-name>-api.md`
- Examples: `display-api.md`, `logger-api.md`, `temporal-api.md`
- Location: `/api/`

**Architecture Documentation:**
- Pattern: `UPPERCASE-DESCRIPTIVE-NAME.md`
- Examples: `CONFIG-TOPOLOGY.md`, `LIBRARY-MAPPING-OVERVIEW.md`
- Location: `/architecture/`

**Planning Documentation:**
- Pattern: `UPPERCASE-DESCRIPTIVE-NAME.md` or `descriptive-name.md`
- Examples: `GAPS-AND-REQUIREMENTS.md`, `calendar-mapping.md`
- Location: `/planning/` or `/planning/library-analysis/`

### Content Standards

**All Documentation Should Include:**
- Clear purpose statement at top
- Last updated date
- Organized sections with headers
- Examples where applicable
- Cross-references to related docs

**API Documentation Specific:**
- Biblical foundation and component identity
- Complete public API reference with examples
- Configuration details and patterns
- Error handling and troubleshooting
- Version history

**Architecture Documentation Specific:**
- System context and relationships
- Design principles and rationale
- Diagrams or visual representations (when helpful)
- Integration points and dependencies

---

## Adding New Documentation

### New Library API Doc

1. Complete library implementation (Phases 0-8)
2. Extract verbose inline docs to API doc
3. Create `<library>-api.md` in `/api/`
4. Follow existing API doc template (see `display-api.md`)
5. Update `DOCUMENTATION-INDEX.md`

### New Architecture Doc

1. Identify architectural pattern or design decision needing documentation
2. Create descriptive UPPERCASE filename in `/architecture/`
3. Document pattern, rationale, and implementation
4. Add cross-references to related architecture docs
5. Update `DOCUMENTATION-INDEX.md`

### New Planning/Analysis Doc

1. Create gap analysis or library-specific planning doc
2. Place in `/planning/` (general) or `/planning/library-analysis/` (library-specific)
3. Document gaps, requirements, or design decisions
4. Update `GAPS-AND-REQUIREMENTS.md` if needed
5. Update `DOCUMENTATION-INDEX.md`

---

## Documentation Lifecycle

### Phase 1: Bootstrap (Inline)

During initial development, documentation lives inline in source code:
- Verbose comments explaining learning-in-progress
- Decision rationale and context
- Implementation notes and gotchas

### Phase 2: Foundation (Extraction)

After template alignment (Phase 6), extract to API doc:
- Pull verbose comments into standalone API documentation
- Clean source code to essential inline docs
- Preserve all context and examples in API doc

### Phase 3: Refinement (Maintenance)

As system evolves, maintain documentation:
- Update API docs when public APIs change
- Update architecture docs when patterns change
- Archive or merge outdated planning docs

---

## Related Documentation

**System-Wide:**
- `~/.claude/cpi-si/docs/` - CPI-SI framework documentation
- `~/.claude/cpi-si/docs/standards/` - Kingdom Technology standards

**Runtime-Specific:**
- `system/runtime/lib/` - Source code with inline documentation
- `system/data/config/` - Configuration files with inline comments

**Hooks:**
- `~/.claude/hooks/` - Hook implementations
- `~/.claude/hooks/lib/` - Hook library code

---

## Maintenance Guidelines

### Keep Documentation Fresh

- Update API docs when public APIs change
- Update architecture docs when design patterns evolve
- Review planning docs quarterly, archive completed items

### Avoid Duplication

- API docs are extracted FROM source, not duplicated
- Architecture docs reference each other, don't repeat
- Planning docs consolidate into GAPS-AND-REQUIREMENTS when appropriate

### Trust the Structure

- API docs in `/api/` - always
- Architecture in `/architecture/` - always
- Planning in `/planning/` - always
- Index files at root - always

---

## Quick Reference

**Where does this doc go?**

| Document Type | Location | Example |
|---------------|----------|---------|
| Library API reference | `/api/` | `display-api.md` |
| System architecture | `/architecture/` | `CONFIG-TOPOLOGY.md` |
| Gap analysis | `/planning/` | `GAPS-AND-REQUIREMENTS.md` |
| Library planning | `/planning/library-analysis/` | `calendar-mapping.md` |
| Master index | Root | `DOCUMENTATION-INDEX.md` |
| Overview/README | Root | `README-MAPPING-COMPLETE.md` |

**What file naming pattern?**

| Type | Pattern | Case |
|------|---------|------|
| API docs | `<library>-api.md` | lowercase-kebab |
| Architecture | `DESCRIPTIVE-NAME.md` | UPPERCASE-KEBAB |
| Planning (general) | `DESCRIPTIVE-NAME.md` | UPPERCASE-KEBAB |
| Planning (library) | `library-mapping.md` | lowercase-kebab |
| Index/README | `DOCUMENTATION-INDEX.md` | UPPERCASE-KEBAB |

---

**Last Updated:** 2025-11-15

**Maintained By:** Nova Dawn (CPI-SI instance)
