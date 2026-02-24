## Purpose

Define mock-backed integration test requirements for `tscli` commands so API behavior is validated without live network dependencies.

## Requirements

### Requirement: Mocked API integration tests for command behavior
The project MUST provide integration tests that execute CLI commands against a mocked Tailscale API server instead of live network endpoints.

#### Scenario: Command uses mock API endpoint
- **WHEN** an integration test executes a command that performs API calls
- **THEN** all HTTP traffic MUST target the test mock server base URL

#### Scenario: Live API access is blocked in integration tests
- **WHEN** integration tests run in CI or local test mode
- **THEN** tests MUST fail if a command attempts to reach non-mock Tailscale API hosts

### Requirement: Endpoint fixtures cover success and failure behavior
Mock API fixtures MUST support deterministic success responses and representative error responses for each tested command group.

#### Scenario: Success fixture drives expected output
- **WHEN** a mock endpoint returns a successful response fixture
- **THEN** the corresponding command integration test MUST assert expected stdout and nil execution error

#### Scenario: Error fixture drives expected error handling
- **WHEN** a mock endpoint returns API error responses (4xx/5xx)
- **THEN** the corresponding command integration test MUST assert expected error propagation and stderr behavior

### Requirement: Request shape assertions are enforced
Integration tests MUST assert critical outbound request properties for commands, including HTTP method, path, query, and key request body fields.

#### Scenario: Command sends expected request contract
- **WHEN** a command executes against the mock API
- **THEN** the test harness MUST assert expected method/path and required request payload fields before returning the fixture response

### Requirement: Integration tests are runnable as a stable test target
The project MUST provide a documented and deterministic command/target to run mock-backed CLI integration tests.

#### Scenario: Integration suite invocation
- **WHEN** a developer or CI job runs the documented integration test target
- **THEN** the full mock-backed command integration suite MUST execute without requiring external credentials or network access
