package boilerserv

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Server is an interface for any server that listens and can be shutdown.
// This interface can be used for testing.
type Server interface {
	ShutdownGracefully(ctx context.Context) error
	ListenAndHandleShutdown(ctx context.Context) error
}

// HttpServer wraps a http.Server and can handle graceful shutdown automatically.
// It implements the Server interface.
type HttpServer struct {
	internalHttpServer *http.Server
	configurator       ServerConfigurator
	port               int
}

// StartServer starts a new HttpServer by providing an implementation of a ServerConfigurator.
func StartServer(ctx context.Context, configurator ServerConfigurator) (*HttpServer, error) {
	port := configurator.Listener().Addr().(*net.TCPAddr).Port
	s := http.Server{
		Handler: configurator.ConfigureRoutes(),
	}

	// Start server
	c := make(chan error)
	go func() {
		configurator.LogInfo("starting http server", port)
		if err := s.Serve(configurator.Listener()); err != nil {
			c <- err
		}
	}()

	// Server object and interruption routine
	server := &HttpServer{
		internalHttpServer: &s,
		configurator:       configurator,
		port:               port,
	}

	go func() {
		<-ctx.Done()
		if err := server.ShutdownGracefully(ctx); err != nil {
			configurator.LogError("could not shutdown http server gracefully", port)
		}
	}()

	select {
	case err := <-c:
		return nil, err
	default:
		return server, nil
	}
}

// ShutdownGracefully will shutdown the HttpServer gracefully.
func (s *HttpServer) ShutdownGracefully(ctx context.Context) error {
	s.configurator.LogInfo("shutting down http server gracefully", s.port)
	return s.internalHttpServer.Shutdown(ctx)
}

// ListenAndHandleShutdown listens on the provided port of the provided net.Listener.
// It will handle shutdown when receiving an interruption signal.
func (s *HttpServer) ListenAndHandleShutdown(ctx context.Context) error {
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-kill
	if err := s.ShutdownGracefully(ctx); err != nil {
		return err
	}

	return nil
}
