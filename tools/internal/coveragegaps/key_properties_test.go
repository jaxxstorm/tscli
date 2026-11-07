package coveragegaps

import "testing"

func TestCreateKeyNodeAttributesPropertyCoverage(t *testing.T) {
	covered, err := expandCoverage([]propertyCoverageEvidence{{Type: "coverage.create_key_operation_request"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := covered["capabilities.devices.create.nodeAttributes"]; !ok {
		t.Fatal("auth-key request coverage must include nodeAttributes")
	}
}
