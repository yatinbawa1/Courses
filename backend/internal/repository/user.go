package repository

import (
	"context"
	"courses/internal/auth"
	"courses/internal/models"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
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

func (r *UserRepo) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	email = strings.ToLower(email)

	query := `SELECT EXISTS (
    SELECT 1 
    FROM "User" 
    WHERE email = $1
	);`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("Failed To Check Email! Internal Error %w", err)
	}

	return exists, nil
}

func (r *UserRepo) GetPasswordForEmail(ctx context.Context, email string) ([]byte, error) {
	email = strings.ToLower(email)
	query := `SELECT hashed_password FROM "User" WHERE email = $1;`

	var storedHash string
	err := r.db.QueryRow(ctx, query, email).Scan(&storedHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []byte{}, auth.ErrUserDoesNotExist
		} else {
			// A Internal Log for Error Can Be Added Here
			return []byte{}, err
		}
	}

	return []byte(storedHash), nil
}

func (r *UserRepo) Add(ctx context.Context, user *models.User) error {
	user.Email = strings.ToLower(user.Email)

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
