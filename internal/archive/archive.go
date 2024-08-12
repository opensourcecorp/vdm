package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensourcecorp/vdm/internal/filetree"
)

// Much of the following functions taken from:
// https://www.arthurkoziel.com/writing-tar-gz-files-in-go/

// CreateArchive writes a gzipped tarball based on the provided root directory
// from which to construct the archive, and its target file name.
func CreateArchive(rootDir string, archivePath string) (err error) {
	if !strings.HasSuffix(archivePath, ".tar.gz") {
		return errors.New("provided archive path must end in .tar.gz")
	}

	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("determining abspath of provided root dir %s: %w", rootDir, err)
	}

	archivePathAbs, err := filepath.Abs(archivePath)
	if err != nil {
		return fmt.Errorf("determining abspath of provided archive path %s: %w", archivePath, err)
	}

	buf, err := os.Create(archivePathAbs)
	if err != nil {
		return fmt.Errorf("opening target archive path: %w", err)
	}
	defer func() {
		if closeErr := buf.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing target archive file: %w", closeErr))
		}
	}()

	// gzip and tar get their own writers, and note that they are chained -- tar
	// writes to gzip, which writes to the buffer. So, we only need to write to
	// the tar writer
	gzipWriter := gzip.NewWriter(buf)
	defer func() {
		if closeErr := gzipWriter.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing gzip writer: %w", closeErr))
		}
	}()

	tarWriter := tar.NewWriter(gzipWriter)
	defer func() {
		if closeErr := tarWriter.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing tar writer: %w", closeErr))
		}
	}()

	files, err := filetree.GetFilePathsInDirectory(rootDirAbs)
	if err != nil {
		return fmt.Errorf("populating list of files from %s: %w", rootDir, err)
	}

	for _, fileName := range files {
		err := addToArchive(tarWriter, rootDirAbs, fileName)
		if err != nil {
			return fmt.Errorf("adding %s to archive: %w", fileName, err)
		}
	}

	return err
}

func addToArchive(tarWriter *tar.Writer, rootDirAbs string, fileName string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("opening file %s: %w", fileName, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing file: %w", closeErr))
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info for %s: %w", fileName, err)
	}

	// Tar needs file headers, so create one from the file info
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return fmt.Errorf("creating tar header for %s: %w", fileName, err)
	}

	topLevelDir, err := maybeGetTopLevelDir(rootDirAbs)
	if err != nil {
		return fmt.Errorf("checking for top-level directory when adding to archive: %w", err)
	}
	fileNameClean := strings.ReplaceAll(fileName, rootDirAbs+string(filepath.Separator), "")
	if topLevelDir != "" {
		fileNameClean = filepath.Join(topLevelDir, fileNameClean)
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = fileNameClean

	// Write file header to the tar archive
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("writing tar header for %s: %w", fileName, err)
	}

	// Copy file content to tar archive
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("adding file %s to tar archive: %w", fileName, err)
	}

	return err
}

func maybeGetTopLevelDir(rootDirAbs string) (topLevelDir string, err error) {
	topLevelFileInfo, err := os.Stat(rootDirAbs)
	if err != nil {
		return "", fmt.Errorf("getting file info for %s: %w", rootDirAbs, err)
	}

	// We only want to include the top-level directory for archive pathing if
	// it's *actually* a directory, obviously
	if !topLevelFileInfo.IsDir() {
		topLevelDir = ""
	} else {
		topLevelDir = filepath.Base(rootDirAbs)
	}

	return topLevelDir, nil
}
