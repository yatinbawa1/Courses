package main

import (
	"context"
	"courses/internal/auth"
	"courses/internal/config"
	"courses/internal/database"
	"courses/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv(l *log.Logger) {
	err := godotenv.Load()

	if err != nil {
		l.Printf("Unable to locate ENV File")
	}
}

func main() {

	logger := log.New(os.Stdout, "\nCourses-", log.Default().Flags())

	// Initialize Everything
	LoadEnv(logger)
	config.Init()
	_, err := database.ConnectDataBase()

	if err != nil {
		logger.Fatalf("%s", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", handlers.NewHomeHandler(logger))
	mux.Handle("POST /", handlers.NewLoginHandler(logger, auth.VerifyToken))

	server := &http.Server{
		Handler:      mux,
		Addr:         config.Port,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Printf("Server Starting on Port %s\n", config.Port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatalf("Error Starting the Server")
		}
	}()

	SigChan := make(chan os.Signal, 1)
	signal.Notify(SigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-SigChan
	logger.Fatalf("Received %s Signal! Initiating Exit Routine", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)
}
