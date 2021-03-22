package htpasswd

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
)

func compareShaHash(plaintext []byte, hash []byte) bool {
	d := sha1.New()
	d.Write(plaintext)
	if subtle.ConstantTimeCompare(hash[5:], []byte(base64.StdEncoding.EncodeToString(d.Sum(nil)))) != 1 {
		return false
	}
	return true
}
