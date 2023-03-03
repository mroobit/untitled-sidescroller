package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

var (
	lastSaveState *SaveData
)

// SaveData holds minimal data to recreate game state at most recent save
type SaveData struct {
	filename   string
	Name       string
	Lives      int
	Score      int
	Count      int
	Complete   map[string]bool // which levels have been completed
	WorldCharX int
	WorldCharY int
	WorldViewX int
	WorldViewY int
}

// NewSaveData creates a new SaveData struct to hold game state
func NewSaveData(g *Game) *SaveData {
	log.Printf("Creating new SaveData")
	saveData := &SaveData{
		filename:   "",
		Name:       playerChar.name,
		Lives:      playerChar.lives,
		Score:      g.score,
		Count:      g.count,
		Complete:   map[string]bool{},
		WorldCharX: worldPlayer.xCoord,
		WorldCharY: worldPlayer.yCoord,
		WorldViewX: worldPlayer.view.xCoord,
		WorldViewY: worldPlayer.view.yCoord,
	}
	var levels []*LevelData
	switch world := g.state["World"].(type) {
	case *World:
		levels = world.levels
	}
	for _, level := range levels {
		if level.Complete == true {
			saveData.Complete[level.Name] = true
		}
	}
	return saveData
}

// Initialize applies base data for game state
func (s *SaveData) Initialize(name string) {
	s.filename = s.Name + ".json"
}

// LoadGame takes game state data from a save file and stores it in active memory to use as game start point for play session
func LoadGame(savefile string) *SaveData {
	var gameData *SaveData
	saveData, err := os.ReadFile("./save/" + savefile)
	if err != nil {
		log.Fatalf("Error when opening file %s: %v", savefile, err)
	}

	err = json.Unmarshal(saveData, &gameData)
	if err != nil {
		log.Fatalf("Error when unmarshalling save data from %s: %v", savefile, err)
	}

	return gameData
	// later complexity: prompt display loading progress bar
}

// Save takes current game state and put it into SaveData and writes to hard memory
func (s *SaveData) Save(g *Game, p *Character, w *WorldChar) {
	log.Printf("Preparing data to save")
	s.filename = strings.ToLower(s.Name) + "0.json"

	save, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error marshalling save data: %v\n", err)
	}

	log.Printf("Writing data to savefile")

	f, err := os.Create("./save/" + s.filename)

	if err != nil {
		log.Fatalf("Error creating save file %s: %v\n", s.filename, err)
	}

	_, err = f.Write(save)
	if err != nil {
		log.Fatalf("Error writing to file %s: %v\n", s.filename, err)
	}

	saveString := string(save)
	log.Printf(saveString)
	//r later complexity: prompt display saving progress bar
	// prompt display confirmation/failure
}
