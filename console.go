package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Console struct {
	game              *Game
	backgroundRect    *ebiten.Image
	visible           bool
	debounceLastTime  int64
	rasterStringArray [CON_LINES]*RasterString
	keysReleased      bool
	currentCommand    string
}

const (
	CON_STARTX            = 100
	CON_STARTY            = 400
	CON_WIDTH             = 400
	CON_HEIGHT            = 100
	CON_ROW_HEIGHT        = 17
	CON_LINES             = 5
	CON_DEBOUNCE_INTERVAL = 120000000
)

var (
	CON_BG_COLOR = color.RGBA{0x7f, 0x5f, 0x5f, 0x7f}
)

func NewConsole(game *Game) *Console {
	con := &Console{}
	con.game = game
	con.backgroundRect = ebiten.NewImage(CON_WIDTH, CON_HEIGHT)
	con.backgroundRect.Fill(CON_BG_COLOR)
	con.visible = false
	now := time.Now()
	con.debounceLastTime = now.UnixNano()
	con.rasterStringArray = [CON_LINES]*RasterString{}
	con.initRasterStringLinesBlank()
	con.keysReleased = false
	return con
}

func (con *Console) Draw(screen *ebiten.Image) {
	if con.visible {

		DrawImageAt(screen, con.backgroundRect, CON_STARTX, CON_STARTY)
		for i, _ := range con.rasterStringArray {
			con.rasterStringArray[i].Draw(screen)
		}
	}

}

func (con *Console) Update() {

}

func (con *Console) enterCommand() {
	//fmt.Println("Entered ", con.currentCommand)
	con.executeCommand(con.currentCommand)
	con.cycleLineHistory()
	con.currentCommand = ""
	con.updateCommandRow()
}

func (con *Console) executeCommand(command string) {
	stringsList := strings.Split(command, " ")
	fmt.Println("Command first part ", stringsList[0])
	argsAmount := len(stringsList)
	functionSelector := stringsList[0]
	editCommandEntered := false
	fillTile := false
	switch functionSelector {
	case "LOCATION":
		gridX := con.game.player.worldX / GAME_TILE_SIZE
		gridY := con.game.player.worldY / GAME_TILE_SIZE
		fmt.Printf("Player grid pos: %d , %d /n", gridX, gridY)

	case "FLY":
		con.game.player.hoverMode = !con.game.player.hoverMode
	case "GOD":
		con.game.godMode = !con.game.godMode
	case "TILE":
		fmt.Println("Set tile")
		con.game.editMode = EditTile
		editCommandEntered = true
	case "ENTITY":
		fmt.Println("Set entity")
		con.game.editMode = EditEntity
		editCommandEntered = true
	case "NONE":
		fmt.Println("Set none")
		con.game.editMode = EditNone
		editCommandEntered = true
	case "PICKUP":
		fmt.Println("Set pickup")
		con.game.editMode = EditPickup
		editCommandEntered = true
	case "DECOR":
		fmt.Println("Set decor")
		con.game.editMode = EditDecor
		editCommandEntered = true
	case "FIDGET":
		fmt.Println("Set decor")
		con.game.editMode = EditFidget
		editCommandEntered = true
	case "PLATFORM":
		fmt.Println("Set platform")
		con.game.editMode = EditPlatform
		editCommandEntered = true
	case "FILL":
		fmt.Println("Fill with tile")
		con.game.editMode = EditTile
		editCommandEntered = true
		fillTile = true
	case "LEVEL":
		if argsAmount == 2 {
			levelNo, err := strconv.Atoi(strings.Trim(stringsList[1], " "))
			fmt.Println("Load level ")
			if err == nil {
				con.game.loadLevel(levelNo)
			} else {
				log.Println("Invalid arg: ", stringsList[1])
				log.Println(err)
			}
		}

	}

	if argsAmount == 2 && editCommandEntered && !fillTile {
		assetID, _ := strconv.Atoi(stringsList[1])
		fmt.Println("Set assetID ", assetID)
		con.game.editor.setActiveComponentAssetID(assetID)
	} else if editCommandEntered && fillTile {
		var assetID int
		if argsAmount == 2 {
			assetID, _ = strconv.Atoi(stringsList[1])
		} else {
			assetID = 0
		}

		fmt.Println("Set assetID fill ", assetID)
		con.game.tileMap.fillWithTile(assetID)
	}
}

func (con *Console) deleteLine() {
	fmt.Println("Delete line ", con.currentCommand)
	con.currentCommand = ""
	con.updateCommandRow()
}

func (con *Console) deleteLast() {
	oldCommand := con.currentCommand
	oldLen := len(oldCommand)
	if oldLen == 0 {
		return
	}
	con.currentCommand = oldCommand[0 : oldLen-1]
	con.updateCommandRow()
}

func keyNameMunger(keyName string) rune {
	nameLength := len(keyName)
	fmt.Println("keyNameMunger ", keyName)
	switch nameLength {
	case 1:
		return rune(keyName[0])
	case 5:
		if keyName[0:5] == "Space" {
			return ' '
		} else if keyName[0:5] == "Comma" {
			return ','
		}

	case 6:
		// Digit
		if keyName[0:5] == "Digit" {
			return rune(keyName[5])
		}
	case 7:
		if keyName[0:7] == "Period" {
			return '.'
		}
	default:
		return rune(keyName[0])
	}
	return rune(keyName[0])

}

func (con *Console) cycleLineHistory() {
	lastLine := con.rasterStringArray[CON_LINES-1].stringContent
	for i := CON_LINES - 1; i >= 0; i-- {
		temp := con.rasterStringArray[i].stringContent
		con.rasterStringArray[i].stringContent = lastLine
		lastLine = temp
	}
}

func (con *Console) consoleHandleKeys(keys []ebiten.Key) {
	keysReleasedLastTime := con.keysReleased
	if len(keys) == 0 {
		con.keysReleased = true
	} else {
		con.keysReleased = false
	}
	var sb strings.Builder
	sb.WriteString(con.currentCommand)
	if keysReleasedLastTime {
		for _, k := range keys {
			switch k {

			case ebiten.KeyF2:
				con.game.tileMap.loadCurrentLevelMapFromFile()
			case ebiten.KeyBackspace:
				con.deleteLast()
			case ebiten.KeyDelete:
				con.deleteLine()
			case ebiten.KeyEnter:
				con.enterCommand()
			case ebiten.KeyBackquote:
				con.game.console.toggleWindow()
			default:
				keyBytes, _ := k.MarshalText()
				keyString := string(keyBytes)
				if len(keyString) != 1 {
					keyString = string(keyNameMunger(keyString))
				}
				//fmt.Println("console key: ", keyString)
				//keyRune :=KeyToRune(k.String())
				sb.WriteString(keyString)
				con.currentCommand = sb.String()
				con.updateCommandRow()

			}
		}

	}

}

func (con *Console) updateCommandRow() {
	con.rasterStringArray[CON_LINES-1].stringContent = con.currentCommand
}

func (con *Console) initRasterStringLinesBlank() {
	currentYLocation := CON_STARTY
	for i, _ := range con.rasterStringArray {
		con.rasterStringArray[i] = NewRasterString(
			con.game, "", CON_STARTX, currentYLocation)
		currentYLocation += CON_ROW_HEIGHT
	}
}

func (con *Console) toggleWindow() {
	now := time.Now()
	timerNowNano := now.UnixNano()

	interval := timerNowNano - con.debounceLastTime
	//fmt.Println("Console last interval, ", interval)
	if interval > CON_DEBOUNCE_INTERVAL {
		con.visible = !con.visible
		con.debounceLastTime = timerNowNano
	}

}
