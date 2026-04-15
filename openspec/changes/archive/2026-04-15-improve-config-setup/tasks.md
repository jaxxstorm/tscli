## 1. Command and setup flow scaffolding

- [x] 1.1 Add a new `tscli config setup` Cobra command and register it under the existing `config` command group without removing current `config profiles` or `config encryption` commands.
- [x] 1.2 Implement a Bubble Tea setup model that covers first-run onboarding, rerun management choices, auth-type selection, profile entry, add-another confirmation, and graceful exit states.
- [x] 1.3 Load existing profiles at setup start so reruns can branch into add or delete flows based on persisted config state.

## 2. Encryption bootstrap and persistence

- [x] 2.1 Add setup actions that ask whether credentials should be encrypted and, when enabled, generate a new `age` identity in-process.
- [x] 2.2 Implement default and custom key-path handling, including parent directory creation and secure key file writing for `~/.tscli/age.txt` or the user-selected path.
- [x] 2.3 Persist generated encryption settings through existing config helpers so `encryption.age.public-key` and `encryption.age.private-key-path` are saved only after key material is written successfully.

## 3. Interactive profile management

- [x] 3.1 Implement API-key and OAuth profile entry screens that collect the required values and call existing profile upsert helpers.
- [x] 3.2 Ensure setup-created profiles use canonical persisted shape and store encrypted secret fields when encryption is enabled.
- [x] 3.3 Implement rerun delete flow that lists removable profiles, blocks active-profile deletion through existing validation, and persists the resulting profile set.

## 4. Validation and compatibility

- [x] 4.1 Extend profile/config validation coverage for encrypted secret field combinations and any setup-specific path validation needed by the new flow.
- [x] 4.2 Verify that `config setup` does not change runtime credential precedence or break existing `config profiles` and `config encryption setup` workflows.
- [x] 4.3 Update command help text or user-facing docs for the new setup behavior, including encryption defaults and rerun management behavior.

## 5. Test coverage

- [x] 5.1 Add unit tests for Bubble Tea state transitions, default key path resolution, custom path handling, and graceful exit behavior.
- [x] 5.2 Add integration tests for first-run plain-text setup and first-run encrypted setup, including directory creation and persisted config assertions.
- [x] 5.3 Add integration tests for rerun add/delete profile flows and encrypted secret persistence for both API-key-backed and OAuth-backed profiles.
