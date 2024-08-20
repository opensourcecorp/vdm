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

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	_ "modernc.org/sqlite"
)

var (
	tableName               = "sums"
	createStatement         = "CREATE"
	insertStatement         = "INSERT"
	checkKeyExistsStatement = "CHECK_EXISTS"
	dbStatements            = map[string]string{
		createStatement: fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				key TEXT UNIQUE PRIMARY KEY,
				source TEXT,
				version TEXT,
				sum TEXT UNIQUE
			);
		`, tableName),
		insertStatement: fmt.Sprintf(`
			INSERT INTO %s (
				key, source, version, sum
			) VALUES (
			 	?, ?, ?, ?
			);
		`, tableName),
		checkKeyExistsStatement: fmt.Sprintf(`
			SELECT COUNT(*)
			FROM %s
			WHERE key = ?
			;
		`, tableName),
	}
)

func CreateSumDB() (err error) {
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

	return err
}

func AddToSumDB(remote vdmspec.Remoter, cacheTargetPath string) (err error) {
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

	f, err := os.Open(cacheTargetPath)
	if err != nil {
		return fmt.Errorf("opening cache target path %q for hashing: %w", cacheTargetPath, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing cache target file %q: %w", cacheTargetPath, closeErr))
		}
	}()

	sum, err := calculateSHASum(f)
	if err != nil {
		return fmt.Errorf("calculating SHA sum for remote source %q: %w", remote.GetSource(), err)
	}
	message.Debugf("sum calculated for remote %q's archive file was %q", remote.GetSource(), sum)

	_, err = db.Exec(
		dbStatements[insertStatement],
		remote.GetSumDBKey(),
		remote.GetSource(),
		remote.GetVersion(),
		sum,
	)
	if err != nil {
		return fmt.Errorf("inserting into sums table: %w", err)
	}

	return err
}

func CheckIfRemoteInSumDB(remote vdmspec.Remoter) (hasKey bool, err error) {
	sumDBPath, err := getSumDBPath()
	if err != nil {
		return false, fmt.Errorf("getting sumdb path %q: %w", sumDBPath, err)
	}

	db, err := sql.Open("sqlite", sumDBPath)
	if err != nil {
		return false, fmt.Errorf("opening sumdb path %q: %w", sumDBPath, err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing sumdb file %q: %w", sumDBPath, closeErr))
		}
	}()

	var numRows int
	err = db.QueryRow(dbStatements[checkKeyExistsStatement], remote.GetSumDBKey()).Scan(&numRows)
	if err != nil {
		return false, fmt.Errorf("querying sumdb: %w", err)
	}
	message.Debugf("number of results from sumdb for remote key %q: %d", remote.GetSumDBKey(), numRows)

	message.Debugf("sumdb query result not yet checked for remote key %q, hasKey: %v", remote.GetSumDBKey(), hasKey)
	if numRows > 0 {
		hasKey = true
	}
	message.Debugf("sumdb query result now checked for remote key %q, hasKey: %v", remote.GetSumDBKey(), hasKey)

	return hasKey, err
}

// calculateSHASum takes an arbitrary [io.Reader] (such as an open
// file handle) and calculates the SHA256 checksum for it.
func calculateSHASum(reader io.Reader) (string, error) {
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
	vdmCachePath, err := vars.GetVDMCacheDir()
	if err != nil {
		return "", fmt.Errorf("determining vdm cache path during sumdb creation: %w", err)
	}

	message.Debugf("vdm cache directory during sumdb path determination was %q", vdmCachePath)
	sumDBPath := filepath.Join(vdmCachePath, "sum.db")
	message.Debugf("sumdb path to be used: %q", sumDBPath)

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
