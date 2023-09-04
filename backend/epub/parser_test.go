package epub

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func assertEq(t *testing.T, a any, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Got %v, found %v", a, b)
	}
}

func TestEpubPath(t *testing.T) {
	e := Epub{}
	assertEq(t, e.bookPath(), "BOOKS")

	e.Name = "Test"
	assertEq(t, e.bookPath(), "BOOKS/Test")

	e.contentFilename = "content.opf"
	assertEq(t, e.bookPath(), "BOOKS/Test")

	assertEq(t, e.bookPath("file.html"), "BOOKS/Test/file.html")

	e.contentFilename = "OEBPS/content.opf"
	assertEq(t, e.bookPath(), "BOOKS/Test/OEBPS")

	assertEq(t, e.bookPath("file.html"), "BOOKS/Test/OEBPS/file.html")
	assertEq(t, e.bookPath("OEBPS/file.html"), "BOOKS/Test/OEBPS/file.html")
}

func TestEpubProcessing(t *testing.T) {
	STORAGE_DIRECTORY = "tests"

	e, err := New("tests/Dune.epub")
	if err != nil {
		panic(err)
	}

	// Check if all paths to files actually exist (Epub.bookPath() doesn't fail)
	for _, file := range e.Files {
		if _, err := os.Stat(file); err != nil {
			t.Errorf(fmt.Sprintf("%v", err))
		}
	}
	for _, link := range e.TableOfContents {
		if _, err := os.Stat(link[1]); err != nil {
			t.Errorf(fmt.Sprintf("%v", err))
		}
	}
	if os.Stat(e.CoverImagePath); err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}

	assertEq(t, e.Info.Language, "UND")
	assertEq(t, e.CoverImagePath, "tests/Dune/titlepage.xhtml")
	assertEq(t, e.Info.Contributor, "calibre (0.6.52) [http://calibre-ebook.com]")
	assertEq(t, e.Info.Identifier, "7e87f9a1-8a4f-459a-8e58-e7032f0c67c6")
	assertEq(t, e.Info.Date, "2010-06-03T04:00:00+00:00")
	assertEq(t, e.Info.Author, "Herbert, Frank")
	assertEq(t, e.Info.Title, "Dune")
	assertEq(t, e.tableOfContentsPath, "tests/Dune/toc.ncx")
	assertEq(t, e.contentFilename, "content.opf")
	assertEq(t, e.TableOfContents, [][2]string{
		{"Dune", "tests/Dune/OEBPS/part1.xhtml"},
		{"Book 1 DUNE", "tests/Dune/OEBPS/part2_split_000.xhtml"},
		{"Book Two MUAD’DIB", "tests/Dune/OEBPS/part3_split_000.xhtml"},
		{"Book Three THE PROPHET", "tests/Dune/OEBPS/part4_split_000.xhtml"},
	})
	assertEq(t, e.Files, []string{"tests/Dune/titlepage.xhtml", "tests/Dune/OEBPS/title.xhtml", "tests/Dune/OEBPS/part1.xhtml", "tests/Dune/OEBPS/part2_split_000.xhtml", "tests/Dune/OEBPS/part2_split_001.xhtml", "tests/Dune/OEBPS/part2_split_002.xhtml", "tests/Dune/OEBPS/part3_split_000.xhtml", "tests/Dune/OEBPS/part3_split_001.xhtml", "tests/Dune/OEBPS/part4_split_000.xhtml", "tests/Dune/OEBPS/part4_split_001.xhtml"})
	assertEq(t, e.IsFixedLayout, false)
}
