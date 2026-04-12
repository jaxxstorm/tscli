## 1. Service response shaping

- [x] 1.1 Inspect the `/services` and `/services/{serviceName}` response shapes and add service-oriented decoding/normalization for `list services` and `get service`.
- [x] 1.2 Update the service commands so pretty and human rendering operate on stable service records while JSON and yaml remain schema-aligned and script-friendly.
- [x] 1.3 Review adjacent service output paths for the same raw-response issue and fix any obvious formatting regressions discovered during implementation.

## 2. Output behavior verification

- [x] 2.1 Expand service fixtures to include representative fields such as addresses, ports, tags, comments, and annotations.
- [x] 2.2 Add command-level example/output tests for `list services` and `get service` that assert readable pretty and human rendering and stable JSON/yaml structure.
- [x] 2.3 Verify the updated service output behavior does not regress command exit semantics or compatibility for existing flags and scripts.

## 3. Documentation and validation

- [x] 3.1 Update any command docs or examples whose rendered output expectations change because of the service pretty/human fix.
- [x] 3.2 Run the relevant service command and docs test suites and capture any follow-up output issues uncovered by the richer fixtures.
