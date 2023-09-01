use std::path::{Path, PathBuf};

/// Set the cwd to the root of the project's directory
pub fn set_root_cwd() {
    std::env::set_current_dir(std::env::var("CARGO_MANIFEST_DIR").unwrap()).unwrap();
}

// TODO: properly propagte errors
pub fn extract_zipfile(file: &str, output_directory: &str) {
    let fullpath = std::path::Path::new(file);
    let base_path = Path::new(output_directory);

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