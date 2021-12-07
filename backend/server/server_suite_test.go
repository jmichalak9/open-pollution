package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jmichalak9/open-pollution/server"
	"github.com/jmichalak9/open-pollution/server/measurement"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var (
	port  string
	ctrl  *gomock.Controller
	cache *measurement.MockCache
	srv   *server.Server
)

var _ = BeforeEach(func() {
	ctrl = gomock.NewController(GinkgoT())
	cache = measurement.NewMockCache(ctrl)
	port = "65123"
	srv = server.NewServer(
		":"+port,
		cache,
	)
	go srv.Run()
	Eventually(func() bool {
		_, err := http.Get(fmt.Sprintf("http://localhost:%v", port))
		return err == nil
	}).Should(BeTrue())
})

var _ = AfterEach(func() {
	err := srv.Shutdown()
	Expect(err).NotTo(HaveOccurred())
})
