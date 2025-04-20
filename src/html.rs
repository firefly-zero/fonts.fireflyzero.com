use crate::extra_fonts::{get_fonts, load_atlases};
use crate::fonts::FONTS;
use anyhow::{Context, Result};
use minijinja::{context, Environment};
use serde::Serialize;
use std::fs;
use std::path::Path;

static ENCODINGS: &[Encoding] = &[
    Encoding {
        slug: "ascii",
        title: "ASCII",
        icon: "🇺🇸",
    },
    Encoding {
        slug: "iso_8859_1",
        title: "Latin-1, Western European.",
        icon: "🇵🇹",
    },
    Encoding {
        slug: "iso_8859_2",
        title: "Latin-2, Central European.",
        icon: "🇩🇪",
    },
    Encoding {
        slug: "iso_8859_3",
        title: "Latin-3, South European.",
        icon: "🇲🇹",
    },
    Encoding {
        slug: "iso_8859_4",
        title: "Latin-4, North European.",
        icon: "🇪🇪",
    },
    Encoding {
        slug: "iso_8859_9",
        title: "Latin-5, Turkish.",
        icon: "🇹🇷",
    },
    Encoding {
        slug: "iso_8859_10",
        title: "Latin-6, Nordic.",
        icon: "🇳🇴",
    },
    Encoding {
        slug: "iso_8859_13",
        title: "Latin-7, Baltic Rim.",
        icon: "🇵🇱",
    },
    Encoding {
        slug: "iso_8859_14",
        title: "Latin-8, Celtic.",
        icon: "🇮🇪",
    },
    Encoding {
        slug: "iso_8859_15",
        title: "Latin-9 (revised Latin-1).",
        icon: "🇵🇹",
    },
    Encoding {
        slug: "iso_8859_16",
        title: "Latin-10: South-East European.",
        icon: "🇷🇴",
    },
    Encoding {
        slug: "iso_8859_5",
        title: "Latin/Cyrillic.",
        icon: "🇷🇺",
    },
    Encoding {
        slug: "iso_8859_7",
        title: "Latin/Greek.",
        icon: "🇬🇷",
    },
    Encoding {
        slug: "jis_x0201",
        title: "Japanese katakana (halfwidth).",
        icon: "🇯🇵",
    },
];

#[derive(Serialize)]
struct Font {
    family: &'static str,
    width: u32,
    height: u32,
    encoding: &'static str,
}

#[derive(Serialize)]
struct Encoding {
    slug: &'static str,
    title: &'static str,
    icon: &'static str,
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

    {
        let out_path = root.join("index.html");
        let tmpl = env.get_template("index.html.j2").context("get template")?;
        let rendered = tmpl
            .render(context!(encodings => ENCODINGS))
            .context("render template")?;
        fs::write(out_path, rendered).context("write html file")?;
    }

    let tmpl = env
        .get_template("encoding.html.j2")
        .context("get template")?;
    for encoding in ENCODINGS {
        let fonts: Vec<_> = all_fonts
            .iter()
            .filter(|f| f.encoding == encoding.slug)
            .collect();
        let out_path = root.join(format!("{}.html", encoding.slug));
        let rendered = tmpl
            .render(context!(fonts => fonts, encoding => encoding))
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

    let atlases = load_atlases().unwrap();
    let fonts = get_fonts(&atlases);
    for (family, font) in fonts {
        result.push(Font {
            family,
            encoding: "ascii",
            width: font.character_size.width,
            height: font.character_size.height,
        })
    }
    result
}
