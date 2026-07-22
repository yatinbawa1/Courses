package handlers

import (
	accounthandler "courses/internal/handlers/account"
	authhandler "courses/internal/handlers/auth"
	courseshandler "courses/internal/handlers/courses"
	"courses/internal/middleware"
	"courses/internal/services/auth"
	"courses/internal/services/course"
	"courses/internal/services/mailer"
	"courses/internal/services/user"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

func RegisterRoutes(
	fileServer http.Handler,
	logger *log.Logger,
	authService *auth.AuthService,
	userService *user.UserService,
	mailer mailer.MailSender,
	frontendFS fs.FS,
	courseService *course.CourseService,
) *http.ServeMux {

	mux := http.NewServeMux()

	// ------------------------
	// Public API
	// ------------------------

	mux.Handle("GET /api/auth/refresh",
		authhandler.NewRefreshHandler(logger, authService))

	mux.Handle("POST /api/auth/login",
		authhandler.NewLoginHandler(logger, authService))

	mux.Handle("POST /api/auth/send-otp",
		authhandler.NewSendOTPHandler(logger, authService, mailer))

	mux.Handle("POST /api/auth/send-otp/verify",
		authhandler.NewVerifyOTP(logger, authService))

	mux.Handle("GET /api/courses/get-top-courses/", courseshandler.NewGetTopCourses(courseService))
	// ------------------------
	// Protected API
	// ------------------------

	mux.Handle("POST /api/auth/logout",
		middleware.CheckAuth(authhandler.NewLogoutHandler(logger, authService)))

	mux.Handle("POST /api/account/update-user",
		middleware.CheckAuth(accounthandler.NewUpdateUserHandler(logger, userService)))

	mux.Handle("GET /api/account/get-profile-image-url",
		middleware.CheckAuth(accounthandler.NewGetProfilePhotoUploadUrl(userService)))

	mux.Handle("GET /api/courses/get-user-courses",
		middleware.CheckAuth(courseshandler.NewGetUserCoursesHandler(courseService)))

	
	// ------------------------
	// Frontend
	// ------------------------

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// API routes should never reach the SPA fallback.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// If the requested file exists in the embedded frontend,
		// let the embedded file server serve it.
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		if _, err := frontendFS.Open(path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Otherwise serve index.html so SvelteKit can handle routing.
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
