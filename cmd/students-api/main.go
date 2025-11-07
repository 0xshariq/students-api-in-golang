package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xshariq/students-api-in-golang/internal/config"
	"github.com/0xshariq/students-api-in-golang/internal/http/handlers/student"
)

func main() {
	// load config
	cfg := config.MustConfig()

	fmt.Printf("Configuration loaded:\n")
	fmt.Printf("  Environment: %s\n", cfg.Env)
	fmt.Printf("  HTTP Server Host: %s\n", cfg.HttpServer.Host)
	fmt.Printf("  HTTP Server Port: %d\n", cfg.HttpServer.Port)

	// database setup
	databaseUrl := cfg.StoragePath
	_ = databaseUrl
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", student.Home())
	router.HandleFunc("POST /api/students", student.NewStudent())
	// setup server
	serverAddr := fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port)
	server := http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	fmt.Printf("Starting server on %s...\n", server.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {

		error := server.ListenAndServe()
		if error != nil {
			log.Fatal("Error starting server:", error)
		}
	}()

	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if error := server.Shutdown(ctx); error != nil {
		slog.Error("Failed to shutdown server", slog.String("error", error.Error()))
	}

	slog.Info("Server Shutdown successfully")
}
