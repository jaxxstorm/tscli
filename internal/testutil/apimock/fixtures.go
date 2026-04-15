package apimock

func Device() map[string]any {
	return map[string]any{
		"id":       "123",
		"nodeId":   "node-123",
		"name":     "device-one",
		"hostname": "device-one",
		"os":       "linux",
		"addresses": []string{
			"100.64.0.1",
		},
		"lastSeen":    "2000-01-01T00:00:00Z",
		"isEphemeral": false,
		"advertisedRoutes": []string{
			"10.0.0.0/24",
		},
		"enabledRoutes": []string{
			"10.0.0.0/24",
		},
		"multipleConnections": true,
		"postureIdentity": map[string]any{
			"serialNumbers": []string{"CP74LFQJXM"},
			"disabled":      false,
		},
	}
}

func DeviceList() map[string]any {
	return map[string]any{
		"devices": []map[string]any{Device()},
	}
}

func DeviceRoutes() map[string]any {
	return map[string]any{
		"advertisedRoutes": []string{"10.0.0.0/24"},
		"enabledRoutes":    []string{"10.0.0.0/24"},
	}
}

func TailnetSettings() map[string]any {
	return map[string]any{
		"devicesApprovalOn":           true,
		"devicesAutoUpdatesOn":        true,
		"devicesKeyDurationDays":      30,
		"usersApprovalOn":             false,
		"networkFlowLoggingOn":        true,
		"regionalRoutingOn":           false,
		"postureIdentityCollectionOn": true,
		"httpsEnabled":                true,
	}
}

func KeyResponse() map[string]any {
	return map[string]any{
		"id":          "k123",
		"description": "test",
		"key":         "tskey-auth-abc123",
	}
}

func KeyList() []map[string]any {
	return []map[string]any{KeyResponse()}
}

func KeyListEnvelope() map[string]any {
	return map[string]any{
		"keys": KeyList(),
	}
}

func Contacts() map[string]any {
	return map[string]any{
		"primary": map[string]any{
			"email": "ops@example.com",
		},
		"billing": map[string]any{
			"email": "billing@example.com",
		},
	}
}

func DNSConfiguration() map[string]any {
	return map[string]any{
		"magicDNS":    true,
		"nameservers": []string{"1.1.1.1"},
		"searchPaths": []string{"corp.example.com"},
		"splitDNS": map[string]any{
			"corp.example.com": []string{"1.1.1.1"},
		},
	}
}

func DNSNameservers() []string {
	return []string{"1.1.1.1", "8.8.8.8"}
}

func DNSSearchPaths() []string {
	return []string{"corp.example.com", "svc.example.com"}
}

func DNSSplitConfig() map[string]any {
	return map[string]any{
		"corp.example.com": []string{"1.1.1.1"},
		"svc.example.com":  []string{"8.8.8.8"},
	}
}

func Invite() map[string]any {
	return map[string]any{
		"id":     "invite-1",
		"email":  "user@example.com",
		"status": "pending",
	}
}

func InviteList() []map[string]any {
	return []map[string]any{Invite()}
}

func DevicePostureResponse() map[string]any {
	return map[string]any{
		"attributes": map[string]any{
			"custom:group": "prod",
		},
		"expiries": map[string]any{},
	}
}

func LogsConfiguration() []map[string]any {
	return []map[string]any{
		{
			"id":     "cfg-1",
			"action": "policy.updated",
		},
	}
}

func LogsNetwork() []map[string]any {
	return []map[string]any{
		{
			"id":    "net-1",
			"srcIP": "100.64.0.1",
			"dstIP": "100.64.0.2",
		},
	}
}

func LogsStream() map[string]any {
	return map[string]any{
		"enabled":  true,
		"endpoint": "https://example.com/logs",
	}
}

func AWSExternalID() map[string]any {
	return map[string]any{
		"externalId": "ext-123",
	}
}

func AWSValidation() map[string]any {
	return map[string]any{
		"valid": true,
	}
}

func Policy() string {
	return "{\n  \"acls\": []\n}\n"
}

func PolicyPreview() map[string]any {
	return map[string]any{
		"matches": []map[string]any{
			{
				"action": "accept",
			},
		},
	}
}

func PolicyValidation() map[string]any {
	return map[string]any{
		"valid": true,
	}
}

func PostureAttributeMap() map[string]any {
	return map[string]any{
		"custom:group": "prod",
	}
}

func PostureIntegration() map[string]any {
	return map[string]any{
		"id":       "pi-1",
		"provider": "falcon",
	}
}

func PostureIntegrationList() map[string]any {
	return map[string]any{
		"integrations": []map[string]any{PostureIntegration()},
	}
}

func Service() map[string]any {
	return map[string]any{
		"name":    "svc:demo-speedtest",
		"addrs":   []string{"100.82.154.103", "fd7a:115c:a1e0::b101:9bb1"},
		"comment": "This Tailscale Service is managed by the Tailscale Kubernetes Operator",
		"ports":   []string{"tcp:443"},
		"tags":    []string{"tag:demo-speedtest"},
		"annotations": map[string]any{
			"tailscale.com/owner-references": `{"ownerRefs":[{"operatorID":"nbFLzCnzKQ11CNTRL"}]}`,
		},
	}
}

func ServiceList() map[string]any {
	return map[string]any{
		"vipServices": []map[string]any{
			Service(),
			{
				"name":    "svc:demo-streamer",
				"addrs":   []string{"100.106.40.81", "fd7a:115c:a1e0::da01:2897"},
				"comment": "This Tailscale Service is managed by the Tailscale Kubernetes Operator",
				"ports":   []string{"tcp:443"},
				"tags":    []string{"tag:demo-streamer"},
				"annotations": map[string]any{
					"tailscale.com/owner-references": `{"ownerRefs":[{"operatorID":"nP1yvuuBKr11CNTRL"}]}`,
				},
			},
		},
	}
}

func ServiceDevices() []map[string]any {
	return []map[string]any{
		{
			"nodeId":   "node-123",
			"approved": true,
		},
	}
}

func ServiceApproval() map[string]any {
	return map[string]any{
		"approved": true,
	}
}

func User() map[string]any {
	return map[string]any{
		"id":                 "user-1",
		"created":            "2025-02-24T19:05:37.458138867Z",
		"currentlyConnected": false,
		"deviceCount":        1,
		"displayName":        "User Example",
		"lastSeen":           "2026-04-14T21:28:51Z",
		"loginName":          "user@example.com",
		"profilePicUrl":      "",
		"role":               "member",
		"status":             "active",
		"tailnetId":          "8193236004369213",
		"type":               "member",
	}
}

func UserList() []map[string]any {
	return []map[string]any{
		User(),
		{
			"id":                 "user-2",
			"created":            "2025-02-24T19:05:37.458138867Z",
			"currentlyConnected": false,
			"deviceCount":        0,
			"displayName":        "Suspended User",
			"lastSeen":           "2026-04-10T21:28:51Z",
			"loginName":          "suspended@example.com",
			"profilePicUrl":      "",
			"role":               "member",
			"status":             "suspended",
			"tailnetId":          "8193236004369213",
			"type":               "member",
		},
	}
}

func UserListEnvelope() map[string]any {
	return map[string]any{
		"users": UserList(),
	}
}

func Webhook() map[string]any {
	return map[string]any{
		"id":          "wh-1",
		"endpointUrl": "https://example.com/hook",
	}
}

func WebhookList() []map[string]any {
	return []map[string]any{Webhook()}
}

func WebhookListEnvelope() map[string]any {
	return map[string]any{
		"webhooks": WebhookList(),
	}
}

func Error(message string) map[string]any {
	return map[string]any{
		"message": message,
	}
}
