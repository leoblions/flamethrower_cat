package main

import (
	"math"
	"path"
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
5 = ant
6 = fly
7 = shark
8 = bird
9 = earwig
10 = boss body part

states:
0 = walk
1 = attack
2 = stand

sprite index:
0 - 3 = left
3 - 7 = right
*/
const (
	EN_MAX_ENTITIES_AT_ONCE     = 10
	IMAGES_IDLE_SHEET           = "jackieD1.png"
	IMAGES_WALK_SHEET           = "jackieD1.png"
	IMAGES_ATTACK_SHEET         = "jackieD1.png"
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
	EN_MAX_ENTITY_KINDS         = 12
	EM_BARNACLEFISH_TYPE        = 11
	EM_BOSS_HEALTH              = 300
)

const (
	FLY_KIND   = 6
	SHARK_KIND = 7
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
	entityImageMap       map[string][]*ebiten.Image
	//AllEntitySpriteCollections
	esCollections [EN_MAX_ENTITY_KINDS]*EntitySpriteCollection
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
	canFly            bool
	motionMethod      func(*Game)
	drawMethod        func(*ebiten.Image, *Entity)
	bossParts         []*BossPart
	direction         rune
	prevDirection     rune
	alive             bool
	onScreen          bool
	isEnemy           bool
	canShoot          bool
	gravityOn         bool

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
	if kind != 0 {
		ent.isEnemy = true
	}
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
	ent.gravityOn = true

	if ent.kind == BOSS_ENT_KIND || ent.kind == EM_BARNACLEFISH_TYPE {
		// regular enemy
		//ent.motionMethod = ent.mobileEntityFollow
		// boss
		ent.motionMethod = ent.bossMotion
		ent.gravityOn = false
		ent.health = EM_BOSS_HEALTH
	} else {
		// regular enemy
		ent.motionMethod = ent.mobileEntityFollow
		// boss
		//ent.motionMethod = ent.bossMotion
	}

	if ent.kind == FLY_KIND || ent.kind == SHARK_KIND {
		ent.canFly = true
	}
	//fmt.Println("NewEntity add entity %d", kind)
	return ent

}

func NewEntityManager(game *Game) *EntityManager {

	fm := &EntityManager{}
	fm.esCollections = [EN_MAX_ENTITY_KINDS]*EntitySpriteCollection{}
	fm.game = game
	fm.filename_base = EN_FILENAME_BASE
	fm.maxEntitys = EN_MAX_ENTITIES_AT_ONCE
	fm.initImages()
	fm.initEntityImages()
	fm.entityList = []*Entity{}
	fm.testRect = &rect{0, 0, EN_SPRITE_W, EN_SPRITE_H}
	fm.assetID = 0
	fm.loadDataFromFile()
	fm.entAttackRect = &rect{}
	fm.enemyLastFiredBullet = time.Now().UnixNano()
	fm.testRect = &rect{}

	return fm
}

func (em *EntityManager) Draw(screen *ebiten.Image) {

	for i, v := range em.entityList {
		entPtr := em.entityList[i]
		if nil != v && true == v.alive {
			// screenX := (v.worldX) - worldOffsetX
			// screenY := (v.worldY) - worldOffsetY
			// entImage := em.selectImage(v.kind, v.state, v.frame)
			// DrawImageAt(screen, entImage, screenX, screenY)
			if v.kind == BOSS_ENT_KIND {
				em.drawBoss(screen, entPtr)
			} else {
				em.drawRegularEntity(screen, entPtr)
			}
		}
	}

}

func (em *EntityManager) drawRegularEntity(screen *ebiten.Image, ent *Entity) {

	screenX := (ent.worldX) - worldOffsetX
	screenY := (ent.worldY) - worldOffsetY
	entImage := em.selectImage(ent.kind, ent.state, ent.frame)
	DrawImageAt(screen, entImage, screenX, screenY)

}

func (em *EntityManager) drawBoss(screen *ebiten.Image, ent *Entity) {
	em.drawRegularEntity(screen, ent)
	for _, v := range ent.bossParts {
		screenX := (v.worldX) - worldOffsetX
		screenY := (v.worldY) - worldOffsetY
		entImage := em.selectImage(ent.kind, ent.state, v.frame)
		DrawImageAt(screen, entImage, screenX, screenY)
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
					em.game.projectileManager.AddProjectile(
						entity.worldX, entity.worldY, 1)
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

	if entity.canFly {
		pposY := game.player.worldY

		if entity.worldY+EN_STOP_FOLLOW_DIST < pposY {
			entity.velY = EN_ENEMY_SPEED_1

		} else if entity.worldY-EN_STOP_FOLLOW_DIST > pposY {
			entity.velY = -EN_ENEMY_SPEED_1

		} else {
			entity.velY = 0
		}
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

func (entity *Entity) mobileEntityFollow(game *Game) {
	if entity.health > 0 {
		entity.entityFollowPlayer(game)
		entity.entityMeleePlayer(game)
	}

	//fmt.Println("EM ", entity.entityDetectPlatformEdge(em.game))
	if !entity.canFly && (entity.entityDetectPlatformEdge(game) ||
		entity.entityDetectAdjacentWall(game)) {
		entity.velX = 0

	}

}

func (entity *Entity) bossBodyPartMotion(game *Game) {

}

func (entity *Entity) bossMotion(game *Game) {
	if entity.health > 0 {
		entity.entityFollowPlayer(game)
		entity.entityMeleePlayer(game)
	}
	if entity.kind == EM_BARNACLEFISH_TYPE {
		entity.velX = 0
		entity.velY = 0
		entity.canFly = true
	}

}

func (em *EntityManager) entityMotion() {

	for _, entity := range em.entityList {

		//modifiable function pointer to perform movement
		entity.motionMethod(em.game)

		//calculate coll rect
		em.testRect.x = entity.worldX
		em.testRect.y = entity.worldY
		em.testRect.width = EN_SPRITE_W
		em.testRect.height = EN_SPRITE_H

		// check collision
		sideCollisions := em.game.tileMap.getSideCollisionData(*em.testRect)
		if !sideCollisions.down && entity.gravityOn {
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

func (em *EntityManager) touchEntityAction(kind, index int) {
	//fmt.Println("Entity touched ", kind)
	//ent.game.incrementScore(1)
	entity := em.entityList[index]
	if entity.isEnemy {
		em.game.player.takeDamageWithDebounce(5)
	}
}

func (em *EntityManager) checkEntitysTouchedPlayer() {

	playerRect := em.game.player.getWorldColliderRect()
	for i, v := range em.entityList {
		if nil != v && true == v.alive {
			em.testRect.x = em.game.tileMap.tileSize * v.startGridX
			em.testRect.y = em.game.tileMap.tileSize * v.startGridY
			em.testRect.width = v.width
			em.testRect.height = v.height

			if collideRect(playerRect, *em.testRect) {
				//v.alive = false
				em.touchEntityAction(v.kind, i)
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
