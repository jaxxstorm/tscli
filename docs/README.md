# tscli Documentation

`tscli` is a Go-based CLI for interacting with the Tailscale API.

This docs site is the source of truth for setup, authentication, configuration, and command reference.

## Recommended reading paths

New users:

- [Getting Started](getting-started.md)
- [Authentication](authentication.md)
- [Configuration](configuration.md)


Existing users:

- [Command Reference](command-reference.md)
- [Configuration](configuration.md)

Contributors:

- [Command Reference](command-reference.md) for generation/check workflow
- [Configuration](configuration.md) for precedence and profile behavior

## Federated credentials

Use `tscli create key --type federated` to provision federated identities that mirror an OIDC issuer/subject pair. Provide `--scope` plus `--issuer` and `--subject`, and optionally `--audience`, `--tags`, and `--claim` to express custom claim rules.

After adding new key types or CLI verbs, keep the OpenAPI coverage mappings in sync by editing `pkg/contract/openapi/command-operation-map.yaml` and rerunning `make coverage-gaps-check`. The generated waterfall is saved in `coverage/coverage-gaps.*` so CI can detect regressions.

For references to the updated create-key documentation, see the generated [tscli create key](commands/tscli_create_key.md) command page.

Contributors, refer back to this file when you need to refresh the docs site or cover new CLI surface area.

Contributors:

- [Command Reference](command-reference.md) for generation/check workflow
- [Configuration](configuration.md) for precedence and profile behavior


For how the documentation site is built and regenerated, see [`docs/docs.md`](docs.md).
