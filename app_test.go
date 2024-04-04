package main

import (
	"mamela/audio"
	"testing"
)

func TestPad(t *testing.T) {
	audio.PadTestSingle(t)
	audio.PadTestDouble(t)
}
