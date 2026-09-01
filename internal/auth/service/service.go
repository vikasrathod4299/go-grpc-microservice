package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	authRepo "github.com/vikasrathod4299/microservice/internal/auth/repository/db"
	tokenauth "github.com/vikasrathod4299/microservice/pkg/auth"
)

type Role string

const (
	RoleRider  Role = "rider"
	RoleDriver Role = "driver"
)

var (
	ErrRepositoryRequired  = errors.New("Auth repository is required")
	ErrJWTSecretRequired   = errors.New("JWT secret must contain at least 32 bytes")
	ErrNameRequired        = errors.New("name is required")
	ErrPhoneRequired       = errors.New("phone is required")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidRole         = errors.New("invalid user role")
	ErrUserIDRequired      = errors.New("user ID is required")
	ErrUserAlreadyExists   = errors.New("a user with that email or phone already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidDatabaseUUID = errors.New("database returned an invalid user ID")
)

type User struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
type RegisterInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Role     Role
}

type Login struct {
	Email    string
	Password string
}

type AuthResult struct {
	User        *User
	AccessToken string
	ExpiresAt   time.Time
}

type AuthSerivce struct {
	repo      authRepo.Querier
	jwtString string
}

func NewAuthService(repo authRepo.Querier, jwtString string) (*AuthSerivce, error) {
	if repo == nil {
		return nil, ErrRepositoryRequired
	}
	if len([]byte(jwtString)) < 32 {
		return nil, ErrJWTSecretRequired
	}

	return &AuthSerivce{
		repo:      repo,
		jwtString: jwtString,
	}, nil
}

func (s *AuthSerivce) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Role = Role(
		strings.ToLower(strings.TrimSpace(string(Role(input.Role)))),
	)
	if input.Name == "" {
		return nil, ErrNameRequired
	}
	if input.Phone == "" {
		return nil, ErrPhoneRequired
	}
	if !validEmail(input.Email) {
		return nil, ErrInvalidEmail
	}
	if !validRole(input.Role) {
		return nil, ErrInvalidPassword
	}

	hashedPassword, err := tokenauth.HashPassword(input.Password)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	createdUser, err := s.repo.CreateUser(ctx, authRepo.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: hashedPassword,
		Role:         authRepo.UserRole(input.Role),
	})
	if err != nil {
		if uniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	user, err := userFromCreateRow(createdUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.issueToken(user)
}

func (s *AuthSerivce) Login(ctx context.Context, input Login) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if !validEmail(email) || input.Password == "" {
		return nil, ErrInvalidEmail
	}

	databaseUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := tokenauth.ComparePassword(input.Password, databaseUser.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := userFromDatabase(databaseUser)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return s.issueToken(user)
}

func (s *AuthSerivce) GetUser(ctx context.Context, userID string) (*User, error) {
	uuid := strings.TrimSpace(userID)
	if uuid == "" {
		return nil, ErrUserIDRequired
	}

	databaseID, err := parseUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}
	databaseUser, err := s.repo.GetUserById(ctx, databaseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by ID: %w", err)
	}
	user, err := userFromGetByIDRow(&databaseUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthSerivce) issueToken(user *User) (*AuthResult, error) {
	token, expiresAt, err := tokenauth.GenerateToken(user.ID, string(user.Role), s.jwtString)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}
	return &AuthResult{
		User:        user,
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}, nil
}

func validEmail(email string) bool {
	if email == "" {
		return false
	}
	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, address.Address)
	return matched && err == nil
}

func validRole(role Role) bool {
	switch role {
	case RoleRider, RoleDriver:
		return true
	default:
		return false
	}
}

func formatUUID(id pgtype.UUID) (string, error) {
	if !id.Valid {
		return "", ErrInvalidDatabaseUUID
	}
	value := id.Bytes
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		value[0:4],
		value[4:6],
		value[6:8],
		value[8:10],
		value[10:16],
	), nil
}

func parseUUID(value string) (pgtype.UUID, error) {
	var id pgtype.UUID

	if err := id.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}

	if !id.Valid {
		return pgtype.UUID{}, ErrInvalidDatabaseUUID
	}

	return id, nil
}

func userFromGetByIDRow(user *authRepo.GetUserByIdRow) (*User, error) {
	id, err := formatUUID(user.ID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      Role(user.Role),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func userFromCreateRow(user authRepo.CreateUserRow) (*User, error) {
	id, err := formatUUID(user.ID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      Role(user.Role),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func userFromDatabase(user authRepo.User) (*User, error) {
	id, err := formatUUID(user.ID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      Role(user.Role),
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func uniqueViolation(err error) bool {
	var postgresError *pgconn.PgError

	return errors.As(err, &postgresError) && postgresError.Code == "23505"
}
