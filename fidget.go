package main

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/*
DATA FILES:
 kind,  gridX,  gridY,  uid

 KINDS:
 0 = DOOR
 1 = BARREL
 2 = SPIKES
 3 = TUNNEL
 4 = TRAFFIC LIGHT
 5 = LADDER

*/

const (
	FM_MAX_FIDGETS              = 10
	DOOR_IMAGE                  = "door1.png"
	TUNNEL_IMAGE                = "door2.png"
	BARREL_IMAGE                = "barrel1.png"
	SPIKES_IMAGE                = "spikes.png"
	TRAFFICLIGHT_IMAGE          = "trafficLight.png"
	LADDER_IMAGE                = "ladder50x100.png"
	FM_FILENAME_BASE            = "fidget"
	FM_FILENAME_END             = ".csv"
	FM_SPRITE_H                 = 100
	FM_SPRITE_W                 = 50
	FM_SPRITE_W_SL              = 100
	FM_CREATE_FILE_IF_NOT_EXIST = true
	FM_TL_CHANGE_TICKS          = 200
)

type FidgetManager struct {
	game               *Game
	maxFidgets         int
	fidgetsArray       [FM_MAX_FIDGETS]*Fidget
	images             []*ebiten.Image
	trafficLightImages []*ebiten.Image
	testRect           *rect
	assetID            int
	filename_base      string
	tlIndex            int
	tlEnabled          bool
	tlCycler           func() int

	touchedtrafficLightLastTick bool
}

type Fidget struct {
	kind     int
	gridX    int
	gridY    int
	uid      int
	alive    bool
	onScreen bool
}

func NewFidgetManager(game *Game) *FidgetManager {

	fm := &FidgetManager{}
	fm.filename_base = FM_FILENAME_BASE
	fm.game = game
	fm.maxFidgets = FM_MAX_FIDGETS
	fm.initImages()
	fm.fidgetsArray = [FM_MAX_FIDGETS]*Fidget{}
	//pum.AddPickup(200, 200, 0)
	//pum.AddPickup(200, 300, 0)
	fm.testRect = &rect{0, 0, FM_SPRITE_W, FM_SPRITE_H}
	fm.assetID = 0

	fm.tlCycler = fm.getTLCycler()
	fm.tlIndex = fm.tlCycler()
	fm.tlEnabled = true

	fm.loadDataFromFile()

	return fm
}
func (tm *FidgetManager) getTLCycler() func() int {
	// closure function for getting state of traffic lights
	currentTick := 0
	currentState := 1
	cyclerClosure := func() int {
		if currentTick >= FM_TL_CHANGE_TICKS {
			currentTick = 0
			if currentState < 3 {
				currentState++
			} else {
				currentState = 1
			}
		} else {
			currentTick++

		}
		if tm.tlEnabled {
			return currentState
		} else {
			return 0 // lights disabled
		}

	}
	return cyclerClosure

}

func (tm *FidgetManager) getAssetID() int {

	return tm.assetID

}

func (pum *FidgetManager) Draw(screen *ebiten.Image) {

	for _, v := range pum.fidgetsArray {
		if nil != v && true == v.alive {
			screenX := (pum.game.tileMap.tileSize * v.gridX) - worldOffsetX
			screenY := (pum.game.tileMap.tileSize * v.gridY) - worldOffsetY

			DrawImageAt(screen, pum.images[v.kind], screenX, screenY)
		}
	}

}

func (pum *FidgetManager) Update() {
	pum.checkFidgetsTouchedPlayer()

	pum.game.activateObject = false
	pum.tlIndex = pum.tlCycler()
	pum.images[4] = pum.trafficLightImages[pum.tlIndex]

}

func (pum *FidgetManager) checkPlayerRanRedlight() {

}

func (pum *FidgetManager) touchFidgetAction(kind, uid int) {

	switch kind {
	case 0, 3:
		if pum.game.activateObject {
			pum.game.warpManager.warpPlayerToWarpID(uid)
			pum.game.activateObject = false
			pum.game.audioPlayer.playSoundByID("dooropen")
		}
	case 1:
		fmt.Println("player touched the barrel")
	case 2:
		pum.game.player.changeHealth(-1)
	case 5:
		pum.game.player.touchingLadder = true
		//fmt.Println("player touched the barrel")

	}

}

func (pum *FidgetManager) checkFidgetsTouchedPlayer() {

	//playerRect := pum.game.player.getWorldColliderRect()
	playerRect := pum.game.player.collRect
	for _, v := range pum.fidgetsArray {
		if nil != v && true == v.alive {
			pum.testRect.x = pum.game.tileMap.tileSize * v.gridX
			pum.testRect.y = pum.game.tileMap.tileSize * v.gridY

			if collideRect(*playerRect, *pum.testRect) {
				//v.alive = false
				pum.touchFidgetAction(v.kind, v.uid)
			}
		}
	}
}

func (pum *FidgetManager) initImages() error {
	pum.images = []*ebiten.Image{}
	// door
	imageDir := path.Join(subdir, DOOR_IMAGE)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	image := copyAndStretchImage(rawImage, FM_SPRITE_W, FM_SPRITE_H)
	pum.images = append(pum.images, image)
	// barrel
	imageDir = path.Join(subdir, BARREL_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	image = copyAndStretchImage(rawImage, FM_SPRITE_W, FM_SPRITE_H)
	pum.images = append(pum.images, image)
	// spikes
	imageDir = path.Join(subdir, SPIKES_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	image = copyAndStretchImage(rawImage, FM_SPRITE_W, FM_SPRITE_H)
	pum.images = append(pum.images, image)
	// tunnel
	imageDir = path.Join(subdir, TUNNEL_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	image = copyAndStretchImage(rawImage, FM_SPRITE_W, FM_SPRITE_H)
	pum.images = append(pum.images, image)
	// stoplight
	imageDir = path.Join(subdir, TRAFFICLIGHT_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	image = getSubImage(rawImage, 0, 0, FM_SPRITE_W_SL, FM_SPRITE_H)
	pum.trafficLightImages = grabImagesRowToList(rawImage, 100, 0, 4)
	pum.images = append(pum.images, image)
	imageDir = path.Join(subdir, LADDER_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	image = getSubImage(rawImage, 0, 0, FM_SPRITE_W_SL, FM_SPRITE_H)
	pum.trafficLightImages = grabImagesRowToList(rawImage, 100, 0, 4)
	pum.images = append(pum.images, image)
	return err

}

func (pum *FidgetManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(pum.fidgetsArray); i++ {
		if nil == pum.fidgetsArray[i] || false == pum.fidgetsArray[i].alive {
			return i
		}
	}
	return -1
}

func (pum *FidgetManager) saveDataToFile() {
	name := pum.getDataFileURL()
	numericData := [][]int{}
	rows := len(pum.fidgetsArray)
	for i := 0; i < rows; i++ {
		pickupObj := pum.fidgetsArray[i]
		if pickupObj != nil {
			record := []int{pickupObj.kind, pickupObj.gridX, pickupObj.gridY, pickupObj.uid}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Fidgets: no data to write, ", name)
	}
}

func (pum *FidgetManager) getDataFileURL() string {
	filename := pum.filename_base + strconv.Itoa(pum.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (pum *FidgetManager) loadDataFromFile() error {
	pum.fidgetsArray = [FM_MAX_FIDGETS]*Fidget{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := pum.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		pum.fidgetsArray = [FM_MAX_FIDGETS]*Fidget{}
		log.Println("Fidget loadDataFromFile no data to load")
		return nil
	}
	if err != nil {
		pum.fidgetsArray = [FM_MAX_FIDGETS]*Fidget{}
		log.Println("Fidget loadDataFromFile error")
		return err
	}

	for i := 0; i < FM_MAX_FIDGETS && i < rows; i++ {
		v := numericData[i]
		pum.fidgetsArray[i] = &Fidget{v[0], v[1], v[2], v[3], true, true}
	}
	return nil
}
func (pum *FidgetManager) getUniqueUID() int {

	return 0
}

func (pum *FidgetManager) AddInstanceToGrid(gridX, gridY, kind int) {
	emptySlot := pum.inactiveSlot()
	if emptySlot != -1 {
		x := gridX
		y := gridY
		uid := pum.getUniqueUID()
		pu := &Fidget{kind, x, y, uid, true, true}
		pum.fidgetsArray[emptySlot] = pu
		log.Println("Added fidget ", kind)
	} else {
		log.Println("Failed to add fidget, no open slots")
	}
}

func (tm *FidgetManager) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}

func (tm *FidgetManager) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected fidget ", tm.assetID)

}

func (tm *FidgetManager) setAssetID(assetID int) {

	if assetID < len(tm.images) && assetID >= 0 {
		tm.assetID = assetID
	}

}
