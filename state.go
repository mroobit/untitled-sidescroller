package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

// Load is a Game State for loading the game
type Load struct {
	splash []*ebiten.Image
	curr   int
}

// Update changes Game State to Title after 200 ticks
func (l *Load) Update(g *Game) error {
	if loaded == false {
		loadFonts()
		g.txtRenderer = newRenderer()
		//loadMenuItems = findSaveFiles()
		initializeMenus()
		initializeTreasures()

		t := NewTitle()
		g.state["Title"] = t
		loaded = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = "Title"
	}
	if g.count > 200 {
		g.mode = "Title"
		log.Printf("Changing state to Title")
	}
	return nil
}

// Draw displays splash screens during game load
func (l *Load) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(l.splash[l.curr], op)
}

// Title is a Game State containing Title Screen and a Menu
type Title struct {
	header string
	menu   *Menu
}

// NewTitle creates a new *Title with default main menu
func NewTitle() *Title {
	title := &Title{
		header: mainHeader,
		menu:   mainMenu,
	}
	return title
}

// Load loads a new Menu into *Title
func (t *Title) Load(m *Menu) {
	t.menu = m
}

// Update changes active selection and selects a MenuItem based on user input
func (t *Title) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		selection := t.menu.Select()

		switch {
		case selection == "New Game":
			log.Printf("Starting New Game")
			// prompt for character name
			// create character with provided name
			playerView = NewViewer()
			worldPlayerView = NewViewer()

			playerChar = NewCharacter("Mona", spriteSheet, playerView, 100)
			worldPlayer = NewWorldChar(spriteSheet, worldPlayerView)

			g.score = 0
			g.count = 0

			//	loadLevels()
			world := NewWorld()
			world.Load(FileSystem)
			g.state["World"] = world
			g.mode = "World"

			// saveData := NewSaveData()
			// saveData.Initialize("Mona")
		case selection == "Load Game":
			log.Printf("Choose a Saved Game")
			if len(loadMenuItems) > 1 {
				t := NewTitle()
				t.Load(loadMenu)
				t.header = loadHeader
				g.state["Title"] = t
			}
			g.mode = "Title"
		case selection == "How To Play":
			//TODO
			log.Printf("Display Instructions -- not yet implemented")
			g.mode = "Title"
		case selection == "Acknowledgements":
			c := NewInfo()
			c.message = infoCredit
			c.previous = "Title"
			g.state["Info"] = c
			g.mode = "Info"
		case selection == "Exit":
			log.Printf("Attempting to Exit Game")
			return ErrExit
		case strings.HasSuffix(selection, ".json"):
			gameData := LoadGame(selection)

			playerView = NewViewer()
			worldPlayerView = NewViewer()
			worldPlayerView.xCoord = gameData.WorldViewX
			worldPlayerView.yCoord = gameData.WorldViewY

			playerChar = NewCharacter(gameData.Name, spriteSheet, playerView, 100)
			playerChar.lives = gameData.Lives

			worldPlayer = NewWorldChar(spriteSheet, worldPlayerView)
			worldPlayer.xCoord = gameData.WorldCharX
			worldPlayer.yCoord = gameData.WorldCharY

			g.score = gameData.Score
			g.count = gameData.Count
			world := NewWorld()
			world.Load(FileSystem)
			for _, level := range world.levels {
				if gameData.Complete[level.Name] {
					level.Complete = true
				}
			}

			g.state["World"] = world
			g.mode = "World"

			/*
				world.Load()
			*/

		case selection == "Main Menu":
			t := NewTitle()
			t.Load(mainMenu)
			t.header = mainHeader
			g.state["Title"] = t
			g.mode = "Title"
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		t.menu.Next()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		t.menu.Prev()
	}
	return nil
}

// Draw displays Title Screen and Menu, highlighting active selection
func (t *Title) Draw(screen *ebiten.Image, g *Game) {
	textColor = menuColorActive
	g.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter) // make sure type is centered (gets changed in Play/Pause)
	g.txtRenderer.SetTarget(screen)

	locY := 80
	g.txtRenderer.SetColor(menuColorInactive)
	g.txtRenderer.Draw(t.header, 300, locY)

	var menuHead = t.menu.head
	locY = 150
	for i := t.menu.length; i > 0; i-- {
		textColor = menuColorInactive
		if menuHead == t.menu.active {
			textColor = menuColorActive
		}
		if menuHead == t.menu.active &&
			((t.menu.active.option == "Load Game" && len(loadMenuItems) < 2) || t.menu.active.option == "How To Play") {
			textColor = menuColorDisabled
		}
		g.txtRenderer.SetColor(textColor)
		g.txtRenderer.Draw(menuHead.option, 300, locY)
		locY += 50
		menuHead = menuHead.next
	}
}

// World is a Game State that holds all level data for active game
type World struct {
	menu   *Menu
	levels []*LevelData
}

// NewWorld creates a new World with all levels not yet completed
func NewWorld() *World {
	world := &World{
		menu: worldMenu,
	}
	world.Load(FileSystem)
	return world
}

// Load loads all default level data into World
func (w *World) Load(fs embed.FS) {
	var levels []*LevelData
	lvlContent, err := fs.ReadFile("levels.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(lvlContent, &levels)
	if err != nil {
		log.Fatal("Error during Unmarshalling: ", err)
	}

	for _, l := range levels {
		l.icon = levelImages[l.Name][0]
		l.iconComplete = levelImages[l.Name][1]
		l.background = levelImages[l.Name][2]
	}

	w.levels = levels
}

// Update changes player location/worldview offset and changes state to Play based on user input
func (w *World) Update(g *Game) error {
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			saveData := NewSaveData(g)
			log.Printf("Name is " + saveData.Name)
			log.Printf(strconv.Itoa(saveData.Lives))
			saveData.Save(g, playerChar, worldPlayer)
		}
	*/
	// set worldPlayer location and view screen: this should go in Menu->Start New Game
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// later, ask confirmation if game not saved since entering World
		// Options: Save, Quit without Saving
		log.Printf("Exiting Game")
		return ErrExit
	}
	// radiusCheck is making sure worldPlayer stays within movement radius of planet
	radiusCheck := math.Sqrt(math.Pow(float64(worldPlayer.xCoord-500-worldPlayer.view.xCoord), 2) + math.Pow(float64(worldPlayer.yCoord-500-worldPlayer.view.yCoord), 2))
	// 4 directions of worldPlayer movement checks
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		worldPlayer.navRight(radiusCheck)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		worldPlayer.navLeft(radiusCheck)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		worldPlayer.navUp(radiusCheck)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		worldPlayer.navDown(radiusCheck)
	}

	worldPlayerBox := image.Rect(worldPlayer.xCoord, worldPlayer.yCoord, worldPlayer.xCoord+worldCharWidth, worldPlayer.yCoord+worldCharHeight)
	// locations of levels on World, checking whether conditions are met to enter the level
	for _, l := range w.levels {
		if worldPlayerBox.Overlaps(image.Rect(l.WorldX+worldPlayer.view.xCoord, l.WorldY+worldPlayer.view.yCoord, l.WorldX+worldPlayer.view.xCoord+150, l.WorldY+worldPlayer.view.yCoord+150)) &&
			ebiten.IsKeyPressed(ebiten.KeyEnter) &&
			l.Complete == false {

			levelWidth, levelHeight = l.background.Size()
			playerChar.resetView()
			playerChar.setLocation(l.PlayerX, l.PlayerY)
			playerChar.hpCurrent = playerChar.hpTotal
			levelSetup(l, playerChar.view.xCoord, playerChar.view.yCoord)
			playLevel := NewPlay(l)
			g.state["Play"] = playLevel
			pauseEntry := NewPause("message", l.Message[0])
			g.state["Pause"] = pauseEntry
			g.mode = "Pause"
		}
	}
	return nil
}

// Draw displays player on World map
func (w *World) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(worldPlayer.view.xCoord), float64(worldPlayer.view.yCoord))
	screen.DrawImage(world, op)
	for _, l := range w.levels {
		levelIcon := l.icon
		if l.Complete == true {
			levelIcon = l.iconComplete
		}
		lop := &ebiten.DrawImageOptions{}
		lop.GeoM.Translate(float64(l.WorldX), float64(l.WorldY))
		world.DrawImage(levelIcon, lop)
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(worldPlayer.xCoord), float64(worldPlayer.yCoord))
	screen.DrawImage(worldPlayer.sprite.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)
}

// Play contains data for active level
type Play struct {
	level *LevelData
	gem   bool
}

// NewPlay creates new Play for a given level on entry
func NewPlay(l *LevelData) *Play {
	play := &Play{
		level: l,
	}
	return play
}

// Update is the main gameplay function. Changes score, player health/lives based on user input and collisions
func (p *Play) Update(g *Game) error {
	// sprite frames for different things -- handle differently later
	portalFrame = (g.count / 5) % portalFrameCount
	hazardFrame = (g.count / 5) % hazardFrameCount
	creatureFrame = (g.count / 5) % creatureFrameCount

	treasureTypeList[3].frame = (g.count / 5) % treasureTypeList[3].frameCt
	treasureTypeList[4].frame = (g.count / 5) % treasureTypeList[4].frameCt

	// player sprite frame reset
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) || inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		currentFrame = defaultFrame
	}

	baseView := [2]int{playerChar.view.xCoord, playerChar.view.yCoord}
	// 2 direction movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		currentFrame = (g.count / 5) % frameCount
		playerChar.moveRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		currentFrame = (g.count / 5) % frameCount
		playerChar.moveLeft()
	}
	// if view x changed, update x location of on-screen objects
	if baseView[0] != playerChar.view.xCoord {
		delta := playerChar.view.xCoord - baseView[0]
		for _, h := range hazardList {
			h.xCoord += delta
		}
		for _, c := range creatureList {
			c.xCoord += delta
		}
		for _, t := range treasureList {
			t.xCoord += delta
		}
	}

	// jump logic
	//		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && playerChar.yVelo == gravity { // && a pixel beneath playerChar is enviroBlock
	//			playerChar.yVelo = -gravity
	//		}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		log.Printf("KeyPress Duration: %d", inpututil.KeyPressDuration(ebiten.KeySpace))
		log.Printf("Character status: %s", playerChar.status)
		playerChar.jump(inpututil.KeyPressDuration(ebiten.KeySpace))
	}

	if playerChar.yVelo < gravity {
		// screen movement vs player movement
		if (playerChar.yCoord < 160 && playerChar.view.yCoord-playerChar.yVelo < 0 && playerChar.yVelo < 0) ||
			(playerChar.yCoord > 160 && playerChar.view.yCoord-playerChar.yVelo > -120 && playerChar.yVelo > 0) {
			playerChar.view.yCoord -= playerChar.yVelo
			for _, h := range hazardList {
				h.yCoord -= playerChar.yVelo
			}
			for _, c := range creatureList {
				c.yCoord -= playerChar.yVelo
			}
			for _, t := range treasureList {
				t.yCoord -= playerChar.yVelo
			}
		} else {
			playerChar.yCoord += playerChar.yVelo
		}

		playerChar.yVelo++

		if playerChar.yVelo >= 0 {
			playerCharBase := (playerChar.yCoord - playerChar.view.yCoord + playerCharHeight + 1) / 50 // checks immediately BELOW base of sprite
			playerCharLeft := (playerChar.xCoord - playerChar.view.xCoord) / 50
			playerCharRight := (playerChar.xCoord - playerChar.view.xCoord + playerCharWidth) / 50
			if levelMap[0][(playerCharBase)*tileXCount+playerCharLeft] == 1 || levelMap[0][(playerCharBase)*tileXCount+playerCharRight] == 1 {
				playerChar.yCoord = (playerCharBase * 50) - 50 + playerChar.view.yCoord
				playerChar.yVelo = gravity
			}
		}
	}
	playerCharBase := (playerChar.yCoord - playerChar.view.yCoord + playerCharHeight + 1) / 50 // checks immediately BELOW base of sprite
	playerCharLeft := (playerChar.xCoord - playerChar.view.xCoord) / 50
	playerCharRight := (playerChar.xCoord - playerChar.view.xCoord + playerCharWidth) / 50
	// gravity fixer
	if playerChar.status != "ground" && levelMap[0][(playerCharBase*tileXCount)+playerCharLeft] != 1 && levelMap[0][(playerCharBase*tileXCount)+playerCharRight] != 1 {
		switch {
		case playerChar.view.yCoord > -120 && playerChar.yCoord > 160:
			playerChar.view.yCoord -= 3
			for _, h := range hazardList {
				h.yCoord -= 3
			}
			for _, c := range creatureList {
				c.yCoord -= 3
			}
			for _, t := range treasureList {
				t.yCoord -= 3
			}
		default:
			playerChar.yCoord += 3
		}
	}

	playerCharFreshBase := (playerChar.yCoord - playerChar.view.yCoord + playerCharHeight + 1) / 50 // checks immediately BELOW base of sprite
	if levelMap[0][(playerCharFreshBase*tileXCount)+playerCharLeft] == 1 || levelMap[0][(playerCharFreshBase*tileXCount)+playerCharRight] == 1 {
		playerChar.status = "ground"
	} else if playerChar.yVelo == gravity {
		playerChar.status = "fall"
	}

	creatureMovement()

	//	smallerBox := 3
	playerBox := image.Rect(playerChar.xCoord, playerChar.yCoord, playerChar.xCoord+playerCharWidth, playerChar.yCoord+playerCharWidth)
	cx, cy := currentFrame*playerCharWidth, playerChar.facing
	playerSubImage := playerChar.sprite.SubImage(image.Rect(cx, cy, cx+playerCharWidth, cy+playerCharHeight)).(*ebiten.Image)
	//playerBox := image.Rect(playerChar.xCoord+smallerBox, playerChar.yCoord+smallerBox, playerChar.xCoord+playerCharWidth-smallerBox, playerChar.yCoord+playerCharWidth-smallerBox)

	for i, t := range treasureList {
		treasureBox := image.Rect(t.xCoord, t.yCoord, t.xCoord+50, t.yCoord+50)
		if playerBox.Overlaps(treasureBox) {
			g.score += t.value
			if t.name == "Portal Gem" {
				p.gem = true
			}
			treasureList = append(treasureList[0:i], treasureList[i+1:]...)
		}
	}

	for _, h := range hazardList {
		hazardBox := image.Rect(h.xCoord, h.yCoord, h.xCoord+50, h.yCoord+50)
		//hazardBox := image.Rect(h.xCoord+smallerBox, h.yCoord+smallerBox, h.xCoord+50-smallerBox, h.yCoord+50-smallerBox)
		if playerBox.Overlaps(hazardBox) {
			// subimage hazardbox
			hx := hazardFrame * 50
			hazardSubImage := h.sprite.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image)
			col := Collides(playerSubImage, hazardSubImage)
			if col {
				playerChar.death()
				g.mode = "Pause"
				g.timer = 30
			}
		}
	}

	for _, c := range creatureList {
		creatureBox := image.Rect(c.xCoord, c.yCoord, c.xCoord+50, c.yCoord+50)
		//creatureBox := image.Rect(c.xCoord+smallerBox, c.yCoord+smallerBox, c.xCoord+50-smallerBox, c.yCoord+50-smallerBox)
		if playerBox.Overlaps(creatureBox) {
			playerChar.death()
			g.mode = "Pause"
			g.timer = 30
		}
	}

	if p.gem &&
		playerBox.Overlaps(image.Rect(p.level.ExitX+playerChar.view.xCoord, p.level.ExitY+playerChar.view.yCoord,
			p.level.ExitX+portalWidth+playerChar.view.xCoord, p.level.ExitY+portalHeight+playerChar.view.yCoord)) {
		p.level.Complete = true
		p.gem = false
		clearLevel()
		log.Print("Just hit the portal")
		//levelComplete()
		log.Printf("Level complete")
		g.mode = "World"
	}
	/*
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			g.state = "Pause"
		}
	*/
	return nil
}

// Draw displays level game play
func (p *Play) Draw(screen *ebiten.Image, g *Game) {
	lvlOp := &ebiten.DrawImageOptions{}
	lvlOp.GeoM.Translate(float64(playerChar.view.xCoord), float64(playerChar.view.yCoord))
	screen.DrawImage(p.level.background, lvlOp)

	switch {
	case playerChar.status == "dying":
		mOp := &ebiten.DrawImageOptions{}
		for i := 0; i < playerCharHeight; i += playerCharHeight / 8 {
			wobble := 30 - g.timer
			if i%12 == 0 {
				wobble *= -1
			}
			mOp.GeoM.Reset()
			mOp.GeoM.Translate(float64(playerChar.xCoord+wobble), float64(playerChar.yCoord+i))
			cx, cy := currentFrame*playerCharWidth, playerChar.facing
			screen.DrawImage(playerChar.sprite.SubImage(image.Rect(cx, cy+i, cx+playerCharWidth, cy+i+6)).(*ebiten.Image), mOp)
		}
	default:
		mOp := &ebiten.DrawImageOptions{}
		mOp.GeoM.Translate(float64(playerChar.xCoord), float64(playerChar.yCoord))
		cx, cy := currentFrame*playerCharWidth, playerChar.facing
		screen.DrawImage(playerChar.sprite.SubImage(image.Rect(cx, cy, cx+playerCharWidth, cy+playerCharHeight)).(*ebiten.Image), mOp)
	}

	for _, e := range enviroList {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(e.xCoord), float64(e.yCoord))
		p.level.background.DrawImage(e.sprite, op)
	}
	if p.gem == true {
		top := &ebiten.DrawImageOptions{}
		top.GeoM.Translate(float64(p.level.ExitX+playerChar.view.xCoord), float64(p.level.ExitY+playerChar.view.yCoord))
		px := portalFrame * 100
		screen.DrawImage(portal.SubImage(image.Rect(px, 0, px+100, 150)).(*ebiten.Image), top)
	}
	for _, h := range hazardList {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(h.xCoord), float64(h.yCoord))
		hx := hazardFrame * 50
		screen.DrawImage(h.sprite.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image), op)
	}

	for _, c := range creatureList {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(c.xCoord), float64(c.yCoord))
		cx, cy := creatureFrame*50, c.facing
		screen.DrawImage(c.sprite.SubImage(image.Rect(cx, cy, cx+50, cy+50)).(*ebiten.Image), op)
	}

	for _, t := range treasureList {
		xOffset := (blockHW - t.width) / 2
		yOffset := (blockHW - t.height) / 2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.xCoord+xOffset), float64(t.yCoord+yOffset))
		tx := t.frame * t.width
		screen.DrawImage(t.sprite.SubImage(image.Rect(tx, 0, tx+t.width, t.height)).(*ebiten.Image), op)
	}

	gx := 0
	if p.gem == true {
		gx = 35
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10.0, 10.0)
	screen.DrawImage(statsBox, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(125.0, 64.0)
	screen.DrawImage(gemCt.SubImage(image.Rect(gx, 0, gx+35, 35)).(*ebiten.Image), op)

	for lx := 0; lx < playerChar.lives-1; lx++ {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(21.0+float64(lx*20), 64.0)
		screen.DrawImage(livesCt, op)
	}

	pointsCt := strconv.Itoa(g.score)
	g.txtRenderer.SetTarget(screen)
	g.txtRenderer.SetColor(scoreDisplayColor)
	g.txtRenderer.SetAlign(etxt.Top, etxt.Right)
	g.txtRenderer.Draw(pointsCt, 160, 16)

	if playerChar.status == "totally dead" {
		overOp := &ebiten.DrawImageOptions{}
		screen.DrawImage(gameOverMessage, overOp)
	}
}

// Pause is a Game State that halts other game logic
type Pause struct {
	mode    string
	message string
	options *Menu
}

// NewPause creates new Pause struct
func NewPause(mod, msg string) *Pause {
	p := &Pause{
		mode:    mod,
		message: msg,
	}
	p.FormatMessage()
	return p
}

// FormatMessage adds newlines to Pause.message to fit into messageBox
func (p *Pause) FormatMessage() {
	maxLineLen := 20 // adjust based on txt size, msgbox width
	if len(p.message) > maxLineLen {
		words := strings.Fields(p.message)
		lines := []string{""}
		curr := 0
		for _, w := range words {
			if len(lines[curr])+len(w)+1 > maxLineLen {
				lines[curr] = lines[curr][:len(lines[curr])-1]
				curr++
				lines = append(lines, "")
			}

			lines[curr] += w + " "
		}
		p.message = ""
		for i, l := range lines {
			p.message += l
			if i != len(lines)-1 {
				p.message += "\n"
			}
		}
	}
}

// Update only updates playerChar status for death animation and Game Over screen
func (p *Pause) Update(g *Game) error {
	switch {
	case p.mode == "message":
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			p.mode = ""
			g.mode = "Play"
		}
	case playerChar.status == "totally dead":
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.mode = "Title"
			clearLevel()
			playerChar.status = "ground"
		}
	case playerChar.status == "dying" && g.timer <= 0:
		if playerChar.lives <= 0 {
			playerChar.status = "totally dead"
		}
		if playerChar.lives > 0 {
			clearLevel()
			g.mode = "World"
		}
	case playerChar.status == "dying" && g.timer > 0:
		g.timer--
	default:
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.mode = "Play"
		}
	}
	return nil
}

// Draw isn't implemented yet
func (p *Pause) Draw(screen *ebiten.Image, g *Game) {
	// TODO
	// overlays, based on what the pause message and options are

	switch {
	case p.mode == "message":
		// draw box image
		boxW, boxH := messageBox.Size()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((winWidth-boxW)/2), float64((winHeight-boxH)/2))
		screen.DrawImage(messageBox, op)
		// draw etxt story message
		g.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
		g.txtRenderer.SetSizePx(28)
		g.txtRenderer.SetTarget(screen)
		g.txtRenderer.SetColor(messageBoxColor)
		g.txtRenderer.Draw(p.message, winWidth/2, winHeight/2)
		// draw menu buttons (define these in Update)
	}
}

// Info is currently for Acknowledgement page information, holds previous game state
type Info struct {
	message  []string
	previous string
}

// NewInfo creates new Info struct
func NewInfo() *Info {
	info := &Info{}
	return info
}

// Update returns to previous game state on [Enter]
func (i *Info) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.mode = i.previous
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println(i.previous)
	}
	return nil
}

// Draw draws Acknowledgements to screen
func (i *Info) Draw(screen *ebiten.Image, g *Game) {
	locY := 100
	g.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	g.txtRenderer.SetTarget(screen)
	g.txtRenderer.SetColor(menuColorInactive)
	g.txtRenderer.Draw("Acknowledgements", winWidth/2, locY)
	locY += 100
	g.txtRenderer.SetSizePx(18)
	for _, m := range i.message {
		g.txtRenderer.Draw(m, winWidth/2, locY)
		locY += 75
	}
	locY += 50
	g.txtRenderer.SetSizePx(32)
	g.txtRenderer.SetColor(menuColorActive)
	g.txtRenderer.Draw("Main Menu", winWidth/2, locY)

}
