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
		yCoord:    380,
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

func (c *Character) moveRight() {
	// only impact c
	playerChar.facing = 0
	switch {
	case playerChar.view.xCoord == 0 && playerChar.xCoord < 290:
		playerChar.xCoord += 5
	case playerChar.view.xCoord == -200 && playerChar.xCoord < 530:
		playerChar.xCoord += 5
	case playerChar.view.xCoord > -200:
		playerChar.view.xCoord -= 5
	}
	playerCharSide := (playerChar.xCoord - playerChar.view.xCoord + playerCharWidth + 1) / 50
	playerCharTop := (playerChar.yCoord - playerChar.view.yCoord) / 50
	if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 /* || levelMap[0][playerCharBase*tileXCount+playerCharSide] == 1*/ {
		playerChar.xCoord -= 5
	}
}

func (c *Character) moveLeft() {
	playerChar.facing = playerCharHeight
	switch {
	case playerChar.view.xCoord == -200 && playerChar.xCoord > 290:
		playerChar.xCoord -= 5
	case playerChar.view.xCoord == 0 && playerChar.xCoord > 40:
		playerChar.xCoord -= 5
	case playerChar.view.xCoord < 0:
		playerChar.view.xCoord += 5
	}
	playerCharSide := (playerChar.xCoord - playerChar.view.xCoord) / 50
	playerCharTop := (playerChar.yCoord - playerChar.view.yCoord) / 50
	if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 {
		playerChar.xCoord += 5
	}
}

func (c *Character) death() {
	c.hpCurrent = 0
	c.lives--
	// initiate character death animation
}

func (w *WorldChar) navRight(radiusCheck float64) {
	worldPlayer.direction = "right"
	switch {
	case worldPlayer.view.xCoord == 0 && worldPlayer.xCoord < 290:
		worldPlayer.xCoord += 5
	case worldPlayer.view.xCoord == -400 && radiusCheck+50 < radius: // worldPlayer.xCoord < 500: // but actually, the arc of the circle
		worldPlayer.xCoord += 5
	case worldPlayer.view.xCoord > -400:
		worldPlayer.view.xCoord -= 5
	}

}
func (w *WorldChar) navLeft(radiusCheck float64) {
	worldPlayer.direction = "left"
	switch {
	case worldPlayer.view.xCoord == -400 && worldPlayer.xCoord > 290:
		worldPlayer.xCoord -= 5
	case worldPlayer.view.xCoord == 0 && radiusCheck < radius: // worldPlayer.xCoord < 500: // but actually, the arc of the circle
		worldPlayer.xCoord -= 5
	case worldPlayer.view.xCoord < 0:
		worldPlayer.view.xCoord += 5
	}
}
func (w *WorldChar) navUp(radiusCheck float64) {
	worldPlayer.direction = "up"
	switch {
	case worldPlayer.view.yCoord == -520 && worldPlayer.yCoord < 230:
		worldPlayer.yCoord -= 5
	case worldPlayer.view.yCoord == 0 && radiusCheck < radius:
		worldPlayer.yCoord -= 5
	case worldPlayer.view.yCoord < 0:
		worldPlayer.view.yCoord += 5
	}
}
func (w *WorldChar) navDown(radiusCheck float64) {
	worldPlayer.direction = "down"
	switch {
	case worldPlayer.view.yCoord == 0 && worldPlayer.yCoord > 250:
		worldPlayer.yCoord += 5
	case worldPlayer.view.yCoord == -520 && radiusCheck+50 < radius:
		worldPlayer.yCoord += 5
	case worldPlayer.view.yCoord > -520:
		worldPlayer.view.yCoord -= 5
	}
}
