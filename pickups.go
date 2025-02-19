package main

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PUM_MAX_PICKUPS   = 10
	PUM_STAR_IMAGE    = "star.png"
	PUM_SKULL_IMAGE   = "skull.png"
	PUM_SKEY_IMAGE    = "steelKey.png"
	PUM_CHICKEN_IMAGE = "chicken.png"
	PUM_FILENAME_BASE = "pickup"
	PUM_FILENAME_END  = ".csv"
	PUM_PICKUP_H      = 50
	PUM_PICKUP_W      = 50
	PUM_BOB_TICKS     = 60
	PUM_BOB_TICKS_MP  = PUM_BOB_TICKS / 2
)

type PickupManager struct {
	game          *Game
	maxPickups    int
	pickupsArray  [PUM_MAX_PICKUPS]*Pickup
	images        []*ebiten.Image
	testRect      *rect
	bobPixels     int
	bobTickCount  int
	currentBob    int
	filename_base string
	assetID       int
}

type Pickup struct {
	kind     int
	gridX    int
	gridY    int
	alive    bool
	onScreen bool
}

func NewPickupManager(game *Game) *PickupManager {

	pum := &PickupManager{}
	pum.game = game
	pum.maxPickups = PUM_MAX_PICKUPS
	pum.assetID = 0
	pum.filename_base = PUM_FILENAME_BASE
	pum.initImages()
	pum.pickupsArray = [PUM_MAX_PICKUPS]*Pickup{}
	//pum.AddPickup(200, 200, 0)
	//pum.AddPickup(200, 300, 0)
	pum.testRect = &rect{0, 0, PUM_PICKUP_W, PUM_PICKUP_H}
	pum.bobPixels = 0
	pum.bobTickCount = 0
	pum.currentBob = PUM_BOB_TICKS
	pum.loadDataFromFile()

	return pum
}

func (pum *PickupManager) Draw(screen *ebiten.Image) {

	for _, v := range pum.pickupsArray {
		if nil != v && true == v.alive {
			screenX := (pum.game.tileMap.tileSize * v.gridX) - worldOffsetX
			screenY := (pum.game.tileMap.tileSize * v.gridY) + pum.bobPixels - worldOffsetY
			DrawImageAt(screen, pum.images[v.kind], screenX, screenY)
		}
	}

}

func (pum *PickupManager) calculateBob() {
	if pum.bobTickCount < PUM_BOB_TICKS {
		pum.bobTickCount++
		if pum.bobTickCount > PUM_BOB_TICKS_MP {
			pum.bobPixels--
		} else {
			pum.bobPixels++
		}

	} else {
		pum.bobTickCount = 0
	}

}

func (pum *PickupManager) Update() {
	pum.checkPickupsTouchedPlayer()
	pum.calculateBob()
}

func (pum *PickupManager) pickupItemAction(kind int) {
	fmt.Println("Pickup item ", kind)
	pum.game.incrementScore(1)
	pum.game.audioPlayer.playSoundByID("canlid_reverb")
	switch kind {
	case 0:
		fmt.Println("Got star ", kind)
		pum.game.player.changeHealth(10)
	case 1:
		fmt.Println("Got skull ", kind)
	case 2:
		fmt.Println("Got chicken ", kind)
		pum.game.player.changeHealth(50)
	case 3:
		fmt.Println("Got key ", kind)
	}
}

func (pum *PickupManager) checkPickupsTouchedPlayer() {

	playerRect := pum.game.player.getWorldColliderRect()
	for _, v := range pum.pickupsArray {
		if nil != v && true == v.alive {
			pum.testRect.x = pum.game.tileMap.tileSize * v.gridX
			pum.testRect.y = pum.game.tileMap.tileSize * v.gridY

			if collideRect(playerRect, *pum.testRect) {
				v.alive = false
				pum.pickupItemAction(v.kind)
			}
		}
	}
}

func (pum *PickupManager) initImages() error {
	pum.images = []*ebiten.Image{}
	var tempImg *ebiten.Image

	// star
	imageDir := path.Join(subdir, PUM_STAR_IMAGE)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	starImage := copyAndStretchImage(rawImage, PUM_PICKUP_W, PUM_PICKUP_H)
	pum.images = append(pum.images, starImage)

	//skull
	imageDir = path.Join(subdir, PUM_SKULL_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	skullImage := copyAndStretchImage(rawImage, PUM_PICKUP_W, PUM_PICKUP_H)
	pum.images = append(pum.images, skullImage)

	//chicken
	imageDir = path.Join(subdir, PUM_CHICKEN_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	tempImg = copyAndStretchImage(rawImage, PUM_PICKUP_W, PUM_PICKUP_H)
	pum.images = append(pum.images, tempImg)

	//key
	imageDir = path.Join(subdir, PUM_SKEY_IMAGE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	tempImg = copyAndStretchImage(rawImage, PUM_PICKUP_W, PUM_PICKUP_H)
	pum.images = append(pum.images, tempImg)
	return err

}

func (pum *PickupManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(pum.pickupsArray); i++ {
		if nil == pum.pickupsArray[i] || false == pum.pickupsArray[i].alive {
			return i
		}
	}
	return -1
}

func (pum *PickupManager) saveDataToFile() {
	name := pum.getDataFileURL()
	numericData := [][]int{}
	rows := len(pum.pickupsArray)
	for i := 0; i < rows; i++ {
		pickupObj := pum.pickupsArray[i]
		if pickupObj != nil {
			record := []int{pickupObj.kind, pickupObj.gridX, pickupObj.gridY}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Pickups: no data to write, ", name)
	}
}

func (pum *PickupManager) getDataFileURL() string {
	filename := pum.filename_base + strconv.Itoa(pum.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (pum *PickupManager) loadDataFromFile() error {
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := pum.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		pum.pickupsArray = [PUM_MAX_PICKUPS]*Pickup{}
		return nil
	}
	if err != nil {
		pum.pickupsArray = [PUM_MAX_PICKUPS]*Pickup{}
		return err
	}
	pum.pickupsArray = [PUM_MAX_PICKUPS]*Pickup{}
	for i := 0; i < PUM_MAX_PICKUPS && i < rows; i++ {
		v := numericData[i]
		pum.pickupsArray[i] = &Pickup{v[0], v[1], v[2], true, true}
	}
	return nil
}

func (tm *PickupManager) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}

func (tm *PickupManager) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected pickup ", tm.assetID)

}

func (pum *PickupManager) AddInstanceToGrid(gridX, gridY, kind int) {
	emptySlot := pum.inactiveSlot()
	if emptySlot != -1 {
		x := gridX
		y := gridY
		pu := &Pickup{kind, x, y, true, true}
		pum.pickupsArray[emptySlot] = pu
		log.Println("Added pickup ", kind)
	} else {
		log.Println("Failed to add pickup, no open slots")
	}
}
func (tm *PickupManager) getAssetID() int {

	return tm.assetID

}

func (tm *PickupManager) setAssetID(assetID int) {

	if assetID < len(tm.images) && assetID >= 0 {
		tm.assetID = assetID
	}

}
