package generator

import (
	"embed"
	"encoding/binary"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/firefly-zero/fonts.fireflyzero.com/generator/bdf"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed templates/*
	templatesFiles embed.FS
	tmplTemplates  = map[string]*template.Template{
		"index":    parseTemplate("index"),
		"encoding": parseTemplate("encoding"),
		"family":   parseTemplate("family"),
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
	Size      string
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

	// Spawn background task for each font.
	fontsChan := make(chan Font)
	group := &errgroup.Group{}
	for _, file := range files {
		name := file.Name()
		if name == "HaxorNarrow-17.bdf" {
			continue
		}
		group.Go(func() error {
			font, err := handleBDFFile("bdf/" + name)
			if err != nil {
				return err
			}
			fontsChan <- *font
			return nil
		})
	}

	// Collect result of all tasks.
	var groupErr error
	go func() {
		groupErr = group.Wait()
		close(fontsChan)
	}()
	fonts := make([]Font, 0, len(files))
	for font := range fontsChan {
		fonts = append(fonts, font)
	}
	if groupErr != nil {
		return groupErr
	}

	// Render HTML pages for the generated fonts.
	slices.SortFunc(fonts, cmpFonts)
	for _, enc := range encodings {
		encFonts := filterEncFonts(enc, fonts)
		err = renderEncodingPage(enc, encFonts)
		if err != nil {
			return err
		}
		for _, fam := range families {
			famFonts := filterFamFonts(fam, encFonts)
			err = renderFamilyPage(enc, fam, famFonts)
			if err != nil {
				return err
			}
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
		Size:      formatSize(len(assets.Firefly)),
		Encodings: found,
	}
	return &info, nil
}

func formatSize(n int) string {
	if n < 2048 {
		return fmt.Sprintf("%d bytes", n)
	}
	return fmt.Sprintf("%.0f Kb", math.Round(float64(n/1024)))
}

func saveAssets(enc Encoding, fontName string, assets Assets) error {
	dirName := "public/" + enc.Slug
	err := savePNG(dirName+"/"+fontName+"-atlas.png", assets.Atlas)
	if err != nil {
		return fmt.Errorf("save atlas: %v", err)
	}
	err = savePNG(dirName+"/"+fontName+"-pangram.png", assets.Pangram)
	if err != nil {
		return fmt.Errorf("save pangram: %v", err)
	}
	err = os.WriteFile(dirName+"/"+fontName+".fff", assets.Firefly, os.ModePerm)
	if err != nil {
		return fmt.Errorf("save firefly font: %v", err)
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
	Firefly []byte
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
	badChar := chars[len(chars)-1]
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
	if font.Foundry == "Raccoon" || strings.HasPrefix(font.FamilyName, "Galmuri") {
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
		ch := getCharFrom(chars, r)
		if ch == nil {
			ch = &badChar
		}
		drawCharAt(pg, *ch, x, (charH - ch.Alpha.Bounds().Dy()))
		x += charW
	}

	return &Assets{
		Atlas:   img,
		Pangram: pg.SubImage(image.Rect(0, 0, x, pgH)),
		Firefly: makeFirefly(enc, chars),
		CharW:   charW,
		CharH:   charH,
	}
}

// Generate Firefly Zero font.
func makeFirefly(enc Encoding, chars []bdf.Character) []byte {
	charW := 0
	charH := 0
	for _, ch := range chars {
		charW = max(charW, ch.Alpha.Rect.Dx())
		charH = max(charH, ch.Alpha.Rect.Dy())
	}
	nChars := len(chars)
	if nChars%8 != 0 {
		panic("invalid char count")
	}

	buf := []byte{
		0x12,
		enc.ID,
		byte(charW),
		byte(charH),
	}
	for _, ch := range chars {
		buf = append(buf, byte(int8(ch.LowerPoint.X)))
		buf = append(buf, byte(int8(ch.LowerPoint.Y)))
		buf = append(buf, byte(ch.Advance.X))
		buf = append(buf, encodeChar(ch, charW, charH)...)
	}
	return buf
}

func encodeChar(ch bdf.Character, charW, charH int) []byte {
	buf := make([]byte, 0)
	dy := charH - ch.Alpha.Rect.Dy()
	bits := 0
	var b byte
	for y := range charH {
		for x := range charW {
			b <<= 1
			if ch.Alpha.AlphaAt(x, y-dy).A != 0 {
				b |= 1
			}
			bits++
			if bits == 8 {
				buf = append(buf, b)
				b = 0
				bits = 0
			}
		}
	}
	if bits != 0 {
		b <<= 8 - bits
		buf = append(buf, b)
	}

	expSize := divCeil(charW*charH, 8)
	if len(buf) != expSize {
		panic(fmt.Sprintf("%d != %d (%dx%d bits)", len(buf), expSize, charW, charH))
	}
	return buf
}

func divCeil(lhs, rhs int) int {
	d := lhs / rhs
	r := lhs % rhs
	if r > 0 {
		return d + 1
	}
	return d
}

func drawCharAt(img *image.Alpha, ch bdf.Character, x, y int) {
	src := ch.Alpha
	rect := image.Rect(
		x+ch.LowerPoint.X,
		y-ch.LowerPoint.Y,
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

	// Some encodings (Tahoma variants) are missing a space character.
	// FOr such encodings, we can create it as an empty space with dimensions of "?".
	if r == ' ' {
		ch := font.CharMap['?']
		return &bdf.Character{
			Name:       "space",
			Encoding:   ' ',
			Advance:    ch.Advance,
			Alpha:      image.NewAlpha(ch.Alpha.Rect),
			LowerPoint: ch.LowerPoint,
		}
	}

	ch = font.CharMap[0x7f]
	if ch != nil {
		return ch
	}
	ch = font.CharMap['?']
	if ch != nil {
		return ch
	}
	panic("character not found")
}

func getCharFrom(chars []bdf.Character, r rune) *bdf.Character {
	for _, ch := range chars {
		if ch.Encoding == r {
			return &ch
		}
	}
	return nil
}

func getFamily(font *bdf.Font) string {
	if font.Encoding == "ibm-437" {
		return "ibm-437"
	}
	name := font.FamilyName
	if strings.HasPrefix(name, "Galmuri") {
		return "galmuri"
	}
	name = strings.Split(name, " ")[0]
	name = strings.ToLower(name)
	if name == "fixed" {
		name = "ucs"
	}
	if name == "fixed4x6" {
		name = "tom-thumb"
	}
	if name == "unifontexmono" {
		name = "unifontex"
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
