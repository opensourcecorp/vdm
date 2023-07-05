package vdmspec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type VDMSpec struct {
	Remote    string `json:"remote"`
	Version   string `json:"version,omitempty"`
	LocalPath string `json:"local_path"`
	Type      string `json:"type,omitempty"`
}

const MetaFileName = "VDMMETA"

func (spec VDMSpec) MakeMetaFilePath() string {
	metaFilePath := filepath.Join(spec.LocalPath, MetaFileName)
	// TODO: this is brittle, but it's the best I can think of right now
	if spec.Type == "file" {
		// converts to e.g. 'VDMMETA_http.proto'
		metaFilePath = fmt.Sprintf("%s_%s", MetaFileName, filepath.Base(spec.LocalPath))
	}

	return metaFilePath
}

func (spec VDMSpec) WriteVDMMeta() error {
	metaFilePath := spec.MakeMetaFilePath()
	vdmMetaContent, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("writing %s: %w", metaFilePath, err)
	}

	vdmMetaContent = append(vdmMetaContent, []byte("\n")...)

	logrus.Debugf("writing metadata file to '%s'", metaFilePath)
	err = os.WriteFile(metaFilePath, vdmMetaContent, 0644)
	if err != nil {
		return fmt.Errorf("writing metadata file: %w", err)
	}

	return nil
}

func (spec VDMSpec) GetVDMMeta() (VDMSpec, error) {
	metaFilePath := spec.MakeMetaFilePath()
	_, err := os.Stat(metaFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return VDMSpec{}, nil // this is ok, because it might literally not exist yet
	} else if err != nil {
		return VDMSpec{}, fmt.Errorf("couldn't check if %s exists at '%s': %w", MetaFileName, metaFilePath, err)
	}

	vdmMetaFile, err := os.ReadFile(filepath.Join(spec.LocalPath, MetaFileName))
	if err != nil {
		logrus.Debugf("error reading VMDMMETA from disk: %v", err)
		return VDMSpec{}, fmt.Errorf("there was a problem reading the %s file from '%s': %w", MetaFileName, metaFilePath, err)
	}
	logrus.Debugf("%s contents read:\n%s", MetaFileName, string(vdmMetaFile))

	var vdmMeta VDMSpec
	err = json.Unmarshal(vdmMetaFile, &vdmMeta)
	if err != nil {
		logrus.Debugf("error during %s unmarshal: %v", MetaFileName, err)
		return VDMSpec{}, fmt.Errorf("there was a problem reading the contents of the %s file at '%s': %v", MetaFileName, metaFilePath, err)
	}
	logrus.Debugf("file %s unmarshalled: %+v", MetaFileName, vdmMeta)

	return vdmMeta, nil
}

func GetSpecsFromFile(specFilePath string) ([]VDMSpec, error) {
	specFile, err := os.ReadFile(specFilePath)
	if err != nil {
		logrus.Debugf("error reading specfile from disk: %v", err)
		return nil, fmt.Errorf(
			strings.Join([]string{
				"there was a problem reading your vdm file from '%s' -- does it not exist?",
				"Either pass the -spec-file flag, or create one in the default location (details in the README).",
				"Error details: %w"},
				" ",
			),
			specFilePath,
			err,
		)
	}
	logrus.Debugf("specfile contents read:\n%s", string(specFile))

	var specs []VDMSpec
	err = json.Unmarshal(specFile, &specs)
	if err != nil {
		logrus.Debugf("error during specfile unmarshal: %v", err)
		return nil, fmt.Errorf("there was a problem reading the contents of your vdm spec file: %w", err)
	}
	logrus.Debugf("vdmSpecs unmarshalled: %+v", specs)

	return specs, nil
}

// OpMsg constructs a loggable message outlining the specific operation being
// performed at the moment
func (spec VDMSpec) OpMsg() string {
	return fmt.Sprintf("%s@%s --> %s", spec.Remote, spec.Version, spec.LocalPath)
}
