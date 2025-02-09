package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PM_MAX_PROJECTILES        = 10
	PM_IMAGE_FILENAME_0       = "bullet.png"
	PM_IMAGE_FILENAME_1       = "greenBullet.png"
	PM_BULLET_SIZE            = 20
	PM_BULLET_LENGTH          = 50
	PM_PROJECTILE_SPEED       = 11
	PM_DELAY_TICKS            = 15
	PM_IMAGE_FIREBALL         = "fireball.png"
	PM_FIREBALL_LIFE          = 100
	PM_FIRE_ANIMATE_SPEED     = 7
	PM_MAX_FIREBALLS          = 10
	PM_DEBOUNCE_INTERVAL      = 120000000
	PM_DEF_DAMAGE_AMOUNT      = 33
	PM_BULLET_PLAYER_X_OFFSET = 25
	PM_FIREBALL_X_OFFSET      = -25
	PM_ENEMY_BULLET_SPEED     = 3.5
	PM_ENEMY_BULLET_X_OFFSET  = 50
	PM_ENEMY_BULLET_Y_OFFSET  = 50
)

type ProjectileManager struct {
	game             *Game
	kind             int
	startWorldX      int
	startWorldY      int
	projectileAmount int
	projectileMax    int
	debounceLastTime int64
	projectileArray  [PM_MAX_PROJECTILES]*Projectile
	//projectileArrayEnemy  [PM_MAX_PROJECTILES]*Projectile
	projectileImage []*ebiten.Image
	fireballImages  [4]*ebiten.Image
	projectileDelay int
	fireballList    []*Fireball
	testRect        *rect
}

type Fireball struct {
	worldX int
	worldY int
	frame  int
	life   int
}

func NewFireball(worldX, worldY int) *Fireball {
	fb := &Fireball{worldX, worldY, 0, PM_FIREBALL_LIFE}
	return fb
}

func (pm *ProjectileManager) addFireball(worldX, worldY int) {
	reusedSlot := false
	for _, v := range pm.fireballList {
		if v == nil || v.life <= 0 {
			v.life = PM_FIREBALL_LIFE
			v.worldX = worldX + PM_FIREBALL_X_OFFSET
			v.worldY = worldY
			reusedSlot = true
			return
		}
	}
	if !reusedSlot && len(pm.fireballList) < PM_MAX_FIREBALLS {
		fb := &Fireball{worldX, worldY, 0, PM_FIREBALL_LIFE}
		pm.fireballList = append(pm.fireballList, fb)
	}
}

func (pm *ProjectileManager) UpdateFireballs() {
	for _, v := range pm.fireballList {
		if v != nil && v.life > 0 {
			change := v.life%PM_FIRE_ANIMATE_SPEED == 0
			if change && v.frame < 3 {
				v.frame++
			} else if change {
				v.frame = 0
			}
			v.life--
		}
	}
}

type Projectile struct {
	MobileObject
	alive   bool
	kind    int
	stopped bool
}

func NewProjectileManager(game *Game, kind int) *ProjectileManager {
	pm := &ProjectileManager{}
	pm.game = game
	pm.kind = kind
	pm.testRect = &rect{}
	pm.projectileDelay = 0
	pm.projectileMax = PM_MAX_PROJECTILES
	pm.debounceLastTime = time.Now().UnixNano()
	pm.projectileArray = [PM_MAX_PROJECTILES]*Projectile{}
	pm.fireballList = []*Fireball{}
	if err := pm.initImages(); err != nil {
		log.Fatal(err)
	}

	return pm

}

func (pm *ProjectileManager) initImages() error {
	// projectile image
	// 0
	imageDir := path.Join(subdir, PM_IMAGE_FILENAME_0)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	imgTemp := copyAndStretchImage(rawImage, PM_BULLET_LENGTH, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)

	imageDir = path.Join(subdir, PM_IMAGE_FILENAME_1)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	imgTemp = copyAndStretchImage(rawImage, PM_BULLET_SIZE, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)

	// fireball images
	imageDir = path.Join(subdir, PM_IMAGE_FIREBALL)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)

	pm.fireballImages = [4]*ebiten.Image{}
	for i := 0; i < 4; i++ {
		x := 50 * i
		// magick numbers welcome here
		tempImage := getSubImage(rawImage, x, 0, 50, 50)
		pm.fireballImages[i] = tempImage
	}

	return err
}

func (pm *ProjectileManager) Draw(screen *ebiten.Image) {
	for _, v := range pm.projectileArray {
		if v != nil && v.alive {
			screenX := v.worldX - worldOffsetX
			screenY := v.worldY - worldOffsetY
			DrawImageAt(screen, pm.projectileImage[v.kind], screenX, screenY)
		}
	}

	pm.DrawFireBalls(screen)
}

func (pm *ProjectileManager) DrawFireBalls(screen *ebiten.Image) {
	for _, v := range pm.fireballList {
		if v.life <= 0 {
			return
		}

		fbCurrImage := pm.fireballImages[v.frame]
		screenX := v.worldX - worldOffsetX
		screenY := v.worldY - worldOffsetY
		DrawImageAt(screen, fbCurrImage, screenX, screenY)
	}
	//fmt.Println("FB list len ", len(pm.fireballList))
}

func (pm *ProjectileManager) Update() {
	pm.UpdateFireballs()
	for _, v := range pm.projectileArray {
		if v != nil && v.alive {
			pm.testRect.x = v.worldX
			pm.testRect.y = v.worldY
			pm.testRect.width = v.width
			pm.testRect.height = v.height
			if v.projectileCollideWall(pm.game) {
				v.alive = false
				pm.addFireball(v.worldX, v.worldY)
				//
				continue
			}

			if v.kind == 0 {
				ent := pm.game.entityManager.rectCollideWithEntity(pm.testRect)
				if ent != nil {
					v.alive = false
					pm.addFireball(v.worldX, v.worldY)
					ent.entityTakeDamage(PM_DEF_DAMAGE_AMOUNT)
					continue
				}

			} else {
				v.worldX += v.velX
				v.worldY += v.velY
				screenX := v.worldX - worldOffsetX
				screenY := v.worldY - worldOffsetY
				if screenY < 0 || screenY > pm.game.screenHeight ||
					screenX < 0 || screenX > pm.game.screenWidth {
					v.alive = false

				}
			}

		}
	}

}

func (pr *Projectile) projectileCollideWall(game *Game) bool {

	return game.tileMap.pointCollidedWithSolidTile(pr.worldX+25, pr.worldY+15)
}

func (pm *ProjectileManager) AddProjectile(startX, startY, kind int) {
	if &pm.projectileArray == nil {
		fmt.Println("Array is null")
	}

	timerNowNano := time.Now().UnixNano()

	interval := timerNowNano - pm.debounceLastTime
	if interval > CON_DEBOUNCE_INTERVAL {
		//fmt.Println("ADDPR ", interval)
		pm.debounceLastTime = timerNowNano
	} else {
		return
	}

	for i := 0; i < len(pm.projectileArray); i++ {
		if pm.projectileArray[i] == nil || pm.projectileArray[i].alive == false {
			Xoffset := 0
			Xvel := 0
			Yvel := 0
			if kind == 0 {
				if pm.game.player.faceLeft {
					Xoffset = 0
					Xvel = -PM_PROJECTILE_SPEED
				} else {
					Xoffset += playerWidth
					Xvel = PM_PROJECTILE_SPEED
				}

			} else {
				xt, yt := pm.enemyProjectileGetVelocity(startX, startY)
				Xvel = int(xt)
				Yvel = int(yt)
			}

			//startX := pm.game.player.worldX + Xoffset
			//startY := pm.game.player.worldY + (playerHeight / 2)
			newProj := &Projectile{}
			newProj.velX = Xvel
			newProj.velY = Yvel
			newProj.worldX = startX
			newProj.worldY = startY
			newProj.alive = true
			newProj.kind = kind
			newProj.stopped = false
			pm.projectileArray[i] = newProj
			return
			//fmt.Println("Added projectile at index ", i)
		}

	}
}

func (pm *ProjectileManager) enemyProjectileGetVelocity(startX, startY int) (float64, float64) {
	dX := float64(startX + PM_ENEMY_BULLET_X_OFFSET - pm.game.player.worldX)
	dY := float64(startY + PM_ENEMY_BULLET_Y_OFFSET - pm.game.player.worldY)
	if dX > dY {
		dY = dY / dX
		dX = 1.0
	} else {
		dX = dX / dY
		dY = 1.0
	}
	scaledX := -PM_ENEMY_BULLET_SPEED * dX
	scaledY := -PM_ENEMY_BULLET_SPEED * dY
	return scaledX, scaledY

}

func (pm *ProjectileManager) AddProjectile_0() {
	if &pm.projectileArray == nil {
		fmt.Println("Array is null")
	}

	timerNowNano := time.Now().UnixNano()

	interval := timerNowNano - pm.debounceLastTime
	if interval > CON_DEBOUNCE_INTERVAL {
		//fmt.Println("ADDPR ", interval)
		pm.debounceLastTime = timerNowNano
	} else {
		return
	}

	for i := 0; i < len(pm.projectileArray); i++ {
		if pm.projectileArray[i] == nil || pm.projectileArray[i].alive == false {
			Xoffset := 0
			Xvel := 0
			if pm.game.player.faceLeft {
				Xoffset = 0
				Xvel = -PM_PROJECTILE_SPEED
			} else {
				Xoffset += playerWidth
				Xvel = PM_PROJECTILE_SPEED
			}
			startX := pm.game.player.worldX + Xoffset
			startY := pm.game.player.worldY + (playerHeight / 2)
			newProj := &Projectile{}
			newProj.velX = Xvel
			newProj.velY = 0
			newProj.worldX = startX
			newProj.worldY = startY
			newProj.alive = true
			newProj.kind = pm.kind
			newProj.stopped = false
			pm.projectileArray[i] = newProj
			return
			//fmt.Println("Added projectile at index ", i)
		}

	}
}
