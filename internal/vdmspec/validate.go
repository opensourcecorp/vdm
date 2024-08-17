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

	protocolRegex := regexp.MustCompile(`(http(s?)://|git://|git@|ftp(s?))`)
	for remoteIndex, remote := range spec.Remotes {
		// Source field
		message.Debugf("Index #%d: validating field 'Source' for %+v", remoteIndex, remote)
		if len(remote.Source) == 0 {
			allErrors = append(allErrors, errors.New("all 'source' fields must be non-zero length"))
		}
		if !protocolRegex.MatchString(remote.Source) {
			allErrors = append(
				allErrors,
				fmt.Errorf("remote #%d provided as %q, but all 'source' fields must begin with a protocol specifier or other valid prefix (e.g. 'https://', '(user|git)@', etc.)", remoteIndex, remote.Source),
			)
		}

		// Version field
		message.Debugf("Index #%d: validating field 'Version' for %+v", remoteIndex, remote)
		if remote.Type == GitType && remote.Version == "" {
			allErrors = append(allErrors, errors.New("all 'version' fields for the 'git' remote type must be non-zero length"))
		}
		if remote.Type == FileType && remote.Version != "" {
			message.Warnf("NOTE: Remote #%d %q specified as type %q, which does not take explicit version info (you provided %q); ignoring version field", remoteIndex, remote.Source, remote.Type, remote.Version)
		}

		// Destination field
		message.Debugf("Index #%d: validating field 'Destination' for %+v", remoteIndex, remote)
		if len(remote.Destination) == 0 {
			allErrors = append(allErrors, errors.New("all 'destination' fields must be non-zero length"))
		}

		// Type field
		message.Debugf("Index #%d: validating field 'Type' for %+v", remoteIndex, remote)
		if remote.Type == "" {
			allErrors = append(allErrors, errors.New("all remotes must specify a 'type' field"))
		}
		typeMap := map[string]int{
			GitType:  1,
			FileType: 2,
		}
		if _, ok := typeMap[remote.Type]; !ok {
			allErrors = append(allErrors, fmt.Errorf("unrecognized remote type %q", remote.Type))
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
