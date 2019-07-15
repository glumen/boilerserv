// Package boilerserv provides a boilerplate implementation for a http server.
//
// It also tries not to be opinionated about your packages (e.g. router, logger),
// so you can use your own desired dependencies.
//
// Example:
//		type MyServerConfigurator struct {
//			listener net.Listener
//		}
//
//		func NewMyServerConfigurator(port int) (*MyServerConfigurator, error) {
//			listener, err := net.Listen("tcp", fmt.Springf(":%d", port))
//			if err != nil {
//				return nil, err
//			}
//
//			return &MyServerConfigurator{
//				listener: listener,
//			}, nil
//		}
//
//		func (sc MyServerConfigurator) Listener() net.Listener {
//			return sc.listener
//		}
//
//		func (MyServerConfigurator) ConfigureRoutes() http.Handler {
//			mux := http.ServeMux{}
//			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//				w.WriteHeader(http.StatusOk)
//			})
//
//			return &mux
//		}
//
//		func (sc *MyServerConfigurator) LogError(message string, port int) {
//			log.Println(message)
//		}
//
//		func (sc *MyServerConfigurator) LogInfo(message string, port int) {
//			log.Println(message)
//		}
//
//		func main() {
//			serverConfigurator, err := NewMyServerConfigurator(8080)
//			if err != nil {
//				log.Fatal(err)
//			}
//
//			server, err := boilerserv.StartServer(context.Background(), serverConfigurator)
//			if err != nil {
//				log.Fatal(err)
//			}
//
//			// Block the routine
//			_ = server.ListenAndHandleShutdown(context.Background())
//		}
package boilerserv
