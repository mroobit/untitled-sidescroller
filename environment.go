package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	blockHW      = 50
	portalWidth  = 100
	portalHeight = 150
)

var (
	brick  *ebiten.Image
	hazard *ebiten.Image
	portal *ebiten.Image
)

var (
	hazardFrame      = 2
	hazardFrameCount = 10
	portalFrame      = 2
	portalFrameCount = 5

	enviroList []*Brick
	hazardList []*Hazard
)

// Brick describes a specific environment object
type Brick struct {
	name         string
	sprite       *ebiten.Image
	xCoord       int
	yCoord       int
	impenetrable bool // can you walk through it
	supportive   bool // can you land on it
	destructible bool // can you destroy it
	//lethal	bool		// will it kill you on contact
	damage int // amount of damage per encounter -- if lethal, set absurdly high
}

// NewBrick creates a new Brick within a level
func NewBrick(name string, sprite *ebiten.Image, x int, y int) *Brick {
	log.Printf("Creating new brick")
	brick := &Brick{
		name:         name,
		sprite:       sprite,
		xCoord:       x,
		yCoord:       y,
		impenetrable: true,
		supportive:   true,
		destructible: false,
		damage:       0,
	}
	return brick
}

// Hazard describes a specific hazardous object
type Hazard struct {
	name       string
	sprite     *ebiten.Image
	frameCurr  int
	frameTotal int
	xCoord     int
	yCoord     int
	damage     int
}

// NewHazard creates a new Hazard within a level
func NewHazard(name string, sprite *ebiten.Image, frames int, x int, y int, damage int) *Hazard {
	log.Printf("Creating new hazard")
	hazard := &Hazard{
		name:       name,
		sprite:     sprite,
		frameCurr:  defaultFrame,
		frameTotal: frames,
		xCoord:     x,
		yCoord:     y,
		damage:     damage,
	}
	return hazard
}
