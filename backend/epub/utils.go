package epub

import (
	"archive/zip"
	"encoding/xml"
	"golang.org/x/net/html"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func contains(a []string, b string) bool {
	for _, s := range a {
		if s == b {
			return true
		}
	}
	return false
}

// Find a node attribute either by expected key and value or by key.
func findAttribute(node *html.Node, key string, val string) string {
	for _, attribute := range node.Attr {
		if attribute.Key == key && val != "" && attribute.Val == val {
			return attribute.Val
		} else if attribute.Key == key && val == "" {
			return attribute.Val
		}
	}
	return ""
}

// Set a node's attr to a value
func setAttribute(node *html.Node, attr string, val string) {
	for i, _ := range node.Attr {
		if node.Attr[i].Key == attr {
			node.Attr[i].Val = val
		}
	}
}

func findNode(root *html.Node, tagName string) *html.Node {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == tagName {
		return root
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		found := findNode(node, tagName)
		if found != nil {
			return found
		}
	}

	return nil
}

func parseHTML(filename string) (*html.Node, error) {
	htmlContents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	document, err := html.Parse(strings.NewReader(string(htmlContents)))
	if err != nil {
		return nil, err
	}

	return document, nil
}

func parseXML[T Container | NCX | Package](filename string) (T, error) {
	var t T

	file, err := os.ReadFile(filename)
	if err != nil {
		return t, err
	}

	err = xml.Unmarshal(file, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// Get the filename without the extention or leading directories.
// ex.
// '/path/to/file.txt' becomes 'file'
func getFileBase(filename string) string {
	i := 0
	parts := strings.Split(filename, string(os.PathSeparator))
	if len(parts) > 0 {
		i = len(parts) - 1
	}
	return strings.Split(parts[i], ".")[0]
}

// Unzip filename into outdir
func unzip(filename, outdir string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		filePath := filepath.Join(outdir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm) // create parent directory
		if err != nil {
			return err
		}

		dest, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		extractedFile, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, extractedFile)
		if err != nil {
			return err
		}

		dest.Close()
		extractedFile.Close()
	}

	return nil
}
