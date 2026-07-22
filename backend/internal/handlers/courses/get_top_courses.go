package courseshandler

import (
	"courses/internal/services/course"
	"encoding/json"
	"net/http"
)

type GetTopCourses struct {
	courseService *course.CourseService
}

func NewGetTopCourses(courseService *course.CourseService) *GetTopCourses {
	return &GetTopCourses{courseService}
}

func (g *GetTopCourses) ServeHTTP(rw http.ResponseWriter, r *http.Request) {	
	courses := g.courseService.GetTopCourses(r.Context())

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(courses)
}
