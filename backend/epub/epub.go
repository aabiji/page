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

type Section struct {
	Path string
	Name string
}

type Epub struct {
	Name                string
	Info                Metadata
	Files               []string
	TableOfContents     []Section
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

	if err := Unzip(filename, e.absolutePath()); err != nil {
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
	fmt.Printf("Cover image: %s\n", e.absolutePath(e.CoverImagePath))
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
		fmt.Printf("URL Path %s | Local Path %s\n", f, e.absolutePath(f))
	}
	fmt.Println("Table of contents: ")
	for _, t := range e.TableOfContents {
		fmt.Printf("%s : %s \n", t.Name, t.Path)
	}
}

// Path to a file inside the extracted epub file directory
func (e *Epub) absolutePath(files ...string) string {
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

func (e *Epub) urlPath(files ...string) string {
	s := e.absolutePath(files...)
	return strings.Replace(s, STORAGE_DIRECTORY+"/", "", -1)
}

func (e *Epub) verifyMimetype() error {
	path := e.absolutePath("mimetype")

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
	c, err := ParseXML[Container](e.absolutePath("META-INF", "container.xml"))
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
			e.coverPath = r.Path
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
		e.coverPath = items[e.coverPath]
	}
}

// Get an absolute path to the epub's cover image.
func (e *Epub) getCoverImagePath() error {
	fileParts := strings.Split(e.coverPath, ".")
	extension := fileParts[len(fileParts)-1]
	if extension != "xhtml" && extension != "html" {
		e.CoverImagePath = e.urlPath(e.coverPath)
		return nil // Cover image path is already found
	}

	e.coverPath = e.absolutePath(e.coverPath)
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
		cssPath := e.absolutePath(relativeCssPath)

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
func (e *Epub) injectCSS(root *html.Node) error {
	head := FindNode(root, "head")
	if head == nil {
		return errors.New(fmt.Sprintf("<head></head> not found"))
	}

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
		var imgSrc string
		if root.Data == "image" {
			imgSrc = "href"
		} else {
			imgSrc = "src"
		}

		relativeImgPath := FindAttribute(root, imgSrc, "")
		SetAttribute(root, imgSrc, e.urlPath(relativeImgPath))
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		e.fixImageLinks(node)
	}

	return nil
}

func (e *Epub) fixFileLinks(root *html.Node) error {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == "a" {
		link := FindAttribute(root, "href", "")
		if link == "" {
			return nil
		}

		urlMatch := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
		regex := regexp.MustCompile(urlMatch)
		matched := regex.Match([]byte(link))
		if matched {
			return nil
		}

		link = e.urlPath(link)
		SetAttribute(root, "href", link)
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		e.fixFileLinks(node)
	}

	return nil
}

// Replace a html file with a html document that embeds all of its styling.
// Replace relative paths to images within the document with absolute paths.
// Replace relative paths to files within the document with the url paths.
func (e *Epub) processFile(relativePath string) (string, error) {
	fileUrlPath := e.urlPath(relativePath)
	path := e.absolutePath(relativePath)

	var err error
	document, err := ParseHTML(path)
	if err != nil {
		return "", err
	}

	err = e.injectCSS(document)
	if err != nil {
		return "", err
	}

	err = e.fixImageLinks(document)
	if err != nil {
		return "", err
	}

	err = e.fixFileLinks(document)
	if err != nil {
		return "", err
	}

	var htmlBytes bytes.Buffer
	err = html.Render(&htmlBytes, document)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(htmlBytes.Bytes())
	if err != nil {
		return "", err
	}

	return fileUrlPath, nil
}

// Remove html tag elements from the epub description
func (e *Epub) cleanDescription() {
	htmlTagRegex := "<[^>]*>"
	regex := regexp.MustCompile(htmlTagRegex)
	e.Info.Description = string(regex.ReplaceAll([]byte(e.Info.Description), []byte{}))
}

func (e *Epub) parseContent() error {
	p, err := ParseXML[Package](e.absolutePath(e.contentFilename))
	if err != nil {
		return err
	}

	// Get the list of ebook files
	items := make(map[string]string)
	for _, i := range p.Manifest.Items {
		items[i.Id] = i.Path
	}

	for _, i := range p.Spine.ITemRefs {
		fileUrlPath, err := e.processFile(items[i.Ref])
		if err != nil {
			return err
		}
		e.Files = append(e.Files, fileUrlPath)
	}

	e.Info = p.Metadata
	e.cleanDescription()
	if len(e.Info.Subjects) == 0 {
		e.Info.Subjects = append(e.Info.Subjects, "")
	}

	e.getCoverPath(p, items)
	e.getCoverImagePath()

	e.tableOfContentsPath = e.absolutePath(items[p.Spine.TableOfContents])
	return nil
}

func (e *Epub) assembleTableOfContents(points []NavPoint) []Section {
	links := []Section{}
	for _, n := range points {
		path := e.urlPath(n.Content.Source)
		entry := Section{Name: n.Label.Text, Path: path}

		links = append(links, entry)
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
