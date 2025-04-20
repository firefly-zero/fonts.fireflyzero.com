use std::{collections::HashMap, path::PathBuf};

use anyhow::{Context, Result};
use embedded_graphics::{
    image::ImageRaw,
    mono_font::{DecorationDimensions, MonoFont},
    prelude::*,
};

type Atlases = HashMap<&'static str, Vec<u8>>;

pub(crate) fn get_fonts(atlases: &Atlases) -> Vec<(&'static str, MonoFont<'_>)> {
    let ibm437b = MonoFont {
        image: ImageRaw::new(&atlases["ibm437b"], 128),
        character_size: Size::new(8, 8),
        character_spacing: 0,
        baseline: 7,
        strikethrough: DecorationDimensions::new(4, 1),
        underline: DecorationDimensions::new(8, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    vec![("ibm437b", ibm437b)]
}

pub(crate) fn load_atlases() -> Result<Atlases> {
    let mut res = HashMap::new();
    res.insert("ibm437b", load_atlas("ibm437b_8x8.png")?);
    Ok(res)
}

fn load_atlas(file_name: &'static str) -> Result<Vec<u8>> {
    let path = PathBuf::new().join("atlas").join(file_name);
    let file = image::ImageReader::open(path).context("open image file")?;
    let img = file.decode().context("decode image")?;
    let img = img.to_luma8();

    let mut raw = Vec::new();
    let mut byte: u8 = 0;
    for (i, color) in img.iter().enumerate() {
        byte = byte << 1 | u8::from(*color != 0);
        if i % 8 == 7 {
            raw.push(byte);
            byte = 0;
        }
    }
    Ok(raw)
}
