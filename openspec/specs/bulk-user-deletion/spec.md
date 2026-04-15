# bulk-user-deletion Specification

## Purpose
Define the `tscli delete users` bulk deletion capability, including validated filters, dry-run behavior, admin safeguards, and structured output expectations.

## Requirements
### Requirement: Bulk user deletion command
The CLI SHALL provide a `tscli delete users` command that deletes multiple tailnet users selected from the user list API and evaluated against validated deletion filters.

#### Scenario: Command is available under delete group
- **WHEN** a user runs `tscli delete --help`
- **THEN** the help output SHALL list `users` as a delete subcommand alongside the existing `user` command

#### Scenario: Command lists matching deletion candidates by default
- **WHEN** a user runs `tscli delete users` with valid selection flags and without `--confirm`
- **THEN** the command SHALL perform a dry run that reports which users match the filters without deleting them

#### Scenario: Command deletes matching users when confirmed
- **WHEN** a user runs `tscli delete users` with valid selection flags and `--confirm`
- **THEN** the command SHALL delete each matching user and report per-user success or failure in structured output

### Requirement: User selection filters
The `tscli delete users` command SHALL support filtering by `--status`, `--last-seen`, and `--devices`, where `--status` and `--last-seen` are mutually exclusive and `--devices` MAY be used either as an additional constraint or as a standalone filter.

#### Scenario: Suspended users can be selected by status
- **WHEN** a user runs `tscli delete users --status suspended`
- **THEN** only users whose API `status` field is `suspended` SHALL be considered deletion candidates

#### Scenario: Inactive users can be selected by status
- **WHEN** a user runs `tscli delete users --status inactive`
- **THEN** only users whose API `status` field is `inactive` SHALL be considered deletion candidates

#### Scenario: Users can be selected by last seen age
- **WHEN** a user runs `tscli delete users --last-seen 24h`
- **THEN** only users whose API `lastSeen` timestamp is older than the provided duration relative to command execution time SHALL be considered deletion candidates

#### Scenario: Users can be constrained by device count
- **WHEN** a user runs `tscli delete users --last-seen 24h --devices 0`
- **THEN** only users matching the primary activity filter and whose API `deviceCount` equals `0` SHALL be considered deletion candidates

#### Scenario: Users can be selected by device count alone
- **WHEN** a user runs `tscli delete users --devices 0`
- **THEN** only users whose API `deviceCount` equals `0` SHALL be considered deletion candidates

#### Scenario: Status and last seen are mutually exclusive
- **WHEN** a user runs `tscli delete users --status suspended --last-seen 24h`
- **THEN** the command SHALL fail before making API mutations and SHALL report that `--status` and `--last-seen` are mutually exclusive

#### Scenario: Invalid status value is rejected
- **WHEN** a user runs `tscli delete users --status active`
- **THEN** the command SHALL fail before making API mutations and SHALL report that only `inactive` and `suspended` are supported `--status` values

### Requirement: Admin users are protected by default
The `tscli delete users` command SHALL exclude admin users from deletion candidates unless the caller explicitly opts in with `--admins=true`.

#### Scenario: Admin users are excluded by default
- **WHEN** a matching user has an API `role` value of `admin` and the caller does not set `--admins`
- **THEN** that user SHALL be excluded from deletion candidates and identified as skipped in command output

#### Scenario: Admin users can be included explicitly
- **WHEN** a user runs `tscli delete users --last-seen 24h --admins=true`
- **THEN** matching users with an API `role` value of `admin` SHALL be eligible deletion candidates

### Requirement: Output is structured and testable
The `tscli delete users` command SHALL produce machine-readable output in supported formats and SHALL include enough detail to distinguish deleted, failed, and skipped users.

#### Scenario: JSON output includes deletion summary
- **WHEN** a user runs `tscli delete users --status suspended --output json`
- **THEN** the JSON output SHALL include aggregate counts and per-user result details for matched, failed, and skipped users

#### Scenario: Unit and integration coverage exists
- **WHEN** the bulk deletion capability is implemented
- **THEN** automated tests SHALL cover filter validation, admin exclusion defaults, dry-run behavior, and successful deletion flows
