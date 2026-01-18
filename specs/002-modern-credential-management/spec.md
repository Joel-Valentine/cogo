# Feature Specification: Modern Credential Management

**Feature ID:** 002  
**Status:** Draft  
**Created:** 2026-01-18  
**Author:** System

## Overview

Modernize cogo's credential management to follow 2026 best practices used by popular CLI tools (gh, aws-vault, doctl). Replace plain-text configuration file storage with secure OS keychain integration while maintaining backward compatibility.

## Problem Statement

### Current Issues
1. **Security Risk**: API tokens stored in plain text in `~/.cogo` JSON file
2. **Limited Flexibility**: No environment variable support despite Viper configuration
3. **Poor UX**: Interactive prompts required on every run if config missing
4. **Outdated Pattern**: Implementation from 2020 doesn't match modern CLI standards
5. **No Automation Support**: Difficult to use in CI/CD or scripts

### User Impact
- Users must choose between convenience (saving token) and security (entering each time)
- No way to use different tokens for different environments
- Token exposure risk if `.cogo` file is accidentally committed or backed up

## User Stories

### Story 1: Secure Default Storage
**As a** developer using cogo on my local machine  
**I want** my DigitalOcean token stored securely in my OS keychain  
**So that** I don't have to worry about accidental token exposure

**Acceptance Criteria:**
- First-time users are prompted to enter token
- Token is automatically stored in OS keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Token is retrieved transparently on subsequent runs
- No plain-text files created by default

### Story 2: Environment Variable Support
**As a** DevOps engineer running cogo in CI/CD  
**I want** to provide credentials via environment variables  
**So that** I can automate deployments without interactive prompts

**Acceptance Criteria:**
- `DIGITALOCEAN_TOKEN` environment variable is checked first
- `COGO_DIGITALOCEAN_TOKEN` is also supported
- Env vars take precedence over keychain/config file
- No prompts when env var is set

### Story 3: Credential Management Commands
**As a** user managing multiple DigitalOcean accounts  
**I want** explicit commands to set/get/delete credentials  
**So that** I can easily switch between accounts

**Acceptance Criteria:**
- `cogo config set-token` stores token securely
- `cogo config get-token` displays masked token (first/last 4 chars)
- `cogo config delete-token` removes token from storage
- `cogo config status` shows where token is coming from

### Story 4: Backward Compatibility
**As an** existing cogo user with `~/.cogo` file  
**I want** my existing setup to continue working  
**So that** I'm not forced to reconfigure immediately

**Acceptance Criteria:**
- Existing `~/.cogo` files are still read
- Warning displayed about insecure storage on first use
- Migration prompt offers to move token to keychain
- Option to disable warnings for CI/CD environments

## Technical Requirements

### Functional Requirements
1. **Priority Order** (highest to lowest):
   - CLI flag: `--token` or `--digitalocean-token`
   - Environment variable: `DIGITALOCEAN_TOKEN`
   - Environment variable: `COGO_DIGITALOCEAN_TOKEN`
   - OS Keychain: `cogo/digitalocean-token`
   - Config file: `~/.cogo` (legacy, with warning)
   - Interactive prompt: (last resort)

2. **Storage Mechanisms**:
   - OS Keychain via `github.com/zalando/go-keyring`
   - Viper for env vars and config file
   - In-memory for single-session tokens

3. **Commands**:
   ```bash
   cogo config set-token <token>      # Store in keychain
   cogo config get-token              # Display (masked)
   cogo config delete-token           # Remove from keychain
   cogo config status                 # Show source and validity
   cogo config migrate                # Move from file to keychain
   ```

4. **Error Handling**:
   - Clear error messages for each source failure
   - Graceful fallback through priority chain
   - Helpful hints when no token found

### Non-Functional Requirements
1. **Security**: No plain-text token storage by default
2. **Performance**: Token retrieval < 50ms
3. **Compatibility**: Works on macOS, Linux, Windows
4. **Backward Compat**: Existing users not disrupted
5. **Testing**: Unit tests for all credential sources

## Key Entities

### `credentials.Provider`
```go
type Provider interface {
    GetToken(ctx context.Context) (string, error)
    SetToken(ctx context.Context, token string) error
    DeleteToken(ctx context.Context) error
    GetSource() string  // Returns where token came from
}
```

### `credentials.Manager`
```go
type Manager struct {
    providers []Provider  // Ordered by priority
}
```

### Provider Implementations
- `FlagProvider` - from CLI flag
- `EnvProvider` - from environment variables
- `KeychainProvider` - from OS keychain
- `FileProvider` - from legacy config file (deprecated)
- `PromptProvider` - interactive prompt

## Edge Cases

1. **Keychain unavailable** (headless Linux): Fall back to env var or prompt
2. **Multiple tokens set**: Follow priority order strictly
3. **Token migration**: Handle partial migration states
4. **Concurrent access**: Handle multiple cogo instances
5. **Invalid tokens**: Validate format before storage
6. **Token rotation**: Easy way to update stored token

## Success Criteria

### Must Have
- ✅ OS keychain integration working on macOS/Linux/Windows
- ✅ Environment variable support with correct priority
- ✅ `cogo config` commands functional
- ✅ Backward compatibility with existing `.cogo` files
- ✅ Security warnings for plain-text storage
- ✅ Migration helper for existing users

### Nice to Have
- 🎯 Token validation before storage
- 🎯 Multiple named credential profiles
- 🎯 Token expiry tracking and warnings
- 🎯 Integration with DigitalOcean OAuth flow

### Out of Scope (Future)
- Multiple cloud provider support (AWS, GCP, Azure)
- Team credential sharing
- Credential encryption keys
- Audit logging of credential access

## Assumptions

1. Users have OS keychain available (or will use env vars)
2. DigitalOcean API tokens don't expire (unless revoked)
3. Token format is `dop_v1_[a-f0-9]{64}`
4. Most users will migrate willingly with proper guidance

## Dependencies

- `github.com/zalando/go-keyring` v1.2.2+ for OS keychain access
- Existing `github.com/spf13/viper` for env vars and config
- Existing `github.com/spf13/cobra` for new commands
- Existing `github.com/manifoldco/promptui` for interactive prompts

## Migration Guide

For existing users:

```bash
# Option 1: Automatic migration
cogo config migrate

# Option 2: Manual
# 1. Note your current token
cat ~/.cogo
# 2. Set it securely
cogo config set-token dop_v1_xxx
# 3. Remove old file
rm ~/.cogo
```

## Security Considerations

1. **Token Masking**: Display only first/last 4 characters
2. **No Logging**: Never log actual token values
3. **Secure Deletion**: Overwrite memory when deleting tokens
4. **File Permissions**: If using file storage, enforce 0600
5. **Warning Display**: Clear warnings about insecure storage methods

## Testing Strategy

1. **Unit Tests**:
   - Each provider implementation
   - Priority ordering logic
   - Token validation

2. **Integration Tests**:
   - End-to-end token retrieval
   - Migration scenarios
   - Cross-platform keychain access

3. **Manual Testing**:
   - Fresh install on clean system
   - Migration from existing setup
   - CI/CD environment usage

