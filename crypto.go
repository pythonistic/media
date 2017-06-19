package media

import (
	"os"
	"golang.org/x/crypto/chacha20poly1305"
	"crypto/rand"
)

func GetRandomBytes(size int) []byte {
	out := make([]byte, size)
	read := 0
	for read < size {
		ct, err := rand.Read(out[read:])
		if err != nil {
			os.Stderr.WriteString("Error getting random bytes: " + err.Error() + "\n")
		}
		read += ct
	}
	return out
}

func Encrypt(plaintext []byte) (nonce []byte, cryptext []byte) {
	cipher, err := chacha20poly1305.New(cookieKey)
	if err != nil {
		os.Stderr.WriteString("Error building cipher: " + err.Error() + "\n")
	}

	nonce = GetRandomBytes(chacha20poly1305.NonceSize)
	dst := []byte{}
	addl := []byte{}
	cryptext = cipher.Seal(dst, nonce, plaintext, addl)
	return
}

func Decrypt(nonce []byte, cryptext []byte) (plaintext []byte) {
	cipher, err := chacha20poly1305.New(cookieKey)
	if err != nil {
		os.Stderr.WriteString("Error building cipher: " + err.Error() + "\n")
	}

	dst := []byte{}
	addl := []byte{}
	plaintext, err = cipher.Open(dst, nonce, cryptext, addl)
	if err != nil {
		os.Stderr.WriteString("Error opening cryptext: " + err.Error() + "\n")
	}
	return
}

