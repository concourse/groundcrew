package drainer

import (
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/groundcrew/ssh"
)

type Drainer struct {
	SSHRunner    ssh.Runner
	IsShutdown   bool
	WatchProcess WatchProcess
	WaitInterval time.Duration
}

func (d *Drainer) Drain(logger lager.Logger) error {
	processIsRunning, err := d.WatchProcess.IsRunning()
	if err != nil {
		logger.Error("failed-to-check-if-process-is-running", err)
		return err
	}

	if !processIsRunning {
		return nil
	}

	if d.IsShutdown {
		return d.SSHRunner.RetireWorker(logger)
	}

	return d.SSHRunner.LandWorker(logger)
}
