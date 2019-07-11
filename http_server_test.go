package boilerserv

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestHttpServer(t *testing.T) {
	t.Run("should start and stop server successfully", func(t *testing.T) {
		tt := NewGomegaWithT(t)

		// Create server configurator and listener
		serverConfigurator := newMockServerConfigurator(tt)
		listener := serverConfigurator.Listener()
		assignedPort := listener.Addr().(*net.TCPAddr).Port
		tt.Expect(assignedPort).To(Not(Equal(0)))

		// Create context to be able to cancel the server
		ctx, cancelFunc := context.WithCancel(context.Background())

		// Start server
		_, err := StartServer(ctx, serverConfigurator)
		tt.Expect(err).To(Not(HaveOccurred()))
		tt.Eventually(func() []byte {
			return serverConfigurator.logger.Bytes()
		}, 100*time.Millisecond).Should(Equal([]byte("starting http server")))

		// Test if port is in use
		assignedPortIsInUse := portIsInUse(assignedPort)
		tt.Expect(assignedPortIsInUse).To(BeTrue())

		// Test Endpoint
		client := http.Client{}
		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", assignedPort), nil)
		tt.Expect(err).To(Not(HaveOccurred()))
		response, err := client.Do(request)
		tt.Expect(err).To(Not(HaveOccurred()))

		tt.Expect(response.StatusCode).To(Equal(http.StatusCreated))

		// Shutdown server
		cancelFunc()
		tt.Eventually(func() bool {
			return portIsInUse(assignedPort)
		}, 100*time.Millisecond).Should(BeFalse())
		tt.Expect(serverConfigurator.logger.Bytes()).To(Equal([]byte("shutting down http server gracefully")))

	})
}

func portIsInUse(port int) bool {
	p := fmt.Sprintf(":%d", port)
	_, err := net.Listen("tcp", p)
	return err != nil
}

type mockServerConfigurator struct {
	tt       *GomegaWithT
	logger   *bytes.Buffer
	listener net.Listener
}

func newMockServerConfigurator(tt *GomegaWithT) *mockServerConfigurator {
	listener, err := net.Listen("tcp", ":0")
	tt.Expect(err).To(Not(HaveOccurred()))

	return &mockServerConfigurator{
		tt:       tt,
		logger:   bytes.NewBuffer(nil),
		listener: listener,
	}
}

func (sc mockServerConfigurator) Listener() net.Listener {
	return sc.listener
}

func (mockServerConfigurator) ConfigureRoutes() http.Handler {
	mux := http.ServeMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	return &mux
}

func (sc *mockServerConfigurator) LogError(message string, port int) {
	sc.logger.Reset()
	sc.logger.Write([]byte(message))
}

func (sc *mockServerConfigurator) LogInfo(message string, port int) {
	sc.logger.Reset()
	sc.logger.Write([]byte(message))
}
