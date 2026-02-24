# API Parity Backlog

All previously uncovered in-scope operations are now implemented and mapped.

## create

### DeviceInvites (1)

- `post /device/{deviceId}/device-invites`

### DevicePosture (1)

- `post /tailnet/{tailnet}/posture/integrations`

### Keys (1)

- `post /tailnet/{tailnet}/keys`

### UserInvites (1)

- `post /tailnet/{tailnet}/user-invites`

### Webhooks (1)

- `post /tailnet/{tailnet}/webhooks`

## set

### Contacts (2)

- `patch /tailnet/{tailnet}/contacts/{contactType}`
- `post /tailnet/{tailnet}/contacts/{contactType}/resend-verification-email`

### DNS (6)

- `patch /tailnet/{tailnet}/dns/split-dns`
- `post /tailnet/{tailnet}/dns/configuration`
- `post /tailnet/{tailnet}/dns/nameservers`
- `post /tailnet/{tailnet}/dns/preferences`
- `post /tailnet/{tailnet}/dns/searchpaths`
- `put /tailnet/{tailnet}/dns/split-dns`

### DeviceInvites (2)

- `post /device-invites/-/accept`
- `post /device-invites/{deviceInviteId}/resend`

### DevicePosture (1)

- `patch /posture/integrations/{id}`

### Devices (9)

- `patch /tailnet/{tailnet}/device-attributes`
- `post /device/{deviceId}/attributes/{attributeKey}`
- `post /device/{deviceId}/authorized`
- `post /device/{deviceId}/expire`
- `post /device/{deviceId}/ip`
- `post /device/{deviceId}/key`
- `post /device/{deviceId}/name`
- `post /device/{deviceId}/routes`
- `post /device/{deviceId}/tags`

### Keys (1)

- `put /tailnet/{tailnet}/keys/{keyId}`

### Logging (1)

- `put /tailnet/{tailnet}/logging/{logType}/stream`

### PolicyFile (1)

- `post /tailnet/{tailnet}/acl`

### Services (2)

- `post /tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved`
- `put /tailnet/{tailnet}/services/{serviceName}`

### TailnetSettings (1)

- `patch /tailnet/{tailnet}/settings`

### UserInvites (1)

- `post /user-invites/{userInviteId}/resend`

### Users (4)

- `post /users/{userId}/approve`
- `post /users/{userId}/restore`
- `post /users/{userId}/role`
- `post /users/{userId}/suspend`

### Webhooks (3)

- `patch /webhooks/{endpointId}`
- `post /webhooks/{endpointId}/rotate`
- `post /webhooks/{endpointId}/test`

## get

### Contacts (1)

- `get /tailnet/{tailnet}/contacts`

### DNS (5)

- `get /tailnet/{tailnet}/dns/configuration`
- `get /tailnet/{tailnet}/dns/nameservers`
- `get /tailnet/{tailnet}/dns/preferences`
- `get /tailnet/{tailnet}/dns/searchpaths`
- `get /tailnet/{tailnet}/dns/split-dns`

### DeviceInvites (1)

- `get /device-invites/{deviceInviteId}`

### DevicePosture (1)

- `get /posture/integrations/{id}`

### Devices (2)

- `get /device/{deviceId}`
- `get /device/{deviceId}/attributes`

### Keys (1)

- `get /tailnet/{tailnet}/keys/{keyId}`

### Logging (4)

- `get /tailnet/{tailnet}/logging/{logType}/stream`
- `get /tailnet/{tailnet}/logging/{logType}/stream/status`
- `post /tailnet/{tailnet}/aws-external-id`
- `post /tailnet/{tailnet}/aws-external-id/{id}/validate-aws-trust-policy`

### PolicyFile (3)

- `get /tailnet/{tailnet}/acl`
- `post /tailnet/{tailnet}/acl/preview`
- `post /tailnet/{tailnet}/acl/validate`

### Services (2)

- `get /tailnet/{tailnet}/services/{serviceName}`
- `get /tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved`

### TailnetSettings (1)

- `get /tailnet/{tailnet}/settings`

### UserInvites (1)

- `get /user-invites/{userInviteId}`

### Users (1)

- `get /users/{userId}`

### Webhooks (1)

- `get /webhooks/{endpointId}`

## list

### DNS (1)

- `get /tailnet/{tailnet}/dns/nameservers`

### DeviceInvites (1)

- `get /device/{deviceId}/device-invites`

### DevicePosture (1)

- `get /tailnet/{tailnet}/posture/integrations`

### Devices (2)

- `get /device/{deviceId}/routes`
- `get /tailnet/{tailnet}/devices`

### Keys (1)

- `get /tailnet/{tailnet}/keys`

### Logging (2)

- `get /tailnet/{tailnet}/logging/configuration`
- `get /tailnet/{tailnet}/logging/network`

### Services (2)

- `get /tailnet/{tailnet}/services`
- `get /tailnet/{tailnet}/services/{serviceName}/devices`

### UserInvites (1)

- `get /tailnet/{tailnet}/user-invites`

### Users (1)

- `get /tailnet/{tailnet}/users`

### Webhooks (1)

- `get /tailnet/{tailnet}/webhooks`

## delete

### DeviceInvites (1)

- `delete /device-invites/{deviceInviteId}`

### DevicePosture (1)

- `delete /posture/integrations/{id}`

### Devices (2)

- `delete /device/{deviceId}`
- `delete /device/{deviceId}/attributes/{attributeKey}`

### Keys (1)

- `delete /tailnet/{tailnet}/keys/{keyId}`

### Logging (1)

- `delete /tailnet/{tailnet}/logging/{logType}/stream`

### Services (1)

- `delete /tailnet/{tailnet}/services/{serviceName}`

### UserInvites (1)

- `delete /user-invites/{userInviteId}`

### Users (1)

- `post /users/{userId}/delete`

### Webhooks (1)

- `delete /webhooks/{endpointId}`

