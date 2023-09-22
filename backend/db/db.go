package db

import (
	"context"
    "encoding/json"
    "errors"
	"github.com/jackc/pgx/v5/pgxpool"
    "github.com/aabiji/page/backend/epub"
)

// SQL Queries:
const create_books_table = `
CREATE TABLE IF NOT EXISTS Books (
    BookId serial PRIMARY KEY,
    CoverImagePath text NOT NULL,
    Files text[] NOT NULL,
    TableOfContents jsonb[] NOT NULL,
    Info jsonb NOT NULL
);`

const create_users_table = `
CREATE TABLE IF NOT EXISTS Users (
    UserId serial PRIMARY KEY,
    Email text NOT NULL,
    Password text NOT NULL
);`

const create_userbooks_table = `
CREATE TABLE IF NOT EXISTS UserBooks (
    UserId integer NOT NULL,
    BookId integer NOT NULL,
    CurrentPage integer NOT NULL,
    ScrollOffsets integer[] NOT NULL,
);`

const create_user = `
INSERT INTO Users (Email, Password) VALUES ($1, $2);`

const read_user = `
SELECT * FROM Users WHERE Email=$1 AND Password=$2;`

const create_book = `
INSERT INTO Books (CoverImagePath, Files, TableOfContents, Info) VALUES ($1, $2, $3, $4);`


type DB struct {
    conns *pgxpool.Pool
    context context.Context
}

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
    if _, err := db.conns.Exec(db.context, create_users_table); err != nil {
        panic(err)
    }
    if _, err := db.conns.Exec(db.context, create_books_table); err != nil {
        panic(err)
    }
    if _, err := db.conns.Exec(db.context, create_userbooks_table); err != nil {
        panic(err)
    }

    return db
}

func (db *DB) CreateUser(email string, hashedPassword string) error {
    _, err := db.conns.Exec(db.context, create_user, email, hashedPassword)
    return err
}

func (db *DB) GetUser(email string, hashedPassword string) (string, string, string, error) {
    var userid string
    rows, err := db.conns.Query(db.context, read_user, email, hashedPassword)
    if err != nil {
        return "", "", "", err
    }

    exists := false
    for rows.Next() {
        exists = true
        if err := rows.Scan(&userid, &email, &hashedPassword); err != nil {
            return "", "", "", err
        }
    }

    if err := rows.Err(); err != nil {
        return "", "", "", err
    }

    if !exists {
        return "", "", "", errors.New("User not found.")
    }

    return userid, email, hashedPassword, nil
}

func (db *DB) CreateBook(e *epub.Epub) error {
    toc, err := json.Marshal(e.TableOfContents)
    if err != nil {
        return err
    }

    info, err := json.Marshal(e.Info)
    if err != nil {
        return err
    }

    _, err = db.conns.Exec(db.context, create_book, e.CoverImagePath, e.Files, toc, info)
    return err
}
