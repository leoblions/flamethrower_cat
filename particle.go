package main

import (
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PAR_MAX_PARTICLES      = 10
	PAR_DEFAULT_LIFE       = 200
	PAR_IMAGE_FILENAME_0   = "smokeSmall.png"
	PAR_IMAGE_FILENAME_1   = "bubbleSmall.png"
	PAR_SMOKE_RISE_AMOUNT  = 1
	PAR_BUBBLE_RISE_AMOUNT = 2
)

type Particle struct {
	worldX int
	worldY int
	kind   int
	frame  int
	life   int
}

type ParticleManager struct {
	game      *Game
	images    []*ebiten.Image
	particles [PAR_MAX_PARTICLES]*Particle
}

func NewParticleManager(game *Game) *ParticleManager {
	pmg := &ParticleManager{}
	pmg.game = game
	pmg.particles = [PAR_MAX_PARTICLES]*Particle{}
	pmg.images = []*ebiten.Image{}
	if err := pmg.initImages(); err != nil {
		log.Fatal(err)
	}
	//pmg.particles[0] = NewParticle(0, 0, 0)
	//pmg.AddParticle(0, 0, 0)
	return pmg

}

func NewParticle(worldX, worldY, kind int) *Particle {
	fb := &Particle{worldX, worldY, kind, 0, PAR_DEFAULT_LIFE}
	return fb
}

func (pmg *ParticleManager) initImages() error {
	// projectile images

	// smoke
	imageDir := path.Join(subdir, PAR_IMAGE_FILENAME_0)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	imgTemp := copyAndStretchImage(rawImage, 50, 50)
	pmg.images = append(pmg.images, imgTemp)
	// bubble

	imageDir = path.Join(subdir, PAR_IMAGE_FILENAME_1)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	imgTemp = copyAndStretchImage(rawImage, 10, 10)
	pmg.images = append(pmg.images, imgTemp)

	return err
}

func (pm *ParticleManager) AddParticle(worldX, worldY, kind int) {

	for i, v := range pm.particles {
		if v == nil || v.life <= 0 {
			pm.particles[i] = &Particle{}
			pm.particles[i].life = PAR_DEFAULT_LIFE
			pm.particles[i].worldX = worldX
			pm.particles[i].worldY = worldY
			pm.particles[i].kind = kind
			return
		}
	}

}

func (pm *ParticleManager) Update() {
	for _, v := range pm.particles {
		if v != nil && v.life > 0 {
			change := v.life%PM_FIRE_ANIMATE_SPEED == 0
			if change && v.frame < 3 {
				v.frame++

			} else if change {
				v.frame = 0
			}
			if v.kind == 0 {
				v.worldY -= PAR_SMOKE_RISE_AMOUNT
			} else if v.kind == 1 {
				v.worldY -= PAR_BUBBLE_RISE_AMOUNT
			}
			v.life--
		}
	}

}

func (pm *ParticleManager) Draw(screen *ebiten.Image) {

	for _, v := range pm.particles {
		if nil != v && v.life > 0 {

			screenX := (v.worldX) - worldOffsetX
			screenY := (v.worldY) - worldOffsetY
			DrawImageAt(screen, pm.images[v.kind], screenX, screenY)
		}
	}
}
