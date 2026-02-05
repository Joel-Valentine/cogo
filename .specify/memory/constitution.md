# Cogo Constitution

## Core Principles

### I. CLI-First Design
Cogo is a command-line interface tool built with Cobra. All functionality must be accessible via CLI commands. Commands should be intuitive, follow the pattern `cogo <action>`, and provide clear help text. Interactive wizards enhance the user experience but must not be the only way to access functionality. Commands should support both interactive and non-interactive modes where possible.

### II. Provider Abstraction
Cogo supports multiple cloud providers through a consistent interface. Each provider implementation lives in its own package (e.g., `digitalocean/`). Provider-specific code must be isolated from core CLI logic. New providers should follow the existing pattern: implement provider-specific operations while maintaining the same command structure. The `utils` package provides shared functionality that works across providers.

### III. Safety First
Destructive operations (like `destroy`) require multiple confirmation steps to prevent accidental data loss. Error messages must be clear and actionable. When operations fail, provide helpful guidance on what went wrong and how to recover. Never silently fail or proceed with destructive operations without explicit user confirmation.

### IV. Test-Driven Development
All new features and bug fixes must include tests. Unit tests for utilities and helper functions are mandatory. Integration tests should cover provider interactions where feasible. Tests must be written before or alongside implementation. Use Go's standard testing package and follow Go testing conventions.

### V. Configuration Management
Cogo supports configuration files in multiple locations (`$HOME/.cogo`, `$HOME/.config/.cogo`, `./.cogo`). Configuration is optional - the tool should work without it by prompting for required values. Configuration files use JSON format. Never commit user-specific configuration files or API tokens to version control.

### VI. Simplicity and Maintainability
Keep the codebase simple and easy to understand. Follow Go idioms and conventions. Avoid unnecessary abstractions - prefer straightforward implementations. Code should be readable by developers learning Go. When adding complexity, document why it's necessary. YAGNI (You Aren't Gonna Need It) principles apply.

### VII. Consistent Navigation and User Experience
All CLI interactions must follow the navigation framework defined in `specs/003-consistent-cli-navigation/`. Multi-step operations must support back navigation, cancellation (Ctrl+C, Esc), and graceful empty state handling. Error messages must be clear, colored (✗ for errors, ✓ for success, ⚠️ for warnings), and actionable. Input validation happens after Enter (not per-keystroke) to prevent spam. All commands must feel identical across cloud providers - users should not be able to tell which provider they're using based on navigation patterns. Empty resource states exit with code 0 (not treated as errors). See `specs/003-consistent-cli-navigation/navigation-patterns.md` for detailed standards.

## Technology Standards

**Language**: Go 1.24+  
**CLI Framework**: Cobra (spf13/cobra)  
**Navigation Framework**: Custom framework (`navigation/` package) with promptui integration  
**Configuration**: JSON files with Viper for parsing  
**Credentials**: Modern multi-source system (keychain → env vars → config file → prompt)  
**Testing**: Standard Go testing package with testify assertions  
**Dependencies**: Minimize external dependencies; prefer standard library when possible

## Development Workflow

1. **Feature Development**: Create feature branches from main. Follow the existing command structure pattern.
2. **Testing**: Run `make test` before committing. Ensure all tests pass.
3. **Code Review**: All changes require review. Code should be clear, tested, and follow Go conventions.
4. **Documentation**: Update README.md for user-facing changes. Document complex logic inline.
5. **Provider Addition**: When adding a new provider, follow the DigitalOcean implementation pattern. Ensure all core commands (create, list, destroy) are supported.

## Governance

This constitution supersedes all other development practices. All pull requests and code reviews must verify compliance with these principles. When violating a principle is necessary, document the justification in the PR description and consider it a technical debt item.

Amendments to this constitution require:
- Clear justification for the change
- Discussion and consensus
- Update to this document with version tracking

**Version**: 2.0.0 | **Ratified**: 2025-01-27 | **Last Amended**: 2026-01-18 | **Major Changes**: Added navigation framework standards (Principle VII), updated Go version to 1.24, modernized credential management
