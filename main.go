package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	//screenWidth           = 640
	//screenHeight          = 480
	GAME_DEFAULT_SCREEN_W = 640
	GAME_DEFAULT_SCREEN_H = 480

	//tileSize            = 32
	titleFontSize                           = fontSize * 1.5
	fontSize                                = 24
	startBallX                              = 250
	startBallY                              = 300
	ballInitVelX                            = 4
	ballInitVelY                            = 4
	TPS                               int64 = 60
	msPerTick                         int64 = 1000 / TPS
	GAME_MAP_ROWS                           = 30
	GAME_MAP_COLS                           = 30
	GAME_FLIP_LIVES_FROM_POINTS_EVERY       = 100
	GAME_START_LEVEL                        = 0
	PLAYER_START_POS_X                      = 250
	PLAYER_START_POS_Y                      = 210
	GAME_LEVEL_DATA_DIR                     = "leveldata"
	GAME_DATA_MATRIX_END                    = ".csv"
	GAME_TILE_SIZE                          = 50
	GAME_HBAR_X                             = 20
	GAME_HBAR_Y                             = 30
	GAME_HBAR_W                             = 200
	GAME_HBAR_H                             = 25
	GAME_INI_FILE                           = "settings.ini"
)

var (
	playerImage         *ebiten.Image
	tileImage           *ebiten.Image
	arcadeFaceSource    *text.GoTextFaceSource
	score               int
	lives               int
	level               int
	timerLastTimeMillis int64
	lastPauseTime       int64 = 0
	pauseIntervalMillis int64 = 1000
	commandLineArgs     []string
	panelHeight         int
	panelWidth          int
	midpointX           int
	midpointY           int
	centerTextX         = panelWidth/2 - 50
	centerTextY         = panelHeight/2 - 50
)

// create Enum for editing maps mode
type EditMode int

const (
	EditNone EditMode = iota
	EditTile
	EditEntity
	EditFidget
	EditDecor
	EditPickup
	EditInteractive
	EditZone
	EditSpawner
	EditPlatform
)

var img *ebiten.Image

func init() {
	var err error
	//img, _, err = ebitenutil.NewImageFromFile("player.png")
	if err != nil {
		log.Fatal(err)
	}
}

type editable interface {
	AddInstanceToGrid(int, int, int)
	CycleAssetKind(int)
	getAssetID() int
	setAssetID(int)
}

type Game struct {
	player *Player
	input  *Input

	tileMap *TileMap
	gameComponents
	score                           int
	lives                           int
	level                           int
	screenWidth                     int
	screenHeight                    int
	levelCompleteScreenActive       bool
	mode                            int
	centerMarqueeEndActionsComplete bool
	editMode                        EditMode
	godMode                         bool
	activateObject                  bool
	exePath                         string
	//sound

}

// mode 0 = gameplay
// mode 1 = level complete
// mode 2 = pause

type gameComponents struct {
	scoreString       *RasterString
	livesString       *RasterString
	stageString       *RasterString
	centerText        *RasterString
	centerMarquee     *Marquee
	projectileManager *ProjectileManager
	pickupManager     *PickupManager
	editor            *Editor
	console           *Console
	fidgetManager     *FidgetManager
	warpManager       *WarpManager
	entityManager     *EntityManager
	platformManager   *PlatformManager
	decorManager      *DecorManager
	audioPlayer       *AudioPlayer
	healthBar         *Bar
	background        *Background
	particleManager   *ParticleManager
}

func (g *Game) Update() error {
	now := time.Now()
	timerNowMillis := now.UnixNano()
	if abs(timerNowMillis-timerLastTimeMillis) > msPerTick {
		//g.activateObject = false
		g.background.Update()
		g.input.Update()
		g.player.Update()
		//g.ball.Update()
		//g.brickGrid.Update()
		g.tileMap.Update()
		g.updateRasterStrings()
		timerLastTimeMillis = timerNowMillis
		g.centerMarquee.Update()
		g.projectileManager.Update()
		g.pickupManager.Update()
		g.console.Update()
		g.fidgetManager.Update()
		g.entityManager.Update()
		g.platformManager.Update()
		g.decorManager.Update()
		g.editor.Update()
		g.particleManager.Update()

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.DrawImage(img, nil)

	//g.ball.Draw(screen)
	//g.brickGrid.Draw(screen)
	g.background.Draw(screen)
	g.tileMap.Draw(screen)
	g.drawRasterStrings(screen)
	g.centerMarquee.Draw(screen)
	g.fidgetManager.Draw(screen)
	g.player.Draw(screen)
	g.projectileManager.Draw(screen)
	g.pickupManager.Draw(screen)
	g.console.Draw(screen)
	g.entityManager.Draw(screen)
	g.platformManager.Draw(screen)
	g.decorManager.Draw(screen)

	g.particleManager.Draw(screen)
	g.healthBar.Draw(screen)
	g.editor.Draw(screen)

}

func (g *Game) Pause() {
	//now := time.Now()
	timerNowMillis := time.Now().UnixNano()
	//fmt.Printf("time diff %d \n", abs(timerNowMillis-lastPauseTime))
	if abs(timerNowMillis-lastPauseTime) > 1000000000 {
		lastPauseTime = timerNowMillis
		if g.mode == 0 { //paused
			g.mode = 2
			g.centerText.stringContent = "Paused"
			g.centerText.visible = true
			g.player.frozen = true
			//g.ball.frozen = true
		} else if g.mode == 2 {
			g.mode = 0
			g.centerText.visible = false
			g.player.frozen = false
			//g.ball.frozen = false
		}

	} else {
		//lastPauseTime = timerNowMillis
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func NewGame() ebiten.Game {
	g := &Game{}
	g.initGameComponents()
	g.newGameSession()
	g.exePath = commandLineArgs[0]
	//g.player.init()

	return g
}

func (g *Game) levelComplete() {
	//g.centerMarquee.centerTextOffset = 25
	g.centerMarquee.UpdateText("Level Complete")
	g.centerMarquee.Start()
	g.mode = 1
	//g.ball.visible = false
	nextLevel := g.level + 1
	if nextLevel > g.level {
		g.level += 1
	}
	//g.level += 1
	oneMoreLife := g.lives + 1
	if oneMoreLife > g.lives {
		g.updateLives(g.lives + 1)

	}

}

func (g *Game) updateLives(lives int) {
	g.lives = lives
	g.livesString.stringContent = fmt.Sprintf("Lives: %d", g.lives)
}

func (g *Game) updateLevel(level int) {
	g.level = level
	g.stageString.stringContent = fmt.Sprintf("Level: %d", g.level)
}

func (g *Game) updateScore(score int) {
	g.score = score
	g.scoreString.stringContent = fmt.Sprintf("Score: %d", g.score)
}

func (g *Game) marqueeMessageComplete() {
	if g.mode == 1 {
		g.mode = 0
		levelText := fmt.Sprintf("Level: %d", g.level)
		//g.centerMarquee.centerTextOffset = 40
		g.centerMarquee.UpdateText(levelText)
		g.gameComponents.stageString.stringContent = levelText
		g.centerMarquee.Start()
		g.centerMarqueeEndActionsComplete = true
		g.levelCompleteScreenActive = false
		//g.brickGrid.Reset()
		//g.ball.reset()
	} else if g.mode == 0 {
		//g.brickGrid.Reset()
		//g.ball.visible = true
		//g.ball.reset()
		//g.ball.frozen = false
		//fmt.Println("reset")
		g.centerMarqueeEndActionsComplete = true
		g.levelCompleteScreenActive = false

	}

}

func (g *Game) initGameComponents() {
	startY := PLAYER_START_POS_Y
	startX := PLAYER_START_POS_X
	g.activateObject = false
	g.screenHeight = panelHeight
	g.screenWidth = panelWidth
	//g.player = &Player{}
	g.level = GAME_START_LEVEL
	g.player = NewPlayer(g, startX, startY)
	g.input = &Input{}
	g.input.init(g)
	//g.ball = &Ball{}
	//g.ball.init(g, startBallX, startBallY, ballInitVelX, ballInitVelY)
	g.tileMap = NewTileMap(g)
	g.projectileManager = NewProjectileManager(g, 0)
	g.pickupManager = NewPickupManager(g)
	g.initRasterStrings()
	g.centerMarquee = NewMarquee(g, 0, 200, "Cooking with gas!")
	g.editor = NewEditor(g)
	g.console = NewConsole(g)
	g.fidgetManager = NewFidgetManager(g)
	g.warpManager = NewWarpManager(g)
	g.entityManager = NewEntityManager(g)
	g.platformManager = NewPlatformManager(g)
	g.audioPlayer = NewAudioPlayer(g)
	g.decorManager = NewDecorManager(g)
	g.background = NewBackground(g)
	g.healthBar = NewBar(GAME_HBAR_X, GAME_HBAR_Y, GAME_HBAR_W, GAME_HBAR_H)
	g.particleManager = NewParticleManager(g)
	g.score = 0
	g.lives = 3

}

func (g *Game) newGameSession() {

	g.mode = 0
	timerLastTimeMillis = time.Now().UnixNano()
	g.editMode = EditTile
	g.godMode = false
	g.loadLevel((GAME_START_LEVEL))

}

func (g *Game) initRasterStrings() {

	g.gameComponents.scoreString = NewRasterString(g, "Score: 0", 10, 10)

	g.gameComponents.centerText = NewRasterString(g, "You won", centerTextX, centerTextY)
	g.gameComponents.centerText.visible = false
	levelStr := fmt.Sprintf("Level: %d", g.level)
	g.gameComponents.stageString = NewRasterString(g, levelStr, g.screenWidth-90, 10)
	g.gameComponents.livesString = NewRasterString(g, "Lives: 3", (g.screenWidth/2)-50, 10)

}

func (g *Game) incrementScore(points int) {
	threshold := GAME_FLIP_LIVES_FROM_POINTS_EVERY
	oldScore := g.score
	newScore := g.score + points
	if fracOld, fracNew := oldScore/threshold, newScore/threshold; fracNew-fracOld == 1 {
		g.updateLives(g.lives + 1)
	}
	if newScore > oldScore {
		g.updateScore(newScore)
	}

}

func (g *Game) drawRasterStrings(screen *ebiten.Image) {
	g.gameComponents.scoreString.Draw(screen)
	g.gameComponents.centerText.Draw(screen)
	g.gameComponents.stageString.Draw(screen)
	g.gameComponents.livesString.Draw(screen)
}
func (g *Game) updateRasterStrings() {
	g.gameComponents.scoreString.Update()
	g.gameComponents.centerText.Update()
	g.gameComponents.stageString.Update()
	g.gameComponents.livesString.Update()
}

func (g *Game) loadLevel(level int) {
	g.updateLevel(level)
	fmt.Println("Load level ", level)
	g.tileMap.loadCurrentLevelMapFromFile()
	g.pickupManager.loadDataFromFile()
	g.fidgetManager.loadDataFromFile()
	g.entityManager.loadDataFromFile()
	g.entityManager.addLevelBoss()
	g.platformManager.loadDataFromFile()
	g.decorManager.loadDataFromFile()

}

func setCheats() {

}

func setResolution() (int, int) {
	const flag = "resolution"
	var foundFlagPos = -1
	argsLen := len(commandLineArgs)
	lastIndex := argsLen - 1
	if argsLen >= 4 {

		for i, _ := range commandLineArgs {
			if strings.TrimSpace(commandLineArgs[i]) == flag {
				foundFlagPos = i
				break
			}
		}
		if lastIndex-2 < foundFlagPos {
			return GAME_DEFAULT_SCREEN_W, GAME_DEFAULT_SCREEN_H
		}
		numA, errA := strconv.Atoi(strings.TrimSpace(commandLineArgs[foundFlagPos+1]))
		numB, errB := strconv.Atoi(strings.TrimSpace(commandLineArgs[foundFlagPos+2]))
		if errA != nil || errB != nil || numA <= 0 || numB <= 0 {

			return GAME_DEFAULT_SCREEN_W, GAME_DEFAULT_SCREEN_H
		} else {
			return numA, numB
		}
	} else {

		return GAME_DEFAULT_SCREEN_W, GAME_DEFAULT_SCREEN_H
	}
}

func main() {
	commandLineArgs = os.Args
	panelWidth, panelHeight = setResolution()
	dataMap, err := openOrCreateDefaultIni(GAME_INI_FILE)
	if err != nil {
		log.Println("failed to open ini")
	}
	_ = dataMap
	pprintMap(dataMap)

	ebiten.SetWindowSize(panelWidth, panelHeight)
	ebiten.SetWindowTitle("Flamethrower Cat")

	var game = NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
