## ADDED Requirements

### Requirement: Service commands render stable structured output
`tscli list services` and `tscli get service` SHALL render service payloads as stable structured records across supported output modes without flattening or collapsing sibling service entries together.

#### Scenario: List services pretty output shows one service per record
- **WHEN** `tscli list services` runs in `pretty` mode against a response containing multiple services
- **THEN** the CLI SHALL render each service as a distinct record separated consistently from the next service
- **AND** fields from one service SHALL NOT be merged into or repeated under another service's rendered block

#### Scenario: List services preserves key service properties
- **WHEN** `tscli list services` receives service objects containing properties such as `name`, `addrs`, `ports`, `tags`, `comment`, and `annotations`
- **THEN** the CLI SHALL preserve those properties in structured output modes
- **AND** nested map fields such as `annotations` SHALL remain associated with the correct service record

#### Scenario: Get service pretty output shows a single service record
- **WHEN** `tscli get service --service <name>` runs in `pretty` or `human` mode
- **THEN** the CLI SHALL render the returned service as one structured service record
- **AND** the rendered output SHALL use stable field formatting consistent with other structured `get` commands

#### Scenario: Structured output remains script-friendly
- **WHEN** `tscli list services` or `tscli get service` runs in `json` or `yaml` mode
- **THEN** the CLI SHALL preserve the API response structure without synthetic summaries or lossy field dropping
- **AND** correcting pretty or human rendering SHALL NOT change command flags or request semantics
