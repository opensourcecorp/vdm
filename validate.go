package main

import (
	"errors"
	"fmt"
	"regexp"
)

func (spec depSpec) Validate(r runFlags) error {
	var allErrors []error

	if r.Debug {
		debugLogger.Printf("validating field 'Remote' for %+v", spec)
	}
	if len(spec.Remote) == 0 {
		allErrors = append(allErrors, errors.New("all 'remote' fields must be non-zero length"))
	}
	protocolRegex := regexp.MustCompile(`(https://|git://|git@)`)
	if !protocolRegex.MatchString(spec.Remote) {
		allErrors = append(
			allErrors,
			fmt.Errorf("remote provided as '%s', but all 'remote' fields must begin with a protocol specifier or other valid prefix (e.g. 'https://', 'git@', etc.)", spec.Remote),
		)
	}

	if r.Debug {
		debugLogger.Printf("validating field 'Version' for %+v", spec)
	}
	if len(spec.Version) == 0 {
		allErrors = append(allErrors, errors.New("all 'version' fields must be non-zero length. If you don't care about the version (even though you should), then use 'latest'"))
	}

	if r.Debug {
		debugLogger.Printf("validating field 'LocalPath' for %+v", spec)
	}
	if len(spec.LocalPath) == 0 {
		allErrors = append(allErrors, errors.New("all 'local_path' fields must be non-zero length"))
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			errLogger.Printf("validation failure: %s", err.Error())
		}
		return fmt.Errorf("%d validation failure(s) found in your vdm spec file", len(allErrors))
	}

	return nil
}
