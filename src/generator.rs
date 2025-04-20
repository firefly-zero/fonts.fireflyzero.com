use crate::extra_fonts::{get_fonts, load_atlases};
use crate::fonts::FONTS;
use anyhow::{Context, Result};
use embedded_graphics::draw_target::DrawTarget;
use embedded_graphics::geometry::OriginDimensions;
use embedded_graphics::mono_font::MonoFont;
use embedded_graphics::pixelcolor::BinaryColor;
use embedded_graphics::prelude::*;
use std::fs::File;
use std::io::{self, Write};
use std::path::Path;

static ENCODINGS: &[&str] = &[
    "ascii",       // 0. ASCII
    "iso_8859_1",  // 1. Latin-1, Western European.
    "iso_8859_2",  // 2. Latin-2, Central European.
    "iso_8859_3",  // 3. Latin-3, South European.
    "iso_8859_4",  // 4. Latin-4, North European.
    "iso_8859_9",  // 5. Latin-5, Turkish.
    "iso_8859_10", // 6. Latin-6, Nordic.
    "iso_8859_13", // 7. Latin-7, Baltic Rim.
    "iso_8859_14", // 8. Latin-8, Celtic.
    "iso_8859_15", // 9. Latin-9 (revised Latin-1).
    "iso_8859_16", // A. Latin-10: South-East European.
    "iso_8859_5",  // B. Latin/Cyrillic.
    "iso_8859_7",  // C. Latin/Greek.
    "jis_x0201",   // D. Japanese katakana (halfwidth).
];

pub(crate) fn save_all_fonts(root: &Path) -> Result<usize> {
    let mut count = 0;
    for (family_name, encoding_name, fonts) in FONTS.iter() {
        let (encoding_index, _) = ENCODINGS
            .iter()
            .enumerate()
            .find(|(_, e)| *e == encoding_name)
            .unwrap();
        let dir_path = root.join(encoding_name);
        std::fs::create_dir_all(&dir_path).context("create encoding dir")?;
        for font in fonts.iter() {
            let size = &font.character_size;
            let file_name = format!("{family_name}_{}x{}.fff", size.width, size.height);
            let path = dir_path.join(file_name);
            dump_font(&path, encoding_index, font).context("dump font")?;
            count += 1
        }
    }

    let atlases = load_atlases().context("load atlases")?;
    let fonts = get_fonts(&atlases);
    let dir_path = root.join("ascii");
    let encoding_index = 0;
    for (family_name, font) in fonts {
        let size = &font.character_size;
        let file_name = format!("{family_name}_{}x{}.fff", size.width, size.height);
        let path = dir_path.join(file_name);
        dump_font(&path, encoding_index, &font).context("dump font")?;
        count += 1
    }

    Ok(count)
}

fn dump_font(path: &Path, encoding_index: usize, font: &MonoFont) -> io::Result<()> {
    let mut file = std::fs::File::create(path)?;
    let f = &mut file;
    write_u8(f, 0x11)?;
    write_u8(f, encoding_index as u8)?;
    write_u8(f, font.character_size.width as u8)?;
    write_u8(f, font.character_size.height as u8)?;
    write_u8(f, font.baseline as u8)?;
    write_u16(f, font.image.size().width as u16)?;
    let mut target = FileWrapper { file };
    font.image.draw(&mut target)
}

fn write_u8(f: &mut File, v: u8) -> io::Result<()> {
    f.write_all(&v.to_le_bytes())
}

fn write_u16(f: &mut File, v: u16) -> io::Result<()> {
    f.write_all(&v.to_le_bytes())
}

struct FileWrapper {
    file: File,
}

impl OriginDimensions for FileWrapper {
    fn size(&self) -> Size {
        unimplemented!("not implemented")
    }
}

impl DrawTarget for FileWrapper {
    type Color = BinaryColor;
    type Error = io::Error;

    fn draw_iter<I>(&mut self, _pixels: I) -> Result<(), Self::Error>
    where
        I: IntoIterator<Item = Pixel<Self::Color>>,
    {
        unimplemented!("use fill_contiguous instead")
    }

    fn fill_contiguous<I>(
        &mut self,
        _area: &embedded_graphics::primitives::Rectangle,
        colors: I,
    ) -> Result<(), Self::Error>
    where
        I: IntoIterator<Item = Self::Color>,
    {
        let mut raw: u8 = 0;
        for (i, color) in colors.into_iter().enumerate() {
            raw = raw << 1 | color.into_storage();
            if i % 8 == 7 {
                self.file.write_all(&[raw])?;
                raw = 0;
            }
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use embedded_graphics::mono_font::ascii;

    const HEADER_SIZE: usize = 7;

    #[test]
    fn test_save_all_fonts() {
        let path = std::env::temp_dir().join("test_save_all_fonts");
        save_all_fonts(&path).unwrap();
        let iter = std::fs::read_dir(&path).unwrap();
        // 14 encodings, a separate dir for each encoding
        assert_eq!(iter.count(), 14);
    }

    #[test]
    fn test_dump_font() {
        let font = ascii::FONT_5X7;
        let path = std::env::temp_dir().join("test_dump_font");
        dump_font(&path, 0, &font).unwrap();
        let dumped = std::fs::read(&path).unwrap();
        assert_eq!(dumped.len(), 420 + HEADER_SIZE);
        assert_eq!(dumped[0], u8::to_le_bytes(0x11)[0]);
    }
}
