package generator

import (
	"cmp"
	"embed"
	"encoding/binary"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

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
	Name      string
	Family    Family
	Style     string
	Width     int
	Height    int
	Encodings map[string]bool
}

func Run() error {
	err := loadEncodings()
	if err != nil {
		return err
	}
	err = renderEncodingsPage()
	if err != nil {
		return err
	}
	for _, enc := range encodings {
		dirName := "public/" + enc.Slug
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create dir: %v", err)
		}
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

	for _, enc := range encodings {
		err = renderEncodingPage(enc, filterFonts(enc, fonts))
		if err != nil {
			return err
		}
	}
	return nil
}

func loadEncodings() error {
	for i := range encodings {
		enc := &encodings[i]
		if enc.Slug == "ascii" {
			continue
		}
		chars := []rune{}
		raw, err := os.ReadFile("generator/encodings/" + enc.Slug + ".bin")
		if err != nil {
			return err
		}
		for i := 0; i < len(raw); i += 4 {
			ch := rune(binary.BigEndian.Uint32(raw[i : i+4]))
			chars = append(chars, ch)
		}
		enc.Chars = chars
	}
	return nil
}

func cmpFonts(a, b Font) int {
	apx := a.Width * a.Height
	bpx := b.Width * b.Height
	if apx != bpx {
		return cmp.Compare(apx, bpx)
	}
	return cmp.Compare(a.Family.Name, b.Family.Name)
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

func filterFonts(enc Encoding, fonts []Font) []Font {
	fonts = slices.Clone(fonts)
	return slices.DeleteFunc(fonts, func(f Font) bool {
		return !f.Encodings[enc.Slug]
	})
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

	assets := makeAssets(encodings[0], font)

	family := getFamily(font)
	style := getStyle(font)
	fontName := fmt.Sprintf("%s-%dx%d%s", family, assets.CharW, assets.CharH, style)

	err = saveAssets(encodings[0], fontName, *assets)
	if err != nil {
		return nil, err
	}

	found := map[string]bool{"ascii": true}
	for _, enc := range encodings[1:] {
		assets := makeAssets(enc, font)
		if assets == nil {
			continue
		}
		found[enc.Slug] = true
		err = saveAssets(enc, fontName, *assets)
		if err != nil {
			return nil, err
		}
	}

	info := Font{
		Name:      fontName,
		Family:    families[family],
		Style:     style,
		Width:     assets.CharW,
		Height:    assets.CharH,
		Encodings: found,
	}
	return &info, nil
}

func saveAssets(enc Encoding, fontName string, assets Assets) error {
	dirName := "public/" + enc.Slug
	err := savePNG(dirName+"/"+fontName+"-atlas.png", assets.Atlas)
	if err != nil {
		return err
	}
	err = savePNG(dirName+"/"+fontName+"-pangram.png", assets.Pangram)
	if err != nil {
		return err
	}
	return nil
}

func savePNG(outPath string, img image.Image) error {
	_, err := os.Stat(outPath)
	if err == nil {
		return fmt.Errorf("file already exists: %s", outPath)
	}
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("open %s on write: %v", outPath, err)
	}
	defer func() { _ = out.Close() }()
	err = png.Encode(out, img)
	if err != nil {
		return fmt.Errorf("save PNG file: %v", err)
	}
	return nil
}

type Assets struct {
	Atlas   image.Image
	Pangram image.Image
	CharW   int
	CharH   int
}

func makeAssets(enc Encoding, font *bdf.Font) *Assets {
	chars := []bdf.Character{}
	charW := 0
	charH := 0
	for r := '\x20'; r <= '\x7F'; r++ {
		ch := getChar(font, r)
		chars = append(chars, *ch)
		charW = max(charW, ch.Alpha.Rect.Dx())
		charH = max(charH, ch.Alpha.Rect.Dy())
	}
	nFound := 0
	for _, r := range enc.Chars {
		ch := getChar(font, r)
		if r == ch.Encoding {
			nFound++
		}
		chars = append(chars, *ch)
		charW = max(charW, ch.Alpha.Rect.Dx())
		charH = max(charH, ch.Alpha.Rect.Dy())
	}
	if len(enc.Chars) > 10 && float32(nFound)/float32(len(enc.Chars)) < .5 {
		return nil
	}

	if charW >= font.BBox.Width-2 {
		charW = font.BBox.Width
	}
	if charH >= font.BBox.Height-2 {
		charH = font.BBox.Height
	}
	if font.Foundry == "Raccoon" {
		charW++
	}

	// Make atlas: PNG image with every font character in order.
	imgW := charW*16 + 2
	imgH := charH*len(chars)/16 + font.Descent + 2
	img := image.NewAlpha(image.Rect(0, 0, imgW, imgH))
	for i, ch := range chars {
		x := charW * (i % 16)
		y := charH*(i/16) + (charH - ch.Alpha.Bounds().Dy())
		drawCharAt(img, ch, x, y)
	}

	// Render pangram: PNG image with a real sentence (using as many letters as possible).
	pgW := utf8.RuneCountInString(enc.Pangram)*charW + 4
	pgH := charH + font.Descent + 2
	pg := image.NewAlpha(image.Rect(0, 0, pgW, pgH))
	x := 2
	for _, r := range enc.Pangram {
		ch := getChar(font, r)
		drawCharAt(pg, *ch, x, (charH - ch.Alpha.Bounds().Dy()))
		x += ch.Advance[0]
	}

	return &Assets{
		Atlas:   img,
		Pangram: pg.SubImage(image.Rect(0, 0, x, pgH)),
		CharW:   charW,
		CharH:   charH,
	}
}

func drawCharAt(img *image.Alpha, ch bdf.Character, x, y int) {
	src := ch.Alpha
	ox := max(ch.LowerPoint[0], 0)
	rect := image.Rect(
		x+ox,
		y-ch.LowerPoint[1],
		img.Rect.Dx(),
		img.Rect.Dy(),
	)
	draw.Draw(img, rect, src, src.Rect.Min, draw.Src)
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
