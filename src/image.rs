use anyhow::{Context, Result};
use image::codecs::png::{CompressionType, FilterType, PngEncoder};
use image::{ExtendedColorType, ImageEncoder};
use std::fs;
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
            let out_path = out_subdir.join(file.file_name()).with_extension(".png");
            font_to_image(&in_path, &out_path).context("convert font to image")?
        }
    }
    Ok(())
}

fn font_to_image(in_path: &Path, out_path: &Path) -> Result<()> {
    let mut output_buffer = Vec::new();
    let encoder = PngEncoder::new_with_quality(
        &mut output_buffer, // buffer
        CompressionType::Best,
        FilterType::NoFilter,
    );
    let raw_font = fs::read(in_path)?;
    let width = i16::from_le_bytes([raw_font[5], raw_font[6]]) as u32;
    let data = raw_font.get(7..).unwrap();
    let height = data.len() as u32 * 8 / width;
    let data = convert_bpp(data);
    encoder.write_image(&data, width, height, ExtendedColorType::L8)?;
    fs::write(out_path, output_buffer)?;
    Ok(())
}

fn convert_bpp(input: &[u8]) -> Vec<u8> {
    let mut output = Vec::new();
    for byte in input {
        for shift in (0..8).rev() {
            let new_byte = (byte >> shift) & 0b1;
            let new_byte = if new_byte != 0 { 0 } else { 255 };
            output.push(new_byte);
        }
    }
    output
}
