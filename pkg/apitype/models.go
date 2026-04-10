package apitype

import tsapi "tailscale.com/client/tailscale/v2"

type DERPRegion struct {
	Preferred           *bool   `json:"preferred,omitempty"`
	LatencyMilliseconds float64 `json:"latencyMs,omitempty"`
}

type ClientSupports struct {
	HairPinning *bool `json:"hairPinning"`
	IPV6        *bool `json:"ipv6"`
	PCP         *bool `json:"pcp"`
	PMP         *bool `json:"pmp"`
	UDP         *bool `json:"udp"`
	UPNP        *bool `json:"upnp"`
}

type ClientConnectivity struct {
	Endpoints             []string              `json:"endpoints,omitempty"`
	MappingVariesByDestIP *bool                 `json:"mappingVariesByDestIP,omitempty"`
	Latency               map[string]DERPRegion `json:"latency,omitempty"`
	ClientSupports        *ClientSupports       `json:"clientSupports,omitempty"`
}

type Distro struct {
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
	CodeName string `json:"codeName,omitempty"`
}

type DevicePostureIdentity struct {
	SerialNumbers []string `json:"serialNumbers,omitempty"`
	Disabled      *bool    `json:"disabled,omitempty"`
}

type Device struct {
	Addresses                 []string               `json:"addresses,omitempty"`
	Name                      string                 `json:"name,omitempty"`
	ID                        string                 `json:"id,omitempty"`
	NodeID                    string                 `json:"nodeId,omitempty"`
	Authorized                *bool                  `json:"authorized,omitempty"`
	User                      string                 `json:"user,omitempty"`
	Tags                      []string               `json:"tags,omitempty"`
	KeyExpiryDisabled         *bool                  `json:"keyExpiryDisabled,omitempty"`
	BlocksIncomingConnections *bool                  `json:"blocksIncomingConnections,omitempty"`
	ClientVersion             string                 `json:"clientVersion,omitempty"`
	Created                   *tsapi.Time            `json:"created,omitempty"`
	Expires                   *tsapi.Time            `json:"expires,omitempty"`
	Hostname                  string                 `json:"hostname,omitempty"`
	IsEphemeral               *bool                  `json:"isEphemeral,omitempty"`
	IsExternal                *bool                  `json:"isExternal,omitempty"`
	ConnectedToControl        *bool                  `json:"connectedToControl,omitempty"`
	LastSeen                  *tsapi.Time            `json:"lastSeen,omitempty"`
	MachineKey                string                 `json:"machineKey,omitempty"`
	NodeKey                   string                 `json:"nodeKey,omitempty"`
	OS                        string                 `json:"os,omitempty"`
	TailnetLockError          string                 `json:"tailnetLockError,omitempty"`
	TailnetLockKey            string                 `json:"tailnetLockKey,omitempty"`
	UpdateAvailable           *bool                  `json:"updateAvailable,omitempty"`
	MultipleConnections       *bool                  `json:"multipleConnections,omitempty"`
	EnabledRoutes             []string               `json:"enabledRoutes,omitempty"`
	AdvertisedRoutes          []string               `json:"advertisedRoutes,omitempty"`
	ClientConnectivity        *ClientConnectivity    `json:"clientConnectivity,omitempty"`
	SSHEnabled                *bool                  `json:"sshEnabled,omitempty"`
	PostureIdentity           *DevicePostureIdentity `json:"postureIdentity,omitempty"`
	Distro                    *Distro                `json:"distro,omitempty"`
}

type DeviceListResponse struct {
	Devices []Device `json:"devices,omitempty"`
}

type DeviceRoutes struct {
	AdvertisedRoutes []string `json:"advertisedRoutes,omitempty"`
	EnabledRoutes    []string `json:"enabledRoutes,omitempty"`
}

type DeviceRoutesUpdateRequest struct {
	Routes []string `json:"routes,omitempty"`
}

type TailnetSettings struct {
	ACLsExternallyManagedOn                *bool   `json:"aclsExternallyManagedOn,omitempty"`
	ACLsExternalLink                       *string `json:"aclsExternalLink,omitempty"`
	DevicesApprovalOn                      *bool   `json:"devicesApprovalOn,omitempty"`
	DevicesAutoUpdatesOn                   *bool   `json:"devicesAutoUpdatesOn,omitempty"`
	DevicesKeyDurationDays                 *int    `json:"devicesKeyDurationDays,omitempty"`
	UsersApprovalOn                        *bool   `json:"usersApprovalOn,omitempty"`
	UsersRoleAllowedToJoinExternalTailnets *string `json:"usersRoleAllowedToJoinExternalTailnets,omitempty"`
	NetworkFlowLoggingOn                   *bool   `json:"networkFlowLoggingOn,omitempty"`
	RegionalRoutingOn                      *bool   `json:"regionalRoutingOn,omitempty"`
	PostureIdentityCollectionOn            *bool   `json:"postureIdentityCollectionOn,omitempty"`
	HTTPSEnabled                           *bool   `json:"httpsEnabled,omitempty"`
}

type UpdateTailnetSettingsRequest struct {
	ACLsExternallyManagedOn                *bool   `json:"aclsExternallyManagedOn,omitempty"`
	ACLsExternalLink                       *string `json:"aclsExternalLink,omitempty"`
	DevicesApprovalOn                      *bool   `json:"devicesApprovalOn,omitempty"`
	DevicesAutoUpdatesOn                   *bool   `json:"devicesAutoUpdatesOn,omitempty"`
	DevicesKeyDurationDays                 *int    `json:"devicesKeyDurationDays,omitempty"`
	UsersApprovalOn                        *bool   `json:"usersApprovalOn,omitempty"`
	UsersRoleAllowedToJoinExternalTailnets *string `json:"usersRoleAllowedToJoinExternalTailnets,omitempty"`
	NetworkFlowLoggingOn                   *bool   `json:"networkFlowLoggingOn,omitempty"`
	RegionalRoutingOn                      *bool   `json:"regionalRoutingOn,omitempty"`
	PostureIdentityCollectionOn            *bool   `json:"postureIdentityCollectionOn,omitempty"`
	HTTPSEnabled                           *bool   `json:"httpsEnabled,omitempty"`
}
