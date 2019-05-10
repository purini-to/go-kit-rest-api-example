package main

import (
	"flag"
	"fmt"
	"github.com/purini-to/go-kit-rest-api-example/middlewares"
	"github.com/purini-to/go-kit-rest-api-example/services"
	"github.com/purini-to/go-kit-rest-api-example/transports"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		port int
		debug bool
	)
	flag.IntVar(&port, "port", 8080,
		`It is a port to listen for HTTP`,
	)
	flag.BoolVar(&debug, "debug", false,
		`Flag to run in the debug environment`,
	)
	flag.Parse()

	run(port, debug)
}

func run(port int, debug bool) {
	var (
		logger *zap.Logger
		err error
	)

	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	var s services.Service
	{
		s = services.NewService()
		s = middlewares.Logging(logger)(s)
	}

	var h http.Handler
	{
		h = transports.MakeHTTPHandler(s, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		port := fmt.Sprintf(":%d", port)
		logger.Info("listen and serve", zap.String("transport", "HTTP"), zap.String("port", port))
		errs <- http.ListenAndServe(port, h)
	}()

	logger.Info("exit", zap.Error(<-errs))
}
