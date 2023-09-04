package epub

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The directory where the epub files will be extracted into
var STORAGE_DIRECTORY = "BOOKS"

type Epub struct {
	Name                string
	Info                Metadata
	Files               []string
	CoverImagePath      string
	TableOfContents     [][2]string
	IsFixedLayout       bool
	tableOfContentsPath string
	contentFilename     string
}

func New(filename string) (Epub, error) {
	if !strings.Contains(filename, ".epub") {
		return Epub{}, errors.New("Invalid epub file")
	}

	e := Epub{Name: getFileBase(filename)}

	if err := unzip(filename, e.bookPath()); err != nil {
		return Epub{}, err
	}
	if err := e.verifyMimetype(); err != nil {
		return Epub{}, err
	}
	if err := e.parseContainer(); err != nil {
		return Epub{}, err
	}
	if err := e.parseContent(); err != nil {
		return Epub{}, err
	}
	if err := e.parseTableOfContents(); err != nil {
		return Epub{}, err
	}

	return e, nil
}

func (e *Epub) Debug() {
	fmt.Printf("%s by %s in %s\n", e.Info.Title, e.Info.Author, e.Info.Date)
	fmt.Printf("Description: %s\n", e.Info.Description)
	fmt.Printf("Cover image: %s\n", e.CoverImagePath)
	fmt.Printf("Files: %v\n", e.Files)
	fmt.Printf("Subjects: %v\n", e.Info.Subjects)
	fmt.Printf("Publisher: %s\n", e.Info.Publisher)
	fmt.Printf("Language: %s\n", e.Info.Language)
	fmt.Printf("Relation: %s\n", e.Info.Relation)
	fmt.Printf("Coverage: %s\n", e.Info.Coverage)
	fmt.Printf("Source: %s\n", e.Info.Source)
	fmt.Printf("Rights: %s\n", e.Info.Rights)
	fmt.Printf("Contributor: %s\n", e.Info.Contributor)
	fmt.Printf("Identifier: %s\n", e.Info.Identifier)
	fmt.Println("Table of contents: ")
	for _, l := range e.TableOfContents {
		fmt.Printf("%s : %s \n", l[0], l[1])
	}
}

// Path to a file inside the extracted epub file directory
func (e *Epub) bookPath(files ...string) string {
	path := []string{STORAGE_DIRECTORY, e.Name} // BOOKS/<BOOK_NAME>

	temp := strings.Split(e.contentFilename, "/")
	internalDirectories := temp[0 : len(temp)-1]
	path = append(path, internalDirectories...)

	for _, f := range files {
		pathParts := strings.Split(f, "/")
		for _, p := range pathParts {
			if !Contains(path, p) {
				path = append(path, p)
			}
		}
	}

	return filepath.Join(path...)
}

func (e *Epub) verifyMimetype() error {
	path := e.bookPath("mimetype")

	mimetype, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if string(mimetype) != "application/epub+zip" {
		return errors.New("Invalid epub file")
	}

	return nil
}

func (e *Epub) parseContainer() error {
	c, err := parseXML[Container](e.bookPath("META-INF", "container.xml"))
	if err != nil {
		return err
	}

	e.contentFilename = c.Rootfiles.Rootfile.FullPath
	if c.Rootfiles.Rootfile.MediaType != "application/oebps-package+xml" {
		return errors.New("Invalid epub file")
	}

	return nil
}

func (e *Epub) getCoverImagePath(p Package, items map[string]string) {
	for _, r := range p.Guide.References {
		if r.Type == "cover" {
			e.CoverImagePath = e.bookPath(r.Path)
			break
		}
	}

	// If there's no references in the guide node check if the meta nodes
	if e.CoverImagePath == "" {
		for _, m := range p.Metadata.Meta {
			if m.Name == "cover" {
				e.CoverImagePath = m.Content
				break
			}
		}
	}

	// If the result from searching the meta nodes isn't a file
	// use the result as key to get a file
	if !strings.Contains(e.CoverImagePath, ".") {
		e.CoverImagePath = e.bookPath(items[e.CoverImagePath])
	} else {
		e.CoverImagePath = e.bookPath(e.CoverImagePath)
	}
}

func (e *Epub) parseContent() error {
	p, err := parseXML[Package](e.bookPath(e.contentFilename))
	if err != nil {
		return err
	}

	// Get the list of ebook files
	items := make(map[string]string)
	for _, i := range p.Manifest.Items {
		items[i.Id] = i.Path
	}

	for _, i := range p.Spine.ITemRefs {
		e.Files = append(e.Files, e.bookPath(items[i.Ref]))
	}

	e.Info = p.Metadata
	e.getCoverImagePath(p, items)
	e.tableOfContentsPath = e.bookPath(items[p.Spine.Toc])

	return nil
}

func (e *Epub) assembleTableOfContents(points []NavPoint) [][2]string {
	links := [][2]string{}
	for _, n := range points {
		path := e.bookPath(n.Content.Source)
		links = append(links, [2]string{n.Label.Text, path})

		links = append(links, e.assembleTableOfContents(n.Children)...)
	}
	return links
}

func (e *Epub) parseTableOfContents() error {
	t, err := parseXML[NCX](e.tableOfContentsPath)
	if err != nil {
		return err
	}

	for _, m := range t.Head.Metadata {
		if m.Name == "dtb:totalPageCount" || m.Name == "dtb:maxPageNumber" {
			if m.Content != "0" {
				e.IsFixedLayout = true
				break
			}
		}
	}

	e.TableOfContents = e.assembleTableOfContents(t.Map.NavPoints)
	return nil
}
