package media

import (
	"testing"
)

func TestGetRandomBytes(t *testing.T) {
	nonce := GetRandomBytes(64)

	if len(nonce) != 64 {
		t.Errorf("Nonce too short, should be 64 bytes, got %d", len(nonce))
	}

	for _, n := range nonce {
		if n != 0 {
			goto pass;
		}
	}
	t.Errorf("Nonce was never initialized, expected random values, got all 0: %x", nonce)

	pass:
}

func TestEncrypt(t *testing.T) {
	plaintext := "Alias molestiae ipsum consequatur aliquam. Debitis quod ut officia. Soluta sapiente voluptate tempora et in. Id incidunt odit ab. Corrupti quaerat harum ullam laboriosam."

	nonce, cryptext := Encrypt([]byte(plaintext))

	if len(nonce) != 12 {
		t.Errorf("Nonce too short, should be 12 bytes, got %d", len(nonce))
	}

	if len(cryptext) == 0 {
		t.Error("Empty cryptext -- failed")
	}

	if string(cryptext) == plaintext {
		t.Error("Cryptext same as plaintext; encryption failed")
	}

	decrypted := Decrypt(nonce, cryptext)

	if string(decrypted) != plaintext {
		t.Errorf("Expected decrypted text to match plaintext, got '%s'", decrypted)
	}
}
