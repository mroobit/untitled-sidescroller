// Package main runs game
package main

import (
	"embed"
	"errors"
	"image"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

const (
	winWidth  = 600
	winHeight = 480

	gravity = 20
	radius  = 375.0
)

var (
	// ErrExit is the "error" that signals to close the game
	ErrExit = errors.New("Exiting Game")

	// FileSystem of images, fonts
	//go:embed imgs
	//go:embed fonts
	//go:embed levels.json
	FileSystem embed.FS

	ebitengineSplash *ebiten.Image

	gemCt      *ebiten.Image
	livesCt    *ebiten.Image
	messageBox *ebiten.Image
	statsBox   *ebiten.Image

	world *ebiten.Image

	blank           *ebiten.Image
	gameOverMessage *ebiten.Image
)

func init() {
	/*
		f, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Printf("Initializing...")
	*/
	rand.Seed(time.Now().UnixNano())
	loadAssets()
	initializeTreasures()
	// if savefiles available, create *Menu of savefiles to choose from
	// for items in savefile location, add to slice
	// loadMenu = NewMenu(saveitemsslice)

}

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Starting up game...")

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("A Pixely Side-Scrolling Game Send-up")

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		if err == ErrExit {
			os.Exit(0)
		}
		log.Fatal(err)
	}
}

// Game contains all relevant data for game
type Game struct {
	mode        Mode
	mainMenu    *Menu
	txtRenderer *etxt.Renderer
	count       int
	lvl         *LevelData
	portalGem   bool // deprecate? Could keep, to cheaply track whether to open portal -- maybe rename levelItem
	score       int
}

// Mode indicates current gameplay state
type Mode int

const (
	Load Mode = iota
	Title
	World
	Play
	Pause
)

// NewGame creates a new Game instance (used once, to run program)
func NewGame() *Game {
	log.Printf("Creating new game")
	game := &Game{}
	game.mainMenu = NewMenu(menuItems)
	game.txtRenderer = newRenderer()
	return game
}

// Update controls all game logic updates. It is part of the main game loop in Ebitengine.
func (g *Game) Update() error {
	g.count++
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = Title
	}
	switch g.mode {
	case Load:
		if g.count > 200 {
			g.mode = Title
			log.Printf("Changing mode to Title")
		}
	case Title:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			selection, err := g.mainMenu.Select()
			if err != nil {
				return err
			}
			g.mode = selection
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.mainMenu.Next()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.mainMenu.Prev()
		}
	case World:
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
		// locations of levels on World, checking whether conditions are met to enter the level
		for _, l := range levelData {
			if ((worldPlayer.xCoord > l.WorldX+worldPlayer.view.xCoord && worldPlayer.xCoord < l.WorldX+150+worldPlayer.view.xCoord ||
				worldPlayer.xCoord+playerCharWidth > l.WorldX+worldPlayer.view.xCoord && worldPlayer.xCoord+playerCharWidth < l.WorldX+150+worldPlayer.view.xCoord) &&
				(worldPlayer.yCoord > l.WorldY+worldPlayer.view.yCoord && worldPlayer.yCoord < l.WorldY+150+worldPlayer.view.yCoord ||
					worldPlayer.yCoord+playerCharHeight > l.WorldY+worldPlayer.view.yCoord && worldPlayer.yCoord+playerCharHeight < l.WorldY+150+worldPlayer.view.yCoord)) &&
				ebiten.IsKeyPressed(ebiten.KeyEnter) &&
				l.Complete == false {

				levelWidth, levelHeight = l.background.Size()
				playerChar.resetView()
				playerChar.setLocation(l.PlayerX, l.PlayerY)
				playerChar.hpCurrent = playerChar.hpTotal
				levelSetup(l, playerChar.view.xCoord, playerChar.view.yCoord)
				g.lvl = l
				g.portalGem = false
				g.mode = Play
			}
		}
	case Play:
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

		playerBox := image.Rect(playerChar.xCoord, playerChar.yCoord, playerChar.xCoord+playerCharWidth, playerChar.yCoord+playerCharWidth)

		for i, t := range treasureList {
			treasureBox := image.Rect(t.xCoord, t.yCoord, t.xCoord+50, t.yCoord+50)
			if playerBox.Overlaps(treasureBox) {
				g.score += t.value
				if t.name == "Portal Gem" {
					g.portalGem = true
				}
				treasureList = append(treasureList[0:i], treasureList[i+1:]...)
			}
		}

		for _, h := range hazardList {
			hazardBox := image.Rect(h.xCoord, h.yCoord, h.xCoord+50, h.yCoord+50)
			if playerBox.Overlaps(hazardBox) {
				playerChar.death()
				clearLevel()
				g.mode = World
			}
		}

		if g.portalGem &&
			(playerChar.xCoord > g.lvl.ExitX+playerChar.view.xCoord && playerChar.xCoord < g.lvl.ExitX+50+playerChar.view.xCoord ||
				playerChar.xCoord+playerCharWidth > g.lvl.ExitX+playerChar.view.xCoord && playerChar.xCoord+playerCharWidth < g.lvl.ExitX+50+playerChar.view.xCoord) &&
			(playerChar.yCoord > g.lvl.ExitY+playerChar.view.yCoord && playerChar.yCoord < g.lvl.ExitY+100+playerChar.view.yCoord ||
				playerChar.yCoord+playerCharHeight > g.lvl.ExitY+playerChar.view.yCoord && playerChar.yCoord+playerCharHeight < g.lvl.ExitY+100+playerChar.view.yCoord) {
			g.lvl.Complete = true
			g.portalGem = false
			clearLevel()
			log.Print("Just hit the portal")
			//levelComplete()
			log.Printf("Level complete")
			g.mode = World
		}

		creatureMovement()

		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			g.mode = Pause
		}
	case Pause:
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.mode = Play
		}
	}
	return nil
}

// Draw contains all code for drawing images to screen. It is part of the main game loop in Ebitengine.
func (g *Game) Draw(screen *ebiten.Image) {
	treasureMsg := ""

	switch g.mode {
	case Load:
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(ebitengineSplash, op)
	case Title:
		textColor = menuColorActive
		g.txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter) // make sure type is centered (gets changed in Play/Pause)

		var menuHead = g.mainMenu.head
		var locY = 100
		for i := g.mainMenu.length; i > 0; i-- {
			textColor = menuColorInactive
			if menuHead == g.mainMenu.active {
				textColor = menuColorActive
			}
			g.txtRenderer.SetTarget(screen)
			g.txtRenderer.SetColor(textColor)
			g.txtRenderer.Draw(menuHead.option, 300, locY)
			locY += 50
			menuHead = menuHead.next
		}

	case World:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(worldPlayer.view.xCoord), float64(worldPlayer.view.yCoord))
		screen.DrawImage(world, op)
		for _, l := range levelData {
			lop := &ebiten.DrawImageOptions{}
			lop.GeoM.Translate(float64(l.WorldX), float64(l.WorldY))
			world.DrawImage(l.icon, lop)
		}
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(worldPlayer.xCoord), float64(worldPlayer.yCoord))
		screen.DrawImage(worldPlayer.sprite.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)

		if playerChar.lives <= 0 {
			overOp := &ebiten.DrawImageOptions{}
			screen.DrawImage(gameOverMessage, overOp)
		}
	case Play, Pause:
		lvlOp := &ebiten.DrawImageOptions{}
		lvlOp.GeoM.Translate(float64(playerChar.view.xCoord), float64(playerChar.view.yCoord))
		screen.DrawImage(g.lvl.background, lvlOp)
		screen.DrawImage(blank, lvlOp)
		mOp := &ebiten.DrawImageOptions{}
		mOp.GeoM.Translate(float64(playerChar.xCoord), float64(playerChar.yCoord))
		cx, cy := currentFrame*playerCharWidth, playerChar.facing
		screen.DrawImage(playerChar.sprite.SubImage(image.Rect(cx, cy, cx+playerCharWidth, cy+playerCharHeight)).(*ebiten.Image), mOp)
		for _, l := range levelMap {
			for i, t := range l {
				switch {
				case t == 2:
					if g.portalGem == true {
						top := &ebiten.DrawImageOptions{}
						top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
						px := portalFrame * 100
						blank.DrawImage(portal.SubImage(image.Rect(px, 0, px+100, 150)).(*ebiten.Image), top)
					}
				}
			}
		}
		for _, e := range enviroList {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(e.xCoord), float64(e.yCoord))
			g.lvl.background.DrawImage(e.sprite, op)
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
		if g.portalGem == true {
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

	}
	// msg := ""
	// ebitenutil.DebugPrint(screen, msg)
}

// Layout controls the game window and scaling. It is part of the main game loop in Ebitengine.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
