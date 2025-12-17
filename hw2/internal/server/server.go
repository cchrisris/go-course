package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	server *http.Server
}

func NewServer(port string) *Server {
	server := Server{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", server.VersionHandler)
	mux.HandleFunc("POST /decode", server.DecodeHandler)
	mux.HandleFunc("GET /hard-op", server.HardOpHandler)

	server.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}
	return &server
}

func (server *Server) Start() {
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := server.server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		fmt.Println(err)
	}

	<-idleConnectionsClosed
}
