package main

import (
	"fmt"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PM_MAX_PROJECTILES  = 10
	PM_IMAGE_FILENAME   = "bullet.png"
	PM_BULLET_SIZE      = 20
	PM_BULLET_LENGTH    = 50
	PM_PROJECTILE_SPEED = 11
	PM_DELAY_TICKS      = 15
)

type ProjectileManager struct {
	game             *Game
	kind             int
	startWorldX      int
	startWorldY      int
	projectileAmount int
	projectileMax    int
	projectileArray  [PM_MAX_PROJECTILES]*Projectile
	projectileImage  *ebiten.Image
	projectileDelay  int
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
	pm.projectileArray = [PM_MAX_PROJECTILES]*Projectile{}
	if err := pm.initImages(); err != nil {
		log.Fatal(err)
	}

	return pm

}

func (pm *ProjectileManager) initImages() error {
	imageDir := path.Join(subdir, PM_IMAGE_FILENAME)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	pm.projectileImage = copyAndStretchImage(rawImage, PM_BULLET_LENGTH, PM_BULLET_SIZE)
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
}

func (pm *ProjectileManager) Update() {
	for _, v := range pm.projectileArray {
		if v != nil && v.alive {
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

func (pm *ProjectileManager) AddProjectile() {
	if &pm.projectileArray == nil {
		fmt.Println("Array is null")
	}
	if pm.projectileDelay > 0 {
		pm.projectileDelay -= 1
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
			pm.projectileDelay = PM_DELAY_TICKS
			//fmt.Println("Added projectile at index ", i)
		}

	}
}
