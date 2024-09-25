package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func GetConnection() *pgx.Conn {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	return conn
}

func DoRequest(conn *pgx.Conn, query string, args ...any) pgx.Rows {
	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	}

	return rows
}
