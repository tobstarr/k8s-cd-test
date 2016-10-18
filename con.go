package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Con interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func connect(dbURL string) (*sql.DB, error) {
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
		if dbURL == "" {
			return nil, fmt.Errorf("DATABASE_URL must be present")
		}
	}
	con, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	err = con.Ping()
	if err != nil {
		return nil, err
	}
	return con, nil
}
