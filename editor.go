package main

import (
	"fmt"
	"log"
)

type ED_CLICK_ENUM int

const (
	ED_CLICK_LEFT ED_CLICK_ENUM = iota
	ED_CLICK_MIDDLE
	ED_CLICK_RIGHT
)

type Editor struct {
	game *Game
	//assetID int
}

func NewEditor(game *Game) *Editor {
	ed := &Editor{}
	ed.game = game
	//ed.assetID = 0
	return ed
}

func (ed *Editor) editHandleWheel(wheelY float64) {
	switch ed.game.editMode {
	case EditTile:
		ed.game.tileMap.CycleAssetKind(int(wheelY))
	case EditPickup:
		ed.game.pickupManager.CycleAssetKind(int(wheelY))
	case EditDecor:
		ed.game.decorManager.CycleAssetKind(int(wheelY))
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

	var component = ed.getActiveEditableComponent()

	ts := ed.game.tileMap.tileSize
	gridX := (mouseX + worldOffsetX) / ts
	gridY := (mouseY + worldOffsetY) / ts
	if component == nil {
		log.Println("Editor editHandleLeftClick no matching component")
		return
	}
	component.AddInstanceToGrid(gridX, gridY, component.getAssetID())

}

func (ed *Editor) editHandleRightClick(mouseX, mouseY int) {

}

func (ed *Editor) setActiveComponentAssetID(assetID int) {
	component := ed.getActiveEditableComponent()
	if component != nil {
		component.setAssetID(assetID)
		fmt.Println("Editor:   component, setAssetID", assetID)
	} else {
		fmt.Println("Editor: Cannot edit nil component, setAssetID")
	}

}

func (ed *Editor) getActiveEditableComponent() editable {
	switch ed.game.editMode {
	case EditNone:
		return nil
	case EditDecor:
		return ed.game.decorManager
	case EditInteractive:
		return ed.game.fidgetManager
	case EditTile:
		return ed.game.tileMap
	case EditSpawner:
		return nil
	case EditPickup:
		return ed.game.pickupManager
	case EditZone:
		return nil
	case EditFidget:
		return ed.game.fidgetManager
	case EditPlatform:
		return ed.game.platformManager
	case EditEntity:
		return ed.game.entityManager

	}
	return nil

}

func (ed *Editor) editHandleMiddleClick(mouseX, mouseY int) {
	component := ed.getActiveEditableComponent()
	if component != nil {
		component.CycleAssetKind(1)
	} else {
		fmt.Println("Editor: Cannot edit nil component")
	}

}
