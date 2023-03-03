package main

import (
	"embed"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

var (
	fontLib *etxt.FontLibrary

	menuColorActive   = color.RGBA{140, 50, 90, 255}
	menuColorInactive = color.RGBA{0xff, 0xff, 0xff, 255}
	menuColorDisabled = color.RGBA{60, 60, 60, 255}
	scoreDisplayColor = color.RGBA{0, 0, 0, 255}

	textColor color.RGBA
)

var (
	loaded = false

	ebitengineSplash *ebiten.Image
	splashImages     []*ebiten.Image
	levelImages      map[string][]*ebiten.Image

	gemCt      *ebiten.Image
	livesCt    *ebiten.Image
	messageBox *ebiten.Image
	statsBox   *ebiten.Image

	world *ebiten.Image

	gameOverMessage *ebiten.Image
)

func loadFonts() {
	log.Printf("Creating new font library")
	fontLib = etxt.NewFontLibrary()

	_, _, err := fontLib.ParseEmbedDirFonts("fonts", FileSystem)
	if err != nil {
		log.Fatalf("Error while loading fonts: %s", err.Error())
	}
}

func newRenderer() *etxt.Renderer {
	log.Printf("...creating new text renderer")
	renderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	renderer.SetCacheHandler(glyphsCache.NewHandler())
	renderer.SetFont(fontLib.GetFont("Consola Mono Bold"))
	renderer.SetAlign(etxt.YCenter, etxt.XCenter)
	renderer.SetSizePx(32)
	return renderer
}

func loadAssets() {
	log.Printf("Loading Images...")
	world = loadImage(FileSystem, "imgs/world--test.png")
	gooAlley = loadImage(FileSystem, "imgs/goo-alley--test.png")
	yikesfulMountain = loadImage(FileSystem, "imgs/yikesful-mountain--test.png")
	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")
	backgroundYikesfulMountain = loadImage(FileSystem, "imgs/level-background-2--test.png")

	ebitengineSplash = loadImage(FileSystem, "imgs/load-ebitengine-splash.png")
	splashImages = append(splashImages, ebitengineSplash)

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")

	gemCt = loadImage(FileSystem, "imgs/gem-count-large.png")
	livesCt = loadImage(FileSystem, "imgs/lives-left.png")
	messageBox = loadImage(FileSystem, "imgs/message-box.png")
	statsBox = loadImage(FileSystem, "imgs/stats-box.png")

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	portal = loadImage(FileSystem, "imgs/portal-b--test.png")
	shinyGreenBall = loadImage(FileSystem, "imgs/treasure--test.png")
	portalGem = loadImage(FileSystem, "imgs/quest-item--test.png")
	hazard = loadImage(FileSystem, "imgs/blob--test.png")
	creature = loadImage(FileSystem, "imgs/creature--test.png")
	gameOverMessage = loadImage(FileSystem, "imgs/game-over.png")

	levelImages = map[string][]*ebiten.Image{
		"Goo Alley":         {gooAlley, levelBG},
		"Yikesful Mountain": {yikesfulMountain, backgroundYikesfulMountain},
	}
}

func findSaveFiles() []string {
	saveFiles := []string{}
	files, err := os.ReadDir("./save/")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range files {
		if string(item.Name()) != "levels.json" {
			saveFiles = append(saveFiles, string(item.Name()))
		}
	}
	saveFiles = append(saveFiles, "Main Menu")
	return saveFiles
}

func loadImage(fs embed.FS, path string) *ebiten.Image {
	log.Printf(" %s", path)
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
