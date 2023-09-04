package epub

import "encoding/xml"

/*
container.xml structure:
  <?xml version=""?>
  <container version="" xmlns="">
    <rootfiles>
      <rootfile full-path="" media-type=""/>
    </rootfiles>
  </container>
*/
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

/*
content.opf structure:
  <?xml version="" encoding="" standalone=""?>
  <package version="" xmlns="" unique-identifier="">
    <metadata xmlns:dc="" xmlns:opf="">
      <meta name="" content="" />
      <dc:title></dc:title>
      <dc:creator></dc:creator>
      <dc:subject></dc:subject>
      <dc:description></dc:description>
      <dc:publisher></dc:publisher>
      <dc:date></dc:date>
      <dc:source></dc:source>
      <dc:relation></dc:relation>
      <dc:coverage></dc:coverage>
      <dc:contributor></dc:contributor>
      <dc:rights></dc:rights> <dc:language></dc:language>
      <dc:identifier id=""></dc:identifier>
    </metadata>
    <manifest>
      <item id="" href="" media-type="" />
      ...
    </manifest>
    <spine toc="">
      <itemref idref="" />
      ...
    </spine>
    <guide>
      <reference type="cover" href="cover.xhtml" title="Cover" />
      ...
    </guide>
  </package>
*/
type Meta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:"content,attr"`
}

type Metadata struct {
	XMLName     xml.Name `xml:"metadata"`
	Language    string   `xml:"language"`
	Author      string   `xml:"creator"`
	Title       string   `xml:"title"`
	Identifier  string   `xml:"identifier"`
	Contributor string   `xml:"contributor"`
	Rights      string   `xml:"rights"`
	Source      string   `xml:"source"`
	Coverage    string   `xml:"coverage"`
	Relation    string   `xml:"relation"`
	Publisher   string   `xml:"publisher"`
	Description string   `xml:"description"`
	Date        string   `xml:"date"`
	Subjects    []string `xml:"subject"`
	Meta        []Meta   `xml:"meta"`
}

type Item struct {
	XMLName   xml.Name `xml:"item"`
	Path      string   `xml:"href,attr"`
	Id        string   `xml:"id,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type Manifest struct {
	XMLName xml.Name `xml:"manifest"`
	Items   []Item   `xml:"item"`
}

type ItemRef struct {
	XMLName xml.Name `xml:"itemref"`
	Ref     string   `xml:"idref,attr"`
}

type Spine struct {
	XMLName  xml.Name  `xml:"spine"`
	Toc      string    `xml:"toc,attr"`
	ITemRefs []ItemRef `xml:"itemref"`
}

type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Path    string   `xml:"href,attr"`
	Type    string   `xml:"type,attr"`
	Title   string   `xml:"title,attr"`
}

type Guide struct {
	XMLName    xml.Name    `xml:"guide"`
	References []Reference `xml:"reference"`
}

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
	Guide    Guide    `xml:"guide"`
}

/*
toc.ncx structure:
  <ncx xmlns="" version="">
    <head>
      <meta name="" content="" />
      ...
    </head>
    <docTitle>
      <text></text>
    </docTitle>
    <navMap>
      <navPoint>
        <navPoint id="" playOrder="">
        <navLabel><text></text></navLabel>
        <content src="" />
        ... (navPoint)
      </navPoint>
      ...
    </navMap>
  </ncx>
*/

type DocTitle struct {
	XMLName xml.Name `xml:"docTitle"`
	Text    string   `xml:"text"`
}

type Head struct {
	Metadata []Meta `xml:"meta"`
}

type NavLabel struct {
	XMLName xml.Name `xml:"navLabel"`
	Text    string   `xml:"text"`
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Source  string   `xml:"src,attr"`
}

type NavPoint struct {
	XMLName   xml.Name   `xml:"navPoint"`
	Id        string     `xml:"id,attr"`
	PlayOrder string     `xml:"playOrder,attr"`
	Label     NavLabel   `xml:"navLabel"`
	Content   Content    `xml:"content"`
	Children  []NavPoint `xml:"navPoint"`
}

type NavMap struct {
	XMLName   xml.Name   `xml:"navMap"`
	NavPoints []NavPoint `xml:"navPoint"`
}

type NCX struct {
	Head  Head     `xml:"head"`
	Title DocTitle `xml:"docTitle"`
	Map   NavMap   `xml:"navMap"`
}
