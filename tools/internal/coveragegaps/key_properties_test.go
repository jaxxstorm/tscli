package coveragegaps

import "testing"

func TestKeyNodeAttributesResponsePropertyCoverage(t *testing.T) {
	for _, root := range []string{"capabilities", "keys[].capabilities"} {
		covered, err := expandCoverage([]propertyCoverageEvidence{{Type: "apitype.KeyCapabilities", Root: root}})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := covered[root+".devices.create.nodeAttributes"]; !ok {
			t.Fatal("key response coverage must include nodeAttributes")
		}
	}
}

func TestCreateKeyNodeAttributesPropertyCoverage(t *testing.T) {
	covered, err := expandCoverage([]propertyCoverageEvidence{{Type: "coverage.create_key_operation_request"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := covered["capabilities.devices.create.nodeAttributes"]; !ok {
		t.Fatal("auth-key request coverage must include nodeAttributes")
	}
}
