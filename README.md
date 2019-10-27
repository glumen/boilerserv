# boilerserv [![Build Status](https://travis-ci.org/glumen/boilerserv.svg?branch=master)](https://travis-ci.org/glumen/boilerserv) [![GoDoc](https://godoc.org/github.com/glumen/boilerserv?status.svg)](https://godoc.org/github.com/glumen/boilerserv)

boilserv is a small library which helps you to setup an absolute boilerplate http server for go.
It can start on any given port provided by a listener and shuts down gracefully.

## Server Configuration

Server is configured by implementing the `ServerConfigurator` interface and providing it the `StartServer` function.

Because the server configuration is implemented by yourself you are not forced to use any specific router or logger. The choice is yours.

### Example

#### MyServerConfigurator
```go
type MyServerConfigurator struct {
    listener net.Listener
}

func NewMyServerConfigurator(port int) (*MyServerConfigurator, error) {
    listener, err := net.Listen("tcp", fmt.Springf(":%d", port))
    if err != nil {
    	return nil, err
    }
	
    return &MyServerConfigurator{
        listener: listener,
    }, nil
}

func (sc MyServerConfigurator) Listener() net.Listener {
    return sc.listener
}

func (MyServerConfigurator) ConfigureRoutes() http.Handler {
    mux := http.ServeMux{}
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOk)
    })
    
    return &mux
}

func (sc *MyServerConfigurator) LogError(message string, port int) {
    log.Println(message)
}

func (sc *MyServerConfigurator) LogInfo(message string, port int) {
    log.Println(message)
}

```

#### Usage

```go
serverConfigurator, err := NewMyServerConfigurator(8080)
if err != nil {
	log.Fatal(err)
}

server, err := boilerserv.StartServer(context.Background(), serverConfigurator)
if err != nil {
	log.Fatal(err)
}

// Block the routine
_ = server.ListenAndHandleShutdown(context.Background())
```

## LICENSE

MIT License

Copyright (c) 2019 Patric Vormstein

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
