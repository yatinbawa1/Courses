package courseshandler

import (
	"courses/internal/services/auth"
	"courses/internal/services/course"
	"encoding/json"
	"net/http"
)

type GetUserCourses struct {
	courseService *course.CourseService
}

func NewGetUserCoursesHandler(courseService *course.CourseService) *GetUserCourses {
	return &GetUserCourses{courseService}
}

func (g *GetUserCourses) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())

	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized"))
		return
	}

	courses := g.courseService.GetAllCoursesForUser(r.Context(), userID)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(courses)
}
