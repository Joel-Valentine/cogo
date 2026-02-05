# CLI Library Research for Back Navigation

**Date**: 2026-01-18  
**Issue**: `promptui` doesn't support Esc/Q for back navigation  
**Goal**: Find a library that supports keyboard shortcuts for multi-step flows

---

## Current Situation with promptui

### What Works ✅
- Arrow keys (↑↓) for navigation
- Enter to select
- Ctrl+C to cancel
- `/` for search
- Custom templates and styling
- "← Back" as menu item (if we add handling)

### What Doesn't Work ❌
- **Esc key** - Not recognized by `promptui.Select`
- **'q' key** - Not recognized as quit
- **'b' key** - Not recognized as back
- No built-in back navigation support

### Why It Doesn't Work
From source code analysis:
- `SelectKeys` only supports: Prev, Next, PageUp, PageDown, Search
- No `Back`, `Cancel`, or custom key fields
- `Run()` either returns selection or error (Ctrl+C only)
- Esc is consumed but ignored
- Cannot customize key handlers without forking

---

## Alternative Libraries

### 1. **survey** (github.com/AlecAivazis/survey)

**Pros**:
- More mature and feature-rich than promptui
- Better keyboard handling
- Built-in support for complex forms
- Good documentation
- Active maintenance

**Cons**:
- Larger dependency
- More opinionated UI
- **Still no native Esc handling for back navigation**
- Would require custom wrapper similar to promptui

**Verdict**: ⚠️ **Same limitations as promptui**

---

### 2. **bubbletea** (github.com/charmbracelet/bubbletea)

**Pros**:
- Full TUI framework with complete key event control
- Can handle ANY keyboard input (Esc, q, b, etc.)
- Model-View-Update architecture
- Beautiful, modern UIs
- Part of Charm ecosystem (lipgloss for styling, bubbles for components)
- **Complete control over keyboard events**

**Cons**:
- **Major rewrite** - completely different paradigm
- Steeper learning curve
- Over-engineered for simple prompts
- More code to maintain
- Would need to rebuild all prompts

**Verdict**: 🎯 **Most powerful, but huge effort**

---

### 3. **huh** (github.com/charmbracelet/huh)

**Pros**:
- Built on bubbletea but simpler API
- Designed specifically for forms and multi-step flows
- **Better keyboard support than promptui**
- Modern, maintained
- Part of Charm ecosystem
- Accessible (screen reader support)

**Cons**:
- Relatively new (less battle-tested)
- Still learning community adoption
- **May still require custom handling for back nav**
- Different API from promptui

**Verdict**: 🤔 **Promising, needs testing**

---

### 4. **pterm** (github.com/pterm/pterm)

**Pros**:
- Beautiful terminal output
- Many components (spinners, progress bars, etc.)
- Good documentation

**Cons**:
- **No interactive prompts** - mainly for output
- Not suitable for our use case

**Verdict**: ❌ **Not applicable**

---

### 5. **Stay with promptui + Workarounds**

**Option A: Restore "← Back" menu items**
- ✅ Works immediately
- ✅ No library change
- ✅ Familiar to users of kubectl, terraform
- ❌ Visual clutter
- ❌ Inconsistent (menu item vs keyboard)

**Option B: Switch to numbered menus (gcloud style)**
- ✅ Works in any terminal
- ✅ Matches our research (gcloud pattern)
- ✅ Type "0" or "back" to go back
- ✅ No library limitations
- ❌ Different UX from current
- ❌ Moderate refactor

---

## Research Findings from Original Study

From `specs/003-consistent-cli-navigation/research.md`:

**gcloud (our inspiration for back navigation)**:
- Uses **numbered menus**, not arrow keys!
- "Enter number (1-5):"
- "Option 0 or 'back' to return to previous screen"
- **No Esc key mentioned**

**Other CLIs with back navigation**:
- **git interactive rebase**: numbered list, `e`/`r`/`s` commands
- **az interactive**: Full TUI shell mode
- **None use Esc for back** in our research!

**Conclusion**: We based our design on gcloud but didn't implement their actual pattern (numbered menus).

---

## Recommendations

### Short Term (Quick Fix)
**Option: Restore "← Back" items + Remove misleading help text**

```go
// Remove from help text:
Old: "↑↓: Navigate | Enter: Select | Esc: Back | Ctrl+C: Cancel"
New: "↑↓: Navigate | Enter: Select | Ctrl+C: Cancel"

// Keep in flows:
items := []string{"← Back", "Option 1", "Option 2"}
```

**Effort**: 30 minutes  
**Risk**: Low  
**UX**: Same as kubectl, terraform (they use "← Back" items too!)

---

### Medium Term (Better Solution)
**Option: Switch to numbered menus (gcloud pattern)**

```
Select Image Type:
  0) ← Back
  1) Distributions
  2) Applications  
  3) Custom

Enter selection (0-3): _
```

**Implementation**:
- Use `promptui.Prompt` (text input) instead of `Select`
- Parse number input
- "0" or "back" goes back
- Validate range
- Clear, traditional UX

**Effort**: 2-3 days  
**Risk**: Medium (changes UX)  
**UX**: Matches gcloud (our research model!)

---

### Long Term (Best but Expensive)
**Option: Switch to bubbletea + huh**

Full TUI framework with complete keyboard control.

**Effort**: 1-2 weeks  
**Risk**: High (major rewrite)  
**UX**: Modern, beautiful, fully customizable

---

## Decision Matrix

| Solution | Effort | Works Now | Future-Proof | UX Quality |
|----------|--------|-----------|--------------|------------|
| Restore "← Back" | 30min | ✅ Yes | ⚠️ OK | ⭐⭐⭐ |
| Numbered menus | 2-3 days | ✅ Yes | ✅ Great | ⭐⭐⭐⭐ |
| bubbletea/huh | 1-2 weeks | ❌ No | ✅ Excellent | ⭐⭐⭐⭐⭐ |
| Stay broken | 0 | ❌ No | ❌ Bad | ⭐ |

---

## Final Recommendation

**Immediate**: Restore "← Back" items (30 min fix)
- Remove misleading "Esc" from help text
- Keep working navigation
- Ship v3.0.1

**Next Release**: Implement numbered menus (gcloud pattern)
- More faithful to research
- Better terminal compatibility
- Traditional, proven UX
- Target v3.1.0

**Future**: Consider bubbletea migration for v4.0
- When we add more complex features
- If we want truly modern UI
- When we have time for major rewrite

---

## Examples from Other CLIs

### kubectl (uses "← back" items)
```
? Select namespace:
  ▸ ← Back
    default
    kube-system
```

### terraform (no back navigation)
```
Do you want to perform these actions?
  Enter a value: _
```

### gcloud (numbered menus)
```
Pick configuration to use:
 [1] default
 [2] production
Please enter your numeric choice: _
```

### gh (uses "← back" items in some prompts)
```
? What would you like to do?
  ▸ ← Back
    Create a pull request
    View pull requests
```

---

## Conclusion

**promptui limitation is real and cannot be fixed without major changes.**

**Best path forward**:
1. Short term: Restore "← Back" items (works, ships fast)
2. Medium term: Switch to numbered menus (matches gcloud research)
3. Long term: Consider bubbletea for future major version

The "Esc for back" dream requires either:
- Custom fork of promptui (maintenance burden)
- Switch to bubbletea (major rewrite)
- Accept that it's not how most CLIs work anyway

**gcloud doesn't use Esc either - they use numbered menus!**

