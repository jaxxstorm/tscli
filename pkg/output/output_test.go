package output

import (
	"reflect"
	"testing"
)

func TestAvailableReturnsSortedFormats(t *testing.T) {
	got := Available()
	want := []string{"human", "json", "pretty", "yaml"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Available() = %v, want %v", got, want)
	}
}
