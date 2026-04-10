## MODIFIED Requirements

### Requirement: Contract validation against pinned upstream API schema/examples
The project MUST validate CLI request/response model handling against a pinned snapshot of official Tailscale API schema or authoritative API examples at the property level, not only at the endpoint or whole-payload shape level.

#### Scenario: Contract snapshot is available
- **WHEN** contract tests run
- **THEN** tests MUST validate key request/response property paths and payload structures against the pinned schema/examples

#### Scenario: Contract mismatch is introduced
- **WHEN** code changes introduce request/response model drift from the pinned contract
- **THEN** contract tests MUST fail with the mismatched property or shape details

#### Scenario: Pinned schema adds a response property
- **WHEN** the pinned schema introduces a new response property for a mapped command model such as `devices[].postureIdentity`
- **THEN** contract validation or coverage checks MUST surface that property until the CLI coverage declaration and representative tests are updated

### Requirement: Command-level model decoding failures are covered
Integration tests MUST verify that invalid or unexpected API payload shapes produce clear command errors instead of silent corruption, and representative command fixtures MUST exercise important covered response properties.

#### Scenario: Unexpected payload field shape
- **WHEN** a mock API response returns incompatible field types for a command model
- **THEN** the command MUST return a surfaced error and integration tests MUST assert the failure path

#### Scenario: Representative covered property is exercised
- **WHEN** a command claims coverage for an important response property in automated coverage data
- **THEN** contract or integration tests MUST include a fixture/assertion that exercises that property through decode or structured output
