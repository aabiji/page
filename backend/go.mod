module github.com/aabiji/read

go 1.18

replace github.com/aabiji/read/epub => ./epub

replace github.com/aabiji/read/server => ./server

require (
	github.com/aabiji/read/epub v0.0.0-00010101000000-000000000000
	github.com/aabiji/read/server v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/rs/cors v1.9.0 // indirect
)
