package integration

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestRoutes(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Routes Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// fire up app (this sleeps a bit)
	gomega.Expect(fireUpApp()).To(gomega.Succeed(), "App should fire up")
})

var _ = ginkgo.AfterSuite(func() {
	// kill the process
	gomega.Expect(killApp()).To(gomega.Succeed(), "App should shutdown gracefully")
})
