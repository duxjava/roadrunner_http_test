package roadrunner_http_test

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// plugin name
const name = "roadrunner_http_test"

// Plugin structure should have exactly the `Plugin` name to be found by RR
type Plugin struct {
	client http.Client
	server *http.Server
}

func (p *Plugin) Init() error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := http.Client{
		Transport: tr,
		Timeout:   60,
	}
	p.client = client

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	r.Methods("GET").HandlerFunc(p.helloWorld).Path("/hello_world")

	// just as sample, we put server here
	server := &http.Server{
		Addr:           ":8083",
		Handler:        c.Handler(r),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	p.server = server

	return nil
}

func (p *Plugin) Serve() chan error {

	errCh := make(chan error, 1)

	f := p.server.Handler

	p.server.Handler = f

	go func() {
		err := p.server.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		} else {
			errCh <- err
		}
	}()
	return errCh
}

func (p *Plugin) Stop() error {
	err := p.server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// sselect just to not collide with select keyword
func (p *Plugin) helloWorld(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("Hello world"))
}

// Name this is not mandatory, but if you implement this interface and provide a plugin name, RR will expose the RPC method of this plugin using this name
func (p *Plugin) Name() string {
	return name
}

// ----------------------------------------------------------------------------
// RPC
// ----------------------------------------------------------------------------

type rpc struct {
	srv *Plugin
}

// RPC interface implementation, RR will find this interface and automatically expose the RPC endpoint with methods (rpc structure)
func (p *Plugin) RPC() interface{} {
	return &rpc{}
}

// Generate this is the function exposed to PHP $rpc->call(), can be any name
func (r *rpc) Generate(input string, output *string) error {
	*output = input
	return nil
}
