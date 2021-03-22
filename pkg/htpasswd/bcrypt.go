package htpasswd

import (
	"golang.org/x/crypto/bcrypt"
)

func compareBcryptHash(plaintext []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, plaintext)
	return err == nil
}
