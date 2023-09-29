module github.com/aabiji/page/backend

go 1.21.1

replace github.com/aabiji/page/backend/epub => ./epub

require (
	github.com/aabiji/page/backend/epub v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v5 v5.4.3
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
