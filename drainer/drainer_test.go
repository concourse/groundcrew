package drainer_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/concourse/groundcrew/drainer"
	"github.com/concourse/groundcrew/drainer/drainerfakes"
	"github.com/concourse/groundcrew/ssh/sshfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Drainer", func() {
	var drainer *Drainer
	var fakeSSHRunner *sshfakes.FakeRunner
	var logger *lagertest.TestLogger
	var fakeWatchProcess *drainerfakes.FakeWatchProcess
	var fakeClock *fakeclock.FakeClock
	var waitInterval time.Duration

	BeforeEach(func() {
		waitInterval = 5 * time.Second
		logger = lagertest.NewTestLogger("drainer")
		fakeSSHRunner = new(sshfakes.FakeRunner)
		fakeWatchProcess = new(drainerfakes.FakeWatchProcess)
		fakeClock = fakeclock.NewFakeClock(time.Unix(0, 123))
	})

	Context("when shutting down", func() {
		BeforeEach(func() {
			drainer = &Drainer{
				SSHRunner:    fakeSSHRunner,
				IsShutdown:   true,
				WatchProcess: fakeWatchProcess,
				Clock:        fakeClock,
			}
		})

		Context("when beacon process is not running", func() {
			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(false, nil)
			})

			It("returns right away", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.RetireWorkerCallCount()).To(Equal(0))
			})
		})

		Context("when failing to check if process is running", func() {
			var disaster = errors.New("disaster")

			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(false, disaster)
			})

			It("returns an error", func() {
				err := drainer.Drain(logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(disaster))
			})
		})

		Context("if watched process is still running", func() {
			BeforeEach(func() {
				callCount := 0
				fakeWatchProcess.IsRunningStub = func(lager.Logger) (bool, error) {
					callCount++
					if callCount > 5 {
						return false, nil
					}

					fakeClock.Increment(waitInterval)
					return true, nil
				}
			})

			It("runs retire-worker until it exits with wait interval", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.RetireWorkerCallCount()).To(Equal(5))
				Expect(fakeSSHRunner.LandWorkerCallCount()).To(Equal(0))
			})
		})
	})

	Context("when not shutting down", func() {
		BeforeEach(func() {
			drainer = &Drainer{
				SSHRunner:    fakeSSHRunner,
				IsShutdown:   false,
				WatchProcess: fakeWatchProcess,
				Clock:        fakeClock,
			}
		})

		Context("when beacon process is not running", func() {
			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(false, nil)
			})

			It("returns right away", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.LandWorkerCallCount()).To(Equal(0))
			})
		})

		Context("when failing to check if process is running", func() {
			var disaster = errors.New("disaster")

			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(false, disaster)
			})

			It("returns an error", func() {
				err := drainer.Drain(logger)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(disaster))
			})
		})

		Context("if watched process is still running", func() {
			BeforeEach(func() {
				callCount := 0
				fakeWatchProcess.IsRunningStub = func(lager.Logger) (bool, error) {
					callCount++
					if callCount > 5 {
						return false, nil
					}

					fakeClock.Increment(waitInterval)
					return true, nil
				}
			})

			It("runs land-worker until it exits with wait interval", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.LandWorkerCallCount()).To(Equal(5))
				Expect(fakeSSHRunner.RetireWorkerCallCount()).To(Equal(0))
			})
		})
	})
})
