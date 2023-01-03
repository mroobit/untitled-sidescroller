// Package main runs game
package main

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	winWidth  = 600
	winHeight = 480

	defaultFrame = 2
	frameCount   = 12
	frameWidth   = 48
	frameHeight  = 48
)

var (
	//go:embed imgs
	FileSystem embed.FS

	levelBG *ebiten.Image
	//	charaSprite *ebiten.Image
	spriteSheet *ebiten.Image

	levelWidth  int
	levelHeight int

	levelView *Viewer

	mona *Character
)

var (
	currentFrame int
)

func init() {
	log.Printf("Initializing...")
	rawFile, err := FileSystem.Open("imgs/level-background--test.png")
	if err != nil {
		log.Fatalf("Error opening file imgs/level-background--test.png: %v\n", err)
	}
	defer rawFile.Close()

	img, err := png.Decode(rawFile)
	if err != nil {
		log.Fatalf("Error decoding file imgs/level-background--test.png: %v\n", err)
	}

	levelBG = ebiten.NewImageFromImage(img)

	// these values are temporarily hard-coded, replace magic numbers later
	levelWidth = 800
	levelHeight = 600

	levelView = NewViewer()
	/*
		rawFile, err = FileSystem.Open("imgs/mona--2023-01-03.png")
		if err != nil {
			log.Fatalf("Error opening file imgs/character-sprite--test.png: %v\n", err)
		}
		defer rawFile.Close()

		img, err = png.Decode(rawFile)
		if err != nil {
			log.Fatalf("Error decoding file imgs/character-sprite--test.png: %v\n", err)
		}

		charaSprite = ebiten.NewImageFromImage(img)

		mona = NewCharacter("Mona", charaSprite, 100)
	*/

	rawFile, err = FileSystem.Open("imgs/walk-test--2023-01-03.png")
	if err != nil {
		log.Fatalf("Error opening file imgs/walk-test--2023-01-03.png: %v\n", err)
	}

	img, err = png.Decode(rawFile)
	if err != nil {
		log.Fatalf("Error decoding file imgs/walk-test.png: %v\n", err)
	}

	spriteSheet = ebiten.NewImageFromImage(img)

	currentFrame = defaultFrame

	mona = NewCharacter("Mona", spriteSheet, 100)

}

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Running level test")

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Mona Game, POC: Movement in Level Space")

	g := NewGame()
	//	g.Setup()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// Game contains all relevant data for game
type Game struct {
	background *ebiten.Image
	view       *Viewer
	count      int
}

// Viewer is the part of the total level that is visible, as indicated by the X,Y of the upper left corner
type Viewer struct {
	xCoord int
	yCoord int
	width  int
	height int
}

type Character struct {
	name       string
	sprite     *ebiten.Image
	facing     string
	xCoord     int
	yCoord     int
	active     bool
	hp_current int
	hp_total   int
}

func NewGame() *Game {
	log.Printf("Creating new game")
	game := &Game{
		background: levelBG,
		view:       levelView,
		count:      0,
	}
	return game
}

func NewViewer() *Viewer {
	log.Printf("Creating new viewer")
	viewer := &Viewer{
		xCoord: 0,
		yCoord: winHeight - levelHeight,
		width:  winWidth,
		height: winHeight,
	}
	return viewer
}

func NewCharacter(name string, sprite *ebiten.Image, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:       name,
		sprite:     sprite,
		facing:     "right",
		xCoord:     40,
		yCoord:     430,
		active:     false,
		hp_current: hp,
		hp_total:   hp,
	}
	return character
}

func (g *Game) Update() error {
	g.count++
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		mona.facing = "right"
		currentFrame = (g.count / 5) % frameCount
		switch {
		case g.view.xCoord == 0 && mona.xCoord < 290:
			mona.xCoord += 5
		case mona.xCoord == 290 && g.view.xCoord > -200:
			g.view.xCoord -= 5
		case g.view.xCoord == -200 && mona.xCoord < 560:
			mona.xCoord += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		mona.facing = "left"
		switch {
		case g.view.xCoord == -200 && mona.xCoord > 290:
			mona.xCoord -= 5
		case mona.xCoord == 290 && g.view.xCoord < 0:
			g.view.xCoord += 5
		case g.view.xCoord == 0 && mona.xCoord > 40:
			mona.xCoord -= 5
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) || inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		currentFrame = defaultFrame
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	lvlOp := &ebiten.DrawImageOptions{}
	lvlOp.GeoM.Translate(float64(g.view.xCoord), float64(g.view.yCoord))
	screen.DrawImage(g.background, lvlOp)
	mOp := &ebiten.DrawImageOptions{}
	mOp.GeoM.Translate(float64(mona.xCoord), float64(mona.yCoord))
	//	if mona.facing == "left" {
	//		mOp.GeoM.Scale(-1.0, 1.0)
	//	}
	cx, cy := currentFrame*frameWidth, 0
	screen.DrawImage(mona.sprite.SubImage(image.Rect(cx, cy, cx+frameWidth, cy+frameHeight)).(*ebiten.Image), mOp)

	msg := ""
	msg += fmt.Sprintf("Mona xCoord: %d", mona.xCoord)
	msg += fmt.Sprintf("Level xCoord: %d", g.view.xCoord)

	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}
