package boilerserv

import (
	"context"
	"net"
	"net/http"
)

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
