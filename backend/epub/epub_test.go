package epub

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func assertEq(t *testing.T, a any, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Found %v, want %v", a, b)
	}
}

func TestEpubProcessing(t *testing.T) {
	EXTRACT_DIRECTORY = "../../test_files"

	e, err := New("../../test_files/Dune.epub")
	if err != nil {
		panic(err)
	}

	// Check if all paths to files actually exist (Epub.bookPath() doesn't fail)
	for _, file := range e.Files {
		realPath := filepath.Join(EXTRACT_DIRECTORY, file)
		if _, err := os.Stat(realPath); err != nil {
			t.Errorf(fmt.Sprintf("%v", err))
		}
	}
	for _, link := range e.TableOfContents {
		realPath := filepath.Join(EXTRACT_DIRECTORY, link.Path)
		if _, err := os.Stat(realPath); err != nil {
			t.Errorf(fmt.Sprintf("%v", err))
		}
	}
	if os.Stat(e.CoverImagePath); err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}

	assertEq(t, e.Info.Language, "UND")
	assertEq(t, e.CoverImagePath, "Dune/cover.jpeg")
	assertEq(t, e.Info.Contributor, "calibre (0.6.52) [http://calibre-ebook.com]")
	assertEq(t, e.Info.Identifier, "7e87f9a1-8a4f-459a-8e58-e7032f0c67c6")
	assertEq(t, e.Info.Date, "2010-06-03T04:00:00+00:00")
	assertEq(t, e.Info.Author, "Herbert, Frank")
	assertEq(t, e.Info.Title, "Dune")
	assertEq(t, e.tableOfContentsPath, EXTRACT_DIRECTORY+"/Dune/toc.ncx")
	assertEq(t, e.contentFilename, "content.opf")
	assertEq(t, e.TableOfContents, []Section{
		{Name: "Dune", Path: "Dune/OEBPS/part1.xhtml"},
		{Name: "Book 1 DUNE", Path: "Dune/OEBPS/part2_split_000.xhtml"},
		{Name: "Book Two MUAD’DIB", Path: "Dune/OEBPS/part3_split_000.xhtml"},
		{Name: "Book Three THE PROPHET", Path: "Dune/OEBPS/part4_split_000.xhtml"},
	})
	assertEq(t, e.Files, []string{"Dune/titlepage.xhtml", "Dune/OEBPS/title.xhtml", "Dune/OEBPS/part1.xhtml", "Dune/OEBPS/part2_split_000.xhtml", "Dune/OEBPS/part2_split_001.xhtml", "Dune/OEBPS/part2_split_002.xhtml", "Dune/OEBPS/part3_split_000.xhtml", "Dune/OEBPS/part3_split_001.xhtml", "Dune/OEBPS/part4_split_000.xhtml", "Dune/OEBPS/part4_split_001.xhtml"})
}
