// Package main runs game
package main

import (
	"embed"
	"fmt"
	"image"
	"log"
	"strconv"

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
	FileSystem embed.FS

	levelBG *ebiten.Image
	//	charaSprite *ebiten.Image
	spriteSheet *ebiten.Image
	brick       *ebiten.Image
	portal      *ebiten.Image
	treasure    *ebiten.Image
	questItem   *ebiten.Image
	blank       *ebiten.Image

	levelWidth  int
	levelHeight int

	levelView *Viewer

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
			0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 4, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	}
)

var (
	currentFrame   int
	treasureFrame  int
	questItemFrame int
)

func init() {
	log.Printf("Initializing...")

	levelBG = loadImage(FileSystem, "imgs/level-background--test.png")
	// these values are temporarily hard-coded, replace magic numbers later
	levelWidth = 800
	levelHeight = 600

	levelView = NewViewer()

	spriteSheet = loadImage(FileSystem, "imgs/walk-test--2023-01-03--lr.png")
	currentFrame = defaultFrame
	treasureFrame = defaultFrame
	questItemFrame = defaultFrame
	mona = NewCharacter("Mona", spriteSheet, 100)

	brick = loadImage(FileSystem, "imgs/brick--test.png")
	basicBrick = NewBrick("basic", brick)

	portal = loadImage(FileSystem, "imgs/portal--test.png")
	treasure = loadImage(FileSystem, "imgs/treasure--test.png")
	questItem = loadImage(FileSystem, "imgs/quest-item--test.png")
	blank = loadImage(FileSystem, "imgs/blank-bg.png")
}

// main sets up game and runs it, or returns error
func main() {

	log.Printf("Running level test")

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Mona Game, POC: Movement in Level Space")

	g := NewGame()
	//	g.Setup()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// Game contains all relevant data for game
type Game struct {
	background    *ebiten.Image
	view          *Viewer
	count         int
	questItem     bool
	treasureCount int
}

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
	facing     int
	xCoord     int
	yCoord     int
	yVelo      int
	active     bool
	hp_current int
	hp_total   int
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

func NewGame() *Game {
	log.Printf("Creating new game")
	game := &Game{
		background:    levelBG,
		view:          levelView,
		count:         0,
		questItem:     false,
		treasureCount: 0,
	}
	return game
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

func NewCharacter(name string, sprite *ebiten.Image, hp int) *Character {
	log.Printf("Creating new character %s", name)
	character := &Character{
		name:       name,
		sprite:     sprite,
		facing:     0,
		xCoord:     20,
		yCoord:     ground,
		yVelo:      gravity,
		active:     false,
		hp_current: hp,
		hp_total:   hp,
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

func (g *Game) Update() error {
	g.count++
	treasureFrame = (g.count / 5) % 7
	questItemFrame = (g.count / 5) % 5
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		mona.facing = 0
		currentFrame = (g.count / 5) % frameCount
		switch {
		case g.view.xCoord == 0 && mona.xCoord < 290:
			mona.xCoord += 5
		case g.view.xCoord == -200 && mona.xCoord < 530:
			mona.xCoord += 5
		case g.view.xCoord > -200:
			g.view.xCoord -= 5
		}
		monaSide := (mona.xCoord - g.view.xCoord + 48 + 1) / 50
		monaTop := (mona.yCoord - g.view.yCoord) / 50
		if levelMap[0][monaTop*tileXCount+monaSide] == 1 /* || levelMap[0][monaBase*tileXCount+monaSide] == 1*/ {
			//		log.Printf("There is a wall here!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			mona.xCoord -= 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		mona.facing = 48
		currentFrame = (g.count / 5) % frameCount
		switch {
		case g.view.xCoord == -200 && mona.xCoord > 290:
			mona.xCoord -= 5
		case g.view.xCoord == 0 && mona.xCoord > 40:
			mona.xCoord -= 5
		case g.view.xCoord < 0:
			g.view.xCoord += 5
		}
		monaSide := (mona.xCoord - g.view.xCoord) / 50
		monaTop := (mona.yCoord - g.view.yCoord) / 50
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
		case mona.yCoord < 160 && g.view.yCoord-mona.yVelo < 0 && mona.yVelo < 0:
			//	mona.yCoord -= mona.yVelo
			g.view.yCoord -= mona.yVelo
		case mona.yCoord > 160 && g.view.yCoord-mona.yVelo > -120 && mona.yVelo > 0:
			//	mona.yCoord -= mona.yVelo
			g.view.yCoord -= mona.yVelo
		default:
			mona.yCoord += mona.yVelo
		}
		mona.yVelo += 1

		if mona.yVelo >= 0 {
			monaBase := (mona.yCoord - g.view.yCoord + 48 + 1) / 50 // checks immediately BELOW base of sprite
			monaLeft := (mona.xCoord - g.view.xCoord) / 50
			monaRight := (mona.xCoord - g.view.xCoord + 48) / 50
			if levelMap[0][(monaBase)*tileXCount+monaLeft] == 1 || levelMap[0][(monaBase)*tileXCount+monaRight] == 1 {
				//			log.Printf("THERE IS A TILE THERE WHILE I AM FALLING")
				mona.yCoord = (monaBase * 50) - 50 + g.view.yCoord
				mona.yVelo = gravity
			}
		}
	}
	monaTop := (mona.yCoord - g.view.yCoord) / 50
	monaBase := (mona.yCoord - g.view.yCoord + 48 + 1) / 50 // checks immediately BELOW base of sprite
	monaLeft := (mona.xCoord - g.view.xCoord) / 50
	monaRight := (mona.xCoord - g.view.xCoord + 48) / 50
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
	if btlVal != 0 {
		switch {
		case btlVal == 3:
			g.questItem = true
		case btlVal == 4:
			g.treasureCount += 1
		}
		levelMap[1][blockTopLeft] = 0
		blank.Clear()
	}
	if btrVal != 0 {
		switch {
		case btrVal == 3:
			g.questItem = true
		case btrVal == 4:
			g.treasureCount += 1
		}
		levelMap[1][blockTopRight] = 0
		blank.Clear()
	}
	if bblVal != 0 {
		switch {
		case bblVal == 3:
			g.questItem = true
		case bblVal == 4:
			g.treasureCount += 1
		}
		levelMap[1][blockBaseLeft] = 0
		blank.Clear()
	}
	if bbrVal != 0 {
		switch {
		case bbrVal == 3:
			g.questItem = true
		case bbrVal == 4:
			g.treasureCount += 1
		}
		levelMap[1][blockBaseRight] = 0
		blank.Clear()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	lvlOp := &ebiten.DrawImageOptions{}
	lvlOp.GeoM.Translate(float64(g.view.xCoord), float64(g.view.yCoord))
	screen.DrawImage(g.background, lvlOp)
	screen.DrawImage(blank, lvlOp)
	mOp := &ebiten.DrawImageOptions{}
	mOp.GeoM.Translate(float64(mona.xCoord), float64(mona.yCoord))
	cx, cy := currentFrame*frameWidth, mona.facing
	screen.DrawImage(mona.sprite.SubImage(image.Rect(cx, cy, cx+frameWidth, cy+frameHeight)).(*ebiten.Image), mOp)
	emptyGridSpot := 0
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
				top := &ebiten.DrawImageOptions{}
				top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
				g.background.DrawImage(portal.SubImage(image.Rect(0, 0, 100, 150)).(*ebiten.Image), top)
			case t == 3:
				top := &ebiten.DrawImageOptions{}
				top.GeoM.Translate(float64((i%tileXCount)*tileSize), float64(i/tileXCount*tileSize))
				qx := questItemFrame * 50
				blank.DrawImage(questItem.SubImage(image.Rect(qx, 0, qx+50, 50)).(*ebiten.Image), top)
			case t == 4:
				top := &ebiten.DrawImageOptions{}
				top.GeoM.Translate(float64((i%tileXCount)*tileSize+5), float64(i/tileXCount*tileSize+5))
				tx := treasureFrame * 40
				blank.DrawImage(treasure.SubImage(image.Rect(tx, 0, tx+40, 40)).(*ebiten.Image), top)
			}
		}
	}

	msg := ""
	msg += fmt.Sprintf("Is screen cleared every frame? %v", ebiten.IsScreenClearedEveryFrame())
	msg += fmt.Sprintf("Empty Grid Spots: %d\n", emptyGridSpot)
	msg += fmt.Sprintf("Mona xCoord: %d\n", mona.xCoord)
	msg += fmt.Sprintf("Mona yCoord: %d\n", mona.yCoord)
	msg += fmt.Sprintf("Viewer xCoord: %d\n", g.view.xCoord)
	msg += fmt.Sprintf("Viewer yCoord: %d\n", g.view.yCoord)
	msg += fmt.Sprintf("Treasure Count: %d\n", g.treasureCount)
	msg += fmt.Sprintf("Quest Item Acquired: %v\n", g.questItem)
	//msg += fmt.Sprintf("Mona Facing: %s\n", mona.facing)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}
