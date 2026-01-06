package models

import (
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    string
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB_CONN       *pgx.Conn
	BytesPerToken int
	Duration      time.Duration
}

func (resetService *PasswordResetService) Create(email string) (*PasswordReset, error) {

	return nil, nil
}

func (resetService *PasswordResetService) Consume(token string) (*User, error) {

	return nil, nil
}
