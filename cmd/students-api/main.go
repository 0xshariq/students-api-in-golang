package main

import (
	"fmt"
	"github.com/0xshariq/students-api-in-golang/internal/config"
	"log"
	"net/http"
)

func main() {
	// load config
	cfg := config.MustConfig()

	fmt.Printf("Configuration loaded:\n")
	fmt.Printf("  Environment: %s\n", cfg.Env)
	fmt.Printf("  Storage Path: %s\n", cfg.StoragePath)
	fmt.Printf("  HTTP Server Host: %s\n", cfg.HttpServer.Host)
	fmt.Printf("  HTTP Server Port: %d\n", cfg.HttpServer.Port)

	// database setup
	databaseUrl := cfg.StoragePath
	_ = databaseUrl
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Welcome to students api in golang"))
	})
	// setup server
	serverAddr := fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port)
	server := http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	fmt.Printf("Starting server on %s...\n", server.Addr)
	error := server.ListenAndServe()
	if error != nil {
		log.Fatal("Error starting server:", error)
	}
}
