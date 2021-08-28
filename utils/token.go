package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func CreateId(title string) (string, error) {
	var strBuild strings.Builder

	strBuild.WriteString(title)
	strBuild.WriteString(time.Now().String())

	return createHash([]byte(strBuild.String())), nil
}

func createHash(seed []byte) string {
	result := sha256.Sum256(seed)
	return hex.EncodeToString(result[:])
}
