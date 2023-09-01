mod extract;

fn main() {
    extract::set_root_cwd();
    extract::extract_zipfile("books/Dune.epub", "Dune");
}