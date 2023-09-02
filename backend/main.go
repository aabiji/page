package main

import "github.com/aabiji/read/epub"

func main() {
    _, err := epub.New("books/Dune.epub")
    if err != nil {
        panic(err)
    }
}
