package repository

import (
	"context"
	"courses/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoursesRepo struct {
	db *pgxpool.Pool
	userRepo *UserRepo
}

func NewCourseRepo(db *pgxpool.Pool, userRepo *UserRepo) *CoursesRepo {
	return &CoursesRepo{db, userRepo}
}


func (c *CoursesRepo) GetCourseDataForUser(ctx context.Context,user_id uuid.UUID) ([]models.Course, error) {
	query := `
		SELECT 
			c.course_id, 
			c.course_name, 
			c.course_description, 
			c.creation_date,
			c.course_thumbnail
		FROM 
			User_Owned_Courses uoc
		INNER JOIN 
			Course c ON uoc.course_id = c.course_id
		WHERE 
			uoc.user_id = $1;
	`
	rows, err := c.db.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	courses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Course])
	if err != nil {
		return nil, err
	}

	return courses, nil
}
