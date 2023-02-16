package main

import "log"

var (
	lastSaveState *SaveData
)

// SaveData holds minimal data to recreate game state at most recent save
type SaveData struct {
	filename string
	count    int
	score    int
	level    [][]string
	charname string
	// name, hpTotal, lives
	// worldChar x,y
	// worldView x,y
}

// NewSaveData creates a new SaveData struct to hold game state
func NewSaveData() *SaveData {
	log.Printf("Creating new SaveData")
	saveData := &SaveData{}
	return saveData
}

// Initialize applies base data for game state
func (s *SaveData) Initialize(name string) {
	s.filename = name + ".json"
}

// Load takes game state data from a save file and stores it in active memory to use as game start point for play session
func (s *SaveData) Load(savefile string) {
	// open json
	// unmarshall json into SavaData struct
	// later complexity: prompt display loading progress bar
}

// Save takes current game state and put it into SaveData and writes to hard memory
func (s *SaveData) Save(g *Game, p *Character, w *WorldChar) {
	// update s with current data from g, p, w
	// marshall into json file, save
	// later complexity: prompt display saving progress bar
	// prompt display confirmation/failure
}
