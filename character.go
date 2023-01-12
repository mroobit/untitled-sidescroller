package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	spriteSheet *ebiten.Image

	monaView *Viewer
	mona     *Character

	worldMonaView *Viewer
	worldMona     *Character
)

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
	view       *Viewer
	facing     int
	xCoord     int
	yCoord     int
	yVelo      int
	active     bool
	hp_current int
	hp_total   int
	lives      int
}

func (c *Character) viewReset() {
	log.Printf("Resetting viewer")
	c.view.xCoord = 0
	c.view.yCoord = winHeight - levelHeight
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

func NewCharacter(name string, sprite *ebiten.Image, view *Viewer, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:       name,
		sprite:     sprite,
		view:       view,
		facing:     0,
		xCoord:     20,
		yCoord:     ground,
		yVelo:      gravity,
		active:     false,
		hp_current: hp,
		hp_total:   hp,
		lives:      3,
	}
	return character
}

func (c *Character) fade() {
	log.Printf("Fade character")
}
func (c *Character) death() {
	c.hp_current = 0
	c.lives--
	// initiate character death animation
}
