## ADDED Requirements

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
