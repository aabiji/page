package main

// TODO: cleanup code

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQL Queries:
const createBooksTableSQL = `
CREATE TABLE IF NOT EXISTS Books (
    BookId serial PRIMARY KEY,
    CoverImagePath text NOT NULL,
    Files text[] NOT NULL,
    TableOfContents jsonb[] NOT NULL,
    Info jsonb NOT NULL
);`

const createUsersTableSQL = `
CREATE TABLE IF NOT EXISTS Users (
    UserId serial PRIMARY KEY,
    Email text NOT NULL,
    Password text NOT NULL
);`

const createUserBooksTableSQL = `
CREATE TABLE IF NOT EXISTS UserBooks (
    UserId integer NOT NULL,
    BookId integer NOT NULL,
    CurrentPage integer NOT NULL,
    ScrollOffsets integer[] NOT NULL
);`

const createUserSQL = `
INSERT INTO Users (Email, Password) VALUES ($1, $2);`

const readUserSQL = `
SELECT * FROM Users WHERE Email=$1 AND Password=$2;`

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
	if _, err := db.conns.Exec(db.context, createUsersTableSQL); err != nil {
		panic(err)
	}
	if _, err := db.conns.Exec(db.context, createBooksTableSQL); err != nil {
		panic(err)
	}
	if _, err := db.conns.Exec(db.context, createUserBooksTableSQL); err != nil {
		panic(err)
	}

	return db
}

func (db *DB) CreateUser(u User) error {
	_, err := db.conns.Exec(db.context, createUserSQL, u.Email, u.Password)
	return err
}

func (db *DB) ReadUser(u User) (User, error) {
	rows, err := db.conns.Query(db.context, readUserSQL, u.Email, u.Password)
	if err != nil {
		return u, err
	}

	exists := false
	for rows.Next() {
		exists = true
		if err := rows.Scan(&u.Id, &u.Email, &u.Password); err != nil {
			return u, err
		}
	}

	if err := rows.Err(); err != nil {
		return u, err
	}

	if !exists {
		return u, errors.New("User with those credentials not found.")
	}

	return u, nil
}
