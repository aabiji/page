use std::path::{Path, PathBuf};
use serde::Deserialize;
use std::io::Read;

/// Read the contents of a file to a String
fn read_file_tostring(filename: &str) -> String {
    let mut contents = String::new();
    let mut file = std::fs::File::open(filename).unwrap();
    file.read_to_string(&mut contents).unwrap();
    return contents;
}

/// Get a filename without its extention
/// ex. '/path/to/file.txt' becomes 'file'
fn strip_filename(filename: &str) -> String {
    let path_parts: Vec<&str> = filename.split("/").collect();
    let file = if path_parts.len() > 0 {
        path_parts[path_parts.len() - 1]
    } else {
        path_parts[0]
    };

    let file_parts: Vec<&str> = file.split(".").collect();
    return String::from(file_parts[0]);
}

// META-INF/container.xml structure
#[derive(Debug, Deserialize)]
#[serde(rename_all = "kebab-case")]
struct Rootfile {
    full_path: String,
    media_type: String,
}
#[derive(Debug, Deserialize)]
struct Rootfiles {
    #[serde(rename = "rootfile")]
    rootfile: Rootfile,
}
#[derive(Debug, Deserialize)]
struct Container {
    #[serde(rename = "rootfiles")]
    rootfiles: Rootfiles,
}

pub struct Epub {
    name: String,
    filename: String,
    content_file: String,
}

const MIMETYPE: &str = "mimetype";
const CONTAINER: &str = "META-INF/container.xml";

impl Epub {
    pub fn new(filename: &str) -> Epub {
        let mut e = Epub{
            name: strip_filename(filename),
            content_file: String::new(),
            filename: String::from(filename),
        };
        e.extract();
        e.verify_mimetype();
        e.get_content_file();
        return e;
    }

    // file paths inside the epub files
    fn _path(&self, file: &str) -> String {
        return format!("{}/{}", self.name, file);
    }

    fn verify_mimetype(&self) {
        let mimetype = read_file_tostring(&self._path(MIMETYPE));
        assert_eq!(mimetype, "application/epub+zip");
    }

    // Read the required META-INF/container.xml to get the
    // full-path to the content file, ex. 'content.opf')
    fn get_content_file(&mut self) {
        let container_xml = read_file_tostring(&self._path(CONTAINER));
        let c: Container = serde_xml_rs::from_str(&container_xml).unwrap();
        assert_eq!(c.rootfiles.rootfile.media_type, "application/oebps-package+xml");
        self.content_file = c.rootfiles.rootfile.full_path;
    }

    // Unzip the epub file
    fn extract(&self) {
        let fullpath = std::path::Path::new(&self.filename);
        let base_path = Path::new(&self.name);

        let zipfile = std::fs::File::open(fullpath).unwrap();
        let mut archive = zip::ZipArchive::new(zipfile).unwrap();

        for i in 0..archive.len() {
            let mut file = archive.by_index(i).unwrap();
            let relative_output_path = match file.enclosed_name() {
                Some(path) => path.to_owned(),
                None => continue,
            };

            // Prepend base directory to relative_output_path, so that
            // the files are extracted into the destination directory
            let mut output_path = PathBuf::from(base_path);
            output_path.push(relative_output_path);

            if (*file.name()).ends_with("/") { // Is a directory ...
                std::fs::create_dir_all(&output_path).unwrap();
            } else {
                if let Some(parent_path) = output_path.parent() {
                    if !parent_path.exists() {
                        std::fs::create_dir_all(parent_path).unwrap();
                    }
                }
                let mut outfile = std::fs::File::create(&output_path).unwrap();
                std::io::copy(&mut file, &mut outfile).unwrap();
            }
        }
    }
}