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

func Contains(a []string, b string) bool {
	for _, s := range a {
		if s == b {
			return true
		}
	}
	return false
}

// Find a node attribute either by expected key and value or by key.
func FindAttribute(node *html.Node, key string, val string) string {
	for _, attribute := range node.Attr {
		if attribute.Key == key && val != "" && attribute.Val == val {
			return attribute.Val
		} else if attribute.Key == key {
			return attribute.Val
		}
	}
	return ""
}

// Set a node's attr to a value
func SetAttribute(node *html.Node, attr string, val string) {
	for i, _ := range node.Attr {
		if node.Attr[i].Key == attr {
			node.Attr[i].Val = val
		}
	}
}

func FindNode(root *html.Node, tagName string) *html.Node {
	if root == nil {
		return nil
	}

	if root.Type == html.ElementNode && root.Data == tagName {
		return root
	}

	for node := root.FirstChild; node != nil; node = node.NextSibling {
		found := FindNode(node, tagName)
		if found != nil {
			return found
		}
	}

	return nil
}

func ParseXML[T Container | NCX | Package](filename string) (T, error) {
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
func GetFileBase(filename string) string {
	i := 0
	parts := strings.Split(filename, string(os.PathSeparator))
	if len(parts) > 0 {
		i = len(parts) - 1
	}
	return strings.Split(parts[i], ".")[0]
}

// Unzip filename into outdir
func Unzip(filename, outdir string) error {
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
