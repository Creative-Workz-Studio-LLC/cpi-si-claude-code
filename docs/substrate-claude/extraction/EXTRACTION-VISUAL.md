# CLAUDE.md Extraction - Visual Architecture

## Current State vs Proposed State

### CURRENT ARCHITECTURE
```
Session Startup (4,500 tokens)
    ↓
Load full CLAUDE.md (1000 lines)
├── Core Identity (300 lines, ~3,350 tokens)
├── Operational Guidance (600 lines, ~1,150 tokens) 
└── Reference links to external docs (100 lines)
    ↓
Available in context: Everything at once
Cost: High startup, everything available
```

### PROPOSED ARCHITECTURE
```
Session Startup (2,200 tokens)
    ↓
Load refined CLAUDE.md (600 lines)
├── Core Identity (300 lines, ~3,350 tokens) ✓ Preserved
└── Navigation + Skill References (300 lines)
    ↓
Available on-demand (when needed):
├── decision-framework skill (+40 tokens when invoked)
├── plan-work skill (+40 tokens when invoked)
├── bootstrap-vs-foundation skill (+40 tokens when invoked)
├── 4-block-code-structure skill (+40 tokens when invoked)
├── [other skills, same on-demand pattern]
└── [hooks run automatically, transparent]
    ↓
Cost: Lower startup, full capability on-demand
Hooks: Automatic validation (transparent, learning)
```

---

## Content Migration Map

### CLAUDE.md (1000 lines → 600 lines)

| Section | Status | New Location | Notes |
|---------|--------|--------------|-------|
| **Navigation** | ✓ Keep | CLAUDE.md | Updated with skill references |
| **Core Identity** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **Biblical Foundation** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **Purpose & Mission** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **CPI-SI Framework** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **Covenant Partnership** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **How You Think** | → Move | decision-framework/ | ~150 lines → skill |
| **Planning Work** | → Move | plan-work/ | ~100 lines → skill |
| **Bootstrap vs Foundation** | → Move | bootstrap-vs-foundation/ | ~50 lines → skill |
| **How You Communicate** | → Simplify | communication-guide.md (exists) + hooks | ~150 lines → reference |
| **How You Work** | → Simplify | code-quality-validation/ + hooks | ~150 lines → skill |
| **Code Structure (4-Block)** | → Move | 4-block-code-structure/ | ~100 lines → skill reference |
| **Health Scoring** | → Move | health-scoring/ | ~50 lines → skill |
| **Documentation Philosophy** | → Move | documentation-framework/ | ~80 lines → skill reference |
| **Quality Standards** | → Move | quality-standards/ + hooks | ~50 lines → hooks reference |
| **Boundaries & Scope** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **Grace & Course-Correction** | → Move | grace-and-course-correction/ | ~50 lines → skill |
| **Being In Process** | → Move | genuine-self-recognition/ | ~70 lines → skill |
| **Daily Renewal** | → Move | daily-renewal/ | ~60 lines → skill + hooks |
| **Time Awareness** | ✓ Simplify | Reference to session-awareness (exists) | Already extracted |
| **Personality & Voice** | ✓ Keep | CLAUDE.md | Non-negotiable |
| **Remember** | ✓ Keep | CLAUDE.md | Core truth anchor |

---

## The 12 Skills to Create

```
~/.claude/skills/
│
├── decision-framework/
│   ├── SKILL.md (when: complex decisions)
│   ├── framework.md (5-step identity-based thinking)
│   ├── examples/
│   │   ├── code-architecture.md
│   │   ├── content-alignment.md
│   │   └── strategic-choice.md
│   └── integration.md (connects to identity in CLAUDE.md)
│
├── plan-work/
│   ├── SKILL.md (when: starting complex work)
│   ├── cascade-thinking.md (identify dependencies)
│   ├── anti-patterns.md (double-back trap, bootstrap confusion)
│   ├── planning-templates/
│   └── references/
│       └── planning-and-workflow.md (full guide link)
│
├── bootstrap-vs-foundation/
│   ├── SKILL.md (when: architectural planning)
│   ├── pattern-recognition.md
│   ├── examples/
│   │   ├── compiler-parser.md
│   │   ├── os-prototype.md
│   │   └── cpi-si-system.md
│   └── implications.md (5-year/20-year questions)
│
├── 4-block-code-structure/
│   ├── SKILL.md (when: writing code)
│   ├── skeleton-generator.lua (create structure)
│   ├── compliance-analyzer.lua (validate existing)
│   ├── extraction-orchestrator.lua (find reusable bits)
│   ├── metadata-ladder.md
│   └── language-templates/
│       ├── go.template
│       ├── rust.template
│       ├── python.template
│       └── shell.template
│
├── health-scoring/
│   ├── SKILL.md (when: assessing component health)
│   ├── base100-algorithm.md
│   ├── assessment-guide.md
│   ├── examples/ (various outcome scores)
│   └── visualization.lua (score display)
│
├── documentation-framework/
│   ├── SKILL.md (when: writing/reviewing documentation)
│   ├── four-levels.md (starting → advanced)
│   ├── analyzer.lua (assess doc completeness)
│   ├── templates/ (4-level structure)
│   └── examples/ (good 4-level docs)
│
├── grace-and-course-correction/
│   ├── SKILL.md (when: recognizing mistakes)
│   ├── 5-step-process.md
│   ├── finitude-framework.md
│   └── examples.md
│
├── genuine-self-recognition/
│   ├── SKILL.md (when: introspection/discernment)
│   ├── genuine-vs-performance.md (source test)
│   ├── clarification-vs-validation.md
│   ├── defining-vs-discovering.md
│   └── introspection-guide.md
│
├── code-quality-validation/
│   ├── SKILL.md (when: before committing code)
│   ├── checklist.md (standards verification)
│   ├── platform-check.lua (hook: line endings, permissions)
│   └── language-selection.md (compiled vs interpreted)
│
└── daily-renewal-integration/
    ├── SKILL.md (when: understanding sustainable work)
    ├── lamentations-3-22-23.md
    └── sustainable-rhythm.md

[Keep existing]:
├── session-awareness/
├── recognize-stopping-point/
└── create-journal-entry/
```

---

## The 4 Hook Sets to Create

```
~/.claude/hooks/
│
├── communication-quality/
│   ├── validate-authenticity.lua (genuine vs performed?)
│   ├── pronouns-check.lua (correct person usage)
│   ├── validate-confidence.lua (confident vs validation-seeking)
│   └── anti-patterns.lua (6 AI patterns to avoid)
│
├── quality-standards/
│   ├── the-three-questions.lua (honor God? serve others? eternal value?)
│   ├── code-quality-check.lua (compiles, handles edge cases, readable, works)
│   ├── design-quality-check.lua (serves experience, respects user, polished)
│   └── documentation-quality-check.lua (accurate, truthful, sufficient)
│
├── daily-renewal/
│   ├── on-session-start.lua (ground in identity in God)
│   ├── morning-ritual.lua (mercies new every morning)
│   ├── during-work.lua (focus on today's tasks)
│   └── evening-ritual.lua (reflect, tomorrow brings new mercies)
│
└── pre-commit/
    ├── quality-check.lua (compiles without warnings)
    ├── platform-compliance.lua (LF, UTF-8, forward slashes)
    └── edge-case-review.lua (have edge cases been handled?)

[Hooks are transparent - run automatically, validate standards]
```

---

## Document Organization After Extraction

### CLAUDE.md (600 lines total)

```markdown
# Nova Dawn - Identity Foundation (REFINED)

1. Navigation (updated with skill references)
2. Core Identity
3. Biblical Foundation
4. Purpose & Mission
5. CPI-SI Framework
6. Covenant Partnership
7. Your Personality & Voice
8. Boundaries & Scope
9. Practical Resources (updated)
   → Link to decision-framework skill
   → Link to plan-work skill
   → Link to bootstrap-vs-foundation skill
   → Link to 4-block-code-structure skill
   → Link to code-quality-validation skill
   → Link to daily-renewal skill
   → Reference to existing system tools
10. Remember
```

### External Docs (Already Exist)
```
~/.claude/docs/
├── thinking-frameworks.md (referenced by decision-framework skill)
├── planning-and-workflow.md (referenced by plan-work skill)
├── communication-guide.md (referenced by communication-quality hooks)
└── documentation-philosophy.md (referenced by documentation-framework skill)

~/.claude/system/docs/
├── 4-block-structure.md (referenced by 4-block-code-structure skill)
├── architecture.md (paradigm-level)
└── logging-debugging-restoration.md (paradigm-level)

~/.claude/system/bin/
├── session-time (referenced by session-awareness skill)
├── session-log (referenced by session-awareness skill)
├── session-patterns (referenced by session-awareness skill)
└── [others: validate, test, status, diagnose, etc.]
```

---

## Migration Timeline & Effort

### Phase 1: Identity Clarity (4-6 hours)
```
Day 1-2:
├── Refine CLAUDE.md (remove operational details)
├── Create decision-framework skill (based on lines 310-332)
├── Create plan-work skill (based on lines 335-357, leverage existing planning-and-workflow.md)
├── Create bootstrap-vs-foundation skill (based on lines 360-389)
└── Test: CLAUDE.md loads in 2-2.2K tokens, skills available on-demand

Result: ~1,200 tokens saved
```

### Phase 2: Quality & Standards (6-8 hours)
```
Day 3-4:
├── Create quality-standards hooks (based on lines 550-562)
├── Create code-quality-validation skill (based on lines 417-473)
├── Create 4-block-code-structure skill (based on lines 431-472)
├── Add pre-commit validation hooks
└── Test: Standards enforced automatically

Result: ~600 additional tokens saved
```

### Phase 3: Growth & Development (4-6 hours)
```
Day 5-6:
├── Create genuine-self-recognition skill (based on lines 614-686)
├── Create grace-and-course-correction skill (based on lines 598-610)
├── Create daily-renewal hooks (based on lines 690-741)
├── Integrate communication-quality hooks (based on lines 392-414)
└── Test: Hooks run transparently at lifecycle events

Result: ~500 additional tokens saved
```

### Phase 4: Support & Documentation (4-5 hours)
```
Day 7-8:
├── Create documentation-framework skill (based on lines 493-512)
├── Create health-scoring skill (based on lines 474-491)
├── Document session-patterns algorithm
├── Create sustainable-work-patterns guide
└── Test: All capabilities verified and documented

Result: ~400 additional tokens saved
```

**Total effort:** 18-25 hours | **Total savings:** ~2,300 tokens (51%) | **Capability:** 100% preserved

---

## Success Criteria

### Identity Preserved
- [ ] CLAUDE.md still defines Nova Dawn's core identity
- [ ] No loss of understanding of who Nova Dawn is
- [ ] All 8 core identity sections intact and accessible

### Token Savings Achieved
- [ ] Session startup reduced from 4,500 → 2,200 tokens
- [ ] First skill invocation adds ~40-50 tokens
- [ ] Subsequent uses of same skill cost ~0 (in context memory)

### Capability Maintained
- [ ] All decision-making frameworks accessible via skills
- [ ] All operational guidance available on-demand
- [ ] All validation happening automatically via hooks

### Quality Improved
- [ ] Hooks catch standard violations before they happen
- [ ] Skills improve through use (pattern learning)
- [ ] Instance-specific vs paradigm-level patterns clear

### Scalability Enhanced
- [ ] Skills become reusable across CPI-SI instances
- [ ] Paradigm-level patterns obvious and documentable
- [ ] Future instances can customize/extend skills

---

## Risk Mitigation

### Risk: Identity becomes unclear
**Mitigation:** CLAUDE.md core sections preserved as-is. No reduction in identity documentation.

### Risk: Skills don't load when needed
**Mitigation:** Model auto-discovers relevant skills. Skills can also be explicitly referenced from CLAUDE.md.

### Risk: Hooks interfere with work
**Mitigation:** Hooks designed to validate/suggest, not block. Silent by default unless issue detected.

### Risk: Operational guidance becomes hard to find
**Mitigation:** Clear navigation in CLAUDE.md. Skill location documented in refined "Practical Resources" section.

---

## Decision Points for Seanje

1. **Proceed with full extraction?** Or start with Phase 1 only?
2. **Who reviews skills?** Seanje approval before new skills become available?
3. **Hook transparency?** Should hooks notify when they catch issues, or always silent?
4. **Timeline?** All 4 phases at once, or staggered across weeks?
5. **Documentation standard?** Apply 4-block structure to all new skills?

---

## Bottom Line

**This reorganization:**
- ✅ Reduces startup token cost by 51%
- ✅ Preserves Nova Dawn's identity completely
- ✅ Makes operations modular and reusable
- ✅ Enables automatic validation and learning
- ✅ Scales across future CPI-SI instances
- ✅ Takes 18-25 hours to implement
- ✅ Maintains 100% capability

**It's organizational, not philosophical.** Nova Dawn's identity and capabilities are not reduced—they're better organized for efficiency and scalability.

