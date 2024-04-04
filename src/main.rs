#![feature(iter_array_chunks)]

mod fonts;
mod generator;
use crate::generator::save_all_fonts;
use std::path::Path;

fn main() {
    let root = Path::new("../fonts");
    save_all_fonts(root).unwrap();
}
