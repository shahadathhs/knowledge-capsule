package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const saltLen = 16

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	h := sha256.Sum256(append(salt, []byte(password)...))
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(h[:]), nil
}

func CheckPassword(password, stored string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		// Legacy format: plain SHA-256 without salt (for backwards compatibility)
		hash := sha256.Sum256([]byte(password))
		return hex.EncodeToString(hash[:]) == stored
	}
	salt, err := hex.DecodeString(parts[0])
	if err != nil || len(salt) != saltLen {
		return false
	}
	h := sha256.Sum256(append(salt, []byte(password)...))
	return hex.EncodeToString(h[:]) == parts[1]
}
