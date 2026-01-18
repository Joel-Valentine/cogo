# Research: DigitalOcean API Modernization

**Feature**: 001-modernize-digitalocean-api  
**Created**: 2026-01-18  
**Purpose**: Document API changes between 2020 (v1.34.0) and current (2026) to guide migration

## Current State Analysis

### Existing Implementation (2020)

**Dependencies** (from `go.mod`):
- `github.com/digitalocean/godo v1.34.0` (Released: ~March 2020)
- `github.com/spf13/cobra v0.0.7`
- `github.com/spf13/viper v1.6.3`
- `github.com/manifoldco/promptui v0.7.0`
- `github.com/fatih/color v1.9.0`

**Primary File**: `digitalocean/digitalocean.go` (~700 lines)

**Functions Requiring Update**:
1. `CreateDroplet()` - Line ~30
2. `DestroyDroplet()` - Line ~149
3. `getToken()` - Line ~257
4. `DisplayDropletList()` - Line ~289
5. `dropletList()` - Line ~319
6. `regionList()` - Line ~352
7. `imageDistributionList()` - Line ~422
8. `imageApplicationList()` - Line ~457
9. `imageCustomList()` - Line ~492
10. `sizeList()` - Line ~387
11. `sshKeyList()` - Line ~527

## Latest SDK Information (2026)

### Godo SDK Updates

**Target Version**: Latest stable (need to verify specific version)

**Known Changes from v1.34.0 to Latest**:

1. **Authentication**:
   - Still uses OAuth2 bearer token (unchanged)
   - `godo.NewFromToken()` method still valid
   - No breaking changes expected

2. **Pagination**:
   - Links header handling remains similar
   - `resp.Links.IsLastPage()` and `resp.Links.CurrentPage()` still valid
   - May have additional pagination helpers

3. **Struct Changes**:
   - `godo.DropletCreateRequest` - verify all fields
   - `godo.Droplet` - may have new fields (tags, monitoring, backups status)
   - `godo.Region` - new regions added, availability fields
   - `godo.Image` - new distributions, status fields
   - `godo.Size` - new size options, pricing updates
   - `godo.Key` - fingerprint handling

4. **Error Handling**:
   - Improved error types and messages
   - Rate limiting information in responses
   - Better 4xx/5xx error differentiation

5. **API Endpoints** (all v2, no changes to paths):
   - POST `/v2/droplets` ✅ Stable
   - GET `/v2/droplets` ✅ Stable
   - DELETE `/v2/droplets/:id` ✅ Stable
   - GET `/v2/regions` ✅ Stable
   - GET `/v2/images` ✅ Stable
   - GET `/v2/sizes` ✅ Stable
   - GET `/v2/account/keys` ✅ Stable

## Compatibility Assessment

### Low Risk Updates

✅ **No Breaking Changes Expected**:
- Core API endpoints unchanged (all v2)
- Authentication method unchanged (bearer token)
- Pagination pattern unchanged (Links header)
- Basic struct fields maintained for backward compatibility

### Medium Risk Updates

⚠️ **Potential Adjustments Needed**:
- New optional fields in `DropletCreateRequest` (tags, monitoring, vpc_uuid)
- Additional response fields in list operations
- New droplet sizes, regions, images to test
- Error response structure enhancements

### Code Changes Required

**Priority 1** (Core Functionality):
1. Update `go.mod` dependencies
2. Test existing API calls with new SDK
3. Verify struct field compatibility
4. Test pagination with large result sets

**Priority 2** (Enhanced Features):
1. Expose new regions in selection
2. Add new OS distributions to image lists
3. Include new droplet sizes
4. Update SSH key handling if needed

**Priority 3** (Error Handling):
1. Enhanced error messages for auth failures
2. Rate limit detection and messaging
3. Service unavailability handling

## Migration Strategy

### Phase 1: Safe Updates
1. Update dependencies in isolated branch
2. Run existing tests to identify breaks
3. Fix compilation errors
4. Verify backward compatibility

### Phase 2: Feature Parity
1. Test against live API
2. Compare available options with DO dashboard
3. Verify all current offerings accessible
4. Test error scenarios

### Phase 3: Enhanced Experience
1. Improve error messages
2. Add rate limit handling
3. Optimize list operations
4. Add any missing validations

## Testing Requirements

### Integration Tests Needed

1. **Authentication**:
   - Valid token → Success
   - Invalid token → Clear error
   - Expired token → Refresh guidance

2. **Droplet Operations**:
   - Create with various configs
   - List with pagination
   - Destroy with confirmations

3. **Resource Listings**:
   - All regions returned
   - All images returned (by type)
   - All sizes returned
   - All SSH keys returned

4. **Error Scenarios**:
   - 401 Unauthorized
   - 429 Rate Limit
   - 500/503 Service Errors
   - Network timeouts

## Dependencies for Implementation

**Required Before Starting**:
- [ ] DigitalOcean test account with API token
- [ ] Access to DigitalOcean dashboard for comparison
- [ ] Go 1.24 development environment
- [ ] Network access to api.digitalocean.com

**Documentation References**:
- DigitalOcean API v2 Documentation: https://docs.digitalocean.com/reference/api/api-reference/
- Godo SDK Repository: https://github.com/digitalocean/godo
- Godo SDK Documentation: https://pkg.go.dev/github.com/digitalocean/godo

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Breaking SDK changes | Review migration guide, test thoroughly, use feature flags if needed |
| API rate limits during testing | Use test account, implement backoff, respect limits |
| Backward compatibility issues | Preserve config format, test with existing configs |
| New required parameters | Review docs carefully, provide defaults, prompt when needed |

## Next Steps

1. ✅ Research document created
2. ⏭️ Verify latest godo SDK version
3. ⏭️ Create API contract documentation
4. ⏭️ Document data model changes
5. ⏭️ Begin Phase 3: Foundational (dependency updates)

---

**Status**: Research phase in progress  
**Blockers**: Need to verify specific godo SDK version number  
**Risk Level**: Low - No major breaking changes expected

