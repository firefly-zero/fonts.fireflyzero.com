package generator

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/firefly-zero/fonts.fireflyzero.com/generator/bdf"
)

func Run() error {
	files, err := os.ReadDir("bdf")
	if err != nil {
		return fmt.Errorf("open bdf dir: %v", err)
	}
	for _, file := range files {
		name := file.Name()
		if name == "HaxorNarrow-17.bdf" {
			continue
		}
		err1 := handleBDFFile("bdf/" + name)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func handleBDFFile(fname string) error {
	raw, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("read BDF file: %v", err)
	}

	font, err := bdf.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse BDF font: %v", err)
	}

	img := makeAtlas(font)

	family := getFamily(font)
	suffix := getStyle(font)
	fontName := fmt.Sprintf("%s-%dx%d%s", family, font.BBox.Width, font.BBox.Height, suffix)

	outName := "public/mono/" + fontName + ".png"
	_, err = os.Stat(outName)
	if err == nil {
		return fmt.Errorf("file already exists: %s", outName)
	}
	out, err := os.Create(outName)
	if err != nil {
		return fmt.Errorf("open %s on write: %v", outName, err)
	}
	defer func() { _ = out.Close() }()
	err = png.Encode(out, img)
	if err != nil {
		return fmt.Errorf("save PNG file: %v", err)
	}
	return nil
}

func makeAtlas(font *bdf.Font) *image.Alpha {
	chars := []bdf.Character{}
	charW := 0
	charH := 0
	for r := '\x20'; r <= '\x7F'; r++ {
		ch := getChar(font, r)
		chars = append(chars, *ch)
		bbox := ch.Alpha.Rect
		charW = max(charW, bbox.Dx())
		charH = max(charH, bbox.Dy())
	}

	imgW := charW * 16
	imgH := charH * len(chars) / 16
	img := image.NewAlpha(image.Rect(0, 0, imgW, imgH))

	for i, ch := range chars {
		src := ch.Alpha
		rect := image.Rect(
			charW*(i%16),
			charH*(i/16),
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
