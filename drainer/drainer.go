package drainer

import (
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/groundcrew/ssh"
)

type Drainer struct {
	SSHRunner    ssh.Runner
	IsShutdown   bool
	WatchProcess WatchProcess
	WaitInterval time.Duration
	Clock        clock.Clock
}

func (d *Drainer) Drain(logger lager.Logger) error {
	for {
		processIsRunning, err := d.WatchProcess.IsRunning()
		if err != nil {
			logger.Error("failed-to-check-if-process-is-running", err)
			return err
		}

		if !processIsRunning {
			return nil
		}

		if d.IsShutdown {
			err := d.SSHRunner.RetireWorker(logger)
			if err != nil {
				logger.Error("failed-to-retire-worker", err)
			}
		}

		err = d.SSHRunner.LandWorker(logger)
		if err != nil {
			logger.Error("failed-to-land-worker", err)
		}

		d.Clock.Sleep(d.WaitInterval)
	}
}
