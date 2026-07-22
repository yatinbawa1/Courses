package main

import (
	"context"
	backend "courses"
	"courses/internal/config"
	"courses/internal/database"
	"courses/internal/handlers"
	"courses/internal/repository"
	"courses/internal/services/auth"
	"courses/internal/services/course"
	"courses/internal/services/mailer"
	"courses/internal/services/user"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var frontendFS embed.FS

func main() {
	logger := log.New(os.Stdout, "\nCourses-", log.Default().Flags())

	// Initialize ENV
	// And Database and Redis
	config.Init(logger)
	pool, err := database.ConnectDataBase()
	if err != nil {
		logger.Fatalf("Unable to connect with Database %s", err)
	}
	defer pool.Close()

	// AWS Client
	s3Client := database.GetS3Client(logger)

	// This is for local redis client
	redisClient, err := database.NewRedisClient(config.REDIS_ADDR, "", 0)
	// redisClient, err := database.NewRedisOnlineClient()
	if err != nil {
		logger.Fatalf("Unable to Connect With Redis %s", err)
	}

	defer redisClient.Close()

	// Initialize Frontend
	strippedFS, err := fs.Sub(backend.FrontendFS, "dist")
	if err != nil {
		log.Fatal("Failed to read embedded folder: ", err)
	}

	fileServer := http.FileServer(http.FS(strippedFS))

	// Initialize Core Services
	resendMailer := mailer.NewResendMailer()
	otpRepo := repository.NewRedisOTPRepo(redisClient)
	userRepo := repository.NewUserRepo(pool)
	courseRepo := repository.NewCourseRepo(pool, userRepo)

	refreshRepo := repository.NewRedisRefreshTokenRepo(redisClient)
	userService := user.NewUserService(userRepo, s3Client)
	authService := auth.NewAuthService(otpRepo, refreshRepo, userService)
	courseService := course.NewCourseService(courseRepo, s3Client)
	
	// Set up a router
	mux := handlers.RegisterRoutes(fileServer, logger, authService,
		 userService, resendMailer, strippedFS, courseService)

	server := &http.Server{
		Handler:      mux,
		Addr:         ":" + config.Port,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("Server Starting on Port %s\n", config.Port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatalf("Error Starting the Server: %s", err)
		}
	}()

	SigChan := make(chan os.Signal, 1)
	signal.Notify(SigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-SigChan
	logger.Printf("Received %s Signal! Initiating Exit Routine", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Shutdown(timeoutContext)
}
