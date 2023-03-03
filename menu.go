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
	mainMenuItems  = []string{"New Game", "Load Game", "How To Play", "Credits", "Exit"}
	worldMenuItems = []string{"Save", "Stats", "Main Menu", "Quit"}
	loadMenuItems  []string
	mainMenu       *Menu
	worldMenu      *Menu
	loadMenu       *Menu
)

func initializeMenus() {
	mainMenu = NewMenu(mainMenuItems)
	worldMenu = NewMenu(worldMenuItems)
	loadMenu = NewMenu(loadMenuItems)
}

// NewMenu creates a Menu from a slice of strings
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
func (m *Menu) Select() string {
	log.Printf("Selecting " + m.active.option)
	return m.active.option
}
