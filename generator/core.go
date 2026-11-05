package generator

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"

	"github.com/firefly-zero/fonts.fireflyzero.com/generator/bdf"
)

func Run() error {
	files, err := os.ReadDir("bdf")
	if err != nil {
		return fmt.Errorf("open bdf dir: %v", err)
	}
	for _, file := range files {
		err1 := handleBDFFile("bdf/" + file.Name())
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

	outName := "public/atlases/" + path.Base(fname) + ".png"
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
