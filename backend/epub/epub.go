package epub

import (
    "archive/zip"
    "path/filepath"
    "os"
    "io"
    "strings"
    "io/ioutil"
    "encoding/xml"
)

type Epub struct {
    Name string

    filename string
    contentFilename string
}

func New(filename string) Epub {
    e := Epub{filename: filename}
    e.Name = strings.Split(filename, ".")[0]

    e.unzip()
    e.readContainer()
    return e
}

func (e *Epub) unzip() {
    archive, err := zip.OpenReader(e.filename)
    if err != nil {
        panic(err)
    }
    defer archive.Close()

    for _, file := range archive.File {
        filePath := filepath.Join(e.Name, file.Name)

        if !strings.HasPrefix(filePath, filepath.Clean(e.Name)+string(os.PathSeparator)) {
            continue // is not prefixed by 'outputDir/'
        }
        if file.FileInfo().IsDir() {
            os.MkdirAll(filePath, os.ModePerm)
            continue
        }

        _ =  os.MkdirAll(filepath.Dir(filePath), os.ModePerm) // create parent directory
        destination, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
        extractedFile, _ := file.Open()
        _, _ = io.Copy(destination, extractedFile)

        destination.Close()
        extractedFile.Close()
    }
}

// Container XML structure
type Rootfile struct {
    XMLName xml.Name `xml:"rootfile"`
    FullPath string `xml:"full-path,attr"`
    MediaType string `xml:"media-type,attr"`
}

type Rootfiles struct {
    XMLName xml.Name `xml:"rootfiles"`
    Rootfile Rootfile `xml:"rootfile"`
}

type Container struct {
    XMLName xml.Name `xml:"container"`
    Rootfiles Rootfiles `xml:"rootfiles"`
}

func (e *Epub) readContainer() {
    file, _ := os.Open("Dune/META-INF/container.xml")
    defer file.Close()
    fileBytes, _ := ioutil.ReadAll(file)

    var c Container
    xml.Unmarshal(fileBytes, &c)

    e.contentFilename = c.Rootfiles.Rootfile.FullPath
    if c.Rootfiles.Rootfile.MediaType != "application/oebps-package+xml" {
        panic("Invalid epub file")
    }
}
