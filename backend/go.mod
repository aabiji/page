module github.com/aabiji/backend/page

go 1.18

replace github.com/aabiji/page/backend/epub => ./epub

replace github.com/aabiji/page/backend/server => ./server

replace github.com/aabiji/page/backend/db => ./db

require github.com/aabiji/page/backend/server v0.0.0-00010101000000-000000000000

require (
	github.com/aabiji/page/backend/db v0.0.0-00010101000000-000000000000 // indirect
	github.com/aabiji/page/backend/epub v0.0.0-20230922023337-0abe6b36a4fc // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/rs/cors v1.10.0 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
