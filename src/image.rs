use anyhow::{Context, Result};
use std::fs;
use std::io::BufWriter;
use std::path::Path;

pub(crate) fn fonts_to_images(in_dir: &Path, out_dir: &Path) -> Result<()> {
    fs::create_dir_all(out_dir).context("create images dir")?;
    let dirs = fs::read_dir(in_dir).context("read input dir")?;
    for subdir in dirs {
        let subdir = subdir.context("acces subdir")?;
        let files = fs::read_dir(subdir.path()).context("read input subdir")?;
        let out_subdir = out_dir.join(subdir.file_name());
        fs::create_dir_all(&out_subdir).context("create images subdir")?;
        for file in files {
            let file = file.context("access font file")?;
            let in_path = file.path();
            let out_path = out_subdir.join(file.file_name()).with_extension("png");
            font_to_image(&in_path, &out_path).context("convert font to image")?
        }
    }
    Ok(())
}

fn font_to_image(in_path: &Path, out_path: &Path) -> Result<()> {
    // read font data
    let raw_font = fs::read(in_path)?;
    let width = i16::from_le_bytes([raw_font[5], raw_font[6]]) as u32;
    let data = raw_font.get(7..).unwrap();
    let height = data.len() as u32 * 8 / width;

    // invert colors
    let mut inv_data = Vec::new();
    for byte in data {
        inv_data.push(!byte)
    }

    // write png
    let file = fs::File::create(out_path).context("create image file")?;
    let buffer = BufWriter::new(file);
    let mut encoder = png::Encoder::new(buffer, width as u32, height as u32);
    encoder.set_color(png::ColorType::Grayscale);
    encoder.set_depth(png::BitDepth::One);
    let mut writer = encoder.write_header().context("write PNG header")?;
    writer
        .write_image_data(&inv_data)
        .context("write image data")?;

    Ok(())
}
