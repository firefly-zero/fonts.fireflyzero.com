mod fonts;
mod generator;
mod html;
use crate::generator::save_all_fonts;
use crate::html::build_html;
use std::path::PathBuf;

fn main() {
    let root = PathBuf::new().join("public");
    std::fs::create_dir_all(&root).unwrap();
    let fonts_path = root.join("fonts");
    std::fs::create_dir_all(&fonts_path).unwrap();
    let count = save_all_fonts(&fonts_path).unwrap();
    println!("Generated {count} fonts");
    build_html(&root).unwrap();
    println!("Generated HTML pages");
}
