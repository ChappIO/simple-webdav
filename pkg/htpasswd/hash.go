package htpasswd

import "bytes"

func getHashMethod(hash []byte) HashMethod {
	if bytes.HasPrefix(hash, []byte("$apr1")) {
		return HashMD5
	}
	if bytes.HasPrefix(hash, []byte("{SHA}")) {
		return HashSha
	}
	if bytes.HasPrefix(hash, []byte("$2")) {
		return HashBcrypt
	}
	panic("unknown hash algo")
}

func compareHash(plaintext []byte, hash []byte, method HashMethod) bool {
	switch method {
	case HashMD5:
		return compareMd5Hash(plaintext, hash)
	case HashSha:
		return compareShaHash(plaintext, hash)
	case HashBcrypt:
		return compareBcryptHash(plaintext, hash)
	default:
		return false
	}
}
