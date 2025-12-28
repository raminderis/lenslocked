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

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		);
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created")

	name := "shahh singh"
	email := "dollr@gmail.com"
	row := conn.QueryRow(context.Background(),
		`INSERT INTO users (name, email)
		VALUES ($1, $2) RETURNING id;`, name, email)
	var id int
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created ", id)

	name = "shahh singh"
	row = conn.QueryRow(context.Background(),
		`SELECT id,name,email FROM users
		WHERE name = $1;`, name)
	var id2 int
	var name2 string
	var email2 string
	err = row.Scan(&id2, &name2, &email2)
	if err == pgx.ErrNoRows {
		fmt.Println("Error no rows error")
		panic(err)
	} else if err != nil {
		fmt.Println("Error other than no rows error")
		panic(err)
	}
	fmt.Printf("Tables created  id=%v name=%s email=%s", id2, name2, email2)

	userId := 2
	for i := 1; i <= 5; i++ {
		amount := i * 342
		desc := fmt.Sprintf("Fake order number # %d", i)
		_, err := conn.Exec(context.Background(),
			`INSERT INTO orders (user_id, amount, description)
			VALUES ($1, $2, $3);`, userId, amount, desc)
		if err != nil {
			fmt.Println("Error inserting")
			panic(err)
		}
	}
	fmt.Println("Added fake order")

	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Description string
	}
	var orders []Order
	rows, err := conn.Query(context.Background(),
		`SELECT id,amount,description FROM orders
		WHERE user_id = $1;`, userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		order.UserID = userId
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	fmt.Println("Orders : ", orders)
}
