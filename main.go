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
	ErrExit = errors.New("Exiting Game")

	//go:embed imgs
	//go:embed fonts
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
	treasureInit()

	playerView = NewViewer()
	worldPlayerView = NewViewer()

	playerChar = NewCharacter("Mona", spriteSheet, playerView, 100)
	worldPlayer = NewWorldChar(spriteSheet, worldPlayerView)

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

type Mode int

const (
	Load Mode = iota
	Title
	World
	Play
	Pause
)

func NewGame() *Game {
	log.Printf("Creating new game")
	game := &Game{}
	game.mainMenu = NewMenu(menuItems)
	game.txtRenderer = newRenderer()
	return game
}

func (g *Game) Update() error {
	g.count++
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = World
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
			worldPlayer.direction = "right"
			switch {
			case worldPlayer.view.xCoord == 0 && worldPlayer.xCoord < 290:
				worldPlayer.xCoord += 5
			case worldPlayer.view.xCoord == -400 && radiusCheck+50 < radius: // worldPlayer.xCoord < 500: // but actually, the arc of the circle
				worldPlayer.xCoord += 5
			case worldPlayer.view.xCoord > -400:
				worldPlayer.view.xCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			worldPlayer.direction = "left"
			switch {
			case worldPlayer.view.xCoord == -400 && worldPlayer.xCoord > 290:
				worldPlayer.xCoord -= 5
			case worldPlayer.view.xCoord == 0 && radiusCheck < radius: // worldPlayer.xCoord < 500: // but actually, the arc of the circle
				worldPlayer.xCoord -= 5
			case worldPlayer.view.xCoord < 0:
				worldPlayer.view.xCoord += 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			worldPlayer.direction = "up"
			switch {
			case worldPlayer.view.yCoord == -520 && worldPlayer.yCoord < 230:
				worldPlayer.yCoord -= 5
			case worldPlayer.view.yCoord == 0 && radiusCheck < radius:
				worldPlayer.yCoord -= 5
			case worldPlayer.view.yCoord < 0:
				worldPlayer.view.yCoord += 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			worldPlayer.direction = "down"
			switch {
			case worldPlayer.view.yCoord == 0 && worldPlayer.yCoord > 250:
				worldPlayer.yCoord += 5
			case worldPlayer.view.yCoord == -520 && radiusCheck+50 < radius:
				worldPlayer.yCoord += 5
			case worldPlayer.view.yCoord > -520:
				worldPlayer.view.yCoord -= 5
			}
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
		// death check
		if playerChar.hpCurrent == 0 {
			if playerChar.lives == 0 {
				//	g.over()
				log.Printf("Player ran out of lives - Game Over")
			}
			//		g.retryLevel()
			log.Printf("Player still has lives -- retry level or return to world?")
		}
		// sprite frames for different things -- handle differently later
		portalFrame = (g.count / 5) % portalFrameCount
		treasureFrame = (g.count / 5) % treasureFrameCount
		portalGemFrame = (g.count / 5) % portalGemFrameCount
		hazardFrame = (g.count / 5) % hazardFrameCount
		creatureFrame = (g.count / 5) % creatureFrameCount
		// 2 direction movement
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			playerChar.facing = 0
			currentFrame = (g.count / 5) % frameCount
			switch {
			case playerChar.view.xCoord == 0 && playerChar.xCoord < 290:
				playerChar.xCoord += 5
			case playerChar.view.xCoord == -200 && playerChar.xCoord < 530:
				playerChar.xCoord += 5
			case playerChar.view.xCoord > -200:
				playerChar.view.xCoord -= 5
				for _, h := range hazardList {
					h.xCoord -= 5
				}
				for _, c := range creatureList {
					c.xCoord -= 5
				}
			}
			playerCharSide := (playerChar.xCoord - playerChar.view.xCoord + playerCharWidth + 1) / 50
			playerCharTop := (playerChar.yCoord - playerChar.view.yCoord) / 50
			if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 /* || levelMap[0][playerCharBase*tileXCount+playerCharSide] == 1*/ {
				playerChar.xCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			playerChar.facing = playerCharHeight
			currentFrame = (g.count / 5) % frameCount
			switch {
			case playerChar.view.xCoord == -200 && playerChar.xCoord > 290:
				playerChar.xCoord -= 5
			case playerChar.view.xCoord == 0 && playerChar.xCoord > 40:
				playerChar.xCoord -= 5
			case playerChar.view.xCoord < 0:
				playerChar.view.xCoord += 5
				for _, c := range creatureList {
					c.xCoord += 5
				}
				for _, h := range hazardList {
					h.xCoord += 5
				}
			}
			playerCharSide := (playerChar.xCoord - playerChar.view.xCoord) / 50
			playerCharTop := (playerChar.yCoord - playerChar.view.yCoord) / 50
			if levelMap[0][playerCharTop*tileXCount+playerCharSide] == 1 {
				playerChar.xCoord += 5
			}
		}
		// player sprite frame reset
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) || inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
			currentFrame = defaultFrame
		}

		// jump logic
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && playerChar.yVelo == gravity {
			playerChar.yVelo = -gravity
		}
		if playerChar.yVelo < gravity {
			// screen movement vs player movement
			switch {
			case playerChar.yCoord < 160 && playerChar.view.yCoord-playerChar.yVelo < 0 && playerChar.yVelo < 0:
				//	playerChar.yCoord -= playerChar.yVelo
				playerChar.view.yCoord -= playerChar.yVelo
				for _, h := range hazardList {
					h.yCoord -= playerChar.yVelo
				}
				for _, c := range creatureList {
					c.yCoord -= playerChar.yVelo
				}
			case playerChar.yCoord > 160 && playerChar.view.yCoord-playerChar.yVelo > -120 && playerChar.yVelo > 0:
				//	playerChar.yCoord -= playerChar.yVelo
				playerChar.view.yCoord -= playerChar.yVelo
				for _, h := range hazardList {
					h.yCoord -= playerChar.yVelo
				}
				for _, c := range creatureList {
					c.yCoord -= playerChar.yVelo
				}
			default:
				playerChar.yCoord += playerChar.yVelo
			}
			playerChar.yVelo += 1

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
		playerCharTop := (playerChar.yCoord - playerChar.view.yCoord) / 50
		playerCharBase := (playerChar.yCoord - playerChar.view.yCoord + playerCharHeight + 1) / 50 // checks immediately BELOW base of sprite
		playerCharLeft := (playerChar.xCoord - playerChar.view.xCoord) / 50
		playerCharRight := (playerChar.xCoord - playerChar.view.xCoord + playerCharWidth) / 50
		// gravity fixer
		if playerChar.yVelo == gravity && levelMap[0][(playerCharBase*tileXCount)+playerCharLeft] != 1 && levelMap[0][(playerCharBase*tileXCount)+playerCharRight] != 1 {
			switch {
			case playerChar.view.yCoord > -120 && playerChar.yCoord > 160:
				playerChar.view.yCoord -= 3
				for _, h := range hazardList {
					h.yCoord -= 3
				}
				for _, c := range creatureList {
					c.yCoord -= 3
				}
			default:
				playerChar.yCoord += 3
			}
		}

		blockTopLeft := playerCharTop*tileXCount + playerCharLeft
		btlVal := levelMap[3][blockTopLeft] + levelMap[4][blockTopLeft] + levelMap[1][blockTopLeft]
		blockTopRight := playerCharTop*tileXCount + playerCharRight
		btrVal := levelMap[3][blockTopRight] + levelMap[4][blockTopRight] + levelMap[1][blockTopRight]
		blockBaseLeft := playerCharBase*tileXCount + playerCharLeft
		bblVal := levelMap[3][blockBaseLeft] + levelMap[4][blockBaseLeft] + levelMap[1][blockBaseLeft]
		blockBaseRight := playerCharBase*tileXCount + playerCharRight
		bbrVal := levelMap[3][blockBaseRight] + levelMap[4][blockBaseRight] + levelMap[1][blockBaseRight]
		if btlVal == 5 || bblVal == 5 || btrVal == 5 || bbrVal == 5 {
			playerChar.death()
			clearLevel()
			if playerChar.lives == 0 {
				log.Printf("Game over, no lives left")
				//g.over()
			}
			log.Printf("Retry Level or return to world map?")
			g.mode = World
			//g.retryLevel()

		}
		if btlVal == 3 || btlVal == 4 {
			switch {
			case btlVal == 3:
				g.portalGem = true
				levelMap[4][blockTopLeft] = 0
			case btlVal == 4:
				g.score += 10
				levelMap[3][blockTopLeft] = 0
			}
			blank.Clear()
		}
		if btrVal == 3 || btrVal == 4 {
			switch {
			case btrVal == 3:
				g.portalGem = true
				levelMap[4][blockTopRight] = 0
			case btrVal == 4:
				g.score += 10
				levelMap[3][blockTopRight] = 0
			}
			blank.Clear()
		}
		if bblVal == 3 || bblVal == 4 {
			switch {
			case bblVal == 3:
				g.portalGem = true
				levelMap[4][blockBaseLeft] = 0
			case bblVal == 4:
				g.score += 10
				levelMap[3][blockBaseLeft] = 0
			}
			blank.Clear()
		}
		if bbrVal == 3 || bbrVal == 4 {
			switch {
			case bbrVal == 3:
				g.portalGem = true
				levelMap[4][blockBaseRight] = 0
			case bbrVal == 4:
				g.score += 10
				levelMap[3][blockBaseRight] = 0
			}
			blank.Clear()
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

func (g *Game) Draw(screen *ebiten.Image) {
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
				case t == 3:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
					qx := portalGemFrame * 50
					blank.DrawImage(portalGem.SubImage(image.Rect(qx, 0, qx+50, 50)).(*ebiten.Image), top)
				case t == 4:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize+5), float64(i/tileXCount*tileSize+5))
					tx := treasureFrame * 40
					blank.DrawImage(treasure.SubImage(image.Rect(tx, 0, tx+40, 40)).(*ebiten.Image), top)
				}
			}
		}
		for _, e := range enviroList {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(e.xCoord), float64(e.yCoord))
			g.lvl.background.DrawImage(e.sprite, op)
		}
		for _, h := range hazardList {
			hop := &ebiten.DrawImageOptions{}
			hop.GeoM.Translate(float64(h.xCoord), float64(h.yCoord))
			hx := hazardFrame * 50
			screen.DrawImage(h.sprite.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image), hop)
		}

		for _, c := range creatureList {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c.xCoord), float64(c.yCoord))
			cx, cy := creatureFrame*50, c.facing
			screen.DrawImage(c.sprite.SubImage(image.Rect(cx, cy, cx+50, cy+50)).(*ebiten.Image), op)
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
	//	msg := ""
	//	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}
