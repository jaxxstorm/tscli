# openapi-property-coverage Specification

## Purpose
Define property-level OpenAPI coverage requirements so `tscli` tracks whether mapped operations cover or intentionally exclude request and response properties derived from the pinned schema.

## Requirements

### Requirement: Property coverage inventory is generated for mapped OpenAPI operations
The project SHALL generate a property coverage inventory for mapped OpenAPI operations from the pinned schema, including JSON request-body properties and successful JSON response-body properties.

#### Scenario: Property inventory is generated from mapped operations
- **WHEN** property coverage analysis runs against the pinned OpenAPI snapshot and command-operation map
- **THEN** the report SHALL enumerate request and response property paths for each mapped operation

#### Scenario: Nested response property is present in inventory
- **WHEN** a mapped operation response schema contains nested properties such as `devices[].postureIdentity`
- **THEN** the property inventory SHALL include that nested property path in the operation's response coverage set

### Requirement: Property coverage declarations are explicit and reviewable
Covered or intentionally excluded properties SHALL be declared in repository data with stable property-path identifiers and reviewer-visible rationale or evidence.

#### Scenario: Property is intentionally excluded
- **WHEN** a request or response property is not expected to be covered by the CLI
- **THEN** the repository SHALL record the property identifier and an exclusion rationale

#### Scenario: Property is marked covered
- **WHEN** a request or response property is handled intentionally by the CLI
- **THEN** the repository SHALL record that property as covered with stable evidence that reviewers can inspect

### Requirement: Property coverage is enforced in automated checks
Automated coverage checks SHALL report uncovered properties and fail when unresolved property gaps or property regressions remain.

#### Scenario: Uncovered property remains
- **WHEN** the property coverage report contains one or more uncovered in-scope properties
- **THEN** the coverage check SHALL fail and list the missing property paths

#### Scenario: Property regression is introduced
- **WHEN** a change increases uncovered properties relative to the baseline
- **THEN** the regression report SHALL identify the new uncovered property paths and CI SHALL fail

### Requirement: Remediated mapped operations no longer rely on blanket property exclusions
The project SHALL move mapped command operation sides selected for property-gap remediation out of blanket request/response exclusion mode and into explicit per-property coverage or narrowly justified exclusions.

#### Scenario: Remediated operation side is explicitly covered
- **WHEN** an existing mapped request or response side is targeted for property-gap remediation
- **THEN** property coverage data SHALL declare explicit coverage evidence for the schema properties the CLI parses instead of relying on the default side exclusion for that operation side

#### Scenario: Concrete blocker still requires an exclusion
- **WHEN** a remediated operation side still cannot support a documented property because of a known blocker
- **THEN** the exclusions file SHALL name that exact property path and rationale instead of excluding the whole request or response side

### Requirement: Shared schema properties are covered consistently across mapped commands
When multiple mapped commands expose the same upstream schema family, the coverage manifest SHALL use schema-aligned models so the same documented properties are evaluated consistently across those command surfaces.

#### Scenario: Device response family is audited
- **WHEN** `list devices`, `get device`, or route-related commands audit shared device or route response properties
- **THEN** fields such as `advertisedRoutes`, `multipleConnections`, `enabledRoutes`, and `postureIdentity` SHALL be covered or explicitly excluded consistently across the mapped operations that return them
