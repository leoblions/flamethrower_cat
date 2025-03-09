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
	PM_MAX_PROJECTILES         = 10
	PM_IMAGE_FILENAME_0        = "bullet.png"
	PM_IMAGE_FILENAME_1        = "greenBullet.png"
	PM_BULLET_SIZE             = 20
	PM_BULLET_LENGTH           = 50
	PM_PROJECTILE_SPEED        = 11
	PM_DELAY_TICKS             = 15
	PM_IMAGE_FIREBALL          = "fireball.png"
	PM_IMAGE_SPLAT             = "splatRing.png"
	PM_IMAGE_BUBBLE            = "bubble.png"
	PM_IMAGE_FIREBALL3         = "projectile3.png"
	PM_IMAGE_EXPLOSION         = "explosion.png"
	PM_FIREBALL_LIFE           = 100
	PM_EXPLOSION_LIFE          = 10
	PM_FIRE_ANIMATE_SPEED      = 7
	PM_EXPLOSION_ANIMATE_SPEED = 3
	PM_MAX_FIREBALLS           = 10
	PM_COMPONENT_VECTOR_MAX    = 10
	PM_DEBOUNCE_INTERVAL       = 120000000
	PM_DEF_DAMAGE_AMOUNT       = 33
	PM_BULLET_PLAYER_X_OFFSET  = 25
	PM_FIREBALL_X_OFFSET       = -25
	PM_ENEMY_BULLET_SPEED      = 3.5
	PM_ENEMY_BULLET_X_OFFSET   = 50
	PM_ENEMY_BULLET_Y_OFFSET   = 50
	PM_DAMEAGE_PLAYER          = 1
	PM_OFFSCREEN_CULL_DISTANCE = 200
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
	projectileImage  []*ebiten.Image
	projectileImage3 *ebiten.Image
	fireballImages   [4]*ebiten.Image
	splatImages      [4]*ebiten.Image
	explosionImages  [5]*ebiten.Image
	projectileDelay  int
	fireballList     []*Fireball
	testRect         *rect
}

type Fireball struct {
	worldX int
	worldY int
	kind   int
	frame  int
	life   int
}

func NewFireball(worldX, worldY, kind int) *Fireball {
	fb := &Fireball{worldX, worldY, kind, 0, PM_FIREBALL_LIFE}
	return fb
}

func (pm *ProjectileManager) addFireball(worldX, worldY, kind int) {
	reusedSlot := false
	var fblife int
	if kind == 2 {
		fblife = PM_EXPLOSION_LIFE
	} else {
		fblife = PM_FIREBALL_LIFE
	}
	for _, v := range pm.fireballList {
		if v == nil || v.life <= 0 {
			v.life = fblife
			v.worldX = worldX + PM_FIREBALL_X_OFFSET
			v.worldY = worldY
			v.kind = kind
			reusedSlot = true
			return
		}
	}
	if !reusedSlot && len(pm.fireballList) < PM_MAX_FIREBALLS {
		fb := &Fireball{worldX, worldY, kind, 0, fblife}
		pm.fireballList = append(pm.fireballList, fb)
	}
}

func (pm *ProjectileManager) UpdateFireballs() {
	for _, v := range pm.fireballList {
		if v != nil && v.life > 0 {
			change := v.life%PM_FIRE_ANIMATE_SPEED == 0
			if v.kind == 2 {
				change = v.life%PM_EXPLOSION_ANIMATE_SPEED == 0
			}
			if change && v.frame < 3 {
				v.frame++
			} else if change {
				v.frame = 0
				pm.game.particleManager.AddParticle(v.worldX, v.worldY, 0)
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
	// projectile images

	// player flamethrower
	imageDir := path.Join(subdir, PM_IMAGE_FILENAME_0)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	imgTemp := copyAndStretchImage(rawImage, PM_BULLET_LENGTH, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)
	// green bullet
	imageDir = path.Join(subdir, PM_IMAGE_FILENAME_1)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	imgTemp = copyAndStretchImage(rawImage, PM_BULLET_SIZE, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)
	// bubble
	imageDir = path.Join(subdir, PM_IMAGE_BUBBLE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	imgTemp = copyAndStretchImage(rawImage, PM_BULLET_SIZE, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)
	// explosion
	imageDir = path.Join(subdir, PM_IMAGE_EXPLOSION)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	imgTemp = copyAndStretchImage(rawImage, PM_BULLET_SIZE, PM_BULLET_SIZE)
	pm.projectileImage = append(pm.projectileImage, imgTemp)
	// fb proj 3
	imageDir = path.Join(subdir, PM_IMAGE_FIREBALL3)
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

	// splat  images
	imageDir = path.Join(subdir, PM_IMAGE_SPLAT)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)

	pm.splatImages = [4]*ebiten.Image{}
	for i := 0; i < 4; i++ {
		x := 50 * i
		// magick numbers welcome here
		tempImage := getSubImage(rawImage, x, 0, 50, 50)
		pm.splatImages[i] = tempImage
	}

	// explosion  images
	imageDir = path.Join(subdir, PM_IMAGE_EXPLOSION)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)

	pm.explosionImages = [5]*ebiten.Image{}
	for i := 0; i < 5; i++ {
		x := 50 * i
		// magick numbers welcome here
		tempImage := getSubImage(rawImage, x, 0, 50, 50)
		tempImage = copyAndStretchImage(tempImage, 100, 100)
		pm.explosionImages[i] = tempImage
	}

	// fireball 3
	imageDir = path.Join(subdir, PM_IMAGE_FIREBALL3)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)

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
		var fbCurrImage *ebiten.Image
		if v.kind == 0 {
			fbCurrImage = pm.fireballImages[v.frame]
		} else if v.kind == 1 {
			fbCurrImage = pm.splatImages[v.frame]
		} else {
			fbCurrImage = pm.explosionImages[v.frame]
		}

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
				fbKind := 0
				if v.kind == 1 {
					fbKind = 1
				}
				pm.addFireball(v.worldX, v.worldY, fbKind)
				//
				continue
			}

			v.worldX += v.velX
			v.worldY += v.velY
			screenX := v.worldX - worldOffsetX
			screenY := v.worldY - worldOffsetY
			// culling
			if screenY < 0-PM_OFFSCREEN_CULL_DISTANCE ||
				screenY > pm.game.screenHeight+PM_OFFSCREEN_CULL_DISTANCE ||
				screenX < 0-PM_OFFSCREEN_CULL_DISTANCE ||
				screenX > pm.game.screenWidth+PM_OFFSCREEN_CULL_DISTANCE {
				v.alive = false
				continue
				// projectile went off screen

			}

			if v.kind == 0 {
				ent := pm.game.entityManager.rectCollideWithEntity(pm.testRect)
				if ent != nil {
					// player projectile hits entity
					v.alive = false
					pm.addFireball(v.worldX, v.worldY, v.kind)
					ent.entityTakeDamage(PM_DEF_DAMAGE_AMOUNT)
					continue
				}

			} else if v.kind == 1 {
				if v.projectileCollidePlayer(pm.game) {
					//fmt.Println("enemy proj hit player")
					pm.game.player.changeHealthRelative(-PM_DAMEAGE_PLAYER)
					continue
				}
			}

		}
	}

}

func (pr *Projectile) projectileCollideWall(game *Game) bool {

	return game.tileMap.pointCollidedWithSolidTile(pr.worldX+25, pr.worldY+15)
}

func (pr *Projectile) projectileCollidePlayer(game *Game) bool {
	game.projectileManager.testRect.x = pr.worldX
	game.projectileManager.testRect.y = pr.worldY
	return collideRect(*game.player.collRect, *game.projectileManager.testRect)
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
			// player projectiles
			if kind == 0 || kind == 2 {
				if pm.game.player.faceLeft {
					Xoffset = 0
					Xvel = -PM_PROJECTILE_SPEED
				} else {
					Xoffset += PL_COLLRECT_W
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
			if kind == 1 {
				newProj.worldX += PM_ENEMY_BULLET_X_OFFSET
				newProj.worldY += PM_ENEMY_BULLET_Y_OFFSET
			}
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
	dX := float64(startX - pm.game.player.worldX)
	dY := float64(startY - pm.game.player.worldY)
	divisor := 1.0
	if abs(dX) > abs(dY) {
		divisor = abs(dX)
	} else if abs(dX) < abs(dY) {
		divisor = abs(dY)
	}

	dX = dX / divisor
	dY = dY / divisor

	dX = clamp(-PM_COMPONENT_VECTOR_MAX, PM_COMPONENT_VECTOR_MAX, dX)
	dY = clamp(-PM_COMPONENT_VECTOR_MAX, PM_COMPONENT_VECTOR_MAX, dY)
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
				Xoffset += PL_COLLRECT_W
				Xvel = PM_PROJECTILE_SPEED
			}
			startX := pm.game.player.worldX + Xoffset
			startY := pm.game.player.worldY + (PL_COLLRECT_H / 2)
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
