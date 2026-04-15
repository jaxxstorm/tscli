## ADDED Requirements

### Requirement: User deletion parity covers bulk deletion workflow
The CLI SHALL provide parity for user deletion workflows by exposing both single-user deletion and filter-driven bulk user deletion commands backed by the in-scope user deletion operations.

#### Scenario: Bulk user deletion is represented in parity coverage
- **WHEN** parity analysis evaluates uncovered user-domain operations and workflows
- **THEN** the `delete users` command SHALL be treated as the supported CLI path for bulk user cleanup based on user-list filtering plus per-user deletion requests

#### Scenario: Bulk user deletion command is mapped explicitly
- **WHEN** the `tscli delete users` command is implemented
- **THEN** command-to-operation mapping and parity checks SHALL include the command and its associated user-domain operations
