module github.com/aabiji/backend/page

go 1.18

replace github.com/aabiji/page/backend/epub => ./epub

replace github.com/aabiji/page/backend/server => ./server

replace github.com/aabiji/page/backend/db => ./db

require github.com/aabiji/page/backend/server v0.0.0-00010101000000-000000000000

require (
	github.com/aabiji/page/backend/epub v0.0.0-20230922023337-0abe6b36a4fc // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/rs/cors v1.10.0 // indirect
	golang.org/x/net v0.15.0 // indirect
)
