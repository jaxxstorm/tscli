## 1. Normalize OpenAPI snapshot refresh inputs

- [x] 1.1 Identify the canonical pinned OpenAPI schema and snapshot metadata paths used by the current coverage-gap workflow.
- [x] 1.2 Add or update the refresh helper logic needed to fetch `https://api.tailscale.com/api/v2?outputOpenapiSchema=true` and rewrite snapshot metadata atomically.
- [x] 1.3 Ensure refresh failures leave the previously pinned snapshot intact and return actionable errors.

## 2. Add make-driven latest-schema coverage validation

- [x] 2.1 Add a dedicated `Makefile` target for refreshing the pinned OpenAPI snapshot and metadata.
- [x] 2.2 Add a composed `Makefile` target that refreshes the snapshot and then runs coverage-gap validation against the refreshed snapshot.
- [x] 2.3 Preserve compatibility for existing `make coverage-gaps` and `make coverage-gaps-check` workflows.

## 3. Verify generated artifacts and coverage behavior

- [x] 3.1 Add or update tests around any new snapshot metadata or refresh helper behavior.
- [x] 3.2 Add or update coverage-gap tests or smoke verification to confirm the latest-schema workflow generates the expected artifacts and failure behavior.
- [x] 3.3 Run the relevant validation commands and inspect the refreshed coverage-gap output for newly exposed parity gaps.

## 4. Document the maintenance workflow

- [x] 4.1 Update `README.md` with the supported refresh-only and refresh-plus-coverage commands.
- [x] 4.2 Document the generated schema, metadata, and coverage-gap artifacts maintainers should review after a refresh.
- [x] 4.3 Note any baseline update or follow-up implementation steps required when the refreshed schema introduces new uncovered operations.
