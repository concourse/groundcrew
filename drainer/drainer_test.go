package drainer_test

import (
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

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("drainer")
		fakeSSHRunner = new(sshfakes.FakeRunner)
		fakeWatchProcess = new(drainerfakes.FakeWatchProcess)
	})

	Context("when shutting down", func() {
		BeforeEach(func() {
			drainer = &Drainer{
				SSHRunner:    fakeSSHRunner,
				IsShutdown:   true,
				WatchProcess: fakeWatchProcess,
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

		Context("when beacon process is running", func() {
			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(true, nil)
			})

			It("runs retire-worker ssh command", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.RetireWorkerCallCount()).To(Equal(1))
			})
		})
	})

	Context("when not shutting down", func() {
		BeforeEach(func() {
			drainer = &Drainer{
				SSHRunner:    fakeSSHRunner,
				IsShutdown:   false,
				WatchProcess: fakeWatchProcess,
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

		Context("when beacon process is running", func() {
			BeforeEach(func() {
				fakeWatchProcess.IsRunningReturns(true, nil)
			})

			It("runs land-worker ssh command", func() {
				err := drainer.Drain(logger)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSSHRunner.LandWorkerCallCount()).To(Equal(1))
			})
		})
	})
})
