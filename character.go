package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	spriteSheet *ebiten.Image
	charSprite  []*ebiten.Image

	playerView       *Viewer
	playerChar       *Character
	playerCharHeight = 48
	playerCharWidth  = 48

	worldPlayerView *Viewer
	worldPlayer     *WorldChar
	worldCharWidth  = 48
	worldCharHeight = 48

	defaultFrame = 2
	currentFrame = defaultFrame
	frameCount   = 12
	charFrameCt  = 12
)

// Viewer is the part of the total level that is visible, as indicated by the X,Y of the upper left corner
type Viewer struct {
	xCoord int
	yCoord int
	width  int
	height int
}

// Character describes the player character's state
type Character struct {
	name      string
	sprite    *ebiten.Image
	sprites   []*ebiten.Image
	view      *Viewer
	facing    int
	xCoord    int
	yCoord    int
	xSpeed    int
	yVelo     int
	active    bool
	status    string
	hpCurrent int
	hpTotal   int
	lives     int
}

// WorldChar describes the player navigation avatar on the main screen
type WorldChar struct { // add to Character struct
	sprite    *ebiten.Image
	view      *Viewer
	direction string
	speed     int
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

// NewViewer creates new Viewer (screen offset)
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

// NewCharacter creates new player character
func NewCharacter(name string, sprite *ebiten.Image, view *Viewer, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:      name,
		sprite:    sprite,
		sprites:   charSprite,
		view:      view,
		facing:    0,
		xCoord:    20,
		yCoord:    380,
		xSpeed:    5,
		yVelo:     gravity,
		active:    false,
		status:    "ground",
		hpCurrent: hp,
		hpTotal:   hp,
		lives:     4,
	}
	return character
}

// NewWorldChar creates new player navigation avatar
func NewWorldChar(sprite *ebiten.Image, view *Viewer) *WorldChar {
	log.Printf("Creating new world-navigation player character")
	wc := &WorldChar{
		sprite:    sprite,
		view:      view,
		direction: "right",
		speed:     5,
		xCoord:    200,
		yCoord:    300,
	}
	wc.view.xCoord = -400
	wc.view.yCoord = -500
	return wc
}

func (c *Character) moveRight() {
	// only impact c
	c.facing = 0
	switch {
	case c.view.xCoord == 0 && c.xCoord < 290: // no offset + player up to just under half-way point of winWidth
		c.xCoord += c.xSpeed
	case c.view.xCoord == -200 && c.xCoord < 530: // full offset + player up to 70 less than winWidth, but 32 off from it because playerWidth = 48
		c.xCoord += c.xSpeed
	case c.view.xCoord > -200: // winSize - levelBGSize
		c.view.xCoord -= c.xSpeed
	}
	playerCharSide := (c.xCoord - c.view.xCoord + playerCharWidth + 1) / 50
	playerCharTop := (c.yCoord - c.view.yCoord) / 50
	if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 /* || levelMap[0][playerCharBase*tileXCount+playerCharSide] == 1*/ {
		c.xCoord -= c.xSpeed
	}
}

func (c *Character) moveLeft() {
	c.facing = 1
	switch {
	case c.view.xCoord == -200 && c.xCoord > 290:
		c.xCoord -= c.xSpeed
	case c.view.xCoord == 0 && c.xCoord > 40:
		c.xCoord -= c.xSpeed
	case c.view.xCoord < 0:
		c.view.xCoord += c.xSpeed
	}
	playerCharSide := (c.xCoord - c.view.xCoord) / 50
	playerCharTop := (c.yCoord - c.view.yCoord) / 50
	if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 {
		c.xCoord += c.xSpeed
	}
}

func (c *Character) jump(duration int) { // strength is keypress duration
	switch {
	case c.status == "ground" && duration == 1:
		c.status = "jump"
		c.yVelo = -gravity
		//		c.yVelo = -10
		//	case c.status == "jump" && duration == 2:
		//		c.yVelo += -6
		//	case c.status == "jump" && duration == 3:
		//		c.yVelo += -4
	}
}

func (c *Character) death() {
	c.hpCurrent = 0
	c.lives--
	c.status = "dying"
}

func (w *WorldChar) navRight(radiusCheck float64) {
	w.direction = "right"
	switch {
	case w.view.xCoord == 0 && w.xCoord < 290:
		w.xCoord += w.speed
	case w.view.xCoord == -400 && radiusCheck+50 < radius: // w.xCoord < 500: // but actually, the arc of the circle
		w.xCoord += w.speed
	case w.view.xCoord > -400:
		w.view.xCoord -= w.speed
	}

}
func (w *WorldChar) navLeft(radiusCheck float64) {
	w.direction = "left"
	switch {
	case w.view.xCoord == -400 && w.xCoord > 290:
		w.xCoord -= w.speed
	case w.view.xCoord == 0 && radiusCheck < radius: // w.xCoord < 500: // but actually, the arc of the circle
		w.xCoord -= w.speed
	case w.view.xCoord < 0:
		w.view.xCoord += w.speed
	}
}
func (w *WorldChar) navUp(radiusCheck float64) {
	w.direction = "up"
	switch {
	case w.view.yCoord == -520 && w.yCoord > 230:
		w.yCoord -= w.speed
	case w.view.yCoord == 0 && radiusCheck < radius:
		w.yCoord -= w.speed
	case w.view.yCoord < 0:
		w.view.yCoord += w.speed
	}
}
func (w *WorldChar) navDown(radiusCheck float64) {
	w.direction = "down"
	switch {
	case w.view.yCoord == 0 && w.yCoord < 250:
		w.yCoord += w.speed
	case w.view.yCoord == -520 && radiusCheck+50 < radius:
		w.yCoord += w.speed
	case w.view.yCoord > -520:
		w.view.yCoord -= w.speed
	}
}
