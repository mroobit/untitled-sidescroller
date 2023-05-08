package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shinyGreenBall *ebiten.Image
	portalGem      *ebiten.Image

	treasureTypeList map[int]*TreasureType
	treasureList     []*Treasure
)

func initializeTreasures() {
	treasureTypeList = map[int]*TreasureType{
		3: {name: "Portal Gem", sprite: portalGem, width: 50, height: 50, value: 0, frame: 0, frameCt: 5},
		4: {name: "Shiny Green Ball", sprite: shinyGreenBall, width: 40, height: 40, value: 10, frame: 0, frameCt: 7},
	}
}

// TreasureType holds general description for a specific type of treasure
type TreasureType struct {
	name    string
	sprite  *ebiten.Image
	sprites []*ebiten.Image
	width   int
	height  int
	value   int
	frame   int
	frameCt int
}

// Treasure describes a specific treasure object in a level
type Treasure struct {
	*TreasureType
	xCoord    int
	yCoord    int
	collected bool
}

// NewTreasure creates a new Treasure object at specific coordinates in a level
func NewTreasure(id int, x int, y int) *Treasure {
	log.Printf("Creating new treasure")
	treasure := &Treasure{
		treasureTypeList[id],
		x,
		y,
		false,
	}
	return treasure
}
