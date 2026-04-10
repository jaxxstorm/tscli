package apimock

func Device() map[string]any {
	return map[string]any{
		"id":       "123",
		"nodeId":   "node-123",
		"hostname": "device-one",
		"os":       "linux",
		"addresses": []string{
			"100.64.0.1",
		},
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

func Error(message string) map[string]any {
	return map[string]any{
		"message": message,
	}
}
