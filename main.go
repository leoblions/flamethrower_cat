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
	GAME_DEFAULT_SCREEN_W = 640
	GAME_DEFAULT_SCREEN_H = 480

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
	GAME_START_LIVES                        = 3
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
	GAME_PAUSED_DEBOUNCE                    = 1000000000
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
}

// mode 0 = gameplay
// mode 1 = level complete
// mode 2 = pause

type gameComponents struct {
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
	hud               *Hud
}

func (g *Game) Update() error {
	now := time.Now()
	timerNowMillis := now.UnixNano()
	if abs(timerNowMillis-timerLastTimeMillis) > msPerTick {
		g.background.Update()
		g.input.Update()
		g.player.Update()
		g.tileMap.Update()
		g.hud.Update()
		timerLastTimeMillis = timerNowMillis
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
	g.background.Draw(screen)
	g.tileMap.Draw(screen)

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
	g.hud.Draw(screen)

}

func (g *Game) Pause() {
	timerNowMillis := time.Now().UnixNano()
	if abs(timerNowMillis-lastPauseTime) > GAME_PAUSED_DEBOUNCE {
		lastPauseTime = timerNowMillis
		if g.mode == 0 { //paused
			g.mode = 2
			g.hud.centerText.updateText("Paused")
			g.hud.centerText.visible = true
			g.player.frozen = true
		} else if g.mode == 2 {
			g.mode = 0
			g.hud.centerText.visible = false
			g.player.frozen = false
		}

	} else {
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

	return g
}

func (g *Game) levelComplete() {
	g.hud.centerMarquee.UpdateText("Level Complete")
	g.hud.centerMarquee.Start()
	g.mode = 1
	nextLevel := g.level + 1
	if nextLevel > g.level {
		g.level += 1
	}
	oneMoreLife := g.lives + 1
	if oneMoreLife > g.lives {
		g.updateLivesRelative(1)

	}

}

func (g *Game) updateLivesAbsolute(lives int) {
	g.lives = lives
	stringContent := fmt.Sprintf("Lives: %d", g.lives)
	g.hud.livesString.updateText(stringContent)
}

func (g *Game) updateLivesRelative(deltaLives int) {
	if livesTemp := deltaLives + g.lives; livesTemp >= 0 {
		g.lives = livesTemp
	}
	stringContent := fmt.Sprintf("Lives: %d", g.lives)
	g.hud.livesString.updateText(stringContent)
}

func (g *Game) updateLevel(level int) {
	g.level = level
	stringContent := fmt.Sprintf("Level: %d", g.level)
	g.hud.levelString.updateText(stringContent)
}

func (g *Game) updateScore(score int) {
	g.score = score
	stringContent := fmt.Sprintf("Score: %d", g.score)
	g.hud.scoreString.updateText(stringContent)
}

func (g *Game) marqueeMessageComplete() {
	if g.mode == 1 {
		g.mode = 0
		levelText := fmt.Sprintf("Level: %d", g.level)
		g.hud.centerMarquee.UpdateText(levelText)
		g.hud.centerMarquee.Start()
		g.centerMarqueeEndActionsComplete = true
		g.levelCompleteScreenActive = false
	} else if g.mode == 0 {
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
	g.level = GAME_START_LEVEL
	g.lives = GAME_START_LIVES
	g.player = NewPlayer(g, startX, startY)
	g.input = &Input{}
	g.input.init(g)
	g.tileMap = NewTileMap(g)
	g.projectileManager = NewProjectileManager(g, 0)
	g.pickupManager = NewPickupManager(g)

	g.editor = NewEditor(g)
	g.console = NewConsole(g)
	g.fidgetManager = NewFidgetManager(g)
	g.warpManager = NewWarpManager(g)
	g.entityManager = NewEntityManager(g)
	g.platformManager = NewPlatformManager(g)
	g.audioPlayer = NewAudioPlayer(g)
	g.decorManager = NewDecorManager(g)
	g.background = NewBackground(g)
	g.hud = NewHud(g)
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

func (g *Game) incrementScore(points int) {
	threshold := GAME_FLIP_LIVES_FROM_POINTS_EVERY
	oldScore := g.score
	newScore := g.score + points
	if fracOld, fracNew := oldScore/threshold, newScore/threshold; fracNew-fracOld == 1 {
		g.updateLivesRelative(1)
	}
	if newScore > oldScore {
		g.updateScore(newScore)
	}

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
