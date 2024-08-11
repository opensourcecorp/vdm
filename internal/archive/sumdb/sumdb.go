package sumdb

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// CalculateSHASum takes an arbitrary [io.Reader] and calculates the SHA256
// checksum for it.
func CalculateSHASum(reader io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("writing reader to hasher: %w", err)
	}
	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	return sum, nil
}
