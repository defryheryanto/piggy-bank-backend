package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/defryheryanto/piggy-bank-backend/config"
	_ "github.com/defryheryanto/piggy-bank-backend/config/env"
	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver"
)

func main() {
	application := buildApp()
	applicationServer := httpserver.NewApplicationServer(config.ListenPort(), config.ListenAddress(), application)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go applicationServer.ServeApp()

	<-quit
}
