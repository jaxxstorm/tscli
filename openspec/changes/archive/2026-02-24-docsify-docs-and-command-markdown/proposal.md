## Why

`tscli` documentation is currently concentrated in `README.md`, which makes command reference discovery and configuration/auth guidance harder to navigate as the CLI surface grows. A dedicated docs site with generated command reference pages will keep docs consistent with the CLI and reduce drift.

## What Changes

- Create a `docs/` site powered by Docsify for end-user and contributor documentation.
- Add automated command reference generation from the Cobra command tree into Markdown files (one page per command path).
- Add docs pages for configuration structure, config precedence, and multi-tailnet profile usage.
- Add docs pages for API key authentication, including flag/env/config approaches and security guidance.
- Add repeatable generation workflow (script/Make target) so command docs can be refreshed as commands change.
- Wire documentation generation and validation into CI-friendly workflows where practical.

## Capabilities

### New Capabilities
- `docsify-documentation-site`: Provide a Docsify-backed documentation site with navigation and structured docs content for tscli.
- `generated-command-reference-markdown`: Auto-generate Markdown reference docs from Cobra commands so command docs track CLI behavior.
- `config-and-auth-documentation`: Document config schema, precedence, active tailnet profiles, and API key authentication paths.

### Modified Capabilities
- None.

## Impact

- Affected code:
  - `docs/**` (new documentation site content, index, sidebar, theme config)
  - documentation generation tooling under project scripts or command utilities
  - command wiring for doc generation (if implemented as a CLI/helper command)
  - CI/Make targets for docs generation checks
- Affected command groups/flags/config/env docs:
  - Generated docs for all `tscli` command groups including `config`, `create`, `get`, `list`, `set`, `delete`, `version`
  - Auth guidance for `--api-key`, `TAILSCALE_API_KEY`, and config `api-key`
  - Config guidance for `tailnet`, `tailnets`, `active-tailnet`, output settings, and precedence
- Backward compatibility:
  - No breaking CLI behavior changes expected; this is a documentation and tooling enhancement.
