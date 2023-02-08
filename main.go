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

	ground  = 380
	gravity = 20

	defaultFrame = 2
	frameCount   = 12

	tileSize   = 50
	tileXCount = 16
	xCount     = winWidth / tileSize
)

var (
	ErrExit = errors.New("Exiting Game")

	//go:embed imgs
	//go:embed fonts
	FileSystem embed.FS

	ebitengineSplash           *ebiten.Image
	gemCt                      *ebiten.Image
	livesCt                    *ebiten.Image
	messageBox                 *ebiten.Image
	statsBox                   *ebiten.Image
	world                      *ebiten.Image
	gooAlley                   *ebiten.Image
	yikesfulMountain           *ebiten.Image
	levelBG                    *ebiten.Image
	backgroundYikesfulMountain *ebiten.Image
	portal                     *ebiten.Image
	creature                   *ebiten.Image
	blank                      *ebiten.Image
	gameOverMessage            *ebiten.Image

	levelWidth  int
	levelHeight int

	radius = 375.0
)

var (
	//	mplusNormalFont font.Face

	currentFrame int
	//portalFrame  int
	//	creatureFrame int

	levelData []*LevelData
	//	creatureList []*Creature

	levelMap [][]int
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

	monaView = NewViewer()
	worldView = NewViewer()

	currentFrame = defaultFrame
	//treasureFrame = defaultFrame
	//portalGemFrame = defaultFrame

	mona = NewCharacter("Mona", spriteSheet, monaView, 100)
	worldMona = NewWorldChar(spriteSheet, worldView)

}

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Starting up Mona Game POC...")

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Mona Game, POC: Movement in Level Space")

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
		// set worldMona location and view screen: this should go in Menu->Start New Game
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			// later, ask confirmation if game not saved since entering World
			// Options: Save, Quit without Saving
			log.Printf("Exiting Game")
			return ErrExit
		}
		// radiusCheck is making sure worldMona stays within movement radius of planet
		radiusCheck := math.Sqrt(math.Pow(float64(worldMona.xCoord-500-worldMona.view.xCoord), 2) + math.Pow(float64(worldMona.yCoord-500-worldMona.view.yCoord), 2))
		// 4 directions of worldMona movement checks
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			worldMona.direction = "right"
			switch {
			case worldMona.view.xCoord == 0 && worldMona.xCoord < 290:
				worldMona.xCoord += 5
			case worldMona.view.xCoord == -400 && radiusCheck+50 < radius: // worldMona.xCoord < 500: // but actually, the arc of the circle
				worldMona.xCoord += 5
			case worldMona.view.xCoord > -400:
				worldMona.view.xCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			worldMona.direction = "left"
			switch {
			case worldMona.view.xCoord == -400 && worldMona.xCoord > 290:
				worldMona.xCoord -= 5
			case worldMona.view.xCoord == 0 && radiusCheck < radius: // worldMona.xCoord < 500: // but actually, the arc of the circle
				worldMona.xCoord -= 5
			case worldMona.view.xCoord < 0:
				worldMona.view.xCoord += 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			worldMona.direction = "up"
			switch {
			case worldMona.view.yCoord == -520 && worldMona.yCoord < 230:
				worldMona.yCoord -= 5
			case worldMona.view.yCoord == 0 && radiusCheck < radius:
				worldMona.yCoord -= 5
			case worldMona.view.yCoord < 0:
				worldMona.view.yCoord += 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			worldMona.direction = "down"
			switch {
			case worldMona.view.yCoord == 0 && worldMona.yCoord > 250:
				worldMona.yCoord += 5
			case worldMona.view.yCoord == -520 && radiusCheck+50 < radius:
				worldMona.yCoord += 5
			case worldMona.view.yCoord > -520:
				worldMona.view.yCoord -= 5
			}
		}
		// locations of levels on World, checking whether conditions are met to enter the level
		for _, l := range levelData {
			if ((worldMona.xCoord > l.WorldX+worldMona.view.xCoord && worldMona.xCoord < l.WorldX+150+worldMona.view.xCoord ||
				worldMona.xCoord+monaWidth > l.WorldX+worldMona.view.xCoord && worldMona.xCoord+monaWidth < l.WorldX+150+worldMona.view.xCoord) &&
				(worldMona.yCoord > l.WorldY+worldMona.view.yCoord && worldMona.yCoord < l.WorldY+150+worldMona.view.yCoord ||
					worldMona.yCoord+monaHeight > l.WorldY+worldMona.view.yCoord && worldMona.yCoord+monaHeight < l.WorldY+150+worldMona.view.yCoord)) &&
				ebiten.IsKeyPressed(ebiten.KeyEnter) &&
				l.Complete == false {

				levelWidth, levelHeight = l.background.Size()
				mona.resetView()
				mona.setLocation(l.PlayerX, l.PlayerY)
				mona.hpCurrent = mona.hpTotal
				levelSetup(l, mona.view.xCoord, mona.view.yCoord)
				g.lvl = l
				g.portalGem = false
				g.mode = Play
			}
		}
	case Play:
		// death check
		if mona.hpCurrent == 0 {
			if mona.lives == 0 {
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
			mona.facing = 0
			currentFrame = (g.count / 5) % frameCount
			switch {
			case mona.view.xCoord == 0 && mona.xCoord < 290:
				mona.xCoord += 5
			case mona.view.xCoord == -200 && mona.xCoord < 530:
				mona.xCoord += 5
			case mona.view.xCoord > -200:
				mona.view.xCoord -= 5
				for _, h := range hazardList {
					h.xCoord -= 5
				}
				for _, c := range creatureList {
					c.xCoord -= 5
				}
			}
			monaSide := (mona.xCoord - mona.view.xCoord + monaWidth + 1) / 50
			monaTop := (mona.yCoord - mona.view.yCoord) / 50
			if levelMap[0][monaTop*tileXCount+monaSide] == 1 /* || levelMap[0][monaBase*tileXCount+monaSide] == 1*/ {
				mona.xCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			mona.facing = monaHeight
			currentFrame = (g.count / 5) % frameCount
			switch {
			case mona.view.xCoord == -200 && mona.xCoord > 290:
				mona.xCoord -= 5
			case mona.view.xCoord == 0 && mona.xCoord > 40:
				mona.xCoord -= 5
			case mona.view.xCoord < 0:
				mona.view.xCoord += 5
				for _, c := range creatureList {
					c.xCoord += 5
				}
				for _, h := range hazardList {
					h.xCoord += 5
				}
			}
			monaSide := (mona.xCoord - mona.view.xCoord) / 50
			monaTop := (mona.yCoord - mona.view.yCoord) / 50
			if levelMap[0][monaTop*tileXCount+monaSide] == 1 {
				mona.xCoord += 5
			}
		}
		// player sprite frame reset
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) || inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
			currentFrame = defaultFrame
		}

		// jump logic
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && mona.yVelo == gravity {
			mona.yVelo = -gravity
		}
		if mona.yVelo < gravity {
			// screen movement vs player movement
			switch {
			case mona.yCoord < 160 && mona.view.yCoord-mona.yVelo < 0 && mona.yVelo < 0:
				//	mona.yCoord -= mona.yVelo
				mona.view.yCoord -= mona.yVelo
				for _, h := range hazardList {
					h.yCoord -= mona.yVelo
				}
				for _, c := range creatureList {
					c.yCoord -= mona.yVelo
				}
			case mona.yCoord > 160 && mona.view.yCoord-mona.yVelo > -120 && mona.yVelo > 0:
				//	mona.yCoord -= mona.yVelo
				mona.view.yCoord -= mona.yVelo
				for _, h := range hazardList {
					h.yCoord -= mona.yVelo
				}
				for _, c := range creatureList {
					c.yCoord -= mona.yVelo
				}
			default:
				mona.yCoord += mona.yVelo
			}
			mona.yVelo += 1

			if mona.yVelo >= 0 {
				monaBase := (mona.yCoord - mona.view.yCoord + monaHeight + 1) / 50 // checks immediately BELOW base of sprite
				monaLeft := (mona.xCoord - mona.view.xCoord) / 50
				monaRight := (mona.xCoord - mona.view.xCoord + monaWidth) / 50
				if levelMap[0][(monaBase)*tileXCount+monaLeft] == 1 || levelMap[0][(monaBase)*tileXCount+monaRight] == 1 {
					mona.yCoord = (monaBase * 50) - 50 + mona.view.yCoord
					mona.yVelo = gravity
				}
			}
		}
		monaTop := (mona.yCoord - mona.view.yCoord) / 50
		monaBase := (mona.yCoord - mona.view.yCoord + monaHeight + 1) / 50 // checks immediately BELOW base of sprite
		monaLeft := (mona.xCoord - mona.view.xCoord) / 50
		monaRight := (mona.xCoord - mona.view.xCoord + monaWidth) / 50
		// gravity fixer
		if mona.yVelo == gravity && levelMap[0][(monaBase*tileXCount)+monaLeft] != 1 && levelMap[0][(monaBase*tileXCount)+monaRight] != 1 {
			switch {
			case mona.view.yCoord > -120 && mona.yCoord > 160:
				mona.view.yCoord -= 3
				for _, h := range hazardList {
					h.yCoord -= 3
				}
				for _, c := range creatureList {
					c.yCoord -= 3
				}
			default:
				mona.yCoord += 3
			}
		}

		blockTopLeft := monaTop*tileXCount + monaLeft
		btlVal := levelMap[3][blockTopLeft] + levelMap[4][blockTopLeft] + levelMap[1][blockTopLeft]
		blockTopRight := monaTop*tileXCount + monaRight
		btrVal := levelMap[3][blockTopRight] + levelMap[4][blockTopRight] + levelMap[1][blockTopRight]
		blockBaseLeft := monaBase*tileXCount + monaLeft
		bblVal := levelMap[3][blockBaseLeft] + levelMap[4][blockBaseLeft] + levelMap[1][blockBaseLeft]
		blockBaseRight := monaBase*tileXCount + monaRight
		bbrVal := levelMap[3][blockBaseRight] + levelMap[4][blockBaseRight] + levelMap[1][blockBaseRight]
		if btlVal == 5 || bblVal == 5 || btrVal == 5 || bbrVal == 5 {
			mona.death()
			clearLevel()
			if mona.lives == 0 {
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
			(mona.xCoord > g.lvl.ExitX+mona.view.xCoord && mona.xCoord < g.lvl.ExitX+50+mona.view.xCoord ||
				mona.xCoord+monaWidth > g.lvl.ExitX+mona.view.xCoord && mona.xCoord+monaWidth < g.lvl.ExitX+50+mona.view.xCoord) &&
			(mona.yCoord > g.lvl.ExitY+mona.view.yCoord && mona.yCoord < g.lvl.ExitY+100+mona.view.yCoord ||
				mona.yCoord+monaHeight > g.lvl.ExitY+mona.view.yCoord && mona.yCoord+monaHeight < g.lvl.ExitY+100+mona.view.yCoord) {
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
		op.GeoM.Translate(float64(worldMona.view.xCoord), float64(worldMona.view.yCoord))
		screen.DrawImage(world, op)
		for _, l := range levelData {
			lop := &ebiten.DrawImageOptions{}
			lop.GeoM.Translate(float64(l.WorldX), float64(l.WorldY))
			world.DrawImage(l.icon, lop)
		}
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(worldMona.xCoord), float64(worldMona.yCoord))
		screen.DrawImage(worldMona.sprite.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)

		if mona.lives <= 0 {
			overOp := &ebiten.DrawImageOptions{}
			screen.DrawImage(gameOverMessage, overOp)
		}
	case Play, Pause:
		lvlOp := &ebiten.DrawImageOptions{}
		lvlOp.GeoM.Translate(float64(mona.view.xCoord), float64(mona.view.yCoord))
		screen.DrawImage(g.lvl.background, lvlOp)
		screen.DrawImage(blank, lvlOp)
		mOp := &ebiten.DrawImageOptions{}
		mOp.GeoM.Translate(float64(mona.xCoord), float64(mona.yCoord))
		cx, cy := currentFrame*monaWidth, mona.facing
		screen.DrawImage(mona.sprite.SubImage(image.Rect(cx, cy, cx+monaWidth, cy+monaHeight)).(*ebiten.Image), mOp)
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

		for lx := 0; lx < mona.lives; lx++ {
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
