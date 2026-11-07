package apitype

import (
	"encoding/json"

	tsapi "tailscale.com/client/tailscale/v2"
)

// Key retains SDK metadata while extending capabilities with posture attributes.
type Key struct {
	tsapi.Key
	Capabilities KeyCapabilities `json:"capabilities"`
}

// KeyCapabilities includes device posture attributes not yet modeled by the SDK.
type KeyCapabilities struct {
	Devices struct {
		Create struct {
			Reusable       bool                       `json:"reusable"`
			Ephemeral      bool                       `json:"ephemeral"`
			Tags           []string                   `json:"tags"`
			Preauthorized  bool                       `json:"preauthorized"`
			NodeAttributes map[string]json.RawMessage `json:"nodeAttributes,omitempty"`
		} `json:"create"`
	} `json:"devices"`
}

// CreateKeyRequest is the auth-key request with optional custom posture attributes.
type CreateKeyRequest struct {
	Capabilities  KeyCapabilities `json:"capabilities"`
	ExpirySeconds int64           `json:"expirySeconds,omitempty"`
	Description   string          `json:"description"`
}
