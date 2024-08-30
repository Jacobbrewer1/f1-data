package utils

import (
	"crypto/sha1"
	"fmt"
)

func GenerateShaToken(value string) string {
	// Store the SHA1-hash of the content as the ID. This is used to detect duplicate submissions.
	crypt := sha1.New()
	crypt.Write([]byte(value))
	return fmt.Sprintf("%x", crypt.Sum(nil))
}
