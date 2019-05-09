package main

import (
	"fmt"
	"github.com/purini-to/go-kit-rest-api-example"
	"github.com/purini-to/go-kit-rest-api-example/middlewares"
	"github.com/purini-to/go-kit-rest-api-example/services"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	port int
}

var opt = &Options{}

func main() {
	var cmd = &cobra.Command{
		Use:   "api",
		Short: "REST APIサーバーを起動します",
		Long:  `REST APIサーバーを起動します`,
		Run:   run,
	}

	cmd.Flags().IntVarP(&opt.port, "port", "p", 8080, "HTTPリッスンポート")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(_ *cobra.Command, _ []string) {
	logger, err := zap.NewDevelopment()
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
		h = go_kit_rest_api_example.MakeHTTPHandler(s, logger)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		port := fmt.Sprintf(":%d", opt.port)
		logger.Info("listen and serve", zap.String("transport", "HTTP"), zap.String("port", port))
		errs <- http.ListenAndServe(port, h)
	}()

	logger.Info("exit", zap.Error(<-errs))
}
