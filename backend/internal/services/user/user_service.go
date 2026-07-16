package user

import (
	"context"
	"courses/internal/models"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/crypto/bcrypt"
)

var (
	hasLower = regexp.MustCompile(`[a-z]`)
	hasUpper = regexp.MustCompile(`[A-Z]`)
	hasNum   = regexp.MustCompile(`\d`)
)

var (
	ErrUnsecurePassword       = errors.New("Password Not Secure Enough")
	ErrPasswordHashGeneration = errors.New("Unable to create a Secure Hash for Password")
	ErrWrongPassword          = errors.New("Wrong Password")
	ErrUserDoesNotExist       = errors.New("User Does Not Exist")
)

type UserDataRepo interface {
	Add(ctx context.Context, user *models.UserAuthCreds) error
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	GetPasswordForEmail(ctx context.Context, email string) ([]byte, error)
	GetUserData(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	CheckIfUserIDExists(ctx context.Context, userID string) (bool, error)
}

type UserService struct {
	userRepo UserDataRepo
	s3Client *s3.Client
}

func NewUserService(repo UserDataRepo, s3Client *s3.Client) *UserService {
	return &UserService{userRepo: repo, s3Client: s3Client}
}

func (s *UserService) ValidateCredentials(ctx context.Context, email string, password string) error {
	pass, err := s.userRepo.GetPasswordForEmail(ctx, email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return ErrWrongPassword
	}

	return nil
}

func (s *UserService) GetUserData(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserData(ctx, email)
}

func (s *UserService) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	return s.userRepo.CheckIfEmailExists(ctx, email)
}

func (s *UserService) SignUpUser(ctx context.Context, email string, password string) error {
	email = strings.ToLower(email)

	if len(password) < 8 || !hasLower.MatchString(password) || !hasUpper.MatchString(password) || !hasNum.MatchString(password) {
		return ErrUnsecurePassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return ErrPasswordHashGeneration
	}

	user := &models.UserAuthCreds{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepo.Add(ctx, user)
}

func (u *UserService) CreatePresignedUploadURLForProfilePhoto(ctx context.Context, userID string) (*v4.PresignedHTTPRequest, error) {
	if ex, _ := u.userRepo.CheckIfUserIDExists(ctx, userID); !ex {
		return nil, fmt.Errorf("Unable to find User!")
	}

	bucketName := "courses-content-portfolio-go-next"
	presignClient := s3.NewPresignClient(u.s3Client)
	key := fmt.Sprintf("users/%s", userID)

	resp, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String("image/png"),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}
