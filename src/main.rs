use sha2::{Digest, Sha256};
use std::{fs::File, io};
use walkdir::WalkDir;

fn main() {
    for entry in get_files(".") {
        println!("{}", entry);
        get_hash(entry)
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

fn get_hash(filename: String) {
    let mut hasher = Sha256::new();
    let mut file = File::open(filename).unwrap();
    let bytes_written = io::copy(&mut file, &mut hasher).unwrap();
    let hash_bytes = hasher.finalize();

    // show only first 7 chars of hash

    println!("hash is: {:x}", hash_bytes);
}
