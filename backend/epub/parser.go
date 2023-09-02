package epub

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

type Epub struct {
	Name            string
	Info            Metadata
	Files           []string
	CoverPath       string
	contentFilename string
}

func New(filename string) (Epub, error) {
	e := Epub{Name: getFileBase(filename)}

	err := e.unzip(filename)
	if err != nil {
		return Epub{}, err
	}

	err = e.parseContainer()
	if err != nil {
		return Epub{}, err
	}

	err = e.parseContent()
	if err != nil {
		return Epub{}, err
	}

	return e, nil
}

func (e *Epub) Debug() {
	fmt.Printf("%s by %s in %s\n", e.Info.Title, e.Info.Author, e.Info.Date)
	fmt.Printf("Description: %s\n", e.Info.Description)
	fmt.Printf("Cover: %s\n", e.CoverPath)
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
}

func (e *Epub) unzip(filename string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		filePath := filepath.Join(e.Name, file.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(e.Name)+string(os.PathSeparator)) {
			continue // invalid file path
		}
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

func (e *Epub) parseContainer() error {
	path := filepath.Join(e.Name, "META-INF", "container.xml")

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var c Container
	err = xml.Unmarshal(file, &c)
	if err != nil {
		return err
	}

	e.contentFilename = c.Rootfiles.Rootfile.FullPath
	if c.Rootfiles.Rootfile.MediaType != "application/oebps-package+xml" {
		return errors.New("Invalid epub file")
	}

	return nil
}

func (e *Epub) parseContent() error {
	path := filepath.Join(e.Name, e.contentFilename)

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var p Package
	err = xml.Unmarshal(file, &p)
	if err != nil {
		return err
	}

	items := make(map[string]string)
	for _, i := range p.Manifest.Items {
		items[i.Id] = i.Path
	}

	for _, i := range p.Spine.ITemRefs {
		e.Files = append(e.Files, items[i.Ref])
	}

	for _, r := range p.Guide.References {
		if r.Type == "cover" {
			e.CoverPath = r.Path
		}
	}

	if e.CoverPath == "" { // no references in <guide></guide>
		for _, m := range p.Metadata.Meta {
			if m.Name == "cover" {
				e.CoverPath = m.Content
			}
		}
	}

	if !strings.Contains(e.CoverPath, ".") { // not a file
		e.CoverPath = items[e.CoverPath]
	}

	e.Info = p.Metadata
	return nil
}
