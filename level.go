package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type LevelData struct {
	Name     string
	Complete bool
	WorldX   int
	WorldY   int
	PlayerX  int
	PlayerY  int
	ExitX    int
	ExitY    int
	Message  []string
	Layout   [][]int

	icon       *ebiten.Image
	background *ebiten.Image // later, this can be []*ebiten.Image, for layered background
}

func populate(lvl *LevelData, vsx int, vsy int) { // pass level name or index number as a parameter, or change to method with *Level as receiver...
	// empty lists first, in case any left over from previous level attempt
	for i, h := range lvl.Layout[1] {
		x := (i%tileXCount)*tileSize - vsx
		y := (i/tileXCount)*tileSize + vsy
		if h == 5 {
			nh := NewHazard("blob", hazard, 10, x, y, 100)
			hazardList = append(hazardList, nh)
		}
		if h == 6 {
			nc := NewCreature("teen yorp", creature, x, y, 100, 100, "teen yorp")
			creatureList = append(creatureList, nc)
		}
	}
}

func clearLevel() { // clear out all hazards, creatures from drawing lists
	hazardList = []*Hazard{}
	creatureList = []*Creature{}
	levelMap = [][]int{}
}

func levelSetup(level *LevelData, viewX int, viewY int) {
	levelMap = layoutCopy(level.Layout)
	populate(level, viewX, viewY)
}

func layoutCopy(layout [][]int) (fresh [][]int) {
	fresh = make([][]int, len(layout))
	for i := range layout {
		fresh[i] = append([]int{}, layout[i]...)
	}
	return
}
