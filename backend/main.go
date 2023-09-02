package main

import (
	"fmt"
	"github.com/aabiji/read/epub"
	"io/ioutil"
)

func main() {
	files, err := ioutil.ReadDir("books")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			e, err := epub.New("books/" + file.Name())
			if err != nil {
				panic(err)
			}

			fmt.Println(file.Name())
			e.Debug()
			fmt.Println()
		}
	}
}
