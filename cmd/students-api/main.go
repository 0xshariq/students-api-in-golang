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

	"github.com/0xshariq/students-api-in-golang/pkg/config"
	"github.com/0xshariq/students-api-in-golang/pkg/http/handlers/student"
	"github.com/0xshariq/students-api-in-golang/pkg/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustConfig()

	fmt.Printf("Configuration loaded:\n")
	fmt.Printf("  Environment: %s\n", cfg.Env)
	fmt.Printf("  HTTP Server Host: %s\n", cfg.HttpServer.Host)
	fmt.Printf("  HTTP Server Port: %d\n", cfg.HttpServer.Port)

	// database setup
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize sqlite:", err)
	}
	_ = db
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", student.Home())
	router.HandleFunc("POST /api/students/create", student.NewStudent())
	// router.HandleFunc("DELETE /api/students/create", student.DeleteStudent())
	// router.HandleFunc("PUT /api/students/create", student.UpdateStudent())
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
