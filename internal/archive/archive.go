package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opensourcecorp/vdm/internal/filetree"
	"github.com/opensourcecorp/vdm/internal/message"
)

var (
	tgzRegex = regexp.MustCompile(`(\.tar\.gz|\.tgz)$`)
	zipRegex = regexp.MustCompile(`\.zip$`)
)

// Much of the following functions taken from:
// https://www.arthurkoziel.com/writing-tar-gz-files-in-go/

// CreateArchive writes a gzipped tarball based on the provided root directory
// from which to construct the archive, and its target file name. It returns an
// open file handle to the archive, which should be closed by the caller.
func CreateArchive(root string, archivePath string) (f *os.File, err error) {
	if !tgzRegex.MatchString(archivePath) {
		return nil, errors.New("provided archive path must have valid gzipped-tar extension")
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("determining abspath of provided root dir %q: %w", root, err)
	}
	message.Debugf("root: %q, rootAbs: %q", root, rootAbs)

	archivePathAbs, err := filepath.Abs(archivePath)
	if err != nil {
		return nil, fmt.Errorf("determining abspath of provided archive path %q: %w", archivePath, err)
	}

	buf, err := os.Create(archivePathAbs)
	if err != nil {
		return nil, fmt.Errorf("opening target archive path: %w", err)
	}
	// NOTE: file not closed because it's returned, open, to the caller

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

	files, err := filetree.GetFilePathsInDirectory(rootAbs)
	if err != nil {
		return nil, fmt.Errorf("populating list of files from %q: %w", root, err)
	}

	for _, fileName := range files {
		err := addToArchive(tarWriter, rootAbs, fileName)
		if err != nil {
			return nil, fmt.Errorf("adding %q to archive: %w", fileName, err)
		}
	}

	return buf, err
}

func ExtractTGZArchive(src, dest string) error {
	if !tgzRegex.MatchString(src) {
		return errors.New("provided archive path must have valid gzipped-tar extension")
	}

	gzipFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening archive file: %w", err)
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Next() failed in tar reader: %w", err)
		}

		filePath := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return fmt.Errorf("making directory during tar extraction: %w", err)
			}
		case tar.TypeReg:
			dirPath := filepath.Dir(filePath)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("making directory during tar extraction: %w", err)
			}

			outFile, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("creating output file %q during tar extraction: %w", filePath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("writing to output file %q during tar extraction: %w", filePath, err)
			}
			closeErr := outFile.Close()
			if closeErr != nil {
				return fmt.Errorf("closing output file used during tar extraction: %w", closeErr)
			}
		default:
			return fmt.Errorf("unknown tar type '%b' in %q", header.Typeflag, filePath)
		}
	}

	return nil
}

func addToArchive(tarWriter *tar.Writer, rootAbs string, fileName string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", fileName, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing file: %w", closeErr))
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info for %q: %w", fileName, err)
	}

	// Tar needs file headers, so create one from the file info
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return fmt.Errorf("creating tar header for %q: %w", fileName, err)
	}

	topLevelDir, err := maybeGetTopLevelDir(rootAbs)
	if err != nil {
		return fmt.Errorf("checking for top-level directory when adding to archive: %w", err)
	}
	fileNameClean := strings.ReplaceAll(fileName, rootAbs+string(filepath.Separator), "")
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
		return fmt.Errorf("writing tar header for %q: %w", fileName, err)
	}

	// Copy file content to tar archive
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("adding file %q to tar archive: %w", fileName, err)
	}

	return err
}

func maybeGetTopLevelDir(rootAbs string) (string, error) {
	topLevelDir := filepath.Base(rootAbs)
	return topLevelDir, nil
}
