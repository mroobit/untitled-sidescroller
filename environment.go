package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	brick  *ebiten.Image
	hazard *ebiten.Image
)

var (
	hazardFrame int
	enviroList  []*Brick
	hazardList  []*Hazard
)

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

type Hazard struct {
	name        string
	sprite      *ebiten.Image
	frame_curr  int
	frame_total int
	xCoord      int
	yCoord      int
	damage      int
}

func NewHazard(name string, sprite *ebiten.Image, frames int, x int, y int, damage int) *Hazard {
	log.Printf("Creating new hazard")
	hazard := &Hazard{
		name:        name,
		sprite:      sprite,
		frame_curr:  defaultFrame,
		frame_total: frames,
		xCoord:      x,
		yCoord:      y,
		damage:      damage,
	}
	return hazard
}
