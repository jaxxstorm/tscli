---
name: openapi-refresh
description: Refresh tscli's pinned Tailscale OpenAPI snapshot before latest-schema implementation or coverage-gap analysis.
license: MIT
---

Use this skill when work depends on the latest upstream Tailscale OpenAPI schema, snapshot metadata, model contract consistency, or newly documented API operations/properties.

## Steps

1. Run the supported OpenAPI refresh workflow before deriving latest-schema gaps:

   ```bash
   make openapi-refresh
   ```

2. Inspect refreshed artifacts:

   - `pkg/contract/openapi/tailscale-v2-openapi.yaml`
   - `pkg/contract/openapi/snapshot-metadata.yaml`

3. Run coverage-gap analysis against the refreshed snapshot:

   ```bash
   make coverage-gaps-check
   ```

4. Use `coverage/coverage-gaps.md`, `coverage/coverage-gaps.json`, and `coverage/coverage-gaps-diff.md` to drive follow-up implementation.

## Failure Handling

- If network access is unavailable, report that the refresh could not fetch the upstream schema and continue only with the pinned in-repo snapshot.
- If refresh changes the pinned schema or metadata, include those artifacts in the review alongside any command, mapping, property coverage, and test changes.
- If the refreshed schema is incompatible with the tooling, fix the refresh or coverage tooling before implementing commands from that schema.
