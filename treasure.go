package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shinyGreenBall *ebiten.Image
	portalGem      *ebiten.Image

	treasureFrame       = 2
	treasureFrameCount  = 7
	portalGemFrame      = 2
	portalGemFrameCount = 5

	treasureType map[int]*TreasureTemplate
	treasureList []*Treasure
)

func initializeTreasures() {
	treasureType = map[int]*TreasureTemplate{
		3: {"Portal Gem", portalGem, 0, 50, 50},
		4: {"Shiny Green Ball", shinyGreenBall, 10, 40, 40},
	}
}

// TreasureTemplate holds general description for a specific type of treasure
type TreasureTemplate struct {
	name   string
	sprite *ebiten.Image
	value  int
	width  int
	height int
}

// Treasure describes a specific treasure object in a level
type Treasure struct {
	name      string
	sprite    *ebiten.Image
	xCoord    int
	yCoord    int
	width     int
	height    int
	value     int
	collected bool
	frame     int
}

// NewTreasure creates a new Treasure object at specific coordinates in a level
func NewTreasure(id int, x int, y int) *Treasure {
	log.Printf("Creating new treasure")
	treasure := &Treasure{
		name:      treasureType[id].name,
		sprite:    treasureType[id].sprite,
		xCoord:    x,
		yCoord:    y,
		width:     treasureType[id].width,
		height:    treasureType[id].height,
		value:     treasureType[id].value,
		collected: false,
	}
	return treasure
}
