package roadrunner_http_test

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Http struct {
	client http.Client
	server *http.Server
}

func (h *Http) Init() error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := http.Client{
		Transport: tr,
		Timeout:   60,
	}
	h.client = client

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	r.Methods("GET").HandlerFunc(h.helloWorld).Path("/hello_world")

	// just as sample, we put server here
	server := &http.Server{
		Addr:           ":8083",
		Handler:        c.Handler(r),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	h.server = server

	return nil
}

func (h *Http) Serve() chan error {
	errCh := make(chan error, 1)

	f := h.server.Handler

	h.server.Handler = f

	go func() {
		err := h.server.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		} else {
			errCh <- err
		}
	}()
	return errCh
}

func (h *Http) Stop() error {
	err := h.server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// sselect just to not collide with select keyword
func (h *Http) helloWorld(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("Hello world"))
}

func (h *Http) Name() string {
	return "super http service"
}
