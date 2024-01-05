package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}
