package apitype

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTailnetSettingsAdminConsoleSessionTimeout(t *testing.T) {
	for _, minutes := range []int{0, 180, 43200} {
		t.Run(fmt.Sprint(minutes), func(t *testing.T) {
			input := fmt.Sprintf(`{"adminConsoleSessionTimeoutMinutes":%d}`, minutes)
			var settings TailnetSettings
			if err := json.Unmarshal([]byte(input), &settings); err != nil {
				t.Fatal(err)
			}
			got, err := json.Marshal(settings)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != input {
				t.Fatalf("settings lost timeout: got %s, want %s", got, input)
			}
		})
	}
}
