## ADDED Requirements

### Requirement: CLI supports persisted default output and debug preferences
The CLI SHALL support top-level persisted `output` and `debug` config keys as default runtime preferences. `output` SHALL accept setup-selected values `json`, `pretty`, and `human`. `debug` SHALL be persisted as a boolean that controls default HTTP request/response logging. For operational commands, effective values SHALL resolve in this order: command flags, environment variables, persisted config values, built-in defaults.

#### Scenario: Persisted output default is used when no override exists
- **WHEN** a config file sets `output: pretty` and a command runs without `--output` or `TSCLI_OUTPUT`
- **THEN** the CLI SHALL use `pretty` as the effective output mode

#### Scenario: Persisted debug default is used when no override exists
- **WHEN** a config file sets `debug: true` and a command runs without `--debug` or `TSCLI_DEBUG`
- **THEN** the CLI SHALL enable debug HTTP request/response logging for that command

#### Scenario: Output flag overrides persisted default
- **WHEN** a config file sets `output: pretty` and a command runs with `--output json`
- **THEN** the CLI SHALL use `json` as the effective output mode

#### Scenario: Output environment overrides persisted default
- **WHEN** a config file sets `output: pretty` and a command runs with `TSCLI_OUTPUT=human`
- **THEN** the CLI SHALL use `human` as the effective output mode

#### Scenario: Debug environment overrides persisted default
- **WHEN** a config file sets `debug: false` and a command runs with `TSCLI_DEBUG=1`
- **THEN** the CLI SHALL enable debug HTTP request/response logging for that command

#### Scenario: Missing persisted preferences fall back to built-in defaults
- **WHEN** `output` and `debug` are absent from config and no flag or environment override is set
- **THEN** the CLI SHALL keep its existing built-in runtime defaults
