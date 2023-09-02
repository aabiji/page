package epub

import (
	"archive/zip"
	"encoding/xml"
    "errors"
	"io"
	"io/ioutil"
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
	Name string

	filename        string
	contentFilename string
}

func New(filename string) (Epub, error) {
	e := Epub{filename: filename}
    e.Name = getFileBase(filename)

    err := e.unzip()
    if err != nil {
        return Epub{}, err
    }

	err = e.readContainer()
    if err != nil {
        return Epub{}, err
    }

	return e, nil
}

func (e *Epub) unzip()  error {
	archive, err := zip.OpenReader(e.filename)
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

type Rootfile struct {
	XMLName   xml.Name `xml:"rootfile"`
	FullPath  string   `xml:"full-path,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type Rootfiles struct {
	XMLName  xml.Name `xml:"rootfiles"`
	Rootfile Rootfile `xml:"rootfile"`
}

type Container struct {
	XMLName   xml.Name  `xml:"container"`
	Rootfiles Rootfiles `xml:"rootfiles"`
}

func (e *Epub) readContainer() error {
    path := filepath.Join(e.Name, "META-INF", "container.xml")

	file, err := os.Open(path)
    if err != nil {
        return err
    }
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

	var c Container
    err = xml.Unmarshal(fileBytes, &c)
    if err != nil {
        return err
    }

	e.contentFilename = c.Rootfiles.Rootfile.FullPath
	if c.Rootfiles.Rootfile.MediaType != "application/oebps-package+xml" {
	    return errors.New("Invalid epub file")
	}

    return nil
}
