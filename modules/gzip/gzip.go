package gzip

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

type Gzip struct {
}

func (gz *Gzip) Init() error {
	return nil
}

func (gz *Gzip) Middleware(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gziphandler.GzipHandler(f).ServeHTTP(w, r)
	}
}

func (gz *Gzip) Name() string {
	return "super Gzip middleware"
}
