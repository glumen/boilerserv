package boilerserv

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	ShutdownGracefully(ctx context.Context) error
	ListenAndHandleShutdown(ctx context.Context) error
}

type HttpServer struct {
	internalHttpServer *http.Server
	configurator       ServerConfigurator
	port               int
}

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

func (s *HttpServer) ShutdownGracefully(ctx context.Context) error {
	s.configurator.LogInfo("shutting down http server gracefully", s.port)
	return s.internalHttpServer.Shutdown(ctx)
}

func (s *HttpServer) ListenAndHandleShutdown(ctx context.Context) error {
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-kill
	if err := s.ShutdownGracefully(ctx); err != nil {
		return err
	}

	return nil
}
