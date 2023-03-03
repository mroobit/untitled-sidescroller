package main

import (
	"encoding/json"
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

type Load struct {
	splash []*ebiten.Image
	curr   int
}

func (l *Load) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = "Title"
	}
	if g.count > 200 {
		g.mode = "Title"
		log.Printf("Changing state to Title")
	}
	return nil
}

func (l *Load) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(l.splash[l.curr], op)
}

type Title struct {
	menu *Menu
}

func NewTitle() *Title {
	title := &Title{
		menu: mainMenu,
	}
	return title
}

func (t *Title) Load(m *Menu) {
	t.menu = m
}

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

			world := NewWorld()
			g.state["World"] = world
			g.mode = "World"

			// Initialize SaveData with character name
			// saveData := NewSaveData()
			// saveData.Initialize("Mona")
		case selection == "Load Game":
			log.Printf("Choose a Saved Game")
			if len(loadMenuItems) > 1 {
				t := NewTitle()
				t.Load(loadMenu)
				g.state["Title"] = t
			}
			g.mode = "Title"
		case selection == "How To Play":
			//TODO
			log.Printf("Display Instructions -- not yet implemented")
			g.mode = "Title"
		case selection == "Credits":
			//TODO
			log.Printf("Display Credits -- not yet implemented")
			g.mode = "Title"
		case selection == "Exit":
			log.Printf("Attempting to Exit Game")
			return ErrExit
		case strings.HasSuffix(selection, ".json"):
			// TODO
			// Need actual save file, capturing not just *LevelData but also player character, score, etc
			world := NewWorld()
			world.Load(selection)
			g.state["World"] = world
			g.mode = "World"
		case selection == "Main Menu":
			t := NewTitle()
			t.Load(mainMenu)
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

func (t *Title) Draw(screen *ebiten.Image, g *Game) {
	textColor = menuColorActive
	g.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter) // make sure type is centered (gets changed in Play/Pause)

	var menuHead = t.menu.head
	var locY = 100
	for i := t.menu.length; i > 0; i-- {
		textColor = menuColorInactive
		if menuHead == t.menu.active {
			textColor = menuColorActive
		}
		if menuHead == t.menu.active && t.menu.active.option == "Load Game" && len(loadMenuItems) < 2 {
			textColor = menuColorDisabled
		}
		g.txtRenderer.SetTarget(screen)
		g.txtRenderer.SetColor(textColor)
		g.txtRenderer.Draw(menuHead.option, 300, locY)
		locY += 50
		menuHead = menuHead.next
	}
}

type World struct {
	menu   *Menu
	levels []*LevelData
}

func NewWorld() *World {
	world := &World{
		menu:   worldMenu,
		levels: levelList,
	}
	return world
}

func (w *World) Load(path string) {
	// currently this just loads the []*LevelData, but ultimately needs to load player, score, etc data
	var levels []*LevelData
	lvlContent, err := os.ReadFile("./save/" + path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(lvlContent, &levels)
	if err != nil {
		log.Fatal("Error during Unmarshalling: ", err)
	}

	for _, l := range levels {
		l.icon = levelImages[l.Name][0]
		l.background = levelImages[l.Name][1]
	}

	w.levels = levels
}

func (w *World) Update(g *Game) error {
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
			g.mode = "Play"
		}
	}
	return nil
}

func (w *World) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(worldPlayer.view.xCoord), float64(worldPlayer.view.yCoord))
	screen.DrawImage(world, op)
	for _, l := range w.levels {
		lop := &ebiten.DrawImageOptions{}
		lop.GeoM.Translate(float64(l.WorldX), float64(l.WorldY))
		world.DrawImage(l.icon, lop)
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(worldPlayer.xCoord), float64(worldPlayer.yCoord))
	screen.DrawImage(worldPlayer.sprite.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)
}

type Play struct {
	level *LevelData
	gem   bool
}

func NewPlay(l *LevelData) *Play {
	play := &Play{
		level: l,
	}
	return play
}

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

	playerBox := image.Rect(playerChar.xCoord, playerChar.yCoord, playerChar.xCoord+playerCharWidth, playerChar.yCoord+playerCharWidth)

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
		if playerBox.Overlaps(hazardBox) {
			playerChar.death()
			g.mode = "Pause"
			g.timer = 30
		}
	}

	for _, c := range creatureList {
		creatureBox := image.Rect(c.xCoord, c.yCoord, c.xCoord+50, c.yCoord+50)
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

	for lx := 0; lx < playerChar.lives; lx++ {
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

type Pause struct {
	message string
	options *Menu
}

func (p *Pause) Update(g *Game) error {
	switch {
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

func (p *Pause) Draw(screen *ebiten.Image, g *Game) {
	// TODO
	// overlays, based on what the pause message and options are
}
