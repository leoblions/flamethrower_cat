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
 0-8 plant 100px

*/

const (
	DM_MAX_DECOR                = 30
	DM_PLANT_1                  = "decorPlant100A.png"
	DM_PLANT_2                  = "decorPlant100B.png"
	DM_FACTORY_1                = "decorFactory100.png"
	DM_SEASIDE                  = "decorSeaside.png"
	DM_ALIEN                    = "decorAlienA.png"
	DM_FILENAME_BASE            = "Decor"
	DM_FILENAME_END             = ".csv"
	DM_SPRITE_H                 = 100
	DM_SPRITE_W                 = 100
	DM_CREATE_FILE_IF_NOT_EXIST = true
)

type DecorManager struct {
	game          *Game
	maxDecor      int
	DecorArray    [DM_MAX_DECOR]*Decor
	images        []*ebiten.Image
	testRect      *rect
	assetID       int
	filename_base string
}

type Decor struct {
	kind     int
	gridX    int
	gridY    int
	uid      int
	alive    bool
	onScreen bool
}

func NewDecorManager(game *Game) *DecorManager {

	fm := &DecorManager{}
	fm.filename_base = DM_FILENAME_BASE
	fm.game = game
	fm.maxDecor = DM_MAX_DECOR
	fm.initImages()
	fm.DecorArray = [DM_MAX_DECOR]*Decor{}
	//dm.AddPickup(200, 200, 0)
	//dm.AddPickup(200, 300, 0)
	fm.testRect = &rect{0, 0, DM_SPRITE_H, DM_SPRITE_H}
	fm.assetID = 0
	fm.loadDataFromFile()

	return fm
}
func (tm *DecorManager) getAssetID() int {

	return tm.assetID

}

func (dm *DecorManager) Draw(screen *ebiten.Image) {

	for _, v := range dm.DecorArray {
		if nil != v && true == v.alive {
			screenX := (dm.game.tileMap.tileSize * v.gridX) - worldOffsetX
			screenY := (dm.game.tileMap.tileSize * v.gridY) - worldOffsetY
			DrawImageAt(screen, dm.images[v.kind], screenX, screenY)
		}
	}

}

func (dm *DecorManager) Update() {
	dm.checkDecorsTouchedPlayer()

	dm.game.activateObject = false

}

func (dm *DecorManager) touchDecorAction(kind, uid int) {
	//fmt.Println("Decor touched ", kind)
	//dm.game.incrementScore(1)
	if kind == 0 && dm.game.activateObject == true {
		dm.game.warpManager.warpPlayerToWarpID(uid)
		dm.game.activateObject = false
	}
}

func (dm *DecorManager) checkDecorsTouchedPlayer() {

	playerRect := dm.game.player.getWorldColliderRect()
	for _, v := range dm.DecorArray {
		if nil != v && true == v.alive {
			dm.testRect.x = dm.game.tileMap.tileSize * v.gridX
			dm.testRect.y = dm.game.tileMap.tileSize * v.gridY

			if collideRect(playerRect, *dm.testRect) {
				//v.alive = false
				dm.touchDecorAction(v.kind, v.uid)
			}
		}
	}
}

func (dm *DecorManager) initImages() error {
	dm.images = []*ebiten.Image{}
	// plant 1
	imageDir := path.Join(subdir, DM_PLANT_1)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	cutImages := cutSpriteSheet(rawImage, DM_SPRITE_W, DM_SPRITE_H, 2, 2)
	dm.images = append(dm.images, cutImages...)
	// plant 2
	imageDir = path.Join(subdir, DM_PLANT_2)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	cutImages = cutSpriteSheet(rawImage, DM_SPRITE_W, DM_SPRITE_H, 2, 2)
	dm.images = append(dm.images, cutImages...)
	// factory 1
	imageDir = path.Join(subdir, DM_FACTORY_1)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	cutImages = cutSpriteSheet(rawImage, DM_SPRITE_W, DM_SPRITE_H, 2, 2)
	dm.images = append(dm.images, cutImages...)
	// seaside
	imageDir = path.Join(subdir, DM_SEASIDE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	cutImages = cutSpriteSheet(rawImage, DM_SPRITE_W, DM_SPRITE_H, 2, 2)
	dm.images = append(dm.images, cutImages...)
	// alien
	imageDir = path.Join(subdir, DM_ALIEN)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	cutImages = cutSpriteSheet(rawImage, DM_SPRITE_W, DM_SPRITE_H, 2, 2)
	dm.images = append(dm.images, cutImages...)
	return err

}

func (dm *DecorManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(dm.DecorArray); i++ {
		if nil == dm.DecorArray[i] || false == dm.DecorArray[i].alive {
			return i
		}
	}
	return -1
}

func (dm *DecorManager) saveDataToFile() {
	name := dm.getDataFileURL()
	numericData := [][]int{}
	rows := len(dm.DecorArray)
	for i := 0; i < rows; i++ {
		pickupObj := dm.DecorArray[i]
		if pickupObj != nil {
			record := []int{pickupObj.kind, pickupObj.gridX, pickupObj.gridY, pickupObj.uid}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Decors: no data to write, ", name)
	}
}

func (dm *DecorManager) getDataFileURL() string {
	filename := dm.filename_base + strconv.Itoa(dm.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (dm *DecorManager) loadDataFromFile() error {
	dm.DecorArray = [DM_MAX_DECOR]*Decor{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := dm.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		dm.DecorArray = [DM_MAX_DECOR]*Decor{}
		log.Println("Decor loadDataFromFile no data to load")
		return nil
	}
	if err != nil {
		dm.DecorArray = [DM_MAX_DECOR]*Decor{}
		log.Println("Decor loadDataFromFile error")
		return err
	}

	for i := 0; i < DM_MAX_DECOR && i < rows; i++ {
		v := numericData[i]
		dm.DecorArray[i] = &Decor{v[0], v[1], v[2], v[3], true, true}
	}
	return nil
}
func (dm *DecorManager) getUniqueUID() int {

	return 0
}

func (dm *DecorManager) AddInstanceToGrid(gridX, gridY, kind int) {
	emptySlot := dm.inactiveSlot()
	if emptySlot != -1 {
		x := gridX
		y := gridY
		uid := dm.getUniqueUID()
		pu := &Decor{kind, x, y, uid, true, true}
		dm.DecorArray[emptySlot] = pu
		log.Println("Added Decor ", kind)
	} else {
		log.Println("Failed to add Decor, no open slots")
	}
}

func (tm *DecorManager) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}

func (tm *DecorManager) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected Decor ", tm.assetID)

}

func (tm *DecorManager) setAssetID(assetID int) {

	if assetID < len(tm.images) && assetID >= 0 {
		tm.assetID = assetID
	}

}
