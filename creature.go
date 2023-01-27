package main

import (
	"log"
	"math/rand"

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

func creatureMovement() {
	for _, c := range creatureList {
		switch {
		case c.movementCtr > 0:
			// keep moving same dir
			c.movementCtr--
			if c.facing == 0 && c.xCoord <= 3 {
				c.movementCtr = 0
			} else if c.facing == 0 && c.xCoord > 3 {
				c.xCoord -= 3
			} else if c.facing == 50 && c.xCoord >= 597 {
				c.movementCtr = 0
			} else if c.xCoord < 597 {
				c.xCoord += 3
			}
		case c.seesChar == true:
			// rampage towards char
			if c.facing == 0 {
				c.xCoord -= 10
			} else {
				c.xCoord += 10
			}
		case c.pauseCtr > 0:
			// pause
			if c.pauseCtr%9 == 0 {
				c.facing = rand.Intn(2) * 50
			}
			c.pauseCtr--
		default:
			// reset random
			c.movementCtr = rand.Intn(50) + 20
			c.pauseCtr = rand.Intn(40) + 20
			c.facing = rand.Intn(2) * 50
		}
	}
}
