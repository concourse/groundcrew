package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc"
	"github.com/concourse/groundcrew"
)

type AddCertificateSymlinksCommand struct {
	LogPath string `long:"log-path" description:"Path to log file." required:"yes"`

	Args struct {
		CertificatesPath string
	} `positional-args:"yes" required:"yes"`
}

func (cmd *AddCertificateSymlinksCommand) Execute(args []string) error {
	logger := lager.NewLogger("groundcrew")
	logFile, err := os.OpenFile(cmd.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.RegisterSink(lager.NewWriterSink(logFile, lager.DEBUG))

	logger = logger.Session("add-ceritificate-symlinks")

	workerJSON, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logger.Error("failed-to-read-stdin", err)
		return err
	}

	var worker atc.Worker
	err = json.Unmarshal(workerJSON, &worker)
	if err != nil {
		logger.Error("failed-to-parse-stdin", err)
		return err
	}

	symlinkFinder := &groundcrew.SymlinkFinder{}
	symlinkedDirs, err := symlinkFinder.Find(logger, cmd.Args.CertificatesPath)
	if err != nil {
		logger.Error("failed-to-find-symlinked-directories", err)
		return err
	}

	logger.Debug("found-symlinked-dirs", lager.Data{"dirs": symlinkedDirs})

	worker.CertificatesSymlinkedPaths = symlinkedDirs

	updatedJSON, err := json.Marshal(worker)
	if err != nil {
		logger.Error("failed-to-marshal-payload", err)
		return err
	}

	os.Stdout.Write(updatedJSON)

	return nil
}
