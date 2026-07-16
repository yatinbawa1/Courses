package repository

import (
	"context"
	"courses/internal/models"
	"courses/internal/services/user"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
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

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	var args []any
	var setClauses []string
	argCounter := 1

	if user.Username != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *user.Username)
		argCounter++
	}

	if user.ProfilePhotoURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("profile_photo_url = $%d", argCounter))
		args = append(args, *user.ProfilePhotoURL)
		argCounter++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, user.User_id)
	whereClause := fmt.Sprintf("WHERE user_id = $%d", argCounter)

	query := fmt.Sprintf(
		`UPDATE "User" SET %s %s;`,
		strings.Join(setClauses, ", "),
		whereClause,
	)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute dynamic user update: %w", err)
	}

	// Optional: Check if the UUID actually existed in the database
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", user.User_id)
	}

	return nil
}

func (r *UserRepo) GetUserData(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	lookupEmail := strings.ToLower(email)

	query := `SELECT user_id, name, profile_photo_url FROM "User" WHERE email = $1;`

	err := r.db.QueryRow(ctx, query, lookupEmail).Scan(
		&user.User_id,
		&user.Username,
		&user.ProfilePhotoURL,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed To Find User! Internal Error %w", err)
	}

	user.Email = lookupEmail

	return user, nil
}

func (r *UserRepo) CheckIfUserIDExists(ctx context.Context, userID string) (bool, error) {
	userID = strings.ToLower(userID)

	query := `SELECT EXISTS (
    SELECT 1 
    FROM "User" 
    WHERE user_id = $1
	);`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("Failed To Check UserID! Internal Error %w", err)
	}

	return exists, nil
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
			return []byte{}, user.ErrUserDoesNotExist
		} else {
			// A Internal Log for Error Can Be Added Here
			return []byte{}, err
		}
	}

	return []byte(storedHash), nil
}

func (r *UserRepo) Add(ctx context.Context, usercreds *models.UserAuthCreds) error {
	usercreds.Email = strings.ToLower(usercreds.Email)

	userId := uuid.New()
	query := `INSERT INTO "User" (user_id,hashed_password, email) values ($1,$2,$3)`

	_, err := r.db.Exec(ctx, query, userId, usercreds.Password, usercreds.Email)

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
