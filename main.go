// Package main runs game
package main

import (
	"embed"
	"errors"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

const (
	winWidth  = 600
	winHeight = 480

	gravity = 20
	radius  = 375.0
)

var (
	// ErrExit is the "error" that signals to close the game
	ErrExit = errors.New("Exiting Game")

	// FileSystem of images, fonts
	//go:embed imgs
	//go:embed fonts
	//go:embed levels.json
	FileSystem embed.FS
)

func main() {
	/*
		f, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	*/
	log.Printf("Starting up game...")
	loadAssets()
	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("A Pixely Side-Scrolling Game Send-up")

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		if err == ErrExit {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}

// Game contains all relevant data for game
type Game struct {
	state       map[string]State
	mode        string
	txtRenderer *etxt.Renderer
	count       int
	timer       int
	score       int
}

// State describes Game State
type State interface {
	Update(g *Game) error
	Draw(screen *ebiten.Image, g *Game)
}

// NewGame creates a new Game instance (used once, to run program)
func NewGame() *Game {
	log.Printf("Generating new game instance")
	game := &Game{
		state: map[string]State{
			"Load":  &Load{splash: splashImages},
			"Title": &Title{},
			"World": &World{},
			"Play":  &Play{},
			"Pause": &Pause{},
			"Info":  &Info{},
		},
		mode: "Load",
	}
	return game
}

// Update controls all game logic updates. It is part of the main game loop in Ebitengine.
func (g *Game) Update() error {
	g.count++
	err := g.state[g.mode].Update(g)
	return err
}

// Draw contains all code for drawing images to screen. It is part of the main game loop in Ebitengine.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case "Pause":
		g.state["Play"].Draw(screen, g)
		g.state[g.mode].Draw(screen, g)
	default:
		g.state[g.mode].Draw(screen, g)
	}
}

// Layout controls the game window and scaling. It is part of the main game loop in Ebitengine.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//	s := ebiten.DeviceScaleFactor()
	//	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
	return winWidth, winHeight
}
