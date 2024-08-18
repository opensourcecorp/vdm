package cache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	_ "modernc.org/sqlite"
)

var (
	createStatement = "CREATE"
	insertStatement = "INSERT"
	dbStatements    = map[string]string{
		createStatement: `
			CREATE TABLE IF NOT EXISTS sums (
				key TEXT PRIMARY KEY,
				source TEXT,
				version TEXT,
				sum TEXT UNIQUE
			);
		`,
		insertStatement: `
			INSERT INTO sums (
				key, source, version, sum
			) VALUES (
			 	?, ?, ?, ?
			);
		`,
	}
)

func AddToSumDB(remote vdmspec.Remoter, reader io.Reader) (err error) {
	sumDBPath, err := getSumDBPath()
	if err != nil {
		return fmt.Errorf("getting sumdb path %q: %w", sumDBPath, err)
	}

	db, err := sql.Open("sqlite", sumDBPath)
	if err != nil {
		return fmt.Errorf("opening sumdb path %q: %w", sumDBPath, err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing sumdb file %q: %w", sumDBPath, closeErr))
		}
	}()

	_, err = db.Exec(dbStatements[createStatement])
	if err != nil {
		return fmt.Errorf("creating sums table: %w", err)
	}

	sum, err := CalculateSHASum(reader)
	if err != nil {
		return fmt.Errorf("calculating SHA sum for remote source %q: %w", remote.GetSource(), err)
	}
	message.Debugf("sum calculated for remote %q's archive file was %q", remote.GetSource(), sum)

	_, err = db.Exec(
		dbStatements[insertStatement],
		remote.GetSourceVersionSum(sum),
		remote.GetSource(),
		remote.GetVersion(),
		sum,
	)
	if err != nil {
		return fmt.Errorf("inserting into sums table: %w", err)
	}

	return err
}

// CalculateSHASum takes an arbitrary [io.Reader] (such as an open
// file handle) and calculates the SHA256 checksum for it.
func CalculateSHASum(reader io.Reader) (string, error) {
	message.Debugf("reader address for calculating SHA sum: %v", reader)
	hasher := sha256.New()
	var n int64
	var err error
	if n, err = io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("writing reader to hasher: %w", err)
	}
	message.Debugf("number of bytes copied to hasher: %d", n)
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

func getSumDBPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determining user homedir: %w", err)
	}
	sumDBPath := filepath.Join(homedir, ".vdm", "cache", "sum.db")
	return sumDBPath, nil
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
