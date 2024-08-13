package cache

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// CalculateSHASum takes an arbitrary [io.Reader] (such as an open file handle)
// and calculates the SHA256 checksum for it.
func CalculateSHASum(reader io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("writing reader to hasher: %w", err)
	}
	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	return sum, nil
}

// StringToBase64 does what it says on the tin. However, since this is intended
// for use as an archive file name, it replaces padding characters ('=') with
// '_PAD'.
func StringToBase64(s string) string {
	// We use base64-encoded names for cache archives, because remote sources
	// have slashes and that can break FS pathing
	out := base64.StdEncoding.EncodeToString([]byte(s))
	out = strings.ReplaceAll(out, "=", "_PAD")
	return out
}

// StringFromBase64 does the inverse of [StringToBase64].
func StringFromBase64(s string) (string, error) {
	fixedString := strings.ReplaceAll(s, "_PAD", "=")
	out, err := base64.StdEncoding.DecodeString(fixedString)
	if err != nil {
		return "", fmt.Errorf("decoding base64 string '%s': %w", s, err)
	}
	return string(out), nil
}
