package handlers

import (
	"courses/internal/auth"
	accounthandler "courses/internal/handlers/account"
	authhandler "courses/internal/handlers/auth"
	"courses/internal/mailer"
	"courses/internal/middleware"
	"io"
	"io/fs"
	"log"
	"net/http"
)

func RegisterRoutes(fileServer http.Handler, logger *log.Logger, authService *auth.AuthService, mailer mailer.MailSender, frontendFS fs.FS) *http.ServeMux {
	mux := http.NewServeMux()	

	// Unsecure Links
	mux.Handle("GET /api/auth/refresh", authhandler.NewRefreshHandler(logger, authService))
	mux.Handle("POST /api/auth/login", authhandler.NewLoginHandler(logger, authService))
	
	mux.Handle("POST /api/auth/send-otp", authhandler.NewSendOTPHandler(logger, authService, mailer))
	mux.Handle("POST /api/auth/send-otp/verify", authhandler.NewVerifyOTP(logger, authService))
	
	// Secure Links
	mux.Handle("POST /api/auth/logout/{user_email}", middleware.CheckAuth(authhandler.NewLogoutHandler(logger, authService)))
	mux.Handle("POST /api/account/update-user", middleware.CheckAuth(accounthandler.NewUpdateUserHandler(logger, authService)))
	

	// Frontend - SPA fallback using embedded FS
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := frontendFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
	})
	return mux
}
