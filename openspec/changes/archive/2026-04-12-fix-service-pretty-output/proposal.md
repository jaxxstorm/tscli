## Why

`tscli list services` and `tscli get service` currently produce incorrect pretty-mode output for real Tailscale service payloads, flattening nested service data into a hard-to-read stream instead of rendering one service record at a time. This needs to be fixed now because the commands are already exposed to users and the current output makes service inspection unreliable in the interactive modes people use most.

## What Changes

- Fix `list services` pretty and human output so service collections render as distinct service records instead of a collapsed nested blob.
- Fix `get service` pretty and human output so a single service renders with the same stable field formatting as other structured commands.
- Normalize service response decoding/output handling so collection and single-service commands preserve service fields such as addresses, ports, tags, comments, and annotations.
- Add output-focused regression coverage for service commands across supported output modes, with emphasis on pretty/human readability and stable JSON/yaml structure.
- Review adjacent service command output behavior for obvious formatting inconsistencies introduced by the same raw-response path and correct them where needed without changing command flags or API semantics.

## Capabilities

### New Capabilities

- `service-command-output`: Define stable user-facing output behavior for `list services` and `get service`, especially in pretty and human modes for nested service payloads.

### Modified Capabilities

- `cli-command-test-coverage`: Tighten command output coverage so service commands exercise supported output modes and catch regressions in rendered structure, not just JSON shape.

## Impact

- Affected command groups: `list services`, `get service`, and any shared service output helpers introduced for pretty/human rendering.
- Affected code: service CLI command implementations under `cmd/tscli/...`, shared output logic under `pkg/output`, and mock/example output tests under `test/cli`.
- Backward compatibility: command paths, flags, config keys, and API requests remain unchanged; the user-visible impact is corrected pretty/human rendering and stronger regression coverage for existing scripts and interactive use.
