package main

import (
	"embed"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadImage(fs embed.FS, path string) *ebiten.Image {
	log.Printf("Loading %s", path)
	rawFile, err := fs.Open(path)
	if err != nil {
		log.Fatalf("Error opening file %s: %v\n", path, err)
	}
	defer rawFile.Close()

	img, err := png.Decode(rawFile)
	if err != nil {
		log.Fatalf("Error decoding file %s: %v\n", path, err)
	}
	loadedImg := ebiten.NewImageFromImage(img)
	return loadedImg
}

/*
func init() {
	fontFile, err := FileSystem.Open("fonts/mplus-1p-regular.ttf")
	if err != nil {
		log.Fatalf("Error opening font: %v\n", err)
	}
	tt, err := opentype.Parse(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
*/
