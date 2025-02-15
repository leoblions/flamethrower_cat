package main

import (
	_ "embed"
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	centerTextX  = screenWidth/2 - 50
	centerTextY  = screenHeight/2 - 50
	//tileSize            = 32
	titleFontSize                           = fontSize * 1.5
	fontSize                                = 24
	startBallX                              = 250
	startBallY                              = 300
	ballInitVelX                            = 4
	ballInitVelY                            = 4
	TPS                               int64 = 60
	msPerTick                         int64 = 1000 / TPS
	mapRows                                 = 30
	mapCols                                 = 30
	GAME_FLIP_LIVES_FROM_POINTS_EVERY       = 100
	GAME_START_LEVEL                        = 0
	PLAYER_START_POS_X                      = 250
	PLAYER_START_POS_Y                      = 210
	GAME_LEVEL_DATA_DIR                     = "leveldata"
	GAME_DATA_MATRIX_END                    = ".csv"
	GAME_TILE_SIZE                          = 50
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
	ball   *Ball
	//brickGrid *BrickGrid
	tileMap *TileMap
	rasterStrings
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
	//sound
	audioContext *audio.Context
	//jumpPlayer   *audio.Player
	//hitPlayer    *audio.Player
	soundEffectPlayers map[string]*audio.Player
}

func playSound(plr *audio.Player) error {
	if err := plr.Rewind(); err != nil {
		return err
	}
	plr.Play()
	return nil
}

// mode 0 = gameplay
// mode 1 = level complete
// mode 2 = pause

type rasterStrings struct {
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
}

func (g *Game) Update() error {
	now := time.Now()
	timerNowMillis := now.UnixNano()
	if abs(timerNowMillis-timerLastTimeMillis) > msPerTick {
		//g.activateObject = false
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
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.DrawImage(img, nil)

	//g.ball.Draw(screen)
	//g.brickGrid.Draw(screen)
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
	g.init()
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
		g.rasterStrings.stageString.stringContent = levelText
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

func (g *Game) init() {
	startY := PLAYER_START_POS_Y
	startX := PLAYER_START_POS_X
	g.activateObject = false
	g.screenHeight = screenHeight
	g.screenWidth = screenWidth
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
	g.initAudioPlayers()
	g.score = 0
	g.lives = 3

	g.mode = 0
	timerLastTimeMillis = time.Now().UnixNano()
	g.editMode = EditTile
	g.godMode = false
	g.loadLevel((GAME_START_LEVEL))

}

func (g *Game) initRasterStrings() {

	g.rasterStrings.scoreString = NewRasterString(g, "Score: 0", 10, 10)

	g.rasterStrings.centerText = NewRasterString(g, "You won", centerTextX, centerTextY)
	g.rasterStrings.centerText.visible = false
	levelStr := fmt.Sprintf("Level: %d", g.level)
	g.rasterStrings.stageString = NewRasterString(g, levelStr, g.screenWidth-90, 10)
	g.rasterStrings.livesString = NewRasterString(g, "Lives: 3", (g.screenWidth/2)-50, 10)

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
	g.rasterStrings.scoreString.Draw(screen)
	g.rasterStrings.centerText.Draw(screen)
	g.rasterStrings.stageString.Draw(screen)
	g.rasterStrings.livesString.Draw(screen)
}
func (g *Game) updateRasterStrings() {
	g.rasterStrings.scoreString.Update()
	g.rasterStrings.centerText.Update()
	g.rasterStrings.stageString.Update()
	g.rasterStrings.livesString.Update()
}

func (g *Game) loadLevel(level int) {
	g.updateLevel(level)
	fmt.Println("Load level ", level)
	g.tileMap.loadCurrentLevelMapFromFile()
	g.pickupManager.loadDataFromFile()
	g.fidgetManager.loadDataFromFile()
	g.entityManager.loadDataFromFile()
	g.platformManager.loadDataFromFile()

}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Flamethrower Cat")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
