package main

import (
	"os"
	"os/signal"
	"syscall"

	"roadrunner_http_test/modules/db"
	"roadrunner_http_test/modules/gzip"
	"roadrunner_http_test/modules/headers"
	"roadrunner_http_test/modules/http"
	"roadrunner_http_test/modules/logger"
	endure "roadrunner_http_test/pkg/container"
)

func main() {
	// no external logger
	container, err := endure.NewContainer(nil, endure.Visualize(endure.StdOut, ""), endure.SetLogLevel(endure.DebugLevel))
	if err != nil {
		panic(err)
	}

	err = container.RegisterAll(
		&http.Http{},
		&db.DB{},
		&logger.Logger{},
		&gzip.Gzip{},
		&headers.Headers{},
	)

	if err != nil {
		panic(err)
	}

	err = container.Init()
	if err != nil {
		panic(err)
	}

	errCh, err := container.Serve()
	if err != nil {
		panic(err)
	}

	// stop by CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGINT)

	for {
		select {
		case e := <-errCh:
			println(e.Error.Error())
			er := container.Stop()
			if er != nil {
				panic(er)
			}
			return
		case <-c:
			er := container.Stop()
			if er != nil {
				panic(er)
			}
			return
		}
	}
}
