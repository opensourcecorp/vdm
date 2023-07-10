package vdmspec

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

func (spec Spec) Validate() error {
	var allErrors []error

	// Remote field
	logrus.Debugf("validating field 'Remote' for %+v", spec)
	if len(spec.Remote) == 0 {
		allErrors = append(allErrors, errors.New("all 'remote' fields must be non-zero length"))
	}
	protocolRegex := regexp.MustCompile(`(http(s?)://|git://|git@)`)
	if !protocolRegex.MatchString(spec.Remote) {
		allErrors = append(
			allErrors,
			fmt.Errorf("remote provided as '%s', but all 'remote' fields must begin with a protocol specifier or other valid prefix (e.g. 'https://', '(user|git)@', etc.)", spec.Remote),
		)
	}

	// Version field
	logrus.Debugf("validating field 'Version' for %+v", spec)
	if spec.Type == "git" && len(spec.Version) == 0 {
		allErrors = append(allErrors, errors.New("all 'version' fields for the 'git' remote type must be non-zero length. If you don't care about the version (even though you probably should), then use 'latest'"))
	}
	if spec.Type == "file" && len(spec.Version) > 0 {
		logrus.Warnf("NOTE: Remote '%s' specified as type '%s', which does not take explicit version info (you provided '%s'); ignoring version field", spec.Remote, spec.Type, spec.Version)
	}

	// LocalPath field
	logrus.Debugf("validating field 'LocalPath' for %+v", spec)
	if len(spec.LocalPath) == 0 {
		allErrors = append(allErrors, errors.New("all 'local_path' fields must be non-zero length"))
	}

	// Type field
	logrus.Debugf("validating field 'Type' for %+v", spec)
	typeMap := map[string]int{
		"git":  1,
		"":     2, // also git
		"file": 3,
	}
	if _, ok := typeMap[spec.Type]; !ok {
		allErrors = append(allErrors, fmt.Errorf("unrecognized remote type '%s'", spec.Type))
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			logrus.Errorf("validation failure: %s", err.Error())
		}
		return fmt.Errorf("%d validation failure(s) found in your vdm spec file", len(allErrors))
	}
	return nil
}
