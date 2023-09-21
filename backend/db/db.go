package db

import (
    "context"
	"github.com/jackc/pgx/v5"
)

// NOTE: pgd.Conn is not thread sage. We should use a connection pool instead

const CREATE_BOOKS_TABLE = `
CREATE TABLE Books (
    BookId serial PRIMARY KEY,
    CoverImagePath text NOT NULL,
    Files text[] NOT NULL,
    TableOfContents jsonb[] NOT NULL,
    Info jsonb NOT NULL,
)
`

const CREATE_USERS_TABLE = `
CREATE TABLE Users (
    UserId serial PRIMARY KEY,
    Email string text NOT NULL,
    Password string text NOT NULL,
)
`
const CREATE_USERBOOKS_TABLE = `
CREATE TABLE UserBooks (
    UserId integer NOT NULL,
    BookId integer NOT NULL,
    CurrentPage integer NOT NULL,
    ScrollOffsets integer[] NOT NULL,
)
`

type DB struct {
    connection *pgx.Conn
}

func NewDatabase() DB {
    db := DB{}
    var err error
    databaseUrl := "postgres://aabiji:0000@localhost:5432/Page"
    db.connection, err = pgx.Connect(context.Background(), databaseUrl)
    if err != nil {
        panic(err)
    }
    return db
}
