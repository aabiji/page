mod epub;

fn main() {
    // Set the cwd to the root of the project's directory
    std::env::set_current_dir(std::env::var("CARGO_MANIFEST_DIR").unwrap()).unwrap();
    let _ = epub::Epub::new("books/Dune.epub");
}