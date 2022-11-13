package util

import (
    "crypto/sha256"
    "encoding/hex"
)

func HashSHA256(rawStr string) string{
    b := []byte(rawStr)
    hashed := sha256.Sum256(b)
    return hex.EncodeToString(hashed[:])
}
