use std::{collections::HashMap, path::PathBuf};

use anyhow::{Context, Result};
use embedded_graphics::{
    image::ImageRaw,
    mono_font::{DecorationDimensions, MonoFont},
    prelude::*,
};
use serde::Serialize;

type Atlases = HashMap<&'static str, Vec<u8>>;

pub(crate) struct Font<'a> {
    pub family: &'static str,
    pub font: MonoFont<'a>,
    pub license: License,
}

#[derive(Copy, Clone, Serialize)]
pub(crate) struct License {
    /// https://spdx.org/licenses/
    pub spdx: &'static str,
    pub url: &'static str,
}

/// Convert atlases into embedded-graphics fonts.
pub(crate) fn get_fonts(atlases: &Atlases) -> Vec<Font<'_>> {
    let ibm437b = MonoFont {
        image: ImageRaw::new(&atlases["ibm437b"], 128),
        character_size: Size::new(8, 8),
        character_spacing: 0,
        baseline: 7,
        strikethrough: DecorationDimensions::new(4, 1),
        underline: DecorationDimensions::new(8, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let ibm437r = MonoFont {
        image: ImageRaw::new(&atlases["ibm437r"], 128),
        character_size: Size::new(8, 8),
        character_spacing: 0,
        baseline: 7,
        strikethrough: DecorationDimensions::new(4, 1),
        underline: DecorationDimensions::new(8, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let pico8 = MonoFont {
        image: ImageRaw::new(&atlases["pico8"], 64),
        character_size: Size::new(4, 6),
        character_spacing: 0,
        baseline: 5,
        strikethrough: DecorationDimensions::new(3, 1),
        underline: DecorationDimensions::new(6, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let kenney11 = MonoFont {
        image: ImageRaw::new(&atlases["kenney11"], 176),
        character_size: Size::new(11, 14),
        character_spacing: 0,
        baseline: 14,
        strikethrough: DecorationDimensions::new(7, 2),
        underline: DecorationDimensions::new(15, 2),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let kenney16 = MonoFont {
        image: ImageRaw::new(&atlases["kenney16"], 256),
        character_size: Size::new(16, 16),
        character_spacing: 0,
        baseline: 14,
        strikethrough: DecorationDimensions::new(8, 2),
        underline: DecorationDimensions::new(15, 2),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let profont = MonoFont {
        image: ImageRaw::new(&atlases["profont"], 80),
        character_size: Size::new(5, 10),
        character_spacing: 0,
        baseline: 9,
        strikethrough: DecorationDimensions::new(5, 1),
        underline: DecorationDimensions::new(10, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let mem = MonoFont {
        image: ImageRaw::new(&atlases["mem"], 64),
        character_size: Size::new(4, 4),
        character_spacing: 0,
        baseline: 3,
        strikethrough: DecorationDimensions::new(2, 1),
        underline: DecorationDimensions::new(4, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    let sixel = MonoFont {
        image: ImageRaw::new(&atlases["sixel"], 16),
        character_size: Size::new(1, 6),
        character_spacing: 0,
        baseline: 6,
        strikethrough: DecorationDimensions::new(3, 1),
        underline: DecorationDimensions::new(7, 1),
        glyph_mapping: &embedded_graphics::mono_font::mapping::ASCII,
    };
    vec![
        Font {
            family: "pico8",
            font: pico8,
            license: License {
                spdx: "CC0-1.0",
                url: "https://www.lexaloffle.com/pico-8.php?page=faq",
            },
        },
        Font {
            family: "profont",
            font: profont,
            license: License {
                spdx: "MIT",
                url: "https://tobiasjung.name/profont/",
            },
        },
        Font {
            family: "ibm437b",
            font: ibm437b,
            license: License {
                spdx: "MIT",
                url: "https://github.com/sbechet/ibm437/blob/master/LICENCE",
            },
        },
        Font {
            family: "ibm437r",
            font: ibm437r,
            license: License {
                spdx: "MIT",
                url: "https://github.com/sbechet/ibm437/blob/master/LICENCE",
            },
        },
        Font {
            family: "mem",
            font: mem,
            license: License {
                spdx: "AGPL-3.0",
                url: "https://github.com/oidoid/mem/blob/main/license.text",
            },
        },
        Font {
            family: "sixel",
            font: sixel,
            license: License {
                spdx: "MIT",
                url: "https://saitoha.github.io/libsixel/",
            },
        },
        Font {
            family: "kenney",
            font: kenney11,
            license: License {
                spdx: "CC0-1.0",
                url: "https://www.kenney.nl/assets/1-bit-pack",
            },
        },
        Font {
            family: "kenney",
            font: kenney16,
            license: License {
                spdx: "CC0-1.0",
                url: "https://www.kenney.nl/assets/1-bit-pack",
            },
        },
    ]
}

/// Load font atlases from filesystem.
pub(crate) fn load_atlases() -> Result<Atlases> {
    let mut res = HashMap::new();
    res.insert("ibm437b", load_atlas("ibm437b_8x8.png")?);
    res.insert("ibm437r", load_atlas("ibm437r_8x8.png")?);
    res.insert("pico8", load_atlas("pico8_4x6.png")?);
    res.insert("profont", load_atlas("profont_5x10.png")?);
    res.insert("mem", load_atlas("mem_4x4.png")?);
    res.insert("sixel", load_atlas("sixel_1x6.png")?);
    res.insert("kenney11", load_atlas("kenney_11x14.png")?);
    res.insert("kenney16", load_atlas("kenney_16x16.png")?);
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
