# Coverage Gaps Report

- OpenAPI operations: `85`
- Excluded operations: `0`
- In-scope operations: `85`
- Manifest commands: `97`
- Excluded commands: `16`
- Covered operations: `85`
- Uncovered operations: `0`
- Covered commands: `81`
- Unmapped commands: `0`
- Unknown mapped commands: `0`
- Covered properties: `163`
- Excluded properties: `637`
- Uncovered properties: `0`

## Uncovered Operations By Domain

- None
## Unmapped Commands

- None

## Unknown Mapped Operations

- None

## Unknown Mapped Commands

- None

## Covered Properties By Operation

### get /device/{deviceId} response

- `addresses`
- `addresses[]`
- `advertisedRoutes`
- `advertisedRoutes[]`
- `authorized`
- `blocksIncomingConnections`
- `clientConnectivity`
- `clientConnectivity.clientSupports`
- `clientConnectivity.clientSupports.hairPinning`
- `clientConnectivity.clientSupports.ipv6`
- `clientConnectivity.clientSupports.pcp`
- `clientConnectivity.clientSupports.pmp`
- `clientConnectivity.clientSupports.udp`
- `clientConnectivity.clientSupports.upnp`
- `clientConnectivity.endpoints`
- `clientConnectivity.endpoints[]`
- `clientConnectivity.latency`
- `clientConnectivity.mappingVariesByDestIP`
- `clientVersion`
- `connectedToControl`
- `created`
- `distro`
- `distro.codeName`
- `distro.name`
- `distro.version`
- `enabledRoutes`
- `enabledRoutes[]`
- `expires`
- `hostname`
- `id`
- `isEphemeral`
- `isExternal`
- `keyExpiryDisabled`
- `lastSeen`
- `machineKey`
- `multipleConnections`
- `name`
- `nodeId`
- `nodeKey`
- `os`
- `postureIdentity`
- `postureIdentity.disabled`
- `postureIdentity.serialNumbers`
- `postureIdentity.serialNumbers[]`
- `sshEnabled`
- `tags`
- `tags[]`
- `tailnetLockError`
- `tailnetLockKey`
- `updateAvailable`
- `user`

### get /device/{deviceId}/routes response

- `advertisedRoutes`
- `advertisedRoutes[]`
- `enabledRoutes`
- `enabledRoutes[]`

### get /tailnet/{tailnet}/devices response

- `devices`
- `devices[]`
- `devices[].addresses`
- `devices[].addresses[]`
- `devices[].advertisedRoutes`
- `devices[].advertisedRoutes[]`
- `devices[].authorized`
- `devices[].blocksIncomingConnections`
- `devices[].clientConnectivity`
- `devices[].clientConnectivity.clientSupports`
- `devices[].clientConnectivity.clientSupports.hairPinning`
- `devices[].clientConnectivity.clientSupports.ipv6`
- `devices[].clientConnectivity.clientSupports.pcp`
- `devices[].clientConnectivity.clientSupports.pmp`
- `devices[].clientConnectivity.clientSupports.udp`
- `devices[].clientConnectivity.clientSupports.upnp`
- `devices[].clientConnectivity.endpoints`
- `devices[].clientConnectivity.endpoints[]`
- `devices[].clientConnectivity.latency`
- `devices[].clientConnectivity.mappingVariesByDestIP`
- `devices[].clientVersion`
- `devices[].connectedToControl`
- `devices[].created`
- `devices[].distro`
- `devices[].distro.codeName`
- `devices[].distro.name`
- `devices[].distro.version`
- `devices[].enabledRoutes`
- `devices[].enabledRoutes[]`
- `devices[].expires`
- `devices[].hostname`
- `devices[].id`
- `devices[].isEphemeral`
- `devices[].isExternal`
- `devices[].keyExpiryDisabled`
- `devices[].lastSeen`
- `devices[].machineKey`
- `devices[].multipleConnections`
- `devices[].name`
- `devices[].nodeId`
- `devices[].nodeKey`
- `devices[].os`
- `devices[].postureIdentity`
- `devices[].postureIdentity.disabled`
- `devices[].postureIdentity.serialNumbers`
- `devices[].postureIdentity.serialNumbers[]`
- `devices[].sshEnabled`
- `devices[].tags`
- `devices[].tags[]`
- `devices[].tailnetLockError`
- `devices[].tailnetLockKey`
- `devices[].updateAvailable`
- `devices[].user`

### get /tailnet/{tailnet}/settings response

- `aclsExternalLink`
- `aclsExternallyManagedOn`
- `devicesApprovalOn`
- `devicesAutoUpdatesOn`
- `devicesKeyDurationDays`
- `httpsEnabled`
- `networkFlowLoggingOn`
- `postureIdentityCollectionOn`
- `regionalRoutingOn`
- `usersApprovalOn`
- `usersRoleAllowedToJoinExternalTailnets`

### patch /tailnet/{tailnet}/settings request

- `devicesApprovalOn`
- `devicesAutoUpdatesOn`
- `devicesKeyDurationDays`
- `networkFlowLoggingOn`
- `postureIdentityCollectionOn`
- `regionalRoutingOn`
- `usersApprovalOn`
- `usersRoleAllowedToJoinExternalTailnets`

### patch /tailnet/{tailnet}/settings response

- `aclsExternalLink`
- `aclsExternallyManagedOn`
- `devicesApprovalOn`
- `devicesAutoUpdatesOn`
- `devicesKeyDurationDays`
- `httpsEnabled`
- `networkFlowLoggingOn`
- `postureIdentityCollectionOn`
- `regionalRoutingOn`
- `usersApprovalOn`
- `usersRoleAllowedToJoinExternalTailnets`

### post /device/{deviceId}/routes request

- `routes`
- `routes[]`

### post /device/{deviceId}/routes response

- `advertisedRoutes`
- `advertisedRoutes[]`
- `enabledRoutes`
- `enabledRoutes[]`

### post /tailnet/{tailnet}/keys request

- `audience`
- `capabilities`
- `capabilities.devices`
- `capabilities.devices.create`
- `capabilities.devices.create.ephemeral`
- `capabilities.devices.create.preauthorized`
- `capabilities.devices.create.reusable`
- `capabilities.devices.create.tags`
- `capabilities.devices.create.tags[]`
- `customClaimRules`
- `description`
- `expirySeconds`
- `issuer`
- `keyType`
- `scopes`
- `scopes[]`
- `subject`
- `tags`
- `tags[]`


## Excluded Properties By Operation

### get /device-invites/{deviceInviteId} response

- `accepted`
- `acceptedBy`
- `acceptedBy.id`
- `acceptedBy.loginName`
- `acceptedBy.profilePicUrl`
- `allowExitNode`
- `created`
- `deviceId`
- `email`
- `id`
- `inviteUrl`
- `lastEmailSentAt`
- `multiUse`
- `sharerId`
- `tailnetId`

### get /device/{deviceId}/attributes response

- `attributes`
- `expiries`

### get /device/{deviceId}/device-invites response

- `[].accepted`
- `[].acceptedBy`
- `[].acceptedBy.id`
- `[].acceptedBy.loginName`
- `[].acceptedBy.profilePicUrl`
- `[].allowExitNode`
- `[].created`
- `[].deviceId`
- `[].email`
- `[].id`
- `[].inviteUrl`
- `[].lastEmailSentAt`
- `[].multiUse`
- `[].sharerId`
- `[].tailnetId`

### get /posture/integrations/{id} response

- `clientId`
- `clientSecret`
- `cloudId`
- `configUpdated`
- `id`
- `provider`
- `status`
- `status.error`
- `status.lastSync`
- `status.matchedCount`
- `status.possibleMatchedCount`
- `status.providerHostCount`
- `tenantId`

### get /tailnet/{tailnet}/contacts response

- `account`
- `account.email`
- `account.fallbackEmail`
- `account.needsVerification`
- `security`
- `security.email`
- `security.fallbackEmail`
- `security.needsVerification`
- `support`
- `support.email`
- `support.fallbackEmail`
- `support.needsVerification`

### get /tailnet/{tailnet}/dns/configuration response

- `nameservers`
- `nameservers[]`
- `nameservers[].address`
- `nameservers[].useWithExitNode`
- `preferences`
- `preferences.magicDNS`
- `preferences.overrideLocalDNS`
- `searchPaths`
- `searchPaths[]`
- `splitDNS`

### get /tailnet/{tailnet}/dns/nameservers response

- `dns`
- `dns[]`

### get /tailnet/{tailnet}/dns/preferences response

- `magicDNS`

### get /tailnet/{tailnet}/dns/searchpaths response

- `searchPaths`
- `searchPaths[]`

### get /tailnet/{tailnet}/keys response

- `keys`
- `keys[]`
- `keys[].audience`
- `keys[].capabilities`
- `keys[].capabilities.devices`
- `keys[].capabilities.devices.create`
- `keys[].capabilities.devices.create.ephemeral`
- `keys[].capabilities.devices.create.preauthorized`
- `keys[].capabilities.devices.create.reusable`
- `keys[].capabilities.devices.create.tags`
- `keys[].capabilities.devices.create.tags[]`
- `keys[].created`
- `keys[].customClaimRules`
- `keys[].description`
- `keys[].expires`
- `keys[].expirySeconds`
- `keys[].id`
- `keys[].invalid`
- `keys[].issuer`
- `keys[].key`
- `keys[].keyType`
- `keys[].revoked`
- `keys[].scopes`
- `keys[].scopes[]`
- `keys[].subject`
- `keys[].tags`
- `keys[].tags[]`
- `keys[].updated`
- `keys[].userId`

### get /tailnet/{tailnet}/keys/{keyId} response

- `audience`
- `capabilities`
- `capabilities.devices`
- `capabilities.devices.create`
- `capabilities.devices.create.ephemeral`
- `capabilities.devices.create.preauthorized`
- `capabilities.devices.create.reusable`
- `capabilities.devices.create.tags`
- `capabilities.devices.create.tags[]`
- `created`
- `customClaimRules`
- `description`
- `expires`
- `expirySeconds`
- `id`
- `invalid`
- `issuer`
- `key`
- `keyType`
- `revoked`
- `scopes`
- `scopes[]`
- `subject`
- `tags`
- `tags[]`
- `updated`
- `userId`

### get /tailnet/{tailnet}/logging/configuration response

- `logs`
- `logs[]`
- `logs[].action`
- `logs[].actionDetails`
- `logs[].actor`
- `logs[].actor.displayName`
- `logs[].actor.id`
- `logs[].actor.loginName`
- `logs[].actor.tags`
- `logs[].actor.tags[]`
- `logs[].actor.type`
- `logs[].deferredAt`
- `logs[].error`
- `logs[].eventGroupID`
- `logs[].eventTime`
- `logs[].new`
- `logs[].new[]`
- `logs[].old`
- `logs[].old[]`
- `logs[].origin`
- `logs[].target`
- `logs[].target.id`
- `logs[].target.isEphemeral`
- `logs[].target.name`
- `logs[].target.property`
- `logs[].target.type`
- `logs[].type`
- `tailnet`
- `version`

### get /tailnet/{tailnet}/logging/network response

- `logs`
- `logs[]`
- `logs[].end`
- `logs[].exitTraffic`
- `logs[].exitTraffic[]`
- `logs[].exitTraffic[].dst`
- `logs[].exitTraffic[].proto`
- `logs[].exitTraffic[].rxBytes`
- `logs[].exitTraffic[].rxPkts`
- `logs[].exitTraffic[].src`
- `logs[].exitTraffic[].txBytes`
- `logs[].exitTraffic[].txPkts`
- `logs[].logged`
- `logs[].nodeId`
- `logs[].physicalTraffic`
- `logs[].physicalTraffic[]`
- `logs[].physicalTraffic[].dst`
- `logs[].physicalTraffic[].proto`
- `logs[].physicalTraffic[].rxBytes`
- `logs[].physicalTraffic[].rxPkts`
- `logs[].physicalTraffic[].src`
- `logs[].physicalTraffic[].txBytes`
- `logs[].physicalTraffic[].txPkts`
- `logs[].start`
- `logs[].subnetTraffic`
- `logs[].subnetTraffic[]`
- `logs[].subnetTraffic[].dst`
- `logs[].subnetTraffic[].proto`
- `logs[].subnetTraffic[].rxBytes`
- `logs[].subnetTraffic[].rxPkts`
- `logs[].subnetTraffic[].src`
- `logs[].subnetTraffic[].txBytes`
- `logs[].subnetTraffic[].txPkts`
- `logs[].virtualTraffic`
- `logs[].virtualTraffic[]`
- `logs[].virtualTraffic[].dst`
- `logs[].virtualTraffic[].proto`
- `logs[].virtualTraffic[].rxBytes`
- `logs[].virtualTraffic[].rxPkts`
- `logs[].virtualTraffic[].src`
- `logs[].virtualTraffic[].txBytes`
- `logs[].virtualTraffic[].txPkts`

### get /tailnet/{tailnet}/logging/{logType}/stream response

- `compressionFormat`
- `destinationType`
- `gcsBucket`
- `gcsCredentials`
- `gcsKeyPrefix`
- `gcsScopes`
- `gcsScopes[]`
- `logType`
- `s3AccessKeyId`
- `s3AuthenticationType`
- `s3Bucket`
- `s3ExternalId`
- `s3KeyPrefix`
- `s3Region`
- `s3RoleArn`
- `s3SecretAccessKey`
- `token`
- `uploadPeriodMinutes`
- `url`
- `user`

### get /tailnet/{tailnet}/logging/{logType}/stream/status response

- `lastActivity`
- `lastError`
- `maxBodySize`
- `numBytesSent`
- `numEntriesSent`
- `numFailedRequests`
- `numSpoofedEntries`
- `numTotalRequests`
- `rateBytesSent`
- `rateEntriesSent`
- `rateFailedRequests`
- `rateTotalRequests`

### get /tailnet/{tailnet}/posture/integrations response

- `integrations`
- `integrations[]`
- `integrations[].clientId`
- `integrations[].clientSecret`
- `integrations[].cloudId`
- `integrations[].configUpdated`
- `integrations[].id`
- `integrations[].provider`
- `integrations[].status`
- `integrations[].status.error`
- `integrations[].status.lastSync`
- `integrations[].status.matchedCount`
- `integrations[].status.possibleMatchedCount`
- `integrations[].status.providerHostCount`
- `integrations[].tenantId`

### get /tailnet/{tailnet}/services response

- `vipServices`
- `vipServices[]`
- `vipServices[].addrs`
- `vipServices[].addrs[]`
- `vipServices[].comment`
- `vipServices[].name`
- `vipServices[].ports`
- `vipServices[].ports[]`
- `vipServices[].tags`
- `vipServices[].tags[]`

### get /tailnet/{tailnet}/services/{serviceName} response

- `addrs`
- `addrs[]`
- `comment`
- `name`
- `ports`
- `ports[]`
- `tags`
- `tags[]`

### get /tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved response

- `approved`
- `autoApproved`

### get /tailnet/{tailnet}/services/{serviceName}/devices response

- `hosts`
- `hosts[]`
- `hosts[].approvalLevel`
- `hosts[].configured`
- `hosts[].stableNodeID`

### get /tailnet/{tailnet}/user-invites response

- `[].email`
- `[].id`
- `[].inviteUrl`
- `[].inviterId`
- `[].lastEmailSentAt`
- `[].role`
- `[].tailnetId`

### get /tailnet/{tailnet}/users response

- `users`
- `users[]`
- `users[].created`
- `users[].currentlyConnected`
- `users[].deviceCount`
- `users[].displayName`
- `users[].id`
- `users[].lastSeen`
- `users[].loginName`
- `users[].profilePicUrl`
- `users[].role`
- `users[].status`
- `users[].tailnetId`
- `users[].type`

### get /tailnet/{tailnet}/webhooks response

- `webhooks`
- `webhooks[]`
- `webhooks[].created`
- `webhooks[].creatorLoginName`
- `webhooks[].endpointId`
- `webhooks[].endpointUrl`
- `webhooks[].lastModified`
- `webhooks[].providerType`
- `webhooks[].secret`
- `webhooks[].subscriptions`
- `webhooks[].subscriptions[]`

### get /user-invites/{userInviteId} response

- `email`
- `id`
- `inviteUrl`
- `inviterId`
- `lastEmailSentAt`
- `role`
- `tailnetId`

### get /users/{userId} response

- `created`
- `currentlyConnected`
- `deviceCount`
- `displayName`
- `id`
- `lastSeen`
- `loginName`
- `profilePicUrl`
- `role`
- `status`
- `tailnetId`
- `type`

### get /webhooks/{endpointId} response

- `created`
- `creatorLoginName`
- `endpointId`
- `endpointUrl`
- `lastModified`
- `providerType`
- `secret`
- `subscriptions`
- `subscriptions[]`

### patch /posture/integrations/{id} request

- `clientId`
- `clientSecret`
- `cloudId`
- `configUpdated`
- `id`
- `provider`
- `status`
- `status.error`
- `status.lastSync`
- `status.matchedCount`
- `status.possibleMatchedCount`
- `status.providerHostCount`
- `tenantId`

### patch /posture/integrations/{id} response

- `clientId`
- `clientSecret`
- `cloudId`
- `configUpdated`
- `id`
- `provider`
- `status`
- `status.error`
- `status.lastSync`
- `status.matchedCount`
- `status.possibleMatchedCount`
- `status.providerHostCount`
- `tenantId`

### patch /tailnet/{tailnet}/contacts/{contactType} request

- `email`

### patch /tailnet/{tailnet}/device-attributes request

- `comment`
- `nodes`

### patch /tailnet/{tailnet}/settings request

- `aclsExternalLink`
- `aclsExternallyManagedOn`
- `httpsEnabled`

### patch /webhooks/{endpointId} request

- `subscriptions`
- `subscriptions[]`

### patch /webhooks/{endpointId} response

- `created`
- `creatorLoginName`
- `endpointId`
- `endpointUrl`
- `lastModified`
- `providerType`
- `secret`
- `subscriptions`
- `subscriptions[]`

### post /device-invites/-/accept request

- `invite`

### post /device-invites/-/accept response

- `acceptedBy`
- `acceptedBy.displayName`
- `acceptedBy.id`
- `acceptedBy.loginName`
- `acceptedBy.profilePicURL`
- `device`
- `device.fqdn`
- `device.id`
- `device.includeExitNode`
- `device.ipv4`
- `device.ipv6`
- `device.name`
- `device.os`
- `sharer`
- `sharer.displayName`
- `sharer.id`
- `sharer.loginName`
- `sharer.profilePicURL`

### post /device/{deviceId}/attributes/{attributeKey} request

- `comment`
- `expiry`
- `value`

### post /device/{deviceId}/attributes/{attributeKey} response

- `attributes`
- `expiries`

### post /device/{deviceId}/authorized request

- `authorized`

### post /device/{deviceId}/device-invites request

- `[].allowExitNode`
- `[].email`
- `[].multiUse`

### post /device/{deviceId}/device-invites response

- `[].accepted`
- `[].acceptedBy`
- `[].acceptedBy.id`
- `[].acceptedBy.loginName`
- `[].acceptedBy.profilePicUrl`
- `[].allowExitNode`
- `[].created`
- `[].deviceId`
- `[].email`
- `[].id`
- `[].inviteUrl`
- `[].lastEmailSentAt`
- `[].multiUse`
- `[].sharerId`
- `[].tailnetId`

### post /device/{deviceId}/ip request

- `ipv4`

### post /device/{deviceId}/key request

- `keyExpiryDisabled`

### post /device/{deviceId}/name request

- `name`

### post /device/{deviceId}/tags request

- `tags`
- `tags[]`

### post /tailnet/{tailnet}/acl/preview response

- `matches`
- `matches[]`
- `matches[].lineNumber`
- `matches[].ports`
- `matches[].ports[]`
- `matches[].users`
- `matches[].users[]`
- `previewFor`
- `type`

### post /tailnet/{tailnet}/acl/validate request

- `[].accept`
- `[].accept[]`
- `[].deny`
- `[].deny[]`
- `[].proto`
- `[].src`
- `[].srcPostureAttrs`

### post /tailnet/{tailnet}/acl/validate response

- `data`
- `data[]`
- `message`

### post /tailnet/{tailnet}/aws-external-id request

- `reusable`

### post /tailnet/{tailnet}/aws-external-id response

- `externalId`
- `tailscaleAwsAccountId`

### post /tailnet/{tailnet}/aws-external-id/{id}/validate-aws-trust-policy request

- `roleArn`

### post /tailnet/{tailnet}/dns/configuration request

- `nameservers`
- `nameservers[]`
- `nameservers[].address`
- `nameservers[].useWithExitNode`
- `preferences`
- `preferences.magicDNS`
- `preferences.overrideLocalDNS`
- `searchPaths`
- `searchPaths[]`
- `splitDNS`

### post /tailnet/{tailnet}/dns/configuration response

- `nameservers`
- `nameservers[]`
- `nameservers[].address`
- `nameservers[].useWithExitNode`
- `preferences`
- `preferences.magicDNS`
- `preferences.overrideLocalDNS`
- `searchPaths`
- `searchPaths[]`
- `splitDNS`

### post /tailnet/{tailnet}/dns/nameservers request

- `dns`
- `dns[]`

### post /tailnet/{tailnet}/dns/nameservers response

- `dns`
- `dns[]`
- `magicDNS`

### post /tailnet/{tailnet}/dns/preferences request

- `magicDNS`

### post /tailnet/{tailnet}/dns/preferences response

- `magicDNS`

### post /tailnet/{tailnet}/dns/searchpaths request

- `searchPaths`
- `searchPaths[]`

### post /tailnet/{tailnet}/dns/searchpaths response

- `searchPaths`
- `searchPaths[]`

### post /tailnet/{tailnet}/keys response

- `audience`
- `capabilities`
- `capabilities.devices`
- `capabilities.devices.create`
- `capabilities.devices.create.ephemeral`
- `capabilities.devices.create.preauthorized`
- `capabilities.devices.create.reusable`
- `capabilities.devices.create.tags`
- `capabilities.devices.create.tags[]`
- `created`
- `customClaimRules`
- `description`
- `expires`
- `expirySeconds`
- `id`
- `invalid`
- `issuer`
- `key`
- `keyType`
- `revoked`
- `scopes`
- `scopes[]`
- `subject`
- `tags`
- `tags[]`
- `updated`
- `userId`

### post /tailnet/{tailnet}/posture/integrations request

- `clientId`
- `clientSecret`
- `cloudId`
- `configUpdated`
- `id`
- `provider`
- `status`
- `status.error`
- `status.lastSync`
- `status.matchedCount`
- `status.possibleMatchedCount`
- `status.providerHostCount`
- `tenantId`

### post /tailnet/{tailnet}/posture/integrations response

- `clientId`
- `clientSecret`
- `cloudId`
- `configUpdated`
- `id`
- `provider`
- `status`
- `status.error`
- `status.lastSync`
- `status.matchedCount`
- `status.possibleMatchedCount`
- `status.providerHostCount`
- `tenantId`

### post /tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved request

- `approved`

### post /tailnet/{tailnet}/services/{serviceName}/device/{deviceId}/approved response

- `approved`
- `autoApproved`

### post /tailnet/{tailnet}/user-invites request

- `[].email`
- `[].role`

### post /tailnet/{tailnet}/user-invites response

- `[].email`
- `[].id`
- `[].inviteUrl`
- `[].inviterId`
- `[].lastEmailSentAt`
- `[].role`
- `[].tailnetId`

### post /tailnet/{tailnet}/webhooks request

- `endpointUrl`
- `providerType`
- `subscriptions`
- `subscriptions[]`

### post /tailnet/{tailnet}/webhooks response

- `created`
- `creatorLoginName`
- `endpointId`
- `endpointUrl`
- `lastModified`
- `providerType`
- `secret`
- `subscriptions`
- `subscriptions[]`

### post /users/{userId}/role request

- `role`

### post /webhooks/{endpointId}/rotate response

- `created`
- `creatorLoginName`
- `endpointId`
- `endpointUrl`
- `lastModified`
- `providerType`
- `secret`
- `subscriptions`
- `subscriptions[]`

### put /tailnet/{tailnet}/keys/{keyId} request

- `audience`
- `customClaimRules`
- `description`
- `issuer`
- `keyType`
- `scopes`
- `scopes[]`
- `subject`
- `tags`
- `tags[]`

### put /tailnet/{tailnet}/keys/{keyId} response

- `audience`
- `capabilities`
- `capabilities.devices`
- `capabilities.devices.create`
- `capabilities.devices.create.ephemeral`
- `capabilities.devices.create.preauthorized`
- `capabilities.devices.create.reusable`
- `capabilities.devices.create.tags`
- `capabilities.devices.create.tags[]`
- `created`
- `customClaimRules`
- `description`
- `expires`
- `expirySeconds`
- `id`
- `invalid`
- `issuer`
- `key`
- `keyType`
- `revoked`
- `scopes`
- `scopes[]`
- `subject`
- `tags`
- `tags[]`
- `updated`
- `userId`

### put /tailnet/{tailnet}/logging/{logType}/stream request

- `compressionFormat`
- `destinationType`
- `gcsBucket`
- `gcsCredentials`
- `gcsKeyPrefix`
- `gcsScopes`
- `gcsScopes[]`
- `logType`
- `s3AccessKeyId`
- `s3AuthenticationType`
- `s3Bucket`
- `s3ExternalId`
- `s3KeyPrefix`
- `s3Region`
- `s3RoleArn`
- `s3SecretAccessKey`
- `token`
- `uploadPeriodMinutes`
- `url`
- `user`

### put /tailnet/{tailnet}/services/{serviceName} request

- `addrs`
- `addrs[]`
- `comment`
- `name`
- `ports`
- `ports[]`
- `tags`
- `tags[]`

### put /tailnet/{tailnet}/services/{serviceName} response

- `addrs`
- `addrs[]`
- `comment`
- `name`
- `ports`
- `ports[]`
- `tags`
- `tags[]`


## Uncovered Properties By Operation

- None
