package handlers

import (
	"courses/internal/auth"
	accounthandler "courses/internal/handlers/account"
	authhandler "courses/internal/handlers/auth"
	"courses/internal/mailer"
	"courses/internal/middleware"
	"log"
	"net/http"
	"path/filepath"
)

func RegisterRoutes(fileServer http.Handler, logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender) *http.ServeMux {
	mux := http.NewServeMux()	

	// Unsecure Links
	mux.Handle("GET /api/auth/refresh", authhandler.NewRefreshHandler(logger, authService))
	mux.Handle("POST /api/auth/login", authhandler.NewLoginHandler(logger, authService))
	
	mux.Handle("POST /api/auth/send-otp", authhandler.NewSendOTPHandler(logger, authService, mailer))
	mux.Handle("POST /api/auth/send-otp/verify", authhandler.NewVerifyOTP(logger, authService))
	
	// Secure Links
	mux.Handle("POST /api/auth/logout/{user_email}", middleware.CheckAuth(authhandler.NewLogoutHandler(logger, authService)))
	mux.Handle("POST /api/account/update-user", middleware.CheckAuth(accounthandler.NewUpdateUserHandler(logger, authService)))
	

	// Frontend - Unsecure
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if filepath.Ext(r.URL.Path) != "" {
            fileServer.ServeHTTP(w, r)
            return
        }

        http.ServeFile(w, r, "./ui/dist/index.html") 
    })
	return mux
}
