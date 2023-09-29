package epub

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The directory where the epub files will be extracted into.
var EXTRACT_DIRECTORY string

type Section struct {
	Path string `json:"Path"`
	Name string `json:"Name"`
}

type Epub struct {
	Name                string
	Info                Metadata
	Files               []string
	TableOfContents     []Section
	CoverImagePath      string
	tableOfContentsPath string
	contentFilename     string
	coverPath           string
}

func New(filename string) (Epub, error) {
	if !strings.Contains(filename, ".epub") {
		return Epub{}, errors.New("Invalid epub file")
	}

	e := Epub{Name: getFileBase(filename)}

	if err := unzip(filename, e.absolutePath("")); err != nil {
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
func (e *Epub) absolutePath(file string) string {
	basePath := filepath.Join(EXTRACT_DIRECTORY, e.Name)
	if file == "" {
		return basePath
	}
	pathParts := strings.Split(file, "/")
	targetFile := pathParts[len(pathParts)-1]

	var foundPath string
	filepath.WalkDir(basePath, func(path string, info fs.DirEntry, err error) error {
		if !info.IsDir() && targetFile == info.Name() {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})

	return foundPath
}

func (e *Epub) urlPath(file string) string {
	s := e.absolutePath(file)
	replace := EXTRACT_DIRECTORY
	if replace != "" {
		replace += "/"
	}
	return strings.Replace(s, replace, "", -1)
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
	c, err := parseXML[Container](e.absolutePath("META-INF/container.xml"))
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
	document, err := parseHTML(e.coverPath)
	if err != nil {
		return err
	}

	imageNode := findNode(document, "img")
	if imageNode == nil {
		imageNode = findNode(document, "image")
	}

	if imageNode == nil {
		return nil // Epub doesn't have cover image
	}

	e.CoverImagePath = findAttribute(imageNode, "src", "")
	if e.CoverImagePath == "" {
		e.CoverImagePath = findAttribute(imageNode, "href", "")
	}

	return nil
}

// Get the contents of the css files linked in a html document's head node.
func (e *Epub) getLinkedCSS(head *html.Node) (string, error) {
	var css string
	var nodesToRemove []*html.Node

	for node := head.FirstChild; node != nil; node = node.NextSibling {
		if node.Data != "link" || findAttribute(node, "rel", "stylesheet") == "" {
			continue
		}

		relativeCssPath := findAttribute(node, "href", "")
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
	head := findNode(root, "head")
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

func (e *Epub) fixLinks(root *html.Node) error {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == "a" {
		link := findAttribute(root, "href", "")
		if link == "" {
			return nil
		}

		urlMatch := `^(https?|ftp)://[^\s/$.?#].[^\s]*$`
		regex := regexp.MustCompile(urlMatch)
		matched := regex.Match([]byte(link))
		if matched {
			return nil
		}
		setAttribute(root, "href", e.urlPath(link))
	} else if root.Type == html.ElementNode && root.Data == "image" || root.Data == "img" {
		var imgSrc string
		if root.Data == "image" {
			imgSrc = "href"
		} else {
			imgSrc = "src"
		}

		relativeImgPath := findAttribute(root, imgSrc, "")
		setAttribute(root, imgSrc, e.urlPath(relativeImgPath))
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		e.fixLinks(node)
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
	document, err := parseHTML(path)
	if err != nil {
		return "", err
	}

	err = e.injectCSS(document)
	if err != nil {
		return "", err
	}

	err = e.fixLinks(document)
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
	p, err := parseXML[Package](e.absolutePath(e.contentFilename))
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
	// TODO: For epub 3.0 archives, the toc.ncx file might not exist.
	// Instead, we'll need to parse the toc.xhtml/toc.html file.
	if !strings.Contains(e.tableOfContentsPath, ".") {
		return nil
	}

	t, err := parseXML[NCX](e.tableOfContentsPath)
	if err != nil {
		return err
	}

	e.TableOfContents = e.assembleTableOfContents(t.Map.NavPoints)
	return nil
}
