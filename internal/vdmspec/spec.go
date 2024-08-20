package vdmspec

import (
	"fmt"
	"os"
	"strings"

	"github.com/opensourcecorp/vdm/internal/message"
	"gopkg.in/yaml.v3"
)

const (
	// GitType represents the string to match against for "git" remote types.
	GitType string = "git"
	// ArchiveType represents the string to match against for "archive" remote
	// types.
	ArchiveType string = "archive"
	// LocalType represents the string to match against for "local" "remote"
	// types.
	LocalType string = "local"
)

// typeMap is used to house more programmatic access to the various remote
// types, such as in tests or [Spec.Validate]
var typeMap = map[string]int{
	GitType:     0,
	ArchiveType: 1,
	LocalType:   2,
}

// Spec defines the overall structure of the vmd specfile.
type Spec struct {
	Remotes []RemoteTemplate `json:"remotes" yaml:"remotes"`
}

// Remoter defines behavior that different remote types must exhibit
type Remoter interface {
	// Cache should retrieve the remote, and cache it as an archive in
	// VDM_HOME's cache
	Cache() (string, error)
	// Sync should unpack the archive from the cache in VDM_HOME to the
	// specified destination
	Sync(src, dest string) error
	// GetRemote should return the [RemoteTemplate.Source] value
	GetSource() string
	// GetRemote should return the [RemoteTemplate.Version] value
	GetVersion() string
	// GetSumDBKey should concatenate the relevant values from the
	// [RemoteTemplate] to produce a string value satisfying the primary key
	// constraint for the sumdb
	GetSumDBKey() string
}

// RemoteTemplate defines the template structure for each potential remote
// configuration in the vdm specfile. This struct can be embedded into other
// remote type structs to provide the common fields between them.
type RemoteTemplate struct {
	// Type is the type of Source, e.g. git, archive, file, etc.
	Type string `json:"type" yaml:"type"`
	// Source is the fully-qualifed location from which the Remote is retrieved,
	// e.g. "https://github.com/some-org/some-repo"
	Source string `json:"source" yaml:"source"`
	// Version states the version requested from Source, and is then later used
	// for tracking purposes. Version can be anything supported by the Type
	// field -- for example, for the "git" Type, this can be a tag, a branch
	// name, or a commit hash.
	Version string `json:"version" yaml:"version"`
	// Destination is the relative or absolute path on disk that Source will be
	// placed at
	Destination string `json:"destination" yaml:"destination"`
	// TryLocalSource helps define behavior driven by the `try-local-sources`
	// CLI flag, which allows checking for a local version of a
	// [RemoteTemplate.Source], and falling back to the other Source field if
	// the local path does not exist. This is especially useful for when you
	// might be developing one of your Remotes in a nearby directory, and want
	// to copy over that version of the Remote and not keep pushing-and-pulling
	// to a Git upstream just to test the changes.
	TryLocalSource string `json:"try_local_source" yaml:"try_local_source"`
}

// GetSpecFromFile reads the specfile from disk (the path of which may be
// determined by the user-supplied flag value), and returns it for further
// processing of remotes.
func GetSpecFromFile(specFilePath string) (Spec, error) {
	specFile, err := os.ReadFile(specFilePath)
	if err != nil {
		message.Debugf("error reading specfile from disk: %v", err)
		return Spec{}, fmt.Errorf(
			strings.Join([]string{
				"there was a problem reading your vdm file from %q -- does it not exist?",
				"Either pass the --spec-file flag, or create one in the default location (details in the README).",
				"Error details: %w"},
				" ",
			),
			specFilePath,
			err,
		)
	}
	message.Debugf("specfile contents read:\n%s", string(specFile))

	var spec Spec
	err = yaml.Unmarshal(specFile, &spec)
	if err != nil {
		message.Debugf("error during specfile unmarshal: w", err)
		return Spec{}, fmt.Errorf("there was a problem reading the contents of your vdm spec file: %w", err)
	}
	message.Debugf("vdmSpecs unmarshalled: %+v", spec)

	return spec, nil
}

// OpMsg constructs a loggable message outlining the specific remote details
// being performed at the moment
func (r RemoteTemplate) OpMsg(msg string) {
	message.Infof("%s@%s --> %s: %s", r.Source, r.Version, r.Destination, msg)
}
