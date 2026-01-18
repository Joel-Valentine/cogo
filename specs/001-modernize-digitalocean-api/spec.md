# Feature Specification: Modernize DigitalOcean API Integration

**Feature Branch**: `001-modernize-digitalocean-api`  
**Created**: 2026-01-18  
**Status**: Draft  
**Input**: User description: "update digitalocean/digitalocean.go to match todays digitalocean api. Its been a long time since this code was last worked on and I assume a lot of the api calls and even urls have changed since then"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Reliable Droplet Operations (Priority: P1)

As a Cogo user, I need all existing droplet operations (create, list, destroy) to continue working with the current DigitalOcean API so that I can manage my servers without interruption.

**Why this priority**: Core functionality must remain operational. If API calls fail due to outdated endpoints or deprecated methods, the tool becomes unusable.

**Independent Test**: Can be fully tested by performing each operation (create, list, destroy) against the current DigitalOcean API and verifying successful responses with correct data.

**Acceptance Scenarios**:

1. **Given** I have valid DigitalOcean API credentials, **When** I run `cogo create`, **Then** a new droplet is successfully created with the selected configuration
2. **Given** I have existing droplets in my account, **When** I run `cogo list`, **Then** all my droplets are displayed with current information (name, IP, status)
3. **Given** I select a droplet to destroy, **When** I confirm the deletion, **Then** the droplet is successfully removed from my account

---

### User Story 2 - Access to New DigitalOcean Features (Priority: P2)

As a Cogo user, I want access to newer DigitalOcean features and options (newer OS images, regions, droplet sizes) that have been added since the code was last updated so that I can take advantage of current offerings.

**Why this priority**: Users expect to access current cloud provider options. Missing new features limits the tool's usefulness and competitive value.

**Independent Test**: Can be tested by comparing available options (images, regions, sizes) in Cogo against the current DigitalOcean dashboard and verifying parity.

**Acceptance Scenarios**:

1. **Given** DigitalOcean has added new regions since 2020, **When** I select a region during droplet creation, **Then** I see all currently available regions
2. **Given** DigitalOcean offers new OS distributions and versions, **When** I select an image, **Then** I see current distributions including recent Ubuntu, Debian, and other OS releases  
3. **Given** DigitalOcean has introduced new droplet sizes and pricing tiers, **When** I select a droplet size, **Then** I see all current size options

---

### User Story 3 - Improved Error Handling (Priority: P3)

As a Cogo user, when API calls fail due to rate limiting, authentication issues, or service problems, I receive clear error messages that help me understand what went wrong and how to fix it.

**Why this priority**: Good error handling improves user experience and reduces support burden, but existing operations can function without enhanced error messages.

**Independent Test**: Can be tested by simulating various API error conditions (invalid token, rate limit, network issues) and verifying clear, actionable error messages are displayed.

**Acceptance Scenarios**:

1. **Given** my API token is invalid or expired, **When** I attempt any operation, **Then** I receive a clear message indicating authentication failure and instructions to update my token
2. **Given** I exceed DigitalOcean's API rate limits, **When** an operation fails, **Then** I receive a message explaining the rate limit and when I can retry
3. **Given** the DigitalOcean API is experiencing service issues, **When** an operation fails, **Then** I receive a message indicating the service is temporarily unavailable

---

### Edge Cases

- What happens when DigitalOcean introduces breaking API changes after the update?
- How does the system handle API deprecation warnings?
- What if new required parameters are added to existing endpoints?
- How are backward compatibility issues with old droplets handled?
- What happens if the API response structure changes significantly?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST successfully authenticate with the current DigitalOcean API v2
- **FR-002**: System MUST create droplets using current API endpoints and parameters
- **FR-003**: System MUST list droplets with all fields available in the current API response
- **FR-004**: System MUST destroy droplets using current API endpoints and handle confirmation workflows
- **FR-005**: System MUST retrieve and display current lists of: regions, images (distributions, applications, custom), sizes, and SSH keys
- **FR-006**: System MUST handle pagination for list operations according to current API specifications
- **FR-007**: System MUST use current godo SDK version or update HTTP client to match current API specifications
- **FR-008**: System MUST handle API errors with appropriate status codes and error messages
- **FR-009**: System MUST respect current API rate limits and provide feedback when limits are approached
- **FR-010**: System MUST maintain existing CLI command structure and user workflow (backward compatible UX)

### Key Entities

- **Droplet**: Represents a virtual machine instance with properties matching current DigitalOcean API schema (name, region, size, image, IP addresses, status, creation date)
- **Region**: Datacenter locations with current availability status and feature support
- **Image**: Operating system distributions, applications, or custom snapshots with current version information
- **Size**: Droplet configurations with current pricing, resources (CPU, RAM, disk), and availability
- **SSH Key**: Authentication credentials for server access with fingerprint and public key data

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All existing Cogo commands (create, list, destroy) complete successfully against the current DigitalOcean API without errors
- **SC-002**: Users can select from all regions, images, and sizes currently offered by DigitalOcean (100% feature parity with current API)
- **SC-003**: API operations complete within reasonable timeframes (list operations under 3 seconds for accounts with up to 100 droplets)
- **SC-004**: Error messages are actionable and clear, reducing user confusion by providing specific next steps for common failure scenarios
- **SC-005**: The update maintains backward compatibility - existing users can continue using Cogo without learning new commands or workflows

## Assumptions

- DigitalOcean API v2 remains the current stable version (v1 was deprecated years ago)
- The godo SDK is the recommended Go client library for DigitalOcean API
- API authentication continues to use bearer token method
- Core concepts (droplets, regions, images, sizes, SSH keys) remain fundamentally similar
- Breaking changes will be documented in DigitalOcean's API changelog
- Existing configuration file format (`.cogo` JSON file) can be maintained

## Dependencies

- Updated godo SDK library (or equivalent HTTP client if not using godo)
- Access to DigitalOcean API documentation and changelog
- Test DigitalOcean account for validation
- Current DigitalOcean API specification

## Out of Scope

- Adding support for new DigitalOcean services (volumes, load balancers, databases, Kubernetes) - only updating existing droplet management
- Changing the CLI command structure or user experience
- Adding new features beyond API modernization
- Supporting DigitalOcean API v1 (deprecated)
- Migration of existing droplets or data

