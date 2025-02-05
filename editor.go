package main

import (
	"log"
)

type ED_CLICK_ENUM int

const (
	ED_CLICK_LEFT ED_CLICK_ENUM = iota
	ED_CLICK_MIDDLE
	ED_CLICK_RIGHT
)

type Editor struct {
	game    *Game
	assetID int
}

func NewEditor(game *Game) *Editor {
	ed := &Editor{}
	ed.game = game
	ed.assetID = 0
	return ed
}

func (ed *Editor) editHandleWheel(wheelY float64) {
	switch ed.game.editMode {
	case EditTile:
		ed.game.tileMap.CycleAssetKind(int(wheelY))
	case EditPickup:
		ed.game.tileMap.CycleAssetKind(int(wheelY))
	}

}

func (ed *Editor) editHandleClick(clickType ED_CLICK_ENUM, mouseX, mouseY int) {

	switch clickType {
	case ED_CLICK_LEFT:
		ed.editHandleLeftClick(mouseX, mouseY)
	case ED_CLICK_RIGHT:
		ed.editHandleLeftClick(mouseX, mouseY)
	case ED_CLICK_MIDDLE:
		ed.editHandleLeftClick(mouseX, mouseY)
	default:
		log.Fatal("editHandleClick invalid argument for click type")

	}

}

func (ed *Editor) editHandleLeftClick(mouseX, mouseY int) {

	var component editable
	switch ed.game.editMode {
	case EditNone:
	case EditDecor:
	case EditEntity:
		component = ed.game.entityManager
	case EditFidget:
		component = ed.game.fidgetManager
	case EditTile:
		component = ed.game.tileMap
	case EditSpawner:
	case EditPickup:
		component = ed.game.pickupManager
	case EditZone:
	case EditPlatform:
		component = ed.game.platformManager

	}
	ts := ed.game.tileMap.tileSize
	gridX := (mouseX + worldOffsetX) / ts
	gridY := (mouseY + worldOffsetY) / ts
	if component == nil {
		log.Println("Editor editHandleLeftClick no matching component")
		return
	}
	component.AddInstanceToGrid(gridX, gridY, ed.assetID)

}

func (ed *Editor) editHandleRightClick(mouseX, mouseY int) {

}

func (ed *Editor) editHandleMiddleClick(mouseX, mouseY int) {
	switch ed.game.editMode {
	case EditNone:
	case EditDecor:
	case EditInteractive:
	case EditTile:
		ed.game.tileMap.CycleAssetKind(1)
	case EditSpawner:
	case EditPickup:
	case EditZone:
	case EditFidget:
		ed.game.fidgetManager.CycleAssetKind(1)

	}

}
