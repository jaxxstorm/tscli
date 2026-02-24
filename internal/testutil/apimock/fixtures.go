package apimock

func Device() map[string]any {
	return map[string]any{
		"id":        "123",
		"nodeId":    "node-123",
		"hostname":  "device-one",
		"os":        "linux",
		"addresses": []string{"100.64.0.1"},
	}
}

func DeviceList() map[string]any {
	return map[string]any{
		"devices": []map[string]any{Device()},
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
