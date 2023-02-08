package main

import "log"

var (
	lastSaveState *SaveData
)

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

func NewSaveData() *SaveData {
	log.Printf("Creating new SaveData")
	saveData := &SaveData{}
	return saveData
}

func (s *SaveData) Initialize(name string) {
	s.filename = name + ".json"
}

func (s *SaveData) Load(savefile string) {
	// open json
	// unmarshall json into SavaData struct
	// later complexity: prompt display loading progress bar
}

func (s *SaveData) Save(g *Game, p *Character, w *WorldChar) {
	// update s with current data from g, p, w
	// marshall into json file, save
	// later complexity: prompt display saving progress bar
	// prompt display confirmation/failure
}
