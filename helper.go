package main

import (
	"embed"
	"encoding/json"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

var (
	fontLib *etxt.FontLibrary

	menuColorActive   = color.RGBA{140, 50, 90, 255}
	menuColorInactive = color.RGBA{0xff, 0xff, 0xff, 255}
	scoreDisplayColor = color.RGBA{0, 0, 0, 255}

	textColor color.RGBA
)

func init() {

	log.Printf("Creating new font library")
	fontLib = etxt.NewFontLibrary()

	_, _, err := fontLib.ParseEmbedDirFonts("fonts", FileSystem)
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}
}

func newRenderer() *etxt.Renderer {
	log.Printf("Creating new text renderer")
	renderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	renderer.SetCacheHandler(glyphsCache.NewHandler())
	renderer.SetFont(fontLib.GetFont("Consola Mono Bold"))
	renderer.SetAlign(etxt.YCenter, etxt.XCenter)
	renderer.SetSizePx(32)
	return renderer
}

func loadAssets() {
	world = loadImage(FileSystem, "imgs/world--test.png")
	gooAlley = loadImage(FileSystem, "imgs/goo-alley--test.png")
	yikesfulMountain = loadImage(FileSystem, "imgs/yikesful-mountain--test.png")
	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")
	backgroundYikesfulMountain = loadImage(FileSystem, "imgs/level-background-2--test.png")

	ebitengineSplash = loadImage(FileSystem, "imgs/load-ebitengine-splash.png")

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")

	gemCt = loadImage(FileSystem, "imgs/gem-count-large.png")
	livesCt = loadImage(FileSystem, "imgs/lives-left.png")
	messageBox = loadImage(FileSystem, "imgs/message-box.png")
	statsBox = loadImage(FileSystem, "imgs/stats-box.png")

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	portal = loadImage(FileSystem, "imgs/portal-b--test.png")
	treasure = loadImage(FileSystem, "imgs/treasure--test.png")
	portalGem = loadImage(FileSystem, "imgs/quest-item--test.png")
	hazard = loadImage(FileSystem, "imgs/blob--test.png")
	creature = loadImage(FileSystem, "imgs/creature--test.png")
	blank = loadImage(FileSystem, "imgs/blank-bg.png")
	gameOverMessage = loadImage(FileSystem, "imgs/game-over.png")

	levelImages := map[string][]*ebiten.Image{
		"Goo Alley":         {gooAlley, levelBG},
		"Yikesful Mountain": {yikesfulMountain, backgroundYikesfulMountain},
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
