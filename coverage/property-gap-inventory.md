# Property Gap Inventory

This inventory captures the mapped operation sides that still relied on default or explicit property exclusions when `fix-openapi-property-gaps` started.

## Device family

- `get /tailnet/{tailnet}/devices response`: shared device schema wrapped in the `devices` envelope, with explicit exclusions for `devices[].advertisedRoutes`, `devices[].advertisedRoutes[]`, and `devices[].multipleConnections`
- `get /device/{deviceId} response`: full device schema returned directly, still default-excluded

## Route family

- `get /device/{deviceId}/routes response`: shared `DeviceRoutes` response, still default-excluded
- `post /device/{deviceId}/routes request`: request body with the `routes[]` payload, still default-excluded
- `post /device/{deviceId}/routes response`: shared `DeviceRoutes` response, still default-excluded because the command returned a synthetic summary

## Settings family

- `get /tailnet/{tailnet}/settings response`: already declared in coverage data, but runtime decoding still depended on SDK models that could flatten nullable values
- `patch /tailnet/{tailnet}/settings request`: request-side property coverage already existed
- `patch /tailnet/{tailnet}/settings response`: still default-excluded because the command printed the outbound request instead of the API response

## Deferred families

- device invites and posture attribute responses
- service approval and service device responses
- other mapped request/response bodies that remain on default exclusions outside this remediation batch
