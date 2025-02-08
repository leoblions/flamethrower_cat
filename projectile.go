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
	PM_MAX_PROJECTILES    = 10
	PM_IMAGE_FILENAME     = "bullet.png"
	PM_BULLET_SIZE        = 20
	PM_BULLET_LENGTH      = 50
	PM_PROJECTILE_SPEED   = 11
	PM_DELAY_TICKS        = 15
	PM_IMAGE_FIREBALL     = "fireball.png"
	PM_FIREBALL_LIFE      = 100
	PM_FIRE_ANIMATE_SPEED = 7
	PM_MAX_FIREBALLS      = 10
	PM_DEBOUNCE_INTERVAL  = 120000000
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
	projectileImage  *ebiten.Image
	fireballImages   [4]*ebiten.Image
	projectileDelay  int
	fireballList     []*Fireball
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
			v.worldX = worldX
			v.worldY = worldY
			reusedSlot = true
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
	velX    int
	velY    int
	worldX  int
	worldY  int
	alive   bool
	kind    int
	stopped bool
}

func NewProjectileManager(game *Game, kind int) *ProjectileManager {
	pm := &ProjectileManager{}
	pm.game = game
	pm.kind = kind
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

	imageDir := path.Join(subdir, PM_IMAGE_FILENAME)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	pm.projectileImage = copyAndStretchImage(rawImage, PM_BULLET_LENGTH, PM_BULLET_SIZE)

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
			DrawImageAt(screen, pm.projectileImage, screenX, screenY)
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
}

func (pm *ProjectileManager) Update() {
	pm.UpdateFireballs()
	for _, v := range pm.projectileArray {
		if v != nil && v.alive {
			if v.projectileCollideWall(pm.game) {
				v.alive = false
				pm.addFireball(v.worldX, v.worldY)
				continue
			}
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

func (pr *Projectile) projectileCollideWall(game *Game) bool {

	return game.tileMap.pointCollidedWithSolidTile(pr.worldX+25, pr.worldY+15)
}

func (pm *ProjectileManager) AddProjectile() {
	if &pm.projectileArray == nil {
		fmt.Println("Array is null")
	}

	timerNowNano := time.Now().UnixNano()

	interval := timerNowNano - pm.debounceLastTime
	if interval > CON_DEBOUNCE_INTERVAL {

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
			pm.projectileArray[i] = &Projectile{Xvel, 0, startX, startY, true, pm.kind, false}

			//fmt.Println("Added projectile at index ", i)
		}

	}
}
