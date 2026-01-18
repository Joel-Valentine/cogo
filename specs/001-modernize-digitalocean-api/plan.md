# Implementation Plan: Modernize DigitalOcean API Integration

**Branch**: `001-modernize-digitalocean-api` | **Date**: 2026-01-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/001-modernize-digitalocean-api/spec.md`

## Summary

Update Cogo's DigitalOcean integration to work with the current DigitalOcean API v2 by upgrading the godo SDK from v1.34.0 (2020) to the latest version, updating API call patterns, and ensuring all droplet management operations (create, list, destroy) function correctly with current API endpoints and response structures. The technical approach involves dependency updates, testing against current API documentation, and maintaining backward compatibility with existing CLI workflows.

## Technical Context

**Language/Version**: Go 1.24  
**Primary Dependencies**: 
- godo SDK (currently v1.34.0, target: latest stable ~v1.130+)
- Cobra v1.8+ (CLI framework - upgrade from v0.0.7)
- Viper v1.19+ (configuration - upgrade from v1.6.3)
- promptui v0.9+ (interactive prompts - upgrade from v0.7.0)
- fatih/color v1.18+ (terminal colors - upgrade from v1.9.0)

**Storage**: JSON configuration files (`.cogo`) in user's home directory or current directory  
**Testing**: Go standard testing package (`go test`)  
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)  
**Project Type**: Single project (CLI tool)  
**Performance Goals**: API list operations complete under 3 seconds for accounts with up to 100 droplets  
**Constraints**: 
- Must maintain backward compatibility with existing CLI commands
- Must preserve current configuration file format
- Must respect DigitalOcean API rate limits (5000 requests/hour per token)
- Must work without configuration file (interactive prompts as fallback)

**Scale/Scope**: 
- Single CLI binary
- ~700 lines of Go code in digitalocean/ package to update
- Primary file: `digitalocean/digitalocean.go`
- Support for DigitalOcean API v2 only (v1 deprecated)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Compliance Status: PASSED

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. CLI-First Design** | ✅ Pass | Maintaining existing `cogo create/list/destroy` command structure. No changes to CLI interface. |
| **II. Provider Abstraction** | ✅ Pass | All changes contained within `digitalocean/` package. Provider isolation maintained. |
| **III. Safety First** | ✅ Pass | Preserving multi-step confirmation for destroy operations. Enhanced error messages improve safety. |
| **IV. Test-Driven Development** | ✅ Pass | Will add integration tests for API operations. Existing tests must continue passing. |
| **V. Configuration Management** | ✅ Pass | Preserving JSON configuration format and multiple location support. No breaking changes. |
| **VI. Simplicity and Maintainability** | ✅ Pass | Straightforward dependency updates. Minimal code changes focused on API compatibility. |

**No constitutional violations detected.** This is a maintenance update that aligns with all core principles.

## Project Structure

### Documentation (this feature)

```text
specs/001-modernize-digitalocean-api/
├── plan.md              # This file
├── spec.md              # Feature specification (completed)
├── research.md          # Phase 0: API changelog research
├── contracts/           # Phase 1: API contract documentation
│   ├── droplets-api.md  # Droplet endpoints and responses
│   ├── images-api.md    # Image listing endpoints
│   └── regions-api.md   # Region and size listings
├── data-model.md        # Phase 1: API response structures
└── checklists/
    └── requirements.md  # Specification quality checklist (completed)
```

### Source Code (repository root)

```text
digitalocean/
└── digitalocean.go     # Primary file requiring updates (API calls, struct handling)

config/
└── config.go           # May need minor updates for new error types

utils/
├── utils.go            # Shared utilities (minimal changes expected)
└── utils_test.go       # Unit tests to verify

cmd/
├── root.go             # CLI commands (no changes expected)
└── version.go          # Version command (no changes)

go.mod                  # Dependency version updates
go.sum                  # Generated checksums

tests/ (to be created if needed)
├── integration/
│   └── digitalocean_integration_test.go
└── contract/
    └── api_contract_test.go
```

**Structure Decision**: Single project structure maintained. All updates focus on the `digitalocean/` package with dependency version bumps in `go.mod`. No architectural changes needed—this is a maintenance update to align with current API specifications.

## Implementation Phases

### Phase 0: Research & Discovery

**Objective**: Understand API changes between godo v1.34.0 (2020) and latest version

**Tasks**:
1. Review DigitalOcean API v2 changelog from 2020 to 2026
2. Identify godo SDK latest stable version and breaking changes
3. Document deprecated methods and their replacements
4. List new required parameters or response fields
5. Review pagination changes and rate limit updates
6. Document authentication flow changes (if any)

**Deliverables**:
- `research.md` - API changelog summary with migration notes
- List of affected functions in `digitalocean/digitalocean.go`
- Dependency version targets for `go.mod`

### Phase 1: Design & Contracts

**Objective**: Define updated API contracts and data structures

**Tasks**:
1. Document current API endpoints Cogo uses:
   - POST `/v2/droplets` (create)
   - GET `/v2/droplets` (list)
   - DELETE `/v2/droplets/:id` (destroy)
   - GET `/v2/regions` (list regions)
   - GET `/v2/images` (list images)
   - GET `/v2/sizes` (list sizes)
   - GET `/v2/account/keys` (list SSH keys)

2. Map each endpoint to current API specification
3. Identify struct field changes in godo SDK
4. Document pagination patterns (Links header handling)
5. Define error response handling patterns
6. Create test fixtures for API responses

**Deliverables**:
- `contracts/` directory with API endpoint documentation
- `data-model.md` with struct definitions and field mappings
- Test fixtures for mock API responses

### Phase 2: Implementation

**Objective**: Update code to work with current DigitalOcean API

**Priority Order** (based on user stories):

#### P1: Core Operations (User Story 1)
1. Update `go.mod` dependencies
   - Upgrade godo SDK to latest
   - Upgrade Cobra, Viper, promptui, color packages
   - Run `go mod tidy` and `go mod download`

2. Update API authentication in `getToken()`
   - Verify bearer token auth pattern still valid
   - Update client initialization if needed

3. Update droplet operations:
   - `CreateDroplet()` - verify request struct, parameters
   - `dropletList()` - update pagination handling
   - `DestroyDroplet()` - verify delete endpoint

4. Run integration tests against live API (test account)

#### P2: Resource Listings (User Story 2)
1. Update resource list functions:
   - `regionList()` - new regions, availability fields
   - `imageDistributionList()` - new OS versions
   - `imageApplicationList()` - new application images
   - `imageCustomList()` - custom image handling
   - `sizeList()` - new droplet sizes, pricing
   - `sshKeyList()` - SSH key fingerprint handling

2. Update selection prompts with new options
3. Test each list operation for completeness

#### P3: Error Handling (User Story 3)
1. Update error handling in all API calls
2. Add specific handling for:
   - 401 Unauthorized (token issues)
   - 429 Too Many Requests (rate limits)
   - 500/503 Service Errors
3. Provide actionable error messages

**Code Changes Required**:
- `digitalocean/digitalocean.go`: Update all API interaction code
- `go.mod`: Dependency version bumps
- `config/config.go`: Minor updates for error handling
- Create integration tests

### Phase 3: Testing & Validation

**Test Strategy**:

1. **Unit Tests** (`utils/utils_test.go`):
   - Verify utility functions still work
   - Test input validation
   - Run existing tests to ensure no regressions

2. **Integration Tests** (new):
   - Create test account on DigitalOcean
   - Test full create/list/destroy workflow
   - Verify all resource list operations
   - Test error scenarios (invalid token, rate limits)

3. **Manual Testing**:
   - Test against real DigitalOcean account
   - Verify all regions, images, sizes appear
   - Confirm backward compatibility with existing config files
   - Test interactive prompts work correctly

**Success Criteria Verification**:
- ✅ SC-001: All commands work without errors
- ✅ SC-002: All current DO offerings visible
- ✅ SC-003: List ops under 3 seconds
- ✅ SC-004: Clear error messages
- ✅ SC-005: Backward compatibility maintained

### Phase 4: Documentation & Deployment

**Tasks**:
1. Update CHANGELOG.md with version bump and changes
2. Update README.md if any new requirements
3. Verify CI pipeline passes (GitHub Actions)
4. Create release notes documenting updates
5. Tag new version following semantic versioning

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking changes in godo SDK | High | Review migration guide, use feature flags if needed, extensive testing |
| API rate limit violations during testing | Medium | Use test account, implement exponential backoff, respect rate limits |
| Backward compatibility breaks | High | Thorough testing with existing configs, maintain config file format |
| New required API parameters | Medium | Review API docs carefully, provide sensible defaults, prompt users when needed |
| Dependency conflicts | Low | Update all dependencies together, test thoroughly |

## Timeline Estimate

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| Phase 0: Research | 2-4 hours | DigitalOcean API docs, godo changelog |
| Phase 1: Design | 2-3 hours | Phase 0 complete |
| Phase 2: Implementation | 6-8 hours | Phase 1 complete, test account ready |
| Phase 3: Testing | 3-4 hours | Phase 2 complete |
| Phase 4: Documentation | 1-2 hours | All phases complete |
| **Total** | **14-21 hours** | Assumes no major blockers |

## Next Steps

1. ✅ Specification created and validated
2. ✅ Implementation plan created
3. ⏭️ Run `/speckit.tasks` to generate actionable task breakdown
4. 📋 Begin Phase 0 research on DigitalOcean API changes
5. 🔨 Proceed with phased implementation

---

**Plan Status**: Ready for task generation  
**Blockers**: None  
**Prerequisites**: Test DigitalOcean account with API token

