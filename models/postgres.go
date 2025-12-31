package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
	}
}

func Open(dbConfig PostgresConfig) (*pgx.Conn, error) {
	url := dbConfig.String()

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("DB Open: %w", err)
	}
	return conn, nil
}
