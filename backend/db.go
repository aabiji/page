package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
)

type User struct {
	Id       string
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DB struct {
	conns   *pgxpool.Pool
	context context.Context
}

// Initialize database instance by creating a series of tables if they weren't
// already created.
func NewDatabase() DB {
	db := DB{context: context.Background()}

	databaseUrl := "postgres://aabiji:0000@localhost:5432/Page"
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		panic(err)
	}

	db.conns, err = pgxpool.NewWithConfig(db.context, config)
	if err != nil {
		panic(err)
	}

	// Create database tables if they don't already exist
	createUsers := `
    CREATE TABLE IF NOT EXISTS Users (
        UserId serial PRIMARY KEY,
        Email text NOT NULL,
        Password text NOT NULL
    );`
	if _, err := db.conns.Exec(db.context, createUsers); err != nil {
		panic(err)
	}

	createBooks := `
    CREATE TABLE IF NOT EXISTS Books (
        BookId serial PRIMARY KEY,
        CoverImagePath text NOT NULL,
        Files text[] NOT NULL,
        TableOfContents jsonb[] NOT NULL,
        Info jsonb NOT NULL
    );`
	if _, err := db.conns.Exec(db.context, createBooks); err != nil {
		panic(err)
	}

	createUserBooks := `
    CREATE TABLE IF NOT EXISTS UserBooks (
        UserId integer NOT NULL,
        BookId integer NOT NULL,
        CurrentPage integer NOT NULL,
        ScrollOffsets integer[] NOT NULL
    );`
	if _, err := db.conns.Exec(db.context, createUserBooks); err != nil {
		panic(err)
	}

	return db
}

// Execute sql query on database.
func (db *DB) Exec(sql string, params ...any) error {
	_, err := db.conns.Exec(db.context, sql, params...)
	return err
}

// Read value from sql database.
// sqlParams is a slice of all the input parameters for the query.
// readParams is a slice of pointers for receiving values of the query.
// Returns a slice containing all the query results.
func (db *DB) Read(sql string, sqlParams []any, readParams []any) ([]any, error) {
	var results []any

	rows, err := db.conns.Query(db.context, sql, sqlParams...)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var values []any
		if err := rows.Scan(readParams...); err != nil {
			return results, err
		}
		for i := 0; i < len(readParams); i++ {
			val := reflect.ValueOf(readParams[i])
			values = append(values, reflect.Indirect(val))
		}
		results = append(results, values)
	}

	if err := rows.Err(); err != nil {
		return results, err
	}

	if len(results) == 0 {
		return results, errors.New("Entries not found.")
	}

	return results, nil
}
