mod extra_fonts;
mod fonts;
mod generator;
mod html;
mod image;
use crate::generator::save_all_fonts;
use crate::html::build_html;
use crate::image::fonts_to_images;
use anyhow::{Context, Result};
use std::path::PathBuf;

fn main() {
    if let Err(err) = run() {
        println!("{err:?}");
        std::process::exit(1);
    }
}

fn run() -> Result<()> {
    let root = PathBuf::new().join("public");
    std::fs::create_dir_all(&root).context("create public dir")?;
    let fonts_path = root.join("fonts");
    let images_path = root.join("images");

    std::fs::create_dir_all(&fonts_path).context("create fonts dir")?;
    let count = save_all_fonts(&fonts_path).context("generate fonts")?;
    println!("Generated {count} fonts");

    fonts_to_images(&fonts_path, &images_path).context("generate images")?;
    println!("Generated images");

    build_html(&root).context("generate HTML")?;
    println!("Generated HTML pages");
    Ok(())
}
