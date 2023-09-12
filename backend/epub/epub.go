package epub

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"path/filepath"
	"strings"
)

// The directory where the epub files will be extracted into.
// Set by Storage struct in the server module.
var STORAGE_DIRECTORY string

var CONTENT_TYPES = map[string]string{
	"html":  "text/html",
	"xhtml": "application/xhtml+xml",
}

type File struct {
	Path         string
	ContentType  string
	ScrollOffset int
	document     *html.Node
}

type Epub struct {
	Name                string
	Info                Metadata
	Files               []File
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

	e := Epub{Name: GetFileBase(filename)}

	if err := Unzip(filename, e.bookPath()); err != nil {
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
	fmt.Printf("Fixed layout? %t\n", e.IsFixedLayout)
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
			if !Contains(path, p) && p != ".." {
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
	c, err := ParseXML[Container](e.bookPath("META-INF", "container.xml"))
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

// Get the contents of the css files linked in a html document's head node.
func (e *Epub) getLinkedCSS(head *html.Node) (string, error) {
	var css string
	var nodesToRemove []*html.Node

	for node := head.FirstChild; node != nil; node = node.NextSibling {
		if node.Data != "link" || FindAttribute(node, "rel", "stylesheet") == "" {
			continue
		}

		relativeCssPath := FindAttribute(node, "href", "")
		cssPath := e.bookPath(relativeCssPath)

		cssFile, err := os.ReadFile(cssPath)
		if err != nil {
			return "", err
		}

		css += string(cssFile)
		nodesToRemove = append(nodesToRemove, node)
	}

	for _, n := range nodesToRemove {
		head.RemoveChild(n)
	}

	return css, nil
}

// Inject a style node containing css into the file's html document
func (e *Epub) injectCSS(f *File) error {
	head := FindNode(f.document, "head")

	css, err := e.getLinkedCSS(head)
	if err != nil {
		return err
	}

	style := html.Node{Type: html.ElementNode, Data: "style"}
	style.AppendChild(&html.Node{Type: html.TextNode, Data: css})
	head.AppendChild(&style)
	return nil
}

// Replace relative paths to images in a file's html document with absolute paths
func (e *Epub) fixImageLinks(root *html.Node) error {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == "image" || root.Data == "img" {
		var attr string
		if root.Data == "image" {
			attr = "href"
		} else {
			attr = "src"
		}

		relativeImgPath := FindAttribute(root, attr, "")
		imgPath := e.bookPath(relativeImgPath)
		imgPath = strings.Replace(imgPath, STORAGE_DIRECTORY, "", -1)
		SetAttribute(root, attr, imgPath)
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		e.fixImageLinks(node)
	}

	return nil
}

// Replace a html file with a html document that embeds all of its styling.
// Replace relative paths to images with absolute paths.
func (e *Epub) updateFile(f *File) error {
	err := e.injectCSS(f)
	if err != nil {
		return err
	}

	err = e.fixImageLinks(f.document)
	if err != nil {
		return err
	}

	var htmlBytes bytes.Buffer
	err = html.Render(&htmlBytes, f.document)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(htmlBytes.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (e *Epub) processFile(relativePath string) (File, error) {
	f := File{Path: e.bookPath(relativePath)}

	filenameParts := strings.Split(f.Path, ".")
	extention := filenameParts[len(filenameParts)-1]
	f.ContentType = CONTENT_TYPES[extention]

	htmlContents, err := os.ReadFile(f.Path)
	if err != nil {
		return File{}, err
	}

	f.document, err = html.Parse(strings.NewReader(string(htmlContents)))
	if err != nil {
		return File{}, err
	}

	err = e.updateFile(&f)
	if err != nil {
		return File{}, err
	}

	return f, nil
}

func (e *Epub) parseContent() error {
	p, err := ParseXML[Package](e.bookPath(e.contentFilename))
	if err != nil {
		return err
	}

	// Get the list of ebook files
	items := make(map[string]string)
	for _, i := range p.Manifest.Items {
		items[i.Id] = i.Path
	}

	for _, i := range p.Spine.ITemRefs {
		file, err := e.processFile(items[i.Ref])
		if err != nil {
			return err
		}
		e.Files = append(e.Files, file)
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
	t, err := ParseXML[NCX](e.tableOfContentsPath)
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
