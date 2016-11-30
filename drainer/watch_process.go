package drainer

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
	return false, nil
}
