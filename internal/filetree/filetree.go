package filetree

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetFilePathsInDirectory traverse a directory tree from the provided root, and
// returns a slice of the files in the tree.
func GetFilePathsInDirectory(rootDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootDir, func(path string, f os.FileInfo, err error) error {
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
