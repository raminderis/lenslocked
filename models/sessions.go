package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/raminderis/lenslocked/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB_CONN       *pgx.Conn
	BytesPerToken int
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("Session Service : %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB_CONN.QueryRow(context.Background(),
		`INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2 RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create session : %w", err)
	}
	fmt.Println("Session inserted ", session.ID)
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenhash := ss.hash(token)
	user := User{}
	row := ss.DB_CONN.QueryRow(context.Background(),
		`SELECT users.id, users.email, users.password_hash
		 FROM sessions 
		 JOIN users on users.id = sessions.user_id
		WHERE sessions.token_hash = $1`, tokenhash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("session not found : %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenhash := ss.hash(token)
	_, err := ss.DB_CONN.Exec(context.Background(),
		`DELETE FROM sessions WHERE token_hash = $1;`, tokenhash)
	if err != nil {
		return fmt.Errorf("Session Delete: %w", err)
	}
	return nil
}
