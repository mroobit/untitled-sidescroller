// Package main runs game
package main

import (
	"embed"
	"fmt"
	"image"
	"log"
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
	levelBG          *ebiten.Image
	//	charaSprite *ebiten.Image
	spriteSheet     *ebiten.Image
	brick           *ebiten.Image
	portal          *ebiten.Image
	treasure        *ebiten.Image
	questItem       *ebiten.Image
	hazard          *ebiten.Image
	creature        *ebiten.Image
	blank           *ebiten.Image
	gameOverMessage *ebiten.Image

	levelWidth  int
	levelHeight int

	monaView *Viewer

	mona       *Character
	basicBrick *Brick
	levelMap   = [][]int{
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
)

var (
	//	mplusNormalFont font.Face

	currentFrame   int
	portalFrame    int
	treasureFrame  int
	questItemFrame int
	hazardFrame    int
	creatureFrame  int

	hazardList   []*Hazard
	creatureList []*Creature
)

func init() {
	log.Printf("Initializing...")

	rand.Seed(time.Now().UnixNano())

	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")
	// these values are temporarily hard-coded, replace magic numbers later
	levelWidth = 800
	levelHeight = 600

	monaView = NewViewer()

	ebitengineSplash = loadImage(FileSystem, "imgs/load-ebitengine-splash.png")

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")
	currentFrame = defaultFrame
	treasureFrame = defaultFrame
	questItemFrame = defaultFrame
	mona = NewCharacter("Mona", spriteSheet, monaView, 100)

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	basicBrick = NewBrick("basic", brick)

	portal = loadImage(FileSystem, "imgs/portal-b--test.png")
	treasure = loadImage(FileSystem, "imgs/treasure--test.png")
	questItem = loadImage(FileSystem, "imgs/quest-item--test.png")
	hazard = loadImage(FileSystem, "imgs/blob--test.png")
	creature = loadImage(FileSystem, "imgs/creature--test.png")
	blank = loadImage(FileSystem, "imgs/blank-bg.png")
	gameOverMessage = loadImage(FileSystem, "imgs/game-over.png")
}

func hazards(vsx int, vsy int) { // rename to indicate that it's populating map with hazards and creatures and things drawn to screen (rather than bg/etc)
	for i, h := range levelMap[1] {
		x := (i%tileXCount)*tileSize - vsx
		y := (i/tileXCount)*tileSize + vsy
		if h == 5 {
			nh := NewHazard("blob", hazard, x, y, 100)
			hazardList = append(hazardList, nh)
		}
		if h == 6 {
			nc := NewCreature("teen yorp", creature, x, y, 100, 100, "teen yorp")
			creatureList = append(creatureList, nc)
		}
	}
}

/*
func init() {
	fontFile, err := FileSystem.Open("fonts/mplus-1p-regular.ttf")
	if err != nil {
		log.Fatalf("Error opening font: %v\n", err)
	}
	tt, err := opentype.Parse(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
*/

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Running level test")

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
	mode          Mode
	background    *ebiten.Image
	count         int
	questItem     bool
	treasureCount int
	active        bool
}

type Mode int

const (
	Load Mode = iota
	Menu
	World
	Level
)

// Viewer is the part of the total level that is visible, as indicated by the X,Y of the upper left corner
type Viewer struct {
	xCoord int
	yCoord int
	width  int
	height int
}

type Character struct {
	name       string
	sprite     *ebiten.Image
	view       *Viewer
	facing     int
	xCoord     int
	yCoord     int
	yVelo      int
	active     bool
	hp_current int
	hp_total   int
	lives      int
}

type Brick struct {
	name         string
	sprite       *ebiten.Image
	impenetrable bool // can you walk through it
	supportive   bool // can you land on it
	destructible bool // can you destroy it
	//lethal	bool		// will it kill you on contact
	damage int // amount of damage per encounter -- if lethal, set absurdly high
}

type Hazard struct {
	name   string
	sprite *ebiten.Image
	xCoord int
	yCoord int
	damage int
}

type Creature struct {
	name        string
	sprite      *ebiten.Image
	facing      int
	xCoord      int
	yCoord      int
	hp_current  int
	hp_total    int
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
		active:        true,
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
func (c *Character) viewReset() {
	log.Printf("Resetting viewer")
	c.view.xCoord = 0
	c.view.yCoord = winHeight - levelHeight
}

func NewViewer() *Viewer {
	log.Printf("Creating new viewer")
	viewer := &Viewer{
		xCoord: 0,
		yCoord: winHeight - levelHeight,
		width:  winWidth,
		height: winHeight,
	}
	return viewer
}

func NewCharacter(name string, sprite *ebiten.Image, view *Viewer, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:       name,
		sprite:     sprite,
		view:       view,
		facing:     0,
		xCoord:     20,
		yCoord:     ground,
		yVelo:      gravity,
		active:     false,
		hp_current: hp,
		hp_total:   hp,
		lives:      3,
	}
	return character
}

func NewBrick(name string, sprite *ebiten.Image) *Brick {
	log.Printf("Creating new brick")
	brick := &Brick{
		name:         name,
		sprite:       sprite,
		impenetrable: true,
		supportive:   true,
		destructible: false,
		damage:       0,
	}
	return brick
}

func NewHazard(name string, sprite *ebiten.Image, x int, y int, damage int) *Hazard {
	log.Printf("Creating new hazard")
	hazard := &Hazard{
		name:   name,
		sprite: sprite,
		xCoord: x,
		yCoord: y,
		damage: damage,
	}
	return hazard
}

func NewCreature(name string, sprite *ebiten.Image, x int, y int, hp int, damage int, movement string) *Creature {
	log.Printf("Creating new creature")
	creature := &Creature{
		name:       name,
		sprite:     sprite,
		facing:     50,
		xCoord:     x,
		yCoord:     y,
		hp_current: hp,
		hp_total:   hp,
		seesChar:   false,
		damage:     damage,
		movement:   name,
	}
	return creature
}

func levelSetup(viewX int, viewY int) {
	hazards(viewX, viewY)
}

func levelComplete() {
	mona.fade()
	end()
}

func (c *Character) fade() {
	log.Printf("Fade character")
}

func end() {
	log.Printf("End Screen")
}

func (g *Game) over() {
	log.Printf("Game Over")
	g.active = false
}

func (g *Game) retryLevel() {
	log.Printf("Retry level")
	// levelReset() -- needs fixing
}

func (c *Character) death() {
	c.hp_current = 0
	c.lives--
	// initiate character death animation
}

func (g *Game) Update() error {
	g.count++
	if inpututil.IsKeyJustPressed(ebiten.KeyF) { // developer skip-ahead
		g.mode = Level
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
		g.mode = Level
	case Level:
		if mona.hp_current == 0 {
			if mona.lives == 0 {
				g.over()
			}
			g.retryLevel()
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
				g.over()
			}
			g.retryLevel()

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
			levelComplete()
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
	case Level:
		lvlOp := &ebiten.DrawImageOptions{}
		lvlOp.GeoM.Translate(float64(mona.view.xCoord), float64(mona.view.yCoord))
		screen.DrawImage(g.background, lvlOp)
		screen.DrawImage(blank, lvlOp)
		//	boh := &ebiten.DrawImageOptions{}
		//	screen.DrawImage(screenSize, boh)
		mOp := &ebiten.DrawImageOptions{}
		mOp.GeoM.Translate(float64(mona.xCoord), float64(mona.yCoord))
		cx, cy := currentFrame*frameWidth, mona.facing
		screen.DrawImage(mona.sprite.SubImage(image.Rect(cx, cy, cx+frameWidth, cy+frameHeight)).(*ebiten.Image), mOp)
		//	emptyGridSpot := 0
		for _, l := range levelMap {
			for i, t := range l {
				switch {
				//			case t == 0:
				//				emptyGridSpot += 1
				//				top := &ebiten.DrawImageOptions{}
				//				g.background.DrawImage(blank, top)
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
					log.Printf("quest loc - x: %v, y: %v", float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
				case t == 4:
					top := &ebiten.DrawImageOptions{}
					top.GeoM.Translate(float64((i%tileXCount)*tileSize+5), float64(i/tileXCount*tileSize+5))
					tx := treasureFrame * 40
					blank.DrawImage(treasure.SubImage(image.Rect(tx, 0, tx+40, 40)).(*ebiten.Image), top)
				}
			}
		}
		/*
			for _, l := range levelMap {
				for i, t := range l {
					if t == 5 {
						top := &ebiten.DrawImageOptions{}
						top.GeoM.Translate(float64((i%tileXCount)*tileSize-g.view.xCoord), float64((i/tileXCount)*tileSize+g.view.yCoord))
						//top.GeoM.Translate(float64(g.view.xCoord-((i%tileXCount)*tileSize)), float64(g.view.yCoord-(i/tileXCount*tileSize)+50))
						log.Printf("hazard loc - x: %v, y: %v", float64(((i%tileXCount)*tileSize)-g.view.xCoord), float64((i/tileXCount*tileSize)+g.view.yCoord))
						hx := hazardFrame * 50
						screen.DrawImage(hazard.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image), top)
					}
				}
			}
		*/
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

		if g.active == false {
			overOp := &ebiten.DrawImageOptions{}
			screen.DrawImage(gameOverMessage, overOp)
		}
		/*
			noOp := &ebiten.DrawImageOptions{}
			noOp.GeoM.Translate(float64(250), float64(380))
			hx := hazardFrame * 50
			screen.DrawImage(hazard.SubImage(image.Rect(hx, 0, hx+50, 50)).(*ebiten.Image), noOp)
		*/
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
