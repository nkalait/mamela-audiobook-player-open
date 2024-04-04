package audio

import "testing"

func PadTestSingle(t *testing.T) {
	padded := pad(0)
	if padded != "00" {
		t.Errorf("pad(0) = %s; want 00", padded)
	}
}
func PadTestDouble(t *testing.T) {
	padded := pad(11)
	if padded != "11" {
		t.Errorf("pad(11) = %s; want 11", padded)
	}
}
