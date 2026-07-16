package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CoursesRepo struct {
	db *pgxpool.Pool
}
