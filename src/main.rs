use walkdir::WalkDir;

fn main() {
    for entry in get_files(".") {
        println!("{}", entry);
    }
}

fn get_files(root: &str) -> Vec<String> {
    WalkDir::new(root)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.file_type().is_file())
        .map(|e| e.path().to_str().unwrap().to_string())
        .collect()
}
