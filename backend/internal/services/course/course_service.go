package course

import (
	"context"
	"courses/internal/models"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type CourseRepo interface {
	GetCourseDataForUser(ctx context.Context, user_id uuid.UUID) ([]models.Course,error)
	GetTopCourses(ctx context.Context)([]models.Course, error)
}

type CourseService struct {
	courseRepo CourseRepo
	s3Client *s3.Client
}

func NewCourseService(c CourseRepo, s *s3.Client) *CourseService {
	return &CourseService{c,s}
}

func (c *CourseService) GetAllCoursesForUser(ctx context.Context, userID uuid.UUID) []models.Course {
	courses, err := c.courseRepo.GetCourseDataForUser(ctx, userID)
	if err != nil {
		return []models.Course{}
	}
	return courses
}

	

func (c *CourseService) GetTopCourses(ctx context.Context) ([]models.Course) {
	courses, err := c.courseRepo.GetTopCourses(ctx)
	
	if err != nil {
		return []models.Course{}
	}

	return courses
}
