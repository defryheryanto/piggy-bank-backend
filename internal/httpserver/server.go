package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/defryheryanto/piggy-bank-backend/internal/app"
)

type ApplicationServer struct {
	port        string
	address     string
	application *app.Application
}

func NewApplicationServer(port, address string, application *app.Application) *ApplicationServer {
	return &ApplicationServer{port, address, application}
}

func (s *ApplicationServer) ServeApp() {
	routes := s.CompileRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.address, s.port),
		Handler: routes,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Listening to %s:%s", s.address, s.port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error %s", err)
		}
	}()
	<-quit
}
