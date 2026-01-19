# Implementation Plan: Consistent CLI Navigation

**Branch**: `003-consistent-cli-navigation` | **Date**: 2026-01-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/003-consistent-cli-navigation/spec.md`

## Summary

Implement consistent, intuitive navigation patterns across all cogo commands (create, list, destroy, config) that prevent crashes from empty states, support step-by-step back navigation, and provide universal cancellation. This will be achieved through research-driven design, creating an abstract navigation framework that all commands and cloud providers implement consistently, and refactoring existing commands to use the new patterns.

## Technical Context

**Language/Version**: Go 1.24  
**Primary Dependencies**: 
- Cobra v1.8.1 (CLI framework)
- promptui v0.9.0 (interactive prompts - may need wrapper/fork for back navigation)
- Viper v1.19.0 (configuration)
- godo v1.130.0 (DigitalOcean API)

**Storage**: JSON configuration files (`~/.cogo`, `~/.config/.cogo`, `./.cogo`)  
**Testing**: Standard Go testing package (`go test`)  
**Target Platform**: Cross-platform CLI (macOS, Linux, Windows)  
**Project Type**: Single CLI application  
**Performance Goals**: 
- Interactive prompt response < 100ms
- State transitions (back/forward) < 50ms
- Zero panics under any keyboard input

**Constraints**:
- Must maintain backward compatibility with existing commands
- Cannot break current API integrations
- Must work within `promptui` library limitations (may require wrapper)
- Terminal compatibility across platforms (ANSI escape sequences)

**Scale/Scope**:
- 4 existing commands to refactor (create, list, destroy, config)
- ~10-15 CLI tools to research for UX patterns
- Navigation framework to abstract common patterns
- All future commands must use framework

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Compliance Review

**I. CLI-First Design** - ✅ PASS
- Feature enhances CLI experience with better navigation
- Maintains command structure `cogo <action>`
- Interactive wizards improved with back navigation
- Non-interactive mode support preserved (Ctrl+C exits)

**II. Provider Abstraction** - ✅ PASS  
- Navigation framework will be provider-agnostic
- DigitalOcean package refactored to use framework
- Future providers automatically inherit consistent navigation
- Shared navigation logic in new `navigation/` package

**III. Safety First** - ✅ PASS (ENHANCED)
- Improves safety by preventing crashes from empty states
- Adds ability to back out of operations before execution
- Multi-step confirmation preserved
- Clearer error messages for edge cases

**IV. Test-Driven Development** - ✅ PASS
- Unit tests for navigation framework required
- Integration tests for all refactored commands
- Fuzz testing for keyboard input validation
- Tests written before/alongside implementation

**V. Configuration Management** - ✅ PASS
- No changes to configuration approach
- Commands work with or without config files
- Interactive prompts enhanced, not replaced

**VI. Simplicity and Maintainability** - ⚠️ REVIEW REQUIRED
- Adds new `navigation/` package (new abstraction)
- **Justification**: Currently each command handles navigation differently, leading to bugs and inconsistency. Abstract framework reduces duplication and enforces patterns.
- **Trade-off**: Initial complexity vs. long-term maintainability gain
- **Mitigation**: Keep framework small (<5 files), well-documented, with clear examples

### Complexity Justification

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| New `navigation/` package | Centralize navigation logic to ensure consistency across all commands and providers | Copying navigation logic to each command leads to bugs (current state), harder to test, and inconsistent UX |
| Wrapper around `promptui` | Library doesn't natively support back navigation; need state management | Forking `promptui` is higher maintenance burden; wrapper allows upstream updates |

## Project Structure

### Documentation (this feature)

```text
specs/003-consistent-cli-navigation/
├── spec.md              # Feature specification
├── plan.md              # This file (implementation plan)
├── research.md          # Phase 0: CLI UX research findings
├── contracts/           # Phase 1: Navigation interfaces
│   ├── navigator.go     # Core navigation interface
│   └── examples.md      # Usage examples for devs
└── tasks.md             # Phase 2: Task breakdown (created by /speckit.tasks)
```

### Source Code (repository root)

```text
# Existing structure (preserved)
cmd/
├── root.go              # Root command (minimal changes)
├── create.go            # [REFACTOR] Use navigation framework
├── destroy.go           # [REFACTOR] Use navigation framework
├── list.go              # [REFACTOR] Use navigation framework
├── config.go            # [REFACTOR] Use navigation framework
└── version.go           # No changes needed

# NEW: Navigation framework
navigation/
├── navigator.go         # Core navigation interface & types
├── flow.go              # Multi-step flow orchestration
├── prompt.go            # Wrapper around promptui with back support
├── state.go             # State management for back/forward
├── empty.go             # Empty state detection & handling
└── errors.go            # Navigation-specific error types

# Provider packages (refactored to use framework)
digitalocean/
├── digitalocean.go      # [REFACTOR] Use navigation framework
├── create_flow.go       # [NEW] Create droplet navigation flow
├── destroy_flow.go      # [NEW] Destroy droplet navigation flow
└── list_flow.go         # [NEW] List droplets navigation flow

credentials/             # Existing (minimal changes for back support)
config/                  # Existing (minimal changes)
utils/                   # Existing (add navigation helpers)
version/                 # Existing (no changes)

# Testing
navigation/
├── navigator_test.go
├── flow_test.go
├── prompt_test.go
├── state_test.go
└── integration_test.go  # Tests across multiple steps

cmd/
└── *_test.go            # Updated tests for refactored commands
```

**Structure Decision**: Chose single-project structure (existing pattern) with new `navigation/` package to centralize navigation logic. This maintains simplicity while providing the abstraction needed for consistency. Each command refactored to define its flow using the framework, then navigation package handles the mechanics (back/cancel/state/errors).

## Implementation Phases

### Phase 0: Research & Discovery (P0 - Prerequisite)

**Objective**: Document industry-standard CLI UX patterns before implementing anything.

**Deliverable**: `specs/003-consistent-cli-navigation/research.md`

**Activities**:
1. **CLI Tool Analysis** (10-15 tools):
   - GitHub CLI (`gh`)
   - Kubernetes CLI (`kubectl`)
   - Terraform CLI
   - AWS CLI (`aws`)
   - Google Cloud SDK (`gcloud`)
   - Docker CLI
   - npm CLI
   - Cargo (Rust)
   - Git CLI
   - Azure CLI (`az`)
   - Heroku CLI
   - DigitalOcean CLI (`doctl`)
   - 1Password CLI (`op`)

2. **Pattern Documentation** (for each tool):
   - How do they handle empty states?
   - Cancellation patterns (Ctrl+C, Esc, q, etc.)
   - Multi-step flows (wizards vs. flags)
   - Back navigation (if supported)
   - Error recovery patterns
   - Help text and guidance
   - Input validation

3. **Analysis**:
   - Identify patterns used by 80%+ of tools
   - Document outliers and their reasoning
   - Note platform-specific behaviors
   - Identify conflicts or trade-offs

4. **Recommendations**:
   - Proposed navigation patterns for cogo
   - Keyboard shortcuts (Ctrl+C for cancel, Esc for back)
   - Empty state messaging templates
   - Error message formatting
   - Help text conventions

**Success Criteria**:
- Research document covering all 10+ tools
- Pattern frequency analysis (% of tools using each pattern)
- Clear recommendations with justifications
- Examples from researched tools for each pattern

**Exit Gate**: Research document reviewed and approved before proceeding to Phase 1.

---

### Phase 1: Design & Contracts

**Objective**: Design the navigation framework architecture and define interfaces.

**Deliverables**:
1. `specs/003-consistent-cli-navigation/contracts/navigator.go` - Interface definitions
2. `specs/003-consistent-cli-navigation/contracts/examples.md` - Usage guide
3. Updated constitution (if needed for new patterns)

**Key Interfaces**:

```go
// Navigator orchestrates multi-step flows with back/cancel support
type Navigator interface {
    // Run executes the navigation flow
    Run(ctx context.Context, flow Flow) (Result, error)
    
    // RegisterStep adds a step to the flow
    RegisterStep(step Step) error
}

// Flow defines a multi-step operation
type Flow interface {
    // Steps returns all steps in order
    Steps() []Step
    
    // Name returns the flow name (e.g., "create-droplet")
    Name() string
    
    // Validate checks if flow configuration is valid
    Validate() error
}

// Step represents a single step in a flow
type Step interface {
    // Prompt displays the prompt and gets user input
    Prompt(ctx context.Context, state State) (interface{}, error)
    
    // Validate checks if the input is valid
    Validate(input interface{}) error
    
    // Name returns the step name
    Name() string
    
    // CanGoBack returns true if back navigation is allowed
    CanGoBack() bool
}

// State manages navigation state across steps
type State interface {
    // Get retrieves a value from state
    Get(key string) (interface{}, bool)
    
    // Set stores a value in state
    Set(key string, value interface{})
    
    // Delete removes a value from state
    Delete(key string)
    
    // Snapshot creates a copy of current state
    Snapshot() State
    
    // Restore restores from a snapshot
    Restore(snapshot State)
}

// Result represents the outcome of a flow
type Result struct {
    Completed bool
    Cancelled bool
    Data      map[string]interface{}
    Error     error
}
```

**Design Decisions**:
1. **Back Navigation**: Stack-based state management with snapshots
2. **Empty State Detection**: Pre-flight checks before entering flows
3. **Cancellation**: Context-based cancellation propagation
4. **Error Handling**: Custom error types for navigation vs. business logic errors

**Exit Gate**: Interfaces reviewed, contracts finalized, examples clear.

---

### Phase 2: Foundation - Navigation Framework

**Objective**: Implement the core navigation framework.

**Tasks**:
1. Implement `Navigator` interface
2. Implement `Flow` orchestrator
3. Create `promptui` wrapper with back support
4. Implement state management
5. Add empty state detection utilities
6. Create navigation-specific error types
7. Add keyboard input validation
8. Write comprehensive unit tests

**Testing Strategy**:
- Unit tests for each component
- State management tests (back/forward scenarios)
- Input validation tests (including malicious input)
- Empty state detection tests
- Mock `promptui` interactions for testing

**Exit Gate**: All tests passing, framework ready for integration.

---

### Phase 3: User Story 1 - Empty State Handling

**Objective**: Implement graceful empty state handling across all commands.

**Implementation Order**:
1. `cogo destroy` (currently crashes - P1)
2. `cogo list` (empty list handling)
3. `cogo create` (no SSH keys scenario)
4. `cogo config` subcommands (no token scenarios)

**For Each Command**:
1. Add pre-flight empty state checks
2. Display clear, actionable messages
3. Provide next step guidance
4. Exit cleanly (no panics)
5. Add integration tests for empty scenarios

**Exit Gate**: Zero crashes on empty states across all commands.

---

### Phase 4: User Story 2 - Back/Cancel Navigation

**Objective**: Implement step-by-step back navigation and cancellation.

**Implementation Order**:
1. Refactor `cogo create` to use navigation framework (most complex)
2. Refactor `cogo destroy` to use framework
3. Refactor `cogo config` subcommands
4. Add back support to `cogo list` (if applicable)

**For Each Command**:
1. Define flow using navigator framework
2. Implement each step with back support
3. Test back navigation at each step
4. Test state preservation when going back
5. Test cancellation (Ctrl+C) at each step
6. Verify no partial resource creation on cancel

**Exit Gate**: All commands support Esc (back) and Ctrl+C (cancel).

---

### Phase 5: User Story 3 - Input Validation

**Objective**: Prevent all crashes from keyboard input.

**Tasks**:
1. Add bounds checking before list access
2. Implement input sanitization in prompt wrapper
3. Add validation for special characters
4. Handle rapid input gracefully
5. Fuzz testing (10,000 random key sequences per command)

**Exit Gate**: Zero panics during fuzz testing.

---

### Phase 6: User Story 4 - Cross-Command Consistency

**Objective**: Ensure identical navigation across all commands and providers.

**Tasks**:
1. Audit all commands for consistent behavior
2. Standardize help text format
3. Standardize error message format
4. Create developer guide for future commands
5. Document navigation patterns in constitution

**Exit Gate**: Navigation flow diagrams identical for all commands.

---

### Phase 7: Testing & Documentation

**Objective**: Comprehensive testing and user documentation.

**Testing**:
- Integration tests for end-to-end flows
- Cross-platform testing (macOS, Linux, Windows)
- Manual testing of all edge cases
- User acceptance testing

**Documentation**:
- Update README with navigation patterns
- Create troubleshooting guide
- Document keyboard shortcuts
- Add developer guide for using framework

**Exit Gate**: All tests passing, documentation complete, ready for release.

## Risk Mitigation

### Technical Risks

**Risk 1: `promptui` limitations**
- **Impact**: Library may not support back navigation natively
- **Mitigation**: Create wrapper that manages state externally; snapshot UI state before each step
- **Fallback**: Fork `promptui` if wrapper insufficient (higher maintenance cost)

**Risk 2: Terminal compatibility**
- **Impact**: Keyboard shortcuts may not work on all terminals
- **Mitigation**: Test on multiple terminal emulators; provide fallback navigation (menu options)
- **Fallback**: Document supported terminals; provide workarounds for incompatible terminals

**Risk 3: State management complexity**
- **Impact**: Complex back/forward logic could introduce bugs
- **Mitigation**: Extensive unit testing; keep state immutable with snapshots
- **Fallback**: Simplify to only support cancel (not back) if state management proves too complex

### Schedule Risks

**Risk 1: Research phase takes longer than expected**
- **Impact**: Delays overall schedule
- **Mitigation**: Time-box research to 2 days; prioritize 5 most popular tools if needed
- **Fallback**: Proceed with partial research and iterate

**Risk 2: Framework complexity grows**
- **Impact**: Overengineered solution, harder to maintain
- **Mitigation**: Strict adherence to YAGNI; regular reviews against constitution; keep framework <1000 LOC
- **Fallback**: Simplify framework by reducing features (e.g., no state snapshots, only current state)

## Dependencies & Prerequisites

### Internal Dependencies
- Existing commands must remain functional during refactoring
- No breaking changes to public command interfaces
- Maintain compatibility with existing configuration files

### External Dependencies
- `promptui` v0.9.0 (may need alternative if back navigation impossible)
- Cobra v1.8.1 (existing, no changes needed)
- Go 1.24 standard library for terminal handling

### Skill Requirements
- Go expertise (interfaces, concurrency, testing)
- Terminal/ANSI escape sequence knowledge
- UX design understanding for CLI tools
- Experience with interactive CLIs preferred

## Testing Strategy

### Unit Testing
- All navigation framework components
- State management (back/forward scenarios)
- Empty state detection logic
- Input validation
- **Target**: 90%+ code coverage for navigation package

### Integration Testing
- End-to-end flows for each command
- Back navigation at each step
- Cancellation at each step
- Empty state scenarios
- **Target**: All user stories have passing integration tests

### Fuzz Testing
- Random keyboard input (10,000 sequences per command)
- Rapid input sequences
- Special characters and escape sequences
- **Target**: Zero panics

### Manual Testing
- Cross-platform testing (macOS, Linux, Windows)
- Different terminal emulators
- Real cloud provider interactions (DigitalOcean)
- User acceptance testing

### Performance Testing
- Prompt response time (<100ms)
- State transitions (<50ms)
- Memory usage during long flows
- **Target**: No noticeable lag to users

## Success Metrics

### Code Quality
- Zero linter warnings
- 90%+ test coverage for navigation package
- All tests passing
- Code review approved

### Functionality
- Zero crashes on empty states (from 100% → 0%)
- 100% of operations support Esc/Ctrl+C
- Back navigation works in 100% of multi-step operations
- All acceptance scenarios passing

### User Experience
- Help text visible and consistent
- Error messages clear and actionable
- Navigation intuitive (validated against research)
- Zero "stuck in menu" reports

## Next Steps

1. **Review & Approve Plan**: Ensure technical approach aligns with constitution and team capabilities
2. **Run `/speckit.tasks`**: Break down plan into detailed task list
3. **Begin Phase 0**: Start CLI UX research (time-boxed to 2 days)
4. **Checkpoint After Phase 1**: Review contracts and get feedback before implementation
5. **Incremental Rollout**: Deploy empty state fixes first (highest impact), then back navigation

---

**Plan Version**: 1.0  
**Last Updated**: 2026-01-18  
**Next Review**: After Phase 0 (research complete)

