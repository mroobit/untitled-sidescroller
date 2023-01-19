package main

import (
	"embed"
	"encoding/json"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadAssets() {
	world = loadImage(FileSystem, "imgs/world--test.png")
	gooAlley = loadImage(FileSystem, "imgs/goo-alley--test.png")
	yikesfulMountain = loadImage(FileSystem, "imgs/yikesful-mountain--test.png")
	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")

	ebitengineSplash = loadImage(FileSystem, "imgs/load-ebitengine-splash.png")

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	portal = loadImage(FileSystem, "imgs/portal-b--test.png")
	treasure = loadImage(FileSystem, "imgs/treasure--test.png")
	questItem = loadImage(FileSystem, "imgs/quest-item--test.png")
	hazard = loadImage(FileSystem, "imgs/blob--test.png")
	creature = loadImage(FileSystem, "imgs/creature--test.png")
	blank = loadImage(FileSystem, "imgs/blank-bg.png")
	gameOverMessage = loadImage(FileSystem, "imgs/game-over.png")

	levelImages := map[string][]*ebiten.Image{
		"Goo Alley":         {gooAlley, levelBG},
		"Yikesful Mountain": {yikesfulMountain, levelBG},
	}

	lvlContent, err := ioutil.ReadFile("./levels.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(lvlContent, &levelData)
	if err != nil {
		log.Fatal("Error during Unmarshalling: ", err)
	}

	for _, l := range levelData {
		l.icon = levelImages[l.Name][0]
		l.background = levelImages[l.Name][1]
	}

}

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
