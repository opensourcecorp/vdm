package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type vdmSpec struct {
	Remote    string `json:"remote"`
	Version   string `json:"version"`
	LocalPath string `json:"local_path"`
	Type      string `json:"type"`
}

func (spec vdmSpec) writeVDMMeta() error {
	metaFilePath := filepath.Join(spec.LocalPath, "VDMMETA")
	vdmMetaContent, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	vdmMetaContent = append(vdmMetaContent, []byte("\n")...)
	os.WriteFile(metaFilePath, vdmMetaContent, 0644)

	return nil
}

func (spec vdmSpec) getVDMMeta() vdmSpec {
	metaFilePath := filepath.Join(spec.LocalPath, "VDMMETA")
	_, err := os.Stat(metaFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return vdmSpec{}
	} else if err != nil {
		logrus.Fatalf("Couldn't check if VDMMETA exists at '%s': %v", metaFilePath, err)
	}

	vdmMetaFile, err := os.ReadFile(filepath.Join(spec.LocalPath, "VDMMETA"))
	if err != nil {
		if debug {
			logrus.Debugf("error reading VMDMMETA from disk: %v", err)
		}
		logrus.Fatalf("There was a problem reading the VDMMETA file from '%s': %v", metaFilePath, err)
	}
	if debug {
		logrus.Debugf("VDMMETA contents read:\n%s", string(vdmMetaFile))
	}

	var vdmMeta vdmSpec
	err = json.Unmarshal(vdmMetaFile, &vdmMeta)
	if err != nil {
		if debug {
			logrus.Debugf("error during VDMMETA unmarshal: %v", err)
		}
		logrus.Fatalf("There was a problem reading the contents of the VDMMETA file at '%s': %v", metaFilePath, err)
	}
	if debug {
		logrus.Debugf("VDMMETA unmarshalled: %+v", vdmMeta)
	}

	return vdmMeta
}

func getSpecsFromFile(ctx context.Context, specFilePath string) []vdmSpec {
	specFile, err := os.ReadFile(specFilePath)
	if err != nil {
		if isDebug(ctx) {
			logrus.Debugf("error reading specFile from disk: %v", err)
		}
		logrus.Fatalf("There was a problem reading your vdm file from '%s' -- does it not exist? Either pass the -spec-file flag, or create one in the default location (details in the README)", specFilePath)
	}
	if debug {
		logrus.Debugf("specFile contents read:\n%s", string(specFile))
	}

	var specs []vdmSpec
	err = json.Unmarshal(specFile, &specs)
	if err != nil {
		if isDebug(ctx) {
			logrus.Debugf("error during specFile unmarshal: %v", err)
		}
		logrus.Fatal("There was a problem reading the contents of your vdm spec file")
	}
	if isDebug(ctx) {
		logrus.Debugf("vdmSpecs unmarshalled: %+v", specs)
	}

	return specs
}
