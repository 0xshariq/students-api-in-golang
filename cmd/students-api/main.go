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
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize sqlite:", err)
	}
	slog.Info("Storage Initiallized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", student.Home())
	router.HandleFunc("GET /api/students/{id}", student.GetStudent(storage))
	router.HandleFunc("GET /api/students", student.GetStudents(storage))
	router.HandleFunc("POST /api/students/create", student.NewStudent(storage))
	router.HandleFunc("DELETE /api/students/delete/{id}", student.DeleteStudent(storage))
	router.HandleFunc("PUT /api/students/update/{id}", student.UpdateStudent(storage))

	// setup server
	serverAddr := fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port)
	server := http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	slog.Info("Server Started", slog.String("host", cfg.Host), slog.Int("port", cfg.Port))
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
