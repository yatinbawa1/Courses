package repository

import (
	"context"
	"courses/internal/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailAlreadyExists = errors.New("This Email Already Exists In Database! Use Another")
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Add(ctx context.Context, user *models.User) error {
	query := `INSERT INTO "User" (user_id,hashed_password, email) values ($1,$2,$3)`

	_, err := r.db.Exec(ctx, query, user.User_id, user.HashedPassword, user.Email)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExists
			}
		}

		return fmt.Errorf("Unable to Add User to Database! %w", err)
	}

	return nil
}
