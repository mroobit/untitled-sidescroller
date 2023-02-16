package main

import "log"

// Menu is a cyclical doubly-linked list with selectable MenuItems and an active option
type Menu struct {
	head   *MenuItem
	tail   *MenuItem
	active *MenuItem
	length int
}

// MenuItem describes a selectable option as a node in a Menu
type MenuItem struct {
	option string
	prev   *MenuItem
	next   *MenuItem
}

var (
	menuItems = []string{"New Game", "Load Game", "How To Play", "Credits", "Exit"}
	mainMenu  *Menu
)

// NewMenu creates a Menu from a slice of select
func NewMenu(items []string) *Menu {
	menu := &Menu{}
	for _, v := range items {
		menu.appendItem(v)
		menu.length++
	}
	menu.active = menu.head
	return menu
}

// NewMenuItem creates a MenuItem node
func NewMenuItem(o string) *MenuItem {
	menuItem := &MenuItem{
		option: o,
	}
	return menuItem
}

func (m *Menu) appendItem(s string) {
	addition := NewMenuItem(s)
	if m.head == nil {
		m.head = addition
		m.tail = addition
	}
	addition.prev = m.tail
	addition.next = m.head
	m.tail.next = addition
	m.tail = addition
	m.head.prev = m.tail
}

// Next changes active menu selection to next item
func (m *Menu) Next() {
	log.Printf("Next Option")
	m.active = m.active.next
}

// Prev changes active menu selection to previous item
func (m *Menu) Prev() {
	log.Printf("Previous Option")
	m.active = m.active.prev
}

// Select changes game mode according to active MenuItem selected
func (m *Menu) Select() (Mode, error) {
	log.Printf("Selecting an Item")
	switch {
	case m.active.option == "New Game":
		log.Printf("Starting New Game")
		// prompt for character name
		// create character with provided name
		playerView = NewViewer()
		worldPlayerView = NewViewer()

		playerChar = NewCharacter("Mona", spriteSheet, playerView, 100)
		worldPlayer = NewWorldChar(spriteSheet, worldPlayerView)

		// Initialize SaveData with character name
		saveData := NewSaveData()
		saveData.Initialize("Mona")
		return World, nil
	case m.active.option == "Load Game":
		log.Printf("Loading Game -- not yet implemented")
		//TODO
		// display save files available
		// selectable-menu
		// on selection [Enter],
		// saveData := NewSaveData()
		// saveData.Load(savefile)
		// return World, nil
		return Title, nil
	case m.active.option == "How To Play":
		//TODO
		log.Printf("Display Instructions -- not yet implemented")
		return Title, nil
	case m.active.option == "Credits":
		//TODO
		log.Printf("Display Credits -- not yet implemented")
		return Title, nil
	case m.active.option == "Exit":
		log.Printf("Attempting to Exit Game")
		return Title, ErrExit
	}
	return Title, nil
}
