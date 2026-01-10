package model

import "testing"

func TestSignatureStable(t *testing.T) {
	r := Report{Port: 1, Proto: "tcp"}
	s1 := r.Signature()
	s2 := r.Signature()
	if s1 != s2 {
		t.Fatalf("signature should be stable; got %q vs %q", s1, s2)
	}
}
