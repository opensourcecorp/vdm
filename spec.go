package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type vdmSpec struct {
	Remote    string `json:"remote"`
	Version   string `json:"version"`
	LocalPath string `json:"local_path"`
	// so we can pass these around for each spec
	runFlags runFlags `json:"-"`
}

func (spec vdmSpec) writeVDMMeta() error {
	metaFilePath := filepath.Join(spec.LocalPath, "VDMMETA")
	vdmMetaContent, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	os.WriteFile(metaFilePath, vdmMetaContent, 0644)

	return nil
}

func (spec vdmSpec) getVDMMeta() vdmSpec {
	metaFilePath := filepath.Join(spec.LocalPath, "VDMMETA")
	_, err := os.Stat(metaFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return vdmSpec{}
	} else if err != nil {
		errLogger.Fatalf("Couldn't check if VDMMETA exists at '%s': %v", metaFilePath, err)
	}

	vdmMetaFile, err := os.ReadFile(filepath.Join(spec.LocalPath, "VDMMETA"))
	if err != nil {
		if debug {
			debugLogger.Printf("error reading VMDMMETA from disk: %v", err)
		}
		errLogger.Fatalf("There was a problem reading the VDMMETA file from '%s': %v", metaFilePath, err)
	}
	if debug {
		debugLogger.Printf("VDMMETA contents read:\n%s", string(vdmMetaFile))
	}

	var vdmMeta vdmSpec
	err = json.Unmarshal(vdmMetaFile, &vdmMeta)
	if err != nil {
		if debug {
			debugLogger.Printf("error during VDMMETA unmarshal: %v", err)
		}
		errLogger.Fatalf("There was a problem reading the contents of the VDMMETA file at '%s': %v", metaFilePath, err)
	}
	if debug {
		debugLogger.Printf("VDMMETA unmarshalled: %+v", vdmMeta)
	}

	return vdmMeta
}

func getSpecsFromFile(specFilePath string, runFlags runFlags) []vdmSpec {
	specFile, err := os.ReadFile(specFilePath)
	if err != nil {
		if debug {
			debugLogger.Printf("error reading specFile from disk: %v", err)
		}
		errLogger.Fatalf("There was a problem reading your vdm file from '%s' -- does it not exist?", specFilePath)
	}
	if debug {
		debugLogger.Printf("specFile contents read:\n%s", string(specFile))
	}

	var specs []vdmSpec
	err = json.Unmarshal(specFile, &specs)
	if err != nil {
		if debug {
			debugLogger.Printf("error during specFile unmarshal: %v", err)
		}
		errLogger.Fatal("There was a problem reading the contents of your vdm spec file")
	}
	if debug {
		debugLogger.Printf("vdmSpecs unmarshalled: %+v", specs)
	}

	for _, spec := range specs {
		spec.runFlags = runFlags
	}

	return specs
}
