package generator

import (
	"cmp"
	"embed"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"os"
	"slices"
	"strings"

	"github.com/firefly-zero/fonts.fireflyzero.com/generator/bdf"
)

var (
	//go:embed templates/*
	templatesFiles embed.FS
	tmplTemplates  = map[string]*template.Template{
		"index":    parseTemplate("index"),
		"encoding": parseTemplate("encoding"),
	}
)

func parseTemplate(name string) *template.Template {
	return template.Must(template.New("root").ParseFS(
		templatesFiles,
		"templates/base.html.tmpl",
		"templates/card.html.tmpl",
		"templates/"+name+".html.tmpl",
	))
}

type Font struct {
	Name   string
	Family string
	Style  string
	Width  int
	Height int
}

func Run() error {
	err := renderEncodingsPage()
	if err != nil {
		return err
	}

	files, err := os.ReadDir("bdf")
	if err != nil {
		return fmt.Errorf("open bdf dir: %v", err)
	}
	fonts := make([]Font, 0, len(files))
	for _, file := range files {
		name := file.Name()
		if name == "HaxorNarrow-17.bdf" {
			continue
		}
		font, err := handleBDFFile("bdf/" + name)
		if err != nil {
			return err
		}
		fonts = append(fonts, *font)
	}

	slices.SortFunc(fonts, cmpFonts)

	err = renderEncodingPage(encodings[0], fonts)
	if err != nil {
		return err
	}
	return nil
}

func cmpFonts(a, b Font) int {
	if a.Width != b.Width {
		return cmp.Compare(a.Width, b.Width)
	}
	if a.Height != b.Height {
		return cmp.Compare(a.Height, b.Height)
	}
	return cmp.Compare(a.Family, b.Family)
}

func renderEncodingsPage() error {
	page := "index.html"
	file, err := os.Create("public/" + page)
	if err != nil {
		return fmt.Errorf("create %s: %v", page, err)
	}
	defer func() { _ = file.Close() }()
	err = tmplTemplates["index"].ExecuteTemplate(file, page+".tmpl", encodings)
	if err != nil {
		return fmt.Errorf("render %s: %v", page, err)
	}
	return nil
}

func renderEncodingPage(enc Encoding, fonts []Font) error {
	outName := fmt.Sprintf("public/%s/index.html", enc.Slug)
	file, err := os.Create(outName)
	if err != nil {
		return fmt.Errorf("create %s: %v", outName, err)
	}
	defer func() { _ = file.Close() }()
	data := map[string]any{
		"Encoding": enc,
		"Fonts":    fonts,
	}
	err = tmplTemplates["encoding"].ExecuteTemplate(file, "encoding.html.tmpl", data)
	if err != nil {
		return fmt.Errorf("render %s: %v", outName, err)
	}
	return nil
}

func handleBDFFile(fname string) (*Font, error) {
	raw, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("read BDF file: %v", err)
	}

	font, err := bdf.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse BDF font: %v", err)
	}
	if font.Foundry == "Raccoon" {
		font.BBox.Width++
	}

	img := makeAtlas(font)

	family := getFamily(font)
	style := getStyle(font)
	fontName := fmt.Sprintf("%s-%dx%d%s", family, font.BBox.Width, font.BBox.Height, style)

	outName := "public/ascii/" + fontName + ".png"
	_, err = os.Stat(outName)
	if err == nil {
		return nil, fmt.Errorf("file already exists: %s", outName)
	}
	out, err := os.Create(outName)
	if err != nil {
		return nil, fmt.Errorf("open %s on write: %v", outName, err)
	}
	defer func() { _ = out.Close() }()
	err = png.Encode(out, img)
	if err != nil {
		return nil, fmt.Errorf("save PNG file: %v", err)
	}

	info := Font{
		Name:   fontName,
		Family: family,
		Style:  style,
		Width:  font.BBox.Width,
		Height: font.BBox.Height,
	}
	return &info, nil
}

func makeAtlas(font *bdf.Font) *image.Alpha {
	chars := []bdf.Character{}
	for r := '\x20'; r <= '\x7F'; r++ {
		ch := getChar(font, r)
		chars = append(chars, *ch)
	}

	charW := font.BBox.Width
	charH := font.BBox.Height
	imgW := charW * 16
	imgH := charH * len(chars) / 16
	img := image.NewAlpha(image.Rect(0, 0, imgW, imgH))

	for i, ch := range chars {
		src := ch.Alpha
		rect := image.Rect(
			charW*(i%16),
			charH*(i/16)+(charH-ch.Alpha.Bounds().Dy()),
			imgW,
			imgH,
		)
		draw.Draw(img, rect, src, src.Rect.Min, draw.Src)
	}
	return img
}

func getChar(font *bdf.Font, r rune) *bdf.Character {
	ch := font.CharMap[r]
	if ch != nil {
		return ch
	}
	ch = font.CharMap[font.DefaultChar]
	if ch != nil {
		return ch
	}
	ch = font.CharMap['?']
	if ch != nil {
		return ch
	}
	panic("character not found")
}

func getFamily(font *bdf.Font) string {
	if font.Encoding == "ibm-437" {
		return "ibm-437"
	}
	name := font.FamilyName
	name = strings.Split(name, " ")[0]
	name = strings.ToLower(name)
	if name == "fixed" {
		name = "ucs"
	}
	if name == "fixed4x6" {
		name = "raccoon"
	}
	return name
}

func getStyle(font *bdf.Font) string {
	if strings.ToLower(font.Weight) == "bold" {
		return "b"
	}
	if font.FamilyName == "Haxor" && font.Weight == "Medium" {
		return "b"
	}
	if font.FamilyName == "CGA" && !strings.Contains(font.Name, "Light") {
		return "b"
	}

	slant := strings.ToUpper(font.Slant)
	if slant == "I" || slant == "O" {
		return "i"
	}
	return ""
}
