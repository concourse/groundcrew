package drainer

import (
	"fmt"
	"io/ioutil"
	"os"
)

//go:generate counterfeiter . WatchProcess

type WatchProcess interface {
	IsRunning() (bool, error)
}

type beaconWatchProcess struct {
	pidFile string
}

func NewBeaconWatchProcess(pidFile string) WatchProcess {
	return &beaconWatchProcess{
		pidFile: pidFile,
	}
}

func (p *beaconWatchProcess) IsRunning() (bool, error) {
	_, err := os.Stat(p.pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	pidContents, err := ioutil.ReadFile(p.pidFile)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%s", string(pidContents)))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
