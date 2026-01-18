# Tasks: Modernize DigitalOcean API Integration

**Input**: Design documents from `/specs/001-modernize-digitalocean-api/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ⏳, data-model.md ⏳, contracts/ ⏳

**Tests**: Tests are required per Cogo constitution (Principle IV: Test-Driven Development)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Research & Discovery

**Purpose**: Understand API changes and prepare for implementation

- [ ] T001 [P] Research DigitalOcean API v2 changelog from 2020-2026 and document breaking changes in `specs/001-modernize-digitalocean-api/research.md`
- [ ] T002 [P] Identify latest stable godo SDK version and review migration guide
- [ ] T003 Document deprecated API methods and their replacements in `research.md`
- [ ] T004 [P] List new required parameters for droplet creation endpoint
- [ ] T005 [P] Review pagination changes in godo SDK (Links header handling)
- [ ] T006 [P] Document rate limit changes and retry strategies
- [ ] T007 Create list of affected functions in `digitalocean/digitalocean.go` with line numbers

**Deliverable**: `research.md` complete with migration notes and affected code locations

---

## Phase 2: Design & API Contracts

**Purpose**: Define updated API contracts and data structures

- [ ] T008 [P] Document POST `/v2/droplets` endpoint in `specs/001-modernize-digitalocean-api/contracts/droplets-api.md`
- [ ] T009 [P] Document GET `/v2/droplets` endpoint with pagination in `contracts/droplets-api.md`
- [ ] T010 [P] Document DELETE `/v2/droplets/:id` endpoint in `contracts/droplets-api.md`
- [ ] T011 [P] Document GET `/v2/regions` endpoint in `contracts/regions-api.md`
- [ ] T012 [P] Document GET `/v2/images` endpoint variants in `contracts/images-api.md`
- [ ] T013 [P] Document GET `/v2/sizes` endpoint in `contracts/regions-api.md`
- [ ] T014 [P] Document GET `/v2/account/keys` endpoint in `contracts/droplets-api.md`
- [ ] T015 Create `specs/001-modernize-digitalocean-api/data-model.md` with godo struct definitions
- [ ] T016 Document error response structures and status codes in `data-model.md`
- [ ] T017 Create test fixtures for mock API responses in `digitalocean/testdata/` directory

**Deliverable**: Complete API contract documentation and data models

---

## Phase 3: Foundational (Blocking Prerequisites)

**Purpose**: Dependency updates that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T018 Update `go.mod` with latest godo SDK version (run `go get -u github.com/digitalocean/godo`)
- [x] T019 Update Cobra to v1.8+ in `go.mod` (run `go get -u github.com/spf13/cobra`)
- [x] T020 [P] Update Viper to v1.19+ in `go.mod` (run `go get -u github.com/spf13/viper`)
- [x] T021 [P] Update promptui to v0.9+ in `go.mod` (run `go get -u github.com/manifoldco/promptui`)
- [x] T022 [P] Update fatih/color to v1.18+ in `go.mod` (run `go get -u github.com/fatih/color`)
- [x] T023 Run `go mod tidy` and resolve any dependency conflicts
- [x] T024 Run `go mod download` to fetch all updated dependencies
- [x] T025 Update imports in `digitalocean/digitalocean.go` if package paths changed
- [x] T026 Run `go build` to verify project compiles with new dependencies
- [x] T027 Run existing tests (`make test`) to identify breaking changes from updates

**Checkpoint**: ✅ All dependencies updated, project compiles, existing test baseline established

---

## Phase 4: User Story 1 - Reliable Droplet Operations (Priority: P1) 🎯 MVP

**Goal**: Ensure core droplet management (create, list, destroy) works with current API

**Independent Test**: Run `cogo create`, `cogo list`, `cogo destroy` against live DigitalOcean test account

### Tests for User Story 1 (Write First - TDD)

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T028 [P] [US1] Create integration test for droplet creation in `tests/integration/droplet_create_test.go`
- [ ] T029 [P] [US1] Create integration test for droplet listing in `tests/integration/droplet_list_test.go`
- [ ] T030 [P] [US1] Create integration test for droplet deletion in `tests/integration/droplet_destroy_test.go`
- [ ] T031 [P] [US1] Create contract test for droplet API responses in `tests/contract/digitalocean_api_test.go`
- [ ] T032 [US1] Run tests to verify they FAIL (no implementation yet)

### Implementation for User Story 1

- [x] T033 [US1] Update `getToken()` function in `digitalocean/digitalocean.go` to verify authentication still works with latest SDK
- [x] T034 [US1] Update `CreateDroplet()` function in `digitalocean/digitalocean.go` - verify godo.DropletCreateRequest struct fields
- [x] T035 [US1] Update droplet creation API call to use current godo client methods
- [x] T036 [US1] Update `dropletList()` pagination logic in `digitalocean/digitalocean.go` for current API response structure
- [x] T037 [US1] Update `DisplayDropletList()` to handle new response fields (status, creation date, etc.)
- [x] T038 [US1] Update `DestroyDroplet()` deletion logic in `digitalocean/digitalocean.go`
- [x] T039 [US1] Verify multi-step confirmation flow still works correctly
- [x] T040 [US1] Test backward compatibility with existing `.cogo` configuration files
- [ ] T041 [US1] Run integration tests to verify they now PASS
- [ ] T042 [US1] Manual test: Create droplet via `cogo create` on test account
- [ ] T043 [US1] Manual test: List droplets via `cogo list` and verify all fields display
- [ ] T044 [US1] Manual test: Destroy droplet via `cogo destroy` with full confirmation flow

**Checkpoint**: User Story 1 complete - All core operations functional with current API

---

## Phase 5: User Story 2 - Access to New DigitalOcean Features (Priority: P2)

**Goal**: Provide access to all current regions, images, and droplet sizes

**Independent Test**: Compare available options in `cogo create` wizard against DigitalOcean dashboard

### Tests for User Story 2 (Write First - TDD)

- [ ] T045 [P] [US2] Create test to verify all current regions are listed in `tests/integration/regions_list_test.go`
- [ ] T046 [P] [US2] Create test to verify current OS images are listed in `tests/integration/images_list_test.go`
- [ ] T047 [P] [US2] Create test to verify current droplet sizes are listed in `tests/integration/sizes_list_test.go`
- [ ] T048 [US2] Run tests to verify they FAIL (no implementation yet)

### Implementation for User Story 2

- [ ] T049 [P] [US2] Update `regionList()` function in `digitalocean/digitalocean.go` for new regions and availability fields
- [ ] T050 [P] [US2] Update `imageDistributionList()` in `digitalocean/digitalocean.go` for new OS distributions
- [ ] T051 [P] [US2] Update `imageApplicationList()` in `digitalocean/digitalocean.go` for new application images
- [ ] T052 [P] [US2] Update `imageCustomList()` in `digitalocean/digitalocean.go` for custom image handling
- [ ] T053 [P] [US2] Update `sizeList()` in `digitalocean/digitalocean.go` for new droplet sizes and pricing
- [ ] T054 [P] [US2] Update `sshKeyList()` in `digitalocean/digitalocean.go` for SSH key fingerprint changes
- [ ] T055 [US2] Update pagination logic in all list functions to match current SDK
- [ ] T056 [US2] Update `ParseRegionListresults()` in `utils/utils.go` if struct fields changed
- [ ] T057 [US2] Update `ParseImageListResults()` in `utils/utils.go` if struct fields changed
- [ ] T058 [US2] Update `ParseSizeListResults()` in `utils/utils.go` if struct fields changed
- [ ] T059 [US2] Update `ParseSSHKeyListResults()` in `utils/utils.go` if struct fields changed
- [ ] T060 [US2] Run integration tests to verify they now PASS
- [ ] T061 [US2] Manual test: Verify new regions appear in region selection (compare with DO dashboard)
- [ ] T062 [US2] Manual test: Verify latest Ubuntu/Debian versions appear in image selection
- [ ] T063 [US2] Manual test: Verify new droplet sizes appear in size selection

**Checkpoint**: User Story 2 complete - All current DO offerings accessible

---

## Phase 6: User Story 3 - Improved Error Handling (Priority: P3)

**Goal**: Provide clear, actionable error messages for common API failures

**Independent Test**: Simulate API errors (invalid token, rate limit) and verify error messages

### Tests for User Story 3 (Write First - TDD)

- [ ] T064 [P] [US3] Create test for invalid token error handling in `tests/integration/error_handling_test.go`
- [ ] T065 [P] [US3] Create test for rate limit error handling in `tests/integration/error_handling_test.go`
- [ ] T066 [P] [US3] Create test for service unavailable error in `tests/integration/error_handling_test.go`
- [ ] T067 [US3] Run tests to verify they FAIL (no implementation yet)

### Implementation for User Story 3

- [ ] T068 [US3] Add error type checking in `getToken()` function in `digitalocean/digitalocean.go`
- [ ] T069 [US3] Add 401 Unauthorized handler with actionable message in `config/config.go`
- [ ] T070 [US3] Add 429 Rate Limit handler with retry guidance across all API calls
- [ ] T071 [US3] Add 500/503 Service Error handler with temporary failure message
- [ ] T072 [US3] Update error returns in `CreateDroplet()` with context
- [ ] T073 [US3] Update error returns in `DestroyDroplet()` with context
- [ ] T074 [US3] Update error returns in all list functions with context
- [ ] T075 [US3] Add error message helper function for consistent formatting
- [ ] T076 [US3] Run integration tests to verify they now PASS
- [ ] T077 [US3] Manual test: Try operation with invalid token, verify clear error message
- [ ] T078 [US3] Manual test: Trigger rate limit (many rapid API calls), verify helpful message

**Checkpoint**: User Story 3 complete - Enhanced error handling implemented

---

## Phase 7: Integration & Validation

**Purpose**: Verify all user stories work together and meet success criteria

- [ ] T079 Run full integration test suite (`go test ./tests/integration/...`)
- [ ] T080 Run contract tests (`go test ./tests/contract/...`)
- [ ] T081 Run existing unit tests (`make test`) to ensure no regressions
- [ ] T082 Verify SC-001: All commands complete successfully without errors
- [ ] T083 Verify SC-002: Compare available options with DO dashboard (100% parity)
- [ ] T084 Verify SC-003: Measure list operation time for 100 droplets (< 3 seconds)
- [ ] T085 Verify SC-004: Review all error messages for clarity and actionability
- [ ] T086 Verify SC-005: Test with old `.cogo` config file (backward compatibility)
- [ ] T087 Test full create/list/destroy workflow on test account
- [ ] T088 Test with no config file (interactive prompts)
- [ ] T089 Test with config file in each supported location
- [ ] T090 Run linter (`golangci-lint run`) and fix any issues
- [ ] T091 Check for any TODO or FIXME comments added during implementation

**Checkpoint**: All tests pass, all success criteria met

---

## Phase 8: Documentation & Polish

**Purpose**: Update documentation and prepare for release

- [ ] T092 [P] Update `CHANGELOG.md` with version bump and change summary
- [ ] T093 [P] Update `README.md` if any new requirements or dependencies
- [ ] T094 [P] Update `version/version.go` with new version number (follow semver)
- [ ] T095 [P] Review and update inline code comments for clarity
- [ ] T096 [P] Create release notes in `specs/001-modernize-digitalocean-api/release-notes.md`
- [ ] T097 Verify CI pipeline passes all checks (GitHub Actions)
- [ ] T098 Create git tag for new version
- [ ] T099 Build release binaries for all platforms (`make build`)
- [ ] T100 Final manual test of release binary on clean system

**Deliverable**: Ready for release

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Research)**: No dependencies - can start immediately
- **Phase 2 (Design)**: Depends on Phase 1 (Research) completion
- **Phase 3 (Foundational)**: Can start in parallel with Phase 2 - BLOCKS all user stories
- **Phase 4-6 (User Stories)**: All depend on Phase 3 (Foundational) completion
  - User stories can proceed sequentially (P1 → P2 → P3) or in parallel if team capacity
- **Phase 7 (Validation)**: Depends on all user stories being complete
- **Phase 8 (Documentation)**: Depends on Phase 7 validation passing

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 3) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 3) - Independent of US1
- **User Story 3 (P3)**: Can start after Foundational (Phase 3) - Independent of US1/US2

### Within Each User Story (TDD Flow)

1. Write tests FIRST (marked with test task numbers)
2. Run tests - verify they FAIL
3. Implement functionality
4. Run tests - verify they PASS
5. Manual validation
6. Story checkpoint

### Parallel Opportunities

- Phase 1: Tasks T001-T006 can all run in parallel (different research topics)
- Phase 2: Tasks T008-T014 can all run in parallel (different API endpoints)
- Phase 3: Tasks T020-T022 can run in parallel (different dependencies)
- Phase 4: Tasks T028-T031 (tests) can run in parallel
- Phase 5: Tasks T045-T047 (tests) can run in parallel
- Phase 5: Tasks T049-T054 (implementations) can run in parallel (different functions)
- Phase 6: Tasks T064-T066 (tests) can run in parallel
- Phase 8: Tasks T092-T096 (documentation) can run in parallel

---

## Implementation Strategy

### MVP First (Recommended)

1. Complete Phase 1: Research (T001-T007)
2. Complete Phase 2: Design & Contracts (T008-T017)
3. Complete Phase 3: Foundational (T018-T027) - CRITICAL
4. Complete Phase 4: User Story 1 Only (T028-T044)
5. **STOP and VALIDATE**: Test core operations thoroughly
6. If US1 works, consider deploying/releasing before adding US2/US3

### Incremental Delivery

1. Foundation (Phases 1-3) → Dependencies updated, project compiles
2. Add User Story 1 (Phase 4) → Test → Tag as v2.0.0-beta1 (MVP!)
3. Add User Story 2 (Phase 5) → Test → Tag as v2.0.0-beta2
4. Add User Story 3 (Phase 6) → Test → Tag as v2.0.0-rc1
5. Complete validation & docs (Phases 7-8) → Release v2.0.0

### Parallel Team Strategy

With 2-3 developers:

1. Team completes Research + Design + Foundational together (Phases 1-3)
2. Once Foundational is done:
   - Developer A: User Story 1 (Phase 4)
   - Developer B: User Story 2 (Phase 5)
   - Developer C: User Story 3 (Phase 6)
3. Merge and validate together (Phase 7)
4. Polish and release (Phase 8)

---

## Time Estimates

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1: Research | T001-T007 | 2-4 hours |
| Phase 2: Design | T008-T017 | 2-3 hours |
| Phase 3: Foundational | T018-T027 | 2-3 hours |
| Phase 4: User Story 1 | T028-T044 | 4-5 hours |
| Phase 5: User Story 2 | T045-T063 | 3-4 hours |
| Phase 6: User Story 3 | T064-T078 | 2-3 hours |
| Phase 7: Validation | T079-T091 | 2-3 hours |
| Phase 8: Documentation | T092-T100 | 1-2 hours |
| **Total** | **100 tasks** | **18-27 hours** |

---

## Notes

- [P] tasks = different files, no dependencies - can run in parallel
- [Story] label maps task to specific user story for traceability (US1, US2, US3)
- Each user story is independently completable and testable
- TDD flow: Write tests → Verify FAIL → Implement → Verify PASS
- Commit after each logical group of tasks (end of each phase checkpoint)
- Stop at any checkpoint to validate independently
- Test account required: Create DigitalOcean test account with API token before Phase 4

