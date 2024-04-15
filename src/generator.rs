use crate::fonts::FONTS;
use embedded_graphics::draw_target::DrawTarget;
use embedded_graphics::geometry::OriginDimensions;
use embedded_graphics::mono_font::MonoFont;
use embedded_graphics::pixelcolor::BinaryColor;
use embedded_graphics::prelude::*;
use std::fs::File;
use std::io::{self, Write};
use std::path::Path;

pub(crate) fn save_all_fonts(root: &Path) -> io::Result<usize> {
    let mut count = 0;
    for (encoding_index, (encoding_name, fonts)) in FONTS.iter().enumerate() {
        let dir_path = root.join(encoding_name);
        std::fs::create_dir_all(&dir_path)?;
        for font in fonts.iter() {
            let size = &font.character_size;
            let file_name = format!("eg_{}x{}.fff", size.width, size.height);
            let path = dir_path.join(file_name);
            dump_font(&path, encoding_index, font)?;
            count += 1
        }
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
        panic!("not implemented")
    }
}

impl DrawTarget for FileWrapper {
    type Color = BinaryColor;
    type Error = io::Error;

    fn draw_iter<I>(&mut self, _pixels: I) -> Result<(), Self::Error>
    where
        I: IntoIterator<Item = Pixel<Self::Color>>,
    {
        panic!("use fill_contiguous instead")
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
