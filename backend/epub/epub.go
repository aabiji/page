package epub

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The directory where the epub files will be extracted into.
// Set by Storage struct in the server module.
var STORAGE_DIRECTORY string

type File struct {
	Path        string
	ContentType string
	document    *html.Node
}

type Epub struct {
	Name                string
	Info                Metadata
	Files               []File
	TableOfContents     [][2]string
	IsFixedLayout       bool
	CoverImagePath      string
	tableOfContentsPath string
	contentFilename     string
	coverPath           string
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
	fmt.Printf("Cover image: %s\n", e.coverPath)
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
	fmt.Println("Files: ")
	for _, f := range e.Files {
		fmt.Printf("Path %s | Content type: %s\n", f.Path, f.ContentType)
	}
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

// Traverse the xml to find a path to the epub's cover.
func (e *Epub) getCoverPath(p Package, items map[string]string) {
	for _, r := range p.Guide.References {
		if r.Type == "cover" {
			e.coverPath = e.bookPath(r.Path)
			break
		}
	}

	// If there's no references in the guide node check if the meta nodes
	if e.coverPath == "" {
		for _, m := range p.Metadata.Meta {
			if m.Name == "cover" {
				e.coverPath = m.Content
				break
			}
		}
	}

	// If the result from searching the meta nodes isn't a file
	// use the result as key to get a file
	if !strings.Contains(e.coverPath, ".") {
		e.coverPath = e.bookPath(items[e.coverPath])
	} else {
		e.coverPath = e.bookPath(e.coverPath)
	}
}

// Get an absolute path to the epub's cover image.
func (e *Epub) getCoverImagePath() error {
	fileParts := strings.Split(e.coverPath, ".")
	extension := fileParts[len(fileParts)-1]
	if extension != "xhtml" && extension != "html" {
		e.CoverImagePath = e.coverPath
		return nil // Cover image path is already found
	}

	document, err := ParseHTML(e.coverPath)
	if err != nil {
		return err
	}

	imageNode := FindNode(document, "img")
	if imageNode == nil {
		imageNode = FindNode(document, "image")
	}

	if imageNode == nil {
		return nil // Epub doesn't have cover image
	}

	e.CoverImagePath = FindAttribute(imageNode, "src", "")
	if e.CoverImagePath == "" {
		e.CoverImagePath = FindAttribute(imageNode, "href", "")
	}

	e.CoverImagePath = e.bookPath(e.CoverImagePath)
	return nil
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
	extension := filenameParts[len(filenameParts)-1]
	if extension == "html" {
		f.ContentType = "text/html"
	} else if extension == "xhtml" {
		f.ContentType = "application/xhtml+xml"
	}

	var err error
	f.document, err = ParseHTML(f.Path)
	if err != nil {
		return File{}, err
	}

	err = e.updateFile(&f)
	if err != nil {
		return File{}, err
	}

	return f, nil
}

// Remove html tag elements from the epub description
func (e *Epub) cleanDescription() {
	htmlTagRegex := "<[^>]*>"
	regex := regexp.MustCompile(htmlTagRegex)
	e.Info.Description = string(regex.ReplaceAll([]byte(e.Info.Description), []byte{}))
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
	e.cleanDescription()

	e.getCoverPath(p, items)
	e.getCoverImagePath()

	e.tableOfContentsPath = e.bookPath(items[p.Spine.TableOfContents])
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
