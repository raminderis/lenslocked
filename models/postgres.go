package models

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

func (cfg PostgresConfig) Migrate(migrationDir string) error {
	connString := cfg.String()
	sqlDB, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("DB Migration error: %w", err)
	}
	defer sqlDB.Close()
	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("DB Migration error : %w", err)
	}
	err = goose.Up(sqlDB, migrationDir)
	if err != nil {
		return fmt.Errorf("DB Migration error : %w", err)
	}
	return nil
}

func (cfg PostgresConfig) MigrateFS(migrationFS fs.FS) error {
	connString := cfg.String()
	entries, err := fs.ReadDir(migrationFS, ".")
	if err != nil {
		return fmt.Errorf("ReadDir error: %w", err)
	}

	for _, e := range entries {
		fmt.Println(" -", e.Name())
	}

	if err != nil {
		return fmt.Errorf("DB MigrateFS error 1: %w", err)
	}
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("DB MigrateFS error 2: %w", err)
	}
	defer db.Close()

	// Create a provider that uses your embedded FS
	p, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrationFS,
	)
	if err != nil {
		return fmt.Errorf("DB MigrateFS error 3: %w", err)
	}

	_, err = p.Up(context.Background())
	if err != nil {
		return fmt.Errorf("DB MigrateFS error 4: %w", err)
	}

	return nil
}
