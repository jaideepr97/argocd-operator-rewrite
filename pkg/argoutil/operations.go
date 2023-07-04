package argoutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomBytes returns a securely generated random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("GenerateRandomBytes: failed to generate bytes: %w", err)
	}
	return bytes, nil
}

// GenerateRandomString returns a securely generated random string.
func GenerateRandomString(s int) (string, error) {
	bytes, err := GenerateRandomBytes(s)
	if err != nil {
		return "", fmt.Errorf("GenerateRandomString: failed to generate string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// IsMergable returns false if any of the extraArgs is already part of the default command Arguments.
func IsMergable(cmd []string, extraArgs []string) bool {
	if len(extraArgs) > 0 {
		for _, arg := range extraArgs {
			if len(arg) > 2 && arg[:2] == "--" {
				if ok := Contains(cmd, arg); ok {
					return false
				}
			}
		}
	}
	return true
}

// Contains returns true if a string is part of the given slice.
func Contains(s []string, g string) bool {
	for _, a := range s {
		if a == g {
			return true
		}
	}
	return false
}
