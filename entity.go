package main

import (
	"fmt"
	"log"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/*
kinds:
0 = Jackie (non enemy)
1 = robodog
2 = worm blob
3 = golem
4 = evil jackie

states:
0 = walk
1 = attack
2 = stand

sprite index:
0 - 3 = left
3 - 7 = right
*/
const (
	FM_MAX_ENTITY_ROOM          = 10
	IMAGES_IDLE_SHEET           = "jackieD1.png"
	IMAGES_WALK_SHEET           = "jackieD1.png"
	IMAGES_ATTACK_SHEET         = "jackieD1.png"
	IMAGES_JACKIE               = "jackieRS.png"
	IMAGES_MONSTER              = "enemyWog.png"
	IMAGES_ANT                  = "entityAnt.png"
	EN_FILENAME_BASE            = "entity"
	EN_FILENAME_END             = ".csv"
	EN_SPRITE_H                 = 100
	EN_SPRITE_W                 = 100
	EN_SPRITE_SIZE              = 100
	EN_CREATE_FILE_IF_NOT_EXIST = true
	EN_FOLLOW_DIST              = 300
	EN_ENEMY_SPEED_1            = 2
	EN_STOP_FOLLOW_DIST         = 100
	EN_SPRITES_PER_ROW          = 4
	EN_FRAME_MAX_VAL            = 3
	EN_FRAME_CHANGE_TICKS       = 30
	EN_DEFAULT_HEALTH           = 100
	EN_SHOOT_FREQUENCY          = 820000000
	EN_ENTITY_SHOOT_RANGE       = 400
	EM_KIND_MAX                 = 4
)

type EntityManager struct {
	game                 *Game
	maxEntitys           int
	entityList           []*Entity
	images               []*ebiten.Image
	testRect             *rect
	assetID              int
	filename_base        string
	frameChangeTickCount int
	enemyLastFiredBullet int64
	entAttackRect        *rect

	AllEntitySpriteCollections
}

type Entity struct {
	lastFiredAtPlayer int64
	kind              int
	startGridX        int
	startGridY        int
	uid               int
	currentImage      *ebiten.Image
	health            int
	frame             int
	state             int

	direction     rune
	prevDirection rune
	alive         bool
	onScreen      bool
	isEnemy       bool
	canShoot      bool

	MobileObject
}

func (ent *Entity) entityTakeDamage(damage int) {
	if ent.kind != 0 {
		newHealth := ent.health - damage
		if newHealth < 0 {
			ent.health = 0
			ent.alive = false
		} else {
			ent.health = newHealth
		}
	}
}

func NewEntity(kind, startGridX, startGridY int) *Entity {
	ent := &Entity{}
	worldX := GAME_TILE_SIZE * startGridX
	worldY := GAME_TILE_SIZE * startGridY
	ent.height = EN_SPRITE_H
	ent.width = EN_SPRITE_W
	ent.kind = kind
	if kind == 1 {
		ent.canShoot = true
	}
	ent.alive = true
	ent.startGridX = startGridX
	ent.startGridY = startGridY
	ent.worldX = worldX
	ent.direction = 'f'
	ent.health = EN_DEFAULT_HEALTH
	ent.worldY = worldY
	//fmt.Println("NewEntity add entity %d", kind)
	return ent

}

func NewEntityManager(game *Game) *EntityManager {

	fm := &EntityManager{}
	fm.game = game
	fm.filename_base = EN_FILENAME_BASE
	fm.maxEntitys = FM_MAX_ENTITY_ROOM
	fm.initImages()
	fm.initEntityImages()
	fm.entityList = []*Entity{}
	//ent.AddPickup(200, 200, 0)
	//ent.AddPickup(200, 300, 0)
	fm.testRect = &rect{0, 0, EN_SPRITE_W, EN_SPRITE_H}
	fm.assetID = 0
	fm.loadDataFromFile()
	fm.entAttackRect = &rect{}
	fm.enemyLastFiredBullet = time.Now().UnixNano()
	fm.testRect = &rect{}

	return fm
}

func (ent *EntityManager) Draw(screen *ebiten.Image) {

	for _, v := range ent.entityList {
		if nil != v && true == v.alive {
			//fmt.Println("ent draw")
			screenX := (v.worldX) - worldOffsetX
			screenY := (v.worldY) - worldOffsetY
			entImage := ent.selectImage(v.kind, v.state, v.frame)
			DrawImageAt(screen, entImage, screenX, screenY)
		}
	}

}
func (em *EntityManager) updateEntityListAnimationFrame() {
	if em.frameChangeTickCount < EN_FRAME_CHANGE_TICKS {
		em.frameChangeTickCount++
	} else {
		em.frameChangeTickCount = 0

		for _, entity := range em.entityList {
			if nil != entity && true == entity.alive {
				em.updateFrame(entity)

			}
		}
	}
}

func (em *EntityManager) Update() {

	em.entityMotion()
	em.checkEntitysTouchedPlayer()

	em.game.activateObject = false
	em.updateEntityListAnimationFrame()
	em.entitiesShootAtPlayer()

}

func (entity *Entity) entityPlayerDistance(game *Game) float64 {
	//fmt.Println("entity folow player")
	pposX := game.player.worldX
	pposY := game.player.worldY

	dX := float64(pposX - entity.worldX)
	dY := float64(pposY - entity.worldY)
	return math.Sqrt((dX * dX) + (dY * dY))

}
func (em *EntityManager) entitiesShootAtPlayer() {
	currTimeNano := time.Now().UnixNano()
	if abs(em.enemyLastFiredBullet-currTimeNano) > EN_SHOOT_FREQUENCY {
		em.enemyLastFiredBullet = currTimeNano
		for _, entity := range em.entityList {
			if nil != entity && true == entity.alive && entity.canShoot {
				epDistance := entity.entityPlayerDistance(em.game)
				if epDistance < EN_ENTITY_SHOOT_RANGE {
					em.game.projectileManager.AddProjectile(entity.worldX, entity.worldY, 1)
				}

			}
		}
	}

}

func (entity *Entity) entityFollowPlayer(game *Game) {
	pposX := game.player.worldX

	if entity.worldX+EN_STOP_FOLLOW_DIST < pposX {
		entity.velX = EN_ENEMY_SPEED_1

	} else if entity.worldX-EN_STOP_FOLLOW_DIST > pposX {
		entity.velX = -EN_ENEMY_SPEED_1

	} else {
		entity.velX = 0
	}

	if entity.worldX < pposX {
		entity.direction = 'l'
	} else if entity.worldX > pposX {
		entity.direction = 'r'

	} else {
		entity.direction = 'f'
	}

}

func (entity *Entity) entityMeleePlayer(game *Game) {
	if !entity.isEnemy {
		return
	}

	game.entityManager.entAttackRect = &rect{
		entity.worldX - EN_STOP_FOLLOW_DIST,
		entity.worldY - EN_STOP_FOLLOW_DIST,
		entity.width + EN_STOP_FOLLOW_DIST,
		entity.height + EN_STOP_FOLLOW_DIST}

	playerInAttackRange := collideRect(*game.entityManager.entAttackRect, game.player.getScreenCollrect())

	if playerInAttackRange {
		entity.state = 1
		//fmt.Println("enemy melee	")

	} else {
		entity.state = 0
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
		if entity.health > 0 {
			entity.entityFollowPlayer(em.game)
			entity.entityMeleePlayer(em.game)
		}

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

		entity.prevDirection = entity.direction

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
	//walk
	imageDir = path.Join(subdir, IMAGES_WALK_SHEET)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	skullImage := copyAndStretchImage(rawImage, EN_SPRITE_W, EN_SPRITE_H)
	ent.images = append(ent.images, skullImage)
	//attack
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

func (ent *EntityManager) rectCollideWithEntity(projectile *rect) *Entity {
	// find usable slot in pickups array, or -1 if there is none
	for _, v := range ent.entityList {
		if v != nil && v.alive {
			ent.testRect.x = v.worldX
			ent.testRect.y = v.worldY
			ent.testRect.height = v.height
			ent.testRect.width = v.width
			if collideRect(*ent.testRect, *projectile) {
				return v
			}

		}

	}

	return nil
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

func (tm *EntityManager) getAssetID() int {
	fmt.Println("EntityManager getAssetID", tm.assetID)
	return tm.assetID

}

func (tm *EntityManager) setAssetID(assetID int) {

	if assetID < EM_KIND_MAX && assetID >= 0 {
		tm.assetID = assetID
	}
	tm.assetID = assetID
	fmt.Println("EntityManager Selected entity type ", tm.assetID)

}
