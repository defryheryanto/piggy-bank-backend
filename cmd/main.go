package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/defryheryanto/piggy-bank-backend/internal/httpserver"
)

func main() {
	applicationServer := httpserver.NewApplicationServer("8080", "localhost")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go applicationServer.ServeApp()

	<-quit
}
