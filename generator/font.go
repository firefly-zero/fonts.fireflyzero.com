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

func handleBDFFile(fname string) error {
	raw, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("read BDF file: %v", err)
	}

	font, err := bdf.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse BDF font: %v", err)
	}

	chars := []bdf.Character{}
	charW := 0
	charH := 0
	for r := '\x20'; r <= '\x7F'; r++ {
		ch := font.CharMap[r]
		if ch == nil {
			ch = font.CharMap[font.DefaultChar]
		}
		chars = append(chars, *ch)
		bbox := ch.Alpha.Rect
		charW = max(charW, bbox.Dx())
		charH = max(charH, bbox.Dy())
	}

	imgW := charW * 16
	imgH := charW * 32
	img := image.NewAlpha(image.Rect(0, 0, imgW, imgH))

	for i, ch := range chars {
		src := ch.Alpha
		rect := image.Rect(
			charW*(i%16),
			charH*(i/16),
			src.Rect.Dx(),
			src.Rect.Dy(),
		)
		draw.Draw(img, rect, src, image.Point{}, draw.Src)
	}

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
