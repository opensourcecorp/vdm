package filetree

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetFilePathsInDirectory traverse a directory tree from the provided root, and
// returns a slice of the files in the tree.
func GetFilePathsInDirectory(root string) ([]string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("determining abspath of root %q: %w", rootAbs, err)
	}

	var files []string
	err = filepath.Walk(rootAbs, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory tree: %w", err)
	}

	return files, nil
}
