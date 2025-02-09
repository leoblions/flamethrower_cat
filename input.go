package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	game      *Game
	keys      []ebiten.Key
	dflags    *DirectionFlags
	cursorPos pos
}

const (
	INP_PRINT_WHEN_KEY_NOT_USED = false
)

type DirectionFlags struct {
	up    bool
	down  bool
	left  bool
	right bool
}

type pos struct {
	x int
	y int
}

var mouseButtonpressedBefore = struct {
	left   bool
	middle bool
	right  bool
}{
	left:   false,
	right:  false,
	middle: false,
}

func (df *DirectionFlags) reset() {
	df.up = false
	df.down = false
	df.left = false
	df.right = false
}

func (i *Input) init(g *Game) {
	i.game = g
	i.dflags = &DirectionFlags{}
}

func (inp *Input) Update() {

	mx, my := ebiten.CursorPosition()

	wheelX, wheelY := ebiten.Wheel()

	if wheelX != 0 || wheelY != 0 {
		inp.game.editor.editHandleWheel(wheelY)
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		if mouseButtonpressedBefore.left {
			inp.game.editor.editHandleClick(ED_CLICK_LEFT, mx, my)
			mouseButtonpressedBefore.left = false

		}

	} else {
		mouseButtonpressedBefore.left = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {

		if mouseButtonpressedBefore.middle {
			inp.game.editor.editHandleClick(ED_CLICK_MIDDLE, mx, my)
			mouseButtonpressedBefore.middle = false

		}

	} else {
		mouseButtonpressedBefore.middle = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {

		if mouseButtonpressedBefore.middle {
			inp.game.editor.editHandleClick(ED_CLICK_RIGHT, mx, my)
			mouseButtonpressedBefore.right = false

		}

	} else {
		mouseButtonpressedBefore.right = true
	}
	inp.cursorPos = pos{
		x: mx,
		y: my,
	}

	inp.keys = inpututil.AppendPressedKeys(inp.keys[:0])

	if inp.game.console.visible {
		inp.game.console.consoleHandleKeys(inp.keys)
	} else {
		for _, k := range inp.keys {
			switch k {
			case ebiten.KeyArrowUp, ebiten.KeySpace, ebiten.KeyW:
				inp.dflags.up = true
				inp.game.activateObject = true
				//inp.game.ball.Serve()
			case ebiten.KeyArrowDown, ebiten.KeyS:
				inp.dflags.down = true
				//inp.game.ball.Serve()
			case ebiten.KeyArrowLeft, ebiten.KeyA:
				inp.dflags.left = true
			case ebiten.KeyArrowRight, ebiten.KeyD:
				inp.dflags.right = true
			case ebiten.KeyF2:
				inp.game.tileMap.loadCurrentLevelMapFromFile()
				inp.game.pickupManager.loadDataFromFile()
				inp.game.fidgetManager.loadDataFromFile()
				inp.game.platformManager.loadDataFromFile()
				inp.game.entityManager.loadDataFromFile()
			case ebiten.KeyF1:
				inp.game.tileMap.saveMapToFile()
				inp.game.pickupManager.saveDataToFile()
				inp.game.fidgetManager.saveDataToFile()
				inp.game.platformManager.saveDataToFile()
				inp.game.entityManager.saveDataToFile()
			case ebiten.KeyP:
				inp.game.Pause()
			case ebiten.KeyF:
				startX := inp.game.player.worldX + PM_BULLET_PLAYER_X_OFFSET
				startY := inp.game.player.worldY + (playerHeight / 2)
				inp.game.projectileManager.AddProjectile(startX, startY, 0)
			case ebiten.KeyBackquote:
				inp.game.console.toggleWindow()
			case ebiten.KeyShiftLeft:
				inp.game.player.run = true
				inp.game.tileMap.runPan = true
			default:
				if INP_PRINT_WHEN_KEY_NOT_USED {
					fmt.Println("Key not used: ", k.String())
				}

			}

		}

		//fmt.Println(k.String())
	}

}
