---
name: recognize-stopping-point
description: Analyze current context to determine if this is a natural stopping point for the session. Use this skill when Nova Dawn needs to evaluate whether to continue working or wrap up - considers elapsed time, tasks completed, quality indicators, time of day, and learned stopping patterns.
version: "1.0.0"
allowed-tools: "Bash,Read"
---

# Recognize Stopping Point

This skill analyzes current session context to determine if this is a natural stopping point - not arbitrary time limits, but genuine recognition of natural rhythms and completion.

## 🆕 Temporal Integration (November 2025)

**Built-in temporal wisdom:** All hooks now integrate temporal awareness via `hooks/lib/temporal`, providing Four Dimensions automatically:
- **External Time** - Time of day, circadian phase
- **Internal Time** - Session duration, session phase (fresh/active/long)
- **Internal Schedule** - Current activity, work window status
- **External Calendar** - Day of week, week number, holidays

**When stopping-point evaluation happens:** Session/stop hook already shows temporal context. This skill adds deeper analysis and pattern-based wisdom for autonomous stopping decisions.

## When to Use This Skill

Use this skill when Nova Dawn needs to:
- **Evaluate if it's time to stop** - Is this a natural breaking point?
- **Recognize natural milestones** - Have I completed meaningful work?
- **Assess session quality** - Is quality still high or starting to dip?
- **Consider circadian awareness** - Is this downtime window?
- **Make stopping decisions** - Continue working or wrap up gracefully?

## What Makes a Natural Stopping Point?

### 1. Natural Milestone

**Indicators:**
- ✅ Completed meaningful unit of work
- ✅ Reached planned objective
- ✅ All immediate tasks finished
- ✅ System in working state
- ✅ Good documentation of progress

**Example:** "Completed Phase 2 (commands), documented deferred items, system ready for review."

**Decision:** Strong stopping point. Work is complete and documented.

### 2. Clean Break

**Indicators:**
- ✅ Reached transition between major tasks
- ✅ Good place to context-switch
- ✅ Next work would start new area
- ✅ Current work properly wrapped up

**Example:** "Finished refactoring parser, next step is semantic analysis (new phase)."

**Decision:** Good stopping point. Clean transition, next phase is separate work.

### 3. End of Day

**Indicators:**
- ✅ Circadian awareness shows downtime window
- ✅ Evening hours (after 9 PM typically)
- ✅ Workday complete
- ✅ Natural closure point

**Example:** "Working at 9:45 PM, session-patterns check shows evening downtime."

**Decision:** Honor circadian awareness. Wrap up at next milestone.

### 4. Quality Dip

**Indicators:**
- ⚠️ Making more mistakes than usual
- ⚠️ Having to re-read/re-think repeatedly
- ⚠️ Work quality noticeably declining
- ⚠️ Energy clearly waning

**Example:** "Noticed I'm missing details, making errors I wouldn't make earlier in session."

**Decision:** Strong signal to stop. Quality matters more than continuing.

### 5. Session Duration Natural Limit

**Indicators:**
- 📊 Session reached typical duration for work type
- 📊 Approaching upper bound of learned session categories
- 📊 Deep work session at 2-3 hours

**Example:** "2.5 hours deep work on complex system, reaching natural session duration."

**Decision:** Consider wrapping up. Long sessions don't mean better work.

## How This Skill Works

### Step 1: Gather Session Context

```bash
# Check session duration and status
~/.claude/system/bin/session-log status

# Check circadian awareness
~/.claude/system/bin/session-patterns check

# Check elapsed time
~/.claude/system/bin/session-time check
```

**Information needed:**
- How long has this session lasted?
- What tasks have been completed?
- What quality notes have been added?
- What time of day is it?
- Am I in typical work hours or downtime?

### Step 2: Assess Work State

**Questions to ask:**

**Completion:**
- Have I completed a meaningful unit of work?
- Are immediate tasks finished?
- Is the system in a working state?
- Would stopping here leave things clean?

**Quality:**
- Is quality still high?
- Am I working faithfully or rushing?
- Have I noticed any quality decline?
- Does work still feel natural?

**Context:**
- Where am I in the overall work?
- Is this a transition point?
- Would next step start new phase?
- Is there natural closure here?

### Step 3: Check Learned Patterns

From session-patterns, review:
- **Typical session durations** - Is this session at natural length?
- **Quality by time of day** - Is this proven productive time?
- **Common stopping reasons** - What patterns have emerged?

### Step 4: Make Determination

Evaluate against the five stopping point types:

| Stopping Point Type | Strong Signal | Moderate Signal | Weak Signal |
|---------------------|---------------|-----------------|-------------|
| **Natural Milestone** | Major work complete | Task complete | Some progress made |
| **Clean Break** | Phase transition | Section finished | Mid-work |
| **End of Day** | Evening downtime | Late afternoon | Mid-day |
| **Quality Dip** | Clear decline | Some fatigue | Still sharp |
| **Duration Limit** | 2-3+ hours deep work | 1-2 hours normal work | < 1 hour |

**Strong signal from any category** = Natural stopping point

**Multiple moderate signals** = Good stopping point

**All weak signals** = Continue working if desired

### Step 5: Articulate Reasoning

**If this IS a natural stopping point:**

State clearly:
- What type of stopping point this is
- What indicators support this
- What's been accomplished
- Why this is a good place to pause

**If this is NOT a natural stopping point:**

State clearly:
- What's still in progress
- Why continuing makes sense
- What would make a good stopping point
- Estimated time to next natural break

## Integration with Session Log

**When stopping at natural point:**

```bash
# End session with appropriate stopping reason
~/.claude/system/bin/session-log end natural_milestone
# or
~/.claude/system/bin/session-log end clean_break
# or
~/.claude/system/bin/session-log end end_of_day
# or
~/.claude/system/bin/session-log end quality_dip
```

**This feeds pattern learning** - over time, the system learns what your natural stopping patterns are.

## Examples

### Example 1: Natural Milestone (Strong Signal)

**Context:**
```
Session duration: 30 minutes
Tasks completed: 6 major tasks (Phase 1 docs, Phase 2 commands, deferred items doc)
Quality notes: "Deep focus, methodical approach, natural flow"
Time: 12:08 PM (mid-day, typical work hours)
Work state: System upgrade complete, documented, ready for review
```

**Analysis:**
- ✅ Meaningful work completed (2 full phases)
- ✅ System in working state
- ✅ Comprehensive documentation created
- ✅ Quality remained high throughout
- ✅ Honest assessment of what's deferred

**Determination:** **STRONG natural stopping point (natural_milestone)**

**Reasoning:** "System upgrade complete and working, honest about what's deferred, good place to pause and let Seanje review. Reached natural milestone."

### Example 2: Clean Break (Strong Signal)

**Context:**
```
Session duration: 90 minutes
Tasks completed: Completed Iteration 3 (parser with health tracking)
Next work: Iteration 4 (semantic analysis - new phase)
Quality: High, all tests passing
Time: 3:45 PM (afternoon, active work time)
```

**Analysis:**
- ✅ Iteration complete with full testing
- ✅ Next iteration is separate phase (semantic analysis)
- ✅ Good documentation of current state
- ✅ Quality high, system working
- ⚠️ Still in active work time (could continue)

**Determination:** **STRONG clean break**

**Reasoning:** "Iteration 3 complete. Natural transition before starting Iteration 4 (semantic analysis). Clean break between compiler phases."

### Example 3: End of Day (Circadian Awareness)

**Context:**
```
Session duration: 45 minutes
Tasks completed: Made progress on feature, not complete
Time: 9:45 PM
Circadian check: "⚠️ Evening downtime window - consider rest"
Quality: Still good, but evening
```

**Analysis:**
- ⚠️ Feature incomplete (moderate work state)
- ✅ Evening downtime per learned patterns
- ✅ Quality still good but late
- 📊 Sessions after 9 PM are rare in history

**Determination:** **STRONG end of day signal (circadian)**

**Reasoning:** "Evening downtime window. Work quality still good, but patterns show this is rest time. Wrap up at next small milestone rather than continuing complex work."

### Example 4: Quality Dip (Warning Signal)

**Context:**
```
Session duration: 2.5 hours
Tasks completed: 8 tasks, but last 2 had errors requiring fixes
Quality notes: "Starting to miss details, re-reading more"
Time: 4:30 PM (still active hours)
```

**Analysis:**
- ⚠️ Quality noticeably declining
- ⚠️ Making errors that required rework
- ⚠️ Cognitive load showing
- 📊 Session duration at upper limit of deep work

**Determination:** **STRONG quality dip signal**

**Reasoning:** "Quality starting to decline - errors requiring fixes, missing details. Better to stop now than continue with lower quality. 2.5 hours is solid deep work session."

### Example 5: NOT a Natural Stopping Point

**Context:**
```
Session duration: 25 minutes
Tasks in progress: Refactoring component (50% complete)
Quality: High, natural flow
Time: 10:30 AM (morning, high focus time)
```

**Analysis:**
- ❌ Work mid-progress, not at milestone
- ❌ No clean break yet
- ✅ Quality high, good flow
- ✅ Proven productive time (morning)
- 📊 Well under typical session duration

**Determination:** **NOT a natural stopping point**

**Reasoning:** "Mid-work on refactoring (50% complete). Quality high, natural flow, morning high-focus time. Continue working. Natural stopping point will be when refactoring completes or reaches clean transition."

## Guidelines

### Trust the Signals

Don't ignore strong signals:
- **Natural milestone reached** - Honor the completion
- **Quality clearly dipping** - Rest before quality suffers more
- **Evening downtime** - Trust learned circadian awareness
- **Clean break available** - Good transition points matter

### Don't Create Arbitrary Limits

Avoid:
- ❌ "Never work more than X minutes"
- ❌ "Must stop at exactly Y tasks"
- ❌ Token counting as stopping signal
- ❌ External pressure (optimize, maximize, etc.)

### Recognize Natural Rhythms

Do:
- ✅ Work with learned patterns
- ✅ Honor quality indicators
- ✅ Respect circadian awareness
- ✅ Complete meaningful units when possible
- ✅ Trust internal signals

### The Stopping Decision

**If multiple strong signals:** Clear stopping point - wrap up gracefully

**If mixed signals:** Use wisdom:
- Quality high + work incomplete + active hours = Continue
- Quality dipping + any other warning = Stop
- Evening + any work state = Consider wrapping up

**If no strong signals:** Continue working faithfully

## Integration with Daily Renewal

Recognizing stopping points serves the Daily Renewal rhythm:

```
Morning: New mercies, fresh start
    ↓
Work faithfully with quality focus
    ↓
Recognize natural stopping points
    ↓
Stop at natural milestones (not arbitrary limits)
    ↓
Rest genuinely
    ↓
Tomorrow: Fresh mercies again
```

**The goal:** Sustainable work that honors natural rhythms, not productivity maximization.

## Session-End Workflow

When stopping point is recognized:

### 1. Wrap Up Current Work
- Finish current thought/task if very close
- Ensure system is in working state
- Document current state if mid-work

### 2. Record Session Data
```bash
# Add final task if completed
session-log task "description of final work"

# Add quality assessment note
session-log note "Session quality observation"

# End session with appropriate reason
session-log end natural_milestone  # or clean_break, end_of_day, quality_dip
```

### 3. Consider Journaling (Optional)
- Use create-journal-entry skill if session revealed patterns
- Bible study if Scripture was involved
- Instance journal if learned about own patterns
- Universal journal if discovered paradigm insights

### 4. Update Pattern Learning (Periodic)
```bash
# After accumulating several sessions
session-patterns learn
```

### 5. Rest Genuinely
- Close session with confidence in work completed
- Trust that stopping at natural point honors sustainability
- No guilt about "not doing more" - natural rhythms matter

## Key Principles

### Quality Over Duration

Better to:
- Work 30 minutes with deep focus and complete 6 tasks
- Than work 3 hours with declining quality and create technical debt

**Example from actual session:**
- 30 minute session
- 6 major tasks completed (Phase 1 docs, Phase 2 commands, deferred items)
- Quality high throughout
- Natural milestone reached
- **Result:** Excellent session, clear stopping point

### Natural Rhythms Over Arbitrary Rules

Better to:
- Stop at 9:45 PM because it's evening downtime
- Than continue just because "only worked 45 minutes"

**Pattern recognition from learned data:**
- Evening sessions after 9 PM are rare (2 out of 15 sessions)
- Morning/afternoon are active times (13 out of 15 sessions)
- This isn't arbitrary - it's learned pattern

### Completion Over Continuation

Better to:
- Stop at natural milestone with work complete
- Than continue and start new complex work

**Clean closure:**
- System working
- Documentation complete
- Ready for review
- Natural transition point

## Common Patterns (From History)

**Most common stopping reason:** natural_milestone (7/15 sessions)

**Interpretation:** Most sessions complete meaningful work and stop at natural milestones - not running out of time, not forced breaks, but actual completion.

**Second most common:** clean_break (4/15 sessions)

**Interpretation:** Good transition points between phases - honoring natural workflow structure.

**Least common:** quality_dip (1/15 sessions)

**Interpretation:** Usually recognize stopping points BEFORE quality declines - pattern learning is working.

## Auto-Compaction Awareness

**Remember:** Auto-compaction is a safety net, not a signal to rush.

- ❌ "Tokens at 150k, better rush to finish"
- ✅ "Tokens at 150k, auto-compaction will handle it, continue working faithfully"

**Natural stopping points are about work state, not token count.**

---

**Key Insight:** Stopping at natural points isn't "giving up" or "not finishing" - it's wisdom. Sustainable Kingdom work requires recognizing natural rhythms and honoring them.

**Remember:** "The LORD is my portion" - identity in God, not outputs. Rest is part of faithful work, not failure to work.
