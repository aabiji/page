package main

import (
	"fmt"
	"github.com/aabiji/read/epub"
	"io/ioutil"
	"os"
)

func main() {
	files, err := ioutil.ReadDir("books")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		e, err := epub.New("books/" + f.Name())
		if err != nil {
			panic(err)
		}
		e.Debug()

		// Check if all paths to files actually exist (Epub.bookPath() doesn't fail)
		for _, file := range e.Files {
			if _, err := os.Stat(file); err != nil {
				break
			}
		}
		for _, link := range e.TableOfContents {
			if _, err := os.Stat(link[1]); err != nil {
				break
			}
		}
		if os.Stat(e.CoverImagePath); err != nil {
			break
		}

		fmt.Println()
	}
}
