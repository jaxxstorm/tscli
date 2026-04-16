## 1. Persisted Preferences

- [x] 1.1 Update `pkg/config` persistence helpers so top-level `debug` values are preserved when config is rewritten while transient keys such as `help` and `base-url` remain excluded.
- [x] 1.2 Add or extend config persistence tests to cover saved `output` and `debug` values alongside canonical profile-backed config writes.
- [x] 1.3 Add or extend runtime precedence tests to verify saved `output` and `debug` defaults remain overridable by `--output`, `TSCLI_OUTPUT`, `--debug`, and `TSCLI_DEBUG`.

## 2. Setup Flow Changes

- [x] 2.1 Extend the `config setup` model with post-profile choice steps for output mode (`json`, `pretty`, `human`) and debug preference.
- [x] 2.2 Gate the new prompts so they run after initial profile creation completes, but do not interrupt existing rerun profile-management flows.
- [x] 2.3 Persist the selected output/debug preferences at the end of initial setup without regressing profile or encryption persistence.
- [x] 2.4 Update setup rendering and any command help or example text needed to reflect the new initial-setup prompts.

## 3. Verification

- [x] 3.1 Add unit tests for setup-model transitions covering the initial-setup path from profile creation through output and debug prompts.
- [x] 3.2 Add command-level integration coverage for first-run `config setup` persisting profile, `output`, and `debug` selections.
- [x] 3.3 Add command-level integration coverage confirming rerun `config setup` profile-management flows do not unexpectedly require the new preference prompts.
