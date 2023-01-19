package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Creature struct {
	name        string
	sprite      *ebiten.Image
	facing      int
	xCoord      int
	yCoord      int
	hpCurrent   int
	hpTotal     int
	damage      int
	movement    string // I have no idea how I'm implementing this -- might just key movement style to name, so all same-type creatures move alike
	seesChar    bool
	movementCtr int
	pauseCtr    int
}

func NewCreature(name string, sprite *ebiten.Image, x int, y int, hp int, damage int, movement string) *Creature {
	log.Printf("Creating new creature")
	creature := &Creature{
		name:      name,
		sprite:    sprite,
		facing:    50,
		xCoord:    x,
		yCoord:    y,
		hpCurrent: hp,
		hpTotal:   hp,
		seesChar:  false,
		damage:    damage,
		movement:  name,
	}
	return creature
}
