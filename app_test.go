package main

import (
	"mamela/audio/tests"
	"testing"
)

func TestPad(t *testing.T) {
	tests.PadTestSingle(t)
	tests.PadTestDouble(t)
}
