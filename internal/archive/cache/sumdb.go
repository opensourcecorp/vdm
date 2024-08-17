package cache

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Caller's job to close file handle
func GetOrCreateSumDBFile() (*os.File, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("determining user homedir: %w", err)
	}

	sumDBPath := filepath.Join(homedir, ".vdm", "cache", "sumdb.json")
	sumDBFile, err := os.OpenFile(sumDBPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("creating/opening sumdb file %q: %w", sumDBPath, err)
	}

	return sumDBFile, err
}

// CalculateSHASum takes an arbitrary [io.Reader] (such as an open
// file handle) and calculates the SHA256 checksum for it.
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
	out = replaceNonAlphaBase64Characters(out)
	return out
}

// StringFromBase64 does the inverse of [StringToBase64].
func StringFromBase64(s string) (string, error) {
	fixedString := restoreNonAlphaBase64Characters(s)
	out, err := base64.StdEncoding.DecodeString(fixedString)
	if err != nil {
		return "", fmt.Errorf("decoding base64 string %q: %w", s, err)
	}
	return string(out), nil
}

// base64NonAlphaMap maps non-alphanumeric base-64 characters to their arbitrary
// replacement values
var base64NonAlphaMap = map[string]string{
	"=": "_EQ",
	"/": "_SLASH",
	"+": "_PLUS",
}

// replaceNonAlphaBase64Characters takes a valid base64 string and replaces
// non-alphanumeric characters
func replaceNonAlphaBase64Characters(s string) string {
	out := s
	for nonAlpha, alpha := range base64NonAlphaMap {
		out = strings.ReplaceAll(out, nonAlpha, alpha)
	}
	return out
}

// restoreNonAlphaBase64Characters takes an invalid base64 string and restores
// non-alphanumeric characters
func restoreNonAlphaBase64Characters(s string) string {
	out := s
	for nonAlpha, alpha := range base64NonAlphaMap {
		out = strings.ReplaceAll(out, alpha, nonAlpha)
	}
	return out
}
