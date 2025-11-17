package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/cchrisris/go-course/hw3/internal/adapters/http"
	frepo "github.com/cchrisris/go-course/hw3/internal/adapters/storage/file"
	mrepo "github.com/cchrisris/go-course/hw3/internal/adapters/storage/memory"
	"github.com/cchrisris/go-course/hw3/internal/usecase"
)

func main() {
	addr := getenv("LISTEN_ADDR", ":8080")
	storage := getenv("STORAGE", "memory")
	filePath := getenv("FILE_PATH", "/data/balances.json")

	var service usecase.Service
	switch storage {
	case "memory":
		service = usecase.NewService(mrepo.NewRepository())
	case "file":
		fileRepository, err := frepo.NewRepository(filePath)
		if err != nil {
			log.Fatalf("file repo init: %v", err)
		}
		service = usecase.NewService(fileRepository)
	default:
		log.Fatalf("unknown STORAGE value: %q (expected 'memory' or 'file')", storage)
	}

	handler := httpapi.NewHandler(service)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (storage=%s)", addr, storage)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}


