package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Much of the following functions taken from:
// https://www.arthurkoziel.com/writing-tar-gz-files-in-go/
func createArchive(files []string, archivePath string) (err error) {
	if !strings.HasSuffix(archivePath, ".tar.gz") {
		return errors.New("provided archive path must end in .tar.gz")
	}

	buf, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("opening target archive path: %w", err)
	}
	defer func() {
		if closeErr := buf.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing target archive file: %w", closeErr))
		}
	}()

	// gzip and tar get their own writers, and note that they are chained -- tar
	// writes to gzip, which writes to the buffer
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

	for _, fileName := range files {
		err := addToArchive(tarWriter, fileName)
		if err != nil {
			return fmt.Errorf("adding %s to archive: %w", fileName, err)
		}
	}

	return err
}

func addToArchive(tarWriter *tar.Writer, fileName string) (err error) {
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

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = fileName

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
