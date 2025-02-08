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
0 = Jackie (non enemy)
1 = robodog
2 = worm blob
3 = ice golem
4 = evil jackie
*/
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
	EN_FOLLOW_DIST              = 300
	EN_ENEMY_SPEED_1            = 2
	EN_STOP_FOLLOW_DIST         = 100
)

type EntityManager struct {
	game          *Game
	maxEntitys    int
	entityList    []*Entity
	images        []*ebiten.Image
	testRect      *rect
	assetID       int
	filename_base string

	EntityManagerImageCollections
}

type Entity struct {
	kind       int
	startGridX int
	startGridY int
	uid        int
	worldX     int
	worldY     int
	velX       int
	velY       int
	health     int
	width      int
	height     int
	direction  rune
	alive      bool
	onScreen   bool
	isEnemy    bool
}

func NewEntity(kind, startGridX, startGridY int) *Entity {
	ent := &Entity{}
	worldX := GAME_TILE_SIZE * startGridX
	worldY := GAME_TILE_SIZE * startGridY
	ent.height = EN_SPRITE_H
	ent.width = EN_SPRITE_W
	ent.kind = kind
	ent.alive = true
	ent.startGridX = startGridX
	ent.startGridY = startGridY
	ent.worldX = worldX
	ent.worldY = worldY
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
	fm.entityList = []*Entity{}
	//ent.AddPickup(200, 200, 0)
	//ent.AddPickup(200, 300, 0)
	fm.testRect = &rect{0, 0, EN_SPRITE_W, EN_SPRITE_H}
	fm.assetID = 0
	fm.loadDataFromFile()
	fm.testRect = &rect{}

	return fm
}

func (ent *EntityManager) Draw(screen *ebiten.Image) {

	for _, v := range ent.entityList {
		if nil != v && true == v.alive {
			//fmt.Println("ent draw")
			screenX := (v.worldX) - worldOffsetX
			screenY := (v.worldY) - worldOffsetY
			DrawImageAt(screen, ent.images[v.kind], screenX, screenY)
		}
	}

}

func (em *EntityManager) Update() {
	em.entityMotion()
	em.checkEntitysTouchedPlayer()

	em.game.activateObject = false

}

func (entity *Entity) entityFollowPlayer(game *Game) {
	//fmt.Println("entity folow player")
	pposX := game.player.worldX
	//pposY := game.player.worldY

	//fmt.Println("EM ", entity.worldX+EN_STOP_FOLLOW_DIST)

	if entity.worldX+EN_STOP_FOLLOW_DIST < pposX {
		entity.velX = EN_ENEMY_SPEED_1
	} else if entity.worldX-EN_STOP_FOLLOW_DIST > pposX {
		entity.velX = -EN_ENEMY_SPEED_1

	} else {
		entity.velX = 0
	}

}

func (entity *Entity) entityDetectPlatformEdge(game *Game) bool {
	checkPointOffset := 5
	checkPointX := 0
	checkPointY := entity.height + entity.worldY + checkPointOffset
	if entity.velX > 0 {
		checkPointX = entity.worldX + entity.width
	} else if entity.velX < 0 {
		checkPointX = entity.worldX
	} else {
		checkPointX = entity.worldX + (entity.width / 2)
	}
	return !game.tileMap.pointCollidedWithSolidTile(checkPointX, checkPointY)

}

func (entity *Entity) entityDetectAdjacentWall(game *Game) bool {
	checkPointOffset := 5
	checkPointX := 0
	checkPointY := (entity.height / 2) + entity.worldY
	if entity.velX > 0 {
		checkPointX = entity.worldX + entity.width + checkPointOffset
	} else if entity.velX < 0 {
		checkPointX = entity.worldX - checkPointOffset
	} else {
		checkPointX = entity.worldX + (entity.width / 2)
	}
	return game.tileMap.pointCollidedWithSolidTile(checkPointX, checkPointY)

}

func (em *EntityManager) entityMotion() {

	for _, entity := range em.entityList {

		entity.entityFollowPlayer(em.game)

		//fmt.Println("EM ", entity.entityDetectPlatformEdge(em.game))
		if entity.entityDetectPlatformEdge(em.game) || entity.entityDetectAdjacentWall(em.game) {
			entity.velX = 0
		}
		//calculate coll rect
		em.testRect.x = entity.worldX
		em.testRect.y = entity.worldY
		em.testRect.width = EN_SPRITE_W
		em.testRect.height = EN_SPRITE_H

		// check collision
		sideCollisions := em.game.tileMap.getSideCollisionData(*em.testRect)
		if !sideCollisions.down {
			entity.velY += PL_GRAVITY_AMOUNT
		} else if entity.velY > 0 {
			entity.velY = 0

		}

		// update position
		entity.worldX += entity.velX
		entity.worldY += entity.velY

	}

}

func (em *EntityManager) touchEntityAction(kind, uid int) {
	//fmt.Println("Entity touched ", kind)
	//ent.game.incrementScore(1)
	if kind == 0 && em.game.activateObject == true {
		em.game.warpManager.warpPlayerToWarpID(uid)
		em.game.activateObject = false
	}
}

func (em *EntityManager) checkEntitysTouchedPlayer() {

	playerRect := em.game.player.getWorldColliderRect()
	for _, v := range em.entityList {
		if nil != v && true == v.alive {
			em.testRect.x = em.game.tileMap.tileSize * v.startGridX
			em.testRect.y = em.game.tileMap.tileSize * v.startGridY

			if collideRect(playerRect, *em.testRect) {
				//v.alive = false
				em.touchEntityAction(v.kind, v.uid)
			}
		}
	}
}

func (ent *EntityManager) initImages() error {
	ent.images = []*ebiten.Image{}
	imageDir := path.Join(subdir, IMAGES_IDLE_SHEET)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	jackieImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	jackieImage = FlipHorizontal(jackieImage)

	ent.images = append(ent.images, jackieImage)
	//skull
	imageDir = path.Join(subdir, IMAGES_WALK_SHEET)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	skullImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	ent.images = append(ent.images, skullImage)
	//spikes
	imageDir = path.Join(subdir, IMAGES_ATTACK_SHEET)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	spikeImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	ent.images = append(ent.images, spikeImage)
	return err

}

func (ent *EntityManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(ent.entityList); i++ {
		if nil == ent.entityList[i] || false == ent.entityList[i].alive {
			return i
		}
	}
	return -1
}

func (ent *EntityManager) saveDataToFile() {
	name := ent.getDataFileURL()
	numericData := [][]int{}
	rows := len(ent.entityList)
	for i := 0; i < rows; i++ {
		entObj := ent.entityList[i]
		if entObj != nil {
			record := []int{entObj.kind, entObj.startGridX, entObj.startGridY, entObj.uid}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Entitys: no data to write, ", name)
	}
}
func (ent *EntityManager) getDataFileURL() string {
	filename := ent.filename_base + strconv.Itoa(ent.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}
func (ent *EntityManager) loadDataFromFile() error {
	ent.entityList = []*Entity{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := ent.getDataFileURL()
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
		//ent.entityList[i] = &Entity{v[0], v[1], v[2], v[3], true, true}
		entityTemp := NewEntity(v[0], v[1], v[2])
		entityTemp.uid = v[3]
		ent.entityList = append(ent.entityList, entityTemp)
		fmt.Println("added entity ")
	}
	return nil
}
func (ent *EntityManager) getUniqueUID() int {

	return 0
}

func (ent *EntityManager) AddInstanceToGrid(gridX, gridY, kind int) {
	//emptySlot := ent.inactiveSlot()
	if 1 == 1 {
		x := gridX
		y := gridY
		//uid := ent.getUniqueUID()
		entity := NewEntity(kind, x, y)
		entity.alive = true
		entity.uid = ent.getUniqueUID()
		entity.width = EN_SPRITE_W
		entity.height = EN_SPRITE_H
		ent.entityList = append(ent.entityList, entity)
		log.Printf("Added Entity %d at %d, %d\n", kind, gridX, gridY)
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

func (tm *EntityManager) setAssetID(assetID int) {

	if assetID < len(tm.images) && assetID >= 0 {
		tm.assetID = assetID
	}

}
