package courses

import (
	"courses/internal/models"

	"github.com/google/uuid"
)

type CourseService struct {
	courseRepo CourseRepo
}

type CourseRepo interface {
	GetCoursesForUser(user_id uuid.UUID) []models.Course
	GetCourseData(course_id string) *models.Course
	GetCourseCounter(course_id string) []models.CourseChapter
	GetChapterContent(chapter_id string) models.AssetType
}
