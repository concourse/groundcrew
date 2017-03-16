package groundcrew_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/lager/lagertest"

	"github.com/concourse/groundcrew"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SymlinkFinder", func() {
	var (
		symlinkFinder    *groundcrew.SymlinkFinder
		tmpDir           string
		certificatesPath string
		logger           *lagertest.TestLogger
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("groundcrew-test")
		var err error
		tmpDir, err = ioutil.TempDir("", "groundcrew")
		Expect(err).NotTo(HaveOccurred())

		certificatesPath = filepath.Join(tmpDir, "certificates-path")
		err = os.MkdirAll(certificatesPath, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		simpleFile := filepath.Join(certificatesPath, "simple-file")
		_, err = os.OpenFile(simpleFile, os.O_RDONLY|os.O_CREATE, 0666)
		Expect(err).NotTo(HaveOccurred())

		anotherSimpleFile := filepath.Join(certificatesPath, "another-simple-file")
		_, err = os.OpenFile(anotherSimpleFile, os.O_RDONLY|os.O_CREATE, 0666)
		Expect(err).NotTo(HaveOccurred())

		err = os.Symlink(simpleFile, filepath.Join(certificatesPath, "symlink-in-current-dir"))
		Expect(err).NotTo(HaveOccurred())

		err = os.Symlink("./another-simple-file", filepath.Join(certificatesPath, "relative-symlink-in-current-dir"))
		Expect(err).NotTo(HaveOccurred())

		symlinkFinder = &groundcrew.SymlinkFinder{}
	})

	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})

	Context("when provided path does not contain symlinked directories", func() {
		It("returns empty list", func() {
			symlinkedDirs, err := symlinkFinder.Find(logger, certificatesPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(symlinkedDirs).To(BeEmpty())
		})
	})

	Context("when provided path contains symlinked directories", func() {
		var (
			symlinkedParentDir  string
			anotherSymlinkedDir string
		)

		BeforeEach(func() {
			symlinkedParentDir = filepath.Join(tmpDir, "symlinked-path")
			err := os.MkdirAll(symlinkedParentDir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			symlinkedFileInParentDir := filepath.Join(symlinkedParentDir, "symlinked-file")

			_, err = os.OpenFile(symlinkedFileInParentDir, os.O_RDONLY|os.O_CREATE, 0666)
			Expect(err).NotTo(HaveOccurred())

			err = os.MkdirAll(filepath.Join(symlinkedParentDir, "subdir"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			symlinkedFileinSubDir := filepath.Join(symlinkedParentDir, "subdir", "symlinked-file-in-subdir")

			_, err = os.OpenFile(symlinkedFileinSubDir, os.O_RDONLY|os.O_CREATE, 0666)
			Expect(err).NotTo(HaveOccurred())

			anotherSymlinkedFileinSubDir := filepath.Join(symlinkedParentDir, "subdir", "another-symlinked-file-in-subdir")

			_, err = os.OpenFile(anotherSymlinkedFileinSubDir, os.O_RDONLY|os.O_CREATE, 0666)
			Expect(err).NotTo(HaveOccurred())

			anotherSymlinkedDir = filepath.Join(tmpDir, "another-symlinked-path")
			err = os.MkdirAll(anotherSymlinkedDir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			anotherSymlinkedFile := filepath.Join(anotherSymlinkedDir, "another-symlinked-file")

			_, err = os.OpenFile(anotherSymlinkedFile, os.O_RDONLY|os.O_CREATE, 0666)
			Expect(err).NotTo(HaveOccurred())

			err = os.Symlink(symlinkedFileInParentDir, filepath.Join(certificatesPath, "symlinked-file-in-parent-dir"))
			Expect(err).NotTo(HaveOccurred())

			err = os.Symlink(symlinkedFileinSubDir, filepath.Join(certificatesPath, "symlinked-file-in-subdir"))
			Expect(err).NotTo(HaveOccurred())

			err = os.Symlink(anotherSymlinkedFileinSubDir, filepath.Join(certificatesPath, "another-symlinked-file-in-subdir"))
			Expect(err).NotTo(HaveOccurred())

			err = os.Symlink(anotherSymlinkedFile, filepath.Join(certificatesPath, "another-symlinked-file"))
			Expect(err).NotTo(HaveOccurred())
		})

		It("includes symlinked directories excluding duplicates and subdirectories", func() {
			symlinkedDirs, err := symlinkFinder.Find(logger, certificatesPath)
			Expect(err).NotTo(HaveOccurred())

			Expect(symlinkedDirs).To(ConsistOf([]string{
				symlinkedParentDir,
				anotherSymlinkedDir,
			}))
		})

		Context("when there is a symlink to parent directory", func() {
			BeforeEach(func() {
				fileInParentDir := filepath.Join(tmpDir, "file-in-parent-dir")
				_, err := os.OpenFile(fileInParentDir, os.O_RDONLY|os.O_CREATE, 0666)
				Expect(err).NotTo(HaveOccurred())

				err = os.Symlink(fileInParentDir, filepath.Join(certificatesPath, "symlink-to-parent-dir"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				_, err := symlinkFinder.Find(logger, certificatesPath)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
