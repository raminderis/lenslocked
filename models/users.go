package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken = errors.New("models: email address is already taken")
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB_CONN *pgx.Conn
}

func (userSvc *UserService) Create(email, passwordPlainText string) (*User, error) {
	email = strings.ToLower(email)
	passwordHashedBytes, err := bcrypt.GenerateFromPassword([]byte(passwordPlainText), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user : %w", err)
	}
	// fmt.Println(string(passwordHashedBytes))
	passwordHashedString := string(passwordHashedBytes)
	user := User{
		ID:           1,
		Email:        email,
		PasswordHash: passwordHashedString,
	}
	row := userSvc.DB_CONN.QueryRow(context.Background(),
		`INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id;`, email, passwordHashedString)
	var id int
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user : %w", err)
	}
	fmt.Println("User inserted ", id)
	return &user, nil
}

func (userSvc *UserService) Authenticate(email, passwordPlainText string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := userSvc.DB_CONN.QueryRow(context.Background(),
		`SELECT id, password_hash FROM users 
		WHERE email = $1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate user : %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordPlainText))
	if err != nil {
		return nil, fmt.Errorf("authenitcate user : %w", err)
	}
	fmt.Println("User authenticated ", email)
	return &user, nil
}

func (userSvc *UserService) UpdatePassword(userID int, password string) error {
	passwordHashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password : %w", err)
	}
	passwordHashedString := string(passwordHashedBytes)
	_, err = userSvc.DB_CONN.Exec(context.Background(), `
		UPDATE users
		SET password_hash = $2
		WHERE id = $1;`, userID, passwordHashedString)
	if err != nil {
		return fmt.Errorf("Update Password: %w", err)
	}
	return nil
}
