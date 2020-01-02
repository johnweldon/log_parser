package main

import "testing"

func TestGenLine(t *testing.T) {
	tcase := genLine(NEWLINE, PUNCT, NUMBER, PUNCT, NUMBER, PUNCT, NUMBER, PUNCT)
	if len(tcase) != 8 {
		t.Fatalf("genLine did not return expected slice len; got (%d): %+v", len(tcase), tcase.DebugString())
	}
}
