package main

import (
	"embed"
	"fmt"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadImage(fs embed.FS, path string) *ebiten.Image {
	fmt.Println("loadImage logic goes here")
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
