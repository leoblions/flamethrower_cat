package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type ED_CLICK_ENUM int

const (
	ED_CLICK_LEFT ED_CLICK_ENUM = iota
	ED_CLICK_MIDDLE
	ED_CLICK_RIGHT
)

const (
	ED_TEXT_LEFT = 20
	ED_TEXT_TOP1 = 55
	ED_TEXT_TOP2 = 70
	ED_TEXT_TOP3 = 85
)

type Editor struct {
	game                  *Game
	drawEditInfo          bool
	holdPaintMouseEnabled bool
	holdPaintMouseActive  bool
	displayAssetID        int
	editString1           *RasterString
	editString2           *RasterString
	editString3           *RasterString
	//assetID int
}

func NewEditor(game *Game) *Editor {
	ed := &Editor{}
	ed.game = game
	ed.drawEditInfo = true
	//ed.assetID = 0
	ed.editString1 = NewRasterString(game, "AssetID ", ED_TEXT_LEFT, ED_TEXT_TOP1)
	ed.editString2 = NewRasterString(game, "Mode ", ED_TEXT_LEFT, ED_TEXT_TOP2)
	ed.editString3 = NewRasterString(game, "PAINT ", ED_TEXT_LEFT, ED_TEXT_TOP3)
	return ed
}

func (ed *Editor) setAssetIDText(assetID int) {
	ed.editString1.stringContent = fmt.Sprintf("AssetID %d", assetID)

}

func (ed *Editor) setModeText(mode string) {
	ed.editString2.stringContent = fmt.Sprintf("Mode %s", mode)

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
		ed.editHandleRightClick(mouseX, mouseY)
	case ED_CLICK_MIDDLE:
		ed.editHandleMiddleClick(mouseX, mouseY)
	default:
		log.Fatal("editHandleClick invalid argument for click type")

	}

}

func (ed *Editor) editHandleLeftClick(mouseX, mouseY int) {
	ed.clickPaintComponent(mouseX, mouseY)

	if ed.holdPaintMouseEnabled {
		ed.holdPaintMouseActive = true
	} else {
		ed.holdPaintMouseActive = false
	}

}

func (ed *Editor) clickPaintComponent(mouseX, mouseY int) {
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

	ed.holdPaintMouseEnabled = !ed.holdPaintMouseEnabled

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

func (ed *Editor) Draw(screen *ebiten.Image) {
	if ed.drawEditInfo {
		ed.editString1.Draw(screen)
		ed.editString2.Draw(screen)
		if ed.holdPaintMouseEnabled {
			ed.editString3.Draw(screen)
		}
	}

}

func (ed *Editor) Update() {
	if ed.holdPaintMouseEnabled && ed.holdPaintMouseActive {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			ed.clickPaintComponent(ebiten.CursorPosition())

		} else {
			// button was released
			ed.holdPaintMouseActive = false
		}
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
