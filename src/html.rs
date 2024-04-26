use crate::fonts::FONTS;
use anyhow::{Context, Result};
use minijinja::{context, Environment};
use serde::Serialize;
use std::collections::HashSet;
use std::fs;
use std::path::Path;

#[derive(Serialize)]
struct Font {
    family:   &'static str,
    width:    u32,
    height:   u32,
    encoding: &'static str,
}

pub(crate) fn build_html(root: &Path) -> Result<()> {
    let mut env = Environment::new();
    let templates = fs::read_dir("templates")?;
    for file_path in templates {
        let file_path = file_path?;
        let content = fs::read_to_string(file_path.path())?;
        let file_name = file_path.file_name();
        let file_name = file_name.to_str().unwrap().to_string();
        env.add_template_owned(file_name, content)?;
    }

    let all_fonts = make_fonts();
    let mut encodings: Vec<&str> = all_fonts
        .iter()
        .map(|f| f.encoding)
        .collect::<HashSet<_>>()
        .into_iter()
        .collect();
    encodings.sort();

    {
        let out_path = root.join("index.html");
        let tmpl = env.get_template("index.html.j2").context("get template")?;
        let rendered = tmpl
            .render(context!(encodings => encodings))
            .context("render template")?;
        fs::write(out_path, rendered).context("write html file")?;
    }

    let tmpl = env
        .get_template("encoding.html.j2")
        .context("get template")?;
    for encoding in encodings {
        let fonts: Vec<_> = all_fonts
            .iter()
            .filter(|f| f.encoding == encoding)
            .collect();
        let out_path = root.join(format!("{encoding}.html"));
        let rendered = tmpl
            .render(context!(fonts => fonts))
            .context("render template")?;
        fs::write(out_path, rendered).context("write html file")?;
    }

    Ok(())
}

fn make_fonts() -> Vec<Font> {
    let mut result = Vec::new();
    for (family, encoding, fonts) in FONTS {
        for font in *fonts {
            result.push(Font {
                family,
                encoding,
                width: font.character_size.width,
                height: font.character_size.height,
            })
        }
    }
    result
}
