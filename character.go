package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	spriteSheet *ebiten.Image

	playerView       *Viewer
	playerChar       *Character
	playerCharHeight = 48
	playerCharWidth  = 48

	worldPlayerView *Viewer
	worldPlayer     *WorldChar

	defaultFrame = 2
	currentFrame = defaultFrame
	frameCount   = 12
)

// Viewer is the part of the total level that is visible, as indicated by the X,Y of the upper left corner
type Viewer struct {
	xCoord int
	yCoord int
	width  int
	height int
}

type Character struct {
	name      string
	sprite    *ebiten.Image
	view      *Viewer
	facing    int
	xCoord    int
	yCoord    int
	yVelo     int
	active    bool
	hpCurrent int
	hpTotal   int
	lives     int
}

type WorldChar struct { // add to Character struct
	sprite    *ebiten.Image
	view      *Viewer
	direction string
	xCoord    int
	yCoord    int
}

func (c *Character) setLocation(x, y int) {
	log.Printf("Resetting x,y coordinates")
	c.xCoord = x
	c.yCoord = y
}
func (c *Character) resetView() {
	log.Printf("Resetting viewer")
	c.view.xCoord = 0
	c.view.yCoord = winHeight - levelHeight
}

func NewViewer() *Viewer {
	log.Printf("Creating new viewer")
	viewer := &Viewer{
		xCoord: 0,
		yCoord: winHeight,
		width:  winWidth,
		height: winHeight,
	}
	return viewer
}

func NewCharacter(name string, sprite *ebiten.Image, view *Viewer, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:      name,
		sprite:    sprite,
		view:      view,
		facing:    0,
		xCoord:    20,
		yCoord:    ground,
		yVelo:     gravity,
		active:    false,
		hpCurrent: hp,
		hpTotal:   hp,
		lives:     3,
	}
	return character
}

func NewWorldChar(sprite *ebiten.Image, view *Viewer) *WorldChar {
	log.Printf("Creating new world-navigation player character")
	wc := &WorldChar{
		sprite:    sprite,
		view:      view,
		direction: "right",
		xCoord:    200,
		yCoord:    300,
	}
	wc.view.xCoord = -400
	wc.view.yCoord = -500
	return wc
}

func (c *Character) death() {
	c.hpCurrent = 0
	c.lives--
	// initiate character death animation
}
