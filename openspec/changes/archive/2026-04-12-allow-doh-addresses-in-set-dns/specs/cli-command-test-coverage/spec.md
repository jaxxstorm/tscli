## MODIFIED Requirements

### Requirement: Command validation behavior is tested
Each command with required inputs MUST have tests that verify missing or invalid input handling and actionable error messaging.

#### Scenario: Required flag is missing
- **WHEN** a command is executed without a required flag or argument
- **THEN** the command test MUST assert a non-nil execution error and an error message that identifies the missing input

#### Scenario: Invalid flag value is provided
- **WHEN** a command is executed with an invalid value (for example malformed route, tag, id format, or unsupported DNS nameserver input where validation exists)
- **THEN** the command test MUST assert failure behavior and the expected validation error

#### Scenario: DNS nameserver DoH input is accepted
- **WHEN** `set dns nameservers` or `set dns split-dns` is executed with a valid DNS-over-HTTPS endpoint address in an existing nameserver flag
- **THEN** command tests MUST assert that validation succeeds and the DoH value is preserved in the outgoing request payload
