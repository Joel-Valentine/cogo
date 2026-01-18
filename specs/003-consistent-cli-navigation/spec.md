# Feature Specification: Consistent CLI Navigation

**Feature Branch**: `003-consistent-cli-navigation`  
**Created**: 2026-01-18  
**Status**: Draft  
**Input**: User description: "a lot of the patterns defined within interactions with the cli have no way of going back for example if I go to destroy and theres no droplets to destroy theres no way to go back. Even worse if I click around on my keyboard it actually errors. can we make sure that all interactions throughout the cli tool have some consistent navigation, perhaps that can apply to all clouds if possible? Lets make a nice interaction abstract interface that will apply to all experiences thats clear and obvious"

**Clarifications Applied**:
- Scope expanded to cover ALL commands (create, list, destroy, config, etc.) not just destroy
- Emphasis on "back" navigation through multi-step flows
- Research component added to understand industry-standard CLI UX patterns before implementation

## User Scenarios & Testing *(mandatory)*

### User Story 0 - CLI UX Pattern Research (Priority: P0)

Before implementing navigation patterns, research and document industry-standard CLI user experience patterns used by popular developer tools to ensure cogo follows familiar conventions.

**Why this priority**: Foundation for all other stories. Understanding what developers expect from CLI tools ensures we build intuitive, familiar patterns rather than inventing potentially confusing ones. This is P0 (prerequisite) because it informs implementation of all other priorities.

**Independent Test**: Can be completed and reviewed independently by examining 10-15 popular CLI tools (gh, kubectl, terraform, aws-cli, gcloud, docker, npm, cargo, git, etc.) and documenting their navigation patterns. Deliverable is a research document with recommendations.

**Acceptance Scenarios**:

1. **Given** need to understand CLI best practices, **When** researching popular developer tools, **Then** document how each handles: empty states, cancellation, back navigation, error recovery, and multi-step flows
2. **Given** research findings, **When** analyzing patterns, **Then** identify common conventions used by 80%+ of tools (e.g., Ctrl+C for cancel, clear error messages)
3. **Given** pattern analysis complete, **When** creating recommendations, **Then** prioritize patterns that are: widely used, intuitive, and compatible with existing cogo design
4. **Given** recommendations documented, **When** proceeding to implementation, **Then** use research findings to inform all navigation design decisions

---

### User Story 1 - Graceful Empty State Handling Across All Commands (Priority: P1)

When a user initiates ANY operation that has no available resources (destroy with no droplets, list with no resources, config operations with no tokens, create with no SSH keys, etc.), they should see a clear message and be able to exit gracefully rather than being trapped or seeing crashes.

**Why this priority**: This is the most critical issue as it currently causes crashes and leaves users with no way to recover. It directly impacts user trust and basic functionality across ALL commands.

**Independent Test**: Can be fully tested by running any cogo command with empty prerequisites (e.g., `cogo destroy` with no droplets, `cogo list` with no resources, `cogo create` with no SSH keys). System should display appropriate message and return to command prompt without errors.

**Acceptance Scenarios**:

1. **Given** no droplets exist, **When** user runs `cogo destroy`, **Then** system displays "No droplets available to destroy" and exits cleanly
2. **Given** no SSH keys configured, **When** user runs `cogo create`, **Then** system displays "No SSH keys found. Please add SSH keys first." and provides clear next steps or offers to continue without SSH
3. **Given** no resources exist, **When** user runs `cogo list`, **Then** system displays "No resources found. Create one with 'cogo create'" and exits cleanly
4. **Given** empty state in ANY interactive menu across ALL commands, **When** user presses any keyboard keys, **Then** system handles input gracefully without crashing

---

### User Story 2 - Universal Back/Cancel Navigation with Step-by-Step Return (Priority: P1)

Users should be able to cancel OR go back step-by-step at any point during multi-step operations using consistent keyboard shortcuts that work across ALL commands (create, destroy, list, config, etc.) and all cloud providers.

**Why this priority**: Provides essential escape mechanisms for users who change their mind or select wrong options. This is elevated to P1 because inability to back out of operations is equally critical to empty state handling - both trap users in the CLI. The "back" functionality is particularly important for complex flows like `create` where users may want to change previous selections.

**Independent Test**: Can be tested by starting ANY multi-step operation (create, destroy, list, config set-token) and pressing `Esc` at each step. System should go back one step. Pressing `Ctrl+C` should exit immediately. Both should work identically across all commands.

**Acceptance Scenarios**:

1. **Given** user is at step 3 of 5 in `cogo create`, **When** user presses `Esc` key, **Then** system returns to step 2 with previous selection retained
2. **Given** user is at step 2 of `cogo create` and presses `Esc` again, **When** at step 1, **Then** pressing `Esc` exits to command prompt
3. **Given** user is in any interactive prompt across ANY command, **When** user presses `Ctrl+C`, **Then** operation cancels immediately and returns to command prompt
4. **Given** user is mid-way through `cogo create`, **When** user goes back and changes region selection, **Then** subsequent steps (like size selection) update to reflect new region's available options
5. **Given** user goes back multiple steps in `cogo destroy`, **When** they change droplet selection, **Then** confirmation prompt updates with new droplet details
6. **Given** user is in `cogo config set-token`, **When** user cancels mid-entry, **Then** no partial configuration is saved
7. **Given** user is in ANY confirmation prompt across ALL commands, **When** user types 'n' or 'no', **Then** system explicitly asks "Go back to previous step? (y/n)" or exits
8. **Given** user completes operation after going back, **When** reviewing actions, **Then** only final choices are executed (no double-execution of backed-out steps)

---

### User Story 3 - Input Validation and Crash Prevention (Priority: P2)

The system should validate all keyboard input and handle unexpected input gracefully without panicking or crashing, especially during interactive selection menus.

**Why this priority**: Prevents crashes that erode user confidence. Critical for production stability but slightly lower priority than providing basic navigation.

**Independent Test**: Can be tested by pressing random keys, special characters, and rapid input during any interactive menu. System should never crash and should handle or ignore invalid input appropriately.

**Acceptance Scenarios**:

1. **Given** user is in an interactive selection menu with no items, **When** user presses arrow keys or any character keys, **Then** system displays "No items available" and does not panic
2. **Given** user is typing in a text prompt, **When** user enters special characters or very long input, **Then** system validates and provides clear error message if invalid
3. **Given** user is in any menu, **When** user presses invalid keys rapidly, **Then** system ignores or queues input appropriately without crashing
4. **Given** an index out of range error could occur, **When** user interacts with empty list, **Then** system checks bounds before access and displays appropriate message

---

### User Story 4 - Consistent Navigation Across ALL Commands and Cloud Providers (Priority: P3)

All commands (create, destroy, list, config, version, etc.) and all cloud provider integrations (DigitalOcean, future AWS, GCP, etc.) should use identical navigation patterns, keyboard shortcuts, and menu behaviors so users have a consistent experience everywhere in cogo.

**Why this priority**: Important for long-term scalability and user experience, but less urgent than fixing current crashes and navigation gaps. However, implementing base navigation infrastructure early prevents technical debt.

**Independent Test**: Can be tested by creating a navigation flow diagram for all commands and verifying consistent behavior. Test matrix should cover: create, destroy, list, config (all subcommands), and any future commands across all providers.

**Acceptance Scenarios**:

1. **Given** user learns navigation in `cogo create`, **When** they use `cogo destroy` or `cogo config`, **Then** navigation keys (`Esc`, arrow keys, `Ctrl+C`) work identically
2. **Given** user is familiar with DigitalOcean operations, **When** they use a different cloud provider, **Then** all navigation patterns work exactly the same
3. **Given** user encounters an empty state in ANY command or provider, **When** they check the message format, **Then** all use consistent formatting and helpful next steps
4. **Given** user is in multi-step operation for ANY command or provider, **When** they cancel or go back, **Then** all handle navigation identically
5. **Given** navigation patterns are documented in an abstract interface, **When** new command or provider is added, **Then** implementation must satisfy the interface contract

---

### Edge Cases

- What happens when user presses `Ctrl+C` during network operation (API call in progress across ANY command)?
- How does system handle rapid key presses that could queue up multiple commands in different operations?
- What if user's terminal doesn't support certain key bindings (e.g., `Esc` key)?
- How to handle graceful degradation on terminals with limited capabilities?
- What if there's a race condition between user input and list loading in ANY command?
- How to prevent panic when `promptui` library encounters unexpected state in ANY interactive prompt?
- What happens when user goes back to step 1 after step 5, changes selection, and previous steps become invalid?
- How to handle state management when user can go back/forward multiple times before completing?
- What if user goes back and forth rapidly - does state get corrupted?
- Should system remember failed attempts if user goes back and tries again?
- What happens to validation state when user returns to a previously validated step?
- How to handle "back" when previous step's options are dynamically loaded from API?

## Requirements *(mandatory)*

### Functional Requirements

**Research & Foundation**
- **FR-000**: System MUST research and document CLI UX patterns from 10+ popular developer tools before implementing navigation
- **FR-001**: System MUST implement navigation patterns based on documented best practices from CLI tool research

**Empty State Handling (All Commands)**
- **FR-002**: System MUST check for empty states before displaying interactive selection menus in ALL commands
- **FR-003**: System MUST display clear, actionable messages when no resources are available in ANY command (e.g., "No droplets found. Create one with 'cogo create'")
- **FR-004**: System MUST validate list/array bounds before accessing elements to prevent index out of range panics in ALL operations

**Back/Cancel Navigation (All Commands)**
- **FR-005**: Users MUST be able to press `Esc` key to go back one step in multi-step operations across ALL commands
- **FR-006**: When at first step, pressing `Esc` MUST exit the operation and return to command prompt
- **FR-007**: Users MUST be able to press `Ctrl+C` to immediately cancel and exit to command prompt from ANY operation at ANY step
- **FR-008**: System MUST preserve user selections when going back, allowing them to review and change previous choices
- **FR-009**: System MUST update subsequent steps when user changes previous selections (e.g., changing region updates available sizes)
- **FR-010**: System MUST prevent partial resource creation when user cancels mid-operation in ANY command

**Input Validation (All Commands)**
- **FR-011**: System MUST handle all keyboard input gracefully without crashing, even for unexpected or invalid input in ANY interactive prompt
- **FR-012**: System MUST validate state transitions when user navigates back and forth to prevent corrupted state

**Consistency (All Commands & Providers)**
- **FR-013**: System MUST provide identical navigation behavior across ALL commands (create, destroy, list, config, etc.)
- **FR-014**: System MUST provide identical navigation behavior across all cloud provider integrations
- **FR-015**: System MUST display consistent help text in ALL interactive menus indicating available navigation keys (e.g., "Press Esc to go back, Ctrl+C to cancel, ↑↓ to navigate")

**User Communication (All Commands)**
- **FR-016**: All confirmation prompts MUST offer clear yes/no options and handle negative responses by explicitly offering to return to previous step
- **FR-017**: System MUST log errors appropriately rather than displaying raw panic traces to users in ALL commands
- **FR-018**: System MUST gracefully handle terminal capabilities and provide fallback navigation methods if standard keys unavailable

### Key Entities

- **Navigation Context**: Represents current position in command flow with history stack (allows back navigation to previous steps while preserving state)
- **Command Flow**: Abstract representation of multi-step operations (create, destroy, config, etc.) with defined entry/exit points and step relationships
- **Cloud Provider Interface**: Abstract definition of navigation behaviors that all providers must implement consistently
- **Input Handler**: Validates and routes keyboard input, prevents invalid operations, manages back/cancel commands
- **Empty State Handler**: Detects and manages empty resource scenarios across all operations and commands
- **State Manager**: Tracks user selections across steps, handles updates when user goes back and changes choices
- **CLI UX Research Document**: Compiled findings from analyzing popular CLI tools, informing navigation design decisions

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: CLI UX research document completed covering 10+ popular tools with documented patterns and recommendations
- **SC-002**: Users can exit ANY operation (create, destroy, list, config, etc.) at any time without force-quitting the application (100% of interactive prompts support Esc/Ctrl+C)
- **SC-003**: Users can go back step-by-step in 100% of multi-step operations and change previous selections
- **SC-004**: System handles empty states without crashing in 100% of scenarios across ALL commands where resources don't exist
- **SC-005**: Zero panic/crash occurrences during random keyboard input testing across ALL commands (fuzz testing with 10,000 random key sequences per command)
- **SC-006**: Navigation patterns are identical across all commands and all cloud providers as measured by navigation flow diagram comparison
- **SC-007**: Users understand how to navigate within 30 seconds of first use (as indicated by help text visibility and consistency)
- **SC-008**: Support tickets related to "stuck in menu", "can't cancel", or "can't go back" reduced to zero
- **SC-009**: 95% of users successfully complete or cancel operations without confusion across ALL commands (measured by completion rate and user feedback)
- **SC-010**: Users who go back and change selections successfully complete operations 100% of the time without state corruption

## Assumptions

- Industry-standard CLI tools have established recognizable patterns that developers expect (will be validated through research)
- Users are familiar with basic terminal navigation conventions (arrow keys, Esc, Ctrl+C) common in popular CLI tools
- The `promptui` library can be configured or wrapped to support custom navigation patterns including back functionality
- Terminal emulators support standard ANSI escape sequences for keyboard input
- Canceling operations will not leave any partial state on cloud provider side (API calls are atomic or can be rolled back)
- Empty state detection can be performed before entering interactive menus across all commands
- Current codebase allows for refactoring to abstract common navigation patterns that apply to all commands
- State can be preserved when going back so users can see and modify previous selections
- Research phase will be completed before implementation begins (P0 prerequisite)

## Out of Scope

- Custom key binding configuration (users cannot remap navigation keys - will follow industry standards from research)
- Mouse/clicking support in terminal
- GUI or web-based interface
- Undo functionality for completed operations (different from back/cancel during operation)
- History or recent commands feature (e.g., bash history integration)
- Progress bars or loading indicators during API calls (separate feature)
- Multi-language support for navigation help text
- Accessibility features beyond standard terminal capabilities
- Auto-save/resume of partially completed operations
- Detailed analytics/telemetry on navigation patterns
- Migration of existing commands beyond refactoring for consistency
