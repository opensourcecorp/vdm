package vdmspec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/vdm/internal/message"
	"gopkg.in/yaml.v3"
)

// Spec defines the overall structure of the vmd specfile.
type Spec struct {
	Remotes []Remote `json:"remotes" yaml:"remotes"`
}

// Remote defines the structure of each remote configuration in the vdm
// specfile.
type Remote struct {
	// Type is the type of Source, e.g. git, archive, etc.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Source is the fully-qualifed location from which the Remote is retrieved,
	// e.g. "https://github.com/some-org/some-repo"
	Source string `json:"source" yaml:"source"`
	// Version states the version requested from Source, and is then later used
	// for tracking purposes. Version can be anything supported by the Type
	// field -- for example, for the "git" Type, this can be a tag, a branch
	// name, or a commit hash. It can also be the word "latest".
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// Destination is the relative or absolute path on disk that Source will be
	// placed at
	Destination string `json:"destination" yaml:"destination"`
	// TryLocalSource helps define behavior driven by the `try-local-sources`
	// CLI flag, which allows checking for a local version of a Remote Source,
	// and falling back to the other Source field if the local path does not
	// exist. This is especially useful for when you might be developing one of
	// your Remotes in a nearby directory, and want to copy over that version of
	// the Remote and not keep pushing-and-pulling to a Git upstream just to
	// test the changes.
	TryLocalSource string `json:"try_local_source" yaml:"try_local_source"`
}

const (
	// MetaFileName is the name of the tracking file that vdm uses to record &
	// track remote statuses on disk.
	MetaFileName string = "VDMMETA"

	// GitType represents the string to match against for git remote types.
	GitType string = "git"
	// FileType represents the string to match against for file remote types.
	FileType string = "file"
)

// MakeMetaFilePath constructs the metafile path that vdm will use to track a
// remote's state on disk.
func (r Remote) MakeMetaFilePath() string {
	metaFilePath := filepath.Join(r.Destination, MetaFileName)
	// TODO: this is brittle, but it's the best I can think of right now
	if r.Type == FileType {
		fileDir := filepath.Dir(r.Destination)
		fileName := filepath.Base(r.Destination)
		// converts to e.g. 'VDMMETA_http.proto'
		metaFilePath = filepath.Join(fileDir, fmt.Sprintf("%s_%s", MetaFileName, fileName))
	}

	return metaFilePath
}

// WriteVDMMeta writes the metafile contents to disk, the path of which is
// determined by [Remote.MakeMetaFilePath].
func (r Remote) WriteVDMMeta() error {
	metaFilePath := r.MakeMetaFilePath()
	vdmMetaContent, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("writing %s: %w", metaFilePath, err)
	}

	vdmMetaContent = append(vdmMetaContent, []byte("\n")...)

	message.Debugf("writing metadata file to '%s'", metaFilePath)
	err = os.WriteFile(metaFilePath, vdmMetaContent, 0644)
	if err != nil {
		return fmt.Errorf("writing metadata file: %w", err)
	}

	return nil
}

// GetVDMMeta reads the metafile from disk, and returns it for further
// processing.
func (r Remote) GetVDMMeta() (Remote, error) {
	metaFilePath := r.MakeMetaFilePath()
	_, err := os.Stat(metaFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return Remote{}, nil // this is ok, because it might literally not exist yet
	} else if err != nil {
		return Remote{}, fmt.Errorf("couldn't check if %s exists at '%s': %w", MetaFileName, metaFilePath, err)
	}

	vdmMetaFile, err := os.ReadFile(metaFilePath)
	if err != nil {
		message.Debugf("error reading VMDMMETA from disk: %w", err)
		return Remote{}, fmt.Errorf("there was a problem reading the %s file from '%s': %w", MetaFileName, metaFilePath, err)
	}
	message.Debugf("%s contents read:\n%s", MetaFileName, string(vdmMetaFile))

	var vdmMeta Remote
	err = yaml.Unmarshal(vdmMetaFile, &vdmMeta)
	if err != nil {
		message.Debugf("error during %s unmarshal: w", MetaFileName, err)
		return Remote{}, fmt.Errorf("there was a problem reading the contents of the %s file at '%s': %w", MetaFileName, metaFilePath, err)
	}
	message.Debugf("file %s unmarshalled: %+v", MetaFileName, vdmMeta)

	return vdmMeta, nil
}

// GetSpecFromFile reads the specfile from disk (the path of which is determined
// by the user-supplied flag value), and returns it for further processing of
// remotes.
func GetSpecFromFile(specFilePath string) (Spec, error) {
	specFile, err := os.ReadFile(specFilePath)
	if err != nil {
		message.Debugf("error reading specfile from disk: %w", err)
		return Spec{}, fmt.Errorf(
			strings.Join([]string{
				"there was a problem reading your vdm file from '%s' -- does it not exist?",
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

// OpMsg constructs a loggable message outlining the specific operation being
// performed at the moment
func (r Remote) OpMsg() string {
	if r.Version != "" {
		return fmt.Sprintf("%s@%s --> %s", r.Source, r.Version, r.Destination)
	}
	return fmt.Sprintf("%s --> %s", r.Source, r.Destination)
}
