package main

import "log"

// menu as cyclical, doubly-linked list

type MenuItem struct {
	option string
	prev   *MenuItem
	next   *MenuItem
}

type Menu struct {
	head   *MenuItem
	tail   *MenuItem
	active *MenuItem
	length int
}

var (
	menuItems = []string{"New Game", "Load Game", "How To Play", "Credits", "Exit"}
	mainMenu  *Menu
)

func NewMenu(items []string) *Menu {
	menu := &Menu{}
	for _, v := range items {
		menu.appendItem(v)
		menu.length++
	}
	menu.active = menu.head
	return menu
}

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

func (m *Menu) Next() {
	log.Printf("Next Option")
	m.active = m.active.next
}

func (m *Menu) Prev() {
	log.Printf("Previous Option")
	m.active = m.active.prev
}

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
