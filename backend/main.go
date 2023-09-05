package main

import (
	"github.com/aabiji/read/epub"
	"github.com/aabiji/read/server"
)

func main() {
	// NOTE: this is only needed for the prototype frontend
	_, err := epub.New("epub/tests/Dune.epub")
	if err != nil {
		panic(err)
	}

	server.Run("localhost:8080")
}
