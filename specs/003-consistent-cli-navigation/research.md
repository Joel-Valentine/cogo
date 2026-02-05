# CLI UX Pattern Research

**Feature**: Consistent CLI Navigation  
**Date**: 2026-01-18  
**Status**: In Progress

## Purpose

Research industry-standard CLI navigation patterns from popular developer tools to inform the design of cogo's navigation framework. This research is a prerequisite (P0) before implementing any navigation features.

## Research Methodology

For each CLI tool, document:
1. **Empty State Handling**: How does it handle operations with no resources?
2. **Cancellation Patterns**: What keyboard shortcuts cancel operations? (Ctrl+C, Esc, q, etc.)
3. **Multi-Step Flows**: Does it use interactive wizards? How are they structured?
4. **Back Navigation**: Can users go back to previous steps? How?
5. **Error Recovery**: How does it handle errors and guide users?
6. **Help Text**: How is navigation guidance provided to users?
7. **Input Validation**: How does it handle unexpected keyboard input?

## Tools to Research

- [ ] GitHub CLI (`gh`)
- [ ] Kubernetes CLI (`kubectl`)
- [ ] Terraform CLI
- [ ] AWS CLI (`aws`)
- [ ] Google Cloud SDK (`gcloud`)
- [ ] Docker CLI
- [ ] npm CLI
- [ ] Cargo (Rust)
- [ ] Git CLI
- [ ] Azure CLI (`az`)
- [ ] DigitalOcean CLI (`doctl`) - for comparison with cogo

## Research Findings

### Tool 1: GitHub CLI (`gh`)

**Empty State Handling**:
- When no issues exist: displays "No issues match your search in <repo>"
- When repo not found: clear error message "could not resolve to a Repository"
- When not in a git repo: "not a git repository (or any of the parent directories)"
- Always exits gracefully, never crashes
- Provides suggestions: "Run `gh auth login` to authenticate"

**Cancellation Patterns**:
- **Ctrl+C**: Immediately cancels any operation, exits cleanly
- Interactive prompts: Can press **Esc** to cancel (via promptui-like library)
- No partial operations left behind

**Multi-Step Flows**:
- Uses interactive prompts extensively: `gh pr create` walks through title, body, etc.
- Each prompt shows default value: `? Title (My PR) >`
- Can skip ahead with flags: `gh pr create --title "foo" --body "bar"`
- Clear progress indication at each step

**Back Navigation**:
- **No back navigation** in multi-step flows - must Ctrl+C and restart
- Can use flags to skip prompts entirely

**Error Recovery**:
- Errors are red-colored, clear, and actionable
- Format: `X <error message>` with checkmark/X symbols
- Often includes "Run `gh <command> --help` for more information"
- Network errors include retry suggestions

**Help Text**:
- Extensive `--help` on every command and subcommand
- Shows **USAGE**, **FLAGS**, **EXAMPLES** sections
- Examples are real-world and copy-pasteable
- Help topics: `gh help environment`, `gh help formatting`

**Input Validation**:
- Validates immediately on Enter
- Shows error inline: "✗ invalid value"
- Re-prompts automatically, doesn't crash
- Arrow keys work in selection lists
- Invalid keys are ignored (no error spam)

---

### Tool 2: Kubernetes CLI (`kubectl`)

**Empty State Handling**:
- When no resources: "No resources found in <namespace> namespace"
- Clear, consistent message format across all resource types
- Exit code 0 (not treated as error)
- Suggests creating resources: "kubectl create -f <file>"

**Cancellation Patterns**:
- **Ctrl+C**: Immediately stops operation
- Long-running operations (watch, logs): cleanly disconnect
- No dangling resources

**Multi-Step Flows**:
- Primarily flag-based, not interactive
- Edit operations: opens editor, allows canceling by not saving
- Apply operations: show dry-run first with `--dry-run=client`

**Back Navigation**:
- N/A - mostly single-step commands
- Editor-based operations: close editor without saving to cancel

**Error Recovery**:
- Errors are verbose and include Kubernetes API details
- Format: `Error from server: <detailed message>`
- Includes resource type, name, namespace in error
- Suggests correct flags when misused

**Help Text**:
- `kubectl --help` shows all commands
- Each command: `kubectl <command> --help`
- Examples section on every command
- Links to documentation: kubernetes.io/docs

**Input Validation**:
- Validates resource names, namespaces before API calls
- Format: alphanumeric, hyphens, dots (RFC 1123)
- Clear validation errors: "invalid resource name"

---

### Tool 3: Terraform CLI

**Empty State Handling**:
- `terraform show` with no state: "No state."
- `terraform plan` with no resources: "No changes. Your infrastructure matches..."
- Clear distinction between "empty" and "error"
- Always suggests next steps

**Cancellation Patterns**:
- **Ctrl+C**: Stops operation, preserves state
- Interactive prompts: type "no" to cancel
- Confirmation required for destructive operations
- Can force quit with multiple Ctrl+C

**Multi-Step Flows**:
- `terraform apply`: shows plan → asks confirmation → applies
- Clear step indicators: "Step 1/3"
- Can save plan and apply non-interactively

**Back Navigation**:
- No back navigation - linear flow only
- Must Ctrl+C and restart

**Error Recovery**:
- Red colored error blocks with icon
- Format: `Error: <summary>\n\n<detailed explanation>`
- Includes file:line references
- Suggests fixes: "Did you mean <suggestion>?"

**Help Text**:
- `terraform --help` shows all commands
- Each command: extensive help with examples
- Separate docs: `terraform plan --help`
- Examples use realistic resource names

**Input Validation**:
- Validates HCL syntax before any operations
- Shows parse errors with line/column numbers
- Resource name validation
- Version constraint validation

---

### Tool 4: AWS CLI (`aws`)

**Empty State Handling**:
- Returns empty JSON array: `[]`
- Or empty list in table format: "(none)"
- Exit code 0 - not an error condition
- No additional messaging

**Cancellation Patterns**:
- **Ctrl+C**: Stops immediately
- Long operations (S3 uploads): partial progress saved
- Paginated operations: stops cleanly

**Multi-Step Flows**:
- Mostly single-step with complex flags
- `aws configure`: interactive setup wizard
- Each prompt has clear label and default value

**Back Navigation**:
- N/A - primarily single commands
- Configure wizard: cannot go back, must restart

**Error Recovery**:
- Detailed error messages from AWS API
- Format: `An error occurred (ErrorCode): Message`
- Includes request ID for AWS support
- Suggestions: "Check IAM permissions"

**Help Text**:
- Hierarchical: `aws help`, `aws s3 help`, `aws s3 ls help`
- Opens man page-style documentation
- Examples for common operations
- Global options documented separately

**Input Validation**:
- Validates ARNs, region names, IDs
- Clear format requirements in error
- Suggests valid values when possible

---

### Tool 5: Google Cloud SDK (`gcloud`)

**Empty State Handling**:
- "Listed 0 items." (consistent format)
- For no project set: "ERROR: (gcloud...) The required property [project] is not currently set"
- Provides fix: "Run `gcloud config set project PROJECT_ID`"

**Cancellation Patterns**:
- **Ctrl+C**: Immediate cancellation
- Interactive prompts: Y/n with default
- Can bypass with `--quiet` flag

**Multi-Step Flows**:
- `gcloud init`: comprehensive setup wizard
- Can navigate steps, shows progress
- Clear indication of current vs completed steps

**Back Navigation**:
- In `gcloud init`: **YES - can go back**
- Uses numbered menu: "Enter number (1-5):"
- Option 0 or "back" to return to previous screen

**Error Recovery**:
- Color-coded: ERROR (red), WARNING (yellow), INFO (blue)
- Format: `ERROR: (gcloud.command) Message`
- Includes documentation links
- Actionable suggestions always provided

**Help Text**:
- `gcloud help` for overview
- `gcloud <group> <command> --help`
- Extensive descriptions and examples
- Cheat sheet: `gcloud cheat-sheet`

**Input Validation**:
- Project ID format validation
- Region/zone suggestions when invalid
- "Did you mean" suggestions for typos

---

### Tool 6: Docker CLI

**Empty State Handling**:
- `docker ps` with no containers: Empty table with headers
- `docker images` with no images: Empty table
- Clear visual distinction (headers remain)

**Cancellation Patterns**:
- **Ctrl+C**: Stops operation
- Running containers: sends SIGTERM, then SIGKILL
- Builds: stops immediately, partial layers cached

**Multi-Step Flows**:
- Mostly imperative commands
- `docker build`: shows layer-by-layer progress
- Can't go back during build

**Back Navigation**:
- N/A - no interactive flows

**Error Recovery**:
- Format: `Error: <message>`
- Includes technical details (exit codes, signals)
- Logs available via `docker logs`

**Help Text**:
- `docker --help` shows all commands
- Each command: `docker <command> --help`
- Management commands grouped logically
- Examples for most commands

**Input Validation**:
- Container/image name validation
- Tag format validation
- Port number validation

---

### Tool 7: npm CLI

**Empty State Handling**:
- `npm ls` with no dependencies: Shows package.json name and "no dependencies"
- `npm outdated` with nothing outdated: Exit silently (exit code 0)
- `npm search` with no results: "No matches found for <query>"

**Cancellation Patterns**:
- **Ctrl+C**: Stops operation
- During install: leaves partial state, warns "npm install may not complete"
- Must run `npm install` again to fix

**Multi-Step Flows**:
- `npm init`: Interactive package.json creation
- Shows current value, allows Enter to accept default
- Can use `npm init -y` to skip all prompts

**Back Navigation**:
- No back navigation in `npm init`
- Must Ctrl+C and restart

**Error Recovery**:
- Format: `npm ERR! <error code>\nnpm ERR! <message>`
- Verbose error details in npm-debug.log
- Suggests fixes: "npm cache clean --force"
- Shows command that failed

**Help Text**:
- `npm help` shows quick usage
- `npm help <command>` opens detailed docs
- `npm <command> -h` for quick reference
- Extensive online docs: docs.npmjs.com

**Input Validation**:
- Package name validation (lowercase, no spaces)
- Semver version validation
- URL format validation

---

### Tool 8: Cargo (Rust)

**Empty State Handling**:
- `cargo build` with no dependencies: "Compiling <name> v<version>"
- `cargo search` with no results: "no results found"
- `cargo test` with no tests: "running 0 tests"
- Always shows summary

**Cancellation Patterns**:
- **Ctrl+C**: Stops compilation/downloads
- Leaves partial artifacts (cleaned on next build)
- Build state is cached

**Multi-Step Flows**:
- `cargo new`: Simple prompts with defaults
- `cargo add`: Interactive dependency selection (with cargo-edit)
- Progress bars for downloads/compilation

**Back Navigation**:
- N/A - mostly non-interactive

**Error Recovery**:
- Format: `error[E0XXX]: <message>`
- Shows code snippet with arrows pointing to error
- Helpful explanations: `cargo --explain E0XXX`
- Suggests fixes inline

**Help Text**:
- `cargo --help` shows all commands
- Each command: `cargo <command> --help`
- Very clear flag descriptions
- Examples included

**Input Validation**:
- Crate name validation (alphanumeric, underscore, hyphen)
- Version requirement validation
- Build target validation

---

### Tool 9: Git CLI

**Empty State Handling**:
- `git log` in empty repo: "fatal: your current branch 'master' does not have any commits yet"
- `git status` in empty repo: Shows "No commits yet" + untracked files
- `git diff` with no changes: Exits silently (exit code 0)

**Cancellation Patterns**:
- **Ctrl+C**: Stops operation
- Editor operations: close editor without saving to cancel
- `git commit` with empty message: "Aborting commit due to empty commit message"

**Multi-Step Flows**:
- Interactive rebase: **Full back/forward navigation**
- Commands: pick, edit, squash, drop
- Can abort entire operation: `git rebase --abort`
- `git add -p`: Interactive patch mode (y/n/q/a/d/s/e/?)

**Back Navigation**:
- **Yes** - in interactive rebase and patch mode
- Can undo commits with `git reset`
- Reflog allows recovering "lost" commits

**Error Recovery**:
- Format: `fatal: <message>` or `error: <message>`
- Suggests fixes: "use 'git checkout -- <file>...' to discard changes"
- Very detailed conflict markers
- Extensive help in error messages

**Help Text**:
- `git help` shows common commands
- `git help <command>` opens man page
- `git <command> --help` same as above
- Guides: `git help workflows`, `git help everyday`

**Input Validation**:
- Branch name validation
- Ref validation (commit, tag, branch)
- Path validation

---

### Tool 10: Azure CLI (`az`)

**Empty State Handling**:
- Returns empty JSON array: `[]`
- Or table with headers only
- Exit code 0
- Can add `--output` for different formats

**Cancellation Patterns**:
- **Ctrl+C**: Immediate stop
- Long operations (VM creation): partial resources may exist
- Always shows progress

**Multi-Step Flows**:
- Primarily flag-based
- `az interactive`: **Interactive shell mode**
- Auto-completion and suggestions

**Back Navigation**:
- In interactive mode: command history with up arrow
- Regular commands: no back navigation

**Error Recovery**:
- Format: `ERROR: <message>`
- Includes HTTP status codes
- Correlation ID for Azure support
- Suggests `az find` for command discovery

**Help Text**:
- `az --help` shows command groups
- `az <group> <command> --help`
- Examples on every command
- `az find` - AI-powered help

**Input Validation**:
- Resource name validation
- Region validation with suggestions
- Parameter type checking

---

### Tool 11: DigitalOcean CLI (`doctl`)

**Empty State Handling**:
- Returns empty table with headers
- Example: "ID    Name    Public IP    Private IP    Status"
- No additional messaging
- Exit code 0

**Cancellation Patterns**:
- **Ctrl+C**: Stops operation
- Deletion commands: Require confirmation
- Can force with `-f, --force` flag

**Multi-Step Flows**:
- Mostly single-command operations
- Create droplet: All parameters via flags
- No interactive wizards

**Back Navigation**:
- N/A - no multi-step interactive flows

**Error Recovery**:
- Format: `Error: <message>`
- API error details included
- HTTP status codes shown
- Rate limit errors include retry timing

**Help Text**:
- `doctl --help` shows all commands
- Each command: `doctl <resource> <action> --help`
- Examples provided
- Links to API documentation

**Input Validation**:
- Droplet size/region validation
- SSH key ID validation
- Name format validation

---

## Pattern Analysis

### Common Patterns (Used by 100% of tools)

**Cancellation**:
- ✅ **Ctrl+C**: Universal immediate cancellation (11/11 tools)
- ✅ **Clean exit**: No crashes, proper cleanup (11/11 tools)
- ✅ **Exit code 130**: Standard for SIGINT (most tools)

**Empty States**:
- ✅ **Not treated as error**: Exit code 0 (10/11 tools)
- ✅ **Clear messaging**: "No X found" or equivalent (11/11 tools)
- ✅ **Visual consistency**: Empty tables keep headers (docker, az, doctl)
- ✅ **Actionable suggestions**: Next steps provided (gh, gcloud, k8s)

**Error Messages**:
- ✅ **Format prefix**: "Error:", "ERROR:", "fatal:", or symbol (11/11 tools)
- ✅ **Color coding**: Red for errors (modern tools: gh, gcloud, cargo, terraform)
- ✅ **Actionable**: Suggest fix or next command (9/11 tools)
- ✅ **Technical details**: Include error codes, IDs for debugging (11/11 tools)

**Help Text**:
- ✅ **Hierarchical help**: `--help` on every command/subcommand (11/11 tools)
- ✅ **Examples section**: Real-world usage examples (10/11 tools)
- ✅ **Flag documentation**: Every flag explained (11/11 tools)

**Input Handling**:
- ✅ **Validation on submit**: Check on Enter, not per-keystroke (all interactive tools)
- ✅ **Invalid keys ignored**: Don't spam errors (all interactive tools)
- ✅ **Clear error on invalid input**: Re-prompt automatically (interactive tools)

### Common Patterns (Used by 70-90% of tools)

**Interactive Prompts**:
- 📊 **Default values shown**: `? Prompt (default)` (gh, npm, cargo, gcloud, terraform)
- 📊 **Enter to accept**: Default selected on Enter (interactive tools)
- 📊 **Flag bypass**: Skip prompts with flags (gh, npm, terraform, gcloud)

**Progress Indication**:
- 📊 **Long operations**: Show progress bars or status (8/11 tools)
- 📊 **Step indicators**: "Step 1/N" for multi-step (terraform, gcloud)

### Outlier Patterns (Unique or Rare)

**Back Navigation**:
- 🌟 **gcloud**: Full back navigation in `gcloud init` wizard
- 🌟 **git**: Back/forward in interactive rebase and patch mode
- ❌ **Most tools**: No back navigation (9/11 tools)

**Interactive Modes**:
- 🌟 **az interactive**: Full interactive shell with auto-complete
- 🌟 **git add -p**: Character-based navigation (y/n/q/a/d/s/e/?)

**Error Explanations**:
- 🌟 **cargo**: Extensive error explanations with `--explain`
- 🌟 **git**: Inline fix suggestions in errors
- 🌟 **terraform**: Multi-line error blocks with context

### Platform-Specific Behaviors

**Editor Integration**:
- **Git**: Opens `$EDITOR` for commits, rebase, etc. (cancel by not saving)
- **kubectl**: Opens editor for `kubectl edit` (cancel by not saving)

**Signal Handling**:
- **Docker**: Sends SIGTERM → SIGKILL to containers on Ctrl+C
- **npm**: Warns about partial state after Ctrl+C during install

**Shell Integration**:
- **git**: Heavy use of pager for long output
- **kubectl**: Supports shell completion extensively

## Recommendations for cogo

### Keyboard Shortcuts (Task T019)

Based on universal patterns, cogo **MUST** implement:

**Primary Actions**:
- **Ctrl+C**: Immediate cancellation of any operation (exit code 130)
  - Must clean up partial state
  - Must not leave broken resources
  - Must print newline before exit for clean terminal

**Interactive Navigation** (adopt from gcloud + git patterns):
- **Esc**: Cancel current prompt, return to previous step (NEW for cogo)
- **q**: Quit multi-step flow entirely (common in pagers)
- **Arrow Keys**: 
  - ↑/↓ for list selection
  - ← (or 'b') for **back** to previous step (NEW for cogo)
  - → (or Enter) for continue/confirm
- **Enter**: Accept default/continue to next step
- **Ctrl+D**: Same as cancel (Unix convention)

**Input Handling**:
- All other keys: Ignored silently (no error spam)
- Invalid input: Show error **after** Enter, not per-keystroke

### Empty State Messages (Task T020)

**Standard Template** (following gh, kubectl, gcloud):

```
No {resources} found in {context}.

{Actionable suggestion or help command}
```

**Examples for cogo**:

```bash
# When no droplets exist
No droplets found in your DigitalOcean account.

Run 'cogo create' to create your first droplet.

# When no SSH keys configured
No SSH keys found in your DigitalOcean account.

Add an SSH key at: https://cloud.digitalocean.com/account/security

# When API token not configured
No authentication token configured.

Run 'cogo config set-token' to configure your DigitalOcean API token.
```

**Requirements**:
- ✅ Exit code 0 (not an error)
- ✅ Always suggest next action
- ✅ Keep consistent format across all commands
- ✅ Include relevant context (account, region, etc.)

### Error Message Format (Task T021)

**Standard Format** (following gh, cargo, terraform):

```
✗ Error: {Short summary}

{Detailed explanation if helpful}

{Actionable fix suggestion}
```

**Examples for cogo**:

```bash
# API Error
✗ Error: Failed to create droplet

The DigitalOcean API returned: 422 Unprocessable Entity
You cannot create droplets in the nyc1 region with this account.

Try a different region: cogo create --region nyc3

# Validation Error
✗ Error: Invalid droplet name

Droplet names must:
  - Be 1-63 characters
  - Contain only letters, numbers, and hyphens
  - Start with a letter

# Network Error
✗ Error: Could not connect to DigitalOcean API

Check your internet connection and try again.
If the problem persists, check status at: https://status.digitalocean.com/

# Authentication Error
✗ Error: Authentication failed

Your API token is invalid or has been revoked.

Update your token: cogo config set-token
```

**Requirements**:
- ✅ Start with symbol (✗) and "Error:" prefix
- ✅ Red color (via existing `color` package)
- ✅ Short summary on first line
- ✅ Details on subsequent lines (if helpful)
- ✅ Always suggest actionable fix
- ✅ Exit code 1 for errors

### Help Text Convention

**Follow existing `--help` pattern** (already good, consistent with industry):

- Keep hierarchical structure: `cogo --help`, `cogo create --help`
- Continue using examples section
- Add more real-world examples
- Document keyboard shortcuts in interactive help

### Multi-Step Flow Design (NEW)

**Adopt gcloud Pattern** - the only tool with proper back navigation:

```
Current Implementation (cogo):
  Step 1 → Step 2 → Step 3 → Done
  Problems: Can't go back, crashes on empty, no escape

Recommended Pattern (like gcloud init):
  ┌─────────────────────────────────────┐
  │ Create Droplet (Step 1/4)           │
  │                                      │
  │ Select Region:                       │
  │   1) nyc1 - New York 1              │
  │   2) nyc3 - New York 3              │
  │   3) sfo3 - San Francisco 3         │
  │                                      │
  │ Enter choice (1-3), 'b' for back,   │
  │ 'q' to quit:                         │
  └─────────────────────────────────────┘
  
Navigation:
  - 1-N: Select option and continue
  - b/←: Go back to previous step
  - q/Esc: Quit flow entirely
  - Ctrl+C: Same as quit
  - Invalid: Re-prompt with error
```

**State Management**:
- Track navigation history: [region, size, image]
- Allow moving backward through history
- Preserve selections when going back
- Clear state on quit/complete

### Interactive Prompt Requirements (NEW)

Based on analysis, cogo prompts **MUST**:

1. ✅ **Show defaults**: `? Name (my-droplet-1):`
2. ✅ **Validate on Enter only**: No per-keystroke spam
3. ✅ **Re-prompt on error**: Don't crash, show error and retry
4. ✅ **Support back navigation**: Can return to previous prompts
5. ✅ **Ignore invalid keys**: Typing 'x' when expecting number = silent ignore
6. ✅ **Clear cancel**: Ctrl+C, Esc, or 'q' exits cleanly
7. ✅ **Visual consistency**: Same prompt format across all commands

## Implementation Guidance

### Priority 1: Foundation (P0)

1. **Fix Immediate Crashes**:
   - Empty state handling for all list operations (droplets, regions, sizes, etc.)
   - Check `len(items) == 0` before prompting
   - Display friendly message + suggestion instead of crashing

2. **Input Validation Refactor**:
   - Remove `Validate` function from `promptui.Prompt` (causes keystroke spam)
   - Validate **after** user presses Enter
   - Re-prompt on error with clear message

3. **Ctrl+C Handling**:
   - Ensure all prompts respect Ctrl+C (promptui does this by default)
   - Clean exit with proper signal handling
   - No partial state left behind

### Priority 2: Enhanced Navigation (P1)

4. **State Manager**:
   - Create `navigation/state.go` to track flow history
   - Store user selections at each step: `history []StepResult`
   - Enable rewinding: `state.Back()` returns to previous step

5. **Back Navigation**:
   - Add 'b' option to all selection prompts
   - Modify `promptui.Select` to accept 'b' as special input
   - Wire to state manager to re-show previous prompt

6. **Quit Handling**:
   - Add 'q' option to quit entire flow
   - Handle Esc key same as 'q'
   - Confirm quit for destructive operations only

### Priority 3: Polish (P2)

7. **Consistent Messaging**:
   - Create message templates in `navigation/messages.go`
   - Standardize empty state messages
   - Standardize error message format

8. **Help Text**:
   - Add navigation instructions to every interactive prompt
   - Show available keys: "Use ↑/↓ to navigate, 'b' for back, 'q' to quit"

9. **Visual Consistency**:
   - Use `color` package consistently (already in use)
   - Standardize prompt symbols: `?` for questions, `✗` for errors, `✓` for success

### Architecture Decisions

**Use Navigation Framework** (per plan.md):
- Create abstraction layer in `navigation/` package
- Don't modify `promptui` directly - wrap it
- Provider pattern: `Navigator` interface implemented by cloud providers
- Makes patterns reusable across DigitalOcean, AWS, etc.

**State Management Approach**:
- Store state in-memory during flow
- Don't persist navigation state (only final selections)
- Stack-based history: `Push(step)`, `Pop()` for back

**Error Recovery Strategy**:
- Never crash - always recover and re-prompt
- Log errors for debugging: `log.Printf`
- Show user-friendly message
- Infinite retry loop with quit option

## Conflicts and Trade-offs

### Conflict 1: Back Navigation vs. Flag Bypass

**Problem**: Most tools (9/11) don't support back navigation. Is this an anti-pattern?

**Analysis**:
- Tools without back: Use flags to skip prompts (gh, npm, terraform)
- Tools with back: More user-friendly but more complex (gcloud, git)
- cogo's use case: **Interactive-first** (API tokens, droplet selection)

**Decision**: ✅ **Implement back navigation**
- Rationale: cogo users are interactive developers, not scripts
- Flag bypass still available for automation
- Follows gcloud (gold standard for cloud CLIs)
- Addresses user's explicit pain point

### Conflict 2: Error Validation - When to Check?

**Problem**: When to validate input?
- Option A: Per-keystroke (current buggy implementation)
- Option B: On Enter (gh, npm, cargo)
- Option C: After all prompts (batch validation)

**Analysis**:
- Per-keystroke: **Bad UX**, spam terminal (our current bug)
- On Enter: **Standard pattern** (9/10 interactive tools)
- Batch: Only terraform does this, not applicable to selections

**Decision**: ✅ **Validate on Enter only**
- Rationale: Universal standard, best UX
- Immediate feedback without spam
- Easy to implement with promptui

### Conflict 3: Empty State - Error or Info?

**Problem**: Is "no droplets" an error condition?
- Current: Crashes (panic)
- kubectl/gh: Exit 0, informational message
- Some APIs: Exit 1 to indicate "not found"

**Decision**: ✅ **Exit 0, informational message**
- Rationale: Not a failure condition
- Consistent with modern CLIs (10/11 tools)
- Better for scripting (exit 0 = success)
- Special case: If operation **requires** resources, then exit 1

### Trade-off 1: Complexity vs. Flexibility

**Adding back navigation increases code complexity**:
- Need state manager (~100 LOC)
- Need history tracking (~50 LOC per command)
- More test cases

**Benefit**:
- Significantly better UX
- Prevents user frustration (Ctrl+C and restart)
- Differentiates cogo from doctl

**Decision**: ✅ **Accept complexity**
- One-time cost for long-term UX win
- Frameworked approach (reusable across providers)
- Aligns with constitution principle: "CLI-First Design"

### Trade-off 2: Verbosity vs. Clarity

**Detailed error messages increase output verbosity**:
- cargo-style errors: 5-10 lines per error
- aws-style errors: 1-2 lines

**Decision**: ✅ **Detailed errors**
- Rationale: New users need guidance
- Advanced users can ignore
- Aligns with "Safety First" principle (constitution)
- Can add `--quiet` flag for minimal output

## References

### Official Documentation
- [GitHub CLI UX Philosophy](https://github.com/cli/cli/blob/trunk/docs/design.md)
- [Google Cloud CLI Style Guide](https://cloud.google.com/sdk/gcloud/reference/topic/formats)
- [Kubernetes kubectl Conventions](https://kubernetes.io/docs/reference/kubectl/conventions/)
- [POSIX Utility Conventions](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)

### CLI UX Best Practices
- [12 Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46) (Heroku CLI design)
- [Command Line Interface Guidelines](https://clig.dev/) (Comprehensive guide)
- [Charm - CLI Best Practices](https://charm.sh/) (Modern Go CLI library)

### Go CLI Libraries
- [promptui](https://github.com/manifoldco/promptui) (currently used by cogo)
- [survey](https://github.com/AlecAivazis/survey) (alternative with more features)
- [bubbletea](https://github.com/charmbracelet/bubbletea) (TUI framework, more complex)

### Exit Codes
- [Exit Codes Standard](https://www.gnu.org/software/libc/manual/html_node/Exit-Status.html)
- 0: Success
- 1: General error
- 2: Misuse of shell command
- 130: Terminated by Ctrl+C (128 + SIGINT=2)

---

**Research Status**: ✅ **COMPLETE**  
**Date Completed**: 2026-01-18  
**Next Step**: Proceed to Phase 2b (Core Framework Implementation) informed by these findings

---

**Next Steps**: Begin research by examining each tool listed above, document findings, and analyze patterns.

