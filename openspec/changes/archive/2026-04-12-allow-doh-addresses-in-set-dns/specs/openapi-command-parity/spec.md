## MODIFIED Requirements

### Requirement: Missing operation is implemented
- **WHEN** an in-scope operation is classified as missing-command
- **THEN** a CLI command SHALL be added that invokes that operation with validated flags and structured output/error handling

#### Scenario: Existing DNS mutation command matches supported nameserver formats
- **WHEN** the upstream DNS nameserver API accepts DNS-over-HTTPS endpoint addresses in addition to literal IP nameservers
- **THEN** the existing `set dns nameservers` and `set dns split-dns` commands SHALL accept those same supported nameserver value formats without requiring a raw API fallback
