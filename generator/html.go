package generator

import (
	"cmp"
	"fmt"
	"os"
	"slices"
)

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

func filterEncFonts(enc Encoding, fonts []Font) []Font {
	fonts = slices.Clone(fonts)
	return slices.DeleteFunc(fonts, func(f Font) bool {
		return !f.Encodings[enc.Slug]
	})
}

func filterFamFonts(fam Family, fonts []Font) []Font {
	fonts = slices.Clone(fonts)
	return slices.DeleteFunc(fonts, func(f Font) bool {
		return f.Family.Name != fam.Name
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

func renderFamilyPage(enc Encoding, fam Family, fonts []Font) error {
	outName := fmt.Sprintf("public/%s/%s.html", enc.Slug, fam.Name)
	file, err := os.Create(outName)
	if err != nil {
		return fmt.Errorf("create %s: %v", outName, err)
	}
	defer func() { _ = file.Close() }()
	data := map[string]any{
		"Encoding": enc,
		"Family":   fam,
		"Fonts":    fonts,
	}
	err = tmplTemplates["family"].ExecuteTemplate(file, "family.html.tmpl", data)
	if err != nil {
		return fmt.Errorf("render %s: %v", outName, err)
	}
	return nil
}
