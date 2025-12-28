package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// export DATABASE_URL = "postgres://user:pass@host:port/db"

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

func main() {
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// conn, err := pgx.Connect(context.Background(), "postgres://baloo:junglebook@127.0.0.1:5432/lenslocked")
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
	}
	url := cfg.String()

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
	err = conn.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Ping Successful!")
}
