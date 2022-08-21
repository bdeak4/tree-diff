use blake2::{Blake2s256, Digest};
use std::path::Path;
use std::{fs, io};

fn main() {
    println!("hi");
    println!("{}", get_hash(Path::new("./src/main.rs")).unwrap());
}

fn get_files() {
    println!("get files")
}

fn get_hash(path: &Path) -> Result<String, io::Error> {
    let mut file = fs::File::open(&path)?;
    let mut hasher = Blake2s256::new();
    io::copy(&mut file, &mut hasher)?;
    let hash = format!("{:x}", hasher.finalize());
    let hash_short = hash[..7].to_string();
    Ok(hash_short)
}

// TODO: walkdir
//
// fn main() -> io::Result<()> {
//     visit_dirs(".", |entry: DirEntry| {
//         let path = entry.path();
//         if path.is_file() {
//             println!("{}", path.display());
//         }
//     })?;

//     // The entries have now been sorted by their path.

//     Ok(())
// }

// fn visit_dirs(dir: &Path, cb: &dyn Fn(&DirEntry)) -> io::Result<()> {
//     if dir.is_dir() {
//         for entry in fs::read_dir(dir)? {
//             let entry = entry?;
//             let path = entry.path();
//             if path.is_dir() {
//                 visit_dirs(&path, cb)?;
//             } else {
//                 cb(&entry);
//             }
//         }
//     }
//     Ok(())
// }
