package vdmspec

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/opensourcecorp/vdm/internal/message"
)

// Validate performs runtime validations on the vdm specfile, and informs the
// caller of any failures encountered.
func (spec Spec) Validate() error {
	var allErrors []error

	for remoteIndex, remote := range spec.Remotes {
		// Remote field
		message.Debugf("Index %d: validating field 'Remote' for %+v", remoteIndex, remote)
		if len(remote.Remote) == 0 {
			allErrors = append(allErrors, errors.New("all 'remote' fields must be non-zero length"))
		}
		protocolRegex := regexp.MustCompile(`(http(s?)://|git://|git@)`)
		if !protocolRegex.MatchString(remote.Remote) {
			allErrors = append(
				allErrors,
				fmt.Errorf("remote #%d provided as '%s', but all 'remote' fields must begin with a protocol specifier or other valid prefix (e.g. 'https://', '(user|git)@', etc.)", remoteIndex, remote.Remote),
			)
		}

		// Version field
		message.Debugf("Index %d: validating field 'Version' for %+v", remoteIndex, remote)
		if remote.Type == "git" && len(remote.Version) == 0 {
			allErrors = append(allErrors, errors.New("all 'version' fields for the 'git' remote type must be non-zero length. If you don't care about the version (even though you probably should), then use 'latest'"))
		}
		if remote.Type == "file" && len(remote.Version) > 0 {
			message.Warnf("NOTE: Remote #%d '%s' specified as type '%s', which does not take explicit version info (you provided '%s'); ignoring version field", remoteIndex, remote.Remote, remote.Type, remote.Version)
		}

		// LocalPath field
		message.Debugf("Index #%d: validating field 'LocalPath' for %+v", remoteIndex, remote)
		if len(remote.LocalPath) == 0 {
			allErrors = append(allErrors, errors.New("all 'local_path' fields must be non-zero length"))
		}

		// Type field
		message.Debugf("Index #%d: validating field 'Type' for %+v", remoteIndex, remote)
		typeMap := map[string]int{
			"git":  1,
			"":     2, // also git
			"file": 3,
		}
		if _, ok := typeMap[remote.Type]; !ok {
			allErrors = append(allErrors, fmt.Errorf("unrecognized remote type '%s'", remote.Type))
		}
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			message.Errorf("validation failure: %s", err.Error())
		}
		return fmt.Errorf("%d validation failure(s) found in your vdm spec file", len(allErrors))
	}
	return nil
}
