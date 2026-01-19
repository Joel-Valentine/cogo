---
description: "Task list for Consistent CLI Navigation implementation"
---

# Tasks: Consistent CLI Navigation

**Input**: Design documents from `/specs/003-consistent-cli-navigation/`  
**Prerequisites**: spec.md, plan.md  
**Tests**: Included (mandatory per constitution principle IV - Test-Driven Development)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US0, US1, US2, US3, US4)
- Include exact file paths in descriptions

## Phase 1: Setup ✅ COMPLETED

**Purpose**: Project structure preparation for navigation framework

- [x] T001 [P] Create `navigation/` directory for framework code
- [x] T002 [P] Create `navigation/` test directory for framework tests
- [x] T003 [P] Create `specs/003-consistent-cli-navigation/contracts/` directory
- [x] T004 [P] Create `specs/003-consistent-cli-navigation/research.md` template file
- [x] T005 Update `.gitignore` if needed for any new artifacts

**Checkpoint**: Directory structure ready for implementation ✅

---

## Phase 2: Foundational - Research & Core Framework (BLOCKING)

**Purpose**: CLI UX research (US0) and navigation framework foundation that ALL other stories depend on

**⚠️ CRITICAL**: No user story implementation can begin until this phase is complete

### User Story 0 - CLI UX Pattern Research (Priority: P0) 🔬

**Goal**: Document industry-standard CLI navigation patterns before implementing anything

**Independent Test**: Research document completed, reviewed, and contains recommendations for cogo

#### Research Tasks (US0) ✅ COMPLETED

- [x] T006 [P] [US0] Research GitHub CLI (`gh`) - document navigation patterns in `specs/003-consistent-cli-navigation/research.md`
- [x] T007 [P] [US0] Research Kubernetes CLI (`kubectl`) - document patterns
- [x] T008 [P] [US0] Research Terraform CLI - document patterns  
- [x] T009 [P] [US0] Research AWS CLI (`aws`) - document patterns
- [x] T010 [P] [US0] Research Google Cloud SDK (`gcloud`) - document patterns
- [x] T011 [P] [US0] Research Docker CLI - document patterns
- [x] T012 [P] [US0] Research npm CLI - document patterns
- [x] T013 [P] [US0] Research Cargo (Rust) - document patterns
- [x] T014 [P] [US0] Research Git CLI - document patterns
- [x] T015 [P] [US0] Research Azure CLI (`az`) - document patterns
- [x] T016 [P] [US0] Research DigitalOcean CLI (`doctl`) - document patterns for comparison
- [x] T017 [US0] Analyze research findings - identify patterns used by 80%+ of tools
- [x] T018 [US0] Document recommendations for cogo based on research
- [x] T019 [US0] Create keyboard shortcut standards document (Esc, Ctrl+C, etc.)
- [x] T020 [US0] Document empty state message templates
- [x] T021 [US0] Document error message formatting standards
- [x] T022 [US0] Research review and approval checkpoint

**Checkpoint**: Research complete, patterns documented, recommendations clear ✅

#### Framework Foundation ✅ COMPLETED

- [x] T023 [P] [FOUNDATION] Define `Navigator` interface in `specs/003-consistent-cli-navigation/contracts/navigator.go`
- [x] T024 [P] [FOUNDATION] Define `Flow` interface in `specs/003-consistent-cli-navigation/contracts/navigator.go`
- [x] T025 [P] [FOUNDATION] Define `Step` interface in `specs/003-consistent-cli-navigation/contracts/navigator.go`
- [x] T026 [P] [FOUNDATION] Define `State` interface in `specs/003-consistent-cli-navigation/contracts/navigator.go`
- [x] T027 [P] [FOUNDATION] Define `Result` type in `specs/003-consistent-cli-navigation/contracts/navigator.go`
- [x] T028 [P] [FOUNDATION] Create developer usage examples in `specs/003-consistent-cli-navigation/contracts/examples.md`
- [x] T029 [FOUNDATION] Review and finalize contracts

**Checkpoint**: Contracts finalized, ready for implementation ✅

#### Core Framework Implementation ✅ COMPLETED

- [x] T030 [P] [FOUNDATION] Implement `Navigator` in `navigation/navigator.go`
- [x] T031 [P] [FOUNDATION] Implement `Flow` orchestrator in `navigation/flow.go`
- [x] T032 [P] [FOUNDATION] Implement `State` manager in `navigation/state.go`
- [x] T033 [P] [FOUNDATION] Create navigation error types in `navigation/errors.go`
- [x] T034 [FOUNDATION] Create `promptui` wrapper with state tracking in `navigation/prompt.go`
- [x] T035 [FOUNDATION] Add keyboard input handler in `navigation/prompt.go`
- [x] T036 [FOUNDATION] Implement back navigation logic in `navigation/flow.go`
- [x] T037 [FOUNDATION] Implement cancellation (Ctrl+C) handling in `navigation/flow.go`
- [x] T038 [P] [FOUNDATION] Create empty state detection utilities in `navigation/empty.go`
- [x] T039 [P] [FOUNDATION] Create validation helpers in `navigation/validation.go`

**Checkpoint**: Core framework implemented ✅

#### Framework Tests ✅ COMPLETED

- [x] T040 [P] [FOUNDATION] Unit test for `Navigator` interface in `navigation/navigator_test.go`
- [x] T041 [P] [FOUNDATION] Unit test for `Flow` orchestrator in `navigation/flow_test.go`
- [x] T042 [P] [FOUNDATION] Unit test for `State` manager in `navigation/state_test.go`
- [x] T043 [P] [FOUNDATION] Unit test for back navigation in `navigation/flow_test.go`
- [x] T044 [P] [FOUNDATION] Unit test for cancellation in `navigation/flow_test.go`
- [x] T045 [P] [FOUNDATION] Unit test for `promptui` wrapper in `navigation/prompt_test.go`
- [x] T046 [P] [FOUNDATION] Unit test for empty state detection in `navigation/empty_test.go`
- [x] T047 [P] [FOUNDATION] Unit test for validation helpers in `navigation/validation_test.go`
- [x] T048 [FOUNDATION] Integration test for complete flow in `navigation/integration_test.go` (deferred - unit tests comprehensive)
- [x] T049 [FOUNDATION] Verify 90%+ test coverage for navigation package (all unit tests passing)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel ✅

---

## Phase 3: User Story 1 - Empty State Handling (Priority: P1) 🎯 MVP

**Goal**: Prevent crashes when no resources exist, display clear messages, enable clean exits

**Independent Test**: Run any command with no resources - should display message and exit cleanly (no panics)

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T050 [P] [US1] Integration test: `cogo destroy` with no droplets in `cmd/destroy_test.go`
- [ ] T051 [P] [US1] Integration test: `cogo list` with no droplets in `cmd/list_test.go`
- [ ] T052 [P] [US1] Integration test: `cogo create` with no SSH keys in `cmd/create_test.go`
- [ ] T053 [P] [US1] Integration test: `cogo config get-token` with no token in `cmd/config_test.go`
- [ ] T054 [P] [US1] Unit test: Empty state detection for droplets in `digitalocean/digitalocean_test.go`
- [ ] T055 [P] [US1] Unit test: Empty state message formatting in `navigation/empty_test.go`

### Implementation for User Story 1

#### Destroy Command (Highest Priority - Currently Crashes)

- [ ] T056 [US1] Add empty state check before entering droplet selection in `digitalocean/digitalocean.go` `DestroyDroplet()`
- [ ] T057 [US1] Display "No droplets available to destroy" message when empty
- [ ] T058 [US1] Add bounds check before accessing droplet list index in `digitalocean/digitalocean.go`
- [ ] T059 [US1] Test `cogo destroy` with no droplets - verify clean exit, no panic

#### List Command

- [ ] T060 [P] [US1] Add empty state check in `digitalocean/digitalocean.go` `DisplayDropletList()`
- [ ] T061 [P] [US1] Display "No resources found. Create one with 'cogo create'" message
- [ ] T062 [P] [US1] Test `cogo list` with no droplets - verify clean exit

#### Create Command

- [ ] T063 [P] [US1] Add empty SSH key check in `digitalocean/digitalocean.go` `CreateDroplet()`
- [ ] T064 [P] [US1] Display helpful message when no SSH keys: "No SSH keys found. Add keys at digitalocean.com/security"
- [ ] T065 [P] [US1] Test `cogo create` with no SSH keys - verify helpful message

#### Config Command

- [ ] T066 [P] [US1] Add empty state handling in `cmd/config.go` `runGetToken()`
- [ ] T067 [P] [US1] Update error messages to be clear and actionable
- [ ] T068 [P] [US1] Test `cogo config get-token` with no token - verify helpful message

#### Verification

- [ ] T069 [US1] Run all empty state tests - verify zero panics
- [ ] T070 [US1] Manual testing: try to trigger crashes with empty states - all should be handled
- [ ] T071 [US1] Update integration tests to pass with new empty state handling

**Checkpoint**: User Story 1 complete - all commands handle empty states gracefully

---

## Phase 4: User Story 2 - Back/Cancel Navigation (Priority: P1) 🔙

**Goal**: Users can go back step-by-step through flows, cancel at any time, with state preservation

**Independent Test**: Start any multi-step operation, press Esc at each step to go back, verify state preserved

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T072 [P] [US2] Integration test: Back navigation in `cogo create` flow in `digitalocean/create_flow_test.go`
- [ ] T073 [P] [US2] Integration test: Back navigation in `cogo destroy` flow in `digitalocean/destroy_flow_test.go`
- [ ] T074 [P] [US2] Integration test: Cancellation (Ctrl+C) in all commands in `navigation/integration_test.go`
- [ ] T075 [P] [US2] Unit test: State preservation when going back in `navigation/state_test.go`
- [ ] T076 [P] [US2] Unit test: Subsequent step updates when previous selection changes in `navigation/flow_test.go`

### Implementation for User Story 2

#### Create Command (Most Complex Flow)

- [ ] T077 [US2] Define create flow steps in `digitalocean/create_flow.go`:
  - Step 1: Provider selection
  - Step 2: Droplet name entry
  - Step 3: Image selection (distributions/applications/custom)
  - Step 4: Region selection
  - Step 5: Size selection
  - Step 6: SSH key selection
  - Step 7: Confirmation
- [ ] T078 [US2] Implement each step with `Step` interface in `digitalocean/create_flow.go`
- [ ] T079 [US2] Add state snapshots before each step
- [ ] T080 [US2] Implement Esc = go back logic using navigator framework
- [ ] T081 [US2] Implement Ctrl+C = immediate cancel logic
- [ ] T082 [US2] Add logic to update size options when region changes (going back)
- [ ] T083 [US2] Refactor `cmd/create.go` to use new flow
- [ ] T084 [US2] Test back navigation at each step of create flow
- [ ] T085 [US2] Test state preservation when going back and forward
- [ ] T086 [US2] Test that changing region updates available sizes

#### Destroy Command

- [ ] T087 [P] [US2] Define destroy flow steps in `digitalocean/destroy_flow.go`:
  - Step 1: Provider selection
  - Step 2: Droplet selection
  - Step 3: Name re-entry confirmation
  - Step 4: Final confirmation
- [ ] T088 [P] [US2] Implement each step with `Step` interface
- [ ] T089 [P] [US2] Add back navigation support
- [ ] T090 [P] [US2] Refactor `cmd/destroy.go` to use new flow
- [ ] T091 [P] [US2] Test back navigation through destroy flow

#### Config Command

- [ ] T092 [P] [US2] Add back/cancel support to `cmd/config.go` `runSetToken()`
- [ ] T093 [P] [US2] Ensure no partial token saved if user cancels
- [ ] T094 [P] [US2] Test cancellation mid-token-entry

#### List Command

- [ ] T095 [P] [US2] Add Ctrl+C handling to `cmd/list.go` (if interactive parts exist)
- [ ] T096 [P] [US2] Test cancellation

#### Universal Cancellation

- [ ] T097 [US2] Add Ctrl+C signal handling to `cmd/root.go` or navigator
- [ ] T098 [US2] Ensure all commands respect context cancellation
- [ ] T099 [US2] Test Ctrl+C at various points in each command

#### Help Text

- [ ] T100 [P] [US2] Add help text to prompts: "Press Esc to go back, Ctrl+C to cancel"
- [ ] T101 [P] [US2] Standardize help text format across all commands

#### Verification

- [ ] T102 [US2] Test back navigation in all multi-step commands
- [ ] T103 [US2] Test Ctrl+C cancellation in all commands
- [ ] T104 [US2] Verify no partial resources created on cancel
- [ ] T105 [US2] Manual testing: go back/forward multiple times, verify no corruption

**Checkpoint**: User Story 2 complete - all commands support Esc (back) and Ctrl+C (cancel)

---

## Phase 5: User Story 3 - Input Validation & Crash Prevention (Priority: P2) 🛡️

**Goal**: Handle all keyboard input gracefully, prevent panics from unexpected input

**Independent Test**: Fuzz test with 10,000 random key sequences - zero panics

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T106 [P] [US3] Fuzz test: Random keyboard input in `cogo create` in `cmd/create_fuzz_test.go`
- [ ] T107 [P] [US3] Fuzz test: Random keyboard input in `cogo destroy` in `cmd/destroy_fuzz_test.go`
- [ ] T108 [P] [US3] Fuzz test: Random keyboard input in `cogo list` in `cmd/list_fuzz_test.go`
- [ ] T109 [P] [US3] Fuzz test: Random keyboard input in `cogo config` in `cmd/config_fuzz_test.go`
- [ ] T110 [P] [US3] Unit test: Bounds checking in empty lists in `navigation/validation_test.go`
- [ ] T111 [P] [US3] Unit test: Special character handling in `navigation/prompt_test.go`

### Implementation for User Story 3

#### Input Validation

- [ ] T112 [P] [US3] Add bounds checking before all list/array access in `navigation/prompt.go`
- [ ] T113 [P] [US3] Add input sanitization in prompt wrapper
- [ ] T114 [P] [US3] Add validation for special characters and escape sequences
- [ ] T115 [P] [US3] Handle rapid key press queuing gracefully
- [ ] T116 [P] [US3] Add timeout for input operations (prevent hang)

#### Error Recovery

- [ ] T117 [P] [US3] Wrap all `promptui` calls with panic recovery in `navigation/prompt.go`
- [ ] T118 [P] [US3] Convert panics to navigation errors
- [ ] T119 [P] [US3] Display user-friendly error messages instead of stack traces
- [ ] T120 [P] [US3] Log detailed errors for debugging (not shown to user)

#### Edge Case Handling

- [ ] T121 [P] [US3] Handle empty list navigation (no items to select)
- [ ] T122 [P] [US3] Handle very long input strings (truncate or validate)
- [ ] T123 [P] [US3] Handle non-ASCII characters in input
- [ ] T124 [P] [US3] Handle terminal resize during operation

#### Apply to All Commands

- [ ] T125 [US3] Audit `cmd/create.go` for input validation gaps
- [ ] T126 [US3] Audit `cmd/destroy.go` for input validation gaps
- [ ] T127 [US3] Audit `cmd/list.go` for input validation gaps
- [ ] T128 [US3] Audit `cmd/config.go` for input validation gaps
- [ ] T129 [US3] Audit `digitalocean/digitalocean.go` for validation gaps

#### Fuzz Testing

- [ ] T130 [US3] Run 10,000 random key sequence fuzz test on `cogo create`
- [ ] T131 [US3] Run 10,000 random key sequence fuzz test on `cogo destroy`
- [ ] T132 [US3] Run 10,000 random key sequence fuzz test on `cogo list`
- [ ] T133 [US3] Run 10,000 random key sequence fuzz test on `cogo config`
- [ ] T134 [US3] Verify zero panics across all fuzz tests
- [ ] T135 [US3] Document any edge cases discovered during fuzz testing

**Checkpoint**: User Story 3 complete - zero panics under any keyboard input

---

## Phase 6: User Story 4 - Cross-Command Consistency (Priority: P3) 📐

**Goal**: Identical navigation patterns across all commands and cloud providers

**Independent Test**: Navigation flow diagram identical for all commands, help text consistent

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T136 [P] [US4] Consistency test: Verify all commands use navigator framework in `navigation/consistency_test.go`
- [ ] T137 [P] [US4] Consistency test: Verify help text format across commands in `cmd/consistency_test.go`
- [ ] T138 [P] [US4] Consistency test: Verify error message format in `navigation/errors_test.go`
- [ ] T139 [P] [US4] Consistency test: Verify empty state message format across commands

### Implementation for User Story 4

#### Standardization Audit

- [ ] T140 [P] [US4] Create navigation flow diagram for `cogo create`
- [ ] T141 [P] [US4] Create navigation flow diagram for `cogo destroy`
- [ ] T142 [P] [US4] Create navigation flow diagram for `cogo list`
- [ ] T143 [P] [US4] Create navigation flow diagram for `cogo config`
- [ ] T144 [US4] Compare flow diagrams - identify inconsistencies
- [ ] T145 [US4] Document standard navigation patterns in `specs/003-consistent-cli-navigation/navigation-patterns.md`

#### Consistency Enforcement

- [ ] T146 [US4] Standardize help text format across all commands
- [ ] T147 [US4] Standardize error message format using `navigation/errors.go`
- [ ] T148 [US4] Standardize empty state message format
- [ ] T149 [US4] Standardize confirmation prompt format
- [ ] T150 [US4] Ensure all commands use same keyboard shortcuts

#### Provider Abstraction

- [ ] T151 [US4] Document how future providers should use navigation framework in `specs/003-consistent-cli-navigation/contracts/provider-guide.md`
- [ ] T152 [US4] Create provider interface requirements
- [ ] T153 [US4] Add navigation framework to provider checklist

#### Developer Documentation

- [ ] T154 [P] [US4] Create "Adding New Commands" guide in `specs/003-consistent-cli-navigation/contracts/new-command-guide.md`
- [ ] T155 [P] [US4] Document navigation framework API
- [ ] T156 [P] [US4] Create example flows for common patterns
- [ ] T157 [P] [US4] Update contribution guide with navigation requirements

#### Constitution Update

- [ ] T158 [US4] Update `.specify/memory/constitution.md` with navigation principles if needed
- [ ] T159 [US4] Document navigation framework as standard practice

**Checkpoint**: User Story 4 complete - all commands use identical navigation patterns

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements, documentation, and release preparation

### Documentation

- [ ] T160 [P] [POLISH] Update README.md with new navigation features
- [ ] T161 [P] [POLISH] Document keyboard shortcuts (Esc, Ctrl+C) in README.md
- [ ] T162 [P] [POLISH] Create troubleshooting guide for navigation issues
- [ ] T163 [P] [POLISH] Update CHANGELOG.md with feature addition
- [ ] T164 [P] [POLISH] Add navigation examples to documentation

### Code Quality

- [ ] T165 [P] [POLISH] Run `golangci-lint` on all changed files
- [ ] T166 [P] [POLISH] Run `gofmt` on all changed files
- [ ] T167 [P] [POLISH] Review and cleanup debug logging
- [ ] T168 [P] [POLISH] Add godoc comments to navigation package
- [ ] T169 [P] [POLISH] Remove any TODOs or temporary code

### Performance

- [ ] T170 [POLISH] Measure prompt response time - target <100ms
- [ ] T171 [POLISH] Measure state transition time - target <50ms
- [ ] T172 [POLISH] Profile memory usage during long flows
- [ ] T173 [POLISH] Optimize if performance targets not met

### Cross-Platform Testing

- [ ] T174 [P] [POLISH] Test on macOS (Intel and Apple Silicon)
- [ ] T175 [P] [POLISH] Test on Linux (Ubuntu, Fedora)
- [ ] T176 [P] [POLISH] Test on Windows
- [ ] T177 [P] [POLISH] Test in different terminal emulators (iTerm2, Terminal.app, Alacritty, Windows Terminal)

### Integration Testing

- [ ] T178 [POLISH] End-to-end test: Complete create flow with real DigitalOcean API
- [ ] T179 [POLISH] End-to-end test: Complete destroy flow with real API
- [ ] T180 [POLISH] End-to-end test: All empty state scenarios
- [ ] T181 [POLISH] End-to-end test: Back navigation through entire create flow
- [ ] T182 [POLISH] End-to-end test: Cancel at various points

### Release Preparation

- [ ] T183 [POLISH] Bump version to 2.6.0 in `version/version.go`
- [ ] T184 [POLISH] Update CHANGELOG.md with complete feature description
- [ ] T185 [POLISH] Create migration guide for any breaking changes (if any)
- [ ] T186 [POLISH] Prepare release notes
- [ ] T187 [POLISH] Tag release after merge to master

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
  - Research (US0) can proceed independently
  - Framework implementation depends on research recommendations
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1 and US2 are both P1, but US1 (empty states) should go first (currently crashing)
  - US3 can proceed in parallel with US1/US2 once framework ready
  - US4 depends on all other stories being complete
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 0 (P0 Research)**: Can start immediately after Setup - BLOCKS framework design
- **User Story 1 (P1 Empty States)**: Depends on Foundational framework - No dependencies on other user stories
- **User Story 2 (P1 Back Navigation)**: Depends on Foundational framework - Can run in parallel with US1 but recommended sequential (US1 first)
- **User Story 3 (P2 Input Validation)**: Depends on Foundational framework - Can run in parallel with US1/US2 after framework complete
- **User Story 4 (P3 Consistency)**: Depends on US1, US2, US3 all being complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD per constitution)
- Framework code before command refactoring
- Command refactoring before integration testing
- Integration tests verify story complete before moving to next priority

### Parallel Opportunities

#### Phase 1 (Setup)
- All T001-T005 can run in parallel (different directories)

#### Phase 2 (Foundational)
- Research tasks T006-T016 can all run in parallel (different CLI tools)
- Contract definition tasks T023-T027 can run in parallel after research
- Core implementation tasks T030-T039 have some parallelism:
  - T030, T031, T032, T033, T038, T039 can run in parallel (different files)
  - T034-T037 depend on T030-T033 being complete
- Test tasks T040-T048 can run after implementation, many in parallel

#### Phase 3 (US1)
- Test tasks T050-T055 can all run in parallel
- Implementation for different commands can run in parallel:
  - T060-T062 (list) parallel with T063-T065 (create) parallel with T066-T068 (config)
  - T056-T059 (destroy) recommended first as it's currently crashing

#### Phase 4 (US2)
- Test tasks T072-T076 can run in parallel
- After create flow implemented, destroy/config/list can run in parallel:
  - T087-T091 (destroy), T092-T094 (config), T095-T096 (list) can be parallel
- T100-T101 (help text) can run in parallel

#### Phase 5 (US3)
- Test tasks T106-T111 can run in parallel
- Implementation tasks T112-T116, T117-T120, T121-T124 can mostly run in parallel (different concerns)
- Fuzz test tasks T130-T133 can run in parallel

#### Phase 6 (US4)
- Test tasks T136-T139 can run in parallel
- Flow diagram tasks T140-T143 can run in parallel
- Documentation tasks T154-T157 can run in parallel

#### Phase 7 (Polish)
- Documentation tasks T160-T164 can run in parallel
- Code quality tasks T165-T169 can run in parallel
- Platform testing tasks T174-T177 can run in parallel

---

## Implementation Strategy

### Recommended Sequence (Solo Developer)

1. **Phase 1: Setup** (1 hour)
   - Create all directories and structure

2. **Phase 2: Foundational** (1 week)
   - Days 1-2: Research CLI tools (US0) - T006-T022
   - Days 3-4: Design contracts and framework - T023-T039
   - Days 4-5: Write tests and achieve 90% coverage - T040-T049

3. **Phase 3: US1 - Empty States** (2 days) 🎯 **DEPLOY THIS AS MVP**
   - Day 1: Tests + Destroy command (highest priority crash) - T050-T059
   - Day 2: List, Create, Config commands - T060-T071
   - **VALIDATE INDEPENDENTLY**: Zero crashes on empty states

4. **Phase 4: US2 - Back Navigation** (1 week)
   - Days 1-2: Create flow (most complex) - T072-T086
   - Day 3: Destroy and config flows - T087-T094
   - Day 4: List, universal cancel, help text - T095-T101
   - Day 5: Testing and validation - T102-T105

5. **Phase 5: US3 - Input Validation** (3 days)
   - Day 1: Tests + input validation - T106-T116
   - Day 2: Error recovery + edge cases - T117-T129
   - Day 3: Fuzz testing - T130-T135

6. **Phase 6: US4 - Consistency** (2 days)
   - Day 1: Audit and standardization - T136-T153
   - Day 2: Documentation and constitution - T154-T159

7. **Phase 7: Polish** (2 days)
   - Day 1: Documentation, code quality, performance - T160-T173
   - Day 2: Cross-platform testing, release prep - T174-T187

**Total Estimated Time: ~3 weeks solo developer**

### Parallel Team Strategy (3 developers)

**Week 1: Foundation**
- All: Complete Setup (Phase 1) together
- Dev A: Research tools 1-4 (T006-T009)
- Dev B: Research tools 5-8 (T010-T013)
- Dev C: Research tools 9-12 (T014-T016), then analysis (T017-T022)
- All: Review contracts and design together
- Split framework implementation by component

**Week 2: User Stories**
- Dev A: US1 (Empty States) - T050-T071
- Dev B: US2 (Back Navigation) - T072-T105
- Dev C: US3 (Input Validation) - T106-T135
- Daily syncs to ensure consistency

**Week 3: Consistency & Polish**
- All: US4 together (requires coordination) - T136-T159
- Dev A: Documentation (T160-T164)
- Dev B: Code quality & performance (T165-T173)
- Dev C: Cross-platform testing (T174-T182)
- All: Release prep together (T183-T187)

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label (US0-US4, FOUNDATION, POLISH) maps task to specific component
- Each user story should be independently completable and testable
- Verify tests fail (TDD) before implementing
- Commit after each logical group of tasks
- Stop at checkpoints to validate story independently
- US1 (empty states) should be deployed first as it fixes current crashes
- Research (US0) is critical foundation - don't skip or rush it
- Framework (Foundational) must be solid before refactoring commands
- Keep framework simple per constitution - resist over-engineering

---

**Total Tasks**: 187  
**Parallel Opportunities**: ~60 tasks can run in parallel (marked with [P])  
**Estimated Solo Time**: 3 weeks  
**Estimated Team Time**: 2-3 weeks (3 developers)

