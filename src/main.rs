mod fonts;
mod generator;
use crate::generator::save_all_fonts;
use std::path::Path;

fn main() {
    let args: Vec<_> = std::env::args().collect();
    let default = "fonts".to_string();
    let root = args.get(1).unwrap_or(&default);
    let root = Path::new(root);
    let count = save_all_fonts(root).unwrap();
    println!("Generated {count} fonts")
}
