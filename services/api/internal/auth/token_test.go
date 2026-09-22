package auth

import "testing"

func TestHashTokenIsStableAndNotRaw(t *testing.T) {
	raw := "example-token"
	hash := HashToken(raw)
	if hash == raw {
		t.Fatal("hash must not equal the raw token")
	}
	if hash != HashToken(raw) {
		t.Fatal("hash must be deterministic")
	}
	if HashToken("other") == hash {
		t.Fatal("different tokens must not hash the same")
	}
}
