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
	FM_MAX_ENTITY_ROOM          = 10
	IMAGES_IDLE_SHEET           = "jackieD1.png"
	IMAGES_WALK_SHEET           = "jackieD1.png"
	IMAGES_ATTACK_SHEET         = "jackieD1.png"
	EN_FILENAME_BASE            = "entity"
	EN_FILENAME_END             = ".csv"
	EN_SPRITE_H                 = 100
	EN_SPRITE_W                 = 100
	EN_CREATE_FILE_IF_NOT_EXIST = true
)

type EntityManager struct {
	game          *Game
	maxEntitys    int
	EntitysArray  [FM_MAX_ENTITY_ROOM]*Entity
	images        []*ebiten.Image
	testRect      *rect
	assetID       int
	filename_base string
}

type Entity struct {
	kind       int
	startGridX int
	startGridY int
	uid        int
	worldX     int
	worldY     int
	health     int
	alive      bool
	onScreen   bool
}

func NewEntity(kind, startGridX, startGridY int) *Entity {
	worldX := GAME_TILE_SIZE * startGridX
	worldY := GAME_TILE_SIZE * startGridY
	ent := &Entity{kind, startGridX, startGridY,
		0, worldX, worldY, 100, true, true}
	return ent

}
func (tm *EntityManager) getAssetID() int {

	return tm.assetID

}

func NewEntityManager(game *Game) *EntityManager {

	fm := &EntityManager{}
	fm.game = game
	fm.filename_base = EN_FILENAME_BASE
	fm.maxEntitys = FM_MAX_ENTITY_ROOM
	fm.initImages()
	fm.EntitysArray = [FM_MAX_ENTITY_ROOM]*Entity{}
	//pum.AddPickup(200, 200, 0)
	//pum.AddPickup(200, 300, 0)
	fm.testRect = &rect{0, 0, EN_SPRITE_W, EN_SPRITE_H}
	fm.assetID = 0
	fm.loadDataFromFile()

	return fm
}

func (pum *EntityManager) Draw(screen *ebiten.Image) {

	for _, v := range pum.EntitysArray {
		if nil != v && true == v.alive {
			screenX := (pum.game.tileMap.tileSize * v.startGridX) - worldOffsetX
			screenY := (pum.game.tileMap.tileSize * v.startGridY) - worldOffsetY
			DrawImageAt(screen, pum.images[v.kind], screenX, screenY)
		}
	}

}

func (pum *EntityManager) Update() {
	pum.checkEntitysTouchedPlayer()

	pum.game.activateObject = false

}

func (pum *EntityManager) touchEntityAction(kind, uid int) {
	//fmt.Println("Entity touched ", kind)
	//pum.game.incrementScore(1)
	if kind == 0 && pum.game.activateObject == true {
		pum.game.warpManager.warpPlayerToWarpID(uid)
		pum.game.activateObject = false
	}
}

func (pum *EntityManager) checkEntitysTouchedPlayer() {

	playerRect := pum.game.player.getWorldColliderRect()
	for _, v := range pum.EntitysArray {
		if nil != v && true == v.alive {
			pum.testRect.x = pum.game.tileMap.tileSize * v.startGridX
			pum.testRect.y = pum.game.tileMap.tileSize * v.startGridY

			if collideRect(playerRect, *pum.testRect) {
				//v.alive = false
				pum.touchEntityAction(v.kind, v.uid)
			}
		}
	}
}

func (pum *EntityManager) initImages() error {
	pum.images = []*ebiten.Image{}
	imageDir := path.Join(subdir, IMAGES_IDLE_SHEET)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	starImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	pum.images = append(pum.images, starImage)
	//skull
	imageDir = path.Join(subdir, IMAGES_WALK_SHEET)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	skullImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	pum.images = append(pum.images, skullImage)
	//spikes
	imageDir = path.Join(subdir, IMAGES_ATTACK_SHEET)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	spikeImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	pum.images = append(pum.images, spikeImage)
	return err

}

func (pum *EntityManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(pum.EntitysArray); i++ {
		if nil == pum.EntitysArray[i] || false == pum.EntitysArray[i].alive {
			return i
		}
	}
	return -1
}

func (pum *EntityManager) saveDataToFile() {
	name := pum.getDataFileURL()
	numericData := [][]int{}
	rows := len(pum.EntitysArray)
	for i := 0; i < rows; i++ {
		pickupObj := pum.EntitysArray[i]
		if pickupObj != nil {
			record := []int{pickupObj.kind, pickupObj.startGridX, pickupObj.startGridY, pickupObj.uid}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Entitys: no data to write, ", name)
	}
}
func (pum *EntityManager) getDataFileURL() string {
	filename := pum.filename_base + strconv.Itoa(pum.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}
func (pum *EntityManager) loadDataFromFile() error {
	pum.EntitysArray = [FM_MAX_ENTITY_ROOM]*Entity{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := pum.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		log.Println("Entity loadDataFromFile no data to load")
		return nil
	}
	if err != nil {
		return err
	}

	for i := 0; i < FM_MAX_ENTITY_ROOM && i < rows; i++ {
		v := numericData[i]
		//pum.EntitysArray[i] = &Entity{v[0], v[1], v[2], v[3], true, true}
		pum.EntitysArray[i] = NewEntity(v[0], v[1], v[2])
		pum.EntitysArray[i].uid = v[3]
	}
	return nil
}
func (pum *EntityManager) getUniqueUID() int {

	return 0
}

func (pum *EntityManager) AddInstanceToGrid(gridX, gridY, kind int) {
	emptySlot := pum.inactiveSlot()
	if emptySlot != -1 {
		x := gridX
		y := gridY
		//uid := pum.getUniqueUID()
		pu := NewEntity(kind, x, y)
		pu.uid = pum.getUniqueUID()
		pum.EntitysArray[emptySlot] = pu
		log.Println("Added Entity ", kind)
	} else {
		log.Println("Failed to add Entity, no open slots")
	}
}

func (tm *EntityManager) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}
func (tm *EntityManager) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected Entity ", tm.assetID)

}
