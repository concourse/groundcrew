package groundcrew

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/lager"
)

var ErrSymlinkedParentDir = errors.New("Cannot include symlinks to parent directory.")

type SymlinkFinder struct {
}

func (sf *SymlinkFinder) Find(logger lager.Logger, sourceDir string) ([]string, error) {
	sourceDir = strings.TrimSuffix(sourceDir, "/")
	foundSymlinkedDirs := map[string]struct{}{}

	err := filepath.Walk(sourceDir, func(path string, f os.FileInfo, err error) error {
		if f.Mode()&os.ModeSymlink != 0 {
			symlinkedPath, err := os.Readlink(path)
			if err != nil {
				logger.Error("failed-to-read-symlink-destination", err)
				return err
			}

			foundDir := strings.TrimSuffix(filepath.Dir(symlinkedPath), "/")
			if !filepath.IsAbs(foundDir) {
				foundDir = filepath.Join(sourceDir, foundDir)
			}

			if foundDir == sourceDir {
				return nil
			}

			rel, err := filepath.Rel(foundDir, sourceDir)
			if err == nil && !strings.HasPrefix(rel, "../") {
				logger.Error("failed-to-get-symlinked-directory-for-parent", ErrSymlinkedParentDir, lager.Data{"dir": sourceDir, "symlinked-dir": foundDir})
				return ErrSymlinkedParentDir
			}

			foundSymlinkedDirs[foundDir] = struct{}{}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	symlinkedDirs := []string{}
	for d := range foundSymlinkedDirs {
		// do not include subdirs
		if hasParentDir(foundSymlinkedDirs, d) {
			continue
		}

		symlinkedDirs = append(symlinkedDirs, d)
	}

	return symlinkedDirs, nil
}

func hasParentDir(otherDirs map[string]struct{}, dir string) bool {
	parentDir := dir

	for parentDir != "/" {
		parentDir = filepath.Dir(parentDir)

		if _, ok := otherDirs[parentDir]; ok {
			return true
		}
	}

	return false
}
