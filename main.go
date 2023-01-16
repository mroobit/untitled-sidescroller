// Package main runs game
package main

import (
	"embed"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	winWidth  = 600
	winHeight = 480

	ground  = 380
	gravity = 20

	defaultFrame = 2
	frameCount   = 12
	frameWidth   = 48
	frameHeight  = 48

	tileSize   = 50
	tileXCount = 16
	xCount     = winWidth / tileSize
)

var (
	//go:embed imgs
	//go:embed fonts
	FileSystem embed.FS

	ebitengineSplash *ebiten.Image
	world            *ebiten.Image
	gooAlley         *ebiten.Image
	yikesfulMountain *ebiten.Image
	levelBG          *ebiten.Image
	portal           *ebiten.Image
	creature         *ebiten.Image
	blank            *ebiten.Image
	gameOverMessage  *ebiten.Image

	levelWidth  int
	levelHeight int

	levelMap = [][]int{
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0,
			0, 0, 0, 4, 0, 5, 0, 0, 6, 4, 0, 0, 0, 5, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	}
	/*
		worldMap = []int{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 2, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 0, 0, 0, 0,
			0, 0, 0, 5094, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}

		worldMapAlt = [][]int{
			{1, 600, 300, 30},
			{2, 400, 500, 50},
		}
	*/

	radius = 375.0
)

var (
	//	mplusNormalFont font.Face

	currentFrame  int
	portalFrame   int
	creatureFrame int

	levels       []*Level
	creatureList []*Creature
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

	world = loadImage(FileSystem, "imgs/world--test.png")
	gooAlley = loadImage(FileSystem, "imgs/level-1--test.png")
	yikesfulMountain = loadImage(FileSystem, "imgs/level-2--test.png")
	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")
	// these values are temporarily hard-coded, replace magic numbers later
	levelWidth = 800
	levelHeight = 600

	monaView = NewViewer()
	worldMonaView = NewViewer()

	ebitengineSplash = loadImage(FileSystem, "imgs/load-ebitengine-splash.png")

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")
	currentFrame = defaultFrame
	treasureFrame = defaultFrame
	questItemFrame = defaultFrame
	mona = NewCharacter("Mona", spriteSheet, monaView, 100)
	worldMona = NewCharacter("World Mona", spriteSheet, worldMonaView, 100)

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	basicBrick = NewBrick("basic", brick)

	portal = loadImage(FileSystem, "imgs/portal-b--test.png")
	treasure = loadImage(FileSystem, "imgs/treasure--test.png")
	questItem = loadImage(FileSystem, "imgs/quest-item--test.png")
	hazard = loadImage(FileSystem, "imgs/blob--test.png")
	creature = loadImage(FileSystem, "imgs/creature--test.png")
	blank = loadImage(FileSystem, "imgs/blank-bg.png")
	gameOverMessage = loadImage(FileSystem, "imgs/game-over.png")

	levels = []*Level{
		{"Goo Alley", false, gooAlley, 500, 600, 625, 375, []string{"Entering Goo Alley", "Goo Alley destroyed you", "With a renewed disgust, you exit Goo Alley."}, levelBG, levelMap},
		{"Yikesful Mountain", false, yikesfulMountain, 300, 300, 625, 375, []string{"Approaching Yikesful Mountain", "...yikes.", "Shaking your head, you successfully leave Yikesful Mountain behind you."}, levelBG, levelMap},
	}

}

func populate(vsx int, vsy int) { // pass level name or index number as a parameter, or change to method with *Level as receiver...
	// empty lists first, in case any left over from previous level attempt
	for i, h := range levelMap[1] {
		x := (i%tileXCount)*tileSize - vsx
		y := (i/tileXCount)*tileSize + vsy
		if h == 5 {
			nh := NewHazard("blob", hazard, 10, x, y, 100)
			hazardList = append(hazardList, nh)
		}
		if h == 6 {
			nc := NewCreature("teen yorp", creature, x, y, 100, 100, "teen yorp")
			creatureList = append(creatureList, nc)
		}
	}
}

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Starting up Mona Game POC...")

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Mona Game, POC: Movement in Level Space")

	g := NewGame()
	levelSetup(mona.view.xCoord, mona.view.yCoord)
	//	g.Setup()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// Game contains all relevant data for game
type Game struct {
	mode       Mode
	background *ebiten.Image
	count      int
	//lvlComplete   []string // list of names of levels that have been completed
	//lvlCurrent    string
	lvlComplete   []int // list of names of levels that have been completed
	lvlCurrent    int
	items         []string // change type, but to track which items have been collected
	questItem     bool     // deprecate? Could keep, to cheaply track whether to open portal -- maybe rename levelItem
	treasureCount int      // to deprecate -- use score instead: add up different types of treasure, different values
	score         int
}

type Level struct {
	name       string
	complete   bool
	icon       *ebiten.Image
	mapX       int
	mapY       int
	exitX      int
	exitY      int
	message    []string      // on entering level, on death, on successful completion
	background *ebiten.Image // later, this can be []*ebiten.Image, for layered background
	layout     [][]int
}

type Mode int

const (
	Load Mode = iota
	Menu
	World
	Play
)

type Creature struct {
	name        string
	sprite      *ebiten.Image
	facing      int
	xCoord      int
	yCoord      int
	hpCurrent   int
	hpTotal     int
	damage      int
	movement    string // I have no idea how I'm implementing this -- might just key movement style to name, so all same-type creatures move alike
	seesChar    bool
	movementCtr int
	pauseCtr    int
}

func NewGame() *Game {
	log.Printf("Creating new game")
	game := &Game{
		background:    levelBG,
		count:         0,
		questItem:     false,
		treasureCount: 0,
	}
	return game
}

/*
	func (g *Game) levelReset() {
		log.Printf("Resetting level")
		mona.viewReset()
		g.count = 0
		g.questItem = false
		g.treasureCount = 0
	}
*/

func NewCreature(name string, sprite *ebiten.Image, x int, y int, hp int, damage int, movement string) *Creature {
	log.Printf("Creating new creature")
	creature := &Creature{
		name:      name,
		sprite:    sprite,
		facing:    50,
		xCoord:    x,
		yCoord:    y,
		hpCurrent: hp,
		hpTotal:   hp,
		seesChar:  false,
		damage:    damage,
		movement:  name,
	}
	return creature
}

func levelSetup(viewX int, viewY int) {
	populate(viewX, viewY)
}

/*
	func levelComplete() {
		mona.fade()
		end()
	}

	func end() {
		log.Printf("End Screen")
	}

	func (g *Game) over() {
		log.Printf("Game Over")
	}

	func (g *Game) retryLevel() {
		log.Printf("Retry level")
		// levelReset() -- needs fixing
	}
*/
func (g *Game) Update() error {
	g.count++
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = Play
	}
	switch g.mode {
	case Load:
		if g.count > 200 {
			g.mode = Menu
			log.Printf("Changing mode to Menu")
		}
	case Menu:
		g.mode = World
	case World:
		if worldMona.xCoord == 20 {
			worldMona.xCoord = 200
			worldMona.yCoord = 300
			worldMona.view.xCoord = -400
			worldMona.view.yCoord = -500
		}
		radiusCheck := math.Sqrt(math.Pow(float64(worldMona.xCoord-500-worldMona.view.xCoord), 2) + math.Pow(float64(worldMona.yCoord-500-worldMona.view.yCoord), 2))
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			worldMona.facing = 0
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
			worldMona.facing = 48
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
			switch {
			case worldMona.view.yCoord == 0 && worldMona.yCoord > 250:
				worldMona.yCoord += 5
			case worldMona.view.yCoord == -520 && radiusCheck+50 < radius:
				worldMona.yCoord += 5
			case worldMona.view.yCoord > -520:
				worldMona.view.yCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = Play
		}
	case Play:
		if mona.hpCurrent == 0 {
			if mona.lives == 0 {
				//	g.over()
				log.Printf("Player ran out of lives - Game Over")
			}
			//		g.retryLevel()
			log.Printf("Player still has lives -- retry level or return to world?")
		}
		portalFrame = (g.count / 5) % 5
		treasureFrame = (g.count / 5) % 7
		questItemFrame = (g.count / 5) % 5
		hazardFrame = (g.count / 5) % 10
		creatureFrame = (g.count / 5) % 5
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
			monaSide := (mona.xCoord - mona.view.xCoord + 48 + 1) / 50
			monaTop := (mona.yCoord - mona.view.yCoord) / 50
			if levelMap[0][monaTop*tileXCount+monaSide] == 1 /* || levelMap[0][monaBase*tileXCount+monaSide] == 1*/ {
				//		log.Printf("There is a wall here!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				mona.xCoord -= 5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			mona.facing = 48
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
				//		log.Printf("Oh a different wall")
				mona.xCoord += 5
			}
		}
		if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) || inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
			currentFrame = defaultFrame
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			// diagnostics to log!
			diagnosticMap := ""
			for _, v := range levelMap[1] {
				sv := strconv.Itoa(v)
				diagnosticMap += sv
			}
			log.Printf(diagnosticMap)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && mona.yVelo == gravity {
			mona.yVelo = -19
			//		log.Printf("JUMP! velo = %v\nyCoord = %v, yVelo = %v", mona.yVelo, mona.yCoord, mona.yVelo)
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
				monaBase := (mona.yCoord - mona.view.yCoord + 48 + 1) / 50 // checks immediately BELOW base of sprite
				monaLeft := (mona.xCoord - mona.view.xCoord) / 50
				monaRight := (mona.xCoord - mona.view.xCoord + 48) / 50
				if levelMap[0][(monaBase)*tileXCount+monaLeft] == 1 || levelMap[0][(monaBase)*tileXCount+monaRight] == 1 {
					//			log.Printf("THERE IS A TILE THERE WHILE I AM FALLING")
					mona.yCoord = (monaBase * 50) - 50 + mona.view.yCoord
					mona.yVelo = gravity
				}
			}
		}
		monaTop := (mona.yCoord - mona.view.yCoord) / 50
		monaBase := (mona.yCoord - mona.view.yCoord + 48 + 1) / 50 // checks immediately BELOW base of sprite
		monaLeft := (mona.xCoord - mona.view.xCoord) / 50
		monaRight := (mona.xCoord - mona.view.xCoord + 48) / 50
		// basic gravity fixer, doesn't address jumping
		if levelMap[0][(monaBase)*tileXCount+monaLeft] != 1 && levelMap[0][(monaBase)*tileXCount+monaRight] != 1 {
			if mona.yVelo == gravity {
				mona.yCoord += 3 // should be gravity, but that lowers too much
				// BST to quickly assess where landing could happen?
			}
		}
		blockTopLeft := monaTop*tileXCount + monaLeft
		btlVal := levelMap[1][blockTopLeft]
		blockTopRight := monaTop*tileXCount + monaRight
		btrVal := levelMap[1][blockTopRight]
		blockBaseLeft := monaBase*tileXCount + monaLeft
		bblVal := levelMap[1][blockBaseLeft]
		blockBaseRight := monaBase*tileXCount + monaRight
		bbrVal := levelMap[1][blockBaseRight]
		if btlVal == 5 || bblVal == 5 || btrVal == 5 || bbrVal == 5 {
			mona.death()
			if mona.lives == 0 {
				log.Printf("Game over, no lives left")
				//g.over()
			}
			log.Printf("Retry Level or return to world map?")
			//g.retryLevel()

		}
		if btlVal == 3 || btlVal == 4 {
			switch {
			case btlVal == 3:
				g.questItem = true
			case btlVal == 4:
				g.treasureCount += 1
			}
			levelMap[1][blockTopLeft] = 0
			blank.Clear()
		}
		if btrVal == 3 || btrVal == 4 {
			switch {
			case btrVal == 3:
				g.questItem = true
			case btrVal == 4:
				g.treasureCount += 1
			}
			levelMap[1][blockTopRight] = 0
			blank.Clear()
		}
		if bblVal == 3 || bblVal == 4 {
			switch {
			case bblVal == 3:
				g.questItem = true
			case bblVal == 4:
				g.treasureCount += 1
			}
			levelMap[1][blockBaseLeft] = 0
			blank.Clear()
		}
		if bbrVal == 3 || bbrVal == 4 {
			switch {
			case bbrVal == 3:
				g.questItem = true
			case bbrVal == 4:
				g.treasureCount += 1
			}
			levelMap[1][blockBaseRight] = 0
			blank.Clear()
		}

		if btlVal == 2 && bblVal == 2 || btrVal == 2 && bbrVal == 2 {
			//levelComplete()
			log.Printf("Level Complete")
		}

		if (mona.xCoord > levels[g.lvlCurrent].exitX && mona.xCoord < levels[g.lvlCurrent].exitX+50 || mona.xCoord+48 > levels[g.lvlCurrent].exitX && mona.xCoord+48 < levels[g.lvlCurrent].exitX+50) && (mona.yCoord > levels[g.lvlCurrent].exitY && mona.yCoord < levels[g.lvlCurrent].exitY+100 || mona.yCoord+48 > levels[g.lvlCurrent].exitY && mona.yCoord+48 < levels[g.lvlCurrent].exitY+100) {
			g.lvlComplete = append(g.lvlComplete, g.lvlCurrent)
			log.Print("Just hit the portal")
			//levelComplete()
			log.Printf("Level complete")
		}

		// creature movements
		for _, c := range creatureList {
			switch {
			case c.movementCtr > 0:
				// keep moving same dir
				c.movementCtr--
				if c.facing == 0 && c.xCoord <= 3 {
					c.movementCtr = 0
				} else if c.facing == 0 && c.xCoord > 3 {
					c.xCoord -= 3
				} else if c.facing == 50 && c.xCoord >= 597 {
					c.movementCtr = 0
				} else if c.xCoord < 597 {
					c.xCoord += 3
				}
			case c.seesChar == true:
				// rampage towards char
				if c.facing == 0 {
					c.xCoord -= 10
				} else {
					c.xCoord += 10
				}
			case c.pauseCtr > 0:
				// pause
				if c.pauseCtr%9 == 0 {
					c.facing = rand.Intn(2) * 50
				}
				c.pauseCtr--
			default:
				// reset random
				c.movementCtr = rand.Intn(50) + 20
				c.pauseCtr = rand.Intn(40) + 20
				c.facing = rand.Intn(2) * 50
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cmsg := "Creatures\n"
	switch g.mode {
	case Load:
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(ebitengineSplash, op)
	case World:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(worldMona.view.xCoord), float64(worldMona.view.yCoord))
		screen.DrawImage(world, op)
		for _, l := range levels {
			lop := &ebiten.DrawImageOptions{}
			lop.GeoM.Translate(float64(l.mapX), float64(l.mapY))
			world.DrawImage(l.icon, lop)
		}
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(worldMona.xCoord), float64(worldMona.yCoord))
		screen.DrawImage(worldMona.sprite.SubImage(image.Rect(0, 0, 50, 50)).(*ebiten.Image), op)
	case Play:
		lvlOp := &ebiten.DrawImageOptions{}
		lvlOp.GeoM.Translate(float64(mona.view.xCoord), float64(mona.view.yCoord))
		screen.DrawImage(g.background, lvlOp)
		screen.DrawImage(blank, lvlOp)
		mOp := &ebiten.DrawImageOptions{}
		mOp.GeoM.Translate(float64(mona.xCoord), float64(mona.yCoord))
		cx, cy := currentFrame*frameWidth, mona.facing
		screen.DrawImage(mona.sprite.SubImage(image.Rect(cx, cy, cx+frameWidth, cy+frameHeight)).(*ebiten.Image), mOp)
		//	emptyGridSpot := 0
		for _, l := range levelMap {
			for i, t := range l {
				switch {
				case t == 1:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64((i/tileXCount)*tileSize))
					g.background.DrawImage(basicBrick.sprite, top)
				case t == 2:
					if g.questItem == true {
						top := &ebiten.DrawImageOptions{}
						top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
						px := portalFrame * 100
						blank.DrawImage(portal.SubImage(image.Rect(px, 0, px+100, 150)).(*ebiten.Image), top)
					}
				case t == 3:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
					qx := questItemFrame * 50
					blank.DrawImage(questItem.SubImage(image.Rect(qx, 0, qx+50, 50)).(*ebiten.Image), top)
					//		log.Printf("quest loc - x: %v, y: %v", float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
				case t == 4:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize+5), float64(i/tileXCount*tileSize+5))
					tx := treasureFrame * 40
					blank.DrawImage(treasure.SubImage(image.Rect(tx, 0, tx+40, 40)).(*ebiten.Image), top)
				}
			}
		}
		hmsg := "Hazards\n"
		for _, h := range hazardList {
			hop := &ebiten.DrawImageOptions{}
			hop.GeoM.Translate(float64(h.xCoord), float64(h.yCoord))
			hx := hazardFrame * 50
			screen.DrawImage(h.sprite.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image), hop)
			hmsg += fmt.Sprintf("x: %v, y: %v\n", h.xCoord, h.yCoord)
		}

		for _, c := range creatureList {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c.xCoord), float64(c.yCoord))
			cx, cy := creatureFrame*50, c.facing
			screen.DrawImage(c.sprite.SubImage(image.Rect(cx, cy, cx+50, cy+50)).(*ebiten.Image), op)
			cmsg += fmt.Sprintf("- x: %v, y: %v, facing: %v\n", c.xCoord, c.yCoord, c.facing)
		}

		if mona.lives <= 0 {
			overOp := &ebiten.DrawImageOptions{}
			screen.DrawImage(gameOverMessage, overOp)
		}
		// upper left - quest item acquired
		// upper right - treasure count

		//	scoreTreasure := "Treasure: " + strconv.Itoa(g.treasureCount)
		//	text.Draw(screen, scoreTreasure, mplusNormalFont, 300, 140, color.White)
	}
	msg := ""
	msg += cmsg
	//msg += hmsg
	//	msg += fmt.Sprintf("Is screen cleared every frame? %v\n", ebiten.IsScreenClearedEveryFrame())
	//	msg += fmt.Sprintf("Empty Grid Spots: %d\n", emptyGridSpot)
	//	msg += fmt.Sprintf("Mona xCoord: %d\n", mona.xCoord)
	//	msg += fmt.Sprintf("Mona yCoord: %d\n", mona.yCoord)
	//	msg += fmt.Sprintf("Viewer xCoord: %d\n", g.view.xCoord)
	//	msg += fmt.Sprintf("Viewer yCoord: %d\n", g.view.yCoord)
	msg += fmt.Sprintf("Treasure Count: %d\n", g.treasureCount)
	msg += fmt.Sprintf("Quest Item Acquired: %v\n", g.questItem)
	msg += fmt.Sprintf("Lives: %v\n", mona.lives)

	//msg += fmt.Sprintf("Mona Facing: %s\n", mona.facing)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}
