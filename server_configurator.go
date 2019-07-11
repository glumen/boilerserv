package boilerserv

import (
	"net"
	"net/http"
)

// ServerConfigurator is an interface which can be used to configure the server
type ServerConfigurator interface {
	Listener() net.Listener
	ConfigureRoutes() http.Handler
	LogError(message string, port int)
	LogInfo(message string, port int)
}
