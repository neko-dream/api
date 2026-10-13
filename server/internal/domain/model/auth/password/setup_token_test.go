package password_auth

import (
	"bytes"
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/neko-dream/api/pkg/hash"
)

func TestSetupTokenHash(t *testing.T) {
	token, err := hash.GenerateRandomToken(32)
	if err != nil {
		t.Fatal(err)
	}
	digest, valid := SetupTokenHash(token)
	expected := sha256.Sum256([]byte(token))
	if !valid || !bytes.Equal(digest, expected[:]) {
		t.Fatal("generated token not accepted")
	}
	for _, invalid := range []string{"", "short", strings.Repeat("A", 44), strings.Repeat("A", 42) + "==", strings.Repeat("A", 42) + "B=", strings.Repeat("A", 42) + "\n="} {
		if _, valid := SetupTokenHash(invalid); valid {
			t.Fatal("malformed or noncanonical token accepted")
		}
	}
}
